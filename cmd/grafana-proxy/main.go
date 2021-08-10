//--------------------------------------------------------------------------
// Copyright 2018 Infinite Devices GmbH
// www.infinimesh.io
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.
//--------------------------------------------------------------------------

package main

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	jwt "github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	cache "github.com/patrickmn/go-cache"

	"time"

	"github.com/infinimesh/infinimesh/pkg/grafana"
	"github.com/infinimesh/infinimesh/pkg/node/nodepb"
)

var (
	signingSecret  []byte
	accountIDClaim = "account_id"

	nodeserver string
	grafanaURL string

	log *zap.Logger

	accountService   nodepb.AccountServiceClient
	namespaceService nodepb.NamespacesClient

	userCache *cache.Cache
	g         *grafana.Client
)

func init() {
	log, _ = zap.NewDevelopment()

	viper.SetDefault("NODE_HOST", "localhost:8082")
	viper.SetDefault("GRAFANA_URL", "http://localhost:3000")
	viper.AutomaticEnv()

	nodeserver = viper.GetString("NODE_HOST")
	grafanaURL = viper.GetString("GRAFANA_URL")

	{
		b64SignSecret := viper.GetString("JWT_SIGNING_KEY")
		if b64SignSecret == "" {
			panic("Invalid signing secret")
		}
		s, err := base64.StdEncoding.DecodeString(b64SignSecret)
		if err != nil {
			panic("Failed to base64 decode sign secret")
		}
		signingSecret = s
	}

	userCache = cache.New(1*time.Second, 5*time.Second)
	g = grafana.NewClient("http://localhost:3000", "admin", "admin")
}

// TODO remove accounts in grafana if deleted in platform
func syncAccounts() {
	resp, err := accountService.ListAccounts(context.Background(), &nodepb.ListAccountsRequest{})
	if err != nil {
		log.Error("Failed to ListAccounts", zap.Error(err))
		return
	}

	for _, account := range resp.Accounts {
		if !account.Enabled {
			log.Debug("Account disabled, skipping", zap.String("id", account.Uid), zap.String("name", account.Name), zap.Error(err))
			continue
		}
		err := g.CreateUser(account.Name)
		if err != nil {
			log.Debug("Could not create account", zap.String("id", account.Uid), zap.String("name", account.Name), zap.Error(err))
		} else {
			log.Info("Created account", zap.String("id", account.Uid), zap.String("name", account.Name))
		}

		if account.IsRoot {
			userID, err := g.GetUserID(account.Name)
			if err != nil {
				log.Error("Failed to get userID of root", zap.Error(err))
				continue
			}

			err = g.MakeUserAdmin(userID)
			if err != nil {
				log.Error("Failed to make root admin", zap.Error(err))
			}

			err = g.AddUserToOrg(1, account.Name, "Admin")
			if err != nil {
				log.Error("Could not add user to main org", zap.String("id", account.Uid), zap.String("name", account.Name), zap.Error(err))
			} else {

				log.Info("Made user root", zap.String("id", account.Uid), zap.String("name", account.Name), zap.Error(err))
			}
		}

	}
}

