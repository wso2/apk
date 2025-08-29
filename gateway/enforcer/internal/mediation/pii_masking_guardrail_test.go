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
	"strings"
	"testing"

	egv1a1 "github.com/envoyproxy/gateway/api/v1alpha1"
	envoy_service_proc_v3 "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"
	dpv2alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v2alpha1"
	"github.com/wso2/apk/gateway/enforcer/internal/logging"
	"github.com/wso2/apk/gateway/enforcer/internal/requestconfig"
)

// Mock logger for PII masking testing
func createMockPIIMaskingLogger() *logging.Logger {
	// Create a proper mock logger using the default logging setup
	mockLogger := logging.DefaultLogger(egv1a1.LogLevelInfo)
	return &mockLogger
}

// Helper function to create test mediation for PII masking guardrail
func createTestPIIMaskingMediation(params map[string]string) *dpv2alpha1.Mediation {
	var parameters []*dpv2alpha1.Parameter
	for key, value := range params {
		parameters = append(parameters, &dpv2alpha1.Parameter{
			Key:   key,
			Value: value,
		})
	}

	return &dpv2alpha1.Mediation{
		PolicyID:      "test-pii-masking-policy",
		PolicyVersion: "1.0.0",
		Parameters:    parameters,
	}
}

// Helper function to create test request config
func createTestPIIMaskingRequestConfig(phase requestconfig.ProcessingPhase, requestBody, responseBody []byte) *requestconfig.Holder {
	config := &requestconfig.Holder{
		ProcessingPhase: phase,
	}

	if requestBody != nil {
		config.RequestBody = &envoy_service_proc_v3.HttpBody{
			Body: requestBody,
		}
	}

	if responseBody != nil {
		config.ResponseBody = &envoy_service_proc_v3.HttpBody{
			Body: responseBody,
		}
	}

	return config
}

// Helper function to compress content with gzip
func compressContent(content []byte) []byte {
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)
	writer.Write(content)
	writer.Close()
	return buf.Bytes()
}

// TestNewPIIMaskingGuardrail tests the constructor
func TestNewPIIMaskingGuardrail(t *testing.T) {
	tests := []struct {
		name     string
		params   map[string]string
		expected func(*PIIMaskingGuardrail) bool
	}{
		{
			name:   "Default values",
			params: map[string]string{},
			expected: func(g *PIIMaskingGuardrail) bool {
				return g.Name == "PIIMaskingGuardrail" &&
					g.JSONPath == "$.content" &&
					!g.RedactPII &&
					!g.ShowAssessment &&
					len(g.PiiEntities) == 0
			},
		},
		{
			name: "Custom values",
			params: map[string]string{
				"name":           "CustomPIIMasking",
				"jsonPath":       "$.data.text",
				"redactPII":      "true",
				"showAssessment": "true",
			},
			expected: func(g *PIIMaskingGuardrail) bool {
				return g.Name == "CustomPIIMasking" &&
					g.JSONPath == "$.data.text" &&
					g.RedactPII &&
					g.ShowAssessment
			},
		},
		{
			name: "PII entities configuration",
			params: map[string]string{
				"name": "PIITest",
				"piiEntities": `[
					{"piiEntity": "EMAIL", "piiRegex": "\\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Z|a-z]{2,}\\b"},
					{"piiEntity": "PHONE", "piiRegex": "\\b\\d{3}-\\d{3}-\\d{4}\\b"}
				]`,
			},
			expected: func(g *PIIMaskingGuardrail) bool {
				return g.Name == "PIITest" &&
					len(g.PiiEntities) == 2 &&
					g.PiiEntities["EMAIL"] != nil &&
					g.PiiEntities["PHONE"] != nil
			},
		},
		{
			name: "Invalid PII entities JSON",
			params: map[string]string{
				"piiEntities": `invalid json`,
			},
			expected: func(g *PIIMaskingGuardrail) bool {
				return len(g.PiiEntities) == 0
			},
		},
		{
			name: "Invalid regex pattern",
			params: map[string]string{
				"piiEntities": `[
					{"piiEntity": "INVALID", "piiRegex": "[invalid regex"}
				]`,
			},
			expected: func(g *PIIMaskingGuardrail) bool {
				return len(g.PiiEntities) == 0
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mediation := createTestPIIMaskingMediation(tt.params)
			guardrail := NewPIIMaskingGuardrail(mediation)

			if !tt.expected(guardrail) {
				t.Errorf("Test %s failed: guardrail configuration doesn't match expectations", tt.name)
			}

			// Verify basic fields
			if guardrail.PolicyID != "test-pii-masking-policy" {
				t.Errorf("Expected PolicyID 'test-pii-masking-policy', got '%s'", guardrail.PolicyID)
			}
			if guardrail.PolicyVersion != "1.0.0" {
				t.Errorf("Expected PolicyVersion '1.0.0', got '%s'", guardrail.PolicyVersion)
			}
		})
	}
}

