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
	"fmt"
	"strings"
	"testing"

	egv1a1 "github.com/envoyproxy/gateway/api/v1alpha1"
	envoy_service_proc_v3 "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"
	dpv2alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v2alpha1"
	"github.com/wso2/apk/gateway/enforcer/internal/logging"
	"github.com/wso2/apk/gateway/enforcer/internal/requestconfig"
)

// Mock logger for testing
func createMockLogger() *logging.Logger {
	// Create a proper mock logger using the default logging setup
	mockLogger := logging.DefaultLogger(egv1a1.LogLevelInfo)
	return &mockLogger
}

// Helper function to create test mediation
func createTestMediation(params map[string]string) *dpv2alpha1.Mediation {
	var parameters []*dpv2alpha1.Parameter
	for key, value := range params {
		parameters = append(parameters, &dpv2alpha1.Parameter{
			Key:   key,
			Value: value,
		})
	}

	return &dpv2alpha1.Mediation{
		PolicyName:    "WordCountGuardrail",
		PolicyVersion: "v1",
		PolicyID:      "test-policy-id",
		Parameters:    parameters,
	}
}

// Helper function to create gzipped content
func createGzippedContent(content string) []byte {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	gz.Write([]byte(content))
	gz.Close()
	return buf.Bytes()
}

