package mediation

import (
	"bytes"
	"compress/gzip"
	"io"
	"strings"

	"github.com/tidwall/gjson"
	dpv2alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v2alpha1"
	"github.com/wso2/apk/common-go-libs/constants"
	"github.com/wso2/apk/gateway/enforcer/internal/config"
	"github.com/wso2/apk/gateway/enforcer/internal/logging"
	"github.com/wso2/apk/gateway/enforcer/internal/requestconfig"
	"google.golang.org/protobuf/types/known/structpb"
)

// AITokenRateLimit represents the configuration for AI token rate limiting in the API Gateway.
type AITokenRateLimit struct {
	PolicyName          string
	PolicyVersion       string
	PolicyID            string
	Enabled             bool
	PromptTokenPath     string
	CompletionTokenPath string
	TotalTokenPath      string
	logger        *logging.Logger
}

const (
	// AITokenRatelimitPolicyKeyEnabled is the key for enabling/disabling the AI token rate limit policy.
	AITokenRatelimitPolicyKeyEnabled = "Enabled"
	// AITokenRatelimitPolicyKeyPromptTokenPath is the key for specifying the path to the prompt token.
	AITokenRatelimitPolicyKeyPromptTokenPath = "PromptTokenPath"
	// AITokenRatelimitPolicyKeyCompletionTokenPath is the key for specifying the path to the completion token.
	AITokenRatelimitPolicyKeyCompletionTokenPath = "CompletionTokenPath"
	// AITokenRatelimitPolicyKeyTotalTokenPath is the key for specifying the path to the total token count.
	AITokenRatelimitPolicyKeyTotalTokenPath = "TotalTokenPath"
)

// NewAITokenRateLimit creates a new AITokenRateLimit instance with default values.
func NewAITokenRateLimit(mediation *dpv2alpha1.Mediation) *AITokenRateLimit {
	enabled := true
	if val, ok := extractPolicyValue(mediation.Parameters, AITokenRatelimitPolicyKeyEnabled); ok {
		if val == "false" {
			enabled = false
		}
	}
	promptTokenPath := ""
	if val, ok := extractPolicyValue(mediation.Parameters, AITokenRatelimitPolicyKeyPromptTokenPath); ok {
		promptTokenPath = val
	}
	completionTokenPath := ""
	if val, ok := extractPolicyValue(mediation.Parameters, AITokenRatelimitPolicyKeyCompletionTokenPath); ok {
		completionTokenPath = val
	}
	totalTokenPath := ""
	if val, ok := extractPolicyValue(mediation.Parameters, AITokenRatelimitPolicyKeyTotalTokenPath); ok {
		totalTokenPath = val
	}
	return &AITokenRateLimit{
		PolicyName:          MediationAITokenRatelimit,
		PolicyVersion:       mediation.PolicyVersion,
		PolicyID:            mediation.PolicyID,
		Enabled:             enabled,
		PromptTokenPath:     promptTokenPath,
		CompletionTokenPath: completionTokenPath,
		TotalTokenPath:      totalTokenPath,
		logger:              &config.GetConfig().Logger,
	}
}

// Process processes the request configuration for AI token rate limiting.
func (a *AITokenRateLimit) Process(requestConfig *requestconfig.Holder) *Result {
	result := &Result{
		StopFurtherProcessing: false,
	}

	if requestConfig == nil || requestConfig.ResponseHeaders == nil || requestConfig.ResponseHeaders.Headers == nil {
		a.logger.Sugar().Debug("No response headers found in requestConfig, skipping analytics processing")
		return result
	}
	isGzipEncoded := false
	for _, header := range requestConfig.ResponseHeaders.Headers.Headers {
		key := header.GetKey()
		value := string(header.GetRawValue())
		if key == "Content-Encoding" {
			if value == "gzip" {
				isGzipEncoded = true
			}
		}
	}
	var br io.Reader
	var err error
	if isGzipEncoded {
		a.logger.Sugar().Debug("Content-Encoding is gzip")
		br, err = gzip.NewReader(bytes.NewReader(requestConfig.ResponseBody.Body))
		if err != nil {
			return result
		} 
	} else {
		br = bytes.NewReader(requestConfig.ResponseBody.Body)
	}

	bodyBytes, err := io.ReadAll(br)
	if err != nil {
		a.logger.Sugar().Errorf("Failed to read response body: %v", err)
		return result
	}
	bodyString := string(bodyBytes)
	a.logger.Sugar().Debugf("Response body: %s", bodyString)

	results := gjson.GetMany(bodyString, removeDollarPrefix(a.PromptTokenPath), removeDollarPrefix(a.CompletionTokenPath), removeDollarPrefix(a.TotalTokenPath))
	if len(results) < 3 {
		a.logger.Sugar().Errorf("Failed to extract token counts from response body: %v", bodyString)
		return result
	}
	promptTokenCount := results[0].Int()
	completionTokenCount := results[1].Int()
	totalTokenCount := results[2].Int()

	result.Metadata[constants.PromptTokenCountIDMetadataKey] = &structpb.Value{Kind: &structpb.Value_NumberValue{NumberValue: float64(promptTokenCount)}}
	result.Metadata[constants.CompletionTokenCountIDMetadataKey] = &structpb.Value{Kind: &structpb.Value_NumberValue{NumberValue: float64(completionTokenCount)}}
	result.Metadata[constants.TotalTokenCountIDMetadataKey] = &structpb.Value{Kind: &structpb.Value_NumberValue{NumberValue: float64(totalTokenCount)}}
	return result
}

func removeDollarPrefix(s string) string {
	if strings.HasPrefix(s, "$.") {
		return s[2:]
	}
	return s
}