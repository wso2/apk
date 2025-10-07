package authorization

import (
	"encoding/json"
	"regexp"

	subscription_model "github.com/wso2/apk/common-go-libs/pkg/server/model"
	"github.com/wso2/apk/gateway/enforcer/internal/authentication/authenticator"
	"github.com/wso2/apk/gateway/enforcer/internal/config"
	"github.com/wso2/apk/gateway/enforcer/internal/datastore"
	"github.com/wso2/apk/gateway/enforcer/internal/dto"
	"github.com/wso2/apk/gateway/enforcer/internal/requestconfig"
	"github.com/wso2/apk/gateway/enforcer/internal/util"
)

var forbiddenMessage = dto.ErrorResponse{Code: 900908, ErrorMessage: "Resource forbidden", ErrorDescription: "User is NOT authorized to access the Resource. API Subscription validation failed."}

// ValidateSubscription validates the subscription.
func ValidateSubscription(rch *requestconfig.Holder, subAppDataStore *datastore.SubscriptionApplicationDataStore, cfg *config.Server) *dto.ImmediateResponse {

	forbiddenJSONMessage, _ := json.MarshalIndent(forbiddenMessage, "", "  ")
	api := rch.MatchedAPI
	if rch.MatchedAPI.SubscriptionValidation {
		cfg.Logger.Sugar().Debugf("Subscription validation enabled for the API %+v", api)
		cfg.Logger.Sugar().Debugf("Authentication Type %+v", rch.AuthenticatedAuthenticationType)
		if rch.AuthenticatedAuthenticationType == authenticator.Oauth2AuthType {
			clientID := rch.JWTValidationInfo.ClientID
			if clientID != "" {
				appKeyInfo := getAppIDUsingConsumerKey(clientID, subAppDataStore, api, "OAuth2", cfg)
				if appKeyInfo != nil {
					appMaps := subAppDataStore.GetApplicationMappings(api.OrganizationID, appKeyInfo.AppID)
					for _, appMap := range appMaps {
						subscriptions := subAppDataStore.GetSubscriptions(api.OrganizationID, appMap.SubscriptionRef)
						for _, subscription := range subscriptions {
							subscribedAPI := subscription.SubscribedAPI
							if subscribedAPI.Name == api.Name {
								versionMatched, err := regexp.MatchString(subscribedAPI.Version, api.Version)
								if err == nil && versionMatched {
									// Check subscription status
									cfg.Logger.Sugar().Debugf("Subscription status %+v, Environment: %s", subscription.SubStatus, appKeyInfo.KeyType)
									if !isSubscriptionActive(subscription.SubStatus, appKeyInfo.KeyType) {
										cfg.Logger.Sugar().Debugf("Subscription is not active. Status: %s, Key Type: %s", subscription.SubStatus, appKeyInfo.KeyType)
										return &dto.ImmediateResponse{
											StatusCode: 403,
											Message:    string(forbiddenJSONMessage),
										}
									}
									cfg.Logger.Sugar().Debugf("Subscription validation successful! Status: %s, Key Type: %s", subscription.SubStatus, appKeyInfo.KeyType)
									rch.MatchedSubscription = subscription
									rch.MatchedApplication = subAppDataStore.GetApplication(api.OrganizationID, appKeyInfo.AppID)
									return nil
								}
							}
						}
					}
				}
			}
			return &dto.ImmediateResponse{
				StatusCode: 403,
				Message:    string(forbiddenJSONMessage),
			}
		} else if rch.AuthenticatedAuthenticationType == authenticator.APIKeyAuthType {
			apiKeyAuthenticationInfo := rch.APIKeyAuthenticationInfo
			cfg.Logger.Sugar().Debugf("API Key Authentication Info %+v", apiKeyAuthenticationInfo)
			if apiKeyAuthenticationInfo != nil {
				application := getApplicationForAPPUUID(api, apiKeyAuthenticationInfo.Application.UUID, subAppDataStore)
				cfg.Logger.Sugar().Debugf("Application %+v", application)
				if application != nil {
					applicationMappings := subAppDataStore.GetApplicationMappings(api.OrganizationID, application.UUID)
					cfg.Logger.Sugar().Debugf("Application Mappings %+v", applicationMappings)
					if applicationMappings != nil && len(applicationMappings) > 0 {
						for _, applicationMapping := range applicationMappings {
							subscriptions := subAppDataStore.GetSubscriptions(api.OrganizationID, applicationMapping.SubscriptionRef)
							cfg.Logger.Sugar().Debugf("Subscriptions %+v", subscriptions)
							for _, subscription := range subscriptions {
								subscribedAPI := subscription.SubscribedAPI
								if subscribedAPI.Name == api.Name {
									versionMatched, err := regexp.MatchString(subscribedAPI.Version, api.Version)
									if err == nil && versionMatched {
										// Check subscription status
										keyType := "DEFAULT" // Default assumption, might need to be determined differently
										cfg.Logger.Sugar().Debugf("Subscription status %+v, Key Type: %s (API Key)", subscription.SubStatus, keyType)
										if !isSubscriptionActive(subscription.SubStatus, keyType) {
											cfg.Logger.Sugar().Debugf("Subscription is not active. Status: %s", subscription.SubStatus)
											return &dto.ImmediateResponse{
												StatusCode: 403,
												Message:    string(forbiddenJSONMessage),
											}
										}
										rch.MatchedSubscription = subscription
										rch.MatchedApplication = application
										cfg.Logger.Sugar().Debugf("Matched Subscription %+v", rch.MatchedSubscription)
										return nil
									}
								}
							}
						}
					}
				}
			}
			return &dto.ImmediateResponse{
				StatusCode: 403,
				Message:    string(forbiddenJSONMessage),
			}
		}
	}
	return nil
}