func TestNewWordCountGuardrail(t *testing.T) {
	tests := []struct {
		name       string
		parameters map[string]string
		expected   WordCountGuardrail
	}{
		{
			name: "Default values",
			parameters: map[string]string{},
			expected: WordCountGuardrail{
				PolicyName:     "WordCountGuardrail",
				PolicyVersion:  "v1",
				PolicyID:       "test-policy-id",
				Name:           "WordCountGuardrail",
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
				"name":           "CustomGuardrail",
				"min":            "10",
				"max":            "50",
				"jsonPath":       "$.message",
				"invert":         "true",
				"showAssessment": "true",
			},
			expected: WordCountGuardrail{
				PolicyName:     "WordCountGuardrail",
				PolicyVersion:  "v1",
				PolicyID:       "test-policy-id",
				Name:           "CustomGuardrail",
				Min:            10,
				Max:            50,
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
			expected: WordCountGuardrail{
				PolicyName:     "WordCountGuardrail",
				PolicyVersion:  "v1",
				PolicyID:       "test-policy-id",
				Name:           "WordCountGuardrail",
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
			mediation := createTestMediation(tt.parameters)
			result := NewWordCountGuardrail(mediation)

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

func TestWordCountGuardrail_Process(t *testing.T) {
	tests := []struct {
		name                string
		guardrail          *WordCountGuardrail
		requestConfig      *requestconfig.Holder
		expectedImmediate  bool
		expectedPassed     bool
	}{
		{
			name: "Request body - valid word count",
			guardrail: &WordCountGuardrail{
				Name:     "TestGuardrail",
				Min:      2,
				Max:      5,
				JSONPath: "$.content",
				logger:   createMockLogger(),
			},
			requestConfig: &requestconfig.Holder{
				ProcessingPhase: requestconfig.ProcessingPhaseRequestBody,
				RequestBody: &envoy_service_proc_v3.HttpBody{
					Body: []byte(`{"content": "hello world test"}`),
				},
			},
			expectedImmediate: false,
			expectedPassed:    true,
		},
		{
			name: "Request body - invalid word count (too few)",
			guardrail: &WordCountGuardrail{
				Name:     "TestGuardrail",
				Min:      5,
				Max:      10,
				JSONPath: "$.content",
				logger:   createMockLogger(),
			},
			requestConfig: &requestconfig.Holder{
				ProcessingPhase: requestconfig.ProcessingPhaseRequestBody,
				RequestBody: &envoy_service_proc_v3.HttpBody{
					Body: []byte(`{"content": "hello"}`),
				},
			},
			expectedImmediate: true,
			expectedPassed:    false,
		},
		{
			name: "Response body - valid word count",
			guardrail: &WordCountGuardrail{
				Name:     "TestGuardrail",
				Min:      2,
				Max:      5,
				JSONPath: "$.content",
				logger:   createMockLogger(),
			},
			requestConfig: &requestconfig.Holder{
				ProcessingPhase: requestconfig.ProcessingPhaseResponseBody,
				ResponseBody: &envoy_service_proc_v3.HttpBody{
					Body: []byte(`{"content": "hello world"}`),
				},
			},
			expectedImmediate: false,
			expectedPassed:    true,
		},
		{
			name: "Response body - invalid word count (too many)",
			guardrail: &WordCountGuardrail{
				Name:     "TestGuardrail",
				Min:      1,
				Max:      2,
				JSONPath: "$.content",
				logger:   createMockLogger(),
			},
			requestConfig: &requestconfig.Holder{
				ProcessingPhase: requestconfig.ProcessingPhaseResponseBody,
				ResponseBody: &envoy_service_proc_v3.HttpBody{
					Body: []byte(`{"content": "hello world test more words"}`),
				},
			},
			expectedImmediate: true,
			expectedPassed:    false,
		},
		{
			name: "No request body",
			guardrail: &WordCountGuardrail{
				Name:     "TestGuardrail",
				Min:      1,
				Max:      5,
				JSONPath: "$.content",
				logger:   createMockLogger(),
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
			guardrail: &WordCountGuardrail{
				Name:     "TestGuardrail",
				Min:      1,
				Max:      5,
				JSONPath: "$.content",
				logger:   createMockLogger(),
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

func TestWordCountGuardrail_validatePayload(t *testing.T) {
	tests := []struct {
		name       string
		guardrail  *WordCountGuardrail
		payload    []byte
		isResponse bool
		expected   bool
		expectErr  bool
	}{
		{
			name: "Valid word count within range",
			guardrail: &WordCountGuardrail{
				Min:      2,
				Max:      5,
				JSONPath: "$.content",
				Inverted: false,
				logger:   createMockLogger(),
			},
			payload:   []byte(`{"content": "hello world test"}`),
			expected:  true,
			expectErr: false,
		},
		{
			name: "Word count below minimum",
			guardrail: &WordCountGuardrail{
				Min:      5,
				Max:      10,
				JSONPath: "$.content",
				Inverted: false,
				logger:   createMockLogger(),
			},
			payload:   []byte(`{"content": "hello"}`),
			expected:  false,
			expectErr: false,
		},
		{
			name: "Word count above maximum",
			guardrail: &WordCountGuardrail{
				Min:      1,
				Max:      2,
				JSONPath: "$.content",
				Inverted: false,
				logger:   createMockLogger(),
			},
			payload:   []byte(`{"content": "hello world test more"}`),
			expected:  false,
			expectErr: false,
		},
		{
			name: "Inverted logic - word count within range (should fail)",
			guardrail: &WordCountGuardrail{
				Min:      2,
				Max:      5,
				JSONPath: "$.content",
				Inverted: true,
				logger:   createMockLogger(),
			},
			payload:   []byte(`{"content": "hello world test"}`),
			expected:  false,
			expectErr: false,
		},
		{
			name: "Inverted logic - word count outside range (should pass)",
			guardrail: &WordCountGuardrail{
				Min:      5,
				Max:      10,
				JSONPath: "$.content",
				Inverted: true,
				logger:   createMockLogger(),
			},
			payload:   []byte(`{"content": "hello"}`),
			expected:  true,
			expectErr: false,
		},
		{
			name: "Invalid JSON path",
			guardrail: &WordCountGuardrail{
				Min:      1,
				Max:      5,
				JSONPath: "$.nonexistent",
				Inverted: false,
				logger:   createMockLogger(),
			},
			payload:   []byte(`{"content": "hello world"}`),
			expected:  false,
			expectErr: true,
		},
		{
			name: "Invalid JSON payload",
			guardrail: &WordCountGuardrail{
				Min:      1,
				Max:      5,
				JSONPath: "$.content",
				Inverted: false,
				logger:   createMockLogger(),
			},
			payload:   []byte(`invalid json`),
			expected:  false,
			expectErr: true,
		},
		{
			name: "Invalid word count range",
			guardrail: &WordCountGuardrail{
				Min:      10,
				Max:      5, // Max < Min
				JSONPath: "$.content",
				Inverted: false,
				logger:   createMockLogger(),
			},
			payload:   []byte(`{"content": "hello world"}`),
			expected:  false,
			expectErr: false,
		},
		{
			name: "Empty content",
			guardrail: &WordCountGuardrail{
				Min:      0,
				Max:      5,
				JSONPath: "$.content",
				Inverted: false,
				logger:   createMockLogger(),
			},
			payload:   []byte(`{"content": ""}`),
			expected:  true,
			expectErr: false,
		},
		{
			name: "Content with punctuation and extra spaces",
			guardrail: &WordCountGuardrail{
				Min:      2,
				Max:      5,
				JSONPath: "$.content",
				Inverted: false,
				logger:   createMockLogger(),
			},
			payload:   []byte(`{"content": "Hello,  world!   Test."}`),
			expected:  true,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.guardrail.validatePayload(tt.payload, tt.isResponse)

			if tt.expectErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Expected result %t, got %t", tt.expected, result)
			}
		})
	}
}

func TestWordCountGuardrail_decompressLLMResp(t *testing.T) {
	guardrail := &WordCountGuardrail{
		logger: createMockLogger(),
	}

	tests := []struct {
		name     string
		payload  []byte
		expected string
	}{
		{
			name:     "Uncompressed content",
			payload:  []byte("hello world"),
			expected: "hello world",
		},
		{
			name:     "Gzipped content",
			payload:  createGzippedContent("hello compressed world"),
			expected: "hello compressed world",
		},
		{
			name:     "Empty payload",
			payload:  []byte{},
			expected: "",
		},
		{
			name:     "Single byte payload",
			payload:  []byte{0x41}, // 'A'
			expected: "A",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := guardrail.decompressLLMResp(tt.payload)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestWordCountGuardrail_extractStringValueFromJsonpath(t *testing.T) {
	guardrail := &WordCountGuardrail{
		logger: createMockLogger(),
	}

	tests := []struct {
		name      string
		payload   []byte
		jsonPath  string
		expected  string
		expectErr bool
	}{
		{
			name:      "Simple field extraction",
			payload:   []byte(`{"content": "hello world"}`),
			jsonPath:  "$.content",
			expected:  "hello world",
			expectErr: false,
		},
		{
			name:      "Nested field extraction",
			payload:   []byte(`{"data": {"message": "nested content"}}`),
			jsonPath:  "$.data.message",
			expected:  "nested content",
			expectErr: false,
		},
		{
			name:      "Non-existent field",
			payload:   []byte(`{"content": "hello world"}`),
			jsonPath:  "$.missing",
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
		{
			name:      "Array element extraction",
			payload:   []byte(`{"items": ["first", "second", "third"]}`),
			jsonPath:  "$.items.0",
			expected:  "first",
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := guardrail.extractStringValueFromJsonpath(tt.payload, tt.jsonPath)

			if tt.expectErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestWordCountGuardrail_buildErrorResponse(t *testing.T) {
	guardrail := &WordCountGuardrail{
		Name:           "TestGuardrail",
		ShowAssessment: true,
		Min:            5,
		Max:            10,
		logger:         createMockLogger(),
	}

	tests := []struct {
		name            string
		isResponse      bool
		validationError error
		expectedCode    int
	}{
		{
			name:            "Request validation error",
			isResponse:      false,
			validationError: nil,
			expectedCode:    GuardrailErrorCode,
		},
		{
			name:            "Response validation error",
			isResponse:      true,
			validationError: nil,
			expectedCode:    GuardrailErrorCode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := guardrail.buildErrorResponse(tt.isResponse, tt.validationError)

			if !result.ImmediateResponse {
				t.Error("Expected ImmediateResponse to be true")
			}
			if int(result.ImmediateResponseCode) != tt.expectedCode {
				t.Errorf("Expected response code %d, got %d", tt.expectedCode, int(result.ImmediateResponseCode))
			}
			if result.ImmediateResponseContentType != "application/json" {
				t.Errorf("Expected content type 'application/json', got '%s'", result.ImmediateResponseContentType)
			}
			if !result.StopFurtherProcessing {
				t.Error("Expected StopFurtherProcessing to be true")
			}

			// Parse the response body to validate structure
			var responseBody map[string]interface{}
			err := json.Unmarshal([]byte(result.ImmediateResponseBody), &responseBody)
			if err != nil {
				t.Errorf("Failed to parse response body as JSON: %v", err)
			}

			if responseBody[ErrorCode] != GuardrailAPIMExceptionCode {
				t.Errorf("Expected error code %s, got %v", GuardrailAPIMExceptionCode, responseBody[ErrorCode])
			}
			if responseBody[ErrorType] != WordCountGuardrailConstant {
				t.Errorf("Expected error type %s, got %v", WordCountGuardrailConstant, responseBody[ErrorType])
			}
		})
	}
}

func TestWordCountGuardrail_buildAssessmentObject(t *testing.T) {
	tests := []struct {
		name            string
		guardrail       *WordCountGuardrail
		isResponse      bool
		validationError error
		expectedFields  map[string]interface{}
	}{
		{
			name: "Request assessment without error",
			guardrail: &WordCountGuardrail{
				Name:           "TestGuardrail",
				ShowAssessment: true,
				Min:            5,
				Max:            10,
				Inverted:       false,
				logger:         createMockLogger(),
			},
			isResponse:      false,
			validationError: nil,
			expectedFields: map[string]interface{}{
				AssessmentAction:     "GUARDRAIL_INTERVENED",
				InterveningGuardrail: "TestGuardrail",
				Direction:            "REQUEST",
				AssessmentReason:     "Violation of applied word count constraints detected.",
			},
		},
		{
			name: "Response assessment with inverted logic",
			guardrail: &WordCountGuardrail{
				Name:           "TestGuardrail",
				ShowAssessment: true,
				Min:            5,
				Max:            10,
				Inverted:       true,
				logger:         createMockLogger(),
			},
			isResponse:      true,
			validationError: nil,
			expectedFields: map[string]interface{}{
				AssessmentAction:     "GUARDRAIL_INTERVENED",
				InterveningGuardrail: "TestGuardrail",
				Direction:            "RESPONSE",
				AssessmentReason:     "Violation of applied word count constraints detected.",
			},
		},
		{
			name: "Assessment with JSONPath error",
			guardrail: &WordCountGuardrail{
				Name:           "TestGuardrail",
				ShowAssessment: true,
				JSONPath:       "$.nonexistent",
				logger:         createMockLogger(),
			},
			isResponse:      false,
			validationError: fmt.Errorf("field not found"),
			expectedFields: map[string]interface{}{
				AssessmentAction:     "GUARDRAIL_INTERVENED",
				InterveningGuardrail: "TestGuardrail",
				Direction:            "REQUEST",
				AssessmentReason:     "Error extracting content from payload using JSONPath.",
			},
		},
		{
			name: "Assessment without showing details",
			guardrail: &WordCountGuardrail{
				Name:           "TestGuardrail",
				ShowAssessment: false,
				Min:            5,
				Max:            10,
				logger:         createMockLogger(),
			},
			isResponse:      false,
			validationError: nil,
			expectedFields: map[string]interface{}{
				AssessmentAction:     "GUARDRAIL_INTERVENED",
				InterveningGuardrail: "TestGuardrail",
				Direction:            "REQUEST",
				AssessmentReason:     "Violation of applied word count constraints detected.",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.guardrail.buildAssessmentObject(tt.isResponse, tt.validationError)

			for key, expected := range tt.expectedFields {
				if result[key] != expected {
					t.Errorf("Expected %s to be %v, got %v", key, expected, result[key])
				}
			}

			// Check that assessments field is present when ShowAssessment is true
			if tt.guardrail.ShowAssessment {
				if _, exists := result[Assessments]; !exists {
					t.Error("Expected assessments field to be present when ShowAssessment is true")
				}
			}
		})
	}
}

// Benchmark tests
func BenchmarkWordCountGuardrail_validatePayload(b *testing.B) {
	guardrail := &WordCountGuardrail{
		Min:      10,
		Max:      50,
		JSONPath: "$.content",
		Inverted: false,
		logger:   createMockLogger(),
	}

	payload := []byte(`{"content": "This is a test payload with multiple words to benchmark the validation performance of the word count guardrail implementation"}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = guardrail.validatePayload(payload, false)
	}
}

func BenchmarkWordCountGuardrail_Process(b *testing.B) {
	guardrail := &WordCountGuardrail{
		Name:     "BenchmarkGuardrail",
		Min:      10,
		Max:      50,
		JSONPath: "$.content",
		Inverted: false,
		logger:   createMockLogger(),
	}

	requestConfig := &requestconfig.Holder{
		ProcessingPhase: requestconfig.ProcessingPhaseRequestBody,
		RequestBody: &envoy_service_proc_v3.HttpBody{
			Body: []byte(`{"content": "This is a test payload with multiple words to benchmark the validation performance"}`),
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = guardrail.Process(requestConfig)
	}
}

// Additional edge case tests
func TestWordCountGuardrail_EdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		guardrail *WordCountGuardrail
		payload   []byte
		expected  bool
		expectErr bool
	}{
		{
			name: "Zero word count with min=0",
			guardrail: &WordCountGuardrail{
				Min:      0,
				Max:      5,
				JSONPath: "$.content",
				Inverted: false,
				logger:   createMockLogger(),
			},
			payload:   []byte(`{"content": ""}`),
			expected:  true,
			expectErr: false,
		},
		{
			name: "Negative min value",
			guardrail: &WordCountGuardrail{
				Min:      -1,
				Max:      5,
				JSONPath: "$.content",
				Inverted: false,
				logger:   createMockLogger(),
			},
			payload:   []byte(`{"content": "hello world"}`),
			expected:  false,
			expectErr: false,
		},
		{
			name: "Zero max value",
			guardrail: &WordCountGuardrail{
				Min:      0,
				Max:      0,
				JSONPath: "$.content",
				Inverted: false,
				logger:   createMockLogger(),
			},
			payload:   []byte(`{"content": "hello"}`),
			expected:  false,
			expectErr: false,
		},
		{
			name: "Very large word count",
			guardrail: &WordCountGuardrail{
				Min:      1,
				Max:      1000,
				JSONPath: "$.content",
				Inverted: false,
				logger:   createMockLogger(),
			},
			payload:   []byte(`{"content": "` + strings.Repeat("word ", 500) + `"}`),
			expected:  true,
			expectErr: false,
		},
		{
			name: "Content with only whitespace",
			guardrail: &WordCountGuardrail{
				Min:      0,
				Max:      5,
				JSONPath: "$.content",
				Inverted: false,
				logger:   createMockLogger(),
			},
			payload:   []byte(`{"content": "   \t  \n  "}`),
			expected:  true,
			expectErr: false,
		},
		{
			name: "Content with Unicode characters",
			guardrail: &WordCountGuardrail{
				Min:      2,
				Max:      5,
				JSONPath: "$.content",
				Inverted: false,
				logger:   createMockLogger(),
			},
			payload:   []byte(`{"content": "Hello 世界 Testing"}`),
			expected:  true,
			expectErr: false,
		},
		{
			name: "Complex JSON path",
			guardrail: &WordCountGuardrail{
				Min:      2,
				Max:      5,
				JSONPath: "$.data.messages.0.text",
				Inverted: false,
				logger:   createMockLogger(),
			},
			payload:   []byte(`{"data": {"messages": [{"text": "hello world"}]}}`),
			expected:  true,
			expectErr: false,
		},
		{
			name: "Numeric content converted to string",
			guardrail: &WordCountGuardrail{
				Min:      1,
				Max:      5,
				JSONPath: "$.value",
				Inverted: false,
				logger:   createMockLogger(),
			},
			payload:   []byte(`{"value": 12345}`),
			expected:  true,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.guardrail.validatePayload(tt.payload, false)

			if tt.expectErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Expected result %t, got %t", tt.expected, result)
			}
		})
	}
}

// Test word counting logic specifically
func TestWordCountGuardrail_WordCounting(t *testing.T) {
	guardrail := &WordCountGuardrail{
		Min:      1,
		Max:      10,
		JSONPath: "$.content",
		Inverted: false,
		logger:   createMockLogger(),
	}

	tests := []struct {
		name          string
		content       string
		expectedWords int
	}{
		{
			name:          "Simple words",
			content:       "hello world",
			expectedWords: 2,
		},
		{
			name:          "Words with punctuation",
			content:       "Hello, world! How are you?",
			expectedWords: 5,
		},
		{
			name:          "Multiple spaces",
			content:       "hello    world     test",
			expectedWords: 3,
		},
		{
			name:          "Leading and trailing spaces",
			content:       "   hello world   ",
			expectedWords: 2,
		},
		{
			name:          "Mixed punctuation and spaces",
			content:       "Well, hello there! How's it going?",
			expectedWords: 6,
		},
		{
			name:          "Single word",
			content:       "hello",
			expectedWords: 1,
		},
		{
			name:          "Empty string",
			content:       "",
			expectedWords: 0,
		},
		{
			name:          "Only punctuation",
			content:       "!@#$%^&*()",
			expectedWords: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := []byte(fmt.Sprintf(`{"content": "%s"}`, tt.content))
			
			// Simulate the word counting logic
			extractedValue, err := guardrail.extractStringValueFromJsonpath(payload, "$.content")
			if err != nil {
				t.Errorf("Failed to extract content: %v", err)
				return
			}

			// Apply the same cleaning and counting logic as in validatePayload
			cleanedValue := TextCleanRegexCompiled.ReplaceAllString(extractedValue, "")
			cleanedValue = strings.TrimSpace(cleanedValue)

			words := WordSplitRegexCompiled.Split(cleanedValue, -1)
			wordCount := 0
			for _, word := range words {
				if word != "" {
					wordCount++
				}
			}

			if wordCount != tt.expectedWords {
				t.Errorf("Expected %d words, got %d for content: '%s'", tt.expectedWords, wordCount, tt.content)
			}
		})
	}
}

// Test gzip compression/decompression specifically
func TestWordCountGuardrail_CompressionHandling(t *testing.T) {
	guardrail := &WordCountGuardrail{
		Min:      2,
		Max:      5,
		JSONPath: "$.content",
		Inverted: false,
		logger:   createMockLogger(),
	}

	originalContent := `{"content": "hello world test"}`
	compressedContent := createGzippedContent(originalContent)

	tests := []struct {
		name       string
		payload    []byte
		isResponse bool
		expected   bool
	}{
		{
			name:       "Uncompressed response",
			payload:    []byte(originalContent),
			isResponse: true,
			expected:   true,
		},
		{
			name:       "Compressed response",
			payload:    compressedContent,
			isResponse: true,
			expected:   true,
		},
		{
			name:       "Uncompressed request (no decompression)",
			payload:    []byte(originalContent),
			isResponse: false,
			expected:   true,
		},
		{
			name:       "Compressed request (no decompression applied)",
			payload:    compressedContent,
			isResponse: false,
			expected:   false, // Should fail because it tries to parse compressed data as JSON
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := guardrail.validatePayload(tt.payload, tt.isResponse)
			if result != tt.expected {
				t.Errorf("Expected result %t, got %t", tt.expected, result)
			}
		})
	}
}

// Test concurrent access safety
func TestWordCountGuardrail_ConcurrentAccess(t *testing.T) {
	guardrail := &WordCountGuardrail{
		Name:     "ConcurrentTestGuardrail",
		Min:      2,
		Max:      5,
		JSONPath: "$.content",
		Inverted: false,
		logger:   createMockLogger(),
	}

	requestConfig := &requestconfig.Holder{
		ProcessingPhase: requestconfig.ProcessingPhaseRequestBody,
		RequestBody: &envoy_service_proc_v3.HttpBody{
			Body: []byte(`{"content": "hello world test"}`),
		},
	}

	// Run multiple goroutines to test concurrent access
	const numGoroutines = 100
	results := make(chan bool, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			result := guardrail.Process(requestConfig)
			results <- !result.ImmediateResponse // Should be true for valid input
		}()
	}

	// Collect all results
	for i := 0; i < numGoroutines; i++ {
		result := <-results
		if !result {
			t.Error("Expected successful processing in concurrent test")
		}
	}
}
