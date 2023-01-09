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
	"io"
	"os"
	"strconv"
	"strings"

	//"github.com/hivelocity/hivelocity-cloud-controller-manager/internal/hcops"
	hv "github.com/hivelocity/hivelocity-client-go/client"
	"github.com/hivelocity/hivelocity-cloud-controller-manager/internal/metrics"

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
	nodeNameENVVar                               = "NODE_NAME"
	providerName                                 = "hivelocity"
	providerVersion                              = "v1.9.1"
)

// cloud implements cloudprovider.Interface for Hivelocity.
type cloud struct {
	client      *hv.APIClient
	authContext *context.Context
	instances   *instances
	instancesV2 *HVInstancesV2
	zones       *zones
	//routes       *routes
	//loadBalancer *loadBalancers
	networkID int
}

var _ = cloudprovider.Interface(&cloud{})

func newCloud(config io.Reader) (cloudprovider.Interface, error) {
	const op = "hivelocity/newCloud"
	metrics.OperationCalled.WithLabelValues(op).Inc()

	apiKey := os.Getenv(hivelocityApiKeyENVVar)
	if apiKey == "" {
		return nil, fmt.Errorf("environment variable %q is required", hivelocityApiKeyENVVar)
	}

	/*
		nodeName := os.Getenv(nodeNameENVVar)
		if nodeName == "" {
			return nil, fmt.Errorf("environment variable %q is required", nodeNameENVVar)
		}
	*/

	/*
		// start metrics server if enabled (enabled by default)
		if os.Getenv(hivelocityMetricsEnabledENVVar) != "false" {
			go metrics.Serve(hivelocityMetricsAddress)

			opts = append(opts, hv.WithInstrumentation(metrics.GetRegistry()))
		}

		if os.Getenv(hivelocityDebugENVVar) == "true" {
			opts = append(opts, hv.WithDebugWriter(os.Stderr))
		}
		if endpoint := os.Getenv(hivelocityEndpointENVVar); endpoint != "" {
			opts = append(opts, hv.WithEndpoint(endpoint))
		}
	*/

	authContext := context.WithValue(context.Background(), hv.ContextAPIKey, hv.APIKey{
		Key: apiKey,
	})
	client := hv.NewAPIClient(hv.NewConfiguration())

	/*
		var networkID int
		if v, ok := os.LookupEnv(hivelocityNetworkENVVar); ok {
			n, _, err := client.Network.Get(context.Background(), v)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", op, err)
			}
			if n == nil {
				return nil, fmt.Errorf("%s: Network %s not found", op, v)
			}
			networkID = n.ID

			networkDisableAttachedCheck, err := getEnvBool(hivelocityNetworkDisableAttachedCheckENVVar)
			if err != nil {
				return nil, fmt.Errorf("%s: checking if server is in Network not possible: %w", op, err)
			}
			if !networkDisableAttachedCheck {
				e, err := serverIsAttachedToNetwork(metadataClient, networkID)
				if err != nil {
					return nil, fmt.Errorf("%s: checking if server is in Network not possible: %w", op, err)
				}
				if !e {
					return nil, fmt.Errorf("%s: This node is not attached to Network %s", op, v)
				}
			}
		}
		if networkID == 0 {
			klog.Infof("%s: %s empty", op, hivelocityNetworkENVVar)
		}

		_, _, err := client.Server.List(context.Background(), hv.ServerListOpts{})
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		lbOpsDefaults, lbDisablePrivateIngress, lbDisableIPv6, err := loadBalancerDefaultsFromEnv()
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	*/

	klog.Infof("Hivelocity Cloud k8s cloud controller %s started\n", providerVersion)

	/*
		lbOps := &hcops.LoadBalancerOps{
			LBClient:      &client.LoadBalancer,
			CertOps:       &hcops.CertificateOps{CertClient: &client.Certificate},
			ActionClient:  &client.Action,
			NetworkClient: &client.Network,
			NetworkID:     networkID,
			Defaults:      lbOpsDefaults,
		}

		loadBalancers := newLoadBalancers(lbOps, &client.Action, lbDisablePrivateIngress, lbDisableIPv6)
		if os.Getenv(hivelocityLoadBalancersEnabledENVVar) == "false" {
			loadBalancers = nil
		}
	*/
	instancesAddressFamily, err := addressFamilyFromEnv()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &cloud{
		client:      client,
		authContext: &authContext,
		//zones:       newZones(client, nodeName),
		instances: newInstances(client, instancesAddressFamily),
		//	loadBalancer: loadBalancers,
		//routes:    nil,
		//networkID: networkID,
	}, nil
}

