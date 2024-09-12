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
package envoyconf_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wso2/apk/adapter/internal/dataholder"
	"github.com/wso2/apk/adapter/internal/discovery/xds"
	"github.com/wso2/apk/adapter/internal/loggers"
	envoy "github.com/wso2/apk/adapter/internal/oasparser/envoyconf"
	"github.com/wso2/apk/adapter/internal/operator/constants"
	"github.com/wso2/apk/adapter/internal/operator/synchronizer"
	operatorutils "github.com/wso2/apk/adapter/internal/operator/utils"
	"github.com/wso2/apk/common-go-libs/apis/dp/v1alpha1"
	"github.com/wso2/apk/common-go-libs/apis/dp/v1alpha2"
	"github.com/wso2/apk/common-go-libs/apis/dp/v1alpha3"
	"github.com/wso2/apk/common-go-libs/apis/dp/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8types "k8s.io/apimachinery/pkg/types"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func TestCreateRoutesWithClustersWithExactAndRegularExpressionRules(t *testing.T) {
	apiState := synchronizer.APIState{}
	apiDefinition := v1alpha3.API{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
			Name:      "test-api-2",
		},
		Spec: v1alpha3.APISpec{
			APIName:    "test-api-2",
			APIVersion: "2.0.0",
			BasePath:   "/test-api/2.0.0",
			Production: []v1alpha3.EnvConfig{
				{
					RouteRefs: []string{
						"test-api-2-prod-http-route",
					},
				},
			},
		},
	}
	apiState.APIDefinition = &apiDefinition
	httpRouteState := synchronizer.HTTPRouteState{}

	methodTypeGet := gwapiv1.HTTPMethodGet
	methodTypePost := gwapiv1.HTTPMethodPost

	httpRoute := gwapiv1.HTTPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
			Name:      "test-api-2-prod-http-route",
		},
		Spec: gwapiv1.HTTPRouteSpec{
			Hostnames:       []gwapiv1.Hostname{"prod.gw.wso2.com"},
			CommonRouteSpec: createDefaultCommonRouteSpec(),
			Rules: []gwapiv1.HTTPRouteRule{
				{
					Matches: []gwapiv1.HTTPRouteMatch{
						{
							Path: &gwapiv1.HTTPPathMatch{
								Type:  operatorutils.PathMatchTypePtr(gwapiv1.PathMatchExact),
								Value: operatorutils.StringPtr("/exact-path-api/2.0.0/(.*)/exact-path"),
							},
							Method: &methodTypeGet,
						},
					},
					BackendRefs: []gwapiv1.HTTPBackendRef{
						createDefaultBackendRef("backend-1"),
					},
				},
				{
					Matches: []gwapiv1.HTTPRouteMatch{
						{
							Path: &gwapiv1.HTTPPathMatch{
								Type:  operatorutils.PathMatchTypePtr(gwapiv1.PathMatchRegularExpression),
								Value: operatorutils.StringPtr("/regex-path/2.0.0/userId/([^/]+)/orderId/([^/]+)"),
							},
							Method: &methodTypePost,
						},
					},
					BackendRefs: []gwapiv1.HTTPBackendRef{
						createDefaultBackendRef("backend-2"),
					},
				},
			},
		},
	}
	hostName := gwapiv1.Hostname("prod.gw.wso2.com")
	gateway := gwapiv1.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
			Name:      "default-gateway",
		},
		Spec: gwapiv1.GatewaySpec{
			Listeners: []gwapiv1.Listener{
				{
					Name:     "httpslistener",
					Hostname: &hostName,
					Protocol: gwapiv1.HTTPSProtocolType,
				},
			},
		},
	}

	dataholder.UpdateGateway(gateway)
	xds.SanitizeGateway("default-gateway", true)
	httpRouteState.HTTPRouteCombined = &httpRoute

	backendMapping := make(map[string]*v1alpha2.ResolvedBackend)
	backendMapping[k8types.NamespacedName{Namespace: "default", Name: "backend-1"}.String()] =
		&v1alpha2.ResolvedBackend{Services: []v1alpha2.Service{{Host: "test-service-1.default", Port: 7001}}, Protocol: v1alpha2.HTTPProtocol}
	backendMapping[k8types.NamespacedName{Namespace: "default", Name: "backend-2"}.String()] =
		&v1alpha2.ResolvedBackend{Services: []v1alpha2.Service{{Host: "test-service-2.default", Port: 7002}}, Protocol: v1alpha2.HTTPProtocol}
	httpRouteState.BackendMapping = backendMapping

	apiState.ProdHTTPRoute = &httpRouteState
	apiState.AIProvider = new(v1alpha3.AIProvider)
	httpRouteState.RuleIdxToAiRatelimitPolicyMapping = make(map[int]*v1alpha3.AIRateLimitPolicy)
	adapterInternalAPI, labels, err := synchronizer.UpdateInternalMapsFromHTTPRoute(apiState, &httpRouteState, constants.Production)
	assert.Equal(t, map[string]struct{}{"default-gateway": {}}, labels, "Labels are incorrect.")
	assert.Nil(t, err, "Error should not be present when apiState is converted to a AdapterInternalAPI object")
	routes, clusters, _, _ := envoy.CreateRoutesWithClusters(adapterInternalAPI, nil, "prod.gw.wso2.com", "carbon.super")
	assert.Equal(t, 3, len(clusters), "Number of production clusters created is incorrect.")

	exactPathCluster := clusters[1]
	clusterName := strings.Split(exactPathCluster.GetName(), "_")

	assert.Equal(t, 5, len(clusterName), "clustername is incorrect. Expected: carbon.super__prod.gw.wso2.com_test-api-22.0.0_<id>, Found: %s", exactPathCluster.GetName())
	assert.Equal(t, clusterName[0], "carbon.super", "Path Level cluster name should contain org carbon.super, but found : %s", clusterName[0])
	assert.Equal(t, clusterName[2], "prod.gw.wso2.com", "Path Level cluster name should contain vhost prod.gw.wso2.com, but found : %s", clusterName[2])
	assert.Equal(t, clusterName[3], "test-api-22.0.0", "Path Level cluster name should contain api ptest-api-22.0.0, but found :  %s", clusterName[3])

	exactPathClusterHost := exactPathCluster.GetLoadAssignment().GetEndpoints()[0].GetLbEndpoints()[0].GetEndpoint().
		GetAddress().GetSocketAddress().GetAddress()
	exactPathClusterPort := exactPathCluster.GetLoadAssignment().GetEndpoints()[0].GetLbEndpoints()[0].GetEndpoint().
		GetAddress().GetSocketAddress().GetPortValue()
	exactPathClusterPriority := exactPathCluster.GetLoadAssignment().GetEndpoints()[0].Priority

	assert.NotEmpty(t, exactPathClusterHost, "Exact path cluster's assigned host should not be null")
	assert.Equal(t, "test-service-1.default", exactPathClusterHost, "Exact path cluster's assigned host is incorrect.")
	assert.NotEmpty(t, exactPathClusterPort, "Exact path cluster's assigned port should not be null")
	assert.Equal(t, uint32(7001), exactPathClusterPort, "Exact path cluster's assigned port is incorrect.")
	assert.Equal(t, uint32(0), exactPathClusterPriority, "Exact path cluster's assigned Priority is incorrect.")

	regexPathCluster := clusters[2]

	regexPathClusterHost := regexPathCluster.GetLoadAssignment().GetEndpoints()[0].GetLbEndpoints()[0].GetEndpoint().
		GetAddress().GetSocketAddress().GetAddress()
	regexPathClusterPort := regexPathCluster.GetLoadAssignment().GetEndpoints()[0].GetLbEndpoints()[0].GetEndpoint().
		GetAddress().GetSocketAddress().GetPortValue()
	regexPathClusterPriority := regexPathCluster.GetLoadAssignment().GetEndpoints()[0].Priority

	assert.NotEmpty(t, regexPathClusterHost, "Regex path cluster's assigned host should not be null")
	assert.Equal(t, "test-service-2.default", regexPathClusterHost, "Regex path cluster's assigned host is incorrect.")
	assert.NotEmpty(t, regexPathClusterPort, "Regex path cluster's assigned port should not be null")
	assert.Equal(t, uint32(7002), regexPathClusterPort, "Regex path cluster's assigned host is incorrect.")
	assert.Equal(t, uint32(0), regexPathClusterPriority, "Regex path cluster's assigned priority is incorrect.")

	assert.Equal(t, 5, len(routes), "Created number of routes are incorrect.")
	assert.Contains(t, []string{"^/test-api/2\\.0\\.0/exact-path-api/2\\.0\\.0/\\(\\.\\*\\)/exact-path([/]{0,1})"}, routes[2].GetMatch().GetSafeRegex().Regex)
	assert.Contains(t, []string{"^/test-api/2\\.0\\.0/regex-path/2.0.0/userId/([^/]+)/orderId/([^/]+)([/]{0,1})"}, routes[3].GetMatch().GetSafeRegex().Regex)
	assert.NotEqual(t, routes[2].GetMatch().GetSafeRegex().Regex, routes[3].GetMatch().GetSafeRegex().Regex,
		"The route regex for the two paths should not be the same")
	for _, route := range routes {
		loggers.LoggerAPKOperator.Infof("routes ==" + route.GetMatch().GetSafeRegex().Regex)
	}
}

