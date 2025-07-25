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
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_service_proc_v3 "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"
	v32 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/wso2/apk/gateway/enforcer/internal/dto"
	"github.com/wso2/apk/gateway/enforcer/internal/logging"
)

// URLGuardrail is a struct that represents a URL guardrail policy.
type URLGuardrail struct {
	dto.BaseInBuiltPolicy
	Name           string
	JSONPath       string
	OnlyDNS        bool
	Timeout        int
	ShowAssessment bool
}

// HandleRequestBody is a method that implements the mediation logic for the URLGuardrail policy on request.
func (r *URLGuardrail) HandleRequestBody(logger *logging.Logger, req *envoy_service_proc_v3.ProcessingRequest, resp *envoy_service_proc_v3.ProcessingResponse, props map[string]interface{}) *envoy_service_proc_v3.ProcessingResponse {
	logger.Sugar().Debugf("Beginning request payload validation for URLGuardrail policy: %s", r.Name)
	validationResult, invalidURLs, err := r.validatePayload(logger, req.GetRequestBody().Body, false)
	if !validationResult {
		logger.Sugar().Debugf("Request payload validation failed for URLGuardrail policy: %s", r.Name)
		return r.buildResponse(logger, false, invalidURLs, err)
	}
	logger.Sugar().Debugf("Request payload validation passed for URLGuardrail policy: %s", r.Name)
	return nil
}

// HandleResponseBody is a method that implements the mediation logic for the URLGuardrail policy on response.
func (r *URLGuardrail) HandleResponseBody(logger *logging.Logger, req *envoy_service_proc_v3.ProcessingRequest, resp *envoy_service_proc_v3.ProcessingResponse, props map[string]interface{}) *envoy_service_proc_v3.ProcessingResponse {
	logger.Sugar().Debugf("Beginning response body validation for URLGuardrail policy: %s", r.Name)
	validationResult, invalidURLs, err := r.validatePayload(logger, req.GetResponseBody().Body, true)
	if !validationResult {
		logger.Sugar().Debugf("Response body validation failed for URLGuardrail policy: %s", r.Name)
		return r.buildResponse(logger, true, invalidURLs, err)
	}
	logger.Sugar().Debugf("Response body validation passed for URLGuardrail policy: %s", r.Name)
	return nil
}

// validatePayload validates the payload against the URLGuardrail policy.
func (r *URLGuardrail) validatePayload(logger *logging.Logger, payload []byte, isResponse bool) (bool, []string, error) {
	if isResponse {
		bodyStr, _, err := DecompressLLMResp(payload)
		if err == nil {
			payload = []byte(bodyStr)
		}
	}
	extractedValue, err := ExtractStringValueFromJsonpath(logger, payload, r.JSONPath)
	if err != nil {
		logger.Error(err, "Error extracting value from JSON using JSONPath")
		return false, []string{}, err
	}

	// Clean and trim
	extractedValue = TextCleanRegexCompiled.ReplaceAllString(extractedValue, "")
	extractedValue = strings.TrimSpace(extractedValue)

	// Extract URLs from the value
	urls := URLRegexCompiled.FindAllString(extractedValue, -1)
	invalidURLs := make([]string, 0)
	validationResult := true
	for _, url := range urls {
		if r.OnlyDNS {
			if r.checkDNS(logger, url) {
				logger.Sugar().Debugf("URL %s is valid via DNS", url)
			} else {
				logger.Sugar().Debugf("URL %s is invalid via DNS", url)
				invalidURLs = append(invalidURLs, url)
				validationResult = false
			}
		} else {
			if r.checkURL(logger, url) {
				logger.Sugar().Debugf("URL %s is reachable via HTTP HEAD request", url)
			} else {
				logger.Sugar().Debugf("URL %s is not reachable via HTTP HEAD request", url)
				invalidURLs = append(invalidURLs, url)
				validationResult = false
			}
		}
	}
	return validationResult, invalidURLs, nil
}

// checkDNS checks if the URL is resolved via DNS using DNS-over-HTTPS.
func (r *URLGuardrail) checkDNS(logger *logging.Logger, target string) bool {
	parsedURL, err := url.Parse(target)
	if err != nil {
		logger.Sugar().Errorf("Failed to parse URL: %v", err)
		return false
	}

	host := parsedURL.Hostname()

	// Create a custom resolver with timeout
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Duration(r.Timeout) * time.Millisecond,
			}
			return d.DialContext(ctx, network, address)
		},
	}

	// Look up IP addresses
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.Timeout)*time.Millisecond)
	defer cancel()

	ips, err := resolver.LookupIP(ctx, "ip", host)
	if err != nil {
		logger.Sugar().Debugf("DNS lookup failed for %s: %v", host, err)
		return false
	}

	return len(ips) > 0
}

// checkURL checks if the URL is reachable via HTTP HEAD request.
func (r *URLGuardrail) checkURL(logger *logging.Logger, target string) bool {
	client := &http.Client{
		Timeout: time.Duration(r.Timeout) * time.Millisecond,
	}

	req, err := http.NewRequest("HEAD", target, nil)
	if err != nil {
		logger.Sugar().Errorf("Failed to create HEAD request: %v", err)
		return false
	}
	req.Header.Set("User-Agent", "URLValidator/1.0")

	resp, err := client.Do(req)
	if err != nil {
		logger.Sugar().Errorf("HEAD request failed: %v", err)
		return false
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode
	return statusCode >= 200 && statusCode < 400
}

// buildResponse is a method that builds the response body for the URLGuardrail policy.
func (r *URLGuardrail) buildResponse(logger *logging.Logger, isResponse bool, invalidURLs []string, validationError error) *envoy_service_proc_v3.ProcessingResponse {
	responseBody := make(map[string]interface{})
	responseBody[ErrorCode] = GuardrailAPIMExceptionCode
	responseBody[ErrorType] = URLGuardrailConstant
	responseBody[ErrorMessage] = r.buildAssessmentObject(logger, isResponse, invalidURLs, validationError)

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

// buildAssessmentObject is a method that builds the assessment object for the URLGuardrail policy.
func (r *URLGuardrail) buildAssessmentObject(logger *logging.Logger, isResponse bool, invalidURLs []string, validationError error) map[string]interface{} {
	logger.Sugar().Debugf("Building assessment object for URLGuardrail policy: %s", r.Name)
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
		assessment[AssessmentReason] = "Violation of url validity detected."
		if r.ShowAssessment {
			assessmentDetails := make(map[string]interface{})
			assessmentDetails["message"] = "One or more URLs in the payload failed validation."
			assessmentDetails["invalidUrls"] = invalidURLs
			assessment[Assessments] = assessmentDetails
		}
	}
	return assessment
}

// NewURLGuardrail initializes the URLGuardrail policy from the given InBuiltPolicy.
func NewURLGuardrail(inBuiltPolicy dto.InBuiltPolicy) *URLGuardrail {
	URLGuardrail := &URLGuardrail{
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
			URLGuardrail.Name = value
		case "jsonPath":
			URLGuardrail.JSONPath = value
		case "onlyDNS":
			URLGuardrail.OnlyDNS = value == "true"
		case "timeout":
			if timeout, err := strconv.Atoi(value); err == nil {
				URLGuardrail.Timeout = timeout
			} else {
				URLGuardrail.Timeout = 3000
			}
		case "showAssessment":
			URLGuardrail.ShowAssessment = value == "true"
		}
	}
	return URLGuardrail
}
