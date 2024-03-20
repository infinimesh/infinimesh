// Code generated by mockery v2.40.1. DO NOT EDIT.

package credentials_mocks

import (
	context "context"

	accounts "github.com/infinimesh/proto/node/accounts"

	credentials "github.com/infinimesh/infinimesh/pkg/credentials"

	driver "github.com/arangodb/go-driver"

	mock "github.com/stretchr/testify/mock"
)

// MockCredentialsController is an autogenerated mock type for the CredentialsController type
type MockCredentialsController struct {
	mock.Mock
}

type MockCredentialsController_Expecter struct {
	mock *mock.Mock
}

func (_m *MockCredentialsController) EXPECT() *MockCredentialsController_Expecter {
	return &MockCredentialsController_Expecter{mock: &_m.Mock}
}

// Authorisable provides a mock function with given fields: ctx, cred
func (_m *MockCredentialsController) Authorisable(ctx context.Context, cred *credentials.Credentials) (*accounts.Account, bool) {
	ret := _m.Called(ctx, cred)

	if len(ret) == 0 {
		panic("no return value specified for Authorisable")
	}

	var r0 *accounts.Account
	var r1 bool
	if rf, ok := ret.Get(0).(func(context.Context, *credentials.Credentials) (*accounts.Account, bool)); ok {
		return rf(ctx, cred)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *credentials.Credentials) *accounts.Account); ok {
		r0 = rf(ctx, cred)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*accounts.Account)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *credentials.Credentials) bool); ok {
		r1 = rf(ctx, cred)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// MockCredentialsController_Authorisable_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Authorisable'
type MockCredentialsController_Authorisable_Call struct {
	*mock.Call
}

// Authorisable is a helper method to define mock.On call
//   - ctx context.Context
//   - cred *credentials.Credentials
func (_e *MockCredentialsController_Expecter) Authorisable(ctx interface{}, cred interface{}) *MockCredentialsController_Authorisable_Call {
	return &MockCredentialsController_Authorisable_Call{Call: _e.mock.On("Authorisable", ctx, cred)}
}

func (_c *MockCredentialsController_Authorisable_Call) Run(run func(ctx context.Context, cred *credentials.Credentials)) *MockCredentialsController_Authorisable_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*credentials.Credentials))
	})
	return _c
}

func (_c *MockCredentialsController_Authorisable_Call) Return(_a0 *accounts.Account, _a1 bool) *MockCredentialsController_Authorisable_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCredentialsController_Authorisable_Call) RunAndReturn(run func(context.Context, *credentials.Credentials) (*accounts.Account, bool)) *MockCredentialsController_Authorisable_Call {
	_c.Call.Return(run)
	return _c
}

// Authorize provides a mock function with given fields: ctx, auth_type, args
func (_m *MockCredentialsController) Authorize(ctx context.Context, auth_type string, args ...string) (*accounts.Account, bool) {
	_va := make([]interface{}, len(args))
	for _i := range args {
		_va[_i] = args[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, auth_type)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Authorize")
	}

	var r0 *accounts.Account
	var r1 bool
	if rf, ok := ret.Get(0).(func(context.Context, string, ...string) (*accounts.Account, bool)); ok {
		return rf(ctx, auth_type, args...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, ...string) *accounts.Account); ok {
		r0 = rf(ctx, auth_type, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*accounts.Account)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, ...string) bool); ok {
		r1 = rf(ctx, auth_type, args...)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// MockCredentialsController_Authorize_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Authorize'
type MockCredentialsController_Authorize_Call struct {
	*mock.Call
}

// Authorize is a helper method to define mock.On call
//   - ctx context.Context
//   - auth_type string
//   - args ...string
func (_e *MockCredentialsController_Expecter) Authorize(ctx interface{}, auth_type interface{}, args ...interface{}) *MockCredentialsController_Authorize_Call {
	return &MockCredentialsController_Authorize_Call{Call: _e.mock.On("Authorize",
		append([]interface{}{ctx, auth_type}, args...)...)}
}

func (_c *MockCredentialsController_Authorize_Call) Run(run func(ctx context.Context, auth_type string, args ...string)) *MockCredentialsController_Authorize_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]string, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(string)
			}
		}
		run(args[0].(context.Context), args[1].(string), variadicArgs...)
	})
	return _c
}

