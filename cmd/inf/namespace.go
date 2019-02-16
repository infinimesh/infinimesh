package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/infinimesh/infinimesh/pkg/node/nodepb"
)

func init() {
	namespacesCmd.Flags().BoolVar(&noHeaderFlag, "no-headers", false, "Hide table headers")
	rootCmd.AddCommand(namespacesCmd)

}

var namespacesCmd = &cobra.Command{
	Use:     "namespaces",
	Short:   "List namespaces",
	Aliases: []string{"ns"},
	Run: func(cmd *cobra.Command, args []string) {
		w := tabwriter.NewWriter(os.Stdout, tabwriterMinWidth, tabwriterWidth, tabwriterPadding, tabwriterPadChar, tabwriterFlags)

		response, err := namespaceClient.ListNamespaces(ctx, &nodepb.ListNamespacesRequest{})
		if err != nil {
			fmt.Println("grpc: failed to fetch data", err)
		}
		_ = response
		if !noHeaderFlag {
			fmt.Fprintf(w, "NAME\tID\t\n")
		}

		for _, ns := range response.Namespaces {
			fmt.Fprintf(w, "%v\t%v\n", ns.GetName(), ns.GetId())
		}

		defer w.Flush()
	},
}
