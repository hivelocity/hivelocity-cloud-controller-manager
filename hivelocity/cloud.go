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
	"fmt"
	"io"
	"os"

	//"github.com/hivelocity/hivelocity-cloud-controller-manager/internal/hcops"
	hv "github.com/hivelocity/hivelocity-client-go/client"
	"github.com/hivelocity/hivelocity-cloud-controller-manager/client"

	//"github.com/hivelocity/hivelocity-client-go/hivelocity/metadata"
	cloudprovider "k8s.io/cloud-provider"
	"k8s.io/klog/v2"
)

const (
	// TODO remove not used
	hivelocityApiKeyENVVar   = "HIVELOCITY_API_KEY"
	hivelocityEndpointENVVar = "HIVELOCITY_ENDPOINT"
	hivelocityNetworkENVVar  = "HIVELOCITY_NETWORK"
	hivelocityDebugENVVar    = "HIVELOCITY_DEBUG"
	// Disable the "master/server is attached to the network" check against the metadata service.
	hivelocityNetworkDisableAttachedCheckENVVar  = "HIVELOCITY_NETWORK_DISABLE_ATTACHED_CHECK"
	hivelocityNetworkRoutesEnabledENVVar         = "HIVELOCITY_NETWORK_ROUTES_ENABLED"
	hivelocityInstancesAddressFamily             = "HIVELOCITY_INSTANCES_ADDRESS_FAMILY"
	hivelocityLoadBalancersEnabledENVVar         = "HIVELOCITY_LOAD_BALANCERS_ENABLED"
	hivelocityLoadBalancersLocation              = "HIVELOCITY_LOAD_BALANCERS_LOCATION"
	hivelocityLoadBalancersNetworkZone           = "HIVELOCITY_LOAD_BALANCERS_NETWORK_ZONE"
	hivelocityLoadBalancersDisablePrivateIngress = "HIVELOCITY_LOAD_BALANCERS_DISABLE_PRIVATE_INGRESS"
	hivelocityLoadBalancersUsePrivateIP          = "HIVELOCITY_LOAD_BALANCERS_USE_PRIVATE_IP"
	hivelocityLoadBalancersDisableIPv6           = "HIVELOCITY_LOAD_BALANCERS_DISABLE_IPV6"
	hivelocityMetricsEnabledENVVar               = "HIVELOCITY_METRICS_ENABLED"
	hivelocityMetricsAddress                     = ":8233"
	providerName                                 = "hivelocity"
	providerVersion                              = "v1.9.1"
)

// cloud implements cloudprovider.Interface for Hivelocity.
type cloud struct {
	client      *hv.APIClient
	instancesV2 *HVInstancesV2
	zones       *zones
	//routes       *cloudprovider.Routes
	//loadBalancer *cloudprovider.LoadBalancer
	networkID int
}

var _ cloudprovider.Interface = (*cloud)(nil)

func newCloud(config io.Reader) (cloudprovider.Interface, error) {
	apiKey := os.Getenv(hivelocityApiKeyENVVar)
	if apiKey == "" {
		return nil, fmt.Errorf("environment variable %q is required", hivelocityApiKeyENVVar)
	}

	apiClientConfig := hv.NewConfiguration()
	apiClientConfig.AddDefaultHeader("X-API-KEY", apiKey)
	apiClient := hv.NewAPIClient(apiClientConfig)

	klog.Infof("Hivelocity cloud controller manager %s started\n", providerVersion)

	i2 := HVInstancesV2{
		API: &client.RealAPI{
			Client: apiClient,
		},
	}

	return &cloud{
		client:      apiClient,
		instancesV2: &i2,
	}, nil
}

func (c *cloud) Initialize(clientBuilder cloudprovider.ControllerClientBuilder, stop <-chan struct{}) {
}

func (c *cloud) Instances() (cloudprovider.Instances, bool) {
	// we only implement InstancesV2
	return nil, false
}

func (c *cloud) InstancesV2() (cloudprovider.InstancesV2, bool) {
	return c.instancesV2, true
}

func (c *cloud) Zones() (cloudprovider.Zones, bool) {
	return c.zones, true
}

func (c *cloud) LoadBalancer() (cloudprovider.LoadBalancer, bool) {
	return nil, false
	/* TODO
	if c.loadBalancer == nil {
		return nil, false
	}
	return c.loadBalancer, true
	*/
}

func (c *cloud) Clusters() (cloudprovider.Clusters, bool) {
	return nil, false
}

func (c *cloud) Routes() (cloudprovider.Routes, bool) {
	/* TODO
	if c.networkID > 0 && os.Getenv(hivelocityNetworkRoutesEnabledENVVar) != "false" {
		r, err := newRoutes(c.client, c.networkID)
		if err != nil {
			klog.ErrorS(err, "create routes provider", "networkID", c.networkID)
			return nil, false
		}
		return r, true
	}
	*/
	return nil, false // If no network is configured, disable the routes part
}

func (c *cloud) ProviderName() string {
	return providerName
}

func (c *cloud) ScrubDNS(nameservers, searches []string) (nsOut, srchOut []string) {
	return nil, nil
}

// HasClusterID is not implemented.
func (c *cloud) HasClusterID() bool {
	return true
}

func init() {
	cloudprovider.RegisterCloudProvider(providerName, func(config io.Reader) (cloudprovider.Interface, error) {
		return newCloud(config)
	})
}
