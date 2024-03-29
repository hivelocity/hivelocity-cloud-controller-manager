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

	"github.com/hivelocity/hivelocity-cloud-controller-manager/client"
	cloudprovider "k8s.io/cloud-provider"
	"k8s.io/klog/v2"
)

// cloud implements cloudprovider.Interface for Hivelocity.
type cloud struct {
	instancesV2 *HVInstancesV2
}

const (
	hivelocityAPIKeyENVVar = "HIVELOCITY_API_KEY" // #nosec G101
	providerName           = "hivelocity"
)

var (
	providerVersion                          = "dev"
	_                cloudprovider.Interface = (*cloud)(nil)
	errEnvVarMissing                         = fmt.Errorf("environment variable %q is missing or empty", hivelocityAPIKeyENVVar)
)

func init() {
	cloudprovider.RegisterCloudProvider(providerName, func(config io.Reader) (cloudprovider.Interface, error) {
		return newCloud()
	})
}

func newCloud() (*cloud, error) {
	apiKey := os.Getenv(hivelocityAPIKeyENVVar)
	if apiKey == "" {
		return nil, errEnvVarMissing
	}

	klog.Infof("Hivelocity cloud controller manager %s started\n", providerVersion)

	i2 := newHVInstanceV2(client.NewClient(apiKey))

	return &cloud{
		instancesV2: i2,
	}, nil
}

// Initialize implements cloudprovider.Interface.Initialize.
func (*cloud) Initialize(cloudprovider.ControllerClientBuilder, <-chan struct{}) {
}

// Instances implements cloudprovider.Interface.Instances.
func (*cloud) Instances() (cloudprovider.Instances, bool) {
	// we only implement InstancesV2
	return nil, false
}

// InstancesV2 implements cloudprovider.Interface.InstancesV2.
func (c *cloud) InstancesV2() (cloudprovider.InstancesV2, bool) {
	return c.instancesV2, true
}

// Zones implements cloudprovider.Interface.Zones.
func (*cloud) Zones() (cloudprovider.Zones, bool) {
	// we only implement InstancesV2
	return nil, false
}

// LoadBalancer implements cloudprovider.Interface.LoadBalancer.
func (*cloud) LoadBalancer() (cloudprovider.LoadBalancer, bool) {
	return nil, false // TODO: Up to now Hivelocity has not API for LoadBalancers.
}

// Clusters implements cloudprovider.Interface.Clusters.
func (*cloud) Clusters() (cloudprovider.Clusters, bool) {
	return nil, false // TODO: Will we implement this optional method?
}

// Routes implements cloudprovider.Interface.Routes.
func (*cloud) Routes() (cloudprovider.Routes, bool) {
	return nil, false // TODO: Will we implement this optional method?
}

// ProviderName implements cloudprovider.Interface.ProviderName.
func (*cloud) ProviderName() string {
	return providerName
}

// HasClusterID implements cloudprovider.Interface.HasClusterID.
func (*cloud) HasClusterID() bool {
	// TODO: The meaning if this method is unclear.
	// Waiting for clarification: https://github.com/kubernetes/cloud-provider/issues/64
	return true
}
