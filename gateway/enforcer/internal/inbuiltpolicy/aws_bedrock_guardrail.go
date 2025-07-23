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
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_service_proc_v3 "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"
	v32 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/wso2/apk/gateway/enforcer/internal/dto"
	"github.com/wso2/apk/gateway/enforcer/internal/logging"
)

// AWSBedrockGuardrail is a struct that represents a AWS Bedrock guardrail policy.
type AWSBedrockGuardrail struct {
	dto.BaseInBuiltPolicy
	Name               string
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	AWSSessionToken    string
	AWSRoleARN         string
	AWSRoleRegion      string
	AWSRoleExternalID  string
	Region             string
	GuardrailID        string
	GuardrailVersion   string
	JSONPath           string
	RedactPII          bool
	PassthroughOnError bool
	ShowAssessment     bool
}

// HandleRequestBody is a method that implements the mediation logic for the AWSBedrockGuardrail policy on request.
func (r *AWSBedrockGuardrail) HandleRequestBody(logger *logging.Logger, req *envoy_service_proc_v3.ProcessingRequest, resp *envoy_service_proc_v3.ProcessingResponse, props map[string]interface{}) *envoy_service_proc_v3.ProcessingResponse {
	logger.Sugar().Debugf("Beginning request payload validation for AWSBedrockGuardrail policy: %s", r.Name)
	result, ok := r.validatePayload(logger, req, false, props)
	if !ok {
		logger.Sugar().Debugf("Request payload validation failed for AWSBedrockGuardrail policy: %s", r.Name)
		return r.buildResponse(logger, result)
	}

	// Check if payload was modified and return the modified content
	if result.ModifiedPayload != nil {
		logger.Sugar().Debugf("Request payload was modified by AWSBedrockGuardrail policy: %s", r.Name)
		r.buildBodyMutationResponse(resp, *result.ModifiedPayload, false)
	}

	logger.Sugar().Debugf("Request payload validation passed for AWSBedrockGuardrail policy: %s", r.Name)
	return nil
}

// HandleResponseBody is a method that implements the mediation logic for the AWSBedrockGuardrail policy on response.
func (r *AWSBedrockGuardrail) HandleResponseBody(logger *logging.Logger, req *envoy_service_proc_v3.ProcessingRequest, resp *envoy_service_proc_v3.ProcessingResponse, props map[string]interface{}) *envoy_service_proc_v3.ProcessingResponse {
	logger.Sugar().Debugf("Beginning response body validation for AWSBedrockGuardrail policy: %s", r.Name)
	result, ok := r.validatePayload(logger, req, true, props)
	if !ok {
		logger.Sugar().Debugf("Response body validation failed for AWSBedrockGuardrail policy: %s", r.Name)
		return r.buildResponse(logger, result)
	}

	// Check if payload was modified and return the modified content
	if result.ModifiedPayload != nil {
		logger.Sugar().Debugf("Response body was modified by AWSBedrockGuardrail policy: %s", r.Name)
		r.buildBodyMutationResponse(resp, *result.ModifiedPayload, true)
	}

	logger.Sugar().Debugf("Response body validation passed for AWSBedrockGuardrail policy: %s", r.Name)
	return nil
}

