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
	"strconv"
	"strings"

	v32 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	dpv2alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v2alpha1"
	"github.com/tidwall/gjson"
	"github.com/wso2/apk/gateway/enforcer/internal/config"
	"github.com/wso2/apk/gateway/enforcer/internal/logging"
	"github.com/wso2/apk/gateway/enforcer/internal/requestconfig"
)

// WordCountGuardrail represents the configuration for Word Count Guardrail policy in the API Gateway.
type WordCountGuardrail struct {
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
	// WordCountGuardrailPolicyKeyName is the key for specifying the name of the guardrail.
	WordCountGuardrailPolicyKeyName = "name"
	// WordCountGuardrailPolicyKeyMin is the key for specifying the minimum word count.
	WordCountGuardrailPolicyKeyMin = "min"
	// WordCountGuardrailPolicyKeyMax is the key for specifying the maximum word count.
	WordCountGuardrailPolicyKeyMax = "max"
	// WordCountGuardrailPolicyKeyJSONPath is the key for specifying the JSON path to extract content.
	WordCountGuardrailPolicyKeyJSONPath = "jsonPath"
	// WordCountGuardrailPolicyKeyInverted is the key for specifying if the validation should be inverted.
	WordCountGuardrailPolicyKeyInverted = "invert"
	// WordCountGuardrailPolicyKeyShowAssessment is the key for specifying if assessment should be shown.
	WordCountGuardrailPolicyKeyShowAssessment = "showAssessment"

	// GuardrailAPIMExceptionCode is the error code for guardrail exceptions.
	GuardrailAPIMExceptionCode   = "GUARDRAIL_API_EXCEPTION"
	// WordCountGuardrailConstant is the constant for word count guardrail errors.
	WordCountGuardrailConstant   = "WORD_COUNT_GUARDRAIL"
	// GuardrailErrorCode is the HTTP status code for guardrail errors.
	GuardrailErrorCode           = 400
	// ErrorCode response constants
	ErrorCode                    = "errorCode"
	// ErrorType response constants
	ErrorType                    = "errorType"
	// ErrorMessage response constants
	ErrorMessage                 = "errorMessage"
	// AssessmentAction response constants
	AssessmentAction             = "action"
	// InterveningGuardrail response constants
	InterveningGuardrail         = "interveningGuardrail"
	// Direction response constants
	Direction                    = "direction"
	// AssessmentReason response constants
	AssessmentReason             = "reason"
	// Assessments response constants
	Assessments                  = "assessments"
)

var (
	// TextCleanRegexCompiled is a compiled regex for cleaning text
	TextCleanRegexCompiled = regexp.MustCompile(`[^\w\s]`)
	// WordSplitRegexCompiled is a compiled regex for splitting words
	WordSplitRegexCompiled = regexp.MustCompile(`\s+`)
)

