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
	"strconv"
	"strings"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_service_proc_v3 "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"
	v32 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/wso2/apk/gateway/enforcer/internal/dto"
	"github.com/wso2/apk/gateway/enforcer/internal/logging"
)

// SentenceCountGuardrail is a struct that represents a sentence count guardrail policy.
type SentenceCountGuardrail struct {
	dto.BaseInBuiltPolicy
	Name           string
	Min            int
	Max            int
	JSONPath       string
	Inverted       bool
	ShowAssessment bool
}

// HandleRequestBody is a method that implements the mediation logic for the SentenceCountGuardrail policy on request.
func (r *SentenceCountGuardrail) HandleRequestBody(logger *logging.Logger, req *envoy_service_proc_v3.ProcessingRequest, resp *envoy_service_proc_v3.ProcessingResponse, props map[string]interface{}) *envoy_service_proc_v3.ProcessingResponse {
	logger.Sugar().Debugf("Beginning request payload validation for SentenceCountGuardrail policy: %s", r.Name)
	validationResult, err := r.validatePayload(logger, req.GetRequestBody().Body, false)
	if !validationResult {
		logger.Sugar().Debugf("Request payload validation failed for SentenceCountGuardrail policy: %s", r.Name)
		return r.buildResponse(logger, false, err)
	}
	logger.Sugar().Debugf("Request payload validation passed for SentenceCountGuardrail policy: %s", r.Name)
	return nil
}

// HandleResponseBody is a method that implements the mediation logic for the SentenceCountGuardrail policy on response.
func (r *SentenceCountGuardrail) HandleResponseBody(logger *logging.Logger, req *envoy_service_proc_v3.ProcessingRequest, resp *envoy_service_proc_v3.ProcessingResponse, props map[string]interface{}) *envoy_service_proc_v3.ProcessingResponse {
	logger.Sugar().Debugf("Beginning response body validation for SentenceCountGuardrail policy: %s", r.Name)
	validationResult, err := r.validatePayload(logger, req.GetResponseBody().Body, true)
	if !validationResult {
		logger.Sugar().Debugf("Response body validation failed for SentenceCountGuardrail policy: %s", r.Name)
		return r.buildResponse(logger, true, err)
	}
	logger.Sugar().Debugf("Response body validation passed for SentenceCountGuardrail policy: %s", r.Name)
	return nil
}

// validatePayload validates the payload against the SentenceCountGuardrail policy.
func (r *SentenceCountGuardrail) validatePayload(logger *logging.Logger, payload []byte, isResponse bool) (bool, error) {
	if isResponse {
		bodyStr, _, err := DecompressLLMResp(payload)
		if err == nil {
			payload = []byte(bodyStr)
		}
	}
	if (r.Min > r.Max) || (r.Min < 0) || (r.Max <= 0) {
		logger.Sugar().Errorf("Invalid sentence count range: min=%d, max=%d", r.Min, r.Max)
		return false, nil
	}

	extractedValue, err := ExtractStringValueFromJsonpath(logger, payload, r.JSONPath)
	if err != nil {
		logger.Error(err, "Error extracting value from JSON using JSONPath")
		return false, err
	}

	// Clean and trim
	extractedValue = TextCleanRegexCompiled.ReplaceAllString(extractedValue, "")
	extractedValue = strings.TrimSpace(extractedValue)

	// Split into sentences and count non-empty
	sentences := SentenceSplitRegexCompiled.Split(extractedValue, -1)
	logger.Sugar().Debugf("Extracted value for SentenceCountGuardrail policy: %v", sentences)
	sentenceCount := 0
	for _, w := range sentences {
		if w != "" {
			sentenceCount++
		}
	}

	isWithinRange := sentenceCount >= r.Min && sentenceCount <= r.Max
	
	if r.Inverted {
		// When inverted, fail if sentence count is within the range
		if isWithinRange {
			logger.Sugar().Debugf("Sentence count validation failed (inverted): %d sentences found, should NOT be between %d and %d sentences", sentenceCount, r.Min, r.Max)
			return false, nil
		}
		logger.Sugar().Debugf("Sentence count validation passed (inverted): %d sentences found, correctly outside range %d-%d", sentenceCount, r.Min, r.Max)
		return true, nil
	}
	// When not inverted, fail if sentence count is outside the range
	if !isWithinRange {
		logger.Sugar().Debugf("Sentence count validation failed: %d sentences found, expected between %d and %d sentences", sentenceCount, r.Min, r.Max)
		return false, nil
	}
	logger.Sugar().Debugf("Sentence count validation passed: %d sentences found, within expected range %d-%d", sentenceCount, r.Min, r.Max)
	return true, nil
}

