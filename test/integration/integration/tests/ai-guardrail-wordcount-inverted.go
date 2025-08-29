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
	IntegrationTests = append(IntegrationTests, AIGuardrailWordCountInverted)
}

// AIGuardrailWordCountInverted test
var AIGuardrailWordCountInverted = suite.IntegrationTest{
	ShortName:   "AIGuardrailWordCountInverted",
	Description: "Tests AI Guardrail Word Count policy with inverted logic (rejects content within specified range)",
	Manifests:   []string{"tests/ai-guardrail-wordcount-inverted.yaml"},
	Test: func(t *testing.T, suite *suite.IntegrationTestSuite) {
		ns := "gateway-integration-test-infra"
		gwAddr := "ai-guardrail-wordcount-inverted.test.gw.wso2.com:9095"
		token := http.GetTestToken(t)

		testCases := []http.ExpectedResponse{
			// Test case 1: Content with 5 words (outside range 10-20) - should pass (inverted)
			{
				TestCaseName: "Content with 5 words (outside range)",
				Request: http.Request{
					Host:   "ai-guardrail-wordcount-inverted.test.gw.wso2.com",
					Path:   "/ai-guardrail-wordcount-inverted/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"content": "This has five words exactly"}`,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path:   "/v2/echo-full",
						Method: "POST",
						Body:   `{"content": "This has five words exactly"}`,
					},
				},
				Response: http.Response{
					StatusCode: 200,
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
			// Test case 2: Content with 25 words (outside range 10-20) - should pass (inverted)
			{
				TestCaseName: "Content with 25 words (outside range)",
				Request: http.Request{
					Host:   "ai-guardrail-wordcount-inverted.test.gw.wso2.com",
					Path:   "/ai-guardrail-wordcount-inverted/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"content": "This is a test message that contains exactly twenty five words to test the inverted word count guardrail functionality and ensure it works correctly"}`,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path:   "/v2/echo-full",
						Method: "POST",
						Body:   `{"content": "This is a test message that contains exactly twenty five words to test the inverted word count guardrail functionality and ensure it works correctly"}`,
					},
				},
				Response: http.Response{
					StatusCode: 200,
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
			// Test case 3: Content with exactly 10 words (within range 10-20) - should fail (inverted)
			{
				TestCaseName: "Content with exactly 10 words (within range)",
				Request: http.Request{
					Host:   "ai-guardrail-wordcount-inverted.test.gw.wso2.com",
					Path:   "/ai-guardrail-wordcount-inverted/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"content": "This message has exactly ten words including these two dot"}`,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
				Response: http.Response{
					StatusCode: 400,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
			// Test case 4: Content with exactly 15 words (within range 10-20) - should fail (inverted)
			{
				TestCaseName: "Content with exactly 15 words (within range)",
				Request: http.Request{
					Host:   "ai-guardrail-wordcount-inverted.test.gw.wso2.com",
					Path:   "/ai-guardrail-wordcount-inverted/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"content": "This is a test message that has exactly fifteen words to test the inverted guardrail"}`,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
				Response: http.Response{
					StatusCode: 400,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
			// Test case 5: Content with exactly 20 words (within range 10-20) - should fail (inverted)
			{
				TestCaseName: "Content with exactly 20 words (within range)",
				Request: http.Request{
					Host:   "ai-guardrail-wordcount-inverted.test.gw.wso2.com",
					Path:   "/ai-guardrail-wordcount-inverted/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"content": "This is a test message that contains exactly twenty words to test the inverted word count guardrail functionality"}`,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
				Response: http.Response{
					StatusCode: 400,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
			// Test case 6: Content with 1 word (outside range 10-20) - should pass (inverted)
			{
				TestCaseName: "Content with 1 word (outside range)",
				Request: http.Request{
					Host:   "ai-guardrail-wordcount-inverted.test.gw.wso2.com",
					Path:   "/ai-guardrail-wordcount-inverted/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"content": "Hello"}`,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path:   "/v2/echo-full",
						Method: "POST",
						Body:   `{"content": "Hello"}`,
					},
				},
				Response: http.Response{
					StatusCode: 200,
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
			// Test case 7: Empty content (0 words, outside range 10-20) - should pass (inverted)
			{
				TestCaseName: "Empty content (0 words, outside range)",
				Request: http.Request{
					Host:   "ai-guardrail-wordcount-inverted.test.gw.wso2.com",
					Path:   "/ai-guardrail-wordcount-inverted/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"content": ""}`,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path:   "/v2/echo-full",
						Method: "POST",
						Body:   `{"content": ""}`,
					},
				},
				Response: http.Response{
					StatusCode: 200,
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
