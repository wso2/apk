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
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
	"testing"
	"time"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	routev3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	extAuthService "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/ext_authz/v3"
	tlsv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
	envoy_type_matcherv3 "github.com/envoyproxy/go-control-plane/envoy/type/matcher/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/stretchr/testify/assert"
	"github.com/wso2/apk/adapter/config"
	"github.com/wso2/apk/adapter/internal/oasparser/model"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestGenerateRoutePaths(t *testing.T) {
	// Tested features
	// 1. Route regex when xWso2BasePath is provided
	xWso2BasePath := "/xWso2BasePath"
	resourcePath := "/resource"

	completeRoutePath := generateRoutePath(xWso2BasePath, resourcePath)
	// TODO: (VirajSalaka) check if it is possible to perform an equals operation instead of prefix
	if !strings.HasPrefix(completeRoutePath, "^/xWso2BasePath/resource") {
		t.Error("The generated path should contain xWso2BasePath as a prefix if xWso2Basepath is available.")
	}
}

func TestCreateRoute(t *testing.T) {
	// Tested features
	// 1. RouteAction (Substitution involved) when xWso2BasePath is provided
	// 2. RouteAction (No substitution) config when xWso2BasePath is empty
	// 3. If HostRewriteSpecifier is set to Auto rewrite
	// 4. Method header regex matcher
	vHost := "localhost"
	xWso2BasePath := "/xWso2BasePath"
	basePath := "/basepath"
	title := "WSO2"
	apiType := "HTTP"
	endpoint := model.Endpoint{
		Host:     "abc.com",
		Basepath: basePath,
		URLType:  "http",
		Port:     80,
		RawURL:   "http://abc.com",
	}
	version := "1.0"
	resourceWithGet := model.CreateMinimalDummyResourceForTests("/resourcePath", []*model.Operation{model.NewOperation("GET", nil, nil)},
		"resource_operation_id", []model.Endpoint{}, []model.Endpoint{})
	resourceWithGetPost := model.CreateMinimalDummyResourceForTests("/resourcePath", []*model.Operation{model.NewOperation("GET", nil, nil), model.NewOperation("POST", nil, nil)},
		"resource_operation_id", []model.Endpoint{}, []model.Endpoint{})
	clusterName := "resource_operation_id"
	hostRewriteSpecifier := &routev3.RouteAction_AutoHostRewrite{
		AutoHostRewrite: &wrapperspb.BoolValue{
			Value: true,
		},
	}
	clusterSpecifier := &routev3.RouteAction_ClusterHeader{
		ClusterHeader: clusterHeaderName,
	}
	regexRewriteWithXWso2BasePath := &envoy_type_matcherv3.RegexMatchAndSubstitute{
		Pattern: &envoy_type_matcherv3.RegexMatcher{
			EngineType: &envoy_type_matcherv3.RegexMatcher_GoogleRe2{
				GoogleRe2: &envoy_type_matcherv3.RegexMatcher_GoogleRE2{
					MaxProgramSize: nil,
				},
			},
			Regex: "^/xWso2BasePath/resourcePath[/]{0,1}",
		},
		Substitution: "/basepath/resourcePath",
	}

	UpgradeConfigsDisabled := []*routev3.RouteAction_UpgradeConfig{{
		UpgradeType: "websocket",
		Enabled:     &wrappers.BoolValue{Value: false},
	}}

	TimeOutConfig := ptypes.DurationProto(60 * time.Second)

	IdleTimeOutConfig := ptypes.DurationProto(300 * time.Second)

	expectedRouteActionWithXWso2BasePath := &routev3.Route_Route{
		Route: &routev3.RouteAction{
			HostRewriteSpecifier: hostRewriteSpecifier,
			RegexRewrite:         regexRewriteWithXWso2BasePath,
			ClusterSpecifier:     clusterSpecifier,
			UpgradeConfigs:       UpgradeConfigsDisabled,
			Timeout:              TimeOutConfig,
			IdleTimeout:          IdleTimeOutConfig,
		},
	}

	generatedRouteArrayWithXWso2BasePath, err := createRoutes(generateRouteCreateParamsForUnitTests(title, apiType, vHost, xWso2BasePath, version,
		endpoint.Basepath, &resourceWithGet, clusterName, "", nil, false))
	assert.Nil(t, err, "Error while creating routes WithXWso2BasePath")
	generatedRouteWithXWso2BasePath := generatedRouteArrayWithXWso2BasePath[0]
	assert.NotNil(t, generatedRouteWithXWso2BasePath, "Route should not be null.")
	assert.Equal(t, expectedRouteActionWithXWso2BasePath, generatedRouteWithXWso2BasePath.Action,
		"Route generation mismatch when xWso2BasePath option is provided.")
	assert.NotNil(t, generatedRouteWithXWso2BasePath.GetMatch().Headers, "Headers property should not be null")
	assert.Equal(t, "^GET|OPTIONS$", generatedRouteWithXWso2BasePath.GetMatch().Headers[0].GetStringMatch().GetSafeRegex().Regex,
		"Assigned HTTP Method Regex is incorrect when single method is available.")

	generatedRouteArrayWithoutXWso2BasePath, err := createRoutes(generateRouteCreateParamsForUnitTests(title, apiType, vHost, "", version,
		endpoint.Basepath, &resourceWithGetPost, clusterName, "", nil, false))
	assert.Nil(t, err, "Error while creating routes WithoutXWso2BasePath")
	generatedRouteWithoutXWso2BasePath := generatedRouteArrayWithoutXWso2BasePath[0]
	assert.NotNil(t, generatedRouteWithoutXWso2BasePath, "Route should not be null")
	assert.NotNil(t, generatedRouteWithoutXWso2BasePath.GetMatch().Headers, "Headers property should not be null")
	assert.Equal(t, "^GET|POST|OPTIONS$", generatedRouteWithoutXWso2BasePath.GetMatch().Headers[0].GetStringMatch().GetSafeRegex().Regex,
		"Assigned HTTP Method Regex is incorrect when multiple methods are available.")

	context := fmt.Sprintf("%s/%s", xWso2BasePath, version)
	generatedRouteWithDefaultVersion, err := createRoutes(generateRouteCreateParamsForUnitTests(title, apiType, vHost, context, version,
		endpoint.Basepath, &resourceWithGetPost, clusterName, "", nil, true))
	assert.Nil(t, err, "Error while creating routes WithDefaultVersion")
	assert.NotNil(t, generatedRouteWithDefaultVersion, "Route should not be null")
	assert.True(t, strings.HasPrefix(generatedRouteWithDefaultVersion[0].GetMatch().GetSafeRegex().Regex, fmt.Sprintf("^(?:%s|%s)", context, xWso2BasePath)),
		"Default version basepath is not generated correctly")
}