func TestExtractAPIDetailsFromHTTPRouteForDefaultCase(t *testing.T) {

	apiState := generateSampleAPI("test-api-1", "1.0.0", "/test-api/1.0.0")
	httpRouteState := synchronizer.HTTPRouteState{}
	httpRouteState = *apiState.ProdHTTPRoute
	apiState.AIProvider = new(v1alpha3.AIProvider)
	httpRouteState.RuleIdxToAiRatelimitPolicyMapping = make(map[int]*v1alpha3.AIRateLimitPolicy)
	xds.SanitizeGateway("default-gateway", true)
	adapterInternalAPI, labels, err := synchronizer.UpdateInternalMapsFromHTTPRoute(apiState, &httpRouteState, constants.Production)
	assert.Equal(t, map[string]struct{}{"default-gateway": {}}, labels, "Labels are incorrect.")
	assert.Nil(t, err, "Error should not be present when apiState is converted to a AdapterInternalAPI object")
	assert.Equal(t, "Default", adapterInternalAPI.GetEnvironment(), "Environment is incorrect.")
}

func TestExtractAPIDetailsFromHTTPRouteForSpecificEnvironment(t *testing.T) {

	apiState := generateSampleAPI("test-api-2", "1.0.0", "/test-api2/1.0.0")
	httpRouteState := synchronizer.HTTPRouteState{}
	httpRouteState = *apiState.ProdHTTPRoute
	apiState.APIDefinition.Spec.Environment = "dev"
	xds.SanitizeGateway("default-gateway", true)

	apiState.AIProvider = new(v1alpha3.AIProvider)
	httpRouteState.RuleIdxToAiRatelimitPolicyMapping = make(map[int]*v1alpha3.AIRateLimitPolicy)
	adapterInternalAPI, labels, err := synchronizer.UpdateInternalMapsFromHTTPRoute(apiState, &httpRouteState, constants.Production)
	assert.Equal(t, map[string]struct{}{"default-gateway": {}}, labels, "Labels are incorrect.")
	assert.Nil(t, err, "Error should not be present when apiState is converted to a AdapterInternalAPI object")
	assert.Equal(t, "dev", adapterInternalAPI.GetEnvironment(), "Environment is incorrect.")
}

