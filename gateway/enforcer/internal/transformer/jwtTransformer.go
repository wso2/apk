package transformer

import (
	"fmt"

	"github.com/wso2/apk/gateway/enforcer/internal/datastore"
	"github.com/wso2/apk/gateway/enforcer/internal/dto"
)

// JWTTransformer represents the JWT transformer
type JWTTransformer struct {
	tokenissuerStore *datastore.JWTIssuerStore
}

// NewJWTTransformer creates a new instance of JWTIssuerStore.
func NewJWTTransformer(jwtIssuerDatastore *datastore.JWTIssuerStore) *JWTTransformer {
	return &JWTTransformer{tokenissuerStore: jwtIssuerDatastore}
}

// TransformJWTClaims transforms the JWT claims
func (transformer *JWTTransformer) TransformJWTClaims(organization string, externalProcessingEnvoyMetadata *dto.ExternalProcessingEnvoyMetadata) dto.JWTValidationInfo {
	if externalProcessingEnvoyMetadata == nil {
		fmt.Printf("External processing envoy metadata is nil\n")
		return dto.JWTValidationInfo{}
	}
	if externalProcessingEnvoyMetadata.JwtAuthenticationData == nil {
		fmt.Printf("JWT authentication data is nil\n")
		return dto.JWTValidationInfo{}
	}
	if externalProcessingEnvoyMetadata.JwtAuthenticationData.Claims == nil {
		fmt.Printf("JWT claims are nil\n")
		return dto.JWTValidationInfo{}
	}
	fmt.Printf("Organization: %v\n", organization)
	fmt.Printf("External processing envoy metadata: %v\n", externalProcessingEnvoyMetadata)
	fmt.Printf("JWT authentication data: %v\n", externalProcessingEnvoyMetadata.JwtAuthenticationData)
	tokenIssuer := transformer.tokenissuerStore.GetJWTIssuerByOrganizationAndIssuer(organization, externalProcessingEnvoyMetadata.JwtAuthenticationData.Issuer)
	jwtValidationInfo := dto.JWTValidationInfo{Issuer: externalProcessingEnvoyMetadata.JwtAuthenticationData.Issuer, Claims: make(map[string]interface{})}
	if tokenIssuer != nil {
		fmt.Printf("Token issuer: %v\n", tokenIssuer)
		remoteClaims := externalProcessingEnvoyMetadata.JwtAuthenticationData.Claims
		if remoteClaims != nil {
			fmt.Printf("Remote claims: %v\n", remoteClaims)
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
					scopes := []string{remoteScopes}
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
	return jwtValidationInfo
}
