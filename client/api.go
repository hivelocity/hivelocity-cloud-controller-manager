package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	hv "github.com/hivelocity/hivelocity-client-go/client"
)

type API interface {
	GetBareMetalDeviceIdResource(deviceId int32) (*hv.BareMetalDevice, error)
}

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