func (_c *MockCredentialsController_Authorize_Call) Return(_a0 *accounts.Account, _a1 bool) *MockCredentialsController_Authorize_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCredentialsController_Authorize_Call) RunAndReturn(run func(context.Context, string, ...string) (*accounts.Account, bool)) *MockCredentialsController_Authorize_Call {
	_c.Call.Return(run)
	return _c
}

// Find provides a mock function with given fields: ctx, auth_type, args
func (_m *MockCredentialsController) Find(ctx context.Context, auth_type string, args ...string) (credentials.Credentials, error) {
	_va := make([]interface{}, len(args))
	for _i := range args {
		_va[_i] = args[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, auth_type)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Find")
	}

	var r0 credentials.Credentials
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, ...string) (credentials.Credentials, error)); ok {
		return rf(ctx, auth_type, args...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, ...string) credentials.Credentials); ok {
		r0 = rf(ctx, auth_type, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(credentials.Credentials)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, ...string) error); ok {
		r1 = rf(ctx, auth_type, args...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCredentialsController_Find_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Find'
type MockCredentialsController_Find_Call struct {
	*mock.Call
}

// Find is a helper method to define mock.On call
//   - ctx context.Context
//   - auth_type string
//   - args ...string
func (_e *MockCredentialsController_Expecter) Find(ctx interface{}, auth_type interface{}, args ...interface{}) *MockCredentialsController_Find_Call {
	return &MockCredentialsController_Find_Call{Call: _e.mock.On("Find",
		append([]interface{}{ctx, auth_type}, args...)...)}
}

func (_c *MockCredentialsController_Find_Call) Run(run func(ctx context.Context, auth_type string, args ...string)) *MockCredentialsController_Find_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]string, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(string)
			}
		}
		run(args[0].(context.Context), args[1].(string), variadicArgs...)
	})
	return _c
}

func (_c *MockCredentialsController_Find_Call) Return(cred credentials.Credentials, err error) *MockCredentialsController_Find_Call {
	_c.Call.Return(cred, err)
	return _c
}

func (_c *MockCredentialsController_Find_Call) RunAndReturn(run func(context.Context, string, ...string) (credentials.Credentials, error)) *MockCredentialsController_Find_Call {
	_c.Call.Return(run)
	return _c
}

// ListCredentials provides a mock function with given fields: ctx, acc
func (_m *MockCredentialsController) ListCredentials(ctx context.Context, acc driver.DocumentID) ([]credentials.ListCredentialsResponse, error) {
	ret := _m.Called(ctx, acc)

	if len(ret) == 0 {
		panic("no return value specified for ListCredentials")
	}

	var r0 []credentials.ListCredentialsResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, driver.DocumentID) ([]credentials.ListCredentialsResponse, error)); ok {
		return rf(ctx, acc)
	}
	if rf, ok := ret.Get(0).(func(context.Context, driver.DocumentID) []credentials.ListCredentialsResponse); ok {
		r0 = rf(ctx, acc)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]credentials.ListCredentialsResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, driver.DocumentID) error); ok {
		r1 = rf(ctx, acc)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCredentialsController_ListCredentials_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListCredentials'
type MockCredentialsController_ListCredentials_Call struct {
	*mock.Call
}

// ListCredentials is a helper method to define mock.On call
//   - ctx context.Context
//   - acc driver.DocumentID
func (_e *MockCredentialsController_Expecter) ListCredentials(ctx interface{}, acc interface{}) *MockCredentialsController_ListCredentials_Call {
	return &MockCredentialsController_ListCredentials_Call{Call: _e.mock.On("ListCredentials", ctx, acc)}
}

func (_c *MockCredentialsController_ListCredentials_Call) Run(run func(ctx context.Context, acc driver.DocumentID)) *MockCredentialsController_ListCredentials_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(driver.DocumentID))
	})
	return _c
}

