package main

import (
	"net"

	"github.com/infinimesh/infinimesh/pkg/apiserver/apipb"
	"github.com/infinimesh/infinimesh/pkg/registry/registrypb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	addr         = ":8080"
	registryHost = "device-registry:8080"
)

func main() {
	srv := grpc.NewServer()

	conn, err := grpc.Dial(registryHost, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	devicesClient := registrypb.NewDevicesClient(conn)

	apipb.RegisterDevicesServer(srv, &deviceAPI{client: devicesClient})

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
