package transformer

import (
	"fmt"
	"strings"

	"github.com/wso2/apk/gateway/enforcer/internal/datastore"
	"github.com/wso2/apk/gateway/enforcer/internal/dto"
)

// JWTTransformer represents the JWT transformer.
type JWTTransformer struct {
	tokenissuerStore *datastore.JWTIssuerStore
}

// NewJWTTransformer creates a new instance of JWTIssuerStore.
func NewJWTTransformer(jwtIssuerDatastore *datastore.JWTIssuerStore) *JWTTransformer {
	return &JWTTransformer{tokenissuerStore: jwtIssuerDatastore}
}

// TransformJWTClaims transforms the JWT claims
func (transformer *JWTTransformer) TransformJWTClaims(organization string, jwtAuthenticationData *dto.JwtAuthenticationData) *dto.JWTValidationInfo {
	if jwtAuthenticationData == nil {
		fmt.Printf("JWT authentication data is nil\n")
		return nil
	}
	tokenissuers := transformer.tokenissuerStore.GetJWTISsuersByOrganization(organization)
	if tokenissuers == nil {
		fmt.Printf("Token issuers are nil\n")
		return nil
	}
	var jwtValidationInfo dto.JWTValidationInfo
	for _, tokenissuer := range tokenissuers {
		fmt.Printf("Token issuer: %v\n", tokenissuer)
		jwtAuthenticationDataSuccess, exists := jwtAuthenticationData.SucessData[tokenissuer.Issuer+"-payload"]
		if exists {
			jwtValidationInfo = dto.JWTValidationInfo{Valid: true, Issuer: jwtAuthenticationDataSuccess.Issuer, Claims: make(map[string]interface{})}
			remoteClaims := jwtAuthenticationDataSuccess.Claims
			if remoteClaims != nil {
				issuedTime := remoteClaims["iat"]
				if issuedTime != nil {
					jwtValidationInfo.IssuedTime = int64(issuedTime.(float64))
				}
				expiryTime := remoteClaims["exp"]
				if expiryTime != nil {
					jwtValidationInfo.ExpiryTime = int64(expiryTime.(float64))
				}
				jti := remoteClaims["jti"]
				if jti != nil {
					jwtValidationInfo.JTI = jti.(string)
				}
				audienceClaim := remoteClaims["aud"]
				if audienceClaim != nil {
					switch audienceClaim.(type) {
					case string:
						audiences := []string{remoteClaims["aud"].(string)}
						jwtValidationInfo.Audiences = audiences
					case []string:
						audiences := remoteClaims["aud"].([]string)
						jwtValidationInfo.Audiences = audiences
					}
				}
				remoteScopes := remoteClaims[tokenissuer.ScopesClaim]
				if remoteScopes != nil {
					switch remoteScopes := remoteScopes.(type) {
					case string:
						scopes := strings.Split(remoteScopes, " ")
						jwtValidationInfo.Scopes = scopes
					case []string:
						scopes := remoteScopes
						jwtValidationInfo.Scopes = scopes
					}
				}
				remoteClientID := remoteClaims[tokenissuer.ConsumerKeyClaim]
				if remoteClientID != nil {
					jwtValidationInfo.ClientID = remoteClientID.(string)
				}
				for claimKey, claimValue := range remoteClaims {
					if localClaim, ok := tokenissuer.ClaimMapping[claimKey]; ok {
						jwtValidationInfo.Claims[localClaim] = claimValue
					} else {
						jwtValidationInfo.Claims[claimKey] = claimValue
					}
				}
			}
			return &jwtValidationInfo
		}
		jwtAuthenticationDataFailure, exists := jwtAuthenticationData.FailedData[tokenissuer.Issuer+"-failed"]
		if exists {
			jwtValidationInfo = dto.JWTValidationInfo{Valid: false, ValidationCode: jwtAuthenticationDataFailure.Code, ValidationMessage: jwtAuthenticationDataFailure.Message}
		}
	}
	return &jwtValidationInfo
}
