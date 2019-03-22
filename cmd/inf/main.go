package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"

	"github.com/infinimesh/infinimesh/pkg/apiserver/apipb"
)

var (
	clientConn      *grpc.ClientConn
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
	// Load cfg
	if cfg, err := ReadConfig(); err == nil {
		cur, err := cfg.GetCurrentContext()
		if err == nil {
			ctx = metadata.AppendToOutgoingContext(context.Background(), "authorization", "bearer "+cur.Token)
		}
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

func connectGRPC() error {
	current, err := config.GetCurrentContext()
	if err != nil {
		return errors.New("no context found")
	}
	var option grpc.DialOption
	pool, _ := x509.SystemCertPool()
	option = grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{RootCAs: pool}))

	if !current.TLS {
		option = grpc.WithInsecure()
	}

	conn, err := grpc.Dial(current.Server, option)
	if err != nil {
		return err
	}

	clientConn = conn
	namespaceClient = apipb.NewNamespacesClient(conn)
	accountClient = apipb.NewAccountsClient(conn)
	deviceClient = apipb.NewDevicesClient(conn)
	objectClient = apipb.NewObjectsClient(conn)
	return nil
}

func disconnectGRPC() error {
	return clientConn.Close()
}

// What to do if unsuccessful connection / no current state unavailable?
func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
