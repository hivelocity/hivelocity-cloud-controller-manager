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

	hv "github.com/hivelocity/hivelocity-client-go/client"
	"github.com/hivelocity/hivelocity-cloud-controller-manager/internal/metrics"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	cloudprovider "k8s.io/cloud-provider"
)

type addressFamily int

const (
	AddressFamilyDualStack addressFamily = iota
	AddressFamilyIPv6
	AddressFamilyIPv4
)

type instances struct {
	client        *hv.APIClient
	addressFamily addressFamily
}

func newInstances(client *hv.APIClient, addressFamily addressFamily) *instances {
	return &instances{client, addressFamily}
}

func (i *instances) NodeAddressesByProviderID(ctx context.Context, providerID string) ([]v1.NodeAddress, error) {
	return nil, fmt.Errorf("TODO, implementNodeAddressesByProviderID()")
	/*
		const op = "hv/instances.NodeAddressesByProviderID"
		metrics.OperationCalled.WithLabelValues(op).Inc()

		id, err := providerIDToServerID(providerID)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		server, err := getServerByID(ctx, i.client, id)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		return i.nodeAddresses(ctx, server), nil
	*/
}

func (i *instances) NodeAddresses(ctx context.Context, nodeName types.NodeName) ([]v1.NodeAddress, error) {
	return nil, fmt.Errorf("TODO: implementNodeAddresses()")
	/*
		const op = "hv/instances.NodeAddresses"
		metrics.OperationCalled.WithLabelValues(op).Inc()

		server, err := getServerByName(ctx, i.client, string(nodeName))
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		return i.nodeAddresses(ctx, server), nil
	*/
}

func (i *instances) ExternalID(ctx context.Context, nodeName types.NodeName) (string, error) {
	const op = "hv/instances.ExternalID"
	metrics.OperationCalled.WithLabelValues(op).Inc()

	id, err := i.InstanceID(ctx, nodeName)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (i *instances) InstanceID(ctx context.Context, nodeName types.NodeName) (string, error) {
	return "", fmt.Errorf("TODO implement InstanceID()")
	/*
		const op = "hv/instances.InstanceID"
		metrics.OperationCalled.WithLabelValues(op).Inc()

		server, err := getServerByName(ctx, i.client, string(nodeName))
		if err != nil {
			return "", fmt.Errorf("%s: %w", op, err)
		}
		return strconv.Itoa(server.ID), nil
	*/
}

func (i *instances) InstanceType(ctx context.Context, nodeName types.NodeName) (string, error) {
	return "", fmt.Errorf("TODO implement InstanceType()")
	/*
		const op = "hv/instances.InstanceType"
		metrics.OperationCalled.WithLabelValues(op).Inc()

		server, err := getServerByName(ctx, i.client, string(nodeName))
		if err != nil {
			return "", fmt.Errorf("%s: %w", op, err)
		}
		return server.ServerType.Name, nil
	*/
}

func (i *instances) InstanceTypeByProviderID(ctx context.Context, providerID string) (string, error) {
	return "", fmt.Errorf("TODO implement InstanceTypeByProviderID()")
	/*
		const op = "hv/instances.InstanceTypeByProviderID"
		metrics.OperationCalled.WithLabelValues(op).Inc()

		id, err := providerIDToServerID(providerID)
		if err != nil {
			return "", fmt.Errorf("%s: %w", op, err)
		}

		server, err := getServerByID(ctx, i.client, id)
		if err != nil {
			return "", fmt.Errorf("%s: %w", op, err)
		}
		return server.ServerType.Name, nil
	*/
}

func (i *instances) AddSSHKeyToAllInstances(ctx context.Context, user string, keyData []byte) error {
	return cloudprovider.NotImplemented
}

func (i *instances) CurrentNodeName(ctx context.Context, hostname string) (types.NodeName, error) {
	return types.NodeName(hostname), nil
}

func (i instances) InstanceExistsByProviderID(ctx context.Context, providerID string) (bool, error) {
	return false, fmt.Errorf("TODO implement InstanceExistsByProviderID()")
	/*
		const op = "hv/instances.InstanceExistsByProviderID"
		metrics.OperationCalled.WithLabelValues(op).Inc()

		id, err := providerIDToServerID(providerID)
		if err != nil {
			return false, fmt.Errorf("%s: %w", op, err)
		}

		server, _, err := i.client.Server.GetByID(ctx, id)
		if err != nil {
			return false, fmt.Errorf("%s: %w", op, err)
		}
		return server != nil, nil
	*/
}

func (i instances) InstanceShutdownByProviderID(ctx context.Context, providerID string) (bool, error) {
	return false, fmt.Errorf("TODO implement InstanceShutdownByProviderID")
	/*
		const op = "hv/instances.InstanceShutdownByProviderID"
		metrics.OperationCalled.WithLabelValues(op).Inc()

		id, err := providerIDToServerID(providerID)
		if err != nil {
			return false, fmt.Errorf("%s: %w", op, err)
		}

		server, _, err := i.client.Server.GetByID(ctx, id)
		if err != nil {
			return false, fmt.Errorf("%s: %w", op, err)
		}
		return server != nil && server.Status == hv.ServerStatusOff, nil
	*/
}

/*
func (i *instances) nodeAddresses(ctx context.Context, server *hv.Server) []v1.NodeAddress {
	var addresses []v1.NodeAddress
	addresses = append(
		addresses,
		v1.NodeAddress{Type: v1.NodeHostName, Address: server.Name},
	)

	if i.addressFamily == AddressFamilyIPv4 || i.addressFamily == AddressFamilyDualStack {
		if !server.PublicNet.IPv4.IP.IsUnspecified() {
			addresses = append(
				addresses,
				v1.NodeAddress{Type: v1.NodeExternalIP, Address: server.PublicNet.IPv4.IP.String()},
			)
		}
	}

	if i.addressFamily == AddressFamilyIPv6 || i.addressFamily == AddressFamilyDualStack {
		if !server.PublicNet.IPv6.IP.IsUnspecified() {
			// For a given IPv6 network of 2001:db8:1234::/64, the instance address is 2001:db8:1234::1
			host_address := server.PublicNet.IPv6.IP
			host_address[len(host_address)-1] |= 0x01

			addresses = append(
				addresses,
				v1.NodeAddress{Type: v1.NodeExternalIP, Address: host_address.String()},
			)
		}
	}

	n := os.Getenv(hvNetworkENVVar)
	if len(n) > 0 {
		network, _, _ := i.client.Network.Get(ctx, n)
		if network != nil {
			for _, privateNet := range server.PrivateNet {
				if privateNet.Network.ID == network.ID {
					addresses = append(
						addresses,
						v1.NodeAddress{Type: v1.NodeInternalIP, Address: privateNet.IP.String()},
					)
				}
			}

		}
	}
	return addresses
}
*/
