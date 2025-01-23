package ratelimit

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	"github.com/wso2/apk/gateway/enforcer/internal/config"
	"github.com/wso2/apk/gateway/enforcer/internal/util"
	subscription_model "github.com/wso2/apk/common-go-libs/pkg/server/model"
)

// AIRatelimitHelper is a helper struct for managing AI rate limiting.
type AIRatelimitHelper struct {
	rlClient *client
}

// TokenCountAndModel is a struct that holds the prompt, completion, and total token counts.
type TokenCountAndModel struct {
	promt      int
	completion int
	total      int
	model      string
}

const (
	// DescriptorKeyForAIPromtTokenCount is the descriptor key for the AI prompt token count.
	DescriptorKeyForAIPromtTokenCount = "airequesttokencount"
	// DescriptorKeyForAICompletionTokenCount is the descriptor key for the AI completion token count.
	DescriptorKeyForAICompletionTokenCount = "airesponsetokencount"
	// DescriptorKeyForAITotalTokenCount is the descriptor key for the AI total token count.
	DescriptorKeyForAITotalTokenCount = "aitotaltokencount"
	// DescriptorKeyForSubscriptionBasedAIRequestTokenCount is the descriptor key for the subscription-based AI request token count.
	DescriptorKeyForSubscriptionBasedAIRequestTokenCount = "airequesttokencountsubs"
	// DescriptorKeyForSubscriptionBasedAIResponseTokenCount is the descriptor key for the subscription-based AI response token count.
	DescriptorKeyForSubscriptionBasedAIResponseTokenCount = "airesponsetokencountsubs"
	// DescriptorKeyForSubscriptionBasedAITotalTokenCount is the descriptor key for the subscription-based AI total token count.
	DescriptorKeyForSubscriptionBasedAITotalTokenCount = "aitotaltokencountsubs"
	// DescriptorKeyForAISubscription is the descriptor key for the AI subscription.
	DescriptorKeyForAISubscription = "subscription"
)

// NewAIRatelimitHelper creates a new instance of the AIRatelimitHelper.
func NewAIRatelimitHelper(cfg *config.Server) *AIRatelimitHelper {
	client := newClient(cfg)
	client.start()
	return &AIRatelimitHelper{
		rlClient: client,
	}
}

// DoAIRatelimit performs AI rate limiting.
func (airl *AIRatelimitHelper) DoAIRatelimit(tokenCount *TokenCountAndModel, doBackendBasedAIRatelimit bool, doSubscriptionBasedAIRatelimit bool, backendBasedAIRatelimitDescriptorValue string, subscription *subscription_model.Subscription, application *subscription_model.Application) {
	go func() {
		configs := []*keyValueHitsAddend{}
		if doBackendBasedAIRatelimit {
			// For promt token count
			configs = append(configs, &keyValueHitsAddend{
				Key:        DescriptorKeyForAIPromtTokenCount,
				Value:      backendBasedAIRatelimitDescriptorValue,
				HitsAddend: tokenCount.promt,
			})
			// For completion token count
			configs = append(configs, &keyValueHitsAddend{
				Key:        DescriptorKeyForAICompletionTokenCount,
				Value:      backendBasedAIRatelimitDescriptorValue,
				HitsAddend: tokenCount.completion,
			})
			// For total token count
			configs = append(configs, &keyValueHitsAddend{
				Key:        DescriptorKeyForAITotalTokenCount,
				Value:      backendBasedAIRatelimitDescriptorValue,
				HitsAddend: tokenCount.total,
			})
		}
		if doSubscriptionBasedAIRatelimit && subscription != nil && application != nil {
			// For promt token count
			configs = append(configs, &keyValueHitsAddend{
				Key:        DescriptorKeyForSubscriptionBasedAIRequestTokenCount,
				Value:      fmt.Sprintf("%s-%s", subscription.Organization, subscription.RatelimitTier),
				KeyValueHitsAddend: &keyValueHitsAddend{
					Key: 	  DescriptorKeyForAISubscription,
					Value: fmt.Sprintf("%s:%s%s", subscription.SubscribedAPI.Name, application.UUID, subscription.UUID),
					HitsAddend: tokenCount.promt,
				},
			})
			// For completion token count
			configs = append(configs, &keyValueHitsAddend{
				Key:        DescriptorKeyForSubscriptionBasedAIResponseTokenCount,
				Value:      fmt.Sprintf("%s-%s", subscription.Organization, subscription.RatelimitTier),
				KeyValueHitsAddend: &keyValueHitsAddend{
					Key: 	  DescriptorKeyForAISubscription,
					Value: fmt.Sprintf("%s:%s%s", subscription.SubscribedAPI.Name, application.UUID, subscription.UUID),
					HitsAddend: tokenCount.completion,
				},
			})
			// For total token count
			configs = append(configs, &keyValueHitsAddend{
				Key:        DescriptorKeyForSubscriptionBasedAITotalTokenCount,
				Value:      fmt.Sprintf("%s-%s", subscription.Organization, subscription.RatelimitTier),
				KeyValueHitsAddend: &keyValueHitsAddend{
					Key: 	  DescriptorKeyForAISubscription,
					Value: fmt.Sprintf("%s:%s%s", subscription.SubscribedAPI.Name, application.UUID, subscription.UUID),
					HitsAddend: tokenCount.total,
				},
			})
		}
		airl.rlClient.shouldRatelimit(configs)
	}()
}

