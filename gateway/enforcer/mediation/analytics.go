package mediation

import (
	dpv2alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v2alpha1"
	"github.com/wso2/apk/gateway/enforcer/internal/requestconfig"
)

type Analytics struct {
	PolicyName    string `json:"policyName"`
	PolicyVersion string `json:"policyVersion"`
	PolicyID      string `json:"policyID"`
	Enabled       bool   `json:"enabled"`
}

const (
	// MediationAnalyticsPolicyKeyEnabled is the key for enabling/disabling the Analytics policy.
	MediationAnalyticsPolicyKeyEnabled = "Enabled"
)

// NewAnalytics creates a new Analytics instance with default values.
func NewAnalytics(mediation *dpv2alpha1.Mediation) *Analytics {
	enabled := true
	if val, ok := extractPolicyValue(mediation.Parameters, MediationAnalyticsPolicyKeyEnabled); ok {
		if val == "false" {
			enabled = false
		}
	}
	return &Analytics{
		PolicyName:    "Analytics",
		PolicyVersion: mediation.PolicyVersion,
		PolicyID:      mediation.PolicyID,
		Enabled:       enabled,
	}
}

// Process processes the request configuration for analytics.
func (a *Analytics) Process(requestConfig *requestconfig.Holder) *MediationResult {
	// Implement the logic to process the requestConfig for analytics
	// This is a placeholder implementation
	result := &MediationResult{}

	// Add logic to handle analytics here

	return result
}