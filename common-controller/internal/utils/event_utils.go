package utils

import (
	time "time"

	"github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/subscription"
	"github.com/wso2/apk/common-controller/internal/loggers"
	cpv1alpha2 "github.com/wso2/apk/common-controller/internal/operator/apis/cp/v1alpha2"
	constants "github.com/wso2/apk/common-controller/internal/operator/constant"
)

// SendAppDeletionEvent sends an application creation event to the enforcer
func SendAppDeletionEvent(applicationUUID string, applicationSpec cpv1alpha2.ApplicationSpec) {
	currentTime := time.Now()
	milliseconds := currentTime.UnixNano() / int64(time.Millisecond)
	event := subscription.Event{
		Uuid:      applicationUUID,
		Type:      constants.APPLICATION_DELETED,
		TimeStamp: milliseconds,
		Application: &subscription.Application{
			Uuid:         applicationUUID,
			Name:         applicationSpec.Name,
			Owner:        applicationSpec.Owner,
			Organization: applicationSpec.Organization,
			Attributes:   applicationSpec.Attributes,
		},
	}
	sendEvent(event)
}

// SendAppUpdateEvent sends an application update event to the enforcer
func SendAppUpdateEvent(applicationUUID string, oldApplicationSpec cpv1alpha2.ApplicationSpec, newApplicationSpec cpv1alpha2.ApplicationSpec) {
	currentTime := time.Now()
	milliseconds := currentTime.UnixNano() / int64(time.Millisecond)
	event := subscription.Event{
		Uuid:      applicationUUID,
		Type:      constants.APPLICATION_UPDATED,
		TimeStamp: milliseconds,
		Application: &subscription.Application{
			Uuid:         applicationUUID,
			Name:         newApplicationSpec.Name,
			Owner:        newApplicationSpec.Owner,
			Organization: newApplicationSpec.Organization,
			Attributes:   newApplicationSpec.Attributes,
		},
	}
	sendEvent(event)
	SendDeleteApplicationKeyMappingEvent(applicationUUID, oldApplicationSpec)
	SendApplicationKeyMappingEvent(applicationUUID, newApplicationSpec)
}

// SendAddApplicationEvent sends an application creation event to the enforcer
func SendAddApplicationEvent(application cpv1alpha2.Application) {
	currentTime := time.Now()
	milliseconds := currentTime.UnixNano() / int64(time.Millisecond)
	event := subscription.Event{
		Uuid:      application.ObjectMeta.Name,
		Type:      constants.APPLICATION_CREATED,
		TimeStamp: milliseconds,
		Application: &subscription.Application{
			Uuid:         application.ObjectMeta.Name,
			Name:         application.Spec.Name,
			Owner:        application.Spec.Owner,
			Organization: application.Spec.Organization,
			Attributes:   application.Spec.Attributes,
		},
	}
	sendEvent(event)
}

// SendAddSubscriptionEvent sends an subscription creation event to the enforcer
func SendAddSubscriptionEvent(sub cpv1alpha2.Subscription) {
	currentTime := time.Now()
	milliseconds := currentTime.UnixNano() / int64(time.Millisecond)
	event := subscription.Event{
		Uuid:      sub.ObjectMeta.Name,
		Type:      constants.SUBSCRIPTION_CREATED,
		TimeStamp: milliseconds,
		Subscription: &subscription.Subscription{
			Uuid:         sub.ObjectMeta.Name,
			SubStatus:    sub.Spec.SubscriptionStatus,
			Organization: sub.Spec.Organization,
			SubscribedApi: &subscription.SubscribedAPI{
				Name:    sub.Spec.API.Name,
				Version: sub.Spec.API.Version,
			},
		},
	}
	sendEvent(event)
}

// SendDeleteSubscriptionEvent sends an subscription deletion event to the enforcer
func SendDeleteSubscriptionEvent(subscriptionUUID string, subscriptionSpec cpv1alpha2.SubscriptionSpec) {
	currentTime := time.Now()
	milliseconds := currentTime.UnixNano() / int64(time.Millisecond)
	event := subscription.Event{
		Uuid:      subscriptionUUID,
		Type:      constants.SUBSCRIPTION_DELETED,
		TimeStamp: milliseconds,
		Subscription: &subscription.Subscription{
			Uuid:         subscriptionUUID,
			SubStatus:    subscriptionSpec.SubscriptionStatus,
			Organization: subscriptionSpec.Organization,
			SubscribedApi: &subscription.SubscribedAPI{
				Name:    subscriptionSpec.API.Name,
				Version: subscriptionSpec.API.Version,
			},
		},
	}
	sendEvent(event)
}

