package graph_test

import (
	"context"
	"testing"

	"github.com/arangodb/go-driver"
	driver_mocks "github.com/infinimesh/infinimesh/mocks/github.com/arangodb/go-driver"
	graph_mocks "github.com/infinimesh/infinimesh/mocks/github.com/infinimesh/infinimesh/pkg/graph"
	"github.com/infinimesh/infinimesh/pkg/graph"
	"github.com/infinimesh/infinimesh/pkg/graph/schema"
	"github.com/infinimesh/proto/node/access"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type icaFixture struct {
	ica graph.InfinimeshCommonActionsRepo

	mocks struct {
		db *driver_mocks.MockDatabase
		ic *graph_mocks.MockInfinimeshController
	}
	data struct {
		ctx context.Context
	}
}

func newIcaFixture(t *testing.T) icaFixture {
	t.Parallel()
	f := icaFixture{}

	f.mocks.db = driver_mocks.NewMockDatabase(t)

	f.mocks.ic = graph_mocks.NewMockInfinimeshController(t)
	f.mocks.ic.EXPECT().DB().Return(f.mocks.db).Maybe()
	f.mocks.ic.EXPECT().Log().Return(zap.NewExample()).Maybe()

	f.ica = graph.NewInfinimeshCommonActionsRepo(zap.NewExample(), f.mocks.db)

	f.data.ctx = context.TODO()

	return f
}

// GetVertexCol
//

func TestGetVertexCol_Success(t *testing.T) {
	f := newIcaFixture(t)

	g := &driver_mocks.MockGraph{}
	f.mocks.db.EXPECT().Graph(f.data.ctx, "some graph").
		Return(g, nil)

	col := &driver_mocks.MockCollection{}
	g.EXPECT().VertexCollection(f.data.ctx, "some collection").
		Return(col, nil)

	actualCol := f.ica.GetVertexCol(f.data.ctx, "some graph", "some collection")
	assert.Equal(t, col, actualCol)
}

// GetEdgeCol
//

func TestGetEdgeCol_Success(t *testing.T) {
	f := newIcaFixture(t)

	g := &driver_mocks.MockGraph{}
	f.mocks.db.EXPECT().Graph(f.data.ctx, schema.PERMISSIONS_GRAPH.Name).
		Return(g, nil)

	col := &driver_mocks.MockCollection{}
	g.EXPECT().EdgeCollection(f.data.ctx, "some collection").
		Return(col, driver.VertexConstraints{}, nil)

	actualCol := f.ica.GetEdgeCol(f.data.ctx, "some collection")
	assert.Equal(t, col, actualCol)
}

// CheckLink
//

func TestCheckLink_Success(t *testing.T) {
	f := newIcaFixture(t)

	edge := &driver_mocks.MockCollection{}
	from := graph.NewBlankAccountDocument("uuid")
	to := graph.NewBlankAccountDocument("uuid")

	edge.EXPECT().DocumentExists(f.data.ctx, "uuid-uuid").
		Return(true, nil)

	actual := f.ica.CheckLink(f.data.ctx, edge, from, to)
	assert.True(t, actual)
}

// Link
//

func TestLink_Remove_FailsOn_RemoveDocument(t *testing.T) {
	f := newIcaFixture(t)

	edge := &driver_mocks.MockCollection{}
	from := graph.NewBlankAccountDocument("uuid")
	to := graph.NewBlankAccountDocument("uuid")

	edge.EXPECT().RemoveDocument(f.data.ctx, "uuid-uuid").
		Return(driver.DocumentMeta{}, assert.AnError)

	err := f.ica.Link(f.data.ctx, edge, from, to, access.Level_NONE, access.Role_UNSET)
	assert.Error(t, err)
	assert.EqualError(t, err, assert.AnError.Error())
}

func TestLink_Update_Success(t *testing.T) {
	f := newIcaFixture(t)

	edge := &driver_mocks.MockCollection{}
	from := graph.NewBlankAccountDocument("uuid")
	to := graph.NewBlankAccountDocument("uuid")

	edge.EXPECT().UpdateDocument(f.data.ctx, "uuid-uuid", graph.Access{
		From:  from.ID(),
		To:    to.ID(),
		Level: access.Level_READ,
		Role:  access.Role_UNSET,
		DocumentMeta: driver.DocumentMeta{
			Key: "uuid-uuid",
		},
	}).
		Return(driver.DocumentMeta{}, nil)

	err := f.ica.Link(f.data.ctx, edge, from, to, access.Level_READ, access.Role_UNSET)
	assert.NoError(t, err)
}

func TestLink_FailsOn_Create(t *testing.T) {
	f := newIcaFixture(t)

	edge := &driver_mocks.MockCollection{}
	from := graph.NewBlankAccountDocument("uuid")
	to := graph.NewBlankAccountDocument("uuid")

	edge.EXPECT().UpdateDocument(f.data.ctx, "uuid-uuid", graph.Access{
		From:  from.ID(),
		To:    to.ID(),
		Level: access.Level_READ,
		Role:  access.Role_UNSET,
		DocumentMeta: driver.DocumentMeta{
			Key: "uuid-uuid",
		},
	}).
		Return(driver.DocumentMeta{}, assert.AnError)

	edge.EXPECT().CreateDocument(f.data.ctx, graph.Access{
		From:  from.ID(),
		To:    to.ID(),
		Level: access.Level_READ,
		Role:  access.Role_UNSET,
		DocumentMeta: driver.DocumentMeta{
			Key: "uuid-uuid",
		},
	}).Return(driver.DocumentMeta{}, assert.AnError)

	err := f.ica.Link(f.data.ctx, edge, from, to, access.Level_READ, access.Role_UNSET)
	assert.Error(t, err)
	assert.EqualError(t, err, assert.AnError.Error())
}

// Move
//

// func TestMove_FailsOn_AccessLevelAndGet_Obj(t *testing.T) {
// 	f := newIcaFixture(t)

// 	obj := graph.NewBlankAccountDocument("uuid")
// 	edge := &driver_mocks.MockCollection{}

// 	f.ica.Move(
// 		f.data.ctx, f.mocks.ic, obj, edge, "ns-uuid",
// 	)
// }

// AccessLevelAndGet
//

func TestAccessLevelAndGet_FailsOn_Query(t *testing.T) {
	f := newIcaFixture(t)

	acc := graph.NewBlankAccountDocument("uuid")
	node := graph.NewBlankAccountDocument("uuid")

	cur := driver_mocks.NewMockCursor(t)

	f.mocks.db.EXPECT().
		Query(
			f.data.ctx, graph.GetWithAccessLevelRoleAndNS,
			map[string]interface{}{
				"account":     acc.ID(),
				"node":        node.ID(),
				"permissions": schema.PERMISSIONS_GRAPH.Name,
			},
		).
		Return(cur, assert.AnError)

	err := f.ica.AccessLevelAndGet(f.data.ctx, acc, node)
	assert.Error(t, err)
	assert.EqualError(t, err, assert.AnError.Error())
}

func TestAccessLevelAndGet_FailsOn_ReadDocument(t *testing.T) {
	f := newIcaFixture(t)

	acc := graph.NewBlankAccountDocument("uuid")
	node := graph.NewBlankAccountDocument("uuid")

	cur := driver_mocks.NewMockCursor(t)
	cur.EXPECT().Close().Return(nil)

	cur.EXPECT().ReadDocument(f.data.ctx, mock.Anything).Return(driver.DocumentMeta{}, assert.AnError)

	f.mocks.db.EXPECT().
		Query(
			f.data.ctx, graph.GetWithAccessLevelRoleAndNS,
			map[string]interface{}{
				"account":     acc.ID(),
				"node":        node.ID(),
				"permissions": schema.PERMISSIONS_GRAPH.Name,
			},
		).
		Return(cur, nil)

	err := f.ica.AccessLevelAndGet(f.data.ctx, acc, node)
	assert.Error(t, err)
	assert.EqualError(t, err, assert.AnError.Error())
}
