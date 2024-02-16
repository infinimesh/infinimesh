package credentials

import (
	"context"
	"errors"

	"github.com/arangodb/go-driver"
	"go.uber.org/zap"
)

type MockCredentials struct {
	Args []string

	log   *zap.Logger
	valid bool
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

	return c.Args[0] == "valid"
}

func (c *MockCredentials) Find(context.Context, driver.Database) bool {
	if len(c.Args) == 0 {
		c.valid = false
		return false
	}

	if c.Args[0] == "valid" {
		c.valid = true
		return true
	}

	c.valid = false
	return false
}

func (cred *MockCredentials) FindByKey(ctx context.Context, col driver.Collection, key string) error {
	if key == "valid" {
		return nil
	}
	return errors.New("not found")
}
