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

package org.wso2.apk.enforcer.analytics.publisher;

import com.azure.core.amqp.AmqpRetryMode;
import com.azure.core.amqp.AmqpRetryOptions;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.testng.Assert;
import org.testng.annotations.BeforeMethod;
import org.testng.annotations.Test;
import org.wso2.apk.enforcer.analytics.publisher.client.EventHubClient;
import org.wso2.apk.enforcer.analytics.publisher.exception.MetricCreationException;
import org.wso2.apk.enforcer.analytics.publisher.exception.MetricReportingException;
import org.wso2.apk.enforcer.analytics.publisher.reporter.MetricEventBuilder;
import org.wso2.apk.enforcer.analytics.publisher.reporter.MetricSchema;
import org.wso2.apk.enforcer.analytics.publisher.reporter.cloud.DefaultCounterMetric;
import org.wso2.apk.enforcer.analytics.publisher.reporter.cloud.EventQueue;
import org.wso2.apk.enforcer.analytics.publisher.util.Constants;

import java.time.Clock;
import java.time.Duration;
import java.time.OffsetDateTime;
import java.util.HashMap;
import java.util.Map;

public class DefaultChoreoFaultMetricBuilderTestCase {

    private static final Logger log = LoggerFactory.getLogger(DefaultChoreoFaultMetricBuilderTestCase.class);
    private MetricEventBuilder builder;

    @BeforeMethod
    public void createBuilder() throws MetricCreationException {

        AmqpRetryOptions retryOptions = new AmqpRetryOptions()
                .setDelay(Duration.ofSeconds(30))
                .setMaxRetries(2)
                .setMaxDelay(Duration.ofSeconds(120))
                .setTryTimeout(Duration.ofSeconds(30))
                .setMode(AmqpRetryMode.FIXED);
        EventHubClient client = new EventHubClient("some_endpoint", "some_token", retryOptions,
                new HashMap<>());
        EventQueue queue = new EventQueue(100, 1, client, 10);
        DefaultCounterMetric metric = new DefaultCounterMetric("test.builder.metric", queue,
                MetricSchema.CHOREO_ERROR);
        builder = metric.getEventBuilder();
    }

    @Test(expectedExceptions = MetricReportingException.class)
    public void testMissingAttributes() throws MetricCreationException, MetricReportingException {

        builder.addAttribute(Constants.REQUEST_TIMESTAMP, System.currentTimeMillis())
                .addAttribute(Constants.CORRELATION_ID, "1234-4567")
                .addAttribute(Constants.KEY_TYPE, "prod")
                .addAttribute(Constants.ERROR_TYPE, "backend")
                .addAttribute(Constants.ERROR_CODE, 401)
                .addAttribute(Constants.ERROR_MESSAGE, "Authentication Error")
                .addAttribute(Constants.API_ID, "9876-54f1")
                .addAttribute(Constants.API_NAME, "PizzaShack")
                .addAttribute(Constants.API_VERSION, "1.0.0")
                .addAttribute(Constants.API_CREATION, "admin")
                .addAttribute(Constants.API_METHOD, "POST")
                .addAttribute(Constants.API_CREATOR_TENANT_DOMAIN, "carbon.super")
                .addAttribute(Constants.APPLICATION_ID, "3445-6778")
                .addAttribute(Constants.APPLICATION_NAME, "default")
                .addAttribute(Constants.APPLICATION_OWNER, "admin")
                .addAttribute(Constants.REGION_ID, "NA")
                .addAttribute(Constants.GATEWAY_TYPE, "Synapse")
                .addAttribute(Constants.PROXY_RESPONSE_CODE, 401)
                .addAttribute(Constants.TARGET_RESPONSE_CODE, "someString")
                .build();
    }

