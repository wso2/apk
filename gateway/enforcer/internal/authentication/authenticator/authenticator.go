package authenticator

import (
	"encoding/json"

	"github.com/wso2/apk/gateway/enforcer/internal/config"
	"github.com/wso2/apk/gateway/enforcer/internal/datastore"
	"github.com/wso2/apk/gateway/enforcer/internal/dto"
	"github.com/wso2/apk/gateway/enforcer/internal/requestconfig"
	"github.com/wso2/apk/gateway/enforcer/internal/transformer"
)

// Authenticator is the main authenticator.
type Authenticator struct {
	subAppDataStore *datastore.SubscriptionApplicationDataStore
	jwtTransformer  *transformer.JWTTransformer
	revokedJTIStore *datastore.RevokedJTIStore
	cfg             *config.Server
}

// NewAuthenticator creates a new Authenticator.
func NewAuthenticator(cfg *config.Server, subAppDataStore *datastore.SubscriptionApplicationDataStore,
	jwtTransformer *transformer.JWTTransformer,
	revokedJTIStore *datastore.RevokedJTIStore) *Authenticator {
	return &Authenticator{cfg: cfg, subAppDataStore: subAppDataStore, jwtTransformer: jwtTransformer, revokedJTIStore: revokedJTIStore}
}

// Authenticate performs the authentication.
func (authenticator *Authenticator) Authenticate(rch *requestconfig.Holder) *dto.ImmediateResponse {

	if rch != nil && rch.MatchedAPI != nil && rch.MatchedAPI.IsGraphQLAPI() {
		applicationSecurity := rch.MatchedAPI.ApplicationSecurity
		var optionalAuthenticationResponse *AuthenticationResponse
		var authenticationResponse AuthenticationResponse
		var authenticationResponses []AuthenticationResponse
		var authenticated bool
		var mandatoryAuthFailed bool
		authenticator.cfg.Logger.Sugar().Debugf("API Security enabled for the API %+v", applicationSecurity)
		for authenticationType, mandatory := range applicationSecurity {
			if authenticationType == "OAuth2" {
				authenticator.cfg.Logger.Sugar().Debugf("OAuth2 authentication is enabled for the API")
				OAuth2Authenticator := NewOAuth2Authenticator(authenticator.jwtTransformer, authenticator.revokedJTIStore, mandatory)
				authenticationResponse = OAuth2Authenticator.Authenticate(rch)
				authenticator.cfg.Logger.Sugar().Debugf("OAuth2 authentication response %+v", authenticationResponse)
			} else if authenticationType == "JWT" {
				authenticator.cfg.Logger.Sugar().Debugf("JWT authentication is enabled for the API")
				JWTAuthenticator := NewJWTAuthenticator(authenticator.jwtTransformer, authenticator.revokedJTIStore, mandatory)
				authenticationResponse = JWTAuthenticator.Authenticate(rch)
				authenticator.cfg.Logger.Sugar().Debugf("JWT authentication response %+v", authenticationResponse)
			} else if authenticationType == "APIKey" {
				authenticator.cfg.Logger.Sugar().Debugf("APIKey authentication is enabled for the API")
				APIKeyAuthenticator := NewAPIKeyAuthenticator(mandatory)
				authenticationResponse = APIKeyAuthenticator.Authenticate(rch)
				authenticator.cfg.Logger.Sugar().Debugf("APIKey authentication response %+v", authenticationResponse)
			}
			if optionalAuthenticationResponse == nil && (authenticationResponse.Authenticated || authenticationResponse.ErrorCode != MissingCredentials) {
				optionalAuthenticationResponse = &authenticationResponse
			}
			if authenticationResponse.MandatoryAuthentication {
				authenticated = authenticationResponse.Authenticated
				mandatoryAuthFailed = !authenticated
			} else if !authenticated && !mandatoryAuthFailed {
				authenticated = authenticationResponse.Authenticated
			}
			if !authenticationResponse.Authenticated {
				authenticationResponses = append(authenticationResponses, authenticationResponse)
			}
			if !authenticationResponse.ContinueToNextAuthenticator {
				break
			}
		}
		if !authenticated {
			if mandatoryAuthFailed || optionalAuthenticationResponse == nil {
				authenticator.cfg.Logger.Sugar().Debugf("Authentication failed for the request. Responses: %+v", authenticationResponses)
				errorResponse := getError(authenticationResponses)
				jsonData, _ := json.MarshalIndent(errorResponse, "", "  ")
				return &dto.ImmediateResponse{StatusCode: 401, Message: string(jsonData)}

			} else if !(optionalAuthenticationResponse.Authenticated) {
				authenticator.cfg.Logger.Sugar().Debugf("Authentication failed for the request. Responses: %+v", authenticationResponses)
				errorResponse := &dto.ErrorResponse{ErrorMessage: optionalAuthenticationResponse.ErrorMessage, Code: optionalAuthenticationResponse.ErrorCode, ErrorDescription: "Make sure you have provided the correct security credentials"}
				jsonData, _ := json.MarshalIndent(errorResponse, "", "  ")
				return &dto.ImmediateResponse{StatusCode: 401, Message: string(jsonData)}
			}
		}
	}
	return nil
}
func getError(authenticationResponses []AuthenticationResponse) *dto.ErrorResponse {
	var immediateResponse *dto.ErrorResponse
	missingCredential := false
	for _, authenticationResponse := range authenticationResponses {
		if !authenticationResponse.ContinueToNextAuthenticator {
			return &dto.ErrorResponse{Code: authenticationResponse.ErrorCode, ErrorMessage: authenticationResponse.ErrorMessage, ErrorDescription: "Make sure you have provided the correct security credentials"}
		}
		if authenticationResponse.MandatoryAuthentication && authenticationResponse.ErrorCode != MissingCredentials {
			immediateResponse = &dto.ErrorResponse{Code: authenticationResponse.ErrorCode, ErrorMessage: authenticationResponse.ErrorMessage, ErrorDescription: "Make sure you have provided the correct security credentials"}
		} else {
			missingCredential = true
		}
	}
	if immediateResponse != nil && missingCredential {
		return &dto.ErrorResponse{Code: InvalidCredentials, ErrorMessage: InvalidCredentialsMessage, ErrorDescription: "Make sure you have provided the correct security credentials"}
	}
	return &dto.ErrorResponse{Code: InvalidCredentials, ErrorMessage: InvalidCredentialsMessage, ErrorDescription: "Make sure you have provided the correct security credentials"}
}
