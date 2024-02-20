package auth_test

import (
	"context"
	"errors"
	"runtime"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/golang-jwt/jwt/v4"
	sessions_mocks "github.com/infinimesh/infinimesh/mocks/github.com/infinimesh/infinimesh/pkg/sessions"
	auth_mocks "github.com/infinimesh/infinimesh/mocks/github.com/infinimesh/infinimesh/pkg/shared/auth"
	infinimesh "github.com/infinimesh/infinimesh/pkg/shared"
	"github.com/infinimesh/infinimesh/pkg/shared/auth"
	"github.com/infinimesh/proto/node"
	"github.com/infinimesh/proto/node/access"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"google.golang.org/grpc/metadata"
)

type interceptorFixture struct {
	interceptor auth.AuthInterceptor

	mocks struct {
		jwth *auth_mocks.MockJWTHandler

		log      *zap.Logger
		observer *observer.ObservedLogs

		sessions *sessions_mocks.MockSessionsHandler
	}
}

func newInterceptorFixture(t *testing.T) interceptorFixture {
	f := interceptorFixture{}

	core, observer := observer.New(zap.DebugLevel)
	f.mocks.log = zap.New(core)
	f.mocks.observer = observer

	f.mocks.jwth = auth_mocks.NewMockJWTHandler(t)
	f.mocks.sessions = sessions_mocks.NewMockSessionsHandler(t)

	f.interceptor = auth.NewAuthInterceptor(f.mocks.log, nil, f.mocks.jwth, nil)
	f.interceptor.SetSessionsHandler(f.mocks.sessions)
	return f
}

func CallUnary(i auth.AuthInterceptor, ctx context.Context, req connect.AnyRequest) (context.Context, connect.AnyRequest, connect.AnyResponse, error) {
	var (
		res_ctx context.Context
		res_ar  connect.AnyRequest
	)
	res, err := i.WrapUnary(func(ctx context.Context, ar connect.AnyRequest) (connect.AnyResponse, error) {
		res_ctx = ctx
		res_ar = ar
		return nil, nil
	})(ctx, req)

	return res_ctx, res_ar, res, err
}

func TestConnectStandardAuthMiddleware_FailsOn_NoToken(t *testing.T) {
	f := newInterceptorFixture(t)

	f.mocks.jwth.EXPECT().Parse("", mock.Anything).Return(nil, errors.New("token contains an invalid number of segments"))

	_, _, res, err := CallUnary(f.interceptor, context.Background(), connect.NewRequest(
		&node.EmptyMessage{},
	))

	assert.Nil(t, res)
	assert.EqualError(t, err, "rpc error: code = Unauthenticated desc = Invalid token format")
}

func TestConnectStandardAuthMiddleware_FailsOn_InvalidToken(t *testing.T) {
	f := newInterceptorFixture(t)

	f.mocks.jwth.EXPECT().Parse("invalid", mock.Anything).Return(nil, errors.New("token contains an invalid number of segments"))

	req := connect.NewRequest(
		&node.EmptyMessage{},
	)
	req.Header().Add("authorization", "Bearer invalid")

	_, _, res, err := CallUnary(f.interceptor, context.Background(), req)

	assert.Nil(t, res)
	assert.EqualError(t, err, "rpc error: code = Unauthenticated desc = Invalid token format")
}

func TestConnectStandardAuthMiddleware_FailsOn_NonValidToken(t *testing.T) {
	f := newInterceptorFixture(t)

	f.mocks.jwth.EXPECT().Parse(mock.Anything, mock.Anything).Return(
		&jwt.Token{
			Claims: jwt.MapClaims{},
			Valid:  false,
		}, nil,
	)

	_, _, res, err := CallUnary(f.interceptor, context.Background(), connect.NewRequest(
		&node.EmptyMessage{},
	))

	assert.Nil(t, res)
	assert.EqualError(t, err, "rpc error: code = Unauthenticated desc = Invalid token format")
}

func TestConnectStandardAuthMiddleware_FailsOn_ClaimsWrongType(t *testing.T) {
	f := newInterceptorFixture(t)

	f.mocks.jwth.EXPECT().Parse(mock.Anything, mock.Anything).Return(
		&jwt.Token{
			Claims: jwt.RegisteredClaims{},
			Valid:  true,
		}, nil,
	)

	_, _, res, err := CallUnary(f.interceptor, context.Background(), connect.NewRequest(
		&node.EmptyMessage{},
	))

	assert.Nil(t, res)
	assert.EqualError(t, err, "rpc error: code = Unauthenticated desc = Invalid token format")
}

