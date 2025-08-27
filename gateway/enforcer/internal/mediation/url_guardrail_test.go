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
	"net/http"
	"net/http/httptest"
	"testing"

	envoy_service_proc_v3 "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"
	dpv2alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v2alpha1"
	"github.com/wso2/apk/gateway/enforcer/internal/requestconfig"
)

// Helper function to create test mediation for URLGuardrail
func createTestURLMediation(params map[string]string) *dpv2alpha1.Mediation {
	var parameters []*dpv2alpha1.Parameter
	for key, value := range params {
		parameters = append(parameters, &dpv2alpha1.Parameter{
			Key:   key,
			Value: value,
		})
	}

	return &dpv2alpha1.Mediation{
		PolicyName:    "URLGuardrail",
		PolicyVersion: "v1",
		PolicyID:      "test-url-policy-id",
		Parameters:    parameters,
	}
}

// Helper function to create gzipped content
func createGzippedURLContent(content string) []byte {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	gz.Write([]byte(content))
	gz.Close()
	return buf.Bytes()
}

func TestNewURLGuardrail(t *testing.T) {
	tests := []struct {
		name       string
		parameters map[string]string
		expected   URLGuardrail
	}{
		{
			name: "Default values",
			parameters: map[string]string{},
			expected: URLGuardrail{
				PolicyName:     "URLGuardrail",
				PolicyVersion:  "v1",
				PolicyID:       "test-url-policy-id",
				Name:           "URLGuardrail",
				JSONPath:       "$.content",
				OnlyDNS:        false,
				Timeout:        3000,
				ShowAssessment: false,
			},
		},
		{
			name: "Custom values",
			parameters: map[string]string{
				"name":           "CustomURLGuardrail",
				"jsonPath":       "$.urls",
				"onlyDNS":        "true",
				"timeout":        "5000",
				"showAssessment": "true",
			},
			expected: URLGuardrail{
				PolicyName:     "URLGuardrail",
				PolicyVersion:  "v1",
				PolicyID:       "test-url-policy-id",
				Name:           "CustomURLGuardrail",
				JSONPath:       "$.urls",
				OnlyDNS:        true,
				Timeout:        5000,
				ShowAssessment: true,
			},
		},
		{
			name: "Invalid timeout falls back to default",
			parameters: map[string]string{
				"timeout": "invalid",
			},
			expected: URLGuardrail{
				PolicyName:     "URLGuardrail",
				PolicyVersion:  "v1",
				PolicyID:       "test-url-policy-id",
				Name:           "URLGuardrail",
				JSONPath:       "$.content",
				OnlyDNS:        false,
				Timeout:        3000,
				ShowAssessment: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mediation := createTestURLMediation(tt.parameters)
			result := NewURLGuardrail(mediation)

			if result.PolicyName != tt.expected.PolicyName {
				t.Errorf("Expected PolicyName %s, got %s", tt.expected.PolicyName, result.PolicyName)
			}
			if result.Name != tt.expected.Name {
				t.Errorf("Expected Name %s, got %s", tt.expected.Name, result.Name)
			}
			if result.JSONPath != tt.expected.JSONPath {
				t.Errorf("Expected JSONPath %s, got %s", tt.expected.JSONPath, result.JSONPath)
			}
			if result.OnlyDNS != tt.expected.OnlyDNS {
				t.Errorf("Expected OnlyDNS %t, got %t", tt.expected.OnlyDNS, result.OnlyDNS)
			}
			if result.Timeout != tt.expected.Timeout {
				t.Errorf("Expected Timeout %d, got %d", tt.expected.Timeout, result.Timeout)
			}
			if result.ShowAssessment != tt.expected.ShowAssessment {
				t.Errorf("Expected ShowAssessment %t, got %t", tt.expected.ShowAssessment, result.ShowAssessment)
			}
		})
	}
}

