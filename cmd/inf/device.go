package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/spf13/cobra"

	"encoding/base64"

	"github.com/infinimesh/infinimesh/pkg/apiserver/apipb"
	"github.com/infinimesh/infinimesh/pkg/registry/registrypb"
)

func init() {
	createDeviceCmd.Flags().StringVarP(&namespaceFlag, "namespace", "n", "", "Namespace")
	lsDeviceCmd.Flags().StringVarP(&namespaceFlag, "namespace", "n", "", "Namespace")
	devicesCmd.AddCommand(lsDeviceCmd)
	devicesCmd.AddCommand(createDeviceCmd)
	rootCmd.AddCommand(devicesCmd)
}

var devicesCmd = &cobra.Command{
	Use:     "device",
	Aliases: []string{"devices", "dev"},
}

var lsDeviceCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Run: func(cmd *cobra.Command, args []string) {
		w := tabwriter.NewWriter(os.Stdout, tabwriterMinWidth, tabwriterWidth, tabwriterPadding, tabwriterPadChar, tabwriterFlags)
		defer w.Flush()

		fmt.Println("ns", namespaceFlag)
		response, err := objectClient.ListObjects(ctx, &apipb.ListObjectsRequest{
			Namespace: namespaceFlag,
		})
		if err != nil {
			fmt.Println("grpc: failed to fetch data", err)
			os.Exit(1)
		}

		if !noHeaderFlag {
			fmt.Fprintf(w, "NAME\tID\t\n")
		}

		for _, object := range response.GetObjects() {
			fmt.Fprintf(w, "%v\t%v\n", object.GetName(), object.GetUid())
		}

	},
}

var createDeviceCmd = &cobra.Command{
	Use:  "create",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := deviceClient.Create(ctx, &registrypb.CreateRequest{
			Device: &registrypb.Device{
				Id:      args[0],
				Enabled: &wrappers.BoolValue{Value: true},
				Certificate: &registrypb.Certificate{
					// TODO cert, don't hardcode ;)
					PemData: `-----BEGIN CERTIFICATE-----
MIIDiDCCAnCgAwIBAgIJAMNNOKhM9eyOMA0GCSqGSIb3DQEBCwUAMFkxCzAJBgNV
BAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYDVQQKDBhJbnRlcm5ldCBX
aWRnaXRzIFB0eSBMdGQxEjAQBgNVBAMMCWxvY2FsaG9zdDAeFw0xODA4MDYyMTU4
NTRaFw0yODA4MDMyMTU4NTRaMFkxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21l
LVN0YXRlMSEwHwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQxEjAQBgNV
BAMMCWxvY2FsaG9zdDCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBALq2
5T2k9R98jWmGXjeFr+iutigtuwI9TQ5CQ1+2Rh9sYpEzyZSeHm2/keMmhfuLD9vv
qN6kHWWArmqLFGZ7MM28wpsXOxMgK5UClmYb95jYUemKQn6opSYCnapvUj6UhuBo
cpg7m6eLysG0WMQZAo1LC2eMIQGTCBmXuVFakRL+0CFjaD5d4+VJUKhvMPM5xpty
qD2Bk9KXNHgS8uX8Yxxe0tB+p6P60Kgv9+yWCrm2RUV/zuSlXX69nUE/VrezSdGn
c/tVSIcspiXTpDlKiHLPoYfL83xwMrwg4Y1EUTDzkAku98upss+GDalkJaSldy67
JJLTs94ZgG5vJTZPJe0CAwEAAaNTMFEwHQYDVR0OBBYEFJOEmob6pthnFZq2lZzf
38wfQZhpMB8GA1UdIwQYMBaAFJOEmob6pthnFZq2lZzf38wfQZhpMA8GA1UdEwEB
/wQFMAMBAf8wDQYJKoZIhvcNAQELBQADggEBAJUiAGJQbHPMeYWi4bOhsuUrvHhP
mN/g4nwtjkAiu6Q5QOHy1xVdGzR7u6rbHZFMmdIrUPQ/5mkqJdZndl5WShbvaG/8
I0U3Uq0B3Xuf0f1Pcn25ioTj+U7PIUYqWQXvjN1YnlsUjcbQ7CQ2EOHKmNA7v2fg
OmWrBAp4qqOaEKWpg0N9fZICb7g4klONQOryAaZYcbeCBwXyg0baCZLXfJzatn41
Xkrr0nVweXiEEk5BosN20FyFZBekpby11th2M1XksArLTWQ41IL1TfWKJALDZgPL
AX99IKELzVTsndkfF8mLVWZr1Oob7soTVXfOI/VBn1e+3qkUrK94JYtYj04=
-----END CERTIFICATE-----`,
				},
			},
			Namespace: namespaceFlag,
		})
		if err != nil {
			panic(err)
		}
		fmt.Printf("Created device.\nFingerprint: %v\n", base64.StdEncoding.EncodeToString(resp.GetFingerprint()))
	},
}
