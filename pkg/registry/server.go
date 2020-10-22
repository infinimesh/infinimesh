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
	"encoding/pem"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/ptypes/wrappers"

	"github.com/infinimesh/infinimesh/pkg/node"
	"github.com/infinimesh/infinimesh/pkg/node/dgraph"
	"github.com/infinimesh/infinimesh/pkg/registry/registrypb"

	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/dgraph-io/dgo"
)

//Server is a Data type for Device Controller file
type Server struct {
	dgo *dgo.Dgraph

	repo node.Repo
}

//NewServer is a method to create the Dgraph Server for Device registry
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

//Create is a method for creating Devices
func (s *Server) Create(ctx context.Context, request *registrypb.CreateRequest) (*registrypb.CreateResponse, error) {

	_, err := s.repo.GetNamespaceID(ctx, request.Device.Namespace)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, "The Namespace provided is not found.")
	}

	if request.Device.Certificate == nil {
		return nil, status.Error(codes.FailedPrecondition, "No certificate provided.")
	}

	resp, err := s.CreateQ(ctx, request)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return resp, nil
}

//Update is a method for updating Devices details
func (s *Server) Update(ctx context.Context, request *registrypb.UpdateRequest) (response *registrypb.UpdateResponse, err error) {

	resp, err := s.UpdateQ(ctx, request)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return resp, nil
}

//GetByFingerprint is a method for get FringerPrint for a Device
func (s *Server) GetByFingerprint(ctx context.Context, request *registrypb.GetByFingerprintRequest) (*registrypb.GetByFingerprintResponse, error) {

	resp, err := s.GetByFingerprintQ(ctx, request)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return resp, nil
}

//Get is a method for get FringerPrint for a Device
func (s *Server) Get(ctx context.Context, request *registrypb.GetRequest) (response *registrypb.GetResponse, err error) {

	resp, err := s.GetQ(ctx, request)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return resp, nil
}

//List is a method that list all Devices for a specific Namespace
func (s *Server) List(ctx context.Context, request *registrypb.ListDevicesRequest) (response *registrypb.ListResponse, err error) {

	resp, err := s.ListQ(ctx, request)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return resp, nil
}

//ListForAccount is a method that list all Devices for a specififc Account
func (s *Server) ListForAccount(ctx context.Context, request *registrypb.ListDevicesRequest) (response *registrypb.ListResponse, err error) {

	resp, err := s.ListForAccountQ(ctx, request)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return resp, nil
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

//Delete is a method that deletes a Device
func (s *Server) Delete(ctx context.Context, request *registrypb.DeleteRequest) (response *registrypb.DeleteResponse, err error) {

	resp, err := s.DeleteQ(ctx, request)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return resp, nil
}