func TestConnectStandardAuthMiddleware_FailsOn_NoRequestor(t *testing.T) {
	f := newInterceptorFixture(t)

	f.mocks.jwth.EXPECT().Parse(mock.Anything, mock.Anything).Return(
		&jwt.Token{
			Claims: jwt.MapClaims{},
			Valid:  true,
		}, nil,
	)

	_, _, res, err := CallUnary(f.interceptor, context.Background(), connect.NewRequest(
		&node.EmptyMessage{},
	))

	assert.Nil(t, res)
	assert.EqualError(t, err, "rpc error: code = Unauthenticated desc = Invalid token format: no requestor ID")
}

func TestConnectStandardAuthMiddleware_FailsOn_RequestorWrongType(t *testing.T) {
	f := newInterceptorFixture(t)

	f.mocks.jwth.EXPECT().Parse(mock.Anything, mock.Anything).Return(
		&jwt.Token{
			Claims: jwt.MapClaims{
				infinimesh.INFINIMESH_ACCOUNT_CLAIM: 666,
			},
			Valid: true,
		}, nil,
	)

	_, _, res, err := CallUnary(f.interceptor, context.Background(), connect.NewRequest(
		&node.EmptyMessage{},
	))

	assert.Nil(t, res)
	assert.EqualError(t, err, "rpc error: code = Unauthenticated desc = Invalid token format: requestor ID isn't string")
}

func TestConnectStandardAuthMiddleware_FailsOn_NoSessionId(t *testing.T) {
	f := newInterceptorFixture(t)

	f.mocks.jwth.EXPECT().Parse(mock.Anything, mock.Anything).Return(
		&jwt.Token{
			Claims: jwt.MapClaims{
				infinimesh.INFINIMESH_ACCOUNT_CLAIM: "test",
			},
			Valid: true,
		}, nil,
	)

	_, _, res, err := CallUnary(f.interceptor, context.Background(), connect.NewRequest(
		&node.EmptyMessage{},
	))

	assert.Nil(t, res)
	assert.EqualError(t, err, "rpc error: code = Unauthenticated desc = Invalid token format: no session ID")
}

func TestConnectStandardAuthMiddleware_FailsOn_SessionIdWrongType(t *testing.T) {
	f := newInterceptorFixture(t)

	f.mocks.jwth.EXPECT().Parse(mock.Anything, mock.Anything).Return(
		&jwt.Token{
			Claims: jwt.MapClaims{
				infinimesh.INFINIMESH_ACCOUNT_CLAIM: "test",
				infinimesh.INFINIMESH_SESSION_CLAIM: 666,
			},
			Valid: true,
		}, nil,
	)

	_, _, res, err := CallUnary(f.interceptor, context.Background(), connect.NewRequest(
		&node.EmptyMessage{},
	))

	assert.Nil(t, res)
	assert.EqualError(t, err, "rpc error: code = Unauthenticated desc = Invalid token format: session ID isn't string")
}

func TestConnectStandardAuthMiddleware_FailsOn_SessionExiredRevokedOrInvalid(t *testing.T) {
	f := newInterceptorFixture(t)

	f.mocks.jwth.EXPECT().Parse(mock.Anything, mock.Anything).Return(
		&jwt.Token{
			Claims: jwt.MapClaims{
				infinimesh.INFINIMESH_ACCOUNT_CLAIM: "test",
				infinimesh.INFINIMESH_SESSION_CLAIM: "test",
			},
			Valid: true,
		}, nil,
	)

	f.mocks.sessions.EXPECT().Check("test", "test").Return(errors.New("session is expired, revoked or invalid"))

	_, _, res, err := CallUnary(f.interceptor, context.Background(), connect.NewRequest(
		&node.EmptyMessage{},
	))

	assert.Nil(t, res)
	assert.EqualError(t, err, "rpc error: code = Unauthenticated desc = Session is expired, revoked or invalid")
}

