/*
 * Copyright (c) 2021, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.wso2.apk.enforcer.config;

import org.apache.commons.lang3.StringUtils;
import org.wso2.apk.enforcer.constants.Constants;

/**
 * Holds and returns the configuration values retrieved from the environment variables.
 */
public class EnvVarConfig {
    private static final String TRUSTED_CA_CERTS_PATH = "TRUSTED_CA_CERTS_PATH";
    private static final String TRUST_DEFAULT_CERTS = "TRUST_DEFAULT_CERTS";
    private static final String ADAPTER_HOST_NAME = "ADAPTER_HOST_NAME";
    private static final String COMMON_CONTROLLER_HOST_NAME = "COMMON_CONTROLLER_HOST_NAME";
    private static final String ENFORCER_PRIVATE_KEY_PATH = "ENFORCER_PRIVATE_KEY_PATH";
    private static final String ENFORCER_PUBLIC_CERT_PATH = "ENFORCER_PUBLIC_CERT_PATH";
    private static final String OPA_CLIENT_PRIVATE_KEY_PATH = "OPA_CLIENT_PRIVATE_KEY_PATH";
    private static final String OPA_CLIENT_PUBLIC_CERT_PATH = "OPA_CLIENT_PUBLIC_CERT_PATH";
    private static final String ADAPTER_HOST = "ADAPTER_HOST";
    private static final String ADAPTER_XDS_PORT = "ADAPTER_XDS_PORT";
    private static final String COMMON_CONTROLLER_HOST = "COMMON_CONTROLLER_HOST";
    private static final String COMMON_CONTROLLER_XDS_PORT = "COMMON_CONTROLLER_XDS_PORT";
    private static final String COMMON_CONTROLLER_REST_PORT = "COMMON_CONTROLLER_REST_PORT";

    private static final String ENFORCER_LABEL = "ENFORCER_LABEL";
    private static final String ENFORCER_REGION_ID = "ENFORCER_REGION";
    public static final String XDS_MAX_MSG_SIZE = "XDS_MAX_MSG_SIZE";
    public static final String XDS_MAX_RETRIES = "XDS_MAX_RETRIES";
    public static final String XDS_RETRY_PERIOD = "XDS_RETRY_PERIOD";
    public static final String HOSTNAME = "HOSTNAME";
    public static final String REDIS_USERNAME = "REDIS_USERNAME";
    public static final String REDIS_PASSWORD = "REDIS_PASSWORD";
    public static final String REDIS_HOST = "REDIS_HOST";
    public static final String REDIS_PORT = "REDIS_PORT";
    public static final String IS_REDIS_TLS_ENABLED = "IS_REDIS_TLS_ENABLED";
    public static final String REDIS_REVOKED_TOKENS_CHANNEL = "REDIS_REVOKED_TOKENS_CHANNEL";
    public static final String REDIS_KEY_FILE = "REDIS_KEY_FILE";
    public static final String REDIS_CERT_FILE = "REDIS_CERT_FILE";
    public static final String REDIS_CA_CERT_FILE = "REDIS_CA_CERT_FILE";
    public static final String REVOKED_TOKEN_CLEANUP_INTERVAL = "REVOKED_TOKEN_CLEANUP_INTERVAL";
    public static final String CHOREO_ANALYTICS_AUTH_TOKEN = "CHOREO_ANALYTICS_AUTH_TOKEN";
    public static final String CHOREO_ANALYTICS_AUTH_URL = "CHOREO_ANALYTICS_AUTH_URL";