func TestCreateRouteClusterSpecifier(t *testing.T) {
	// Tested features
	// 1. If the cluster for route is defined correctly depending on prodution only, sandbox only or both.

	// In this test case, the extAuthz context variables are not tested
	prodClusterName := "prodCluster"
	sandClusterName := "sandCluster"

	vHost := "localhost"
	xWso2BasePath := "/xWso2BasePath"
	endpointBasePath := "/basepath"
	title := "WSO2"
	version := "1.0.0"
	apiType := "HTTP"

	resourceWithGet := model.CreateMinimalDummyResourceForTests("/resourcePath", []*model.Operation{model.NewOperation("GET", nil, nil)},
		"resource_operation_id", []model.Endpoint{}, []model.Endpoint{})

	routeWithProdEp, err := createRoutes(generateRouteCreateParamsForUnitTests(title, apiType, vHost, xWso2BasePath, version, endpointBasePath,
		&resourceWithGet, prodClusterName, "", nil, false))
	assert.Nil(t, err, "Error while creating routeWithProdEp")
	assert.NotNil(t, routeWithProdEp[0], "Route should not be null")
	assert.NotNil(t, routeWithProdEp[0].GetRoute().GetClusterHeader(), "Route Cluster Header should not be null.")
	assert.Empty(t, routeWithProdEp[0].GetRoute().GetCluster(), "Route Cluster Name should be empty.")
	assert.Equal(t, clusterHeaderName, routeWithProdEp[0].GetRoute().GetClusterHeader(), "Route Cluster Name mismatch.")

	routeWithSandEp, err := createRoutes(generateRouteCreateParamsForUnitTests(title, apiType, vHost, xWso2BasePath, version, endpointBasePath,
		&resourceWithGet, "", sandClusterName, nil, false))
	assert.Nil(t, err, "Error while creating routeWithSandEp")
	assert.NotNil(t, routeWithSandEp[0], "Route should not be null")
	assert.NotNil(t, routeWithSandEp[0].GetRoute().GetClusterHeader(), "Route Cluster Header should not be null.")
	assert.Empty(t, routeWithSandEp[0].GetRoute().GetCluster(), "Route Cluster Name should be empty.")
	assert.Equal(t, clusterHeaderName, routeWithSandEp[0].GetRoute().GetClusterHeader(), "Route Cluster Name mismatch.")

	routeWithProdSandEp, err := createRoutes(generateRouteCreateParamsForUnitTests(title, apiType, vHost, xWso2BasePath, version, endpointBasePath,
		&resourceWithGet, prodClusterName, sandClusterName, nil, false))
	assert.Nil(t, err, "Error while creating routeWithProdSandEp")
	assert.NotNil(t, routeWithProdSandEp[0], "Route should not be null")
	assert.NotNil(t, routeWithProdSandEp[0].GetRoute().GetClusterHeader(), "Route Cluster Header should not be null.")
	assert.Empty(t, routeWithProdSandEp[0].GetRoute().GetCluster(), "Route Cluster Name should be empty.")
	assert.Equal(t, clusterHeaderName, routeWithProdSandEp[0].GetRoute().GetClusterHeader(), "Route Cluster Name mismatch.")
}

