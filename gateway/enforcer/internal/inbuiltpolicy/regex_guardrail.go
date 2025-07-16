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
	"regexp"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_service_proc_v3 "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"
	v32 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/wso2/apk/gateway/enforcer/internal/dto"
	"github.com/wso2/apk/gateway/enforcer/internal/logging"
)

// RegexGuardrail is a struct that represents a regex guardrail policy.
type RegexGuardrail struct {
	dto.BaseInBuiltPolicy
	Name           string
	Regex          string
	JSONPath       string
	Inverted       bool
	ShowAssessment bool
}

// HandleRequestBody is a method that implements the mediation logic for the RegexGuardrail policy on request.
func (r *RegexGuardrail) HandleRequestBody(logger *logging.Logger, req *envoy_service_proc_v3.ProcessingRequest, resp *envoy_service_proc_v3.ProcessingResponse, props map[string]interface{}) *envoy_service_proc_v3.ProcessingResponse {
	logger.Sugar().Debugf("Beginning request payload validation for RegexGuardrail policy: %s", r.Name)
	validationResult := r.validatePayload(logger, req.GetRequestBody().Body)
	if !validationResult {
		logger.Sugar().Debugf("Request payload validation failed for RegexGuardrail policy: %s", r.Name)
		return r.buildResponse(logger, false)
	}
	logger.Sugar().Debugf("Request payload validation passed for RegexGuardrail policy: %s", r.Name)
	return nil
}

// HandleResponseBody is a method that implements the mediation logic for the RegexGuardrail policy on response.
func (r *RegexGuardrail) HandleResponseBody(logger *logging.Logger, req *envoy_service_proc_v3.ProcessingRequest, resp *envoy_service_proc_v3.ProcessingResponse, props map[string]interface{}) *envoy_service_proc_v3.ProcessingResponse {
	logger.Sugar().Debugf("Beginning response body validation for RegexGuardrail policy: %s", r.Name)
	validationResult := r.validatePayload(logger, req.GetResponseBody().Body)
	if !validationResult {
		logger.Sugar().Debugf("Response body validation failed for RegexGuardrail policy: %s", r.Name)
		return r.buildResponse(logger, true)
	}
	logger.Sugar().Debugf("Response body validation passed for RegexGuardrail policy: %s", r.Name)
	return nil
}

// validatePayload is a method that returns the name of the policy for validation purposes.
func (r *RegexGuardrail) validatePayload(logger *logging.Logger, payload []byte) bool {
	extractedValue, err := ExtractStringValueFromJsonpath(logger, payload, r.JSONPath)
	if err != nil {
		logger.Error(err, "Error extracting value from JSON using JSONPath")
		return false
	}
	// Perform regex matching
	matched, err := regexp.MatchString(r.Regex, extractedValue)
	if err != nil {
		logger.Error(err, "Error matching regex")
		return false
	}
	if matched && r.Inverted {
		logger.Sugar().Debugf("Regex matched and inverted condition is true, returning false")
		return false
	} else if !matched && !r.Inverted {
		logger.Sugar().Debugf("Regex did not match and inverted condition is false, returning false")
		return false
	}
	logger.Sugar().Debugf("Regex matched successfully, returning true")
	return true
}

// buildResponse is a method that builds the response body for the RegexGuardrail policy.
func (r *RegexGuardrail) buildResponse(logger *logging.Logger, isResponse bool) *envoy_service_proc_v3.ProcessingResponse {
	responseBody := make(map[string]interface{})
	responseBody[ErrorCode] = GuardrailAPIMExceptionCode
	responseBody[ErrorType] = RegexGuardrailConstant
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

// buildAssessmentObject is a method that builds the assessment object for the RegexGuardrail policy.
func (r *RegexGuardrail) buildAssessmentObject(logger *logging.Logger, isResponse bool) map[string]interface{} {
	logger.Sugar().Debugf("Building assessment object for RegexGuardrail policy: %s", r.Name)
	assessment := make(map[string]interface{})
	assessment[AssessmentAction] = "GUARDRAIL_INTERVENED"
	assessment[InterveningGuardrail] = r.Name
	if isResponse {
		assessment[Direction] = "RESPONSE"
	} else {
		assessment[Direction] = "REQUEST"
	}
	assessment[AssessmentReason] = "Violation of regular expression detected."

	if r.ShowAssessment {
		assessment[Assessments] = "Violated regular expression: " + r.Regex
	}
	return assessment
}

// NewRegexGuardrail initializes the RegexGuardrail policy from the given InBuiltPolicy.
func NewRegexGuardrail(inBuiltPolicy dto.InBuiltPolicy) *RegexGuardrail {
	regexGuardrail := &RegexGuardrail{
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
			regexGuardrail.Name = value
		case "regex":
			regexGuardrail.Regex = value
		case "jsonPath":
			regexGuardrail.JSONPath = value
		case "invert":
			regexGuardrail.Inverted = value == "true"
		case "showAssessment":
			regexGuardrail.ShowAssessment = value == "true"
		}
	}
	return regexGuardrail
}