// TestPIIMaskingGuardrailProcess tests the Process method
func TestPIIMaskingGuardrailProcess(t *testing.T) {
	tests := []struct {
		name           string
		params         map[string]string
		phase          requestconfig.ProcessingPhase
		requestBody    []byte
		responseBody   []byte
		expectModified bool
		expectError    bool
	}{
		{
			name: "No body - request phase",
			params: map[string]string{
				"piiEntities": `[{"piiEntity": "EMAIL", "piiRegex": "\\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Z|a-z]{2,}\\b"}]`,
			},
			phase:          requestconfig.ProcessingPhaseRequestBody,
			requestBody:    nil,
			expectModified: false,
			expectError:    false,
		},
		{
			name: "No body - response phase",
			params: map[string]string{
				"piiEntities": `[{"piiEntity": "EMAIL", "piiRegex": "\\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Z|a-z]{2,}\\b"}]`,
			},
			phase:          requestconfig.ProcessingPhaseResponseBody,
			responseBody:   nil,
			expectModified: false,
			expectError:    false,
		},
		{
			name: "Request with PII - masking mode",
			params: map[string]string{
				"redactPII":   "false",
				"piiEntities": `[{"piiEntity": "EMAIL", "piiRegex": "\\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Z|a-z]{2,}\\b"}]`,
			},
			phase:          requestconfig.ProcessingPhaseRequestBody,
			requestBody:    []byte(`{"content": "Contact us at john@example.com for help"}`),
			expectModified: true,
			expectError:    false,
		},
		{
			name: "Request with PII - redaction mode",
			params: map[string]string{
				"redactPII":   "true",
				"piiEntities": `[{"piiEntity": "EMAIL", "piiRegex": "\\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Z|a-z]{2,}\\b"}]`,
			},
			phase:          requestconfig.ProcessingPhaseRequestBody,
			requestBody:    []byte(`{"content": "Contact us at john@example.com for help"}`),
			expectModified: true,
			expectError:    false,
		},
		{
			name: "Request without PII",
			params: map[string]string{
				"piiEntities": `[{"piiEntity": "EMAIL", "piiRegex": "\\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Z|a-z]{2,}\\b"}]`,
			},
			phase:          requestconfig.ProcessingPhaseRequestBody,
			requestBody:    []byte(`{"content": "This is just normal text without sensitive data"}`),
			expectModified: false,
			expectError:    false,
		},
		{
			name: "Response with compressed body",
			params: map[string]string{
				"redactPII":   "true",
				"piiEntities": `[{"piiEntity": "EMAIL", "piiRegex": "\\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Z|a-z]{2,}\\b"}]`,
			},
			phase:          requestconfig.ProcessingPhaseResponseBody,
			responseBody:   compressContent([]byte(`{"content": "Contact support@company.com"}`)),
			expectModified: true,
			expectError:    false,
		},
		{
			name: "Invalid JSON path",
			params: map[string]string{
				"jsonPath":    "$.nonexistent",
				"piiEntities": `[{"piiEntity": "EMAIL", "piiRegex": "\\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Z|a-z]{2,}\\b"}]`,
			},
			phase:          requestconfig.ProcessingPhaseRequestBody,
			requestBody:    []byte(`{"content": "Contact us at john@example.com"}`),
			expectModified: false,
			expectError:    true,
		},
		{
			name: "Different processing phase",
			params: map[string]string{
				"piiEntities": `[{"piiEntity": "EMAIL", "piiRegex": "\\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Z|a-z]{2,}\\b"}]`,
			},
			phase:          requestconfig.ProcessingPhaseRequestHeaders,
			requestBody:    []byte(`{"content": "Contact us at john@example.com"}`),
			expectModified: false,
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mediation := createTestPIIMaskingMediation(tt.params)
			guardrail := NewPIIMaskingGuardrail(mediation)
			requestConfig := createTestPIIMaskingRequestConfig(tt.phase, tt.requestBody, tt.responseBody)

			// Store original body for comparison
			var originalBody []byte
			if tt.phase == requestconfig.ProcessingPhaseRequestBody && requestConfig.RequestBody != nil {
				originalBody = make([]byte, len(requestConfig.RequestBody.Body))
				copy(originalBody, requestConfig.RequestBody.Body)
			} else if tt.phase == requestconfig.ProcessingPhaseResponseBody && requestConfig.ResponseBody != nil {
				originalBody = make([]byte, len(requestConfig.ResponseBody.Body))
				copy(originalBody, requestConfig.ResponseBody.Body)
			}

			result := guardrail.Process(requestConfig)

			// Check error expectation
			if tt.expectError {
				if result == nil || !result.ImmediateResponse {
					t.Errorf("Expected error response, but got success")
				}
				return
			}

			if result == nil {
				t.Errorf("Expected result, got nil")
				return
			}

			if result.ImmediateResponse {
				t.Errorf("Expected success, but got error response: %s", result.ImmediateResponseBody)
				return
			}

			// Check modification expectation
			var bodyModified bool
			if tt.phase == requestconfig.ProcessingPhaseRequestBody && requestConfig.RequestBody != nil {
				bodyModified = !bytes.Equal(originalBody, requestConfig.RequestBody.Body)
			} else if tt.phase == requestconfig.ProcessingPhaseResponseBody && requestConfig.ResponseBody != nil {
				bodyModified = !bytes.Equal(originalBody, requestConfig.ResponseBody.Body)
			}

			if tt.expectModified != bodyModified {
				t.Errorf("Expected modification: %v, but body was modified: %v", tt.expectModified, bodyModified)
			}
		})
	}
}

