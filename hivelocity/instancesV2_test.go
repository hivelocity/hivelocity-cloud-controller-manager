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

	"github.com/hexops/autogold"
	hv "github.com/hivelocity/hivelocity-client-go/client"
	"github.com/hivelocity/hivelocity-cloud-controller-manager/client"
	"github.com/hivelocity/hivelocity-cloud-controller-manager/mocks"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	cloudprovider "k8s.io/cloud-provider"
)

func Test_GetHivelocityDeviceIdFromNode(t *testing.T) {
	type args struct {
		node *corev1.Node
	}
	tests := []struct {
		name    string
		args    args
		want    int32
		wantErr bool
	}{
		{
			name: "empty deviceId should fail",
			args: args{
				node: &corev1.Node{},
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "Correct deviceId should get parsed",
			args: args{
				node: &corev1.Node{
					Spec: corev1.NodeSpec{
						ProviderID: "12345",
					},
				},
			},
			want:    12345,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetHivelocityDeviceIdFromNode(tt.args.node)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetHivelocityDeviceIdFromNode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetHivelocityDeviceIdFromNode() = %v, want %v", got, tt.want)
			}
		})
	}
}

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

var mockDeviceId int = 14730

func newHVInstanceV2(t *testing.T) (*HVInstancesV2, *mocks.API) {
	var i2 HVInstancesV2
	client := getAPIClient()
	api := mocks.NewAPI(t)
	i2.API = api
	i2.Client = client
	return &i2, api
}

func newNode() *corev1.Node {
	return &corev1.Node{
		Spec: corev1.NodeSpec{
			ProviderID: strconv.Itoa(mockDeviceId),
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
	i2, _ := newHVInstanceV2(t)
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
	i2, m := newHVInstanceV2(t)
	node := newNode()
	ctx := context.Background()
	standardMocks(m)
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
		Zone:   "LAX2",
		Region: "LAX2",
	}).Equal(t, metaData)

	node.Spec.ProviderID = "9999999"
	metaData, err = i2.InstanceMetadata(ctx, node)
	require.Error(t, err)
	require.Nil(t, metaData)
}

func Test_getInstanceTypeFromTags(t *testing.T) {
	type args struct {
		tags     []string
		deviceId int32
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"empty slice returns empty string", args{[]string{}, 1}, ""},
		{"invalid label value will be skipped", args{[]string{"instance-type=&"}, 1}, ""},
		{"valid label value will be used", args{[]string{"instance-type=abc"}, 1}, "abc"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getInstanceTypeFromTags(tt.args.tags, tt.args.deviceId); got != tt.want {
				t.Errorf("getInstanceTypeFromTags() = %v, want %v", got, tt.want)
			}
		})
	}
}
