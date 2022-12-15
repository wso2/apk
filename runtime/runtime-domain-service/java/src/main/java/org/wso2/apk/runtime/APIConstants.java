/*
 * Copyright (c) 2022, WSO2 LLC. (http://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.wso2.apk.runtime;

import java.util.Arrays;
import java.util.Collections;
import java.util.HashSet;
import java.util.Set;

/**
 * This class represents the constants that are used for APIManager implementation
 */
public final class APIConstants {
    public static final String AUTH_APPLICATION_OR_USER_LEVEL_TOKEN = "Any";
    public static final String DEFAULT_SUB_POLICY_UNLIMITED = "Unlimited";
    public static final String HTTP_POST = "POST";
    //Swagger v2.0 constants
    public static final String SWAGGER_X_SCOPE = "x-scope";
    public static final String SWAGGER_X_AMZN_RESOURCE_NAME = "x-amzn-resource-name";
    public static final String SWAGGER_X_AMZN_RESOURCE_TIMEOUT = "x-amzn-resource-timeout";
    public static final String SWAGGER_X_AUTH_TYPE = "x-auth-type";
    public static final String SWAGGER_X_THROTTLING_TIER = "x-throttling-tier";
    public static final String SWAGGER_X_THROTTLING_BANDWIDTH = "x-throttling-bandwidth";
    public static final String SWAGGER_X_MEDIATION_SCRIPT = "x-mediation-script";
    public static final String SWAGGER_X_WSO2_SECURITY = "x-wso2-security";
    public static final String SWAGGER_X_WSO2_SCOPES = "x-wso2-scopes";
    public static final String SWAGGER_X_EXAMPLES = "x-examples";
    public static final String SWAGGER_SCOPE_KEY = "key";
    public static final String SWAGGER_NAME = "name";
    public static final String SWAGGER_DESCRIPTION = "description";
    public static final String SWAGGER_ROLES = "roles";
    public static final String SWAGGER_OBJECT_NAME_APIM = "apim";
    public static final String SWAGGER_RESPONSE_200 = "200";
    public static final String SWAGGER_APIM_DEFAULT_SECURITY = "default";
    public static final String OPENAPI_IS_MISSING_MSG = "openapi is missing";
    public static final String SWAGGER_X_SCOPES_BINDINGS = "x-scopes-bindings";
    public static final String DEFAULT_API_SECURITY_OAUTH2 = "oauth2";

    public static final String STRING = "string";
    public static final String OBJECT = "object";

    public static final String GRAPHQL_API = "GRAPHQL";
    public static final String APPLICATION_JSON_MEDIA_TYPE = "application/json";

    // registry location for OpenAPI files
    public static final String OPENAPI_ARCHIVES_TEMP_FOLDER = "OPENAPI-archives";
    public static final String OPENAPI_EXTRACTED_DIRECTORY = "extracted";
    public static final String OPENAPI_ARCHIVE_ZIP_FILE = "openapi-archive.zip";
    public static final String OPENAPI_MASTER_JSON = "swagger.json";
    public static final String OPENAPI_MASTER_YAML = "swagger.yaml";

    public static final String HTTPS_PROTOCOL = "https";
    public static final String HTTPS_PROTOCOL_URL_PREFIX = "https://";
    public static final String HTTP_PROTOCOL = "http";
    public static final String HTTP_PROTOCOL_URL_PREFIX = "http://";