    // Since the container is running in linux container, path separator is not needed.
    private static final String DEFAULT_TRUSTED_CA_CERTS_PATH = "/home/wso2/security/truststore";
    private static final String DEFAULT_TRUST_DEFAULT_CERTS = "true";
    private static final String DEFAULT_ADAPTER_HOST_NAME = "adapter";
    private static final String DEFAULT_COMMON_CONTROLLER_HOST_NAME = "common-controller";
    private static final String DEFAULT_ENFORCER_PRIVATE_KEY_PATH = "/home/wso2/security/keystore/mg.key";
    private static final String DEFAULT_ENFORCER_PUBLIC_CERT_PATH = "/home/wso2/security/keystore/mg.pem";
    private static final String DEFAULT_ENFORCER_REGION_ID = "UNKNOWN";
    private static final String DEFAULT_ADAPTER_HOST = "adapter";
    private static final String DEFAULT_ADAPTER_XDS_PORT = "18000";
    private static final String DEFAULT_COMMON_CONTROLLER_HOST = "common-controller";
    private static final String DEFAULT_COMMON_CONTROLLER_XDS_PORT = "18002";
    private static final String DEFAULT_COMMON_CONTROLLER_REST_PORT = "18003";
    private static final String DEFAULT_ENFORCER_LABEL = "enforcer";
    public static final String DEFAULT_XDS_MAX_MSG_SIZE = "4194304";
    public static final String DEFAULT_XDS_MAX_RETRIES = Integer.toString(Constants.MAX_XDS_RETRIES);
    public static final String DEFAULT_XDS_RETRY_PERIOD = Integer.toString(Constants.XDS_DEFAULT_RETRY);
    public static final String DEFAULT_HOSTNAME = "Unassigned";
    public static final String DEFAULT_REDIS_USERNAME = "default";
    public static final String DEFAULT_REDIS_PASSWORD = "";
    public static final String DEFAULT_REDIS_HOST = "redis-master";
    public static final int DEFAULT_REDIS_PORT = 6379;
    public static final String DEFAULT_IS_REDIS_TLS_ENABLED = "false";
    public static final String DEFAULT_REDIS_REVOKED_TOKENS_CHANNEL = "wso2-apk-revoked-tokens-channel";
    public static final String DEFAULT_REDIS_KEY_FILE = "/home/wso2/security/redis/redis.key";
    public static final String DEFAULT_REDIS_CERT_FILE = "/home/wso2/security/redis/redis.crt";
    public static final String DEFAULT_REDIS_CA_CERT_FILE = "/home/wso2/security/redis/ca.crt";
    public static final int DEFAULT_REVOKED_TOKEN_CLEANUP_INTERVAL = 60*60; // In seconds

    public static final String DEFAULT_CHOREO_ANALYTICS_AUTH_TOKEN = "";
    public static final String DEFAULT_CHOREO_ANALYTICS_AUTH_URL = "";

    private static EnvVarConfig instance;
    private final String trustedAdapterCertsPath;
    private final String trustDefaultCerts;
    private final String enforcerPrivateKeyPath;
    private final String enforcerPublicKeyPath;
    private final String opaClientPrivateKeyPath;
    private final String opaClientPublicKeyPath;
    private final String adapterHost;
    private final String commonControllerHost;
    private final String enforcerLabel;
    private final String adapterXdsPort;
    private final String commonControllerXdsPort;
    private final String commonControllerRestPort;
    private final String adapterHostname;
    private final String commonControllerHostname;
    // TODO: (VirajSalaka) Enforcer ID should be picked from router once envoy 1.18.0 is released and microgateway
    // is updated.
    private final String enforcerRegionId;
    private final String xdsMaxMsgSize;
    private final String xdsMaxRetries;
    private final String xdsRetryPeriod;
    private final String instanceIdentifier;
    private final String redisUsername;
    private final String redisPassword;
    private final String redisHost;
    private final int redisPort;
    private final boolean isRedisTlsEnabled;
    private final String revokedTokensRedisChannel;
    private final String redisKeyFile;
    private final String redisCertFile;
    private final String redisCaCertFile;

    private final String choreoAnalyticsAuthToken;
    private final String choreoAnalyticsAuthUrl;
    private final int revokedTokenCleanupInterval;

