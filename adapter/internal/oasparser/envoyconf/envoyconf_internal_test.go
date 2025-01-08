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
	"io/ioutil"
	"regexp"
	"strings"
	"testing"
	"time"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	routev3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	cors_filter_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/cors/v3"
	// extAuthService "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/ext_authz/v3"
	tlsv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
	envoy_type_matcherv3 "github.com/envoyproxy/go-control-plane/envoy/type/matcher/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/stretchr/testify/assert"
	"github.com/wso2/apk/adapter/config"
	"github.com/wso2/apk/adapter/internal/oasparser/constants"
	"github.com/wso2/apk/adapter/internal/oasparser/model"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func TestCreateRoute(t *testing.T) {
	// Tested features
	// 1. RouteAction (Substitution involved) when xWso2BasePath is provided
	// 2. RouteAction (No substitution) config when xWso2BasePath is empty
	// 3. If HostRewriteSpecifier is set to Auto rewrite
	// 4. Method header regex matcher
	vHost := "localhost"
	xWso2BasePath := "/xWso2BasePath"
	title := "WSO2"
	apiType := "HTTP"
	endpoint := model.Endpoint{
		Host:    "abc.com",
		URLType: "http",
		Port:    80,
		RawURL:  "http://abc.com",
	}
	version := "1.0"

	// Creating path rewrite policy
	var policies = model.OperationPolicies{}
	policyParameters := make(map[string]interface{})
	policyParameters[constants.RewritePathType] = gwapiv1.PrefixMatchHTTPPathModifier
	policyParameters[constants.IncludeQueryParams] = true
	policyParameters[constants.RewritePathResourcePath] = "/basepath/resourcePath"
	policies.Request = append(policies.Request, model.Policy{
		PolicyName: string(gwapiv1.HTTPRouteFilterURLRewrite),
		Action:     constants.ActionRewritePath,
		Parameters: policyParameters,
	})

	resourceWithGet := model.CreateMinimalDummyResourceForTests("/xWso2BasePath/resourcePath",
		[]*model.Operation{model.NewOperationWithPolicies("GET", policies, "")},
		"resource_operation_id", []model.Endpoint{endpoint}, true, false)
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
			Regex: "^/xWso2BasePath/resourcePath((?:/.*)*)",
		},
		Substitution: "/basepath/resourcePath\\1",
	}

	UpgradeConfigsDisabled := []*routev3.RouteAction_UpgradeConfig{{
		UpgradeType: "websocket",
		Enabled:     &wrappers.BoolValue{Value: false},
	}}

	IdleTimeOutConfig := durationpb.New(time.Duration(300) * time.Second)

	expectedRouteActionWithXWso2BasePath := &routev3.Route_Route{
		Route: &routev3.RouteAction{
			HostRewriteSpecifier: hostRewriteSpecifier,
			RegexRewrite:         regexRewriteWithXWso2BasePath,
			ClusterSpecifier:     clusterSpecifier,
			UpgradeConfigs:       UpgradeConfigsDisabled,
			IdleTimeout:          IdleTimeOutConfig,
		},
	}

	resourceWithGet.GetEndpoints().Config = &model.EndpointConfig{
		IdleTimeoutInSeconds: 300,
	}
	routeParams := generateRouteCreateParamsForUnitTests(title, apiType, vHost, xWso2BasePath, version,
		endpoint.Basepath, &resourceWithGet, clusterName, nil, false)

	generatedRouteArrayWithXWso2BasePath, err := createRoutes(routeParams)
	assert.Nil(t, err, "Error while creating routes WithXWso2BasePath")
	generatedRouteWithXWso2BasePath := generatedRouteArrayWithXWso2BasePath[0]
	assert.NotNil(t, generatedRouteWithXWso2BasePath, "Route should not be null.")
	assert.Equal(t, expectedRouteActionWithXWso2BasePath, generatedRouteWithXWso2BasePath.Action,
		"Route generation mismatch when xWso2BasePath option is provided.")
	assert.NotNil(t, generatedRouteWithXWso2BasePath.GetMatch().Headers, "Headers property should not be null")
	assert.Equal(t, "^GET$", generatedRouteWithXWso2BasePath.GetMatch().Headers[0].GetStringMatch().GetSafeRegex().Regex,
		"Assigned HTTP Method Regex is incorrect when single method is available.")
}

