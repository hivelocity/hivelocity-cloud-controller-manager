/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	swagger "github.com/hivelocity/hivelocity-client-go/client"
	mock "github.com/stretchr/testify/mock"
)

// API is an autogenerated mock type for the API type
type API struct {
	mock.Mock
}

// GetBareMetalDeviceIdResource provides a mock function with given fields: deviceId
func (_m *API) GetBareMetalDeviceIdResource(deviceId int32) (*swagger.BareMetalDevice, error) {
	ret := _m.Called(deviceId)

	var r0 *swagger.BareMetalDevice
	if rf, ok := ret.Get(0).(func(int32) *swagger.BareMetalDevice); ok {
		r0 = rf(deviceId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*swagger.BareMetalDevice)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int32) error); ok {
		r1 = rf(deviceId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewAPI interface {
	mock.TestingT
	Cleanup(func())
}

// NewAPI creates a new instance of API. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewAPI(t mockConstructorTestingTNewAPI) *API {
	mock := &API{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
