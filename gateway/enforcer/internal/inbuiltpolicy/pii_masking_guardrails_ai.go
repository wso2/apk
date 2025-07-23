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
	"io"
	"net/http"
	"strings"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_service_proc_v3 "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"
	v32 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/wso2/apk/gateway/enforcer/internal/dto"
	"github.com/wso2/apk/gateway/enforcer/internal/logging"
	"github.com/wso2/apk/gateway/enforcer/internal/util"
)

// PIIServiceRequest represents the request payload for the PII masking service
type PIIServiceRequest struct {
	Text        string   `json:"text"`
	Redact      bool     `json:"redact"`
	PiiEntities []string `json:"piiEntities"`
}

// PIIAssessment represents individual PII assessment from the service
type PIIAssessment struct {
	PiiEntity string `json:"piiEntity"`
	PiiValue  string `json:"piiValue"`
}

// PIIServiceResponse represents the response from the PII masking service
type PIIServiceResponse struct {
	AnonymizedText string          `json:"anonymizedText"`
	Assessment     []PIIAssessment `json:"assessment"`
}

// PIIMaskingGuardrailsAI is a struct that represents a REST-based PII masking policy.
type PIIMaskingGuardrailsAI struct {
	dto.BaseInBuiltPolicy
	Name        string
	PiiEntities []string
	JSONPath    string
	RedactPII   bool
}

const (
	piiServiceURL = "http://52.230.120.80:9447/validate"
)

// HandleRequestBody is a method that implements the mediation logic for the PIIMaskingGuardrailsAI policy on request.
func (r *PIIMaskingGuardrailsAI) HandleRequestBody(logger *logging.Logger, req *envoy_service_proc_v3.ProcessingRequest, resp *envoy_service_proc_v3.ProcessingResponse, props map[string]interface{}) *envoy_service_proc_v3.ProcessingResponse {
	logger.Sugar().Debugf("Beginning request payload validation for PIIMaskingGuardrailsAI policy: %s", r.Name)
	result, ok := r.validatePayload(logger, req, false, props)
	if !ok {
		logger.Sugar().Debugf("Request payload validation failed for PIIMaskingGuardrailsAI policy: %s", r.Name)
		return r.buildResponse(logger, result)
	}

	// Check if payload was modified and return the modified content
	if result.ModifiedPayload != nil {
		logger.Sugar().Debugf("Request payload was modified by PIIMaskingGuardrailsAI policy: %s", r.Name)
		r.buildBodyMutationResponse(resp, *result.ModifiedPayload, false)
	}

	logger.Sugar().Debugf("Request payload validation passed for PIIMaskingGuardrailsAI policy: %s", r.Name)
	return nil
}

// HandleResponseBody is a method that implements the mediation logic for the PIIMaskingGuardrailsAI policy on response.
func (r *PIIMaskingGuardrailsAI) HandleResponseBody(logger *logging.Logger, req *envoy_service_proc_v3.ProcessingRequest, resp *envoy_service_proc_v3.ProcessingResponse, props map[string]interface{}) *envoy_service_proc_v3.ProcessingResponse {
	logger.Sugar().Debugf("Beginning response body validation for PIIMaskingGuardrailsAI policy: %s", r.Name)
	result, ok := r.validatePayload(logger, req, true, props)
	if !ok {
		logger.Sugar().Debugf("Response body validation failed for PIIMaskingGuardrailsAI policy: %s", r.Name)
		return r.buildResponse(logger, result)
	}

	// Check if payload was modified and return the modified content
	if result.ModifiedPayload != nil {
		logger.Sugar().Debugf("Response body was modified by PIIMaskingGuardrailsAI policy: %s", r.Name)
		r.buildBodyMutationResponse(resp, *result.ModifiedPayload, true)
	}

	logger.Sugar().Debugf("Response body validation passed for PIIMaskingGuardrailsAI policy: %s", r.Name)
	return nil
}

