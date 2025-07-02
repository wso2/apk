package server

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wso2/apk/common-controller/internal/cache"
	"github.com/wso2/apk/common-controller/internal/config"
	dpv2alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v2alpha1"
	"github.com/wso2/apk/common-go-libs/pkg/server/model"
)

var applicationMap = make(map[string]model.Application)
var subscriptionMap = make(map[string]model.Subscription)
var applicationMappingMap = make(map[string]model.ApplicationMapping)
var applicationKeyMappingMap = make(map[string]model.ApplicationKeyMapping)

// StartInternalServer starts the internal server
func StartInternalServer() {

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.GET("/applications", func(c *gin.Context) {
		applicationList := []model.Application{}
		for _, application := range applicationMap {
			applicationList = append(applicationList, application)
		}
		c.JSON(http.StatusOK, model.ApplicationList{List: applicationList})
	})
	r.GET("/subscriptions", func(c *gin.Context) {
		subscriptionList := []model.Subscription{}
		for _, subscription := range subscriptionMap {
			subscriptionList = append(subscriptionList, subscription)
		}
		c.JSON(http.StatusOK, model.SubscriptionList{List: subscriptionList})
	})
	r.GET("/applicationmappings", func(c *gin.Context) {
		applicationMappingList := []model.ApplicationMapping{}
		for _, applicationMapping := range applicationMappingMap {
			applicationMappingList = append(applicationMappingList, applicationMapping)
		}
		c.JSON(http.StatusOK, model.ApplicationMappingList{List: applicationMappingList})
	})
	r.GET("/applicationkeymappings", func(c *gin.Context) {
		applicationKeyMappingList := []model.ApplicationKeyMapping{}
		for _, applicationKeyMapping := range applicationKeyMappingMap {
			applicationKeyMappingList = append(applicationKeyMappingList, applicationKeyMapping)
		}
		c.JSON(http.StatusOK, model.ApplicationKeyMappingList{List: applicationKeyMappingList})
	})
	r.GET("/routepolicies", func(c *gin.Context) {
		log.Printf("Received request from %s", c.ClientIP())
		rps := cache.GetRoutePolicyDataStore().GetRoutePolicies()
		rpList := dpv2alpha1.RoutePolicyList{
			Items: []dpv2alpha1.RoutePolicy{},
		}
		for _, rp := range rps {
			rpList.Items = append(rpList.Items, rp)
		}
		log.Printf("Returning %d route policies", len(rpList.Items))
		c.JSON(http.StatusOK, rpList)
	})
	r.GET("/routemetadata", func(c *gin.Context) {
		log.Printf("Received request from %s", c.ClientIP())
		rmds := cache.GetRouteMetadataDataStore().GetRouteMetadatas()
		rmdList := dpv2alpha1.RouteMetadataList{
			Items: []dpv2alpha1.RouteMetadata{},
		}
		for _, rmd := range rmds {
			rmdList.Items = append(rmdList.Items, rmd)
		}
		log.Printf("Returning %d route metadata", len(rmdList.Items))
		c.JSON(http.StatusOK, rmdList)
	})
	conf := config.ReadConfigs()
	certPath := conf.CommonController.Keystore.CertPath
	keyPath := conf.CommonController.Keystore.KeyPath
	port := conf.CommonController.InternalAPIServer.Port
	r.RunTLS(fmt.Sprintf(":%d", port), certPath, keyPath)
}

// AddApplication adds an application to the application list
func AddApplication(application model.Application) {
	applicationMap[application.UUID] = application
}

// DeleteAllApplications deletes all applications from the application list
func DeleteAllApplications() {
	applicationMap = make(map[string]model.Application)
}

// DeleteAllSubscriptions deletes all subscriptions from the subscription list
func DeleteAllSubscriptions() {
	subscriptionMap = make(map[string]model.Subscription)
}

// DeleteAllApplicationMappings deletes all application mappings from the application mapping list
func DeleteAllApplicationMappings() {
	applicationMappingMap = make(map[string]model.ApplicationMapping)
}

// DeleteAllApplicationKeyMappings deletes all application key mappings from the application key mapping list
func DeleteAllApplicationKeyMappings() {
	applicationKeyMappingMap = make(map[string]model.ApplicationKeyMapping)
}

// AddSubscription adds a subscription to the subscription list
func AddSubscription(subscription model.Subscription) {
	subscriptionMap[subscription.UUID] = subscription
}

// AddApplicationMapping adds an application mapping to the application mapping list
func AddApplicationMapping(applicationMapping model.ApplicationMapping) {
	applicationMappingMap[applicationMapping.UUID] = applicationMapping
}

// AddApplicationKeyMapping adds an application key mapping to the application key mapping list
func AddApplicationKeyMapping(applicationKeyMapping model.ApplicationKeyMapping) {
	applicationMappingKey := strings.Join([]string{applicationKeyMapping.ApplicationUUID, applicationKeyMapping.EnvID, applicationKeyMapping.SecurityScheme, applicationKeyMapping.KeyType}, ":")
	applicationKeyMappingMap[applicationMappingKey] = applicationKeyMapping
}

// DeleteApplicationKeyMapping deletes an application key mapping from the application key mapping list
func DeleteApplicationKeyMapping(applicationKeyMapping model.ApplicationKeyMapping) {
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
func GetApplicationMappingFromStore(applicationMappingUUID string) model.ApplicationMapping {
	return applicationMappingMap[applicationMappingUUID]
}
