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
	//IntegrationTests = append(IntegrationTests, APIWithRequestHeaderModify)
}

// APIWithRequestHeaderModify test
var APIWithRequestHeaderModify = suite.IntegrationTest{
	ShortName:   "APIWithRequestHeaderModify",
	Description: "An API with request header modify",
	Manifests:   []string{"tests/api-with-request-header-modify.yaml"},
	Test: func(t *testing.T, suite *suite.IntegrationTestSuite) {
		ns := "gateway-integration-test-infra"
		gwAddr := "backend-base-path.test.gw.wso2.com:9095"
		token := http.GetTestToken(t)

		testCases := []http.ExpectedResponse{
			{
				Request: http.Request{
					Host: "backend-base-path.test.gw.wso2.com",
					Path: "/request-header-modify/1.0.0/users",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path: "/users",
					},
					AbsentHeaders: []string{"X-Header-Remove"},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
			{
				Request: http.Request{
					Host: "backend-base-path.test.gw.wso2.com",
					Path: "/request-header-modify/1.0.0/orders",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path:    "/orders",
						Headers: map[string]string{"test-header": "test"},
					},
					AbsentHeaders: []string{"X-Header-add"},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
			{
				Request: http.Request{
					Host: "backend-base-path.test.gw.wso2.com",
					Path: "/request-header-modify/users",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path: "/users",
					},
					AbsentHeaders: []string{"X-Header-Remove"},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
			{
				Request: http.Request{
					Host: "backend-base-path.test.gw.wso2.com",
					Path: "/request-header-modify/orders",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path:    "/orders",
						Headers: map[string]string{"test-header": "test"},
					},
					AbsentHeaders: []string{"X-Header-add"},
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