// validatePayload validates the payload against the AWSBedrockGuardrail policy.
func (r *AWSBedrockGuardrail) validatePayload(logger *logging.Logger, req *envoy_service_proc_v3.ProcessingRequest, isResponse bool, props map[string]interface{}) (AssessmentResult, bool) {
	var result AssessmentResult
	result.IsResponse = isResponse

	var payload []byte
	var compressionType string
	if isResponse {
		var bodyStr string
		var err error
		payload = req.GetResponseBody().Body
		bodyStr, compressionType, err = DecompressLLMResp(payload)
		if err == nil {
			payload = []byte(bodyStr)
		}
	} else {
		payload = req.GetRequestBody().Body
	}

	// Transform response if redactPII is disabled and PIIs identified in request
	if !r.RedactPII && isResponse {
		if maskedPII, exists := props["awsBedrockGuardrailPIIEntities"]; exists {
			if maskedPIIMap, ok := maskedPII.(map[string]string); ok {
				// For response flow, always transform the entire payload (JSONPath is not applicable)
				transformedContent := r.identifyPIIAndTransform(string(payload), maskedPIIMap, logger)
				result.InspectedContent = transformedContent
				modifiedPayload, err := CompressLLMResp([]byte(transformedContent), compressionType)
				if err != nil {
					result.Error = "Error compressing modified payload: " + err.Error()
					logger.Error(err, result.Error)
					return result, false
				}
				result.ModifiedPayload = &modifiedPayload
				return result, true // Continue processing after PII restoration
			}
		}
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

	// Store the inspected content for assessment reporting
	result.InspectedContent = extractedValue

	// Validate content using AWS Bedrock Guardrail
	violation, guardrailOutput, err := r.applyBedrockGuardrailWithDetails(context.Background(), extractedValue, logger)
	if err != nil {
		if r.PassthroughOnError {
			logger.Sugar().Warnf("AWS Bedrock Guardrail validation failed, but PassthroughOnError is enabled: %v", err)
			return result, true
		}
		result.Error = "Error calling AWS Bedrock Guardrail: " + err.Error()
		logger.Sugar().Error(err, result.Error)
		return result, false
	}

	// Store guardrail output for assessment building
	if guardrailOutput != nil {
		result.GuardrailOutput = guardrailOutput
	}

	// Handle guardrail intervention cases
	if guardrailOutput != nil && guardrailOutput.Action == types.GuardrailActionGuardrailIntervened {
		reason := aws.ToString(guardrailOutput.ActionReason)

		// Check if guardrail blocked the request
		if reason == "Guardrail blocked." {
			logger.Sugar().Debug("Guardrail blocked the content")
			return result, false // Violation detected, block the request
		}

		// Check if guardrail masked any PII
		maskApplied := reason == "Guardrail masked."
		if maskApplied {
			logger.Sugar().Debug("Guardrail applied PII masking")

			// Handle PII masking if redactPII is disabled and this is a request
			if !r.RedactPII && !isResponse {
				logger.Sugar().Debug("PII masking applied by Bedrock service. Masking PII in request.")
				maskedContent, maskedPII, err := r.processPIIEntitiesForMasking(guardrailOutput, extractedValue, logger)
				if err != nil {
					logger.Sugar().Errorf("Error processing PII entities: %v", err)
					return result, false
				}
				if len(maskedPII) > 0 {
					dynamicMetadataKeyValuePairs, ok := props["dynamicMetadataMap"].(map[string]interface{})
					if ok {
						dynamicMetadataKeyValuePairs[awsBedrockGuardrailPIIEntitiesKey] = maskedPII
					}
					dynamicMetadataKeyValuePairs[awsBedrockGuardrailPIIEntitiesKey] = maskedPII
					logger.Sugar().Debugf("PII entities masked: %v", maskedPII)
				}
				result.InspectedContent = maskedContent
				// Update the original payload with masked content
				modifiedPayload := r.updatePayloadWithMaskedContent(payload, extractedValue, maskedContent, logger)
				result.ModifiedPayload = &modifiedPayload
				return result, true // Continue processing after masking PII
			}

			// Handle PII redaction if enabled
			if r.RedactPII {
				logger.Sugar().Debug("PII redaction is enabled, processing redacted content")
				redactedContent := r.extractRedactedContent(guardrailOutput, logger)
				if redactedContent != "" {
					result.InspectedContent = redactedContent
					// Update the original payload with redacted content
					modifiedPayload := r.updatePayloadWithMaskedContent(payload, extractedValue, redactedContent, logger)
					result.ModifiedPayload = &modifiedPayload
				}
				return result, true // Continue processing after redacting PII
			}

			return result, true // Continue processing after handling PII
		}

		// Other intervention reasons - block by default
		logger.Sugar().Debugf("Guardrail intervened with reason: %s - blocking content", reason)
		return result, false // Violation detected, block content
	}

	// If violation detected, block the content
	if violation {
		logger.Sugar().Debugf("AWS Bedrock Guardrail detected violation for content: %s", extractedValue)
		return result, false
	}

	logger.Sugar().Debugf("AWS Bedrock Guardrail validation passed for content length: %d", len(extractedValue))
	return result, true
}

// applyBedrockGuardrailWithDetails calls AWS Bedrock Guardrail ApplyGuardrail API and returns detailed output
func (r *AWSBedrockGuardrail) applyBedrockGuardrailWithDetails(ctx context.Context, content string, logger *logging.Logger) (bool, *bedrockruntime.ApplyGuardrailOutput, error) {
	// Load AWS configuration with custom credentials
	cfg, err := r.loadAWSConfig(ctx, logger)
	if err != nil {
		return false, nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create Bedrock Runtime client
	client := bedrockruntime.NewFromConfig(cfg)

	// Prepare ApplyGuardrail input
	input := &bedrockruntime.ApplyGuardrailInput{
		GuardrailIdentifier: aws.String(r.GuardrailID),
		GuardrailVersion:    aws.String(r.GuardrailVersion),
		Source:              types.GuardrailContentSourceInput,
		Content: []types.GuardrailContentBlock{
			&types.GuardrailContentBlockMemberText{
				Value: types.GuardrailTextBlock{
					Text: aws.String(content),
				},
			},
		},
	}

	// Call ApplyGuardrail API
	output, err := client.ApplyGuardrail(ctx, input)
	if err != nil {
		return false, nil, fmt.Errorf("ApplyGuardrail API call failed: %w", err)
	}

	// Evaluate the guardrail response
	violation, err := r.evaluateGuardrailResponse(output, logger)
	return violation, output, err
}

// loadAWSConfig creates AWS configuration with custom credentials and role assumption
// This method supports three authentication modes:
// 1. Role-based authentication: Uses AssumeRole with optional external ID
// 2. Static credentials: Uses provided access key, secret key, and optional session token
// 3. Default credential chain: Falls back to AWS SDK default credential providers
func (r *AWSBedrockGuardrail) loadAWSConfig(ctx context.Context, logger *logging.Logger) (aws.Config, error) {
	var cfg aws.Config
	var err error

	// Check if role-based authentication should be used
	if r.AWSRoleARN != "" && r.AWSRoleRegion != "" {
		logger.Sugar().Debugf("Using role-based authentication with ARN: %s", r.AWSRoleARN)
		cfg, err = r.loadAWSConfigWithAssumeRole(ctx, logger)
	} else if r.AWSAccessKeyID != "" && r.AWSSecretAccessKey != "" {
		logger.Sugar().Debugf("Using direct AWS credentials for authentication")
		cfg, err = r.loadAWSConfigWithStaticCredentials(ctx, logger)
	} else {
		logger.Sugar().Debugf("Using default AWS credential chain")
		cfg, err = config.LoadDefaultConfig(ctx, config.WithRegion(r.Region))
	}

	return cfg, err
}

// loadAWSConfigWithStaticCredentials creates AWS config with static credentials
func (r *AWSBedrockGuardrail) loadAWSConfigWithStaticCredentials(ctx context.Context, logger *logging.Logger) (aws.Config, error) {
	// Create static credentials provider
	var credsProvider aws.CredentialsProvider
	if r.AWSSessionToken != "" {
		// With session token
		credsProvider = credentials.NewStaticCredentialsProvider(
			r.AWSAccessKeyID,
			r.AWSSecretAccessKey,
			r.AWSSessionToken,
		)
	} else {
		// Without session token
		credsProvider = credentials.NewStaticCredentialsProvider(
			r.AWSAccessKeyID,
			r.AWSSecretAccessKey,
			"",
		)
	}

	// Load config with static credentials
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(r.Region),
		config.WithCredentialsProvider(credsProvider),
	)
	if err != nil {
		return aws.Config{}, fmt.Errorf("failed to load AWS config with static credentials: %w", err)
	}

	return cfg, nil
}

// loadAWSConfigWithAssumeRole creates AWS config with role assumption
func (r *AWSBedrockGuardrail) loadAWSConfigWithAssumeRole(ctx context.Context, logger *logging.Logger) (aws.Config, error) {
	// First, create config for the base credentials (to assume the role)
	var baseCfg aws.Config
	var err error

	if r.AWSAccessKeyID != "" && r.AWSSecretAccessKey != "" {
		// Use provided credentials as base
		var baseCredsProvider aws.CredentialsProvider
		if r.AWSSessionToken != "" {
			baseCredsProvider = credentials.NewStaticCredentialsProvider(
				r.AWSAccessKeyID,
				r.AWSSecretAccessKey,
				r.AWSSessionToken,
			)
		} else {
			baseCredsProvider = credentials.NewStaticCredentialsProvider(
				r.AWSAccessKeyID,
				r.AWSSecretAccessKey,
				"",
			)
		}

		baseCfg, err = config.LoadDefaultConfig(ctx,
			config.WithRegion(r.AWSRoleRegion),
			config.WithCredentialsProvider(baseCredsProvider),
		)
	} else {
		// Use default credential chain as base
		baseCfg, err = config.LoadDefaultConfig(ctx, config.WithRegion(r.AWSRoleRegion))
	}

	if err != nil {
		return aws.Config{}, fmt.Errorf("failed to load base AWS config for role assumption: %w", err)
	}

	// Create STS client for role assumption
	stsClient := sts.NewFromConfig(baseCfg)

	// Create assume role credentials provider
	assumeRoleProvider := stscreds.NewAssumeRoleProvider(stsClient, r.AWSRoleARN, func(o *stscreds.AssumeRoleOptions) {
		if r.AWSRoleExternalID != "" {
			o.ExternalID = aws.String(r.AWSRoleExternalID)
		}
		o.RoleSessionName = "bedrock-guardrail-session"
	})

	// Load final config with assumed role credentials for the target region
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(r.Region),
		config.WithCredentialsProvider(assumeRoleProvider),
	)
	if err != nil {
		return aws.Config{}, fmt.Errorf("failed to load AWS config with assume role: %w", err)
	}

	return cfg, nil
}