func generateSampleAPI(apiName string, apiVersion string, basePath string) synchronizer.APIState {

	apiState := synchronizer.APIState{}
	apiDefinition := v1alpha3.API{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
			Name:      apiName,
		},
		Spec: v1alpha3.APISpec{
			APIName:    apiName,
			APIVersion: apiVersion,
			BasePath:   basePath,
			Production: []v1alpha3.EnvConfig{
				{
					RouteRefs: []string{
						apiName + "-prod-http-route",
					},
				},
			},
		},
	}
	apiState.APIDefinition = &apiDefinition
	httpRouteState := synchronizer.HTTPRouteState{}
	methodTypeGet := gwapiv1.HTTPMethodGet

	httpRoute := gwapiv1.HTTPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
			Name:      apiName + "-prod-http-route",
		},
		Spec: gwapiv1.HTTPRouteSpec{
			Hostnames:       []gwapiv1.Hostname{"prod.gw.wso2.com"},
			CommonRouteSpec: createDefaultCommonRouteSpec(),
			Rules: []gwapiv1.HTTPRouteRule{
				{
					Matches: []gwapiv1.HTTPRouteMatch{
						{
							Path: &gwapiv1.HTTPPathMatch{
								Type:  operatorutils.PathMatchTypePtr(gwapiv1.PathMatchExact),
								Value: operatorutils.StringPtr("/exact-path-api/2.0.0/(.*)/exact-path"),
							},
							Method: &methodTypeGet,
						},
					},
					BackendRefs: []gwapiv1.HTTPBackendRef{
						createDefaultBackendRef(apiName + "backend-1"),
					},
				},
			},
		},
	}

	httpRouteState.HTTPRouteCombined = &httpRoute

	backendMapping := make(map[string]*v1alpha2.ResolvedBackend)
	backendMapping[k8types.NamespacedName{Namespace: "default", Name: apiName + "backend-1"}.String()] =
		&v1alpha2.ResolvedBackend{Services: []v1alpha2.Service{{Host: "test-service-1.default", Port: 7001}}, Protocol: v1alpha2.HTTPProtocol}
	httpRouteState.BackendMapping = backendMapping

	apiState.ProdHTTPRoute = &httpRouteState

	apiState.AIProvider = new(v1alpha3.AIProvider)
	httpRouteState.RuleIdxToAiRatelimitPolicyMapping = make(map[int]*v1alpha3.AIRateLimitPolicy)
	return apiState
}

