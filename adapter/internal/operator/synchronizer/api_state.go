/*
 *  Copyright (c) 2022, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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

package synchronizer

import (
	"github.com/wso2/apk/common-go-libs/apis/dp/v1alpha1"
	"github.com/wso2/apk/common-go-libs/apis/dp/v1alpha2"
	"github.com/wso2/apk/common-go-libs/apis/dp/v1alpha3"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
)

// APIState holds the state of the deployed APIs. This state is compared with
// the state of the Kubernetes controller cache to detect updates.
// +k8s:deepcopy-gen=true
type APIState struct {
	APIDefinition             *v1alpha2.API
	ProdHTTPRoute             *HTTPRouteState
	SandHTTPRoute             *HTTPRouteState
	ProdGQLRoute              *GQLRouteState
	SandGQLRoute              *GQLRouteState
	Authentications           map[string]v1alpha2.Authentication
	RateLimitPolicies         map[string]v1alpha3.RateLimitPolicy
	ResourceAuthentications   map[string]v1alpha2.Authentication
	ResourceRateLimitPolicies map[string]v1alpha3.RateLimitPolicy
	ResourceAPIPolicies       map[string]v1alpha3.APIPolicy
	APIPolicies               map[string]v1alpha3.APIPolicy
	AIProvider                *v1alpha3.AIProvider
	InterceptorServiceMapping map[string]v1alpha1.InterceptorService
	BackendJWTMapping         map[string]v1alpha1.BackendJWT
	APIDefinitionFile         []byte
	SubscriptionValidation    bool
	MutualSSL                 *v1alpha2.MutualSSL
}

// HTTPRouteState holds the state of the deployed httpRoutes. This state is compared with
// the state of the Kubernetes controller cache to detect updates.
// +k8s:deepcopy-gen=true
type HTTPRouteState struct {
	HTTPRouteCombined   *gwapiv1.HTTPRoute
	HTTPRoutePartitions map[string]*gwapiv1.HTTPRoute
	BackendMapping      map[string]*v1alpha2.ResolvedBackend
	Scopes              map[string]v1alpha1.Scope
	RuleIdxToAiRatelimitPolicyMapping map[int]*v1alpha3.AIRateLimitPolicy
}

// GQLRouteState holds the state of the deployed gqlRoutes. This state is compared with
// the state of the Kubernetes controller cache to detect updates.
// +k8s:deepcopy-gen=true
type GQLRouteState struct {
	GQLRouteCombined   *v1alpha2.GQLRoute
	GQLRoutePartitions map[string]*v1alpha2.GQLRoute
	BackendMapping     map[string]*v1alpha2.ResolvedBackend
	Scopes             map[string]v1alpha1.Scope
}
