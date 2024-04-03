// Code generated by mockery v2.42.1. DO NOT EDIT.

package node_mocks

import (
	context "context"

	grpc "google.golang.org/grpc"

	mock "github.com/stretchr/testify/mock"

	node "github.com/infinimesh/proto/node"
)

// MockInternalServiceClient is an autogenerated mock type for the InternalServiceClient type
type MockInternalServiceClient struct {
	mock.Mock
}

type MockInternalServiceClient_Expecter struct {
	mock *mock.Mock
}

func (_m *MockInternalServiceClient) EXPECT() *MockInternalServiceClient_Expecter {
	return &MockInternalServiceClient_Expecter{mock: &_m.Mock}
}

// GetLDAPProviders provides a mock function with given fields: ctx, in, opts
func (_m *MockInternalServiceClient) GetLDAPProviders(ctx context.Context, in *node.EmptyMessage, opts ...grpc.CallOption) (*node.LDAPProviders, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for GetLDAPProviders")
	}

	var r0 *node.LDAPProviders
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *node.EmptyMessage, ...grpc.CallOption) (*node.LDAPProviders, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *node.EmptyMessage, ...grpc.CallOption) *node.LDAPProviders); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*node.LDAPProviders)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *node.EmptyMessage, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockInternalServiceClient_GetLDAPProviders_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetLDAPProviders'
type MockInternalServiceClient_GetLDAPProviders_Call struct {
	*mock.Call
}

// GetLDAPProviders is a helper method to define mock.On call
//   - ctx context.Context
//   - in *node.EmptyMessage
//   - opts ...grpc.CallOption
func (_e *MockInternalServiceClient_Expecter) GetLDAPProviders(ctx interface{}, in interface{}, opts ...interface{}) *MockInternalServiceClient_GetLDAPProviders_Call {
	return &MockInternalServiceClient_GetLDAPProviders_Call{Call: _e.mock.On("GetLDAPProviders",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MockInternalServiceClient_GetLDAPProviders_Call) Run(run func(ctx context.Context, in *node.EmptyMessage, opts ...grpc.CallOption)) *MockInternalServiceClient_GetLDAPProviders_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*node.EmptyMessage), variadicArgs...)
	})
	return _c
}

func (_c *MockInternalServiceClient_GetLDAPProviders_Call) Return(_a0 *node.LDAPProviders, _a1 error) *MockInternalServiceClient_GetLDAPProviders_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockInternalServiceClient_GetLDAPProviders_Call) RunAndReturn(run func(context.Context, *node.EmptyMessage, ...grpc.CallOption) (*node.LDAPProviders, error)) *MockInternalServiceClient_GetLDAPProviders_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockInternalServiceClient creates a new instance of MockInternalServiceClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockInternalServiceClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockInternalServiceClient {
	mock := &MockInternalServiceClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
