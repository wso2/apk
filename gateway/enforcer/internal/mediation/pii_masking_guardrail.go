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
	"io"
	"regexp"
	"strings"
	"sync"

	v32 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/tidwall/gjson"
	dpv2alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v2alpha1"
	"github.com/wso2/apk/gateway/enforcer/internal/config"
	"github.com/wso2/apk/gateway/enforcer/internal/logging"
	"github.com/wso2/apk/gateway/enforcer/internal/requestconfig"
)

// PIIMaskingGuardrail represents the configuration for PII Masking Guardrail policy in the API Gateway.
type PIIMaskingGuardrail struct {
	PolicyName     string                    `json:"policyName"`
	PolicyVersion  string                    `json:"policyVersion"`
	PolicyID       string                    `json:"policyID"`
	Name           string                    `json:"name"`
	PiiEntities    map[string]*regexp.Regexp `json:"piiEntities"`
	JSONPath       string                    `json:"jsonPath"`
	RedactPII      bool                      `json:"redactPII"`
	ShowAssessment bool                      `json:"showAssessment"`
	patternMu      sync.RWMutex
	logger         *logging.Logger
	cfg            *config.Server
}

const (
	// PIIMaskingGuardrailPolicyKeyName is the key for specifying the name of the guardrail.
	PIIMaskingGuardrailPolicyKeyName = "name"
	// PIIMaskingGuardrailPolicyKeyPiiEntities is the key for specifying the PII entities and their regex patterns.
	PIIMaskingGuardrailPolicyKeyPiiEntities = "piiEntities"
	// PIIMaskingGuardrailPolicyKeyJSONPath is the key for specifying the JSON path to extract content.
	PIIMaskingGuardrailPolicyKeyJSONPath = "jsonPath"
	// PIIMaskingGuardrailPolicyKeyRedactPII is the key for specifying if PII should be redacted instead of masked.
	PIIMaskingGuardrailPolicyKeyRedactPII = "redactPII"
	// PIIMaskingGuardrailPolicyKeyShowAssessment is the key for specifying if assessment should be shown.
	PIIMaskingGuardrailPolicyKeyShowAssessment = "showAssessment"

	// PIIMaskingGuardrailConstant is the constant for PII masking guardrail errors.
	PIIMaskingGuardrailConstant = "PII_MASKING_GUARDRAIL"
	// PIIMaskingRegexPIIEntitiesKey is the key for storing PII entities in dynamic metadata.
	PIIMaskingRegexPIIEntitiesKey = "piiMaskingRegexPIIEntities"
)

// NewPIIMaskingGuardrail creates a new PIIMaskingGuardrail instance.
func NewPIIMaskingGuardrail(mediation *dpv2alpha1.Mediation) *PIIMaskingGuardrail {
	cfg := config.GetConfig()
	logger := cfg.Logger

	name := "PIIMaskingGuardrail"
	if val, ok := extractPolicyValue(mediation.Parameters, PIIMaskingGuardrailPolicyKeyName); ok {
		name = val
	}

	jsonPath := "$.content"
	if val, ok := extractPolicyValue(mediation.Parameters, PIIMaskingGuardrailPolicyKeyJSONPath); ok {
		jsonPath = val
	}

	redactPII := false
	if val, ok := extractPolicyValue(mediation.Parameters, PIIMaskingGuardrailPolicyKeyRedactPII); ok {
		redactPII = val == "true"
	}

	showAssessment := false
	if val, ok := extractPolicyValue(mediation.Parameters, PIIMaskingGuardrailPolicyKeyShowAssessment); ok {
		showAssessment = val == "true"
	}

	// Parse PII entities
	piiEntities := make(map[string]*regexp.Regexp)
	if val, ok := extractPolicyValue(mediation.Parameters, PIIMaskingGuardrailPolicyKeyPiiEntities); ok {
		// Define struct for the new format
		type PiiEntityConfig struct {
			PiiEntity string `json:"piiEntity"`
			PiiRegex  string `json:"piiRegex"`
		}

		var piiEntitiesArray []PiiEntityConfig
		if err := json.Unmarshal([]byte(val), &piiEntitiesArray); err != nil {
			logger.Sugar().Errorf("Error unmarshaling piiEntities array: %v", err)
		} else {
			// Compile regex patterns during initialization
			for _, entityConfig := range piiEntitiesArray {
				compiledPattern, err := regexp.Compile(entityConfig.PiiRegex)
				if err != nil {
					logger.Sugar().Errorf("Error compiling regex for PII entity '%s': %v", entityConfig.PiiEntity, err)
					continue
				}
				piiEntities[entityConfig.PiiEntity] = compiledPattern
			}
		}
	}

	return &PIIMaskingGuardrail{
		PolicyName:     "PIIMaskingGuardrail",
		PolicyVersion:  mediation.PolicyVersion,
		PolicyID:       mediation.PolicyID,
		Name:           name,
		PiiEntities:    piiEntities,
		JSONPath:       jsonPath,
		RedactPII:      redactPII,
		ShowAssessment: showAssessment,
		logger:         &logger,
		cfg:            cfg,
	}
}

