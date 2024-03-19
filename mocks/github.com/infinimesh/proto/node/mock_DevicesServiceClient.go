// Code generated by mockery v2.42.0. DO NOT EDIT.

package node_mocks

import (
	context "context"

	access "github.com/infinimesh/proto/node/access"

	devices "github.com/infinimesh/proto/node/devices"

	grpc "google.golang.org/grpc"

	mock "github.com/stretchr/testify/mock"

	node "github.com/infinimesh/proto/node"
)

// MockDevicesServiceClient is an autogenerated mock type for the DevicesServiceClient type
type MockDevicesServiceClient struct {
	mock.Mock
}

type MockDevicesServiceClient_Expecter struct {
	mock *mock.Mock
}

func (_m *MockDevicesServiceClient) EXPECT() *MockDevicesServiceClient_Expecter {
	return &MockDevicesServiceClient_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: ctx, in, opts
func (_m *MockDevicesServiceClient) Create(ctx context.Context, in *devices.CreateRequest, opts ...grpc.CallOption) (*devices.CreateResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 *devices.CreateResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *devices.CreateRequest, ...grpc.CallOption) (*devices.CreateResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *devices.CreateRequest, ...grpc.CallOption) *devices.CreateResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*devices.CreateResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *devices.CreateRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDevicesServiceClient_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type MockDevicesServiceClient_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - in *devices.CreateRequest
//   - opts ...grpc.CallOption
func (_e *MockDevicesServiceClient_Expecter) Create(ctx interface{}, in interface{}, opts ...interface{}) *MockDevicesServiceClient_Create_Call {
	return &MockDevicesServiceClient_Create_Call{Call: _e.mock.On("Create",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MockDevicesServiceClient_Create_Call) Run(run func(ctx context.Context, in *devices.CreateRequest, opts ...grpc.CallOption)) *MockDevicesServiceClient_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*devices.CreateRequest), variadicArgs...)
	})
	return _c
}

func (_c *MockDevicesServiceClient_Create_Call) Return(_a0 *devices.CreateResponse, _a1 error) *MockDevicesServiceClient_Create_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDevicesServiceClient_Create_Call) RunAndReturn(run func(context.Context, *devices.CreateRequest, ...grpc.CallOption) (*devices.CreateResponse, error)) *MockDevicesServiceClient_Create_Call {
	_c.Call.Return(run)
	return _c
}

// Delete provides a mock function with given fields: ctx, in, opts
func (_m *MockDevicesServiceClient) Delete(ctx context.Context, in *devices.Device, opts ...grpc.CallOption) (*node.DeleteResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 *node.DeleteResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *devices.Device, ...grpc.CallOption) (*node.DeleteResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *devices.Device, ...grpc.CallOption) *node.DeleteResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*node.DeleteResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *devices.Device, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDevicesServiceClient_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type MockDevicesServiceClient_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - ctx context.Context
//   - in *devices.Device
//   - opts ...grpc.CallOption
func (_e *MockDevicesServiceClient_Expecter) Delete(ctx interface{}, in interface{}, opts ...interface{}) *MockDevicesServiceClient_Delete_Call {
	return &MockDevicesServiceClient_Delete_Call{Call: _e.mock.On("Delete",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MockDevicesServiceClient_Delete_Call) Run(run func(ctx context.Context, in *devices.Device, opts ...grpc.CallOption)) *MockDevicesServiceClient_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*devices.Device), variadicArgs...)
	})
	return _c
}