// evaluateGuardrailResponse processes the AWS Bedrock Guardrail response
func (r *AWSBedrockGuardrail) evaluateGuardrailResponse(output *bedrockruntime.ApplyGuardrailOutput, logger *logging.Logger) (bool, error) {
	if output == nil {
		if r.PassthroughOnError {
			logger.Sugar().Warn("AWS Bedrock Guardrail returned empty response, but PassthroughOnError is enabled")
			return false, nil // No violation, continue processing
		}
		return true, fmt.Errorf("AWS Bedrock Guardrails API returned an invalid response") // Block due to error
	}

	// Check if guardrail intervened
	if output.Action == types.GuardrailActionGuardrailIntervened {
		logger.Sugar().Debugf("AWS Bedrock Guardrail has intervened")

		reason := aws.ToString(output.ActionReason)
		logger.Sugar().Debugf("Guardrail intervention reason: %s", reason)

		// Check if guardrail blocked the request
		if reason == "Guardrail blocked." {
			logger.Sugar().Debugf("Guardrail blocked the content")
			return true, nil // Violation detected, block the request
		}

		// Check if guardrail masked any PII
		maskApplied := reason == "Guardrail masked."
		if maskApplied {
			logger.Sugar().Debugf("Guardrail applied PII masking")

			// Handle PII redaction if enabled
			if r.RedactPII && len(output.Outputs) > 0 {
				logger.Sugar().Debugf("PII redaction is enabled, processing redacted content")
				r.processPIIRedaction(output.Outputs, logger)
			} else if !r.RedactPII {
				logger.Sugar().Debugf("PII masking applied by Bedrock service, but redactPII is disabled")
				// You might want to process PII masking here if needed
			}

			return false, nil // No violation, continue processing after handling PII
		}

		// Other intervention reasons - block by default
		logger.Sugar().Debugf("Guardrail intervened with reason: %s - blocking content", reason)
		return true, nil // Violation detected, block content
	}

	// Check for no intervention
	if output.Action == types.GuardrailActionNone {
		logger.Sugar().Debugf("No guardrail intervention - content is safe")
		return false, nil // No violation, continue processing
	}

	// Unexpected response
	logger.Sugar().Warnf("AWS Bedrock Guardrails returned unexpected action: %s", string(output.Action))
	if r.PassthroughOnError {
		return false, nil // No violation, continue processing
	}
	return true, fmt.Errorf("AWS Bedrock Guardrails returned unexpected response action: %s", string(output.Action)) // Block due to error
}

