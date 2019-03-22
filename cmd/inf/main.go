package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

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
	rootCmd.AddCommand(loginCmd)

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

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to account",
	Run: func(cmd *cobra.Command, args []string) {
		response, err := accountClient.Token(context.Background(), &apipb.TokenRequest{Username: args[0], Password: args[1]})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Login failed: %v\n", err)
			os.Exit(1)
		}

		cfg, err := config.GetCurrentContext()
		if err != nil {
			panic(err)
		}

		cfg.Token = response.Token

		tokenPayload, _ := base64.RawURLEncoding.DecodeString(strings.Split(response.GetToken(), ".")[1])

		var tokenData struct {
			AccountID string `json:"account_id"`
			DefaultNS string `json:"default_ns"`
		}

		err = json.Unmarshal([]byte(tokenPayload), &tokenData)
		if err == nil {
			config.DefaultNamespace = tokenData.DefaultNS
		}

		err = config.Write()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write config: %v\n", err)
		}

		fmt.Println("Logged in successfully.")
	},
	Args: cobra.ExactArgs(2),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return connectGRPC()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		return disconnectGRPC()
	},
}