func TestCreateRouteClusterSpecifier(t *testing.T) {
	// Tested features
	// In this test case, the extAuthz context variables are not tested
	clusterName := "cluster"

	vHost := "localhost"
	xWso2BasePath := "/xWso2BasePath"
	endpointBasePath := "/basepath"
	title := "WSO2"
	version := "1.0.0"
	apiType := "HTTP"

	endpoint := model.Endpoint{
		Host:    "abc.com",
		URLType: "http",
		Port:    80,
		RawURL:  "http://abc.com",
	}
	resourceWithGet := model.CreateMinimalDummyResourceForTests("/resourcePath", []*model.Operation{model.NewOperation("GET", nil, nil, "")},
		"resource_operation_id", []model.Endpoint{endpoint}, false, false)

	route, err := createRoutes(generateRouteCreateParamsForUnitTests(title, apiType, vHost, xWso2BasePath, version, endpointBasePath,
		&resourceWithGet, clusterName, nil, false))
	assert.Nil(t, err, "Error while creating route")
	assert.NotNil(t, route[0], "Route should not be null")
	assert.NotNil(t, route[0].GetRoute().GetClusterHeader(), "Route Cluster Header should not be null.")
	assert.Empty(t, route[0].GetRoute().GetCluster(), "Route Cluster Name should be empty.")
	assert.Equal(t, clusterHeaderName, route[0].GetRoute().GetClusterHeader(), "Route Cluster Name mismatch.")
}

// func TestCreateRouteExtAuthzContext(t *testing.T) {
// 	// Tested features
// 	// 1. The context variables inside extAuthzPerRoute configuration including
// 	// (clustername, method regex, basePath, resourcePath, title, version)
// 	clusterName := "cluster"

// 	vHost := "localhost"
// 	xWso2BasePath := "/xWso2BasePath"
// 	endpointBasePath := "/basepath"
// 	title := "WSO2"
// 	version := "1.0.0"
// 	apiType := "HTTP"

// 	endpoint := model.Endpoint{
// 		Host:    "abc.com",
// 		URLType: "http",
// 		Port:    80,
// 		RawURL:  "http://abc.com",
// 	}
// 	resourceWithGet := model.CreateMinimalDummyResourceForTests("/resourcePath", []*model.Operation{model.NewOperation("GET", nil, nil, "")},
// 		"resource_operation_id", []model.Endpoint{endpoint}, false, false)

// 	route, err := createRoutes(generateRouteCreateParamsForUnitTests(title, apiType, vHost, xWso2BasePath, version,
// 		endpointBasePath, &resourceWithGet, clusterName, nil, false))
// 	assert.Nil(t, err, "Error while creating route")
// 	assert.NotNil(t, route[0], "Route should not be null")
// 	assert.NotNil(t, route[0].GetTypedPerFilterConfig(), "TypedPerFilter config should not be null")
// 	assert.NotNil(t, route[0].GetTypedPerFilterConfig()[wellknown.HTTPExternalAuthorization],
// 		"ExtAuthzPerRouteConfig should not be empty")

// 	extAuthPerRouteConfig := &extAuthService.ExtAuthzPerRoute{}
// 	err = route[0].TypedPerFilterConfig[wellknown.HTTPExternalAuthorization].UnmarshalTo(extAuthPerRouteConfig)
// 	assert.Nilf(t, err, "Error while parsing ExtAuthzPerRouteConfig %v", extAuthPerRouteConfig)
// 	assert.NotEmpty(t, extAuthPerRouteConfig.GetCheckSettings(), "Check Settings per ext authz route should not be empty")
// 	assert.NotEmpty(t, extAuthPerRouteConfig.GetCheckSettings().ContextExtensions,
// 		"ContextExtensions per ext authz route should not be empty")

// 	contextExtensionMap := extAuthPerRouteConfig.GetCheckSettings().ContextExtensions
// 	assert.Equal(t, title, contextExtensionMap[apiNameContextExtension], "Title mismatch in route ext authz context.")
// 	assert.Equal(t, xWso2BasePath, contextExtensionMap[basePathContextExtension], "Basepath mismatch in route ext authz context.")
// 	assert.Equal(t, version, contextExtensionMap[apiVersionContextExtension], "Version mismatch in route ext authz context.")
// 	assert.Equal(t, "GET", contextExtensionMap[methodContextExtension], "Method mismatch in route ext authz context.")
// 	assert.Equal(t, clusterName, contextExtensionMap[clusterNameContextExtension], "Cluster mismatch in route ext authz context.")
// }

