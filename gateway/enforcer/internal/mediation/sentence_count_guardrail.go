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
	"github.com/tidwall/gjson"
	dpv2alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v2alpha1"
	"github.com/wso2/apk/gateway/enforcer/internal/config"
	"github.com/wso2/apk/gateway/enforcer/internal/logging"
	"github.com/wso2/apk/gateway/enforcer/internal/requestconfig"
)

// SentenceCountGuardrail represents the configuration for Sentence Count Guardrail policy in the API Gateway.
type SentenceCountGuardrail struct {
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
	// SentenceCountGuardrailPolicyKeyName is the key for specifying the name of the guardrail.
	SentenceCountGuardrailPolicyKeyName = "name"
	// SentenceCountGuardrailPolicyKeyMin is the key for specifying the minimum sentence count.
	SentenceCountGuardrailPolicyKeyMin = "min"
	// SentenceCountGuardrailPolicyKeyMax is the key for specifying the maximum sentence count.
	SentenceCountGuardrailPolicyKeyMax = "max"
	// SentenceCountGuardrailPolicyKeyJSONPath is the key for specifying the JSON path to extract content.
	SentenceCountGuardrailPolicyKeyJSONPath = "jsonPath"
	// SentenceCountGuardrailPolicyKeyInverted is the key for specifying if the validation should be inverted.
	SentenceCountGuardrailPolicyKeyInverted = "invert"
	// SentenceCountGuardrailPolicyKeyShowAssessment is the key for specifying if assessment should be shown.
	SentenceCountGuardrailPolicyKeyShowAssessment = "showAssessment"

	// SentenceCountGuardrailConstant is the constant for sentence count guardrail errors.
	SentenceCountGuardrailConstant = "SENTENCE_COUNT_GUARDRAIL"
)

var (
	// SentenceSplitRegexCompiled is a compiled regex for splitting sentences
	SentenceSplitRegexCompiled = regexp.MustCompile(`[.!?]+`)
)

