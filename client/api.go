/*
Copyright 2023 The Kubernetes Authors.

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

package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	hv "github.com/hivelocity/hivelocity-client-go/client"
)

// API is a wrapper of hv.APIClient. This way mocking (for tests)
// is easier.
type API interface {
	GetBareMetalDeviceIdResource(deviceId int32) (*hv.BareMetalDevice, error)
}

// RealAPI implements API.
type RealAPI struct {
	Client *hv.APIClient
}

var _ API = (*RealAPI)(nil)

var ErrNoSuchDevice = errors.New("no such device")

func (api *RealAPI) GetBareMetalDeviceIdResource(deviceId int32) (*hv.BareMetalDevice, error) {
	device, response, err := api.Client.BareMetalDevicesApi.GetBareMetalDeviceIdResource(
		context.Background(), deviceId, nil)
	if err != nil {
		err, ok := err.(hv.GenericSwaggerError)
		if !ok {
			return nil, fmt.Errorf(
				"unknown error during GetBareMetalDeviceIdResource StatusCode %d deviceId %q. %w",
				response.StatusCode, deviceId, err)
		}
		var result struct {
			Code    int
			Message string
		}
		if err2 := json.Unmarshal(err.Body(), &result); err2 != nil {
			return nil, fmt.Errorf(
				"GetBareMetalDeviceIdResource failed to parse response body %s. StatusCode %d deviceId %q. %w",
				err.Body(),
				response.StatusCode, deviceId, err2)
		}

		if result.Message == "Device not found" {
			return nil, ErrNoSuchDevice
		}
		return nil, fmt.Errorf("GetBareMetalDeviceIdResource failed with %d. deviceId %q. %w",
			response.StatusCode, deviceId, err)
	}
	return &device, nil
}
