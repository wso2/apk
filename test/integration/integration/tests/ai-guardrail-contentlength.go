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
	//IntegrationTests = append(IntegrationTests, AIGuardrailContentLength)
}

// AIGuardrailContentLength test
var AIGuardrailContentLength = suite.IntegrationTest{
	ShortName:   "AIGuardrailContentLength",
	Description: "Tests AI Guardrail Content Length policy with various content length scenarios",
	Manifests:   []string{"tests/ai-guardrail-contentlength.yaml"},
	Test: func(t *testing.T, suite *suite.IntegrationTestSuite) {
		ns := "gateway-integration-test-infra"
		gwAddr := "all-http-methods-for-wildcard-contentlength.test.gw.wso2.com:9095"
		token := http.GetTestToken(t)

		testCases := []http.ExpectedResponse{
			// Test case 1: Valid content length within range (10-1000 bytes) - should pass
			{
				TestCaseName: "Valid content length within range",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard-contentlength.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard-contentlength/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"content": "This is a test message with sufficient length to pass the content length validation."}`,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path:   "/v2/echo-full",
						Method: "POST",
						Body:   `{"content": "This is a test message with sufficient length to pass the content length validation."}`,
					},
				},
				Response: http.Response{
					StatusCode: 200,
				},
				//Backend:   "infra-backend-v1-contentlength",
				Namespace: ns,
			},
			// Test case 2: Empty content (0 bytes) - should fail (below minimum)
			{
				TestCaseName: "Empty content (0 bytes)",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard-contentlength.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard-contentlength/v1.0.0/v2/echo-full",
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
				//Backend:   "infra-backend-v1-contentlength",
				Namespace: ns,
			},
			// Test case 3: Content with exactly 10 bytes - should pass (minimum boundary)
			{
				TestCaseName: "Content with exactly 10 bytes",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard-contentlength.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard-contentlength/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"content": "1234567890"}`,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path:   "/v2/echo-full",
						Method: "POST",
						Body:   `{"content": "1234567890"}`,
					},
				},
				Response: http.Response{
					StatusCode: 200,
				},
				//Backend:   "infra-backend-v1-contentlength",
				Namespace: ns,
			},
			// Test case 4: Content with 9 bytes - should fail (below minimum)
			{
				TestCaseName: "Content with 9 bytes (below minimum)",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard-contentlength.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard-contentlength/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"content": "123456789"}`,
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
				//Backend:   "infra-backend-v1-contentlength",
				Namespace: ns,
			},
			// Test case 5: Content with exactly 1000 bytes - should pass (maximum boundary)
			{
				TestCaseName: "Content with exactly 1000 bytes",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard-contentlength.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard-contentlength/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"content": "` + generateStringOfLength(1000) + `"}`,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path:   "/v2/echo-full",
						Method: "POST",
						Body:   `{"content": "` + generateStringOfLength(1000) + `"}`,
					},
				},
				Response: http.Response{
					StatusCode: 200,
				},
				//Backend:   "infra-backend-v1-contentlength",
				Namespace: ns,
			},
			// Test case 6: Content with 1001 bytes - should fail (above maximum)
			{
				TestCaseName: "Content with 1001 bytes (above maximum)",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard-contentlength.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard-contentlength/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"content": "` + generateStringOfLength(1001) + `"}`,
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
				//Backend:   "infra-backend-v1-contentlength",
				Namespace: ns,
			},
			// Test case 7: Content with Unicode characters - should count bytes correctly
			{
				TestCaseName: "Content with Unicode characters",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard-contentlength.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard-contentlength/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"content": "Hello 世界! This content has Unicode characters. 你好世界! 🌍🌎🌏"}`,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path:   "/v2/echo-full",
						Method: "POST",
						Body:   `{"content": "Hello 世界! This content has Unicode characters. 你好世界! 🌍🌎🌏"}`,
					},
				},
				Response: http.Response{
					StatusCode: 200,
				},
				//Backend:   "infra-backend-v1-contentlength",
				Namespace: ns,
			},
			// Test case 8: Content with special characters and punctuation
			{
				TestCaseName: "Content with special characters and punctuation",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard-contentlength.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard-contentlength/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"content": "Special chars: !@#$%^&*()_+-=[]{}|;':\",./<>? and numbers 123456789"}`,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path:   "/v2/echo-full",
						Method: "POST",
						Body:   `{"content": "Special chars: !@#$%^&*()_+-=[]{}|;':\",./<>? and numbers 123456789"}`,
					},
				},
				Response: http.Response{
					StatusCode: 200,
				},
				//Backend:   "infra-backend-v1-contentlength",
				Namespace: ns,
			},
			// Test case 9: Content with whitespace and newlines - should handle correctly
			{
				TestCaseName: "Content with whitespace and newlines",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard-contentlength.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard-contentlength/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"content": "Line one\nLine two\tTabbed text\nThird line with sufficient content for validation"}`,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path:   "/v2/echo-full",
						Method: "POST",
						Body:   `{"content": "Line one\nLine two\tTabbed text\nThird line with sufficient content for validation"}`,
					},
				},
				Response: http.Response{
					StatusCode: 200,
				},
				//Backend:   "infra-backend-v1-contentlength",
				Namespace: ns,
			},
			// Test case 10: Missing content field - should fail with JSONPath error
			{
				TestCaseName: "Missing content field",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard-contentlength.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard-contentlength/v1.0.0/v2/echo-full",
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
				//Backend:   "infra-backend-v1-contentlength",
				Namespace: ns,
			},
			// Test case 11: Invalid JSON - should fail
			{
				TestCaseName: "Invalid JSON",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard-contentlength.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard-contentlength/v1.0.0/v2/echo-full",
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
				//Backend:   "infra-backend-v1-contentlength",
				Namespace: ns,
			},
			// Test case 12: No body request - should be allowed (no validation performed)
			{
				TestCaseName: "No body request",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard-contentlength.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard-contentlength/v1.0.0/v2/echo-full",
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
				//Backend:   "infra-backend-v1-contentlength",
				Namespace: ns,
			},
			// Test case 13: Whitespace-only content - should fail (0 bytes after cleaning)
			{
				TestCaseName: "Whitespace-only content",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard-contentlength.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard-contentlength/v1.0.0/v2/echo-full",
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
				//Backend:   "infra-backend-v1-contentlength",
				Namespace: ns,
			},
			// Test case 14: Content with HTML tags - should count after cleaning
			{
				TestCaseName: "Content with HTML tags",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard-contentlength.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard-contentlength/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"content": "<p>This is HTML content with tags.</p><br><strong>Bold text</strong> and <em>italic text</em>."}`,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path:   "/v2/echo-full",
						Method: "POST",
						Body:   `{"content": "<p>This is HTML content with tags.</p><br><strong>Bold text</strong> and <em>italic text</em>."}`,
					},
				},
				Response: http.Response{
					StatusCode: 200,
				},
				//Backend:   "infra-backend-v1-contentlength",
				Namespace: ns,
			},
			// Test case 15: Content with escape sequences
			{
				TestCaseName: "Content with escape sequences",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard-contentlength.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard-contentlength/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"content": "Content with escape sequences: \\n \\t \\r \\\" \\\\ and more text to meet minimum length"}`,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path:   "/v2/echo-full",
						Method: "POST",
						Body:   `{"content": "Content with escape sequences: \\n \\t \\r \\\" \\\\ and more text to meet minimum length"}`,
					},
				},
				Response: http.Response{
					StatusCode: 200,
				},
				//Backend:   "infra-backend-v1-contentlength",
				Namespace: ns,
			},
			// Test case 16: Mixed content with medium length (500 bytes)
			{
				TestCaseName: "Mixed content with medium length",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard-contentlength.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard-contentlength/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"content": "` + generateMixedContent(500) + `"}`,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path:   "/v2/echo-full",
						Method: "POST",
						Body:   `{"content": "` + generateMixedContent(500) + `"}`,
					},
				},
				Response: http.Response{
					StatusCode: 200,
				},
				//Backend:   "infra-backend-v1-contentlength",
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

// generateStringOfLength generates a string of the specified byte length
func generateStringOfLength(length int) string {
	if length <= 0 {
		return ""
	}

	// Use a mix of characters to ensure proper byte counting
	pattern := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)

	for i := 0; i < length; i++ {
		result[i] = pattern[i%len(pattern)]
	}

	return string(result)
}

// generateMixedContent generates mixed content with various character types
func generateMixedContent(targetLength int) string {
	base := "This is mixed content with numbers 123, special chars !@#$%, and unicode 世界🌍. "
	content := ""

	for len(content) < targetLength {
		content += base
	}

	// Trim to exact length if needed
	if len(content) > targetLength {
		content = content[:targetLength]
	}

	return content
}
