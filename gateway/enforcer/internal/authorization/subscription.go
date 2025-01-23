package authorization

import (
	"github.com/wso2/apk/gateway/enforcer/internal/datastore"
	"github.com/wso2/apk/gateway/enforcer/internal/dto"
	"github.com/wso2/apk/gateway/enforcer/internal/requestconfig"
	"github.com/wso2/apk/gateway/enforcer/internal/util"
)

const (
	forbiddenMessage = "Resource forbidden"
)

// validateSubscription validates the subscription.
func validateSubscription(appID string, subAppDatastore *datastore.SubscriptionApplicationDataStore, rch *requestconfig.Holder) *dto.ImmediateResponse{
	api := rch.MatchedAPI
	appMaps := subAppDatastore.GetApplicationMappings(api.OrganizationID, appID)
	for _, appMap := range appMaps {
		subscriptions := subAppDatastore.GetSubscriptions(api.OrganizationID, appMap.SubscriptionRef)
		for _, subscription := range subscriptions {
			subscribedAPI := subscription.SubscribedAPI
			if subscribedAPI.Name == api.Name && subscribedAPI.Version == api.Version {
				rch.MatchedSubscription = subscription
				rch.MatchedApplication = subAppDatastore.GetApplication(api.OrganizationID, appID)
				return nil
			}
		}
		
	}
	return &dto.ImmediateResponse{
		StatusCode: 403,
		Message: forbiddenMessage,
	}
}

func getAppIDUsingConsumerKey(consumerKey string, subAppDatastore *datastore.SubscriptionApplicationDataStore, api *requestconfig.API, securityScheme string) string {
	appKeyMapKey :=  util.PrepareApplicationKeyMappingCacheKey(consumerKey, api.EnvType, securityScheme, api.Environment)
	appKeyMap := subAppDatastore.GetApplicationKeyMapping(api.OrganizationID,appKeyMapKey)
	if appKeyMap != nil {
		return appKeyMap.ApplicationIdentifier
	}
	return ""
}