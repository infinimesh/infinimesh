//--------------------------------------------------------------------------
// Copyright 2018 Infinite Devices GmbH
// www.infinimesh.io
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.
//--------------------------------------------------------------------------

package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/spf13/cobra"

	"io/ioutil"

	"github.com/infinimesh/infinimesh/pkg/apiserver/apipb"
	"github.com/infinimesh/infinimesh/pkg/registry/registrypb"
)

var (
	allNamespaces bool
	certFile      string
)

func init() {
	createDeviceCmd.Flags().StringVarP(&namespaceFlag, "namespace", "n", "", "Namespace")
	lsDeviceCmd.Flags().StringVarP(&namespaceFlag, "namespace", "n", "", "Namespace")
	lsDeviceCmd.Flags().BoolVar(&allNamespaces, "all-namespaces", false, "Show devices in all namespaces")
	devicesCmd.AddCommand(lsDeviceCmd)
	devicesCmd.AddCommand(createDeviceCmd)
	devicesCmd.AddCommand(deleteDeviceCmd)
	rootCmd.AddCommand(devicesCmd)

	createDeviceCmd.Flags().StringVar(&certFile, "cert-file", "", "Path to X509 certificate file of device")
}

var devicesCmd = &cobra.Command{
	Use:     "device",
	Aliases: []string{"devices", "dev"},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return connectGRPC()
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		return disconnectGRPC()
	},
}

var lsDeviceCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Run: func(cmd *cobra.Command, args []string) {
		w := tabwriter.NewWriter(os.Stdout, tabwriterMinWidth, tabwriterWidth, tabwriterPadding, tabwriterPadChar, tabwriterFlags)
		defer w.Flush()

		if namespace := getNamespace(); namespace == "" {
			fmt.Println("Please set namespace for you id.")
			os.Exit(1)
		}

		response, err := deviceClient.List(ctx, &apipb.ListDevicesRequest{
			Namespaceid: getNamespace(),
		})
		if err != nil {
			fmt.Println("grpc: failed to fetch data", err)
			os.Exit(1)
		}

		if !noHeaderFlag {
			fmt.Fprintf(w, "ID\tNAME\tENABLED\t\n")
		}

		for _, device := range response.Devices {
			fmt.Fprintf(w, "%v\t%v\t%v\t\n", device.Id, device.Name, device.Enabled)
		}

	},
}

var createDeviceCmd = &cobra.Command{
	Use:  "create",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pem, err := ioutil.ReadFile(certFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not read certificate: %v", err)
			os.Exit(1)
		}

		resp, err := deviceClient.Create(ctx, &registrypb.CreateRequest{
			Device: &registrypb.Device{
				Name:    args[0],
				Enabled: &wrappers.BoolValue{Value: true},
				Certificate: &registrypb.Certificate{
					PemData: string(pem),
				},
				Namespace: getNamespace(),
			},
		})
		if err != nil {
			panic(err)
		}
		fmt.Printf("Created device [%v].\n", resp.Device.Id)
	},
}

var deleteDeviceCmd = &cobra.Command{
	Use:  "delete",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		_, err := deviceClient.Delete(ctx, &registrypb.DeleteRequest{
			Id: args[0],
		})
		if err != nil {
			fmt.Println("grpc: failed to delete device", err)
		}
		fmt.Println("Deleted device.")
	},
}
