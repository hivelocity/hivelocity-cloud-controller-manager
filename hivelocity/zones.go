/*
Copyright 2018 Hivelocity Cloud GmbH.

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
	"k8s.io/apimachinery/pkg/types"
	cloudprovider "k8s.io/cloud-provider"
)

type zones struct {
	client   *hv.APIClient
	nodeName string // name of the node the programm is running on
}

func newZones(client *hv.APIClient, nodeName string) *zones {
	return &zones{client, nodeName}
}

func (z zones) GetZone(ctx context.Context) (cloudprovider.Zone, error) {
	return cloudprovider.Zone{}, fmt.Errorf("TODO implement GetZone")
	/*
		const op = "hv/zones.GetZone"
		metrics.OperationCalled.WithLabelValues(op).Inc()

		server, err := getServerByName(ctx, z.client, z.nodeName)
		if err != nil {
			return cloudprovider.Zone{}, fmt.Errorf("%s: %w", op, err)
		}
		return zoneFromServer(server), nil
	*/
}

func (z zones) GetZoneByProviderID(ctx context.Context, providerID string) (cloudprovider.Zone, error) {
	return cloudprovider.Zone{}, fmt.Errorf("TODO implement GetZoneByProviderID")
	/*
		const op = "hv/zones.GetZoneByProviderID"
		metrics.OperationCalled.WithLabelValues(op).Inc()

		id, err := providerIDToServerID(providerID)
		if err != nil {
			return cloudprovider.Zone{}, fmt.Errorf("%s: %w", op, err)
		}

		server, err := getServerByID(ctx, z.client, id)
		if err != nil {
			return cloudprovider.Zone{}, fmt.Errorf("%s: %w", op, err)
		}

		return zoneFromServer(server), nil
	*/
}

func (z zones) GetZoneByNodeName(ctx context.Context, nodeName types.NodeName) (cloudprovider.Zone, error) {
	return cloudprovider.Zone{}, fmt.Errorf("TODO: implement GetZoneByNodeName()")
	/*
		const op = "hv/zones.GetZoneByNodeName"
		metrics.OperationCalled.WithLabelValues(op).Inc()

		server, err := getServerByName(ctx, z.client, string(nodeName))
		if err != nil {
			return cloudprovider.Zone{}, fmt.Errorf("%s: %w", op, err)
		}

		return zoneFromServer(server), nil
	*/
}

/*
func zoneFromServer(server *hv.Server) cloudprovider.Zone {
	return cloudprovider.Zone{
		Region:        server.Datacenter.Location.Name,
		FailureDomain: server.Datacenter.Name,
	}
}
*/
