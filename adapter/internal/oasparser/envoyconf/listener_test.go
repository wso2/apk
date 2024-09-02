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

package envoyconf

import (
	"testing"

	routev3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	envoy_config_trace_v3 "github.com/envoyproxy/go-control-plane/envoy/config/trace/v3"
	cors_filter_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/cors/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	"github.com/stretchr/testify/assert"
	"github.com/wso2/apk/adapter/config"
	"github.com/wso2/apk/adapter/internal/oasparser/model"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func TestCreateListenerWithRds(t *testing.T) {
	// TODO: (Vajira) Add more test scenarios
	gateway := new(gwapiv1.Gateway)
	gateway.Name = "default"
	listenerObj := new(gwapiv1.Listener)
	listenerObj.Name = "httpslistener"
	var hostname gwapiv1.Hostname
	hostname = "0.0.0.0"
	listenerObj.Hostname = &hostname
	listenerObj.Port = 9095
	listenerObj.Protocol = "HTTPS"
	gateway.Spec.Listeners = append(gateway.Spec.Listeners, *listenerObj)
	resolvedListenerCerts := make(map[string]map[string][]byte)
	listeners := CreateListenerByGateway(gateway, resolvedListenerCerts, "")
	assert.NotEmpty(t, listeners, "Listeners creation has been failed")

	securedListener := listeners[0]
	if securedListener.Validate() != nil {
		t.Error("Listener validation failed")
	}
	assert.Equal(t, "0.0.0.0", securedListener.GetAddress().GetSocketAddress().GetAddress(),
		"Address mismatch for secured Listener.")
	assert.Equal(t, uint32(9095), securedListener.GetAddress().GetSocketAddress().GetPortValue(),
		"Address mismatch for secured Listener.")
	assert.NotEmpty(t, securedListener.FilterChains, "Filter chain for listener should not be null.")
	assert.NotNil(t, securedListener.FilterChains[0].GetTransportSocket(),
		"Transport Socket should not be null for secured listener")
}

func TestCreateVirtualHost(t *testing.T) {
	// TODO: (Vajira) Add more test scenarios

	vhostToRouteArrayMap := map[string][]*routev3.Route{
		"*":           testCreateRoutesForUnitTests(t),
		"mg.wso2.com": testCreateRoutesForUnitTests(t),
	}
	vHosts := CreateVirtualHosts(vhostToRouteArrayMap, nil, make(map[string]bool), make(map[string]bool))

	if len(vHosts) != 2 {
		t.Error("Virtual Host creation failed")
	}

	for _, vHost := range vHosts {
		_, found := vhostToRouteArrayMap[vHost.Name]
		if found {
			if vHost.Domains[0] != vHost.Name {
				t.Errorf("Virtual Host domain mismatched, expected %s but found %s",
					vHost.Name, vHost.Domains[0])
			}
		} else {
			t.Errorf("Invalid additional Virtual Host: %s", vHost.Name)
		}
	}
}

func TestCreateRoutesConfigForRds(t *testing.T) {
	// TODO: (Vajira) Add more test scenarios
	vhostToRouteArrayMap := map[string][]*routev3.Route{
		"*":           testCreateRoutesForUnitTests(t),
		"mg.wso2.com": testCreateRoutesForUnitTests(t),
	}
	httpListeners := "httpslistener"
	vHosts := CreateVirtualHosts(vhostToRouteArrayMap, nil, make(map[string]bool), make(map[string]bool))
	rConfig := CreateRoutesConfigForRds(vHosts, httpListeners)

	assert.NotNil(t, rConfig, "CreateRoutesConfigForRds is failed")
	if rConfig.Validate() != nil {
		t.Errorf("rConfig Validation failed")
	}
}

func TestGetTracingOTLPForSuccessPath(t *testing.T) {

	conf := config.ReadConfigs()
	conf.Tracing.Enabled = true
	conf.Tracing.Type = "otlp"
	conf.Tracing.ConfigProperties = map[string]string{
		"endpoint":               "localhost:55680",
		"maxPathLength":          "512",
		"host":                   "localhost",
		"port":                   "55680",
		"instrumentationName":    "CHOREO-CONNECT",
		"maximumTracesPerSecond": "2",
		"connectionTimeout":      "10",
	}
	config.SetDefaultConfig()
	config.SetConfig(conf)

	tracing, err := getTracingOTLP(conf)
	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, tracing, "Tracing should not be nil")

	assert.Equal(t, "envoy.tracers.opentelemetry", tracing.GetProvider().Name, "Name should be envoy.tracers.opentelemetry")
	assert.Equal(t, uint32(512), tracing.GetMaxPathTagLength().GetValue(),
		"MaxPathTagLength should be 512")

	opConf := &envoy_config_trace_v3.OpenTelemetryConfig{}
	err = tracing.GetProvider().GetTypedConfig().UnmarshalTo(opConf)
	assert.Nilf(t, err, "Error while parsing Open Telemetry Config %v", opConf)
	assert.Equal(t, tracerServiceNameRouter, opConf.GetServiceName(), "Service Name should be "+tracerServiceNameRouter)
	assert.Equal(t, int64(10), opConf.GetGrpcService().GetTimeout().GetSeconds(), "Timeout should be 10 seconds.")
	assert.Equal(t, tracingClusterName, opConf.GetGrpcService().GetEnvoyGrpc().GetClusterName(),
		"Tracing cluster name should be "+tracingClusterName)
}