func TestCreateRouteExtAuthzContext(t *testing.T) {
	// Tested features
	// 1. The context variables inside extAuthzPerRoute configuration including
	// (prod/sand clustername, method regex, basePath, resourcePath, title, version)
	prodClusterName := "prodCluster"
	sandClusterName := "sandCluster"

	vHost := "localhost"
	xWso2BasePath := "/xWso2BasePath"
	endpointBasePath := "/basepath"
	title := "WSO2"
	version := "1.0.0"
	apiType := "HTTP"

	resourceWithGet := model.CreateMinimalDummyResourceForTests("/resourcePath", []*model.Operation{model.NewOperation("GET", nil, nil)},
		"resource_operation_id", []model.Endpoint{}, []model.Endpoint{})

	routeWithProdEp, err := createRoutes(generateRouteCreateParamsForUnitTests(title, apiType, vHost, xWso2BasePath, version,
		endpointBasePath, &resourceWithGet, prodClusterName, sandClusterName, nil, false))
	assert.Nil(t, err, "Error while creating routeWithProdEp")
	assert.NotNil(t, routeWithProdEp[0], "Route should not be null")
	assert.NotNil(t, routeWithProdEp[0].GetTypedPerFilterConfig(), "TypedPerFilter config should not be null")
	assert.NotNil(t, routeWithProdEp[0].GetTypedPerFilterConfig()[wellknown.HTTPExternalAuthorization],
		"ExtAuthzPerRouteConfig should not be empty")

	extAuthPerRouteConfig := &extAuthService.ExtAuthzPerRoute{}
	err = ptypes.UnmarshalAny(routeWithProdEp[0].TypedPerFilterConfig[wellknown.HTTPExternalAuthorization],
		extAuthPerRouteConfig)
	assert.Nilf(t, err, "Error while parsing ExtAuthzPerRouteConfig %v", extAuthPerRouteConfig)
	assert.NotEmpty(t, extAuthPerRouteConfig.GetCheckSettings(), "Check Settings per ext authz route should not be empty")
	assert.NotEmpty(t, extAuthPerRouteConfig.GetCheckSettings().ContextExtensions,
		"ContextExtensions per ext authz route should not be empty")

	contextExtensionMap := extAuthPerRouteConfig.GetCheckSettings().ContextExtensions
	assert.Equal(t, title, contextExtensionMap[apiNameContextExtension], "Title mismatch in route ext authz context.")
	assert.Equal(t, xWso2BasePath, contextExtensionMap[basePathContextExtension], "Basepath mismatch in route ext authz context.")
	assert.Equal(t, version, contextExtensionMap[apiVersionContextExtension], "Version mismatch in route ext authz context.")
	assert.Equal(t, "GET", contextExtensionMap[methodContextExtension], "Method mismatch in route ext authz context.")
	assert.Equal(t, prodClusterName, contextExtensionMap[prodClusterNameContextExtension],
		"Production Cluster mismatch in route ext authz context.")
	assert.Equal(t, sandClusterName, contextExtensionMap[sandClusterNameContextExtension],
		"Sandbox Cluster mismatch in route ext authz context.")
}