    private EnvVarConfig() {
        trustedAdapterCertsPath = retrieveEnvVarOrDefault(TRUSTED_CA_CERTS_PATH,
                DEFAULT_TRUSTED_CA_CERTS_PATH);
        trustDefaultCerts = retrieveEnvVarOrDefault(TRUST_DEFAULT_CERTS,
                DEFAULT_TRUST_DEFAULT_CERTS);
        enforcerPrivateKeyPath = retrieveEnvVarOrDefault(ENFORCER_PRIVATE_KEY_PATH,
                DEFAULT_ENFORCER_PRIVATE_KEY_PATH);
        enforcerPublicKeyPath = retrieveEnvVarOrDefault(ENFORCER_PUBLIC_CERT_PATH,
                DEFAULT_ENFORCER_PUBLIC_CERT_PATH);
        opaClientPrivateKeyPath = retrieveEnvVarOrDefault(OPA_CLIENT_PRIVATE_KEY_PATH,
                DEFAULT_ENFORCER_PRIVATE_KEY_PATH);
        opaClientPublicKeyPath = retrieveEnvVarOrDefault(OPA_CLIENT_PUBLIC_CERT_PATH,
                DEFAULT_ENFORCER_PUBLIC_CERT_PATH);
        enforcerLabel = retrieveEnvVarOrDefault(ENFORCER_LABEL, DEFAULT_ENFORCER_LABEL);
        adapterHost = retrieveEnvVarOrDefault(ADAPTER_HOST, DEFAULT_ADAPTER_HOST);
        adapterHostname = retrieveEnvVarOrDefault(ADAPTER_HOST_NAME, DEFAULT_ADAPTER_HOST_NAME);
        adapterXdsPort = retrieveEnvVarOrDefault(ADAPTER_XDS_PORT, DEFAULT_ADAPTER_XDS_PORT);
        commonControllerHost = retrieveEnvVarOrDefault(COMMON_CONTROLLER_HOST, DEFAULT_COMMON_CONTROLLER_HOST);
        commonControllerHostname = retrieveEnvVarOrDefault(COMMON_CONTROLLER_HOST_NAME,
                DEFAULT_COMMON_CONTROLLER_HOST_NAME);
        commonControllerXdsPort = retrieveEnvVarOrDefault(COMMON_CONTROLLER_XDS_PORT,
                DEFAULT_COMMON_CONTROLLER_XDS_PORT);
        commonControllerRestPort = retrieveEnvVarOrDefault(COMMON_CONTROLLER_REST_PORT,
                DEFAULT_COMMON_CONTROLLER_REST_PORT);
        xdsMaxMsgSize = retrieveEnvVarOrDefault(XDS_MAX_MSG_SIZE, DEFAULT_XDS_MAX_MSG_SIZE);
        enforcerRegionId = retrieveEnvVarOrDefault(ENFORCER_REGION_ID, DEFAULT_ENFORCER_REGION_ID);
        xdsMaxRetries = retrieveEnvVarOrDefault(XDS_MAX_RETRIES, DEFAULT_XDS_MAX_RETRIES);
        xdsRetryPeriod = retrieveEnvVarOrDefault(XDS_RETRY_PERIOD, DEFAULT_XDS_RETRY_PERIOD);
        // HOSTNAME environment property is readily available in docker and kubernetes, and it represents the Pod
        // name in Kubernetes context, containerID in docker context.
        instanceIdentifier = retrieveEnvVarOrDefault(HOSTNAME, DEFAULT_HOSTNAME);
        redisUsername = retrieveEnvVarOrDefault(REDIS_USERNAME, DEFAULT_REDIS_USERNAME);
        redisPassword = retrieveEnvVarOrDefault(REDIS_PASSWORD, DEFAULT_REDIS_PASSWORD);
        redisHost = retrieveEnvVarOrDefault(REDIS_HOST, DEFAULT_REDIS_HOST);
        redisPort = getRedisPortFromEnv();
        isRedisTlsEnabled = retrieveEnvVarOrDefault(IS_REDIS_TLS_ENABLED, DEFAULT_IS_REDIS_TLS_ENABLED).toLowerCase()
                .equals(DEFAULT_IS_REDIS_TLS_ENABLED)? false:true;
        revokedTokensRedisChannel = retrieveEnvVarOrDefault(REDIS_REVOKED_TOKENS_CHANNEL, DEFAULT_REDIS_REVOKED_TOKENS_CHANNEL);
        redisKeyFile = retrieveEnvVarOrDefault(REDIS_KEY_FILE, DEFAULT_REDIS_KEY_FILE);
        redisCertFile = retrieveEnvVarOrDefault(REDIS_CERT_FILE, DEFAULT_REDIS_CERT_FILE);
        redisCaCertFile = retrieveEnvVarOrDefault(REDIS_CA_CERT_FILE, DEFAULT_REDIS_CA_CERT_FILE);
        revokedTokenCleanupInterval = getRevokedTokenCleanupIntervalFromEnv();
        choreoAnalyticsAuthToken = retrieveEnvVarOrDefault(CHOREO_ANALYTICS_AUTH_TOKEN, DEFAULT_CHOREO_ANALYTICS_AUTH_TOKEN);
        choreoAnalyticsAuthUrl = retrieveEnvVarOrDefault(CHOREO_ANALYTICS_AUTH_URL, DEFAULT_CHOREO_ANALYTICS_AUTH_URL);
    }

