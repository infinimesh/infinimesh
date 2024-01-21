// Code generated by mockery v2.40.1. DO NOT EDIT.

package auth

import (
	jwt "github.com/golang-jwt/jwt/v4"
	mock "github.com/stretchr/testify/mock"
)

// MockJWTHandler is an autogenerated mock type for the JWTHandler type
type MockJWTHandler struct {
	mock.Mock
}

type MockJWTHandler_Expecter struct {
	mock *mock.Mock
}

func (_m *MockJWTHandler) EXPECT() *MockJWTHandler_Expecter {
	return &MockJWTHandler_Expecter{mock: &_m.Mock}
}

// Parse provides a mock function with given fields: tokenString, keyFunc, options
func (_m *MockJWTHandler) Parse(tokenString string, keyFunc jwt.Keyfunc, options ...jwt.ParserOption) (*jwt.Token, error) {
	_va := make([]interface{}, len(options))
	for _i := range options {
		_va[_i] = options[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, tokenString, keyFunc)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Parse")
	}

	var r0 *jwt.Token
	var r1 error
	if rf, ok := ret.Get(0).(func(string, jwt.Keyfunc, ...jwt.ParserOption) (*jwt.Token, error)); ok {
		return rf(tokenString, keyFunc, options...)
	}
	if rf, ok := ret.Get(0).(func(string, jwt.Keyfunc, ...jwt.ParserOption) *jwt.Token); ok {
		r0 = rf(tokenString, keyFunc, options...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*jwt.Token)
		}
	}

	if rf, ok := ret.Get(1).(func(string, jwt.Keyfunc, ...jwt.ParserOption) error); ok {
		r1 = rf(tokenString, keyFunc, options...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockJWTHandler_Parse_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Parse'
type MockJWTHandler_Parse_Call struct {
	*mock.Call
}

// Parse is a helper method to define mock.On call
//   - tokenString string
//   - keyFunc jwt.Keyfunc
//   - options ...jwt.ParserOption
func (_e *MockJWTHandler_Expecter) Parse(tokenString interface{}, keyFunc interface{}, options ...interface{}) *MockJWTHandler_Parse_Call {
	return &MockJWTHandler_Parse_Call{Call: _e.mock.On("Parse",
		append([]interface{}{tokenString, keyFunc}, options...)...)}
}

func (_c *MockJWTHandler_Parse_Call) Run(run func(tokenString string, keyFunc jwt.Keyfunc, options ...jwt.ParserOption)) *MockJWTHandler_Parse_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]jwt.ParserOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(jwt.ParserOption)
			}
		}
		run(args[0].(string), args[1].(jwt.Keyfunc), variadicArgs...)
	})
	return _c
}

func (_c *MockJWTHandler_Parse_Call) Return(_a0 *jwt.Token, _a1 error) *MockJWTHandler_Parse_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockJWTHandler_Parse_Call) RunAndReturn(run func(string, jwt.Keyfunc, ...jwt.ParserOption) (*jwt.Token, error)) *MockJWTHandler_Parse_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockJWTHandler creates a new instance of MockJWTHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockJWTHandler(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockJWTHandler {
	mock := &MockJWTHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
