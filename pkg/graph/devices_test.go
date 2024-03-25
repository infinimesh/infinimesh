package graph_test

import (
	"context"
	"errors"
	"testing"

	"github.com/infinimesh/proto/node/accounts"

	"connectrpc.com/connect"
	"github.com/arangodb/go-driver"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	driver_mocks "github.com/infinimesh/infinimesh/mocks/github.com/arangodb/go-driver"
	graph_mocks "github.com/infinimesh/infinimesh/mocks/github.com/infinimesh/infinimesh/pkg/graph"
	handsfree_mocks "github.com/infinimesh/infinimesh/mocks/github.com/infinimesh/proto/handsfree"
	"github.com/infinimesh/infinimesh/pkg/graph"
	"github.com/infinimesh/infinimesh/pkg/graph/schema"
	inf "github.com/infinimesh/infinimesh/pkg/shared"
	"github.com/infinimesh/proto/handsfree"
	"github.com/infinimesh/proto/node"
	"github.com/infinimesh/proto/node/access"
	devpb "github.com/infinimesh/proto/node/devices"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/structpb"
)

type devicesControllerFixture struct {
	ctrl *graph.DevicesController

	mocks struct {
		db      *driver_mocks.MockDatabase
		col     *driver_mocks.MockCollection
		ns2dev  *driver_mocks.MockCollection
		acc2dev *driver_mocks.MockCollection

		hfc       *handsfree_mocks.MockHandsfreeServiceClient
		ica_repo  *graph_mocks.MockInfinimeshCommonActionsRepo
		repo      *graph_mocks.MockInfinimeshGenericActionsRepo[*devpb.Device]
		accs_repo *graph_mocks.MockInfinimeshGenericActionsRepo[*accounts.Account]
		bus       *graph_mocks.MockEventBusService
	}

	data struct {
		ctx      context.Context
		acc_uuid string
		dev_uuid string
		ns_uuid  string

		cert          string
		create_req    devpb.CreateRequest
		create_hf_req devpb.CreateRequest
		patch_req     devpb.Device
	}
}

