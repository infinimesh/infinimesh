/*
Copyright © 2021-2022 Nikita Ivanovski info@slnt-opp.xyz

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/infinimesh/infinimesh/pkg/node/proto"
	accpb "github.com/infinimesh/infinimesh/pkg/node/proto/accounts"
)

var VERSION string

func getVersion() string {
	if VERSION == "" {
		return "dev"
	}
	return VERSION
}

// contextCmd represents the context command
var contextCmd = &cobra.Command{
	Use:   "context",
	Aliases: []string{"ctx"},
	Short: "Print current infinimesh CLI Context | Aliases: ctx",
	RunE: func(cmd *cobra.Command, args []string) error {
		data := make(map[string]interface{})
		data["version"] = getVersion()

		data["host"] = viper.Get("infinimesh")
		if data["host"] == nil {
			data = map[string]interface{}{
				"error": "No infinimesh context found",
			}
		}

		if insec := viper.GetBool("insecure"); insec {
			data["insecure"] = insec
		}

		if printJson, _ := cmd.Flags().GetBool("json"); printJson {
			data, err := json.Marshal(data)
			if err != nil {
				return err
			}
			fmt.Println(string(data))
			return nil
		}

		for k, v := range data {
			fmt.Printf("%s: %v\n", strings.Title(k), v)
		}

		return nil
	},
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Aliases: []string{"l", "auth", "a"},
	Short: "Authorize in infinimesh and store credentials",
	RunE: func(cmd *cobra.Command, args []string) error {
		creds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})
		insec, _ := cmd.Flags().GetBool("insecure")
		if insec {
			creds = insecure.NewCredentials()
		}
		conn, err := grpc.Dial(args[0], grpc.WithTransportCredentials(creds))
		if err != nil {
			return err
		}

		client := pb.NewAccountsServiceClient(conn)
		req := &pb.TokenRequest{
			Auth: &accpb.Credentials{
				Type: "standard", Data: []string{args[1], args[2]},
			},
		}

		res, err := client.Token(context.Background(), req)
		if err != nil {
			return err
		}
		token := res.GetToken()
		printToken, _ := cmd.Flags().GetBool("print-token")
		if printToken {
			fmt.Println(token)
		}

		viper.Set("infinimesh", args[0])
		viper.Set("token", token)
		viper.Set("insecure", insec)

		err = viper.WriteConfig()
		return err
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print infinimesh CLI version",
	RunE: func(cmd *cobra.Command, args []string) error {
		if printJson, _ := cmd.Flags().GetBool("json"); printJson {
			data, err := json.Marshal(map[string]string{
				"version": getVersion(),
			})
			if err != nil {
				return err
			}
			fmt.Println(string(data))
			return nil
		}

		fmt.Println("CLI Version:", getVersion())
		return nil
	},
}


func init() {
	loginCmd.Flags().Bool("print-token", false, "")
	loginCmd.Flags().Bool("insecure", false, "Use WithInsecure instead of TLS")

	rootCmd.AddCommand(contextCmd)
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(versionCmd)
}
