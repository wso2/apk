/*
 *  Copyright (c) 2025, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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

package mediation

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"testing"

	egv1a1 "github.com/envoyproxy/gateway/api/v1alpha1"
	envoy_service_proc_v3 "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"
	dpv2alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v2alpha1"
	"github.com/wso2/apk/gateway/enforcer/internal/logging"
	"github.com/wso2/apk/gateway/enforcer/internal/requestconfig"
)

// Mock logger for testing
func createMockLoggerForRegex() *logging.Logger {
	// Create a proper mock logger using the default logging setup
	mockLogger := logging.DefaultLogger(egv1a1.LogLevelInfo)
	return &mockLogger
}

// Helper function to create test mediation for RegexGuardrail
func createTestRegexMediation(params map[string]string) *dpv2alpha1.Mediation {
	var parameters []*dpv2alpha1.Parameter
	for key, value := range params {
		parameters = append(parameters, &dpv2alpha1.Parameter{
			Key:   key,
			Value: value,
		})
	}

	return &dpv2alpha1.Mediation{
		PolicyName:    "RegexGuardrail",
		PolicyVersion: "v1",
		PolicyID:      "test-regex-policy-id",
		Parameters:    parameters,
	}
}

// Helper function to create gzipped content for regex tests
func createGzippedContentForRegex(content string) []byte {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	gz.Write([]byte(content))
	gz.Close()
	return buf.Bytes()
}

func TestNewRegexGuardrail(t *testing.T) {
	tests := []struct {
		name       string
		parameters map[string]string
		expected   RegexGuardrail
	}{
		{
			name:       "Default values",
			parameters: map[string]string{},
			expected: RegexGuardrail{
				PolicyName:     "RegexGuardrail",
				PolicyVersion:  "v1",
				PolicyID:       "test-regex-policy-id",
				Name:           "RegexGuardrail",
				Regex:          "",
				JSONPath:       "$.content",
				Inverted:       false,
				ShowAssessment: false,
			},
		},
		{
			name: "Custom values",
			parameters: map[string]string{
				"name":           "CustomRegexGuardrail",
				"regex":          "^[a-zA-Z0-9]+$",
				"jsonPath":       "$.message",
				"invert":         "true",
				"showAssessment": "true",
			},
			expected: RegexGuardrail{
				PolicyName:     "RegexGuardrail",
				PolicyVersion:  "v1",
				PolicyID:       "test-regex-policy-id",
				Name:           "CustomRegexGuardrail",
				Regex:          "^[a-zA-Z0-9]+$",
				JSONPath:       "$.message",
				Inverted:       true,
				ShowAssessment: true,
			},
		},
		{
			name: "Email regex pattern",
			parameters: map[string]string{
				"name":     "EmailValidator",
				"regex":    "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$",
				"jsonPath": "$.email",
			},
			expected: RegexGuardrail{
				PolicyName:     "RegexGuardrail",
				PolicyVersion:  "v1",
				PolicyID:       "test-regex-policy-id",
				Name:           "EmailValidator",
				Regex:          "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$",
				JSONPath:       "$.email",
				Inverted:       false,
				ShowAssessment: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mediation := createTestRegexMediation(tt.parameters)
			result := NewRegexGuardrail(mediation)

			if result.PolicyName != tt.expected.PolicyName {
				t.Errorf("Expected PolicyName %s, got %s", tt.expected.PolicyName, result.PolicyName)
			}
			if result.Name != tt.expected.Name {
				t.Errorf("Expected Name %s, got %s", tt.expected.Name, result.Name)
			}
			if result.Regex != tt.expected.Regex {
				t.Errorf("Expected Regex %s, got %s", tt.expected.Regex, result.Regex)
			}
			if result.JSONPath != tt.expected.JSONPath {
				t.Errorf("Expected JSONPath %s, got %s", tt.expected.JSONPath, result.JSONPath)
			}
			if result.Inverted != tt.expected.Inverted {
				t.Errorf("Expected Inverted %t, got %t", tt.expected.Inverted, result.Inverted)
			}
			if result.ShowAssessment != tt.expected.ShowAssessment {
				t.Errorf("Expected ShowAssessment %t, got %t", tt.expected.ShowAssessment, result.ShowAssessment)
			}
		})
	}
}

func TestRegexGuardrail_Process(t *testing.T) {
	tests := []struct {
		name              string
		guardrail         *RegexGuardrail
		requestConfig     *requestconfig.Holder
		expectedImmediate bool
		expectedPassed    bool
	}{
		{
			name: "Request body - regex matches (normal mode)",
			guardrail: &RegexGuardrail{
				Name:     "TestRegexGuardrail",
				Regex:    "^hello.*",
				JSONPath: "$.content",
				Inverted: false,
				logger:   createMockLoggerForRegex(),
			},
			requestConfig: &requestconfig.Holder{
				ProcessingPhase: requestconfig.ProcessingPhaseRequestBody,
				RequestBody: &envoy_service_proc_v3.HttpBody{
					Body: []byte(`{"content": "hello world"}`),
				},
			},
			expectedImmediate: false,
			expectedPassed:    true,
		},
		{
			name: "Request body - regex doesn't match (normal mode)",
			guardrail: &RegexGuardrail{
				Name:     "TestRegexGuardrail",
				Regex:    "^hello.*",
				JSONPath: "$.content",
				Inverted: false,
				logger:   createMockLoggerForRegex(),
			},
			requestConfig: &requestconfig.Holder{
				ProcessingPhase: requestconfig.ProcessingPhaseRequestBody,
				RequestBody: &envoy_service_proc_v3.HttpBody{
					Body: []byte(`{"content": "goodbye world"}`),
				},
			},
			expectedImmediate: true,
			expectedPassed:    false,
		},
		{
			name: "Request body - regex matches (inverted mode)",
			guardrail: &RegexGuardrail{
				Name:     "TestRegexGuardrail",
				Regex:    "^hello.*",
				JSONPath: "$.content",
				Inverted: true,
				logger:   createMockLoggerForRegex(),
			},
			requestConfig: &requestconfig.Holder{
				ProcessingPhase: requestconfig.ProcessingPhaseRequestBody,
				RequestBody: &envoy_service_proc_v3.HttpBody{
					Body: []byte(`{"content": "hello world"}`),
				},
			},
			expectedImmediate: true,
			expectedPassed:    false,
		},
		{
			name: "Request body - regex doesn't match (inverted mode)",
			guardrail: &RegexGuardrail{
				Name:     "TestRegexGuardrail",
				Regex:    "^hello.*",
				JSONPath: "$.content",
				Inverted: true,
				logger:   createMockLoggerForRegex(),
			},
			requestConfig: &requestconfig.Holder{
				ProcessingPhase: requestconfig.ProcessingPhaseRequestBody,
				RequestBody: &envoy_service_proc_v3.HttpBody{
					Body: []byte(`{"content": "goodbye world"}`),
				},
			},
			expectedImmediate: false,
			expectedPassed:    true,
		},
		{
			name: "Response body - regex matches",
			guardrail: &RegexGuardrail{
				Name:     "TestRegexGuardrail",
				Regex:    ".*success.*",
				JSONPath: "$.status",
				Inverted: false,
				logger:   createMockLoggerForRegex(),
			},
			requestConfig: &requestconfig.Holder{
				ProcessingPhase: requestconfig.ProcessingPhaseResponseBody,
				ResponseBody: &envoy_service_proc_v3.HttpBody{
					Body: []byte(`{"status": "operation successful"}`),
				},
			},
			expectedImmediate: false,
			expectedPassed:    true,
		},
		{
			name: "Response body - gzipped content",
			guardrail: &RegexGuardrail{
				Name:     "TestRegexGuardrail",
				Regex:    ".*success.*",
				JSONPath: "$.status",
				Inverted: false,
				logger:   createMockLoggerForRegex(),
			},
			requestConfig: &requestconfig.Holder{
				ProcessingPhase: requestconfig.ProcessingPhaseResponseBody,
				ResponseBody: &envoy_service_proc_v3.HttpBody{
					Body: createGzippedContentForRegex(`{"status": "operation successful"}`),
				},
			},
			expectedImmediate: false,
			expectedPassed:    true,
		},
		{
			name: "Email validation - valid email",
			guardrail: &RegexGuardrail{
				Name:     "EmailValidator",
				Regex:    "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$",
				JSONPath: "$.email",
				Inverted: false,
				logger:   createMockLoggerForRegex(),
			},
			requestConfig: &requestconfig.Holder{
				ProcessingPhase: requestconfig.ProcessingPhaseRequestBody,
				RequestBody: &envoy_service_proc_v3.HttpBody{
					Body: []byte(`{"email": "user@example.com"}`),
				},
			},
			expectedImmediate: false,
			expectedPassed:    true,
		},
		{
			name: "Email validation - invalid email",
			guardrail: &RegexGuardrail{
				Name:     "EmailValidator",
				Regex:    "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$",
				JSONPath: "$.email",
				Inverted: false,
				logger:   createMockLoggerForRegex(),
			},
			requestConfig: &requestconfig.Holder{
				ProcessingPhase: requestconfig.ProcessingPhaseRequestBody,
				RequestBody: &envoy_service_proc_v3.HttpBody{
					Body: []byte(`{"email": "invalid-email"}`),
				},
			},
			expectedImmediate: true,
			expectedPassed:    false,
		},
		{
			name: "No request body",
			guardrail: &RegexGuardrail{
				Name:     "TestRegexGuardrail",
				Regex:    ".*",
				JSONPath: "$.content",
				logger:   createMockLoggerForRegex(),
			},
			requestConfig: &requestconfig.Holder{
				ProcessingPhase: requestconfig.ProcessingPhaseRequestBody,
				RequestBody:     nil,
			},
			expectedImmediate: false,
			expectedPassed:    true,
		},
		{
			name: "Different processing phase",
			guardrail: &RegexGuardrail{
				Name:     "TestRegexGuardrail",
				Regex:    ".*",
				JSONPath: "$.content",
				logger:   createMockLoggerForRegex(),
			},
			requestConfig: &requestconfig.Holder{
				ProcessingPhase: requestconfig.ProcessingPhaseRequestHeaders,
			},
			expectedImmediate: false,
			expectedPassed:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.guardrail.Process(tt.requestConfig)

			if result.ImmediateResponse != tt.expectedImmediate {
				t.Errorf("Expected ImmediateResponse %t, got %t", tt.expectedImmediate, result.ImmediateResponse)
			}

			// If we expect immediate response, check that the response contains error information
			if tt.expectedImmediate {
				if result.ImmediateResponseBody == "" {
					t.Error("Expected non-empty ImmediateResponseBody for error case")
				}
				if result.ImmediateResponseContentType != "application/json" {
					t.Errorf("Expected ContentType 'application/json', got '%s'", result.ImmediateResponseContentType)
				}

				// Parse the response body to check error structure
				var responseBody map[string]interface{}
				err := json.Unmarshal([]byte(result.ImmediateResponseBody), &responseBody)
				if err != nil {
					t.Errorf("Failed to parse response body JSON: %v", err)
				}

				if responseBody[RegexErrorCode] != RegexGuardrailAPIMExceptionCode {
					t.Errorf("Expected error code %s, got %v", RegexGuardrailAPIMExceptionCode, responseBody[RegexErrorCode])
				}

				if responseBody[RegexErrorType] != RegexGuardrailConstant {
					t.Errorf("Expected error type %s, got %v", RegexGuardrailConstant, responseBody[RegexErrorType])
				}
			}
		})
	}
}

func TestRegexGuardrail_validatePayload(t *testing.T) {
	tests := []struct {
		name          string
		guardrail     *RegexGuardrail
		payload       []byte
		isResponse    bool
		expectedValid bool
		expectedError bool
	}{
		{
			name: "Valid alphanumeric string",
			guardrail: &RegexGuardrail{
				Regex:    "^[a-zA-Z0-9]+$",
				JSONPath: "$.content",
				Inverted: false,
				logger:   createMockLoggerForRegex(),
			},
			payload:       []byte(`{"content": "hello123"}`),
			isResponse:    false,
			expectedValid: true,
			expectedError: false,
		},
		{
			name: "Invalid alphanumeric string (contains special chars)",
			guardrail: &RegexGuardrail{
				Regex:    "^[a-zA-Z0-9]+$",
				JSONPath: "$.content",
				Inverted: false,
				logger:   createMockLoggerForRegex(),
			},
			payload:       []byte(`{"content": "hello@world!"}`),
			isResponse:    false,
			expectedValid: false,
			expectedError: false,
		},
		{
			name: "Inverted validation - matches but inverted",
			guardrail: &RegexGuardrail{
				Regex:    "^[a-zA-Z0-9]+$",
				JSONPath: "$.content",
				Inverted: true,
				logger:   createMockLoggerForRegex(),
			},
			payload:       []byte(`{"content": "hello123"}`),
			isResponse:    false,
			expectedValid: false,
			expectedError: false,
		},
		{
			name: "Inverted validation - doesn't match and inverted",
			guardrail: &RegexGuardrail{
				Regex:    "^[a-zA-Z0-9]+$",
				JSONPath: "$.content",
				Inverted: true,
				logger:   createMockLoggerForRegex(),
			},
			payload:       []byte(`{"content": "hello@world!"}`),
			isResponse:    false,
			expectedValid: true,
			expectedError: false,
		},
		{
			name: "Invalid JSON path",
			guardrail: &RegexGuardrail{
				Regex:    ".*",
				JSONPath: "$.nonexistent",
				Inverted: false,
				logger:   createMockLoggerForRegex(),
			},
			payload:       []byte(`{"content": "hello"}`),
			isResponse:    false,
			expectedValid: true, // Empty string matches ".*" regex
			expectedError: false,
		},
		{
			name: "Invalid regex pattern",
			guardrail: &RegexGuardrail{
				Regex:    "[",
				JSONPath: "$.content",
				Inverted: false,
				logger:   createMockLoggerForRegex(),
			},
			payload:       []byte(`{"content": "hello"}`),
			isResponse:    false,
			expectedValid: false,
			expectedError: true,
		},
		{
			name: "Response with gzipped content",
			guardrail: &RegexGuardrail{
				Regex:    ".*success.*",
				JSONPath: "$.status",
				Inverted: false,
				logger:   createMockLoggerForRegex(),
			},
			payload:       createGzippedContentForRegex(`{"status": "operation successful"}`),
			isResponse:    true,
			expectedValid: true,
			expectedError: false,
		},
		{
			name: "Phone number validation - valid",
			guardrail: &RegexGuardrail{
				Regex:    "^\\+?[1-9]\\d{1,14}$",
				JSONPath: "$.phone",
				Inverted: false,
				logger:   createMockLoggerForRegex(),
			},
			payload:       []byte(`{"phone": "+1234567890"}`),
			isResponse:    false,
			expectedValid: true,
			expectedError: false,
		},
		{
			name: "Phone number validation - invalid",
			guardrail: &RegexGuardrail{
				Regex:    "^\\+?[1-9]\\d{1,14}$",
				JSONPath: "$.phone",
				Inverted: false,
				logger:   createMockLoggerForRegex(),
			},
			payload:       []byte(`{"phone": "invalid-phone"}`),
			isResponse:    false,
			expectedValid: false,
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, err := tt.guardrail.validatePayload(tt.payload, tt.isResponse)

			if (err != nil) != tt.expectedError {
				t.Errorf("Expected error %t, got error: %v", tt.expectedError, err)
			}

			if valid != tt.expectedValid {
				t.Errorf("Expected valid %t, got %t", tt.expectedValid, valid)
			}
		})
	}
}

func TestRegexGuardrail_extractStringValueFromJsonpath(t *testing.T) {
	guardrail := &RegexGuardrail{
		logger: createMockLoggerForRegex(),
	}

	tests := []struct {
		name        string
		payload     []byte
		jsonPath    string
		expectedVal string
		expectedErr bool
	}{
		{
			name:        "Simple string extraction",
			payload:     []byte(`{"content": "hello world"}`),
			jsonPath:    "$.content",
			expectedVal: "hello world",
			expectedErr: false,
		},
		{
			name:        "Nested object extraction",
			payload:     []byte(`{"user": {"name": "John Doe"}}`),
			jsonPath:    "$.user.name",
			expectedVal: "John Doe",
			expectedErr: false,
		},
		{
			name:        "Array element extraction",
			payload:     []byte(`{"items": ["first", "second", "third"]}`),
			jsonPath:    "$.items[0]",
			expectedVal: "first",
			expectedErr: false,
		},
		{
			name:        "Non-existent path",
			payload:     []byte(`{"content": "hello"}`),
			jsonPath:    "$.nonexistent",
			expectedVal: "",
			expectedErr: false,
		},
		{
			name:        "Number as string",
			payload:     []byte(`{"age": 25}`),
			jsonPath:    "$.age",
			expectedVal: "25",
			expectedErr: false,
		},
		{
			name:        "Boolean as string",
			payload:     []byte(`{"active": true}`),
			jsonPath:    "$.active",
			expectedVal: "true",
			expectedErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := guardrail.extractStringValueFromJsonpath(tt.payload, tt.jsonPath)

			if (err != nil) != tt.expectedErr {
				t.Errorf("Expected error %t, got error: %v", tt.expectedErr, err)
			}

			if result != tt.expectedVal {
				t.Errorf("Expected value '%s', got '%s'", tt.expectedVal, result)
			}
		})
	}
}

func TestRegexGuardrail_buildAssessmentObject(t *testing.T) {
	tests := []struct {
		name              string
		guardrail         *RegexGuardrail
		isResponse        bool
		validationError   error
		expectedReason    string
		expectedDirection string
	}{
		{
			name: "Request validation error",
			guardrail: &RegexGuardrail{
				Name:           "TestGuardrail",
				Regex:          "^test.*",
				ShowAssessment: true,
				logger:         createMockLoggerForRegex(),
			},
			isResponse:        false,
			validationError:   nil,
			expectedReason:    "Violation of regular expression detected.",
			expectedDirection: "REQUEST",
		},
		{
			name: "Response validation error",
			guardrail: &RegexGuardrail{
				Name:           "TestGuardrail",
				Regex:          "^test.*",
				ShowAssessment: true,
				logger:         createMockLoggerForRegex(),
			},
			isResponse:        true,
			validationError:   nil,
			expectedReason:    "Violation of regular expression detected.",
			expectedDirection: "RESPONSE",
		},
		{
			name: "JSONPath extraction error",
			guardrail: &RegexGuardrail{
				Name:           "TestGuardrail",
				JSONPath:       "$.invalid.path",
				ShowAssessment: true,
				logger:         createMockLoggerForRegex(),
			},
			isResponse:        false,
			validationError:   json.Unmarshal([]byte("invalid"), &map[string]interface{}{}),
			expectedReason:    "Error extracting content from payload using JSONPath.",
			expectedDirection: "REQUEST",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assessment := tt.guardrail.buildAssessmentObject(tt.isResponse, tt.validationError)

			if assessment[RegexAssessmentAction] != "GUARDRAIL_INTERVENED" {
				t.Errorf("Expected action 'GUARDRAIL_INTERVENED', got %v", assessment[RegexAssessmentAction])
			}

			if assessment[RegexInterveningGuardrail] != tt.guardrail.Name {
				t.Errorf("Expected guardrail name '%s', got %v", tt.guardrail.Name, assessment[RegexInterveningGuardrail])
			}

			if assessment[RegexDirection] != tt.expectedDirection {
				t.Errorf("Expected direction '%s', got %v", tt.expectedDirection, assessment[RegexDirection])
			}

			if assessment[RegexAssessmentReason] != tt.expectedReason {
				t.Errorf("Expected reason '%s', got %v", tt.expectedReason, assessment[RegexAssessmentReason])
			}

			// Check if assessment details are included when ShowAssessment is true
			if tt.guardrail.ShowAssessment {
				if _, exists := assessment[RegexAssessments]; !exists {
					t.Error("Expected assessments field when ShowAssessment is true")
				}
			}
		})
	}
}
