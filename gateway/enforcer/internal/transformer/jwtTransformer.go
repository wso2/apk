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
	if jwtAuthenticationData.Status != nil {
		return &dto.JWTValidationInfo{Valid: false, ValidationCode: jwtAuthenticationData.Status.Code, ValidationMessage: jwtAuthenticationData.Status.Message}
	}
	if jwtAuthenticationData.Claims == nil {
		fmt.Printf("JWT claims are nil\n")
		return nil
	}
	fmt.Printf("Organization: %v\n", organization)
	fmt.Printf("JWT authentication data: %v\n", jwtAuthenticationData)
	tokenIssuer := transformer.tokenissuerStore.GetJWTIssuerByOrganizationAndIssuer(organization, jwtAuthenticationData.Issuer)
	jwtValidationInfo := dto.JWTValidationInfo{Valid: true, Issuer: jwtAuthenticationData.Issuer, Claims: make(map[string]interface{})}
	if tokenIssuer != nil {
		fmt.Printf("Token issuer: %v\n", tokenIssuer)
		remoteClaims := jwtAuthenticationData.Claims
		if remoteClaims != nil {
			fmt.Printf("Remote claims: %v\n", remoteClaims)
			issuedTime := remoteClaims["iat"]
			if issuedTime != nil {
				fmt.Printf("Issued time: %v\n", issuedTime)
				jwtValidationInfo.IssuedTime = int64(issuedTime.(float64))
			}
			expiryTime := remoteClaims["exp"]
			if expiryTime != nil {
				fmt.Printf("Expiry time: %v\n", expiryTime)
				jwtValidationInfo.ExpiryTime = int64(expiryTime.(float64))
			}
			jti := remoteClaims["jti"]
			if jti != nil {
				fmt.Printf("JTI: %v\n", jti)
				jwtValidationInfo.JTI = jti.(string)
			}
			audienceClaim := remoteClaims["aud"]
			if audienceClaim != nil {
				fmt.Printf("Audience claim: %v\n", audienceClaim)
				switch audienceClaim.(type) {
				case string:
					audiences := []string{remoteClaims["aud"].(string)}
					jwtValidationInfo.Audiences = audiences
				case []string:
					audiences := remoteClaims["aud"].([]string)
					jwtValidationInfo.Audiences = audiences
				}
			}
			remoteScopes := remoteClaims[tokenIssuer.ScopesClaim]
			if remoteScopes != nil {
				fmt.Printf("Remote scopes: %v\n", remoteScopes)
				switch remoteScopes := remoteScopes.(type) {
				case string:
					scopes := strings.Split(remoteScopes, " ")
					jwtValidationInfo.Scopes = scopes
				case []string:
					scopes := remoteScopes
					jwtValidationInfo.Scopes = scopes
				}
			}
			remoteClientID := remoteClaims[tokenIssuer.ConsumerKeyClaim]
			if remoteClientID != nil {
				fmt.Printf("Remote client ID: %v\n", remoteClientID)
				jwtValidationInfo.ClientID = remoteClientID.(string)
			}
			for claimKey, claimValue := range remoteClaims {
				fmt.Printf("Claim key: %v, Claim value: %v\n", claimKey, claimValue)
				if localClaim, ok := tokenIssuer.ClaimMapping[claimKey]; ok {
					jwtValidationInfo.Claims[localClaim] = claimValue
				} else {
					jwtValidationInfo.Claims[claimKey] = claimValue
				}
			}
		}
	} else {
		fmt.Printf("Token issuer is nil\n")
	}
	fmt.Printf("JWT validation info: %v\n", jwtValidationInfo)
	return &jwtValidationInfo
}
