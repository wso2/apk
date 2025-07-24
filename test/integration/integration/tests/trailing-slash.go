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
	//IntegrationTests = append(IntegrationTests, TrailingSlash)
}

// TrailingSlash test
var TrailingSlash = suite.IntegrationTest{
	ShortName:   "TrailingSlash",
	Description: "Invoking API with and without trailing slash",
	Manifests:   []string{"tests/trailing-slash.yaml"},
	Test: func(t *testing.T, suite *suite.IntegrationTestSuite) {
		ns := "gateway-integration-test-infra"
		gwAddr := "trailing-slash.test.gw.wso2.com:9095"
		token := http.GetTestToken(t, gwAddr)

		testCases := []http.ExpectedResponse{
			// test path with trailing slash but without path parameters
			{
				Request: http.Request{
					Host: "trailing-slash.test.gw.wso2.com",
					Path: "/trailing-slash/v1.0.0/echo-full/no-slash/findByStatus",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path: "/v2/echo-full/no-slash/findByStatus",
					},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
			{
				Request: http.Request{
					Host: "trailing-slash.test.gw.wso2.com",
					Path: "/trailing-slash/v1.0.0/echo-full/no-slash/findByStatus/",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path: "/v2/echo-full/no-slash/findByStatus",
					},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
			// test path with trailing slash with path parameter
			{
				Request: http.Request{
					Host: "trailing-slash.test.gw.wso2.com",
					Path: "/trailing-slash/v1.0.0/echo-full/no-slash/1",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path: "/v2/echo-full/no-slash/1",
					},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
			{
				Request: http.Request{
					Host: "trailing-slash.test.gw.wso2.com",
					Path: "/trailing-slash/v1.0.0/echo-full/no-slash/1/",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path: "/v2/echo-full/no-slash/1",
					},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
			// test path with trailing slash with multiple path parameters
			{
				Request: http.Request{
					Host: "trailing-slash.test.gw.wso2.com",
					Path: "/trailing-slash/v1.0.0/echo-full/no-slash/1/pet/123",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path: "/v2/echo-full/no-slash/1/pet/123",
					},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
			{
				Request: http.Request{
					Host: "trailing-slash.test.gw.wso2.com",
					Path: "/trailing-slash/v1.0.0/echo-full/no-slash/1/pet/123/",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path: "/v2/echo-full/no-slash/1/pet/123",
					},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
			// test path with trailing slash but without path parameters
			{
				Request: http.Request{
					Host: "trailing-slash.test.gw.wso2.com",
					Path: "/trailing-slash/v1.0.0/echo-full/with-slash/findByStatus",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path: "/v2/echo-full/with-slash/findByStatus/",
					},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
			{
				Request: http.Request{
					Host: "trailing-slash.test.gw.wso2.com",
					Path: "/trailing-slash/v1.0.0/echo-full/with-slash/findByStatus/",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path: "/v2/echo-full/with-slash/findByStatus/",
					},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
			// test path with trailing slash but without path parameters
			{
				Request: http.Request{
					Host: "trailing-slash.test.gw.wso2.com",
					Path: "/trailing-slash/v1.0.0/echo-full/with-slash/1",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path: "/v2/echo-full/with-slash/1/",
					},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
			{
				Request: http.Request{
					Host: "trailing-slash.test.gw.wso2.com",
					Path: "/trailing-slash/v1.0.0/echo-full/with-slash/1/",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path: "/v2/echo-full/with-slash/1/",
					},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
			// test path with trailing slash with multiple path parameters
			{
				Request: http.Request{
					Host: "trailing-slash.test.gw.wso2.com",
					Path: "/trailing-slash/v1.0.0/echo-full/with-slash/1/pet/123",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path: "/v2/echo-full/with-slash/1/pet/123/",
					},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
			{
				Request: http.Request{
					Host: "trailing-slash.test.gw.wso2.com",
					Path: "/trailing-slash/v1.0.0/echo-full/with-slash/1/pet/123/",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path: "/v2/echo-full/with-slash/1/pet/123/",
					},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
			// test path with additional chars
			{
				Request: http.Request{
					Host: "trailing-slash.test.gw.wso2.com",
					Path: "/trailing-slash/v1.0.0/echo-full/chars",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path: "/v2/echo-full/chars",
					},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
			{
				Request: http.Request{
					Host: "trailing-slash.test.gw.wso2.com",
					Path: "/trailing-slash/v1.0.0/echo-full/charsAdditional",
				},
				Response: http.Response{StatusCode: 404},
			},
			// test path with additional chars after trailing slash
			{
				Request: http.Request{
					Host: "trailing-slash.test.gw.wso2.com",
					Path: "/trailing-slash/v1.0.0/echo-full/with-param/1/",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path: "/v2/echo-full/with-param/1",
					},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
			{
				Request: http.Request{
					Host: "trailing-slash.test.gw.wso2.com",
					Path: "/trailing-slash/v1.0.0/echo-full/with-param/1/additional",
				},
				Response: http.Response{StatusCode: 404},
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