func (_c *MockCredentialsController_ListCredentials_Call) Return(r []credentials.ListCredentialsResponse, err error) *MockCredentialsController_ListCredentials_Call {
	_c.Call.Return(r, err)
	return _c
}

func (_c *MockCredentialsController_ListCredentials_Call) RunAndReturn(run func(context.Context, driver.DocumentID) ([]credentials.ListCredentialsResponse, error)) *MockCredentialsController_ListCredentials_Call {
	_c.Call.Return(run)
	return _c
}

// ListCredentialsAndEdges provides a mock function with given fields: ctx, account
func (_m *MockCredentialsController) ListCredentialsAndEdges(ctx context.Context, account driver.DocumentID) ([]string, error) {
	ret := _m.Called(ctx, account)

	if len(ret) == 0 {
		panic("no return value specified for ListCredentialsAndEdges")
	}

	var r0 []string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, driver.DocumentID) ([]string, error)); ok {
		return rf(ctx, account)
	}
	if rf, ok := ret.Get(0).(func(context.Context, driver.DocumentID) []string); ok {
		r0 = rf(ctx, account)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, driver.DocumentID) error); ok {
		r1 = rf(ctx, account)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCredentialsController_ListCredentialsAndEdges_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListCredentialsAndEdges'
type MockCredentialsController_ListCredentialsAndEdges_Call struct {
	*mock.Call
}

// ListCredentialsAndEdges is a helper method to define mock.On call
//   - ctx context.Context
//   - account driver.DocumentID
func (_e *MockCredentialsController_Expecter) ListCredentialsAndEdges(ctx interface{}, account interface{}) *MockCredentialsController_ListCredentialsAndEdges_Call {
	return &MockCredentialsController_ListCredentialsAndEdges_Call{Call: _e.mock.On("ListCredentialsAndEdges", ctx, account)}
}

func (_c *MockCredentialsController_ListCredentialsAndEdges_Call) Run(run func(ctx context.Context, account driver.DocumentID)) *MockCredentialsController_ListCredentialsAndEdges_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(driver.DocumentID))
	})
	return _c
}

func (_c *MockCredentialsController_ListCredentialsAndEdges_Call) Return(nodes []string, err error) *MockCredentialsController_ListCredentialsAndEdges_Call {
	_c.Call.Return(nodes, err)
	return _c
}

func (_c *MockCredentialsController_ListCredentialsAndEdges_Call) RunAndReturn(run func(context.Context, driver.DocumentID) ([]string, error)) *MockCredentialsController_ListCredentialsAndEdges_Call {
	_c.Call.Return(run)
	return _c
}

// MakeCredentials provides a mock function with given fields: _a0
func (_m *MockCredentialsController) MakeCredentials(_a0 *accounts.Credentials) (credentials.Credentials, error) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for MakeCredentials")
	}

	var r0 credentials.Credentials
	var r1 error
	if rf, ok := ret.Get(0).(func(*accounts.Credentials) (credentials.Credentials, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(*accounts.Credentials) credentials.Credentials); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(credentials.Credentials)
		}
	}

	if rf, ok := ret.Get(1).(func(*accounts.Credentials) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCredentialsController_MakeCredentials_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'MakeCredentials'
type MockCredentialsController_MakeCredentials_Call struct {
	*mock.Call
}

// MakeCredentials is a helper method to define mock.On call
//   - _a0 *accounts.Credentials
func (_e *MockCredentialsController_Expecter) MakeCredentials(_a0 interface{}) *MockCredentialsController_MakeCredentials_Call {
	return &MockCredentialsController_MakeCredentials_Call{Call: _e.mock.On("MakeCredentials", _a0)}
}

func (_c *MockCredentialsController_MakeCredentials_Call) Run(run func(_a0 *accounts.Credentials)) *MockCredentialsController_MakeCredentials_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*accounts.Credentials))
	})
	return _c
}

