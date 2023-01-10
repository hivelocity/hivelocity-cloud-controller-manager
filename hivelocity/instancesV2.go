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
	"fmt"
	"strconv"

	hv "github.com/hivelocity/hivelocity-client-go/client"
	"github.com/hivelocity/hivelocity-cloud-controller-manager/client"
	corev1 "k8s.io/api/core/v1"
	cloudprovider "k8s.io/cloud-provider"
)


// HVInstancesV2 implements cloudprovider.InstanceV2
type HVInstancesV2 struct {
	Client *hv.APIClient
	API    client.API
}

var _ cloudprovider.InstancesV2 = &HVInstancesV2{}

func GetHivelocityDeviceIdFromNode(node *corev1.Node) (int32, error) {
	deviceId, err := strconv.ParseInt(node.Spec.ProviderID, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("failed to convert node.Spec.ProviderID %q to int32",
			node.Spec.ProviderID)
	}
	return int32(deviceId), nil
}

// InstanceExists returns true if the instance for the given node exists according to the cloud provider.
// Use the node.name or node.spec.providerID field to find the node in the cloud provider.
func (i2 *HVInstancesV2) InstanceExists(ctx context.Context, node *corev1.Node) (bool, error) {
	deviceId, err := GetHivelocityDeviceIdFromNode(node)
	if err != nil {
		return false, err
	}
	_, err = i2.API.GetBareMetalDeviceIdResource(deviceId)
	if err == client.ErrNoSuchDevice {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// InstanceShutdown returns true if the instance is shutdown according to the cloud provider.
// Use the node.name or node.spec.providerID field to find the node in the cloud provider.
func (i2 *HVInstancesV2) InstanceShutdown(ctx context.Context, node *corev1.Node) (bool, error) {
	deviceId, err := GetHivelocityDeviceIdFromNode(node)
	if err != nil {
		return false, err
	}
	device, _, err := i2.Client.BareMetalDevicesApi.GetBareMetalDeviceIdResource(ctx, deviceId, nil)
	if err != nil {
		return false, err
	}
	switch device.PowerStatus {
	case "ON":
		return false, nil
	case "OFF":
		return true, nil
	default:
		return false, fmt.Errorf("device with ID %q has unknown PowerStatus %q", deviceId, device.PowerStatus)
	}
}

func (i2 *HVInstancesV2) InstanceMetadata(ctx context.Context, node *corev1.Node) (*cloudprovider.InstanceMetadata, error) {
	deviceId, err := GetHivelocityDeviceIdFromNode(node)
	if err != nil {
		return nil, err
	}
	device, err := i2.API.GetBareMetalDeviceIdResource(deviceId)
	if err != nil {
		return nil, err
	}

	addr := corev1.NodeAddress{
		Type:    "ExternalIP",
		Address: device.PrimaryIp,
	}

	var metaData = cloudprovider.InstanceMetadata{
		ProviderID:    strconv.Itoa(int(deviceId)),
		InstanceType:  device.ProductName,
		NodeAddresses: []corev1.NodeAddress{addr},
		Zone:          device.LocationName,
		Region:        "",
	}
	return &metaData, nil
}
