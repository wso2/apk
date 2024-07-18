/*
 *  Copyright (c) 2022, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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

package org.wso2.apk.enforcer.analytics.publisher.reporter;

import java.util.AbstractMap;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;
import java.util.stream.Stream;

import static org.wso2.apk.enforcer.analytics.publisher.util.Constants.API_CONTEXT;
import static org.wso2.apk.enforcer.analytics.publisher.util.Constants.API_CREATION;
import static org.wso2.apk.enforcer.analytics.publisher.util.Constants.API_CREATOR_TENANT_DOMAIN;
import static org.wso2.apk.enforcer.analytics.publisher.util.Constants.API_ID;
import static org.wso2.apk.enforcer.analytics.publisher.util.Constants.API_METHOD;
import static org.wso2.apk.enforcer.analytics.publisher.util.Constants.API_NAME;
import static org.wso2.apk.enforcer.analytics.publisher.util.Constants.API_RESOURCE_TEMPLATE;
import static org.wso2.apk.enforcer.analytics.publisher.util.Constants.API_TYPE;
import static org.wso2.apk.enforcer.analytics.publisher.util.Constants.API_VERSION;
import static org.wso2.apk.enforcer.analytics.publisher.util.Constants.APPLICATION_ID;
import static org.wso2.apk.enforcer.analytics.publisher.util.Constants.APPLICATION_NAME;
import static org.wso2.apk.enforcer.analytics.publisher.util.Constants.APPLICATION_OWNER;
import static org.wso2.apk.enforcer.analytics.publisher.util.Constants.BACKEND_LATENCY;
import static org.wso2.apk.enforcer.analytics.publisher.util.Constants.CORRELATION_ID;
import static org.wso2.apk.enforcer.analytics.publisher.util.Constants.DESTINATION;
import static org.wso2.apk.enforcer.analytics.publisher.util.Constants.ENVIRONMENT_ID;
import static org.wso2.apk.enforcer.analytics.publisher.util.Constants.ERROR_CODE;
import static org.wso2.apk.enforcer.analytics.publisher.util.Constants.ERROR_MESSAGE;
import static org.wso2.apk.enforcer.analytics.publisher.util.Constants.ERROR_TYPE;
import static org.wso2.apk.enforcer.analytics.publisher.util.Constants.GATEWAY_TYPE;
import static org.wso2.apk.enforcer.analytics.publisher.util.Constants.KEY_TYPE;
import static org.wso2.apk.enforcer.analytics.publisher.util.Constants.ORGANIZATION_ID;
import static org.wso2.apk.enforcer.analytics.publisher.util.Constants.PROPERTIES;
import static org.wso2.apk.enforcer.analytics.publisher.util.Constants.PROXY_RESPONSE_CODE;
import static org.wso2.apk.enforcer.analytics.publisher.util.Constants.REGION_ID;
import static org.wso2.apk.enforcer.analytics.publisher.util.Constants.REQUEST_MEDIATION_LATENCY;
import static org.wso2.apk.enforcer.analytics.publisher.util.Constants.REQUEST_TIMESTAMP;
import static org.wso2.apk.enforcer.analytics.publisher.util.Constants.RESPONSE_CACHE_HIT;
import static org.wso2.apk.enforcer.analytics.publisher.util.Constants.RESPONSE_LATENCY;
import static org.wso2.apk.enforcer.analytics.publisher.util.Constants.RESPONSE_MEDIATION_LATENCY;
import static org.wso2.apk.enforcer.analytics.publisher.util.Constants.TARGET_RESPONSE_CODE;
import static org.wso2.apk.enforcer.analytics.publisher.util.Constants.USER_AGENT_HEADER;
import static org.wso2.apk.enforcer.analytics.publisher.util.Constants.USER_IP;
import static org.wso2.apk.enforcer.analytics.publisher.util.Constants.USER_NAME;

/**
 * Generic Input Validator for any event metric reporter. In default required attributes include userName.
 * Validator holds all required attributes against which
 * inputs will be validated.
 */
public class GenericInputValidator {

