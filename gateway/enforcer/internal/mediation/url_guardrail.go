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
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	v32 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/tidwall/gjson"
	dpv2alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v2alpha1"
	"github.com/wso2/apk/gateway/enforcer/internal/config"
	"github.com/wso2/apk/gateway/enforcer/internal/logging"
	"github.com/wso2/apk/gateway/enforcer/internal/requestconfig"
)

// URLGuardrail represents the configuration for URL Guardrail policy in the API Gateway.
type URLGuardrail struct {
	PolicyName     string `json:"policyName"`
	PolicyVersion  string `json:"policyVersion"`
	PolicyID       string `json:"policyID"`
	Name           string `json:"name"`
	JSONPath       string `json:"jsonPath"`
	OnlyDNS        bool   `json:"onlyDNS"`
	Timeout        int    `json:"timeout"`
	ShowAssessment bool   `json:"showAssessment"`
	logger         *logging.Logger
	cfg            *config.Server
}

const (
	// URLGuardrailPolicyKeyName is the key for specifying the name of the guardrail.
	URLGuardrailPolicyKeyName = "name"
	// URLGuardrailPolicyKeyJSONPath is the key for specifying the JSON path to extract content.
	URLGuardrailPolicyKeyJSONPath = "jsonPath"
	// URLGuardrailPolicyKeyOnlyDNS is the key for specifying if only DNS validation should be performed.
	URLGuardrailPolicyKeyOnlyDNS = "onlyDNS"
	// URLGuardrailPolicyKeyTimeout is the key for specifying the timeout for URL validation.
	URLGuardrailPolicyKeyTimeout = "timeout"
	// URLGuardrailPolicyKeyShowAssessment is the key for specifying if assessment should be shown.
	URLGuardrailPolicyKeyShowAssessment = "showAssessment"

	// URLGuardrailAPIMExceptionCode is the error code used when an API-level exception occurs due to URL guardrails.
	URLGuardrailAPIMExceptionCode = "GUARDRAIL_API_EXCEPTION"
	// URLGuardrailConstant is the identifier for URL guardrail constants.
	URLGuardrailConstant = "URL_GUARDRAIL"
	// URLGuardrailErrorCode is the HTTP status code returned when a URL guardrail violation occurs.
	URLGuardrailErrorCode = 400
	// URLErrorCode represents the JSON key for the error code in URL guardrail responses.
	URLErrorCode = "errorCode"
	// URLErrorType represents the JSON key for the error type in URL guardrail responses.
	URLErrorType = "errorType"
	// URLErrorMessage represents the JSON key for the error message in URL guardrail responses.
	URLErrorMessage = "errorMessage"
	// URLAssessmentAction represents the JSON key for the action in URL guardrail assessments.
	URLAssessmentAction = "action"
	// URLInterveningGuardrail represents the JSON key for the intervening guardrail in URL guardrail responses.
	URLInterveningGuardrail = "interveningGuardrail"
	// URLDirection represents the JSON key for the direction in URL guardrail responses.
	URLDirection = "direction"
	// URLAssessmentReason represents the JSON key for the reason in URL guardrail assessments.
	URLAssessmentReason = "reason"
	// URLAssessments represents the JSON key for the list of assessments in URL guardrail responses.
	URLAssessments = "assessments"
)

var (
	// URLRegexCompiled is a compiled regex for extracting URLs
	URLRegexCompiled = regexp.MustCompile(`https?://[^\s]+`)
	// URLTextCleanRegexCompiled is a compiled regex for cleaning text
	URLTextCleanRegexCompiled = regexp.MustCompile(`[^\w\s:/.?=&-]`)
)

// NewURLGuardrail creates a new URLGuardrail instance.
func NewURLGuardrail(mediation *dpv2alpha1.Mediation) *URLGuardrail {
	cfg := config.GetConfig()
	logger := cfg.Logger

	name := "URLGuardrail"
	if val, ok := extractPolicyValue(mediation.Parameters, URLGuardrailPolicyKeyName); ok {
		name = val
	}

	jsonPath := "$.content"
	if val, ok := extractPolicyValue(mediation.Parameters, URLGuardrailPolicyKeyJSONPath); ok {
		jsonPath = val
	}

	onlyDNS := false
	if val, ok := extractPolicyValue(mediation.Parameters, URLGuardrailPolicyKeyOnlyDNS); ok {
		onlyDNS = val == "true"
	}

	timeout := 3000 // default 3 seconds
	if val, ok := extractPolicyValue(mediation.Parameters, URLGuardrailPolicyKeyTimeout); ok {
		if intValue, err := strconv.Atoi(val); err == nil {
			timeout = intValue
		}
	}

	showAssessment := false
	if val, ok := extractPolicyValue(mediation.Parameters, URLGuardrailPolicyKeyShowAssessment); ok {
		showAssessment = val == "true"
	}

	return &URLGuardrail{
		PolicyName:     "URLGuardrail",
		PolicyVersion:  mediation.PolicyVersion,
		PolicyID:       mediation.PolicyID,
		Name:           name,
		JSONPath:       jsonPath,
		OnlyDNS:        onlyDNS,
		Timeout:        timeout,
		ShowAssessment: showAssessment,
		logger:         &logger,
		cfg:            cfg,
	}
}

