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

	"strings"

	"github.com/infinimesh/infinimesh/pkg/shadow/shadowpb"

	_struct "github.com/golang/protobuf/ptypes/struct"
)

var (
	watch bool
)

func init() {
	stateCmd.AddCommand(stateGetCmd)
	stateCmd.AddCommand(stateSetCmd)
	rootCmd.AddCommand(stateCmd)

	stateGetCmd.Flags().BoolVarP(&watch, "watch", "w", false, "Watch for changes")
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

var stateSetCmd = &cobra.Command{
	Use:  "set <deviceID> <JSON State>",
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		u := &jsonpb.Unmarshaler{}
		var state _struct.Value
		err := u.Unmarshal(strings.NewReader(args[1]), &state)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid state: %v\n", err)
		}

		_, err = shadowClient.PatchDesiredState(ctx, &shadowpb.PatchDesiredStateRequest{
			Id:   args[0],
			Data: &state,
		})

		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to set state: %v", err)
		}
		fmt.Println("Successfully set state.")
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

		if watch {
			for {
				resp, err := shadowClient.StreamReportedStateChanges(ctx, &shadowpb.StreamReportedStateChangesRequest{
					Id: args[0],
				})
				if err != nil {
					fmt.Fprintf(os.Stderr, "Failed to get state: %v\n", err)
					os.Exit(1)
				}

				for {
					msg, err := resp.Recv()
					if err != nil {
						break
					}

					if msg.ReportedState != nil {
						fmt.Fprintf(w, "Reported State:")
						printState(w, msg.ReportedState)
					}

					if msg.DesiredState != nil {
						fmt.Fprintf(w, "Desired State:")
						printState(w, msg.DesiredState)
					}

					w.Flush()
				}
			}

		}
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
