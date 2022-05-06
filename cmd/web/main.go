/*
Copyright © 2021-2022 Nikita Ivanovski info@slnt-opp.xyz

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
	"context"
	"crypto/tls"
	"net/http"
	"strings"

	logger "github.com/infinimesh/infinimesh/pkg/log"
	pb "github.com/infinimesh/infinimesh/pkg/node/proto"
	"github.com/tmc/grpc-websocket-proxy/wsproxy"

	"github.com/gorilla/handlers"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	log 			*zap.Logger
	
	apiserver 		string
	corsAllowed 	[]string
	secure 				bool
	with_block 		bool
)

func init() {
	viper.AutomaticEnv()
	Log, err := logger.NewProdOrDev()
	if err != nil {
		panic(err)
	}
	log = Log

	viper.SetDefault("CORS_ALLOWED", []string{"*"})
	viper.SetDefault("APISERVER_HOST", "proxy:8000")
	viper.SetDefault("SECURE", false)
	viper.SetDefault("WITH_BLOCK", false)

	apiserver   = viper.GetString("APISERVER_HOST")
	corsAllowedIn := viper.GetString("CORS_ALLOWED")
	if corsAllowedIn != "" {
		corsAllowed = strings.Split(corsAllowedIn, ",")
	}
	secure      = viper.GetBool("SECURE")
	with_block  = viper.GetBool("WITH_BLOCK")
}

func main() {
	defer func() {
		_ = log.Sync()
	}()

	log.Info("Starting REST-API Server")
	log.Info("Registering Endpoints", zap.String("server", apiserver))
	var err error

	gwmux := runtime.NewServeMux()
	creds := insecure.NewCredentials()
	if secure {
		creds = credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})
	}
	opts := []grpc.DialOption{grpc.WithTransportCredentials(creds)}

	if with_block {
		opts = append(opts, grpc.WithBlock())
	}

	log.Info("Registering Accounts Service")
	err = pb.RegisterAccountsServiceHandlerFromEndpoint(context.Background(), gwmux, apiserver, opts)
	if err != nil {
		log.Fatal("Failed to register AccountsService gateway", zap.Error(err))
	}

	log.Info("Registering Namespaces Service")
	err = pb.RegisterNamespacesServiceHandlerFromEndpoint(context.Background(), gwmux, apiserver, opts)
	if err != nil {
		log.Fatal("Failed to register NamespacesService gateway", zap.Error(err))
	}

	log.Info("Registering Devices Service")
	err = pb.RegisterDevicesServiceHandlerFromEndpoint(context.Background(), gwmux, apiserver, opts)
	if err != nil {
		log.Fatal("Failed to register DevicesService gateway", zap.Error(err))
	}

	log.Info("Registering Shadow Service")
	err = pb.RegisterShadowServiceHandlerFromEndpoint(context.Background(), gwmux, apiserver, opts)
	if err != nil {
		log.Fatal("Failed to register ShadowService gateway", zap.Error(err))
	}

	log.Info("Allowed Origins", zap.Strings("hosts", corsAllowed))
	handler := handlers.CORS(
		handlers.AllowedOrigins(corsAllowed),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "PUT", "PATCH", "OPTIONS", "HEAD"}),
	)(gwmux)

	log.Info("Serving gRPC-Gateway on http://0.0.0.0:8000")
	log.Fatal("Failed to Listen and Serve Gateway-Server", zap.Error(http.ListenAndServe(":8000", wsproxy.WebsocketProxy(handler))))
}