/*
Copyright © 2021-2022 Infinite Devices GmbH, Nikita Ivanovski info@slnt-opp.xyz

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
package graph

import (
	"context"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"

	"github.com/arangodb/go-driver"
	"github.com/infinimesh/infinimesh/pkg/graph/schema"
	pb "github.com/infinimesh/infinimesh/pkg/node/proto"
	devpb "github.com/infinimesh/infinimesh/pkg/node/proto/devices"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Device struct {
	*devpb.Device
	driver.DocumentMeta	
}

func (o *Device) ID() (driver.DocumentID) {
	return o.DocumentMeta.ID
}

func NewBlankDeviceDocument(key string) (*Device) {
	return &Device{
		Device: &devpb.Device{
			Uuid: key,
		},
		DocumentMeta: NewBlankDocument(schema.DEVICES_COL, key),
	}
}


type DevicesController struct {
	pb.UnimplementedDevicesServiceServer
	log *zap.Logger

	col driver.Collection // Devices Collection
	db driver.Database

	ns2dev driver.Collection // Namespaces to Devices permissions edge collection

	SIGNING_KEY []byte
}

func NewDevicesController(log *zap.Logger, db driver.Database) DevicesController {
	ctx := context.TODO()
	col, _ := db.Collection(ctx, schema.DEVICES_COL)

	return DevicesController{
		log: log.Named("DevicesController"), col: col, db: db,
		ns2dev: GetEdgeCol(ctx, db, schema.NS2DEV),
		SIGNING_KEY: []byte("just-an-init-thing-replace-me"),
	}
}

func sha256Fingerprint(cert *devpb.Certificate) (err error) {
	block, _ := pem.Decode([]byte(cert.PemData))
	if block == nil {
		return errors.New("coudn't decode PEM data")
	}

	parsed, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return err
	}

	s := sha256.New()
	_, err = s.Write(parsed.Raw)
	if err != nil {
		return err
	}

	cert.Algorithm = "sha256"
	cert.Fingerprint = s.Sum(nil)

	return nil
}

func (c *DevicesController) Create(ctx context.Context, req *devpb.CreateRequest) (*devpb.CreateResponse, error) {
	log := c.log.Named("Create")
	log.Debug("Create request received", zap.Any("request", req), zap.Any("context", ctx))
	
	//Get metadata from context and perform validation
	_, requestor, err := Validate(ctx, log)
	if err != nil {
		return nil, err
	}
	log.Debug("Requestor", zap.String("id", requestor))

	ns_id := req.GetNamespace()
	if ns_id == "" {
		ns_id = schema.ROOT_NAMESPACE_KEY
	}

	ns := NewBlankNamespaceDocument(ns_id)

	ok, level := AccessLevel(ctx, c.db, NewBlankAccountDocument(requestor), ns)
	if !ok || level < int32(schema.ADMIN) {
		return nil, status.Errorf(codes.PermissionDenied, "No Access to Namespace %s", ns_id)
	}

	device := Device{Device: req.GetDevice()}
	err = sha256Fingerprint(device.Certificate)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Can't generate fingerprint: %v", err)
	}

	meta, err := c.col.CreateDocument(ctx, device)
	if err != nil {
		log.Error("Error creating Device", zap.Error(err))
		return nil, status.Error(codes.Internal, "Error while creating Device")
	}
	device.Uuid = meta.ID.Key()
	device.DocumentMeta = meta

	err = Link(ctx, log, c.ns2dev, ns, &device, schema.ADMIN)
	if err != nil {
		log.Error("Error creating edge", zap.Error(err))
		c.col.RemoveDocument(ctx, device.Uuid)
		return nil, status.Error(codes.Internal, "error creating Permission")
	}

	return &devpb.CreateResponse{
		Device: device.Device,
	}, nil
}

func (c *DevicesController) Get(ctx context.Context, dev *devpb.Device) (*devpb.Device, error) {
	log := c.log.Named("Create")
	log.Debug("Get request received", zap.Any("request", dev), zap.Any("context", ctx))

	//Get metadata from context and perform validation
	_, requestor, err := Validate(ctx, log)
	if err != nil {
		return nil, err
	}
	log.Debug("Requestor", zap.String("id", requestor))

	// Getting Account from DB
	// and Check requestor access
	device := *NewBlankDeviceDocument(dev.GetUuid())
	ok, level := AccessLevelAndGet(ctx, log, c.db, NewBlankAccountDocument(requestor), &device)
	if !ok {
		return nil, status.Error(codes.NotFound, "Account not found or not enough Access Rights")
	}
	if level < 1 {
		return nil, status.Error(codes.PermissionDenied, "Not enough Access Rights")
	}

	return device.Device, nil
}

func (c *DevicesController) List(ctx context.Context, _ *pb.EmptyMessage) (*devpb.DevicesPool, error) {
	log := c.log.Named("List")

	//Get metadata from context and perform validation
	_, requestor, err := Validate(ctx, log)
	if err != nil {
		return nil, err
	}
	log.Debug("Requestor", zap.String("id", requestor))

	cr, err := ListQuery(ctx, log, c.db, NewBlankAccountDocument(requestor), schema.DEVICES_COL, 4)
	if err != nil {
		log.Error("Error executing query", zap.Error(err))
		return nil, status.Error(codes.Internal, "Couldn't execute query")
	}
	defer cr.Close()

	var r []*devpb.Device
	for {
		var dev devpb.Device
		meta, err := cr.ReadDocument(ctx, &dev)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			log.Error("Error unmarshalling Document", zap.Error(err))
			return nil, status.Error(codes.Internal, "Couldn't execute query")
		}
		dev.Uuid = meta.ID.Key()
		log.Debug("Got document", zap.Any("device", &dev))
		r = append(r, &dev)
	}

	return &devpb.DevicesPool{
		Devices: r,
	}, nil
}

func (c *DevicesController) Delete(ctx context.Context, req *devpb.Device) (*pb.DeleteResponse, error) {
	log := c.log.Named("Delete")

	//Get metadata from context and perform validation
	_, requestor, err := Validate(ctx, log)
	if err != nil {
		return nil, err
	}
	log.Debug("Requestor", zap.String("id", requestor))

	acc := *NewBlankAccountDocument(requestor)
	dev := *NewBlankDeviceDocument(req.GetUuid())
	ok, level := AccessLevelAndGet(ctx, log, c.db, &acc, &dev)
	if !ok {
		return nil, status.Error(codes.NotFound, "Account not found or not enough Access Rights")
	}
	if level < 3 {
		return nil, status.Error(codes.PermissionDenied, "Not enough Access Rights")
	}

	_, err = c.col.RemoveDocument(ctx, dev.ID().Key())
	if err != nil {
		log.Error("Error removing document", zap.Error(err))
		return nil, status.Error(codes.Internal, "Error deleting Device")
	}

	return &pb.DeleteResponse{}, nil
}

const findByFingerprintQuery = 
`FOR device IN @@devices
FILTER device.certificate.fingerprint == @fingerprint
RETURN device`

func (c *DevicesController) GetByFingerprint(ctx context.Context, req *devpb.GetByFingerprintRequest) (*devpb.Device, error) {
	log := c.log.Named("GetByFingerprint")
	log.Debug("GetByFingerprint request received", zap.Any("request", req), zap.Any("context", ctx))

	//Get metadata from context and perform validation
	_, requestor, err := Validate(ctx, log)
	if err != nil {
		return nil, err
	}
	log.Debug("Requestor", zap.String("id", requestor))

	cr, err := c.db.Query(ctx, findByFingerprintQuery, map[string]interface{}{
		"@devices": schema.DEVICES_COL,
		"fingerprint": req.GetFingerprint(),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "Error executing query")
	}
	defer cr.Close()

	var r devpb.Device
	meta, err := cr.ReadDocument(ctx, &r)
	if driver.IsNoMoreDocuments(err) {
		return nil, status.Error(codes.NotFound, "Device not found")
	}
	if err != nil {
		log.Error("Error unmarshalling Document", zap.Error(err))
		return nil, status.Error(codes.Internal, "Error executing query")
	}
	r.Uuid = meta.ID.Key()

	return &r, nil
}