func (_c *MockCredentialsController_MakeCredentials_Call) Return(_a0 credentials.Credentials, _a1 error) *MockCredentialsController_MakeCredentials_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCredentialsController_MakeCredentials_Call) RunAndReturn(run func(*accounts.Credentials) (credentials.Credentials, error)) *MockCredentialsController_MakeCredentials_Call {
	_c.Call.Return(run)
	return _c
}

// MakeListable provides a mock function with given fields: r
func (_m *MockCredentialsController) MakeListable(r credentials.ListCredentialsResponse) (credentials.ListableCredentials, error) {
	ret := _m.Called(r)

	if len(ret) == 0 {
		panic("no return value specified for MakeListable")
	}

	var r0 credentials.ListableCredentials
	var r1 error
	if rf, ok := ret.Get(0).(func(credentials.ListCredentialsResponse) (credentials.ListableCredentials, error)); ok {
		return rf(r)
	}
	if rf, ok := ret.Get(0).(func(credentials.ListCredentialsResponse) credentials.ListableCredentials); ok {
		r0 = rf(r)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(credentials.ListableCredentials)
		}
	}

	if rf, ok := ret.Get(1).(func(credentials.ListCredentialsResponse) error); ok {
		r1 = rf(r)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCredentialsController_MakeListable_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'MakeListable'
type MockCredentialsController_MakeListable_Call struct {
	*mock.Call
}

// MakeListable is a helper method to define mock.On call
//   - r credentials.ListCredentialsResponse
func (_e *MockCredentialsController_Expecter) MakeListable(r interface{}) *MockCredentialsController_MakeListable_Call {
	return &MockCredentialsController_MakeListable_Call{Call: _e.mock.On("MakeListable", r)}
}

func (_c *MockCredentialsController_MakeListable_Call) Run(run func(r credentials.ListCredentialsResponse)) *MockCredentialsController_MakeListable_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(credentials.ListCredentialsResponse))
	})
	return _c
}

func (_c *MockCredentialsController_MakeListable_Call) Return(_a0 credentials.ListableCredentials, _a1 error) *MockCredentialsController_MakeListable_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCredentialsController_MakeListable_Call) RunAndReturn(run func(credentials.ListCredentialsResponse) (credentials.ListableCredentials, error)) *MockCredentialsController_MakeListable_Call {
	_c.Call.Return(run)
	return _c
}

// SetCredentials provides a mock function with given fields: ctx, acc, c
func (_m *MockCredentialsController) SetCredentials(ctx context.Context, acc driver.DocumentID, c credentials.Credentials) error {
	ret := _m.Called(ctx, acc, c)

	if len(ret) == 0 {
		panic("no return value specified for SetCredentials")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, driver.DocumentID, credentials.Credentials) error); ok {
		r0 = rf(ctx, acc, c)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockCredentialsController_SetCredentials_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetCredentials'
type MockCredentialsController_SetCredentials_Call struct {
	*mock.Call
}

// SetCredentials is a helper method to define mock.On call
//   - ctx context.Context
//   - acc driver.DocumentID
//   - c credentials.Credentials
func (_e *MockCredentialsController_Expecter) SetCredentials(ctx interface{}, acc interface{}, c interface{}) *MockCredentialsController_SetCredentials_Call {
	return &MockCredentialsController_SetCredentials_Call{Call: _e.mock.On("SetCredentials", ctx, acc, c)}
}

func (_c *MockCredentialsController_SetCredentials_Call) Run(run func(ctx context.Context, acc driver.DocumentID, c credentials.Credentials)) *MockCredentialsController_SetCredentials_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(driver.DocumentID), args[2].(credentials.Credentials))
	})
	return _c
}

func (_c *MockCredentialsController_SetCredentials_Call) Return(_a0 error) *MockCredentialsController_SetCredentials_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockCredentialsController_SetCredentials_Call) RunAndReturn(run func(context.Context, driver.DocumentID, credentials.Credentials) error) *MockCredentialsController_SetCredentials_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockCredentialsController creates a new instance of MockCredentialsController. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockCredentialsController(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockCredentialsController {
	mock := &MockCredentialsController{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
