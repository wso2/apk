package transformer

import (
	"strings"

	"github.com/wso2/apk/gateway/enforcer/internal/config"
	"github.com/wso2/apk/gateway/enforcer/internal/datastore"
	"github.com/wso2/apk/gateway/enforcer/internal/dto"
)

// JWTTransformer represents the JWT transformer.
type JWTTransformer struct {
	tokenissuerStore *datastore.JWTIssuerStore
	cfg              *config.Server
}

// NewJWTTransformer creates a new instance of JWTIssuerStore.
func NewJWTTransformer(cfg *config.Server, jwtIssuerDatastore *datastore.JWTIssuerStore) *JWTTransformer {
	return &JWTTransformer{cfg: cfg, tokenissuerStore: jwtIssuerDatastore}
}

// TransformJWTClaims transforms the JWT claims
func (transformer *JWTTransformer) TransformJWTClaims(organization string, jwtAuthenticationData *dto.AuthenticationData, tokenType string) *dto.JWTValidationInfo {
	if jwtAuthenticationData == nil {
		return nil
	}
	tokenissuers := transformer.tokenissuerStore.GetJWTISsuersByOrganization(organization)
	if tokenissuers == nil {
		return nil
	}
	var jwtValidationInfoSucess *dto.JWTValidationInfo
	var jwtValidationInfoFailure *dto.JWTValidationInfo
	for _, tokenissuer := range tokenissuers {
		jwtAuthenticationDataSuccess, exists := jwtAuthenticationData.SucessData[tokenissuer.Issuer+"-"+tokenType+"-payload"]
		if exists {
			jwtValidationInfoSucess = &dto.JWTValidationInfo{Valid: true, Issuer: jwtAuthenticationDataSuccess.Issuer, Claims: make(map[string]interface{})}
			remoteClaims := jwtAuthenticationDataSuccess.Claims
			if remoteClaims != nil {
				issuedTime := remoteClaims["iat"]
				if issuedTime != nil {
					jwtValidationInfoSucess.IssuedTime = int64(issuedTime.(float64))
				}
				expiryTime := remoteClaims["exp"]
				if expiryTime != nil {
					jwtValidationInfoSucess.ExpiryTime = int64(expiryTime.(float64))
				}
				jti := remoteClaims["jti"]
				if jti != nil {
					jwtValidationInfoSucess.JTI = jti.(string)
				}
				audienceClaim := remoteClaims["aud"]
				if audienceClaim != nil {
					switch audienceClaim.(type) {
					case string:
						audiences := []string{remoteClaims["aud"].(string)}
						jwtValidationInfoSucess.Audiences = audiences
					case []string:
						audiences := remoteClaims["aud"].([]string)
						jwtValidationInfoSucess.Audiences = audiences
					}
				}
				remoteScopes := remoteClaims[tokenissuer.ScopesClaim]
				if remoteScopes != nil {
					switch remoteScopes := remoteScopes.(type) {
					case string:
						scopes := strings.Split(remoteScopes, " ")
						jwtValidationInfoSucess.Scopes = scopes
					case []string:
						scopes := remoteScopes
						jwtValidationInfoSucess.Scopes = scopes
					}
				}
				remoteClientID := remoteClaims[tokenissuer.ConsumerKeyClaim]
				if remoteClientID != nil {
					jwtValidationInfoSucess.ClientID = remoteClientID.(string)
				}
				for claimKey, claimValue := range remoteClaims {
					if localClaim, ok := tokenissuer.ClaimMapping[claimKey]; ok {
						jwtValidationInfoSucess.Claims[localClaim] = claimValue
					} else {
						jwtValidationInfoSucess.Claims[claimKey] = claimValue
					}
				}
			}
			transformer.cfg.Logger.Sugar().Debugf("JWT validation success for the issuer %s", jwtValidationInfoSucess)
			return jwtValidationInfoSucess
		}
		jwtAuthenticationDataFailure, exists := jwtAuthenticationData.FailedData[tokenissuer.Issuer+"-oauth2-failed"]
		if exists {
			jwtValidationInfoFailure = &dto.JWTValidationInfo{Valid: false, ValidationCode: jwtAuthenticationDataFailure.Code, ValidationMessage: jwtAuthenticationDataFailure.Message}
		}
	}
	if jwtValidationInfoFailure != nil {
		return jwtValidationInfoFailure
	}
	return nil
}

// GetTokenIssuerCount obtains the total token issuer count for metrics purposes.
func (transformer *JWTTransformer) GetTokenIssuerCount() int {
	return transformer.tokenissuerStore.GetJWTIssuerCount()
}