func newDevicesControllerFixture(t *testing.T) *devicesControllerFixture {
	t.Parallel()
	f := &devicesControllerFixture{}

	f.mocks.db = driver_mocks.NewMockDatabase(t)
	f.mocks.col = driver_mocks.NewMockCollection(t)
	f.mocks.ns2dev = driver_mocks.NewMockCollection(t)
	f.mocks.acc2dev = driver_mocks.NewMockCollection(t)
	f.mocks.db.EXPECT().Collection(context.TODO(), schema.DEVICES_COL).Return(f.mocks.col, nil)

	f.mocks.hfc = handsfree_mocks.NewMockHandsfreeServiceClient(t)
	f.mocks.ica_repo = graph_mocks.NewMockInfinimeshCommonActionsRepo(t)
	f.mocks.ica_repo.EXPECT().GetEdgeCol(context.TODO(), schema.NS2DEV).Return(f.mocks.ns2dev).Maybe()
	f.mocks.ica_repo.EXPECT().GetEdgeCol(context.TODO(), schema.ACC2DEV).Return(f.mocks.acc2dev).Maybe()
	f.mocks.bus = graph_mocks.NewMockEventBusService(t)

	f.mocks.repo = graph_mocks.NewMockInfinimeshGenericActionsRepo[*devpb.Device](t)
	f.mocks.accs_repo = graph_mocks.NewMockInfinimeshGenericActionsRepo[*accounts.Account](t)

	f.data.acc_uuid = uuid.New().String()
	f.data.dev_uuid = uuid.New().String()
	f.data.ns_uuid = uuid.New().String()
	f.data.ctx = context.WithValue(context.Background(), inf.InfinimeshAccountCtxKey, f.data.acc_uuid)

	f.data.cert = `-----BEGIN CERTIFICATE-----
MIIEljCCAn4CCQC7oNynkLPhTjANBgkqhkiG9w0BAQsFADANMQswCQYDVQQGEwJk
ZTAeFw0yMTA2MTYxMTMyNDRaFw0yMjA2MTYxMTMyNDRaMA0xCzAJBgNVBAYTAmRl
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEA0hk6i+PxRW7XAy21QAsR
Dlyz60ojkDU5q2BfXzmo5GPGaAXuEwwT+AJGFAgIvSIWh7SBDY3re75YbShfbLEP
biHDtNKzr0v+RmNiZ66qZy7lVPyTcDe4Aj9iOsdAiocKXBECgpdvDPM2SPVsL915
oajg2RAp/VmvtHdENBjgD0e7xVXV4hKwn2UDMQbw1KBfIXVj6n7fwMvouovcmdc+
A107+HTudDqvhrkevAJXDmxTWRKz3anoU/dCcV4d1aHLys29L/vnlF0q29KEfSLJ
Ov9H/9mX/NjcmMqr4tsqjmu5ZepORhtGqq0Rmcg++FbCA4f68OchTPopvYKz7ExN
CPzgxufqduBdThIwNzdtXctm0othphQ3ADxnxCqDfAhqr02w7qaCd/c1KBK6EKvJ
uIWiqaVV3ipqre+T98AuzJ7il+mhIsRsXpBt3o7LBCgyl8rri+ZLEDRj3hOu3UN5
pS71R0xm62P8psKY0xtDneReUQ1CGObQS7XZDCJ0qlHDGUMTBwvGbcqrTwpA1udu
cP1GGDhRsdlx0NgJEemSojEiMKmSc1McNsubczfJCZAZRNNvR7pn4MyyS20aMNnd
1rRkX6ikyvRA96dJD0M4iI2f6asNpGe8SplwPJweNv/avwYiWKFVO5neuVEdiAcw
XjFL9u8OK0ID8Uid3TWV4psCAwEAATANBgkqhkiG9w0BAQsFAAOCAgEALKx4BlYg
dizAl5jVICrswgVlS/Ec8dw3hTmuDodhA5jP5NLFIrzWHp6voythjhFIdXHI+8nW
y0V1TVviW73qFP9ib5LnLn30QVajwFRjBIOt4qsrIvMFDvwtQ940pUgR1iVGphV4
ahlCwNeZStdxMV8M4/5o78wP7uvyhleIaYrF7dLfFoszT4PfyRC2UEXtTknz1hH8
kOFwiZCio5sIzWNsAzHlOKbf2Rl0WtC9YWcKpdS1MrWi6E/jAJQ1/GyhUOEZHE/Z
fY1heN2YXPacYtFQTRmkp/oPzsIvwgfx6OKJe8RGa7EErQUVGTMYkZue7lpIOyJD
8m37TUVNizW2+OrQb/NUK9uwEBkGlpavTdK7eKAw0+KnlPqMpmQx7Vs5oE0ejy7y
GuMpc8AeJXUX9lHMJIT+lwkKzrVReC+jgyvO0QyRN7PTwRW8+9SNOeHRiC9Fj7Zg
fLCCa/hdALN6ECHn3JsQGiAbY6JS8LOdiLpnlR+cOQSQ3HnaBkpPeBmWfRvlvGeU
r+vyP3YimFBE9AbM5GgfUHGRBJBpC40aVaE7HtHapE4JJNit4NfBvfDotNUs6shJ
6Y893NPueYB4PfvC+1kgZFjXFEMDURaGUeEwl481Zn/rGXM4ev5qGPQgJ4fhmI68
cgSqKFgDFRxlHXLo9TZnxyBrIvN/siE+ZQI=
-----END CERTIFICATE-----`

	f.data.create_req = devpb.CreateRequest{
		Device: &devpb.Device{
			Title:   faker.Username(),
			Enabled: false,
			Certificate: &devpb.Certificate{
				PemData: f.data.cert,
			},
			Tags:         []string{"sample", "tags", "as:well", "with-special", "chars!"},
			BasicEnabled: false,
			Config:       &structpb.Struct{},
		},
		Namespace: f.data.ns_uuid,
	}

	f.data.create_hf_req = devpb.CreateRequest{
		Device: &devpb.Device{
			Title:        faker.Username(),
			Enabled:      false,
			Tags:         []string{"sample", "tags", "as:well", "with-special", "chars!"},
			BasicEnabled: false,
			Config:       &structpb.Struct{},
		},
		Namespace: f.data.ns_uuid,
		Handsfree: &devpb.HandsfreeCreate{
			Code:    "123456",
			Payload: []string{f.data.cert},
		},
	}

	f.data.patch_req = devpb.Device{
		Config: &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"Message": structpb.NewStringValue("Hello !!!"),
			},
		},
	}

	f.ctrl = graph.NewDevicesController(
		zap.NewExample(), f.mocks.db,
		f.mocks.hfc, f.mocks.ica_repo,
		f.mocks.repo,
		f.mocks.bus,
	)

	return f
}

