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
	namespaceClient apipb.NamespacesClient
	objectClient    apipb.ObjectsClient
	accountClient   apipb.AccountsClient
	deviceClient    apipb.DevicesClient
	ctx             context.Context

	noHeaderFlag bool

	namespaceFlag string

	config *Config
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

	namespaceClient = apipb.NewNamespacesClient(conn)
	accountClient = apipb.NewAccountsClient(conn)
	deviceClient = apipb.NewDevicesClient(conn)
	objectClient = apipb.NewObjectsClient(conn)

	// Load cfg
	if cfg, err := ReadConfig(); err == nil {
		ctx = metadata.AppendToOutgoingContext(context.Background(), "authorization", "bearer "+cfg.Token)
		config = cfg

	} else {
		ctx = context.Background()
	}

}

func getNamespace() string {
	if allNamespaces {
		return ""
	}
	if namespaceFlag != "" {
		return namespaceFlag
	}
	if config != nil && config.DefaultNamespace != "" {
		return config.DefaultNamespace
	}
	return ""
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
