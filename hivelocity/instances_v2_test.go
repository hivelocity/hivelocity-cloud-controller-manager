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
		node := newNode(tt.providerID)
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

func newNode(providerID string) *corev1.Node {
	node := corev1.Node{
		Spec: corev1.NodeSpec{
			ProviderID: providerID,
		},
	}
	node.SetName("myNode")
	return &node
}

func standardMocks(m *mocks.Interface) {
	m.On("GetBareMetalDevice", mock.Anything, int32(dummyDeviceID)).Return(
		&hv.BareMetalDevice{
			Hostname:                 "",
			PrimaryIp:                "66.165.243.74",
			CustomIPXEScriptURL:      "",
			LocationName:             "LAX2",
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
			Tags:                     []string{"caphv-device-type=bare-metal-x"},
		},
		nil)

	m.On("GetBareMetalDevice", mock.Anything, int32(unknownDeviceID)).Return(
		nil, client.ErrNoSuchDevice)
}

func Test_InstanceExists(t *testing.T) {
	t.Parallel()
	m := mocks.NewInterface(t)

	ctx := context.Background()
	standardMocks(m)
	i2 := newHVInstanceV2(m)

	tests := []struct {
		deviceID int64
		wantBool bool
		wantErr  error
	}{
		{
			deviceID: dummyDeviceID,
			wantBool: true,
			wantErr:  nil,
		},
		{
			deviceID: unknownDeviceID,
			wantBool: false,
			wantErr:  nil,
		},
		{
			deviceID: invalidDeviceID,
			wantBool: false,
			wantErr:  errFailedToConvertProviderID,
		},
	}
	for _, tt := range tests {
		node := newNode(fmt.Sprintf("hivelocity://%d", tt.deviceID))
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
			wantErr:  client.ErrNoSuchDevice,
		},
	}
	for _, tt := range tests {
		node := newNode(fmt.Sprintf("hivelocity://%d", tt.deviceID))
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
				Zone:         "LAX2",
				Region:       "LAX2",
				InstanceType: "bare-metal-x",
			},
			wantErr: nil,
		},
		{
			deviceID:     unknownDeviceID,
			wantMetaData: nil,
			wantErr:      client.ErrNoSuchDevice,
		},
	}
	for _, tt := range tests {
		node := newNode(fmt.Sprintf("hivelocity://%d", tt.deviceID))
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
