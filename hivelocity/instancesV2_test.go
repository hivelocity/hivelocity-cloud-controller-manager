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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/joho/godotenv"
	corev1 "k8s.io/api/core/v1"
)

func Test_getHivelocityDeviceIdFromNode(t *testing.T) {
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
			got, err := getHivelocityDeviceIdFromNode(tt.args.node)
			if (err != nil) != tt.wantErr {
				// TODO: lieber assert
				t.Errorf("getHivelocityDeviceIdFromNode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getHivelocityDeviceIdFromNode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getAPIClient() (*hv.APIClient, context.Context) {
	err := godotenv.Overload("../.envrc")
	if err != nil {
		panic(err)
	}

	apiKey := os.Getenv("HIVELOCITY_API_KEY")
	if apiKey == "" {
		panic("Missing environment variable HIVELOCITY_API_KEY")
	}
	ctx := context.WithValue(context.Background(), hv.ContextAPIKey, hv.APIKey{
		Key: apiKey,
	})
	return hv.NewAPIClient(hv.NewConfiguration()), ctx
}

var deviceID int = 14730

func Test_InstanceExists(t *testing.T) {
	var i2 hvInstancesV2
	client, ctx := getAPIClient()
	i2.client = client
	node := corev1.Node{
		Spec: corev1.NodeSpec{
			ProviderID: strconv.Itoa(deviceID),
		},
	}
	myBool, err := i2.InstanceExists(ctx, &node)
	require.NoError(t, err)
	assert.Equal(t, nil, err) // TODO require --> abbruch.
	assert.Equal(t, true, myBool)

	node.Spec.ProviderID = "9999999"
	myBool, err = i2.InstanceExists(ctx, &node)
	assert.Equal(t, false, myBool)
	assert.Equal(t, nil, err)

	node.Spec.ProviderID = "9999999999999999999999999999"
	myBool, err = i2.InstanceExists(ctx, &node)
	assert.Equal(t, false, myBool)
	assert.Equal(t, "failed to convert node.Spec.ProviderID \"9999999999999999999999999999\" to int32", err.Error())
}
