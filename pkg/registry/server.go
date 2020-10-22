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
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/ptypes/wrappers"

	"github.com/infinimesh/infinimesh/pkg/node"
	"github.com/infinimesh/infinimesh/pkg/node/dgraph"
	"github.com/infinimesh/infinimesh/pkg/registry/registrypb"

	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"

	"encoding/base64"
)

//Server is a Data type for Device Controller file
type Server struct {
	dgo *dgo.Dgraph

	repo node.Repo
}

func NewServer(dg *dgo.Dgraph) *Server {
	return &Server{
		dgo: dg,
		repo: &dgraph.DGraphRepo{
			Dg: dg,
		},
	}
}

func (s *Server) getFingerprint(pemCert []byte, certType string) (fingerprint []byte, err error) {
	pemBlock, _ := pem.Decode(pemCert)
	if pemBlock == nil {
		return nil, errors.New("Could not decode PEM")
	}
	cert, err := x509.ParseCertificate(pemBlock.Bytes)
	if err != nil {
		return nil, err
	}

	return sha256Sum(cert.Raw), nil
}

func sha256Sum(c []byte) []byte {
	s := sha256.New()
	_, err := s.Write(c)
	if err != nil {
		panic(err)
	}
	return s.Sum(nil)
}

func (s *Server) Create(ctx context.Context, request *registrypb.CreateRequest) (*registrypb.CreateResponse, error) {
	txn := s.dgo.NewTxn()
	defer txn.Discard(ctx) // nolint
	if exists := dgraph.NameExists(ctx, txn, request.Device.Name, request.Device.Namespace, ""); exists {
		return nil, status.Error(codes.FailedPrecondition, "The device name exists already. Please provide a different name.")
	}

	ns, err := s.repo.GetNamespaceID(ctx, request.Device.Namespace)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, "The Namespace provided is not found.")
	}

	if request.Device.Certificate == nil {
		return nil, status.Error(codes.FailedPrecondition, "No certificate provided.")
	}

	fp, err := s.getFingerprint([]byte(request.Device.Certificate.PemData), request.Device.Certificate.Algorithm)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, "Invalid Certificate is provided.")
	}

	//To check if the fingerprint already exists, in this case creating new device is not permissible
	if exists := dgraph.FingerprintExists(ctx, txn, fp); exists {
		return nil, status.Error(codes.FailedPrecondition, "Certificate already exists. Please provide a different certificate.")
	}

	var enabled bool
	if request.Device.Enabled != nil {
		enabled = request.Device.Enabled.GetValue()
	}

	d := &Device{
		Object: dgraph.Object{
			Node: dgraph.Node{
				UID:  "_:new",
				Type: "object",
			},
			Name: request.Device.Name,
			Kind: node.KindDevice,
		},
		Enabled: enabled,
		Tags:    request.Device.Tags,
		Certificates: []*X509Cert{
			&X509Cert{
				PemData:              request.Device.Certificate.PemData,
				Algorithm:            request.Device.Certificate.Algorithm,
				Fingerprint:          fp,
				FingerprintAlgorithm: "sha256",
			},
		},
	}

	js, err := json.Marshal(d)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to create device")
	}

	mutRes, err := txn.Mutate(ctx, &api.Mutation{
		SetJson: js,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Failed to create object: %v", err))
	}

	newUID := mutRes.GetUids()["new"]

	nsMut := &api.NQuad{
		Subject:   ns.GetId(),
		Predicate: "owns",
		ObjectId:  newUID,
	}

	_, err = txn.Mutate(ctx, &api.Mutation{
		Set: []*api.NQuad{
			nsMut,
		},
		CommitNow: true,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Failed to create object: %v", err))
	}

	request.Device.Certificate.Fingerprint = fp
	request.Device.Certificate.FingerprintAlgorithm = "sha256"

	return &registrypb.CreateResponse{
		Device: &registrypb.Device{
			Id:          newUID,
			Name:        request.Device.Name,
			Enabled:     request.Device.Enabled,
			Tags:        request.Device.Tags,
			Namespace:   request.Device.Namespace,
			Certificate: request.Device.Certificate,
		},
	}, nil
}

// TODO tags can currently only be added due to the non-idempotent behavior of dgraph with list types
func (s *Server) Update(ctx context.Context, request *registrypb.UpdateRequest) (response *registrypb.UpdateResponse, err error) {
	txn := s.dgo.NewTxn()

	//Query to get the device details from the Dgraph DB
	const q = `query devices($id: string){
		device(func: uid($id)) @filter(eq(kind, "device")) {
		  uid
		  name
		  ~owns {
			uid
		  }
		  tags
		  enabled
		  certificates {
			uid
			pem_data
			algorithm
			fingerprint
			fingerprint.algorithm
		  }
		}
	  }`

	//Execute the Query to get device details
	resp, err := txn.QueryWithVars(ctx, q, map[string]string{
		"$id": request.Device.Id,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Failed to patch device: %v", err))
	}

	var result struct {
		Devices []Device `json:"device"`
	}

	err = json.Unmarshal(resp.Json, &result)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Failed to patch device: %v", err))
	}

	if len(result.Devices) != 1 {
		return nil, status.Error(codes.NotFound, "Device not found")
	}

	d := &Device{
		Object: dgraph.Object{
			Node: dgraph.Node{
				UID: result.Devices[0].UID,
			},
			Name: result.Devices[0].Name,
		},
		Enabled: result.Devices[0].Enabled,
		Tags:    result.Devices[0].Tags,
		Certificates: []*X509Cert{
			&X509Cert{
				Node: dgraph.Node{
					UID: result.Devices[0].Certificates[0].UID,
				},
				PemData:              result.Devices[0].Certificates[0].PemData,
				Algorithm:            result.Devices[0].Certificates[0].Algorithm,
				Fingerprint:          result.Devices[0].Certificates[0].Fingerprint,
				FingerprintAlgorithm: result.Devices[0].Certificates[0].FingerprintAlgorithm,
			},
		},
	}

	//Update the device details based on the data available.
	for _, field := range request.FieldMask.GetPaths() {
		switch strings.ToLower(field) {

		//Update the device details
		case "enabled":
			d.Enabled = request.Device.GetEnabled().Value
		case "tags":
			d.Tags = request.Device.Tags
		case "name":
			if exists := dgraph.NameExists(ctx, txn, request.Device.Name, request.Device.Namespace, ""); exists {
				return nil, status.Error(codes.FailedPrecondition, "The device name exists already. Please provide a different name.")
			}
			d.Name = request.Device.Name
		case "certificate":
			//Pre-check for updating certificates
			if request.Device.Certificate == nil {
				return nil, status.Error(codes.FailedPrecondition, "No certificate provided.")
			}

			fp, err := s.getFingerprint([]byte(request.Device.Certificate.PemData), request.Device.Certificate.Algorithm)
			if err != nil {
				return nil, status.Error(codes.FailedPrecondition, "Invalid Certificate is provided.")
			}

			//To check if the fingerprint already exists, in this case creating new device is not permissible
			if exists := dgraph.FingerprintExists(ctx, txn, fp); exists {
				return nil, status.Error(codes.FailedPrecondition, "Certificate already exists. Please provide a different certificate.")
			}

			//update the certificate
			d.Certificates[0].PemData = request.Device.Certificate.PemData
			d.Certificates[0].Algorithm = request.Device.Certificate.Algorithm
			d.Certificates[0].Fingerprint = request.Device.Certificate.Fingerprint
			d.Certificates[0].FingerprintAlgorithm = request.Device.Certificate.FingerprintAlgorithm
		}
	}

	js, err := json.Marshal(&d)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Failed to patch device: %v", err))
	}

	_, err = txn.Mutate(ctx, &api.Mutation{
		SetJson:   js,
		CommitNow: true,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Failed to patch device: %v", err))
	}

	return &registrypb.UpdateResponse{}, nil
}

func (s *Server) GetByFingerprint(ctx context.Context, request *registrypb.GetByFingerprintRequest) (*registrypb.GetByFingerprintResponse, error) {
	txn := s.dgo.NewReadOnlyTxn()

	const q = `query devices($fingerprint: string){
  devices(func: eq(fingerprint, $fingerprint)) @normalize {
    ~certificates {
      uid : uid
      name : name
      enabled : enabled
      ~owns {
        namespace: name
      }
    }
  }
}
  `

	vars := map[string]string{
		"$fingerprint": base64.StdEncoding.EncodeToString(request.Fingerprint),
	}

	resp, err := txn.QueryWithVars(ctx, q, vars)
	if err != nil {
		return nil, err
	}

	var res struct {
		Devices []struct {
			Uid       string `json:"uid"`
			Name      string `json:"name"`
			Enabled   bool   `json:"enabled"`
			Namespace string `json:"namespace"`
		} `json:"devices"`
	}

	err = json.Unmarshal(resp.Json, &res)
	if err != nil {
		return nil, err
	}

	var devices []*registrypb.Device
	for _, device := range res.Devices {
		devices = append(devices, &registrypb.Device{
			Id:        device.Uid,
			Name:      device.Name,
			Namespace: device.Namespace,
			Enabled:   &wrappers.BoolValue{Value: device.Enabled},
		})
	}

	return &registrypb.GetByFingerprintResponse{
		Devices: devices,
	}, nil
}

func (s *Server) Get(ctx context.Context, request *registrypb.GetRequest) (response *registrypb.GetResponse, err error) {
	txn := s.dgo.NewReadOnlyTxn()

	const q = `query devices($id: string){
  device(func: uid($id)) @filter(eq(kind, "device")) {
    uid
    name
    tags
    enabled
    certificates {
      pem_data
      algorithm
      fingerprint
      fingerprint.algorithm
    }
    ~owns {
      name
    }
  }
}`

	vars := map[string]string{
		"$id": request.Id,
	}

	resp, err := txn.QueryWithVars(ctx, q, vars)
	if err != nil {
		return nil, err
	}

	var res struct {
		Devices []*Device `json:"device"`
	}

	err = json.Unmarshal(resp.Json, &res)
	if err != nil {
		return nil, err
	}

	if len(res.Devices) == 0 {
		return &registrypb.GetResponse{}, status.Error(codes.NotFound, "Device not found")
	}

	return &registrypb.GetResponse{
		Device: toProto(res.Devices[0]),
	}, nil
}

//List is a method that list details of the all devices
func (s *Server) List(ctx context.Context, request *registrypb.ListDevicesRequest) (response *registrypb.ListResponse, err error) {
	return s.List(ctx, request)
}

func (s *Server) ListForAccount(ctx context.Context, request *registrypb.ListDevicesRequest) (response *registrypb.ListResponse, err error) {
	txn := s.dgo.NewReadOnlyTxn()

	// TODO direct access!

	var q = `query list($account: string, $namespaceid: string){
		var(func: uid($account)) {
		  access.to.namespace %v {
			owns {
			  OBJs as uid
			} @filter(eq(kind, "device"))
		  }
		}

		nodes(func: uid(OBJs)) @recurse {
		  children{} 
		  uid
		  name
		  kind
		  enabled
		  tags
		  ~owns {
			name
		  }
		}
	  }`

	if request.Namespace != "" {
		q = fmt.Sprintf(q, "@filter(uid($namespaceid))")
	} else {
		q = fmt.Sprintf(q, "")
	}

	vars := map[string]string{
		"$account":     request.Account,
		"$namespaceid": request.Namespace,
	}

	resp, err := txn.QueryWithVars(ctx, q, vars)
	if err != nil {
		return nil, err
	}

	var res struct {
		Nodes []Device `json:"nodes"`
	}

	err = json.Unmarshal(resp.Json, &res)
	if err != nil {
		return nil, err
	}

	var devices []*registrypb.Device
	for _, device := range res.Nodes {
		devices = append(devices, toProto(&device))
	}

	return &registrypb.ListResponse{
		Devices: devices,
	}, nil
}

// TODO make registrypb.Device have multiple certs, ... we also ignore valid_to for now
func toProto(device *Device) *registrypb.Device {
	res := &registrypb.Device{
		Id:      device.UID,
		Name:    device.Name,
		Enabled: &wrappers.BoolValue{Value: device.Enabled},
		Tags:    device.Tags,
		// TODO cert etc
	}

	if len(device.OwnedBy) == 1 {
		res.Namespace = device.OwnedBy[0].Name
	}

	if len(device.Certificates) > 0 {
		res.Certificate = &registrypb.Certificate{
			PemData:              device.Certificates[0].PemData,
			Algorithm:            device.Certificates[0].Algorithm,
			FingerprintAlgorithm: device.Certificates[0].FingerprintAlgorithm,
			Fingerprint:          device.Certificates[0].Fingerprint,
		}
	}
	return res
}

func (s *Server) Delete(ctx context.Context, request *registrypb.DeleteRequest) (response *registrypb.DeleteResponse, err error) {
	txn := s.dgo.NewTxn()
	m := &api.Mutation{CommitNow: true}

	//Query to get the device to be deleted with all the related edges
	const q = `query delete($device: string){
		objects(func: uid($device)) @filter(eq(kind, "device")) {
			uid
		  ~owns {
			uid
		  }
		  ~children {
			uid
		  }
		 certificates {
			uid
        type
		  }
		}
	  }`

	res, err := txn.QueryWithVars(ctx, q, map[string]string{"$device": request.Id})
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to delete node "+err.Error())
	}

	var result struct {
		//Get the Device edge details from the query response and build JSON
		Objects []*Device `json:"objects"`
	}

	err = json.Unmarshal(res.Json, &result)
	if err != nil {
		return nil, err
	}

	if len(result.Objects) != 1 {
		return nil, status.Error(codes.NotFound, "The Device is not found")
	}

	//Append edge if there is a owns edge
	if len(result.Objects[0].OwnedBy) == 1 {
		m.Del = append(m.Del, &api.NQuad{
			Subject:   result.Objects[0].OwnedBy[0].UID,
			Predicate: "owns",
			ObjectId:  request.Id,
		})
	}

	//Append edge if there is a children edge
	if len(result.Objects[0].Parent) == 1 {
		m.Del = append(m.Del, &api.NQuad{
			Subject:   result.Objects[0].Parent[0].UID,
			Predicate: "children",
			ObjectId:  request.Id,
		})
	}

	//Delete all the edges appended in mutation m
	dgo.DeleteEdges(m, request.Id, "_STAR_ALL")

	//Append node if there is a certificate edge to delete the certificate node
	if len(result.Objects[0].Certificates) == 1 {
		m.Del = append(m.Del, &api.NQuad{
			Subject:     result.Objects[0].Certificates[0].UID,
			Predicate:   "_STAR_ALL",
			ObjectId:    "_STAR_ALL",
			ObjectValue: &api.Value{Val: &api.Value_DefaultVal{DefaultVal: "_STAR_ALL"}},
		})
	}

	_, err = txn.Mutate(context.Background(), m)
	if err != nil {
		return nil, err
	}

	return &registrypb.DeleteResponse{}, nil
}
