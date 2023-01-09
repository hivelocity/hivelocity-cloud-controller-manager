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

package hivelocity

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	hv "github.com/hivelocity/hivelocity-client-go/client"
	v1 "k8s.io/api/core/v1"
)

var NoSuchDeviceError = errors.New("No such device")

type RemoteAPI interface{
	GetBareMetalDeviceIdResource (client *hv.APIClient, deviceId int32) (*hv.BareMetalDevice, error)
}


type RealRemoteAPI struct {}

var _ RemoteAPI = (*RealRemoteAPI)(nil)

func (remote *RealRemoteAPI) GetBareMetalDeviceIdResource (client *hv.APIClient, deviceId int32) (*hv.BareMetalDevice, error){
	device, response, err := client.BareMetalDevicesApi.GetBareMetalDeviceIdResource(
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
			return nil, NoSuchDeviceError
		}
		return nil, fmt.Errorf("GetBareMetalDeviceIdResource failed with %d. deviceId %q. %w",
			response.StatusCode, deviceId, err)
	}
	return &device, nil
}

// HVInstancesV2 implements cloudprovider.InstanceV2
type HVInstancesV2 struct {
	Client *hv.APIClient
	Remote RemoteAPI
}

func GetHivelocityDeviceIdFromNode(node *v1.Node) (int32, error) {
	deviceId, err := strconv.ParseInt(node.Spec.ProviderID, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("failed to convert node.Spec.ProviderID %q to int32",
			node.Spec.ProviderID)
	}
	return int32(deviceId), nil
}

// InstanceExists returns true if the instance for the given node exists according to the cloud provider.
// Use the node.name or node.spec.providerID field to find the node in the cloud provider.
func (i2 *HVInstancesV2) InstanceExists(ctx context.Context, node *v1.Node) (bool, error) {
	deviceID, err := GetHivelocityDeviceIdFromNode(node)
	if err != nil {
		return false, err
	}
	_, err = i2.Remote.GetBareMetalDeviceIdResource(i2.Client, deviceID)
	if err == NoSuchDeviceError {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// InstanceShutdown returns true if the instance is shutdown according to the cloud provider.
// Use the node.name or node.spec.providerID field to find the node in the cloud provider.
func (i2 *HVInstancesV2) InstanceShutdown(ctx context.Context, node *v1.Node) (bool, error) {
	deviceID, err := GetHivelocityDeviceIdFromNode(node)
	if err != nil {
		return false, err
	}
	device, _, err := i2.Client.BareMetalDevicesApi.GetBareMetalDeviceIdResource(ctx, deviceID, nil)
	if err != nil {
		return false, err
	}
	switch device.PowerStatus {
	case "ON":
		return false, nil
	case "OFF":
		return true, nil
	default:
		return false, fmt.Errorf("device with ID %q has unknown PowerStatus %q", deviceID, device.PowerStatus)
	}
}