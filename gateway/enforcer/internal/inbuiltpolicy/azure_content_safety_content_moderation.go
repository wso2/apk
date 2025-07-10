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

	"github.com/wso2/apk/gateway/enforcer/internal/util"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_service_proc_v3 "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"
	v32 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/wso2/apk/gateway/enforcer/internal/dto"
	"github.com/wso2/apk/gateway/enforcer/internal/logging"
)

const (
	azureContentSafetyContentModerationEndpoint = "/contentsafety/text:analyze?api-version=2024-09-01"
)

// AzureContentSafetyContentModeration is a struct that represents a Azure Content Safety Content Moderation policy.
type AzureContentSafetyContentModeration struct {
	dto.BaseInBuiltPolicy
	Name                      string
	AzureContentSafetyEnpoint string
	AzureContentSafetyKey     string
	HateCategory              int
	SexualCategory            int
	SelfHarmCategory          int
	ViolenceCategory          int
	JSONPath                  string
	PassthroughOnError        bool
	ShowAssessment            bool
}

// HandleRequestBody is a method that implements the mediation logic for the AzureContentSafetyContentModeration policy on request.
func (r *AzureContentSafetyContentModeration) HandleRequestBody(logger *logging.Logger, req *envoy_service_proc_v3.ProcessingRequest) *envoy_service_proc_v3.ProcessingResponse {
	logger.Sugar().Debugf("Beginning request payload validation for AzureContentSafetyContentModeration policy: %s", r.Name)
	result, ok := r.validatePayload(logger, req.GetRequestBody().Body, false)
	if !ok {
		logger.Sugar().Debugf("Request payload validation failed for AzureContentSafetyContentModeration policy: %s", r.Name)
		return r.buildResponse(logger, result)
	}
	logger.Sugar().Debugf("Request payload validation passed for AzureContentSafetyContentModeration policy: %s", r.Name)
	return nil
}

// HandleResponseBody is a method that implements the mediation logic for the AzureContentSafetyContentModeration policy on response.
func (r *AzureContentSafetyContentModeration) HandleResponseBody(logger *logging.Logger, req *envoy_service_proc_v3.ProcessingRequest) *envoy_service_proc_v3.ProcessingResponse {
	logger.Sugar().Debugf("Beginning response body validation for AzureContentSafetyContentModeration policy: %s", r.Name)
	result, ok := r.validatePayload(logger, req.GetResponseBody().Body, true)
	if !ok {
		logger.Sugar().Debugf("Response body validation failed for AzureContentSafetyContentModeration policy: %s", r.Name)
		return r.buildResponse(logger, result)
	}
	logger.Sugar().Debugf("Response body validation passed for AzureContentSafetyContentModeration policy: %s", r.Name)
	return nil
}

// validatePayload validates the payload against the AzureContentSafetyContentModeration policy.
func (r *AzureContentSafetyContentModeration) validatePayload(logger *logging.Logger, payload []byte, isResponse bool) (AssessmentResult, bool) {
	var result AssessmentResult
	result.IsResponse = isResponse
	result.CategoryMap = map[string]int{
		"Hate":     r.HateCategory,
		"Sexual":   r.SexualCategory,
		"SelfHarm": r.SelfHarmCategory,
		"Violence": r.ViolenceCategory,
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
	result.InspectedContent = extractedValue

	categories := []string{}
	for name, val := range result.CategoryMap {
		if val >= 0 && val <= 7 {
			categories = append(categories, name)
		} else {
			logger.Sugar().Debugf("Invalid %s Category: %d. It should be between 0 and 7.", name, val)
		}
	}
	if len(categories) == 0 {
		logger.Sugar().Debug("No valid categories provided for Azure Content Safety API.")
		return result, true
	}

	requestBody := map[string]interface{}{
		"text":               extractedValue,
		"categories":         categories,
		"haltOnBlocklistHit": true,
		"outputType":         "EightSeverityLevels",
	}

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		result.Error = "Failed to marshal request body for Azure Content Safety API: " + err.Error()
		logger.Error(err, result.Error)
		return result, false
	}

	headers := map[string]string{
		"Content-Type":              "application/json",
		"Ocp-Apim-Subscription-Key": r.AzureContentSafetyKey,
	}

	serviceURL := r.AzureContentSafetyEnpoint + azureContentSafetyContentModerationEndpoint
	resp, err := util.MakeHTTPRequestWithRetry("POST", serviceURL, nil, headers, bodyBytes, 30000, 5, 1000)
	if err != nil {
		result.Error = "Failed to call Azure Content Safety API: " + err.Error()
		logger.Error(err, result.Error)
		if r.PassthroughOnError {
			return result, true
		}
		return result, false
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		result.Error = "Azure Content Safety API returned non-200 status code: " + strconv.Itoa(resp.StatusCode)
		logger.Sugar().Debugf(result.Error)
		if r.PassthroughOnError {
			return result, true
		}
		return result, false
	}

	responseBody := make(map[string]interface{})
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		result.Error = "Failed to decode response body from Azure Content Safety API: " + err.Error()
		logger.Error(err, result.Error)
		if r.PassthroughOnError {
			return result, true
		}
		return result, false
	}

	categoriesAnalysis, ok := responseBody["categoriesAnalysis"].([]interface{})
	if !ok {
		result.Error = "categoriesAnalysis missing or invalid in Azure Content Safety API response"
		logger.Sugar().Debugf(result.Error)
		if r.PassthroughOnError {
			return result, true
		}
		return result, false
	}

	// Convert []interface{} to []map[string]interface{} for easier handling
	var assessments []map[string]interface{}
	for _, item := range categoriesAnalysis {
		if analysis, ok := item.(map[string]interface{}); ok {
			assessments = append(assessments, analysis)
		}
	}
	result.CategoriesAnalysis = assessments

	// Check for violations
	for _, analysis := range assessments {
		category, _ := analysis["category"].(string)
		severityFloat, _ := analysis["severity"].(float64)
		severity := int(severityFloat)
		threshold := result.CategoryMap[category]
		if threshold >= 0 && severity >= threshold {
			return result, false
		}
	}
	return result, true
}

