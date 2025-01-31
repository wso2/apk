/*
 *  Copyright (c) 2024, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 *  This file contains code derived from Envoy Gateway,
 *  https://github.com/envoyproxy/gateway
 *  and is provided here subject to the following:
 *  Copyright Project Envoy Gateway Authors
 *
 */

package provider

import (
	dpv1alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha1"
	dpv1alpha2 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha2"
	dpv1alpha3 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha3"
	dpv1alpha4 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha4"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
	gwapiv1a2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gwapiv1b1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

// GetScheme returns a scheme with types supported by the Kubernetes provider.
func GetScheme() *runtime.Scheme {
	// todo(amali) move this to init method once we remove the old operator
	scheme := runtime.NewScheme()
	if err := clientgoscheme.AddToScheme(scheme); err != nil {
		panic(err)
	}
	// Add Gateway API types.
	if err := gwapiv1.AddToScheme(scheme); err != nil {
		panic(err)
	}
	if err := gwapiv1b1.AddToScheme(scheme); err != nil {
		panic(err)
	}
	if err := gwapiv1a2.AddToScheme(scheme); err != nil {
		panic(err)
	}
	if err := dpv1alpha1.AddToScheme(scheme); err != nil {
		panic(err)
	}
	if err := dpv1alpha2.AddToScheme(scheme); err != nil {
		panic(err)
	}
	if err := dpv1alpha3.AddToScheme(scheme); err != nil {
		panic(err)
	}
	if err := dpv1alpha4.AddToScheme(scheme); err != nil {
		panic(err)
	}

	return scheme
}