func (_c *MockDevicesServiceClient_Delete_Call) Return(_a0 *node.DeleteResponse, _a1 error) *MockDevicesServiceClient_Delete_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDevicesServiceClient_Delete_Call) RunAndReturn(run func(context.Context, *devices.Device, ...grpc.CallOption) (*node.DeleteResponse, error)) *MockDevicesServiceClient_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: ctx, in, opts
func (_m *MockDevicesServiceClient) Get(ctx context.Context, in *devices.Device, opts ...grpc.CallOption) (*devices.Device, error) {
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

	var r0 *devices.Device
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *devices.Device, ...grpc.CallOption) (*devices.Device, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *devices.Device, ...grpc.CallOption) *devices.Device); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*devices.Device)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *devices.Device, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDevicesServiceClient_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type MockDevicesServiceClient_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx context.Context
//   - in *devices.Device
//   - opts ...grpc.CallOption
func (_e *MockDevicesServiceClient_Expecter) Get(ctx interface{}, in interface{}, opts ...interface{}) *MockDevicesServiceClient_Get_Call {
	return &MockDevicesServiceClient_Get_Call{Call: _e.mock.On("Get",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MockDevicesServiceClient_Get_Call) Run(run func(ctx context.Context, in *devices.Device, opts ...grpc.CallOption)) *MockDevicesServiceClient_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*devices.Device), variadicArgs...)
	})
	return _c
}

func (_c *MockDevicesServiceClient_Get_Call) Return(_a0 *devices.Device, _a1 error) *MockDevicesServiceClient_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDevicesServiceClient_Get_Call) RunAndReturn(run func(context.Context, *devices.Device, ...grpc.CallOption) (*devices.Device, error)) *MockDevicesServiceClient_Get_Call {
	_c.Call.Return(run)
	return _c
}

// GetByFingerprint provides a mock function with given fields: ctx, in, opts
func (_m *MockDevicesServiceClient) GetByFingerprint(ctx context.Context, in *devices.GetByFingerprintRequest, opts ...grpc.CallOption) (*devices.Device, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for GetByFingerprint")
	}

	var r0 *devices.Device
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *devices.GetByFingerprintRequest, ...grpc.CallOption) (*devices.Device, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *devices.GetByFingerprintRequest, ...grpc.CallOption) *devices.Device); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*devices.Device)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *devices.GetByFingerprintRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDevicesServiceClient_GetByFingerprint_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByFingerprint'
type MockDevicesServiceClient_GetByFingerprint_Call struct {
	*mock.Call
}

// GetByFingerprint is a helper method to define mock.On call
//   - ctx context.Context
//   - in *devices.GetByFingerprintRequest
//   - opts ...grpc.CallOption
func (_e *MockDevicesServiceClient_Expecter) GetByFingerprint(ctx interface{}, in interface{}, opts ...interface{}) *MockDevicesServiceClient_GetByFingerprint_Call {
	return &MockDevicesServiceClient_GetByFingerprint_Call{Call: _e.mock.On("GetByFingerprint",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MockDevicesServiceClient_GetByFingerprint_Call) Run(run func(ctx context.Context, in *devices.GetByFingerprintRequest, opts ...grpc.CallOption)) *MockDevicesServiceClient_GetByFingerprint_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*devices.GetByFingerprintRequest), variadicArgs...)
	})
	return _c
}

func (_c *MockDevicesServiceClient_GetByFingerprint_Call) Return(_a0 *devices.Device, _a1 error) *MockDevicesServiceClient_GetByFingerprint_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDevicesServiceClient_GetByFingerprint_Call) RunAndReturn(run func(context.Context, *devices.GetByFingerprintRequest, ...grpc.CallOption) (*devices.Device, error)) *MockDevicesServiceClient_GetByFingerprint_Call {
	_c.Call.Return(run)
	return _c
}

// GetByToken provides a mock function with given fields: ctx, in, opts
func (_m *MockDevicesServiceClient) GetByToken(ctx context.Context, in *devices.Device, opts ...grpc.CallOption) (*devices.Device, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for GetByToken")
	}

	var r0 *devices.Device
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *devices.Device, ...grpc.CallOption) (*devices.Device, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *devices.Device, ...grpc.CallOption) *devices.Device); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*devices.Device)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *devices.Device, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDevicesServiceClient_GetByToken_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByToken'
type MockDevicesServiceClient_GetByToken_Call struct {
	*mock.Call
}

