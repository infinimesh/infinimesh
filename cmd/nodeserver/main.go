package main

import (
	"context"
	"net"
	"syscall"

	"os"
	"os/signal"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/infinimesh/infinimesh/pkg/log"
	"github.com/infinimesh/infinimesh/pkg/node"
	"github.com/infinimesh/infinimesh/pkg/node/nodepb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	dgraphURL   = "localhost:9080"
	grpcAddress = ":8082"
)

func main() {
	log, err := log.NewProdOrDev()
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = log.Sync()
	}()

	log.Info("Connecting to dgraph", zap.String("URL", dgraphURL))
	conn, err := grpc.Dial(dgraphURL, grpc.WithInsecure())
	if err != nil {
		log.Fatal("Failed to connect to dgraph", zap.Error(err))
	}
	log.Info("Connected to dgraph")
	defer conn.Close()

	dg := dgo.NewDgraphClient(api.NewDgraphClient(conn))

	lis, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		log.Fatal("Failed to listen", zap.String("address", grpcAddress), zap.Error(err))
	}

	exampleAuthFunc := func(ctx context.Context) (context.Context, error) {
		token, err := grpc_auth.AuthFromMD(ctx, "bearer")
		if err != nil {
			return nil, err
		}

		log.Info("Extracted bearer token", zap.String("token", token))

		// TODO parse JWT

		newCtx := context.WithValue(ctx, node.ContextKeyAccount, "0xeabb")
		return newCtx, nil
	}

	srv := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_auth.UnaryServerInterceptor(exampleAuthFunc)),
	)

	serverHandler := &node.Server{
		Repo:   node.NewDGraphRepo(dg),
		Dgraph: dg,
		Log:    log.Named("server"),
	}

	nodepb.RegisterObjectServiceServer(srv, serverHandler)
	nodepb.RegisterAccountServer(srv, serverHandler)
	reflection.Register(srv)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT)

	go func() {
		log.Info("Starting gRPC server", zap.String("address", grpcAddress))
		if err := srv.Serve(lis); err != nil {
			log.Fatal("Failed to serve gRPC", zap.Error(err))
		}
	}()

	<-signals
	log.Info("Stopping gRPC server")
	srv.GracefulStop()
	log.Info("Stopped gRPC server")
}
