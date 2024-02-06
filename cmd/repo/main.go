/*
Copyright © 2021-2023 Infinite Devices GmbH

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
package main

import (
	"fmt"
	"github.com/infinimesh/infinimesh/pkg/oauth"
	"gopkg.in/yaml.v2"
	"net/http"
	"os"
	"strings"

	"connectrpc.com/connect"
	"connectrpc.com/grpchealth"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/infinimesh/infinimesh/pkg/graph"
	"github.com/infinimesh/infinimesh/pkg/graph/schema"
	logger "github.com/infinimesh/infinimesh/pkg/log"
	auth "github.com/infinimesh/infinimesh/pkg/shared/auth"
	"github.com/infinimesh/proto/handsfree"
	"github.com/infinimesh/proto/node/nodeconnect"
	"github.com/infinimesh/proto/plugins/pluginsconnect"
	shadowpb "github.com/infinimesh/proto/shadow"
	"github.com/rs/cors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	log *zap.Logger

	port string

	arangodbHost string
	arangodbCred string

	accountsConnection   string
	namespacesConnection string
	configs              map[string]oauth.Config

	rootPass string

	redisHost string

	SIGNING_KEY []byte
	services    map[string]bool
)

func init() {
	viper.AutomaticEnv()

	log = logger.NewLogger()

	viper.SetDefault("PORT", "8000")
	viper.SetDefault("DB_HOST", "db:8529")
	viper.SetDefault("DB_CRED", "root:openSesame")
	viper.SetDefault("SIGNING_KEY", "seeeecreet")
	viper.SetDefault("INF_DEFAULT_ROOT_PASS", "infinimesh")
	viper.SetDefault("REDIS_HOST", "redis:6379")

	viper.SetDefault("SERVICES", "accounts,namespaces,sessions,devices,shadow,plugins,internal,oauth")

	port = viper.GetString("PORT")

	arangodbHost = viper.GetString("DB_HOST")
	arangodbCred = viper.GetString("DB_CRED")

	SIGNING_KEY = []byte(viper.GetString("SIGNING_KEY"))
	rootPass = viper.GetString("INF_DEFAULT_ROOT_PASS")

	redisHost = viper.GetString("REDIS_HOST")

	services = make(map[string]bool)
	for _, s := range strings.Split(viper.GetString("SERVICES"), ",") {
		services[s] = true
	}

	viper.SetDefault("REGISTRY", "repo:8000")
	viper.SetDefault("NAMESPACES", "repo:8000")

	accountsConnection = viper.GetString("REGISTRY")
	namespacesConnection = viper.GetString("NAMESPACES")

	configs = map[string]oauth.Config{}
	file, err := os.ReadFile("oauth2_config.yaml")
	if err == nil {
		err = yaml.Unmarshal(file, &configs)
		if err != nil {
			log.Error("Failed to parse oauth config", zap.Error(err))
		}
	} else {
		log.Error("Failed to open file", zap.Error(err))
	}
}

func main() {
	defer func() {
		_ = log.Sync()
	}()

	log.Info("Connecting to DB", zap.String("URL", arangodbHost))
	db := schema.InitDB(log, arangodbHost, arangodbCred, rootPass, false)
	log.Info("DB connection established")

	log.Info("Connecting to Redis", zap.String("URL", redisHost))
	rdb := redis.NewClient(&redis.Options{
		Addr: redisHost,
		DB:   0, // use default DB
	})
	log.Info("Redis connection established")

	router := mux.NewRouter()
	router.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Debug("Request", zap.String("method", r.Method), zap.String("path", r.URL.Path))
			h.ServeHTTP(w, r)
		})
	})

	authInterceptor := auth.NewAuthInterceptor(log, rdb, nil, SIGNING_KEY)

	interceptors := connect.WithInterceptors(authInterceptor)

	log.Debug("Registering services", zap.Any("services", services))

	if _, ok := services["oauth"]; ok {
		accClient := nodeconnect.NewAccountsServiceClient(http.DefaultClient, accountsConnection)
		nsClient := nodeconnect.NewNamespacesServiceClient(http.DefaultClient, namespacesConnection)

		token, err := authInterceptor.MakeToken(schema.ROOT_ACCOUNT_KEY)
		if err != nil {
			log.Fatal("Failed to create token", zap.Error(err))
		}

		service := oauth.NewOauthService(log)
		service.Register(router, configs, accClient, nsClient, token)
	}

	ensure_root := false
	if _, ok := services["accounts"]; ok {
		log.Info("Registering accounts service")
		acc_ctrl := graph.NewAccountsController(log, db, rdb)
		acc_ctrl.SIGNING_KEY = SIGNING_KEY
		path, handler := nodeconnect.NewAccountsServiceHandler(acc_ctrl, interceptors)
		router.PathPrefix(path).Handler(handler)

		ensure_root = true
	}
	if _, ok := services["namespaces"]; ok {
		log.Info("Registering namespaces service")
		ns_ctrl := graph.NewNamespacesController(log, db)
		path, handler := nodeconnect.NewNamespacesServiceHandler(ns_ctrl, interceptors)
		router.PathPrefix(path).Handler(handler)

		ensure_root = true
	}

	if ensure_root {
		ica := graph.NewInfinimeshCommonActionsRepo(db)
		err := ica.EnsureRootExists(log, rdb, rootPass)
		if err != nil {
			log.Warn("Failed to ensure root exists", zap.Error(err))
		}
	}

	if _, ok := services["sessions"]; ok {
		log.Info("Registering sessions service")
		sess_ctrl := graph.NewSessionsController(log, rdb)
		path, handler := nodeconnect.NewSessionsServiceHandler(sess_ctrl, interceptors)
		router.PathPrefix(path).Handler(handler)
	}

	if _, ok := services["devices"]; ok {
		log.Info("Registering devices service")
		viper.SetDefault("HANDSFREE_HOST", "handsfree:8000")
		host := viper.GetString("HANDSFREE_HOST")
		conn, err := grpc.Dial(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatal("Failed to connect to handsfree", zap.String("address", host), zap.Error(err))
		}

		dev_ctrl := graph.NewDevicesControllerModule(log, db, handsfree.NewHandsfreeServiceClient(conn))
		dev_ctrl.SetSigningKey(SIGNING_KEY)

		path, handler := nodeconnect.NewDevicesServiceHandler(dev_ctrl.Handler(), interceptors)
		router.PathPrefix(path).Handler(handler)
	}
	if _, ok := services["shadow"]; ok {
		log.Info("Registering shadow service")
		viper.SetDefault("SHADOW_HOST", "shadow-api:8000")
		host := viper.GetString("SHADOW_HOST")
		conn, err := grpc.Dial(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatal("Failed to connect to shadow", zap.String("address", host), zap.Error(err))
		}
		client := shadowpb.NewShadowServiceClient(conn)

		path, handler := nodeconnect.NewShadowServiceHandler(NewShadowAPI(log, client), interceptors)
		router.PathPrefix(path).Handler(handler)
	}

	if _, ok := services["plugins"]; ok {
		log.Info("Registering plugins service")
		plug_ctrl := graph.NewPluginsController(log, db)

		path, handler := pluginsconnect.NewPluginsServiceHandler(plug_ctrl, interceptors)
		router.PathPrefix(path).Handler(handler)
	}

	if _, ok := services["internal"]; ok {
		log.Info("Registering Internal service")
		is := graph.InternalService{}
		path, handler := nodeconnect.NewInternalServiceHandler(&is, interceptors)
		router.PathPrefix(path).Handler(handler)
	}

	checker := grpchealth.NewStaticChecker()
	path, handler := grpchealth.NewHandler(checker)
	router.PathPrefix(path).Handler(handler)

	host := fmt.Sprintf("0.0.0.0:%s", port)

	handler = cors.New(cors.Options{
		AllowedOrigins:      []string{"*"},
		AllowedMethods:      []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders:      []string{"*", "Connect-Protocol-Version"},
		AllowCredentials:    true,
		AllowPrivateNetwork: true,
	}).Handler(h2c.NewHandler(router, &http2.Server{}))

	log.Info("Serving", zap.String("host", host))
	err := http.ListenAndServe(host, handler)
	if err != nil {
		log.Fatal("Failed to start server", zap.Error(err))
	}
}
