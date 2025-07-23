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

package inbuiltpolicy

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"sync"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_service_proc_v3 "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"
	v32 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/wso2/apk/gateway/enforcer/internal/dto"
	"github.com/wso2/apk/gateway/enforcer/internal/logging"
)

// PIIMaskingRegex is a struct that represents a regex-based PII masking policy.
type PIIMaskingRegex struct {
	dto.BaseInBuiltPolicy
	Name        string
	PiiEntities map[string]*regexp.Regexp
	JSONPath    string
	RedactPII   bool
	patternMu   sync.RWMutex
}

// HandleRequestBody is a method that implements the mediation logic for the PIIMaskingRegex policy on request.
func (r *PIIMaskingRegex) HandleRequestBody(logger *logging.Logger, req *envoy_service_proc_v3.ProcessingRequest, resp *envoy_service_proc_v3.ProcessingResponse, props map[string]interface{}) *envoy_service_proc_v3.ProcessingResponse {
	logger.Sugar().Debugf("Beginning request payload validation for PIIMaskingRegex policy: %s", r.Name)
	result, ok := r.validatePayload(logger, req, false, props)
	if !ok {
		logger.Sugar().Debugf("Request payload validation failed for PIIMaskingRegex policy: %s", r.Name)
		return r.buildResponse(logger, result)
	}

	// Check if payload was modified and return the modified content
	if result.ModifiedPayload != nil {
		logger.Sugar().Debugf("Request payload was modified by PIIMaskingRegex policy: %s", r.Name)
		r.buildBodyMutationResponse(resp, *result.ModifiedPayload, false)
	}

	logger.Sugar().Debugf("Request payload validation passed for PIIMaskingRegex policy: %s", r.Name)
	return nil
}

// HandleResponseBody is a method that implements the mediation logic for the PIIMaskingRegex policy on response.
func (r *PIIMaskingRegex) HandleResponseBody(logger *logging.Logger, req *envoy_service_proc_v3.ProcessingRequest, resp *envoy_service_proc_v3.ProcessingResponse, props map[string]interface{}) *envoy_service_proc_v3.ProcessingResponse {
	logger.Sugar().Debugf("Beginning response body validation for PIIMaskingRegex policy: %s", r.Name)
	result, ok := r.validatePayload(logger, req, true, props)
	if !ok {
		logger.Sugar().Debugf("Response body validation failed for PIIMaskingRegex policy: %s", r.Name)
		return r.buildResponse(logger, result)
	}

	// Check if payload was modified and return the modified content
	if result.ModifiedPayload != nil {
		logger.Sugar().Debugf("Response body was modified by PIIMaskingRegex policy: %s", r.Name)
		r.buildBodyMutationResponse(resp, *result.ModifiedPayload, true)
	}

	logger.Sugar().Debugf("Response body validation passed for PIIMaskingRegex policy: %s", r.Name)
	return nil
}

// validatePayload validates the payload against the PIIMaskingRegex policy.
func (r *PIIMaskingRegex) validatePayload(logger *logging.Logger, req *envoy_service_proc_v3.ProcessingRequest, isResponse bool, props map[string]interface{}) (AssessmentResult, bool) {
	var result AssessmentResult
	result.IsResponse = isResponse

	var payload []byte
	var compressionType string
	if isResponse {
		var bodyStr string
		var err error
		bodyStr, compressionType, err = DecompressLLMResp(payload)
		if err == nil {
			payload = []byte(bodyStr)
		}
	} else {
		payload = req.GetRequestBody().Body
	}

	// Transform response if redactPII is disabled and PIIs identified in request
	if !r.RedactPII && isResponse {
		if maskedPII, exists := props["PIIMaskingRegexPIIEntities"]; exists {
			if maskedPIIMap, ok := maskedPII.(map[string]string); ok {
				// For response flow, always transform the entire payload (JSONPath is not applicable)
				transformedContent := r.restorePIIInResponse(string(payload), maskedPIIMap, logger)
				result.InspectedContent = transformedContent
				modifiedPayload, err := CompressLLMResp([]byte(transformedContent), compressionType)
				if err != nil {
					result.Error = "Error compressing modified payload: " + err.Error()
					logger.Error(err, result.Error)
					return result, false
				}
				result.ModifiedPayload = &modifiedPayload
				return result, true // Continue processing after PII restoration
			}
		}
	}

	extractedValue, err := ExtractStringValueFromJsonpath(logger, payload, r.JSONPath)
	if err != nil {
		result.Error = "Error extracting value from JSON using JSONPath: " + err.Error()
		logger.Error(err, result.Error)
		return result, false
	}
	// Clean and trim
	extractedValue = TextCleanRegexCompiled.ReplaceAllString(extractedValue, "")
	extractedValue = strings.TrimSpace(extractedValue)

	// Store the inspected content for assessment reporting
	result.InspectedContent = extractedValue

	// First check regex patterns for PII
	if len(r.PiiEntities) > 0 {
		if r.RedactPII {
			redactedContent := r.redactPIIFromContent(extractedValue, logger)
			if redactedContent != "" {
				result.InspectedContent = redactedContent
				modifiedPayload := r.updatePayloadWithMaskedContent(payload, extractedValue, redactedContent, logger)
				result.ModifiedPayload = &modifiedPayload
				return result, true
			}
		} else {
			maskedContent := r.maskPIIFromContent(extractedValue, isResponse, props, logger)
			if maskedContent != "" {
				result.InspectedContent = maskedContent
				modifiedPayload := r.updatePayloadWithMaskedContent(payload, extractedValue, maskedContent, logger)
				result.ModifiedPayload = &modifiedPayload
				return result, true
			}
		}
	}

	// If no regex patterns matched, proceed with existing logic
	logger.Sugar().Debugf("PII validation passed for content length: %d", len(extractedValue))
	return result, true
}

