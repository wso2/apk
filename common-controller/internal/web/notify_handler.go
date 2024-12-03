/*
 *  Copyright (c) 2023, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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

package web

import (
	"fmt"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"net/http"
	"time"
	"crypto/tls"
	"os"
	"crypto/x509"
	loggers "github.com/wso2/apk/common-controller/internal/loggers"
	"github.com/wso2/apk/adapter/pkg/logging"
	config "github.com/wso2/apk/common-controller/internal/config"
	"io/ioutil"
	"strings"
	"encoding/base64"
	"encoding/json"	
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
)

const tokenRevocationType = "TOKEN_REVOCATION"

type revokeRequest struct {
	Token string `json:"token"`
	Jti string `json:"jti"`
	Expiry int64 `json:"expiry"`
}

type jWTClaims struct {
	Jti string `json:"jti"`
	Exp int64  `json:"exp"`
}

var (
	redisAddr       string
	redisUsername   string
	redisPassword   string
	redisUserCertPath   string
	redisUserKeyPath    string
	redisCACertPath string
	isTLSEnabled    bool
	redisRevokedTokenChannel string
	tokenExpiryDivider = "_##_"
	authKeyPath string
	authKeyHeader string
	rdb *redis.Client
)

func init() {
	conf := config.ReadConfigs()
	redisHost := conf.CommonController.Redis.Host
	redisPort := conf.CommonController.Redis.Port
	redisAddr = redisHost + ":" + redisPort
	redisUsername = conf.CommonController.Redis.Username
	redisPassword = conf.CommonController.Redis.Password
	redisUserCertPath = conf.CommonController.Redis.UserCertPath
	redisUserKeyPath = conf.CommonController.Redis.UserKeyPath
	redisCACertPath = conf.CommonController.Redis.CACertPath
	isTLSEnabled = conf.CommonController.Redis.TLSEnabled
	redisRevokedTokenChannel = conf.CommonController.Redis.RevokedTokenChannel
	authKeyPath = conf.CommonController.Sts.AuthKeyPath
	authKeyHeader = conf.CommonController.Sts.AuthKeyHeader
	utilruntime.Must(initRedisClient())
}

// initRedisClient initializes the redis connection
func initRedisClient() error {
	if isTLSEnabled {
		cert, err := tls.LoadX509KeyPair(redisUserCertPath, redisUserKeyPath)
		if err != nil {
			return err;
		}
		caCert, err := os.ReadFile(redisCACertPath)
		if err != nil {
			return err;
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		options := &redis.Options{
			Addr:     redisAddr,
			Password: redisPassword,
			TLSConfig: &tls.Config{
				MinVersion:   tls.VersionTLS12,
				Certificates: []tls.Certificate{cert},
				RootCAs:      caCertPool,
				InsecureSkipVerify: true,
			},
		}
		if redisUsername != "" {
			options.Username = redisUsername
		}
		rdb = redis.NewClient(options)
	} else {
		options := &redis.Options{
			Addr:     redisAddr,
			Password: redisPassword,
		}
		// Only set Username if it's not empty
		if redisUsername != "" {
			options.Username = redisUsername
		}
		rdb = redis.NewClient(options)
	}
	return nil;
}

// NotifyHandler handles notify requests
func NotifyHandler(c *gin.Context) {
	_type := c.Query("type")
	if _type == tokenRevocationType {
		revokeToken(c)
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid type"})
	return
}

func revokeToken(c *gin.Context) {
	if !authenticateTokenRevocationRequest(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized request"})
		return
	}

	var request revokeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		loggers.LoggerAPI.ErrorC(logging.PrintError(logging.Error3200, logging.MAJOR, "Error while parsing body: %v", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while parsing json payload"})
		return
	}
  var jti string;
	var expiry int64;
	if request.Token != "" {
		claims, err := extractClaimsFromJWT(request.Token)
		if err != nil {
			loggers.LoggerAPI.ErrorC(logging.PrintError(logging.Error3203, logging.MAJOR, "Error decoding token: %v", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error decoding token"})
			return
		}
		jti = claims.Jti
		expiry = claims.Exp
	} else {
		jti = request.Jti
		expiry = request.Expiry
	}
	if expiry <= time.Now().Unix() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token is already expired"})
		return
	}
	err := storeTokenInRedis(jti, expiry)
	if (err != nil) {
		loggers.LoggerAPI.ErrorC(logging.PrintError(logging.Error3202, logging.MAJOR, "Error adding revoked tokens to redis: %v", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to store the token in Redis cache"})
		return
	} 
	c.JSON(http.StatusOK, gin.H{"message": "Token revoked successfully"})
}

func generateKey(jti string) string {
	return fmt.Sprintf("%s:%s", "wso2:apk:revoked_token", jti)
}

func storeTokenInRedis(token string, expiry int64) error {
	key := generateKey(token)
	err := rdb.Do(context.Background(), "set", key, expiry, "EXAT", expiry).Err()
	if err != nil {
		loggers.LoggerAPI.Warnf("Error occured while trying to set key with expiry. Error: %+v. \n Trying to use SET and EXPIREAT command...", err)
		err = rdb.Do(context.Background(), "set", key, expiry).Err()
		if err != nil { 
			loggers.LoggerAPI.Errorf("Error occured while setting the key. Error %+v", err)
			return err
		} 
		err = rdb.Do(context.Background(), "expireat", key, expiry).Err()
		if err != nil {
			loggers.LoggerAPI.Errorf("Error occured while setting the expiry. Error %+v", err)
			return err
		}
	}
	publishValue := fmt.Sprintf("%s%s%d", token, tokenExpiryDivider, expiry)
	err = rdb.Do(context.Background(), "publish", redisRevokedTokenChannel, publishValue).Err()
	if err != nil {
		return err
	}
	return nil
}

func authenticateTokenRevocationRequest(c *gin.Context) bool {
	fileContent, err := ioutil.ReadFile(authKeyPath)
	if err != nil {
		loggers.LoggerAPI.ErrorC(logging.PrintError(logging.Error3204, logging.MAJOR, "Error reading shared key file: %v", err))
		return false
	}
	headerValue := c.GetHeader(authKeyHeader)
	if headerValue == "" {
		loggers.LoggerAPI.ErrorC(logging.PrintError(logging.Error3205, logging.MAJOR, "Unauthorized: Missing header"))
		return false
	}
	if !strings.EqualFold(string(fileContent), headerValue) {
		loggers.LoggerAPI.ErrorC(logging.PrintError(logging.Error3206, logging.MAJOR, "Unauthorized: Invalid header value"))
		return false
	}
	return true
}

func extractClaimsFromJWT(jwtToken string) (jWTClaims, error) {
	var claims jWTClaims
	parts := strings.Split(jwtToken, ".")
	if len(parts) != 3 {
		loggers.LoggerAPI.ErrorC(logging.PrintError(logging.Error3205, logging.MAJOR, "Invalid JWT"))
		return claims, fmt.Errorf("Invalid JWT token")
	}
	payloadStr := strings.TrimRight(parts[1], "=")
	payload, err := base64.RawURLEncoding.DecodeString(payloadStr)
	if err != nil {
		loggers.LoggerAPI.ErrorC(logging.PrintError(logging.Error3205, logging.MAJOR, "Invalid JWT"))
		return claims, fmt.Errorf("Error decoding payload: %v", err)
	}
	err = json.Unmarshal(payload, &claims)
	if err != nil {
		return claims, fmt.Errorf("Error parsing JSON: %v", err)
	}
	return claims, nil
}
