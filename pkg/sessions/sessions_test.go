package sessions_test

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	redis_mocks "github.com/infinimesh/infinimesh/mocks/github.com/go-redis/redis/v8"
	"github.com/infinimesh/infinimesh/pkg/sessions"
	sess_pb "github.com/infinimesh/proto/node/sessions"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type sessionsFixture struct {
	sh sessions.SessionsHandler

	mocks struct {
		rdb *redis_mocks.MockCmdable
	}

	data struct {
		storeSession *sess_pb.Session
	}
}

func newSessionsFixture(
	t *testing.T,
	protoMarshal func(m protoreflect.ProtoMessage) ([]byte, error),
	protoUnmarshal func(b []byte, m protoreflect.ProtoMessage) error,
) (f *sessionsFixture) {

	t.Parallel()
	f = &sessionsFixture{}
	f.mocks.rdb = redis_mocks.NewMockCmdable(t)

	f.sh = sessions.NewSessionsHandler(
		f.mocks.rdb,
		protoMarshal,
		protoUnmarshal,
	)

	f.data.storeSession = &sess_pb.Session{
		Id:      "test-id",
		Expires: timestamppb.New(time.Now().Add(time.Hour)),
	}

	return f
}

// New
func TestNew_NoExp_Success(t *testing.T) {
	client := "Linux / Test / Client"

	f := newSessionsFixture(t, nil, nil)

	session := f.sh.New(0, client)

	assert.Equal(t, client, session.Client)
	assert.Nil(t, session.Expires)
	assert.NotEmpty(t, session.Id)
	assert.NotEmpty(t, session.Created.AsTime())
}

func TestNew_Success(t *testing.T) {
	// Exp is set to random int64
	exp := rand.Int63()
	client := "Linux / Test / Client"

	f := newSessionsFixture(t, nil, nil)

	session := f.sh.New(exp, client)

	assert.Equal(t, client, session.Client)
	assert.Equal(t, exp, session.Expires.AsTime().Unix())
	assert.NotEmpty(t, session.Id)
	assert.NotEmpty(t, session.Created.AsTime())
}

// Store
//

func TestStore_FailsOn_Marshal(t *testing.T) {
	f := newSessionsFixture(t,
		func(m protoreflect.ProtoMessage) ([]byte, error) {
			return nil, assert.AnError
		}, nil)

	err := f.sh.Store("account", f.data.storeSession)
	assert.Equal(t, assert.AnError, err)
}

func TestStore_FailsOn_Set(t *testing.T) {
	f := newSessionsFixture(t, func(m protoreflect.ProtoMessage) ([]byte, error) {
		return []byte("data"), nil
	}, nil)

	res := redis.NewStatusCmd(context.Background())
	res.SetErr(assert.AnError)

	f.mocks.rdb.On(
		"Set", context.Background(),
		mock.Anything, []byte("data"),
		mock.MatchedBy(func(ret time.Duration) bool {
			return ret < time.Hour
		})).
		Return(res)

	err := f.sh.Store("account", f.data.storeSession)
	assert.Equal(t, assert.AnError, err)
}

func TestStore_Success(t *testing.T) {
	f := newSessionsFixture(t, func(m protoreflect.ProtoMessage) ([]byte, error) {
		return []byte("data"), nil
	}, nil)

	res := redis.NewStatusCmd(context.Background())

	f.mocks.rdb.On(
		"Set", context.Background(),
		mock.Anything, []byte("data"),
		mock.MatchedBy(func(ret time.Duration) bool {
			return ret < time.Hour
		})).
		Return(res)

	err := f.sh.Store("account", f.data.storeSession)
	assert.Nil(t, err)
}

// Check
//

func TestCheck_FailsOn_Get(t *testing.T) {
	f := newSessionsFixture(t, nil, nil)

	res := redis.NewStringCmd(context.Background())
	res.SetErr(assert.AnError)

	f.mocks.rdb.On(
		"Get", context.Background(),
		mock.Anything).
		Return(res)

	err := f.sh.Check("account", "test-id")
	assert.Equal(t, assert.AnError, err)
}

func TestCheck_FailsOn_Unmarshal(t *testing.T) {
	f := newSessionsFixture(t, nil, func(b []byte, m protoreflect.ProtoMessage) error {
		return assert.AnError
	})

	res := redis.NewStringCmd(context.Background())
	res.SetVal("data")

	f.mocks.rdb.On(
		"Get", context.Background(),
		mock.Anything).
		Return(res)

	err := f.sh.Check("account", "test-id")
	assert.Equal(t, assert.AnError, err)
}

