package credentials

import (
	"context"
	"errors"

	"github.com/arangodb/go-driver"
	"go.uber.org/zap"
)

// MockCredentials is a mock implementation of the Credentials interface
// Method name can be specified in the first argument to make it fail
type MockCredentials struct {
	Args []string

	log   *zap.Logger
	valid bool
}

func NewMockCredentials(args ...string) (Credentials, error) {
	if args[0] == "NewMockCredentials" {
		return nil, errors.New("invalid")
	}
	return &MockCredentials{Args: args}, nil
}

func (c *MockCredentials) SetLogger(l *zap.Logger) {
	c.log = l
}

func (c *MockCredentials) Type() string {
	return "mock"
}

func (c *MockCredentials) Key() string {
	return "mock"
}

func (c *MockCredentials) Authorize(args ...string) bool {
	if len(c.Args) == 0 {
		return false
	}

	return c.Args[0] != "Authorize"
}

func (c *MockCredentials) Find(context.Context, driver.Database) bool {
	if len(c.Args) == 0 {
		c.valid = false
		return false
	}

	if c.Args[0] != "Find" {
		c.valid = true
		return true
	}

	c.valid = false
	return false
}

func (cred *MockCredentials) FindByKey(ctx context.Context, col driver.Collection, key string) error {
	if key != "FindByKey" {
		return nil
	}
	return errors.New("not found")
}