// Create
//

func TestCreate_FailsOn_NoNamespace(t *testing.T) {
	f := newDevicesControllerFixture(t)
	f.data.create_req.Namespace = ""
	res, err := f.ctrl.Create(f.data.ctx, connect.NewRequest(&f.data.create_req))
	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Namespace ID is required")
}

func TestCreate_FailsOn_NoAccessToNamespace(t *testing.T) {
	f := newDevicesControllerFixture(t)
	f.mocks.ica_repo.EXPECT().AccessLevel(f.data.ctx, mock.Anything, graph.NewBlankNamespaceDocument(f.data.ns_uuid)).Return(
		false, access.Level_NONE,
	)
	res, err := f.ctrl.Create(f.data.ctx, connect.NewRequest(&f.data.create_req))
	assert.Nil(t, res)
	assert.Error(t, err)
	assert.EqualError(t, err, "permission_denied: No Access to Namespace "+f.data.ns_uuid)
}

func TestCreate_FailsOn_GenerateFingerprint(t *testing.T) {
	f := newDevicesControllerFixture(t)
	f.mocks.ica_repo.EXPECT().AccessLevel(f.data.ctx, mock.Anything, graph.NewBlankNamespaceDocument(f.data.ns_uuid)).Return(
		true, access.Level_ADMIN,
	)
	f.data.create_req.Device.Certificate.PemData = "invalid"
	res, err := f.ctrl.Create(f.data.ctx, connect.NewRequest(&f.data.create_req))
	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Can't generate fingerprint")
}

func TestCreate_FailsOn_CreateDocument(t *testing.T) {
	f := newDevicesControllerFixture(t)
	f.mocks.ica_repo.EXPECT().AccessLevel(f.data.ctx, mock.Anything, graph.NewBlankNamespaceDocument(f.data.ns_uuid)).Return(
		true, access.Level_ADMIN,
	)
	f.mocks.col.EXPECT().CreateDocument(f.data.ctx, mock.Anything).Return(driver.DocumentMeta{}, assert.AnError)
	res, err := f.ctrl.Create(f.data.ctx, connect.NewRequest(&f.data.create_req))
	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Error while creating Device")
}

func TestCreate_FailsOn_Link(t *testing.T) {
	f := newDevicesControllerFixture(t)
	f.mocks.ica_repo.EXPECT().AccessLevel(f.data.ctx, mock.Anything, graph.NewBlankNamespaceDocument(f.data.ns_uuid)).Return(
		true, access.Level_ADMIN,
	)
	f.mocks.col.EXPECT().CreateDocument(f.data.ctx, mock.Anything).Return(driver.DocumentMeta{
		ID: driver.NewDocumentID(schema.DEVICES_COL, f.data.dev_uuid),
	}, nil)
	f.mocks.ica_repo.EXPECT().Link(
		f.data.ctx, f.mocks.ns2dev,
		graph.NewBlankNamespaceDocument(f.data.ns_uuid),
		mock.Anything,
		access.Level_ADMIN, access.Role_OWNER,
	).Return(assert.AnError)
	f.mocks.col.EXPECT().RemoveDocument(f.data.ctx, f.data.dev_uuid).Return(driver.DocumentMeta{}, nil)

	f.mocks.repo.EXPECT().UpdateDeviceModifyDate(f.data.ctx, mock.Anything, f.data.dev_uuid).Return(nil)

	res, err := f.ctrl.Create(f.data.ctx, connect.NewRequest(&f.data.create_req))
	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error creating Permission")

	f.mocks.col.AssertNumberOfCalls(t, "RemoveDocument", 1)
}

