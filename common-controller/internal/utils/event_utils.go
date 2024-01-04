package utils

import (
	time "time"

	"github.com/google/uuid"
	"github.com/wso2/apk/common-controller/internal/loggers"
	constants "github.com/wso2/apk/common-controller/internal/operator/constant"
	cpv1alpha2 "github.com/wso2/apk/common-go-libs/apis/cp/v1alpha2"
	apkmgt "github.com/wso2/apk/common-go-libs/pkg/discovery/api/wso2/discovery/service/apkmgt"
	"github.com/wso2/apk/common-go-libs/pkg/discovery/api/wso2/discovery/subscription"
)

// SendAppDeletionEvent sends an application creation event to the enforcer
func SendAppDeletionEvent(applicationUUID string, applicationSpec cpv1alpha2.ApplicationSpec) {
	currentTime := time.Now()
	milliseconds := currentTime.UnixNano() / int64(time.Millisecond)
	event := subscription.Event{
		Uuid:      uuid.New().String(),
		Type:      constants.ApplicationDeleted,
		TimeStamp: milliseconds,
		Application: &subscription.Application{
			Uuid:         applicationUUID,
			Name:         applicationSpec.Name,
			Owner:        applicationSpec.Owner,
			Organization: applicationSpec.Organization,
			Attributes:   applicationSpec.Attributes,
		},
	}
	sendEvent(&event)
}

// SendAppUpdateEvent sends an application update event to the enforcer
func SendAppUpdateEvent(applicationUUID string, oldApplicationSpec cpv1alpha2.ApplicationSpec, newApplicationSpec cpv1alpha2.ApplicationSpec) {
	currentTime := time.Now()
	milliseconds := currentTime.UnixNano() / int64(time.Millisecond)
	event := subscription.Event{
		Uuid:      uuid.New().String(),
		Type:      constants.ApplicationUpdated,
		TimeStamp: milliseconds,
		Application: &subscription.Application{
			Uuid:         applicationUUID,
			Name:         newApplicationSpec.Name,
			Owner:        newApplicationSpec.Owner,
			Organization: newApplicationSpec.Organization,
			Attributes:   newApplicationSpec.Attributes,
		},
	}
	loggers.LoggerAPKOperator.Debugf("Sending event to all clients: %v", &event)
	sendEvent(&event)
	sendDeleteApplicationKeyMappingEvent(applicationUUID, oldApplicationSpec)
	sendApplicationKeyMappingEvent(applicationUUID, newApplicationSpec)
}

// SendAddApplicationEvent sends an application creation event to the enforcer
func SendAddApplicationEvent(application cpv1alpha2.Application) {
	currentTime := time.Now()
	milliseconds := currentTime.UnixNano() / int64(time.Millisecond)
	event := subscription.Event{
		Uuid:      uuid.New().String(),
		Type:      constants.ApplicationCreated,
		TimeStamp: milliseconds,
		Application: &subscription.Application{
			Uuid:         application.ObjectMeta.Name,
			Name:         application.Spec.Name,
			Owner:        application.Spec.Owner,
			Organization: application.Spec.Organization,
			Attributes:   application.Spec.Attributes,
		},
	}
	sendEvent(&event)
	sendApplicationKeyMappingEvent(application.ObjectMeta.Name, application.Spec)
}

