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

	"github.com/joho/godotenv"
	corev1 "k8s.io/api/core/v1"
	"github.com/hivelocity/hivelocity-cloud-controller-manager/hivelocity"
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
			name: "empty deviceID should fail",
			args: args{
				node: &corev1.Node{},
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "Correct deviceID should get parsed",
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
			got, err := hivelocity.GetHivelocityDeviceIdFromNode(tt.args.node)
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

var deviceID int = 14730

func Test_InstanceExists(t *testing.T) {
	var i2 hivelocity.HVInstancesV2
	client := getAPIClient()
	i2.Client = client
	node := corev1.Node{
		Spec: corev1.NodeSpec{
			ProviderID: strconv.Itoa(deviceID),
		},
	}
	ctx := context.Background()
	myBool, err := i2.InstanceExists(ctx, &node)
	require.NoError(t, err)
	require.Equal(t, true, myBool)

	node.Spec.ProviderID = "9999999"
	myBool, err = i2.InstanceExists(ctx, &node)
	require.Equal(t, false, myBool)
	require.Equal(t, nil, err)

	node.Spec.ProviderID = "9999999999999999999999999999"
	myBool, err = i2.InstanceExists(ctx, &node)
	require.Equal(t, false, myBool)
	require.Equal(t, "failed to convert node.Spec.ProviderID \"9999999999999999999999999999\" to int32", err.Error())
}