// maskPIIFromContent masks PII from content using regex patterns
func (r *PIIMaskingRegex) maskPIIFromContent(jsonContent string, isResponse bool, props map[string]interface{}, logger *logging.Logger) string {
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

		for key, pattern := range patterns {
			matches := pattern.FindAllString(maskedContent, -1)
			for _, match := range matches {
				// Reuse if already seen
				if _, exists := maskedPIIEntities[match]; !exists {
					// Generate unique placeholder like <Person_0001>
					masked := fmt.Sprintf("<%s_%04x>", key, counter)
					maskedPIIEntities[match] = masked
					counter++
				}
				maskedContent = strings.ReplaceAll(maskedContent, match, maskedPIIEntities[match])
				foundAndMasked = true
			}
		}

		// Store PII_ENTITIES for later reversal
		if len(maskedPIIEntities) > 0 {
			if dynamicMetadataKeyValuePairs, ok := props["dynamicMetadataMap"].(map[string]interface{}); ok {
				dynamicMetadataKeyValuePairs[piiMaskingRegexPIIEntitiesKey] = maskedPIIEntities
			}
		}
	} else {
		// Response flow: restore original PII
		if maskedPII, exists := props["piiMaskingRegexPIIEntities"]; exists {
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
		logger.Sugar().Debugf("Masked content: %s", maskedContent)
		return maskedContent
	}

	return ""
}

// redactPIIFromContent redacts PII from content using regex patterns
func (r *PIIMaskingRegex) redactPIIFromContent(jsonContent string, logger *logging.Logger) string {
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
		logger.Sugar().Debugf("Redacted content: %s", maskedContent)
		return maskedContent
	}

	return ""
}

// restorePIIInResponse handles PII restoration in responses when redactPII is disabled
func (r *PIIMaskingRegex) restorePIIInResponse(originalContent string, maskedPIIEntities map[string]string, logger *logging.Logger) string {
	if maskedPIIEntities == nil || len(maskedPIIEntities) == 0 {
		logger.Sugar().Debug("No PII entities found in request. No response transformation needed.")
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
		logger.Sugar().Debug("PII entities found in request. Replacing masked PIIs back in response.")
	} else {
		logger.Sugar().Debug("No masked PII entities found in response content.")
	}

	return transformedContent
}

// buildResponse is a method that builds the response body for the PIIMaskingRegex policy.
func (r *PIIMaskingRegex) buildResponse(logger *logging.Logger, result AssessmentResult) *envoy_service_proc_v3.ProcessingResponse {
	responseBody := make(map[string]interface{})
	responseBody[ErrorCode] = APIMInternalExceptionCode
	responseBody[ErrorMessage] = "Error occurred during PIIMaskingRegex mediation."

	bodyBytes, err := json.Marshal(responseBody)
	if err != nil {
		logger.Sugar().Error(err, "Error marshaling response body to JSON")
		return nil
	}

	headers := &envoy_service_proc_v3.HeaderMutation{
		SetHeaders: []*corev3.HeaderValueOption{
			{
				Header: &corev3.HeaderValue{
					Key:      "Content-Type",
					RawValue: []byte("Application/json"),
				},
			},
		},
	}

	return &envoy_service_proc_v3.ProcessingResponse{
		Response: &envoy_service_proc_v3.ProcessingResponse_ImmediateResponse{
			ImmediateResponse: &envoy_service_proc_v3.ImmediateResponse{
				Status: &v32.HttpStatus{
					Code: v32.StatusCode(APIMInternalErrorCode),
				},
				Body:    bodyBytes,
				Headers: headers,
			},
		},
	}
}