// buildResponse is a method that builds the response body for the AzureContentSafetyContentModeration policy.
func (r *AzureContentSafetyContentModeration) buildResponse(logger *logging.Logger, result AssessmentResult) *envoy_service_proc_v3.ProcessingResponse {
	responseBody := make(map[string]interface{})
	responseBody[ErrorCode] = GuardrailAPIMExceptionCode
	responseBody[ErrorType] = AzureContentSafetyContentModerationConstant
	responseBody[ErrorMessage] = r.buildAssessmentObject(logger, result)

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

// buildAssessmentObject builds a detailed assessment object for the AzureContentSafetyContentModeration policy.
func (r *AzureContentSafetyContentModeration) buildAssessmentObject(logger *logging.Logger, result AssessmentResult) map[string]interface{} {
	logger.Sugar().Debugf("Building assessment object for AzureContentSafetyContentModeration policy: %s", r.Name)
	assessment := make(map[string]interface{})
	assessment[AssessmentAction] = "GUARDRAIL_INTERVENED"
	assessment[InterveningGuardrail] = r.Name
	if result.IsResponse {
		assessment[Direction] = "RESPONSE"
	} else {
		assessment[Direction] = "REQUEST"
	}

	assessment[AssessmentReason] = "Violation of Azure content safety content moderation detected."

	if r.ShowAssessment {
		if result.Error != "" {
			assessment[Assessments] = result.Error
			return assessment
		}
		if len(result.CategoriesAnalysis) > 0 && len(result.CategoryMap) > 0 {
			assessmentsWrapper := make(map[string]interface{})
			assessmentsWrapper["inspectedContent"] = result.InspectedContent
			var assessmentsArray []map[string]interface{}
			for _, analysis := range result.CategoriesAnalysis {
				category, _ := analysis["category"].(string)
				severityFloat, _ := analysis["severity"].(float64)
				severity := int(severityFloat)
				threshold := result.CategoryMap[category]
				categoryAssessment := map[string]interface{}{
					"category":  category,
					"severity":  severity,
					"threshold": threshold,
					"result": func() string {
						if threshold >= 0 && severity >= threshold {
							return "FAIL"
						}
						return "PASS"
					}(),
				}
				assessmentsArray = append(assessmentsArray, categoryAssessment)
			}
			assessmentsWrapper["categories"] = assessmentsArray
			assessment[Assessments] = assessmentsWrapper
		} else {
			assessment[Assessments] = result.CategoriesAnalysis
		}
	}
	return assessment
}

// NewAzureContentSafetyContentModeration initializes the AzureContentSafetyContentModeration policy from the given InBuiltPolicy.
func NewAzureContentSafetyContentModeration(inBuiltPolicy dto.InBuiltPolicy) *AzureContentSafetyContentModeration {
	// Set default values
	azureContentSafetyContentModeration := &AzureContentSafetyContentModeration{
		BaseInBuiltPolicy: dto.BaseInBuiltPolicy{
			PolicyName:    inBuiltPolicy.GetPolicyName(),
			PolicyID:      inBuiltPolicy.GetPolicyID(),
			PolicyVersion: inBuiltPolicy.GetPolicyVersion(),
			Parameters:    inBuiltPolicy.GetParameters(),
			PolicyOrder:   inBuiltPolicy.GetPolicyOrder(),
		},
		Name:                      AzureContentSafetyContentModerationName,
		AzureContentSafetyEnpoint: "",
		AzureContentSafetyKey:     "",
		HateCategory:              -1,
		SexualCategory:            -1,
		SelfHarmCategory:          -1,
		ViolenceCategory:          -1,
		JSONPath:                  "",
		PassthroughOnError:        false,
		ShowAssessment:            false,
	}

	for key, value := range inBuiltPolicy.GetParameters() {
		switch key {
		case "name":
			azureContentSafetyContentModeration.Name = value
		case "azureContentSafetyEndpoint":
			if strings.HasSuffix(value, "/") {
				value = strings.TrimSuffix(value, "/")
			}
			azureContentSafetyContentModeration.AzureContentSafetyEnpoint = value
		case "azureContentSafetyKey":
			azureContentSafetyContentModeration.AzureContentSafetyKey = value
		case "hateCategory":
			if val, err := strconv.Atoi(value); err == nil {
				azureContentSafetyContentModeration.HateCategory = val
			}
		case "sexualCategory":
			if val, err := strconv.Atoi(value); err == nil {
				azureContentSafetyContentModeration.SexualCategory = val
			}
		case "selfHarmCategory":
			if val, err := strconv.Atoi(value); err == nil {
				azureContentSafetyContentModeration.SelfHarmCategory = val
			}
		case "violenceCategory":
			if val, err := strconv.Atoi(value); err == nil {
				azureContentSafetyContentModeration.ViolenceCategory = val
			}
		case "jsonPath":
			azureContentSafetyContentModeration.JSONPath = value
		case "passthroughOnError":
			azureContentSafetyContentModeration.PassthroughOnError = value == "true"
		case "showAssessment":
			azureContentSafetyContentModeration.ShowAssessment = value == "true"
		}
	}
	return azureContentSafetyContentModeration
}