// TODO: Fix this test case
func TestCreateRoutesWithClustersWithMultiplePathPrefixRules(t *testing.T) {
	apiState := synchronizer.APIState{}
	apiDefinition := v1alpha3.API{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
			Name:      "test-api-1",
		},
		Spec: v1alpha3.APISpec{
			APIName:    "test-api",
			APIVersion: "1.0.0",
			BasePath:   "/test-api/1.0.0",
			Production: []v1alpha3.EnvConfig{
				{
					RouteRefs: []string{
						"test-api-1-prod-http-route",
					},
				},
			},
		},
	}
	apiState.APIDefinition = &apiDefinition
	httpRouteState := synchronizer.HTTPRouteState{}

	apiState.AIProvider = new(v1alpha3.AIProvider)
	httpRouteState.RuleIdxToAiRatelimitPolicyMapping = make(map[int]*v1alpha3.AIRateLimitPolicy)
	httpRoute := gwapiv1.HTTPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
			Name:      "test-api-1-prod-http-route",
		},
		Spec: gwapiv1.HTTPRouteSpec{
			Hostnames:       []gwapiv1.Hostname{"prod.gw.wso2.com"},
			CommonRouteSpec: createDefaultCommonRouteSpec(),
			Rules: []gwapiv1.HTTPRouteRule{
				{
					Matches: []gwapiv1.HTTPRouteMatch{
						{
							Path: &gwapiv1.HTTPPathMatch{
								Type:  operatorutils.PathMatchTypePtr(gwapiv1.PathMatchPathPrefix),
								Value: operatorutils.StringPtr("/orders"),
							},
						},
					},
					BackendRefs: []gwapiv1.HTTPBackendRef{
						createDefaultBackendRef("order-backend"),
					},
				},
				{
					Matches: []gwapiv1.HTTPRouteMatch{
						{
							Path: &gwapiv1.HTTPPathMatch{
								Type:  operatorutils.PathMatchTypePtr(gwapiv1.PathMatchPathPrefix),
								Value: operatorutils.StringPtr("/users"),
							},
						},
					},
					BackendRefs: []gwapiv1.HTTPBackendRef{
						createDefaultBackendRef("user-backend"),
					},
				},
			},
		},
	}

	httpRouteState.HTTPRouteCombined = &httpRoute

	backendMapping := make(map[string]*v1alpha2.ResolvedBackend)
	backendMapping[k8types.NamespacedName{Namespace: "default", Name: "order-backend"}.String()] =
		&v1alpha2.ResolvedBackend{Services: []v1alpha2.Service{
			{Host: "order-service.default", Port: 80},
			{Host: "order-service-2.default", Port: 8080}},
			Protocol: v1alpha2.HTTPProtocol}
	backendMapping[k8types.NamespacedName{Namespace: "default", Name: "user-backend"}.String()] =
		&v1alpha2.ResolvedBackend{Services: []v1alpha2.Service{
			{Host: "user-service.default", Port: 8081},
			{Host: "user-service-2.default", Port: 8081}},
			Protocol: v1alpha2.HTTPProtocol}
	httpRouteState.BackendMapping = backendMapping

	apiState.ProdHTTPRoute = &httpRouteState
	xds.SanitizeGateway("default-gateway", true)

	adapterInternalAPI, labels, err := synchronizer.UpdateInternalMapsFromHTTPRoute(apiState, &httpRouteState, constants.Production)
	assert.Equal(t, map[string]struct{}{"default-gateway": {}}, labels, "Labels are incorrect.")
	assert.Nil(t, err, "Error should not be present when apiState is converted to a AdapterInternalAPI object")
	routes, clusters, _, _ := envoy.CreateRoutesWithClusters(adapterInternalAPI, nil, "prod.gw.wso2.com", "carbon.super")
	assert.Equal(t, 3, len(clusters), "Number of production clusters created is incorrect.")

	orderServiceCluster := clusters[1]
	clusterName := strings.Split(orderServiceCluster.GetName(), "_")

	assert.Equal(t, 5, len(clusterName), "clustername is incorrect. Expected: carbon.super__prod.gw.wso2.com_test-api1.0.0_<id>, Found: %s", orderServiceCluster.GetName())
	assert.Equal(t, clusterName[0], "carbon.super", "Path Level cluster name should contain org carbon.super, but found : %s", clusterName[0])
	assert.Equal(t, clusterName[2], "prod.gw.wso2.com", "Path Level cluster name should contain vhost prod.gw.wso2.com, but found : %s", clusterName[2])
	assert.Equal(t, clusterName[3], "test-api1.0.0", "Path Level cluster name should contain api test-api1.0.0, but found :  %s", clusterName[3])

	orderServiceClusterHost0 := orderServiceCluster.GetLoadAssignment().GetEndpoints()[0].GetLbEndpoints()[0].GetEndpoint().
		GetAddress().GetSocketAddress().GetAddress()
	orderServiceClusterPort0 := orderServiceCluster.GetLoadAssignment().GetEndpoints()[0].GetLbEndpoints()[0].GetEndpoint().
		GetAddress().GetSocketAddress().GetPortValue()
	orderServiceClusterPriority0 := orderServiceCluster.GetLoadAssignment().GetEndpoints()[0].Priority
	orderServiceClusterHost1 := orderServiceCluster.GetLoadAssignment().GetEndpoints()[1].GetLbEndpoints()[0].GetEndpoint().
		GetAddress().GetSocketAddress().GetAddress()
	orderServiceClusterPort1 := orderServiceCluster.GetLoadAssignment().GetEndpoints()[1].GetLbEndpoints()[0].GetEndpoint().
		GetAddress().GetSocketAddress().GetPortValue()
	orderServiceClusterPriority1 := orderServiceCluster.GetLoadAssignment().GetEndpoints()[1].Priority

	assert.NotEmpty(t, orderServiceClusterHost0, "Order Service Cluster's assigned host should not be null")
	assert.Equal(t, "order-service.default", orderServiceClusterHost0, "Order Service Cluster's assigned host is incorrect.")
	assert.NotEmpty(t, orderServiceClusterPort0, "Order Service Cluster's assigned port should not be null")
	assert.Equal(t, uint32(80), orderServiceClusterPort0, "Order Service Cluster's assigned port is incorrect.")
	assert.Equal(t, uint32(0), orderServiceClusterPriority0, "Order Service Cluster's assigned Priority is incorrect.")

	assert.NotEmpty(t, orderServiceClusterHost1, "Order Service Cluster's second endpoint host should not be null")
	assert.Equal(t, "order-service-2.default", orderServiceClusterHost1, "Order Service Cluster's second endpoint host is incorrect.")
	assert.NotEmpty(t, orderServiceClusterPort1, "Order Service Cluster's second endpoint port should not be null")
	assert.Equal(t, uint32(8080), orderServiceClusterPort1, "Order Service Cluster's second endpoint port is incorrect.")
	assert.Equal(t, uint32(0), orderServiceClusterPriority1, "Order Service Cluster's second endpoint Priority is incorrect.")

	userServiceCluster := clusters[2]

	userServiceClusterHost0 := userServiceCluster.GetLoadAssignment().GetEndpoints()[0].GetLbEndpoints()[0].GetEndpoint().
		GetAddress().GetSocketAddress().GetAddress()
	userServiceClusterPort0 := userServiceCluster.GetLoadAssignment().GetEndpoints()[0].GetLbEndpoints()[0].GetEndpoint().
		GetAddress().GetSocketAddress().GetPortValue()
	userServiceClusterPriority0 := userServiceCluster.GetLoadAssignment().GetEndpoints()[0].Priority
	userServiceClusterHost1 := userServiceCluster.GetLoadAssignment().GetEndpoints()[1].GetLbEndpoints()[0].GetEndpoint().
		GetAddress().GetSocketAddress().GetAddress()
	userServiceClusterPort1 := userServiceCluster.GetLoadAssignment().GetEndpoints()[1].GetLbEndpoints()[0].GetEndpoint().
		GetAddress().GetSocketAddress().GetPortValue()
	userServiceClusterPriority1 := userServiceCluster.GetLoadAssignment().GetEndpoints()[1].Priority

	assert.NotEmpty(t, userServiceClusterHost0, "User Service Cluster's assigned host should not be null")
	assert.Equal(t, "user-service.default", userServiceClusterHost0, "User Service Cluster's assigned host is incorrect.")
	assert.NotEmpty(t, userServiceClusterPort0, "User Service Cluster's assigned port should not be null")
	assert.Equal(t, uint32(8081), userServiceClusterPort0, "User Service Cluster's assigned host is incorrect.")
	assert.Equal(t, uint32(0), userServiceClusterPriority0, "User Service Cluster's assigned priority is incorrect.")

	assert.NotEmpty(t, userServiceClusterHost1, "User Service Cluster's second endpoint host should not be null")
	assert.Equal(t, "user-service-2.default", userServiceClusterHost1, "User Service Cluster's second endpoint host is incorrect.")
	assert.NotEmpty(t, userServiceClusterPort1, "User Service Cluster's second endpoint port should not be null")
	assert.Equal(t, uint32(8081), userServiceClusterPort1, "User Service Cluster's second endpoint port is incorrect.")
	assert.Equal(t, uint32(0), userServiceClusterPriority1, "API Level Cluster's second endpoint Priority is incorrect.")

	assert.Equal(t, 15, len(routes), "Created number of routes are incorrect.")
	assert.Contains(t, []string{"^/test-api/1\\.0\\.0/orders((?:/.*)*)"}, routes[2].GetMatch().GetSafeRegex().Regex)
	assert.Contains(t, []string{"^/test-api/1\\.0\\.0/users((?:/.*)*)"}, routes[9].GetMatch().GetSafeRegex().Regex)
	assert.NotEqual(t, routes[1].GetMatch().GetSafeRegex().Regex, routes[8].GetMatch().GetSafeRegex().Regex,
		"The route regex for the two paths should not be the same")
}