func TestConnectStandardAuthMiddleware_Success(t *testing.T) {
	f := newInterceptorFixture(t)

	f.mocks.jwth.EXPECT().Parse(mock.Anything, mock.Anything).Return(
		&jwt.Token{
			Claims: jwt.MapClaims{
				infinimesh.INFINIMESH_ACCOUNT_CLAIM: "test",
				infinimesh.INFINIMESH_SESSION_CLAIM: "test",
				infinimesh.INFINIMESH_ROOT_CLAIM:    true,
				"exp":                               float64(777),
			},
			Valid: true,
		}, nil,
	)

	f.mocks.sessions.EXPECT().Check("test", "test").Return(nil)
	f.mocks.sessions.EXPECT().LogActivity("test", "test", mock.Anything).Return(nil).WaitUntil(time.After(1 * time.Second))

	ctx, _, res, err := CallUnary(f.interceptor, context.Background(), connect.NewRequest(
		&node.EmptyMessage{},
	))

	runtime.Gosched() // To ensure that goroutine is executed

	assert.Nil(t, res)
	assert.NoError(t, err)

	assert.Equal(t, "test", ctx.Value(infinimesh.InfinimeshAccountCtxKey))
	assert.Equal(t, "test", ctx.Value(infinimesh.InfinimeshSessionCtxKey))
	assert.Equal(t, int64(777), ctx.Value(infinimesh.ContextKey("exp")))
	assert.Equal(t, true, ctx.Value(infinimesh.InfinimeshRootCtxKey))

	md, ok := metadata.FromOutgoingContext(ctx)
	assert.True(t, ok)
	assert.Equal(t, []string{"test"}, md.Get(infinimesh.INFINIMESH_ACCOUNT_CLAIM))
	assert.Equal(t, []string{"Bearer "}, md.Get("authorization"))
}

// ConnectDeviceAuthMiddleware

func TestConnectDeviceAuthMiddleware_FailsOn_NoToken(t *testing.T) {
	f := newInterceptorFixture(t)

	f.mocks.jwth.EXPECT().Parse("test", mock.Anything).Return(nil, errors.New("token contains an invalid number of segments"))

	_, _, err := f.interceptor.ConnectDeviceAuthMiddleware(context.Background(), []byte{}, ("test"))

	assert.EqualError(t, err, "token contains an invalid number of segments")
}

func TestConnectDeviceAuthMiddleware_FailsOn_NoDevicesScope(t *testing.T) {
	f := newInterceptorFixture(t)

	f.mocks.jwth.EXPECT().Parse("test", mock.Anything).Return(
		&jwt.Token{
			Claims: jwt.MapClaims{},
			Valid:  true,
		}, nil,
	)

	_, _, err := f.interceptor.ConnectDeviceAuthMiddleware(context.Background(), []byte{}, ("test"))

	assert.EqualError(t, err, "rpc error: code = Unauthenticated desc = Invalid token format: no devices scope")
}

func TestConnectDeviceAuthMiddleware_FailsOn_WrongType(t *testing.T) {
	f := newInterceptorFixture(t)

	f.mocks.jwth.EXPECT().Parse("test", mock.Anything).Return(
		&jwt.Token{
			Claims: jwt.MapClaims{
				infinimesh.INFINIMESH_DEVICES_CLAIM: 666,
			},
			Valid: true,
		}, nil,
	)

	_, _, err := f.interceptor.ConnectDeviceAuthMiddleware(context.Background(), []byte{}, ("test"))

	assert.EqualError(t, err, "rpc error: code = Unauthenticated desc = Invalid token format: devices scope isn't a map")
}

func TestConnectDeviceAuthMiddleware_FailsOn_WrongValueType(t *testing.T) {
	f := newInterceptorFixture(t)

	f.mocks.jwth.EXPECT().Parse("test", mock.Anything).Return(
		&jwt.Token{
			Claims: jwt.MapClaims{
				infinimesh.INFINIMESH_DEVICES_CLAIM: map[string]any{
					"uuid": "invalid",
				},
			},
			Valid: true,
		}, nil,
	)

	_, _, err := f.interceptor.ConnectDeviceAuthMiddleware(context.Background(), []byte{}, ("test"))

	assert.EqualError(t, err, "rpc error: code = Unauthenticated desc = Invalid token format: element invalid is not a number")
}

func TestConnectDeviceAuthMiddleware_Success(t *testing.T) {
	f := newInterceptorFixture(t)

	f.mocks.jwth.EXPECT().Parse("test", mock.Anything).Return(
		&jwt.Token{
			Claims: jwt.MapClaims{
				infinimesh.INFINIMESH_DEVICES_CLAIM: map[string]any{
					"uuid": float64(1),
				},
			},
			Valid: true,
		}, nil,
	)

	ctx, log, err := f.interceptor.ConnectDeviceAuthMiddleware(context.Background(), []byte{}, ("test"))

	assert.NoError(t, err)
	assert.Equal(t, map[string]access.Level{"uuid": 1}, ctx.Value(infinimesh.InfinimeshDevicesCtxKey))
	assert.Equal(t, false, log)
}