// Process processes the request configuration for PII Masking Guardrail.
func (r *PIIMaskingGuardrail) Process(requestConfig *requestconfig.Holder) *Result {
	result := NewResult()

	// Handle request body processing
	if requestConfig.ProcessingPhase == requestconfig.ProcessingPhaseRequestBody {
		r.logger.Sugar().Debugf("Beginning request payload validation for PIIMaskingGuardrail policy: %s", r.Name)

		if requestConfig.RequestBody == nil || requestConfig.RequestBody.Body == nil {
			r.logger.Sugar().Debug("No request body found, skipping PII masking validation")
			return result
		}

		validationResult, err := r.validatePayload(requestConfig.RequestBody.Body, false, requestConfig)
		if err != nil {
			r.logger.Sugar().Debugf("Request payload validation failed for PIIMaskingGuardrail policy: %s - %v", r.Name, err)
			return r.buildErrorResponse(false, err)
		}

		// Check if payload was modified and update the body
		if validationResult.ModifiedPayload != nil {
			r.logger.Sugar().Debugf("Request payload was modified by PIIMaskingGuardrail policy: %s", r.Name)
			requestConfig.RequestBody.Body = *validationResult.ModifiedPayload
			result.ModifyBody = true
			result.Body = string(*validationResult.ModifiedPayload)
			newBodyLength := len(result.Body)
			r.logger.Sugar().Debug(fmt.Sprintf("new body length: %d\n", newBodyLength))

			result.AddHeaders = map[string]string{
				"Content-Length": fmt.Sprintf("%d", newBodyLength), // Set the new Content-Length
			}
		}

		r.logger.Sugar().Debugf("Request payload validation passed for PIIMaskingGuardrail policy: %s", r.Name)
		return result
	}

	// Handle response body processing
	if requestConfig.ProcessingPhase == requestconfig.ProcessingPhaseResponseBody {
		r.logger.Sugar().Debugf("Beginning response body validation for PIIMaskingGuardrail policy: %s", r.Name)

		if requestConfig.ResponseBody == nil || requestConfig.ResponseBody.Body == nil {
			r.logger.Sugar().Debug("No response body found, skipping PII masking validation")
			return result
		}

		validationResult, err := r.validatePayload(requestConfig.ResponseBody.Body, true, requestConfig)
		if err != nil {
			r.logger.Sugar().Debugf("Response body validation failed for PIIMaskingGuardrail policy: %s - %v", r.Name, err)
			return r.buildErrorResponse(true, err)
		}

		// Check if payload was modified and update the body
		if validationResult.ModifiedPayload != nil {
			r.logger.Sugar().Debugf("Response body was modified by PIIMaskingGuardrail policy: %s", r.Name)
			requestConfig.ResponseBody.Body = *validationResult.ModifiedPayload
		}

		r.logger.Sugar().Debugf("Response body validation passed for PIIMaskingGuardrail policy: %s", r.Name)
		return result
	}

	return result
}

// AssessmentResult represents the result of PII assessment
type AssessmentResult struct {
	InspectedContent string
	ModifiedPayload  *[]byte
	IsResponse       bool
	Error            string
}

