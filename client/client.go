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

// Package client provides the interfaces to communicate with the
// API of Hivelocity.
// Creating interfaces makes unit testing via mocking possible.
package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	hv "github.com/hivelocity/hivelocity-client-go/client"
)

// Interface is a wrapper of hv.APIClient. In this way, mocking (for tests)
// is easier.
type Interface interface {
	GetBareMetalDevice(ctx context.Context, deviceID int32) (*hv.BareMetalDevice, error)
}

// Client implements the Interface interface.
type Client struct {
	client *hv.APIClient
}

var _ Interface = (*Client)(nil)

// NewClient creates a struct which implements the Client interface.
func NewClient(client *hv.APIClient) *Client {
	return &Client{client: client}
}

// ErrNoSuchDevice means that no device was found via the Hivelocity API.
var ErrNoSuchDevice = errors.New("no such device")

// GetBareMetalDevice returns the device fetched via the Hivelocity API.
func (c *Client) GetBareMetalDevice(
	ctx context.Context,
	deviceID int32,
) (*hv.BareMetalDevice, error) {
	device, response, err := c.client.BareMetalDevicesApi.GetBareMetalDeviceIdResource(
		ctx, deviceID, nil)
	if err == nil {
		return &device, nil
	}

	// Analyze the error that has been returned.

	var swaggerErr *hv.GenericSwaggerError
	if !errors.As(err, swaggerErr) {
		return nil, fmt.Errorf(
			"[GetBareMetalDevice] unknown error during GetBareMetalDeviceIdResource. StatusCode %d, deviceID %q: %w",
			response.StatusCode,
			deviceID,
			err,
		)
	}
	var result struct {
		Code    int
		Message string
	}

	if unmarshalErr := json.Unmarshal(swaggerErr.Body(), &result); unmarshalErr != nil {
		return nil, fmt.Errorf(
			"[GetBareMetalDevice] GetBareMetalDeviceIdResource failed to parse response body %s."+
				"StatusCode %d, deviceID %q: %w",
			swaggerErr.Body(),
			response.StatusCode,
			deviceID,
			unmarshalErr,
		)
	}

	if result.Message == "Device not found" {
		return nil, ErrNoSuchDevice
	}

	return nil, fmt.Errorf(
		"[GetBareMetalDevice] GetBareMetalDeviceIdResource failed. StatusCode %d, deviceID %q: %w",
		response.StatusCode,
		deviceID,
		err,
	)
}
