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
	"io"
	"regexp"

	v32 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/tidwall/gjson"
	dpv2alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v2alpha1"
	"github.com/wso2/apk/gateway/enforcer/internal/config"
	"github.com/wso2/apk/gateway/enforcer/internal/logging"
	"github.com/wso2/apk/gateway/enforcer/internal/requestconfig"
)

// RegexGuardrail represents the configuration for Regex Guardrail policy in the API Gateway.
type RegexGuardrail struct {
	PolicyName     string `json:"policyName"`
	PolicyVersion  string `json:"policyVersion"`
	PolicyID       string `json:"policyID"`
	Name           string `json:"name"`
	Regex          string `json:"regex"`
	JSONPath       string `json:"jsonPath"`
	Inverted       bool   `json:"inverted"`
	ShowAssessment bool   `json:"showAssessment"`
	logger         *logging.Logger
	cfg            *config.Server
}

const (
	// RegexGuardrailPolicyKeyName is the key for specifying the name of the guardrail.
	RegexGuardrailPolicyKeyName = "name"
	// RegexGuardrailPolicyKeyRegex is the key for specifying the regex pattern.
	RegexGuardrailPolicyKeyRegex = "regex"
	// RegexGuardrailPolicyKeyJSONPath is the key for specifying the JSON path to extract content.
	RegexGuardrailPolicyKeyJSONPath = "jsonPath"
	// RegexGuardrailPolicyKeyInverted is the key for specifying if the validation should be inverted.
	RegexGuardrailPolicyKeyInverted = "invert"
	// RegexGuardrailPolicyKeyShowAssessment is the key for specifying if assessment should be shown.
	RegexGuardrailPolicyKeyShowAssessment = "showAssessment"
	// RegexGuardrailAPIMExceptionCode is the error code used when an API-level exception occurs due to regex guardrails.
	RegexGuardrailAPIMExceptionCode = "GUARDRAIL_API_EXCEPTION"
	// RegexGuardrailConstant is the identifier for regex guardrail constants.
	RegexGuardrailConstant = "REGEX_GUARDRAIL"
	// RegexGuardrailErrorCode is the HTTP status code returned when a regex guardrail violation occurs.
	RegexGuardrailErrorCode = 400
	// RegexErrorCode represents the JSON key for the error code in regex guardrail responses.
	RegexErrorCode = "errorCode"
	// RegexErrorType represents the JSON key for the error type in regex guardrail responses.
	RegexErrorType = "errorType"
	// RegexErrorMessage represents the JSON key for the error message in regex guardrail responses.
	RegexErrorMessage = "errorMessage"
	// RegexAssessmentAction represents the JSON key for the action in regex guardrail assessments.
	RegexAssessmentAction = "action"
	// RegexInterveningGuardrail represents the JSON key for the intervening guardrail in regex guardrail responses.
	RegexInterveningGuardrail = "interveningGuardrail"
	// RegexDirection represents the JSON key for the direction in regex guardrail responses.
	RegexDirection = "direction"
	// RegexAssessmentReason represents the JSON key for the reason in regex guardrail assessments.
	RegexAssessmentReason = "reason"
	// RegexAssessments represents the JSON key for the list of assessments in regex guardrail responses.
	RegexAssessments = "assessments"
)

// NewRegexGuardrail creates a new RegexGuardrail instance.
func NewRegexGuardrail(mediation *dpv2alpha1.Mediation) *RegexGuardrail {
	cfg := config.GetConfig()
	logger := cfg.Logger

	name := "RegexGuardrail"
	if val, ok := extractPolicyValue(mediation.Parameters, RegexGuardrailPolicyKeyName); ok {
		name = val
	}

	regex := ""
	if val, ok := extractPolicyValue(mediation.Parameters, RegexGuardrailPolicyKeyRegex); ok {
		regex = val
	}

	jsonPath := "$.content"
	if val, ok := extractPolicyValue(mediation.Parameters, RegexGuardrailPolicyKeyJSONPath); ok {
		jsonPath = val
	}

	inverted := false
	if val, ok := extractPolicyValue(mediation.Parameters, RegexGuardrailPolicyKeyInverted); ok {
		inverted = val == "true"
	}

	showAssessment := false
	if val, ok := extractPolicyValue(mediation.Parameters, RegexGuardrailPolicyKeyShowAssessment); ok {
		showAssessment = val == "true"
	}

	return &RegexGuardrail{
		PolicyName:     "RegexGuardrail",
		PolicyVersion:  mediation.PolicyVersion,
		PolicyID:       mediation.PolicyID,
		Name:           name,
		Regex:          regex,
		JSONPath:       jsonPath,
		Inverted:       inverted,
		ShowAssessment: showAssessment,
		logger:         &logger,
		cfg:            cfg,
	}
}