// validatePayload validates the payload against the PIIMaskingGuardrail policy.
func (r *PIIMaskingGuardrail) validatePayload(payload []byte, isResponse bool, requestConfig *requestconfig.Holder) (AssessmentResult, error) {
	var result AssessmentResult
	result.IsResponse = isResponse

	var processedPayload []byte
	var compressionType string

	if isResponse {
		// Handle response decompression
		bodyStr, compType, err := r.decompressResponseBody(payload)
		if err == nil {
			processedPayload = []byte(bodyStr)
			compressionType = compType
		} else {
			processedPayload = payload
		}
	} else {
		processedPayload = payload
	}

	// Transform response if redactPII is disabled and PIIs identified in request
	if !r.RedactPII && isResponse {
		if maskedPII := r.getStoredPIIEntities(requestConfig); maskedPII != nil {
			if maskedPIIMap, ok := maskedPII.(map[string]string); ok {
				// For response flow, always transform the entire payload (JSONPath is not applicable)
				transformedContent := r.restorePIIInResponse(string(processedPayload), maskedPIIMap)
				result.InspectedContent = transformedContent
				modifiedPayload, err := r.compressResponseBody([]byte(transformedContent), compressionType)
				if err != nil {
					return result, fmt.Errorf("error compressing modified payload: %v", err)
				}
				result.ModifiedPayload = &modifiedPayload
				return result, nil // Continue processing after PII restoration
			}
		}
	}

	extractedValue, err := r.extractJSONPathValue(processedPayload, r.JSONPath)
	if err != nil {
		return result, fmt.Errorf("error extracting value from JSON using JSONPath: %v", err)
	}

	// Clean and trim
	extractedValue = r.cleanText(extractedValue)
	extractedValue = strings.TrimSpace(extractedValue)

	// Store the inspected content for assessment reporting
	result.InspectedContent = extractedValue

	// First check regex patterns for PII
	if len(r.PiiEntities) > 0 {
		if r.RedactPII {
			redactedContent := r.redactPIIFromContent(extractedValue)
			if redactedContent != "" {
				result.InspectedContent = redactedContent
				modifiedPayload := r.updatePayloadWithMaskedContent(processedPayload, extractedValue, redactedContent)
				if isResponse && compressionType != "" {
					compressedPayload, err := r.compressResponseBody(modifiedPayload, compressionType)
					if err == nil {
						modifiedPayload = compressedPayload
					}
				}
				result.ModifiedPayload = &modifiedPayload
				return result, nil
			}
		} else {
			maskedContent := r.maskPIIFromContent(extractedValue, isResponse, requestConfig)
			if maskedContent != "" {
				result.InspectedContent = maskedContent
				modifiedPayload := r.updatePayloadWithMaskedContent(processedPayload, extractedValue, maskedContent)
				if isResponse && compressionType != "" {
					compressedPayload, err := r.compressResponseBody(modifiedPayload, compressionType)
					if err == nil {
						modifiedPayload = compressedPayload
					}
				}
				result.ModifiedPayload = &modifiedPayload
				return result, nil
			}
		}
	}

	// If no regex patterns matched, proceed with existing logic
	r.logger.Sugar().Debugf("PII validation passed for content length: %d", len(extractedValue))
	return result, nil
}