func TestCreateRoutesWithClustersWithBackendTLSConfigs(t *testing.T) {
	apiState := synchronizer.APIState{}
	apiDefinition := v1alpha3.API{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
			Name:      "test-api-3",
		},
		Spec: v1alpha3.APISpec{
			APIName:    "test-api-3",
			APIVersion: "1.0.0",
			BasePath:   "/test-api-3/1.0.0",
			Production: []v1alpha3.EnvConfig{
				{
					RouteRefs: []string{
						"test-api-3-prod-http-route",
					},
				},
			},
		},
	}
	apiState.APIDefinition = &apiDefinition
	httpRouteState := synchronizer.HTTPRouteState{}
	methodTypeGet := gwapiv1.HTTPMethodGet

	apiState.AIProvider = new(v1alpha3.AIProvider)
	httpRouteState.RuleIdxToAiRatelimitPolicyMapping = make(map[int]*v1alpha3.AIRateLimitPolicy)
	httpRoute := gwapiv1.HTTPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
			Name:      "test-api-3-prod-http-route",
		},
		Spec: gwapiv1.HTTPRouteSpec{
			Hostnames:       []gwapiv1.Hostname{"prod.gw.wso2.com"},
			CommonRouteSpec: createDefaultCommonRouteSpec(),
			Rules: []gwapiv1.HTTPRouteRule{
				{
					Matches: []gwapiv1.HTTPRouteMatch{
						{
							Path: &gwapiv1.HTTPPathMatch{
								Type:  operatorutils.PathMatchTypePtr(gwapiv1.PathMatchExact),
								Value: operatorutils.StringPtr("/resource-path"),
							},
							Method: &methodTypeGet,
						},
					},
					BackendRefs: []gwapiv1.HTTPBackendRef{
						createDefaultBackendRef("test-backend-3"),
					},
				},
			},
		},
	}

	httpRouteState.HTTPRouteCombined = &httpRoute

	backendMapping := make(map[string]*v1alpha2.ResolvedBackend)
	backendMapping[k8types.NamespacedName{Namespace: "default", Name: "test-backend-3"}.String()] =
		&v1alpha2.ResolvedBackend{Services: []v1alpha2.Service{{Host: "webhook.site", Port: 443}},
			Protocol: v1alpha2.HTTPSProtocol,
			TLS: v1alpha2.ResolvedTLSConfig{
				ResolvedCertificate: `-----BEGIN CERTIFICATE-----test-cert-data-----END CERTIFICATE-----`,
			}}
	httpRouteState.BackendMapping = backendMapping

	apiState.ProdHTTPRoute = &httpRouteState
	xds.SanitizeGateway("default-gateway", true)

	adapterInternalAPI, labels, err := synchronizer.UpdateInternalMapsFromHTTPRoute(apiState, &httpRouteState, constants.Production)
	assert.Equal(t, map[string]struct{}{"default-gateway": {}}, labels, "Labels are incorrect.")
	assert.Nil(t, err, "Error should not be present when apiState is converted to a AdapterInternalAPI object")
	_, clusters, _, _ := envoy.CreateRoutesWithClusters(adapterInternalAPI, nil, "prod.gw.wso2.com", "carbon.super")
	assert.Equal(t, 2, len(clusters), "Number of production clusters created is incorrect.")

	exactPathCluster := clusters[1]

	assert.True(t, strings.HasPrefix(exactPathCluster.GetName(), "carbon.super__prod.gw.wso2.com_test-api-31.0.0_"),
		"Exact path cluster name mismatch, %v", exactPathCluster.GetName())

	exactPathClusterHost := exactPathCluster.GetLoadAssignment().GetEndpoints()[0].GetLbEndpoints()[0].GetEndpoint().
		GetAddress().GetSocketAddress().GetAddress()
	exactPathClusterPort := exactPathCluster.GetLoadAssignment().GetEndpoints()[0].GetLbEndpoints()[0].GetEndpoint().
		GetAddress().GetSocketAddress().GetPortValue()
	exactPathClusterPriority := exactPathCluster.GetLoadAssignment().GetEndpoints()[0].Priority

	assert.NotEmpty(t, exactPathClusterHost, "Exact path cluster's assigned host should not be null")
	assert.Equal(t, "webhook.site", exactPathClusterHost, "Exact path cluster's assigned host is incorrect.")
	assert.NotEmpty(t, exactPathClusterPort, "Exact path cluster's assigned port should not be null")
	assert.Equal(t, uint32(443), exactPathClusterPort, "Exact path cluster's assigned port is incorrect.")
	assert.Equal(t, uint32(0), exactPathClusterPriority, "Exact path cluster's assigned Priority is incorrect.")
}

