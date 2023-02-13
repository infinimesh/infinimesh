/*
Copyright © 2021-2023 Infinite Devices GmbH, Nikita Ivanovski info@slnt-opp.xyz

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
	"fmt"

	"github.com/arangodb/go-driver"
	"github.com/golang-jwt/jwt"
	"github.com/infinimesh/infinimesh/pkg/graph/schema"
	inf "github.com/infinimesh/infinimesh/pkg/shared"
	"github.com/infinimesh/proto/handsfree"
	pb "github.com/infinimesh/proto/node"
	access "github.com/infinimesh/proto/node/access"
	devpb "github.com/infinimesh/proto/node/devices"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Device struct {
	*devpb.Device
	driver.DocumentMeta
}

func (o *Device) ID() driver.DocumentID {
	return o.DocumentMeta.ID
}

func (o *Device) SetAccessLevel(level access.Level) {
	if o.Access == nil {
		o.Access = &access.Access{
			Level: level,
		}
		return
	}
	o.Access.Level = level
}

func (o *Device) GetAccess() *access.Access {
	if o.Access == nil {
		return &access.Access{
			Level: access.Level_NONE,
		}
	}
	return o.Access
}

func NewBlankDeviceDocument(key string) *Device {
	return &Device{
		Device: &devpb.Device{
			Uuid: key,
		},
		DocumentMeta: NewBlankDocument(schema.DEVICES_COL, key),
	}
}

func NewDeviceFromPB(dev *devpb.Device) (res *Device) {
	return &Device{
		Device:       dev,
		DocumentMeta: NewBlankDocument(schema.DEVICES_COL, dev.Uuid),
	}
}

type DevicesController struct {
	pb.UnimplementedDevicesServiceServer
	InfinimeshBaseController

	col driver.Collection // Devices Collection
	hfc handsfree.HandsfreeServiceClient

	ns2dev  driver.Collection // Namespaces to Devices permissions edge collection
	acc2dev driver.Collection // Accounts to Devices permissions edge collection

	SIGNING_KEY []byte
}

func NewDevicesController(log *zap.Logger, db driver.Database, hfc handsfree.HandsfreeServiceClient) *DevicesController {
	ctx := context.TODO()
	col, _ := db.Collection(ctx, schema.DEVICES_COL)

	return &DevicesController{
		InfinimeshBaseController: InfinimeshBaseController{
			log: log.Named("DevicesController"), db: db,
		},
		col: col, hfc: hfc,
		ns2dev:      GetEdgeCol(ctx, db, schema.NS2DEV),
		acc2dev:     GetEdgeCol(ctx, db, schema.ACC2DEV),
		SIGNING_KEY: []byte("just-an-init-thing-replace-me"),
	}
}

func sha256Fingerprint(cert *devpb.Certificate) (err error) {
	if cert == nil {
		return errors.New("certificate is nil")
	}
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

	if req.Handsfree != nil {
		return c._HandsfreeCreate(ctx, req)
	}

	requestor := ctx.Value(inf.InfinimeshAccountCtxKey).(string)
	log.Debug("Requestor", zap.String("id", requestor))

	ns_id := req.GetNamespace()
	if ns_id == "" {
		return nil, status.Error(codes.InvalidArgument, "Namespace ID is required")
	}

	ns := NewBlankNamespaceDocument(ns_id)

	ok, level := AccessLevel(ctx, c.db, NewBlankAccountDocument(requestor), ns)
	if !ok || level < access.Level_ADMIN {
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
		log.Warn("Error creating Device", zap.Error(err))
		return nil, status.Error(codes.Internal, "Error while creating Device")
	}
	device.Uuid = meta.ID.Key()
	device.DocumentMeta = meta

	err = Link(ctx, log, c.ns2dev, ns, &device, access.Level_ADMIN, access.Role_OWNER)
	if err != nil {
		log.Warn("Error creating edge", zap.Error(err))
		c.col.RemoveDocument(ctx, device.Uuid)
		return nil, status.Error(codes.Internal, "error creating Permission")
	}

	return &devpb.CreateResponse{
		Device: device.Device,
	}, nil
}

func (c *DevicesController) _HandsfreeCreate(ctx context.Context, req *devpb.CreateRequest) (*devpb.CreateResponse, error) {
	log := c.log.Named("HandsfreeCreate")
	log.Debug("Request received")

	requestor := ctx.Value(inf.InfinimeshAccountCtxKey).(string)
	log.Debug("Requestor", zap.String("id", requestor))

	ns_id := req.GetNamespace()
	if ns_id == "" {
		return nil, status.Error(codes.InvalidArgument, "Namespace ID is required")
	}

	ns := NewBlankNamespaceDocument(ns_id)

	ok, level := AccessLevel(ctx, c.db, NewBlankAccountDocument(requestor), ns)
	if !ok || level < access.Level_ADMIN {
		return nil, status.Errorf(codes.PermissionDenied, "No Access to Namespace %s", ns_id)
	}

	dev := req.GetDevice()
	device := Device{Device: dev}
	device.Token = ""

	meta, err := c.col.CreateDocument(ctx, device)
	if err != nil {
		log.Warn("Error creating Device", zap.Error(err))
		return nil, status.Error(codes.Internal, "Error while creating Device")
	}
	device.Uuid = meta.ID.Key()
	device.DocumentMeta = meta

	err = Link(ctx, log, c.ns2dev, ns, &device, access.Level_ADMIN, access.Role_OWNER)
	if err != nil {
		log.Warn("Error creating edge", zap.Error(err))
		c.col.RemoveDocument(ctx, device.Uuid)
		return nil, status.Error(codes.Internal, "error creating Permission")
	}

	cleanup := func(err error) (*devpb.CreateResponse, error) {
		if _, d_err := c.Delete(ctx, device.Device); d_err != nil {
			log.Warn("Couldn't delete Device", zap.Error(d_err))
			return &devpb.CreateResponse{
				Device: device.Device,
			}, status.Error(codes.OK, "Couldn't delete freshly created device as well as set the certificate")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	cp, err := c.hfc.Send(ctx, &handsfree.ControlPacket{
		Payload: append([]string{req.GetHandsfree().GetCode(), device.Uuid}, req.GetHandsfree().GetPayload()...),
	})
	if err != nil {
		log.Warn("Couldn't obtain certificate from Handsfree", zap.Error(err))
		return cleanup(err)
	}

	if len(cp.GetPayload()) == 0 {
		log.Warn("Handsfree connection Payload is empty")
		return cleanup(fmt.Errorf("issue with Handsfree Payload: is empty"))
	}

	dev.Certificate = &devpb.Certificate{
		PemData: cp.GetPayload()[0],
	}
	dev.Tags = append(dev.Tags, cp.GetAppId())

	err = sha256Fingerprint(dev.Certificate)
	if err != nil {
		log.Warn("Couldn't generate certificate Hash", zap.Error(err))
		return cleanup(err)
	}

	_, err = c.col.ReplaceDocument(ctx, device.Uuid, dev)
	if err != nil {
		log.Warn("Couldn't set Device Certificate", zap.Error(err))
		return cleanup(err)
	}

	return &devpb.CreateResponse{
		Device: dev,
	}, nil
}

func (c *DevicesController) Update(ctx context.Context, dev *devpb.Device) (*devpb.Device, error) {
	log := c.log.Named("Update")
	log.Debug("Update request received", zap.Any("device", dev), zap.Any("context", ctx))

	curr, err := c.Get(ctx, dev)
	if err != nil {
		return nil, err
	}

	if curr.GetAccess().GetLevel() < access.Level_MGMT {
		return nil, status.Errorf(codes.PermissionDenied, "No Access to Device %s", dev.Uuid)
	}

	if dev.Title == "" {
		return nil, status.Error(codes.InvalidArgument, "Device Title cannot be empty")
	}

	curr.Tags = dev.Tags
	curr.Title = dev.Title

	_, err = c.col.ReplaceDocument(ctx, dev.Uuid, curr)
	if err != nil {
		log.Warn("Error updating Device", zap.Error(err))
		return nil, status.Error(codes.Internal, "Error while updating Device")
	}

	return curr, nil
}

func (c *DevicesController) Toggle(ctx context.Context, dev *devpb.Device) (*devpb.Device, error) {
	log := c.log.Named("Update")
	log.Debug("Update request received", zap.Any("device", dev), zap.Any("context", ctx))

	curr, err := c.Get(ctx, dev)
	if err != nil {
		return nil, err
	}

	if curr.GetAccess().GetLevel() < access.Level_MGMT {
		return nil, status.Errorf(codes.PermissionDenied, "No Access to Device %s", dev.Uuid)
	}

	res := NewDeviceFromPB(curr)
	err = Toggle(ctx, c.db, res, "enabled")
	if err != nil {
		log.Warn("Error updating Device", zap.Error(err))
		return nil, status.Error(codes.Internal, "Error while updating Device")
	}

	return curr, nil
}

func (c *DevicesController) ToggleBasic(ctx context.Context, dev *devpb.Device) (*devpb.Device, error) {
	log := c.log.Named("Update")
	log.Debug("Update request received", zap.Any("device", dev), zap.Any("context", ctx))

	curr, err := c.Get(ctx, dev)
	if err != nil {
		return nil, err
	}

	if curr.GetAccess().GetLevel() < access.Level_MGMT {
		return nil, status.Errorf(codes.PermissionDenied, "No Access to Device %s", dev.Uuid)
	}

	res := NewDeviceFromPB(curr)
	err = Toggle(ctx, c.db, res, "basic_enabled")
	if err != nil {
		log.Warn("Error updating Device", zap.Error(err))
		return nil, status.Error(codes.Internal, "Error while updating Device")
	}

	return curr, nil
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
	if device.Access.Level < access.Level_READ {
		return nil, status.Error(codes.PermissionDenied, "Not enough Access Rights")
	}

	post := false
	if device.Access.Level > access.Level_READ {
		post = true
	} else {
		device.Certificate = nil
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

func (c *DevicesController) List(ctx context.Context, _ *pb.EmptyMessage) (*devpb.Devices, error) {
	log := c.log.Named("List")

	requestor := ctx.Value(inf.InfinimeshAccountCtxKey).(string)
	log.Debug("Requestor", zap.String("id", requestor))

	cr, err := ListQuery(ctx, log, c.db, NewBlankAccountDocument(requestor), schema.DEVICES_COL, 10)
	if err != nil {
		log.Warn("Error executing query", zap.Error(err))
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
			log.Warn("Error unmarshalling Document", zap.Error(err))
			return nil, status.Error(codes.Internal, "Couldn't execute query")
		}
		dev.Uuid = meta.ID.Key()
		if dev.Access.Level < access.Level_MGMT {
			dev.Certificate = nil
		}

		log.Debug("Got document", zap.Any("device", &dev))
		r = append(r, &dev)
	}

	return &devpb.Devices{
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
	if dev.Access.Level < access.Level_ADMIN {
		return nil, status.Error(codes.PermissionDenied, "Not enough Access Rights")
	}

	_, err = c.col.RemoveDocument(ctx, dev.ID().Key())
	if err != nil {
		log.Warn("Error removing document", zap.Error(err))
		return nil, status.Error(codes.Internal, "Error deleting Device")
	}

	err = Link(
		ctx, log, c.ns2dev,
		NewBlankNamespaceDocument(*dev.Access.Namespace),
		&dev, access.Level_NONE, access.Role_UNSET,
	)
	if err != nil {
		log.Warn("Error removing device from namespace", zap.Error(err))
	}

	return &pb.DeleteResponse{}, nil
}

const findByFingerprintQuery = `FOR device IN @@devices
FILTER device.certificate.fingerprint == @fingerprint
RETURN device`

func (c *DevicesController) GetByFingerprint(ctx context.Context, req *devpb.GetByFingerprintRequest) (*devpb.Device, error) {
	log := c.log.Named("GetByFingerprint")
	log.Debug("GetByFingerprint request received", zap.Any("request", req), zap.Any("context", ctx))

	requestor := ctx.Value(inf.InfinimeshAccountCtxKey).(string)
	log.Debug("Requestor", zap.String("id", requestor))

	cr, err := c.db.Query(ctx, findByFingerprintQuery, map[string]interface{}{
		"@devices":    schema.DEVICES_COL,
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
		log.Warn("Error unmarshalling Document", zap.Error(err))
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
	_access := access.Level_READ
	if req.GetPost() {
		_access = access.Level_MGMT
	}

	for _, uuid := range req.GetDevices() {
		ok, level := AccessLevel(ctx, c.db, &acc, NewBlankDeviceDocument(uuid))
		if !ok {
			return nil, status.Errorf(codes.NotFound, "Account not found or not enough Access Rights to device: %s", uuid)
		}
		if level < _access {
			return nil, status.Errorf(codes.PermissionDenied, "Not enough Access Rights to device: %s", uuid)
		}
	}

	token_string, err := c._MakeToken(req.GetDevices(), req.GetPost(), req.GetExp())
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to issue token")
	}

	return &pb.TokenResponse{Token: token_string}, nil
}

func (c *DevicesController) _MakeToken(devices []string, post bool, exp int64) (string, error) {
	claims := jwt.MapClaims{}
	claims[inf.INFINIMESH_DEVICES_CLAIM] = devices
	claims[inf.INFINIMESH_POST_STATE_ALLOWED_CLAIM] = post
	claims["exp"] = exp

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(c.SIGNING_KEY)
}

const listDeviceJoinsQuery = `
FOR node, edge, path IN 1 INBOUND @device
GRAPH Permissions
FILTER edge.level > 0 && edge.role != 1
RETURN {
    node: node._id,
    access: KEEP(edge, ["level", "role"])
}
`

func (c *DevicesController) Joins(ctx context.Context, req *devpb.Device) (*access.Nodes, error) {
	log := c.log.Named("Joins")
	log.Debug("Fetch Joins request received", zap.String("device", req.GetUuid()))

	requestor := ctx.Value(inf.InfinimeshAccountCtxKey).(string)
	log.Debug("Requestor", zap.String("id", requestor))

	dev, err := c.Get(ctx, req)
	if err != nil {
		return nil, err
	}
	if dev.Access == nil || dev.Access.Level < access.Level_ADMIN {
		return nil, status.Error(codes.PermissionDenied, "Must be device Admin to fetch Joins")
	}

	cr, err := c.db.Query(ctx, listDeviceJoinsQuery, map[string]interface{}{
		"device": NewBlankDeviceDocument(dev.Uuid).ID(),
	})
	if err != nil {
		log.Warn("Error querying for joins", zap.Error(err))
		return nil, status.Error(codes.Internal, "Error querying for joins")
	}
	defer cr.Close()

	var r []*access.Node
	for {
		var node access.Node
		_, err := cr.ReadDocument(ctx, &node)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			log.Warn("Error unmarshalling Document", zap.Error(err))
			return nil, status.Error(codes.Internal, "Couldn't execute query")
		}
		log.Debug("Got document", zap.Any("node", &node))
		r = append(r, &node)
	}

	return &access.Nodes{Nodes: r}, nil
}

func (c *DevicesController) Join(ctx context.Context, req *pb.JoinGeneralRequest) (*access.Node, error) {
	log := c.log.Named("Join")

	requestor_id := ctx.Value(inf.InfinimeshAccountCtxKey).(string)
	log.Debug("Requestor", zap.String("id", requestor_id))

	requestor := NewBlankAccountDocument(requestor_id)
	dev := NewBlankDeviceDocument(req.Node)

	err := AccessLevelAndGet(ctx, log, c.db, requestor, dev)
	if err != nil {
		log.Warn("Error getting Device and access level", zap.Error(err))
		return nil, status.Error(codes.NotFound, "Device not found or not enough Access Rights")
	}
	if dev.Access.Role != access.Role_OWNER && dev.Access.Level < access.Level_ADMIN {
		return nil, status.Error(codes.PermissionDenied, "Not enough Access Rights")
	}

	var obj InfinimeshGraphNode
	var edge driver.Collection

	col, key := SplitDocID(req.Join)
	switch col {
	case "Accounts":
		obj = NewBlankAccountDocument(key)
		edge = c.acc2dev
	case "Namespaces":
		obj = NewBlankNamespaceDocument(key)
		edge = c.ns2dev
	}

	if obj == nil {
		return nil, status.Error(codes.InvalidArgument, "Unable to determine Object type")
	}

	err = AccessLevelAndGet(ctx, log, c.db, requestor, obj)
	if err != nil {
		log.Warn("Error getting Object and access level", zap.String("id", req.Join), zap.Error(err))
		return nil, status.Error(codes.NotFound, "Object not found or not enough Access Rights")
	}

	lvl := req.Access
	if lvl >= access.Level_ADMIN {
		return nil, status.Error(codes.InvalidArgument, "Not allowed to share Admin or Root priviliges")
	}

	err = Link(ctx, log, edge, obj, dev, req.Access, access.Role_UNSET)
	if err != nil {
		log.Warn("Error creating edge", zap.Error(err))
		return nil, status.Error(codes.Internal, "error creating Permission")
	}

	return &access.Node{
		Node: req.Join,
		Access: &access.Access{
			Level: lvl,
		},
	}, nil
}
