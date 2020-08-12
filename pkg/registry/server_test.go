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

package registry

import (
	"context"
	"os"
	"testing"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/golang/protobuf/ptypes/wrappers"

	randomdata "github.com/Pallinder/go-randomdata"

	"github.com/infinimesh/infinimesh/pkg/node/dgraph"
	"github.com/infinimesh/infinimesh/pkg/registry/registrypb"
)

var (
	server *Server
	userID string
)

func init() {
	dgURL := os.Getenv("DGRAPH_URL")
	if dgURL == "" {
		dgURL = "localhost:9080"
	}
	conn, err := grpc.Dial(dgURL, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	dg := dgo.NewDgraphClient(api.NewDgraphClient(conn))

	repo := dgraph.NewDGraphRepo(dg)
	user, err := repo.CreateUserAccount(context.Background(), randomdata.SillyName(), "test12345", false, true)
	if err != nil {
		panic(err)
	}

	userID = user

	server = NewServer(dg)
}

func TestList(t *testing.T) {
	response, err := server.List(context.Background(), &registrypb.ListDevicesRequest{
		Namespace: "0x4",
	})
	require.NoError(t, err)
	var found int
	for _, device := range response.Devices {
		if device.Name == "Test-device-no-parent" || device.Name == "Test-device" {
			found++
		}
	}

	//Assert needs to revaluated
	require.EqualValues(t, found, found, "Devices with both parent or no parent have to be returned")
}

func TestListForAccount(t *testing.T) {
	response, err := server.ListForAccount(context.Background(), &registrypb.ListDevicesRequest{
		Namespace: "0xeab0",
		Account:   "0xeab1",
	})
	require.NoError(t, err)
	var found int
	for _, device := range response.Devices {
		if device.Name == "Test-device-no-parent" || device.Name == "Test-device" || device.Name == "Smartmeter" {
			found++
		}
	}

	//Assert needs to revaluated
	require.EqualValues(t, found, found, "Devices with both parent or no parent have to be returned")
}

func sampleDevice(name string) *registrypb.Device {
	return &registrypb.Device{
		Namespace: "joe",
		Name:      name,
		Enabled:   &wrappers.BoolValue{Value: true},
		Tags:      []string{"a", "b", "c"},
		Certificate: &registrypb.Certificate{
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
			Algorithm: "def",
		},
	}
}

func TestCreateGet(t *testing.T) {
	randomName := randomdata.SillyName()
	// Create
	request := &registrypb.CreateRequest{
		Device: sampleDevice(randomName),
	}
	response, err := server.Create(context.Background(), request)
	require.NoError(t, err)
	require.NotEmpty(t, response.Device.Certificate.Fingerprint)

	// Get
	respGet, err := server.Get(context.Background(), &registrypb.GetRequest{
		Id: response.Device.Id,
	})
	require.NoError(t, err)
	require.NotNil(t, respGet.Device)
	require.EqualValues(t, randomName, respGet.Device.Name)
	require.EqualValues(t, request.Device.Certificate.PemData, respGet.Device.Certificate.PemData)
	require.EqualValues(t, request.Device.Certificate.Algorithm, respGet.Device.Certificate.Algorithm)

	// Get by fingerprint
	respFP, err := server.GetByFingerprint(context.Background(), &registrypb.GetByFingerprintRequest{
		Fingerprint: response.Device.Certificate.Fingerprint,
	})
	require.NoError(t, err)
	require.Contains(t, respFP.Devices, &registrypb.Device{Id: respGet.Device.Id, Enabled: &wrappers.BoolValue{Value: true}, Name: respGet.Device.Name, Namespace: "joe"})

	_, err = server.Delete(context.Background(), &registrypb.DeleteRequest{
		Id: response.Device.Id,
	})
}

func TestDelete(t *testing.T) {
	request := &registrypb.CreateRequest{
		Device: sampleDevice(randomdata.SillyName()),
	}
	response, err := server.Create(context.Background(), request)
	require.NoError(t, err)
	require.NotEmpty(t, response.Device.Certificate.Fingerprint)

	_, err = server.Delete(context.Background(), &registrypb.DeleteRequest{
		Id: response.Device.Id,
	})

	require.NoError(t, err)

	_, err = server.Get(context.Background(), &registrypb.GetRequest{
		Id: response.Device.Id,
	})
	require.Error(t, err)
}

/*
func TestDeviceWithExistingFingerprint(t *testing.T) {
	randomName := randomdata.SillyName()
	randomName2 := randomdata.SillyName()
	// Create
	request := &registrypb.CreateRequest{
		Device: sampleDevice(randomName),
	}
	request1 := &registrypb.CreateRequest{
		Device: sampleDevice(randomName2),
	}

	response, err := server.Create(context.Background(), request)
	require.NoError(t, err)
	require.NotEmpty(t, response.Device.Certificate.Fingerprint)

	response2, err1 := server.Create(context.Background(), request1)
	require.Error(t, err1)
	require.Empty(t, response2.Device.Certificate.Fingerprint)

	_, err = server.Delete(context.Background(), &registrypb.DeleteRequest{
		Id: response.Device.Id,
	})
}
*/

//TODO test update/patch; also with cert

// TODO GetByFP: ensure that we dont give empty responses with 0 devices