func createDefaultCommonRouteSpec() gwapiv1.CommonRouteSpec {
	return gwapiv1.CommonRouteSpec{
		ParentRefs: []gwapiv1.ParentReference{
			{
				Group:       operatorutils.GroupPtr("gateway.networking.k8s.io"),
				Kind:        operatorutils.KindPtr("Gateway"),
				Name:        gwapiv1.ObjectName("default-gateway"),
				SectionName: (*gwapiv1.SectionName)(operatorutils.StringPtr("httpslistener")),
			},
		},
	}
}

func createDefaultBackendRef(backendName string) gwapiv1.HTTPBackendRef {
	return gwapiv1.HTTPBackendRef{
		BackendRef: gwapiv1.BackendRef{
			BackendObjectReference: gwapiv1.BackendObjectReference{
				Group: (*gwapiv1.Group)(&v1alpha1.GroupVersion.Group),
				Kind:  operatorutils.KindPtr("Backend"),
				Name:  gwapiv1.ObjectName(backendName),
			},
		},
	}
}

func TestCreateHealthEndpoint(t *testing.T) {
	route := envoy.CreateHealthEndpoint()
	assert.NotNil(t, route, "Health Endpoint Route should not be null.")
	assert.Equal(t, "/health", route.Name, "Health Route Name is incorrect.")
	assert.Equal(t, "/health", route.GetMatch().GetPath(), "Health route path is incorrect.")
	assert.Equal(t, "{\"status\": \"healthy\"}", route.GetDirectResponse().GetBody().GetInlineString(), "Health response message is incorrect.")
	assert.Equal(t, uint32(200), route.GetDirectResponse().GetStatus(), "Health response status is incorrect.")
}

