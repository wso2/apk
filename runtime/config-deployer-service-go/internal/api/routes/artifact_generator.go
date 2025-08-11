/* Copyright (c) 2025 WSO2 LLC. (http://www.wso2.com) All Rights Reserved.
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package routes

import (
	"github.com/wso2/apk/config-deployer-service-go/internal/api/handlers"
	"github.com/wso2/apk/config-deployer-service-go/internal/config"
	"github.com/wso2/apk/config-deployer-service-go/internal/dto"
	"os"

	"net/http"

	"github.com/gin-gonic/gin"
	// "github.com/wso2/apk/gateway/enforcer/internal/util"
)

// StartArtifactGeneratorServer sets up and starts the HTTP server for artifact generation APIs.
func StartArtifactGeneratorServer(cfg *config.Server) {
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	r.GET("/health", func(c *gin.Context) {
		status := gin.H{
			"health": "Ok",
		}
		c.JSON(http.StatusOK, status)
	})

	api := r.Group("/api/configurator")
	{
		// Create API configuration file from api specification.
		api.POST("/apis/generate-configuration", func(c *gin.Context) {
			handlers.GetGeneratedAPKConf(c)
		})

		// Generate K8s Resources
		api.POST("/apis/generate-k8s-resources", func(c *gin.Context) {
			organization := c.Query("organization")
			if organization == "" {
				organization = "default"
			}
			cpInitiated := c.Query("cpInitiated")
			if cpInitiated == "" {
				cpInitiated = "false"
			}
			organizationObj := dto.NewOrganization("", organization, "default",
				"default", true)
			handlers.GetGeneratedK8sResources(c, organizationObj, cpInitiated)
		})
	}

	// Get certificate paths from environment or config
	certPath := os.Getenv("TLS_CERT_PATH")
	keyPath := os.Getenv("TLS_KEY_PATH")

	if certPath == "" {
		certPath = "/home/wso2kgw/security/config.pem"
	}
	if keyPath == "" {
		keyPath = "/home/wso2kgw/security/config.key"
	}

	// Start HTTPS server
	if err := r.RunTLS(":9443", certPath, keyPath); err != nil {
		panic("Failed to start HTTPS server: " + err.Error())
	}

	//Start HTTP server
	//if err := r.Run(":9444"); err != nil {
	//	panic("Failed to start API server: " + err.Error())
	//}
}
