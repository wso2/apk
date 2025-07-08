package mediation

import (
	dpv2alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v2alpha1"
	"github.com/wso2/apk/gateway/enforcer/internal/requestconfig"
)

// SubscriptionValidation represents the configuration for subscription validation in the API Gateway.
type SubscriptionValidation struct {
	PolicyName    string
	PolicyVersion string
	PolicyID      string
	Enabled       bool
}

const (
	// SubscriptionValidationPolicyKeyEnabled is the key for enabling/disabling the Subscription Validation policy.
	SubscriptionValidationPolicyKeyEnabled = "Enabled"
)

// NewSubscriptionValidation creates a new SubscriptionValidation instance with default values.
func NewSubscriptionValidation(mediation *dpv2alpha1.Mediation) *SubscriptionValidation {
	enabled := true
	if val, ok := extractPolicyValue(mediation.Parameters, SubscriptionValidationPolicyKeyEnabled); ok {
		if val == "false" {
			enabled = false
		}
	}
	return &SubscriptionValidation{
		PolicyName:    "SubscriptionValidation",
		PolicyVersion: "v1",
		PolicyID:      "subscription-validation",
		Enabled:       enabled,
	}
}

// Process processes the request configuration for Subscription Validation.
func (s *SubscriptionValidation) Process(requestConfig *requestconfig.Holder) *Result {
	// Implement the logic to process the requestConfig for Subscription Validation
	// This is a placeholder implementation
	result := &Result{}

	// Add logic to handle subscription validation here

	return result
}