func TestGenerateTLSCert(t *testing.T) {
	publicKeyPath := config.GetMgwHome() + "/adapter/security/localhost.pem"
	privateKeyPath := config.GetMgwHome() + "/adapter/security/localhost.key"

	tlsCert := generateTLSCert(privateKeyPath, publicKeyPath)

	assert.NotNil(t, tlsCert, "TLS Certificate should not be null")

	assert.NotNil(t, tlsCert.GetPrivateKey(), "Private Key should not be null in the TLS certificate")
	assert.NotNil(t, tlsCert.GetCertificateChain(), "Certificate chain should not be null in the TLS certificate")

	assert.Equal(t, tlsCert.GetPrivateKey().GetFilename(), privateKeyPath, "Private Key Value mismatch in the TLS Certificate")
	assert.Equal(t, tlsCert.GetCertificateChain().GetFilename(), publicKeyPath, "Certificate Chain Value mismatch in the TLS Certificate")
}

func TestGenerateRegex(t *testing.T) {

	type generateRegexTestItem struct {
		basePath      string
		resourcePath  string
		userInputPath string
		message       string
		isMatched     bool
	}
	dataItems := []generateRegexTestItem{
		{
			basePath:      "/v2",
			resourcePath:  "/pet/{petId}",
			userInputPath: "/v2/pet/5",
			message:       "when path parameter is provided end of the path",
			isMatched:     true,
		},
		{
			basePath:      `/v2`,
			resourcePath:  `/pet/{petId}`,
			userInputPath: "/v2/pet/5/",
			message:       "when the input path does not have tailing slash and user input path has trailing slash",
			isMatched:     true,
		},
		{
			basePath:      "/v2",
			resourcePath:  "/pet/{petId}/",
			userInputPath: "/v2/pet/5",
			message:       "when the input path has tailing slash and user input path does not have trailing slash",
			isMatched:     true,
		},
		{
			basePath:      "/v2",
			resourcePath:  "/pet/{petId}/",
			userInputPath: "/v2/pet/5/",
			message:       "when both the input path and user input path has trailing slash",
			isMatched:     true,
		},
		{
			basePath:      "/v2",
			resourcePath:  "/pet/{petId}/info",
			userInputPath: "/v2/pet/5/info",
			message:       "when path parameter is provided in the middle of the path",
			isMatched:     true,
		},
		{
			basePath:      "/v2",
			resourcePath:  "/pet/{petId}",
			userInputPath: "/v2/pet/5",
			message:       "when path parameter is provided end of the path",
			isMatched:     true,
		},
		{
			basePath:      "/v2",
			resourcePath:  "/pet/{petId}/tst/{petId}",
			userInputPath: "/v2/pet/5/tst/3",
			message:       "when multiple path parameter are provided",
			isMatched:     true,
		},
		{
			basePath:      "/v2",
			resourcePath:  "/pet/{petId}",
			userInputPath: "/v2/pet/5/test",
			message:       "when path parameter is provided end of the path and provide incorrect path",
			isMatched:     false,
		},
		{
			basePath:      "/v2",
			resourcePath:  "/pet/5",
			userInputPath: "/v2/pett/5",
			message:       "when provide a incorrect path",
			isMatched:     false,
		},
		{
			basePath:      "/v2",
			resourcePath:  "/pet/*",
			userInputPath: "/v2/pet/",
			message:       "when the resource ends with *, empty string with / substitution fails.",
			isMatched:     true,
		},
		{
			basePath:      "/v2",
			resourcePath:  "/pet/*",
			userInputPath: "/v2/pet",
			message:       "when the resource ends with *, empty string substitution fails.",
			isMatched:     true,
		},
		{
			basePath:      "/v2",
			resourcePath:  "/pet/*",
			userInputPath: "/v2/pet/foo/bar",
			message:       "when the resource ends with *, multiple trailing slashes substitution fails.",
			isMatched:     true,
		},
		{
			basePath:      "/v2",
			resourcePath:  "/pet/*",
			userInputPath: "/v2/pet123",
			message:       "when the resource ends with *, trailing characters substitution passes",
			isMatched:     false,
		},
		{
			basePath:      "/v2",
			resourcePath:  "/pet/{petId}.api",
			userInputPath: "/v2/pet/findByIdstatus=availabe",
			message:       "when the resource path param suffixed",
			isMatched:     false,
		},
		{
			basePath:      "/v2",
			resourcePath:  "/pet/{petId}.api",
			userInputPath: "/v2/pet/pet1.api",
			message:       "when the resource path param suffixed",
			isMatched:     true,
		},
		{
			basePath:      "/v2",
			resourcePath:  "/pet/pet{petId}",
			userInputPath: "/v2/pet/findByIdstatus=availabe",
			message:       "when the resource ends with *",
			isMatched:     false,
		},
		{
			basePath:      "/v2",
			resourcePath:  "/pet/pet{petId}",
			userInputPath: "/v2/pet/pet1",
			message:       "when the resource ends with *",
			isMatched:     true,
		},
	}

	for _, item := range dataItems {
		resultPattern := generateRoutePath(item.basePath, item.resourcePath)
		// regexp.MatchString also returns true for partial matches. Therefore, an additional $ is added
		// below to replicate the behavior of envoy proxy. As per the doc,
		// "The entire path (without the query string) must match the regex.
		// The rule will not match if only a subsequence of the :path header matches the regex."
		// https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/route/v3/route_components.proto#envoy-v3-api-field-config-route-v3-routematch-safe-regex
		resultIsMatching, err := regexp.MatchString(resultPattern+"$", item.userInputPath)
		assert.Equal(t, item.isMatched, resultIsMatching, resultPattern)
		assert.Nil(t, err)
	}
}

