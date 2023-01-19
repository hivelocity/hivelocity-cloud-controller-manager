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

package hvutils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getInstanceTypeFromTags(t *testing.T) {
	tests := []struct {
		name string
		tags []string
		want string
		err  error
	}{
		{
			name: "empty slice returns empty string",
			tags: []string{},
			want: "",
			err:  fmt.Errorf("no instance-type tag found on deviceID=1"),
		},
		{
			name: "invalid label value will be skipped",
			tags: []string{"instance-type=&"},
			want: "",
			err:  fmt.Errorf("deviceID=1 has invalid tag \"&\" a valid label must be an empty string or consist of alphanumeric characters, '-', '_' or '.', and must start and end with an alphanumeric character (e.g. 'MyValue',  or 'my_value',  or '12345', regex used for validation is '(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?')"),
		},
		{
			name: "valid label value will be used",
			tags: []string{"foo", "instance-type=abc", "bar", "key=value"},
			want: "abc",
			err:  nil,
		},
		{
			name: "two labels",
			tags: []string{"instance-type=abc", "instance-type=abc"},
			want: "",
			err:  fmt.Errorf("more than one instance-type tag found on deviceID=1: [abc abc]"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetInstanceTypeFromTags(tt.tags, 1)
			assert.Equal(t, tt.want, got, fmt.Sprintf("tags: %v", tt.tags))
			assert.Equal(t, tt.err, err, fmt.Sprintf("tags: %v", tt.tags))
		})
	}
}