func TestCreate_Success(t *testing.T) {
	f := newDevicesControllerFixture(t)
	f.mocks.ica_repo.EXPECT().AccessLevel(f.data.ctx, mock.Anything, graph.NewBlankNamespaceDocument(f.data.ns_uuid)).Return(
		true, access.Level_ADMIN,
	)
	f.mocks.col.EXPECT().CreateDocument(f.data.ctx, mock.Anything).Return(driver.DocumentMeta{
		ID: driver.NewDocumentID(schema.DEVICES_COL, f.data.dev_uuid),
	}, nil)
	f.mocks.ica_repo.EXPECT().Link(
		f.data.ctx, f.mocks.ns2dev,
		graph.NewBlankNamespaceDocument(f.data.ns_uuid),
		mock.Anything,
		access.Level_ADMIN, access.Role_OWNER,
	).Return(nil)

	f.mocks.bus.EXPECT().Notify(f.data.ctx, mock.Anything).Return(func() error {
		return nil
	}, nil)

	f.mocks.repo.EXPECT().UpdateDeviceModifyDate(f.data.ctx, mock.Anything, f.data.dev_uuid).Return(nil)

	res, err := f.ctrl.Create(f.data.ctx, connect.NewRequest(&f.data.create_req))
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, f.data.dev_uuid, res.Msg.GetDevice().GetUuid())
}

// Create with Handfree
//

func TestCreateHf_FailsOn_NoNamespace(t *testing.T) {
	f := newDevicesControllerFixture(t)
	f.data.create_hf_req.Namespace = ""
	res, err := f.ctrl.Create(f.data.ctx, connect.NewRequest(&f.data.create_hf_req))
	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Namespace ID is required")
}

func TestCreateHf_FailsOn_NoAccessToNamespace(t *testing.T) {
	f := newDevicesControllerFixture(t)
	f.mocks.ica_repo.EXPECT().AccessLevel(f.data.ctx, mock.Anything, graph.NewBlankNamespaceDocument(f.data.ns_uuid)).Return(
		false, access.Level_NONE,
	)
	res, err := f.ctrl.Create(f.data.ctx, connect.NewRequest(&f.data.create_hf_req))
	assert.Nil(t, res)
	assert.Error(t, err)
	assert.EqualError(t, err, "permission_denied: No Access to Namespace "+f.data.ns_uuid)
}

func TestCreateHf_FailsOn_CreateDocument(t *testing.T) {
	f := newDevicesControllerFixture(t)
	f.mocks.ica_repo.EXPECT().AccessLevel(f.data.ctx, mock.Anything, graph.NewBlankNamespaceDocument(f.data.ns_uuid)).Return(
		true, access.Level_ADMIN,
	)
	f.mocks.col.EXPECT().CreateDocument(f.data.ctx, mock.Anything).Return(driver.DocumentMeta{}, assert.AnError)
	res, err := f.ctrl.Create(f.data.ctx, connect.NewRequest(&f.data.create_hf_req))
	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Error while creating Device")
}

func TestCreateHf_FailsOn_Link(t *testing.T) {
	f := newDevicesControllerFixture(t)
	f.mocks.ica_repo.EXPECT().AccessLevel(f.data.ctx, mock.Anything, graph.NewBlankNamespaceDocument(f.data.ns_uuid)).Return(
		true, access.Level_ADMIN,
	)
	f.mocks.col.EXPECT().CreateDocument(f.data.ctx, mock.Anything).Return(driver.DocumentMeta{
		ID: driver.NewDocumentID(schema.DEVICES_COL, f.data.dev_uuid),
	}, nil)
	f.mocks.ica_repo.EXPECT().
		Link(f.data.ctx,
			f.mocks.ns2dev, graph.NewBlankNamespaceDocument(f.data.ns_uuid),
			mock.Anything,
			access.Level_ADMIN, access.Role_OWNER,
		).Return(assert.AnError)
	f.mocks.col.EXPECT().RemoveDocument(f.data.ctx, f.data.dev_uuid).Return(driver.DocumentMeta{}, nil)

	res, err := f.ctrl.Create(f.data.ctx, connect.NewRequest(&f.data.create_hf_req))
	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error creating Permission")

	f.mocks.col.AssertNumberOfCalls(t, "RemoveDocument", 1)
}

