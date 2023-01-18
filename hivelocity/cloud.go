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

	hv "github.com/hivelocity/hivelocity-client-go/client"
	"github.com/hivelocity/hivelocity-cloud-controller-manager/client"
	cloudprovider "k8s.io/cloud-provider"
	"k8s.io/klog/v2"
)

const (
	hivelocityApiKeyENVVar = "HIVELOCITY_API_KEY"
	providerName           = "hivelocity"
	providerVersion        = "v0.0.1"
)

// cloud implements cloudprovider.Interface for Hivelocity.
type cloud struct {
	client      *hv.APIClient
	instancesV2 *HVInstancesV2
}

var _ cloudprovider.Interface = (*cloud)(nil)

func newCloud(config io.Reader) (cloudprovider.Interface, error) {
	apiKey := os.Getenv(hivelocityApiKeyENVVar)
	if apiKey == "" {
		return nil, fmt.Errorf("environment variable %q is missing or empty", hivelocityApiKeyENVVar)
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
	return nil, false
}

func (c *cloud) LoadBalancer() (cloudprovider.LoadBalancer, bool) {
	return nil, false
}

func (c *cloud) Clusters() (cloudprovider.Clusters, bool) {
	return nil, false
}

func (c *cloud) Routes() (cloudprovider.Routes, bool) {
	return nil, false
}

func (c *cloud) ProviderName() string {
	return providerName
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
