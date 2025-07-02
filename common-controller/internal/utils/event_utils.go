package utils

import (
	time "time"

	"github.com/google/uuid"
	"github.com/wso2/apk/common-controller/internal/loggers"
	cpv1alpha2 "github.com/wso2/apk/common-go-libs/apis/cp/v1alpha2"
	cpv1alpha3 "github.com/wso2/apk/common-go-libs/apis/cp/v1alpha3"
	"github.com/wso2/apk/common-go-libs/constants"
	"github.com/wso2/apk/common-go-libs/pkg/discovery/api/wso2/discovery/subscription"
)

// SendAppDeletionEvent sends an application creation event to the enforcer
func SendAppDeletionEvent(applicationUUID string, applicationSpec cpv1alpha2.ApplicationSpec) {
	SendApplicationEvent(constants.ApplicationDeleted, applicationUUID, applicationSpec.Name, applicationSpec.Owner,
		applicationSpec.Organization, applicationSpec.Attributes)
}

// SendApplicationEvent sends an application deletion event to the enforcer
func SendApplicationEvent(eventType, applicationUUID, applicationName, applicationOwner, organization string, appAttribute map[string]string) {
	currentTime := time.Now()
	milliseconds := currentTime.UnixNano() / int64(time.Millisecond)
	event := subscription.Event{
		Uuid:      applicationUUID,
		Type:      eventType,
		TimeStamp: milliseconds,
		Application: &subscription.Application{
			Uuid:         applicationUUID,
			Name:         applicationName,
			Owner:        applicationOwner,
			Attributes:   appAttribute,
			Organization: organization,
		},
	}
	loggers.LoggerAPKOperator.Debugf("Sending event to all clients: %v", &event)
	sendEvent(&event)
}

// SendAppUpdateEvent sends an application update event to the enforcer
func SendAppUpdateEvent(applicationUUID string, oldApplicationSpec cpv1alpha2.ApplicationSpec, newApplicationSpec cpv1alpha2.ApplicationSpec) {
	SendApplicationEvent(constants.ApplicationUpdated, applicationUUID, oldApplicationSpec.Name, oldApplicationSpec.Owner,
		oldApplicationSpec.Organization, oldApplicationSpec.Attributes)
	if oldApplicationSpec.SecuritySchemes != nil {
		sendDeleteApplicationKeyMappingEvent(applicationUUID, oldApplicationSpec)
	}
	if newApplicationSpec.SecuritySchemes != nil {
		sendApplicationKeyMappingEvent(applicationUUID, newApplicationSpec)
	}
}

// SendAddApplicationEvent sends an application creation event to the enforcer
func SendAddApplicationEvent(application cpv1alpha2.Application) {
	SendApplicationEvent(constants.ApplicationCreated, application.ObjectMeta.Name, application.Spec.Name, application.Spec.Owner,
		application.Spec.Organization, application.Spec.Attributes)
	if application.Spec.SecuritySchemes != nil {
		sendApplicationKeyMappingEvent(application.ObjectMeta.Name, application.Spec)
	}
}

// SendAddSubscriptionEvent sends an subscription creation event to the enforcer
func SendAddSubscriptionEvent(sub cpv1alpha3.Subscription) {
	SendSubscriptionEvent(constants.SubscriptionCreated, sub.ObjectMeta.Name, sub.Spec.SubscriptionStatus,
		sub.Spec.Organization, sub.Spec.API.Name, sub.Spec.API.Version, sub.Spec.RatelimitRef.Name)
}

// SendSubscriptionEvent sends an subscription creation event to the enforcer
func SendSubscriptionEvent(eventType, subscriptionID, subscriptionStatus, organization, apiName, apiVersion string, ratelimit string) {
	currentTime := time.Now()
	milliseconds := currentTime.UnixNano() / int64(time.Millisecond)
	event := subscription.Event{
		Uuid:      uuid.New().String(),
		Type:      eventType,
		TimeStamp: milliseconds,
		Subscription: &subscription.Subscription{
			Uuid:         subscriptionID,
			SubStatus:    subscriptionStatus,
			Organization: organization,
			SubscribedApi: &subscription.SubscribedAPI{
				Name:    apiName,
				Version: apiVersion,
			},
			RatelimitTier: ratelimit,
		},
	}
	sendEvent(&event)
}

// SendDeleteSubscriptionEvent sends an subscription deletion event to the enforcer
func SendDeleteSubscriptionEvent(subscriptionUUID string, sub cpv1alpha3.Subscription) {
	SendSubscriptionEvent(constants.SubscriptionDeleted, subscriptionUUID, sub.Spec.SubscriptionStatus,
		sub.Spec.Organization, sub.Spec.API.Name, sub.Spec.API.Version, sub.Spec.RatelimitRef.Name)
}

// SendCreateApplicationMappingEvent sends an application mapping event to the enforcer
func SendCreateApplicationMappingEvent(applicationMapping cpv1alpha2.ApplicationMapping, application cpv1alpha2.Application, subscriptionCr cpv1alpha3.Subscription) {
	SendApplicationMappingEvent(constants.ApplicationMappingCreated, applicationMapping.ObjectMeta.Name, applicationMapping.Spec.ApplicationRef,
		applicationMapping.Spec.SubscriptionRef, application.Spec.Organization)
}