// processPIIRedaction handles PII redaction from guardrail outputs
func (r *AWSBedrockGuardrail) processPIIRedaction(outputs []types.GuardrailOutputContent, logger *logging.Logger) {
	for _, outputContent := range outputs {
		if outputContent.Text != nil {
			redactedText := aws.ToString(outputContent.Text)
			logger.Sugar().Debugf("Redacted content available: %s", redactedText)
			// Note: This method is now primarily for logging purposes.
			// Actual payload modification is handled in validatePayload method
			// through the ModifiedPayload field in AssessmentResult
		}
	}
}

// processPIIEntitiesForMasking handles PII masking when redactPII is disabled
func (r *AWSBedrockGuardrail) processPIIEntitiesForMasking(output *bedrockruntime.ApplyGuardrailOutput, originalContent string, logger *logging.Logger) (string, map[string]string, error) {
	if output == nil || len(output.Assessments) == 0 {
		return originalContent, nil, nil
	}

	maskedPII := make(map[string]string)
	updatedContent := originalContent
	counter := 0

	for _, assessment := range output.Assessments {
		if assessment.SensitiveInformationPolicy != nil {
			// Process PII entities
			if len(assessment.SensitiveInformationPolicy.PiiEntities) > 0 {
				for _, entity := range assessment.SensitiveInformationPolicy.PiiEntities {
					if entity.Action == types.GuardrailSensitiveInformationPolicyActionAnonymized {
						match := aws.ToString(entity.Match)
						if match != "" && maskedPII[match] == "" {
							entityType := string(entity.Type)
							replacement := fmt.Sprintf("%s_%04x", entityType, counter)
							updatedContent = strings.ReplaceAll(updatedContent, match, replacement)
							maskedPII[match] = replacement
							counter++
							logger.Sugar().Debugf("Masked PII entity: %s -> %s", match, replacement)
						}
					}
				}
			}

			// Process regex matches
			if len(assessment.SensitiveInformationPolicy.Regexes) > 0 {
				for _, regex := range assessment.SensitiveInformationPolicy.Regexes {
					if regex.Action == types.GuardrailSensitiveInformationPolicyActionAnonymized {
						match := aws.ToString(regex.Match)
						name := aws.ToString(regex.Name)
						if match != "" && maskedPII[match] == "" {
							replacement := fmt.Sprintf("%s_%04x", name, counter)
							updatedContent = strings.ReplaceAll(updatedContent, match, replacement)
							maskedPII[match] = replacement
							counter++
							logger.Sugar().Debugf("Masked regex match: %s -> %s", match, replacement)
						}
					}
				}
			}
		}
	}

	return updatedContent, maskedPII, nil
}