func TestCreateHf_FailsOn_Send(t *testing.T) {
	f := newDevicesControllerFixture(t)
	f.mocks.ica_repo.EXPECT().AccessLevel(f.data.ctx, mock.Anything, graph.NewBlankNamespaceDocument(f.data.ns_uuid)).Return(
		true, access.Level_ADMIN,
	)
	f.mocks.col.EXPECT().CreateDocument(f.data.ctx, mock.Anything).Return(driver.DocumentMeta{
		ID: driver.NewDocumentID(schema.DEVICES_COL, f.data.dev_uuid),
	}, nil)
	f.mocks.ica_repo.EXPECT().
		Link(f.data.ctx,
			f.mocks.ns2dev, graph.NewBlankNamespaceDocument(f.data.ns_uuid),
			mock.Anything,
			access.Level_ADMIN, access.Role_OWNER,
		).Return(nil)
	f.mocks.col.EXPECT().RemoveDocument(f.data.ctx, f.data.dev_uuid).Return(driver.DocumentMeta{}, nil)

	f.mocks.ica_repo.EXPECT().AccessLevelAndGet(f.data.ctx, mock.Anything, mock.MatchedBy(func(d *graph.Device) bool {
		d.Access = &access.Access{
			Level:     access.Level_ADMIN,
			Namespace: &f.data.ns_uuid,
		}
		return true
	})).Return(nil)

	f.mocks.bus.EXPECT().Notify(f.data.ctx, mock.Anything).Return(func() error {
		return nil
	}, nil)

	f.mocks.col.EXPECT().RemoveDocument(f.data.ctx, f.data.dev_uuid).Return(driver.DocumentMeta{}, nil)
	f.mocks.ica_repo.EXPECT().
		Link(f.data.ctx, f.mocks.ns2dev,
			mock.Anything, mock.Anything, access.Level_NONE, access.Role_UNSET,
		).Return(assert.AnError)

	f.mocks.hfc.EXPECT().Send(f.data.ctx, mock.Anything).Return(nil, assert.AnError)

	res, err := f.ctrl.Create(f.data.ctx, connect.NewRequest(&f.data.create_hf_req))
	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), assert.AnError.Error())
}

func TestCreateHf_FailsOn_EmptyPayloadAndDelete(t *testing.T) {
	f := newDevicesControllerFixture(t)
	f.mocks.ica_repo.EXPECT().AccessLevel(f.data.ctx, mock.Anything, graph.NewBlankNamespaceDocument(f.data.ns_uuid)).Return(
		true, access.Level_ADMIN,
	)
	f.mocks.col.EXPECT().CreateDocument(f.data.ctx, mock.Anything).Return(driver.DocumentMeta{
		ID: driver.NewDocumentID(schema.DEVICES_COL, f.data.dev_uuid),
	}, nil)
	f.mocks.ica_repo.EXPECT().
		Link(f.data.ctx,
			f.mocks.ns2dev, graph.NewBlankNamespaceDocument(f.data.ns_uuid),
			mock.Anything,
			access.Level_ADMIN, access.Role_OWNER,
		).Return(nil)

	f.mocks.ica_repo.EXPECT().AccessLevelAndGet(f.data.ctx, mock.Anything, mock.Anything).Return(assert.AnError)

	f.mocks.hfc.EXPECT().Send(f.data.ctx, mock.Anything).Return(&handsfree.ControlPacket{
		Payload: []string{},
	}, nil)

	res, err := f.ctrl.Create(f.data.ctx, connect.NewRequest(&f.data.create_hf_req))
	assert.NotNil(t, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "issue with Handsfree Payload: is empty")
}

func TestCreateHf_FailsOn_GenerateFingerprint(t *testing.T) {
	f := newDevicesControllerFixture(t)
	f.mocks.ica_repo.EXPECT().AccessLevel(f.data.ctx, mock.Anything, graph.NewBlankNamespaceDocument(f.data.ns_uuid)).Return(
		true, access.Level_ADMIN,
	)
	f.mocks.col.EXPECT().CreateDocument(f.data.ctx, mock.Anything).Return(driver.DocumentMeta{
		ID: driver.NewDocumentID(schema.DEVICES_COL, f.data.dev_uuid),
	}, nil)
	f.mocks.ica_repo.EXPECT().
		Link(f.data.ctx,
			f.mocks.ns2dev, graph.NewBlankNamespaceDocument(f.data.ns_uuid),
			mock.Anything,
			access.Level_ADMIN, access.Role_OWNER,
		).Return(nil)
	f.mocks.col.EXPECT().RemoveDocument(f.data.ctx, f.data.dev_uuid).Return(driver.DocumentMeta{}, nil)

	f.mocks.bus.EXPECT().Notify(f.data.ctx, mock.Anything).Return(func() error {
		return nil
	}, nil)

	f.mocks.ica_repo.EXPECT().AccessLevelAndGet(f.data.ctx, mock.Anything, mock.MatchedBy(func(d *graph.Device) bool {
		d.Access = &access.Access{
			Level:     access.Level_ADMIN,
			Namespace: &f.data.ns_uuid,
		}
		return true
	})).Return(nil)

	f.mocks.col.EXPECT().RemoveDocument(f.data.ctx, f.data.dev_uuid).Return(driver.DocumentMeta{}, assert.AnError)
	f.mocks.ica_repo.EXPECT().
		Link(f.data.ctx, f.mocks.ns2dev,
			mock.Anything, mock.Anything, access.Level_NONE, access.Role_UNSET,
		).Return(nil)

	app_id := "some-app"
	f.mocks.hfc.EXPECT().Send(f.data.ctx, mock.Anything).Return(&handsfree.ControlPacket{
		Payload: []string{"invalid cert"},
		AppId:   &app_id,
	}, nil)

	res, err := f.ctrl.Create(f.data.ctx, connect.NewRequest(&f.data.create_hf_req))
	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "coudn't decode PEM data")
}

