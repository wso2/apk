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

// WordCountGuardrail is a struct that represents a word count guardrail policy.
type WordCountGuardrail struct {
	dto.BaseInBuiltPolicy
	Name           string
	Min            int
	Max            int
	JSONPath       string
	Inverted       bool
	ShowAssessment bool
}

// HandleRequestBody is a method that implements the mediation logic for the WordCountGuardrail policy on request.
func (r *WordCountGuardrail) HandleRequestBody(logger *logging.Logger, req *envoy_service_proc_v3.ProcessingRequest, resp *envoy_service_proc_v3.ProcessingResponse, props map[string]interface{}) *envoy_service_proc_v3.ProcessingResponse {
	logger.Sugar().Debugf("Beginning request payload validation for WordCountGuardrail policy: %s", r.Name)
	validationResult := r.validatePayload(logger, req.GetRequestBody().Body)
	if !validationResult {
		logger.Sugar().Debugf("Request payload validation failed for WordCountGuardrail policy: %s", r.Name)
		return r.buildResponse(logger, false)
	}
	logger.Sugar().Debugf("Request payload validation passed for WordCountGuardrail policy: %s", r.Name)
	return nil
}

// HandleResponseBody is a method that implements the mediation logic for the WordCountGuardrail policy on response.
func (r *WordCountGuardrail) HandleResponseBody(logger *logging.Logger, req *envoy_service_proc_v3.ProcessingRequest, resp *envoy_service_proc_v3.ProcessingResponse, props map[string]interface{}) *envoy_service_proc_v3.ProcessingResponse {
	logger.Sugar().Debugf("Beginning response body validation for WordCountGuardrail policy: %s", r.Name)
	validationResult := r.validatePayload(logger, req.GetResponseBody().Body)
	if !validationResult {
		logger.Sugar().Debugf("Response body validation failed for WordCountGuardrail policy: %s", r.Name)
		return r.buildResponse(logger, true)
	}
	logger.Sugar().Debugf("Response body validation passed for WordCountGuardrail policy: %s", r.Name)
	return nil
}

// validatePayload validates the payload against the WordCountGuardrail policy.
func (r *WordCountGuardrail) validatePayload(logger *logging.Logger, payload []byte) bool {
	if (r.Min > r.Max) || (r.Min < 0) || (r.Max <= 0) {
		logger.Sugar().Errorf("Invalid word count range: min=%d, max=%d", r.Min, r.Max)
		return false
	}

	extractedValue, err := ExtractStringValueFromJsonpath(logger, payload, r.JSONPath)
	if err != nil {
		logger.Error(err, "Error extracting value from JSON using JSONPath")
		return false
	}

	// Clean and trim
	extractedValue = TextCleanRegexCompiled.ReplaceAllString(extractedValue, "")
	extractedValue = strings.TrimSpace(extractedValue)

	// Split into words and count non-empty
	words := WordSplitRegexCompiled.Split(extractedValue, -1)
	wordCount := 0
	for _, w := range words {
		if w != "" {
			wordCount++
		}
	}

	if wordCount < r.Min || wordCount > r.Max {
		logger.Sugar().Debugf("Word count validation failed: %d words found, expected between %d and %d words", wordCount, r.Min, r.Max)
		if r.Inverted {
			logger.Sugar().Debugf("Inverted condition is true, returning true")
			return true
		}
		logger.Sugar().Debugf("Inverted condition is false, returning false")
		return false
	}
	return true
}

// buildResponse is a method that builds the response body for the WordCountGuardrail policy.
func (r *WordCountGuardrail) buildResponse(logger *logging.Logger, isResponse bool) *envoy_service_proc_v3.ProcessingResponse {
	responseBody := make(map[string]interface{})
	responseBody[ErrorCode] = GuardrailAPIMExceptionCode
	responseBody[ErrorType] = WordCountGuardrailConstant
	responseBody[ErrorMessage] = r.buildAssessmentObject(logger, isResponse)

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

// buildAssessmentObject is a method that builds the assessment object for the WordCountGuardrail policy.
func (r *WordCountGuardrail) buildAssessmentObject(logger *logging.Logger, isResponse bool) map[string]interface{} {
	logger.Sugar().Debugf("Building assessment object for WordCountGuardrail policy: %s", r.Name)
	assessment := make(map[string]interface{})
	assessment[AssessmentAction] = "GUARDRAIL_INTERVENED"
	assessment[InterveningGuardrail] = r.Name
	if isResponse {
		assessment[Direction] = "RESPONSE"
	} else {
		assessment[Direction] = "REQUEST"
	}
	assessment[AssessmentReason] = "Violation of applied word count constraints detected."

	if r.ShowAssessment {
		var minStr, maxStr string
		if r.Inverted {
			minStr = "less than"
			maxStr = "or more than"
		} else {
			minStr = "between"
			maxStr = "and"
		}
		assessment[Assessments] = "Violation of word count detected. Expected " + minStr + " " + strconv.Itoa(r.Min) + " " + maxStr + " " + strconv.Itoa(r.Max) + " words."
	}
	return assessment
}

// NewWordCountGuardrail initializes the WordCountGuardrail policy from the given InBuiltPolicy.
func NewWordCountGuardrail(inBuiltPolicy dto.InBuiltPolicy) *WordCountGuardrail {
	wordCountGuardrail := &WordCountGuardrail{
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
			wordCountGuardrail.Name = value
		case "min":
			if intValue, err := strconv.Atoi(value); err == nil {
				wordCountGuardrail.Min = intValue
			}
		case "max":
			if intValue, err := strconv.Atoi(value); err == nil {
				wordCountGuardrail.Max = intValue
			}
		case "jsonPath":
			wordCountGuardrail.JSONPath = value
		case "invert":
			wordCountGuardrail.Inverted = value == "true"
		case "showAssessment":
			wordCountGuardrail.ShowAssessment = value == "true"
		}
	}
	return wordCountGuardrail
}