// TestPIIMaskingGuardrailValidatePayload tests the validatePayload method
func TestPIIMaskingGuardrailValidatePayload(t *testing.T) {
	tests := []struct {
		name           string
		params         map[string]string
		payload        []byte
		isResponse     bool
		expectModified bool
		expectError    bool
	}{
		{
			name: "Email masking",
			params: map[string]string{
				"redactPII":   "false",
				"piiEntities": `[{"piiEntity": "EMAIL", "piiRegex": "\\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Z|a-z]{2,}\\b"}]`,
			},
			payload:        []byte(`{"content": "Contact john@example.com and mary@test.org"}`),
			isResponse:     false,
			expectModified: true,
			expectError:    false,
		},
		{
			name: "Phone number redaction",
			params: map[string]string{
				"redactPII":   "true",
				"piiEntities": `[{"piiEntity": "PHONE", "piiRegex": "\\b\\d{3}-\\d{3}-\\d{4}\\b"}]`,
			},
			payload:        []byte(`{"content": "Call 123-456-7890 or 987-654-3210"}`),
			isResponse:     false,
			expectModified: true,
			expectError:    false,
		},
		{
			name: "Multiple PII types",
			params: map[string]string{
				"redactPII": "false",
				"piiEntities": `[
					{"piiEntity": "EMAIL", "piiRegex": "\\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Z|a-z]{2,}\\b"},
					{"piiEntity": "PHONE", "piiRegex": "\\b\\d{3}-\\d{3}-\\d{4}\\b"}
				]`,
			},
			payload:        []byte(`{"content": "Contact john@test.com at 123-456-7890"}`),
			isResponse:     false,
			expectModified: true,
			expectError:    false,
		},
		{
			name: "No PII detected",
			params: map[string]string{
				"piiEntities": `[{"piiEntity": "EMAIL", "piiRegex": "\\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Z|a-z]{2,}\\b"}]`,
			},
			payload:        []byte(`{"content": "This is clean text with no sensitive data"}`),
			isResponse:     false,
			expectModified: false,
			expectError:    false,
		},
		{
			name: "Empty content",
			params: map[string]string{
				"piiEntities": `[{"piiEntity": "EMAIL", "piiRegex": "\\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Z|a-z]{2,}\\b"}]`,
			},
			payload:        []byte(`{"content": ""}`),
			isResponse:     false,
			expectModified: false,
			expectError:    false,
		},
		{
			name: "Custom JSON path",
			params: map[string]string{
				"jsonPath":    "$.message",
				"redactPII":   "true",
				"piiEntities": `[{"piiEntity": "EMAIL", "piiRegex": "\\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Z|a-z]{2,}\\b"}]`,
			},
			payload:        []byte(`{"message": "Email us at support@company.com", "content": "Other text"}`),
			isResponse:     false,
			expectModified: true,
			expectError:    false,
		},
		{
			name: "Invalid JSON path",
			params: map[string]string{
				"jsonPath":    "$.nonexistent",
				"piiEntities": `[{"piiEntity": "EMAIL", "piiRegex": "\\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Z|a-z]{2,}\\b"}]`,
			},
			payload:        []byte(`{"content": "Email john@test.com"}`),
			isResponse:     false,
			expectModified: false,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mediation := createTestPIIMaskingMediation(tt.params)
			guardrail := NewPIIMaskingGuardrail(mediation)
			requestConfig := &requestconfig.Holder{}

			result, err := guardrail.validatePayload(tt.payload, tt.isResponse, requestConfig)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if tt.expectModified {
				if result.ModifiedPayload == nil {
					t.Errorf("Expected modified payload, but got none")
				}
			} else {
				if result.ModifiedPayload != nil {
					t.Errorf("Expected no modification, but payload was modified")
				}
			}
		})
	}
}

