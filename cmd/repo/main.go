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
	"net"
	"strings"

	"github.com/go-redis/redis/v8"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/infinimesh/infinimesh/pkg/graph"
	"github.com/infinimesh/infinimesh/pkg/graph/schema"
	logger "github.com/infinimesh/infinimesh/pkg/log"
	auth "github.com/infinimesh/infinimesh/pkg/shared/auth"
	"github.com/infinimesh/proto/handsfree"
	"github.com/infinimesh/proto/plugins"
	shadowpb "github.com/infinimesh/proto/shadow"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	pb "github.com/infinimesh/proto/node"
)

var (
	log *zap.Logger

	port string

	arangodbHost string
	arangodbCred string

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

	viper.SetDefault("SERVICES", "accounts,namespaces,devices,shadow,plugins,internal")

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

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		log.Fatal("Failed to listen", zap.String("address", port), zap.Error(err))
	}

	auth.SetContext(log, SIGNING_KEY)
	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_zap.UnaryServerInterceptor(log),
			grpc.UnaryServerInterceptor(auth.JWT_AUTH_INTERCEPTOR),
		)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_zap.StreamServerInterceptor(log),
			grpc_auth.StreamServerInterceptor(auth.JwtDeviceAuthMiddleware),
		)),
	)

	log.Debug("Registering services", zap.Any("services", services))

	ensure_root := false
	if _, ok := services["accounts"]; ok {
		log.Info("Registering accounts service")
		acc_ctrl := graph.NewAccountsController(log, db, rdb)
		acc_ctrl.SIGNING_KEY = SIGNING_KEY
		pb.RegisterAccountsServiceServer(s, acc_ctrl)

		ensure_root = true
	}
	if _, ok := services["namespaces"]; ok {
		log.Info("Registering namespaces service")
		ns_ctrl := graph.NewNamespacesController(log, db)
		pb.RegisterNamespacesServiceServer(s, ns_ctrl)

		ensure_root = true
	}

	if ensure_root {
		err := graph.EnsureRootExists(log, db, rdb, rootPass)
		if err != nil {
			log.Warn("Failed to ensure root exists", zap.Error(err))
		}
	}

	if _, ok := services["devices"]; ok {
		log.Info("Registering devices service")
		viper.SetDefault("HANDSFREE_HOST", "handsfree:8000")
		host := viper.GetString("HANDSFREE_HOST")
		conn, err := grpc.Dial(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatal("Failed to connect to handsfree", zap.String("address", host), zap.Error(err))
		}

		dev_ctrl := graph.NewDevicesController(log, db, handsfree.NewHandsfreeServiceClient(conn))
		dev_ctrl.SIGNING_KEY = SIGNING_KEY

		pb.RegisterDevicesServiceServer(s, dev_ctrl)
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
		pb.RegisterShadowServiceServer(s, NewShadowAPI(log, client))
	}

	if _, ok := services["plugins"]; ok {
		log.Info("Registering plugins service")
		plug_ctrl := graph.NewPluginsController(log, db)
		plugins.RegisterPluginsServiceServer(s, plug_ctrl)
	}

	if _, ok := services["internal"]; ok {
		log.Info("Registering Internal service")
		is := graph.InternalService{}
		pb.RegisterInternalServiceServer(s, &is)
	}

	healthpb.RegisterHealthServer(s, health.NewServer())

	log.Info(fmt.Sprintf("Serving gRPC on 0.0.0.0:%v", port))
	log.Fatal("Failed to serve gRPC", zap.Error(s.Serve(lis)))
}
