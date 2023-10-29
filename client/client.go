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
	"net/http"
	"regexp"
	"runtime/debug"
	"strings"

	"github.com/go-logr/logr"
	hv "github.com/hivelocity/hivelocity-client-go/client"
	ctrl "sigs.k8s.io/controller-runtime"
)

// Interface is a wrapper of hv.APIClient. In this way, mocking (for tests)
// is easier.
type Interface interface {
	GetBareMetalDevice(ctx context.Context, deviceID int32) (*hv.BareMetalDevice, error)
	ListDevices(context.Context) ([]hv.BareMetalDevice, error)
}

// Client implements the Interface interface.
type Client struct {
	client *hv.APIClient
}

var _ Interface = (*Client)(nil)

// ErrNoSuchDevice means that no device was found via the Hivelocity API.
var ErrNoSuchDevice = errors.New("no such device")

// LoggingTransport is a struct for creating new logger for Hivelocity API.
type LoggingTransport struct {
	roundTripper http.RoundTripper
	log          logr.Logger
}

var replaceHex = regexp.MustCompile(`0x[0123456789abcdef]+`)

// RoundTrip is used for logging api calls to Hivelocity API.
func (lt *LoggingTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	stack := replaceHex.ReplaceAllString(string(debug.Stack()), "0xX")
	stack = strings.ReplaceAll(stack, "\n", "\\n")

	resp, err = lt.roundTripper.RoundTrip(req)
	if err != nil {
		lt.log.V(1).Info("hivelocity API. Error.", "err", err, "method", req.Method, "url", req.URL, "stack", stack)
		return resp, fmt.Errorf("failed to RoundTrip: %w", err)
	}
	lt.log.V(1).Info("hivelocity API called.", "statusCode", resp.StatusCode, "method", req.Method, "url", req.URL, "stack", stack)
	return resp, nil
}

// NewClient creates a struct which implements the Client interface.
func NewClient(apiKey string) *Client {
	config := hv.NewConfiguration()
	config.HTTPClient = &http.Client{
		Transport: &LoggingTransport{
			roundTripper: http.DefaultTransport,
			log:          ctrl.Log.WithName("Hivelocity-api"),
		},
	}

	config.AddDefaultHeader("X-API-KEY", apiKey)
	apiClient := hv.NewAPIClient(config)
	return &Client{client: apiClient}
}

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

	if err := response.Body.Close(); err != nil {
		return nil, fmt.Errorf("failed to close response body: %w", err)
	}

	return nil, fmt.Errorf(
		"[GetBareMetalDevice] GetBareMetalDeviceIdResource failed. StatusCode %d, deviceID %q: %w",
		response.StatusCode,
		deviceID,
		err,
	)
}

// ListDevices lists all devices via Hivelocity API.
func (c *Client) ListDevices(ctx context.Context) ([]hv.BareMetalDevice, error) {
	devices, response, err := c.client.BareMetalDevicesApi.GetBareMetalDeviceResource(ctx, nil)
	if err == nil {
		return devices, nil
	}

	var swaggerErr *hv.GenericSwaggerError
	if !errors.As(err, swaggerErr) {
		return nil, fmt.Errorf(
			"[ListDevices] unknown error during GetBareMetalDeviceResource. StatusCode %d: %w",
			response.StatusCode,
			err,
		)
	}

	if err := response.Body.Close(); err != nil {
		return nil, fmt.Errorf("failed to close response body: %w", err)
	}

	return nil, fmt.Errorf(
		"[ListDevices] GetBareMetalDeviceResource failed. StatusCode %d: %w",
		response.StatusCode,
		err,
	)
}
