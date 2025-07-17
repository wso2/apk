package mediation

import (
	dpv2alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v2alpha1"
	"github.com/wso2/apk/gateway/enforcer/internal/config"
	"github.com/wso2/apk/gateway/enforcer/internal/logging"
	"github.com/wso2/apk/gateway/enforcer/internal/requestconfig"
	"github.com/wso2/apk/common-go-libs/constants"
)

// SubscriptionRatelimit represents the configuration for subscription rate limiting in the API Gateway.
type SubscriptionRatelimit struct {
	PolicyName    string
	PolicyVersion string
	PolicyID      string
	Enabled       bool
	logger 	  *logging.Logger
}

const (
	// SubscriptionRatelimitPolicyKeyEnabled is the key for enabling/disabling the Subscription Rate Limit policy.
	SubscriptionRatelimitPolicyKeyEnabled = "Enabled"
)

// NewSubscriptionRatelimit creates a new SubscriptionValidation instance with default values.
func NewSubscriptionRatelimit(meidation *dpv2alpha1.Mediation) *SubscriptionRatelimit {
	enabled := true
	if val, ok := extractPolicyValue(meidation.Parameters, SubscriptionRatelimitPolicyKeyEnabled); ok {
		if val == "false" {
			enabled = false
		}
	}
	cfg := config.GetConfig()
	logger := cfg.Logger
	return &SubscriptionRatelimit{
		PolicyName:    "SubscriptionValidation",
		PolicyVersion: "v1",
		PolicyID:      "subscription-validation",
		Enabled:       enabled,
		logger:        &logger,
	}
}

// Process processes the request configuration for Subscription Rate Limit.
func (s *SubscriptionRatelimit) Process(requestConfig *requestconfig.Holder) *Result {
	result := &Result{}

	if requestConfig.MatchedSubscription != nil {
		// Add subscription rate limit headers to the requestConfig
		result.AddHeaders[constants.SubscriptionUUIDHeaderName] = requestConfig.MatchedSubscription.UUID
	} else {
		s.logger.Sugar().Errorf("No subscription found for the request. Hence not adding any subscription rate limit headers.")
	}

	return result
}
