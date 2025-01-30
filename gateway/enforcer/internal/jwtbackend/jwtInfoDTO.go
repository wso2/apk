package jwtbackend

import "github.com/wso2/apk/gateway/enforcer/internal/dto"

// JWTInfoDto holds information related to JWT tokens.
type JWTInfoDto struct {
	ApplicationTier   string                     `json:"applicationTier"`
	KeyType           string                     `json:"keyType"`
	Version           string                     `json:"version"`
	ApplicationName   string                     `json:"applicationName"`
	EndUser           string                     `json:"endUser"`
	EndUserTenantID   int                        `json:"endUserTenantId"`
	ApplicationUUID   string                     `json:"applicationUUId"`
	Subscriber        string                     `json:"subscriber"`
	SubscriptionTier  string                     `json:"subscriptionTier"`
	ApplicationID     string                     `json:"applicationId"`
	APIContext        string                     `json:"apiContext"`
	APIName           string                     `json:"apiName"`
	JwtValidationInfo *JWTValidationInfo         `json:"jwtValidationInfo"`
	AppAttributes     map[string]string          `json:"appAttributes"`
	Sub               string                     `json:"sub"`
	Organizations     []string                   `json:"organizations"`
	Claims            map[string]*dto.ClaimValue `json:"claims"`
}

// NewJWTInfoDto creates a new JWTInfoDto instance.
func NewJWTInfoDto() *JWTInfoDto {
	return &JWTInfoDto{
		AppAttributes: make(map[string]string),
		Claims:        make(map[string]*dto.ClaimValue),
		Organizations: make([]string, 0),
	}
}

// Clone creates a deep copy of the JWTInfoDto.
func (j *JWTInfoDto) Clone() *JWTInfoDto {
	clone := *j
	clone.AppAttributes = CloneStringMap(j.AppAttributes)
	clone.Claims = CloneClaimsMap(j.Claims)
	clone.Organizations = copyStringSlice(j.Organizations)
	if j.JwtValidationInfo != nil {
		clone.JwtValidationInfo = j.JwtValidationInfo.Clone()
	}
	return &clone
}

// CloneStringMap creates a copy of a string map.
func CloneStringMap(original map[string]string) map[string]string {
	clone := make(map[string]string)
	for k, v := range original {
		clone[k] = v
	}
	return clone
}

// CloneClaimsMap creates a deep copy of a map with ClaimValueDTO pointers.
func CloneClaimsMap(original map[string]*dto.ClaimValue) map[string]*dto.ClaimValue {
	clonedMap := make(map[string]*dto.ClaimValue)
	for k, v := range original {
		if v != nil {
			valueCopy := *v // Create a copy of the struct value
			clonedMap[k] = &valueCopy
		}
	}
	return clonedMap
}

// Helper function to copy a string slice.
func copyStringSlice(original []string) []string {
	return append([]string{}, original...)
}
