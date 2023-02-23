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
	"github.com/wso2/apk/test/integration/integration/utils/kubernetes"
	"github.com/wso2/apk/test/integration/integration/utils/suite"
)

func init() {
	IntegrationTests = append(IntegrationTests, ALLHTTPMethodsForWildCard)
}

// APIWithPathParams test
var ALLHTTPMethodsForWildCard = suite.IntegrationTest{
	ShortName:   "ALLHTTPMethodsForWildCard",
	Description: "Tests an API with wild card path using path prefix match and unspecified HTTP method",
	Manifests:   []string{"tests/all-http-methods-for-wildcard.yaml"},
	Test: func(t *testing.T, suite *suite.IntegrationTestSuite) {
		ns := "gateway-integration-test-infra"
		gwAddr := kubernetes.WaitForGatewayAddress(t, suite.Client, suite.TimeoutConfig)
		token := http.GetTestToken(t, gwAddr)

		testCases := []http.ExpectedResponse{
			// test path with trailing slash for GET
			{
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard/v1.0.0/",
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
			// {
			// 	Request: http.Request{
			// 		Host:   "all-http-methods-for-wildcard.test.gw.wso2.com",
			// 		Path:   "/all-http-methods-for-wildcard/v1.0.0",
			// 		Method: "GET",
			// 	},
			// 	ExpectedRequest: &http.ExpectedRequest{
			// 		Request: http.Request{
			// 			Path: "/v2/echo-full",
			// 		},
			// 	},
			// 	Backend:   "infra-backend-v1",
			// 	Namespace: ns,
			// },
			{
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard/v1.0.0/test",
					Method: "GET",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path: "/v2/echo-full/test",
					},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
			// {
			// 	Request: http.Request{
			// 		Host:   "all-http-methods-for-wildcard.test.gw.wso2.com",
			// 		Path:   "/all-http-methods-for-wildcard/v1.0.0/test1/test2/test3?foo=foo1&bar=bar1",
			// 		Method: "GET",
			// 	},
			// 	ExpectedRequest: &http.ExpectedRequest{
			// 		Request: http.Request{
			// 			Path: "/v2/echo-full/test1/test2/test3",
			// 		},
			// 	},
			// 	Backend:   "infra-backend-v1",
			// 	Namespace: ns,
			// },
			// test path with trailing slash for POST
			{
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard/v1.0.0/",
					Method: "POST",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path:   "/v2/echo-full/",
						Method: "POST",
					},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
			// {
			// 	Request: http.Request{
			// 		Host:   "all-http-methods-for-wildcard.test.gw.wso2.com",
			// 		Path:   "/all-http-methods-for-wildcard/v1.0.0",
			// 		Method: "POST",
			// 	},
			// 	ExpectedRequest: &http.ExpectedRequest{
			// 		Request: http.Request{
			// 			Path: "/v2/echo-full",
			// 			Method: "POST",
			// 		},
			// 	},
			// 	Backend:   "infra-backend-v1",
			// 	Namespace: ns,
			// },
			{
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard/v1.0.0/test",
					Method: "POST",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path:   "/v2/echo-full/test",
						Method: "POST",
					},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
			// test path with trailing slash for PUT
			{
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard/v1.0.0/",
					Method: "PUT",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path:   "/v2/echo-full/",
						Method: "PUT",
					},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
			// {
			// 	Request: http.Request{
			// 		Host:   "all-http-methods-for-wildcard.test.gw.wso2.com",
			// 		Path:   "/all-http-methods-for-wildcard/v1.0.0",
			// 		Method: "PUT",
			// 	},
			// 	ExpectedRequest: &http.ExpectedRequest{
			// 		Request: http.Request{
			// 			Path: "/v2/echo-full",
			// 			Method: "PUT",
			// 		},
			// 	},
			// 	Backend:   "infra-backend-v1",
			// 	Namespace: ns,
			// },
			{
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard/v1.0.0/test",
					Method: "PUT",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path:   "/v2/echo-full/test",
						Method: "PUT",
					},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
			// test path with trailing slash for PATCH
			{
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard/v1.0.0/",
					Method: "PATCH",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path:   "/v2/echo-full/",
						Method: "PATCH",
					},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
			// {
			// 	Request: http.Request{
			// 		Host:   "all-http-methods-for-wildcard.test.gw.wso2.com",
			// 		Path:   "/all-http-methods-for-wildcard/v1.0.0",
			// 		Method: "PATCH",
			// 	},
			// 	ExpectedRequest: &http.ExpectedRequest{
			// 		Request: http.Request{
			// 			Path: "/v2/echo-full",
			// 			Method: "PATCH",
			// 		},
			// 	},
			// 	Backend:   "infra-backend-v1",
			// 	Namespace: ns,
			// },
			{
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard/v1.0.0/test",
					Method: "PATCH",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path:   "/v2/echo-full/test",
						Method: "PATCH",
					},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
			// test path with trailing slash for DELETE
			{
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard/v1.0.0/",
					Method: "DELETE",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path:   "/v2/echo-full/",
						Method: "DELETE",
					},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
			// {
			// 	Request: http.Request{
			// 		Host:   "all-http-methods-for-wildcard.test.gw.wso2.com",
			// 		Path:   "/all-http-methods-for-wildcard/v1.0.0",
			// 		Method: "DELETE",
			// 	},
			// 	ExpectedRequest: &http.ExpectedRequest{
			// 		Request: http.Request{
			// 			Path: "/v2/echo-full",
			// 			Method: "DELETE",
			// 		},
			// 	},
			// 	Backend:   "infra-backend-v1",
			// 	Namespace: ns,
			// },
			{
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard/v1.0.0/test",
					Method: "DELETE",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path:   "/v2/echo-full/test",
						Method: "DELETE",
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
