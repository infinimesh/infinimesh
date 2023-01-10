/*
Copyright © 2023 Infinite Devices GmbH

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package credentials

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/arangodb/go-driver"
	"github.com/go-ldap/ldap/v3"
	"github.com/infinimesh/infinimesh/pkg/graph/schema"
	logger "github.com/infinimesh/infinimesh/pkg/log"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

var (
	LDAP_CONFIGURED = false
	LDAP            LDAPConfFile
)

func init() {
	viper.AutomaticEnv()
	viper.SetDefault("LDAP_CONF", "/home/slnt_opp/repos/infinimesh/e2e/ldap.yml")

	log := logger.NewLogger().Named("LDAP Credentials Init")
	path := viper.GetString("LDAP_CONF")
	if path == "" {
		log.Info("No LDAP config found, skipping...")
		return
	}

	data, err := os.ReadFile(path)
	if err != nil {
		log.Error("Can't read LDAP config, skipping", zap.Error(err))
		return
	}

	err = yaml.Unmarshal(data, &LDAP)
	if err != nil {
		log.Error("Error while parsing LDAP config, skipping", zap.Error(err))
	}

	conns := 0
	for key, provider := range LDAP.Providers {
		conn, err := ldap.DialURL(provider.URL)
		if err != nil {
			log.Error("Couldn't Dial LDAP provider, skipping",
				zap.String("url", provider.URL), zap.Error(err))
			delete(LDAP.Providers, key)
			continue
		}
		log.Info("Success dialing LDAP provider", zap.String("url", provider.URL))
		provider.Conn = conn
		LDAP.Providers[key] = provider
		conns++
	}

	if conns == 0 {
		log.Warn("Every LDAP provider failed on dial, LDAP not configured")
		return
	}
	LDAP_CONFIGURED = true

	exit := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-exit
		log.Info("Exit signal received, closing connections", zap.Any("signal", sig))
		for _, provider := range LDAP.Providers {
			provider.Conn.Close()
		}
	}()
}

type LDAPConfFile struct {
	Providers map[string]LDAPProvider `yaml:"providers"`
}

type LDAPProvider struct {
	Conn *ldap.Conn
	URL  string `yaml:"url"`

	BaseDN       string `yaml:"base_dn"`
	Scope        int    `yaml:"scope"`
	DerefAliases int    `yaml:"deref_aliases"`
	SizeLimit    int    `yaml:"size_limit"`
	TimeLimit    int    `yaml:"time_limit"`
}

type LDAPCredentials struct {
	Username    string `json:"username"`
	ProviderKey string `json:"key"`

	log *zap.Logger
	driver.DocumentMeta
}

func (c *LDAPCredentials) Listable() []string {
	return []string{
		c.Username, c.ProviderKey,
	}
}

func NewLDAPCredentials(username, key string) (Credentials, error) {
	if _, ok := LDAP.Providers[key]; !ok {
		return nil, fmt.Errorf("requested Provider Key(%s) is not registered", key)
	}
	return &LDAPCredentials{
		Username: username, ProviderKey: key,
	}, nil
}

func LDAPFromMap(d map[string]interface{}) (ListableCredentials, error) {
	c := &LDAPCredentials{}

	iuser, ok := d["username"]
	if !ok {
		return c, errors.New("'username' is not present")
	}

	user, ok := iuser.(string)
	if !ok {
		return c, errors.New("'username' is not string")
	}

	c.Username = user

	iprovider, ok := d["key"]
	if !ok {
		return c, errors.New("provider 'key' is not present")
	}

	provider, ok := iprovider.(string)
	if !ok {
		return c, errors.New("provider 'key' is not string")
	}
	c.ProviderKey = provider

	return c, nil
}

func (c *LDAPCredentials) Authorize(args ...string) bool {
	provider, ok := LDAP.Providers[c.ProviderKey]
	if !ok {
		c.log.Warn("Existent Credentials have wrong Provider Key", zap.String("key", c.ProviderKey))
		return false
	}

	req := ldap.NewSearchRequest(
		provider.BaseDN, provider.Scope, provider.DerefAliases,
		provider.SizeLimit, provider.TimeLimit, false,
		fmt.Sprintf("(&(objectClass=organizationalPerson)(uid=%s))", c.Username),
		[]string{"dn"}, nil,
	)
	res, err := provider.Conn.Search(req)
	if err != nil {
		c.log.Warn("Error while executing LDAP Tree search", zap.Error(err))
		return false
	}
	if len(res.Entries) != 1 {
		c.log.Warn("Result has none or too many results", zap.Int("length", len(res.Entries)))
	}

	err = provider.Conn.Bind(res.Entries[0].DN, args[1])
	if err != nil {
		c.log.Warn("Error Binding", zap.Error(err))
		return false
	}
	return true
}

func (*LDAPCredentials) Type() string {
	return "ldap"
}

func (sc *LDAPCredentials) Key() string {
	return sc.ID.Key()
}

func (c *LDAPCredentials) SetLogger(log *zap.Logger) {
	c.log = log.Named("LDAP Auth")
}

func (cred *LDAPCredentials) Find(ctx context.Context, db driver.Database) bool {
	query := `FOR cred IN @@credentials FILTER cred.username == @username RETURN cred`
	c, err := db.Query(ctx, query, map[string]interface{}{
		"username":     cred.Username,
		"@credentials": schema.CREDENTIALS_COL,
	})
	if err != nil {
		return false
	}
	defer c.Close()

	_, err = c.ReadDocument(ctx, &cred)
	return err == nil
}

func (cred *LDAPCredentials) FindByKey(ctx context.Context, col driver.Collection, key string) error {
	_, err := col.ReadDocument(ctx, key, cred)
	return err
}