    @Test(expectedExceptions = MetricReportingException.class)
    public void testAttributesWithInvalidTypes() throws MetricCreationException, MetricReportingException {

        builder.addAttribute(Constants.REQUEST_TIMESTAMP, System.currentTimeMillis())
                .addAttribute(Constants.ORGANIZATION_ID, "wso2.com")
                .addAttribute(Constants.ENVIRONMENT_ID, "abcd-1234-5678-efgh")
                .addAttribute(Constants.CORRELATION_ID, "1234-4567")
                .addAttribute(Constants.KEY_TYPE, "prod")
                .addAttribute(Constants.ERROR_TYPE, "backend")
                .addAttribute(Constants.ERROR_CODE, 401)
                .addAttribute(Constants.ERROR_MESSAGE, "Authentication Error")
                .addAttribute(Constants.API_ID, "9876-54f1")
                .addAttribute(Constants.API_NAME, "PizzaShack")
                .addAttribute(Constants.API_VERSION, "1.0.0")
                .addAttribute(Constants.API_CREATION, "admin")
                .addAttribute(Constants.API_METHOD, "POST")
                .addAttribute(Constants.API_CREATOR_TENANT_DOMAIN, "carbon.super")
                .addAttribute(Constants.APPLICATION_ID, "3445-6778")
                .addAttribute(Constants.APPLICATION_NAME, "default")
                .addAttribute(Constants.APPLICATION_OWNER, "admin")
                .addAttribute(Constants.REGION_ID, "NA")
                .addAttribute(Constants.GATEWAY_TYPE, "Synapse")
                .addAttribute(Constants.PROXY_RESPONSE_CODE, 401)
                .addAttribute(Constants.TARGET_RESPONSE_CODE, "someString")
                .build();
    }

    @Test
    public void testMetricBuilder() throws MetricCreationException, MetricReportingException {

        Map<String, Object> eventMap = builder
                .addAttribute(Constants.REQUEST_TIMESTAMP, OffsetDateTime.now(Clock.systemUTC()).toString())
                .addAttribute(Constants.CORRELATION_ID, "1234-4567")
                .addAttribute(Constants.ORGANIZATION_ID, "wso2.com")
                .addAttribute(Constants.ENVIRONMENT_ID, "abcd-1234-5678-efgh")
                .addAttribute(Constants.KEY_TYPE, "prod")
                .addAttribute(Constants.ERROR_TYPE, "backend")
                .addAttribute(Constants.ERROR_CODE, 401)
                .addAttribute(Constants.ERROR_MESSAGE, "Authentication Error")
                .addAttribute(Constants.API_ID, "9876-54f1")
                .addAttribute(Constants.API_TYPE, "HTTP")
                .addAttribute(Constants.API_NAME, "PizzaShack")
                .addAttribute(Constants.API_VERSION, "1.0.0")
                .addAttribute(Constants.API_CREATION, "admin")
                .addAttribute(Constants.API_METHOD, "POST")
                .addAttribute(Constants.API_CREATOR_TENANT_DOMAIN, "carbon.super")
                .addAttribute(Constants.APPLICATION_ID, "3445-6778")
                .addAttribute(Constants.APPLICATION_NAME, "default")
                .addAttribute(Constants.APPLICATION_OWNER, "admin")
                .addAttribute(Constants.REGION_ID, "NA")
                .addAttribute(Constants.GATEWAY_TYPE, "Synapse")
                .addAttribute(Constants.PROXY_RESPONSE_CODE, 401)
                .addAttribute(Constants.TARGET_RESPONSE_CODE, 401)
                .build();

        Assert.assertFalse(eventMap.isEmpty());
        Assert.assertEquals(eventMap.size(), 23, "Some attributes are missing from the resulting event map");
        Assert.assertEquals(eventMap.get(Constants.EVENT_TYPE), "fault", "Event type should be set to fault");
        Assert.assertEquals(eventMap.get(Constants.ORGANIZATION_ID), "wso2.com",
                "Organization ID should be wso2.com");
    }
}
