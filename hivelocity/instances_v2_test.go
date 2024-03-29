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

const (
	dummyDeviceID   = 12345
	unknownDeviceID = 9999999
	invalidDeviceID = 999999999999999999

	region   = "LAX2"
	nodeName = "myNode"
)

func Test_getHivelocityDeviceIDFromNode(t *testing.T) {
	t.Parallel()
	tests := []struct {
		providerID   string
		wantDeviceID int32
		wantErr      error
	}{
		{
			providerID:   "",
			wantDeviceID: 0,
			wantErr:      errMissingProviderPrefix,
		},
		{
			providerID:   "hivelocity://12345",
			wantDeviceID: dummyDeviceID,
			wantErr:      nil,
		},
	}
	for _, tt := range tests {
		node := newNode(tt.providerID, nodeName)
		gotProviderID, gotErr := getHivelocityDeviceIDFromNode(node)
		msg := fmt.Sprintf("Input: providerID=%q", tt.providerID)
		if tt.wantErr == nil {
			require.NoError(t, gotErr, msg)
		} else {
			require.ErrorIs(t, gotErr, tt.wantErr, msg)
		}
		require.Equal(t, tt.wantDeviceID, gotProviderID, msg)
	}
}

func newNode(providerID, nodeName string) *corev1.Node {
	var node corev1.Node

	if providerID != "" {
		node = corev1.Node{
			Spec: corev1.NodeSpec{
				ProviderID: providerID,
			},
		}
	}
	node.SetName(nodeName)

	return &node
}

func standardMocks(m *mocks.Interface) {
	m.On("GetBareMetalDevice", mock.Anything, int32(dummyDeviceID)).Return(
		&hv.BareMetalDevice{
			Hostname:                 "",
			PrimaryIp:                "66.165.243.74",
			CustomIPXEScriptURL:      "",
			LocationName:             region,
			ServiceId:                0,
			DeviceId:                 dummyDeviceID,
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
			Tags:                     []string{"caphv-device-type=bare-metal-x", "caphv-machine-name=myNode"},
		},
		nil)

	m.On("GetBareMetalDevice", mock.Anything, int32(unknownDeviceID)).Return(
		nil, client.ErrNoSuchDevice)

	m.On("ListDevices", mock.Anything).Return(
		[]hv.BareMetalDevice{
			{
				Hostname:                 "",
				PrimaryIp:                "66.165.243.74",
				CustomIPXEScriptURL:      "",
				LocationName:             region,
				ServiceId:                0,
				DeviceId:                 dummyDeviceID,
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
				Tags:                     []string{"caphv-machine-name=myNode", "caphv-device-type=bare-metal-x"},
			},
		},
		nil)
}

func Test_InstanceExists(t *testing.T) {
	t.Parallel()
	m := mocks.NewInterface(t)

	ctx := context.Background()
	standardMocks(m)
	i2 := newHVInstanceV2(m)

	tests := []struct {
		deviceID int64
		nodeName string
		wantBool bool
		wantErr  error
	}{
		{
			deviceID: dummyDeviceID,
			nodeName: nodeName,
			wantBool: true,
			wantErr:  nil,
		},
		{
			deviceID: unknownDeviceID,
			nodeName: nodeName,
			wantBool: false,
			wantErr:  nil,
		},
		{
			deviceID: invalidDeviceID,
			nodeName: nodeName,
			wantBool: false,
			wantErr:  errFailedToConvertProviderID,
		},
		{
			deviceID: 0,
			nodeName: nodeName,
			wantBool: true,
			wantErr:  nil,
		},
		{
			deviceID: dummyDeviceID,
			nodeName: "unknown",
			wantBool: false,
			wantErr:  nil,
		},
	}
	for _, tt := range tests {
		var node *corev1.Node
		if tt.deviceID == 0 {
			node = newNode("", tt.nodeName)
		} else {
			node = newNode(fmt.Sprintf("hivelocity://%d", tt.deviceID), tt.nodeName)
		}

		gotBool, gotErr := i2.InstanceExists(ctx, node)
		msg := fmt.Sprintf("Input: deviceID=%+v", tt.deviceID)

		if tt.wantErr == nil {
			require.NoError(t, gotErr, msg)
		} else {
			require.ErrorIs(t, gotErr, tt.wantErr, msg)
		}

		require.Equal(t, tt.wantBool, gotBool, msg)
	}
}