// ExtractTokenCountFromExternalProcessingResponseHeaders extracts token counts from external processing response headers.
func ExtractTokenCountFromExternalProcessingResponseHeaders(headerValues []*v3.HeaderValue, promptHeader, completionHeader, totalHeader, modelHeader string) (*TokenCountAndModel, error) {
	tokenCount := &TokenCountAndModel{}
	promtFlag, completionFlag, totalFlag := false, false, false
	for _, headerValue := range headerValues {
		if headerValue.Key == promptHeader {
			if headerValue.Value != "" {
				value, err := util.ConvertStringToInt(headerValue.Value)
				if err != nil {
					tokenCount.promt = value
					promtFlag = true
				} else {
					return nil, err
				}
			} else if len(headerValue.RawValue) != 0 {
				value, err := util.ConvertBytesToInt(headerValue.RawValue)
				if err != nil {
					tokenCount.promt = value
					promtFlag = true
				} else {
					return nil, err
				}
			}

		} else if headerValue.Key == completionHeader {
			if headerValue.Value != "" {
				value, err := strconv.Atoi(headerValue.Value)
				if err != nil {
					tokenCount.completion = value
					completionFlag = true
				} else {
					return nil, err
				}
			} else if len(headerValue.RawValue) != 0 {
				value, err := util.ConvertBytesToInt(headerValue.RawValue)
				if err != nil {
					tokenCount.completion = value
					completionFlag = true
				} else {
					return nil, err
				}
			}
		} else if headerValue.Key == totalHeader {
			if headerValue.Value != "" {
				value, err := strconv.Atoi(headerValue.Value)
				if err != nil {
					tokenCount.total = value
					totalFlag = true
				} else {
					return nil, err
				}
			} else if len(headerValue.RawValue) != 0 {
				value, err := util.ConvertBytesToInt(headerValue.RawValue)
				if err != nil {
					tokenCount.total = value
					totalFlag = true
				} else {
					return nil, err
				}
			}
		} else if headerValue.Key == modelHeader {
			if headerValue.Value != "" {
				tokenCount.model = headerValue.Value
			} else if len(headerValue.RawValue) != 0 {
				tokenCount.model = string(headerValue.RawValue)
			}
		}
	}
	if !promtFlag || !completionFlag || !totalFlag {
		return nil, fmt.Errorf("missing token headers from the AI response headers")
	}
	return tokenCount, nil
}

// ExtractTokenCountFromExternalProcessingResponseBody extracts token counts from external processing response body.
func ExtractTokenCountFromExternalProcessingResponseBody(body []byte, promptPath, completionPath, totalPath, modelPath string) (*TokenCountAndModel, error) {
	bodyStr := string(body)
	sanitizedBody := sanitize(bodyStr)
	tokenCount, err := extractUsageFromBody(sanitizedBody, promptPath, completionPath, totalPath, "model")
	if err != nil {
		return nil, fmt.Errorf("failed to extract token count from the AI response body: %w", err)
	}
	return tokenCount, nil

}

func sanitize(input string) string {
	// Define a regex to match all newline characters and tabs
	re := regexp.MustCompile(`[\t\n\r]+`)
	// Replace matches with a space and trim the result
	return strings.TrimSpace(re.ReplaceAllString(input, " "))
}

// extractValueFromPath extracts a value from a nested JSON structure based on a dot-separated path.
func extractValueFromPath(data map[string]interface{}, path string) (interface{}, error) {
	keys := strings.Split(path, ".")
	if len(keys) > 0 && keys[0] == "$" {
		keys = keys[1:]
	}

	var current interface{} = data
	for _, key := range keys {
		if node, ok := current.(map[string]interface{}); ok {
			if val, exists := node[key]; exists {
				current = val
			} else {
				return nil, errors.New("key not found: " + key)
			}
		} else {
			return nil, errors.New("invalid structure for key: " + key)
		}
	}
	return current, nil
}

// extractUsageFromBody extracts usage data from the JSON body based on the provided paths.
func extractUsageFromBody(body, completionTokenPath, promptTokenPath, totalTokenPath, modelPath string) (*TokenCountAndModel, error) {
	body = sanitize(body)
	var rootNode map[string]interface{}
	if err := json.Unmarshal([]byte(body), &rootNode); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	usage := &TokenCountAndModel{}

	// Extract prompt tokens
	promt, err := extractValueFromPath(rootNode, promptTokenPath)
	if err != nil {
		return nil, fmt.Errorf("failed to extract prompt tokens: %w", err)
	}
	if pt, ok := promt.(float64); ok { // JSON numbers are decoded as float64
		usage.promt = int(pt)
	} else {
		return nil, errors.New("invalid type for prompt tokens")
	}

	// Extract completion tokens
	completion, err := extractValueFromPath(rootNode, completionTokenPath)
	if err != nil {
		return nil, fmt.Errorf("failed to extract completion tokens: %w", err)
	}
	if ct, ok := completion.(float64); ok {
		usage.completion = int(ct)
	} else {
		return nil, errors.New("invalid type for completion tokens")
	}

	// Extract total tokens
	total, err := extractValueFromPath(rootNode, totalTokenPath)
	if err != nil {
		return nil, fmt.Errorf("failed to extract total tokens: %w", err)
	}
	if tt, ok := total.(float64); ok {
		usage.total = int(tt)
	} else {
		return nil, errors.New("invalid type for total tokens")
	}

	// Extract model
	// model, err := extractValueFromPath(rootNode, modelPath)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to extract model: %w", err)
	// }
	// if m, ok := model.(string); ok {
	// 	usage.model = m
	// } else {
	// 	return nil, errors.New("invalid type for model")
	// }

	return usage, nil
}