func TestGenerateSubstitutionString(t *testing.T) {
	type generateSubsStringTestItem struct {
		inputPath          string
		expectedSubsString string
		message            string
		shouldEqual        bool
	}
	dataItems := []generateSubsStringTestItem{
		{
			"/v2/pet",
			"/basepath/v2/pet",
			"when input path does not have a trailing slash",
			true,
		},
		{
			"/v2/pet/",
			"/basepath/v2/pet/",
			"when input path has a trailing slash",
			true,
		},
		{
			"/v2/pet/",
			"/basepath/v2/pet",
			"when input path has a trailing slash",
			false,
		},
		{
			"/v2/pet/{petId}",
			"/basepath/v2/pet/\\1",
			"when input path has a path param",
			true,
		},
		{
			"/v2/pet/{petId}/",
			"/basepath/v2/pet/\\1/",
			"when input path has a path param and trailing slash",
			true,
		},
		{
			"/v2/pet/{petId}/info",
			"/basepath/v2/pet/\\1/info",
			"when input path has a path param in the middle of the path",
			true,
		},
		{
			"/v2/pet/{petId}/test/{petId}",
			"/basepath/v2/pet/\\1/test/\\2",
			"when input path has a two path params",
			true,
		},
		{
			"/v2/*",
			"/basepath/v2\\1",
			"when input path has a wildcard at the end",
			true,
		},
		{
			"/v2/{petId}/*",
			"/basepath/v2/\\1\\2",
			"when input path has a path param and a wildcard at the end",
			true,
		},
	}
	for _, item := range dataItems {
		generatedSubstitutionString := generateSubstitutionString("/basepath", item.inputPath)
		if item.shouldEqual {
			assert.Equal(t, item.expectedSubsString, generatedSubstitutionString, item.message)
		} else {
			assert.NotEqual(t, item.expectedSubsString, generatedSubstitutionString, item.message)
		}
	}
}

