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
	"testing"

	hv "github.com/hivelocity/hivelocity-client-go/client"
	"github.com/hivelocity/hivelocity-cloud-controller-manager/client"
	"github.com/hivelocity/hivelocity-cloud-controller-manager/mocks"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	cloudprovider "k8s.io/cloud-provider"
)

func Test_GetHivelocityDeviceIdFromNode(t *testing.T) {
	// Empty ProviderID should fail
	i, err := GetHivelocityDeviceIdFromNode(&corev1.Node{})
	require.Equal(t, int32(0), i)
	require.Error(t, err)

	// Correct ProviderID should get parsed
	i, err = GetHivelocityDeviceIdFromNode(&corev1.Node{
		Spec: corev1.NodeSpec{
			ProviderID: "hivelocity://12345",
		}})
	require.Equal(t, int32(12345), i)
	require.NoError(t, err)
}

var mockDeviceId int = 14730

func newHVInstanceV2(t *testing.T) (*HVInstancesV2, *mocks.API) {
	api := mocks.NewAPI(t)
	return &HVInstancesV2{
		API: api,
	}, api
}

func newNode() *corev1.Node {
	return &corev1.Node{
		Spec: corev1.NodeSpec{
			ProviderID: fmt.Sprintf("hivelocity://%d", mockDeviceId),
		},
	}
}

func standardMocks(m *mocks.API) {
	m.On("GetBareMetalDeviceIdResource", int32(14730)).Return(
		&hv.BareMetalDevice{
			Hostname:                 "",
			PrimaryIp:                "66.165.243.74",
			CustomIPXEScriptURL:      "",
			LocationName:             "LAX2",
			ServiceId:                0,
			DeviceId:                 14730,
			ProductName:              "",
			VlanId:                   0,
			Period:                   "",
			PublicSshKeyId:           0,
			Script:                   "",
			PowerStatus:              "ON",
			CustomIPXEScriptContents: "",
			OrderId:                  0,
			OsName:                   "",
			ProductId:                0,
			Tags:                     []string{"instance-type=bare-metal-x"},
		},
		nil)

	m.On("GetBareMetalDeviceIdResource", int32(9999999)).Return(
		nil, client.ErrNoSuchDevice)
}
func Test_InstanceExists(t *testing.T) {
	i2, m := newHVInstanceV2(t)
	node := newNode()
	ctx := context.Background()
	standardMocks(m)
	tests := []struct {
		providerId string
		wantBool bool
		wantErrString string
	}{
		{fmt.Sprintf("hivelocity://%d", mockDeviceId), true, ""},
		{"hivelocity://9999999", false, ""},
		{"hivelocity://9999999999999999999999999999", false, "GetHivelocityDeviceIdFromNode(node) failed: failed to convert node.Spec.ProviderID \"hivelocity://9999999999999999999999999999\" to int32"},
	}
	for _, row := range tests {
		node.Spec.ProviderID = row.providerId
		gotBool, gotErr := i2.InstanceExists(ctx, node)
		require.Equal(t, row.wantBool, gotBool, row.providerId)
		if row.wantErrString != ""{
			require.Error(t, gotErr)
			require.Equal(t, row.wantErrString, gotErr.Error())
		} else {
			require.NoError(t, gotErr)
		}
	}
}

func Test_InstanceShutdown(t *testing.T) {
	i2, m := newHVInstanceV2(t)
	standardMocks(m)
	node := newNode()
	ctx := context.Background()
	isDown, err := i2.InstanceShutdown(ctx, node)
	require.False(t, isDown)
	require.NoError(t, err)

	node.Spec.ProviderID = "hivelocity://9999999"
	_, err = i2.InstanceShutdown(ctx, node)
	require.Error(t, err)
	require.Equal(t, "i2.API.GetBareMetalDeviceIdResource(deviceId) failed: no such device", err.Error())
}

func Test_InstanceMetadata(t *testing.T) {
	i2, m := newHVInstanceV2(t)
	node := newNode()
	ctx := context.Background()
	standardMocks(m)
	metaData, err := i2.InstanceMetadata(ctx, node)
	require.NoError(t, err)
	require.Equal(t, &cloudprovider.InstanceMetadata{
		ProviderID: "14730",
		NodeAddresses: []corev1.NodeAddress{
			{
				Type:    corev1.NodeAddressType("ExternalIP"),
				Address: "66.165.243.74",
			},
		},
		Zone:   "LAX2",
		Region: "LAX2",
		InstanceType: "bare-metal-x",
	}, metaData)

	node.Spec.ProviderID = "hivelocity://9999999"
	metaData, err = i2.InstanceMetadata(ctx, node)
	require.Error(t, err)
	require.Equal(t, "i2.API.GetBareMetalDeviceIdResource(deviceId) failed: no such device", err.Error())
	require.Nil(t, metaData)
}