func TestGenerateTLSCert(t *testing.T) {
	publicKeyPath := config.GetApkHome() + "/adapter/security/localhost.pem"
	privateKeyPath := config.GetApkHome() + "/adapter/security/localhost.key"

	tlsCert := generateTLSCert(privateKeyPath, publicKeyPath)

	assert.NotNil(t, tlsCert, "TLS Certificate should not be null")

	assert.NotNil(t, tlsCert.GetPrivateKey(), "Private Key should not be null in the TLS certificate")
	assert.NotNil(t, tlsCert.GetCertificateChain(), "Certificate chain should not be null in the TLS certificate")

	assert.Equal(t, tlsCert.GetPrivateKey().GetFilename(), privateKeyPath, "Private Key Value mismatch in the TLS Certificate")
	assert.Equal(t, tlsCert.GetCertificateChain().GetFilename(), publicKeyPath, "Certificate Chain Value mismatch in the TLS Certificate")
}

func TestGenerateRegex(t *testing.T) {

	type generateRegexTestItem struct {
		pathMatchType gwapiv1.PathMatchType
		resourcePath  string
		userInputPath string
		message       string
		isMatched     bool
	}
	dataItems := []generateRegexTestItem{
		{
			pathMatchType: gwapiv1.PathMatchRegularExpression,
			resourcePath:  "/pet/[\\d]{1}",
			userInputPath: "/pet/5",
			message:       "when regex is provided end of the path",
			isMatched:     true,
		},
		{
			pathMatchType: gwapiv1.PathMatchRegularExpression,
			resourcePath:  "/pet/[\\d]{1}",
			userInputPath: "/pet/5/",
			message:       "when the input path does not have tailing slash and user input path has trailing slash",
			isMatched:     true,
		},
		{
			pathMatchType: gwapiv1.PathMatchRegularExpression,
			resourcePath:  "/pet/[\\d]{1}/",
			userInputPath: "/pet/5",
			message:       "when the input path has tailing slash and user input path does not have trailing slash",
			isMatched:     true,
		},
		{
			pathMatchType: gwapiv1.PathMatchRegularExpression,
			resourcePath:  "/pet/[\\d]{1}/",
			userInputPath: "/pet/5/",
			message:       "when both the input path and user input path has trailing slash",
			isMatched:     true,
		},
		{
			pathMatchType: gwapiv1.PathMatchRegularExpression,
			resourcePath:  "/pet/[\\d]{1}/info",
			userInputPath: "/pet/5/info",
			message:       "when regex is provided in the middle of the path",
			isMatched:     true,
		},
		{
			pathMatchType: gwapiv1.PathMatchRegularExpression,
			resourcePath:  "/pet/[\\d]{1}/tst/[\\d]{1}",
			userInputPath: "/pet/5/tst/3",
			message:       "when multiple regex match sections are provided",
			isMatched:     true,
		},
		{
			pathMatchType: gwapiv1.PathMatchRegularExpression,
			resourcePath:  "/pet/[\\d]{1}",
			userInputPath: "/pet/5/test",
			message:       "when path parameter is provided end of the path and provide incorrect path",
			isMatched:     false,
		},
		{
			pathMatchType: gwapiv1.PathMatchExact,
			resourcePath:  "/pet/5",
			userInputPath: "/pet/5",
			message:       "when using an exact match type",
			isMatched:     true,
		},
		{
			pathMatchType: gwapiv1.PathMatchExact,
			resourcePath:  "/pet/5",
			userInputPath: "/pett/5",
			message:       "when provide a incorrect path with exact match",
			isMatched:     false,
		},
		{
			pathMatchType: gwapiv1.PathMatchExact,
			resourcePath:  "/pet/5",
			userInputPath: "/pet/5/",
			message:       "when using an exact match with a trailing slash in user input only",
			isMatched:     true,
		},
		{
			pathMatchType: gwapiv1.PathMatchExact,
			resourcePath:  "/pet/5/",
			userInputPath: "/pet/5",
			message:       "when using an exact match with a trailing slash in resource path only",
			isMatched:     true,
		},
		{
			pathMatchType: gwapiv1.PathMatchExact,
			resourcePath:  "/pet/5/",
			userInputPath: "/pet/5/",
			message:       "when using an exact match with a trailing slash in user input and resource path",
			isMatched:     true,
		},
		{
			pathMatchType: gwapiv1.PathMatchPathPrefix,
			resourcePath:  "/pet",
			userInputPath: "/pet/",
			message:       "when using path prefix type match, a single trailing slash is allowed",
			isMatched:     true,
		},
		{
			pathMatchType: gwapiv1.PathMatchPathPrefix,
			resourcePath:  "/pet",
			userInputPath: "/pet",
			message:       "when using path prefix type, it can match any value after a single slash",
			isMatched:     true,
		},
		{
			pathMatchType: gwapiv1.PathMatchPathPrefix,
			resourcePath:  "/pet",
			userInputPath: "/pet/foo/bar",
			message:       "when using path prefix type, it can match several slash and value sections",
			isMatched:     true,
		},
		{
			pathMatchType: gwapiv1.PathMatchPathPrefix,
			resourcePath:  "/pet",
			userInputPath: "/pet123",
			message:       "cannot have a value without starting with a trailing slash",
			isMatched:     false,
		},
		{
			pathMatchType: gwapiv1.PathMatchPathPrefix,
			resourcePath:  "/pet/[\\d]{1}.api",
			userInputPath: "/pet/findByIdstatus=availabe",
			message:       "when the resource regex section is suffixed",
			isMatched:     false,
		},
		{
			pathMatchType: gwapiv1.PathMatchRegularExpression,
			resourcePath:  "/pet/[a-z0-9]+.api",
			userInputPath: "/pet/pet1.api",
			message:       "when the resource path param suffixed",
			isMatched:     true,
		},
		{
			pathMatchType: gwapiv1.PathMatchRegularExpression,
			resourcePath:  "/pet/pet[a-z0-9]+",
			userInputPath: "/pet/findByIdstatus=availabe",
			message:       "when the resource ends with regex section",
			isMatched:     false,
		},
		{
			pathMatchType: gwapiv1.PathMatchRegularExpression,
			resourcePath:  "/pet/pet[a-z0-9]+",
			userInputPath: "/pet/pet1",
			message:       "when the resource ends with regex section",
			isMatched:     true,
		},
	}

	for _, item := range dataItems {
		resultPattern := generateRoutePath(item.resourcePath, item.pathMatchType)
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
		pathMatchType      gwapiv1.PathMatchType
		inputPath          string
		expectedSubsString string
		message            string
		shouldEqual        bool
	}
	dataItems := []generateSubsStringTestItem{
		{
			gwapiv1.PathMatchExact,
			"/v2/pet/",
			"/v2/pet/",
			"when using exact type",
			true,
		},
		{
			gwapiv1.PathMatchPathPrefix,
			"/v2/pet/",
			"/v2/pet\\1",
			"when using path prefix type",
			true,
		},
		{
			gwapiv1.PathMatchRegularExpression,
			"/v2/pet/(dog-[\\d]{2})",
			"\\1",
			"when using regex type",
			true,
		},
		{
			gwapiv1.PathMatchExact,
			"/v2/pet",
			"/v2/pet",
			"when using exact type without a trailing slash in the path",
			true,
		},
		{
			gwapiv1.PathMatchPathPrefix,
			"/v2/pet",
			"/v2/pet\\1",
			"when using path prefix type without a trailing slash in the path",
			true,
		},
		{
			gwapiv1.PathMatchRegularExpression,
			"/v2/pet/(dog-[\\d]{2})/",
			"\\1",
			"when using regex type with a trailing slash in the path",
			true,
		},
	}
	for _, item := range dataItems {
		generatedSubstitutionString := generateSubstitutionString(item.inputPath, item.pathMatchType)
		if item.shouldEqual {
			assert.Equal(t, item.expectedSubsString, generatedSubstitutionString, item.message)
		} else {
			assert.NotEqual(t, item.expectedSubsString, generatedSubstitutionString, item.message)
		}
	}
}

