/*
 * Copyright (c) 2020, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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
package org.wso2.apk.enforcer.constants;

import java.util.List;

/**
 * Holds the common set of constants for the enforcer package.
 */
public class APIConstants {

    public static final String DEFAULT = "default";
    public static final String GW_VHOST_PARAM = "vHost";
    public static final String ROUTE_NAME_PARAM = "route-name";
    public static final String GW_BASE_PATH_PARAM = "basePath";
    public static final String GW_RES_PATH_PARAM = "path";
    public static final String GW_VERSION_PARAM = "version";
    public static final String GW_API_NAME_PARAM = "name";
    public static final String PROTOTYPED_LIFE_CYCLE_STATUS = "PROTOTYPED";
    public static final String UNLIMITED_TIER = "Unlimited";
    public static final String UNAUTHENTICATED_TIER = "Unauthenticated";
    public static final String END_USER_ANONYMOUS = "anonymous";
    public static final String END_USER_UNKNOWN = "unknown";
    public static final String ANONYMOUS_PREFIX = "anon:";
    public static final String PUBLISHER_CERTIFICATE_ALIAS = "publisher_certificate_alias";
    public static final String APIKEY_CERTIFICATE_ALIAS = "apikey_certificate_alias";
    public static final String WSO2_PUBLIC_CERTIFICATE_ALIAS = "wso2carbon";
    public static final String HTTPS_PROTOCOL = "https";
    public static final String SUPER_TENANT_DOMAIN_NAME = "carbon.super";
    public static final String BANDWIDTH_TYPE = "bandwidthVolume";
    public static final String AUTHORIZATION_HEADER_DEFAULT = "Authorization";
    public static final String AUTHORIZATION_BEARER = "Bearer ";
    public static final String API_KEY_TYPE_PRODUCTION = "PRODUCTION";
    public static final String API_KEY_TYPE_SANDBOX = "SANDBOX";
    public static final String DEFAULT_ENVIRONMENT_NAME = "Default";

    public static final String AUTHORIZATION_HEADER_BASIC = "Basic";
    public static final String API_SECURITY_OAUTH2 = "OAuth2";
    public static final String API_SECURITY_BASIC_AUTH = "basic_auth";
    public static final String API_SECURITY_API_KEY = "\"API Key\"";
    public static final String SWAGGER_API_KEY_IN_HEADER = "Header";
    public static final String SWAGGER_API_KEY_IN_QUERY = "Query";
    public static final String API_SECURITY_MUTUAL_SSL_NAME = "mtls";
    public static final String CLIENT_CERTIFICATE_HEADER_DEFAULT = "X-WSO2-CLIENT-CERTIFICATE";
    public static final String WWW_AUTHENTICATE = "WWW-Authenticate";
    public static final String TEST_CONSOLE_KEY_HEADER = "internal-key";

    public static final String BEGIN_CERTIFICATE_STRING = "-----BEGIN CERTIFICATE-----";
    public static final String END_CERTIFICATE_STRING = "-----END CERTIFICATE-----";
    public static final String EVENT_PAYLOAD = "event";
    public static final String EVENT_PAYLOAD_DATA = "payloadData";

    public static final String NOT_FOUND_MESSAGE = "Not Found";
    public static final String NOT_FOUND_DESCRIPTION = "The requested resource is not available.";
    public static final String NOT_IMPLEMENTED_MESSAGE = "Not Implemented";
    public static final String BAD_REQUEST_MESSAGE = "Bad Request";
    public static final String INTERNAL_SERVER_ERROR_MESSAGE = "Internal Server Error";

    //headers and values
    public static final String CONTENT_TYPE_HEADER = "content-type";
    public static final String SOAP_ACTION_HEADER_NAME = "soapaction";
    public static final String ACCEPT_HEADER = "accept";
    public static final String PREFER_HEADER = "prefer";
    public static final String PREFER_CODE = "code";
    public static final String PREFER_EXAMPLE = "example";
    public static final List<String> PREFER_KEYS = List.of(PREFER_CODE, PREFER_EXAMPLE);
    public static final String APPLICATION_JSON = "application/json";
    public static final String CONTENT_TYPE_TEXT_XML = "text/xml";
    public static final String CONTENT_TYPE_SOAP_XML = "application/soap+xml";
    public static final String APPLICATION_GRAPHQL = "application/graphql";
    public static final String X_FORWARDED_FOR = "x-forwarded-for";
    public static final String PATH_HEADER = ":path";
    public static final String UPGRADE_HEADER = "upgrade";
    public static final String WEBSOCKET = "websocket";

    public static final String LOG_TRACE_ID = "traceId";

    // SOAP protocol versions
    public static final String SOAP11_PROTOCOL = "SOAP 1.1 Protocol";
    public static final String SOAP12_PROTOCOL = "SOAP 1.2 Protocol";
    public static final String API_UUID = "API_UUID";
    public static final String HTTP_PROTOCOL = "http";

    /**
     * Holds the constants related to denied response types.
     */
    public static class ErrorResponseTypes {

        public static final String SOAP11 = "SOAP11";
        public static final String SOAP12 = "SOAP12";
        public static final String JSON = "JSON";
    }

    /**
     * Holds the common set of constants related to the output status codes of the security validations.
     */
    public static class KeyValidationStatus {

        public static final int API_AUTH_INVALID_CREDENTIALS = 900901;
        public static final int API_BLOCKED = 900907;
        public static final int API_AUTH_RESOURCE_FORBIDDEN = 900908;
        public static final int SUBSCRIPTION_INACTIVE = 900909;
        public static final int INVALID_SCOPE = 900910;
        public static final int SUBSCRIPTION_ON_HOLD = 900911;
        public static final int SUBSCRIPTION_REJECTED = 900912;
        public static final int SUBSCRIPTION_BLOCKED = 900913;
        public static final int SUBSCRIPTION_PROD_BLOCKED = 900914;


