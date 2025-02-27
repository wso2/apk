package authorization

import (
	"regexp"

	"github.com/wso2/apk/gateway/enforcer/internal/config"
	"github.com/wso2/apk/gateway/enforcer/internal/datastore"
	"github.com/wso2/apk/gateway/enforcer/internal/dto"
	"github.com/wso2/apk/gateway/enforcer/internal/requestconfig"
	"github.com/wso2/apk/gateway/enforcer/internal/util"
)

const (
	forbiddenMessage = "Resource forbidden"
)

// ValidateSubscription validates the subscription.
func ValidateSubscription(rch *requestconfig.Holder, subAppDataStore *datastore.SubscriptionApplicationDataStore, cfg *config.Server) *dto.ImmediateResponse {
	api := rch.MatchedAPI
	clientID := rch.JWTValidationInfo.ClientID
	if rch.MatchedAPI.SubscriptionValidation && rch.AuthenticatedAuthenticationType == "oauth2" {
		if rch.JWTValidationInfo.ClientID != "" {
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
			Message:    forbiddenMessage,
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
