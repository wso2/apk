package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wso2/apk/common-controller/internal/config"
)

var applicationMap = make(map[string]Application)
var subscriptionMap = make(map[string]Subscription)
var applicationMappingMap = make(map[string]ApplicationMapping)
var applicationKeyMappingMap = make(map[string]ApplicationKeyMapping)

// StartInternalServer starts the internal server
func StartInternalServer() {
	r := gin.Default()

	r.GET("/applications", func(c *gin.Context) {
		applicationList := []Application{}
		for _, application := range applicationMap {
			applicationList = append(applicationList, application)
		}
		c.JSON(http.StatusOK, ApplicationList{List: applicationList})
	})
	r.GET("/subscriptions", func(c *gin.Context) {
		subscriptionList := []Subscription{}
		for _, subscription := range subscriptionMap {
			subscriptionList = append(subscriptionList, subscription)
		}
		c.JSON(http.StatusOK, SubscriptionList{List: subscriptionList})
	})
	r.GET("/applicationmappings", func(c *gin.Context) {
		applicationMappingList := []ApplicationMapping{}
		for _, applicationMapping := range applicationMappingMap {
			applicationMappingList = append(applicationMappingList, applicationMapping)
		}
		c.JSON(http.StatusOK, ApplicationMappingList{List: applicationMappingList})
	})
	r.GET("/applicationkeymappings", func(c *gin.Context) {
		applicationKeyMappingList := []ApplicationKeyMapping{}
		for _, applicationKeyMapping := range applicationKeyMappingMap {
			applicationKeyMappingList = append(applicationKeyMappingList, applicationKeyMapping)
		}
		c.JSON(http.StatusOK, ApplicationKeyMappingList{List: applicationKeyMappingList})
	})
	gin.SetMode(gin.ReleaseMode)
	conf := config.ReadConfigs()
	certPath := conf.CommonController.Keystore.CertPath
	keyPath := conf.CommonController.Keystore.KeyPath
	port := conf.CommonController.InternalAPIServer.Port
	r.RunTLS(fmt.Sprintf(":%d", port), certPath, keyPath)
}

// AddApplication adds an application to the application list
func AddApplication(application Application) {
	applicationMap[application.UUID] = application
}

// DeleteAllApplications deletes all applications from the application list
func DeleteAllApplications() {
	applicationMap = make(map[string]Application)
}

// DeleteAllSubscriptions deletes all subscriptions from the subscription list
func DeleteAllSubscriptions() {
	subscriptionMap = make(map[string]Subscription)
}

// DeleteAllApplicationMappings deletes all application mappings from the application mapping list
func DeleteAllApplicationMappings() {
	applicationMappingMap = make(map[string]ApplicationMapping)
}

// DeleteAllApplicationKeyMappings deletes all application key mappings from the application key mapping list
func DeleteAllApplicationKeyMappings() {
	applicationKeyMappingMap = make(map[string]ApplicationKeyMapping)
}

// AddSubscription adds a subscription to the subscription list
func AddSubscription(subscription Subscription) {
	subscriptionMap[subscription.UUID] = subscription
}

// AddApplicationMapping adds an application mapping to the application mapping list
func AddApplicationMapping(applicationMapping ApplicationMapping) {
	applicationMappingMap[applicationMapping.UUID] = applicationMapping
}

// AddApplicationKeyMapping adds an application key mapping to the application key mapping list
func AddApplicationKeyMapping(applicationKeyMapping ApplicationKeyMapping) {
	applicationMappingKey := strings.Join([]string{applicationKeyMapping.ApplicationUUID, applicationKeyMapping.EnvID, applicationKeyMapping.SecurityScheme, applicationKeyMapping.KeyType}, ":")
	applicationKeyMappingMap[applicationMappingKey] = applicationKeyMapping
}

// DeleteApplicationKeyMapping deletes an application key mapping from the application key mapping list
func DeleteApplicationKeyMapping(applicationKeyMapping ApplicationKeyMapping) {
	applicationMappingKey := strings.Join([]string{applicationKeyMapping.ApplicationUUID, applicationKeyMapping.EnvID, applicationKeyMapping.SecurityScheme, applicationKeyMapping.KeyType}, ":")
	delete(applicationKeyMappingMap, applicationMappingKey)
}

// DeleteApplication deletes an application from the application list
func DeleteApplication(applicationUUID string) {
	delete(applicationMap, applicationUUID)
	for key := range applicationKeyMappingMap {
		if strings.HasPrefix(key, applicationUUID) {
			delete(applicationKeyMappingMap, key)
		}
	}
}

// DeleteSubscription deletes a subscription from the subscription list
func DeleteSubscription(subscriptionUUID string) {
	delete(subscriptionMap, subscriptionUUID)
}

// DeleteApplicationMapping deletes an application mapping from the application mapping list
func DeleteApplicationMapping(applicationMappingUUID string) {
	delete(applicationMappingMap, applicationMappingUUID)
}

// GetApplicationMappingFromStore returns an application mapping from the application mapping list
func GetApplicationMappingFromStore(applicationMappingUUID string) ApplicationMapping {
	return applicationMappingMap[applicationMappingUUID]
}