func TestURLGuardrail_Process(t *testing.T) {
	// Create a test HTTP server for testing URL validation
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer testServer.Close()

	tests := []struct {
		name               string
		guardrail          *URLGuardrail
		requestConfig      *requestconfig.Holder
		expectedImmediate  bool
		expectedPassed     bool
	}{
		{
			name: "Request body - valid URL",
			guardrail: &URLGuardrail{
				Name:     "TestURLGuardrail",
				JSONPath: "$.content",
				OnlyDNS:  false,
				Timeout:  5000,
				logger:   createMockLogger(),
			},
			requestConfig: &requestconfig.Holder{
				ProcessingPhase: requestconfig.ProcessingPhaseRequestBody,
				RequestBody: &envoy_service_proc_v3.HttpBody{
					Body: []byte(fmt.Sprintf(`{"content": "Check this URL: %s"}`, testServer.URL)),
				},
			},
			expectedImmediate: false,
			expectedPassed:    true,
		},
		{
			name: "Request body - invalid URL",
			guardrail: &URLGuardrail{
				Name:     "TestURLGuardrail",
				JSONPath: "$.content",
				OnlyDNS:  false,
				Timeout:  1000,
				logger:   createMockLogger(),
			},
			requestConfig: &requestconfig.Holder{
				ProcessingPhase: requestconfig.ProcessingPhaseRequestBody,
				RequestBody: &envoy_service_proc_v3.HttpBody{
					Body: []byte(`{"content": "Check this URL: http://invalid-nonexistent-domain-12345.com"}`),
				},
			},
			expectedImmediate: true,
			expectedPassed:    false,
		},
		{
			name: "Response body - valid URL with DNS check",
			guardrail: &URLGuardrail{
				Name:     "TestURLGuardrail",
				JSONPath: "$.content",
				OnlyDNS:  true,
				Timeout:  5000,
				logger:   createMockLogger(),
			},
			requestConfig: &requestconfig.Holder{
				ProcessingPhase: requestconfig.ProcessingPhaseResponseBody,
				ResponseBody: &envoy_service_proc_v3.HttpBody{
					Body: []byte(`{"content": "Visit https://google.com for more info"}`),
				},
			},
			expectedImmediate: false,
			expectedPassed:    true,
		},
		{
			name: "Response body - invalid URL with DNS check",
			guardrail: &URLGuardrail{
				Name:     "TestURLGuardrail",
				JSONPath: "$.content",
				OnlyDNS:  true,
				Timeout:  1000,
				logger:   createMockLogger(),
			},
			requestConfig: &requestconfig.Holder{
				ProcessingPhase: requestconfig.ProcessingPhaseResponseBody,
				ResponseBody: &envoy_service_proc_v3.HttpBody{
					Body: []byte(`{"content": "Visit https://invalid-nonexistent-domain-12345.com"}`),
				},
			},
			expectedImmediate: true,
			expectedPassed:    false,
		},
		{
			name: "No request body",
			guardrail: &URLGuardrail{
				Name:     "TestURLGuardrail",
				JSONPath: "$.content",
				OnlyDNS:  false,
				Timeout:  3000,
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
			guardrail: &URLGuardrail{
				Name:     "TestURLGuardrail",
				JSONPath: "$.content",
				OnlyDNS:  false,
				Timeout:  3000,
				logger:   createMockLogger(),
			},
			requestConfig: &requestconfig.Holder{
				ProcessingPhase: requestconfig.ProcessingPhaseRequestHeaders,
			},
			expectedImmediate: false,
			expectedPassed:    true,
		},
		{
			name: "No URLs in content",
			guardrail: &URLGuardrail{
				Name:     "TestURLGuardrail",
				JSONPath: "$.content",
				OnlyDNS:  false,
				Timeout:  3000,
				logger:   createMockLogger(),
			},
			requestConfig: &requestconfig.Holder{
				ProcessingPhase: requestconfig.ProcessingPhaseRequestBody,
				RequestBody: &envoy_service_proc_v3.HttpBody{
					Body: []byte(`{"content": "This is just plain text with no URLs"}`),
				},
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

func TestURLGuardrail_validatePayload(t *testing.T) {
	// Create a test HTTP server for testing URL validation
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer testServer.Close()

	tests := []struct {
		name           string
		guardrail      *URLGuardrail
		payload        []byte
		isResponse     bool
		expectedValid  bool
		expectedCount  int
		expectErr      bool
	}{
		{
			name: "Valid URL - HTTP check",
			guardrail: &URLGuardrail{
				JSONPath: "$.content",
				OnlyDNS:  false,
				Timeout:  5000,
				logger:   createMockLogger(),
			},
			payload:       []byte(fmt.Sprintf(`{"content": "Visit %s for info"}`, testServer.URL)),
			expectedValid: true,
			expectedCount: 0,
			expectErr:     false,
		},
		{
			name: "Invalid URL - HTTP check",
			guardrail: &URLGuardrail{
				JSONPath: "$.content",
				OnlyDNS:  false,
				Timeout:  1000,
				logger:   createMockLogger(),
			},
			payload:       []byte(`{"content": "Visit http://invalid-nonexistent-domain-12345.com"}`),
			expectedValid: false,
			expectedCount: 1,
			expectErr:     false,
		},
		{
			name: "Valid URL - DNS check",
			guardrail: &URLGuardrail{
				JSONPath: "$.content",
				OnlyDNS:  true,
				Timeout:  5000,
				logger:   createMockLogger(),
			},
			payload:       []byte(`{"content": "Visit https://google.com for info"}`),
			expectedValid: true,
			expectedCount: 0,
			expectErr:     false,
		},
		{
			name: "Invalid URL - DNS check",
			guardrail: &URLGuardrail{
				JSONPath: "$.content",
				OnlyDNS:  true,
				Timeout:  1000,
				logger:   createMockLogger(),
			},
			payload:       []byte(`{"content": "Visit https://invalid-nonexistent-domain-12345.com"}`),
			expectedValid: false,
			expectedCount: 1,
			expectErr:     false,
		},
		{
			name: "Multiple URLs - mixed validity",
			guardrail: &URLGuardrail{
				JSONPath: "$.content",
				OnlyDNS:  true,
				Timeout:  3000,
				logger:   createMockLogger(),
			},
			payload:       []byte(`{"content": "Visit https://google.com and https://invalid-nonexistent-domain-12345.com"}`),
			expectedValid: false,
			expectedCount: 1, // Only one invalid URL
			expectErr:     false,
		},
		{
			name: "No URLs in content",
			guardrail: &URLGuardrail{
				JSONPath: "$.content",
				OnlyDNS:  false,
				Timeout:  3000,
				logger:   createMockLogger(),
			},
			payload:       []byte(`{"content": "This is just plain text"}`),
			expectedValid: true,
			expectedCount: 0,
			expectErr:     false,
		},
		{
			name: "Invalid JSON path",
			guardrail: &URLGuardrail{
				JSONPath: "$.nonexistent",
				OnlyDNS:  false,
				Timeout:  3000,
				logger:   createMockLogger(),
			},
			payload:       []byte(`{"content": "Visit https://google.com"}`),
			expectedValid: false,
			expectedCount: 0,
			expectErr:     true,
		},
		{
			name: "Invalid JSON payload",
			guardrail: &URLGuardrail{
				JSONPath: "$.content",
				OnlyDNS:  false,
				Timeout:  3000,
				logger:   createMockLogger(),
			},
			payload:       []byte(`invalid json`),
			expectedValid: false,
			expectedCount: 0,
			expectErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, invalidURLs, err := tt.guardrail.validatePayload(tt.payload, tt.isResponse)

			if tt.expectErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if valid != tt.expectedValid {
				t.Errorf("Expected valid %t, got %t", tt.expectedValid, valid)
			}
			if len(invalidURLs) != tt.expectedCount {
				t.Errorf("Expected %d invalid URLs, got %d", tt.expectedCount, len(invalidURLs))
			}
		})
	}
}

func TestURLGuardrail_checkDNS(t *testing.T) {
	guardrail := &URLGuardrail{
		Timeout: 5000,
		logger:  createMockLogger(),
	}

	tests := []struct {
		name     string
		url      string
		expected bool
	}{
		{
			name:     "Valid domain",
			url:      "https://google.com",
			expected: true,
		},
		{
			name:     "Invalid domain",
			url:      "https://invalid-nonexistent-domain-12345.com",
			expected: false,
		},
		{
			name:     "Malformed URL",
			url:      "not-a-url",
			expected: false,
		},
		{
			name:     "URL without scheme",
			url:      "google.com",
			expected: false, // Should fail parsing
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := guardrail.checkDNS(tt.url)
			if result != tt.expected {
				t.Errorf("Expected %t for URL %s, got %t", tt.expected, tt.url, result)
			}
		})
	}
}

func TestURLGuardrail_checkURL(t *testing.T) {
	// Create test servers
	validServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer validServer.Close()

	errorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer errorServer.Close()

	guardrail := &URLGuardrail{
		Timeout: 5000,
		logger:  createMockLogger(),
	}

	tests := []struct {
		name     string
		url      string
		expected bool
	}{
		{
			name:     "Valid reachable URL",
			url:      validServer.URL,
			expected: true,
		},
		{
			name:     "URL returning 404",
			url:      errorServer.URL,
			expected: false,
		},
		{
			name:     "Invalid URL",
			url:      "http://invalid-nonexistent-domain-12345.com",
			expected: false,
		},
		{
			name:     "Malformed URL",
			url:      "not-a-url",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := guardrail.checkURL(tt.url)
			if result != tt.expected {
				t.Errorf("Expected %t for URL %s, got %t", tt.expected, tt.url, result)
			}
		})
	}
}

func TestURLGuardrail_decompressLLMResp(t *testing.T) {
	guardrail := &URLGuardrail{
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
			payload:  createGzippedURLContent("hello compressed world"),
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

func TestURLGuardrail_extractStringValueFromJsonpath(t *testing.T) {
	guardrail := &URLGuardrail{
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
			payload:   []byte(`{"content": "Visit https://google.com"}`),
			jsonPath:  "$.content",
			expected:  "Visit https://google.com",
			expectErr: false,
		},
		{
			name:      "Nested field extraction",
			payload:   []byte(`{"data": {"urls": "https://example.com"}}`),
			jsonPath:  "$.data.urls",
			expected:  "https://example.com",
			expectErr: false,
		},
		{
			name:      "Non-existent field",
			payload:   []byte(`{"content": "hello"}`),
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
			payload:   []byte(`{"urls": ["https://first.com", "https://second.com"]}`),
			jsonPath:  "$.urls.0",
			expected:  "https://first.com",
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

func TestURLGuardrail_buildErrorResponse(t *testing.T) {
	guardrail := &URLGuardrail{
		Name:           "TestURLGuardrail",
		ShowAssessment: true,
		OnlyDNS:        false,
		logger:         createMockLogger(),
	}

	tests := []struct {
		name            string
		isResponse      bool
		invalidURLs     []string
		validationError error
		expectedCode    int
	}{
		{
			name:            "Request validation error",
			isResponse:      false,
			invalidURLs:     []string{"http://invalid.com"},
			validationError: nil,
			expectedCode:    URLGuardrailErrorCode,
		},
		{
			name:            "Response validation error",
			isResponse:      true,
			invalidURLs:     []string{"http://bad1.com", "http://bad2.com"},
			validationError: nil,
			expectedCode:    URLGuardrailErrorCode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := guardrail.buildErrorResponse(tt.isResponse, tt.invalidURLs, tt.validationError)

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

			if responseBody[URLErrorCode] != URLGuardrailAPIMExceptionCode {
				t.Errorf("Expected error code %s, got %v", URLGuardrailAPIMExceptionCode, responseBody[URLErrorCode])
			}
			if responseBody[URLErrorType] != URLGuardrailConstant {
				t.Errorf("Expected error type %s, got %v", URLGuardrailConstant, responseBody[URLErrorType])
			}
		})
	}
}

func TestURLGuardrail_buildAssessmentObject(t *testing.T) {
	tests := []struct {
		name            string
		guardrail       *URLGuardrail
		isResponse      bool
		invalidURLs     []string
		validationError error
		expectedFields  map[string]interface{}
	}{
		{
			name: "Request assessment without error",
			guardrail: &URLGuardrail{
				Name:           "TestURLGuardrail",
				ShowAssessment: true,
				OnlyDNS:        false,
				logger:         createMockLogger(),
			},
			isResponse:      false,
			invalidURLs:     []string{"http://invalid.com"},
			validationError: nil,
			expectedFields: map[string]interface{}{
				URLAssessmentAction:     "GUARDRAIL_INTERVENED",
				URLInterveningGuardrail: "TestURLGuardrail",
				URLDirection:            "REQUEST",
				URLAssessmentReason:     "Violation of URL validity detected.",
			},
		},
		{
			name: "Response assessment with DNS validation",
			guardrail: &URLGuardrail{
				Name:           "TestURLGuardrail",
				ShowAssessment: true,
				OnlyDNS:        true,
				logger:         createMockLogger(),
			},
			isResponse:      true,
			invalidURLs:     []string{"http://bad1.com", "http://bad2.com"},
			validationError: nil,
			expectedFields: map[string]interface{}{
				URLAssessmentAction:     "GUARDRAIL_INTERVENED",
				URLInterveningGuardrail: "TestURLGuardrail",
				URLDirection:            "RESPONSE",
				URLAssessmentReason:     "Violation of URL validity detected.",
			},
		},
		{
			name: "Assessment with JSONPath error",
			guardrail: &URLGuardrail{
				Name:           "TestURLGuardrail",
				ShowAssessment: true,
				JSONPath:       "$.nonexistent",
				logger:         createMockLogger(),
			},
			isResponse:      false,
			invalidURLs:     []string{},
			validationError: fmt.Errorf("field not found"),
			expectedFields: map[string]interface{}{
				URLAssessmentAction:     "GUARDRAIL_INTERVENED",
				URLInterveningGuardrail: "TestURLGuardrail",
				URLDirection:            "REQUEST",
				URLAssessmentReason:     "Error extracting content from payload using JSONPath.",
			},
		},
		{
			name: "Assessment without showing details",
			guardrail: &URLGuardrail{
				Name:           "TestURLGuardrail",
				ShowAssessment: false,
				OnlyDNS:        false,
				logger:         createMockLogger(),
			},
			isResponse:      false,
			invalidURLs:     []string{"http://invalid.com"},
			validationError: nil,
			expectedFields: map[string]interface{}{
				URLAssessmentAction:     "GUARDRAIL_INTERVENED",
				URLInterveningGuardrail: "TestURLGuardrail",
				URLDirection:            "REQUEST",
				URLAssessmentReason:     "Violation of URL validity detected.",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.guardrail.buildAssessmentObject(tt.isResponse, tt.invalidURLs, tt.validationError)

			for key, expected := range tt.expectedFields {
				if result[key] != expected {
					t.Errorf("Expected %s to be %v, got %v", key, expected, result[key])
				}
			}

			// Check that assessments field is present when ShowAssessment is true
			if tt.guardrail.ShowAssessment {
				if _, exists := result[URLAssessments]; !exists {
					t.Error("Expected assessments field to be present when ShowAssessment is true")
				}

				// For URL validation errors, check the structure
				if tt.validationError == nil {
					if assessments, ok := result[URLAssessments].(map[string]interface{}); ok {
						if _, exists := assessments["invalidUrls"]; !exists {
							t.Error("Expected invalidUrls field in assessments")
						}
						if _, exists := assessments["validationType"]; !exists {
							t.Error("Expected validationType field in assessments")
						}
					}
				}
			}
		})
	}
}

// Test URL regex pattern
func TestURLGuardrail_URLRegexPattern(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []string
	}{
		{
			name:     "Single HTTP URL",
			content:  "Visit http://example.com",
			expected: []string{"http://example.com"},
		},
		{
			name:     "Single HTTPS URL",
			content:  "Check https://secure.example.com",
			expected: []string{"https://secure.example.com"},
		},
		{
			name:     "Multiple URLs",
			content:  "Visit http://example.com and https://secure.com",
			expected: []string{"http://example.com", "https://secure.com"},
		},
		{
			name:     "URL with path and query",
			content:  "API endpoint: https://api.example.com/v1/users?page=1",
			expected: []string{"https://api.example.com/v1/users?page=1"},
		},
		{
			name:     "No URLs",
			content:  "This is just plain text",
			expected: []string{},
		},
		{
			name:     "URL with port",
			content:  "Local server: http://localhost:8080",
			expected: []string{"http://localhost:8080"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urls := URLRegexCompiled.FindAllString(tt.content, -1)
			
			if len(urls) != len(tt.expected) {
				t.Errorf("Expected %d URLs, got %d", len(tt.expected), len(urls))
				return
			}

			for i, expectedURL := range tt.expected {
				if i < len(urls) && urls[i] != expectedURL {
					t.Errorf("Expected URL %s, got %s", expectedURL, urls[i])
				}
			}
		})
	}
}

// Benchmark tests
func BenchmarkURLGuardrail_validatePayload(b *testing.B) {
	guardrail := &URLGuardrail{
		JSONPath: "$.content",
		OnlyDNS:  true, // Use DNS validation for faster benchmarking
		Timeout:  1000,
		logger:   createMockLogger(),
	}

	payload := []byte(`{"content": "Visit https://google.com and https://github.com for more info"}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = guardrail.validatePayload(payload, false)
	}
}

func BenchmarkURLGuardrail_Process(b *testing.B) {
	guardrail := &URLGuardrail{
		Name:     "BenchmarkURLGuardrail",
		JSONPath: "$.content",
		OnlyDNS:  true,
		Timeout:  1000,
		logger:   createMockLogger(),
	}

	requestConfig := &requestconfig.Holder{
		ProcessingPhase: requestconfig.ProcessingPhaseRequestBody,
		RequestBody: &envoy_service_proc_v3.HttpBody{
			Body: []byte(`{"content": "Visit https://google.com for more information"}`),
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = guardrail.Process(requestConfig)
	}
}

// Test concurrent access safety
func TestURLGuardrail_ConcurrentAccess(t *testing.T) {
	guardrail := &URLGuardrail{
		Name:     "ConcurrentTestURLGuardrail",
		JSONPath: "$.content",
		OnlyDNS:  true,
		Timeout:  3000,
		logger:   createMockLogger(),
	}

	requestConfig := &requestconfig.Holder{
		ProcessingPhase: requestconfig.ProcessingPhaseRequestBody,
		RequestBody: &envoy_service_proc_v3.HttpBody{
			Body: []byte(`{"content": "Visit https://google.com"}`),
		},
	}

	// Run multiple goroutines to test concurrent access
	const numGoroutines = 50
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
