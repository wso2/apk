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
func createMockSentenceLogger() *logging.Logger {
	// Create a proper mock logger using the default logging setup
	mockLogger := logging.DefaultLogger(egv1a1.LogLevelInfo)
	return &mockLogger
}

// Helper function to create test mediation for sentence count guardrail
func createTestSentenceMediation(params map[string]string) *dpv2alpha1.Mediation {
	var parameters []*dpv2alpha1.Parameter
	for key, value := range params {
		parameters = append(parameters, &dpv2alpha1.Parameter{
			Key:   key,
			Value: value,
		})
	}

	return &dpv2alpha1.Mediation{
		PolicyName:    "SentenceCountGuardrail",
		PolicyVersion: "v1",
		PolicyID:      "test-sentence-policy-id",
		Parameters:    parameters,
	}
}

// Helper function to create gzipped content for sentence tests
func createGzippedSentenceContent(content string) []byte {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	gz.Write([]byte(content))
	gz.Close()
	return buf.Bytes()
}

func TestNewSentenceCountGuardrail(t *testing.T) {
	tests := []struct {
		name       string
		parameters map[string]string
		expected   SentenceCountGuardrail
	}{
		{
			name:       "Default values",
			parameters: map[string]string{},
			expected: SentenceCountGuardrail{
				PolicyName:     "SentenceCountGuardrail",
				PolicyVersion:  "v1",
				PolicyID:       "test-sentence-policy-id",
				Name:           "SentenceCountGuardrail",
				Min:            0,
				Max:            100,
				JSONPath:       "$.content",
				Inverted:       false,
				ShowAssessment: false,
			},
		},
		{
			name: "Custom values",
			parameters: map[string]string{
				"name":           "CustomSentenceGuardrail",
				"min":            "3",
				"max":            "10",
				"jsonPath":       "$.message",
				"invert":         "true",
				"showAssessment": "true",
			},
			expected: SentenceCountGuardrail{
				PolicyName:     "SentenceCountGuardrail",
				PolicyVersion:  "v1",
				PolicyID:       "test-sentence-policy-id",
				Name:           "CustomSentenceGuardrail",
				Min:            3,
				Max:            10,
				JSONPath:       "$.message",
				Inverted:       true,
				ShowAssessment: true,
			},
		},
		{
			name: "Invalid numeric values fall back to defaults",
			parameters: map[string]string{
				"min": "invalid",
				"max": "also-invalid",
			},
			expected: SentenceCountGuardrail{
				PolicyName:     "SentenceCountGuardrail",
				PolicyVersion:  "v1",
				PolicyID:       "test-sentence-policy-id",
				Name:           "SentenceCountGuardrail",
				Min:            0,
				Max:            100,
				JSONPath:       "$.content",
				Inverted:       false,
				ShowAssessment: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mediation := createTestSentenceMediation(tt.parameters)
			result := NewSentenceCountGuardrail(mediation)

			if result.PolicyName != tt.expected.PolicyName {
				t.Errorf("Expected PolicyName %s, got %s", tt.expected.PolicyName, result.PolicyName)
			}
			if result.Name != tt.expected.Name {
				t.Errorf("Expected Name %s, got %s", tt.expected.Name, result.Name)
			}
			if result.Min != tt.expected.Min {
				t.Errorf("Expected Min %d, got %d", tt.expected.Min, result.Min)
			}
			if result.Max != tt.expected.Max {
				t.Errorf("Expected Max %d, got %d", tt.expected.Max, result.Max)
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

func TestSentenceCountGuardrail_Process(t *testing.T) {
	tests := []struct {
		name              string
		guardrail         *SentenceCountGuardrail
		requestConfig     *requestconfig.Holder
		expectedImmediate bool
		expectedPassed    bool
	}{
		{
			name: "Request body - valid sentence count",
			guardrail: &SentenceCountGuardrail{
				Name:     "TestSentenceGuardrail",
				Min:      2,
				Max:      4,
				JSONPath: "$.content",
				logger:   createMockSentenceLogger(),
			},
			requestConfig: &requestconfig.Holder{
				ProcessingPhase: requestconfig.ProcessingPhaseRequestBody,
				RequestBody: &envoy_service_proc_v3.HttpBody{
					Body: []byte(`{"content": "Hello world. This is a test. Good day!"}`),
				},
			},
			expectedImmediate: false,
			expectedPassed:    true,
		},
		{
			name: "Request body - invalid sentence count (too few)",
			guardrail: &SentenceCountGuardrail{
				Name:     "TestSentenceGuardrail",
				Min:      5,
				Max:      10,
				JSONPath: "$.content",
				logger:   createMockSentenceLogger(),
			},
			requestConfig: &requestconfig.Holder{
				ProcessingPhase: requestconfig.ProcessingPhaseRequestBody,
				RequestBody: &envoy_service_proc_v3.HttpBody{
					Body: []byte(`{"content": "Hello world."}`),
				},
			},
			expectedImmediate: true,
			expectedPassed:    false,
		},
		{
			name: "Response body - valid sentence count",
			guardrail: &SentenceCountGuardrail{
				Name:     "TestSentenceGuardrail",
				Min:      1,
				Max:      3,
				JSONPath: "$.content",
				logger:   createMockSentenceLogger(),
			},
			requestConfig: &requestconfig.Holder{
				ProcessingPhase: requestconfig.ProcessingPhaseResponseBody,
				ResponseBody: &envoy_service_proc_v3.HttpBody{
					Body: []byte(`{"content": "Hello! How are you?"}`),
				},
			},
			expectedImmediate: false,
			expectedPassed:    true,
		},
		{
			name: "Response body - invalid sentence count (too many)",
			guardrail: &SentenceCountGuardrail{
				Name:     "TestSentenceGuardrail",
				Min:      1,
				Max:      2,
				JSONPath: "$.content",
				logger:   createMockSentenceLogger(),
			},
			requestConfig: &requestconfig.Holder{
				ProcessingPhase: requestconfig.ProcessingPhaseResponseBody,
				ResponseBody: &envoy_service_proc_v3.HttpBody{
					Body: []byte(`{"content": "Hello world! This is a test. How are you? Good day!"}`),
				},
			},
			expectedImmediate: true,
			expectedPassed:    false,
		},
		{
			name: "No request body",
			guardrail: &SentenceCountGuardrail{
				Name:     "TestSentenceGuardrail",
				Min:      1,
				Max:      5,
				JSONPath: "$.content",
				logger:   createMockSentenceLogger(),
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
			guardrail: &SentenceCountGuardrail{
				Name:     "TestSentenceGuardrail",
				Min:      1,
				Max:      5,
				JSONPath: "$.content",
				logger:   createMockSentenceLogger(),
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
			}
		})
	}
}

func TestSentenceCountGuardrail_validatePayload(t *testing.T) {
	tests := []struct {
		name       string
		guardrail  *SentenceCountGuardrail
		payload    []byte
		isResponse bool
		expected   bool
		expectErr  bool
	}{
		{
			name: "Valid sentence count within range",
			guardrail: &SentenceCountGuardrail{
				Min:      2,
				Max:      4,
				JSONPath: "$.content",
				Inverted: false,
				logger:   createMockSentenceLogger(),
			},
			payload:   []byte(`{"content": "Hello world. This is a test. Good day!"}`),
			expected:  true,
			expectErr: false,
		},
		{
			name: "Sentence count below minimum",
			guardrail: &SentenceCountGuardrail{
				Min:      5,
				Max:      10,
				JSONPath: "$.content",
				Inverted: false,
				logger:   createMockSentenceLogger(),
			},
			payload:   []byte(`{"content": "Hello world."}`),
			expected:  false,
			expectErr: false,
		},
		{
			name: "Sentence count above maximum",
			guardrail: &SentenceCountGuardrail{
				Min:      1,
				Max:      2,
				JSONPath: "$.content",
				Inverted: false,
				logger:   createMockSentenceLogger(),
			},
			payload:   []byte(`{"content": "Hello world! This is a test. How are you? Good day!"}`),
			expected:  false,
			expectErr: false,
		},
		{
			name: "Inverted logic - sentence count within range (should fail)",
			guardrail: &SentenceCountGuardrail{
				Min:      2,
				Max:      4,
				JSONPath: "$.content",
				Inverted: true,
				logger:   createMockSentenceLogger(),
			},
			payload:   []byte(`{"content": "Hello world. This is a test. Good day!"}`),
			expected:  false,
			expectErr: false,
		},
		{
			name: "Inverted logic - sentence count outside range (should pass)",
			guardrail: &SentenceCountGuardrail{
				Min:      5,
				Max:      10,
				JSONPath: "$.content",
				Inverted: true,
				logger:   createMockSentenceLogger(),
			},
			payload:   []byte(`{"content": "Hello world."}`),
			expected:  true,
			expectErr: false,
		},
		{
			name: "Invalid JSON path",
			guardrail: &SentenceCountGuardrail{
				Min:      1,
				Max:      5,
				JSONPath: "$.nonexistent",
				Inverted: false,
				logger:   createMockSentenceLogger(),
			},
			payload:   []byte(`{"content": "Hello world! This is a test."}`),
			expected:  false,
			expectErr: true,
		},
		{
			name: "Invalid sentence count range (min > max)",
			guardrail: &SentenceCountGuardrail{
				Min:      10,
				Max:      5, // Max < Min
				JSONPath: "$.content",
				Inverted: false,
				logger:   createMockSentenceLogger(),
			},
			payload:   []byte(`{"content": "Hello world! This is a test."}`),
			expected:  false,
			expectErr: false,
		},
		{
			name: "Invalid sentence count range (negative min)",
			guardrail: &SentenceCountGuardrail{
				Min:      -1,
				Max:      5,
				JSONPath: "$.content",
				Inverted: false,
				logger:   createMockSentenceLogger(),
			},
			payload:   []byte(`{"content": "Hello world! This is a test."}`),
			expected:  false,
			expectErr: false,
		},
		{
			name: "Invalid sentence count range (zero max)",
			guardrail: &SentenceCountGuardrail{
				Min:      1,
				Max:      0,
				JSONPath: "$.content",
				Inverted: false,
				logger:   createMockSentenceLogger(),
			},
			payload:   []byte(`{"content": "Hello world! This is a test."}`),
			expected:  false,
			expectErr: false,
		},
		{
			name: "Empty content after cleaning",
			guardrail: &SentenceCountGuardrail{
				Min:      1,
				Max:      5,
				JSONPath: "$.content",
				Inverted: false,
				logger:   createMockSentenceLogger(),
			},
			payload:   []byte(`{"content": "   "}`), // Only whitespace
			expected:  false,
			expectErr: false,
		},
		{
			name: "Content with punctuation but no sentences",
			guardrail: &SentenceCountGuardrail{
				Min:      1,
				Max:      5,
				JSONPath: "$.content",
				Inverted: false,
				logger:   createMockSentenceLogger(),
			},
			payload:   []byte(`{"content": ";;;"}`), // Only punctuation
			expected:  false,
			expectErr: false,
		},
		{
			name: "Mixed sentence endings",
			guardrail: &SentenceCountGuardrail{
				Min:      3,
				Max:      3,
				JSONPath: "$.content",
				Inverted: false,
				logger:   createMockSentenceLogger(),
			},
			payload:   []byte(`{"content": "Hello world! How are you? Great day."}`),
			expected:  true,
			expectErr: false,
		},
		{
			name: "Response body - compressed content",
			guardrail: &SentenceCountGuardrail{
				Min:      2,
				Max:      4,
				JSONPath: "$.content",
				Inverted: false,
				logger:   createMockSentenceLogger(),
			},
			payload:    createGzippedSentenceContent(`{"content": "Hello world! This is compressed. Good day!"}`),
			isResponse: true,
			expected:   true,
			expectErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.guardrail.validatePayload(tt.payload, tt.isResponse)

			if tt.expectErr {
				if err == nil {
					t.Error("Expected an error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}

			if result != tt.expected {
				t.Errorf("Expected validation result %t, got %t", tt.expected, result)
			}
		})
	}
}

func TestSentenceCountGuardrail_buildErrorResponse(t *testing.T) {
	tests := []struct {
		name            string
		guardrail       *SentenceCountGuardrail
		isResponse      bool
		validationError error
		expectedFields  map[string]interface{}
	}{
		{
			name: "Request error without validation error",
			guardrail: &SentenceCountGuardrail{
				Name:           "TestGuardrail",
				Min:            2,
				Max:            5,
				ShowAssessment: true,
				Inverted:       false,
				logger:         createMockSentenceLogger(),
			},
			isResponse:      false,
			validationError: nil,
			expectedFields: map[string]interface{}{
				"errorCode": "GUARDRAIL_API_EXCEPTION",
				"errorType": "SENTENCE_COUNT_GUARDRAIL",
			},
		},
		{
			name: "Response error with validation error",
			guardrail: &SentenceCountGuardrail{
				Name:           "TestGuardrail",
				JSONPath:       "$.invalid",
				ShowAssessment: true,
				logger:         createMockSentenceLogger(),
			},
			isResponse:      true,
			validationError: &SentenceValidationError{message: "field not found"},
			expectedFields: map[string]interface{}{
				"errorCode": "GUARDRAIL_API_EXCEPTION",
				"errorType": "SENTENCE_COUNT_GUARDRAIL",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.guardrail.buildErrorResponse(tt.isResponse, tt.validationError)

			if !result.ImmediateResponse {
				t.Error("Expected ImmediateResponse to be true")
			}

			if result.ImmediateResponseContentType != "application/json" {
				t.Errorf("Expected ContentType 'application/json', got '%s'", result.ImmediateResponseContentType)
			}

			if result.ImmediateResponseBody == "" {
				t.Error("Expected non-empty response body")
			}

			// Parse the response body
			var responseBody map[string]interface{}
			err := json.Unmarshal([]byte(result.ImmediateResponseBody), &responseBody)
			if err != nil {
				t.Fatalf("Failed to parse response body: %v", err)
			}

			// Check expected fields
			for key, expectedValue := range tt.expectedFields {
				if actualValue, exists := responseBody[key]; !exists {
					t.Errorf("Expected field '%s' not found in response", key)
				} else if actualValue != expectedValue {
					t.Errorf("Expected field '%s' to be '%v', got '%v'", key, expectedValue, actualValue)
				}
			}

			// Check that errorMessage exists and is a map
			if errorMessage, exists := responseBody["errorMessage"]; !exists {
				t.Error("Expected 'errorMessage' field in response")
			} else {
				if _, isMap := errorMessage.(map[string]interface{}); !isMap {
					t.Error("Expected 'errorMessage' to be a map")
				}
			}
		})
	}
}

func TestSentenceCountGuardrail_buildAssessmentObject(t *testing.T) {
	tests := []struct {
		name            string
		guardrail       *SentenceCountGuardrail
		isResponse      bool
		validationError error
		expectedAction  string
		expectedDir     string
	}{
		{
			name: "Request assessment without error",
			guardrail: &SentenceCountGuardrail{
				Name:           "TestGuardrail",
				Min:            2,
				Max:            5,
				ShowAssessment: true,
				Inverted:       false,
				logger:         createMockSentenceLogger(),
			},
			isResponse:     false,
			expectedAction: "GUARDRAIL_INTERVENED",
			expectedDir:    "REQUEST",
		},
		{
			name: "Response assessment with error",
			guardrail: &SentenceCountGuardrail{
				Name:           "TestGuardrail",
				JSONPath:       "$.invalid",
				ShowAssessment: false,
				logger:         createMockSentenceLogger(),
			},
			isResponse:      true,
			validationError: &SentenceValidationError{message: "field not found"},
			expectedAction:  "GUARDRAIL_INTERVENED",
			expectedDir:     "RESPONSE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assessment := tt.guardrail.buildAssessmentObject(tt.isResponse, tt.validationError)

			if action, exists := assessment["action"]; !exists {
				t.Error("Expected 'action' field in assessment")
			} else if action != tt.expectedAction {
				t.Errorf("Expected action '%s', got '%s'", tt.expectedAction, action)
			}

			if direction, exists := assessment["direction"]; !exists {
				t.Error("Expected 'direction' field in assessment")
			} else if direction != tt.expectedDir {
				t.Errorf("Expected direction '%s', got '%s'", tt.expectedDir, direction)
			}

			if guardrail, exists := assessment["interveningGuardrail"]; !exists {
				t.Error("Expected 'interveningGuardrail' field in assessment")
			} else if guardrail != tt.guardrail.Name {
				t.Errorf("Expected interveningGuardrail '%s', got '%s'", tt.guardrail.Name, guardrail)
			}

			if _, exists := assessment["reason"]; !exists {
				t.Error("Expected 'reason' field in assessment")
			}
		})
	}
}

func TestSentenceCountGuardrail_decompressLLMResp(t *testing.T) {
	tests := []struct {
		name     string
		payload  []byte
		expected string
		hasError bool
	}{
		{
			name:     "Non-compressed content",
			payload:  []byte("Hello world! This is a test."),
			expected: "Hello world! This is a test.",
			hasError: false,
		},
		{
			name:     "Gzipped content",
			payload:  createGzippedSentenceContent("Hello world! This is compressed."),
			expected: "Hello world! This is compressed.",
			hasError: false,
		},
		{
			name:     "Empty payload",
			payload:  []byte{},
			expected: "",
			hasError: false,
		},
		{
			name:     "Small payload (less than 2 bytes)",
			payload:  []byte{0x1f},
			expected: "\x1f",
			hasError: false,
		},
	}

	guardrail := &SentenceCountGuardrail{
		logger: createMockSentenceLogger(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := guardrail.decompressLLMResp(tt.payload)

			if tt.hasError {
				if err == nil {
					t.Error("Expected an error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}

			if result != tt.expected {
				t.Errorf("Expected result '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestSentenceCountGuardrail_extractStringValueFromJsonpath(t *testing.T) {
	tests := []struct {
		name      string
		payload   []byte
		jsonPath  string
		expected  string
		expectErr bool
	}{
		{
			name:      "Valid JSON path",
			payload:   []byte(`{"content": "Hello world! This is a test."}`),
			jsonPath:  "$.content",
			expected:  "Hello world! This is a test.",
			expectErr: false,
		},
		{
			name:      "Nested JSON path",
			payload:   []byte(`{"data": {"message": "Hello world! How are you?"}}`),
			jsonPath:  "$.data.message",
			expected:  "Hello world! How are you?",
			expectErr: false,
		},
		{
			name:      "Non-existent path",
			payload:   []byte(`{"content": "Hello world!"}`),
			jsonPath:  "$.nonexistent",
			expected:  "",
			expectErr: true,
		},
		{
			name:      "Invalid JSON",
			payload:   []byte(`invalid json`),
			jsonPath:  "$.content",
			expected:  "",
			expectErr: true,
		},
	}

	guardrail := &SentenceCountGuardrail{
		logger: createMockSentenceLogger(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := guardrail.extractStringValueFromJsonpath(tt.payload, tt.jsonPath)

			if tt.expectErr {
				if err == nil {
					t.Error("Expected an error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}

			if result != tt.expected {
				t.Errorf("Expected result '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

// Custom error type for testing
type SentenceValidationError struct {
	message string
}

func (e *SentenceValidationError) Error() string {
	return e.message
}
