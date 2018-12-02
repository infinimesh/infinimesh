package main

import (
	"net"

	"github.com/infinimesh/infinimesh/pkg/apiserver/apipb"
	"google.golang.org/grpc"
)

var addr = ":8080"

func main() {
	srv := grpc.NewServer()

	apipb.RegisterDevicesServer(srv, &deviceAPI{})

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	err = srv.Serve(listener)
	if err != nil {
		panic(err)
	}
}
