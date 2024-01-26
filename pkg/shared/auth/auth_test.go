package auth_test

import (
	"testing"

	auth_mocks "github.com/infinimesh/infinimesh/mocks/github.com/infinimesh/infinimesh/pkg/shared/auth"
	"github.com/infinimesh/infinimesh/pkg/shared/auth"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

type interceptorFixture struct {
	interceptor auth.AuthInterceptor

	mocks struct {
		jwth     *auth_mocks.MockJWTHandler
		log      *zap.Logger
		observer *observer.ObservedLogs
	}
}

func newInterceptorFixture(t *testing.T) *interceptorFixture {
	f := &interceptorFixture{}

	core, observer := observer.New(zap.DebugLevel)
	f.mocks.log = zap.New(core)
	f.mocks.observer = observer
	f.mocks.jwth = auth_mocks.NewMockJWTHandler(t)

	f.interceptor = auth.NewAuthInterceptor(f.mocks.log, nil, f.mocks.jwth, nil)
	return f
}
