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

func Test_GetHivelocityDeviceIDFromNode(t *testing.T) {
	t.Parallel()
	tests := []struct {
		providerID    string
		wantDeviceID  int32
		wantErrString string
	}{
		{
			providerID:    "",
			wantDeviceID:  0,
			wantErrString: "missing prefix \"hivelocity://\" in node.Spec.ProviderID \"\"",
		},
		{
			providerID:    "hivelocity://12345",
			wantDeviceID:  12345,
			wantErrString: "",
		},
	}
	node := newNode()
	for _, tt := range tests {
		node.Spec.ProviderID = tt.providerID
		gotProviderID, gotErr := getHivelocityDeviceIDFromNode(node)
		msg := fmt.Sprintf("Input: providerID=%q", tt.providerID)
		if tt.wantErrString == "" {
			require.NoError(t, gotErr, msg)
		} else {
			require.Error(t, gotErr, msg)
			require.Equal(t, tt.wantErrString, gotErr.Error(), msg)
		}
		require.Equal(t, tt.wantDeviceID, gotProviderID, msg)
	}
}

func newNode() *corev1.Node {
	return &corev1.Node{
		Spec: corev1.NodeSpec{
			ProviderID: "hivelocity://14730",
		},
	}
}

func standardMocks(m *mocks.Client) {
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
	t.Parallel()
	m := mocks.NewClient(t)

	node := newNode()
	ctx := context.Background()
	standardMocks(m)
	i2 := NewHVInstanceV2(m)

	tests := []struct {
		providerID    int64
		wantBool      bool
		wantErrString string
	}{
		{
			providerID:    14730,
			wantBool:      true,
			wantErrString: "",
		},
		{
			providerID:    9999999,
			wantBool:      false,
			wantErrString: "",
		},
		{
			providerID: 999999999999999999,
			wantBool:   false,
			wantErrString: "GetHivelocityDeviceIDFromNode(node) failed: " +
				"failed to convert node.Spec.ProviderID \"hivelocity://999999999999999999\" to int32",
		},
	}
	for _, tt := range tests {
		node.Spec.ProviderID = fmt.Sprintf("hivelocity://%d", tt.providerID)
		gotBool, gotErr := i2.InstanceExists(ctx, node)
		msg := fmt.Sprintf("Input: providerID=%+v", tt.providerID)
		if tt.wantErrString == "" {
			require.NoError(t, gotErr, msg)
		} else {
			require.Error(t, gotErr, msg)
			require.Equal(t, tt.wantErrString, gotErr.Error(), msg)
		}
		require.Equal(t, tt.wantBool, gotBool, tt.providerID, msg)
	}
}

func Test_InstanceShutdown(t *testing.T) {
	t.Parallel()
	m := mocks.NewClient(t)
	node := newNode()
	ctx := context.Background()
	standardMocks(m)
	i2 := NewHVInstanceV2(m)

	tests := []struct {
		providerID    int
		wantBool      bool
		wantErrString string
	}{
		{
			providerID:    14730,
			wantBool:      false,
			wantErrString: "",
		},
		{
			providerID:    9999999,
			wantBool:      false,
			wantErrString: "i2.API.GetBareMetalDeviceIdResource(deviceID) failed: no such device",
		},
	}
	for _, tt := range tests {
		node.Spec.ProviderID = fmt.Sprintf("hivelocity://%d", tt.providerID)
		gotBool, gotErr := i2.InstanceShutdown(ctx, node)
		msg := fmt.Sprintf("Input: providerID=%+v", tt.providerID)
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
	t.Parallel()
	m := mocks.NewClient(t)
	node := newNode()
	ctx := context.Background()
	standardMocks(m)
	i2 := NewHVInstanceV2(m)
	tests := []struct {
		providerID    int
		wantMetaData  *cloudprovider.InstanceMetadata
		wantErrString string
	}{
		{
			providerID: 14730,
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
			providerID:    9999999,
			wantMetaData:  nil,
			wantErrString: "i2.API.GetBareMetalDeviceIdResource(deviceID) failed: no such device",
		},
	}
	for _, tt := range tests {
		node.Spec.ProviderID = fmt.Sprintf("hivelocity://%d", tt.providerID)
		gotMetaData, gotErr := i2.InstanceMetadata(ctx, node)
		msg := fmt.Sprintf("Input: providerID=%+v", tt.providerID)
		if tt.wantErrString == "" {
			require.NoError(t, gotErr, msg)
		} else {
			require.Error(t, gotErr, msg)
			require.Equal(t, tt.wantErrString, gotErr.Error(), msg)
		}
		require.Equal(t, tt.wantMetaData, gotMetaData, msg)
	}
}
