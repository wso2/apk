package mediation

import (
	dpv2alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v2alpha1"
	"github.com/wso2/apk/gateway/enforcer/internal/config"
	"github.com/wso2/apk/gateway/enforcer/internal/logging"
	"github.com/wso2/apk/gateway/enforcer/internal/requestconfig"
)

// Analytics represents the configuration for Analytics policy in the API Gateway.
type Analytics struct {
	PolicyName    string `json:"policyName"`
	PolicyVersion string `json:"policyVersion"`
	PolicyID      string `json:"policyID"`
	Enabled       bool   `json:"enabled"`
	logger        *logging.Logger
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
	logger := config.GetConfig().Logger
	return &Analytics{
		PolicyName:    "Analytics",
		PolicyVersion: mediation.PolicyVersion,
		PolicyID:      mediation.PolicyID,
		Enabled:       enabled,
		logger:        &logger,
	}
}

// Process processes the request configuration for analytics.
func (a *Analytics) Process(requestConfig *requestconfig.Holder) *Result {
	// Implement the logic to process the requestConfig for analytics
	// This is a placeholder implementation
	result := &Result{
		StopFurtherProcessing: false,
	}

	
	// Add logic to handle analytics here

	return result
}
