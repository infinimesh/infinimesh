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
	"google.golang.org/genproto/protobuf/field_mask"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/golang/protobuf/ptypes/wrappers"

	randomdata "github.com/Pallinder/go-randomdata"

	"github.com/infinimesh/infinimesh/pkg/node/dgraph"
	"github.com/infinimesh/infinimesh/pkg/node/nodepb"
	"github.com/infinimesh/infinimesh/pkg/repo"

	logger "github.com/infinimesh/infinimesh/pkg/log"
	"github.com/infinimesh/infinimesh/pkg/registry/registrypb"
)

var (
	server *Server
	userID string
)

var rep1 repo.Repo

func init() {
	dgURL := os.Getenv("DGRAPH_URL")
	if dgURL == "" {
		dgURL = "localhost:9080"
	}
	dbURL := os.Getenv("DB_ADDR")
	if dbURL == "" {
		dbURL = ":6379"
	}
	r, err := repo.NewRedisRepo(dbURL)
	if err != nil {
		panic(err)
	}
	rep1 = r

	conn, err := grpc.Dial(dgURL, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	dg := dgo.NewDgraphClient(api.NewDgraphClient(conn))

	drepo := dgraph.NewDGraphRepo(dg)
	user, err := drepo.CreateUserAccount(context.Background(), randomdata.SillyName(), "test12345", false, false, true)
	if err != nil {
		panic(err)
	}

	log, err := logger.NewProdOrDev()
	if err != nil {
		panic(err)
	}

	userID = user

	server = NewServer(dg, repo.Server{
		Repo: rep1,
		Log:  log.Named("RepoController"),
	})
	server.Log = log.Named("Device Registry Test")
}

func TestList(t *testing.T) {

	ctx := context.Background()

	randomName := randomdata.SillyName()

	accid, err := server.repo.CreateUserAccount(ctx, randomName, "password", true, false, true)
	require.NoError(t, err)

	ctx = metadata.NewIncomingContext(ctx, metadata.New(map[string]string{"requestorid": accid}))

	ns, err := server.repo.GetNamespace(ctx, randomName)
	require.NoError(t, err)

	// Create
	request := &registrypb.CreateRequest{
		Device: sampleDevice(randomName, ns.Id),
	}
	resp, err := server.Create(ctx, request)
	require.NoError(t, err)

	response, err := server.List(ctx, &registrypb.ListDevicesRequest{
		Namespaceid: ns.Id,
	})
	require.NoError(t, err)
	var found int
	for _, device := range response.Devices {
		if device.Name == randomName {
			found++
		}
	}

	//Assert needs to revaluated
	require.EqualValues(t, found, 1, "Devices with both parent or no parent have to be returned")

	_, err = server.Delete(ctx, &registrypb.DeleteRequest{
		Id: resp.Device.Id,
	})

	//Delete the Account created
	_ = server.repo.DeleteAccount(ctx, &nodepb.DeleteAccountRequest{Uid: accid})
}

func TestListForAccount(t *testing.T) {

	ctx := context.Background()

	randomName := randomdata.SillyName()

	accid, err := server.repo.CreateUserAccount(ctx, randomName, "password", false, false, true)
	require.NoError(t, err)

	ctx = metadata.NewIncomingContext(ctx, metadata.New(map[string]string{"requestorid": accid}))

	ns, err := server.repo.GetNamespace(ctx, randomName)
	require.NoError(t, err)

	// Create the device
	request := &registrypb.CreateRequest{
		Device: sampleDevice(randomName, ns.Id),
	}
	resp, err := server.Create(ctx, request)
	require.NoError(t, err)

	response, err := server.List(ctx, &registrypb.ListDevicesRequest{
		Namespaceid: ns.Id,
		Account:     accid,
	})
	require.NoError(t, err)
	var found int
	for _, device := range response.Devices {
		if device.Name == randomName {
			found++
		}
	}

	//Assert needs to revaluated
	require.EqualValues(t, found, 1, "Devices with both parent or no parent have to be returned")

	_, err = server.Delete(ctx, &registrypb.DeleteRequest{
		Id: resp.Device.Id,
	})

	//Delete the Account created
	_ = server.repo.DeleteAccount(ctx, &nodepb.DeleteAccountRequest{Uid: accid})
}

func sampleDevice(name string, namespaceid string) *registrypb.Device {
	return &registrypb.Device{
		Namespace: namespaceid,
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
	ctx := context.Background()

	randomName := randomdata.SillyName()

	accid, err := server.repo.CreateUserAccount(ctx, randomName, "password", false, true, true)
	require.NoError(t, err)

	ctx = metadata.NewIncomingContext(ctx, metadata.New(map[string]string{"requestorid": accid}))

	ns, err := server.repo.GetNamespace(ctx, randomName)
	require.NoError(t, err)

	// Create
	request := &registrypb.CreateRequest{
		Device: sampleDevice(randomName, ns.Id),
	}
	response, err := server.Create(ctx, request)
	require.NoError(t, err)
	require.NotEmpty(t, response.Device.Certificate.Fingerprint)

	// Get
	respGet, err := server.Get(ctx, &registrypb.GetRequest{
		Id: response.Device.Id,
	})
	require.NoError(t, err)
	require.NotNil(t, respGet.Device)
	require.EqualValues(t, randomName, respGet.Device.Name)
	require.EqualValues(t, request.Device.Certificate.PemData, respGet.Device.Certificate.PemData)
	require.EqualValues(t, request.Device.Certificate.Algorithm, respGet.Device.Certificate.Algorithm)

	// Get by fingerprint
	respFP, err := server.GetByFingerprint(ctx, &registrypb.GetByFingerprintRequest{
		Fingerprint: response.Device.Certificate.Fingerprint,
	})
	require.NoError(t, err)
	require.Contains(t, respFP.Devices, &registrypb.Device{Id: respGet.Device.Id, Enabled: &wrappers.BoolValue{Value: true}, Name: respGet.Device.Name})

	_, err = server.Delete(ctx, &registrypb.DeleteRequest{
		Id: response.Device.Id,
	})

	//Delete the Account created
	_ = server.repo.DeleteAccount(ctx, &nodepb.DeleteAccountRequest{Uid: accid})
}

func TestUpdate(t *testing.T) {
	ctx := context.Background()

	randomName := randomdata.SillyName()

	accid, err := server.repo.CreateUserAccount(ctx, randomName, "password", false, true, true)
	require.NoError(t, err)

	ctx = metadata.NewIncomingContext(ctx, metadata.New(map[string]string{"requestorid": accid}))

	ns, err := server.repo.GetNamespace(ctx, randomName)
	require.NoError(t, err)

	// Create the device
	request := &registrypb.CreateRequest{
		Device: sampleDevice(randomName, ns.Id),
	}
	response, err := server.Create(ctx, request)
	require.NoError(t, err)
	require.NotEmpty(t, response.Device.Certificate.Fingerprint)

	// Get the device
	respGet, err := server.Get(ctx, &registrypb.GetRequest{
		Id: response.Device.Id,
	})

	//Validate the device
	require.NoError(t, err)
	require.NotNil(t, respGet.Device)
	require.EqualValues(t, randomName, respGet.Device.Name)
	require.EqualValues(t, request.Device.Certificate.PemData, respGet.Device.Certificate.PemData)
	require.EqualValues(t, request.Device.Certificate.Algorithm, respGet.Device.Certificate.Algorithm)

	// Get by fingerprint
	respFP, err := server.GetByFingerprint(ctx, &registrypb.GetByFingerprintRequest{
		Fingerprint: response.Device.Certificate.Fingerprint,
	})
	require.NoError(t, err)
	require.Contains(t, respFP.Devices, &registrypb.Device{Id: respGet.Device.Id, Enabled: &wrappers.BoolValue{Value: true}, Name: respGet.Device.Name})

	//Set new values
	NewName := randomdata.SillyName()

	var NewStatus *wrappers.BoolValue
	NewTag := []string{"d"}

	//Update the device
	_, err = server.Update(ctx, &registrypb.UpdateRequest{
		Device: &registrypb.Device{
			Id:      response.Device.Id,
			Name:    NewName,
			Enabled: NewStatus,
			Tags:    NewTag,
		},
		FieldMask: &field_mask.FieldMask{
			Paths: []string{"Name", "Tags"},
		},
	})
	require.NoError(t, err)

	// Get the updated device details
	respGet, err = server.Get(ctx, &registrypb.GetRequest{
		Id: response.Device.Id,
	})
	require.NoError(t, err)

	//Validate the updated device
	require.NoError(t, err)
	require.EqualValues(t, NewName, respGet.Device.Name)
	require.EqualValues(t, []string{"d", "c", "b", "a"}, respGet.Device.Tags)

	_, err = server.Delete(ctx, &registrypb.DeleteRequest{
		Id: response.Device.Id,
	})

	//Delete the Account created
	_ = server.repo.DeleteAccount(ctx, &nodepb.DeleteAccountRequest{Uid: accid})
}

func TestDelete(t *testing.T) {
	ctx := context.Background()

	randomName := randomdata.SillyName()

	//Create account for test
	accid, err := server.repo.CreateUserAccount(ctx, randomName, "password", false, false, true)
	require.NoError(t, err)

	//Set metadata for context
	ctx = metadata.NewIncomingContext(ctx, metadata.New(map[string]string{"requestorid": accid}))

	request := &registrypb.CreateRequest{
		Device: sampleDevice(randomdata.SillyName(), "0x1"),
	}
	response, err := server.Create(ctx, request)
	require.NoError(t, err)
	require.NotEmpty(t, response.Device.Certificate.Fingerprint)

	_, err = server.Delete(ctx, &registrypb.DeleteRequest{
		Id: response.Device.Id,
	})

	require.NoError(t, err)

	_, err = server.Get(ctx, &registrypb.GetRequest{
		Id: response.Device.Id,
	})
	require.Error(t, err)

	//Delete the Account created
	_ = server.repo.DeleteAccount(ctx, &nodepb.DeleteAccountRequest{Uid: accid})
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

	_, err1 := server.Create(context.Background(), request1)
	require.Error(t, err1)
	//require.Empty(t, response2.Device.Certificate.Fingerprint)

	_, err = server.Delete(context.Background(), &registrypb.DeleteRequest{
		Id: response.Device.Id,
	})
}
*/

//TODO test update/patch; also with cert

// TODO GetByFP: ensure that we dont give empty responses with 0 devices