// Helper function to check if subscription status allows API access
func isSubscriptionActive(subStatus string, keyType string) bool {
	switch subStatus {
	case "UNBLOCKED":
		return true
	case "BLOCKED":
		return false
	case "PROD_ONLY_BLOCKED":
		return keyType == "SANDBOX"
	default:
		return true
	}
}

// AppKeyInfo stores the application ID and key type
type AppKeyInfo struct {
	AppID   string
	KeyType string
}

func getAppIDUsingConsumerKey(consumerKey string, subAppDatastore *datastore.SubscriptionApplicationDataStore, api *requestconfig.API, securityScheme string, cfg *config.Server) *AppKeyInfo {
	// Try both possible key types since we don't know which one the client is using
	keyTypes := []string{"PRODUCTION", "SANDBOX"}

	cfg.Logger.Sugar().Debugf("Looking up application for consumerKey=%s, api.EnvType=%s, api.Environment=%s, securityScheme=%s",
		consumerKey, api.EnvType, api.Environment, securityScheme)

	for _, keyType := range keyTypes {
		appKeyMapKey := util.PrepareApplicationKeyMappingCacheKey(consumerKey, keyType, securityScheme, api.Environment)
		cfg.Logger.Sugar().Debugf("Trying cache key: %s (keyType=%s)", appKeyMapKey, keyType)
		appKeyMap := subAppDatastore.GetApplicationKeyMapping(api.OrganizationID, appKeyMapKey)
		if appKeyMap != nil {
			cfg.Logger.Sugar().Debugf("Found app mapping: UUID=%s, KeyType=%s", appKeyMap.ApplicationUUID, appKeyMap.KeyType)
			return &AppKeyInfo{
				AppID:   appKeyMap.ApplicationUUID,
				KeyType: appKeyMap.KeyType,
			}
		}
	}
	cfg.Logger.Sugar().Warnf("No application mapping found for consumerKey=%s with any key type", consumerKey)
	return nil

}
func getApplicationForAPPUUID(api *requestconfig.API, applicationUUID string, subAppDatastore *datastore.SubscriptionApplicationDataStore) *subscription_model.Application {
	return subAppDatastore.GetApplication(api.OrganizationID, applicationUUID)
}