// SendCreateApplicationMappingEvent sends an application mapping event to the enforcer
func SendCreateApplicationMappingEvent(applicationMapping cpv1alpha2.ApplicationMapping) {
	currentTime := time.Now()
	milliseconds := currentTime.UnixNano() / int64(time.Millisecond)
	event := subscription.Event{
		Uuid:      applicationMapping.ObjectMeta.Name,
		Type:      constants.APPLICATION_MAPPING_CREATED,
		TimeStamp: milliseconds,
		ApplicationMapping: &subscription.ApplicationMapping{
			Uuid:            applicationMapping.ObjectMeta.Name,
			ApplicationRef:  applicationMapping.Spec.ApplicationRef,
			SubscriptionRef: applicationMapping.Spec.SubscriptionRef,
		},
	}
	sendEvent(event)
}

// SendDeleteApplicationMappingEvent sends an application mapping deletion event to the enforcer
func SendDeleteApplicationMappingEvent(applicationMappingUUID string, applicationMappingSpec cpv1alpha2.ApplicationMappingSpec) {
	currentTime := time.Now()
	milliseconds := currentTime.UnixNano() / int64(time.Millisecond)
	event := subscription.Event{
		Uuid:      applicationMappingUUID,
		Type:      constants.APPLICATION_DELETED,
		TimeStamp: milliseconds,
		ApplicationMapping: &subscription.ApplicationMapping{
			Uuid:            applicationMappingUUID,
			ApplicationRef:  applicationMappingSpec.ApplicationRef,
			SubscriptionRef: applicationMappingSpec.SubscriptionRef,
		},
	}
	sendEvent(event)
}
func SendDeleteApplicationKeyMappingEvent(applicationUUID string, applicationKeyMapping cpv1alpha2.ApplicationSpec) {
	currentTime := time.Now()
	milliseconds := currentTime.UnixNano() / int64(time.Millisecond)
	var oauth2SecurityScheme = applicationKeyMapping.SecuritySchemes.OAuth2
	if oauth2SecurityScheme != nil {
		for _, env := range oauth2SecurityScheme.Environments {
			event := subscription.Event{
				Uuid:      applicationUUID,
				Type:      constants.APPLICATION_KEY_MAPPING_DELETED,
				TimeStamp: milliseconds,
				ApplicationKeyMapping: &subscription.ApplicationKeyMapping{
					ApplicationUUID:       applicationUUID,
					SecurityScheme:        constants.OAuth2,
					ApplicationIdentifier: env.AppID,
					KeyType:               env.KeyType,
					EnvID:                 env.EnvID,
				},
			}
			sendEvent(event)
		}
	}
}
func SendApplicationKeyMappingEvent(applicationUUID string, applicationSpec cpv1alpha2.ApplicationSpec) {
	currentTime := time.Now()
	milliseconds := currentTime.UnixNano() / int64(time.Millisecond)
	var oauth2SecurityScheme = applicationSpec.SecuritySchemes.OAuth2
	if oauth2SecurityScheme != nil {
		for _, env := range oauth2SecurityScheme.Environments {
			event := subscription.Event{
				Uuid:      applicationUUID,
				Type:      constants.APPLICATION_KEY_MAPPING_CREATED,
				TimeStamp: milliseconds,
				ApplicationKeyMapping: &subscription.ApplicationKeyMapping{
					ApplicationUUID:       applicationUUID,
					SecurityScheme:        constants.OAuth2,
					ApplicationIdentifier: env.AppID,
					KeyType:               env.KeyType,
					EnvID:                 env.EnvID,
				},
			}
			sendEvent(event)
		}
	}
}
func sendEvent(event subscription.Event) {
	for clientId, stream := range GetAllClientConnections() {
		err := stream.Send(&event)
		if err != nil {
			loggers.LoggerAPK.Errorf("Error sending event to client %s: %v", clientId, err)
		} else {
			loggers.LoggerAPK.Debugf("Event sent to client %s", clientId)
		}
	}
}
