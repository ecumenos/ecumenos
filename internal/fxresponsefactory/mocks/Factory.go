// Code generated by mockery v2.39.1. DO NOT EDIT.

package mocks

import (
	http "net/http"

	fxresponsefactory "github.com/ecumenos/ecumenos/internal/fxresponsefactory"

	mock "github.com/stretchr/testify/mock"
)

// Factory is an autogenerated mock type for the Factory type
type Factory struct {
	mock.Mock
}

// NewWriter provides a mock function with given fields: rw
func (_m *Factory) NewWriter(rw http.ResponseWriter) fxresponsefactory.Writer {
	ret := _m.Called(rw)

	if len(ret) == 0 {
		panic("no return value specified for NewWriter")
	}

	var r0 fxresponsefactory.Writer
	if rf, ok := ret.Get(0).(func(http.ResponseWriter) fxresponsefactory.Writer); ok {
		r0 = rf(rw)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(fxresponsefactory.Writer)
		}
	}

	return r0
}

// NewFactory creates a new instance of Factory. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewFactory(t interface {
	mock.TestingT
	Cleanup(func())
}) *Factory {
	mock := &Factory{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
