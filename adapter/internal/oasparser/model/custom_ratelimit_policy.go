package model

import dpv1alpha1 "github.com/wso2/apk/adapter/pkg/operator/apis/dp/v1alpha1"

type RateLimit struct {
	// RequestPerUnit is the number of requests allowed per unit time
	//
	RequestsPerUnit int `json:"requestsPerUnit,omitempty"`

	// Unit is the unit of the requestPerUnit
	//
	Unit string `json:"unit,omitempty"`
}

type CustomRateLimitPolicy struct {
	Key  string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
	RateLimit RateLimit `json:"rateLimit,omitempty"`
}

func ParseCustomRateLimitPolicy (customRateLimitCR dpv1alpha1.RateLimitPolicy) *CustomRateLimitPolicy {
	rlPolicy := concatRateLimitPolicies(&customRateLimitCR, nil)
	return &CustomRateLimitPolicy{
		Key: rlPolicy.Spec.Override.Custom.Key,
		Value: rlPolicy.Spec.Override.Custom.Value,
		RateLimit: RateLimit{
			RequestsPerUnit: rlPolicy.Spec.Override.Custom.RateLimit.RequestsPerUnit,
			Unit: rlPolicy.Spec.Override.Custom.RateLimit.Unit,
		},
	}
}