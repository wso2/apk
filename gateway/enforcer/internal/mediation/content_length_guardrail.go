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
	"strconv"
	"strings"

	v32 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/tidwall/gjson"
	dpv2alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v2alpha1"
	"github.com/wso2/apk/gateway/enforcer/internal/config"
	"github.com/wso2/apk/gateway/enforcer/internal/logging"
	"github.com/wso2/apk/gateway/enforcer/internal/requestconfig"
)

// ContentLengthGuardrail represents the configuration for Content Length Guardrail policy in the API Gateway.
type ContentLengthGuardrail struct {
	PolicyName     string `json:"policyName"`
	PolicyVersion  string `json:"policyVersion"`
	PolicyID       string `json:"policyID"`
	Name           string `json:"name"`
	Min            int    `json:"min"`
	Max            int    `json:"max"`
	JSONPath       string `json:"jsonPath"`
	Inverted       bool   `json:"inverted"`
	ShowAssessment bool   `json:"showAssessment"`
	logger         *logging.Logger
	cfg            *config.Server
}

const (
	// ContentLengthGuardrailPolicyKeyName is the key for specifying the name of the guardrail.
	ContentLengthGuardrailPolicyKeyName = "name"
	// ContentLengthGuardrailPolicyKeyMin is the key for specifying the minimum content length.
	ContentLengthGuardrailPolicyKeyMin = "min"
	// ContentLengthGuardrailPolicyKeyMax is the key for specifying the maximum content length.
	ContentLengthGuardrailPolicyKeyMax = "max"
	// ContentLengthGuardrailPolicyKeyJSONPath is the key for specifying the JSON path to extract content.
	ContentLengthGuardrailPolicyKeyJSONPath = "jsonPath"
	// ContentLengthGuardrailPolicyKeyInverted is the key for specifying if the validation should be inverted.
	ContentLengthGuardrailPolicyKeyInverted = "invert"
	// ContentLengthGuardrailPolicyKeyShowAssessment is the key for specifying if assessment should be shown.
	ContentLengthGuardrailPolicyKeyShowAssessment = "showAssessment"

	// ContentLengthGuardrailConstant is the constant for content length guardrail errors.
	ContentLengthGuardrailConstant = "CONTENT_LENGTH_GUARDRAIL"
)

// NewContentLengthGuardrail creates a new ContentLengthGuardrail instance.
func NewContentLengthGuardrail(mediation *dpv2alpha1.Mediation) *ContentLengthGuardrail {
	cfg := config.GetConfig()
	logger := cfg.Logger

	name := "ContentLengthGuardrail"
	if val, ok := extractPolicyValue(mediation.Parameters, ContentLengthGuardrailPolicyKeyName); ok {
		name = val
	}

	minL := 0
	if val, ok := extractPolicyValue(mediation.Parameters, ContentLengthGuardrailPolicyKeyMin); ok {
		if intValue, err := strconv.Atoi(val); err == nil {
			minL = intValue
		}
	}

	maxL := 10000
	if val, ok := extractPolicyValue(mediation.Parameters, ContentLengthGuardrailPolicyKeyMax); ok {
		if intValue, err := strconv.Atoi(val); err == nil {
			maxL = intValue
		}
	}

	jsonPath := "$.content"
	if val, ok := extractPolicyValue(mediation.Parameters, ContentLengthGuardrailPolicyKeyJSONPath); ok {
		jsonPath = val
	}

	inverted := false
	if val, ok := extractPolicyValue(mediation.Parameters, ContentLengthGuardrailPolicyKeyInverted); ok {
		inverted = val == "true"
	}

	showAssessment := false
	if val, ok := extractPolicyValue(mediation.Parameters, ContentLengthGuardrailPolicyKeyShowAssessment); ok {
		showAssessment = val == "true"
	}

	return &ContentLengthGuardrail{
		PolicyName:     "ContentLengthGuardrail",
		PolicyVersion:  mediation.PolicyVersion,
		PolicyID:       mediation.PolicyID,
		Name:           name,
		Min:            minL,
		Max:            maxL,
		JSONPath:       jsonPath,
		Inverted:       inverted,
		ShowAssessment: showAssessment,
		logger:         &logger,
		cfg:            cfg,
	}
}

