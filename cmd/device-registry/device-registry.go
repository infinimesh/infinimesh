package main

import (
	"fmt"
	"net"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/infinimesh/infinimesh/pkg/registry"
	"github.com/infinimesh/infinimesh/pkg/registry/registrypb"

	logger "github.com/infinimesh/infinimesh/pkg/log"
)

var (
	port string

	dgraphURL string
)

func init() {
	viper.AutomaticEnv()
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("DGRAPH_HOST", "localhost:9080")

	dgraphURL = viper.GetString("DGRAPH_HOST")
	port = viper.GetString("PORT")
}

func main() {
	log, err := logger.NewProdOrDev()
	if err != nil {
		panic(err)
	}

	log.Info("Connecting to dgraph", zap.String("URL", dgraphURL))
	conn, err := grpc.Dial(dgraphURL, grpc.WithInsecure())
	if err != nil {
		log.Fatal("Failed to connect to dgraph", zap.Error(err))
	}
	log.Info("Connected to dgraph")
	defer conn.Close()

	dg := dgo.NewDgraphClient(api.NewDgraphClient(conn))

	server := registry.NewServer(dg)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		log.Fatal("Failed to listen", zap.String("address", port), zap.Error(err))
	}
	s := grpc.NewServer()
	registrypb.RegisterDevicesServer(s, server)
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatal("Failed to serve gRPC", zap.Error(err))
	}
}