// updatePayloadWithMaskedContent updates the original payload by replacing the extracted content
// with the masked/redacted content, preserving the JSON structure if JSONPath is used (request flow only)
func (r *PIIMaskingRegex) updatePayloadWithMaskedContent(originalPayload []byte, extractedValue, modifiedContent string, logger *logging.Logger) []byte {
	if r.JSONPath == "" {
		// If no JSONPath, the entire payload was processed, return the modified content
		logger.Sugar().Debug("No JSONPath specified, replacing entire payload")
		return []byte(modifiedContent)
	}

	// If JSONPath is specified, update only the specific field in the JSON structure (request flow only)
	logger.Sugar().Debugf("Updating JSONPath field '%s' with masked content", r.JSONPath)

	var jsonData map[string]interface{}
	if err := json.Unmarshal(originalPayload, &jsonData); err != nil {
		logger.Sugar().Errorf("Error unmarshaling JSON payload for update: %v", err)
		// Fallback to returning the modified content as-is
		return []byte(modifiedContent)
	}

	// Set the new value at the JSONPath location using the jsonpath utility function
	err := setValueAtJSONPath(jsonData, r.JSONPath, modifiedContent)
	if err != nil {
		logger.Sugar().Errorf("Error setting value at JSONPath '%s': %v", r.JSONPath, err)
		// Fallback to returning the original payload
		return originalPayload
	}

	// Marshal back to JSON to get the full modified payload
	updatedPayload, err := json.Marshal(jsonData)
	if err != nil {
		logger.Sugar().Errorf("Error marshaling updated JSON payload: %v", err)
		// Fallback to returning the original payload
		return originalPayload
	}

	logger.Sugar().Debugf("Successfully updated payload with masked content at JSONPath '%s'", r.JSONPath)
	return updatedPayload
}

// buildBodyMutationResponse creates a response that modifies the request/response body
func (r *PIIMaskingRegex) buildBodyMutationResponse(resp *envoy_service_proc_v3.ProcessingResponse, modifiedBody []byte, isResponse bool) {
	// Calculate the new body length
	newBodyLength := len(modifiedBody)

	// Update the Content-Length header
	headers := &envoy_service_proc_v3.HeaderMutation{
		SetHeaders: []*corev3.HeaderValueOption{
			{
				Header: &corev3.HeaderValue{
					Key:      "Content-Length",
					RawValue: []byte(fmt.Sprintf("%d", newBodyLength)),
				},
			},
		},
	}

	bodyResponse := &envoy_service_proc_v3.BodyResponse{
		Response: &envoy_service_proc_v3.CommonResponse{
			Status:         envoy_service_proc_v3.CommonResponse_CONTINUE_AND_REPLACE,
			HeaderMutation: headers,
			BodyMutation: &envoy_service_proc_v3.BodyMutation{
				Mutation: &envoy_service_proc_v3.BodyMutation_Body{
					Body: modifiedBody,
				},
			},
		},
	}

	if isResponse {
		resp.Response = &envoy_service_proc_v3.ProcessingResponse_ResponseBody{
			ResponseBody: bodyResponse,
		}
	} else {
		resp.Response = &envoy_service_proc_v3.ProcessingResponse_RequestBody{
			RequestBody: bodyResponse,
		}
	}
}

// NewPIIMaskingRegex initializes the PIIMaskingRegex policy from the given InBuiltPolicy.
func NewPIIMaskingRegex(logger *logging.Logger, inBuiltPolicy dto.InBuiltPolicy) *PIIMaskingRegex {
	// Set default values
	PIIMaskingRegex := &PIIMaskingRegex{
		BaseInBuiltPolicy: dto.BaseInBuiltPolicy{
			PolicyName:    inBuiltPolicy.GetPolicyName(),
			PolicyID:      inBuiltPolicy.GetPolicyID(),
			PolicyVersion: inBuiltPolicy.GetPolicyVersion(),
			Parameters:    inBuiltPolicy.GetParameters(),
			PolicyOrder:   inBuiltPolicy.GetPolicyOrder(),
		},
		Name:      PIIMaskingRegexName,
		JSONPath:  "",
		RedactPII: false,
	}

	for key, value := range inBuiltPolicy.GetParameters() {
		switch key {
		case "name":
			PIIMaskingRegex.Name = value
		case "piiEntities":
			var piiEntities map[string]string
			if err := json.Unmarshal([]byte(value), &piiEntities); err != nil {
				PIIMaskingRegex.PiiEntities = make(map[string]*regexp.Regexp)
				// Skip adding error pattern for invalid format
			} else {
				// Compile regex patterns during initialization
				PIIMaskingRegex.PiiEntities = make(map[string]*regexp.Regexp)
				for entityKey, pattern := range piiEntities {
					compiledPattern, err := regexp.Compile(pattern)
					if err != nil {
						logger.Sugar().Errorf("Error compiling regex for PII entity '%s': %v", entityKey, err)
						continue
					}
					PIIMaskingRegex.PiiEntities[entityKey] = compiledPattern
				}
			}
		case "jsonPath":
			PIIMaskingRegex.JSONPath = value
		case "redactPII":
			PIIMaskingRegex.RedactPII = value == "true"
		}
	}

	return PIIMaskingRegex
}
