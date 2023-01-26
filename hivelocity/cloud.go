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

// Package hivelocity implements the cloud controller manager.
// The interfaces are from: https://github.com/kubernetes/cloud-provider
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
	hivelocityAPIKeyENVVar = "HIVELOCITY_API_KEY" // #nosec G101
	providerName           = "hivelocity"
	providerVersion        = "v0.0.1"
)

// cloud implements cloudprovider.Interface for Hivelocity.
type cloud struct {
	client      *hv.APIClient
	instancesV2 *HVInstancesV2
}

var _ cloudprovider.Interface = (*cloud)(nil)

func newCloud() (*cloud, error) {
	apiKey := os.Getenv(hivelocityAPIKeyENVVar)
	if apiKey == "" {
		return nil, errEnvVarMissing
	}

	apiClientConfig := hv.NewConfiguration()
	apiClientConfig.AddDefaultHeader("X-API-KEY", apiKey)
	apiClient := hv.NewAPIClient(apiClientConfig)

	klog.Infof("Hivelocity cloud controller manager %s started\n", providerVersion)

	i2 := NewHVInstanceV2(client.NewClient(apiClient))

	return &cloud{
		client:      apiClient,
		instancesV2: i2,
	}, nil
}

var errEnvVarMissing = fmt.Errorf("environment variable %q is missing or empty", hivelocityAPIKeyENVVar)

func (c *cloud) Initialize(clientBuilder cloudprovider.ControllerClientBuilder, stop <-chan struct{}) {
}

func (c *cloud) Instances() (cloudprovider.Instances, bool) { //nolint:ireturn // implements cloudprovider.Interface
	// we only implement InstancesV2
	return nil, false
}

func (c *cloud) InstancesV2() (cloudprovider.InstancesV2, bool) { //nolint:ireturn // implements cloudprovider.Interface
	return c.instancesV2, true
}

func (c *cloud) Zones() (cloudprovider.Zones, bool) { //nolint:ireturn // implements cloudprovider.Interface
	// we only implement InstancesV2
	return nil, false
}

func (c *cloud) LoadBalancer() (cloudprovider.LoadBalancer, bool) { //nolint:ireturn,lll // implements cloudprovider.Interface
	return nil, false // TODO: Up to now Hivelocity has not API for LoadBalancers.
}

func (c *cloud) Clusters() (cloudprovider.Clusters, bool) { //nolint:ireturn // implements cloudprovider.Interface
	return nil, false // TODO: Will we implement this optional method?
}

func (c *cloud) Routes() (cloudprovider.Routes, bool) { //nolint:ireturn // implements cloudprovider.Interface
	return nil, false // TODO: Will we implement this optional method?
}

func (c *cloud) ProviderName() string {
	return providerName
}

// HasClusterID is not implemented.
func (c *cloud) HasClusterID() bool {
	// TODO: The meaning if this method is unclear.
	// Waiting for clarification: https://github.com/kubernetes/cloud-provider/issues/64
	return true
}

func init() {
	cloudprovider.RegisterCloudProvider(providerName, func(config io.Reader) (cloudprovider.Interface, error) {
		return newCloud()
	})
}
