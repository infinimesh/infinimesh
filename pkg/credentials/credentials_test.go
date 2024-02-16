package credentials_test

import (
	"testing"

	driver_mocks "github.com/infinimesh/infinimesh/mocks/github.com/arangodb/go-driver"
	"github.com/infinimesh/infinimesh/pkg/credentials"
	"go.uber.org/zap"
)

type credentialsControllerFixture struct {
	ctrl credentials.CredentialsController

	mocks struct {
		db *driver_mocks.MockDatabase
	}
}

func newCredentialsControllerFixture(t *testing.T) *credentialsControllerFixture {
	f := &credentialsControllerFixture{}

	f.mocks.db = &driver_mocks.MockDatabase{}

	f.ctrl = credentials.NewCredentialsController(zap.NewExample(), f.mocks.db)

	return f
}
