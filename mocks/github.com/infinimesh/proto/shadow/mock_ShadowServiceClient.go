// Code generated by mockery v2.40.1. DO NOT EDIT.

package shadow_mocks

import (
	context "context"

	grpc "google.golang.org/grpc"

	mock "github.com/stretchr/testify/mock"

	shadow "github.com/infinimesh/proto/shadow"
)

// MockShadowServiceClient is an autogenerated mock type for the ShadowServiceClient type
type MockShadowServiceClient struct {
	mock.Mock
}

type MockShadowServiceClient_Expecter struct {
	mock *mock.Mock
}

func (_m *MockShadowServiceClient) EXPECT() *MockShadowServiceClient_Expecter {
	return &MockShadowServiceClient_Expecter{mock: &_m.Mock}
}

// Get provides a mock function with given fields: ctx, in, opts
func (_m *MockShadowServiceClient) Get(ctx context.Context, in *shadow.GetRequest, opts ...grpc.CallOption) (*shadow.GetResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *shadow.GetResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *shadow.GetRequest, ...grpc.CallOption) (*shadow.GetResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *shadow.GetRequest, ...grpc.CallOption) *shadow.GetResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*shadow.GetResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *shadow.GetRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockShadowServiceClient_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type MockShadowServiceClient_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx context.Context
//   - in *shadow.GetRequest
//   - opts ...grpc.CallOption
func (_e *MockShadowServiceClient_Expecter) Get(ctx interface{}, in interface{}, opts ...interface{}) *MockShadowServiceClient_Get_Call {
	return &MockShadowServiceClient_Get_Call{Call: _e.mock.On("Get",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MockShadowServiceClient_Get_Call) Run(run func(ctx context.Context, in *shadow.GetRequest, opts ...grpc.CallOption)) *MockShadowServiceClient_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*shadow.GetRequest), variadicArgs...)
	})
	return _c
}

func (_c *MockShadowServiceClient_Get_Call) Return(_a0 *shadow.GetResponse, _a1 error) *MockShadowServiceClient_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockShadowServiceClient_Get_Call) RunAndReturn(run func(context.Context, *shadow.GetRequest, ...grpc.CallOption) (*shadow.GetResponse, error)) *MockShadowServiceClient_Get_Call {
	_c.Call.Return(run)
	return _c
}

// Patch provides a mock function with given fields: ctx, in, opts
func (_m *MockShadowServiceClient) Patch(ctx context.Context, in *shadow.Shadow, opts ...grpc.CallOption) (*shadow.Shadow, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Patch")
	}

	var r0 *shadow.Shadow
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *shadow.Shadow, ...grpc.CallOption) (*shadow.Shadow, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *shadow.Shadow, ...grpc.CallOption) *shadow.Shadow); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*shadow.Shadow)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *shadow.Shadow, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockShadowServiceClient_Patch_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Patch'
type MockShadowServiceClient_Patch_Call struct {
	*mock.Call
}

// Patch is a helper method to define mock.On call
//   - ctx context.Context
//   - in *shadow.Shadow
//   - opts ...grpc.CallOption
func (_e *MockShadowServiceClient_Expecter) Patch(ctx interface{}, in interface{}, opts ...interface{}) *MockShadowServiceClient_Patch_Call {
	return &MockShadowServiceClient_Patch_Call{Call: _e.mock.On("Patch",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MockShadowServiceClient_Patch_Call) Run(run func(ctx context.Context, in *shadow.Shadow, opts ...grpc.CallOption)) *MockShadowServiceClient_Patch_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*shadow.Shadow), variadicArgs...)
	})
	return _c
}

func (_c *MockShadowServiceClient_Patch_Call) Return(_a0 *shadow.Shadow, _a1 error) *MockShadowServiceClient_Patch_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockShadowServiceClient_Patch_Call) RunAndReturn(run func(context.Context, *shadow.Shadow, ...grpc.CallOption) (*shadow.Shadow, error)) *MockShadowServiceClient_Patch_Call {
	_c.Call.Return(run)
	return _c
}

