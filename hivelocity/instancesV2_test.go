package hivelocity

import (
	"context"
	"log"
	"os"
	"strconv"
	"testing"

	hv "github.com/hivelocity/hivelocity-client-go/client"
	"github.com/stretchr/testify/assert"

	"github.com/joho/godotenv"
	v1 "k8s.io/api/core/v1"
)

func Test_getHivelocityDeviceIdFromNode(t *testing.T) {
	type args struct {
		node *v1.Node
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
				node: &v1.Node{},
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "Correct deviceID should get parsed",
			args: args{
				node: &v1.Node{
					Spec: v1.NodeSpec{
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
				t.Errorf("getHivelocityDeviceIdFromNode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getHivelocityDeviceIdFromNode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getAPIClient() (hv.APIClient, context.Context) {
	err := godotenv.Overload("../.envrc")
	if err != nil {
		panic(err)
	}

	apiKey := os.Getenv("HIVELOCITY_API_KEY")
	if apiKey == "" {
		log.Fatalln("Missing environment variable HIVELOCITY_API_KEY")
		os.Exit(1)
	}
	ctx := context.WithValue(context.Background(), hv.ContextAPIKey, hv.APIKey{
		Key: apiKey,
	})
	return *hv.NewAPIClient(hv.NewConfiguration()), ctx
}

var deviceID int = 14730

func Test_InstanceExists(t *testing.T) {
	var i2 hvInstancesV2
	client, ctx := getAPIClient()
	i2.client = &client
	node := v1.Node{
		Spec: v1.NodeSpec{
			ProviderID: strconv.Itoa(deviceID),
		},
	}
	myBool, err := i2.InstanceExists(ctx, &node)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, myBool)

	node.Spec.ProviderID = "9999999"
	myBool, err = i2.InstanceExists(ctx, &node)
	assert.Equal(t, false, myBool)
	assert.Equal(t, nil, err)

	node.Spec.ProviderID = "9999999999999999999999999999"
	myBool, err = i2.InstanceExists(ctx, &node)
	assert.Equal(t, false, myBool)
	assert.Equal(t, "failed to convert node.Spec.ProviderID \"9999999999999999999999999999\" to int32.", err.Error())
}