// NewWordCountGuardrail creates a new WordCountGuardrail instance.
func NewWordCountGuardrail(mediation *dpv2alpha1.Mediation) *WordCountGuardrail {
	cfg := config.GetConfig()
	logger := cfg.Logger

	name := "WordCountGuardrail"
	if val, ok := extractPolicyValue(mediation.Parameters, WordCountGuardrailPolicyKeyName); ok {
		name = val
	}

	minL := 0
	if val, ok := extractPolicyValue(mediation.Parameters, WordCountGuardrailPolicyKeyMin); ok {
		if intValue, err := strconv.Atoi(val); err == nil {
			minL = intValue
		}
	}

	maxL := 100
	if val, ok := extractPolicyValue(mediation.Parameters, WordCountGuardrailPolicyKeyMax); ok {
		if intValue, err := strconv.Atoi(val); err == nil {
			maxL = intValue
		}
	}

	jsonPath := "$.content"
	if val, ok := extractPolicyValue(mediation.Parameters, WordCountGuardrailPolicyKeyJSONPath); ok {
		jsonPath = val
	}

	inverted := false
	if val, ok := extractPolicyValue(mediation.Parameters, WordCountGuardrailPolicyKeyInverted); ok {
		inverted = val == "true"
	}

	showAssessment := false
	if val, ok := extractPolicyValue(mediation.Parameters, WordCountGuardrailPolicyKeyShowAssessment); ok {
		showAssessment = val == "true"
	}

	return &WordCountGuardrail{
		PolicyName:     "WordCountGuardrail",
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

// Process processes the request configuration for Word Count Guardrail.
func (w *WordCountGuardrail) Process(requestConfig *requestconfig.Holder) *Result {
	result := NewResult()

	// Handle request body processing
	if requestConfig.ProcessingPhase == requestconfig.ProcessingPhaseRequestBody {
		w.logger.Sugar().Debugf("Beginning request payload validation for WordCountGuardrail policy: %s", w.Name)
		
		if requestConfig.RequestBody == nil || requestConfig.RequestBody.Body == nil {
			w.logger.Sugar().Debug("No request body found, skipping word count validation")
			return result
		}

		validationResult, err := w.validatePayload(requestConfig.RequestBody.Body, false)
		if !validationResult {
			w.logger.Sugar().Debugf("Request payload validation failed for WordCountGuardrail policy: %s", w.Name)
			return w.buildErrorResponse(false, err)
		}
		w.logger.Sugar().Debugf("Request payload validation passed for WordCountGuardrail policy: %s", w.Name)
		return result
	}

	// Handle response body processing
	if requestConfig.ProcessingPhase == requestconfig.ProcessingPhaseResponseBody {
		w.logger.Sugar().Debugf("Beginning response body validation for WordCountGuardrail policy: %s", w.Name)

		if requestConfig.ResponseBody == nil || requestConfig.ResponseBody.Body == nil {
			w.logger.Sugar().Debug("No response body found, skipping word count validation")
			return result
		}

		validationResult, err := w.validatePayload(requestConfig.ResponseBody.Body, true)
		if !validationResult {
			w.logger.Sugar().Debugf("Response body validation failed for WordCountGuardrail policy: %s", w.Name)
			return w.buildErrorResponse(true, err)
		}
		w.logger.Sugar().Debugf("Response body validation passed for WordCountGuardrail policy: %s", w.Name)
		return result
	}

	return result
}

// validatePayload validates the payload against the WordCountGuardrail policy.
func (w *WordCountGuardrail) validatePayload(payload []byte, isResponse bool) (bool, error) {
	// Decompress response body if needed
	if isResponse {
		bodyStr, err := w.decompressLLMResp(payload)
		if err == nil {
			payload = []byte(bodyStr)
		}
	}

	// Validate word count range
	if (w.Min > w.Max) || (w.Min < 0) || (w.Max <= 0) {
		w.logger.Sugar().Errorf("Invalid word count range: min=%d, max=%d", w.Min, w.Max)
		return false, nil
	}

	// Extract value from JSON using JSONPath
	extractedValue, err := w.extractStringValueFromJsonpath(payload, w.JSONPath)
	if err != nil {
		w.logger.Error(err, "Error extracting value from JSON using JSONPath")
		return false, err
	}

	// Clean and trim the extracted text
	extractedValue = TextCleanRegexCompiled.ReplaceAllString(extractedValue, "")
	extractedValue = strings.TrimSpace(extractedValue)

	// Split into words and count non-empty words
	words := WordSplitRegexCompiled.Split(extractedValue, -1)
	wordCount := 0
	for _, word := range words {
		if word != "" {
			wordCount++
		}
	}

	isWithinRange := wordCount >= w.Min && wordCount <= w.Max

	if w.Inverted {
		// When inverted, fail if word count is within the range
		if isWithinRange {
			w.logger.Sugar().Debugf("Word count validation failed (inverted): %d words found, should NOT be between %d and %d words", wordCount, w.Min, w.Max)
			return false, nil
		}
		w.logger.Sugar().Debugf("Word count validation passed (inverted): %d words found, correctly outside range %d-%d", wordCount, w.Min, w.Max)
		return true, nil
	}

	// When not inverted, fail if word count is outside the range
	if !isWithinRange {
		w.logger.Sugar().Debugf("Word count validation failed: %d words found, expected between %d and %d words", wordCount, w.Min, w.Max)
		return false, nil
	}

	w.logger.Sugar().Debugf("Word count validation passed: %d words found, within expected range %d-%d", wordCount, w.Min, w.Max)
	return true, nil
}

// buildErrorResponse builds the error response for the WordCountGuardrail policy.
func (w *WordCountGuardrail) buildErrorResponse(isResponse bool, validationError error) *Result {
	result := NewResult()

	responseBody := make(map[string]interface{})
	responseBody[ErrorCode] = GuardrailAPIMExceptionCode
	responseBody[ErrorType] = WordCountGuardrailConstant
	responseBody[ErrorMessage] = w.buildAssessmentObject(isResponse, validationError)

	bodyBytes, err := json.Marshal(responseBody)
	if err != nil {
		w.logger.Error(err, "Error marshaling response body to JSON")
		return result
	}

	result.ImmediateResponse = true
	result.ImmediateResponseCode = v32.StatusCode(GuardrailErrorCode)
	result.ImmediateResponseBody = string(bodyBytes)
	result.ImmediateResponseContentType = "application/json"
	result.StopFurtherProcessing = true

	return result
}

// buildAssessmentObject builds the assessment object for the WordCountGuardrail policy.
func (w *WordCountGuardrail) buildAssessmentObject(isResponse bool, validationError error) map[string]interface{} {
	w.logger.Sugar().Debugf("Building assessment object for WordCountGuardrail policy: %s", w.Name)
	assessment := make(map[string]interface{})
	assessment[AssessmentAction] = "GUARDRAIL_INTERVENED"
	assessment[InterveningGuardrail] = w.Name

	if isResponse {
		assessment[Direction] = "RESPONSE"
	} else {
		assessment[Direction] = "REQUEST"
	}

	// Check if this is a JSONPath extraction error
	if validationError != nil {
		assessment[AssessmentReason] = "Error extracting content from payload using JSONPath."
		if w.ShowAssessment {
			assessmentMessage := "JSONPath extraction failed: " + validationError.Error() + ". Please check the JSONPath configuration: " + w.JSONPath
			assessment[Assessments] = assessmentMessage
		}
	} else {
		assessment[AssessmentReason] = "Violation of applied word count constraints detected."
		if w.ShowAssessment {
			var assessmentMessage string
			if w.Inverted {
				assessmentMessage = "Violation of word count detected. Expected word count to be outside the range of " + strconv.Itoa(w.Min) + " to " + strconv.Itoa(w.Max) + " words."
			} else {
				assessmentMessage = "Violation of word count detected. Expected word count to be between " + strconv.Itoa(w.Min) + " and " + strconv.Itoa(w.Max) + " words."
			}
			assessment[Assessments] = assessmentMessage
		}
	}
	return assessment
}

// decompressLLMResp decompresses the LLM response if it's compressed.
func (w *WordCountGuardrail) decompressLLMResp(payload []byte) (string, error) {
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
func (w *WordCountGuardrail) extractStringValueFromJsonpath(payload []byte, jsonPath string) (string, error) {
	bodyString := string(payload)
	result := gjson.Get(bodyString, removeDollarPrefix(jsonPath))
	
	if !result.Exists() {
		return "", fmt.Errorf("field not found: %s", jsonPath)
	}

	return result.String(), nil
}
