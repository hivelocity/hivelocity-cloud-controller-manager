package hivelocity

import (
	"context"
	"fmt"
	"strconv"

	hv "github.com/hivelocity/hivelocity-client-go/client"
	v1 "k8s.io/api/core/v1"
)

// hvInstancesV2 implements cloudprovider.InstanceV2
type hvInstancesV2 struct {
	client        *hv.APIClient
	addressFamily addressFamily
}

func getHivelocityDeviceIdFromNode(node *v1.Node) (int32, error){
	deviceId, err := strconv.ParseInt(node.Spec.ProviderID, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("failed to convert node.Spec.ProviderID %q to int32.",
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
	if response.StatusCode == 404 {
		// todo: ggf "schöner" err mit (no device)?
		// methode schreiben, die diese err nachricht ausliest.
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("GetBareMetalDeviceIdResource failed node.Spec.ProviderID %q. %w",
			node.Spec.ProviderID, err)
	}
	return true, nil
}

// InstanceShutdown returns true if the instance is shutdown according to the cloud provider.
// Use the node.name or node.spec.providerID field to find the node in the cloud provider.
func (i2 *hvInstancesV2) InstanceShutdown(ctx context.Context, node *v1.Node) (bool, error){
	deviceID, err := getHivelocityDeviceIdFromNode(node)
	if err != nil {
		return false, err
	}
	device, _, err := i2.client.BareMetalDevicesApi.GetBareMetalDeviceIdResource(ctx, deviceID, nil)
	if err != nil {
		return false, err
	}
	switch device.PowerStatus {
	case "ON": return false, nil
	case "OFF": return true, nil
	default: return false, fmt.Errorf("Device with ID %q has unknown PowerStatus %q.", deviceID, device.PowerStatus)
	}
}