// Process processes the request configuration for Regex Guardrail.
func (r *RegexGuardrail) Process(requestConfig *requestconfig.Holder) *Result {
	result := NewResult()

	// Handle request body processing
	if requestConfig.ProcessingPhase == requestconfig.ProcessingPhaseRequestBody {
		r.logger.Sugar().Debugf("Beginning request payload validation for RegexGuardrail policy: %s", r.Name)

		if requestConfig.RequestBody == nil || requestConfig.RequestBody.Body == nil {
			r.logger.Sugar().Debug("No request body found, skipping regex validation")
			return result
		}

		validationResult, err := r.validatePayload(requestConfig.RequestBody.Body, false)
		if !validationResult {
			r.logger.Sugar().Debugf("Request payload validation failed for RegexGuardrail policy: %s", r.Name)
			return r.buildErrorResponse(false, err)
		}
		r.logger.Sugar().Debugf("Request payload validation passed for RegexGuardrail policy: %s", r.Name)
		return result
	}

	// Handle response body processing
	if requestConfig.ProcessingPhase == requestconfig.ProcessingPhaseResponseBody {
		r.logger.Sugar().Debugf("Beginning response body validation for RegexGuardrail policy: %s", r.Name)

		if requestConfig.ResponseBody == nil || requestConfig.ResponseBody.Body == nil {
			r.logger.Sugar().Debug("No response body found, skipping regex validation")
			return result
		}

		validationResult, err := r.validatePayload(requestConfig.ResponseBody.Body, true)
		if !validationResult {
			r.logger.Sugar().Debugf("Response body validation failed for RegexGuardrail policy: %s", r.Name)
			return r.buildErrorResponse(true, err)
		}
		r.logger.Sugar().Debugf("Response body validation passed for RegexGuardrail policy: %s", r.Name)
		return result
	}

	return result
}

// validatePayload validates the payload against the RegexGuardrail policy.
func (r *RegexGuardrail) validatePayload(payload []byte, isResponse bool) (bool, error) {
	// Decompress response body if needed
	if isResponse {
		bodyStr, err := r.decompressLLMResp(payload)
		if err == nil {
			payload = []byte(bodyStr)
		}
	}

	// Extract value using JSONPath
	extractedValue, err := r.extractStringValueFromJsonpath(payload, r.JSONPath)
	if err != nil {
		r.logger.Error(err, "Error extracting value from JSON using JSONPath")
		return false, err
	}

	// Perform regex matching
	matched, err := regexp.MatchString(r.Regex, extractedValue)
	if err != nil {
		r.logger.Error(err, "Error matching regex")
		return false, err
	}

	// Apply inversion logic
	if matched && r.Inverted {
		r.logger.Sugar().Debugf("Regex matched and inverted condition is true, returning false")
		return false, nil
	} else if !matched && !r.Inverted {
		r.logger.Sugar().Debugf("Regex did not match and inverted condition is false, returning false")
		return false, nil
	}

	r.logger.Sugar().Debugf("Regex matched successfully, returning true")
	return true, nil
}

// decompressLLMResp decompresses the response body if it's gzip compressed.
func (r *RegexGuardrail) decompressLLMResp(body []byte) (string, error) {
	reader := bytes.NewReader(body)
	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		// If it's not gzip compressed, return the original body as string
		return string(body), nil
	}
	defer gzipReader.Close()

	decompressed, err := io.ReadAll(gzipReader)
	if err != nil {
		return "", err
	}
	return string(decompressed), nil
}

// extractStringValueFromJsonpath extracts a string value from JSON using the provided JSONPath.
func (r *RegexGuardrail) extractStringValueFromJsonpath(payload []byte, jsonPath string) (string, error) {
	bodyString := string(payload)
	// Convert JSONPath to gjson compatible path
	gjsonPath := convertJSONPathToGjsonPath(removeDollarPrefix(jsonPath))
	result := gjson.Get(bodyString, gjsonPath)

	if !result.Exists() {
		return "", nil
	}

	return result.String(), nil
}

// convertJSONPathToGjsonPath converts JSONPath array notation to gjson compatible path
func convertJSONPathToGjsonPath(path string) string {
	// Convert [0] style array indexing to .0 style for gjson
	re := regexp.MustCompile(`\[(\d+)\]`)
	return re.ReplaceAllString(path, ".$1")
}

// buildErrorResponse builds an error response for the RegexGuardrail policy.
func (r *RegexGuardrail) buildErrorResponse(isResponse bool, validationError error) *Result {
	result := NewResult()

	responseBody := make(map[string]interface{})
	responseBody[RegexErrorCode] = RegexGuardrailAPIMExceptionCode
	responseBody[RegexErrorType] = RegexGuardrailConstant
	responseBody[RegexErrorMessage] = r.buildAssessmentObject(isResponse, validationError)

	bodyBytes, err := json.Marshal(responseBody)
	if err != nil {
		r.logger.Error(err, "Error marshaling response body to JSON")
		return result
	}

	result.ImmediateResponse = true
	result.ImmediateResponseCode = v32.StatusCode(RegexGuardrailErrorCode)
	result.ImmediateResponseBody = string(bodyBytes)
	result.ImmediateResponseContentType = "application/json"
	result.StopFurtherProcessing = true

	return result
}

// buildAssessmentObject builds the assessment object for the RegexGuardrail policy.
func (r *RegexGuardrail) buildAssessmentObject(isResponse bool, validationError error) map[string]interface{} {
	r.logger.Sugar().Debugf("Building assessment object for RegexGuardrail policy: %s", r.Name)
	assessment := make(map[string]interface{})
	assessment[RegexAssessmentAction] = "GUARDRAIL_INTERVENED"
	assessment[RegexInterveningGuardrail] = r.Name
	if isResponse {
		assessment[RegexDirection] = "RESPONSE"
	} else {
		assessment[RegexDirection] = "REQUEST"
	}

	// Check if this is a JSONPath extraction error
	if validationError != nil {
		assessment[RegexAssessmentReason] = "Error extracting content from payload using JSONPath."
		if r.ShowAssessment {
			assessmentMessage := "JSONPath extraction failed: " + validationError.Error() + ". Please check the JSONPath configuration: " + r.JSONPath
			assessment[RegexAssessments] = assessmentMessage
		}
	} else {
		assessment[RegexAssessmentReason] = "Violation of regular expression detected."
		if r.ShowAssessment {
			assessment[RegexAssessments] = "Violated regular expression: " + r.Regex
		}
	}
	return assessment
}