func TestCreateHf_FailsOn_ReplaceDocument(t *testing.T) {
	f := newDevicesControllerFixture(t)
	f.mocks.ica_repo.EXPECT().AccessLevel(f.data.ctx, mock.Anything, graph.NewBlankNamespaceDocument(f.data.ns_uuid)).Return(
		true, access.Level_ADMIN,
	)
	f.mocks.col.EXPECT().CreateDocument(f.data.ctx, mock.Anything).Return(driver.DocumentMeta{
		ID: driver.NewDocumentID(schema.DEVICES_COL, f.data.dev_uuid),
	}, nil)
	f.mocks.ica_repo.EXPECT().
		Link(f.data.ctx,
			f.mocks.ns2dev, graph.NewBlankNamespaceDocument(f.data.ns_uuid),
			mock.Anything,
			access.Level_ADMIN, access.Role_OWNER,
		).Return(nil)
	f.mocks.col.EXPECT().RemoveDocument(f.data.ctx, f.data.dev_uuid).Return(driver.DocumentMeta{}, nil)

	f.mocks.ica_repo.EXPECT().AccessLevelAndGet(f.data.ctx, mock.Anything, mock.MatchedBy(func(d *graph.Device) bool {
		d.Access = &access.Access{
			Level:     access.Level_ADMIN,
			Namespace: &f.data.ns_uuid,
		}
		return true
	})).Return(nil)

	f.mocks.col.EXPECT().RemoveDocument(f.data.ctx, f.data.dev_uuid).Return(driver.DocumentMeta{}, assert.AnError)
	f.mocks.ica_repo.EXPECT().
		Link(f.data.ctx, f.mocks.ns2dev,
			mock.Anything, mock.Anything, access.Level_NONE, access.Role_UNSET,
		).Return(nil)

	app_id := "some-app"
	f.mocks.hfc.EXPECT().Send(f.data.ctx, mock.Anything).Return(&handsfree.ControlPacket{
		Payload: []string{f.data.cert},
		AppId:   &app_id,
	}, nil)

	f.mocks.col.EXPECT().ReplaceDocument(f.data.ctx, mock.Anything, mock.Anything).Return(driver.DocumentMeta{}, assert.AnError)

	f.mocks.bus.EXPECT().Notify(f.data.ctx, mock.Anything).Return(func() error {
		return nil
	}, nil)

	res, err := f.ctrl.Create(f.data.ctx, connect.NewRequest(&f.data.create_hf_req))
	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), assert.AnError.Error())
}

func TestCreateHf_Success(t *testing.T) {
	f := newDevicesControllerFixture(t)
	f.mocks.ica_repo.EXPECT().AccessLevel(f.data.ctx, mock.Anything, graph.NewBlankNamespaceDocument(f.data.ns_uuid)).Return(
		true, access.Level_ADMIN,
	)
	f.mocks.col.EXPECT().CreateDocument(f.data.ctx, mock.Anything).Return(driver.DocumentMeta{
		ID: driver.NewDocumentID(schema.DEVICES_COL, f.data.dev_uuid),
	}, nil)
	f.mocks.ica_repo.EXPECT().
		Link(f.data.ctx,
			f.mocks.ns2dev, mock.Anything,
			mock.Anything,
			access.Level_ADMIN, access.Role_OWNER,
		).Return(nil)

	app_id := "some-app"
	f.mocks.hfc.EXPECT().Send(f.data.ctx, mock.Anything).Return(&handsfree.ControlPacket{
		Payload: []string{f.data.cert},
		AppId:   &app_id,
	}, nil)

	f.mocks.col.EXPECT().ReplaceDocument(f.data.ctx, mock.Anything, mock.Anything).Return(driver.DocumentMeta{}, nil)

	res, err := f.ctrl.Create(f.data.ctx, connect.NewRequest(&f.data.create_hf_req))
	assert.NoError(t, err)
	assert.NotNil(t, res)
}

