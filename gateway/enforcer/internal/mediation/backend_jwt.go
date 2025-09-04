package mediation

import (
	"encoding/json"
	"strconv"
	"time"

	"crypto/rsa"

	"github.com/golang-jwt/jwt/v5"
	dpv2alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v2alpha1"
	"github.com/wso2/apk/gateway/enforcer/internal/config"
	"github.com/wso2/apk/gateway/enforcer/internal/jwtbackend"
	"github.com/wso2/apk/gateway/enforcer/internal/requestconfig"
	"github.com/wso2/apk/gateway/enforcer/internal/util"
)

var restrictedClaims = []string{"iss", "sub", "aud", "exp", "nbf", "iat", "jti", "application", "tierInfo", "subscribedAPIs", "aut"}

// BackendJWT represents the configuration for Backend JWT policy in the API Gateway.
type BackendJWT struct {
	PolicyName       string                 `json:"policyName"`
	PolicyVersion    string                 `json:"policyVersion"`
	PolicyID         string                 `json:"policyID"`
	Enabled          bool                   `json:"enabled"`
	Encoding         string                 `json:"encoding"`
	Header           string                 `json:"header"`
	SigningAlgorithm string                 `json:"signingAlgorithm"`
	TokenTTL         int                    `json:"tokenTTL"`
	CustomClaims     map[string]interface{} `json:"customClaims"`
	ClaimMapping     map[string]string      `json:"claimMapping"`
	UseKid           bool                   `json:"useKid"`
	PrivateKey       *rsa.PrivateKey
	cfg              *config.Server
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
	cfg := config.GetConfig()
	cfg.Logger.Sugar().Infof("Creating BackendJWT policy with mediation: %p", mediation)
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
	customClaims := make(map[string]interface{})
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
	useKid := false
	if val, ok := extractPolicyValue(mediation.Parameters, "UseKid"); ok {
		if val == "true" {
			useKid = true
		}
	}
	privateKey, err := util.LoadPrivateKey(cfg.JWTGeneratorPrivateKeyPath)
	if err != nil {
		logger.Errorf("Failed to load private key for Backend JWT: %v", err)
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
		UseKid:           useKid,
		PrivateKey:       privateKey,
		cfg:              cfg,
	}
}

// Process processes the request configuration for Backend JWT.
func (b *BackendJWT) Process(requestConfig *requestconfig.Holder) *Result {
	// Implement the logic to process the requestConfig for Backend JWT
	// This is a placeholder implementation
	result := NewResult()

	jwt := b.createJWT(requestConfig)
	if jwt != "" {
		if b.Header != "" {
			result.AddHeaders[b.Header] = "Bearer " + jwt
		} else {
			result.AddHeaders["Authorization"] = "Bearer " + jwt
		}
	} else {
		b.cfg.Logger.Sugar().Error("Failed to create JWT token")
	}

	return result
}

const (
	apiGatewayID  = "wso2.org/products/am"
	dialectURI    = "http://wso2.org/claims/"
	sha256WithRSA = "SHA256withRSA"
)

func (b *BackendJWT) createJWT(rch *requestconfig.Holder) string {
	application := rch.MatchedApplication
	subscription := rch.MatchedSubscription
	jwtClaims := jwt.MapClaims{}

	customClaims := b.CustomClaims
	jwtClaims["iss"] = apiGatewayID
	currentTime := time.Now().Unix()
	expireIn := currentTime + int64(b.TokenTTL)
	jwtClaims["exp"] = expireIn
	jwtClaims["iat"] = currentTime
	jwtClaims[dialectURI+"apiname"] = rch.RouteMetadata.Spec.API.Name
	jwtClaims[dialectURI+"apicontext"] = rch.RouteMetadata.Spec.API.Context
	jwtClaims[dialectURI+"version"] = rch.RouteMetadata.Spec.API.Version
	jwtClaims[dialectURI+"keytype"] = rch.RouteMetadata.Spec.API.Environment
	if application != nil {
		jwtClaims[dialectURI+"subscriber"] = application.Owner
		jwtClaims[dialectURI+"applicationid"] = application.UUID
		jwtClaims[dialectURI+"applicationname"] = application.Name
	}
	if subscription != nil {
		jwtClaims[dialectURI+"tier"] = subscription.RatelimitTier
	}
	// if rch.JWTValidationInfo != nil {
	// 	if sub, exists := rch.JWTValidationInfo.Claims["sub"]; exists {
	// 		jwtClaims["sub"] = sub.(string)
	// 	}
	// 	for claim, claimValue := range rch.JWTValidationInfo.Claims {
	// 		if !util.Contains(restrictedClaims, claim) {
	// 			if claimValue, ok := claimValue.(string); ok {
	// 				jwtClaims[claim] = claimValue
	// 			}
	// 		}
	// 	}
	// }
	if customClaims != nil {
		for claim, claimValue := range customClaims {
			jwtClaims[claim] = claimValue
		}
	}
	var signingMethod jwt.SigningMethod
	signatureAlgorithm := b.SigningAlgorithm
	if signatureAlgorithm != "NONE" && signatureAlgorithm != sha256WithRSA {
		signingMethod = jwt.SigningMethodRS256
	} else if signatureAlgorithm == sha256WithRSA {
		signingMethod = jwt.SigningMethodRS256
	} else {
		signingMethod = jwt.SigningMethodNone
	}
	token := jwt.NewWithClaims(signingMethod, jwtClaims)
	if b.UseKid {
		token.Header["kid"] = jwtbackend.JWKKEy.KeyID()
	}

	signedToken, err := token.SignedString(b.PrivateKey)
	if err != nil {
		b.cfg.Logger.Sugar().Errorf("Failed to sign the JWT token: %v", err)
		return ""
	}
	return signedToken
}
