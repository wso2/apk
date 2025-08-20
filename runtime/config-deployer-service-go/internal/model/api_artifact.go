/* Copyright (c) 2025 WSO2 LLC. (http://www.wso2.com) All Rights Reserved.
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package model

import (
	egv1a1 "github.com/envoyproxy/gateway/api/v1alpha1"
	dpv2alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v2alpha1"
	corev1 "k8s.io/api/core/v1"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
	gwapiv1alpha3 "sigs.k8s.io/gateway-api/apis/v1alpha3"
)

// APIArtifact represents the API artifact containing all the necessary information
type APIArtifact struct {
	Name                        string                                     `json:"name" yaml:"name"`
	Version                     string                                     `json:"version" yaml:"version"`
	RouteMetadata               *dpv2alpha1.RouteMetadata                  `json:"routeMetadata" yaml:"routeMetadata"`
	ProductionHttpRoutes        []*gwapiv1.HTTPRoute                       `json:"productionHttpRoutes" yaml:"productionHttpRoutes"`
	SandboxHttpRoutes           []*gwapiv1.HTTPRoute                       `json:"sandboxHttpRoutes" yaml:"sandboxHttpRoutes"`
	ProductionGqlRoutes         []*gwapiv1.HTTPRoute                       `json:"productionGqlRoutes" yaml:"productionGqlRoutes"`
	ProductionGqlRoutePolicies  []*dpv2alpha1.RoutePolicy                  `json:"productionGqlRoutePolicies" yaml:"productionGqlRoutePolicies"`
	SandboxGqlRoutes            []*gwapiv1.HTTPRoute                       `json:"sandboxGqlRoutes" yaml:"sandboxGqlRoutes"`
	SandboxGqlRoutePolicies     []*dpv2alpha1.RoutePolicy                  `json:"sandboxGqlRoutePolicies" yaml:"sandboxGqlRoutePolicies"`
	ProductionGrpcRoutes        []*gwapiv1.GRPCRoute                       `json:"productionGrpcRoutes" yaml:"productionGrpcRoutes"`
	SandboxGrpcRoutes           []*gwapiv1.GRPCRoute                       `json:"sandboxGrpcRoutes" yaml:"sandboxGrpcRoutes"`
	Definition                  *corev1.ConfigMap                          `json:"definition,omitempty" yaml:"definition,omitempty"`
	EndpointCertificates        map[string]*corev1.ConfigMap               `json:"endpointCertificates" yaml:"endpointCertificates"`
	CertificateMap              map[string]string                          `json:"certificateMap" yaml:"certificateMap"`
	BackendServices             map[string]*egv1a1.Backend                 `json:"backendServices" yaml:"backendServices"`
	BackendTLSPolicies          map[string]*gwapiv1alpha3.BackendTLSPolicy `json:"backendTLSPolicies" yaml:"backendTLSPolicies"`
	// BackendSecurity             map[string]*v1alpha1.BackendJWT            `json:"backendSecurity" yaml:"backendSecurity"`
	AuthenticationMap           map[string]*egv1a1.SecurityPolicy          `json:"authenticationMap" yaml:"authenticationMap"`
	Scopes                      map[string]*egv1a1.SecurityPolicy          `json:"scopes" yaml:"scopes"`
	RateLimitPolicies           map[string]*egv1a1.BackendTrafficPolicy    `json:"rateLimitPolicies" yaml:"rateLimitPolicies"`
	AIRatelimitPolicies         map[string]*egv1a1.BackendTrafficPolicy    `json:"aiRatelimitPolicies" yaml:"aiRatelimitPolicies"`
	AIRatelimitRoutePolicies    map[string]*dpv2alpha1.RoutePolicy         `json:"aiRatelimitRoutePolicies" yaml:"aiRatelimitRoutePolicies"`
	ApiPolicies                 map[string]*dpv2alpha1.RoutePolicy         `json:"apiPolicies" yaml:"apiPolicies"`
	InterceptorServices         map[string]*egv1a1.EnvoyExtensionPolicy    `json:"interceptorServices" yaml:"interceptorServices"`
	SandboxEndpointAvailable    bool                                       `json:"sandboxEndpointAvailable" yaml:"sandboxEndpointAvailable"`
	ProductionUrl               []string                                   `json:"productionUrl,omitempty" yaml:"productionUrl,omitempty"`
	SandboxUrl                  []string                                   `json:"sandboxUrl,omitempty" yaml:"sandboxUrl,omitempty"`
	ProductionEndpointAvailable bool                                       `json:"productionEndpointAvailable" yaml:"productionEndpointAvailable"`
	UniqueID                    string                                     `json:"uniqueId" yaml:"uniqueId"`
	Secrets                     *corev1.Secret                             `json:"secrets" yaml:"secrets"`
	BackendJwt                  *dpv2alpha1.RoutePolicy                    `json:"backendJwt,omitempty" yaml:"backendJwt,omitempty"`
	Namespace                   string                                     `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Organization                string                                     `json:"organization" yaml:"organization"`
}
