/*
 *  Copyright (c) 2022, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
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
 */

package utils

import (
	"k8s.io/apimachinery/pkg/types"
	gwapiv1b1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

// ExtractExtensions extract extensions of the http route.
func ExtractExtensions(httpRoute *gwapiv1b1.HTTPRoute) []types.NamespacedName {
	authentications := []types.NamespacedName{}
	for _, rule := range httpRoute.Spec.Rules {
		for _, filter := range rule.Filters {
			if filter.ExtensionRef != nil && filter.ExtensionRef.Kind == "Authentication" {
				namespacedName := types.NamespacedName{
					Name:      string(filter.ExtensionRef.Name),
					Namespace: httpRoute.Namespace}
				authentications = append(authentications, namespacedName)
			}
		}
	}
	return authentications
}
