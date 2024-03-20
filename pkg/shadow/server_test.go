package shadow_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	redis_mocks "github.com/infinimesh/infinimesh/mocks/github.com/go-redis/redis/v8"
	pubsub_mocks "github.com/infinimesh/infinimesh/mocks/github.com/infinimesh/infinimesh/pkg/pubsub"
	shadow_mocks "github.com/infinimesh/infinimesh/mocks/github.com/infinimesh/proto/shadow"
	"github.com/infinimesh/infinimesh/pkg/shadow"
	pb "github.com/infinimesh/proto/shadow"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"google.golang.org/protobuf/types/known/structpb"
)

type shadowServiceServerFixture struct {
	service *shadow.ShadowServiceServer
	mocks   struct {
		ps  *pubsub_mocks.MockPubSub
		rdb *redis_mocks.MockCmdable

		srv *shadow_mocks.MockShadowService_StreamShadowServer

		log      *zap.Logger
		observer *observer.ObservedLogs
	}
	data struct {
		ctx               context.Context
		uuid              string
		sample_state      string
		sample_connection string
	}
}

func newShadowServiceServerFixture(t *testing.T, args ...bool) *shadowServiceServerFixture {
	t.Parallel()
	f := &shadowServiceServerFixture{}

	core, observer := observer.New(zap.DebugLevel)
	f.mocks.log = zap.New(core)
	f.mocks.observer = observer

	f.mocks.ps = pubsub_mocks.NewMockPubSub(t)
	f.mocks.rdb = redis_mocks.NewMockCmdable(t)
	f.mocks.srv = shadow_mocks.NewMockShadowService_StreamShadowServer(t)

	f.service = shadow.NewShadowServiceServer(
		f.mocks.log, f.mocks.rdb, f.mocks.ps,
	)

	f.data.ctx = context.Background()
	f.data.uuid = uuid.New().String()
	f.data.sample_state = `{
		"data": {
		  "diff": 2
		},
		"timestamp": {
		  "nanos": 369620756,
		  "seconds": 1687185838
		}
	  }`
	f.data.sample_connection = `{
		"timestamp": {
			"nanos": 369620756,
			"seconds": 1687185838
		},
		"connected": true,
		"client_id": "device1"
	  }`

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
		f.data.sample_state,
		f.data.sample_state,
		f.data.sample_connection,
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

// Patch

func TestPatch_FailsOn_NoDevice(t *testing.T) {
	f := newShadowServiceServerFixture(t)

	_, err := f.service.Patch(f.data.ctx, &pb.Shadow{})

	assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = no device specified")
}

func TestPatch_Success(t *testing.T) {
	f := newShadowServiceServerFixture(t)

	request := &pb.Shadow{
		Device: uuid.New().String(),
		Reported: &pb.State{
			Data: &structpb.Struct{
				Fields: map[string]*structpb.Value{
					"diff": structpb.NewNumberValue(2),
				},
			},
		},
		Desired: &pb.State{
			Data: &structpb.Struct{
				Fields: map[string]*structpb.Value{
					"diff": structpb.NewNumberValue(2),
				},
			},
		},
	}

	f.mocks.ps.EXPECT().TryPub(request, "mqtt.incoming", "mqtt.outgoing").
		Return()

	resp, err := f.service.Patch(f.data.ctx, request)

	assert.NoError(t, err)
	f.mocks.ps.AssertNumberOfCalls(t, "TryPub", 1)

	assert.Equal(t, request.Device, resp.Device)

	assert.Equal(t, request.Reported.Data, resp.Reported.Data)
	assert.NotNil(t, resp.Reported.Timestamp)

	assert.Equal(t, request.Desired.Data, resp.Desired.Data)
	assert.NotNil(t, resp.Desired.Timestamp)
}

// Remove

func TestRemove_FailsOn_NoDevice(t *testing.T) {
	f := newShadowServiceServerFixture(t)

	_, err := f.service.Remove(f.data.ctx, &pb.RemoveRequest{})

	assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = no device specified")
}

func TestRemove_FailsOn_NoKey(t *testing.T) {
	f := newShadowServiceServerFixture(t)

	_, err := f.service.Remove(f.data.ctx, &pb.RemoveRequest{
		Device: uuid.New().String(),
	})

	assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = key not specified")
}

func TestRemove_FailsOn_RedisGet(t *testing.T) {
	f := newShadowServiceServerFixture(t)

	f.mocks.rdb.EXPECT().Get(f.data.ctx, "device1:reported").Return(redis.NewStringResult("", assert.AnError))

	_, err := f.service.Remove(f.data.ctx, &pb.RemoveRequest{
		Device: "device1",
		Key:    "diff",
	})

	assert.EqualError(t, err, "rpc error: code = Internal desc = failed to get Shadow")
}

func TestRemove_FailsOn_Unmarshal(t *testing.T) {
	f := newShadowServiceServerFixture(t)

	f.mocks.rdb.EXPECT().Get(f.data.ctx, "device1:reported").Return(redis.NewStringResult("invalid", nil))

	_, err := f.service.Remove(f.data.ctx, &pb.RemoveRequest{
		Device: "device1",
		Key:    "diff",
	})

	assert.EqualError(t, err, "rpc error: code = Internal desc = cannot Unmarshal state")
}

func TestRemove_Reported_Success(t *testing.T) {
	f := newShadowServiceServerFixture(t)

	f.mocks.rdb.EXPECT().Get(f.data.ctx, "device1:reported").Return(redis.NewStringResult(f.data.sample_state, nil))
	f.mocks.rdb.EXPECT().Set(f.data.ctx, "device1:reported", mock.MatchedBy(func(s string) bool {
		state := pb.State{}
		err := json.Unmarshal([]byte(s), &state)
		if err != nil {
			t.Logf("Error unmarshalling state: %s", err)
			return false
		}

		if len(state.Data.Fields) != 0 {
			t.Logf("Expected state to be empty, got: %s", s)
			return false
		}

		return true
	}), time.Duration(0)).Return(redis.NewStatusResult("", nil))

	res, err := f.service.Remove(f.data.ctx, &pb.RemoveRequest{
		Device: "device1",
		Key:    "diff",
	})

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotNil(t, res.Reported)
	assert.NotNil(t, res.Reported.Data)
	assert.Len(t, res.Reported.Data.Fields, 0)
}

func TestRemove_Desired_Success(t *testing.T) {
	f := newShadowServiceServerFixture(t)

	f.mocks.rdb.EXPECT().Get(f.data.ctx, "device1:desired").Return(redis.NewStringResult(f.data.sample_state, nil))
	f.mocks.rdb.EXPECT().Set(f.data.ctx, "device1:desired", mock.MatchedBy(func(s string) bool {
		state := pb.State{}
		err := json.Unmarshal([]byte(s), &state)
		if err != nil {
			t.Logf("Error unmarshalling state: %s", err)
			return false
		}

		if len(state.Data.Fields) != 0 {
			t.Logf("Expected state to be empty, got: %s", s)
			return false
		}

		return true
	}), time.Duration(0)).Return(redis.NewStatusResult("", nil))

	res, err := f.service.Remove(f.data.ctx, &pb.RemoveRequest{
		Device:   "device1",
		Key:      "diff",
		StateKey: pb.StateKey_DESIRED,
	})

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotNil(t, res.Desired)
	assert.NotNil(t, res.Desired.Data)
	assert.Len(t, res.Desired.Data.Fields, 0)
}

// StreamShadow

func TestStreamShadow_FailsOn_NoDevices(t *testing.T) {
	f := newShadowServiceServerFixture(t)

	err := f.service.StreamShadow(&pb.StreamShadowRequest{}, f.mocks.srv)

	assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = no devices specified")
}

func TestStreamShadow_Success_ServerClosed(t *testing.T) {
	f := newShadowServiceServerFixture(t)

	f.mocks.srv.EXPECT().Context().Return(f.data.ctx)
	f.mocks.srv.EXPECT().Send(mock.Anything).Return(assert.AnError)

	mget_result := redis.NewSliceCmd(f.data.ctx)
	mget_result.SetVal([]interface{}{f.data.sample_state, f.data.sample_state, f.data.sample_connection})
	f.mocks.rdb.EXPECT().MGet(
		f.data.ctx, shadow.Key(f.data.uuid, pb.StateKey_REPORTED),
		shadow.Key(f.data.uuid, pb.StateKey_DESIRED),
		shadow.Key(f.data.uuid, pb.StateKey_CONNECTION),
	).Return(mget_result)

	f.mocks.ps.EXPECT().AddSub(mock.MatchedBy(func(ch chan interface{}) bool {
		ch <- &pb.Shadow{
			Device: f.data.uuid,
		}
		return true
	}), "mqtt.incoming", "mqtt.outgoing").Return()
	f.mocks.ps.EXPECT().Unsub(mock.MatchedBy(func(ch chan interface{}) bool {
		close(ch)
		return true
	})).Return()

	done := make(chan bool, 1)
	go func(done chan bool) {
		err := f.service.StreamShadow(&pb.StreamShadowRequest{
			Devices: []string{f.data.uuid},
			Sync:    true,
		}, f.mocks.srv)
		assert.NoError(t, err)
		done <- true
	}(done)

	timer := time.After(time.Second * 5)
	ticker := time.Tick(10 * time.Millisecond)

timer_loop:
	for {
		select {
		case <-timer:
			t.Fatal("Server did not close")
			break timer_loop
		case <-done:
			break timer_loop
		case <-ticker:
			time.Sleep(time.Millisecond * 10)
		}
	}

	f.mocks.srv.AssertNumberOfCalls(t, "Send", 2)
}

// Store

func TestStore_FailsOn_Marshal(t *testing.T) {
	f := newShadowServiceServerFixture(t)

	key := f.data.uuid + ":reported"
	res, ok := f.service.Store(zap.NewExample(), f.data.uuid, pb.StateKey_REPORTED, make(chan int)) // ensure Marshal fails

	assert.Equal(t, false, ok)
	assert.Equal(t, key, res)
}

func TestStore_FailsOn_RedisSet(t *testing.T) {
	f := newShadowServiceServerFixture(t)

	key := f.data.uuid + ":reported"
	f.mocks.rdb.EXPECT().Set(
		f.data.ctx, key, "{}", time.Duration(0),
	).Return(redis.NewStatusResult("", assert.AnError))

	res, ok := f.service.Store(zap.NewExample(), f.data.uuid, pb.StateKey_REPORTED, &pb.State{}) // ensure Marshal fails

	assert.Equal(t, false, ok)
	assert.Equal(t, key, res)
	f.mocks.rdb.AssertNumberOfCalls(t, "Set", 1)
}

func TestStore_Success(t *testing.T) {
	f := newShadowServiceServerFixture(t)

	key := f.data.uuid + ":reported"
	f.mocks.rdb.EXPECT().Set(
		f.data.ctx, key, "{}", time.Duration(0),
	).Return(redis.NewStatusResult("", nil))

	res, ok := f.service.Store(zap.NewExample(), f.data.uuid, pb.StateKey_REPORTED, &pb.State{}) // ensure Marshal fails

	assert.Equal(t, true, ok)
	assert.Equal(t, key, res)
	f.mocks.rdb.AssertNumberOfCalls(t, "Set", 1)
}

// StoreConnectionState

func TestStoreConnectionState_FailsOn_Store(t *testing.T) {
	f := newShadowServiceServerFixture(t)

	key := f.data.uuid + ":connection"
	f.mocks.rdb.EXPECT().Set(
		f.data.ctx, key, "{}", time.Duration(0),
	).Return(redis.NewStatusResult("", assert.AnError))

	f.service.StoreConnectionState(zap.NewExample(), f.data.uuid, &pb.ConnectionState{})

	f.mocks.rdb.AssertNumberOfCalls(t, "Set", 1)
}

func TestStoreConnectionState_Success(t *testing.T) {
	f := newShadowServiceServerFixture(t)

	key := f.data.uuid + ":connection"
	f.mocks.rdb.EXPECT().Set(
		f.data.ctx, key, "{}", time.Duration(0),
	).Return(redis.NewStatusResult("", nil))
	f.mocks.rdb.EXPECT().Expire(
		f.data.ctx, key, time.Hour*24,
	).Return(redis.NewBoolResult(false, assert.AnError))

	f.service.StoreConnectionState(zap.NewExample(), f.data.uuid, &pb.ConnectionState{})

	f.mocks.rdb.AssertNumberOfCalls(t, "Set", 1)
	f.mocks.rdb.AssertNumberOfCalls(t, "Expire", 1)
}

// MergeAndStore

func TestMergeAndStore_FailsOn_RedisGetAndSet(t *testing.T) {
	f := newShadowServiceServerFixture(t)

	key := f.data.uuid + ":reported"
	f.mocks.rdb.EXPECT().Get(
		f.data.ctx, key,
	).Return(redis.NewStringResult("", assert.AnError))

	f.mocks.rdb.EXPECT().Set(
		f.data.ctx, key, "{}", time.Duration(0),
	).Return(redis.NewStatusResult("", assert.AnError))

	f.service.MergeAndStore(zap.NewExample(), f.data.uuid, pb.StateKey_REPORTED, &pb.State{})

	f.mocks.rdb.AssertNumberOfCalls(t, "Get", 1)
	f.mocks.rdb.AssertNumberOfCalls(t, "Set", 1)
}

func TestMergeAndStore_FailsOn_MergeOldIsInvalid(t *testing.T) {
	f := newShadowServiceServerFixture(t)

	key := f.data.uuid + ":reported"
	f.mocks.rdb.EXPECT().Get(
		f.data.ctx, key,
	).Return(redis.NewStringResult("invalid", nil))

	f.service.MergeAndStore(zap.NewExample(), f.data.uuid, pb.StateKey_REPORTED, &pb.State{})

	f.mocks.rdb.AssertNumberOfCalls(t, "Get", 1)
}

func TestMergeAndStore_SuccessWithMerge(t *testing.T) {
	f := newShadowServiceServerFixture(t)

	key := f.data.uuid + ":reported"
	f.mocks.rdb.EXPECT().Get(
		f.data.ctx, key,
	).Return(redis.NewStringResult(f.data.sample_state, nil))

	f.mocks.rdb.EXPECT().Set(
		f.data.ctx, key, mock.MatchedBy(func(s string) bool {
			state := pb.State{}
			err := json.Unmarshal([]byte(s), &state)
			if err != nil {
				t.Errorf("Error unmarshalling state: %s", err)
				return false
			}

			assert.Len(t, state.Data.Fields, 2)
			assert.Equal(t, float64(2), state.Data.Fields["diff"].GetNumberValue())
			assert.Equal(t, "bar", state.Data.Fields["foo"].GetStringValue())

			return true
		}), time.Duration(0),
	).Return(redis.NewStatusResult("", nil))

	f.service.MergeAndStore(zap.NewExample(), f.data.uuid, pb.StateKey_REPORTED, &pb.State{
		Data: &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"foo": structpb.NewStringValue("bar"),
			},
		},
	})

	f.mocks.rdb.AssertNumberOfCalls(t, "Get", 1)
	f.mocks.rdb.AssertNumberOfCalls(t, "Set", 1)
}

