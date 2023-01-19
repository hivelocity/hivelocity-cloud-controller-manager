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
	"strings"

	"k8s.io/apimachinery/pkg/util/validation"
)

func GetInstanceTypeFromTags(tags []string, deviceId int32) (string, error) {
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
		return "", fmt.Errorf("No instance-type tag found on deviceId=%d", deviceId)
	}
	if len(instanceTypes) > 1 {
		return "", fmt.Errorf("More than one instance-type tag found on deviceId=%d: %v", deviceId,
			instanceTypes)
	}
	instanceType := instanceTypes[0]

	if errs := validation.IsValidLabelValue(instanceType); len(errs) != 0 {
		return "", fmt.Errorf("deviceID=%d has invalid tag %q %s", deviceId, instanceType,
			strings.Join(errs, "; "))
	}
	return instanceType, nil
}
