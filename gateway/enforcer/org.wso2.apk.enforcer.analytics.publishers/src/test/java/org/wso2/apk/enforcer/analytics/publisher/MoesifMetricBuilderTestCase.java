/*
 * Copyright (c) 2023, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 *
 * WSO2 Inc. licenses this file to you under the Apache License,
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

package org.wso2.apk.enforcer.analytics.publisher;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.testng.Assert;
import org.testng.annotations.BeforeMethod;
import org.testng.annotations.Test;
import org.wso2.apk.enforcer.analytics.publisher.exception.MetricCreationException;
import org.wso2.apk.enforcer.analytics.publisher.exception.MetricReportingException;
import org.wso2.apk.enforcer.analytics.publisher.reporter.MetricEventBuilder;
import org.wso2.apk.enforcer.analytics.publisher.reporter.MetricSchema;
import org.wso2.apk.enforcer.analytics.publisher.reporter.moesif.EventQueue;
import org.wso2.apk.enforcer.analytics.publisher.reporter.moesif.MoesifCounterMetric;
import org.wso2.apk.enforcer.analytics.publisher.reporter.moesif.util.MoesifMicroserviceConstants;
import org.wso2.apk.enforcer.analytics.publisher.retriever.MoesifKeyRetrieverChoreoClient;
import org.wso2.apk.enforcer.analytics.publisher.util.Constants;

import java.time.Clock;
import java.time.OffsetDateTime;
import java.util.HashMap;
import java.util.LinkedHashMap;
import java.util.Map;

public class MoesifMetricBuilderTestCase {

    private static final Logger log = LoggerFactory.getLogger(MoesifMetricBuilderTestCase.class);

    private MetricEventBuilder builder;

    @BeforeMethod
    public void createBuilder() throws MetricCreationException {

        Map<String, String> properties = new HashMap<>();
        properties.put(MoesifMicroserviceConstants.MS_USERNAME_CONFIG_KEY, "some_username");
        properties.put(MoesifMicroserviceConstants.MS_PWD_CONFIG_KEY, "some_password");
        properties.put(MoesifMicroserviceConstants.MOESIF_PROTOCOL_WITH_FQDN_KEY, "abcde");
        properties.put(MoesifMicroserviceConstants.MOESIF_MS_VERSIONING_KEY, "aaa");
        MoesifKeyRetrieverChoreoClient keyRetriever =
                MoesifKeyRetrieverChoreoClient.getInstance(properties);
        EventQueue queue = new EventQueue(100, 1, keyRetriever);
        MoesifCounterMetric metric =
                new MoesifCounterMetric("test.builder.metric", queue, MetricSchema.CHOREO_RESPONSE);
        builder = metric.getEventBuilder();
    }

    @Test(expectedExceptions = MetricReportingException.class)
    public void testMissingAttributes() throws MetricCreationException, MetricReportingException {

        builder.addAttribute(Constants.REQUEST_TIMESTAMP, System.currentTimeMillis())
                .addAttribute(Constants.CORRELATION_ID, "1234-4567")
                .addAttribute(Constants.KEY_TYPE, "prod")
                .addAttribute(Constants.API_ID, "9876-54f1")
                .addAttribute(Constants.API_NAME, "PizzaShack")
                .addAttribute(Constants.API_VERSION, "1.0.0")
                .addAttribute(Constants.API_CREATION, "admin")
                .addAttribute(Constants.API_METHOD, "POST")
                .addAttribute(Constants.API_METHOD, "POST")
                .addAttribute(Constants.API_RESOURCE_TEMPLATE, "/resource/{value}")
                .addAttribute(Constants.API_CREATOR_TENANT_DOMAIN, "carbon.super")
                .addAttribute(Constants.DESTINATION, "localhost:8080")
                .addAttribute(Constants.APPLICATION_ID, "3445-6778")
                .addAttribute(Constants.APPLICATION_NAME, "default")
                .addAttribute(Constants.APPLICATION_OWNER, "admin")
                .addAttribute(Constants.REGION_ID, "NA")
                .addAttribute(Constants.GATEWAY_TYPE, "Synapse")
                .addAttribute(Constants.USER_AGENT, "Mozilla")
                .addAttribute(Constants.PROXY_RESPONSE_CODE, 401)
                .addAttribute(Constants.TARGET_RESPONSE_CODE, "someString")
                .addAttribute(Constants.RESPONSE_CACHE_HIT, true)
                .addAttribute(Constants.RESPONSE_LATENCY, 2000)
                .addAttribute(Constants.BACKEND_LATENCY, 3000)
                .addAttribute(Constants.REQUEST_MEDIATION_LATENCY, "1000")
                .addAttribute(Constants.RESPONSE_MEDIATION_LATENCY, 1000)
                .addAttribute(Constants.USER_IP, "127.0.0.1")
                .build();
    }

    @Test(expectedExceptions = MetricReportingException.class)
    public void testAttributesWithInvalidTypes() throws MetricCreationException, MetricReportingException {

        builder.addAttribute(Constants.REQUEST_TIMESTAMP, System.currentTimeMillis())
                .addAttribute(Constants.CORRELATION_ID, "1234-4567")
                .addAttribute(Constants.ORGANIZATION_ID, "wso2.com")
                .addAttribute(Constants.ENVIRONMENT_ID, "abcd-1234-5678-efgh")
                .addAttribute(Constants.KEY_TYPE, "prod")
                .addAttribute(Constants.API_ID, "9876-54f1")
                .addAttribute(Constants.API_NAME, "PizzaShack")
                .addAttribute(Constants.API_VERSION, "1.0.0")
                .addAttribute(Constants.API_CREATION, "admin")
                .addAttribute(Constants.API_METHOD, "POST")
                .addAttribute(Constants.API_METHOD, "POST")
                .addAttribute(Constants.API_CONTEXT, "/v1/")
                .addAttribute(Constants.API_RESOURCE_TEMPLATE, "/resource/{value}")
                .addAttribute(Constants.API_CREATOR_TENANT_DOMAIN, "carbon.super")
                .addAttribute(Constants.DESTINATION, "localhost:8080")
                .addAttribute(Constants.APPLICATION_ID, "3445-6778")
                .addAttribute(Constants.APPLICATION_NAME, "default")
                .addAttribute(Constants.APPLICATION_OWNER, "admin")
                .addAttribute(Constants.REGION_ID, "NA")
                .addAttribute(Constants.GATEWAY_TYPE, "Synapse")
                .addAttribute(Constants.USER_AGENT, "Mozilla")
                .addAttribute(Constants.USER_NAME, "admin")
                .addAttribute(Constants.PROXY_RESPONSE_CODE, 401)
                .addAttribute(Constants.TARGET_RESPONSE_CODE, "someString")
                .addAttribute(Constants.RESPONSE_CACHE_HIT, true)
                .addAttribute(Constants.RESPONSE_LATENCY, 2000)
                .addAttribute(Constants.BACKEND_LATENCY, 3000)
                .addAttribute(Constants.REQUEST_MEDIATION_LATENCY, "1000")
                .addAttribute(Constants.RESPONSE_MEDIATION_LATENCY, 1000)
                .addAttribute(Constants.USER_IP, "127.0.0.1")
                .build();
    }

    @Test
    public void testMetricBuilder() throws MetricCreationException, MetricReportingException {

        String uaString = "Mozilla/5.0 (iPhone; CPU iPhone OS 5_1_1 like Mac OS X) AppleWebKit/534.46 (KHTML, "
                + "like Gecko) Version/5.1 Mobile/9B206 Safari/7534.48.3";

        LinkedHashMap values = new LinkedHashMap<>();
        values.put("x-original-gw-url", "foo");
        LinkedHashMap properties = new LinkedHashMap();
        properties.put("properties", values);

        Map<String, Object> eventMap = builder
                .addAttribute(Constants.REQUEST_TIMESTAMP, OffsetDateTime.now(Clock.systemUTC()).toString())
                .addAttribute(Constants.CORRELATION_ID, "1234-4567")
                .addAttribute(Constants.ORGANIZATION_ID, "wso2.com")
                .addAttribute(Constants.ENVIRONMENT_ID, "abcd-1234-5678-efgh")
                .addAttribute(Constants.KEY_TYPE, "prod")
                .addAttribute(Constants.API_ID, "9876-54f1")
                .addAttribute(Constants.API_TYPE, "HTTP")
                .addAttribute(Constants.API_NAME, "PizzaShack")
                .addAttribute(Constants.API_VERSION, "1.0.0")
                .addAttribute(Constants.API_CREATION, "admin")
                .addAttribute(Constants.API_METHOD, "POST")
                .addAttribute(Constants.API_CONTEXT, "/v1/")
                .addAttribute(Constants.USER_NAME, "admin")
                .addAttribute(Constants.API_RESOURCE_TEMPLATE, "/resource/{value}")
                .addAttribute(Constants.API_CREATOR_TENANT_DOMAIN, "carbon.super")
                .addAttribute(Constants.DESTINATION, "localhost:8080")
                .addAttribute(Constants.APPLICATION_ID, "3445-6778")
                .addAttribute(Constants.APPLICATION_NAME, "default")
                .addAttribute(Constants.APPLICATION_OWNER, "admin")
                .addAttribute(Constants.REGION_ID, "NA")
                .addAttribute(Constants.GATEWAY_TYPE, "Synapse")
                .addAttribute(Constants.USER_AGENT_HEADER, uaString)
                .addAttribute(Constants.PROXY_RESPONSE_CODE, 401)
                .addAttribute(Constants.TARGET_RESPONSE_CODE, 401)
                .addAttribute(Constants.RESPONSE_CACHE_HIT, true)
                .addAttribute(Constants.RESPONSE_LATENCY, 2000L)
                .addAttribute(Constants.BACKEND_LATENCY, 3000L)
                .addAttribute(Constants.REQUEST_MEDIATION_LATENCY, 1000L)
                .addAttribute(Constants.RESPONSE_MEDIATION_LATENCY, 1000L)
                .addAttribute(Constants.USER_IP, "127.0.0.1")
                .addAttribute(Constants.PROPERTIES, properties)
                .build();

        Assert.assertFalse(eventMap.isEmpty());
        // We expect only 30 attributes in Moesif scenario unlike in choreo scenario.
        // In choreo scenario we parse user agent header,
        // and build one additional attribute.
        Assert.assertEquals(eventMap.size(), 31, "Some attributes are missing from the resulting event map");
        Assert.assertEquals(eventMap.get(Constants.ORGANIZATION_ID), "wso2.com",
                "Organization ID should be wso2.com");
    }
}
