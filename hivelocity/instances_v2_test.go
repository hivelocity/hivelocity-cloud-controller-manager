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
	"github.com/hivelocity/hivelocity-cloud-controller-manager/client/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	cloudprovider "k8s.io/cloud-provider"
)

func Test_GetHivelocityDeviceIdFromNode(t *testing.T) {
	var tests = []struct {
		providerId    string
		wantDeviceId  int32
		wantErrString string
	}{
		{
			providerId:    "",
			wantDeviceId:  0,
			wantErrString: "missing prefix \"hivelocity://\" in node.Spec.ProviderID \"\"",
		},
		{
			providerId:    "hivelocity://12345",
			wantDeviceId:  12345,
			wantErrString: "",
		},
	}
	var node = &corev1.Node{}
	for _, tt := range tests {
		node.Spec.ProviderID = tt.providerId
		gotProviderId, gotErr := GetHivelocityDeviceIdFromNode(node)
		msg := fmt.Sprintf("Input: providerId=%q", tt.providerId)
		if tt.wantErrString == "" {
			require.NoError(t, gotErr, msg)
		} else {
			require.Error(t, gotErr, msg)
			require.Equal(t, tt.wantErrString, gotErr.Error(), msg)
		}
		require.Equal(t, tt.wantDeviceId, gotProviderId, msg)
	}
}

var mockDeviceId int = 14730

func newTestData(m *mocks.API) (*HVInstancesV2, *corev1.Node, context.Context) {
	node := newNode()
	ctx := context.Background()
	standardMocks(m)
	return &HVInstancesV2{
		API: m,
	}, node, ctx
}

func newNode() *corev1.Node {
	return &corev1.Node{
		Spec: corev1.NodeSpec{
			ProviderID: fmt.Sprintf("hivelocity://%d", mockDeviceId),
		},
	}
}

func standardMocks(m *mocks.API) {
	m.On("GetBareMetalDevice", mock.Anything, int32(14730)).Return(
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

	m.On("GetBareMetalDevice", mock.Anything, int32(9999999)).Return(
		nil, client.ErrNoSuchDevice)
}
func Test_InstanceExists(t *testing.T) {
	m := mocks.NewAPI(t)
	i2, node, ctx := newTestData(m)

	tests := []struct {
		providerId    int64
		wantBool      bool
		wantErrString string
	}{
		{
			providerId:    int64(mockDeviceId),
			wantBool:      true,
			wantErrString: "",
		},
		{
			providerId:    9999999,
			wantBool:      false,
			wantErrString: "",
		},
		{
			providerId:    999999999999999999,
			wantBool:      false,
			wantErrString: "GetHivelocityDeviceIdFromNode(node) failed: failed to convert node.Spec.ProviderID \"hivelocity://999999999999999999\" to int32",
		},
	}
	for _, tt := range tests {
		node.Spec.ProviderID = fmt.Sprintf("hivelocity://%d", tt.providerId)
		gotBool, gotErr := i2.InstanceExists(ctx, node)
		msg := fmt.Sprintf("Input: providerId=%+v", tt.providerId)
		if tt.wantErrString == "" {
			require.NoError(t, gotErr, msg)
		} else {
			require.Error(t, gotErr, msg)
			require.Equal(t, tt.wantErrString, gotErr.Error(), msg)
		}
		require.Equal(t, tt.wantBool, gotBool, tt.providerId, msg)
	}
}

func Test_InstanceShutdown(t *testing.T) {
	m := mocks.NewAPI(t)
	i2, node, ctx := newTestData(m)
	tests := []struct {
		providerId    int
		wantBool      bool
		wantErrString string
	}{
		{
			providerId:    mockDeviceId,
			wantBool:      false,
			wantErrString: "",
		},
		{
			providerId:    9999999,
			wantBool:      false,
			wantErrString: "i2.API.GetBareMetalDeviceIdResource(deviceId) failed: no such device",
		},
	}
	for _, tt := range tests {
		node.Spec.ProviderID = fmt.Sprintf("hivelocity://%d", tt.providerId)
		gotBool, gotErr := i2.InstanceShutdown(ctx, node)
		msg := fmt.Sprintf("Input: providerId=%+v", tt.providerId)
		if tt.wantErrString == "" {
			require.NoError(t, gotErr, msg)
		} else {
			require.Error(t, gotErr, msg)
			require.Equal(t, tt.wantErrString, gotErr.Error(), msg)
		}
		require.Equal(t, tt.wantBool, gotBool, msg)
	}
}

func Test_InstanceMetadata(t *testing.T) {
	m := mocks.NewAPI(t)
	i2, node, ctx := newTestData(m)
	tests := []struct {
		providerId    int
		wantMetaData  *cloudprovider.InstanceMetadata
		wantErrString string
	}{
		{
			providerId: mockDeviceId,
			wantMetaData: &cloudprovider.InstanceMetadata{
				ProviderID: "14730",
				NodeAddresses: []corev1.NodeAddress{
					{
						Type:    corev1.NodeAddressType("ExternalIP"),
						Address: "66.165.243.74",
					},
				},
				Zone:         "LAX2",
				Region:       "LAX2",
				InstanceType: "bare-metal-x",
			},
			wantErrString: "",
		},
		{
			providerId:    9999999,
			wantMetaData:  nil,
			wantErrString: "i2.API.GetBareMetalDeviceIdResource(deviceId) failed: no such device",
		},
	}
	for _, tt := range tests {
		node.Spec.ProviderID = fmt.Sprintf("hivelocity://%d", tt.providerId)
		gotMetaData, gotErr := i2.InstanceMetadata(ctx, node)
		msg := fmt.Sprintf("Input: providerId=%+v", tt.providerId)
		if tt.wantErrString == "" {
			require.NoError(t, gotErr, msg)
		} else {
			require.Error(t, gotErr, msg)
			require.Equal(t, tt.wantErrString, gotErr.Error(), msg)
		}
		require.Equal(t, tt.wantMetaData, gotMetaData, msg)
	}
}