func (c *cloud) Initialize(clientBuilder cloudprovider.ControllerClientBuilder, stop <-chan struct{}) {
}

func (c *cloud) Instances() (cloudprovider.Instances, bool) {
	return c.instances, true
}

func (c *cloud) InstancesV2() (cloudprovider.InstancesV2, bool) {
	// TODO enable InstancesV2 and disable old way.
	return nil, false
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

func (c *cloud) HasClusterID() bool {
	return false
}

/*
func loadBalancerDefaultsFromEnv() (hcops.LoadBalancerDefaults, bool, bool, error) {
	defaults := hcops.LoadBalancerDefaults{
		Location:    os.Getenv(hivelocityLoadBalancersLocation),
		NetworkZone: os.Getenv(hivelocityLoadBalancersNetworkZone),
	}

	if defaults.Location != "" && defaults.NetworkZone != "" {
		return defaults, false, false, errors.New(
			"HIVELOCITY_LOAD_BALANCERS_LOCATION/HIVELOCITY_LOAD_BALANCERS_NETWORK_ZONE: Only one of these can be set")
	}

	disablePrivateIngress, err := getEnvBool(hivelocityLoadBalancersDisablePrivateIngress)
	if err != nil {
		return defaults, false, false, err
	}

	disableIPv6, err := getEnvBool(hivelocityLoadBalancersDisableIPv6)
	if err != nil {
		return defaults, false, false, err
	}

	defaults.UsePrivateIP, err = getEnvBool(hivelocityLoadBalancersUsePrivateIP)
	if err != nil {
		return defaults, false, false, err
	}

	return defaults, disablePrivateIngress, disableIPv6, nil
}
*/

// serverIsAttachedToNetwork checks if the server where the master is running on is attached to the configured private network
// We use this measurement to protect users against some parts of misconfiguration, like configuring a master in a not attached
// network.
/*
func serverIsAttachedToNetwork(metadataClient *metadata.Client, networkID int) (bool, error) {
	const op = "serverIsAttachedToNetwork"
	metrics.OperationCalled.WithLabelValues(op).Inc()

	serverPrivateNetworks, err := metadataClient.PrivateNetworks()
	if err != nil {
		return false, fmt.Errorf("%s: %s", op, err)
	}
	return strings.Contains(serverPrivateNetworks, fmt.Sprintf("network_id: %d\n", networkID)), nil
}
*/

// addressFamilyFromEnv returns the address family for the instance address from the environment
// variable. Returns AddressFamilyIPv4 if unset.
func addressFamilyFromEnv() (addressFamily, error) {
	family, ok := os.LookupEnv(hivelocityInstancesAddressFamily)
	if !ok {
		return AddressFamilyIPv4, nil
	}

	switch strings.ToLower(family) {
	case "ipv6":
		return AddressFamilyIPv6, nil
	case "ipv4":
		return AddressFamilyIPv4, nil
	case "dualstack":
		return AddressFamilyDualStack, nil
	default:
		return -1, fmt.Errorf(
			"%v: Invalid value, expected one of: ipv4,ipv6,dualstack", hivelocityInstancesAddressFamily)
	}
}

// getEnvBool returns the boolean parsed from the environment variable with the given key and a potential error
// parsing the var. Returns false if the env var is unset.
func getEnvBool(key string) (bool, error) {
	v, ok := os.LookupEnv(key)
	if !ok {
		return false, nil
	}

	b, err := strconv.ParseBool(v)
	if err != nil {
		return false, fmt.Errorf("%s: %v", key, err)
	}

	return b, nil
}

func init() {
	cloudprovider.RegisterCloudProvider(providerName, func(config io.Reader) (cloudprovider.Interface, error) {
		return newCloud(config)
	})
}
