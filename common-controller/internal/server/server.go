package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wso2/apk/common-controller/internal/config"
)

var applicationList *ApplicationList
var subscriptionList *SubscriptionList
var applicationMappingList *ApplicationMappingList
var applicationKeyMappingList *ApplicationKeyMappingList

// StartInternalServer starts the internal server
func StartInternalServer() {
	r := gin.Default()

	r.GET("/applications", func(c *gin.Context) {
		if applicationList == nil {
			c.JSON(http.StatusOK, ApplicationList{List: make([]Application, 0)})
		}
		c.JSON(http.StatusOK, applicationList)
	})
	r.GET("/subscriptions", func(c *gin.Context) {
		if subscriptionList == nil {
			c.JSON(http.StatusOK, SubscriptionList{List: make([]Subscription, 0)})
		}
		c.JSON(http.StatusOK, subscriptionList)
	})
	r.GET("/applicationmappings", func(c *gin.Context) {
		if applicationMappingList == nil {
			c.JSON(http.StatusOK, ApplicationMappingList{List: make([]ApplicationMapping, 0)})
		}
		c.JSON(http.StatusOK, applicationMappingList)
	})
	r.GET("/applicationkeymappings", func(c *gin.Context) {
		if applicationKeyMappingList == nil {
			c.JSON(http.StatusOK, ApplicationKeyMappingList{List: make([]ApplicationKeyMapping, 0)})
		}
		c.JSON(http.StatusOK, applicationKeyMappingList)
	})
	gin.SetMode(gin.ReleaseMode)
	conf := config.ReadConfigs()
	certPath := conf.CommonController.Keystore.CertPath
	keyPath := conf.CommonController.Keystore.KeyPath
	port := conf.CommonController.InternalAPIServer.Port
	r.RunTLS(fmt.Sprintf(":%d", port), certPath, keyPath)
}

// AddApplication adds an application to the application list
func AddApplication(appList *ApplicationList) {
	applicationList = appList
}

// AddSubscription adds a subscription to the subscription list
func AddSubscription(subList *SubscriptionList) {
	subscriptionList = subList
}

// AddApplicationMapping adds an application mapping to the application mapping list
func AddApplicationMapping(appMappingList ApplicationMappingList) {
	applicationMappingList = &appMappingList
}

// AddApplicationKeyMapping adds an application key mapping to the application key mapping list
func AddApplicationKeyMapping(appKeyMappingList ApplicationKeyMappingList) {
	applicationKeyMappingList = &appKeyMappingList
}
