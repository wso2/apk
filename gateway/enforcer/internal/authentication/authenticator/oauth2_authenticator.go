package authenticator

import (
	"github.com/wso2/apk/gateway/enforcer/internal/datastore"
	"github.com/wso2/apk/gateway/enforcer/internal/requestconfig"
	"github.com/wso2/apk/gateway/enforcer/internal/transformer"
)

// OAuth2Authenticator is the main authenticator.
type OAuth2Authenticator struct {
	mandatory       bool
	jwtTransformer  *transformer.JWTTransformer
	revokedJTIStore *datastore.RevokedJTIStore
}

// NewOAuth2Authenticator creates a new OAuth2Authenticator.
func NewOAuth2Authenticator(jwtTransformer *transformer.JWTTransformer, revokedJTIStore *datastore.RevokedJTIStore, mandatory bool) *OAuth2Authenticator {
	return &OAuth2Authenticator{jwtTransformer: jwtTransformer, mandatory: mandatory, revokedJTIStore: revokedJTIStore}
}

const (
	// Oauth2AuthType is the Oauth2 authentication type.
	Oauth2AuthType = "oauth2"
)

// Authenticate performs the authentication.
func (authenticator *OAuth2Authenticator) Authenticate(rch *requestconfig.Holder) AuthenticationResponse {
	if rch != nil && rch.ExternalProcessingEnvoyMetadata != nil && rch.ExternalProcessingEnvoyMetadata.AuthenticationData != nil {
		jwtValidationInfo := authenticator.jwtTransformer.TransformJWTClaims(rch.MatchedAPI.OrganizationID, rch.ExternalProcessingEnvoyMetadata.AuthenticationData, Oauth2AuthType)
		if jwtValidationInfo != nil {
			if authenticator.revokedJTIStore != nil && authenticator.revokedJTIStore.IsJTIRevoked(jwtValidationInfo.JTI) {
				return AuthenticationResponse{Authenticated: false, MandatoryAuthentication: authenticator.mandatory, ErrorCode: ExpiredToken, ErrorMessage: ExpiredTokenMessage, ContinueToNextAuthenticator: false}
			}
			if jwtValidationInfo.Valid {
				rch.JWTValidationInfo = jwtValidationInfo
				rch.AuthenticatedAuthenticationType = Oauth2AuthType
				return AuthenticationResponse{Authenticated: true, MandatoryAuthentication: authenticator.mandatory, ContinueToNextAuthenticator: false}
			}
			return AuthenticationResponse{Authenticated: false, MandatoryAuthentication: authenticator.mandatory, ContinueToNextAuthenticator: false, ErrorCode: InvalidCredentials, ErrorMessage: InvalidCredentialsMessage}
		}
	}
	return AuthenticationResponse{Authenticated: false, MandatoryAuthentication: authenticator.mandatory, ErrorCode: MissingCredentials, ErrorMessage: MissingCredentialsMesage, ContinueToNextAuthenticator: !authenticator.mandatory}
}
