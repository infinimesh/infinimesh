package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/infinimesh/infinimesh/pkg/node/nodepb"
)

func init() {
	listNamespacesCmd.Flags().BoolVar(&noHeaderFlag, "no-headers", false, "Hide table headers")
	namespaceCmd.AddCommand(describeNamespace)
	namespaceCmd.AddCommand(listNamespacesCmd)
	namespaceCmd.AddCommand(createNamespaceCmd)
	rootCmd.AddCommand(namespaceCmd)

}

var namespaceCmd = &cobra.Command{
	Use:     "namespace",
	Aliases: []string{"ns", "namespaces"},
}

var createNamespaceCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a namespace",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		_, err := namespaceClient.CreateNamespace(ctx, &nodepb.CreateNamespaceRequest{Name: args[0]})
		if err != nil {
			fmt.Println("grpc: failed to create namespace", err)
			os.Exit(1)
		}
		fmt.Printf("Created namespace %v.\n", args[0])
	},
}

var listNamespacesCmd = &cobra.Command{
	Use:     "list",
	Short:   "List namespaces",
	Aliases: []string{"ls"},
	Run: func(cmd *cobra.Command, args []string) {
		w := tabwriter.NewWriter(os.Stdout, tabwriterMinWidth, tabwriterWidth, tabwriterPadding, tabwriterPadChar, tabwriterFlags)
		defer w.Flush()

		response, err := namespaceClient.ListNamespaces(ctx, &nodepb.ListNamespacesRequest{})
		if err != nil {
			fmt.Println("grpc: failed to fetch data", err)
			os.Exit(1)
		}
		_ = response
		if !noHeaderFlag {
			fmt.Fprintf(w, "NAME\tID\t\n")
		}

		for _, ns := range response.Namespaces {
			fmt.Fprintf(w, "%v\t%v\n", ns.GetName(), ns.GetId())
		}
	},
}

var describeNamespace = &cobra.Command{
	Use:     "describe",
	Short:   "Describe namespace",
	Aliases: []string{"desc"},
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("desc")
		w := tabwriter.NewWriter(os.Stdout, tabwriterMinWidth, tabwriterWidth, tabwriterPadding, tabwriterPadChar, tabwriterFlags)

		response, err := namespaceClient.GetNamespace(ctx, &nodepb.GetNamespaceRequest{Namespace: args[0]})
		if err != nil {
			fmt.Println("grpc: failed to fetch data", err)
			os.Exit(1)
		}
		fmt.Println(response)
		// _ = response
		// if !noHeaderFlag {
		// 	fmt.Fprintf(w, "NAME\tID\t\n")
		// }

		// for _, ns := range response.Namespaces {
		// 	fmt.Fprintf(w, "%v\t%v\n", ns.GetName(), ns.GetId())
		// }

		defer w.Flush()
	},
}
