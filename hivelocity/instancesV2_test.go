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
		providerId     string
		wantProviderId int32
		wantErrString  string
	}{
		{
			providerId:     "",
			wantProviderId: 0,
			wantErrString:  "xxmissing prefix \"hivelocity://\" in node.Spec.ProviderID \"\"",
		},
		{
			providerId:     "hivelocity://12345",
			wantProviderId: 12345,
			wantErrString:  "",
		},
	}
	var node = &corev1.Node{}
	for _, row := range tests {
		node.Spec.ProviderID = row.providerId
		gotProviderId, gotErr := GetHivelocityDeviceIdFromNode(node)
		msg := fmt.Sprintf("Input: providerId=%q", row.providerId)
		if row.wantErrString == "" {
			require.NoError(t, gotErr, msg)
		} else {
			require.Error(t, gotErr, msg)
			require.Equal(t, row.wantErrString, gotErr.Error(), msg)
		}
		require.Equal(t, row.wantProviderId, gotProviderId, msg)
	}
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
	i2, m := newHVInstanceV2(t)
	node := newNode()
	ctx := context.Background()
	standardMocks(m)
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
			wantErrString: "GxetHivelocityDeviceIdFromNode(node) failed: failed to convert node.Spec.ProviderID \"hivelocity://999999999999999999\" to int32",
		},
	}
	for _, row := range tests {
		node.Spec.ProviderID = fmt.Sprintf("hivelocity://%d", row.providerId)
		gotBool, gotErr := i2.InstanceExists(ctx, node)
		msg := fmt.Sprintf("Input: providerId=%+v", row.providerId)
		if row.wantErrString == "" {
			require.NoError(t, gotErr, msg)
		} else {
			require.Error(t, gotErr, msg)
			require.Equal(t, row.wantErrString, gotErr.Error(), msg)
		}
		require.Equal(t, row.wantBool, gotBool, row.providerId, msg)
	}
}

func Test_InstanceShutdown(t *testing.T) {
	i2, m := newHVInstanceV2(t)
	standardMocks(m)
	node := newNode()
	ctx := context.Background()
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
	for _, row := range tests {
		node.Spec.ProviderID = fmt.Sprintf("hivelocity://%d", row.providerId)
		gotBool, gotErr := i2.InstanceShutdown(ctx, node)
		msg := fmt.Sprintf("Input: providerId=%+v", row.providerId)
		if row.wantErrString == "" {
			require.NoError(t, gotErr, msg)
		} else {
			require.Error(t, gotErr, msg)
			require.Equal(t, row.wantErrString, gotErr.Error(), msg)
		}
		require.Equal(t, row.wantBool, gotBool, msg)
	}
}

func Test_InstanceMetadata(t *testing.T) {
	i2, m := newHVInstanceV2(t)
	node := newNode()
	ctx := context.Background()
	standardMocks(m)
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
	for _, row := range tests {
		node.Spec.ProviderID = fmt.Sprintf("hivelocity://%d", row.providerId)
		gotMetaData, gotErr := i2.InstanceMetadata(ctx, node)
		msg := fmt.Sprintf("Input: providerId=%+v", row.providerId)
		if row.wantErrString == "" {
			require.NoError(t, gotErr, msg)
		} else {
			require.Error(t, gotErr, msg)
			require.Equal(t, row.wantErrString, gotErr.Error(), msg)
		}
		require.Equal(t, row.wantMetaData, gotMetaData, msg)
	}
}