func TestCheck_FailsOn_SessionExpired(t *testing.T) {
	f := newSessionsFixture(t, nil, func(b []byte, m protoreflect.ProtoMessage) error {
		desc := m.ProtoReflect().Descriptor()
		field := desc.Fields().ByName("expires")

		m.ProtoReflect().Set(field, protoreflect.ValueOfMessage(
			timestamppb.New(time.Now().Add(-time.Hour)).ProtoReflect(),
		))

		return nil
	})

	res := redis.NewStringCmd(context.Background())
	res.SetVal("data")

	f.mocks.rdb.On(
		"Get", context.Background(),
		mock.Anything).
		Return(res)

	err := f.sh.Check("account", "test-id")
	assert.Equal(t, "session expired", err.Error())
}

func TestCheck_Success(t *testing.T) {
	f := newSessionsFixture(t, nil, func(b []byte, m protoreflect.ProtoMessage) error {
		desc := m.ProtoReflect().Descriptor()
		field := desc.Fields().ByName("expires")

		m.ProtoReflect().Set(field, protoreflect.ValueOfMessage(
			timestamppb.New(time.Now().Add(time.Hour)).ProtoReflect(),
		))

		return nil
	})

	res := redis.NewStringCmd(context.Background())
	res.SetVal("data")

	f.mocks.rdb.On(
		"Get", context.Background(),
		mock.Anything).
		Return(res)

	err := f.sh.Check("account", "test-id")
	assert.Nil(t, err)
}

// LogActivity
//

func TestLogActivity_FailsOn_Set(t *testing.T) {
	f := newSessionsFixture(t, nil, nil)

	res := redis.NewStatusCmd(context.Background())
	res.SetErr(assert.AnError)

	f.mocks.rdb.On(
		"Set", context.Background(),
		mock.Anything, mock.Anything, mock.Anything).
		Return(res)

	err := f.sh.LogActivity("account", "test-id", 0)
	assert.Equal(t, assert.AnError, err)
}

func TestLogActivity_Success(t *testing.T) {
	f := newSessionsFixture(t, nil, nil)

	res := redis.NewStatusCmd(context.Background())

	f.mocks.rdb.On(
		"Set", context.Background(),
		mock.Anything, mock.Anything, mock.Anything).
		Return(res)

	err := f.sh.LogActivity("account", "test-id", 0)
	assert.Nil(t, err)
}

// GetActivity
//

func TestGetActivity_FailsOn_Keys(t *testing.T) {
	f := newSessionsFixture(t, nil, nil)

	res := redis.NewStringSliceCmd(context.Background())
	res.SetErr(assert.AnError)

	f.mocks.rdb.On(
		"Keys", context.Background(),
		mock.Anything).
		Return(res)

	_, err := f.sh.GetActivity("account")
	assert.Equal(t, assert.AnError, err)
}

func TestGetActivity_FailsOn_MGet(t *testing.T) {
	f := newSessionsFixture(t, nil, nil)

	keysCmd := redis.NewStringSliceCmd(context.Background())
	keysCmd.SetVal([]string{"key"})

	res := redis.NewSliceCmd(context.Background())
	res.SetErr(assert.AnError)

	f.mocks.rdb.On(
		"Keys", context.Background(),
		mock.Anything).
		Return(keysCmd)

	f.mocks.rdb.On(
		"MGet", context.Background(),
		mock.Anything).
		Return(res)

	_, err := f.sh.GetActivity("account")
	assert.Equal(t, assert.AnError, err)
}

func TestGetActivity_FailsOn_InvalidDataType(t *testing.T) {
	f := newSessionsFixture(t, nil, nil)

	keysCmd := redis.NewStringSliceCmd(context.Background())
	keysCmd.SetVal([]string{"key"})

	res := redis.NewSliceCmd(context.Background())
	res.SetVal([]interface{}{1})

	f.mocks.rdb.On(
		"Keys", context.Background(),
		mock.Anything).
		Return(keysCmd)

	f.mocks.rdb.On(
		"MGet", context.Background(),
		mock.Anything).
		Return(res)

	_, err := f.sh.GetActivity("account")
	assert.Equal(t, "invalid data type: key", err.Error())
}

func TestGetActivity_FailsOn_InvalidDataTypeOn_Atoi(t *testing.T) {
	f := newSessionsFixture(t, nil, nil)

	keysCmd := redis.NewStringSliceCmd(context.Background())
	keysCmd.SetVal([]string{"key"})

	res := redis.NewSliceCmd(context.Background())
	res.SetVal([]interface{}{"string"})

	f.mocks.rdb.On(
		"Keys", context.Background(),
		mock.Anything).
		Return(keysCmd)

	f.mocks.rdb.On(
		"MGet", context.Background(),
		mock.Anything).
		Return(res)

	_, err := f.sh.GetActivity("account")
	assert.Equal(t, "invalid data type: key | strconv.Atoi: parsing \"string\": invalid syntax", err.Error())
}