func TestCreateUpstreamTLSContext(t *testing.T) {
	certFilePath := config.GetApkHome() + "/test-resources/testcrt.crt"
	certByteArr, err := ioutil.ReadFile(certFilePath)
	assert.Nil(t, err, "Error while reading the certificate : "+certFilePath)
	defaultAPKKeyPath := "/home/wso2/security/keystore/mg.key"
	defaultAPKCertPath := "/home/wso2/security/keystore/mg.pem"
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

	tlsCert := generateTLSCert(defaultAPKKeyPath, defaultAPKCertPath)
	upstreamTLSContextWithCerts := createUpstreamTLSContext(certByteArr, nil, hostNameAddress, false)
	upstreamTLSContextWithoutCerts := createUpstreamTLSContext(nil, nil, hostNameAddress, false)
	upstreamTLSContextWithIP := createUpstreamTLSContext(certByteArr, nil, hostNameAddressWithIP, false)

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
	endpoint := model.Endpoint{
		Host:    "abc.com",
		URLType: "http",
		Port:    80,
		RawURL:  "http://abc.com",
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

	resourceWithGet := model.CreateMinimalDummyResourceForTests("/resourcePath", []*model.Operation{model.NewOperation("GET", nil, nil, "")},
		"resource_operation_id", []model.Endpoint{endpoint}, false, false)

	// Route without CORS configuration
	routeWithoutCors, err := createRoutes(generateRouteCreateParamsForUnitTests("test", "HTTP", "localhost", "/test", "1.0.0", "/test",
		&resourceWithGet, "test-cluster", nil, false))
	assert.Nil(t, err, "Error while creating routeWithoutCors")

	corsConfig1 := &cors_filter_v3.CorsPolicy{}
	err = routeWithoutCors[0].GetTypedPerFilterConfig()[wellknown.CORS].UnmarshalTo(corsConfig1)

	assert.Nilf(t, err, "Error while parsing Cors Configuration %v", corsConfig1)
	assert.Empty(t, corsConfig1.GetAllowHeaders(), "Cors AllowHeaders should be empty.")
	assert.Empty(t, corsConfig1.GetAllowCredentials(), "Cors AllowCredentials should be empty.")
	assert.Empty(t, corsConfig1.GetAllowMethods(), "Cors AllowMethods should be empty.")
	assert.Empty(t, corsConfig1.GetAllowOriginStringMatch(), "Cors AllowOriginStringMatch should be empty.")
	assert.Empty(t, corsConfig1.GetExposeHeaders(), "Cors ExposeHeaders should be empty.")

	// Route with CORS configuration
	routeWithCors, err := createRoutes(generateRouteCreateParamsForUnitTests("test", "HTTP", "localhost", "/test", "1.0.0", "/test",
		&resourceWithGet, "test-cluster", corsConfigModel3, false))
	assert.Nil(t, err, "Error while creating routeWithCors")

	corsConfig2 := &cors_filter_v3.CorsPolicy{}
	err = routeWithCors[0].GetTypedPerFilterConfig()[wellknown.CORS].UnmarshalTo(corsConfig2)

	assert.Nilf(t, err, "Error while parsing Cors Configuration %v", corsConfig2)
	assert.NotEmpty(t, corsConfig2.GetAllowOriginStringMatch(), "Cors AllowOriginStringMatch should not be empty.")
	assert.NotEmpty(t, corsConfig2.GetAllowMethods(), "Cors AllowMethods should not be empty.")
	assert.Empty(t, corsConfig2.GetAllowHeaders(), "Cors AllowHeaders should be empty.")
	assert.Empty(t, corsConfig2.GetExposeHeaders(), "Cors ExposeHeaders should be empty.")
	assert.Empty(t, corsConfig2.GetAllowCredentials(), "Cors AllowCredentials should be empty.")
}

func generateRouteCreateParamsForUnitTests(title string, apiType string, vhost string, xWso2Basepath string, version string, endpointBasepath string,
	resource *model.Resource, clusterName string, corsConfig *model.CorsConfig, isDefaultVersion bool) *routeCreateParams {
	return &routeCreateParams{
		title:            title,
		apiType:          apiType,
		version:          version,
		vHost:            vhost,
		xWSO2BasePath:    xWso2Basepath,
		resource:         resource,
		clusterName:      clusterName,
		endpointBasePath: endpointBasepath,
		corsPolicy:       corsConfig,
		isDefaultVersion: isDefaultVersion,
	}
}
