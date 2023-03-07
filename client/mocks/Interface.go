// Code generated by mockery v2.21.1. DO NOT EDIT.

package mocks

import (
	context "context"

	swagger "github.com/hivelocity/hivelocity-client-go/client"
	mock "github.com/stretchr/testify/mock"
)

// Interface is an autogenerated mock type for the Interface type
type Interface struct {
	mock.Mock
}

// GetBareMetalDevice provides a mock function with given fields: ctx, deviceID
func (_m *Interface) GetBareMetalDevice(ctx context.Context, deviceID int32) (*swagger.BareMetalDevice, error) {
	ret := _m.Called(ctx, deviceID)

	var r0 *swagger.BareMetalDevice
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int32) (*swagger.BareMetalDevice, error)); ok {
		return rf(ctx, deviceID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int32) *swagger.BareMetalDevice); ok {
		r0 = rf(ctx, deviceID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*swagger.BareMetalDevice)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int32) error); ok {
		r1 = rf(ctx, deviceID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewInterface interface {
	mock.TestingT
	Cleanup(func())
}

// NewInterface creates a new instance of Interface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewInterface(t mockConstructorTestingTNewInterface) *Interface {
	mock := &Interface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
