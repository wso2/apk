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
	IntegrationTests = append(IntegrationTests, Subscription)
}

// Subscription test
var Subscription = suite.IntegrationTest{
	ShortName:   "Subscription",
	Description: "Tests the subscription functionality",
	Manifests:   []string{"tests/subscription.yaml"},
	Test: func(t *testing.T, suite *suite.IntegrationTestSuite) {
		ns := "gateway-integration-test-infra"
		gwAddr := "all-http-methods-for-wildcard.test.gw.wso2.com:9095"
		clientToken := http.GetClientCredentialsToken(t, "45f1c5c8-a92e-11ed-afa1-0242ac120065")
		appToken := http.GetAppClientCredentialsToken(t, "llm-app")
		token := http.GetTestToken(t)
		testCases := []http.ExpectedResponse{
			// test path with trailing slash for GET
			{
				TestCaseName: "test path with trailing slash for GET",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard/v1.0.0/v2/echo-full/",
					Method: "GET",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path: "/v2/echo-full/",
					},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
		}

		testCasesWithNoSubsToken := []http.ExpectedResponse{
			// test path with trailing slash for GET
			{
				TestCaseName: "test path with trailing slash for GET with no subscription token",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard/v1.0.0/v2/echo-full/",
					Method: "GET",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path: "/v2/echo-full/",
					},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
				Response: http.Response{
					StatusCode: 403,
				},
			},
		}

		for i := range testCases {
			tc := testCases[i]
			tc.Request.Headers = http.AddBearerTokenToHeader(clientToken, tc.Request.Headers)
			t.Run(tc.GetTestCaseName(i), func(t *testing.T) {
				t.Parallel()
				http.MakeRequestAndExpectEventuallyConsistentResponse(t, suite.RoundTripper, suite.TimeoutConfig, gwAddr, tc)
			})
		}
		for i := range testCases {
			tc := testCases[i]
			tc.Request.Headers = http.AddBearerTokenToHeader(appToken, tc.Request.Headers)
			t.Run(tc.GetTestCaseName(i), func(t *testing.T) {
				t.Parallel()
				http.MakeRequestAndExpectEventuallyConsistentResponse(t, suite.RoundTripper, suite.TimeoutConfig, gwAddr, tc)
			})
		}
		for i := range testCasesWithNoSubsToken {
			tc := testCasesWithNoSubsToken[i]
			tc.Request.Headers = http.AddBearerTokenToHeader(token, tc.Request.Headers)
			t.Run(tc.GetTestCaseName(i), func(t *testing.T) {
				t.Parallel()
				http.MakeRequestAndExpectEventuallyConsistentResponse(t, suite.RoundTripper, suite.TimeoutConfig, gwAddr, tc)
			})
		}
	},
}