// Process processes the request configuration for URL Guardrail.
func (u *URLGuardrail) Process(requestConfig *requestconfig.Holder) *Result {
	result := NewResult()

	// Handle request body processing
	if requestConfig.ProcessingPhase == requestconfig.ProcessingPhaseRequestBody {
		u.logger.Sugar().Debugf("Beginning request payload validation for URLGuardrail policy: %s", u.Name)

		if requestConfig.RequestBody == nil || requestConfig.RequestBody.Body == nil {
			u.logger.Sugar().Debug("No request body found, skipping URL validation")
			return result
		}

		validationResult, invalidURLs, err := u.validatePayload(requestConfig.RequestBody.Body, false)
		if !validationResult {
			u.logger.Sugar().Debugf("Request payload validation failed for URLGuardrail policy: %s", u.Name)
			return u.buildErrorResponse(false, invalidURLs, err)
		}
		u.logger.Sugar().Debugf("Request payload validation passed for URLGuardrail policy: %s", u.Name)
		return result
	}

	// Handle response body processing
	if requestConfig.ProcessingPhase == requestconfig.ProcessingPhaseResponseBody {
		u.logger.Sugar().Debugf("Beginning response body validation for URLGuardrail policy: %s", u.Name)

		if requestConfig.ResponseBody == nil || requestConfig.ResponseBody.Body == nil {
			u.logger.Sugar().Debug("No response body found, skipping URL validation")
			return result
		}

		validationResult, invalidURLs, err := u.validatePayload(requestConfig.ResponseBody.Body, true)
		if !validationResult {
			u.logger.Sugar().Debugf("Response body validation failed for URLGuardrail policy: %s", u.Name)
			return u.buildErrorResponse(true, invalidURLs, err)
		}
		u.logger.Sugar().Debugf("Response body validation passed for URLGuardrail policy: %s", u.Name)
		return result
	}

	return result
}

// validatePayload validates the payload against the URLGuardrail policy.
func (u *URLGuardrail) validatePayload(payload []byte, isResponse bool) (bool, []string, error) {
	// Decompress response body if needed
	if isResponse {
		bodyStr, err := u.decompressLLMResp(payload)
		if err == nil {
			payload = []byte(bodyStr)
		}
	}

	// Extract value from JSON using JSONPath
	extractedValue, err := u.extractStringValueFromJsonpath(payload, u.JSONPath)
	if err != nil {
		u.logger.Error(err, "Error extracting value from JSON using JSONPath")
		return false, []string{}, err
	}

	// Clean and trim the extracted text
	extractedValue = URLTextCleanRegexCompiled.ReplaceAllString(extractedValue, "")
	extractedValue = strings.TrimSpace(extractedValue)

	// Extract URLs from the value
	urls := URLRegexCompiled.FindAllString(extractedValue, -1)
	if len(urls) == 0 {
		u.logger.Sugar().Debug("No URLs found in the extracted content")
		return true, []string{}, nil
	}

	invalidURLs := make([]string, 0)
	validationResult := true

	for _, urlStr := range urls {
		var isValid bool
		if u.OnlyDNS {
			isValid = u.checkDNS(urlStr)
			if isValid {
				u.logger.Sugar().Debugf("URL %s is valid via DNS", urlStr)
			} else {
				u.logger.Sugar().Debugf("URL %s is invalid via DNS", urlStr)
			}
		} else {
			isValid = u.checkURL(urlStr)
			if isValid {
				u.logger.Sugar().Debugf("URL %s is reachable via HTTP HEAD request", urlStr)
			} else {
				u.logger.Sugar().Debugf("URL %s is not reachable via HTTP HEAD request", urlStr)
			}
		}

		if !isValid {
			invalidURLs = append(invalidURLs, urlStr)
			validationResult = false
		}
	}

	return validationResult, invalidURLs, nil
}