func TestGenerateRegexSegment(t *testing.T) {

	type generateRegexSegmentTestItem struct {
		inputPath    string
		regexSegment string
		message      string
		shouldEqual  bool
	}
	dataItems := []generateRegexSegmentTestItem{
		{
			inputPath:    "/v2/pet/",
			regexSegment: "/v2/pet/",
			message:      "when the input path has a trailing slash",
			shouldEqual:  true,
		},
		{
			inputPath:    "/v2/pet",
			regexSegment: "/v2/pet",
			message:      "when the input path does not have a trailing slash",
			shouldEqual:  true,
		},
		{
			inputPath:    "/v2/pet/{petId}",
			regexSegment: "/v2/pet/([^/]+)",
			message:      "when the input path has a path param and does not have a trailing slash",
			shouldEqual:  true,
		},
		{
			inputPath:    "/v2/pet/{petId}/",
			regexSegment: "/v2/pet/([^/]+)/",
			message:      "when the input path has a path param and a trailing slash",
			shouldEqual:  true,
		},
		{
			inputPath:    "/v2/pet/{petId}/test/{petId}",
			regexSegment: "/v2/pet/([^/]+)/test/([^/]+)",
			message:      "when the input path has two path params",
			shouldEqual:  true,
		},
		{
			inputPath:    "/v2/pet/*",
			regexSegment: "/v2/pet/*",
			message:      "when the input path ends with *",
			shouldEqual:  true,
		},
	}

	for _, item := range dataItems {
		generatedPathRegexSegment := replacePathParamsWithCaptureGroups(item.inputPath)
		isEqual := generatedPathRegexSegment == item.regexSegment
		assert.Equal(t, item.shouldEqual, isEqual, item.message)
	}
}