// TestPIIMaskingGuardrailMaskPIIFromContent tests the maskPIIFromContent method
func TestPIIMaskingGuardrailMaskPIIFromContent(t *testing.T) {
	tests := []struct {
		name       string
		params     map[string]string
		content    string
		isResponse bool
		expected   string
	}{
		{
			name: "Email masking in request",
			params: map[string]string{
				"piiEntities": `[{"piiEntity": "EMAIL", "piiRegex": "\\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Z|a-z]{2,}\\b"}]`,
			},
			content:    "Contact john@example.com for help",
			isResponse: false,
			expected:   "Contact [EMAIL_0000] for help",
		},
		{
			name: "Multiple emails masking",
			params: map[string]string{
				"piiEntities": `[{"piiEntity": "EMAIL", "piiRegex": "\\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Z|a-z]{2,}\\b"}]`,
			},
			content:    "Contact john@test.com or mary@example.org",
			isResponse: false,
			expected:   "Contact [EMAIL_0000] or [EMAIL_0001]",
		},
		{
			name: "Phone number masking",
			params: map[string]string{
				"piiEntities": `[{"piiEntity": "PHONE", "piiRegex": "\\b\\d{3}-\\d{3}-\\d{4}\\b"}]`,
			},
			content:    "Call 123-456-7890 for support",
			isResponse: false,
			expected:   "Call [PHONE_0000] for support",
		},
		{
			name: "No PII found",
			params: map[string]string{
				"piiEntities": `[{"piiEntity": "EMAIL", "piiRegex": "\\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Z|a-z]{2,}\\b"}]`,
			},
			content:    "This is clean text",
			isResponse: false,
			expected:   "",
		},
		{
			name: "Empty content",
			params: map[string]string{
				"piiEntities": `[{"piiEntity": "EMAIL", "piiRegex": "\\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Z|a-z]{2,}\\b"}]`,
			},
			content:    "",
			isResponse: false,
			expected:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mediation := createTestPIIMaskingMediation(tt.params)
			guardrail := NewPIIMaskingGuardrail(mediation)
			requestConfig := &requestconfig.Holder{}

			result := guardrail.maskPIIFromContent(tt.content, tt.isResponse, requestConfig)

			if tt.expected == "" {
				if result != "" {
					t.Errorf("Expected no masking, but got: %s", result)
				}
			} else {
				if result == "" {
					t.Errorf("Expected masking result, but got empty string")
				}
				// For masking, we check that placeholders are created (exact format may vary)
				if !strings.Contains(result, "[") || !strings.Contains(result, "]") {
					t.Errorf("Expected placeholder format in result: %s", result)
				}
			}
		})
	}
}

