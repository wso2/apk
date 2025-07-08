package mediation

import (
	dpv2alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v2alpha1"
	"github.com/wso2/apk/gateway/enforcer/internal/requestconfig"
)

// SubscriptionRatelimit represents the configuration for subscription rate limiting in the API Gateway.
type SubscriptionRatelimit struct {
	PolicyName    string
	PolicyVersion string
	PolicyID      string
	Enabled       bool
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
	return &SubscriptionRatelimit{
		PolicyName:    "SubscriptionValidation",
		PolicyVersion: "v1",
		PolicyID:      "subscription-validation",
		Enabled:       enabled,
	}
}

// Process processes the request configuration for Subscription Rate Limit.
func (s *SubscriptionRatelimit) Process(requestConfig *requestconfig.Holder) *MediationResult {
	// Implement the logic to process the requestConfig for Subscription Rate Limit
	// This is a placeholder implementation
	result := &MediationResult{}

	// Add logic to handle subscription rate limiting here

	return result
}