func TestGetTracingOTLPForInvalidMaxPath(t *testing.T) {

	conf := config.ReadConfigs()
	conf.Tracing.Enabled = true
	conf.Tracing.Type = "otlp"
	conf.Tracing.ConfigProperties = map[string]string{
		"endpoint": "localhost:55680",
	}
	config.SetDefaultConfig()
	config.SetConfig(conf)

	tracing, err := getTracingOTLP(conf)
	assert.NotNil(t, err, "Error should not be nil")
	assert.EqualErrorf(t, err, "invalid max path length provided for tracing endpoint",
		"Error should be invalid max path length provided for tracing endpoint")
	assert.Nil(t, tracing, "Tracing should be nil")
}

func TestGetTracingOTLPForInvalidTimeout(t *testing.T) {

	conf := config.ReadConfigs()
	conf.Tracing.Enabled = true
	conf.Tracing.Type = "otlp"
	conf.Tracing.ConfigProperties = map[string]string{
		"maxPathLength":          "512",
		"maximumTracesPerSecond": "2",
		"connectionTimeout":      "10s",
	}
	config.SetDefaultConfig()
	config.SetConfig(conf)

	tracing, err := getTracingOTLP(conf)
	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, tracing, "Tracing should not be nil")

	assert.Equal(t, "envoy.tracers.opentelemetry", tracing.GetProvider().Name, "Name should be envoy.tracers.opentelemetry")
	assert.Equal(t, uint32(512), tracing.GetMaxPathTagLength().GetValue(),
		"MaxPathTagLength should be 512")

	opConf := &envoy_config_trace_v3.OpenTelemetryConfig{}
	err = tracing.GetProvider().GetTypedConfig().UnmarshalTo(opConf)
	assert.Nilf(t, err, "Error while parsing Open Telemetry Config %v", opConf)
	assert.Equal(t, int64(20), opConf.GetGrpcService().GetTimeout().GetSeconds(), "Timeout should be 20 seconds.")
}

// Create some routes to perform unit tests
func testCreateRoutesForUnitTests(t *testing.T) []*routev3.Route {
	//cors configuration
	corsConfigModel3 := &model.CorsConfig{
		Enabled:                   true,
		AccessControlAllowMethods: []string{"GET"},
		AccessControlAllowOrigins: []string{"http://test1.com", "http://test2.com"},
	}

	endpoint := model.Endpoint{
		Host:    "abc.com",
		URLType: "http",
		Port:    80,
		RawURL:  "http://abc.com",
	}

	operationGet := model.NewOperation("GET", nil, nil, "")
	operationPost := model.NewOperation("POST", nil, nil, "")
	operationPut := model.NewOperation("PUT", nil, nil, "")
	resourceWithGet := model.CreateMinimalDummyResourceForTests("/resourcePath", []*model.Operation{operationGet},
		"resource_operation_id", []model.Endpoint{endpoint}, false, false)
	resourceWithPost := model.CreateMinimalDummyResourceForTests("/resourcePath", []*model.Operation{operationPost},
		"resource_operation_id", []model.Endpoint{endpoint}, false, false)
	resourceWithPut := model.CreateMinimalDummyResourceForTests("/resourcePath", []*model.Operation{operationPut},
		"resource_operation_id", []model.Endpoint{endpoint}, false, false)
	resourceWithMultipleOperations := model.CreateMinimalDummyResourceForTests("/resourcePath", []*model.Operation{operationGet, operationPut},
		"resource_operation_id", []model.Endpoint{endpoint}, false, false)

	route1, err := createRoutes(generateRouteCreateParamsForUnitTests("test", "HTTP", "localhost", "/test", "1.0.0", "/test",
		&resourceWithGet, "test-cluster", corsConfigModel3, false))
	assert.Nil(t, err, "Error while creating routes for resourceWithGet")
	route2, err := createRoutes(generateRouteCreateParamsForUnitTests("test", "HTTP", "localhost", "/test", "1.0.0", "/test",
		&resourceWithPost, "test-cluster", corsConfigModel3, false))
	assert.Nil(t, err, "Error while creating routes for resourceWithPost")
	route3, err := createRoutes(generateRouteCreateParamsForUnitTests("test", "HTTP", "localhost", "/test", "1.0.0", "/test",
		&resourceWithPut, "test-cluster", corsConfigModel3, false))
	assert.Nil(t, err, "Error while creating routes for resourceWithPut")
	route4, err := createRoutes(generateRouteCreateParamsForUnitTests("test", "HTTP", "localhost", "/test", "1.0.0", "/test",
		&resourceWithMultipleOperations, "test-cluster", corsConfigModel3, false))
	assert.Nil(t, err, "Error while creating routes for resourceWithMultipleOperations")

	routes := []*routev3.Route{route1[0], route2[0], route3[0], route4[0]}

	// check cors after creating routes
	for _, r := range routes {
		corsConfig := &cors_filter_v3.CorsPolicy{}
		err := r.GetTypedPerFilterConfig()[wellknown.CORS].UnmarshalTo(corsConfig)
		assert.Nilf(t, err, "Error while parsing Cors Configuration %v", corsConfig)
		assert.NotEmpty(t, corsConfig.GetAllowMethods(), "Cors AllowMethods should not be empty.")
		assert.NotEmpty(t, corsConfig.GetAllowOriginStringMatch(), "Cors AllowOriginStringMatch should not be empty.")
	}

	return routes
}
