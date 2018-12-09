package main

import (
	"net"

	"github.com/infinimesh/infinimesh/pkg/apiserver/apipb"
	"github.com/infinimesh/infinimesh/pkg/registry/registrypb"
	"github.com/infinimesh/infinimesh/pkg/shadow/shadowpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	addr         = ":8080"
	registryHost = "device-registry:8080"
	shadowHost   = "shadow-api:8096"
)

func main() {
	srv := grpc.NewServer()

	registryConn, err := grpc.Dial(registryHost, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	devicesClient := registrypb.NewDevicesClient(registryConn)

	shadowConn, err := grpc.Dial(shadowHost, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	shadowClient := shadowpb.NewShadowClient(shadowConn)

	apipb.RegisterDevicesServer(srv, &deviceAPI{client: devicesClient})
	apipb.RegisterShadowServer(srv, &shadowAPI{client: shadowClient})
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	reflection.Register(srv)

	err = srv.Serve(listener)
	if err != nil {
		panic(err)
	}

}