// Delete
//

func TestDelete_FailsOn_AccessLevelAndGet(t *testing.T) {
	f := newDevicesControllerFixture(t)

	f.mocks.ica_repo.EXPECT().AccessLevelAndGet(f.data.ctx, mock.Anything, mock.Anything).Return(
		assert.AnError,
	)
	res, err := f.ctrl.Delete(f.data.ctx, connect.NewRequest(&devpb.Device{
		Uuid: f.data.dev_uuid,
	}))

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, "not_found: Device not found or not enough Access Rights", err.Error())
}

func TestDelete_FailsOn_AccessLevel_NotEnoughAccess(t *testing.T) {
	f := newDevicesControllerFixture(t)

	f.mocks.ica_repo.EXPECT().AccessLevelAndGet(f.data.ctx, mock.Anything, mock.MatchedBy(func(d *graph.Device) bool {
		d.Access = &access.Access{
			Level: access.Level_READ,
		}
		return true
	})).Return(
		nil,
	)

	res, err := f.ctrl.Delete(f.data.ctx, connect.NewRequest(&devpb.Device{
		Uuid: f.data.dev_uuid,
	}))

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, "permission_denied: Not enough Access Rights", err.Error())
}

func TestDelete_FailsOn_DeleteDocument(t *testing.T) {
	f := newDevicesControllerFixture(t)

	f.mocks.ica_repo.EXPECT().AccessLevelAndGet(f.data.ctx, mock.Anything, mock.MatchedBy(func(d *graph.Device) bool {
		d.Access = &access.Access{
			Level: access.Level_ADMIN,
		}
		return true
	})).Return(nil)

	f.mocks.bus.EXPECT().Notify(f.data.ctx, mock.Anything).Return(func() error {
		return nil
	}, nil)

	f.mocks.col.EXPECT().RemoveDocument(f.data.ctx, f.data.dev_uuid).Return(driver.DocumentMeta{}, assert.AnError)

	res, err := f.ctrl.Delete(f.data.ctx, connect.NewRequest(&devpb.Device{
		Uuid: f.data.dev_uuid,
	}))

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, "internal: Error while deleting Device", err.Error())
}

func TestDelete_Success(t *testing.T) {
	f := newDevicesControllerFixture(t)

	f.mocks.ica_repo.EXPECT().AccessLevelAndGet(f.data.ctx, mock.Anything, mock.MatchedBy(func(d *graph.Device) bool {
		d.Access = &access.Access{
			Level:     access.Level_ADMIN,
			Namespace: &f.data.ns_uuid,
		}
		return true
	})).Return(nil)

	f.mocks.col.EXPECT().RemoveDocument(f.data.ctx, f.data.dev_uuid).Return(driver.DocumentMeta{}, nil)
	f.mocks.ica_repo.EXPECT().
		Link(f.data.ctx, f.mocks.ns2dev,
			mock.Anything, mock.Anything, access.Level_NONE, access.Role_UNSET,
		).Return(assert.AnError)

	f.mocks.bus.EXPECT().Notify(f.data.ctx, mock.Anything).Return(func() error {
		return nil
	}, nil)

	res, err := f.ctrl.Delete(f.data.ctx, connect.NewRequest(&devpb.Device{
		Uuid: f.data.dev_uuid,
	}))

	assert.NoError(t, err)
	assert.NotNil(t, res)
}

// MakeDevicesToken
//

func TestMakeDevicesToken_FailsOn_AccessLevel_NotFound(t *testing.T) {
	f := newDevicesControllerFixture(t)

	f.mocks.ica_repo.EXPECT().AccessLevel(f.data.ctx, mock.Anything, mock.Anything).Return(
		false, access.Level_NONE,
	)
	res, err := f.ctrl.MakeDevicesToken(f.data.ctx, connect.NewRequest(&node.DevicesTokenRequest{
		Devices: map[string]access.Level{f.data.dev_uuid: access.Level_NONE},
	}))

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Account not found or not enough Access Rights to device:")
}