// Process processes the request configuration for Content Length Guardrail.
func (c *ContentLengthGuardrail) Process(requestConfig *requestconfig.Holder) *Result {
	result := NewResult()

	// Handle request body processing
	if requestConfig.ProcessingPhase == requestconfig.ProcessingPhaseRequestBody {
		c.logger.Sugar().Debugf("Beginning request payload validation for ContentLengthGuardrail policy: %s", c.Name)

		if requestConfig.RequestBody == nil || requestConfig.RequestBody.Body == nil {
			c.logger.Sugar().Debug("No request body found, skipping content length validation")
			return result
		}

		validationResult, err := c.validatePayload(requestConfig.RequestBody.Body, false)
		if !validationResult {
			c.logger.Sugar().Debugf("Request payload validation failed for ContentLengthGuardrail policy: %s", c.Name)
			return c.buildErrorResponse(false, err)
		}
		c.logger.Sugar().Debugf("Request payload validation passed for ContentLengthGuardrail policy: %s", c.Name)
		return result
	}

	// Handle response body processing
	if requestConfig.ProcessingPhase == requestconfig.ProcessingPhaseResponseBody {
		c.logger.Sugar().Debugf("Beginning response body validation for ContentLengthGuardrail policy: %s", c.Name)

		if requestConfig.ResponseBody == nil || requestConfig.ResponseBody.Body == nil {
			c.logger.Sugar().Debug("No response body found, skipping content length validation")
			return result
		}

		validationResult, err := c.validatePayload(requestConfig.ResponseBody.Body, true)
		if !validationResult {
			c.logger.Sugar().Debugf("Response body validation failed for ContentLengthGuardrail policy: %s", c.Name)
			return c.buildErrorResponse(true, err)
		}
		c.logger.Sugar().Debugf("Response body validation passed for ContentLengthGuardrail policy: %s", c.Name)
		return result
	}

	return result
}

// validatePayload validates the payload against the ContentLengthGuardrail policy.
func (c *ContentLengthGuardrail) validatePayload(payload []byte, isResponse bool) (bool, error) {
	// Decompress response body if needed
	if isResponse {
		bodyStr, err := c.decompressLLMResp(payload)
		if err == nil {
			payload = []byte(bodyStr)
		}
	}

	// Validate content length range
	if (c.Min > c.Max) || (c.Min < 0) || (c.Max <= 0) {
		c.logger.Sugar().Errorf("Invalid content length range: min=%d, max=%d", c.Min, c.Max)
		return false, nil
	}

	// Extract value from JSON using JSONPath
	extractedValue, err := c.extractStringValueFromJsonpath(payload, c.JSONPath)
	if err != nil {
		c.logger.Error(err, "Error extracting value from JSON using JSONPath")
		return false, err
	}

	// Clean and trim the extracted text
	extractedValue = TextCleanRegexCompiled.ReplaceAllString(extractedValue, "")
	extractedValue = strings.TrimSpace(extractedValue)

	// Count the bytes in the extracted value
	byteCount := len([]byte(extractedValue))

	isWithinRange := byteCount >= c.Min && byteCount <= c.Max

	if c.Inverted {
		// When inverted, fail if content length is within the range
		if isWithinRange {
			c.logger.Sugar().Debugf("Content length validation failed (inverted): %d bytes found, should NOT be between %d and %d bytes", byteCount, c.Min, c.Max)
			return false, nil
		}
		c.logger.Sugar().Debugf("Content length validation passed (inverted): %d bytes found, correctly outside range %d-%d", byteCount, c.Min, c.Max)
		return true, nil
	}

	// When not inverted, fail if content length is outside the range
	if !isWithinRange {
		c.logger.Sugar().Debugf("Content length validation failed: %d bytes found, expected between %d and %d bytes", byteCount, c.Min, c.Max)
		return false, nil
	}

	c.logger.Sugar().Debugf("Content length validation passed: %d bytes found, within expected range %d-%d", byteCount, c.Min, c.Max)
	return true, nil
}