// checkDNS checks if the URL is resolved via DNS.
func (u *URLGuardrail) checkDNS(target string) bool {
	parsedURL, err := url.Parse(target)
	if err != nil {
		u.logger.Sugar().Errorf("Failed to parse URL: %v", err)
		return false
	}

	host := parsedURL.Hostname()
	if host == "" {
		u.logger.Sugar().Errorf("No hostname found in URL: %s", target)
		return false
	}

	// Create a custom resolver with timeout
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Duration(u.Timeout) * time.Millisecond,
			}
			return d.DialContext(ctx, network, address)
		},
	}

	// Look up IP addresses
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(u.Timeout)*time.Millisecond)
	defer cancel()

	ips, err := resolver.LookupIP(ctx, "ip", host)
	if err != nil {
		u.logger.Sugar().Debugf("DNS lookup failed for %s: %v", host, err)
		return false
	}

	return len(ips) > 0
}

// checkURL checks if the URL is reachable via HTTP HEAD request.
func (u *URLGuardrail) checkURL(target string) bool {
	client := &http.Client{
		Timeout: time.Duration(u.Timeout) * time.Millisecond,
	}

	req, err := http.NewRequest("HEAD", target, nil)
	if err != nil {
		u.logger.Sugar().Errorf("Failed to create HEAD request: %v", err)
		return false
	}
	req.Header.Set("User-Agent", "URLValidator/1.0")

	resp, err := client.Do(req)
	if err != nil {
		u.logger.Sugar().Errorf("HEAD request failed: %v", err)
		return false
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode
	return statusCode >= 200 && statusCode < 400
}

// buildErrorResponse builds the error response for the URLGuardrail policy.
func (u *URLGuardrail) buildErrorResponse(isResponse bool, invalidURLs []string, validationError error) *Result {
	result := NewResult()

	responseBody := make(map[string]interface{})
	responseBody[URLErrorCode] = URLGuardrailAPIMExceptionCode
	responseBody[URLErrorType] = URLGuardrailConstant
	responseBody[URLErrorMessage] = u.buildAssessmentObject(isResponse, invalidURLs, validationError)

	bodyBytes, err := json.Marshal(responseBody)
	if err != nil {
		u.logger.Error(err, "Error marshaling response body to JSON")
		return result
	}

	result.ImmediateResponse = true
	result.ImmediateResponseCode = v32.StatusCode(URLGuardrailErrorCode)
	result.ImmediateResponseBody = string(bodyBytes)
	result.ImmediateResponseContentType = "application/json"
	result.StopFurtherProcessing = true

	return result
}

// buildAssessmentObject builds the assessment object for the URLGuardrail policy.
func (u *URLGuardrail) buildAssessmentObject(isResponse bool, invalidURLs []string, validationError error) map[string]interface{} {
	u.logger.Sugar().Debugf("Building assessment object for URLGuardrail policy: %s", u.Name)
	assessment := make(map[string]interface{})
	assessment[URLAssessmentAction] = "GUARDRAIL_INTERVENED"
	assessment[URLInterveningGuardrail] = u.Name

	if isResponse {
		assessment[URLDirection] = "RESPONSE"
	} else {
		assessment[URLDirection] = "REQUEST"
	}

	// Check if this is a JSONPath extraction error
	if validationError != nil {
		assessment[URLAssessmentReason] = "Error extracting content from payload using JSONPath."
		if u.ShowAssessment {
			assessmentMessage := "JSONPath extraction failed: " + validationError.Error() + ". Please check the JSONPath configuration: " + u.JSONPath
			assessment[URLAssessments] = assessmentMessage
		}
	} else {
		assessment[URLAssessmentReason] = "Violation of URL validity detected."
		if u.ShowAssessment {
			assessmentDetails := make(map[string]interface{})
			assessmentDetails["message"] = "One or more URLs in the payload failed validation."
			assessmentDetails["invalidUrls"] = invalidURLs
			if u.OnlyDNS {
				assessmentDetails["validationType"] = "DNS lookup"
			} else {
				assessmentDetails["validationType"] = "HTTP HEAD request"
			}
			assessment[URLAssessments] = assessmentDetails
		}
	}
	return assessment
}

// decompressLLMResp decompresses the LLM response if it's compressed.
func (u *URLGuardrail) decompressLLMResp(payload []byte) (string, error) {
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
func (u *URLGuardrail) extractStringValueFromJsonpath(payload []byte, jsonPath string) (string, error) {
	bodyString := string(payload)
	result := gjson.Get(bodyString, removeDollarPrefix(jsonPath))

	if !result.Exists() {
		return "", fmt.Errorf("field not found: %s", jsonPath)
	}

	return result.String(), nil
}
