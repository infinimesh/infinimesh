package shadow_test

import (
	"context"
	"testing"

	"github.com/go-redis/redis/v8"
	redis_mocks "github.com/infinimesh/infinimesh/mocks/github.com/go-redis/redis/v8"
	pubsub_mocks "github.com/infinimesh/infinimesh/mocks/github.com/infinimesh/infinimesh/pkg/pubsub"
	"github.com/infinimesh/infinimesh/pkg/shadow"
	pb "github.com/infinimesh/proto/shadow"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type shadowServiceServerFixture struct {
	service *shadow.ShadowServiceServer
	mocks   struct {
		ps  *pubsub_mocks.MockPubSub
		rdb *redis_mocks.MockCmdable
	}
	data struct {
		ctx context.Context
	}
}

func newShadowServiceServerFixture(t *testing.T) *shadowServiceServerFixture {
	f := &shadowServiceServerFixture{}

	f.mocks.ps = pubsub_mocks.NewMockPubSub(t)
	f.mocks.rdb = redis_mocks.NewMockCmdable(t)

	f.service = shadow.NewShadowServiceServer(
		zap.NewExample(), f.mocks.rdb, f.mocks.ps,
	)

	f.data.ctx = context.Background()

	return f
}

// Get

func TestGet_FailsOn_NoDevices(t *testing.T) {
	f := newShadowServiceServerFixture(t)

	_, err := f.service.Get(f.data.ctx, &pb.GetRequest{
		Pool: []string{},
	})

	assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = no devices specified")
}

func TestGet_FailsOn_RedisError(t *testing.T) {
	f := newShadowServiceServerFixture(t)

	mget_cmd := redis.NewSliceCmd(f.data.ctx)
	mget_cmd.SetErr(assert.AnError)
	f.mocks.rdb.
		EXPECT().
		MGet(f.data.ctx, "device1:reported", "device1:desired", "device1:connection").
		Return(mget_cmd)

	_, err := f.service.Get(f.data.ctx, &pb.GetRequest{
		Pool: []string{"device1"},
	})

	assert.EqualError(t, err, "rpc error: code = Internal desc = failed to get Shadows")
}

func TestGet_Success(t *testing.T) {
	f := newShadowServiceServerFixture(t)

	mget_cmd := redis.NewSliceCmd(f.data.ctx)
	mget_cmd.SetVal([]interface{}{
		`{
			"data": {
			  "diff": 2
			},
			"timestamp": {
			  "nanos": 369620756,
			  "seconds": 1687185838
			}
		  }`, `{
			"data": {
			  "diff": 2
			},
			"timestamp": {
			  "nanos": 369620756,
			  "seconds": 1687185838
			}
		  }`, `{
			"timestamp": {
				"nanos": 369620756,
				"seconds": 1687185838
			},
			"connected": true,
			"client_id": "device1"
		  }`,
	})

	f.mocks.rdb.
		EXPECT().
		MGet(f.data.ctx, "device1:reported", "device1:desired", "device1:connection").
		Return(mget_cmd)

	resp, err := f.service.Get(f.data.ctx, &pb.GetRequest{
		Pool: []string{"device1"},
	})

	assert.NoError(t, err)

	assert.Len(t, resp.Shadows, 1)

	shadow := resp.Shadows[0]
	assert.Equal(t, "device1", shadow.Device)

	reported := shadow.Reported
	assert.NotNil(t, reported)
	assert.NotNil(t, reported.Data)
	assert.Len(t, reported.Data.Fields, 1)
	assert.Contains(t, reported.Data.Fields, "diff")
	assert.Equal(t, float64(2), reported.Data.Fields["diff"].GetNumberValue())

	desired := shadow.Desired
	assert.NotNil(t, desired)
	assert.NotNil(t, desired.Data)
	assert.Len(t, desired.Data.Fields, 1)
	assert.Contains(t, desired.Data.Fields, "diff")
	assert.Equal(t, float64(2), desired.Data.Fields["diff"].GetNumberValue())

	connection := shadow.Connection
	assert.NotNil(t, connection)
	assert.Equal(t, true, connection.Connected)
	assert.Equal(t, "device1", connection.ClientId)
}
