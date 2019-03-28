package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	prettyjson "github.com/hokaccha/go-prettyjson"
	"github.com/spf13/cobra"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"

	"github.com/infinimesh/infinimesh/pkg/shadow/shadowpb"
)

func init() {
	stateCmd.AddCommand(stateGetCmd)
	rootCmd.AddCommand(stateCmd)
}

var stateCmd = &cobra.Command{
	Use: "state",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return connectGRPC()
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		return disconnectGRPC()
	},
}

var stateGetCmd = &cobra.Command{
	Use:  "get <deviceID>",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		w := tabwriter.NewWriter(os.Stdout, tabwriterMinWidthNested, 4, 2, tabwriterPadChar, tabwriterFlags)
		defer w.Flush()

		response, err := shadowClient.Get(ctx, &shadowpb.GetRequest{
			Id: args[0],
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get state: %v\n", err)
			os.Exit(1)
		}

		fmt.Fprintf(w, "Reported State:")
		printState(w, response.Shadow.Reported)

		fmt.Fprintf(w, "Desired State:")
		printState(w, response.Shadow.Desired)

		fmt.Fprintf(w, "Configuration:")
		printState(w, response.Shadow.Config)
	},
}

func printState(w io.Writer, state *shadowpb.VersionedValue) {
	if state == nil || state.Data == nil || state.Version == 0 {
		fmt.Fprintln(w, " <none>")
		return
	}
	fmt.Fprintf(w, "\n")
	buf := &bytes.Buffer{}
	m := jsonpb.Marshaler{}
	err := m.Marshal(buf, state.Data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get state: %v\n", err)
		os.Exit(1)
	}

	formatter := prettyjson.NewFormatter()
	formatter.Newline = "\n\t\t"
	pretty, err := formatter.Format(buf.Bytes())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get state: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(w, "\tVersion:\t%v\n", state.Version)
	fmt.Fprintf(w, "\tTimestamp:\t%v\n", convertTimestamp(state.Timestamp))
	fmt.Fprintf(w, "\tData:\n\t\t%v\n", string(pretty))

}

func convertTimestamp(ts *timestamp.Timestamp) string {
	res, err := ptypes.Timestamp(ts)
	if err != nil {
		return "<unknown>"
	}

	return res.Local().String()

}