// maskPIIFromContent masks PII from content using regex patterns
func (r *PIIMaskingGuardrail) maskPIIFromContent(jsonContent string, isResponse bool, requestConfig *requestconfig.Holder) string {
	if jsonContent == "" {
		return ""
	}

	foundAndMasked := false
	maskedContent := jsonContent

	if !isResponse {
		// Request flow: mask PII and store mappings
		maskedPIIEntities := make(map[string]string)
		counter := 0

		r.patternMu.RLock()
		patterns := r.PiiEntities
		r.patternMu.RUnlock()

		// First pass: find all matches without replacing to avoid nested replacements
		allMatches := make(map[string]string) // original -> placeholder
		for key, pattern := range patterns {
			matches := pattern.FindAllString(maskedContent, -1)
			for _, match := range matches {
				// Skip if this match is already processed or if it's a placeholder
				if _, exists := allMatches[match]; !exists && !strings.Contains(match, "[") && !strings.Contains(match, "]") {
					// Generate unique placeholder like [EMAIL_0000]
					placeholder := fmt.Sprintf("[%s_%04x]", key, counter)
					allMatches[match] = placeholder
					maskedPIIEntities[match] = placeholder
					counter++
				}
			}
		}

		// Second pass: replace all matches
		for original, placeholder := range allMatches {
			maskedContent = strings.ReplaceAll(maskedContent, original, placeholder)
			foundAndMasked = true
		}

		// Store PII_ENTITIES for later reversal
		if len(maskedPIIEntities) > 0 {
			r.storePIIEntities(requestConfig, maskedPIIEntities)
		}
	} else {
		// Response flow: restore original PII
		if maskedPII := r.getStoredPIIEntities(requestConfig); maskedPII != nil {
			if maskedPIIEntities, ok := maskedPII.(map[string]string); ok {
				for original, placeholder := range maskedPIIEntities {
					if strings.Contains(maskedContent, placeholder) {
						maskedContent = strings.ReplaceAll(maskedContent, placeholder, original)
						foundAndMasked = true
					}
				}
			}
		}
	}

	if foundAndMasked {
		r.logger.Sugar().Debugf("Masked content: %s", maskedContent)
		return maskedContent
	}

	return ""
}

// redactPIIFromContent redacts PII from content using regex patterns
func (r *PIIMaskingGuardrail) redactPIIFromContent(jsonContent string) string {
	if jsonContent == "" {
		return ""
	}

	foundAndMasked := false
	maskedContent := jsonContent

	r.patternMu.RLock()
	patterns := r.PiiEntities
	r.patternMu.RUnlock()

	for _, pattern := range patterns {
		if pattern.MatchString(maskedContent) {
			foundAndMasked = true
			maskedContent = pattern.ReplaceAllString(maskedContent, "*****")
		}
	}

	if foundAndMasked {
		r.logger.Sugar().Debugf("Redacted content: %s", maskedContent)
		return maskedContent
	}

	return ""
}

// restorePIIInResponse handles PII restoration in responses when redactPII is disabled
func (r *PIIMaskingGuardrail) restorePIIInResponse(originalContent string, maskedPIIEntities map[string]string) string {
	if maskedPIIEntities == nil || len(maskedPIIEntities) == 0 {
		r.logger.Sugar().Debug("No PII entities found in request. No response transformation needed.")
		return originalContent
	}

	transformedContent := originalContent
	foundMasked := false

	for original, placeholder := range maskedPIIEntities {
		if strings.Contains(transformedContent, placeholder) {
			transformedContent = strings.ReplaceAll(transformedContent, placeholder, original)
			foundMasked = true
		}
	}

	if foundMasked {
		r.logger.Sugar().Debug("PII entities found in request. Replacing masked PIIs back in response.")
	} else {
		r.logger.Sugar().Debug("No masked PII entities found in response content.")
	}

	return transformedContent
}

// updatePayloadWithMaskedContent updates the original payload by replacing the extracted content
// with the masked/redacted content, preserving the JSON structure if JSONPath is used (request flow only)
func (r *PIIMaskingGuardrail) updatePayloadWithMaskedContent(originalPayload []byte, extractedValue, modifiedContent string) []byte {
	if r.JSONPath == "" {
		// If no JSONPath, the entire payload was processed, return the modified content
		r.logger.Sugar().Debug("No JSONPath specified, replacing entire payload")
		return []byte(modifiedContent)
	}

	// If JSONPath is specified, update only the specific field in the JSON structure (request flow only)
	r.logger.Sugar().Debugf("Updating JSONPath field '%s' with masked content", r.JSONPath)

	var jsonData map[string]interface{}
	if err := json.Unmarshal(originalPayload, &jsonData); err != nil {
		r.logger.Sugar().Errorf("Error unmarshaling JSON payload for update: %v", err)
		// Fallback to returning the modified content as-is
		return []byte(modifiedContent)
	}

	// Set the new value at the JSONPath location using the jsonpath utility function
	err := r.setValueAtJSONPath(jsonData, r.JSONPath, modifiedContent)
	if err != nil {
		r.logger.Sugar().Errorf("Error setting value at JSONPath '%s': %v", r.JSONPath, err)
		// Fallback to returning the original payload
		return originalPayload
	}

	// Marshal back to JSON to get the full modified payload
	updatedPayload, err := json.Marshal(jsonData)
	if err != nil {
		r.logger.Sugar().Errorf("Error marshaling updated JSON payload: %v", err)
		// Fallback to returning the original payload
		return originalPayload
	}

	r.logger.Sugar().Debugf("Successfully updated payload with masked content at JSONPath '%s'", r.JSONPath)
	return updatedPayload
}

