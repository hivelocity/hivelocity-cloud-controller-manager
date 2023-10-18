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

	hv "github.com/hivelocity/hivelocity-client-go/client"
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

var (
	errNodeIsNil = errors.New("node is nil")

	errUnknownPowerStatus = errors.New("unknown PowerStatus")

	errNoDeviceFound = errors.New("no device found")
)

var (
	errMissingProviderPrefix = fmt.Errorf(
		"missing prefix %q in node.Spec.ProviderID",
		providerName,
	)
	errFailedToConvertProviderID = fmt.Errorf("failed to convert node.Spec.ProviderID")
)

// newHVInstanceV2 creates a new HVInstancesV2 struct.
func newHVInstanceV2(c client.Interface) *HVInstancesV2 {
	return &HVInstancesV2{client: c}
}

// getHivelocityDeviceIDFromNode returns the deviceID from a Node.
// Example: If Node.Spec.ProviderID is "hivelocity://123", then 123
// will be returned.
func getHivelocityDeviceIDFromNode(node *corev1.Node) (int32, error) {
	providerPrefix := providerName + "://"
	if !strings.HasPrefix(node.Spec.ProviderID, providerPrefix) {
		return 0, fmt.Errorf(
			"[getHivelocityDeviceIDFromNode] HasPrefix() failed. Node %q, ProviderID %q: %w",
			node.GetName(),
			node.Spec.ProviderID,
			errMissingProviderPrefix,
		)
	}
	deviceID, err := strconv.ParseInt(
		strings.TrimPrefix(node.Spec.ProviderID, providerPrefix),
		10,
		32,
	)
	if err != nil {
		return 0, fmt.Errorf(
			"[getHivelocityDeviceIDFromNode] ParseInt() failed. Node %q, ProviderID %q: %w",
			node.GetName(),
			node.Spec.ProviderID,
			errFailedToConvertProviderID,
		)
	}
	return int32(deviceID), nil
}

// lookUpDevice looks for device via Hivelocity API if provider ID is present otherwise look for machine name label
// present in the devices (caphv-machine-name=foo).
func (i2 *HVInstancesV2) lookUpDevice(ctx context.Context, node *corev1.Node) (device *hv.BareMetalDevice, err error) {
	if node.Spec.ProviderID != "" {
		deviceID, err := getHivelocityDeviceIDFromNode(node)
		if err != nil {
			return nil, fmt.Errorf(
				"[lookUpDevice] getHivelocityDeviceIDFromNode() failed. node %q: %w",
				node.GetName(),
				err,
			)
		}
		device, err = i2.client.GetBareMetalDevice(ctx, deviceID)
		if errors.Is(err, client.ErrNoSuchDevice) {
			return nil, nil //nolint:nilnil // we ignore the error if no device is found.
		}
		if err != nil {
			return nil, fmt.Errorf(
				"[InstanceExists] GetBareMetalDevice() failed. node %q, deviceID %d: %w",
				node.GetName(),
				deviceID,
				err,
			)
		}
	} else {
		devices, err := i2.client.ListDevices(ctx)
		if err != nil {
			return nil, fmt.Errorf(
				"[InstanceExists] ListDevices() failed. node %q: %w",
				node.GetName(),
				err,
			)
		}

		for i := range devices {
			device := devices[i]
			name, err := hvutils.GetMachineNameFromTags(device.Tags)
			if err != nil {
				continue
			}

			if name == node.GetName() {
				return &device, nil
			}
		}
	}

	return device, nil
}

// InstanceExists returns true if the instance for the given node exists according to the cloud provider.
// Use the node.name or node.spec.providerID field to find the node in the cloud provider.
// Implements cloudprovider.InstancesV2.InstanceExists.
func (i2 *HVInstancesV2) InstanceExists(ctx context.Context, node *corev1.Node) (bool, error) {
	const op = "hivelocity/instancesv2.InstanceExists"

	if node == nil {
		return false, errNodeIsNil
	}

	device, err := i2.lookUpDevice(ctx, node)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	if device == nil {
		return false, nil
	}

	name, err := hvutils.GetMachineNameFromTags(device.Tags)
	if err != nil {
		return false, nil //nolint:nilerr // we ignore the device if there is no such label available.
	}

	return name == node.GetName(), nil
}

// InstanceShutdown returns true if the instance is shutdown according to the cloud provider.
// Use the node.name or node.spec.providerID field to find the node in the cloud provider.
// Implements cloudprovider.InstancesV2.InstanceShutdown.
func (i2 *HVInstancesV2) InstanceShutdown(ctx context.Context, node *corev1.Node) (bool, error) {
	const op = "hivelocity/instancesv2.InstanceShutdown"

	if node == nil {
		return false, errNodeIsNil
	}

	device, err := i2.lookUpDevice(ctx, node)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	if device == nil {
		return false, errNoDeviceFound
	}

	switch device.PowerStatus {
	case "ON":
		return false, nil
	case "OFF":
		return true, nil
	default:
		return false, fmt.Errorf(
			"[InstanceShutdown] unknown PowerStatus %q. deviceID %d, node %q: %w",
			device.PowerStatus,
			device.DeviceId,
			node.GetName(),
			errUnknownPowerStatus,
		)
	}
}

// InstanceMetadata returns the instance's metadata. The values returned in InstanceMetadata are
// translated into specific fields and labels in the Node object on registration.
// Implementations should always check node.spec.providerID first when trying to discover the instance
// for a given node. In cases where node.spec.providerID is empty, implementations can use other
// properties of the node like its name, labels and annotations.
// Implements cloudprovider.InstancesV2.InstanceMetadata.
func (i2 *HVInstancesV2) InstanceMetadata(
	ctx context.Context,
	node *corev1.Node,
) (*cloudprovider.InstanceMetadata, error) {
	const op = "hivelocity/instancesv2.InstanceMetadata"

	if node == nil {
		return nil, errNodeIsNil
	}

	device, err := i2.lookUpDevice(ctx, node)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if device == nil {
		return nil, errNoDeviceFound
	}

	// HV tag. Example "caphv-device-type=abc".
	instanceType, err := hvutils.GetInstanceTypeFromTags(device.Tags)
	if err != nil {
		return nil, fmt.Errorf(
			"InstanceMetadata(): GetInstanceTypeFromTags() failed. node %q, deviceID %d: %w",
			node.GetName(),
			device.DeviceId,
			err,
		)
	}

	metaData := cloudprovider.InstanceMetadata{
		ProviderID:   strconv.Itoa(int(device.DeviceId)),
		InstanceType: instanceType,
		NodeAddresses: []corev1.NodeAddress{{
			Type:    "ExternalIP",
			Address: device.PrimaryIp,
		}},
		Zone:   device.LocationName, // for example LAX1
		Region: device.LocationName, // for example LAX1
	}
	return &metaData, nil
}
