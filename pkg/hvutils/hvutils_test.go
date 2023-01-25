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
	t.Parallel()
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
			err:  ErrMoreThanOneTagFound,
		},
		{
			name: "invalid label value will be skipped",
			tags: []string{"instance-type=&"},
			want: "",
			err:  ErrInvalidLabelValue,
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
			err:  ErrMoreThanOneTagFound,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := GetInstanceTypeFromTags(tt.tags)
			assert.Equal(t, tt.want, got, fmt.Sprintf("tags: %v", tt.tags))
			assert.Equal(t, tt.err, err, fmt.Sprintf("tags: %v", tt.tags))
		})
	}
}
