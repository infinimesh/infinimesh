//--------------------------------------------------------------------------
// Copyright 2018 infinimesh
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
	"net/http"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"github.com/spf13/viper"
	"github.com/tmc/grpc-websocket-proxy/wsproxy"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/slntopp/infinimesh/pkg/apiserver/apipb"
)

var (
	apiserverEndpoint string
)

func init() {
	viper.SetDefault("APISERVER_ENDPOINT", "localhost:8080")
	viper.AutomaticEnv()
	apiserverEndpoint = viper.GetString("APISERVER_ENDPOINT")
}

func main() {
	defer glog.Flush()

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux(
		runtime.WithDisablePathLengthFallback(),
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true, EmitUnpopulated: true,
		}}),
	)
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := apipb.RegisterDevicesHandlerFromEndpoint(ctx, mux, apiserverEndpoint, opts)
	if err != nil {
		return err
	}

	err = apipb.RegisterStatesHandlerFromEndpoint(ctx, mux, apiserverEndpoint, opts)
	if err != nil {
		return err
	}

	err = apipb.RegisterAccountsHandlerFromEndpoint(ctx, mux, apiserverEndpoint, opts)
	if err != nil {
		return err
	}

	err = apipb.RegisterObjectsHandlerFromEndpoint(ctx, mux, apiserverEndpoint, opts)
	if err != nil {
		return err
	}

	err = apipb.RegisterNamespacesHandlerFromEndpoint(ctx, mux, apiserverEndpoint, opts)
	if err != nil {
		return err
	}

	corsMiddleware := cors.AllowAll().Handler(mux)
	return http.ListenAndServe(":8081", wsproxy.WebsocketProxy(corsMiddleware))
}