// buildErrorResponse builds the error response for the ContentLengthGuardrail policy.
func (c *ContentLengthGuardrail) buildErrorResponse(isResponse bool, validationError error) *Result {
	result := NewResult()

	responseBody := make(map[string]interface{})
	responseBody[ErrorCode] = GuardrailAPIMExceptionCode
	responseBody[ErrorType] = ContentLengthGuardrailConstant
	responseBody[ErrorMessage] = c.buildAssessmentObject(isResponse, validationError)

	bodyBytes, err := json.Marshal(responseBody)
	if err != nil {
		c.logger.Error(err, "Error marshaling response body to JSON")
		return result
	}

	result.ImmediateResponse = true
	result.ImmediateResponseCode = v32.StatusCode(GuardrailErrorCode)
	result.ImmediateResponseBody = string(bodyBytes)
	result.ImmediateResponseContentType = "application/json"
	result.StopFurtherProcessing = true

	return result
}

// buildAssessmentObject builds the assessment object for the ContentLengthGuardrail policy.
func (c *ContentLengthGuardrail) buildAssessmentObject(isResponse bool, validationError error) map[string]interface{} {
	c.logger.Sugar().Debugf("Building assessment object for ContentLengthGuardrail policy: %s", c.Name)
	assessment := make(map[string]interface{})
	assessment[AssessmentAction] = "GUARDRAIL_INTERVENED"
	assessment[InterveningGuardrail] = c.Name

	if isResponse {
		assessment[Direction] = "RESPONSE"
	} else {
		assessment[Direction] = "REQUEST"
	}

	// Check if this is a JSONPath extraction error
	if validationError != nil {
		assessment[AssessmentReason] = "Error extracting content from payload using JSONPath."
		if c.ShowAssessment {
			assessmentMessage := "JSONPath extraction failed: " + validationError.Error() + ". Please check the JSONPath configuration: " + c.JSONPath
			assessment[Assessments] = assessmentMessage
		}
	} else {
		assessment[AssessmentReason] = "Violation of applied content length constraints detected."
		if c.ShowAssessment {
			var assessmentMessage string
			if c.Inverted {
				assessmentMessage = "Violation of content length detected. Expected content length to be outside the range of " + strconv.Itoa(c.Min) + " to " + strconv.Itoa(c.Max) + " bytes."
			} else {
				assessmentMessage = "Violation of content length detected. Expected content length to be between " + strconv.Itoa(c.Min) + " and " + strconv.Itoa(c.Max) + " bytes."
			}
			assessment[Assessments] = assessmentMessage
		}
	}
	return assessment
}

// decompressLLMResp decompresses the LLM response if it's compressed.
func (c *ContentLengthGuardrail) decompressLLMResp(payload []byte) (string, error) {
	// Try to detect if it's gzipped by checking for gzip header
	if len(payload) < 2 {
		return string(payload), nil
	}

	// Check for gzip magic numbers
	if payload[0] == 0x1f && payload[1] == 0x8b {
		reader, err := gzip.NewReader(bytes.NewReader(payload))
		if err != nil {
			return string(payload), err // Return original if decompression fails
		}
		defer reader.Close()

		decompressed, err := io.ReadAll(reader)
		if err != nil {
			return string(payload), err // Return original if decompression fails
		}
		return string(decompressed), nil
	}

	// Not compressed, return as is
	return string(payload), nil
}

// extractStringValueFromJsonpath extracts a string value from JSON using JSONPath.
func (c *ContentLengthGuardrail) extractStringValueFromJsonpath(payload []byte, jsonPath string) (string, error) {
	bodyString := string(payload)
	result := gjson.Get(bodyString, removeDollarPrefix(jsonPath))

	if !result.Exists() {
		return "", fmt.Errorf("field not found: %s", jsonPath)
	}

	return result.String(), nil
}