// SendApplicationMappingEvent sends an application mapping event to the enforcer
func SendApplicationMappingEvent(eventType, id, applicationRef, subscriptionRef, organization string) {
	currentTime := time.Now()
	milliseconds := currentTime.UnixNano() / int64(time.Millisecond)
	event := subscription.Event{
		Uuid:      uuid.New().String(),
		Type:      eventType,
		TimeStamp: milliseconds,
		ApplicationMapping: &subscription.ApplicationMapping{
			Uuid:            id,
			ApplicationRef:  applicationRef,
			SubscriptionRef: subscriptionRef,
			Organization:    organization,
		},
	}
	sendEvent(&event)
}

// SendDeleteApplicationMappingEvent sends an application mapping deletion event to the enforcer
func SendDeleteApplicationMappingEvent(applicationMappingUUID string,
	applicationMappingSpec cpv1alpha2.ApplicationMappingSpec, organization string) {
	SendApplicationMappingEvent(constants.ApplicationMappingDeleted, applicationMappingUUID,
		applicationMappingSpec.ApplicationRef, applicationMappingSpec.SubscriptionRef, organization)
}

func sendDeleteApplicationKeyMappingEvent(applicationUUID string, applicationKeyMapping cpv1alpha2.ApplicationSpec) {
	var oauth2SecurityScheme = applicationKeyMapping.SecuritySchemes.OAuth2
	if oauth2SecurityScheme != nil {
		for _, env := range oauth2SecurityScheme.Environments {
			SendApplicationKeyMappingEvent(constants.ApplicationKeyMappingDeleted, applicationUUID, constants.OAuth2,
				env.AppID, env.KeyType, env.EnvID, applicationKeyMapping.Organization)
		}
	}
}

func sendApplicationKeyMappingEvent(applicationUUID string, applicationSpec cpv1alpha2.ApplicationSpec) {
	var oauth2SecurityScheme = applicationSpec.SecuritySchemes.OAuth2
	if oauth2SecurityScheme != nil {
		for _, env := range oauth2SecurityScheme.Environments {
			SendApplicationKeyMappingEvent(constants.ApplicationKeyMappingCreated, applicationUUID, constants.OAuth2,
				env.AppID, env.KeyType, env.EnvID, applicationSpec.Organization)
		}
	}
}

// SendApplicationKeyMappingEvent sends an application key mapping event to the enforcer
func SendApplicationKeyMappingEvent(eventType, applicationUUID, securityScheme, applicationIdentifier, keyType, envID,
	organization string) {
	currentTime := time.Now()
	milliseconds := currentTime.UnixNano() / int64(time.Millisecond)
	event := subscription.Event{
		Uuid:      uuid.New().String(),
		Type:      eventType,
		TimeStamp: milliseconds,
		ApplicationKeyMapping: &subscription.ApplicationKeyMapping{
			ApplicationUUID:       applicationUUID,
			SecurityScheme:        securityScheme,
			ApplicationIdentifier: applicationIdentifier,
			KeyType:               keyType,
			EnvID:                 envID,
			Organization:          organization,
		},
	}
	sendEvent(&event)
}

// SendRoutePolicyCreatedOrUpdatedEvent sends a route policy creation or update event to the enforcer
func SendRoutePolicyCreatedOrUpdatedEvent(routePolicy string) {
	currentTime := time.Now()
	milliseconds := currentTime.UnixNano() / int64(time.Millisecond)
	event := subscription.Event{
		Uuid:      uuid.New().String(),
		Type:      constants.RoutePolicyCreatedOrUpdated,
		TimeStamp: milliseconds,
		RoutePolicy: routePolicy,
	}
	sendEvent(&event)
}

// SendRoutePolicyDeletedEvent sends a route policy deletion event to the enforcer
func SendRoutePolicyDeletedEvent(routePolicy string) {
	currentTime := time.Now()
	milliseconds := currentTime.UnixNano() / int64(time.Millisecond)
	event := subscription.Event{
		Uuid:      uuid.New().String(),
		Type:      constants.RoutePolicyDeleted,
		TimeStamp: milliseconds,
		RoutePolicy: routePolicy,
	}
	sendEvent(&event)
}

// SendRouteMetadataCreatedOrUpdatedEvent sends a route metadata creation or update event to the enforcer
func SendRouteMetadataCreatedOrUpdatedEvent(routeMetadata string) {
	currentTime := time.Now()
	milliseconds := currentTime.UnixNano() / int64(time.Millisecond)
	event := subscription.Event{
		Uuid:      uuid.New().String(),
		Type:      constants.RouteMetadataCreatedOrUpdated,
		TimeStamp: milliseconds,
		RouteMetadata: routeMetadata,
	}
	sendEvent(&event)
}

// SendRouteMetadataDeletedEvent sends a route metadata deletion event to the enforcer
func SendRouteMetadataDeletedEvent(routeMetadata string) {
	currentTime := time.Now()
	milliseconds := currentTime.UnixNano() / int64(time.Millisecond)
	event := subscription.Event{
		Uuid:      uuid.New().String(),
		Type:      constants.RouteMetadataDeleted,
		TimeStamp: milliseconds,
		RouteMetadata: routeMetadata,
	}
	sendEvent(&event)
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

// SendResetEvent sends initial event to the enforcer
func SendResetEvent() {
	currentTime := time.Now()
	milliseconds := currentTime.UnixNano() / int64(time.Millisecond)
	event := &subscription.Event{
		Uuid:      uuid.New().String(),
		Type:      constants.AllEvents,
		TimeStamp: milliseconds,
	}
	for clientID, stream := range GetAllClientConnections() {
		err := stream.Send(event)
		if err != nil {
			loggers.LoggerAPKOperator.Errorf("Error sending reset event to client %s: %v", clientID, err)
		} else {
			loggers.LoggerAPKOperator.Debugf("Reset event sent to client %s", clientID)
		}
	}
}
