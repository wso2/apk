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
	//IntegrationTests = append(IntegrationTests, DefaultAPIVersionRatelimit)
}

// DefaultAPIVersionRatelimit tests default api version feature
var DefaultAPIVersionRatelimit = suite.IntegrationTest{
	ShortName:   "DefaultAPIVersionRatelimitTest",
	Description: "Tests an API with default version",
	Manifests:   []string{"tests/default-api-version-ratelimit.yaml"},
	Test: func(t *testing.T, testSuite *suite.IntegrationTestSuite) {
		ns := "gateway-integration-test-infra"
		gwAddr := "default-api-version-ratelimit.test.gw.wso2.com:9095"
		token := http.GetTestToken(t)

		testCases := []http.ExpectedResponse{
			{
				Request: http.Request{
					Host:   "default-api-version-ratelimit.test.gw.wso2.com",
					Path:   "/default-api-version-ratelimit/v2/echo-full",
					Method: "GET",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path: "/v2/echo-full",
					},
				},
				Backend:              "infra-backend-v1",
				Namespace:            ns,
				Response:             http.Response{StatusCode: 200},
				UnacceptableStatuses: []int{429},
			},
			{
				Request: http.Request{
					Host:   "default-api-version-ratelimit.test.gw.wso2.com",
					Path:   "/default-api-version-ratelimit/v1.0.0/v2/echo-full",
					Method: "GET",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path: "/v2/echo-full",
					},
				},
				Backend:              "infra-backend-v1",
				Namespace:            ns,
				Response:             http.Response{StatusCode: 200},
				UnacceptableStatuses: []int{429},
			},
			{
				Request: http.Request{
					Host:   "default-api-version-ratelimit.test.gw.wso2.com",
					Path:   "/default-api-version-ratelimit/v1.0.0/v2/echo-full",
					Method: "GET",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path: "/v2/echo-full",
					},
				},
				Backend:              "infra-backend-v1",
				Namespace:            ns,
				Response:             http.Response{StatusCode: 429},
				UnacceptableStatuses: []int{200},
			},
			{
				Request: http.Request{
					Host:   "default-api-version-ratelimit.test.gw.wso2.com",
					Path:   "/default-api-version-ratelimit/v2/echo-full",
					Method: "GET",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path: "/v2/echo-full",
					},
				},
				Backend:              "infra-backend-v1",
				Namespace:            ns,
				Response:             http.Response{StatusCode: 429},
				UnacceptableStatuses: []int{200},
			},
		}
		suite.WaitForNextMinute(t)
		for i := range testCases {
			tc := testCases[i]
			tc.Request.Headers = http.AddBearerTokenToHeader(token, tc.Request.Headers)
			t.Run(tc.GetTestCaseName(i), func(t *testing.T) {
				http.MakeRequestAndExpectEventuallyConsistentResponse(t, testSuite.RoundTripper, testSuite.TimeoutConfig, gwAddr, tc)
			})
		}
	},
}
