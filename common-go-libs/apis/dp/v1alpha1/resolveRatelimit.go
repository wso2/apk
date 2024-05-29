/*
Copyright 2023.

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
package v1alpha1

import (
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
)

// ResolveRateLimitAPIPolicy defines the desired state of Policy
type ResolveRateLimitAPIPolicy struct {
	API          ResolveRateLimit  `json:"api,omitempty"`
	Resources    []ResolveResource `json:"resourceList,omitempty"`
	Organization string            `json:"organization,omitempty"`
	BasePath     string            `json:"basePath,omitempty"`
	UUID         string            `json:"uuid,omitempty"`
	Environment  string            `json:"environment,omitempty"`
}

// ResolveRateLimit is the rate limit value for the applied policy
type ResolveRateLimit struct {
	// RequestPerUnit is the number of requests allowed per unit time
	//
	RequestsPerUnit uint32 `json:"requestsPerUnit,omitempty"`

	// Unit is the unit of the requestsPerUnit
	//
	// +kubebuilder:validation:Enum=Minute;Hour;Day
	Unit string `json:"unit,omitempty"`
}

// ResolveResource defines the desired state of Resource
type ResolveResource struct {
	ResourceRatelimit ResolveRateLimit      `json:"resourceRatelimit,omitempty"`
	Path              string                `json:"path,omitempty"`
	PathMatchType     gwapiv1.PathMatchType `json:"pathMatchType,omitempty"`
	Method            string                `json:"method,omitempty"`
}