// identifyPIIAndTransform handles PII restoration in responses when redactPII is disabled
func (r *AWSBedrockGuardrail) identifyPIIAndTransform(originalContent string, maskedPIIEntities map[string]string, logger *logging.Logger) string {
	if maskedPIIEntities == nil || len(maskedPIIEntities) == 0 {
		logger.Sugar().Debug("No PII entities found in request. No response transformation needed.")
		return originalContent
	}

	transformedContent := originalContent
	foundMasked := false

	for original, placeholder := range maskedPIIEntities {
		if strings.Contains(transformedContent, placeholder) {
			transformedContent = strings.ReplaceAll(transformedContent, placeholder, original)
			foundMasked = true
		}
	}

	if foundMasked {
		logger.Sugar().Debug("PII entities found in request. Replacing masked PIIs back in response.")
	} else {
		logger.Sugar().Debug("No masked PII entities found in response content.")
	}

	return transformedContent
}

// extractRedactedContent extracts redacted content from guardrail outputs
func (r *AWSBedrockGuardrail) extractRedactedContent(output *bedrockruntime.ApplyGuardrailOutput, logger *logging.Logger) string {
	if output == nil || len(output.Outputs) == 0 {
		return ""
	}

	// Get the first output text
	if output.Outputs[0].Text != nil {
		redactedText := aws.ToString(output.Outputs[0].Text)
		logger.Sugar().Debugf("Extracted redacted content of length: %d", len(redactedText))
		return redactedText
	}

	return ""
}