// GetByToken is a helper method to define mock.On call
//   - ctx context.Context
//   - in *devices.Device
//   - opts ...grpc.CallOption
func (_e *MockDevicesServiceClient_Expecter) GetByToken(ctx interface{}, in interface{}, opts ...interface{}) *MockDevicesServiceClient_GetByToken_Call {
	return &MockDevicesServiceClient_GetByToken_Call{Call: _e.mock.On("GetByToken",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MockDevicesServiceClient_GetByToken_Call) Run(run func(ctx context.Context, in *devices.Device, opts ...grpc.CallOption)) *MockDevicesServiceClient_GetByToken_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*devices.Device), variadicArgs...)
	})
	return _c
}

func (_c *MockDevicesServiceClient_GetByToken_Call) Return(_a0 *devices.Device, _a1 error) *MockDevicesServiceClient_GetByToken_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDevicesServiceClient_GetByToken_Call) RunAndReturn(run func(context.Context, *devices.Device, ...grpc.CallOption) (*devices.Device, error)) *MockDevicesServiceClient_GetByToken_Call {
	_c.Call.Return(run)
	return _c
}

// Join provides a mock function with given fields: ctx, in, opts
func (_m *MockDevicesServiceClient) Join(ctx context.Context, in *node.JoinGeneralRequest, opts ...grpc.CallOption) (*access.Node, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Join")
	}

	var r0 *access.Node
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *node.JoinGeneralRequest, ...grpc.CallOption) (*access.Node, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *node.JoinGeneralRequest, ...grpc.CallOption) *access.Node); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*access.Node)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *node.JoinGeneralRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDevicesServiceClient_Join_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Join'
type MockDevicesServiceClient_Join_Call struct {
	*mock.Call
}

// Join is a helper method to define mock.On call
//   - ctx context.Context
//   - in *node.JoinGeneralRequest
//   - opts ...grpc.CallOption
func (_e *MockDevicesServiceClient_Expecter) Join(ctx interface{}, in interface{}, opts ...interface{}) *MockDevicesServiceClient_Join_Call {
	return &MockDevicesServiceClient_Join_Call{Call: _e.mock.On("Join",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MockDevicesServiceClient_Join_Call) Run(run func(ctx context.Context, in *node.JoinGeneralRequest, opts ...grpc.CallOption)) *MockDevicesServiceClient_Join_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*node.JoinGeneralRequest), variadicArgs...)
	})
	return _c
}

func (_c *MockDevicesServiceClient_Join_Call) Return(_a0 *access.Node, _a1 error) *MockDevicesServiceClient_Join_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDevicesServiceClient_Join_Call) RunAndReturn(run func(context.Context, *node.JoinGeneralRequest, ...grpc.CallOption) (*access.Node, error)) *MockDevicesServiceClient_Join_Call {
	_c.Call.Return(run)
	return _c
}

// Joins provides a mock function with given fields: ctx, in, opts
func (_m *MockDevicesServiceClient) Joins(ctx context.Context, in *devices.Device, opts ...grpc.CallOption) (*access.Nodes, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Joins")
	}

	var r0 *access.Nodes
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *devices.Device, ...grpc.CallOption) (*access.Nodes, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *devices.Device, ...grpc.CallOption) *access.Nodes); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*access.Nodes)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *devices.Device, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDevicesServiceClient_Joins_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Joins'
type MockDevicesServiceClient_Joins_Call struct {
	*mock.Call
}

// Joins is a helper method to define mock.On call
//   - ctx context.Context
//   - in *devices.Device
//   - opts ...grpc.CallOption
func (_e *MockDevicesServiceClient_Expecter) Joins(ctx interface{}, in interface{}, opts ...interface{}) *MockDevicesServiceClient_Joins_Call {
	return &MockDevicesServiceClient_Joins_Call{Call: _e.mock.On("Joins",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MockDevicesServiceClient_Joins_Call) Run(run func(ctx context.Context, in *devices.Device, opts ...grpc.CallOption)) *MockDevicesServiceClient_Joins_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*devices.Device), variadicArgs...)
	})
	return _c
}

func (_c *MockDevicesServiceClient_Joins_Call) Return(_a0 *access.Nodes, _a1 error) *MockDevicesServiceClient_Joins_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDevicesServiceClient_Joins_Call) RunAndReturn(run func(context.Context, *devices.Device, ...grpc.CallOption) (*access.Nodes, error)) *MockDevicesServiceClient_Joins_Call {
	_c.Call.Return(run)
	return _c
}

