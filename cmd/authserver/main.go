package main

import (
	"net"
	"syscall"

	"os"
	"os/signal"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"github.com/infinimesh/infinimesh/pkg/auth"
	"github.com/infinimesh/infinimesh/pkg/auth/authpb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	dgraphURL = "127.0.0.1:9080"
)

func main() {
	log, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer log.Sync()
	conn, err := grpc.Dial(dgraphURL, grpc.WithInsecure())
	if err != nil {
		log.Fatal("Failed to connect to dgraph", zap.Error(err))
	}
	defer conn.Close()

	dg := dgo.NewDgraphClient(api.NewDgraphClient(conn))

	lis, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatal("Failed to listen", zap.Error(err))
	}

	srv := grpc.NewServer()

	serverHandler := &auth.Server{
		Dgraph: dg,
	}

	authpb.RegisterAuthServer(srv, serverHandler)
	reflection.Register(srv)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT)

	go func() {
		log.Info("Starting gRPC server")
		if err := srv.Serve(lis); err != nil {
			log.Fatal("Failed to serve gRPC", zap.Error(err))
		}
	}()

	<-signals
	log.Info("Exiting")
}