    public static EnvVarConfig getInstance() {
        if (instance == null) {
            synchronized (EnvVarConfig.class) {
                if (instance == null) {
                    instance = new EnvVarConfig();
                }
            }
        }
        return instance;
    }

    private int getRedisPortFromEnv() {
        String portStr = retrieveEnvVarOrDefault(REDIS_PORT, String.valueOf(DEFAULT_REDIS_PORT));
        try {
            return Integer.parseInt(portStr);
        } catch (Exception e) {
            return DEFAULT_REDIS_PORT;
        }
    }

    private int getRevokedTokenCleanupIntervalFromEnv() {
        String intervalStr = retrieveEnvVarOrDefault(REVOKED_TOKEN_CLEANUP_INTERVAL, String.valueOf(DEFAULT_REVOKED_TOKEN_CLEANUP_INTERVAL));
        try {
            return Integer.parseInt(intervalStr);
        } catch (Exception e) {
            return DEFAULT_REVOKED_TOKEN_CLEANUP_INTERVAL;
        }
    }

    private String retrieveEnvVarOrDefault(String variable, String defaultValue) {
        if (StringUtils.isEmpty(System.getenv(variable))) {
            return defaultValue;
        }
        return System.getenv(variable);
    }

    public String getTrustedAdapterCertsPath() {
        return trustedAdapterCertsPath;
    }

    public boolean isTrustDefaultCerts() {
        return Boolean.valueOf(trustDefaultCerts);
    }

    public String getEnforcerPrivateKeyPath() {
        return enforcerPrivateKeyPath;
    }

    public String getEnforcerPublicKeyPath() {
        return enforcerPublicKeyPath;
    }

    public String getOpaClientPrivateKeyPath() {
        return opaClientPrivateKeyPath;
    }

    public String getOpaClientPublicKeyPath() {
        return opaClientPublicKeyPath;
    }

    public String getAdapterHost() {
        return adapterHost;
    }

    public String getCommonControllerHost() {
        return commonControllerHost;
    }

    public String getEnforcerLabel() {
        return enforcerLabel;
    }

    public String getAdapterXdsPort() {
        return adapterXdsPort;
    }

    public String getCommonControllerXdsPort() {
        return commonControllerXdsPort;
    }

    public String getAdapterHostname() {
        return adapterHostname;
    }

    public String getCommonControllerHostname() {
        return commonControllerHostname;
    }

    public String getXdsMaxMsgSize() {
        return xdsMaxMsgSize;
    }


    public String getEnforcerRegionId() {
        return enforcerRegionId;
    }

    public String getXdsMaxRetries() {
        return xdsMaxRetries;
    }

    public String getXdsRetryPeriod() {
        return xdsRetryPeriod;
    }

    public String getInstanceIdentifier() {
        return instanceIdentifier;
    }

    public String getRedisUsername() {

        return redisUsername;
    }

    public String getRedisPassword() {

        return redisPassword;
    }

    public String getRedisHost() {

        return redisHost;
    }

    public boolean isRedisTlsEnabled() {

        return isRedisTlsEnabled;
    }

    public int getRedisPort() {
        return redisPort;
    }

    public String getRedisKeyFile() {
        return redisKeyFile;
    }

    public String getRedisCertFile() {
        return redisCertFile;
    }

    public String getRedisCaCertFile() {
        return redisCaCertFile;
    }

    public String getRevokedTokensRedisChannel() {
        return revokedTokensRedisChannel;
    }

    public int getRevokedTokenCleanupInterval() {
        return revokedTokenCleanupInterval;
    }

    public String getCommonControllerRestPort() {

        return commonControllerRestPort;
    }

    public String getChoreoAnalyticsAuthToken() {
        return choreoAnalyticsAuthToken;
    }

    public String getChoreoAnalyticsAuthUrl() {
        return choreoAnalyticsAuthUrl;
    }
}