func TestMakeDevicesToken_FailsOn_AccessLevel_NotEnoughAccess(t *testing.T) {
	f := newDevicesControllerFixture(t)

	f.mocks.ica_repo.EXPECT().AccessLevel(f.data.ctx, mock.Anything, mock.Anything).Return(
		true, access.Level_READ,
	)
	res, err := f.ctrl.MakeDevicesToken(f.data.ctx, connect.NewRequest(&node.DevicesTokenRequest{
		Devices: map[string]access.Level{f.data.dev_uuid: access.Level_ADMIN},
	}))

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Not enough Access Rights to device:")
}

// TODO: Fix this test by adding mock

// func TestMakeDevicesToken_FailsOn_MakeToken(t *testing.T) {
// 	f := newDevicesControllerFixture(t)

// 	f.mocks.ica_repo.EXPECT().AccessLevel(f.data.ctx, mock.Anything, mock.Anything).Return(
// 		true, access.Level_ADMIN,
// 	)

// 	// Mock MakeToken or cause signing error

// 	res, err := f.ctrl.MakeDevicesToken(f.data.ctx, connect.NewRequest(&node.DevicesTokenRequest{
// 		Devices: []string{f.data.uuid},
// 		Post:    true,
// 	}))

// 	assert.Nil(t, res)
// 	assert.Error(t, err)
// 	assert.Contains(t, err.Error(), "Not enough Access Rights to device:")
// }

func TestMakeDevicesToken_Success(t *testing.T) {
	f := newDevicesControllerFixture(t)

	f.mocks.ica_repo.EXPECT().AccessLevel(f.data.ctx, mock.Anything, mock.Anything).Return(
		true, access.Level_ADMIN,
	)

	res, err := f.ctrl.MakeDevicesToken(f.data.ctx, connect.NewRequest(&node.DevicesTokenRequest{
		Devices: map[string]access.Level{f.data.dev_uuid: access.Level_NONE},
	}))

	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestList_Success(t *testing.T) {
	f := newDevicesControllerFixture(t)

	count := 0
	result := &graph.ListQueryResult[*devpb.Device]{
		Result: []*devpb.Device{},
		Count:  count,
	}
	f.mocks.repo.EXPECT().ListQuery(mock.Anything, mock.Anything, mock.Anything).Return(result, nil)

	resp, err := f.ctrl.List(f.data.ctx, connect.NewRequest(&node.QueryRequest{}))

	assert.NoError(t, err)
	assert.Equal(t, count, len(resp.Msg.Devices))
}

func TestList_FailsOn_ListQuery(t *testing.T) {
	f := newDevicesControllerFixture(t)

	f.mocks.repo.EXPECT().ListQuery(mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("Error"))

	_, err := f.ctrl.List(f.data.ctx, connect.NewRequest(&node.QueryRequest{}))

	assert.Error(t, err)
}

func TestPatchConfig_Success(t *testing.T) {
	f := newDevicesControllerFixture(t)

	f.mocks.ica_repo.EXPECT().AccessLevelAndGet(f.data.ctx, mock.Anything, mock.MatchedBy(func(d *graph.Device) bool {
		d.Access = &access.Access{
			Level:     access.Level_ROOT,
			Namespace: &f.data.ns_uuid,
		}
		return true
	})).Return(nil)

	f.mocks.col.EXPECT().ReplaceDocument(f.data.ctx, mock.Anything, mock.MatchedBy(func(d *devpb.Device) bool {
		return true
	})).Return(driver.DocumentMeta{}, nil)

	f.mocks.bus.EXPECT().Notify(f.data.ctx, mock.Anything).Return(func() error {
		return nil
	}, nil)

	res, err := f.ctrl.PatchConfig(f.data.ctx, connect.NewRequest(&f.data.patch_req))

	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestPatchConfig_FailsOn_NoAccess(t *testing.T) {
	f := newDevicesControllerFixture(t)

	f.mocks.ica_repo.EXPECT().AccessLevelAndGet(f.data.ctx, mock.Anything, mock.MatchedBy(func(d *graph.Device) bool {
		d.Access = &access.Access{
			Level:     access.Level_MGMT,
			Namespace: &f.data.ns_uuid,
		}
		return true
	})).Return(nil)

	res, err := f.ctrl.PatchConfig(f.data.ctx, connect.NewRequest(&f.data.patch_req))

	assert.Error(t, err)
	assert.Nil(t, res)
}
