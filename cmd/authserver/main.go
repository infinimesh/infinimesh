package main

import (
	"fmt"
	"log"
	"net"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"github.com/infinimesh/infinimesh/pkg/auth"
	"github.com/infinimesh/infinimesh/pkg/auth/authpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	dgraphURL = "127.0.0.1:9080"
)

func main() {
	conn, err := grpc.Dial(dgraphURL, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	dg := dgo.NewDgraphClient(api.NewDgraphClient(conn))

	lis, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()

	serverHandler := &auth.Server{
		Dgraph: dg,
	}

	authpb.RegisterAuthServer(srv, serverHandler)
	reflection.Register(srv)
	fmt.Println("Serving")
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