func TestMergeAndStore_SuccessWithMergeOldEmpty(t *testing.T) {
	f := newShadowServiceServerFixture(t)

	key := f.data.uuid + ":reported"
	f.mocks.rdb.EXPECT().Get(
		f.data.ctx, key,
	).Return(redis.NewStringResult("", nil))

	f.mocks.rdb.EXPECT().Set(
		f.data.ctx, key, mock.MatchedBy(func(s string) bool {
			state := pb.State{}
			err := json.Unmarshal([]byte(s), &state)
			if err != nil {
				t.Errorf("Error unmarshalling state: %s", err)
				return false
			}

			assert.Len(t, state.Data.Fields, 1)
			assert.Equal(t, "bar", state.Data.Fields["foo"].GetStringValue())

			return true
		}), time.Duration(0),
	).Return(redis.NewStatusResult("", nil))

	f.service.MergeAndStore(zap.NewExample(), f.data.uuid, pb.StateKey_REPORTED, &pb.State{
		Data: &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"foo": structpb.NewStringValue("bar"),
			},
		},
	})

	f.mocks.rdb.AssertNumberOfCalls(t, "Get", 1)
	f.mocks.rdb.AssertNumberOfCalls(t, "Set", 1)
}

// MergeJSON

type mergeJsonCase struct {
	msg         string
	old         string
	new         string
	expected    string
	expectedErr bool

	intransitive bool
}

