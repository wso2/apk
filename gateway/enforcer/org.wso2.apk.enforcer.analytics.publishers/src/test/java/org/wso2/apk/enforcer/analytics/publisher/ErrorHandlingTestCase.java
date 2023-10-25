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

import org.apache.logging.log4j.Level;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.apache.logging.log4j.core.LoggerContext;
import org.apache.logging.log4j.core.config.Configuration;
import org.testng.Assert;
import org.testng.annotations.Test;
import org.wso2.apk.enforcer.analytics.publisher.client.EventHubClient;
import org.wso2.apk.enforcer.analytics.publisher.exception.MetricCreationException;
import org.wso2.apk.enforcer.analytics.publisher.exception.MetricReportingException;
import org.wso2.apk.enforcer.analytics.publisher.reporter.CounterMetric;
import org.wso2.apk.enforcer.analytics.publisher.reporter.MetricEventBuilder;
import org.wso2.apk.enforcer.analytics.publisher.reporter.MetricReporter;
import org.wso2.apk.enforcer.analytics.publisher.reporter.MetricReporterFactory;
import org.wso2.apk.enforcer.analytics.publisher.reporter.MetricSchema;
import org.wso2.apk.enforcer.analytics.publisher.util.Constants;
import org.wso2.apk.enforcer.analytics.publisher.util.TestUtils;
import org.wso2.apk.enforcer.analytics.publisher.util.UnitTestAppender;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class ErrorHandlingTestCase {

    @Test
    public void testConnectionInvalidURL() throws MetricCreationException, MetricReportingException {

        Logger log = LogManager.getLogger(EventHubClient.class);
        LoggerContext context = LoggerContext.getContext(false);
        Configuration config = context.getConfiguration();
        UnitTestAppender appender = config.getAppender("UnitTestAppender");

        Map<String, String> configMap = new HashMap<>();
        configMap.put(Constants.AUTH_API_URL, "some_url");
        configMap.put(Constants.AUTH_API_TOKEN, "some_token");
        MetricReporter metricReporter = MetricReporterFactory.getInstance().createMetricReporter(configMap);
        CounterMetric metric = metricReporter.createCounterMetric("test-connection-counter", MetricSchema.RESPONSE);
        List<String> messages = new ArrayList<>(appender.getMessages());
        Assert.assertTrue(TestUtils.isContains(messages, "Unrecoverable error occurred when creating Eventhub "
                + "Client"), "Expected error hasn't logged in the "
                + "EventHubClientClass");
    }

    @Test(dependsOnMethods = {"testConnectionInvalidURL"})
    public void testConnectionUnavailability() throws Exception {

        Logger log = LogManager.getLogger(EventHubClient.class);
        LoggerContext context = LoggerContext.getContext(false);
        Configuration config = context.getConfiguration();
        UnitTestAppender appender = config.getAppender("UnitTestAppender");
        log.atLevel(Level.DEBUG);

        Map<String, String> configMap = new HashMap<>();
        configMap.put(Constants.AUTH_API_URL, "https://localhost:1234/non-existance");
        configMap.put(Constants.AUTH_API_TOKEN, "some_token");
        MetricReporterFactory factory = MetricReporterFactory.getInstance();
        factory.reset();
        MetricReporter metricReporter = factory.createMetricReporter(configMap);
        CounterMetric metric = metricReporter.createCounterMetric("test-connection-counter", MetricSchema.RESPONSE);
        List<String> messages = appender.getMessages();
        Thread.sleep(1000);
        Assert.assertTrue(TestUtils.isContains(messages, "Recoverable error occurred when creating Eventhub Client. "
                + "Retry attempts will be made"));
        Assert.assertTrue(TestUtils.isContains(messages, "Provided authentication endpoint "
                + "https://localhost:1234/non-existance is not "
                + "reachable."));
        MetricEventBuilder builder = metric.getEventBuilder();
        TestUtils.populateBuilder(builder);
        metric.incrementCount(builder);
        Thread.sleep(1000);
        messages = appender.getMessages();
        Assert.assertTrue(TestUtils.isContains(messages, "will be parked as EventHub Client is inactive."), "Thread "
                + "waiting log entry has not printed.");
    }
}
