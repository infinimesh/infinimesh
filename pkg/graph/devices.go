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
	"github.com/golang-jwt/jwt"
	"github.com/infinimesh/infinimesh/pkg/graph/schema"
	pb "github.com/infinimesh/infinimesh/pkg/node/proto"
	devpb "github.com/infinimesh/infinimesh/pkg/node/proto/devices"
	inf "github.com/infinimesh/infinimesh/pkg/shared"
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

func (o *Device) SetAccessLevel(level schema.InfinimeshAccessLevel) {
	il := int32(level)
	o.AccessLevel = &il
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

func NewDevicesController(log *zap.Logger, db driver.Database) *DevicesController {
	ctx := context.TODO()
	col, _ := db.Collection(ctx, schema.DEVICES_COL)

	return &DevicesController{
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
	
	requestor := ctx.Value(inf.InfinimeshAccountCtxKey).(string)
	log.Debug("Requestor", zap.String("id", requestor))

	ns_id := req.GetNamespace()
	if ns_id == "" {
		return nil, status.Error(codes.InvalidArgument, "Namespace ID is required")
	}

	ns := NewBlankNamespaceDocument(ns_id)

	ok, level := AccessLevel(ctx, c.db, NewBlankAccountDocument(requestor), ns)
	if !ok || level < int32(schema.ADMIN) {
		return nil, status.Errorf(codes.PermissionDenied, "No Access to Namespace %s", ns_id)
	}

	device := Device{Device: req.GetDevice()}
	err := sha256Fingerprint(device.Device.Certificate)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Can't generate fingerprint: %v", err)
	}
	device.Token = ""

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
	log := c.log.Named("Get")
	log.Debug("Get request received", zap.Any("request", dev), zap.Any("context", ctx))

	requestor := ctx.Value(inf.InfinimeshAccountCtxKey).(string)
	log.Debug("Requestor", zap.String("id", requestor))

	// Getting Account from DB
	// and Check requestor access
	device := *NewBlankDeviceDocument(dev.GetUuid())
	err := AccessLevelAndGet(ctx, log, c.db, NewBlankAccountDocument(requestor), &device)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Account not found or not enough Access Rights")
	}
	if *device.AccessLevel < 1 {
		return nil, status.Error(codes.PermissionDenied, "Not enough Access Rights")
	}

	post := false
	if *device.AccessLevel > 1 {
		post = true
	}
	token, err := c._MakeToken([]string{device.Uuid}, post, 0)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to issue token")
	}
	device.Token = token

	return device.Device, nil
}

func (c *DevicesController) GetByToken(ctx context.Context, dev *devpb.Device) (*devpb.Device, error) {
	log := c.log.Named("GetByToken")
	log.Debug("Get by Token request received", zap.String("device", dev.Uuid), zap.Any("context", ctx))

	devices_scope := ctx.Value(inf.InfinimeshDevicesCtxKey).([]string)
	log.Debug("Devices Scope", zap.Any("devices", devices_scope))

	found := false
	for _, device := range devices_scope {
		if device == dev.GetUuid() {
			found = true
			break
		}
	}
	if !found {
		return nil, status.Error(codes.Unauthenticated, "Requested device is outside of token scope")
	}

	var device devpb.Device
	meta, err := c.col.ReadDocument(ctx, dev.GetUuid(), &device)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Device not found")
	}
	device.Uuid = meta.ID.Key()

	if !ctx.Value(inf.InfinimeshPostAllowedCtxKey).(bool) {
		device.Certificate = nil
	}

	return &device, nil
}

func (c *DevicesController) List(ctx context.Context, _ *pb.EmptyMessage) (*devpb.DevicesPool, error) {
	log := c.log.Named("List")

	requestor := ctx.Value(inf.InfinimeshAccountCtxKey).(string)
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

	requestor := ctx.Value(inf.InfinimeshAccountCtxKey).(string)
	log.Debug("Requestor", zap.String("id", requestor))

	acc := *NewBlankAccountDocument(requestor)
	dev := *NewBlankDeviceDocument(req.GetUuid())
	err := AccessLevelAndGet(ctx, log, c.db, &acc, &dev)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Account not found or not enough Access Rights")
	}
	if *dev.AccessLevel < 3 {
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

	requestor := ctx.Value(inf.InfinimeshAccountCtxKey).(string)
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

	token, err := c._MakeToken([]string{r.Uuid}, false, 0)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to issue token")
	}
	r.Token = token

	return &r, nil
}

func (c *DevicesController) MakeDevicesToken(ctx context.Context, req *pb.DevicesTokenRequest) (*pb.TokenResponse, error) {
	log := c.log.Named("MakeDevicesToken")
	log.Debug("MakeDevicesToken request received", zap.Any("request", req), zap.Any("context", ctx))

	requestor := ctx.Value(inf.InfinimeshAccountCtxKey).(string)
	log.Debug("Requestor", zap.String("id", requestor))

	acc := *NewBlankAccountDocument(requestor)
	access := int32(schema.READ)
	if req.GetPost() {
		access = int32(schema.MGMT)
	}

	for _, uuid := range req.GetDevices() {
		ok, level := AccessLevel(ctx, c.db, &acc, NewBlankDeviceDocument(uuid))
		if !ok {
			return nil, status.Errorf(codes.NotFound, "Account not found or not enough Access Rights to device: %s", uuid)
		}
		if level < access {
			return nil, status.Errorf(codes.PermissionDenied, "Not enough Access Rights to device: %s", uuid)
		}
	}

	token_string, err := c._MakeToken(req.GetDevices(), req.GetPost(), req.GetExp())
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to issue token")
	}

	return &pb.TokenResponse{Token: token_string}, nil
}

func (c *DevicesController) _MakeToken(devices []string, post bool, exp int32) (string, error) {
	claims := jwt.MapClaims{}
	claims[inf.INFINIMESH_DEVICES_CLAIM] = devices
	claims[inf.INFINIMESH_POST_STATE_ALLOWED_CLAIM] = post
	claims["exp"] = exp

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(c.SIGNING_KEY)
}