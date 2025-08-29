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
	//IntegrationTests = append(IntegrationTests, AIGuardrailSentenceCount)
}

// AIGuardrailSentenceCount test
var AIGuardrailSentenceCount = suite.IntegrationTest{
	ShortName:   "AIGuardrailSentenceCount",
	Description: "Tests AI Guardrail Sentence Count policy with various sentence count scenarios",
	Manifests:   []string{"tests/ai-guardrail-sentencecount.yaml"},
	Test: func(t *testing.T, suite *suite.IntegrationTestSuite) {
		ns := "gateway-integration-test-infra"
		gwAddr := "all-http-methods-for-wildcard-sentencecount.test.gw.wso2.com:9095"
		token := http.GetTestToken(t)

		testCases := []http.ExpectedResponse{
			// Test case 1: Valid sentence count within range (1-10 sentences) - should pass
			{
				TestCaseName: "Valid sentence count within range",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard-sentencecount.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard-sentencecount/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"content": "This is the first sentence. This is the second sentence. This is the third sentence."}`,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path:   "/v2/echo-full",
						Method: "POST",
						Body:   `{"content": "This is the first sentence. This is the second sentence. This is the third sentence."}`,
					},
				},
				Response: http.Response{
					StatusCode: 200,
				},
				// Backend:   "infra-backend-v1-sentencecount",
				Namespace: ns,
			},
			// Test case 2: Empty content (0 sentences) - should fail (below minimum)
			{
				TestCaseName: "Empty content (0 sentences)",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard-sentencecount.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard-sentencecount/v1.0.0/v2/echo-full",
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
				// Backend:   "infra-backend-v1-sentencecount",
				Namespace: ns,
			},
			// Test case 3: Single sentence (1 sentence) - should pass (minimum boundary)
			{
				TestCaseName: "Single sentence (1 sentence)",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard-sentencecount.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard-sentencecount/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"content": "This is a single sentence."}`,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path:   "/v2/echo-full",
						Method: "POST",
						Body:   `{"content": "This is a single sentence."}`,
					},
				},
				Response: http.Response{
					StatusCode: 200,
				},
				// Backend:   "infra-backend-v1-sentencecount",
				Namespace: ns,
			},
			// Test case 4: Exactly 10 sentences - should pass (maximum boundary)
			{
				TestCaseName: "Exactly 10 sentences",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard-sentencecount.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard-sentencecount/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"content": "First sentence. Second sentence. Third sentence. Fourth sentence. Fifth sentence. Sixth sentence. Seventh sentence. Eighth sentence. Ninth sentence. Tenth sentence."}`,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path:   "/v2/echo-full",
						Method: "POST",
						Body:   `{"content": "First sentence. Second sentence. Third sentence. Fourth sentence. Fifth sentence. Sixth sentence. Seventh sentence. Eighth sentence. Ninth sentence. Tenth sentence."}`,
					},
				},
				Response: http.Response{
					StatusCode: 200,
				},
				// Backend:   "infra-backend-v1-sentencecount",
				Namespace: ns,
			},
			// Test case 5: 11 sentences - should fail (above maximum)
			{
				TestCaseName: "11 sentences (above maximum)",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard-sentencecount.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard-sentencecount/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"content": "First sentence. Second sentence. Third sentence. Fourth sentence. Fifth sentence. Sixth sentence. Seventh sentence. Eighth sentence. Ninth sentence. Tenth sentence. Eleventh sentence."}`,
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
				// Backend:   "infra-backend-v1-sentencecount",
				Namespace: ns,
			},
			// Test case 6: Content with mixed punctuation - should count correctly
			{
				TestCaseName: "Content with mixed punctuation",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard-sentencecount.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard-sentencecount/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"content": "This is a question? This is an exclamation! This is a regular sentence. Another question? Another exclamation!"}`,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path:   "/v2/echo-full",
						Method: "POST",
						Body:   `{"content": "This is a question? This is an exclamation! This is a regular sentence. Another question? Another exclamation!"}`,
					},
				},
				Response: http.Response{
					StatusCode: 200,
				},
				// Backend:   "infra-backend-v1-sentencecount",
				Namespace: ns,
			},
			// Test case 7: Content with ellipsis and multiple punctuation - should handle correctly
			{
				TestCaseName: "Content with ellipsis and multiple punctuation",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard-sentencecount.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard-sentencecount/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"content": "First sentence... Second sentence!!! Third sentence??? Fourth sentence with mixed... punctuation!!!"}`,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path:   "/v2/echo-full",
						Method: "POST",
						Body:   `{"content": "First sentence... Second sentence!!! Third sentence??? Fourth sentence with mixed... punctuation!!!"}`,
					},
				},
				Response: http.Response{
					StatusCode: 200,
				},
				// Backend:   "infra-backend-v1-sentencecount",
				Namespace: ns,
			},
			// Test case 8: Missing content field - should fail with JSONPath error
			{
				TestCaseName: "Missing content field",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard-sentencecount.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard-sentencecount/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"message": "This field is not content."}`,
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
				// Backend:   "infra-backend-v1-sentencecount",
				Namespace: ns,
			},
			// Test case 9: Invalid JSON - should fail
			{
				TestCaseName: "Invalid JSON",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard-sentencecount.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard-sentencecount/v1.0.0/v2/echo-full",
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
				// Backend:   "infra-backend-v1-sentencecount",
				Namespace: ns,
			},
			// Test case 10: No body request - should be allowed (no validation performed)
			{
				TestCaseName: "No body request",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard-sentencecount.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard-sentencecount/v1.0.0/v2/echo-full",
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
				// Backend:   "infra-backend-v1-sentencecount",
				Namespace: ns,
			},
			// Test case 11: Whitespace-only content - should fail (0 sentences)
			{
				TestCaseName: "Whitespace-only content",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard-sentencecount.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard-sentencecount/v1.0.0/v2/echo-full",
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
				// Backend:   "infra-backend-v1-sentencecount",
				Namespace: ns,
			},
			// Test case 12: Content with newlines and tabs between sentences - should handle correctly
			{
				TestCaseName: "Content with newlines and tabs between sentences",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard-sentencecount.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard-sentencecount/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"content": "First sentence.\nSecond sentence with newline.\tThird sentence with tab. Fourth sentence."}`,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path:   "/v2/echo-full",
						Method: "POST",
						Body:   `{"content": "First sentence.\nSecond sentence with newline.\tThird sentence with tab. Fourth sentence."}`,
					},
				},
				Response: http.Response{
					StatusCode: 200,
				},
				// Backend:   "infra-backend-v1-sentencecount",
				Namespace: ns,
			},
			// Test case 13: Content with incomplete sentences (no ending punctuation) - should count correctly
			{
				TestCaseName: "Content with incomplete sentences",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard-sentencecount.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard-sentencecount/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"content": "Complete sentence. Incomplete sentence without punctuation Another complete sentence."}`,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path:   "/v2/echo-full",
						Method: "POST",
						Body:   `{"content": "Complete sentence. Incomplete sentence without punctuation Another complete sentence."}`,
					},
				},
				Response: http.Response{
					StatusCode: 200,
				},
				// Backend:   "infra-backend-v1-sentencecount",
				Namespace: ns,
			},
			// Test case 14: Content with only punctuation marks - should fail (0 sentences)
			{
				TestCaseName: "Content with only punctuation marks",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard-sentencecount.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard-sentencecount/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"content": "...!!!???"}`,
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
				// Backend:   "infra-backend-v1-sentencecount",
				Namespace: ns,
			},
			// Test case 15: Content with abbreviations and decimals - should handle correctly
			{
				TestCaseName: "Content with abbreviations and decimals",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard-sentencecount.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard-sentencecount/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"content": "Mr. Smith visited the U.S.A. in 2023. The temperature was 98.6 degrees. He paid $1,234.56 for the trip."}`,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path:   "/v2/echo-full",
						Method: "POST",
						Body:   `{"content": "Mr. Smith visited the U.S.A. in 2023. The temperature was 98.6 degrees. He paid $1,234.56 for the trip."}`,
					},
				},
				Response: http.Response{
					StatusCode: 200,
				},
				// Backend:   "infra-backend-v1-sentencecount",
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