func Test_InstanceShutdown(t *testing.T) {
	t.Parallel()
	m := mocks.NewInterface(t)
	ctx := context.Background()
	standardMocks(m)
	i2 := newHVInstanceV2(m)

	tests := []struct {
		deviceID int
		wantBool bool
		wantErr  error
	}{
		{
			deviceID: dummyDeviceID,
			wantBool: false,
			wantErr:  nil,
		},
		{
			deviceID: unknownDeviceID,
			wantBool: false,
			wantErr:  errNoDeviceFound,
		},
		{
			deviceID: 0,
			wantBool: false,
			wantErr:  nil,
		},
	}
	for _, tt := range tests {
		var node *corev1.Node
		if tt.deviceID == 0 {
			node = newNode("", nodeName)
		} else {
			node = newNode(fmt.Sprintf("hivelocity://%d", tt.deviceID), nodeName)
		}

		gotBool, gotErr := i2.InstanceShutdown(ctx, node)
		msg := fmt.Sprintf("Input: deviceID=%+v", tt.deviceID)

		if tt.wantErr == nil {
			require.NoError(t, gotErr, msg)
		} else {
			require.Error(t, gotErr, msg)
			require.ErrorIs(t, gotErr, tt.wantErr, msg)
		}

		require.Equal(t, tt.wantBool, gotBool, msg)
	}
}

func Test_InstanceMetadata(t *testing.T) {
	t.Parallel()
	m := mocks.NewInterface(t)
	ctx := context.Background()
	standardMocks(m)
	i2 := newHVInstanceV2(m)
	tests := []struct {
		deviceID     int
		wantMetaData *cloudprovider.InstanceMetadata
		wantErr      error
	}{
		{
			deviceID: dummyDeviceID,
			wantMetaData: &cloudprovider.InstanceMetadata{
				ProviderID: fmt.Sprint(dummyDeviceID),
				NodeAddresses: []corev1.NodeAddress{
					{
						Type:    corev1.NodeAddressType("ExternalIP"),
						Address: "66.165.243.74",
					},
				},
				Zone:         region,
				Region:       region,
				InstanceType: "bare-metal-x",
			},
			wantErr: nil,
		},
		{
			deviceID:     unknownDeviceID,
			wantMetaData: nil,
			wantErr:      errNoDeviceFound,
		},
		{
			deviceID: 0,
			wantMetaData: &cloudprovider.InstanceMetadata{
				ProviderID: fmt.Sprint(dummyDeviceID),
				NodeAddresses: []corev1.NodeAddress{
					{
						Type:    corev1.NodeAddressType("ExternalIP"),
						Address: "66.165.243.74",
					},
				},
				Zone:         region,
				Region:       region,
				InstanceType: "bare-metal-x",
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		var node *corev1.Node
		if tt.deviceID == 0 {
			node = newNode("", nodeName)
		} else {
			node = newNode(fmt.Sprintf("hivelocity://%d", tt.deviceID), nodeName)
		}

		gotMetaData, gotErr := i2.InstanceMetadata(ctx, node)
		msg := fmt.Sprintf("Input: deviceID=%+v", tt.deviceID)

		if tt.wantErr == nil {
			require.NoError(t, gotErr, msg)
		} else {
			require.ErrorIs(t, gotErr, tt.wantErr, msg)
		}

		require.Equal(t, tt.wantMetaData, gotMetaData, msg)
	}
}