func TestCreateUpstreamTLSContext(t *testing.T) {
	certFilePath := config.GetMgwHome() + "/../adapter/test-resources/envoycodegen/certs/testcrt.crt"
	certByteArr, err := ioutil.ReadFile(certFilePath)
	assert.Nil(t, err, "Error while reading the certificate : "+certFilePath)
	defaultMgwKeyPath := "/home/wso2/security/keystore/mg.key"
	defaultMgwCertPath := "/home/wso2/security/keystore/mg.pem"
	defaultCipherArray := "ECDHE-ECDSA-AES128-GCM-SHA256, ECDHE-RSA-AES128-GCM-SHA256, ECDHE-ECDSA-AES128-SHA," +
		" ECDHE-RSA-AES128-SHA, AES128-GCM-SHA256, AES128-SHA, ECDHE-ECDSA-AES256-GCM-SHA384, ECDHE-RSA-AES256-GCM-SHA384," +
		" ECDHE-ECDSA-AES256-SHA, ECDHE-RSA-AES256-SHA, AES256-GCM-SHA384, AES256-SHA"
	defaultCACertPath := "/etc/ssl/certs/ca-certificates.crt"
	hostNameAddress := &corev3.Address{Address: &corev3.Address_SocketAddress{
		SocketAddress: &corev3.SocketAddress{
			Address:  "abc.com",
			Protocol: corev3.SocketAddress_TCP,
			PortSpecifier: &corev3.SocketAddress_PortValue{
				PortValue: uint32(2384),
			},
		},
	}}

	hostNameAddressWithIP := &corev3.Address{Address: &corev3.Address_SocketAddress{
		SocketAddress: &corev3.SocketAddress{
			Address:  "10.10.10.10",
			Protocol: corev3.SocketAddress_TCP,
			PortSpecifier: &corev3.SocketAddress_PortValue{
				PortValue: uint32(2384),
			},
		},
	}}

	tlsCert := generateTLSCert(defaultMgwKeyPath, defaultMgwCertPath)
	upstreamTLSContextWithCerts := createUpstreamTLSContext(certByteArr, hostNameAddress, false)
	upstreamTLSContextWithoutCerts := createUpstreamTLSContext(nil, hostNameAddress, false)
	upstreamTLSContextWithIP := createUpstreamTLSContext(certByteArr, hostNameAddressWithIP, false)

	assert.NotEmpty(t, upstreamTLSContextWithCerts, "Upstream TLS Context should not be null when certs provided")
	assert.NotEmpty(t, upstreamTLSContextWithCerts.CommonTlsContext, "CommonTLSContext should not be "+
		"null when certs provided")
	assert.NotEmpty(t, upstreamTLSContextWithCerts.CommonTlsContext.TlsParams, "TlsParams should not be "+
		"null when certs provided")
	// tested against default TLS Parameters
	assert.Equal(t, tlsv3.TlsParameters_TLSv1_2, upstreamTLSContextWithCerts.CommonTlsContext.TlsParams.TlsMaximumProtocolVersion,
		"TLS maximum parameter mismatch")
	assert.Equal(t, tlsv3.TlsParameters_TLSv1_1, upstreamTLSContextWithCerts.CommonTlsContext.TlsParams.TlsMinimumProtocolVersion,
		"TLS minimum parameter mismatch")

	assert.Equal(t, defaultCipherArray, strings.Join(upstreamTLSContextWithCerts.CommonTlsContext.TlsParams.CipherSuites, ", "), "cipher suites mismatch")
	// the microgateway's certificate will be provided all the time. (For mutualSSL when required)
	assert.NotEmpty(t, upstreamTLSContextWithCerts.CommonTlsContext.TlsCertificates, "TLScerts should not be null")
	assert.Equal(t, tlsCert, upstreamTLSContextWithCerts.CommonTlsContext.TlsCertificates[0], "TLScert mismatch")
	assert.Equal(t, certByteArr, upstreamTLSContextWithCerts.CommonTlsContext.GetValidationContext().GetTrustedCa().GetInlineBytes(),
		"validation context certificate mismatch")
	assert.Equal(t, defaultCACertPath, upstreamTLSContextWithoutCerts.CommonTlsContext.GetValidationContext().GetTrustedCa().GetFilename(),
		"validation context certificate filepath mismatch")
	assert.NotEmpty(t, upstreamTLSContextWithCerts.CommonTlsContext.GetValidationContext().GetMatchTypedSubjectAltNames(),
		"Subject Alternative Names Should not be empty.")
	assert.Equal(t, "abc.com", upstreamTLSContextWithCerts.CommonTlsContext.GetValidationContext().GetMatchTypedSubjectAltNames()[0].GetMatcher().GetExact(),
		"Upstream SAN mismatch.")
	assert.Equal(t, tlsv3.SubjectAltNameMatcher_DNS, upstreamTLSContextWithCerts.CommonTlsContext.GetValidationContext().GetMatchTypedSubjectAltNames()[0].SanType,
		"Upstream SAN type mismatch.")

	assert.NotEmpty(t, upstreamTLSContextWithIP.CommonTlsContext.GetValidationContext().GetMatchTypedSubjectAltNames(),
		"Subject Alternative Names Should not be empty.")
	assert.Equal(t, "10.10.10.10", upstreamTLSContextWithIP.CommonTlsContext.GetValidationContext().GetMatchTypedSubjectAltNames()[0].GetMatcher().GetExact(),
		"Upstream SAN mismatch.")
	assert.Equal(t, tlsv3.SubjectAltNameMatcher_IP_ADDRESS, upstreamTLSContextWithIP.CommonTlsContext.GetValidationContext().GetMatchTypedSubjectAltNames()[0].SanType,
		"Upstream SAN type mismatch.")
}

