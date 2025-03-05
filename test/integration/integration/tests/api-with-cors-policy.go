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
	IntegrationTests = append(IntegrationTests, APIWithCORSPolicy)
}

// APIWithCORSPolicy test
var APIWithCORSPolicy = suite.IntegrationTest{
	ShortName:   "APIWithCORSPolicy",
	Description: "Tests API with CORS policy",
	Manifests:   []string{"tests/api-with-cors-policy.yaml"},
	Test: func(t *testing.T, suite *suite.IntegrationTestSuite) {

		gwAddr := "cors-policy.test.gw.wso2.com:9095"

		testCases := []http.ExpectedResponse{
			{
				Request: http.Request{
					Host: "cors-policy.test.gw.wso2.com",
					Path: "/cors-policy-api/1.0.0/test",
					Headers: map[string]string{
						"origin":                        "apk.wso2.com",
						"access-control-request-method": "GET",
					},
					Method: "OPTIONS",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Host:   "",
						Method: "OPTIONS",
					},
				},
				Response: http.Response{
					Headers: map[string]string{
						"access-control-allow-origin":      "apk.wso2.com",
						"access-control-allow-credentials": "true",
						"access-control-allow-methods":     "GET, POST",
						"access-control-allow-headers":     "authorization",
						"access-control-expose-headers":    "*",
					},
					StatusCode: 200,
				},
			},
			// {
			// 	Request: http.Request{
			// 		Host: "cors-policy.test.gw.wso2.com",
			// 		Path: "/cors-policy-api/1.0.0/test",
			// 		Headers: map[string]string{
			// 			"origin":                        "apk.wso2.org",
			// 			"access-control-request-method": "GET",
			// 		},
			// 		Method: "OPTIONS",
			// 	},
			// 	ExpectedRequest: &http.ExpectedRequest{
			// 		Request: http.Request{
			// 			Host:   "",
			// 			Method: "OPTIONS",
			// 		},
			// 	},
			// 	Response: http.Response{
			// 		Headers: map[string]string{
			// 			"allow": "OPTIONS, GET",
			// 		},
			// 		StatusCode: 200,
			// 	},
			// },
			{
				Request: http.Request{
					Host: "cors-policy.test.gw.wso2.com",
					Path: "/no-cors-policy-api/1.0.0/test",
					Headers: map[string]string{
						"origin":                        "apk.wso2.com",
						"access-control-request-method": "GET",
					},
					Method: "OPTIONS",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Host:   "",
						Method: "OPTIONS",
					},
				},
				Response: http.Response{
					StatusCode: 404,
				},
			},
			// Check for default api path
			{
				Request: http.Request{
					Host: "cors-policy.test.gw.wso2.com",
					Path: "/cors-policy-api/test",
					Headers: map[string]string{
						"origin":                        "apk.wso2.com",
						"access-control-request-method": "GET",
					},
					Method: "OPTIONS",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Host:   "",
						Method: "OPTIONS",
					},
				},
				Response: http.Response{
					Headers: map[string]string{
						"access-control-allow-origin":      "apk.wso2.com",
						"access-control-allow-credentials": "true",
						"access-control-allow-methods":     "GET, POST",
						"access-control-allow-headers":     "authorization",
						"access-control-expose-headers":    "*",
					},
					StatusCode: 200,
				},
			},
			// {
			// 	Request: http.Request{
			// 		Host: "cors-policy.test.gw.wso2.com",
			// 		Path: "/cors-policy-api/test",
			// 		Headers: map[string]string{
			// 			"origin":                        "apk.wso2.org",
			// 			"access-control-request-method": "GET",
			// 		},
			// 		Method: "OPTIONS",
			// 	},
			// 	ExpectedRequest: &http.ExpectedRequest{
			// 		Request: http.Request{
			// 			Host:   "",
			// 			Method: "OPTIONS",
			// 		},
			// 	},
			// 	Response: http.Response{
			// 		Headers: map[string]string{
			// 			"allow": "OPTIONS, GET",
			// 		},
			// 		StatusCode: 200,
			// 	},
			// },
			{
				Request: http.Request{
					Host: "cors-policy.test.gw.wso2.com",
					Path: "/no-cors-policy-api/test",
					Headers: map[string]string{
						"origin":                        "apk.wso2.com",
						"access-control-request-method": "GET",
					},
					Method: "OPTIONS",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Host:   "",
						Method: "OPTIONS",
					},
				},
				Response: http.Response{
					StatusCode: 404,
				},
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
