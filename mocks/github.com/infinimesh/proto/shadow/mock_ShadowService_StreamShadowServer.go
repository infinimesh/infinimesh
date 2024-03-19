// Code generated by mockery v2.42.0. DO NOT EDIT.

package shadow_mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	metadata "google.golang.org/grpc/metadata"

	shadow "github.com/infinimesh/proto/shadow"
)

// MockShadowService_StreamShadowServer is an autogenerated mock type for the ShadowService_StreamShadowServer type
type MockShadowService_StreamShadowServer struct {
	mock.Mock
}

type MockShadowService_StreamShadowServer_Expecter struct {
	mock *mock.Mock
}

func (_m *MockShadowService_StreamShadowServer) EXPECT() *MockShadowService_StreamShadowServer_Expecter {
	return &MockShadowService_StreamShadowServer_Expecter{mock: &_m.Mock}
}

// Context provides a mock function with given fields:
func (_m *MockShadowService_StreamShadowServer) Context() context.Context {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Context")
	}

	var r0 context.Context
	if rf, ok := ret.Get(0).(func() context.Context); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(context.Context)
		}
	}

	return r0
}

// MockShadowService_StreamShadowServer_Context_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Context'
type MockShadowService_StreamShadowServer_Context_Call struct {
	*mock.Call
}

// Context is a helper method to define mock.On call
func (_e *MockShadowService_StreamShadowServer_Expecter) Context() *MockShadowService_StreamShadowServer_Context_Call {
	return &MockShadowService_StreamShadowServer_Context_Call{Call: _e.mock.On("Context")}
}

func (_c *MockShadowService_StreamShadowServer_Context_Call) Run(run func()) *MockShadowService_StreamShadowServer_Context_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockShadowService_StreamShadowServer_Context_Call) Return(_a0 context.Context) *MockShadowService_StreamShadowServer_Context_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockShadowService_StreamShadowServer_Context_Call) RunAndReturn(run func() context.Context) *MockShadowService_StreamShadowServer_Context_Call {
	_c.Call.Return(run)
	return _c
}

// RecvMsg provides a mock function with given fields: m
func (_m *MockShadowService_StreamShadowServer) RecvMsg(m interface{}) error {
	ret := _m.Called(m)

	if len(ret) == 0 {
		panic("no return value specified for RecvMsg")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}) error); ok {
		r0 = rf(m)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockShadowService_StreamShadowServer_RecvMsg_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RecvMsg'
type MockShadowService_StreamShadowServer_RecvMsg_Call struct {
	*mock.Call
}

// RecvMsg is a helper method to define mock.On call
//   - m interface{}
func (_e *MockShadowService_StreamShadowServer_Expecter) RecvMsg(m interface{}) *MockShadowService_StreamShadowServer_RecvMsg_Call {
	return &MockShadowService_StreamShadowServer_RecvMsg_Call{Call: _e.mock.On("RecvMsg", m)}
}

func (_c *MockShadowService_StreamShadowServer_RecvMsg_Call) Run(run func(m interface{})) *MockShadowService_StreamShadowServer_RecvMsg_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(interface{}))
	})
	return _c
}

func (_c *MockShadowService_StreamShadowServer_RecvMsg_Call) Return(_a0 error) *MockShadowService_StreamShadowServer_RecvMsg_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockShadowService_StreamShadowServer_RecvMsg_Call) RunAndReturn(run func(interface{}) error) *MockShadowService_StreamShadowServer_RecvMsg_Call {
	_c.Call.Return(run)
	return _c
}

// Send provides a mock function with given fields: _a0
func (_m *MockShadowService_StreamShadowServer) Send(_a0 *shadow.Shadow) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for Send")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*shadow.Shadow) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockShadowService_StreamShadowServer_Send_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Send'
type MockShadowService_StreamShadowServer_Send_Call struct {
	*mock.Call
}

// Send is a helper method to define mock.On call
//   - _a0 *shadow.Shadow
func (_e *MockShadowService_StreamShadowServer_Expecter) Send(_a0 interface{}) *MockShadowService_StreamShadowServer_Send_Call {
	return &MockShadowService_StreamShadowServer_Send_Call{Call: _e.mock.On("Send", _a0)}
}

func (_c *MockShadowService_StreamShadowServer_Send_Call) Run(run func(_a0 *shadow.Shadow)) *MockShadowService_StreamShadowServer_Send_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*shadow.Shadow))
	})
	return _c
}

func (_c *MockShadowService_StreamShadowServer_Send_Call) Return(_a0 error) *MockShadowService_StreamShadowServer_Send_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockShadowService_StreamShadowServer_Send_Call) RunAndReturn(run func(*shadow.Shadow) error) *MockShadowService_StreamShadowServer_Send_Call {
	_c.Call.Return(run)
	return _c
}

// SendHeader provides a mock function with given fields: _a0
func (_m *MockShadowService_StreamShadowServer) SendHeader(_a0 metadata.MD) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for SendHeader")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(metadata.MD) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockShadowService_StreamShadowServer_SendHeader_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SendHeader'
type MockShadowService_StreamShadowServer_SendHeader_Call struct {
	*mock.Call
}

