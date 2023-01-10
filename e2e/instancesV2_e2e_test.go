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
	"os"
	"strconv"
	"testing"

	hv "github.com/hivelocity/hivelocity-client-go/client"
	"github.com/stretchr/testify/require"

	"github.com/hexops/autogold"
	"github.com/hivelocity/hivelocity-cloud-controller-manager/client"
	"github.com/hivelocity/hivelocity-cloud-controller-manager/hivelocity"
	"github.com/joho/godotenv"
	corev1 "k8s.io/api/core/v1"
	cloudprovider "k8s.io/cloud-provider"
)

func getAPIClient() *hv.APIClient {
	err := godotenv.Overload("../.envrc")
	if err != nil {
		panic(err)
	}

	apiKey := os.Getenv("HIVELOCITY_API_KEY")
	if apiKey == "" {
		panic("Missing environment variable HIVELOCITY_API_KEY")
	}
	config := hv.NewConfiguration()
	config.AddDefaultHeader("X-API-KEY", apiKey)
	return hv.NewAPIClient(config)
}

var e2eDeviceId int = 14730

func newHVInstanceV2() *hivelocity.HVInstancesV2 {
	var i2 hivelocity.HVInstancesV2
	apiClient := getAPIClient()
	i2.API = &client.RealAPI{
		Client: apiClient,
	}
	i2.Client = apiClient
	return &i2
}

func newNode() *corev1.Node {
	return &corev1.Node{
		Spec: corev1.NodeSpec{
			ProviderID: strconv.Itoa(e2eDeviceId),
		},
	}
}

func Test_InstanceExists(t *testing.T) {
	i2 := newHVInstanceV2()
	node := newNode()
	ctx := context.Background()
	myBool, err := i2.InstanceExists(ctx, node)
	require.NoError(t, err)
	require.Equal(t, true, myBool)

	node.Spec.ProviderID = "9999999"
	myBool, err = i2.InstanceExists(ctx, node)
	require.Equal(t, false, myBool)
	require.NoError(t, err)

	node.Spec.ProviderID = "9999999999999999999999999999"
	myBool, err = i2.InstanceExists(ctx, node)
	require.Equal(t, false, myBool)
	require.Equal(t, "failed to convert node.Spec.ProviderID \"9999999999999999999999999999\" to int32", err.Error())
}

func Test_InstanceShutdown(t *testing.T) {
	i2 := newHVInstanceV2()
	node := newNode()
	ctx := context.Background()
	isDown, err := i2.InstanceShutdown(ctx, node)
	require.False(t, isDown)
	require.NoError(t, err)

	node.Spec.ProviderID = "9999999"
	_, err = i2.InstanceShutdown(ctx, node)
	require.Error(t, err)
}

func Test_InstanceMetadata(t *testing.T) {
	i2 := newHVInstanceV2()
	node := newNode()
	ctx := context.Background()
	metaData, err := i2.InstanceMetadata(ctx, node)
	require.NoError(t, err)
	autogold.Want("metaData", &cloudprovider.InstanceMetadata{
		ProviderID: "14730",
		NodeAddresses: []corev1.NodeAddress{
			{
				Type:    corev1.NodeAddressType("ExternalIP"),
				Address: "66.165.243.74",
			},
		},
		Zone: "LAX2",
	}).Equal(t, metaData)
}