func TestCreateRoutesWithClustersDifferentBackendRefs(t *testing.T) {
	apiState := synchronizer.APIState{}
	apiDefinition := v1alpha3.API{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
			Name:      "test-api-different-backendrefs",
		},
		Spec: v1alpha3.APISpec{
			APIName:    "test-api-different-backendrefs",
			APIVersion: "1.0.0",
			BasePath:   "/test-api-different-backendrefs/1.0.0",
			Production: []v1alpha3.EnvConfig{
				{
					RouteRefs: []string{
						"test-api-different-backendrefs-prod-http-route",
					},
				},
			},
		},
	}
	apiState.APIDefinition = &apiDefinition
	httpRouteState := synchronizer.HTTPRouteState{}
	methodTypeGet := gwapiv1.HTTPMethodGet

	apiState.AIProvider = new(v1alpha3.AIProvider)
	httpRouteState.RuleIdxToAiRatelimitPolicyMapping = make(map[int]*v1alpha3.AIRateLimitPolicy)
	httpRoute := gwapiv1.HTTPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
			Name:      "test-api-different-backendrefs-prod-http-route",
		},
		Spec: gwapiv1.HTTPRouteSpec{
			Hostnames:       []gwapiv1.Hostname{"prod.gw.wso2.com"},
			CommonRouteSpec: createDefaultCommonRouteSpec(),
			Rules: []gwapiv1.HTTPRouteRule{
				{
					Matches: []gwapiv1.HTTPRouteMatch{
						{
							Path: &gwapiv1.HTTPPathMatch{
								Type:  operatorutils.PathMatchTypePtr(gwapiv1.PathMatchExact),
								Value: operatorutils.StringPtr("/resource-path-1"),
							},
							Method: &methodTypeGet,
						},
					},
					BackendRefs: []gwapiv1.HTTPBackendRef{
						createDefaultBackendRef("test-backend-1"),
					},
				},
				{
					Matches: []gwapiv1.HTTPRouteMatch{
						{
							Path: &gwapiv1.HTTPPathMatch{
								Type:  operatorutils.PathMatchTypePtr(gwapiv1.PathMatchExact),
								Value: operatorutils.StringPtr("/resource-path-2"),
							},
							Method: &methodTypeGet,
						},
					},
					BackendRefs: []gwapiv1.HTTPBackendRef{
						createDefaultBackendRef("test-backend-2"),
					},
				},
			},
		},
	}

	httpRouteState.HTTPRouteCombined = &httpRoute

	backendMapping := make(map[string]*v1alpha2.ResolvedBackend)
	backendMapping[k8types.NamespacedName{Namespace: "default", Name: "test-backend-1"}.String()] =
		&v1alpha2.ResolvedBackend{Services: []v1alpha2.Service{{Host: "webhook.site.1", Port: 443}},
			Protocol: v1alpha2.HTTPSProtocol,
			TLS: v1alpha2.ResolvedTLSConfig{
				ResolvedCertificate: `-----BEGIN CERTIFICATE-----test-cert-data-----END CERTIFICATE-----`,
			}}
	backendMapping[k8types.NamespacedName{Namespace: "default", Name: "test-backend-2"}.String()] =
		&v1alpha2.ResolvedBackend{Services: []v1alpha2.Service{{Host: "webhook.site.2", Port: 443}},
			Protocol: v1alpha2.HTTPSProtocol,
			TLS: v1alpha2.ResolvedTLSConfig{
				ResolvedCertificate: `-----BEGIN CERTIFICATE-----test-cert-data-----END CERTIFICATE-----`,
			}}
	httpRouteState.BackendMapping = backendMapping

	apiState.ProdHTTPRoute = &httpRouteState
	xds.SanitizeGateway("default-gateway", true)

	adapterInternalAPI, labels, err := synchronizer.UpdateInternalMapsFromHTTPRoute(apiState, &httpRouteState, constants.Production)
	assert.Equal(t, map[string]struct{}{"default-gateway": {}}, labels, "Labels are incorrect.")
	assert.Nil(t, err, "Error should not be present when apiState is converted to a AdapterInternalAPI object")
	_, clusters, _, _ := envoy.CreateRoutesWithClusters(adapterInternalAPI, nil, "prod.gw.wso2.com", "carbon.super")
	assert.Equal(t, 3, len(clusters), "Number of production clusters created is incorrect.")
}

