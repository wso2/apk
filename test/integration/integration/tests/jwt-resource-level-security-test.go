/*
 *  Copyright (c) 2024, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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
	IntegrationTests = append(IntegrationTests, ResourceLevelJWT)
}

// ResourceLevelJWT test
var ResourceLevelJWT = suite.IntegrationTest{
	ShortName:   "ResourceLevelJWT",
	Description: "Tests resource level jwt security",
	Manifests:   []string{"tests/jwt-resource-level-security.yaml"},
	Test: func(t *testing.T, suite *suite.IntegrationTestSuite) {
		gwAddr := "resource-level-jwt.test.gw.wso2.com:9095"
		token := http.GetTestToken(t)
		ns := "gateway-integration-test-infra"
		testCases := []http.ExpectedResponse{
			{
				Request: http.Request{
					Host: "resource-level-jwt.test.gw.wso2.com",
					Path: "/resource-level-jwt/v1.0.0/v2/echo-full",
					Headers: map[string]string{
						"content-type": "application/json",
						"internal-key": token,
					},
					Method: "GET",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path: "/v2/echo-full",
					},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
				Response: http.Response{StatusCode: 200},
			},
			{
				Request: http.Request{
					Host: "resource-level-jwt.test.gw.wso2.com",
					Path: "/resource-level-jwt/v1.0.0/v2/echo-full",
					Headers: map[string]string{
						"content-type": "application/json",
					},
					Method: "GET",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path: "/v2/echo-full",
					},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
				Response: http.Response{StatusCode: 401},
			},
			{
				Request: http.Request{
					Host: "resource-level-jwt.test.gw.wso2.com",
					Path: "/resource-level-jwt/v1.0.0/v2/echo-full",
					Headers: map[string]string{
						"content-type": "application/json",
						"internal-key": "invalid",
					},
					Method: "GET",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path: "/v2/echo-full",
					},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
				Response:  http.Response{StatusCode: 401},
			},
			// Test wrong audience
			{
				Request: http.Request{
					Host: "resource-level-jwt.test.gw.wso2.com",
					Path: "/resource-level-jwt/v1.0.0/v2/echo-1",
					Headers: map[string]string{
						"content-type": "application/json",
						"internal-key": token,
					},
					Method: "GET",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path: "/v2/echo-1",
					},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
				Response:  http.Response{StatusCode: 401},
			},
			// Test correct audience
			{
				Request: http.Request{
					Host: "resource-level-jwt.test.gw.wso2.com",
					Path: "/resource-level-jwt/v1.0.0/v2/echo-2",
					Headers: map[string]string{
						"content-type": "application/json",
						"internal-key": token,
					},
					Method: "GET",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path: "/v2/echo-2",
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
