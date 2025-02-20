/*
 *  Copyright (c) 2025, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */
package jwtbackend

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/wso2/apk/common-go-libs/loggers"
	"github.com/wso2/apk/gateway/enforcer/internal/config"
	"github.com/wso2/apk/gateway/enforcer/internal/util"
)

// StartJWKSServer starts a server that serves JWKS.
func StartJWKSServer(cfg *config.Server) {
	r := gin.Default()

	r.GET("/jwks", func(c *gin.Context) {
		jwks, err := readAndConvertToJwks(cfg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWKS"})
			return
		}
		c.JSON(http.StatusOK, jwks)
	})
	r.RunTLS(":9092", cfg.EnforcerPublicKeyPath, cfg.EnforcerPrivateKeyPath)
}
func readAndConvertToJwks(cfg *config.Server) (jwk.Set, error) {
	// Decode the PEM data
	publicCert, _ := util.LoadCertificate(cfg.JWTGeneratorPublicKeyPath)

	// Extract the public key from the certificate
	pubKey := publicCert.PublicKey

	// Convert the public key to a JWK
	jwkKey, err := jwk.FromRaw(pubKey)
	if err != nil {
		loggers.LoggerAPKOperator.Errorf("failed to create JWK: %s", err)
	}
	jwks := jwk.NewSet()
	jwks.AddKey(jwkKey)

	return jwks, nil
}
