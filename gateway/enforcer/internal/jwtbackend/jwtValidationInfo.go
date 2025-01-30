package jwtbackend

import (
	"github.com/golang-jwt/jwt/v4"
)

// JWTValidationInfo holds JWT validation related information.
type JWTValidationInfo struct {
	User           string                 `json:"user"`
	ExpiryTime     int64                  `json:"expiryTime"`
	ConsumerKey    string                 `json:"consumerKey"`
	Valid          bool                   `json:"valid"`
	Scopes         []string               `json:"scopes"`
	Claims         map[string]interface{} `json:"claims"`
	ValidationCode int                    `json:"validationCode"`
	KeyManager     string                 `json:"keyManager"`
	Identifier     string                 `json:"identifier"`
	JWTClaimsSet   *jwt.MapClaims         `json:"jwtClaimsSet"`
	Token          string                 `json:"token"`
	Audience       []string               `json:"audience"`
}

// NewJWTValidationInfo creates a new instance of JWTValidationInfo.
func NewJWTValidationInfo() *JWTValidationInfo {
	return &JWTValidationInfo{
		Scopes:   make([]string, 0),
		Claims:   make(map[string]interface{}),
		Audience: make([]string, 0),
	}
}

// Clone creates a copy of an existing JWTValidationInfo.
func (j *JWTValidationInfo) Clone() *JWTValidationInfo {
	return &JWTValidationInfo{
		User:           j.User,
		ExpiryTime:     j.ExpiryTime,
		ConsumerKey:    j.ConsumerKey,
		Valid:          j.Valid,
		Scopes:         append([]string{}, j.Scopes...),
		Claims:         CloneMap(j.Claims),
		ValidationCode: j.ValidationCode,
		KeyManager:     j.KeyManager,
		Identifier:     j.Identifier,
		JWTClaimsSet:   j.JWTClaimsSet,
		Token:          j.Token,
		Audience:       append([]string{}, j.Audience...),
	}
}

// CloneMap creates a shallow copy of a map[string]interface{}.
func CloneMap(original map[string]interface{}) map[string]interface{} {
	clonedMap := make(map[string]interface{})
	for k, v := range original {
		clonedMap[k] = v
	}
	return clonedMap
}
