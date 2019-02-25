package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/infinimesh/infinimesh/pkg/apiserver/apipb"
	"github.com/infinimesh/infinimesh/pkg/node/nodepb"
)

func init() {
	accountCmd.AddCommand(accountLoginCmd)
	accountCmd.AddCommand(accountCreateCmd)
	accountCmd.AddCommand(accountListCmd)
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

		cfg := &Config{
			Token: response.Token,
		}

		tokenPayload, _ := base64.RawURLEncoding.DecodeString(strings.Split(response.GetToken(), ".")[1])

		var tokenData struct {
			AccountID string `json:"account_id"`
			DefaultNS string `json:"default_ns"`
		}

		err = json.Unmarshal([]byte(tokenPayload), &tokenData)
		if err == nil {
			cfg.DefaultNamespace = tokenData.DefaultNS
		}

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

var accountListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List user accounts",
	Aliases: []string{"ls"},
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		w := tabwriter.NewWriter(os.Stdout, tabwriterMinWidth, tabwriterWidth, tabwriterPadding, tabwriterPadChar, tabwriterFlags)
		defer w.Flush()

		response, err := accountClient.ListAccounts(ctx, &nodepb.ListAccountsRequest{})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to list accounts: %v\n", err)
			os.Exit(1)
		}

		if !noHeaderFlag {
			fmt.Fprintf(w, "NAME\tID\t\n")
		}

		for _, account := range response.GetAccounts() {
			fmt.Fprintf(w, "%v\t%v\n", account.GetName(), account.GetUid())
		}

	},
}