// List provides a mock function with given fields: ctx, in, opts
func (_m *MockDevicesServiceClient) List(ctx context.Context, in *node.QueryRequest, opts ...grpc.CallOption) (*devices.Devices, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for List")
	}

	var r0 *devices.Devices
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *node.QueryRequest, ...grpc.CallOption) (*devices.Devices, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *node.QueryRequest, ...grpc.CallOption) *devices.Devices); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*devices.Devices)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *node.QueryRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDevicesServiceClient_List_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'List'
type MockDevicesServiceClient_List_Call struct {
	*mock.Call
}

// List is a helper method to define mock.On call
//   - ctx context.Context
//   - in *node.QueryRequest
//   - opts ...grpc.CallOption
func (_e *MockDevicesServiceClient_Expecter) List(ctx interface{}, in interface{}, opts ...interface{}) *MockDevicesServiceClient_List_Call {
	return &MockDevicesServiceClient_List_Call{Call: _e.mock.On("List",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MockDevicesServiceClient_List_Call) Run(run func(ctx context.Context, in *node.QueryRequest, opts ...grpc.CallOption)) *MockDevicesServiceClient_List_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*node.QueryRequest), variadicArgs...)
	})
	return _c
}

func (_c *MockDevicesServiceClient_List_Call) Return(_a0 *devices.Devices, _a1 error) *MockDevicesServiceClient_List_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDevicesServiceClient_List_Call) RunAndReturn(run func(context.Context, *node.QueryRequest, ...grpc.CallOption) (*devices.Devices, error)) *MockDevicesServiceClient_List_Call {
	_c.Call.Return(run)
	return _c
}

// MakeDevicesToken provides a mock function with given fields: ctx, in, opts
func (_m *MockDevicesServiceClient) MakeDevicesToken(ctx context.Context, in *node.DevicesTokenRequest, opts ...grpc.CallOption) (*node.TokenResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for MakeDevicesToken")
	}

	var r0 *node.TokenResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *node.DevicesTokenRequest, ...grpc.CallOption) (*node.TokenResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *node.DevicesTokenRequest, ...grpc.CallOption) *node.TokenResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*node.TokenResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *node.DevicesTokenRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDevicesServiceClient_MakeDevicesToken_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'MakeDevicesToken'
type MockDevicesServiceClient_MakeDevicesToken_Call struct {
	*mock.Call
}

// MakeDevicesToken is a helper method to define mock.On call
//   - ctx context.Context
//   - in *node.DevicesTokenRequest
//   - opts ...grpc.CallOption
func (_e *MockDevicesServiceClient_Expecter) MakeDevicesToken(ctx interface{}, in interface{}, opts ...interface{}) *MockDevicesServiceClient_MakeDevicesToken_Call {
	return &MockDevicesServiceClient_MakeDevicesToken_Call{Call: _e.mock.On("MakeDevicesToken",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MockDevicesServiceClient_MakeDevicesToken_Call) Run(run func(ctx context.Context, in *node.DevicesTokenRequest, opts ...grpc.CallOption)) *MockDevicesServiceClient_MakeDevicesToken_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*node.DevicesTokenRequest), variadicArgs...)
	})
	return _c
}

func (_c *MockDevicesServiceClient_MakeDevicesToken_Call) Return(_a0 *node.TokenResponse, _a1 error) *MockDevicesServiceClient_MakeDevicesToken_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDevicesServiceClient_MakeDevicesToken_Call) RunAndReturn(run func(context.Context, *node.DevicesTokenRequest, ...grpc.CallOption) (*node.TokenResponse, error)) *MockDevicesServiceClient_MakeDevicesToken_Call {
	_c.Call.Return(run)
	return _c
}