// buildResponse is a method that builds the response body for the AWSBedrockGuardrail policy.
func (r *AWSBedrockGuardrail) buildResponse(logger *logging.Logger, result AssessmentResult) *envoy_service_proc_v3.ProcessingResponse {
	responseBody := make(map[string]interface{})
	responseBody[ErrorCode] = GuardrailAPIMExceptionCode
	responseBody[ErrorType] = AWSBedrockGuardrailConstant
	responseBody[ErrorMessage] = r.buildAssessmentObject(logger, result)

	bodyBytes, err := json.Marshal(responseBody)
	if err != nil {
		logger.Sugar().Error(err, "Error marshaling response body to JSON")
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

// buildAssessmentObject builds a detailed assessment object for the AWSBedrockGuardrail policy.
func (r *AWSBedrockGuardrail) buildAssessmentObject(logger *logging.Logger, result AssessmentResult) map[string]interface{} {
	logger.Sugar().Debugf("Building assessment object for AWSBedrockGuardrail policy: %s", r.Name)
	assessment := make(map[string]interface{})
	assessment[AssessmentAction] = "GUARDRAIL_INTERVENED"
	assessment[InterveningGuardrail] = r.Name
	if result.IsResponse {
		assessment[Direction] = "RESPONSE"
	} else {
		assessment[Direction] = "REQUEST"
	}

	assessment[AssessmentReason] = "Violation of AWS Bedrock Guardrails detected."

	if r.ShowAssessment {
		if result.Error != "" {
			assessment[Assessments] = result.Error
			return assessment
		}

		// Handle AWS Bedrock Guardrail specific assessment data
		if result.GuardrailOutput != nil {
			if guardrailOutput, ok := result.GuardrailOutput.(*bedrockruntime.ApplyGuardrailOutput); ok {
				// Extract first assessment directly from Bedrock response, similar to Java implementation
				if len(guardrailOutput.Assessments) > 0 {
					firstAssessment := r.convertBedrockAssessmentToMap(guardrailOutput.Assessments[0])
					assessment[Assessments] = firstAssessment
				}
			}
		}
	}
	return assessment
}

// convertBedrockAssessmentToMap converts a Bedrock assessment to a map structure,
// similar to the Java implementation that extracts the first assessment
func (r *AWSBedrockGuardrail) convertBedrockAssessmentToMap(assessment types.GuardrailAssessment) map[string]interface{} {
	assessmentMap := make(map[string]interface{})

	// Handle content policy assessment
	if assessment.ContentPolicy != nil {
		contentPolicy := make(map[string]interface{})
		if len(assessment.ContentPolicy.Filters) > 0 {
			filters := make([]map[string]interface{}, 0, len(assessment.ContentPolicy.Filters))
			for _, filter := range assessment.ContentPolicy.Filters {
				filterMap := map[string]interface{}{
					"action":     string(filter.Action),
					"confidence": string(filter.Confidence),
					"type":       string(filter.Type),
				}
				filters = append(filters, filterMap)
			}
			contentPolicy["filters"] = filters
		}
		assessmentMap["contentPolicy"] = contentPolicy
	}

	// Handle topic policy assessment
	if assessment.TopicPolicy != nil {
		topicPolicy := make(map[string]interface{})
		if len(assessment.TopicPolicy.Topics) > 0 {
			topics := make([]map[string]interface{}, 0, len(assessment.TopicPolicy.Topics))
			for _, topic := range assessment.TopicPolicy.Topics {
				topicMap := map[string]interface{}{
					"action": string(topic.Action),
					"name":   aws.ToString(topic.Name),
					"type":   string(topic.Type),
				}
				topics = append(topics, topicMap)
			}
			topicPolicy["topics"] = topics
		}
		assessmentMap["topicPolicy"] = topicPolicy
	}

	// Handle word policy assessment
	if assessment.WordPolicy != nil {
		wordPolicy := make(map[string]interface{})
		if len(assessment.WordPolicy.CustomWords) > 0 {
			customWords := make([]map[string]interface{}, 0, len(assessment.WordPolicy.CustomWords))
			for _, word := range assessment.WordPolicy.CustomWords {
				wordMap := map[string]interface{}{
					"action": string(word.Action),
					"match":  aws.ToString(word.Match),
				}
				customWords = append(customWords, wordMap)
			}
			wordPolicy["customWords"] = customWords
		}
		if len(assessment.WordPolicy.ManagedWordLists) > 0 {
			managedWords := make([]map[string]interface{}, 0, len(assessment.WordPolicy.ManagedWordLists))
			for _, word := range assessment.WordPolicy.ManagedWordLists {
				wordMap := map[string]interface{}{
					"action": string(word.Action),
					"match":  aws.ToString(word.Match),
					"type":   string(word.Type),
				}
				managedWords = append(managedWords, wordMap)
			}
			wordPolicy["managedWordLists"] = managedWords
		}
		assessmentMap["wordPolicy"] = wordPolicy
	}

	// Handle sensitive information policy assessment
	if assessment.SensitiveInformationPolicy != nil {
		sipPolicy := make(map[string]interface{})
		if len(assessment.SensitiveInformationPolicy.PiiEntities) > 0 {
			piiEntities := make([]map[string]interface{}, 0, len(assessment.SensitiveInformationPolicy.PiiEntities))
			for _, entity := range assessment.SensitiveInformationPolicy.PiiEntities {
				entityMap := map[string]interface{}{
					"action": string(entity.Action),
					"match":  aws.ToString(entity.Match),
					"type":   string(entity.Type),
				}
				piiEntities = append(piiEntities, entityMap)
			}
			sipPolicy["piiEntities"] = piiEntities
		}
		if len(assessment.SensitiveInformationPolicy.Regexes) > 0 {
			regexes := make([]map[string]interface{}, 0, len(assessment.SensitiveInformationPolicy.Regexes))
			for _, regex := range assessment.SensitiveInformationPolicy.Regexes {
				regexMap := map[string]interface{}{
					"action": string(regex.Action),
					"match":  aws.ToString(regex.Match),
					"name":   aws.ToString(regex.Name),
				}
				regexes = append(regexes, regexMap)
			}
			sipPolicy["regexes"] = regexes
		}
		assessmentMap["sensitiveInformationPolicy"] = sipPolicy
	}

	return assessmentMap
}

// updatePayloadWithMaskedContent updates the original payload by replacing the extracted content
// with the masked/redacted content, preserving the JSON structure if JSONPath is used (request flow only)
func (r *AWSBedrockGuardrail) updatePayloadWithMaskedContent(originalPayload []byte, extractedValue, modifiedContent string, logger *logging.Logger) []byte {
	if r.JSONPath == "" {
		// If no JSONPath, the entire payload was processed, return the modified content
		logger.Sugar().Debug("No JSONPath specified, replacing entire payload")
		return []byte(modifiedContent)
	}

	// If JSONPath is specified, update only the specific field in the JSON structure (request flow only)
	logger.Sugar().Debugf("Updating JSONPath field '%s' with masked content", r.JSONPath)

	var jsonData map[string]interface{}
	if err := json.Unmarshal(originalPayload, &jsonData); err != nil {
		logger.Sugar().Errorf("Error unmarshaling JSON payload for update: %v", err)
		// Fallback to returning the modified content as-is
		return []byte(modifiedContent)
	}

	// Set the new value at the JSONPath location using the jsonpath utility function
	err := setValueAtJSONPath(jsonData, r.JSONPath, modifiedContent)
	if err != nil {
		logger.Sugar().Errorf("Error setting value at JSONPath '%s': %v", r.JSONPath, err)
		// Fallback to returning the original payload
		return originalPayload
	}

	// Marshal back to JSON to get the full modified payload
	updatedPayload, err := json.Marshal(jsonData)
	if err != nil {
		logger.Sugar().Errorf("Error marshaling updated JSON payload: %v", err)
		// Fallback to returning the original payload
		return originalPayload
	}

	logger.Sugar().Debugf("Successfully updated payload with masked content at JSONPath '%s'", r.JSONPath)
	return updatedPayload
}

// buildBodyMutationResponse creates a response that modifies the request/response body
func (r *AWSBedrockGuardrail) buildBodyMutationResponse(resp *envoy_service_proc_v3.ProcessingResponse, modifiedBody []byte, isResponse bool) {
	// Calculate the new body length
	newBodyLength := len(modifiedBody)

	// Update the Content-Length header
	headers := &envoy_service_proc_v3.HeaderMutation{
		SetHeaders: []*corev3.HeaderValueOption{
			{
				Header: &corev3.HeaderValue{
					Key:      "Content-Length",
					RawValue: []byte(fmt.Sprintf("%d", newBodyLength)),
				},
			},
		},
	}

	bodyResponse := &envoy_service_proc_v3.BodyResponse{
		Response: &envoy_service_proc_v3.CommonResponse{
			Status:         envoy_service_proc_v3.CommonResponse_CONTINUE_AND_REPLACE,
			HeaderMutation: headers,
			BodyMutation: &envoy_service_proc_v3.BodyMutation{
				Mutation: &envoy_service_proc_v3.BodyMutation_Body{
					Body: modifiedBody,
				},
			},
		},
	}

	if isResponse {
		resp.Response = &envoy_service_proc_v3.ProcessingResponse_ResponseBody{
			ResponseBody: bodyResponse,
		}
	} else {
		resp.Response = &envoy_service_proc_v3.ProcessingResponse_RequestBody{
			RequestBody: bodyResponse,
		}
	}
}

// NewAWSBedrockGuardrail initializes the AWSBedrockGuardrail policy from the given InBuiltPolicy.
func NewAWSBedrockGuardrail(inBuiltPolicy dto.InBuiltPolicy) *AWSBedrockGuardrail {
	// Set default values
	awsBedrockGuardrail := &AWSBedrockGuardrail{
		BaseInBuiltPolicy: dto.BaseInBuiltPolicy{
			PolicyName:    inBuiltPolicy.GetPolicyName(),
			PolicyID:      inBuiltPolicy.GetPolicyID(),
			PolicyVersion: inBuiltPolicy.GetPolicyVersion(),
			Parameters:    inBuiltPolicy.GetParameters(),
			PolicyOrder:   inBuiltPolicy.GetPolicyOrder(),
		},
		Name:               AWSBedrockGuardrailName,
		JSONPath:           "",
		RedactPII:          false,
		PassthroughOnError: false,
		ShowAssessment:     false,
	}

	for key, value := range inBuiltPolicy.GetParameters() {
		switch key {
		case "name":
			awsBedrockGuardrail.Name = value
		case "awsAccessKeyID":
			awsBedrockGuardrail.AWSAccessKeyID = value
		case "awsSecretAccessKey":
			awsBedrockGuardrail.AWSSecretAccessKey = value
		case "awsSessionToken":
			awsBedrockGuardrail.AWSSessionToken = value
		case "awsRoleARN":
			awsBedrockGuardrail.AWSRoleARN = value
		case "awsRoleRegion":
			awsBedrockGuardrail.AWSRoleRegion = value
		case "awsRoleExternalID":
			awsBedrockGuardrail.AWSRoleExternalID = value
		case "region":
			awsBedrockGuardrail.Region = value
		case "guardrailID":
			awsBedrockGuardrail.GuardrailID = value
		case "guardrailVersion":
			awsBedrockGuardrail.GuardrailVersion = value
		case "jsonPath":
			awsBedrockGuardrail.JSONPath = value
		case "redactPII":
			awsBedrockGuardrail.RedactPII = value == "true"
		case "passthroughOnError":
			awsBedrockGuardrail.PassthroughOnError = value == "true"
		case "showAssessment":
			awsBedrockGuardrail.ShowAssessment = value == "true"
		}
	}
	return awsBedrockGuardrail
}