func TestGetActivity_Success(t *testing.T) {
	f := newSessionsFixture(t, nil, nil)

	keysCmd := redis.NewStringSliceCmd(context.Background())
	keysCmd.SetVal([]string{"sessions:account:session:key"})

	res := redis.NewSliceCmd(context.Background())
	res.SetVal([]interface{}{"1"})

	f.mocks.rdb.On(
		"Keys", context.Background(),
		mock.Anything).
		Return(keysCmd)

	f.mocks.rdb.On(
		"MGet", context.Background(),
		mock.Anything).
		Return(res)

	_, err := f.sh.GetActivity("account")
	assert.Nil(t, err)
}

// Get
//

func TestGet_FailsOn_Keys(t *testing.T) {
	f := newSessionsFixture(t, nil, nil)

	res := redis.NewStringSliceCmd(context.Background())
	res.SetErr(assert.AnError)

	f.mocks.rdb.On(
		"Keys", context.Background(),
		mock.Anything).
		Return(res)

	_, err := f.sh.Get("account")
	assert.Equal(t, assert.AnError, err)
}

func TestGet_FailsOn_MGet(t *testing.T) {
	f := newSessionsFixture(t, nil, nil)

	keysCmd := redis.NewStringSliceCmd(context.Background())
	keysCmd.SetVal([]string{"key"})

	res := redis.NewSliceCmd(context.Background())
	res.SetErr(assert.AnError)

	f.mocks.rdb.On(
		"Keys", context.Background(),
		mock.Anything).
		Return(keysCmd)

	f.mocks.rdb.On(
		"MGet", context.Background(),
		mock.Anything).
		Return(res)

	_, err := f.sh.Get("account")
	assert.Equal(t, assert.AnError, err)
}

func TestGet_FailsOn_InvalidDataType(t *testing.T) {
	f := newSessionsFixture(t, nil, nil)

	keysCmd := redis.NewStringSliceCmd(context.Background())
	keysCmd.SetVal([]string{"key"})

	res := redis.NewSliceCmd(context.Background())
	res.SetVal([]interface{}{1})

	f.mocks.rdb.On(
		"Keys", context.Background(),
		mock.Anything).
		Return(keysCmd)

	f.mocks.rdb.On(
		"MGet", context.Background(),
		mock.Anything).
		Return(res)

	_, err := f.sh.Get("account")
	assert.Equal(t, "invalid data type: key", err.Error())
}

func TestGet_FailsOn_Unmarshal(t *testing.T) {
	f := newSessionsFixture(t, nil, func(b []byte, m protoreflect.ProtoMessage) error {
		return assert.AnError
	})

	keysCmd := redis.NewStringSliceCmd(context.Background())
	keysCmd.SetVal([]string{"key"})

	res := redis.NewSliceCmd(context.Background())
	res.SetVal([]interface{}{"string"})

	f.mocks.rdb.On(
		"Keys", context.Background(),
		mock.Anything).
		Return(keysCmd)

	f.mocks.rdb.On(
		"MGet", context.Background(),
		mock.Anything).
		Return(res)

	_, err := f.sh.Get("account")
	assert.Equal(t, assert.AnError, err)
}

func TestGet_Success(t *testing.T) {
	f := newSessionsFixture(t, nil, func(b []byte, m protoreflect.ProtoMessage) error {
		return nil
	})

	keysCmd := redis.NewStringSliceCmd(context.Background())
	keysCmd.SetVal([]string{"key"})

	res := redis.NewSliceCmd(context.Background())
	res.SetVal([]interface{}{"string"})

	f.mocks.rdb.On(
		"Keys", context.Background(),
		mock.Anything).
		Return(keysCmd)

	f.mocks.rdb.On(
		"MGet", context.Background(),
		mock.Anything).
		Return(res)

	_, err := f.sh.Get("account")
	assert.Nil(t, err)
}

// Revoke
//

func TestRevoke_FailsOn_Del(t *testing.T) {
	f := newSessionsFixture(t, nil, nil)

	res := redis.NewIntCmd(context.Background())
	res.SetErr(assert.AnError)

	f.mocks.rdb.On(
		"Del", context.Background(),
		mock.Anything).
		Return(res)

	err := f.sh.Revoke("account", "test-id")
	assert.Equal(t, assert.AnError, err)
}

func TestRevoke_Success(t *testing.T) {
	f := newSessionsFixture(t, nil, nil)

	res := redis.NewIntCmd(context.Background())

	f.mocks.rdb.On(
		"Del", context.Background(),
		mock.Anything).
		Return(res)

	err := f.sh.Revoke("account", "test-id")
	assert.Nil(t, err)
}
