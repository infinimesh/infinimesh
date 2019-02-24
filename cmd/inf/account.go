package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/infinimesh/infinimesh/pkg/apiserver/apipb"
	"github.com/infinimesh/infinimesh/pkg/node/nodepb"
)

func init() {
	accountCmd.AddCommand(accountLoginCmd)
	accountCmd.AddCommand(accountCreateCmd)
	rootCmd.AddCommand(accountCmd)
}

var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "Manage accounts",
}

var accountLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to account",
	Run: func(cmd *cobra.Command, args []string) {
		response, err := accountClient.Token(context.Background(), &apipb.TokenRequest{Username: args[0], Password: args[1]})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Login failed: %v\n", err)
			os.Exit(1)
		}

		cfg := &Config{Token: response.Token}
		err = cfg.Write()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write config: %v\n", err)
		}

		fmt.Println("Logged in successfully.")
	},
	Args: cobra.ExactArgs(2),
}

var accountCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create user account",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		response, err := accountClient.CreateUserAccount(ctx, &nodepb.CreateUserAccountRequest{
			Name:     args[0],
			Password: args[1],
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create user: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Created user %v with id %v.\n", args[0], response.GetUid())
	},
}