func TestMergeJSON_Cases(t *testing.T) {
	cases := []mergeJsonCase{
		{
			msg:      "Message is empty",
			old:      "",
			new:      `{ "foo": "bar" }`,
			expected: `{ "foo": "bar" }`,
		},
		{
			msg:         "Message is invalid",
			old:         "invalid",
			new:         `{ "foo": "bar" }`,
			expected:    `{ "foo": "bar" }`,
			expectedErr: true,
		},
		{
			msg:      "Message is empty object",
			old:      `{}`,
			new:      `{ "foo": "bar" }`,
			expected: `{ "foo": "bar" }`,
		},
		{
			msg:          "Message (old) is empty array",
			old:          `[]`,
			new:          `{ "foo": "bar" }`,
			expected:     `{ "foo": "bar" }`,
			intransitive: true,
		},
		{
			msg:          "Message is valid",
			old:          `{ "foo": "bar" }`,
			new:          `{ "foo": "baz" }`,
			expected:     `{ "foo": "baz" }`,
			intransitive: true,
		},
	}

	test_func := func(c mergeJsonCase) func(t *testing.T) {
		return func(t *testing.T) {
			merged, err := shadow.MergeJSON([]byte(c.old), []byte(c.new))
			if c.expectedErr {
				assert.Errorf(t, err, "expected error, got nil and `%s`", string(merged))
				return
			}
			assert.NoError(t, err)
			assert.JSONEq(t, c.expected, string(merged))
		}
	}

	for _, c := range cases {
		t.Run(c.msg+"(testing old message)", test_func(c))
		if c.intransitive {
			continue
		}
		c.old, c.new = c.new, c.old
		t.Run(c.msg+"(testing new message)", test_func(c))
	}
}

