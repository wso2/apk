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
 
 package config

import (
	"sync"

	egv1a1 "github.com/envoyproxy/gateway/api/v1alpha1"
	"github.com/kelseyhightower/envconfig"
	"github.com/wso2/apk/gateway/enforcer/internal/logging"
)

// Server holds the configuration parameters for the application.
type Server struct {
	TrustedAdapterCertsPath     string `envconfig:"TRUSTED_CA_CERTS_PATH" default:"/home/wso2/security/truststore"`
	TrustDefaultCerts           string `envconfig:"TRUST_DEFAULT_CERTS" default:"true"`
	EnforcerPrivateKeyPath      string `envconfig:"ENFORCER_PRIVATE_KEY_PATH" default:"/home/wso2/security/keystore/mg.key"`
	EnforcerPublicKeyPath       string `envconfig:"ENFORCER_PUBLIC_CERT_PATH" default:"/home/wso2/security/keystore/mg.pem"`
	OpaClientPrivateKeyPath     string `envconfig:"OPA_CLIENT_PRIVATE_KEY_PATH" default:"/home/wso2/security/keystore/mg.key"`
	OpaClientPublicKeyPath      string `envconfig:"OPA_CLIENT_PUBLIC_CERT_PATH" default:"/home/wso2/security/keystore/mg.pem"`
	AdapterHost                 string `envconfig:"ADAPTER_HOST" default:"adapter"`
	RatelimiterHost             string `envconfig:"RATELIMITER_HOST" default:"apk-test-wso2-apk-ratelimiter-service.apk.svc"`
	RatelimiterPort             string `envconfig:"RATELIMITER_PORT" default:"8091"`
	AdapterHostname             string `envconfig:"ADAPTER_HOST_NAME" default:"adapter"`
	AdapterXdsPort              string `envconfig:"ADAPTER_XDS_PORT" default:"18000"`
	CommonControllerHost        string `envconfig:"COMMON_CONTROLLER_HOST" default:"common-controller"`
	CommonControllerHostname    string `envconfig:"COMMON_CONTROLLER_HOST_NAME" default:"common-controller"`
	CommonControllerXdsPort     string `envconfig:"COMMON_CONTROLLER_XDS_PORT" default:"18002"`
	CommonControllerRestPort    string `envconfig:"COMMON_CONTROLLER_REST_PORT" default:"18003"`
	XdsMaxMsgSize               int    `envconfig:"XDS_MAX_MSG_SIZE" default:"4194304"`
	EnforcerLabel               string `envconfig:"ENFORCER_LABEL" default:"enforcer"`
	EnforcerRegionID            string `envconfig:"ENFORCER_REGION" default:"UNKNOWN"`
	XdsMaxRetries               int    `envconfig:"XDS_MAX_RETRIES" default:"3"`
	XdsRetryPeriod              int    `envconfig:"XDS_RETRY_PERIOD" default:"5000"` // milliseconds
	InstanceIdentifier          string `envconfig:"HOSTNAME" default:"Unassigned"`
	RedisUsername               string `envconfig:"REDIS_USERNAME" default:""`
	RedisPassword               string `envconfig:"REDIS_PASSWORD" default:""`
	RedisHost                   string `envconfig:"REDIS_HOST" default:"redis-master"`
	RedisPort                   int    `envconfig:"REDIS_PORT" default:"6379"`
	IsRedisTLSEnabled           bool   `envconfig:"IS_REDIS_TLS_ENABLED" default:"false"`
	RevokedTokensRedisChannel   string `envconfig:"REDIS_REVOKED_TOKENS_CHANNEL" default:"wso2-apk-revoked-tokens-channel"`
	RedisKeyFile                string `envconfig:"REDIS_KEY_FILE" default:"/home/wso2/security/redis/redis.key"`
	RedisCertFile               string `envconfig:"REDIS_CERT_FILE" default:"/home/wso2/security/redis/redis.crt"`
	RedisCaCertFile             string `envconfig:"REDIS_CA_CERT_FILE" default:"/home/wso2/security/redis/ca.crt"`
	RevokedTokenCleanupInterval int    `envconfig:"REVOKED_TOKEN_CLEANUP_INTERVAL" default:"3600"` // seconds
	ChoreoAnalyticsAuthToken    string `envconfig:"CHOREO_ANALYTICS_AUTH_TOKEN" default:""`
	ChoreoAnalyticsAuthURL      string `envconfig:"CHOREO_ANALYTICS_AUTH_URL" default:""`
	MoesifToken                 string `envconfig:"MOESIF_TOKEN" default:""`
	LogLevel                    string `envconfig:"LOG_LEVEL" default:"INFO"`
	ExternalProcessingPort      string `envconfig:"EXTERNAL_PROCESSING_PORT" default:"8080"`
	Logger                      logging.Logger
}

// package-level variable and mutex for thread safety
var (
	processOnce     sync.Once
	settingInstance *Server
)

// GetConfig initializes and returns a singleton instance of the Settings struct.
// It uses sync.Once to ensure that the initialization logic is executed only once,
// making it safe for concurrent use. If there is an error during the initialization,
// the function will panic.
//
// Returns:
//
//	*Settings - A pointer to the singleton instance of the Settings struct. from environment variables.
func GetConfig() *Server {
	var err error
	processOnce.Do(func() {
		settingInstance = &Server{}
		err = envconfig.Process("", settingInstance)
	})
	if err != nil {
		panic(err)
	}
	// Create Logger based on the env var
	settingInstance.Logger = logging.NewLogger(&egv1a1.EnvoyGatewayLogging{
		Level: map[egv1a1.EnvoyGatewayLogComponent]egv1a1.LogLevel{
			egv1a1.LogComponentGatewayDefault: egv1a1.LogLevel(settingInstance.LogLevel),
		},
	})
	return settingInstance
}
