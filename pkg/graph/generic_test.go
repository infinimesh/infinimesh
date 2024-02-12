package graph_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/arangodb/go-driver"
	"github.com/google/uuid"
	driver_mocks "github.com/infinimesh/infinimesh/mocks/github.com/arangodb/go-driver"
	"github.com/infinimesh/infinimesh/pkg/graph"
	"github.com/infinimesh/proto/node/access"
	devpb "github.com/infinimesh/proto/node/devices"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type serviceFixture[T graph.InfinimeshProtobufEntity] struct {
	repo  graph.InfinimeshGenericActionsRepo[T]
	mocks struct {
		db     *driver_mocks.MockDatabase
		cursor *driver_mocks.MockCursor
		log    *zap.Logger
	}
	data struct {
		ctx    context.Context
		device *devpb.Device
	}
}

func newServiceFixture[T graph.InfinimeshProtobufEntity](t *testing.T) *serviceFixture[T] {
	f := &serviceFixture[T]{}
	f.mocks.db = driver_mocks.NewMockDatabase(t)
	f.mocks.cursor = driver_mocks.NewMockCursor(t)
	f.mocks.log = zap.NewExample()
	f.repo = graph.NewGenericRepo[T](f.mocks.db)

	f.data.ctx = context.TODO()
	f.data.device = &devpb.Device{
		Uuid:    uuid.New().String(),
		Title:   "test",
		Enabled: true,
		Access: &access.Access{
			Level: access.Level_ADMIN,
			Role:  access.Role_OWNER,
		},
	}

	return f
}

func TestListOwned_Devices_Success(t *testing.T) {
	f := newServiceFixture[*devpb.Device](t)

	f.mocks.db.EXPECT().Query(
		f.data.ctx, mock.Anything,
		map[string]interface{}{
			"depth":             10,
			"from":              driver.DocumentID("Accounts/infinimesh"),
			"permissions_graph": "Permissions",
			"@kind":             "Devices",
			"offset":            int64(0),
			"limit":             int64(0),
		},
	).Return(f.mocks.cursor, nil)

	f.mocks.cursor.EXPECT().ReadDocument(f.data.ctx, mock.MatchedBy(func(r *graph.ListQueryResult[*devpb.Device]) bool {
		t.Logf("ReadDocument(ctx, %v)", r)
		err := json.Unmarshal([]byte(`{
			"count": 1,
			"result": [{
				"uuid": "123",
				"title": "test",
				"enabled": true,
				"access": {
					"level": 3,
					"role": 1
				}
			}]
		}`), r)
		return assert.NoError(t, err)
	})).Return(driver.DocumentMeta{}, nil)
	f.mocks.cursor.EXPECT().Close().Return(nil)

	res, err := f.repo.ListQuery(f.data.ctx, f.mocks.log, graph.NewBlankAccountDocument("infinimesh"))
	assert.NoError(t, err)

	assert.Equal(t, 1, res.Count)
	assert.Len(t, res.Result, 1)
	assert.Equal(t, "123", res.Result[0].Uuid)
	assert.Equal(t, "test", res.Result[0].Title)
	assert.True(t, res.Result[0].Enabled)
	assert.Equal(t, access.Level_ADMIN, res.Result[0].Access.Level)
	assert.Equal(t, access.Role_OWNER, res.Result[0].Access.Role)
}
