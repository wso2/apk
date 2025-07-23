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
	// "crypto/sha256"
	// "encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
	// "github.com/wso2/apk/gateway/enforcer/internal/util"
)

// StartArtifactGeneratorServer sets up and starts the HTTP server for artifact generation APIs.
// It defines API routes under the /api/configurator base path.
func StartArtifactGeneratorServer(cfg *config.Server) {
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	api := r.Group("/api/configurator")
	{
		// Create API configuration file from api specification.
		api.POST("/apis/generate-configuration", func(c *gin.Context) {
			handlers.GetGeneratedAPKConf(c)
		})

		// Generate K8s Resources
		api.POST("/apis/generate-k8s-resources", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Generate K8s Resources API called"})
		})
	}

	// r.RunTLS(":9443", cfg.EnforcerPublicKeyPath, cfg.EnforcerPrivateKeyPath)

	// Start HTTP server
	if err := r.Run(":9444"); err != nil {
		panic("Failed to start API server: " + err.Error())
	}
}

func createJWKSSet() string {
	return "Response with JWKS set" // Placeholder response, replace with actual JWKS set
}
