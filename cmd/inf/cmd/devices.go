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
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/infinimesh/infinimesh/pkg/convert"
	pb "github.com/infinimesh/infinimesh/pkg/node/proto"
	devpb "github.com/infinimesh/infinimesh/pkg/node/proto/devices"
	"github.com/jedib0t/go-pretty/v6/table"

	"github.com/spf13/cobra"
)

func makeDevicesServiceClient(ctx context.Context) (pb.DevicesServiceClient, error) {
	conn, err := makeConnection(ctx)
	if err != nil {
		return nil, err
	}
	return pb.NewDevicesServiceClient(conn), nil
}

// devicesCmd represents the devices command
var devicesCmd = &cobra.Command{
	Use:   "devices",
	Short: "Manage infinimesh devices",
	Aliases: []string{"device", "dev"},
	RunE: listDevicesCmd.RunE,
}

var listDevicesCmd = &cobra.Command{
	Use:   "list",
	Short: "List infinimesh devices",
	Aliases: []string{"ls", "l"},
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := makeContextWithBearerToken()
		client, err := makeDevicesServiceClient(ctx)
		if err != nil {
			return err
		}

		r, err := client.List(ctx, &pb.EmptyMessage{})
		if err != nil {
			return err
		}

		if printJson, _ := cmd.Flags().GetBool("json"); printJson {
			return printJsonResponse(r)
		}

		PrintDevicesPool(r.Devices)
		return nil
	},
}

var getDeviceCmd = &cobra.Command{
	Use:   "get",
	Short: "Get infinimesh device",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := makeContextWithBearerToken()
		client, err := makeDevicesServiceClient(ctx)
		if err != nil {
			return err
		}

		r, err := client.Get(ctx, &devpb.Device{Uuid: args[0]})
		if err != nil {
			return err
		}

		if printJson, _ := cmd.Flags().GetBool("json"); printJson {
			return printJsonResponse(r)
		}

		PrintSingleDevice(r)
		return nil
	},
}

var makeDeviceTokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Make device token",
	Aliases: []string{"tok", "t"},
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := makeContextWithBearerToken()
		client, err := makeDevicesServiceClient(ctx)
		if err != nil {
			return err
		}

		allowPost, _ := cmd.Flags().GetBool("allow-post")
		r, err := client.MakeDevicesToken(ctx, &pb.DevicesTokenRequest{
			Devices: args,
			Post: allowPost,
		})
		if err != nil {
			return err
		}

		if printJson, _ := cmd.Flags().GetBool("json"); printJson {
			return printJsonResponse(r)
		}

		fmt.Println(r.Token)
		return nil
	},
}

var createDeviceCmd = &cobra.Command{
	Use:   "create",
	Short: "Create infinimesh device",
	Aliases: []string{"add", "a", "new", "crt"},
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := makeContextWithBearerToken()
		client, err := makeDevicesServiceClient(ctx)
		if err != nil {
			return err
		}

		if _, err := os.Stat(args[0]); os.IsNotExist(err) {
			return errors.New("Template doesn't exist at path " + args[0])
		}

		var format string
		{
			pathSlice := strings.Split(args[0], ".")
			format = pathSlice[len(pathSlice) - 1]
		}

		template, err := os.ReadFile(args[0])
		if err != nil {
			fmt.Println("Error while reading template")
			return err
		}

		switch format {
		case "json":
		case "yml", "yaml":
			template, err = convert.ConvertBytes(template)
		default:
			return errors.New("Unsupported template format " + format)
		}

		if err != nil {
			fmt.Println("Error while parsing template")
			return err
		}
		
		fmt.Println("Template", string(template))

		var device devpb.Device
		err = json.Unmarshal(template, &device)
		if err != nil {
			fmt.Println("Error while parsing template")
			return err
		}

		certPath, _ := cmd.Flags().GetString("crt")
		if _, err := os.Stat(certPath); os.IsNotExist(err) {
			return errors.New("Certificate doesn't exist at path " + certPath)
		}
		pem, err := os.ReadFile(certPath)
		if err != nil {
			fmt.Println("Error while reading certificate")
			return err
		}

		cert := &devpb.Certificate{
			PemData: string(pem),
		}
		device.Certificate = cert

		ns, _ := cmd.Flags().GetString("namespace")

		res, err := client.Create(ctx, &devpb.CreateRequest{
			Device: &device,
			Namespace: ns,
		})
		if err != nil {
			return err
		}

		fmt.Println("Device Created, UUID:", res.Device.Uuid)
		return nil
	},
}

func PrintSingleDevice(d *devpb.Device) {
	fmt.Printf("UUID: %s\n", d.Uuid)
	fmt.Printf("Title: %s\n", d.Title)
	fmt.Printf("Enabled: %t\n", d.Enabled)

	tags := strings.Join(d.Tags, ",")
	if tags == "" {
		tags = "-"
	}
	fmt.Printf("Tags: %s\n", tags)

	fingerprint := hex.EncodeToString(d.Certificate.Fingerprint)
	fmt.Printf("Fingerprint:\n  Algorythm: %s\n  Hash: %s\n", d.Certificate.Algorithm, fingerprint)
}

func PrintDevicesPool(pool []*devpb.Device) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"UUID", "Title", "Enabled", "Tags"})

	rows := make([]table.Row, len(pool))
	for i, dev := range pool {
		tags := strings.Join(dev.Tags, ",")
		if tags == "" {
			tags = "-"
		}
		rows[i] = table.Row{dev.Uuid, dev.Title, dev.Enabled, tags}
	}
	t.AppendRows(rows)

	t.SortBy([]table.SortBy{
		{Name: "UUID", Mode: table.Asc},
	})

	t.AppendFooter(table.Row{"", "", "Total Found", len(pool)}, table.RowConfig{AutoMerge: true})
	t.Render()
}

func init() {
	devicesCmd.AddCommand(listDevicesCmd)
	devicesCmd.AddCommand(getDeviceCmd)

	makeDeviceTokenCmd.Flags().Bool("allow-post", false, "Allow posting devices states")
	devicesCmd.AddCommand(makeDeviceTokenCmd)

	createDeviceCmd.Flags().String("crt", "", "Path to certificate file")
	createDeviceCmd.Flags().StringP("namespace", "n", "", "Namespace to create device in")
	devicesCmd.AddCommand(createDeviceCmd)

	rootCmd.AddCommand(devicesCmd)
}