    private static final org.wso2.apk.enforcer.analytics.publisher.reporter.GenericInputValidator INSTANCE =
            new org.wso2.apk.enforcer.analytics.publisher.reporter.GenericInputValidator();
    private static final Map<String, Class> defaultResponseEventSchema = Stream.of(
                    new AbstractMap.SimpleImmutableEntry<>(REQUEST_TIMESTAMP, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(CORRELATION_ID, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(KEY_TYPE, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_ID, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_TYPE, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_NAME, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_VERSION, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_CREATION, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_METHOD, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_CONTEXT, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_RESOURCE_TEMPLATE, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_CREATOR_TENANT_DOMAIN, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(DESTINATION, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(APPLICATION_ID, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(APPLICATION_NAME, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(APPLICATION_OWNER, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(ORGANIZATION_ID, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(ENVIRONMENT_ID, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(REGION_ID, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(GATEWAY_TYPE, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(USER_AGENT_HEADER, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(USER_NAME, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(PROXY_RESPONSE_CODE, Integer.class),
                    new AbstractMap.SimpleImmutableEntry<>(TARGET_RESPONSE_CODE, Integer.class),
                    new AbstractMap.SimpleImmutableEntry<>(RESPONSE_CACHE_HIT, Boolean.class),
                    new AbstractMap.SimpleImmutableEntry<>(RESPONSE_LATENCY, Long.class),
                    new AbstractMap.SimpleImmutableEntry<>(BACKEND_LATENCY, Long.class),
                    new AbstractMap.SimpleImmutableEntry<>(REQUEST_MEDIATION_LATENCY, Long.class),
                    new AbstractMap.SimpleImmutableEntry<>(RESPONSE_MEDIATION_LATENCY, Long.class),
                    new AbstractMap.SimpleImmutableEntry<>(USER_IP, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(PROPERTIES, LinkedHashMap.class))
            .collect(Collectors.toMap(Map.Entry::getKey, Map.Entry::getValue));
    private static final Map<String, Class> elkResponseEventSchema = Stream.of(
                    new AbstractMap.SimpleImmutableEntry<>(REQUEST_TIMESTAMP, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(CORRELATION_ID, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(KEY_TYPE, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_ID, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_TYPE, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_NAME, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_VERSION, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_CREATION, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_METHOD, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_CONTEXT, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_RESOURCE_TEMPLATE, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_CREATOR_TENANT_DOMAIN, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(DESTINATION, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(APPLICATION_ID, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(APPLICATION_NAME, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(APPLICATION_OWNER, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(REGION_ID, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(GATEWAY_TYPE, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(ORGANIZATION_ID, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(ENVIRONMENT_ID, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(USER_AGENT_HEADER, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(USER_NAME, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(PROXY_RESPONSE_CODE, Integer.class),
                    new AbstractMap.SimpleImmutableEntry<>(TARGET_RESPONSE_CODE, Integer.class),
                    new AbstractMap.SimpleImmutableEntry<>(RESPONSE_CACHE_HIT, Boolean.class),
                    new AbstractMap.SimpleImmutableEntry<>(RESPONSE_LATENCY, Long.class),
                    new AbstractMap.SimpleImmutableEntry<>(BACKEND_LATENCY, Long.class),
                    new AbstractMap.SimpleImmutableEntry<>(REQUEST_MEDIATION_LATENCY, Long.class),
                    new AbstractMap.SimpleImmutableEntry<>(RESPONSE_MEDIATION_LATENCY, Long.class),
                    new AbstractMap.SimpleImmutableEntry<>(USER_IP, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(PROPERTIES, LinkedHashMap.class))
            .collect(Collectors.toMap(Map.Entry::getKey, Map.Entry::getValue));
    private static final Map<String, Class> defaultFaultEventSchema = Stream.of(
                    new AbstractMap.SimpleImmutableEntry<>(REQUEST_TIMESTAMP, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(CORRELATION_ID, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(KEY_TYPE, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(ERROR_TYPE, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(ERROR_CODE, Integer.class),
                    new AbstractMap.SimpleImmutableEntry<>(ERROR_MESSAGE, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_ID, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_TYPE, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_NAME, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_VERSION, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_METHOD, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_CREATION, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_CREATOR_TENANT_DOMAIN, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(APPLICATION_ID, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(APPLICATION_NAME, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(APPLICATION_OWNER, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(REGION_ID, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(GATEWAY_TYPE, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(ORGANIZATION_ID, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(ENVIRONMENT_ID, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(PROXY_RESPONSE_CODE, Integer.class),
                    new AbstractMap.SimpleImmutableEntry<>(TARGET_RESPONSE_CODE, Integer.class),
                    new AbstractMap.SimpleImmutableEntry<>(RESPONSE_LATENCY, Long.class),
                    new AbstractMap.SimpleImmutableEntry<>(USER_IP, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(USER_AGENT_HEADER, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(PROPERTIES, LinkedHashMap.class))
            .collect(Collectors.toMap(Map.Entry::getKey, Map.Entry::getValue));
    private static final Map<String, Class> elkFaultEventSchema = Stream.of(
                    new AbstractMap.SimpleImmutableEntry<>(REQUEST_TIMESTAMP, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(CORRELATION_ID, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(KEY_TYPE, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(ERROR_TYPE, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(ERROR_CODE, Integer.class),
                    new AbstractMap.SimpleImmutableEntry<>(ERROR_MESSAGE, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_ID, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_TYPE, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_NAME, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_VERSION, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_METHOD, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_CREATION, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_CREATOR_TENANT_DOMAIN, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(APPLICATION_ID, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(APPLICATION_NAME, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(APPLICATION_OWNER, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(REGION_ID, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(GATEWAY_TYPE, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(PROXY_RESPONSE_CODE, Integer.class),
                    new AbstractMap.SimpleImmutableEntry<>(TARGET_RESPONSE_CODE, Integer.class),
                    new AbstractMap.SimpleImmutableEntry<>(USER_IP, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(USER_AGENT_HEADER, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(PROPERTIES, LinkedHashMap.class))
            .collect(Collectors.toMap(Map.Entry::getKey, Map.Entry::getValue));
    private static final Map<String, Class> choreoResponseEventSchema = Stream.of(
                    new AbstractMap.SimpleImmutableEntry<>(REQUEST_TIMESTAMP, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(CORRELATION_ID, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(KEY_TYPE, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_ID, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_TYPE, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_NAME, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_VERSION, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_CREATION, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_METHOD, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_RESOURCE_TEMPLATE, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_CREATOR_TENANT_DOMAIN, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_CONTEXT, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(DESTINATION, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(APPLICATION_ID, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(APPLICATION_NAME, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(APPLICATION_OWNER, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(REGION_ID, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(ORGANIZATION_ID, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(ENVIRONMENT_ID, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(GATEWAY_TYPE, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(USER_AGENT_HEADER, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(USER_NAME, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(PROXY_RESPONSE_CODE, Integer.class),
                    new AbstractMap.SimpleImmutableEntry<>(TARGET_RESPONSE_CODE, Integer.class),
                    new AbstractMap.SimpleImmutableEntry<>(RESPONSE_CACHE_HIT, Boolean.class),
                    new AbstractMap.SimpleImmutableEntry<>(RESPONSE_LATENCY, Long.class),
                    new AbstractMap.SimpleImmutableEntry<>(BACKEND_LATENCY, Long.class),
                    new AbstractMap.SimpleImmutableEntry<>(REQUEST_MEDIATION_LATENCY, Long.class),
                    new AbstractMap.SimpleImmutableEntry<>(RESPONSE_MEDIATION_LATENCY, Long.class),
                    new AbstractMap.SimpleImmutableEntry<>(USER_IP, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(PROPERTIES, LinkedHashMap.class))
            .collect(Collectors.toMap(Map.Entry::getKey, Map.Entry::getValue));
    private static final Map<String, Class> choreoFaultEventSchema = Stream.of(
                    new AbstractMap.SimpleImmutableEntry<>(REQUEST_TIMESTAMP, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(CORRELATION_ID, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(KEY_TYPE, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(ERROR_TYPE, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(ERROR_CODE, Integer.class),
                    new AbstractMap.SimpleImmutableEntry<>(ERROR_MESSAGE, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_ID, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_TYPE, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_NAME, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_VERSION, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_CREATION, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_METHOD, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_CREATOR_TENANT_DOMAIN, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(APPLICATION_ID, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(APPLICATION_NAME, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(APPLICATION_OWNER, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(REGION_ID, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(ORGANIZATION_ID, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(ENVIRONMENT_ID, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(GATEWAY_TYPE, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(PROXY_RESPONSE_CODE, Integer.class),
                    new AbstractMap.SimpleImmutableEntry<>(TARGET_RESPONSE_CODE, Integer.class),
                    new AbstractMap.SimpleImmutableEntry<>(RESPONSE_LATENCY, Long.class),
                    new AbstractMap.SimpleImmutableEntry<>(USER_IP, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(USER_AGENT_HEADER, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(PROPERTIES, LinkedHashMap.class))
            .collect(Collectors.toMap(Map.Entry::getKey, Map.Entry::getValue));

    private static final Map<String, Class> prometheusResponseEventSchema = Stream.of(
                    new AbstractMap.SimpleImmutableEntry<>(REQUEST_TIMESTAMP, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(CORRELATION_ID, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(KEY_TYPE, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_ID, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_TYPE, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_NAME, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_VERSION, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_CREATION, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_METHOD, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_RESOURCE_TEMPLATE, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_CREATOR_TENANT_DOMAIN, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_CONTEXT, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(DESTINATION, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(APPLICATION_ID, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(APPLICATION_NAME, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(APPLICATION_OWNER, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(REGION_ID, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(ORGANIZATION_ID, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(ENVIRONMENT_ID, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(GATEWAY_TYPE, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(USER_AGENT_HEADER, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(USER_NAME, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(PROXY_RESPONSE_CODE, Integer.class),
                    new AbstractMap.SimpleImmutableEntry<>(TARGET_RESPONSE_CODE, Integer.class),
                    new AbstractMap.SimpleImmutableEntry<>(RESPONSE_CACHE_HIT, Boolean.class),
                    new AbstractMap.SimpleImmutableEntry<>(RESPONSE_LATENCY, Long.class),
                    new AbstractMap.SimpleImmutableEntry<>(BACKEND_LATENCY, Long.class),
                    new AbstractMap.SimpleImmutableEntry<>(REQUEST_MEDIATION_LATENCY, Long.class),
                    new AbstractMap.SimpleImmutableEntry<>(RESPONSE_MEDIATION_LATENCY, Long.class),
                    new AbstractMap.SimpleImmutableEntry<>(USER_IP, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(PROPERTIES, LinkedHashMap.class))
            .collect(Collectors.toMap(Map.Entry::getKey, Map.Entry::getValue));

    private static final Map<String, Class> prometheusFaultEventSchema = Stream.of(
                    new AbstractMap.SimpleImmutableEntry<>(REQUEST_TIMESTAMP, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(CORRELATION_ID, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(KEY_TYPE, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(ERROR_TYPE, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(ERROR_CODE, Integer.class),
                    new AbstractMap.SimpleImmutableEntry<>(ERROR_MESSAGE, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_ID, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_TYPE, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_NAME, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_VERSION, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_CREATION, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_METHOD, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(API_CREATOR_TENANT_DOMAIN, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(APPLICATION_ID, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(APPLICATION_NAME, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(APPLICATION_OWNER, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(REGION_ID, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(ORGANIZATION_ID, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(ENVIRONMENT_ID, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(GATEWAY_TYPE, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(USER_AGENT_HEADER, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(USER_NAME, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(PROXY_RESPONSE_CODE, Integer.class),
                    new AbstractMap.SimpleImmutableEntry<>(TARGET_RESPONSE_CODE, Integer.class),
                    new AbstractMap.SimpleImmutableEntry<>(RESPONSE_CACHE_HIT, Boolean.class),
                    new AbstractMap.SimpleImmutableEntry<>(RESPONSE_LATENCY, Long.class),
                    new AbstractMap.SimpleImmutableEntry<>(BACKEND_LATENCY, Long.class),
                    new AbstractMap.SimpleImmutableEntry<>(REQUEST_MEDIATION_LATENCY, Long.class),
                    new AbstractMap.SimpleImmutableEntry<>(RESPONSE_MEDIATION_LATENCY, Long.class),
                    new AbstractMap.SimpleImmutableEntry<>(USER_IP, String.class),
                    new AbstractMap.SimpleImmutableEntry<>(PROPERTIES, LinkedHashMap.class))
            .collect(Collectors.toMap(Map.Entry::getKey, Map.Entry::getValue));
    private static final List<String> configProperties = new ArrayList<>();

    private GenericInputValidator() {
        //private constructor
    }

    public static org.wso2.apk.enforcer.analytics.publisher.reporter.GenericInputValidator getInstance() {

        return INSTANCE;
    }

    public Map<String, Class> getEventProperties(MetricSchema schema) {

        switch (schema) {
            case RESPONSE:
                return defaultResponseEventSchema;
            case ERROR:
                return defaultFaultEventSchema;
            case CHOREO_RESPONSE:
                return choreoResponseEventSchema;
            case CHOREO_ERROR:
                return choreoFaultEventSchema;
            case ELK_RESPONSE:
                return elkResponseEventSchema;
            case ELK_ERROR:
                return elkFaultEventSchema;
            case PROMETHEUS_RESPONSE:
                return prometheusResponseEventSchema;
            case PROMETHEUS_ERROR:
                return prometheusFaultEventSchema;
            default:
                return new HashMap<>();
        }
    }

    public List<String> getConfigProperties() {

        return configProperties;
    }
}
