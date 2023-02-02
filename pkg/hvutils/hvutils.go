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

// Package hvutils provies utility methods to access the Hivelocity API.
package hvutils

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/util/validation"
)

var (

	// ErrMoreThanOneTagFound gets returned if more than one instance-type tag was found via the HV API.
	ErrMoreThanOneTagFound = fmt.Errorf("more than one instance-type tag found")

	// ErrInvalidLabelValue gets returned if the HV tag contains a value which is an invalid K8s label.
	ErrInvalidLabelValue = fmt.Errorf("invalid label value")

	// ErrNoInstanceTypeFound gets returned if no instance-type tag was found via the HV API.
	ErrNoInstanceTypeFound = fmt.Errorf("no instance-type tag found")
)

// GetInstanceTypeFromTags is a utility method to read the instance-type
// from a slice of strings.
// The slice is usually from the Hivelocity API of a device.
// Example: {"instance-type=foo", "other-label"} would return "foo".
func GetInstanceTypeFromTags(tags []string) (string, error) {
	prefix := "instance-type="
	instanceTypes := make([]string, 0, 1)
	for _, tag := range tags {
		if !strings.HasPrefix(tag, prefix) {
			continue
		}

		instanceType := strings.TrimSpace(strings.TrimPrefix(tag, prefix))
		instanceTypes = append(instanceTypes, instanceType)
	}
	if len(instanceTypes) == 0 {
		return "", ErrNoInstanceTypeFound
	}
	if len(instanceTypes) > 1 {
		return "", fmt.Errorf(
			"[GetInstanceTypeFromTags] more than one instance type. instanceTypes %v: %w",
			instanceTypes,
			ErrMoreThanOneTagFound,
		)
	}
	instanceType := instanceTypes[0]

	if errs := validation.IsValidLabelValue(instanceType); len(errs) != 0 {
		return "", fmt.Errorf("[GetInstanceTypeFromTags] Hivelocity tag is no valid K8s label. "+
			"errors %q, instance-type %q: %w",
			strings.Join(errs, "; "),
			instanceType,
			ErrInvalidLabelValue,
		)
	}
	return instanceType, nil
}