// Remove provides a mock function with given fields: ctx, in, opts
func (_m *MockShadowServiceClient) Remove(ctx context.Context, in *shadow.RemoveRequest, opts ...grpc.CallOption) (*shadow.Shadow, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Remove")
	}

	var r0 *shadow.Shadow
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *shadow.RemoveRequest, ...grpc.CallOption) (*shadow.Shadow, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *shadow.RemoveRequest, ...grpc.CallOption) *shadow.Shadow); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*shadow.Shadow)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *shadow.RemoveRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockShadowServiceClient_Remove_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Remove'
type MockShadowServiceClient_Remove_Call struct {
	*mock.Call
}

// Remove is a helper method to define mock.On call
//   - ctx context.Context
//   - in *shadow.RemoveRequest
//   - opts ...grpc.CallOption
func (_e *MockShadowServiceClient_Expecter) Remove(ctx interface{}, in interface{}, opts ...interface{}) *MockShadowServiceClient_Remove_Call {
	return &MockShadowServiceClient_Remove_Call{Call: _e.mock.On("Remove",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MockShadowServiceClient_Remove_Call) Run(run func(ctx context.Context, in *shadow.RemoveRequest, opts ...grpc.CallOption)) *MockShadowServiceClient_Remove_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*shadow.RemoveRequest), variadicArgs...)
	})
	return _c
}

func (_c *MockShadowServiceClient_Remove_Call) Return(_a0 *shadow.Shadow, _a1 error) *MockShadowServiceClient_Remove_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockShadowServiceClient_Remove_Call) RunAndReturn(run func(context.Context, *shadow.RemoveRequest, ...grpc.CallOption) (*shadow.Shadow, error)) *MockShadowServiceClient_Remove_Call {
	_c.Call.Return(run)
	return _c
}

// StreamShadow provides a mock function with given fields: ctx, in, opts
func (_m *MockShadowServiceClient) StreamShadow(ctx context.Context, in *shadow.StreamShadowRequest, opts ...grpc.CallOption) (shadow.ShadowService_StreamShadowClient, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for StreamShadow")
	}

	var r0 shadow.ShadowService_StreamShadowClient
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *shadow.StreamShadowRequest, ...grpc.CallOption) (shadow.ShadowService_StreamShadowClient, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *shadow.StreamShadowRequest, ...grpc.CallOption) shadow.ShadowService_StreamShadowClient); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(shadow.ShadowService_StreamShadowClient)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *shadow.StreamShadowRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockShadowServiceClient_StreamShadow_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'StreamShadow'
type MockShadowServiceClient_StreamShadow_Call struct {
	*mock.Call
}

// StreamShadow is a helper method to define mock.On call
//   - ctx context.Context
//   - in *shadow.StreamShadowRequest
//   - opts ...grpc.CallOption
func (_e *MockShadowServiceClient_Expecter) StreamShadow(ctx interface{}, in interface{}, opts ...interface{}) *MockShadowServiceClient_StreamShadow_Call {
	return &MockShadowServiceClient_StreamShadow_Call{Call: _e.mock.On("StreamShadow",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MockShadowServiceClient_StreamShadow_Call) Run(run func(ctx context.Context, in *shadow.StreamShadowRequest, opts ...grpc.CallOption)) *MockShadowServiceClient_StreamShadow_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*shadow.StreamShadowRequest), variadicArgs...)
	})
	return _c
}

func (_c *MockShadowServiceClient_StreamShadow_Call) Return(_a0 shadow.ShadowService_StreamShadowClient, _a1 error) *MockShadowServiceClient_StreamShadow_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockShadowServiceClient_StreamShadow_Call) RunAndReturn(run func(context.Context, *shadow.StreamShadowRequest, ...grpc.CallOption) (shadow.ShadowService_StreamShadowClient, error)) *MockShadowServiceClient_StreamShadow_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockShadowServiceClient creates a new instance of MockShadowServiceClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockShadowServiceClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockShadowServiceClient {
	mock := &MockShadowServiceClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
