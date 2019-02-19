package main

import (
	"fmt"
	"log"
	"net"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/infinimesh/infinimesh/pkg/node/nodepb"
	"github.com/infinimesh/infinimesh/pkg/registry"
	"github.com/infinimesh/infinimesh/pkg/registry/registrypb"
)

var (
	dbAddr   string
	nodeHost string
	port     string
)

func init() {
	viper.SetDefault("DB_ADDR", "postgresql://root@localhost:26257/postgres?sslmode=disable")
	viper.SetDefault("NODE_HOST", "nodeserver:8082")
	viper.SetDefault("PORT", "8080")

	viper.AutomaticEnv()

	dbAddr = viper.GetString("DB_ADDR")
	nodeHost = viper.GetString("NODE_HOST")
	port = viper.GetString("PORT")
}

func main() {
	nodeConn, err := grpc.Dial(nodeHost, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	objectClient := nodepb.NewObjectServiceClient(nodeConn)

	server := registry.NewServer(dbAddr, objectClient)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
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
