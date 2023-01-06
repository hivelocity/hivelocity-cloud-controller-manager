package hivelocity

import (
	"testing"

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
			want: 0,
			wantErr: true,
		},
		{
			name: "Correct deviceID should get parsed",
			args: args{
				node: &v1.Node{
					Spec:       v1.NodeSpec{
						ProviderID: "12345",
					},
				},
			},
			want: 12345,
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
