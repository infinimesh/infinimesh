package main

import (
	"context"
	"flag"
	"net/http"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/infinimesh/infinimesh/pkg/apiserver/apipb"
	"google.golang.org/grpc"
)

var (
	apiserverEndpoint = flag.String("apiserver", "localhost:8080", "gRPC APIServer Host:Port pair")
)

func main() {
	flag.Parse()
	defer glog.Flush()

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := apipb.RegisterDevicesHandlerFromEndpoint(ctx, mux, *apiserverEndpoint, opts)
	if err != nil {
		return err
	}

	err = apipb.RegisterShadowHandlerFromEndpoint(ctx, mux, *apiserverEndpoint, opts)
	if err != nil {
		return err
	}

	return http.ListenAndServe(":8081", mux)
}
