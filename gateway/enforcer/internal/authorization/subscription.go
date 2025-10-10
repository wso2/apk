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
				appID := getAppIDUsingConsumerKey(clientID, subAppDataStore, api, "OAuth2")
				if appID != "" {
					appMaps := subAppDataStore.GetApplicationMappings(api.OrganizationID, appID)
					for _, appMap := range appMaps {
						subscriptions := subAppDataStore.GetSubscriptions(api.OrganizationID, appMap.SubscriptionRef)
						for _, subscription := range subscriptions {
							subscribedAPI := subscription.SubscribedAPI
							if subscribedAPI.Name == api.Name {
								versionMatched, err := regexp.MatchString(subscribedAPI.Version, api.Version)
								if err == nil && versionMatched {
									rch.MatchedSubscription = subscription
									rch.MatchedApplication = subAppDataStore.GetApplication(api.OrganizationID, appID)
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

func getAppIDUsingConsumerKey(consumerKey string, subAppDatastore *datastore.SubscriptionApplicationDataStore, api *requestconfig.API, securityScheme string) string {
	appKeyMapKey := util.PrepareApplicationKeyMappingCacheKey(consumerKey, api.EnvType, securityScheme, api.Environment)
	appKeyMap := subAppDatastore.GetApplicationKeyMapping(api.OrganizationID, appKeyMapKey)
	if appKeyMap != nil {
		return appKeyMap.ApplicationUUID
	}
	return ""
}
func getApplicationForAPPUUID(api *requestconfig.API, applicationUUID string, subAppDatastore *datastore.SubscriptionApplicationDataStore) *subscription_model.Application {
	return subAppDatastore.GetApplication(api.OrganizationID, applicationUUID)
}
