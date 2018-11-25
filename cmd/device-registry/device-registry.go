package main

import (
	"log"
	"net"

	"github.com/infinimesh/infinimesh/pkg/registry"
	"github.com/infinimesh/infinimesh/pkg/registry/registrypb"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const port = ":8080"

var dbAddr string

func init() {
	viper.SetDefault("DB_HOST", "localhost")
	viper.AutomaticEnv()

	dbAddr = viper.GetString("DB_ADDR")
}

func main() {
	server := registry.NewServer(dbAddr)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	registrypb.RegisterDevicesServer(s, server)
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
