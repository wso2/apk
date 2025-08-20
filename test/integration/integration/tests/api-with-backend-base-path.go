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
	// We no longer support backend base path as a special feature. If someone wants a base path they need to use the path rewrite filters
	//IntegrationTests = append(IntegrationTests, APIWithBackendBasePath)
}

// APIWithBackendBasePath test
var APIWithBackendBasePath = suite.IntegrationTest{
	ShortName:   "APIWithBackendBasePath",
	Description: "An API with a backend base path should be able to route requests to the backend",
	Manifests:   []string{"tests/api-with-backend-base-path.yaml"},
	Test: func(t *testing.T, suite *suite.IntegrationTestSuite) {
		ns := "gateway-integration-test-infra"
		gwAddr := "backend-base-path.test.gw.wso2.com:9095"
		token := http.GetTestToken(t)

		testCases := []http.ExpectedResponse{
			{
				Request: http.Request{
					Host: "backend-base-path.test.gw.wso2.com",
					Path: "/test-api-with-backend-base-path/1.0.0/orders",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path: "/backend-base-path/orders",
					},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
			{
				Request: http.Request{
					Host: "backend-base-path.test.gw.wso2.com",
					Path: "/test-api-with-backend-base-path/1.0.0/users",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path: "/backend-base-path/users",
					},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
			{
				Request: http.Request{
					Host: "backend-base-path.test.gw.wso2.com",
					Path: "/test-api-with-backend-base-path/orders",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path: "/backend-base-path/orders",
					},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
			{
				Request: http.Request{
					Host: "backend-base-path.test.gw.wso2.com",
					Path: "/test-api-with-backend-base-path/users",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path: "/backend-base-path/users",
					},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
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
