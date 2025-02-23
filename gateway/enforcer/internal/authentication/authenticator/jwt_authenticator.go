package authenticator

import (
	"github.com/wso2/apk/gateway/enforcer/internal/datastore"
	"github.com/wso2/apk/gateway/enforcer/internal/requestconfig"
	"github.com/wso2/apk/gateway/enforcer/internal/transformer"
)

// JWTAuthenticator is the main authenticator.
type JWTAuthenticator struct {
	mandatory       bool
	jwtTransformer  *transformer.JWTTransformer
	revokedJTIStore *datastore.RevokedJTIStore
}

// NewJWTAuthenticator creates a new JWTAuthenticator.
func NewJWTAuthenticator(jwtTransformer *transformer.JWTTransformer, revokedJTIStore *datastore.RevokedJTIStore, mandatory bool) *JWTAuthenticator {
	return &JWTAuthenticator{jwtTransformer: jwtTransformer, mandatory: mandatory, revokedJTIStore: revokedJTIStore}
}

// Authenticate performs the authentication.
func (authenticator *JWTAuthenticator) Authenticate(rch *requestconfig.Holder) AuthenticationResponse {
	if rch != nil && rch.ExternalProcessingEnvoyMetadata != nil && rch.ExternalProcessingEnvoyMetadata.AuthenticationData != nil {
		jwtValidationInfo := authenticator.jwtTransformer.TransformJWTClaims(rch.MatchedAPI.OrganizationID, rch.ExternalProcessingEnvoyMetadata.AuthenticationData)
		if jwtValidationInfo != nil {
			if authenticator.revokedJTIStore != nil && authenticator.revokedJTIStore.IsJTIRevoked(jwtValidationInfo.JTI) {
				return AuthenticationResponse{Authenticated: false, MandatoryAuthentication: authenticator.mandatory, ErrorCode: ExpiredToken, ErrorMessage: ExpiredTokenMessage, ContinueToNextAuthenticator: false}
			}
			if jwtValidationInfo.Valid {
				rch.JWTValidationInfo = jwtValidationInfo
				return AuthenticationResponse{Authenticated: true, MandatoryAuthentication: authenticator.mandatory, ContinueToNextAuthenticator: false}
			}
			return AuthenticationResponse{Authenticated: false, MandatoryAuthentication: authenticator.mandatory, ContinueToNextAuthenticator: false, ErrorCode: InvalidCredentials, ErrorMessage: InvalidCredentialsMessage}
		}
	}
	return AuthenticationResponse{Authenticated: false, MandatoryAuthentication: authenticator.mandatory, ErrorCode: MissingCredentials, ErrorMessage: MissingCredentialsMesage, ContinueToNextAuthenticator: true}
}