// NewSentenceCountGuardrail creates a new SentenceCountGuardrail instance.
func NewSentenceCountGuardrail(mediation *dpv2alpha1.Mediation) *SentenceCountGuardrail {
	cfg := config.GetConfig()
	logger := cfg.Logger

	name := "SentenceCountGuardrail"
	if val, ok := extractPolicyValue(mediation.Parameters, SentenceCountGuardrailPolicyKeyName); ok {
		name = val
	}

	minL := 0
	if val, ok := extractPolicyValue(mediation.Parameters, SentenceCountGuardrailPolicyKeyMin); ok {
		if intValue, err := strconv.Atoi(val); err == nil {
			minL = intValue
		}
	}

	maxL := 100
	if val, ok := extractPolicyValue(mediation.Parameters, SentenceCountGuardrailPolicyKeyMax); ok {
		if intValue, err := strconv.Atoi(val); err == nil {
			maxL = intValue
		}
	}

	jsonPath := "$.content"
	if val, ok := extractPolicyValue(mediation.Parameters, SentenceCountGuardrailPolicyKeyJSONPath); ok {
		jsonPath = val
	}

	inverted := false
	if val, ok := extractPolicyValue(mediation.Parameters, SentenceCountGuardrailPolicyKeyInverted); ok {
		inverted = val == "true"
	}

	showAssessment := false
	if val, ok := extractPolicyValue(mediation.Parameters, SentenceCountGuardrailPolicyKeyShowAssessment); ok {
		showAssessment = val == "true"
	}

	return &SentenceCountGuardrail{
		PolicyName:     "SentenceCountGuardrail",
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

// Process processes the request configuration for Sentence Count Guardrail.
func (s *SentenceCountGuardrail) Process(requestConfig *requestconfig.Holder) *Result {
	result := NewResult()

	// Handle request body processing
	if requestConfig.ProcessingPhase == requestconfig.ProcessingPhaseRequestBody {
		s.logger.Sugar().Debugf("Beginning request payload validation for SentenceCountGuardrail policy: %s", s.Name)

		if requestConfig.RequestBody == nil || requestConfig.RequestBody.Body == nil {
			s.logger.Sugar().Debug("No request body found, skipping sentence count validation")
			return result
		}

		validationResult, err := s.validatePayload(requestConfig.RequestBody.Body, false)
		if !validationResult {
			s.logger.Sugar().Debugf("Request payload validation failed for SentenceCountGuardrail policy: %s", s.Name)
			return s.buildErrorResponse(false, err)
		}
		s.logger.Sugar().Debugf("Request payload validation passed for SentenceCountGuardrail policy: %s", s.Name)
		return result
	}

	// Handle response body processing
	if requestConfig.ProcessingPhase == requestconfig.ProcessingPhaseResponseBody {
		s.logger.Sugar().Debugf("Beginning response body validation for SentenceCountGuardrail policy: %s", s.Name)

		if requestConfig.ResponseBody == nil || requestConfig.ResponseBody.Body == nil {
			s.logger.Sugar().Debug("No response body found, skipping sentence count validation")
			return result
		}

		validationResult, err := s.validatePayload(requestConfig.ResponseBody.Body, true)
		if !validationResult {
			s.logger.Sugar().Debugf("Response body validation failed for SentenceCountGuardrail policy: %s", s.Name)
			return s.buildErrorResponse(true, err)
		}
		s.logger.Sugar().Debugf("Response body validation passed for SentenceCountGuardrail policy: %s", s.Name)
		return result
	}

	return result
}

// validatePayload validates the payload against the SentenceCountGuardrail policy.
func (s *SentenceCountGuardrail) validatePayload(payload []byte, isResponse bool) (bool, error) {
	// Decompress response body if needed
	if isResponse {
		bodyStr, err := s.decompressLLMResp(payload)
		if err == nil {
			payload = []byte(bodyStr)
		}
	}

	// Validate sentence count range
	if (s.Min > s.Max) || (s.Min < 0) || (s.Max <= 0) {
		s.logger.Sugar().Errorf("Invalid sentence count range: min=%d, max=%d", s.Min, s.Max)
		return false, nil
	}

	// Extract value from JSON using JSONPath
	extractedValue, err := s.extractStringValueFromJsonpath(payload, s.JSONPath)
	if err != nil {
		s.logger.Error(err, "Error extracting value from JSON using JSONPath")
		return false, err
	}

	// Trim the extracted text
	extractedValue = strings.TrimSpace(extractedValue)

	// Split into sentences and count non-empty sentences
	sentences := SentenceSplitRegexCompiled.Split(extractedValue, -1)
	s.logger.Sugar().Debugf("Extracted value for SentenceCountGuardrail policy: %v", sentences)
	sentenceCount := 0
	for _, sentence := range sentences {
		// Clean and trim each sentence, then check if it's non-empty
		cleanedSentence := TextCleanRegexCompiled.ReplaceAllString(sentence, "")
		cleanedSentence = strings.TrimSpace(cleanedSentence)
		if cleanedSentence != "" {
			sentenceCount++
		}
	}

	isWithinRange := sentenceCount >= s.Min && sentenceCount <= s.Max

	if s.Inverted {
		// When inverted, fail if sentence count is within the range
		if isWithinRange {
			s.logger.Sugar().Debugf("Sentence count validation failed (inverted): %d sentences found, should NOT be between %d and %d sentences", sentenceCount, s.Min, s.Max)
			return false, nil
		}
		s.logger.Sugar().Debugf("Sentence count validation passed (inverted): %d sentences found, correctly outside range %d-%d", sentenceCount, s.Min, s.Max)
		return true, nil
	}

	// When not inverted, fail if sentence count is outside the range
	if !isWithinRange {
		s.logger.Sugar().Debugf("Sentence count validation failed: %d sentences found, expected between %d and %d sentences", sentenceCount, s.Min, s.Max)
		return false, nil
	}

	s.logger.Sugar().Debugf("Sentence count validation passed: %d sentences found, within expected range %d-%d", sentenceCount, s.Min, s.Max)
	return true, nil
}

// buildErrorResponse builds the error response for the SentenceCountGuardrail policy.
func (s *SentenceCountGuardrail) buildErrorResponse(isResponse bool, validationError error) *Result {
	result := NewResult()

	responseBody := make(map[string]interface{})
	responseBody[ErrorCode] = GuardrailAPIMExceptionCode
	responseBody[ErrorType] = SentenceCountGuardrailConstant
	responseBody[ErrorMessage] = s.buildAssessmentObject(isResponse, validationError)

	bodyBytes, err := json.Marshal(responseBody)
	if err != nil {
		s.logger.Error(err, "Error marshaling response body to JSON")
		return result
	}

	result.ImmediateResponse = true
	result.ImmediateResponseCode = v32.StatusCode(GuardrailErrorCode)
	result.ImmediateResponseBody = string(bodyBytes)
	result.ImmediateResponseContentType = "application/json"
	result.StopFurtherProcessing = true

	return result
}

// buildAssessmentObject builds the assessment object for the SentenceCountGuardrail policy.
func (s *SentenceCountGuardrail) buildAssessmentObject(isResponse bool, validationError error) map[string]interface{} {
	s.logger.Sugar().Debugf("Building assessment object for SentenceCountGuardrail policy: %s", s.Name)
	assessment := make(map[string]interface{})
	assessment[AssessmentAction] = "GUARDRAIL_INTERVENED"
	assessment[InterveningGuardrail] = s.Name

	if isResponse {
		assessment[Direction] = "RESPONSE"
	} else {
		assessment[Direction] = "REQUEST"
	}

	// Check if this is a JSONPath extraction error
	if validationError != nil {
		assessment[AssessmentReason] = "Error extracting content from payload using JSONPath."
		if s.ShowAssessment {
			assessmentMessage := "JSONPath extraction failed: " + validationError.Error() + ". Please check the JSONPath configuration: " + s.JSONPath
			assessment[Assessments] = assessmentMessage
		}
	} else {
		assessment[AssessmentReason] = "Violation of applied sentence count constraints detected."
		if s.ShowAssessment {
			var assessmentMessage string
			if s.Inverted {
				assessmentMessage = "Violation of sentence count detected. Expected sentence count to be outside the range of " + strconv.Itoa(s.Min) + " to " + strconv.Itoa(s.Max) + " sentences."
			} else {
				assessmentMessage = "Violation of sentence count detected. Expected sentence count to be between " + strconv.Itoa(s.Min) + " and " + strconv.Itoa(s.Max) + " sentences."
			}
			assessment[Assessments] = assessmentMessage
		}
	}
	return assessment
}

// decompressLLMResp decompresses the LLM response if it's compressed.
func (s *SentenceCountGuardrail) decompressLLMResp(payload []byte) (string, error) {
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
func (s *SentenceCountGuardrail) extractStringValueFromJsonpath(payload []byte, jsonPath string) (string, error) {
	bodyString := string(payload)
	result := gjson.Get(bodyString, removeDollarPrefix(jsonPath))

	if !result.Exists() {
		return "", fmt.Errorf("field not found: %s", jsonPath)
	}

	return result.String(), nil
}