// validatePayload validates the payload against the PIIMaskingGuardrailsAI policy.
func (r *PIIMaskingGuardrailsAI) validatePayload(logger *logging.Logger, req *envoy_service_proc_v3.ProcessingRequest, isResponse bool, props map[string]interface{}) (AssessmentResult, bool) {
	var result AssessmentResult
	result.IsResponse = isResponse

	var payload []byte
	var compressionType string
	if isResponse {
		var bodyStr string
		var err error
		payload = req.GetResponseBody().Body
		bodyStr, compressionType, err = DecompressLLMResp(payload)
		if err == nil {
			payload = []byte(bodyStr)
		}
	} else {
		payload = req.GetRequestBody().Body
	}

	// Transform response if redactPII is disabled and PIIs identified in request
	if !r.RedactPII && isResponse {
		if maskedPII, exists := props["piiMaskingGuardrailsAIPIIEntities"]; exists {
			if maskedPIIMap, ok := maskedPII.(map[string]string); ok {
				// For response flow, treat the entire payload as a string and replace placeholders
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

	// Call PII service to identify and process PII entities
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

	// If no PII entities found or processed, continue with validation
	logger.Sugar().Debugf("PII validation passed for content length: %d", len(extractedValue))
	return result, true
}

// maskPIIFromContent masks PII from content using REST API calls
func (r *PIIMaskingGuardrailsAI) maskPIIFromContent(jsonContent string, isResponse bool, props map[string]interface{}, logger *logging.Logger) string {
	if jsonContent == "" {
		return ""
	}

	if !isResponse {
		// Request flow: call PII service to identify and mask PII
		return r.callPIIService(jsonContent, r.RedactPII, props, logger)
	}
	if maskedPII, exists := props["piiMaskingGuardrailsAIPIIEntities"]; exists {
		if maskedPIIEntities, ok := maskedPII.(map[string]string); ok {
			return r.restorePIIInResponse(jsonContent, maskedPIIEntities, logger)
		}
	}
	return ""
}

// redactPIIFromContent redacts PII from content using REST API calls
func (r *PIIMaskingGuardrailsAI) redactPIIFromContent(jsonContent string, logger *logging.Logger) string {
	if jsonContent == "" {
		return ""
	}

	logger.Sugar().Debugf("Redacting PII from content: %s", jsonContent)
	logger.Sugar().Debugf("Using %d PII entities for redaction", len(r.PiiEntities))

	if len(r.PiiEntities) == 0 {
		logger.Sugar().Debug("No PII entities defined, skipping redaction")
		return ""
	}

	// Call PII service for redaction using the same method, just pass redact=true
	return r.callPIIService(jsonContent, true, nil, logger)
}

// callPIIService makes HTTP request to PII service for masking (request flow)
func (r *PIIMaskingGuardrailsAI) callPIIService(jsonContent string, redact bool, props map[string]interface{}, logger *logging.Logger) string {
	// Prepare request payload - always send redact: false to get original PII values
	requestPayload := PIIServiceRequest{
		Text:        jsonContent,
		Redact:      false, // Always false to get original PII values for local processing
		PiiEntities: r.PiiEntities,
	}

	jsonData, err := json.Marshal(requestPayload)
	if err != nil {
		logger.Sugar().Errorf("Failed to marshal PII service request: %v", err)
		return ""
	}

	// Prepare headers
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	// Make HTTP request with retry
	resp, err := util.MakeHTTPRequestWithRetry("POST", piiServiceURL, nil, headers, jsonData, 30000, 3, 1000)
	if err != nil {
		logger.Sugar().Errorf("Failed to call PII masking service: %v", err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		logger.Sugar().Errorf("Unexpected status code %d from PII service: %s", resp.StatusCode, string(bodyBytes))
		return ""
	}

	// Parse response
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Sugar().Errorf("Failed to read response body: %v", err)
		return ""
	}

	var piiResponse PIIServiceResponse
	if err := json.Unmarshal(responseBody, &piiResponse); err != nil {
		logger.Sugar().Errorf("Failed to unmarshal PII service response: %v", err)
		return ""
	}

	// Process the response based on redact flag
	processedContent := r.processPIIResponse(jsonContent, piiResponse, redact, props, logger)

	logger.Sugar().Debugf("PII service returned processed text: %s", processedContent)
	return processedContent
}

// processPIIResponse processes the PII service response based on redact/mask mode
func (r *PIIMaskingGuardrailsAI) processPIIResponse(originalContent string, piiResponse PIIServiceResponse, redact bool, props map[string]interface{}, logger *logging.Logger) string {
	if len(piiResponse.Assessment) == 0 {
		logger.Sugar().Debug("No PII entities found in assessment")
		return originalContent
	}

	processedContent := originalContent
	counter := 0

	if redact {
		// Redaction mode: replace all PII values with *****
		for _, assessment := range piiResponse.Assessment {
			if strings.Contains(processedContent, assessment.PiiValue) {
				processedContent = strings.ReplaceAll(processedContent, assessment.PiiValue, "*****")
				logger.Sugar().Debugf("Redacted PII value '%s' with *****", assessment.PiiValue)
			}
		}
	} else {
		// Masking mode: replace PII values with piiEntity + hexID and store mappings
		maskedPIIEntities := make(map[string]string)

		for _, assessment := range piiResponse.Assessment {
			if strings.Contains(processedContent, assessment.PiiValue) {
				// Generate unique placeholder like [PERSON_0001]
				placeholder := fmt.Sprintf("[%s_%04x]", strings.ToUpper(assessment.PiiEntity), counter)
				processedContent = strings.ReplaceAll(processedContent, assessment.PiiValue, placeholder)

				// Store original value for response restoration
				maskedPIIEntities[assessment.PiiValue] = placeholder
				counter++

				logger.Sugar().Debugf("Masked PII value '%s' with '%s'", assessment.PiiValue, placeholder)
			}
		}

		// Store PII mappings for response restoration
		if len(maskedPIIEntities) > 0 {
			if dynamicMetadataKeyValuePairs, ok := props["dynamicMetadataMap"].(map[string]interface{}); ok {
				dynamicMetadataKeyValuePairs[piiMaskingGuardrailsAIPIIEntitiesKey] = maskedPIIEntities
			}
		}
	}

	return processedContent
}

// restorePIIInResponse handles PII restoration in responses when redactPII is disabled
func (r *PIIMaskingGuardrailsAI) restorePIIInResponse(originalContent string, maskedPIIEntities map[string]string, logger *logging.Logger) string {
	if maskedPIIEntities == nil || len(maskedPIIEntities) == 0 {
		logger.Sugar().Debug("No PII entities found in request. No response transformation needed.")
		return originalContent
	}

	transformedContent := originalContent
	foundMasked := false

	// The map structure is originalValue -> placeholder, so we need to iterate and replace placeholder with original
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

// buildResponse is a method that builds the response body for the PIIMaskingGuardrailsAI policy.
func (r *PIIMaskingGuardrailsAI) buildResponse(logger *logging.Logger, result AssessmentResult) *envoy_service_proc_v3.ProcessingResponse {
	responseBody := make(map[string]interface{})
	responseBody[ErrorCode] = APIMInternalExceptionCode
	responseBody[ErrorMessage] = "Error occurred during PIIMaskingGuardrailsAI mediation."

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
func (r *PIIMaskingGuardrailsAI) updatePayloadWithMaskedContent(originalPayload []byte, extractedValue, modifiedContent string, logger *logging.Logger) []byte {
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
func (r *PIIMaskingGuardrailsAI) buildBodyMutationResponse(resp *envoy_service_proc_v3.ProcessingResponse, modifiedBody []byte, isResponse bool) {
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

// NewPIIMaskingGuardrailsAI initializes the PIIMaskingGuardrailsAI policy from the given InBuiltPolicy.
func NewPIIMaskingGuardrailsAI(logger *logging.Logger, inBuiltPolicy dto.InBuiltPolicy) *PIIMaskingGuardrailsAI {
	// Set default values
	PIIMaskingGuardrailsAI := &PIIMaskingGuardrailsAI{
		BaseInBuiltPolicy: dto.BaseInBuiltPolicy{
			PolicyName:    inBuiltPolicy.GetPolicyName(),
			PolicyID:      inBuiltPolicy.GetPolicyID(),
			PolicyVersion: inBuiltPolicy.GetPolicyVersion(),
			Parameters:    inBuiltPolicy.GetParameters(),
			PolicyOrder:   inBuiltPolicy.GetPolicyOrder(),
		},
		Name:      PIIMaskingGuardrailsAIName,
		JSONPath:  "",
		RedactPII: false,
	}

	for key, value := range inBuiltPolicy.GetParameters() {
		switch key {
		case "name":
			PIIMaskingGuardrailsAI.Name = value
		case "piiEntities": // Expecting a string with , as a separator
			if value != "" {
				PIIMaskingGuardrailsAI.PiiEntities = strings.Split(value, ",")
				// Trim whitespace from each entity
				for i, entity := range PIIMaskingGuardrailsAI.PiiEntities {
					PIIMaskingGuardrailsAI.PiiEntities[i] = strings.TrimSpace(entity)
				}
			} else {
				logger.Sugar().Warn("No PII entities defined for PIIMaskingGuardrailsAI policy")
			}
		case "jsonPath":
			PIIMaskingGuardrailsAI.JSONPath = value
		case "redactPII":
			PIIMaskingGuardrailsAI.RedactPII = value == "true"
		}
	}

	logger.Sugar().Debugf("PIIMaskingGuardrailsAI initialized with %d PII entities", len(PIIMaskingGuardrailsAI.PiiEntities))
	for _, entity := range PIIMaskingGuardrailsAI.PiiEntities {
		logger.Sugar().Debugf("Loaded PII entity: %s", entity)
	}

	return PIIMaskingGuardrailsAI
}
