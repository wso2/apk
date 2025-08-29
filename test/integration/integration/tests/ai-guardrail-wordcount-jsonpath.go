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
	//IntegrationTests = append(IntegrationTests, AIGuardrailWordCountJSONPath)
}

// AIGuardrailWordCountJSONPath test
var AIGuardrailWordCountJSONPath = suite.IntegrationTest{
	ShortName:   "AIGuardrailWordCountJSONPath",
	Description: "Tests AI Guardrail Word Count policy with different JSONPath configurations",
	Manifests:   []string{"tests/ai-guardrail-wordcount-jsonpath.yaml"},
	Test: func(t *testing.T, suite *suite.IntegrationTestSuite) {
		ns := "gateway-integration-test-infra"
		gwAddr := "ai-guardrail-wordcount-jsonpath.test.gw.wso2.com:9095"
		token := http.GetTestToken(t)

		testCases := []http.ExpectedResponse{
			// Test case 1: Valid nested JSON with correct path ($.message.text) - should pass
			{
				TestCaseName: "Valid nested JSON with correct path",
				Request: http.Request{
					Host:   "ai-guardrail-wordcount-jsonpath.test.gw.wso2.com",
					Path:   "/ai-guardrail-wordcount-jsonpath/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"message": {"text": "This is a test message"}}`,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path:   "/v2/echo-full",
						Method: "POST",
						Body:   `{"message": {"text": "This is a test message"}}`,
					},
				},
				Response: http.Response{
					StatusCode: 200,
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
			// Test case 2: Nested JSON with too many words (11 words, max is 10) - should fail
			{
				TestCaseName: "Nested JSON with too many words",
				Request: http.Request{
					Host:   "ai-guardrail-wordcount-jsonpath.test.gw.wso2.com",
					Path:   "/ai-guardrail-wordcount-jsonpath/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"message": {"text": "This is a test message with too many words exceeding limit"}}`,
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
			// Test case 3: Wrong JSONPath - field exists but at different path - should fail
			{
				TestCaseName: "Wrong JSONPath with existing field",
				Request: http.Request{
					Host:   "ai-guardrail-wordcount-jsonpath.test.gw.wso2.com",
					Path:   "/ai-guardrail-wordcount-jsonpath/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"content": "This field is not at the configured path"}`,
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
			// Test case 4: Missing nested field - should fail with JSONPath error
			{
				TestCaseName: "Missing nested field",
				Request: http.Request{
					Host:   "ai-guardrail-wordcount-jsonpath.test.gw.wso2.com",
					Path:   "/ai-guardrail-wordcount-jsonpath/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"message": {"content": "Wrong field name in nested object"}}`,
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
			// Test case 5: Complex nested JSON with valid path - should pass
			{
				TestCaseName: "Complex nested JSON with valid path",
				Request: http.Request{
					Host:   "ai-guardrail-wordcount-jsonpath.test.gw.wso2.com",
					Path:   "/ai-guardrail-wordcount-jsonpath/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"message": {"text": "Valid message", "metadata": {"id": 123, "timestamp": "2023-01-01"}}, "other": "ignored"}`,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path:   "/v2/echo-full",
						Method: "POST",
						Body:   `{"message": {"text": "Valid message", "metadata": {"id": 123, "timestamp": "2023-01-01"}}, "other": "ignored"}`,
					},
				},
				Response: http.Response{
					StatusCode: 200,
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
			// Test case 6: Empty text in nested field - should fail (0 words, min is 1)
			{
				TestCaseName: "Empty text in nested field",
				Request: http.Request{
					Host:   "ai-guardrail-wordcount-jsonpath.test.gw.wso2.com",
					Path:   "/ai-guardrail-wordcount-jsonpath/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"message": {"text": ""}}`,
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
			// Test case 7: Exactly at boundary (10 words) - should pass
			{
				TestCaseName: "Exactly at boundary (10 words)",
				Request: http.Request{
					Host:   "ai-guardrail-wordcount-jsonpath.test.gw.wso2.com",
					Path:   "/ai-guardrail-wordcount-jsonpath/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"message": {"text": "This message has exactly ten words to test boundary conditions"}}`,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path:   "/v2/echo-full",
						Method: "POST",
						Body:   `{"message": {"text": "This message has exactly ten words to test boundary conditions"}}`,
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