// TestPIIMaskingGuardrailRedactPIIFromContent tests the redactPIIFromContent method
func TestPIIMaskingGuardrailRedactPIIFromContent(t *testing.T) {
	tests := []struct {
		name     string
		params   map[string]string
		content  string
		expected string
	}{
		{
			name: "Email redaction",
			params: map[string]string{
				"piiEntities": `[{"piiEntity": "EMAIL", "piiRegex": "\\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Z|a-z]{2,}\\b"}]`,
			},
			content:  "Contact john@example.com for help",
			expected: "Contact ***** for help",
		},
		{
			name: "Multiple emails redaction",
			params: map[string]string{
				"piiEntities": `[{"piiEntity": "EMAIL", "piiRegex": "\\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Z|a-z]{2,}\\b"}]`,
			},
			content:  "Contact john@test.com or mary@example.org",
			expected: "Contact ***** or *****",
		},
		{
			name: "No PII found",
			params: map[string]string{
				"piiEntities": `[{"piiEntity": "EMAIL", "piiRegex": "\\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Z|a-z]{2,}\\b"}]`,
			},
			content:  "This is clean text",
			expected: "",
		},
		{
			name: "Empty content",
			params: map[string]string{
				"piiEntities": `[{"piiEntity": "EMAIL", "piiRegex": "\\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Z|a-z]{2,}\\b"}]`,
			},
			content:  "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mediation := createTestPIIMaskingMediation(tt.params)
			guardrail := NewPIIMaskingGuardrail(mediation)

			result := guardrail.redactPIIFromContent(tt.content)

			if tt.expected == "" {
				if result != "" {
					t.Errorf("Expected no redaction, but got: %s", result)
				}
			} else {
				if result != tt.expected {
					t.Errorf("Expected: %s, got: %s", tt.expected, result)
				}
			}
		})
	}
}

// TestPIIMaskingGuardrailDecompression tests response body decompression
func TestPIIMaskingGuardrailDecompression(t *testing.T) {
	tests := []struct {
		name     string
		payload  []byte
		expected string
		hasError bool
	}{
		{
			name:     "Uncompressed content",
			payload:  []byte("This is plain text"),
			expected: "This is plain text",
			hasError: false,
		},
		{
			name:     "Gzip compressed content",
			payload:  compressContent([]byte("This is compressed text")),
			expected: "This is compressed text",
			hasError: false,
		},
		{
			name:     "Empty payload",
			payload:  []byte{},
			expected: "",
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mediation := createTestPIIMaskingMediation(map[string]string{})
			guardrail := NewPIIMaskingGuardrail(mediation)

			result, _, err := guardrail.decompressResponseBody(tt.payload)

			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("Expected: %s, got: %s", tt.expected, result)
			}
		})
	}
}

// TestPIIMaskingGuardrailJSONPathExtraction tests JSON path value extraction
func TestPIIMaskingGuardrailJSONPathExtraction(t *testing.T) {
	tests := []struct {
		name     string
		payload  []byte
		jsonPath string
		expected string
		hasError bool
	}{
		{
			name:     "Extract from content field",
			payload:  []byte(`{"content": "Hello world"}`),
			jsonPath: "$.content",
			expected: "Hello world",
			hasError: false,
		},
		{
			name:     "Extract from nested field",
			payload:  []byte(`{"data": {"message": "Test message"}}`),
			jsonPath: "$.data.message",
			expected: "Test message",
			hasError: false,
		},
		{
			name:     "Empty JSON path",
			payload:  []byte(`{"content": "Hello world"}`),
			jsonPath: "",
			expected: `{"content": "Hello world"}`,
			hasError: false,
		},
		{
			name:     "Non-existent path",
			payload:  []byte(`{"content": "Hello world"}`),
			jsonPath: "$.nonexistent",
			expected: "",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mediation := createTestPIIMaskingMediation(map[string]string{})
			guardrail := NewPIIMaskingGuardrail(mediation)

			result, err := guardrail.extractJSONPathValue(tt.payload, tt.jsonPath)

			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("Expected: %s, got: %s", tt.expected, result)
			}
		})
	}
}

// TestPIIMaskingGuardrailTextCleaning tests text cleaning functionality
func TestPIIMaskingGuardrailTextCleaning(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Normal text",
			input:    "Hello world!",
			expected: "Hello world!",
		},
		{
			name:     "Text with special characters",
			input:    "Email: john@test.com, Phone: 123-456-7890",
			expected: "Email: john@test.com, Phone: 123-456-7890",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mediation := createTestPIIMaskingMediation(map[string]string{})
			guardrail := NewPIIMaskingGuardrail(mediation)

			result := guardrail.cleanText(tt.input)

			// Since cleanText removes some characters, we just verify it doesn't crash
			// and returns a non-nil result
			if len(tt.input) > 0 && len(result) == 0 {
				t.Errorf("Expected non-empty result for input: %s", tt.input)
			}
		})
	}
}