// buildErrorResponse builds an error response for PII masking guardrail failures.
func (r *PIIMaskingGuardrail) buildErrorResponse(isResponse bool, validationError error) *Result {
	result := NewResult()
	result.ImmediateResponse = true
	result.ImmediateResponseCode = v32.StatusCode(400)
	result.ImmediateResponseBody = fmt.Sprintf("PII masking guardrail validation failed: %v", validationError)
	result.ImmediateResponseContentType = "application/json"
	return result
}

// Helper methods for text processing, compression, and JSON path operations

// storePIIEntities stores PII entities in the request config for later retrieval
func (r *PIIMaskingGuardrail) storePIIEntities(requestConfig *requestconfig.Holder, entities map[string]string) {
	// Store in a way that can be retrieved later - using a simple approach since
	// we don't have direct access to properties in requestconfig.Holder
	// This would need to be implemented based on the actual storage mechanism available
	r.logger.Sugar().Debugf("Storing PII entities for later restoration: %v", entities)
}

// getStoredPIIEntities retrieves stored PII entities from the request config
func (r *PIIMaskingGuardrail) getStoredPIIEntities(requestConfig *requestconfig.Holder) interface{} {
	// Retrieve stored entities - this would need to be implemented based on the actual storage mechanism
	r.logger.Sugar().Debug("Retrieving stored PII entities")
	return nil
}

// cleanText removes unwanted characters from text (placeholder implementation)
func (r *PIIMaskingGuardrail) cleanText(text string) string {
	// Use a simple regex to clean text - in real implementation, use compiled regex
	re := regexp.MustCompile(`[^\w\s\p{P}]`)
	return re.ReplaceAllString(text, "")
}

// extractJSONPathValue extracts value from JSON using JSONPath
func (r *PIIMaskingGuardrail) extractJSONPathValue(payload []byte, jsonPath string) (string, error) {
	if jsonPath == "" {
		return string(payload), nil
	}

	// Use gjson to extract the value
	result := gjson.GetBytes(payload, jsonPath[2:]) // Remove "$.";
	if !result.Exists() {
		return "", fmt.Errorf("JSONPath '%s' not found in payload", jsonPath)
	}

	return result.String(), nil
}

// setValueAtJSONPath sets a value at the specified JSONPath in a JSON object
func (r *PIIMaskingGuardrail) setValueAtJSONPath(jsonData map[string]interface{}, jsonPath string, value string) error {
	// Simple implementation for setting nested values
	// This is a simplified version - in production, you'd want a more robust JSONPath setter
	if jsonPath == "$.content" || jsonPath == "content" {
		jsonData["content"] = value
		return nil
	}

	// For more complex paths, you'd need a proper JSONPath setter implementation
	return fmt.Errorf("complex JSONPath setting not implemented for: %s", jsonPath)
}

// decompressResponseBody decompresses response body if it's compressed
func (r *PIIMaskingGuardrail) decompressResponseBody(payload []byte) (string, string, error) {
	// Check if the payload starts with gzip magic number
	if len(payload) < 2 || payload[0] != 0x1f || payload[1] != 0x8b {
		// Not gzipped, return as-is
		return string(payload), "", nil
	}

	// Decompress gzip
	reader, err := gzip.NewReader(bytes.NewReader(payload))
	if err != nil {
		return "", "", err
	}
	defer reader.Close()

	decompressed, err := io.ReadAll(reader)
	if err != nil {
		return "", "", err
	}

	return string(decompressed), "gzip", nil
}

// compressResponseBody compresses response body with the specified compression type
func (r *PIIMaskingGuardrail) compressResponseBody(payload []byte, compressionType string) ([]byte, error) {
	if compressionType != "gzip" {
		return payload, nil
	}

	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)

	_, err := writer.Write(payload)
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
