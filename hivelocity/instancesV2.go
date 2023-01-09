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
	"fmt"
	"strconv"

	hv "github.com/hivelocity/hivelocity-client-go/client"
	v1 "k8s.io/api/core/v1"
)

// hvInstancesV2 implements cloudprovider.InstanceV2
type hvInstancesV2 struct {
	client *hv.APIClient
}

func getHivelocityDeviceIdFromNode(node *v1.Node) (int32, error) {
	deviceId, err := strconv.ParseInt(node.Spec.ProviderID, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("failed to convert node.Spec.ProviderID %q to int32",
			node.Spec.ProviderID)
	}
	return int32(deviceId), nil
}

// InstanceExists returns true if the instance for the given node exists according to the cloud provider.
// Use the node.name or node.spec.providerID field to find the node in the cloud provider.
func (i2 *hvInstancesV2) InstanceExists(ctx context.Context, node *v1.Node) (bool, error) {
	deviceID, err := getHivelocityDeviceIdFromNode(node)
	if err != nil {
		return false, err
	}
	_, response, err := i2.client.BareMetalDevicesApi.GetBareMetalDeviceIdResource(ctx, deviceID, nil)
	if err != nil {
		err, ok := err.(hv.GenericSwaggerError)
		if !ok {
			return false, fmt.Errorf(
				"unknown error during GetBareMetalDeviceIdResource StatusCode %d node.Spec.ProviderID %q. %w",
				response.StatusCode, node.Spec.ProviderID, err)
		}
		var result struct {
			Code int
			Message string
		}
		if err2 := json.Unmarshal(err.Body(), &result); err2 != nil {
			return false, fmt.Errorf(
				"GetBareMetalDeviceIdResource failed to parse response body %s. StatusCode %d node.Spec.ProviderID %q. %w",
				err.Body(),
				response.StatusCode, node.Spec.ProviderID, err2)
		}

		if result.Message == "Device not found" {
			return false, nil
		}
		return false, fmt.Errorf("GetBareMetalDeviceIdResource failed with %d. node.Spec.ProviderID %q. %w",
			response.StatusCode, node.Spec.ProviderID, err)
	}
	return true, nil
}

// InstanceShutdown returns true if the instance is shutdown according to the cloud provider.
// Use the node.name or node.spec.providerID field to find the node in the cloud provider.
func (i2 *hvInstancesV2) InstanceShutdown(ctx context.Context, node *v1.Node) (bool, error) {
	deviceID, err := getHivelocityDeviceIdFromNode(node)
	if err != nil {
		return false, err
	}
	device, _, err := i2.client.BareMetalDevicesApi.GetBareMetalDeviceIdResource(ctx, deviceID, nil)
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