// buildResponse is a method that builds the response body for the SentenceCountGuardrail policy.
func (r *SentenceCountGuardrail) buildResponse(logger *logging.Logger, isResponse bool, validationError error) *envoy_service_proc_v3.ProcessingResponse {
	responseBody := make(map[string]interface{})
	responseBody[ErrorCode] = GuardrailAPIMExceptionCode
	responseBody[ErrorType] = SentenceCountGuardrailConstant
	responseBody[ErrorMessage] = r.buildAssessmentObject(logger, isResponse, validationError)

	bodyBytes, err := json.Marshal(responseBody)
	if err != nil {
		logger.Error(err, "Error marshaling response body to JSON")
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
					Code: v32.StatusCode(GuardrailErrorCode),
				},
				Body:    bodyBytes,
				Headers: headers,
			},
		},
	}
}

// buildAssessmentObject is a method that builds the assessment object for the SentenceCountGuardrail policy.
func (r *SentenceCountGuardrail) buildAssessmentObject(logger *logging.Logger, isResponse bool, validationError error) map[string]interface{} {
	logger.Sugar().Debugf("Building assessment object for SentenceCountGuardrail policy: %s", r.Name)
	assessment := make(map[string]interface{})
	assessment[AssessmentAction] = "GUARDRAIL_INTERVENED"
	assessment[InterveningGuardrail] = r.Name
	if isResponse {
		assessment[Direction] = "RESPONSE"
	} else {
		assessment[Direction] = "REQUEST"
	}
	
	// Check if this is a JSONPath extraction error
	if validationError != nil {
		assessment[AssessmentReason] = "Error extracting content from payload using JSONPath."
		if r.ShowAssessment {
			assessmentMessage := "JSONPath extraction failed: " + validationError.Error() + ". Please check the JSONPath configuration: " + r.JSONPath
			assessment[Assessments] = assessmentMessage
		}
	} else {
		assessment[AssessmentReason] = "Violation of applied sentence count constraints detected."
		if r.ShowAssessment {
			var assessmentMessage string
			if r.Inverted {
				assessmentMessage = "Violation of sentence count detected. Expected sentence count to be outside the range of " + strconv.Itoa(r.Min) + " to " + strconv.Itoa(r.Max) + " sentences."
			} else {
				assessmentMessage = "Violation of sentence count detected. Expected sentence count to be between " + strconv.Itoa(r.Min) + " and " + strconv.Itoa(r.Max) + " sentences."
			}
			assessment[Assessments] = assessmentMessage
		}
	}
	return assessment
}

// NewSentenceCountGuardrail initializes the SentenceCountGuardrail policy from the given InBuiltPolicy.
func NewSentenceCountGuardrail(inBuiltPolicy dto.InBuiltPolicy) *SentenceCountGuardrail {
	sentenceCountGuardrail := &SentenceCountGuardrail{
		BaseInBuiltPolicy: dto.BaseInBuiltPolicy{
			PolicyName:    inBuiltPolicy.GetPolicyName(),
			PolicyID:      inBuiltPolicy.GetPolicyID(),
			PolicyVersion: inBuiltPolicy.GetPolicyVersion(),
			Parameters:    inBuiltPolicy.GetParameters(),
			PolicyOrder:   inBuiltPolicy.GetPolicyOrder(),
		},
	}

	for key, value := range inBuiltPolicy.GetParameters() {
		switch key {
		case "name":
			sentenceCountGuardrail.Name = value
		case "min":
			if intValue, err := strconv.Atoi(value); err == nil {
				sentenceCountGuardrail.Min = intValue
			}
		case "max":
			if intValue, err := strconv.Atoi(value); err == nil {
				sentenceCountGuardrail.Max = intValue
			}
		case "jsonPath":
			sentenceCountGuardrail.JSONPath = value
		case "invert":
			sentenceCountGuardrail.Inverted = value == "true"
		case "showAssessment":
			sentenceCountGuardrail.ShowAssessment = value == "true"
		}
	}
	return sentenceCountGuardrail
}
