/*
Copyright 2024.

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

// ResolveSubscriptionRatelimitPolicy defines the structure to resolve subscription rate limit policies.
type ResolveSubscriptionRatelimitPolicy struct {
	Name             string              `json:"name"`
	StopOnQuotaReach bool                `json:"stopOnQuotaReach"`
	Organization     string              `json:"organization"`
	RequestCount     ResolveRequestCount `json:"requestCount,omitempty"`
	BurstControl     ResolveBurstControl `json:"burstControl,omitempty"`
}

// ResolveRequestCount defines the rule for request count quota.
type ResolveRequestCount struct {
	RequestsPerUnit uint32 `json:"requestsPerUnit,omitempty"`
	Unit            string `json:"unit,omitempty"`
}

// ResolveBurstControl defines the rule for token count quota.
type ResolveBurstControl struct {
	RequestsPerUnit uint32 `json:"requestsPerUnit,omitempty"`
	Unit            string `json:"unit,omitempty"`
}
