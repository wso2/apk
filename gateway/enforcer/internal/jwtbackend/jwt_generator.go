package jwtbackend

import (
	"fmt"
	"time"

	"github.com/wso2/apk/gateway/enforcer/internal/config"
	"github.com/wso2/apk/gateway/enforcer/internal/dto"
	"github.com/wso2/apk/gateway/enforcer/internal/requestconfig"
	"github.com/wso2/apk/gateway/enforcer/internal/util"
)

const (
	apiGatewayID   = "wso2.org/products/am"
	dialectURI     = "http://wso2.org/claims/"
	sha256WithRSA = "SHA256withRSA"
)

// CreateBackendJWT creates a JWT token for the backend.
func CreateBackendJWT(rch *requestconfig.Holder, cfg *config.Server) string {
	api := rch.MatchedAPI
	application := rch.MatchedApplication
	subscription := rch.MatchedSubscription

	if api != nil && api.BackendJwtConfiguration != nil && api.BackendJwtConfiguration.Enabled {
		bjc := api.BackendJwtConfiguration
		customClaims := bjc.CustomClaims
		if customClaims == nil {
			customClaims = make(map[string]*dto.ClaimValue)
		}
		customClaims["iss"] = &dto.ClaimValue{
			Value: apiGatewayID,
			Type:  "string",
		}
		currentTime := time.Now().Unix()
		expireIn := currentTime + bjc.TTL
		customClaims["exp"] = &dto.ClaimValue{
			Value: fmt.Sprintf("%d", expireIn),
			Type:  "int",
		}
		customClaims["iat"] = &dto.ClaimValue{
			Value: fmt.Sprintf("%d", currentTime),
			Type:  "int",
		}
		customClaims[dialectURI+"apiname"] = &dto.ClaimValue{
			Value: api.Name,
			Type:  "string",
		}
		customClaims[dialectURI+"apicontext"] = &dto.ClaimValue{
			Value: api.BasePath,
			Type:  "string",
		}
		customClaims[dialectURI+"version"] = &dto.ClaimValue{
			Value: api.Version,
			Type:  "string",
		}
		customClaims[dialectURI+"keytype"] = &dto.ClaimValue{
			Value: api.EnvType,
			Type:  "string",
		}
		if application != nil {
			customClaims[dialectURI+"subscriber"] = &dto.ClaimValue{
				Value: application.Owner,
				Type:  "string",
			}
			customClaims[dialectURI+"applicationid"] = &dto.ClaimValue{
				Value: application.UUID,
				Type:  "string",
			}
			customClaims[dialectURI+"applicationname"] = &dto.ClaimValue{
				Value: application.Name,
				Type:  "string",
			}
			customClaims[dialectURI+"applicationtier"] = &dto.ClaimValue{
				Value: subscription.RatelimitTier,
				Type:  "string",
			}
		}
		if subscription != nil {
			customClaims[dialectURI+"tier"] = &dto.ClaimValue{
				Value: subscription.RatelimitTier,
				Type:  "string",
			}
		}
		signatureAlgorithm := bjc.SignatureAlgorithm
		if signatureAlgorithm !=  "NONE" && signatureAlgorithm != sha256WithRSA {
			signatureAlgorithm = sha256WithRSA
		}

		return util.GenerateJWTToken(signatureAlgorithm, true, bjc.PublicCert, customClaims, bjc.PrivateKey)
	}
	return ""
}
