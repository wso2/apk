package mediation

import (
	"encoding/json"
	"regexp"

	dpv2alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v2alpha1"
	"github.com/wso2/apk/common-go-libs/constants"
	subscription_model "github.com/wso2/apk/common-go-libs/pkg/server/model"
	"github.com/wso2/apk/gateway/enforcer/internal/config"
	"github.com/wso2/apk/gateway/enforcer/internal/datastore"
	"github.com/wso2/apk/gateway/enforcer/internal/dto"
	"github.com/wso2/apk/gateway/enforcer/internal/logging"
	"github.com/wso2/apk/gateway/enforcer/internal/requestconfig"
	"github.com/wso2/apk/gateway/enforcer/internal/util"
)

var forbiddenMessage = dto.ErrorResponse{Code: 900908, ErrorMessage: "Resource forbidden", ErrorDescription: "User is NOT authorized to access the Resource. API Subscription validation failed."}
var forbiddenJSONMessage string
func init() {
	forbiddenJSONMessageBytes, _ := json.MarshalIndent(forbiddenMessage, "", "  ")
	forbiddenJSONMessage = string(forbiddenJSONMessageBytes)
}

// SubscriptionValidation represents the configuration for subscription validation in the API Gateway.
type SubscriptionValidation struct {
	PolicyName    string
	PolicyVersion string
	PolicyID      string
	Enabled       bool
	logger *logging.Logger
	cfg    *config.Server
}

const (
	// SubscriptionValidationPolicyKeyEnabled is the key for enabling/disabling the Subscription Validation policy.
	SubscriptionValidationPolicyKeyEnabled = "Enabled"
)

// NewSubscriptionValidation creates a new SubscriptionValidation instance with default values.
func NewSubscriptionValidation(mediation *dpv2alpha1.Mediation) *SubscriptionValidation {
	enabled := true
	if val, ok := extractPolicyValue(mediation.Parameters, SubscriptionValidationPolicyKeyEnabled); ok {
		if val == "false" {
			enabled = false
		}
	}
	cfg := config.GetConfig()
	logger := cfg.Logger
	return &SubscriptionValidation{
		PolicyName:    "SubscriptionValidation",
		PolicyVersion: "v1",
		PolicyID:      "subscription-validation",
		Enabled:       enabled,
		logger:        &logger,
		cfg:           cfg,
	}
}

// Process processes the request configuration for Subscription Validation.
func (s *SubscriptionValidation) Process(requestConfig *requestconfig.Holder) *Result {
	// Implement the logic to process the requestConfig for Subscription Validation
	// This is a placeholder implementation
	result := NewResult()
	if !s.Enabled {
		s.logger.Sugar().Debugf("Subscription Validation policy is disabled. Skipping processing.")
		return result
	}
	// get the consumer key.
	clientID := "" 
	for _, header := range requestConfig.RequestHeaders.Headers.Headers {
		if header.Key == constants.ClientIDHeaderKey {
			clientID = header.Value
			if clientID == "" {
				clientID = string(header.RawValue)
			}
			break
		}
	}
	subAppDataStore := datastore.GetSubAppDataStore(s.cfg)
	if clientID != "" {
		appID := s.getAppIDUsingConsumerKey(clientID, requestConfig.RouteMetadata.Spec.API.EnvType, requestConfig.RouteMetadata.Spec.API.Environment, requestConfig.RouteMetadata.Spec.API.Organization, "OAuth2")
		if appID != "" {
			
			appMaps := subAppDataStore.GetApplicationMappings(requestConfig.RouteMetadata.Spec.API.Organization, appID)
			for _, appMap := range appMaps {
				subscriptions := subAppDataStore.GetSubscriptions(requestConfig.RouteMetadata.Spec.API.Organization, appMap.SubscriptionRef)
				for _, subscription := range subscriptions {
					subscribedAPI := subscription.SubscribedAPI
					if subscribedAPI.Name == requestConfig.RouteMetadata.Spec.API.Name {
						versionMatched, err := regexp.MatchString(subscribedAPI.Version, requestConfig.RouteMetadata.Spec.API.Version)
						if err == nil && versionMatched {
							requestConfig.MatchedSubscription = subscription
							requestConfig.MatchedApplication = subAppDataStore.GetApplication(requestConfig.RouteMetadata.Spec.API.Organization, appID)
							return result
						}
					}
				}
			}
		}
	} else {
		application := requestConfig.JWTAuthnPayloaClaims["application"]
		if application != nil {
			if applicationMap, ok := application.(map[string]interface{}); ok {
				applicationUUID := applicationMap["uuid"].(string)
				application := s.getApplicationForAPPUUID(requestConfig.RouteMetadata.Spec.API.Organization, applicationUUID)
				s.cfg.Logger.Sugar().Debugf("Application %+v", application)
				if application != nil {
					applicationMappings := subAppDataStore.GetApplicationMappings(requestConfig.RouteMetadata.Spec.API.Organization, application.UUID)
					s.cfg.Logger.Sugar().Debugf("Application Mappings %+v", applicationMappings)
					if applicationMappings != nil && len(applicationMappings) > 0 {
						for _, applicationMapping := range applicationMappings {
							subscriptions := subAppDataStore.GetSubscriptions(requestConfig.RouteMetadata.Spec.API.Organization, applicationMapping.SubscriptionRef)
							s.cfg.Logger.Sugar().Debugf("Subscriptions %+v", subscriptions)
							for _, subscription := range subscriptions {
								subscribedAPI := subscription.SubscribedAPI
								if subscribedAPI.Name == requestConfig.RouteMetadata.Spec.API.Name {
									versionMatched, err := regexp.MatchString(subscribedAPI.Version, requestConfig.RouteMetadata.Spec.API.Version)
									if err == nil && versionMatched {
										requestConfig.MatchedSubscription = subscription
										requestConfig.MatchedApplication = application
										s.cfg.Logger.Sugar().Debugf("Matched Subscription %+v", requestConfig.MatchedSubscription)
										return result
									}
								}
							}
						}
					}
				}
			}
		}
	}

	result.StopFurtherProcessing = true
	result.ImmediateResponse = true
	result.ImmediateResponseCode = 401
	result.ImmediateResponseBody = forbiddenJSONMessage
	result.ImmediateResponseDetail = "User is NOT authorized to access the Resource. API Subscription validation failed."

	return result
}


func (s *SubscriptionValidation) getAppIDUsingConsumerKey(consumerKey string,  envType, environment, organization, securityScheme string) string {
	subAppDatastore :=  datastore.GetSubAppDataStore(s.cfg)
	appKeyMapKey := util.PrepareApplicationKeyMappingCacheKey(consumerKey, envType, securityScheme, environment)
	appKeyMap := subAppDatastore.GetApplicationKeyMapping(organization, appKeyMapKey)
	if appKeyMap != nil {
		return appKeyMap.ApplicationUUID
	}
	return ""
}
func (s *SubscriptionValidation) getApplicationForAPPUUID(organization string, applicationUUID string) *subscription_model.Application {
	subAppDatastore :=  datastore.GetSubAppDataStore(s.cfg)
	return subAppDatastore.GetApplication(organization, applicationUUID)
}