func TestGetCorsPolicy(t *testing.T) {

	corsConfigModel1 := &model.CorsConfig{
		Enabled: false,
	}

	corsConfigModel2 := &model.CorsConfig{
		Enabled:                       true,
		AccessControlAllowMethods:     []string{"GET", "POST"},
		AccessControlAllowHeaders:     []string{"X-TEST-HEADER1", "X-TEST-HEADER2"},
		AccessControlExposeHeaders:    []string{"X-Custom-Header"},
		AccessControlAllowOrigins:     []string{"http://test.com"},
		AccessControlAllowCredentials: true,
	}

	corsConfigModel3 := &model.CorsConfig{
		Enabled:                   true,
		AccessControlAllowMethods: []string{"GET"},
		AccessControlAllowOrigins: []string{"http://test1.com", "http://test2.com"},
	}

	// Test the configuration when cors is disabled.
	corsPolicy1 := getCorsPolicy(corsConfigModel1)
	assert.Nil(t, corsPolicy1, "Cors Policy should be null.")

	// Test configuration when all the fields are provided.
	corsPolicy2 := getCorsPolicy(corsConfigModel2)
	assert.NotNil(t, corsPolicy2, "Cors Policy should not be null.")
	assert.NotEmpty(t, corsPolicy2.GetAllowOriginStringMatch(), "Cors Allowded Origins should not be null.")
	assert.Equal(t, regexp.QuoteMeta("http://test.com"),
		corsPolicy2.GetAllowOriginStringMatch()[0].GetSafeRegex().GetRegex(),
		"Cors Allowed Origin Header mismatch")
	assert.NotNil(t, corsPolicy2.GetAllowMethods())
	assert.Equal(t, "GET, POST", corsPolicy2.GetAllowMethods(), "Cors allow methods mismatch.")
	assert.NotNil(t, corsPolicy2.GetAllowHeaders(), "Cors Allowed headers should not be null.")
	assert.Equal(t, "X-TEST-HEADER1, X-TEST-HEADER2", corsPolicy2.GetAllowHeaders(), "Cors Allow headers mismatch")
	assert.NotNil(t, corsPolicy2.GetExposeHeaders(), "Cors Expose headers should not be null.")
	assert.Equal(t, "X-Custom-Header", corsPolicy2.GetExposeHeaders(), "Cors Expose headers mismatch")
	assert.True(t, corsPolicy2.GetAllowCredentials().GetValue(), "Cors Access Allow Credentials should be true")

	// Test the configuration when headers configuration is not provided.
	corsPolicy3 := getCorsPolicy(corsConfigModel3)
	assert.NotNil(t, corsPolicy3, "Cors Policy should not be null.")
	assert.Empty(t, corsPolicy3.GetAllowHeaders(), "Cors Allow headers should be null.")
	assert.Empty(t, corsPolicy3.GetExposeHeaders(), "Cors Expose Headers should be null.")
	assert.NotEmpty(t, corsPolicy3.GetAllowOriginStringMatch(), "Cors Allowded Origins should not be null.")
	assert.Equal(t, regexp.QuoteMeta("http://test1.com"),
		corsPolicy3.GetAllowOriginStringMatch()[0].GetSafeRegex().GetRegex(),
		"Cors Allowed Origin Header mismatch")
	assert.Equal(t, regexp.QuoteMeta("http://test2.com"),
		corsPolicy3.GetAllowOriginStringMatch()[1].GetSafeRegex().GetRegex(),
		"Cors Allowed Origin Header mismatch")
	assert.Empty(t, corsPolicy3.GetAllowCredentials(), "Allow Credential property should not be assigned.")

	resourceWithGet := model.CreateMinimalDummyResourceForTests("/resourcePath", []*model.Operation{model.NewOperation("GET", nil, nil)},
		"resource_operation_id", []model.Endpoint{}, []model.Endpoint{})

	// Route without CORS configuration
	routeWithoutCors, err := createRoutes(generateRouteCreateParamsForUnitTests("test", "HTTP", "localhost", "/test", "1.0.0", "/test",
		&resourceWithGet, "test-cluster", "", nil, false))
	assert.Nil(t, err, "Error while creating routeWithoutCors")
	assert.Nil(t, routeWithoutCors[0].GetRoute().Cors, "Cors Configuration should be null.")

	// Route with CORS configuration
	routeWithCors, err := createRoutes(generateRouteCreateParamsForUnitTests("test", "HTTP", "localhost", "/test", "1.0.0", "/test",
		&resourceWithGet, "test-cluster", "", corsConfigModel3, false))
	assert.Nil(t, err, "Error while creating routeWithCors")
	assert.NotNil(t, routeWithCors[0].GetRoute().Cors, "Cors Configuration should not be null.")
}

func generateRouteCreateParamsForUnitTests(title string, apiType string, vhost string, xWso2Basepath string, version string, endpointBasepath string,
	resource *model.Resource, prodClusterName string, sandClusterName string,
	corsConfig *model.CorsConfig, isDefaultVersion bool) *routeCreateParams {
	return &routeCreateParams{
		title:            title,
		apiType:          apiType,
		version:          version,
		vHost:            vhost,
		xWSO2BasePath:    xWso2Basepath,
		resource:         resource,
		prodClusterName:  prodClusterName,
		sandClusterName:  sandClusterName,
		endpointBasePath: endpointBasepath,
		corsPolicy:       corsConfig,
		isDefaultVersion: isDefaultVersion,
	}
}