// Persister

func TestPersister_Success(t *testing.T) {
	f := newShadowServiceServerFixture(t)

	f.mocks.ps.EXPECT().AddSub(mock.MatchedBy(func(ch chan interface{}) bool {
		ch <- &pb.Shadow{
			Device: f.data.uuid,
		}
		go func() {
			time.Sleep(time.Millisecond * 200)
			close(ch)
		}()
		return true
	}), "mqtt.incoming", "mqtt.outgoing").Return()
	f.mocks.ps.EXPECT().Unsub(mock.Anything).Return()

	done := make(chan bool, 1)
	go func(done chan bool) {
		f.service.Persister()
		done <- true
	}(done)

	timer := time.After(time.Second * 5)
	ticker := time.Tick(10 * time.Millisecond)

timer_loop:
	for {
		select {
		case <-timer:
			t.Fatal("Server did not close")
			break timer_loop
		case <-done:
			break timer_loop
		case <-ticker:
			time.Sleep(time.Millisecond * 10)
		}
	}

	alleged_exited_logs := f.mocks.observer.FilterMessage("Exited").All()
	assert.Len(t, alleged_exited_logs, 1)
	assert.Equal(t, zap.WarnLevel, alleged_exited_logs[0].Level)

	f.mocks.ps.AssertNumberOfCalls(t, "AddSub", 1)
	f.mocks.ps.AssertNumberOfCalls(t, "Unsub", 1)
}