// SendHeader is a helper method to define mock.On call
//   - _a0 metadata.MD
func (_e *MockShadowService_StreamShadowServer_Expecter) SendHeader(_a0 interface{}) *MockShadowService_StreamShadowServer_SendHeader_Call {
	return &MockShadowService_StreamShadowServer_SendHeader_Call{Call: _e.mock.On("SendHeader", _a0)}
}

func (_c *MockShadowService_StreamShadowServer_SendHeader_Call) Run(run func(_a0 metadata.MD)) *MockShadowService_StreamShadowServer_SendHeader_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(metadata.MD))
	})
	return _c
}

func (_c *MockShadowService_StreamShadowServer_SendHeader_Call) Return(_a0 error) *MockShadowService_StreamShadowServer_SendHeader_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockShadowService_StreamShadowServer_SendHeader_Call) RunAndReturn(run func(metadata.MD) error) *MockShadowService_StreamShadowServer_SendHeader_Call {
	_c.Call.Return(run)
	return _c
}

// SendMsg provides a mock function with given fields: m
func (_m *MockShadowService_StreamShadowServer) SendMsg(m interface{}) error {
	ret := _m.Called(m)

	if len(ret) == 0 {
		panic("no return value specified for SendMsg")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}) error); ok {
		r0 = rf(m)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockShadowService_StreamShadowServer_SendMsg_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SendMsg'
type MockShadowService_StreamShadowServer_SendMsg_Call struct {
	*mock.Call
}

// SendMsg is a helper method to define mock.On call
//   - m interface{}
func (_e *MockShadowService_StreamShadowServer_Expecter) SendMsg(m interface{}) *MockShadowService_StreamShadowServer_SendMsg_Call {
	return &MockShadowService_StreamShadowServer_SendMsg_Call{Call: _e.mock.On("SendMsg", m)}
}

func (_c *MockShadowService_StreamShadowServer_SendMsg_Call) Run(run func(m interface{})) *MockShadowService_StreamShadowServer_SendMsg_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(interface{}))
	})
	return _c
}

func (_c *MockShadowService_StreamShadowServer_SendMsg_Call) Return(_a0 error) *MockShadowService_StreamShadowServer_SendMsg_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockShadowService_StreamShadowServer_SendMsg_Call) RunAndReturn(run func(interface{}) error) *MockShadowService_StreamShadowServer_SendMsg_Call {
	_c.Call.Return(run)
	return _c
}

// SetHeader provides a mock function with given fields: _a0
func (_m *MockShadowService_StreamShadowServer) SetHeader(_a0 metadata.MD) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for SetHeader")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(metadata.MD) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockShadowService_StreamShadowServer_SetHeader_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetHeader'
type MockShadowService_StreamShadowServer_SetHeader_Call struct {
	*mock.Call
}

// SetHeader is a helper method to define mock.On call
//   - _a0 metadata.MD
func (_e *MockShadowService_StreamShadowServer_Expecter) SetHeader(_a0 interface{}) *MockShadowService_StreamShadowServer_SetHeader_Call {
	return &MockShadowService_StreamShadowServer_SetHeader_Call{Call: _e.mock.On("SetHeader", _a0)}
}

func (_c *MockShadowService_StreamShadowServer_SetHeader_Call) Run(run func(_a0 metadata.MD)) *MockShadowService_StreamShadowServer_SetHeader_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(metadata.MD))
	})
	return _c
}

func (_c *MockShadowService_StreamShadowServer_SetHeader_Call) Return(_a0 error) *MockShadowService_StreamShadowServer_SetHeader_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockShadowService_StreamShadowServer_SetHeader_Call) RunAndReturn(run func(metadata.MD) error) *MockShadowService_StreamShadowServer_SetHeader_Call {
	_c.Call.Return(run)
	return _c
}

// SetTrailer provides a mock function with given fields: _a0
func (_m *MockShadowService_StreamShadowServer) SetTrailer(_a0 metadata.MD) {
	_m.Called(_a0)
}

// MockShadowService_StreamShadowServer_SetTrailer_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetTrailer'
type MockShadowService_StreamShadowServer_SetTrailer_Call struct {
	*mock.Call
}

// SetTrailer is a helper method to define mock.On call
//   - _a0 metadata.MD
func (_e *MockShadowService_StreamShadowServer_Expecter) SetTrailer(_a0 interface{}) *MockShadowService_StreamShadowServer_SetTrailer_Call {
	return &MockShadowService_StreamShadowServer_SetTrailer_Call{Call: _e.mock.On("SetTrailer", _a0)}
}

func (_c *MockShadowService_StreamShadowServer_SetTrailer_Call) Run(run func(_a0 metadata.MD)) *MockShadowService_StreamShadowServer_SetTrailer_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(metadata.MD))
	})
	return _c
}

func (_c *MockShadowService_StreamShadowServer_SetTrailer_Call) Return() *MockShadowService_StreamShadowServer_SetTrailer_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockShadowService_StreamShadowServer_SetTrailer_Call) RunAndReturn(run func(metadata.MD)) *MockShadowService_StreamShadowServer_SetTrailer_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockShadowService_StreamShadowServer creates a new instance of MockShadowService_StreamShadowServer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockShadowService_StreamShadowServer(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockShadowService_StreamShadowServer {
	mock := &MockShadowService_StreamShadowServer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
