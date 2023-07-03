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

	"github.com/go-redis/redis/v8"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/infinimesh/infinimesh/pkg/handsfree"
	"github.com/infinimesh/infinimesh/pkg/shared/auth"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	logger "github.com/infinimesh/infinimesh/pkg/log"

	pb "github.com/infinimesh/proto/handsfree"
)

var (
	log *zap.Logger

	port string

	redisHost string

	SIGNING_KEY []byte
)

func init() {
	viper.AutomaticEnv()

	log = logger.NewLogger()

	viper.SetDefault("PORT", "8000")
	viper.SetDefault("SIGNING_KEY", "seeeecreet")
	viper.SetDefault("REDIS_HOST", "redis:6379")

	port = viper.GetString("PORT")
	SIGNING_KEY = []byte(viper.GetString("SIGNING_KEY"))

	redisHost = viper.GetString("REDIS_HOST")
}

func main() {
	defer func() {
		_ = log.Sync()
	}()

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

	auth.SetContext(log, rdb, SIGNING_KEY)
	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_zap.UnaryServerInterceptor(log),
			grpc.UnaryServerInterceptor(auth.JWT_AUTH_INTERCEPTOR),
		)),
	)

	srv := handsfree.NewHandsfreeServer(log)
	pb.RegisterHandsfreeServiceServer(s, srv)
	healthpb.RegisterHealthServer(s, health.NewServer())

	log.Info(fmt.Sprintf("Serving gRPC on 0.0.0.0:%v", port))
	log.Fatal("Failed to serve gRPC", zap.Error(s.Serve(lis)))

}