// Move provides a mock function with given fields: ctx, in, opts
func (_m *MockDevicesServiceClient) Move(ctx context.Context, in *node.MoveRequest, opts ...grpc.CallOption) (*node.EmptyMessage, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Move")
	}

	var r0 *node.EmptyMessage
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *node.MoveRequest, ...grpc.CallOption) (*node.EmptyMessage, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *node.MoveRequest, ...grpc.CallOption) *node.EmptyMessage); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*node.EmptyMessage)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *node.MoveRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDevicesServiceClient_Move_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Move'
type MockDevicesServiceClient_Move_Call struct {
	*mock.Call
}

// Move is a helper method to define mock.On call
//   - ctx context.Context
//   - in *node.MoveRequest
//   - opts ...grpc.CallOption
func (_e *MockDevicesServiceClient_Expecter) Move(ctx interface{}, in interface{}, opts ...interface{}) *MockDevicesServiceClient_Move_Call {
	return &MockDevicesServiceClient_Move_Call{Call: _e.mock.On("Move",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MockDevicesServiceClient_Move_Call) Run(run func(ctx context.Context, in *node.MoveRequest, opts ...grpc.CallOption)) *MockDevicesServiceClient_Move_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*node.MoveRequest), variadicArgs...)
	})
	return _c
}

func (_c *MockDevicesServiceClient_Move_Call) Return(_a0 *node.EmptyMessage, _a1 error) *MockDevicesServiceClient_Move_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDevicesServiceClient_Move_Call) RunAndReturn(run func(context.Context, *node.MoveRequest, ...grpc.CallOption) (*node.EmptyMessage, error)) *MockDevicesServiceClient_Move_Call {
	_c.Call.Return(run)
	return _c
}

// PatchConfig provides a mock function with given fields: ctx, in, opts
func (_m *MockDevicesServiceClient) PatchConfig(ctx context.Context, in *devices.Device, opts ...grpc.CallOption) (*devices.Device, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for PatchConfig")
	}

	var r0 *devices.Device
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *devices.Device, ...grpc.CallOption) (*devices.Device, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *devices.Device, ...grpc.CallOption) *devices.Device); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*devices.Device)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *devices.Device, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDevicesServiceClient_PatchConfig_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PatchConfig'
type MockDevicesServiceClient_PatchConfig_Call struct {
	*mock.Call
}

// PatchConfig is a helper method to define mock.On call
//   - ctx context.Context
//   - in *devices.Device
//   - opts ...grpc.CallOption
func (_e *MockDevicesServiceClient_Expecter) PatchConfig(ctx interface{}, in interface{}, opts ...interface{}) *MockDevicesServiceClient_PatchConfig_Call {
	return &MockDevicesServiceClient_PatchConfig_Call{Call: _e.mock.On("PatchConfig",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MockDevicesServiceClient_PatchConfig_Call) Run(run func(ctx context.Context, in *devices.Device, opts ...grpc.CallOption)) *MockDevicesServiceClient_PatchConfig_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*devices.Device), variadicArgs...)
	})
	return _c
}

func (_c *MockDevicesServiceClient_PatchConfig_Call) Return(_a0 *devices.Device, _a1 error) *MockDevicesServiceClient_PatchConfig_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDevicesServiceClient_PatchConfig_Call) RunAndReturn(run func(context.Context, *devices.Device, ...grpc.CallOption) (*devices.Device, error)) *MockDevicesServiceClient_PatchConfig_Call {
	_c.Call.Return(run)
	return _c
}

// Toggle provides a mock function with given fields: ctx, in, opts
func (_m *MockDevicesServiceClient) Toggle(ctx context.Context, in *devices.Device, opts ...grpc.CallOption) (*devices.Device, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Toggle")
	}

	var r0 *devices.Device
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *devices.Device, ...grpc.CallOption) (*devices.Device, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *devices.Device, ...grpc.CallOption) *devices.Device); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*devices.Device)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *devices.Device, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDevicesServiceClient_Toggle_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Toggle'
type MockDevicesServiceClient_Toggle_Call struct {
	*mock.Call
}

