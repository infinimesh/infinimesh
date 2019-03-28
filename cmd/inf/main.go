package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	input "github.com/tcnksm/go-input"
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
	shadowClient    apipb.ShadowsClient
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
	shadowClient = apipb.NewShadowsClient(conn)
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
		scanner := bufio.NewScanner(os.Stdin)

		fmt.Print("Username: ")
		scanner.Scan()
		username := scanner.Text()

		ui := &input.UI{
			Writer: os.Stdout,
			Reader: os.Stdin,
		}

		password, err := ui.Ask("Password", &input.Options{
			Required:  true,
			Mask:      true,
			HideOrder: true,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read password: %v\n", err)
		}

		response, err := accountClient.Token(context.Background(), &apipb.TokenRequest{Username: username, Password: password})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Login failed: %v\n", err)
			os.Exit(1)
		}

		cfg, err := config.GetCurrentContext()
		if err != nil {
			panic(err)
		}

		cfg.Token = response.Token

		err = config.Write()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write config: %v\n", err)
		}

		fmt.Println("Logged in successfully.")
	},
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return connectGRPC()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		return disconnectGRPC()
	},
}
