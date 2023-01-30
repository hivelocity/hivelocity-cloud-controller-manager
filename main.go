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

// Based on: https://github.com/kubernetes/cloud-provider/blob/master/sample/basic_main.go

// The external controller manager is responsible for running controller loops that
// are cloud provider dependent. It uses the API to listen to new events on resources.

// Package main provides the executable to start the cloud controller manager.
package main

import (
	"os"

	_ "github.com/hivelocity/hivelocity-cloud-controller-manager/hivelocity"
	"k8s.io/apimachinery/pkg/util/wait"
	cloudprovider "k8s.io/cloud-provider"
	"k8s.io/cloud-provider/app"
	cloudcontrollerconfig "k8s.io/cloud-provider/app/config"
	"k8s.io/cloud-provider/options"
	"k8s.io/component-base/cli"
	cliflag "k8s.io/component-base/cli/flag"
	_ "k8s.io/component-base/metrics/prometheus/clientgo" // load all the prometheus client-go plugins
	_ "k8s.io/component-base/metrics/prometheus/version"  // for version metric registration
	"k8s.io/klog/v2"
)

func main() {
	ccmOptions, err := options.NewCloudControllerManagerOptions()
	if err != nil {
		klog.Fatalf("unable to initialize command options: %v", err)
	}

	controllerInitializers := app.DefaultInitFuncConstructors

	fss := cliflag.NamedFlagSets{}
	command := app.NewCloudControllerManagerCommand(
		ccmOptions,
		cloudInitializer,
		controllerInitializers,
		fss,
		wait.NeverStop,
	)
	code := cli.Run(command)
	os.Exit(code)
}

func cloudInitializer(config *cloudcontrollerconfig.CompletedConfig) cloudprovider.Interface { //nolint:ireturn,revive // implements InitCloudFunc
	cloudConfig := config.ComponentConfig.KubeCloudShared.CloudProvider
	cloud, err := cloudprovider.InitCloudProvider(cloudConfig.Name, cloudConfig.CloudConfigFile)
	if err != nil {
		klog.Fatalf("Cloud provider could not be initialized: %v", err)
	}
	if cloud == nil {
		klog.Fatalf("Cloud provider is nil")
	}

	/*
		AFAIK this HasClusterID() is not needed
		https://github.com/kubernetes/cloud-provider/issues/12
		if !cloud.HasClusterID() {
			if config.ComponentConfig.KubeCloudShared.AllowUntaggedCloud {
				klog.Warning("detected a cluster without a ClusterID. " +
				    "A ClusterID will be required in the future. "+
					"Please tag your cluster to avoid any future issues")
			} else {
				klog.Fatalf("no ClusterID found. A ClusterID is required " +
				    "for the cloud provider to function properly. " +
					"This check can be bypassed by setting the allow-untagged-cloud option")
			}
		}
	*/

	return cloud
}
