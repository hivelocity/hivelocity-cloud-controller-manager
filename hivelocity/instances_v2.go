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
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hivelocity/hivelocity-cloud-controller-manager/client"
	"github.com/hivelocity/hivelocity-cloud-controller-manager/pkg/hvutils"
	corev1 "k8s.io/api/core/v1"
	cloudprovider "k8s.io/cloud-provider"
)

// HVInstancesV2 implements cloudprovider.InstanceV2.
type HVInstancesV2 struct {
	client client.Interface
}

var _ cloudprovider.InstancesV2 = &HVInstancesV2{}

// NewHVInstanceV2 creates a new HVInstancesV2 struct.
func NewHVInstanceV2(c client.Interface) *HVInstancesV2 {
	return &HVInstancesV2{client: c}
}

// getHivelocityDeviceIDFromNode returns the deviceID from a Node.
// Example: If Node.Spec.ProviderID is "hivelocity://123", then 123
// will be returned.
func getHivelocityDeviceIDFromNode(node *corev1.Node) (int32, error) {
	providerPrefix := providerName + "://"
	if !strings.HasPrefix(node.Spec.ProviderID, providerPrefix) {
		return 0, fmt.Errorf("ProviderID: %q: %w", node.Spec.ProviderID, errMissingProviderPrefix)
	}
	deviceID, err := strconv.ParseInt(strings.TrimPrefix(node.Spec.ProviderID, providerPrefix), 10, 32)
	if err != nil {
		return 0, fmt.Errorf("ParseInt failed. node.Spec.ProviderID %q: %w",
			node.Spec.ProviderID, errFailedToConvertProviderID)
	}
	return int32(deviceID), nil
}

var (
	errMissingProviderPrefix     = fmt.Errorf("missing prefix %q in node.Spec.ProviderID", providerName)
	errFailedToConvertProviderID = fmt.Errorf("failed to convert node.Spec.ProviderID")
)

// InstanceExists returns true if the instance for the given node exists according to the cloud provider.
// Use the node.name or node.spec.providerID field to find the node in the cloud provider.
func (i2 *HVInstancesV2) InstanceExists(ctx context.Context, node *corev1.Node) (bool, error) {
	deviceID, err := getHivelocityDeviceIDFromNode(node)
	if err != nil {
		return false, err
	}
	_, err = i2.client.GetBareMetalDevice(ctx, deviceID)
	if errors.Is(err, client.ErrNoSuchDevice) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// ErrUnknownPowerStatus .
var ErrUnknownPowerStatus = errors.New("unknown PowerStatus")

// InstanceShutdown returns true if the instance is shutdown according to the cloud provider.
// Use the node.name or node.spec.providerID field to find the node in the cloud provider.
func (i2 *HVInstancesV2) InstanceShutdown(ctx context.Context, node *corev1.Node) (bool, error) {
	deviceID, err := getHivelocityDeviceIDFromNode(node)
	if err != nil {
		return false, fmt.Errorf("getHivelocityDeviceIDFromNode(node) failed: %w", err)
	}
	device, err := i2.client.GetBareMetalDevice(ctx, deviceID)
	if err != nil {
		return false, fmt.Errorf("i2.API.GetBareMetalDeviceIdResource(deviceID) failed: %w", err)
	}
	switch device.PowerStatus {
	case "ON":
		return false, nil
	case "OFF":
		return true, nil
	default:
		return false, fmt.Errorf("device with ID %q has unknown PowerStatus %q: %w",
			deviceID, device.PowerStatus, ErrUnknownPowerStatus)
	}
}

// InstanceMetadata returns the instance's metadata. The values returned in InstanceMetadata are
// translated into specific fields and labels in the Node object on registration.
// Implementations should always check node.spec.providerID first when trying to discover the instance
// for a given node. In cases where node.spec.providerID is empty, implementations can use other
// properties of the node like its name, labels and annotations.
func (i2 *HVInstancesV2) InstanceMetadata(ctx context.Context, node *corev1.Node) (
	*cloudprovider.InstanceMetadata, error,
) {
	deviceID, err := getHivelocityDeviceIDFromNode(node)
	if err != nil {
		return nil, fmt.Errorf("getHivelocityDeviceIDFromNode(node) failed: %w", err)
	}
	device, err := i2.client.GetBareMetalDevice(ctx, deviceID)
	if err != nil {
		return nil, fmt.Errorf("i2.API.GetBareMetalDeviceIdResource(deviceID) failed: %w", err)
	}

	addr := corev1.NodeAddress{
		Type:    "ExternalIP",
		Address: device.PrimaryIp,
	}

	// HV tag. Example "instance-type=abc".
	instanceType, err := hvutils.GetInstanceTypeFromTags(device.Tags)
	if err != nil {
		return nil, fmt.Errorf("GetInstanceTypeFromTags() failed. deviceID=%d. %w", deviceID,
			err)
	}

	metaData := cloudprovider.InstanceMetadata{
		ProviderID:    strconv.Itoa(int(deviceID)),
		InstanceType:  instanceType,
		NodeAddresses: []corev1.NodeAddress{addr},
		Zone:          device.LocationName, // for example LAX1
		Region:        device.LocationName, // for example LAX1
	}
	return &metaData, nil
}