// Toggle is a helper method to define mock.On call
//   - ctx context.Context
//   - in *devices.Device
//   - opts ...grpc.CallOption
func (_e *MockDevicesServiceClient_Expecter) Toggle(ctx interface{}, in interface{}, opts ...interface{}) *MockDevicesServiceClient_Toggle_Call {
	return &MockDevicesServiceClient_Toggle_Call{Call: _e.mock.On("Toggle",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MockDevicesServiceClient_Toggle_Call) Run(run func(ctx context.Context, in *devices.Device, opts ...grpc.CallOption)) *MockDevicesServiceClient_Toggle_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*devices.Device), variadicArgs...)
	})
	return _c
}

func (_c *MockDevicesServiceClient_Toggle_Call) Return(_a0 *devices.Device, _a1 error) *MockDevicesServiceClient_Toggle_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDevicesServiceClient_Toggle_Call) RunAndReturn(run func(context.Context, *devices.Device, ...grpc.CallOption) (*devices.Device, error)) *MockDevicesServiceClient_Toggle_Call {
	_c.Call.Return(run)
	return _c
}

// ToggleBasic provides a mock function with given fields: ctx, in, opts
func (_m *MockDevicesServiceClient) ToggleBasic(ctx context.Context, in *devices.Device, opts ...grpc.CallOption) (*devices.Device, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for ToggleBasic")
	}

	var r0 *devices.Device
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *devices.Device, ...grpc.CallOption) (*devices.Device, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *devices.Device, ...grpc.CallOption) *devices.Device); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*devices.Device)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *devices.Device, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDevicesServiceClient_ToggleBasic_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ToggleBasic'
type MockDevicesServiceClient_ToggleBasic_Call struct {
	*mock.Call
}

// ToggleBasic is a helper method to define mock.On call
//   - ctx context.Context
//   - in *devices.Device
//   - opts ...grpc.CallOption
func (_e *MockDevicesServiceClient_Expecter) ToggleBasic(ctx interface{}, in interface{}, opts ...interface{}) *MockDevicesServiceClient_ToggleBasic_Call {
	return &MockDevicesServiceClient_ToggleBasic_Call{Call: _e.mock.On("ToggleBasic",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MockDevicesServiceClient_ToggleBasic_Call) Run(run func(ctx context.Context, in *devices.Device, opts ...grpc.CallOption)) *MockDevicesServiceClient_ToggleBasic_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*devices.Device), variadicArgs...)
	})
	return _c
}

func (_c *MockDevicesServiceClient_ToggleBasic_Call) Return(_a0 *devices.Device, _a1 error) *MockDevicesServiceClient_ToggleBasic_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDevicesServiceClient_ToggleBasic_Call) RunAndReturn(run func(context.Context, *devices.Device, ...grpc.CallOption) (*devices.Device, error)) *MockDevicesServiceClient_ToggleBasic_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: ctx, in, opts
func (_m *MockDevicesServiceClient) Update(ctx context.Context, in *devices.Device, opts ...grpc.CallOption) (*devices.Device, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 *devices.Device
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *devices.Device, ...grpc.CallOption) (*devices.Device, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *devices.Device, ...grpc.CallOption) *devices.Device); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*devices.Device)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *devices.Device, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDevicesServiceClient_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type MockDevicesServiceClient_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - ctx context.Context
//   - in *devices.Device
//   - opts ...grpc.CallOption
func (_e *MockDevicesServiceClient_Expecter) Update(ctx interface{}, in interface{}, opts ...interface{}) *MockDevicesServiceClient_Update_Call {
	return &MockDevicesServiceClient_Update_Call{Call: _e.mock.On("Update",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MockDevicesServiceClient_Update_Call) Run(run func(ctx context.Context, in *devices.Device, opts ...grpc.CallOption)) *MockDevicesServiceClient_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*devices.Device), variadicArgs...)
	})
	return _c
}

func (_c *MockDevicesServiceClient_Update_Call) Return(_a0 *devices.Device, _a1 error) *MockDevicesServiceClient_Update_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDevicesServiceClient_Update_Call) RunAndReturn(run func(context.Context, *devices.Device, ...grpc.CallOption) (*devices.Device, error)) *MockDevicesServiceClient_Update_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockDevicesServiceClient creates a new instance of MockDevicesServiceClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockDevicesServiceClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockDevicesServiceClient {
	mock := &MockDevicesServiceClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