// SendAddSubscriptionEvent sends an subscription creation event to the enforcer
func SendAddSubscriptionEvent(sub cpv1alpha2.Subscription) {
	currentTime := time.Now()
	milliseconds := currentTime.UnixNano() / int64(time.Millisecond)
	event := subscription.Event{
		Uuid:      uuid.New().String(),
		Type:      constants.SubscriptionCreated,
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
	sendEvent(&event)
}

// SendDeleteSubscriptionEvent sends an subscription deletion event to the enforcer
func SendDeleteSubscriptionEvent(subscriptionUUID string, subscriptionSpec cpv1alpha2.SubscriptionSpec) {
	currentTime := time.Now()
	milliseconds := currentTime.UnixNano() / int64(time.Millisecond)
	event := subscription.Event{
		Uuid:      subscriptionUUID,
		Type:      constants.SubscriptionDeleted,
		TimeStamp: milliseconds,
		Subscription: &subscription.Subscription{
			Uuid:         uuid.New().String(),
			SubStatus:    subscriptionSpec.SubscriptionStatus,
			Organization: subscriptionSpec.Organization,
			SubscribedApi: &subscription.SubscribedAPI{
				Name:    subscriptionSpec.API.Name,
				Version: subscriptionSpec.API.Version,
			},
		},
	}
	sendEvent(&event)
}

// SendCreateApplicationMappingEvent sends an application mapping event to the enforcer
func SendCreateApplicationMappingEvent(applicationMapping cpv1alpha2.ApplicationMapping, application cpv1alpha2.Application, subscriptionCr cpv1alpha2.Subscription) {
	currentTime := time.Now()
	milliseconds := currentTime.UnixNano() / int64(time.Millisecond)
	event := subscription.Event{
		Uuid:      uuid.New().String(),
		Type:      constants.ApplicationMappingCreated,
		TimeStamp: milliseconds,
		ApplicationMapping: &subscription.ApplicationMapping{
			Uuid:            applicationMapping.ObjectMeta.Name,
			ApplicationRef:  applicationMapping.Spec.ApplicationRef,
			SubscriptionRef: applicationMapping.Spec.SubscriptionRef,
			Organization:    application.Spec.Organization,
		},
	}
	sendEvent(&event)
}

// SendDeleteApplicationMappingEvent sends an application mapping deletion event to the enforcer
func SendDeleteApplicationMappingEvent(applicationMappingUUID string, applicationMappingSpec cpv1alpha2.ApplicationMappingSpec) {
	currentTime := time.Now()
	milliseconds := currentTime.UnixNano() / int64(time.Millisecond)
	event := subscription.Event{
		Uuid:      uuid.New().String(),
		Type:      constants.ApplicationMappingDeleted,
		TimeStamp: milliseconds,
		ApplicationMapping: &subscription.ApplicationMapping{
			Uuid:            applicationMappingUUID,
			ApplicationRef:  applicationMappingSpec.ApplicationRef,
			SubscriptionRef: applicationMappingSpec.SubscriptionRef,
		},
	}
	sendEvent(&event)
}
func sendDeleteApplicationKeyMappingEvent(applicationUUID string, applicationKeyMapping cpv1alpha2.ApplicationSpec) {
	currentTime := time.Now()
	milliseconds := currentTime.UnixNano() / int64(time.Millisecond)
	var oauth2SecurityScheme = applicationKeyMapping.SecuritySchemes.OAuth2
	if oauth2SecurityScheme != nil {
		for _, env := range oauth2SecurityScheme.Environments {
			event := subscription.Event{
				Uuid:      uuid.New().String(),
				Type:      constants.ApplicationKeyMappingDeleted,
				TimeStamp: milliseconds,
				ApplicationKeyMapping: &subscription.ApplicationKeyMapping{
					ApplicationUUID:       applicationUUID,
					SecurityScheme:        constants.OAuth2,
					ApplicationIdentifier: env.AppID,
					KeyType:               env.KeyType,
					EnvID:                 env.EnvID,
					Organization:          applicationKeyMapping.Organization,
				},
			}
			sendEvent(&event)
		}
	}
}
func sendApplicationKeyMappingEvent(applicationUUID string, applicationSpec cpv1alpha2.ApplicationSpec) {
	currentTime := time.Now()
	milliseconds := currentTime.UnixNano() / int64(time.Millisecond)
	var oauth2SecurityScheme = applicationSpec.SecuritySchemes.OAuth2
	if oauth2SecurityScheme != nil {
		for _, env := range oauth2SecurityScheme.Environments {
			event := subscription.Event{
				Uuid:      uuid.New().String(),
				Type:      constants.ApplicationKeyMappingCreated,
				TimeStamp: milliseconds,
				ApplicationKeyMapping: &subscription.ApplicationKeyMapping{
					ApplicationUUID:       applicationUUID,
					SecurityScheme:        constants.OAuth2,
					ApplicationIdentifier: env.AppID,
					KeyType:               env.KeyType,
					EnvID:                 env.EnvID,
				},
			}
			sendEvent(&event)
		}
	}
}
func sendEvent(event *subscription.Event) {
	loggers.LoggerAPKOperator.Debugf("Sending event to all clients: %v", event)
	for clientID, stream := range GetAllClientConnections() {
		err := stream.Send(event)
		if err != nil {
			loggers.LoggerAPKOperator.Errorf("Error sending event to client %s: %v", clientID, err)
		} else {
			loggers.LoggerAPKOperator.Debugf("Event sent to client %s", clientID)
		}
	}
}

// SendInitialEvent sends initial event to the enforcer
func SendInitialEvent(srv apkmgt.EventStreamService_StreamEventsServer) {
	currentTime := time.Now()
	milliseconds := currentTime.UnixNano() / int64(time.Millisecond)

	event := subscription.Event{
		Uuid:      uuid.New().String(),
		Type:      constants.AllEvnts,
		TimeStamp: milliseconds,
	}
	loggers.LoggerAPKOperator.Debugf("Sending initial event to client: %v", &event)
	srv.Send(&event)
}
