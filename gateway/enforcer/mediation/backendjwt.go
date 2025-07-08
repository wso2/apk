package mediation

import (
	"encoding/json"
	"strconv"

	dpv2alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v2alpha1"
	"github.com/wso2/apk/gateway/enforcer/internal/config"
	"github.com/wso2/apk/gateway/enforcer/internal/requestconfig"
)

type BackendJWT struct {
	PolicyName       string            `json:"policyName"`
	PolicyVersion    string            `json:"policyVersion"`
	PolicyID         string            `json:"policyID"`
	Enabled          bool              `json:"enabled"`
	Encoding         string            `json:"encoding"`
	Header           string            `json:"header"`
	SigningAlgorithm string            `json:"signingAlgorithm"`
	TokenTTL         int               `json:"tokenTTL"`
	CustomClaims     map[string]string `json:"customClaims"`
	ClaimMapping     map[string]string `json:"claimMapping"`
}

const (
	// BackendJWTPolicyKeyEnabled is the key for enabling/disabling the Backend JWT policy.
	BackendJWTPolicyKeyEnabled = "Enabled"
	// BackendJWTPolicyKeyEncoding is the key for specifying the encoding type (e.g., "HS256").
	BackendJWTPolicyKeyEncoding = "Encoding"
	// BackendJWTPolicyKeyHeader is the key for specifying the JWT header.
	BackendJWTPolicyKeyHeader = "Header"
	// BackendJWTPolicyKeySigningAlgorithm is the key for specifying the signing algorithm (e.g., "HS256").
	BackendJWTPolicyKeySigningAlgorithm = "SigningAlgorithm"
	// BackendJWTPolicyKeyTokenTTL is the key for specifying the token time-to-live (TTL) in seconds.
	BackendJWTPolicyKeyTokenTTL = "TokenTTL"
	// BackendJWTPolicyKeyCustomClaims is the key for specifying custom claims in the JWT.
	BackendJWTPolicyKeyCustomClaims = "CustomClaims"
	// BackendJWTPolicyKeyClaimMapping is the key for specifying claim mapping in the JWT.
	BackendJWTPolicyKeyClaimMapping = "ClaimMapping"
)

// NewBackendJWT creates a new BackendJWT instance with default values.
func NewBackendJWT(mediation *dpv2alpha1.Mediation) *BackendJWT {
	enabled := true
	if val, ok := extractPolicyValue(mediation.Parameters, BackendJWTPolicyKeyEnabled); ok {
		if val == "false" {
			enabled = false
		}
	}
	encoding := "HS256"
	if val, ok := extractPolicyValue(mediation.Parameters, BackendJWTPolicyKeyEncoding); ok {
		encoding = val
	}
	header := ""
	if val, ok := extractPolicyValue(mediation.Parameters, BackendJWTPolicyKeyHeader); ok {
		header = val
	}
	signingAlgorithm := "HS256"
	if val, ok := extractPolicyValue(mediation.Parameters, BackendJWTPolicyKeySigningAlgorithm); ok {
		signingAlgorithm = val
	}
	logger := config.GetConfig().Logger.Sugar()
	tokenTTL := 3600 // Default to 1 hour
	if val, ok := extractPolicyValue(mediation.Parameters, BackendJWTPolicyKeyTokenTTL); ok {
		if ttl, err := strconv.Atoi(val); err == nil {
			tokenTTL = ttl
		} else {
			logger.Errorf("Invalid TokenTTL value: %s, using default value of 3600 seconds", val)
		}
	}
	customClaims := make(map[string]string)
	if val, ok := extractPolicyValue(mediation.Parameters, BackendJWTPolicyKeyCustomClaims); ok {
		// Assuming val is a JSON string representing a map of custom claims
		if err := json.Unmarshal([]byte(val), &customClaims); err != nil {
			// Handle error, possibly log it
			logger.Errorf("Failed to unmarshal CustomClaims: %v, error: %v", val, err)
		}
	}
	claimMapping := make(map[string]string)
	if val, ok := extractPolicyValue(mediation.Parameters, BackendJWTPolicyKeyClaimMapping); ok {
		// Assuming val is a JSON string representing a map of claim mappings
		if err := json.Unmarshal([]byte(val), &claimMapping); err != nil {
			// Handle error, possibly log it
			logger.Errorf("Failed to unmarshal ClaimMapping: %v, error: %v", val, err)
		}
	}
	return &BackendJWT{
		PolicyName:       "BackendJWT",
		PolicyVersion:    mediation.PolicyVersion,
		PolicyID:         mediation.PolicyID,
		Enabled:          enabled,
		Encoding:         encoding,
		Header:           header,
		SigningAlgorithm: signingAlgorithm,
		TokenTTL:         tokenTTL,
		CustomClaims:     customClaims,
		ClaimMapping:     claimMapping,
	}
}

// Process processes the request configuration for Backend JWT.
func (b *BackendJWT) Process(requestConfig *requestconfig.Holder) *MediationResult {
	// Implement the logic to process the requestConfig for Backend JWT
	// This is a placeholder implementation
	result := &MediationResult{}

	// Add logic to handle JWT creation and signing here

	return result
}
