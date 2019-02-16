package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"google.golang.org/grpc/metadata"

	"github.com/infinimesh/infinimesh/pkg/apiserver/apipb"
)

var (
	namespaceClient apipb.NamespaceServiceClient
	accountClient   apipb.AccountClient
	ctx             context.Context

	noHeaderFlag bool
)

var rootCmd = &cobra.Command{
	Use:   "inf",
	Short: "Official commandline client for Infinimesh IoT",
}

func init() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	namespaceClient = apipb.NewNamespaceServiceClient(conn)
	accountClient = apipb.NewAccountClient(conn)

	// Load cfg
	if cfg, err := ReadConfig(); err == nil {
		ctx = metadata.AppendToOutgoingContext(context.Background(), "authorization", "bearer "+cfg.Token)

	} else {
		ctx = context.Background()
	}

}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
