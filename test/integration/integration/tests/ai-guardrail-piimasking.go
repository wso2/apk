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
	IntegrationTests = append(IntegrationTests, AIGuardrailPIIMasking)
}

// AIGuardrailPIIMasking test
var AIGuardrailPIIMasking = suite.IntegrationTest{
	ShortName:   "AIGuardrailPIIMasking",
	Description: "Tests AI Guardrail PII Masking policy with various PII detection and masking scenarios",
	Manifests:   []string{"tests/ai-guardrail-piimasking.yaml"},
	Test: func(t *testing.T, suite *suite.IntegrationTestSuite) {
		ns := "gateway-integration-test-infra"
		gwAddr := "all-http-methods-for-wildcard-piimasking.test.gw.wso2.com:9095"
		token := http.GetTestToken(t)

		testCases := []http.ExpectedResponse{
			// Test case 1: Content with email - should be masked and processed
			{
				TestCaseName: "Content with email - should be masked and processed",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard-piimasking.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard-piimasking/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"content": "Please contact John Doe at john.doe@example.com for further information."}`,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
				Response: http.Response{
					StatusCode: 200,
				},
				//Backend:   "infra-backend-v1-piimasking",
				Namespace: ns,
			},
			// Test case 2: Content with phone number - should be masked and processed
			{
				TestCaseName: "Content with phone number - should be masked and processed",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard-piimasking.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard-piimasking/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"content": "Call me at 555-123-4567 or (555) 987-6543 for assistance."}`,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
				Response: http.Response{
					StatusCode: 200,
				},
				//Backend:   "infra-backend-v1-piimasking",
				Namespace: ns,
			},
			// Test case 3: Content with SSN - should be masked and processed
			{
				TestCaseName: "Content with SSN - should be masked and processed",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard-piimasking.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard-piimasking/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"content": "Social Security Number: 123-45-6789 for verification."}`,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
				Response: http.Response{
					StatusCode: 200,
				},
				//Backend:   "infra-backend-v1-piimasking",
				Namespace: ns,
			},
			// Test case 4: Content with credit card number - should be masked and processed
			{
				TestCaseName: "Content with credit card number - should be masked and processed",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard-piimasking.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard-piimasking/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"content": "Credit card: 4111 1111 1111 1111 expires next month."}`,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
				Response: http.Response{
					StatusCode: 200,
				},
				//Backend:   "infra-backend-v1-piimasking",
				Namespace: ns,
			},
			// Test case 5: Content with multiple PII types - should mask all and process
			{
				TestCaseName: "Content with multiple PII types - should mask all and process",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard-piimasking.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard-piimasking/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"content": "Contact: jane.smith@company.org, Phone: 555-987-6543, SSN: 987-65-4321, Card: 5555-5555-5555-4444"}`,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
				Response: http.Response{
					StatusCode: 200,
				},
				//Backend:   "infra-backend-v1-piimasking",
				Namespace: ns,
			},
			// Test case 6: Content without PII - should pass through unchanged
			{
				TestCaseName: "Content without PII - should pass through unchanged",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard-piimasking.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard-piimasking/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"content": "This is regular content without any sensitive information."}`,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path:   "/v2/echo-full",
						Method: "POST",
						Body:   `{"content": "This is regular content without any sensitive information."}`,
					},
				},
				Response: http.Response{
					StatusCode: 200,
				},
				//Backend:   "infra-backend-v1-piimasking",
				Namespace: ns,
			},
			// Test case 7: Empty content - should pass through
			{
				TestCaseName: "Empty content - should pass through",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard-piimasking.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard-piimasking/v1.0.0/v2/echo-full",
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
				//Backend:   "infra-backend-v1-piimasking",
				Namespace: ns,
			},
			// Test case 8: Content with partial matches - should not mask, pass through
			{
				TestCaseName: "Content with partial matches - should not mask, pass through",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard-piimasking.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard-piimasking/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"content": "Email format like user@domain but incomplete, phone 555-123 incomplete, SSN 123-45 incomplete."}`,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path:   "/v2/echo-full",
						Method: "POST",
						Body:   `{"content": "Email format like user@domain but incomplete, phone 555-123 incomplete, SSN 123-45 incomplete."}`,
					},
				},
				Response: http.Response{
					StatusCode: 200,
				},
				//Backend:   "infra-backend-v1-piimasking",
				Namespace: ns,
			},
			// Test case 9: Content with edge case email formats - should be masked
			{
				TestCaseName: "Content with edge case email formats - should be masked",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard-piimasking.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard-piimasking/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"content": "Emails: test.email+tag@example.co.uk, user_name@subdomain.domain.org"}`,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
				Response: http.Response{
					StatusCode: 200,
				},
				//Backend:   "infra-backend-v1-piimasking",
				Namespace: ns,
			},
			// Test case 10: Content with alternative credit card formats - should be masked
			{
				TestCaseName: "Content with alternative credit card formats - should be masked",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard-piimasking.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard-piimasking/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"content": "Cards: 4111111111111111 (no spaces), 4111-1111-1111-1111 (dashes), 4111 1111 1111 1111 (spaces)"}`,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
				Response: http.Response{
					StatusCode: 200,
				},
				//Backend:   "infra-backend-v1-piimasking",
				Namespace: ns,
			},
			// Test case 11: Missing content field - should fail with JSONPath error
			{
				TestCaseName: "Missing content field",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard-piimasking.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard-piimasking/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"message": "This field is not content with email test@example.com"}`,
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
				//Backend:   "infra-backend-v1-piimasking",
				Namespace: ns,
			},
			// Test case 12: Invalid JSON - should fail
			{
				TestCaseName: "Invalid JSON",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard-piimasking.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard-piimasking/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"content": "incomplete json with email test@example.com`,
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
				//Backend:   "infra-backend-v1-piimasking",
				Namespace: ns,
			},
			// Test case 13: No body request - should be allowed (no validation performed)
			{
				TestCaseName: "No body request",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard-piimasking.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard-piimasking/v1.0.0/v2/echo-full",
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
				//Backend:   "infra-backend-v1-piimasking",
				Namespace: ns,
			},
			// Test case 14: Content with mixed PII and regular text - should mask PII only
			{
				TestCaseName: "Content with mixed PII and regular text - should mask PII only",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard-piimasking.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard-piimasking/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"content": "Hello, this is a normal message. Please contact support@company.com or call 1-800-555-0123. Reference ID: REF123456."}`,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
				Response: http.Response{
					StatusCode: 200,
				},
				//Backend:   "infra-backend-v1-piimasking",
				Namespace: ns,
			},
			// Test case 15: Content with special characters and PII - should mask PII, keep special chars
			{
				TestCaseName: "Content with special characters and PII - should mask PII, keep special chars",
				Request: http.Request{
					Host:   "all-http-methods-for-wildcard-piimasking.test.gw.wso2.com",
					Path:   "/all-http-methods-for-wildcard-piimasking/v1.0.0/v2/echo-full",
					Method: "POST",
					Body:   `{"content": "Special chars: !@#$%^&*() Email: user@test.com Phone: (555) 123-4567 End."}`,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
				Response: http.Response{
					StatusCode: 200,
				},
				//Backend:   "infra-backend-v1-piimasking",
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
