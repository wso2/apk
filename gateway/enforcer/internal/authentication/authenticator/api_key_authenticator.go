package authenticator

import (
	"strings"

	"github.com/wso2/apk/gateway/enforcer/internal/dto"
	"github.com/wso2/apk/gateway/enforcer/internal/requestconfig"
)

// APIKeyAuthenticator is the main authenticator.
type APIKeyAuthenticator struct {
	mandatory bool
}

// NewAPIKeyAuthenticator creates a new APIKeyAuthenticator.
func NewAPIKeyAuthenticator(mandatory bool) *APIKeyAuthenticator {
	return &APIKeyAuthenticator{mandatory: mandatory}
}

const (
	// APIKeyAuthType is the APIKey authentication type.
	APIKeyAuthType = "apikey"
)

// Authenticate performs the authentication.
func (authenticator *APIKeyAuthenticator) Authenticate(rch *requestconfig.Holder) AuthenticationResponse {
	if rch != nil && rch.ExternalProcessingEnvoyMetadata != nil && rch.ExternalProcessingEnvoyMetadata.AuthenticationData != nil {
		if rch.ExternalProcessingEnvoyMetadata.AuthenticationData.SucessData != nil && len(rch.ExternalProcessingEnvoyMetadata.AuthenticationData.SucessData) > 0 {
			apiKeyAuthenticationSucessData, exists := rch.ExternalProcessingEnvoyMetadata.AuthenticationData.SucessData["apikey-payload"]
			if exists {
				apiKeyInfo := extractAPIKeyAuthenticationInfo(apiKeyAuthenticationSucessData, nil)
				rch.APIKeyAuthenticationInfo = &apiKeyInfo
				return AuthenticationResponse{Authenticated: true, MandatoryAuthentication: authenticator.mandatory, ContinueToNextAuthenticator: false}
			}
		} else if rch.ExternalProcessingEnvoyMetadata.AuthenticationData.FailedData != nil && len(rch.ExternalProcessingEnvoyMetadata.AuthenticationData.FailedData) > 0 {
			apiKeyAuthenticationFailedData, exists := rch.ExternalProcessingEnvoyMetadata.AuthenticationData.FailedData["apikey-failed"]
			if exists {
				apiKeyInfo := extractAPIKeyAuthenticationInfo(nil, apiKeyAuthenticationFailedData)
				rch.APIKeyAuthenticationInfo = &apiKeyInfo
				rch.AuthenticatedAuthenticationType = APIKeyAuthType
				return AuthenticationResponse{Authenticated: false, MandatoryAuthentication: authenticator.mandatory, ErrorCode: InvalidCredentials, ErrorMessage: InvalidCredentialsMessage, ContinueToNextAuthenticator: false}
			}
		}

	}
	return AuthenticationResponse{Authenticated: false, MandatoryAuthentication: authenticator.mandatory, ErrorCode: MissingCredentials, ErrorMessage: MissingCredentialsMesage, ContinueToNextAuthenticator: true}
}

func extractAPIKeyAuthenticationInfo(authenticationData *dto.AuthenticationSuccessData, authenticationFailureData *dto.AuthenticationFailureData) dto.APIKeyAuthenticationInfo {
	apiKeyAuthenticationInfo := dto.APIKeyAuthenticationInfo{}
	if authenticationData != nil {
		keyType, exists := authenticationData.Claims["keytype"]
		if exists {
			if keyTypeStr, ok := keyType.(string); ok {
				apiKeyAuthenticationInfo.Keytype = keyTypeStr
			}
		}
		issuedTime, exists := authenticationData.Claims["iat"]
		if exists {
			if issuedTimeFloat, ok := issuedTime.(float64); ok {
				apiKeyAuthenticationInfo.IssuedTime = int64(issuedTimeFloat)
			}
		}
		expiryTime, exists := authenticationData.Claims["exp"]
		if exists {
			if expiryTimeFloat, ok := expiryTime.(float64); ok {
				apiKeyAuthenticationInfo.ExpiryTime = int64(expiryTimeFloat)
			}
		}
		permittedIP, exists := authenticationData.Claims["permittedIP"]
		if exists {
			if permittedIPStr, ok := permittedIP.(string); ok {
				apiKeyAuthenticationInfo.PermittedIP = strings.Split(permittedIPStr, ",")
			}
		}
		permittedReferer, exists := authenticationData.Claims["permittedReferer"]
		if exists {
			if permittedRefererStr, ok := permittedReferer.(string); ok {
				apiKeyAuthenticationInfo.PermittedReferer = strings.Split(permittedRefererStr, ",")
			}
		}
		application := authenticationData.Claims["application"]
		if application != nil {
			if applicationMap, ok := application.(map[string]interface{}); ok {
				apiKeyAuthenticationInfo.Application = &dto.ApplicationInfo{
					ID:   applicationMap["id"].(float64),
					UUID: applicationMap["uuid"].(string),
				}
			}
		}
	} else if authenticationFailureData != nil {
		apiKeyAuthenticationInfo.Valid = false
		apiKeyAuthenticationInfo.ValidationCode = authenticationFailureData.Code
		apiKeyAuthenticationInfo.ValidationMessage = authenticationFailureData.Message
	}
	return apiKeyAuthenticationInfo
}
