package jwtbackend

import (

	// "github.com/golang-jwt/jwt/v5"
	"github.com/wso2/apk/gateway/enforcer/internal/config"
	"github.com/wso2/apk/gateway/enforcer/internal/requestconfig"
	// "github.com/wso2/apk/gateway/enforcer/internal/util"
)

const (
	apiGatewayID  = "wso2.org/products/am"
	dialectURI    = "http://wso2.org/claims/"
	sha256WithRSA = "SHA256withRSA"
)

var restrictedClaims = []string{"iss", "sub", "aud", "exp", "nbf", "iat", "jti", "application", "tierInfo", "subscribedAPIs", "aut"}

// CreateBackendJWT creates a JWT token for the backend.
func CreateBackendJWT(rch *requestconfig.Holder, cfg *config.Server) string {
	// api := rch.MatchedAPI
	// application := rch.MatchedApplication
	// subscription := rch.MatchedSubscription
	// jwtClaims := jwt.MapClaims{}
	// if api != nil && api.BackendJwtConfiguration != nil && api.BackendJwtConfiguration.Enabled {
	// 	bjc := api.BackendJwtConfiguration
	// 	customClaims := bjc.CustomClaims
	// 	jwtClaims["iss"] = apiGatewayID
	// 	currentTime := time.Now().Unix()
	// 	expireIn := currentTime + bjc.TTL
	// 	jwtClaims["exp"] = expireIn
	// 	jwtClaims["iat"] = currentTime
	// 	jwtClaims[dialectURI+"apiname"] = api.Name
	// 	jwtClaims[dialectURI+"apicontext"] = api.BasePath
	// 	jwtClaims[dialectURI+"version"] = api.Version
	// 	jwtClaims[dialectURI+"keytype"] = api.EnvType
	// 	if application != nil {
	// 		jwtClaims[dialectURI+"subscriber"] = application.Owner
	// 		jwtClaims[dialectURI+"applicationid"] = application.UUID
	// 		jwtClaims[dialectURI+"applicationname"] = application.Name
	// 	}
	// 	if subscription != nil {
	// 		jwtClaims[dialectURI+"tier"] = subscription.RatelimitTier
	// 	}
	// 	if rch.JWTValidationInfo != nil {
	// 		if sub, exists := rch.JWTValidationInfo.Claims["sub"]; exists {
	// 			jwtClaims["sub"] = sub.(string)
	// 		}
	// 		for claim, claimValue := range rch.JWTValidationInfo.Claims {
	// 			if !util.Contains(restrictedClaims, claim) {
	// 				if claimValue, ok := claimValue.(string); ok {
	// 					jwtClaims[claim] = claimValue
	// 				}
	// 			}
	// 		}
	// 	}
	// 	if customClaims != nil {
	// 		for claim, claimValue := range customClaims {
	// 			jwtClaims[claim] = claimValue.Value
	// 		}
	// 	}
	// 	var signingMethod jwt.SigningMethod
	// 	signatureAlgorithm := bjc.SignatureAlgorithm
	// 	if signatureAlgorithm != "NONE" && signatureAlgorithm != sha256WithRSA {
	// 		signingMethod = jwt.SigningMethodRS256
	// 	} else if signatureAlgorithm == sha256WithRSA {
	// 		signingMethod = jwt.SigningMethodRS256
	// 	} else {
	// 		signingMethod = jwt.SigningMethodNone
	// 	}
	// 	token := jwt.NewWithClaims(signingMethod, jwtClaims)
	// 	if bjc.UseKid {
	// 		token.Header["kid"] = JWKKEy.KeyID()
	// 	}
	// 	signedToken, err := token.SignedString(bjc.PrivateKey)
	// 	if err != nil {
	// 		cfg.Logger.Sugar().Errorf("Failed to sign the JWT token: %v", err)
	// 		return ""
	// 	}
	// 	return signedToken
	// }
	return ""
}
