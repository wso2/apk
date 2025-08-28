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

// Mock logger for content length testing
func createMockContentLengthLogger() *logging.Logger {
	// Create a proper mock logger using the default logging setup
	mockLogger := logging.DefaultLogger(egv1a1.LogLevelInfo)
	return &mockLogger
}

// Helper function to create test mediation for content length guardrail
func createTestContentLengthMediation(params map[string]string) *dpv2alpha1.Mediation {
	var parameters []*dpv2alpha1.Parameter
	for key, value := range params {
		parameters = append(parameters, &dpv2alpha1.Parameter{
			Key:   key,
			Value: value,
		})
	}

	return &dpv2alpha1.Mediation{
		PolicyName:    "ContentLengthGuardrail",
		PolicyVersion: "v1",
		PolicyID:      "test-content-length-policy-id",
		Parameters:    parameters,
	}
}

// Helper function to create gzipped content for content length tests
func createGzippedContentLengthContent(content string) []byte {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	gz.Write([]byte(content))
	gz.Close()
	return buf.Bytes()
}

func TestNewContentLengthGuardrail(t *testing.T) {
	tests := []struct {
		name       string
		parameters map[string]string
		expected   ContentLengthGuardrail
	}{
		{
			name:       "Default values",
			parameters: map[string]string{},
			expected: ContentLengthGuardrail{
				PolicyName:     "ContentLengthGuardrail",
				PolicyVersion:  "v1",
				PolicyID:       "test-content-length-policy-id",
				Name:           "ContentLengthGuardrail",
				Min:            0,
				Max:            10000,
				JSONPath:       "$.content",
				Inverted:       false,
				ShowAssessment: false,
			},
		},
		{
			name: "Custom values",
			parameters: map[string]string{
				"name":           "CustomContentLengthGuardrail",
				"min":            "50",
				"max":            "500",
				"jsonPath":       "$.message",
				"invert":         "true",
				"showAssessment": "true",
			},
			expected: ContentLengthGuardrail{
				PolicyName:     "ContentLengthGuardrail",
				PolicyVersion:  "v1",
				PolicyID:       "test-content-length-policy-id",
				Name:           "CustomContentLengthGuardrail",
				Min:            50,
				Max:            500,
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
			expected: ContentLengthGuardrail{
				PolicyName:     "ContentLengthGuardrail",
				PolicyVersion:  "v1",
				PolicyID:       "test-content-length-policy-id",
				Name:           "ContentLengthGuardrail",
				Min:            0,
				Max:            10000,
				JSONPath:       "$.content",
				Inverted:       false,
				ShowAssessment: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mediation := createTestContentLengthMediation(tt.parameters)
			result := NewContentLengthGuardrail(mediation)

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

func TestContentLengthGuardrail_Process(t *testing.T) {
	tests := []struct {
		name              string
		guardrail         *ContentLengthGuardrail
		requestConfig     *requestconfig.Holder
		expectedImmediate bool
		expectedPassed    bool
	}{
		{
			name: "Request body - valid content length",
			guardrail: &ContentLengthGuardrail{
				Name:     "TestContentLengthGuardrail",
				Min:      5,
				Max:      20,
				JSONPath: "$.content",
				logger:   createMockContentLengthLogger(),
			},
			requestConfig: &requestconfig.Holder{
				ProcessingPhase: requestconfig.ProcessingPhaseRequestBody,
				RequestBody: &envoy_service_proc_v3.HttpBody{
					Body: []byte(`{"content": "Hello world"}`), // 11 characters
				},
			},
			expectedImmediate: false,
			expectedPassed:    true,
		},
		{
			name: "Request body - content too short",
			guardrail: &ContentLengthGuardrail{
				Name:     "TestContentLengthGuardrail",
				Min:      20,
				Max:      100,
				JSONPath: "$.content",
				logger:   createMockContentLengthLogger(),
			},
			requestConfig: &requestconfig.Holder{
				ProcessingPhase: requestconfig.ProcessingPhaseRequestBody,
				RequestBody: &envoy_service_proc_v3.HttpBody{
					Body: []byte(`{"content": "Hi"}`), // 2 characters
				},
			},
			expectedImmediate: true,
			expectedPassed:    false,
		},
		{
			name: "Response body - valid content length",
			guardrail: &ContentLengthGuardrail{
				Name:     "TestContentLengthGuardrail",
				Min:      1,
				Max:      15,
				JSONPath: "$.content",
				logger:   createMockContentLengthLogger(),
			},
			requestConfig: &requestconfig.Holder{
				ProcessingPhase: requestconfig.ProcessingPhaseResponseBody,
				ResponseBody: &envoy_service_proc_v3.HttpBody{
					Body: []byte(`{"content": "Hello world"}`), // 11 characters
				},
			},
			expectedImmediate: false,
			expectedPassed:    true,
		},
		{
			name: "Response body - content too long",
			guardrail: &ContentLengthGuardrail{
				Name:     "TestContentLengthGuardrail",
				Min:      1,
				Max:      5,
				JSONPath: "$.content",
				logger:   createMockContentLengthLogger(),
			},
			requestConfig: &requestconfig.Holder{
				ProcessingPhase: requestconfig.ProcessingPhaseResponseBody,
				ResponseBody: &envoy_service_proc_v3.HttpBody{
					Body: []byte(`{"content": "This is a very long message that exceeds the limit"}`), // ~49 characters
				},
			},
			expectedImmediate: true,
			expectedPassed:    false,
		},
		{
			name: "No request body",
			guardrail: &ContentLengthGuardrail{
				Name:     "TestContentLengthGuardrail",
				Min:      1,
				Max:      100,
				JSONPath: "$.content",
				logger:   createMockContentLengthLogger(),
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
			guardrail: &ContentLengthGuardrail{
				Name:     "TestContentLengthGuardrail",
				Min:      1,
				Max:      100,
				JSONPath: "$.content",
				logger:   createMockContentLengthLogger(),
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

func TestContentLengthGuardrail_validatePayload(t *testing.T) {
	tests := []struct {
		name       string
		guardrail  *ContentLengthGuardrail
		payload    []byte
		isResponse bool
		expected   bool
		expectErr  bool
	}{
		{
			name: "Valid content length within range",
			guardrail: &ContentLengthGuardrail{
				Min:      5,
				Max:      20,
				JSONPath: "$.content",
				Inverted: false,
				logger:   createMockContentLengthLogger(),
			},
			payload:   []byte(`{"content": "Hello world"}`), // 11 characters
			expected:  true,
			expectErr: false,
		},
		{
			name: "Content length below minimum",
			guardrail: &ContentLengthGuardrail{
				Min:      20,
				Max:      100,
				JSONPath: "$.content",
				Inverted: false,
				logger:   createMockContentLengthLogger(),
			},
			payload:   []byte(`{"content": "Hi"}`), // 2 characters
			expected:  false,
			expectErr: false,
		},
		{
			name: "Content length above maximum",
			guardrail: &ContentLengthGuardrail{
				Min:      1,
				Max:      10,
				JSONPath: "$.content",
				Inverted: false,
				logger:   createMockContentLengthLogger(),
			},
			payload:   []byte(`{"content": "This is a very long message"}`), // 26 characters
			expected:  false,
			expectErr: false,
		},
		{
			name: "Inverted logic - content length within range (should fail)",
			guardrail: &ContentLengthGuardrail{
				Min:      5,
				Max:      20,
				JSONPath: "$.content",
				Inverted: true,
				logger:   createMockContentLengthLogger(),
			},
			payload:   []byte(`{"content": "Hello world"}`), // 11 characters
			expected:  false,
			expectErr: false,
		},
		{
			name: "Inverted logic - content length outside range (should pass)",
			guardrail: &ContentLengthGuardrail{
				Min:      20,
				Max:      100,
				JSONPath: "$.content",
				Inverted: true,
				logger:   createMockContentLengthLogger(),
			},
			payload:   []byte(`{"content": "Hi"}`), // 2 characters
			expected:  true,
			expectErr: false,
		},
		{
			name: "Invalid JSON path",
			guardrail: &ContentLengthGuardrail{
				Min:      1,
				Max:      100,
				JSONPath: "$.nonexistent",
				Inverted: false,
				logger:   createMockContentLengthLogger(),
			},
			payload:   []byte(`{"content": "Hello world"}`),
			expected:  false,
			expectErr: true,
		},
		{
			name: "Invalid content length range (min > max)",
			guardrail: &ContentLengthGuardrail{
				Min:      100,
				Max:      50, // Max < Min
				JSONPath: "$.content",
				Inverted: false,
				logger:   createMockContentLengthLogger(),
			},
			payload:   []byte(`{"content": "Hello world"}`),
			expected:  false,
			expectErr: false,
		},
		{
			name: "Invalid content length range (negative min)",
			guardrail: &ContentLengthGuardrail{
				Min:      -1,
				Max:      100,
				JSONPath: "$.content",
				Inverted: false,
				logger:   createMockContentLengthLogger(),
			},
			payload:   []byte(`{"content": "Hello world"}`),
			expected:  false,
			expectErr: false,
		},
		{
			name: "Invalid content length range (zero max)",
			guardrail: &ContentLengthGuardrail{
				Min:      1,
				Max:      0,
				JSONPath: "$.content",
				Inverted: false,
				logger:   createMockContentLengthLogger(),
			},
			payload:   []byte(`{"content": "Hello world"}`),
			expected:  false,
			expectErr: false,
		},
		{
			name: "Empty content after cleaning",
			guardrail: &ContentLengthGuardrail{
				Min:      1,
				Max:      100,
				JSONPath: "$.content",
				Inverted: false,
				logger:   createMockContentLengthLogger(),
			},
			payload:   []byte(`{"content": "   "}`), // Only whitespace
			expected:  false,
			expectErr: false,
		},
		{
			name: "Unicode content length",
			guardrail: &ContentLengthGuardrail{
				Min:      5,
				Max:      20,
				JSONPath: "$.content",
				Inverted: false,
				logger:   createMockContentLengthLogger(),
			},
			payload:   []byte(`{"content": "Hello 世界"}`), // Mixed ASCII and Unicode
			expected:  true,
			expectErr: false,
		},
		{
			name: "Response body - compressed content",
			guardrail: &ContentLengthGuardrail{
				Min:      5,
				Max:      20,
				JSONPath: "$.content",
				Inverted: false,
				logger:   createMockContentLengthLogger(),
			},
			payload:    createGzippedContentLengthContent(`{"content": "Hello world"}`),
			isResponse: true,
			expected:   true,
			expectErr:  false,
		},
		{
			name: "Exact boundary - minimum",
			guardrail: &ContentLengthGuardrail{
				Min:      11,
				Max:      20,
				JSONPath: "$.content",
				Inverted: false,
				logger:   createMockContentLengthLogger(),
			},
			payload:   []byte(`{"content": "Hello world"}`), // Exactly 11 characters
			expected:  true,
			expectErr: false,
		},
		{
			name: "Exact boundary - maximum",
			guardrail: &ContentLengthGuardrail{
				Min:      5,
				Max:      11,
				JSONPath: "$.content",
				Inverted: false,
				logger:   createMockContentLengthLogger(),
			},
			payload:   []byte(`{"content": "Hello world"}`), // Exactly 11 characters
			expected:  true,
			expectErr: false,
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

func TestContentLengthGuardrail_buildErrorResponse(t *testing.T) {
	tests := []struct {
		name            string
		guardrail       *ContentLengthGuardrail
		isResponse      bool
		validationError error
		expectedFields  map[string]interface{}
	}{
		{
			name: "Request error without validation error",
			guardrail: &ContentLengthGuardrail{
				Name:           "TestGuardrail",
				Min:            5,
				Max:            100,
				ShowAssessment: true,
				Inverted:       false,
				logger:         createMockContentLengthLogger(),
			},
			isResponse:      false,
			validationError: nil,
			expectedFields: map[string]interface{}{
				"errorCode": "GUARDRAIL_API_EXCEPTION",
				"errorType": "CONTENT_LENGTH_GUARDRAIL",
			},
		},
		{
			name: "Response error with validation error",
			guardrail: &ContentLengthGuardrail{
				Name:           "TestGuardrail",
				JSONPath:       "$.invalid",
				ShowAssessment: true,
				logger:         createMockContentLengthLogger(),
			},
			isResponse:      true,
			validationError: &ContentLengthValidationError{message: "field not found"},
			expectedFields: map[string]interface{}{
				"errorCode": "GUARDRAIL_API_EXCEPTION",
				"errorType": "CONTENT_LENGTH_GUARDRAIL",
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

func TestContentLengthGuardrail_buildAssessmentObject(t *testing.T) {
	tests := []struct {
		name            string
		guardrail       *ContentLengthGuardrail
		isResponse      bool
		validationError error
		expectedAction  string
		expectedDir     string
	}{
		{
			name: "Request assessment without error",
			guardrail: &ContentLengthGuardrail{
				Name:           "TestGuardrail",
				Min:            5,
				Max:            100,
				ShowAssessment: true,
				Inverted:       false,
				logger:         createMockContentLengthLogger(),
			},
			isResponse:     false,
			expectedAction: "GUARDRAIL_INTERVENED",
			expectedDir:    "REQUEST",
		},
		{
			name: "Response assessment with error",
			guardrail: &ContentLengthGuardrail{
				Name:           "TestGuardrail",
				JSONPath:       "$.invalid",
				ShowAssessment: false,
				logger:         createMockContentLengthLogger(),
			},
			isResponse:      true,
			validationError: &ContentLengthValidationError{message: "field not found"},
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

func TestContentLengthGuardrail_decompressLLMResp(t *testing.T) {
	tests := []struct {
		name     string
		payload  []byte
		expected string
		hasError bool
	}{
		{
			name:     "Non-compressed content",
			payload:  []byte("Hello world content"),
			expected: "Hello world content",
			hasError: false,
		},
		{
			name:     "Gzipped content",
			payload:  createGzippedContentLengthContent("Hello world compressed content"),
			expected: "Hello world compressed content",
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

	guardrail := &ContentLengthGuardrail{
		logger: createMockContentLengthLogger(),
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

func TestContentLengthGuardrail_extractStringValueFromJsonpath(t *testing.T) {
	tests := []struct {
		name      string
		payload   []byte
		jsonPath  string
		expected  string
		expectErr bool
	}{
		{
			name:      "Valid JSON path",
			payload:   []byte(`{"content": "Hello world content"}`),
			jsonPath:  "$.content",
			expected:  "Hello world content",
			expectErr: false,
		},
		{
			name:      "Nested JSON path",
			payload:   []byte(`{"data": {"message": "Hello nested content"}}`),
			jsonPath:  "$.data.message",
			expected:  "Hello nested content",
			expectErr: false,
		},
		{
			name:      "Non-existent path",
			payload:   []byte(`{"content": "Hello world"}`),
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

	guardrail := &ContentLengthGuardrail{
		logger: createMockContentLengthLogger(),
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
type ContentLengthValidationError struct {
	message string
}

func (e *ContentLengthValidationError) Error() string {
	return e.message
}
