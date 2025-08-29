package mediation

import (
	dpv2alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v2alpha1"
	"github.com/wso2/apk/gateway/enforcer/internal/config"
	"github.com/wso2/apk/gateway/enforcer/internal/requestconfig"

	"strings"
)

// BackendAPIKey represents the configuration for Backend API Key policy in the API Gateway.
type BackendAPIKey struct {
	PolicyName    string `json:"policyName"`
	PolicyVersion string `json:"policyVersion"`
	PolicyID      string `json:"policyID"`
	Enabled       bool   `json:"enabled"`
	In            string `json:"in"`
	InValue       string `json:"inValue"`
	APIKey        string `json:"apiKey"`
	cfg           *config.Server
}

const (
	// BackendAPIKeyPolicyKeyEnabled is the key for enabling/disabling the Backend API Key policy.
	BackendAPIKeyPolicyKeyEnabled = "Enabled"
	// BackendAPIKeyPolicyKeyIn is the key for specifying the location of the API Key (e.g., "header", "query").
	BackendAPIKeyPolicyKeyIn = "In"
	// BackendAPIKeyPolicyKeyInValue is the key for specifying the value of the API Key location.
	BackendAPIKeyPolicyKeyInValue = "InValue"
	// BackendAPIKeyPolicyKeyAPIKey is the key for specifying the API Key.
	BackendAPIKeyPolicyKeyAPIKey = "APIKey"
)

// NewBackendAPIKey creates a new BackendAPIKey instance with default values.
func NewBackendAPIKey(mediation *dpv2alpha1.Mediation) *BackendAPIKey {
	cfg := config.GetConfig()
	cfg.Logger.Sugar().Infof("Creating BackendAPIKey policy with mediation: %p", mediation)
	enabled := true
	if val, ok := extractPolicyValue(mediation.Parameters, BackendAPIKeyPolicyKeyEnabled); ok {
		if val == "false" {
			enabled = false
		}
	}
	in := "header"
	if val, ok := extractPolicyValue(mediation.Parameters, BackendAPIKeyPolicyKeyIn); ok {
		in = val
	}
	inValue := ""
	if val, ok := extractPolicyValue(mediation.Parameters, BackendAPIKeyPolicyKeyInValue); ok {
		inValue = val
	}
	apiKey := ""
	if val, ok := extractPolicyValue(mediation.Parameters, BackendAPIKeyPolicyKeyAPIKey); ok {
		apiKey = val
	}
	return &BackendAPIKey{
		PolicyName:    "BackendAPIKey",
		PolicyVersion: "v1",
		PolicyID:      mediation.PolicyID,
		Enabled:       enabled,
		In:            in,
		InValue:       inValue,
		APIKey:        apiKey,
		cfg:           cfg,
	}
}

// Process processes the request configuration for Backend API Key.
func (b *BackendAPIKey) Process(requestConfig *requestconfig.Holder) *Result {
	// Implement the logic to process the requestConfig for Backend API Key
	// This is a placeholder implementation
	result := NewResult()
	if !b.Enabled {
		b.cfg.Logger.Sugar().Debugf("Backend API Key policy is disabled. Skipping processing.")
		return result
	}
	if strings.ToLower(b.In) == "header" {
		result.AddHeaders[b.InValue] = b.APIKey
	}
	return result
}
