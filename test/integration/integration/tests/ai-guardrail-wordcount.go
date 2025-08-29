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
	//IntegrationTests = append(IntegrationTests, AIGuardrailWordCount)
}

// AIGuardrailWordCount test
var AIGuardrailWordCount = suite.IntegrationTest{
	ShortName:   "AIGuardrailWordCount",
	Description: "Tests AI Guardrail Word Count policy with various word count scenarios",
	Manifests:   []string{"tests/ai-guardrail-wordcount.yaml"},
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
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path:   "/v2/echo-full",
						Method: "POST",
						Body:   `{"content": "This is a test message with exactly ten words here."}`,
					},
				},
				Response: http.Response{
					StatusCode: 200,
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
			// Test case 2: Empty content (0 words) - should fail (below minimum)
			{
				TestCaseName: "Empty content (0 words)",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"content": ""}`,
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
			// Test case 3: Single word (1 word) - should pass (minimum boundary)
			{
				TestCaseName: "Single word (1 word)",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard/v1.0.0/v2/echo-full",
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
			// Test case 4: Exactly 100 words - should pass (maximum boundary)
			{
				TestCaseName: "Exactly 100 words",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"content": "word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word "}`,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path:   "/v2/echo-full",
						Method: "POST",
						Body:   `{"content": "word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word "}`,
					},
				},
				Response: http.Response{
					StatusCode: 200,
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
			// Test case 5: 101 words - should fail (above maximum)
			{
				TestCaseName: "101 words (above maximum)",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"content": "word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word word extra"}`,
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
			// Test case 6: Content with punctuation and special characters - should clean and count correctly
			{
				TestCaseName: "Content with punctuation and special characters",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"content": "Hello, world! This is a test... with punctuation & special characters!!! Should count as twelve words."}`,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path:   "/v2/echo-full",
						Method: "POST",
						Body:   `{"content": "Hello, world! This is a test... with punctuation & special characters!!! Should count as twelve words."}`,
					},
				},
				Response: http.Response{
					StatusCode: 200,
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
			// Test case 7: Content with extra whitespaces - should handle correctly
			{
				TestCaseName: "Content with extra whitespaces",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"content": "   Multiple    spaces     between    words    should   be   handled   correctly   "}`,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path:   "/v2/echo-full",
						Method: "POST",
						Body:   `{"content": "   Multiple    spaces     between    words    should   be   handled   correctly   "}`,
					},
				},
				Response: http.Response{
					StatusCode: 200,
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
			// Test case 8: Missing content field - should fail with JSONPath error
			{
				TestCaseName: "Missing content field",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"message": "This field is not content"}`,
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
			// Test case 9: Invalid JSON - should fail
			{
				TestCaseName: "Invalid JSON",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"content": "incomplete json`,
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
			// Test case 10: No body request - should be allowed (no validation performed)
			{
				TestCaseName: "No body request",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard/v1.0.0/v2/echo-full",
					Method: "GET",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path:   "/v2/echo-full",
						Method: "GET",
					},
				},
				Response: http.Response{
					StatusCode: 200,
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
			// Test case 11: Whitespace-only content - should fail (0 words)
			{
				TestCaseName: "Whitespace-only content",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"content": "   \t   \n   "}`,
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
			// Test case 12: Content with mixed case and numbers - should count correctly
			{
				TestCaseName: "Content with mixed case and numbers",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"content": "Testing 123 Mixed Case WORDS with Numbers 456 and Symbols"}`,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path:   "/v2/echo-full",
						Method: "POST",
						Body:   `{"content": "Testing 123 Mixed Case WORDS with Numbers 456 and Symbols"}`,
					},
				},
				Response: http.Response{
					StatusCode: 200,
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
			// Test case 13: Content with newlines and tabs - should handle correctly
			{
				TestCaseName: "Content with newlines and tabs",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"content": "Line one\nLine two\tTabbed text\nThird line with words"}`,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path:   "/v2/echo-full",
						Method: "POST",
						Body:   `{"content": "Line one\nLine two\tTabbed text\nThird line with words"}`,
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