        private KeyValidationStatus() {

        }
    }

    /**
     * Holds the common set of constants for output of the subscription validation.
     */
    public static class SubscriptionStatus {

        public static final String BLOCKED = "BLOCKED";
        public static final String PROD_ONLY_BLOCKED = "PROD_ONLY_BLOCKED";
        public static final String ON_HOLD = "ON_HOLD";
        public static final String REJECTED = "REJECTED";
        public static final String INACTIVE = "INACTIVE";

        private SubscriptionStatus() {

        }
    }

    /**
     * Holds the common set of constants related to life cycle states.
     */
    public static class LifecycleStatus {

        public static final String BLOCKED = "BLOCKED";
    }

    /**
     * Holds the common set of constants for validating the JWT tokens.
     */
    public static class JwtTokenConstants {

        public static final String APPLICATION = "application";
        public static final String APPLICATION_ID = "id";
        public static final String APPLICATION_UUID = "uuid";
        public static final String APPLICATION_NAME = "name";
        public static final String APPLICATION_TIER = "tier";
        public static final String APPLICATION_OWNER = "owner";
        public static final String SUBSCRIPTION_TIER = "subscriptionTier";
        public static final String SUBSCRIBER_TENANT_DOMAIN = "subscriberTenantDomain";
        public static final String TIER_INFO = "tierInfo";
        public static final String STOP_ON_QUOTA_REACH = "stopOnQuotaReach";
        public static final String SCOPE = "scope";
        public static final String SCOPE_DELIMITER = " ";
        public static final String TOKEN_TYPE = "token_type";
        public static final String SUBSCRIBED_APIS = "subscribedAPIs";
        public static final String API_CONTEXT = "context";
        public static final String API_VERSION = "version";
        public static final String API_PUBLISHER = "publisher";
        public static final String API_NAME = "name";
        public static final String QUOTA_TYPE = "tierQuotaType";
        public static final String QUOTA_TYPE_BANDWIDTH = "bandwidthVolume";
        public static final String PERMITTED_IP = "permittedIP";
        public static final String PERMITTED_REFERER = "permittedReferer";
        public static final String INTERNAL_KEY_TOKEN_TYPE = "InternalKey";
        public static final String PARAM_SEPARATOR = "&";
        public static final String PARAM_VALUE_SEPARATOR = "=";
        public static final String INTERNAL_KEY_APP_NAME = "internal-key-app";

    }

    /**
     * Holds the common set of constants for validating the JWT tokens
     */
    public static class KeyManager {

        public static final String DEFAULT_KEY_MANAGER = "Resident Key Manager";
        public static final String APIM_PUBLISHER_ISSUER = "APIM Publisher";
        public static final String APIM_APIKEY_ISSUER = "APIM APIkey";

        // APIM_APIKEY_ISSUER_URL is intentionally different from the Resident Key Manager
        // to avoid conflicts with the access token issuer configs.
        public static final String APIM_APIKEY_ISSUER_URL = "https://host:9443/apikey";

        public static final String ISSUER = "issuer";
        public static final String JWKS_ENDPOINT = "jwks_endpoint";
        public static final String SELF_VALIDATE_JWT = "self_validate_jwt";
        public static final String CLAIM_MAPPING = "claim_mappings";
        public static final String CONSUMER_KEY_CLAIM = "consumer_key_claim";
        public static final String SCOPES_CLAIM = "scopes_claim";
        public static final String CERTIFICATE_TYPE = "certificate_type";
        public static final String CERTIFICATE_VALUE = "certificate_value";
        public static final String CERTIFICATE_TYPE_JWKS_ENDPOINT = "JWKS";
    }

    /**
     * Supported policy types.
     */
    public enum PolicyType {
        API,
        APPLICATION,
        SUBSCRIPTION
    }

    /**
     * Holds the constants related to attributes to be sent in the response in case of an error
     * scenario raised within the enforcer.
     */
    public static class MessageFormat {

        public static final String STATUS_CODE = "status_code";
        public static final String ERROR_CODE = "code";
        public static final String ERROR_MESSAGE = "error_message";
        public static final String ERROR_DESCRIPTION = "error_description";
    }

    /**
     * Holds the values related http status codes.
     */
    public enum StatusCodes {
        OK("200", 200),
        UNAUTHENTICATED("401", 401),
        UNAUTHORIZED("403", 403),
        NOTFOUND("404", 404),
        THROTTLED("429", 429),
        SERVICE_UNAVAILABLE("503", 503),
        INTERNAL_SERVER_ERROR("500", 500),
        BAD_REQUEST_ERROR("400", 400),
        NOT_IMPLEMENTED_ERROR("501", 501);

        private final String value;
        private final int code;

        StatusCodes(String value, int code) {

            this.value = value;
            this.code = code;
        }

        public String getValue() {

            return this.value;
        }

        public int getCode() {

            return this.code;
        }
    }

    /**
     * Holds the values for API types
     */
    public static class ApiType {

        public static final String WEB_SOCKET = "WS";
        public static final String GRAPHQL = "GraphQL";
    }

    /**
     * Holds values for optionality.
     */
    public static class Optionality {

        public static final String MANDATORY = "mandatory";
        public static final String OPTIONAL = "optional";
    }

    /**
     * Holds values related to readiness check.
     */
    public static class ReadinessCheck {
        public static final String ENDPOINT = "/ready";
        public static final String RESPONSE_KEY = "status";
        public static final String RESPONSE_VALUE = "ready";
    }

}