func TestCreateRoutesWithClustersSameBackendRefs(t *testing.T) {
	apiState := synchronizer.APIState{}
	apiDefinition := v1alpha3.API{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
			Name:      "test-api-same-backendrefs",
		},
		Spec: v1alpha3.APISpec{
			APIName:    "test-api-same-backendrefs",
			APIVersion: "1.0.0",
			BasePath:   "/test-api-same-backendrefs/1.0.0",
			Production: []v1alpha3.EnvConfig{
				{
					RouteRefs: []string{
						"test-api-same-backendrefs-prod-http-route",
					},
				},
			},
		},
	}
	apiState.APIDefinition = &apiDefinition
	httpRouteState := synchronizer.HTTPRouteState{}
	methodTypeGet := gwapiv1.HTTPMethodGet

	apiState.AIProvider = new(v1alpha3.AIProvider)
	httpRouteState.RuleIdxToAiRatelimitPolicyMapping = make(map[int]*v1alpha3.AIRateLimitPolicy)
	httpRoute := gwapiv1.HTTPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
			Name:      "test-api-same-backendrefs-prod-http-route",
		},
		Spec: gwapiv1.HTTPRouteSpec{
			Hostnames:       []gwapiv1.Hostname{"prod.gw.wso2.com"},
			CommonRouteSpec: createDefaultCommonRouteSpec(),
			Rules: []gwapiv1.HTTPRouteRule{
				{
					Matches: []gwapiv1.HTTPRouteMatch{
						{
							Path: &gwapiv1.HTTPPathMatch{
								Type:  operatorutils.PathMatchTypePtr(gwapiv1.PathMatchExact),
								Value: operatorutils.StringPtr("/resource-path-1"),
							},
							Method: &methodTypeGet,
						},
					},
					BackendRefs: []gwapiv1.HTTPBackendRef{
						createDefaultBackendRef("test-backend-1"),
					},
				},
				{
					Matches: []gwapiv1.HTTPRouteMatch{
						{
							Path: &gwapiv1.HTTPPathMatch{
								Type:  operatorutils.PathMatchTypePtr(gwapiv1.PathMatchExact),
								Value: operatorutils.StringPtr("/resource-path-2"),
							},
							Method: &methodTypeGet,
						},
					},
					BackendRefs: []gwapiv1.HTTPBackendRef{
						createDefaultBackendRef("test-backend-1"),
					},
				},
			},
		},
	}

	httpRouteState.HTTPRouteCombined = &httpRoute

	backendMapping := make(map[string]*v1alpha2.ResolvedBackend)
	backendMapping[k8types.NamespacedName{Namespace: "default", Name: "test-backend-1"}.String()] =
		&v1alpha2.ResolvedBackend{Services: []v1alpha2.Service{{Host: "webhook.site", Port: 443}},
			Protocol: v1alpha2.HTTPSProtocol,
			TLS: v1alpha2.ResolvedTLSConfig{
				ResolvedCertificate: `-----BEGIN CERTIFICATE-----test-cert-data-----END CERTIFICATE-----`,
			}}
	httpRouteState.BackendMapping = backendMapping

	apiState.ProdHTTPRoute = &httpRouteState
	xds.SanitizeGateway("default-gateway", true)

	adapterInternalAPI, labels, err := synchronizer.UpdateInternalMapsFromHTTPRoute(apiState, &httpRouteState, constants.Production)
	assert.Equal(t, map[string]struct{}{"default-gateway": {}}, labels, "Labels are incorrect.")
	assert.Nil(t, err, "Error should not be present when apiState is converted to a AdapterInternalAPI object")
	_, clusters, _, _ := envoy.CreateRoutesWithClusters(adapterInternalAPI, nil, "prod.gw.wso2.com", "carbon.super")
	assert.Equal(t, 2, len(clusters), "Number of production clusters created is incorrect.")
}
