package model

import dpv1alpha1 "github.com/wso2/apk/adapter/pkg/operator/apis/dp/v1alpha1"

// RateLimit is the rate limit values for a policy
type RateLimit struct {
	// RequestPerUnit is the number of requests allowed per unit time
	//
	RequestsPerUnit int `json:"requestsPerUnit,omitempty"`

	// Unit is the unit of the requestPerUnit
	//
	Unit string `json:"unit,omitempty"`
}

// CustomRateLimitPolicy defines the desired state of CustomPolicy
type CustomRateLimitPolicy struct {
	Key  string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
	RateLimit RateLimit `json:"rateLimit,omitempty"`
	Organization string `json:"organization,omitempty"`
}

// ParseCustomRateLimitPolicy parses the custom rate limit policy
func ParseCustomRateLimitPolicy (customRateLimitCR dpv1alpha1.RateLimitPolicy) *CustomRateLimitPolicy {
	rlPolicy := concatRateLimitPolicies(&customRateLimitCR, nil)
	return &CustomRateLimitPolicy{
		Key: rlPolicy.Spec.Override.Custom.Key,
		Value: rlPolicy.Spec.Override.Custom.Value,
		RateLimit: RateLimit{
			RequestsPerUnit: rlPolicy.Spec.Override.Custom.RateLimit.RequestsPerUnit,
			Unit: rlPolicy.Spec.Override.Custom.RateLimit.Unit,
		},
		Organization: rlPolicy.Spec.Override.Organization,
	}
}