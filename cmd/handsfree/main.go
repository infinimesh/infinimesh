/*
Copyright © 2021-2022 Infinite Devices GmbH

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

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/infinimesh/infinimesh/pkg/handsfree"
	"github.com/infinimesh/infinimesh/pkg/shared/auth"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	logger "github.com/infinimesh/infinimesh/pkg/log"

	pb "github.com/infinimesh/proto/handsfree"
)

var (
	log *zap.Logger

	port string

	SIGNING_KEY []byte
)

func init() {
	viper.AutomaticEnv()

	log = logger.NewLogger()

	viper.SetDefault("PORT", "8000")
	viper.SetDefault("SIGNING_KEY", "seeeecreet")

	port = viper.GetString("PORT")
	SIGNING_KEY = []byte(viper.GetString("SIGNING_KEY"))
}

func main() {
	defer func() {
		_ = log.Sync()
	}()

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
	)

	srv := handsfree.NewHandsfreeServer(log)
	pb.RegisterHandsfreeServiceServer(s, srv)

	log.Info(fmt.Sprintf("Serving gRPC on 0.0.0.0:%v", port))
	log.Fatal("Failed to serve gRPC", zap.Error(s.Serve(lis)))

}
