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
	"crypto/sha256"
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/wso2/apk/gateway/enforcer/internal/config"
	"github.com/wso2/apk/gateway/enforcer/internal/util"
)

// JWKKEy represents the JWK key.
var JWKKEy jwk.Key

// StartJWKSServer starts a server that serves JWKS.
func StartJWKSServer(cfg *config.Server) {
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	r.GET("/jwks", func(c *gin.Context) {
		jwks := createJWKSSet()
		c.JSON(http.StatusOK, jwks)
	})
	r.RunTLS(":9092", cfg.EnforcerPublicKeyPath, cfg.EnforcerPrivateKeyPath)
}

// ReadAndConvertToJwks reads the public key from the given path and converts it to a JWK.
func ReadAndConvertToJwks(cfg *config.Server) (jwk.Key, error) {
	// Decode the PEM data
	publicCert, _ := util.LoadCertificate(cfg.JWTGeneratorPublicKeyPath)

	// Extract the public key from the certificate
	pubKey := publicCert.PublicKey

	// Convert the public key to a JWK
	jwkKey, err := jwk.FromRaw(pubKey)
	cfg = config.GetConfig()
	if err != nil {
		cfg.Logger.Sugar().Errorf("Failed to convert public key to JWK: %v", err)
		return nil, err
	}
	hash := sha256.Sum256(publicCert.RawSubjectPublicKeyInfo)
	kid := base64.RawURLEncoding.EncodeToString(hash[:])
	jwkKey.Set(jwk.KeyIDKey, kid)
	return jwkKey, nil
}
func createJWKSSet() jwk.Set {
	jwks := jwk.NewSet()
	jwks.AddKey(JWKKEy)
	return jwks
}
