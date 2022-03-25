/*
Copyright © 2021-2022 Infinite Devices GmbH, Nikita Ivanovski info@slnt-opp.xyz

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
	"github.com/infinimesh/infinimesh/pkg/graph"
	logger "github.com/infinimesh/infinimesh/pkg/log"
	auth "github.com/infinimesh/infinimesh/pkg/shared/auth"
	connectdb "github.com/infinimesh/infinimesh/pkg/shared/connectdb"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	pb "github.com/infinimesh/infinimesh/pkg/node/proto"
)

var (
	port string

	arangodbHost  string
	arangodbCred  string
	SIGNING_KEY   []byte
)

func init() {
	viper.AutomaticEnv()

	viper.SetDefault("PORT", "8080")
	viper.SetDefault("DB_HOST", "db:8529")
	viper.SetDefault("DB_CRED", "root:openSesame")
	viper.SetDefault("SIGNING_KEY", "seeeecreet")

	port = viper.GetString("PORT")

	arangodbHost = viper.GetString("DB_HOST")
	arangodbCred = viper.GetString("DB_CRED")
	SIGNING_KEY = []byte(viper.GetString("SIGNING_KEY"))
}

func main() {
	log, err := logger.NewProdOrDev()
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = log.Sync()
	}()

	log.Info("Connecting to DB", zap.String("URL", arangodbHost))
	db := connectdb.MakeDBConnection(log, arangodbHost, arangodbCred)
	log.Info("DB connection established")

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

	acc_ctrl := graph.NewAccountsController(log, db)
	pb.RegisterAccountsServiceServer(s, acc_ctrl)
	ns_ctrl  := graph.NewNamespacesController(log, db)
	pb.RegisterNamespacesServiceServer(s, ns_ctrl)
	dev_ctrl := graph.NewDevicesController(log, db)
	pb.RegisterDevicesServiceServer(s, dev_ctrl)

	log.Info(fmt.Sprintf("Serving gRPC on 0.0.0.0:%v", port))
	log.Fatal("Failed to serve gRPC", zap.Error(s.Serve(lis)))
}