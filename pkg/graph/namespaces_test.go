package graph_test

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	"github.com/arangodb/go-driver"
	"github.com/google/uuid"
	driver_mocks "github.com/infinimesh/infinimesh/mocks/github.com/arangodb/go-driver"
	graph_mocks "github.com/infinimesh/infinimesh/mocks/github.com/infinimesh/infinimesh/pkg/graph"
	"github.com/infinimesh/infinimesh/pkg/graph"
	"github.com/infinimesh/infinimesh/pkg/graph/schema"
	inf "github.com/infinimesh/infinimesh/pkg/shared"
	"github.com/infinimesh/proto/node/access"
	"github.com/infinimesh/proto/node/namespaces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

type namespacesServiceFixture struct {
	ctrl *graph.NamespacesController

	mocks struct {
		db *driver_mocks.MockDatabase

		ns_col  *driver_mocks.MockCollection
		acc_col *driver_mocks.MockCollection
		acc2ns  *driver_mocks.MockCollection
		ns2acc  *driver_mocks.MockCollection

		bus *graph_mocks.MockEventBusService

		log      *zap.Logger
		observer *observer.ObservedLogs

		ica  *graph_mocks.MockInfinimeshCommonActionsRepo
		repo *graph_mocks.MockInfinimeshGenericActionsRepo[*namespaces.Namespace]
	}
	data struct {
		ctx context.Context

		acc_uuid string
		dev_uuid string
		ns_uuid  string
	}
}

func newNamespacesServiceFixture(t *testing.T) *namespacesServiceFixture {
	f := &namespacesServiceFixture{}

	f.mocks.db = driver_mocks.NewMockDatabase(t)
	f.mocks.bus = graph_mocks.NewMockEventBusService(t)

	f.mocks.ns_col = driver_mocks.NewMockCollection(t)
	f.mocks.db.EXPECT().Collection(context.TODO(), schema.NAMESPACES_COL).Return(f.mocks.ns_col, nil)

	f.mocks.acc_col = driver_mocks.NewMockCollection(t)
	f.mocks.db.EXPECT().Collection(context.TODO(), schema.ACCOUNTS_COL).Return(f.mocks.acc_col, nil)

	f.mocks.acc2ns = driver_mocks.NewMockCollection(t)
	f.mocks.ns2acc = driver_mocks.NewMockCollection(t)

	core, observer := observer.New(zap.DebugLevel)
	f.mocks.log = zap.New(core)
	f.mocks.observer = observer

	f.mocks.ica = graph_mocks.NewMockInfinimeshCommonActionsRepo(t)
	f.mocks.ica.EXPECT().GetEdgeCol(context.TODO(), schema.ACC2NS).Return(f.mocks.acc2ns)
	f.mocks.ica.EXPECT().GetEdgeCol(context.TODO(), schema.NS2ACC).Return(f.mocks.ns2acc)

	f.mocks.repo = graph_mocks.NewMockInfinimeshGenericActionsRepo[*namespaces.Namespace](t)

	f.ctrl = graph.NewNamespacesController(f.mocks.log, f.mocks.db, f.mocks.bus, f.mocks.ica, f.mocks.repo)

	f.data.acc_uuid = uuid.New().String()
	f.data.dev_uuid = uuid.New().String()
	f.data.ns_uuid = uuid.New().String()
	f.data.ctx = context.WithValue(context.Background(), inf.InfinimeshAccountCtxKey, f.data.acc_uuid)

	return f
}

// Create
//

func TestNamespaceCreate_FailsOn_NoTitle(t *testing.T) {
	f := newNamespacesServiceFixture(t)

	_, err := f.ctrl.Create(f.data.ctx, connect.NewRequest(&namespaces.Namespace{}))
	assert.Error(t, err)
	assert.EqualError(t, err, "invalid_argument: Title is required")
}

func TestNamespaceCreate_FailsOn_CreateDocument(t *testing.T) {
	f := newNamespacesServiceFixture(t)

	f.mocks.ns_col.EXPECT().CreateDocument(f.data.ctx, graph.Namespace{
		Namespace: &namespaces.Namespace{Title: "test"},
	}).Return(driver.DocumentMeta{}, assert.AnError)

	_, err := f.ctrl.Create(f.data.ctx, connect.NewRequest(&namespaces.Namespace{Title: "test"}))
	assert.Error(t, err)
	assert.EqualError(t, err, "internal: Error while creating namespace")
}

func TestNamespaceCreate_FailsOn_Link(t *testing.T) {
	f := newNamespacesServiceFixture(t)

	f.mocks.ns_col.EXPECT().CreateDocument(f.data.ctx, graph.Namespace{
		Namespace: &namespaces.Namespace{Title: "test"},
	}).Return(driver.DocumentMeta{Key: f.data.ns_uuid}, nil)

	f.mocks.ica.EXPECT().Link(f.data.ctx, mock.Anything, mock.Anything, mock.Anything, access.Level_ADMIN, access.Role_OWNER).
		Return(assert.AnError)

	f.mocks.ns_col.EXPECT().RemoveDocument(f.data.ctx, mock.Anything).Return(driver.DocumentMeta{}, nil)

	_, err := f.ctrl.Create(f.data.ctx, connect.NewRequest(&namespaces.Namespace{Title: "test"}))
	assert.Error(t, err)
	assert.EqualError(t, err, "internal: Error creating Permission")
}

func TestNamespaceCreate_Success_FailsToCreateNotifier(t *testing.T) {
	f := newNamespacesServiceFixture(t)

	f.mocks.ns_col.EXPECT().CreateDocument(f.data.ctx, graph.Namespace{
		Namespace: &namespaces.Namespace{Title: "test"},
	}).Return(driver.DocumentMeta{Key: f.data.ns_uuid}, nil)

	f.mocks.ica.EXPECT().Link(f.data.ctx, mock.Anything, mock.Anything, mock.Anything, access.Level_ADMIN, access.Role_OWNER).
		Return(nil)

	f.mocks.bus.EXPECT().Notify(f.data.ctx, mock.Anything).Return(nil, assert.AnError)

	ns, err := f.ctrl.Create(f.data.ctx, connect.NewRequest(&namespaces.Namespace{Title: "test"}))
	assert.NoError(t, err)
	assert.Equal(t, "test", ns.Msg.Title)
}

func TestNamespaceCreate_Success_FailsToNotify(t *testing.T) {
	f := newNamespacesServiceFixture(t)

	f.mocks.ns_col.EXPECT().CreateDocument(f.data.ctx, graph.Namespace{
		Namespace: &namespaces.Namespace{Title: "test"},
	}).Return(driver.DocumentMeta{Key: f.data.ns_uuid}, nil)

	f.mocks.ica.EXPECT().Link(f.data.ctx, mock.Anything, mock.Anything, mock.Anything, access.Level_ADMIN, access.Role_OWNER).
		Return(nil)

	f.mocks.bus.EXPECT().Notify(f.data.ctx, mock.Anything).Return(func() error {
		return assert.AnError
	}, nil)

	ns, err := f.ctrl.Create(f.data.ctx, connect.NewRequest(&namespaces.Namespace{Title: "test"}))
	assert.NoError(t, err)
	assert.Equal(t, "test", ns.Msg.Title)
}