func syncPermissions(namespace string) {
	// Get Permissions per NS
	r, err := namespaceService.ListPermissions(context.Background(), &nodepb.ListPermissionsRequest{Namespace: namespace})
	if err != nil {
		log.Error("Failed to fetch perms", zap.Error(err))
	}

	orgID, err := g.GetOrgID(namespace)
	if err != nil {
		log.Error("Could not get ID for namespace", zap.String("namespace", namespace))
	}

	if orgID == 1 {
		log.Error("Ignoring Org with ID 1")
		return
	}
	for _, permission := range r.Permissions {
		log.Info("Found permission", zap.Any("perm", permission))

		role, err := actionToGrafanaRole(permission.Action)
		if err != nil {
			return
		}

		err = g.AddUserToOrg(orgID, permission.AccountName, role)
		if err != nil {
			log.Info("Failed to add user to org", zap.String("org/ns", namespace), zap.String("account", permission.AccountName), zap.Error(err))
		} else {
			log.Info("Added user to org", zap.String("org", namespace), zap.String("account", permission.AccountName))

			if permission.AccountName == namespace {

				userID, err := g.GetUserID(permission.AccountName)
				if err != nil {
					log.Info("Failed to get userID", zap.String("org/ns", namespace), zap.String("account", permission.AccountName), zap.Error(err))
				}

				err = g.SwitchUserOrg(userID, orgID)
				if err != nil {
					log.Info("Failed to switch user org", zap.String("org/ns", namespace), zap.String("account", permission.AccountName), zap.Error(err))
				}
				log.Info("Failed to set org", zap.String("org/ns", namespace), zap.String("account", permission.AccountName), zap.Error(err))
			}
		}
	}

}

func actionToGrafanaRole(action nodepb.Action) (string, error) {
	switch action {
	case nodepb.Action_READ:
		return "Viewer", nil
	case nodepb.Action_WRITE:
		return "Editor", nil
	}
	return "", errors.New("Invalid action")
}

// TODO remove namespaces & permissions deleted in platform
func syncNamespaces() {
	resp, err := namespaceService.ListNamespaces(context.Background(), &nodepb.ListNamespacesRequest{})
	if err != nil {
		log.Error("Failed to ListNamespaces", zap.Error(err))
	}

	for _, namespace := range resp.Namespaces {
		err := g.CreateOrg(namespace.Name)
		if err != nil {
			log.Debug("Could not create Org", zap.Error(err))
		} else {
			log.Info("Created Org", zap.String("name", namespace.Name))
		}

		syncPermissions(namespace.Name)

	}
}

func main() {
	u, err := url.Parse(grafanaURL)
	if err != nil {
		panic(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(u)

	c, err := grpc.Dial(nodeserver, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	accountService = nodepb.NewAccountServiceClient(c)
	namespaceService = nodepb.NewNamespacesClient(c)

	ticker := time.NewTicker(time.Minute * 5)

	go func() {
		time.Sleep(time.Second * 30)
		log.Info("Start Sync")
		syncAccounts()
		syncNamespaces()
		for range ticker.C {
			syncAccounts()
			syncNamespaces()
		}
	}()

	http.HandleFunc("/", handler(proxy))
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func handler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Request", zap.String("URL", r.URL.RequestURI()))
		cookie, err := r.Cookie("token")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userID, err := ParseJWT(cookie.Value)
		if err != nil {
			log.Error("Failed")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		var name string
		if cachedUser, ok := userCache.Get(userID); ok {
			if userString, ok := cachedUser.(string); ok {
				log.Debug("Using cached value", zap.String("userID", userID), zap.String("username", userString))
				name = userString
			}
		} else {
			acc, err := accountService.GetAccount(context.Background(), &nodepb.GetAccountRequest{
				Id: userID,
			})
			if err != nil {
				log.Error("Failed to fetch account", zap.String("userID", userID), zap.Error(err))
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			name = acc.Name

			_ = userCache.Add(userID, acc.Name, time.Minute*1)
		}

		r.Header.Set("X-WEBAUTH-USER", name)
		p.ServeHTTP(w, r)
	}
}

func ParseJWT(tokenString string) (user string, err error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, status.Error(codes.Unauthenticated, fmt.Sprintf("Unexpected signing method: %v", t.Header["alg"]))
		}
		return signingSecret, nil
	})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("Invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		log.Info("Validated token", zap.Any("claims", claims))

		if accountID, ok := claims[accountIDClaim]; ok {
			if accountIDStr, ok := accountID.(string); ok {
				return accountIDStr, nil
			}

		}
		log.Info("Token does not contain account id field", zap.Any("token", token))
	}

	return "", errors.New("Invalid token")

}