    //URI Authentication Schemes
    public static final Set<String> SUPPORTED_METHODS =
            Collections.unmodifiableSet(new HashSet<String>(
                    Arrays.asList(new String[]{"get", "put", "post", "delete", "patch", "head", "options"})));
    public static final String TYPE = "Type";
    public static final String JAVA_IO_TMPDIR = "java.io.tmpdir";
    public static final String WSO2_GATEWAY_ENVIRONMENT = "wso2";
    public static final String HTTP_VERB_PUBLISH = "PUBLISH";
    public static final String HTTP_VERB_SUBSCRIBE = "SUBSCRIBE";
    public static final String API_TYPE_WS = "WS";
    // Protocol variables
    public static final String HTTP_TRANSPORT_PROTOCOL_NAME = "http";
    public static final String HTTPS_TRANSPORT_PROTOCOL_NAME = "https";
    public static final String WS_TRANSPORT_PROTOCOL_NAME = "ws";
    public static final String KAFKA_TRANSPORT_PROTOCOL_NAME = "kafka";
    public static final String AMQP_TRANSPORT_PROTOCOL_NAME = "amqp";
    public static final String AMQPS_TRANSPORT_PROTOCOL_NAME = "amqps";
    public static final String AMQP1_TRANSPORT_PROTOCOL_NAME = "amqp1";
    public static final String MQTT_TRANSPORT_PROTOCOL_NAME = "mqtt";
    public static final String SECURE_MQTT_TRANSPORT_PROTOCOL_NAME = "secure-mqtt";
    public static final String WS_MQTT_TRANSPORT_PROTOCOL_NAME = "ws-mqtt";
    public static final String WSS_MQTT_TRANSPORT_PROTOCOL_NAME = "wss-mqtt";
    public static final String MQTT5_TRANSPORT_PROTOCOL_NAME = "mqtt5";
    public static final String NATS_TRANSPORT_PROTOCOL_NAME = "nats";
    public static final String JMS_TRANSPORT_PROTOCOL_NAME = "jms";
    public static final String SNS_TRANSPORT_PROTOCOL_NAME = "sns";
    public static final String SQS_TRANSPORT_PROTOCOL_NAME = "sqs";
    public static final String STOMP_TRANSPORT_PROTOCOL_NAME = "stomp";
    public static final String REDIS_TRANSPORT_PROTOCOL_NAME = "redis";
    public static final String SMF_TRANSPORT_PROTOCOL_NAME = "smf";
    public static final String SMF_TRANSPORT_PROTOCOL_VERSION = "smf";
    public static final String SMFS_TRANSPORT_PROTOCOL_NAME = "smfs";
    public static final String SMFS_TRANSPORT_PROTOCOL_VERSION = "smfs";
    // GraphQL related constants
    public static final String GRAPHQL_QUERY = "Query";
    public static final String GRAPHQL_MUTATION = "Mutation";
    public static final String GRAPHQL_SUBSCRIPTION = "Subscription";
    public static final String SCOPE_ROLE_MAPPING = "WSO2ScopeRoleMapping";
    public static final String SCOPE_OPERATION_MAPPING = "WSO2ScopeOperationMapping";
    public static final String OPERATION_THROTTLING_MAPPING = "WSO2OperationThrottlingMapping";
    public static final String OPERATION_AUTH_SCHEME_MAPPING = "WSO2OperationAuthSchemeMapping";
    public static final String OPERATION_SECURITY_ENABLED = "Enabled";
    public static final String OPERATION_SECURITY_DISABLED = "Disabled";
    public static final String GRAPHQL_PAYLOAD = "GRAPHQL_PAYLOAD";
    public static final String GRAPHQL_SCHEMA = "GRAPHQL_SCHEMA";
    public static final String GRAPHQL_ACCESS_CONTROL_POLICY = "WSO2GraphQLAccessControlPolicy";
    public static final String QUERY_ANALYSIS_COMPLEXITY = "complexity";
    public static final String MAXIMUM_QUERY_COMPLEXITY = "max_query_complexity";
    public static final String MAXIMUM_QUERY_DEPTH = "max_query_depth";
    public static final String GRAPHQL_MAX_DEPTH = "graphQLMaxDepth";
    public static final String GRAPHQL_MAX_COMPLEXITY = "graphQLMaxComplexity";
    public static final String GRAPHQL_ADDITIONAL_TYPE_PREFIX = "WSO2";
    public static final String AUTH_NO_AUTHENTICATION = "None";

    public enum ParserType {
        OAS3, OAS2, ASYNC
    }


    public static class OperationParameter {

        public static final String PAYLOAD_PARAM_NAME = "Payload";

        private OperationParameter() {

        }
    }

}

