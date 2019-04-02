package main

import (
	"context"
	"net/http"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/rs/cors"
	"github.com/spf13/viper"
	"github.com/tmc/grpc-websocket-proxy/wsproxy"
	"google.golang.org/grpc"

	"github.com/infinimesh/infinimesh/pkg/apiserver/apipb"
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
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{OrigName: true, EmitDefaults: true}),
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
