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
	"time"

	"github.com/wso2/apk/test/integration/integration/utils/http"
	"github.com/wso2/apk/test/integration/integration/utils/suite"
)

func init() {
	IntegrationTests = append(IntegrationTests, ExternalCustomMediation)
}

// ExternalCustomMediation test
var ExternalCustomMediation = suite.IntegrationTest{
	ShortName:   "ExternalCustomMediation",
	Description: "Tests External Custom Mediation policy with various scenarios",
	Manifests:   []string{"tests/external-custom-mediation.yaml"},
	Test: func(t *testing.T, suite *suite.IntegrationTestSuite) {
		ns := "gateway-integration-test-infra"
		gwAddr := "all-http-methods-for-wildcard.test.gw.wso2.com:9095"
		token := http.GetTestToken(t)

		testCases := []http.ExpectedResponse{
			// Test case 1: Valid word count within range (1-100 words) - should pass
			{
				TestCaseName: "Valid word count within range",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"content": "This is a test message with exactly ten words here."}`,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
				// ExpectedRequest: &http.ExpectedRequest{
				// 	Request: http.Request{
				// 		Path:   "/v2/echo-full",
				// 		Method: "POST",
				// 		Body:   `{"content": "This is a test message with exactly ten words here."}`,
				// 	},
				// },
				Response: http.Response{
					StatusCode: 211,
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
		}
		// Wait for 10 seconds to ensure resources are ready
		time.Sleep(10 * time.Second)
		for i := range testCases {
			tc := testCases[i]
			tc.Request.Headers = http.AddBearerTokenToHeader(token, tc.Request.Headers)
			tc.Request.Headers = http.AddBearerTokenToHeader(token, tc.Request.Headers)

			t.Run(tc.GetTestCaseName(i), func(t *testing.T) {
				t.Parallel()
				http.MakeRequestAndExpectEventuallyConsistentResponse(t, suite.RoundTripper, suite.TimeoutConfig, gwAddr, tc)
			})
		}
	},
}
