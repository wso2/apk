/*
 * Copyright (c) 2023, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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

package org.wso2.apk.enforcer.analytics.publisher.util;

/**
 * Class to hold String constants.
 */
public class Constants {
    //Event attribute names
    public static final String CORRELATION_ID = "correlationId";
    public static final String KEY_TYPE = "keyType";
    public static final String API_ID = "apiId";
    public static final String API_NAME = "apiName";
    public static final String API_CONTEXT = "apiContext";
    public static final String USER_NAME = "userName";
    public static final String API_VERSION = "apiVersion";
    public static final String API_CREATION = "apiCreator";
    public static final String API_METHOD = "apiMethod";
    public static final String API_RESOURCE_TEMPLATE = "apiResourceTemplate";
    public static final String API_CREATOR_TENANT_DOMAIN = "apiCreatorTenantDomain";
    public static final String DESTINATION = "destination";
    public static final String APPLICATION_ID = "applicationId";
    public static final String APPLICATION_NAME = "applicationName";
    public static final String APPLICATION_OWNER = "applicationOwner";
    public static final String REGION_ID = "regionId";
    public static final String ORGANIZATION_ID = "organizationId";
    public static final String ENVIRONMENT_ID = "environmentId";
    public static final String GATEWAY_TYPE = "gatewayType";
    public static final String USER_AGENT_HEADER = "userAgentHeader";
    public static final String USER_AGENT = "userAgent";
    public static final String PLATFORM = "platform";
    public static final String PROXY_RESPONSE_CODE = "proxyResponseCode";
    public static final String TARGET_RESPONSE_CODE = "targetResponseCode";
    public static final String RESPONSE_CACHE_HIT = "responseCacheHit";
    public static final String RESPONSE_LATENCY = "responseLatency";
    public static final String BACKEND_LATENCY = "backendLatency";
    public static final String REQUEST_MEDIATION_LATENCY = "requestMediationLatency";
    public static final String RESPONSE_MEDIATION_LATENCY = "responseMediationLatency";
    public static final String REQUEST_TIMESTAMP = "requestTimestamp";
    public static final String EVENT_TYPE = "eventType";
    public static final String API_TYPE = "apiType";
    public static final String USER_IP = "userIp";
    public static final String ERROR_TYPE = "errorType";
    public static final String ERROR_CODE = "errorCode";
    public static final String ERROR_MESSAGE = "errorMessage";
    public static final String PROPERTIES = "properties";

    //Builder event types
    public static final String RESPONSE_EVENT_TYPE = "response";
    public static final String FAULT_EVENT_TYPE = "fault";

    //Reporter config properties
    public static final String AUTH_API_URL = "authURL";
    public static final String AUTH_API_TOKEN = "authToken";
    public static final String MOESIF_TOKEN = "moesifToken";

    //Proxy configs
    public static final String PROXY_ENABLE = "proxy_config_enable";
    public static final String PROXY_HOST = "proxy_config_host";
    public static final String PROXY_PORT = "proxy_config_port";
    public static final String PROXY_USERNAME = "proxy_config_username";
    public static final String PROXY_PASSWORD = "proxy_config_password";
    public static final String PROXY_PROTOCOL = "proxy_config_protocol";
    public static final String KEYSTORE_LOCATION = "keystore_location";
    public static final String KEYSTORE_PASSWORD = "keystore_password";
    public static final String HTTPS_PROTOCOL = "https";
    public static final String HTTP_PROTOCOL = "http";
    public static final String KEYSTORE_TYPE = "JKS";

    //Reporter constants
    public static final String DEFAULT_REPORTER = "default";
    public static final String ELK_REPORTER = "elk";
    public static final String MOESIF_REPORTER = "moesif";
    public static final String PROMETHEUS_REPORTER = "prometheus";

    //EventHub Client retry options constants
    public static final int DEFAULT_MAX_RETRIES = 2;
    public static final int DEFAULT_DELAY = 15;
    public static final int DEFAULT_MAX_DELAY = 30;
    public static final int DEFAULT_TRY_TIMEOUT = 30;
    public static final String EVENTHUB_CLIENT_MAX_RETRIES = "eventhub.client.max.retries";
    public static final String EVENTHUB_CLIENT_DELAY = "eventhub.client.delay";
    public static final String EVENTHUB_CLIENT_MAX_DELAY = "eventhub.client.max.delay";
    public static final String EVENTHUB_CLIENT_TRY_TIMEOUT = "eventhub.client.try.timeout";
    public static final String EVENTHUB_CLIENT_RETRY_MODE = "eventhub.client.retry.mode";
    public static final String FIXED = "fixed";
    public static final String EXPONENTIAL = "exponential";

    public static final String UNKNOWN_VALUE = "UNKNOWN";
    public static final String WORKER_THREAD_COUNT = "worker.thread.count";
    public static final String QUEUE_SIZE = "queue.size";
    public static final String CLIENT_FLUSHING_DELAY = "client.flushing.delay";
    public static final int DEFAULT_QUEUE_SIZE = 20000;
    public static final int DEFAULT_WORKER_THREADS = 1;
    public static final int DEFAULT_FLUSHING_DELAY = 15;
    public static final int USER_AGENT_DEFAULT_CACHE_SIZE = 50;

    // Moesif sdk related constants
    public static final String MOESIF_CONTENT_TYPE_HEADER = "application/json";
    public static final String MOESIF_CONTENT_TYPE_KEY = "Content-Type";
    public static final String MOESIF_USER_AGENT_KEY = "User-Agent";
    public static final String GATEWAY_URL = "x-original-gw-url";
    public static final String DEPLOYMENT_TYPE = "deployment-type";
    public static final String PRODUCTION = "PRODUCTION";

    public static final String MOESIF_KEY_RETRIEVER_CLIENT_TYPE = "moesifKeyRetrieverClientType";

    public static final String MOESIF_KEY_RETRIEVER_CHOREO_CLIENT = "Choreo";
    public static final String MOESIF_KEY_VALUE = "moesifToken";
    public static final String DEFAULT_ENVIRONMENT = "Default";
}
