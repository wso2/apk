/*
 *  Copyright (c) 2023, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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

package tests

import (
	"testing"

	"github.com/wso2/apk/test/integration/integration/utils/http"
	"github.com/wso2/apk/test/integration/integration/utils/suite"
)

func init() {
	// //IntegrationTests = append(IntegrationTests, DisableResourceLevelJWT)
	//IntegrationTests = append(IntegrationTests, DisableResourceLevelJWTWithFalseValueTest)
	//IntegrationTests = append(IntegrationTests, DisableResourceLevelJWTWithNoOtherAuth)
}

// DisableResourceLevelJWT tests disabling and enabling jwt feature resource level with disabled = true value
var DisableResourceLevelJWT = suite.IntegrationTest{
	ShortName:   "DisableResourceLevelJWTTest",
	Description: "Tests disabled JWT in resource level",
	Manifests:   []string{"tests/disable-resource-level-jwt.yaml"},
	Test: func(t *testing.T, suite *suite.IntegrationTestSuite) {
		ns := "gateway-integration-test-infra"
		gwAddr := "disable-resource-level-jwt.test.gw.wso2.com:9095"
		token := http.GetTestToken(t)

		testCases := []http.ExpectedResponse{
			{
				Request: http.Request{
					Host:   "disable-resource-level-jwt.test.gw.wso2.com",
					Path:   "/disable-resource-level-jwt/v1.0.0/v2/echo-full/",
					Method: "GET",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path: "/v2/echo-full/",
					},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
				Response:  http.Response{StatusCode: 401},
			},
			{
				Request: http.Request{
					Host:   "disable-resource-level-jwt.test.gw.wso2.com",
					Path:   "/disable-resource-level-jwt/v2/echo-full/",
					Method: "GET",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path: "/v2/echo-full/",
					},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
				Response:  http.Response{StatusCode: 401},
			},
		}
		for i := range testCases {
			tc := testCases[i]
			tc.Request.Headers = http.AddBearerTokenToHeader(token, tc.Request.Headers)
			t.Run(tc.GetTestCaseName(i), func(t *testing.T) {
				t.Parallel()
				http.MakeRequestAndExpectEventuallyConsistentResponse(t, suite.RoundTripper, suite.TimeoutConfig, gwAddr, tc)
			})
		}
	},
}

// DisableResourceLevelJWTWithFalseValueTest tests disabling and enabling jwt feature resource level with disabled = false value
var DisableResourceLevelJWTWithFalseValueTest = suite.IntegrationTest{
	ShortName:   "DisableResourceLevelJWTWithFalseValueTest",
	Description: "Tests enabled JWT in API level",
	Manifests:   []string{"tests/disable-resource-level-jwt.yaml"},
	Test: func(t *testing.T, suite *suite.IntegrationTestSuite) {
		ns := "gateway-integration-test-infra"
		gwAddr := "disable-resource-level-jwt1.test.gw.wso2.com:9095"
		token := http.GetTestToken(t)

		testCases := []http.ExpectedResponse{
			{
				Request: http.Request{
					Host:   "disable-resource-level-jwt1.test.gw.wso2.com",
					Path:   "/disable-resource-level-jwt1/v1.0.0/v2/echo-full/",
					Method: "GET",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path: "/v2/echo-full/",
					},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
				Response:  http.Response{StatusCode: 200},
			},
			{
				Request: http.Request{
					Host:   "disable-resource-level-jwt1.test.gw.wso2.com",
					Path:   "/disable-resource-level-jwt1/v2/echo-full/",
					Method: "GET",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path: "/v2/echo-full/",
					},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
				Response:  http.Response{StatusCode: 200},
			},
		}
		for i := range testCases {
			tc := testCases[i]
			tc.Request.Headers = http.AddBearerTokenToHeader(token, tc.Request.Headers)
			t.Run(tc.GetTestCaseName(i), func(t *testing.T) {
				t.Parallel()
				http.MakeRequestAndExpectEventuallyConsistentResponse(t, suite.RoundTripper, suite.TimeoutConfig, gwAddr, tc)
			})
		}
	},
}

// DisableResourceLevelJWTWithNoOtherAuth tests disabling and enabling jwt feature resource level with disabled = false value
var DisableResourceLevelJWTWithNoOtherAuth = suite.IntegrationTest{
	ShortName:   "DisableResourceLevelJWTWithNoOtherAuthTest",
	Description: "Tests enabled JWT in API level with no other auth",
	Manifests:   []string{"tests/disable-resource-level-jwt.yaml"},
	Test: func(t *testing.T, suite *suite.IntegrationTestSuite) {
		ns := "gateway-integration-test-infra"
		gwAddr := "disable-resource-level-jwt2.test.gw.wso2.com:9095"

		testCases := []http.ExpectedResponse{
			{
				Request: http.Request{
					Host:   "disable-resource-level-jwt2.test.gw.wso2.com",
					Path:   "/disable-resource-level-jwt2/v1.0.0/v2/echo-full/",
					Method: "GET",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path: "/v2/echo-full/",
					},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
				Response:  http.Response{StatusCode: 200},
			},
			{
				Request: http.Request{
					Host:   "disable-resource-level-jwt2.test.gw.wso2.com",
					Path:   "/disable-resource-level-jwt2/v2/echo-full/",
					Method: "GET",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path: "/v2/echo-full/",
					},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
				Response:  http.Response{StatusCode: 200},
			},
		}
		for i := range testCases {
			tc := testCases[i]
			t.Run(tc.GetTestCaseName(i), func(t *testing.T) {
				t.Parallel()
				http.MakeRequestAndExpectEventuallyConsistentResponse(t, suite.RoundTripper, suite.TimeoutConfig, gwAddr, tc)
			})
		}
	},
}
