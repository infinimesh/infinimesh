/*
Copyright © 2022 Infinite Devices GmbH

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
