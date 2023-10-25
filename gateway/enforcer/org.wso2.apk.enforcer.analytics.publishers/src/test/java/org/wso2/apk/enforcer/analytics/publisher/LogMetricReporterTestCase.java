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

import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.apache.logging.log4j.core.LoggerContext;
import org.apache.logging.log4j.core.config.Configuration;
import org.testng.Assert;
import org.testng.annotations.Test;
import org.wso2.apk.enforcer.analytics.publisher.exception.MetricCreationException;
import org.wso2.apk.enforcer.analytics.publisher.exception.MetricReportingException;
import org.wso2.apk.enforcer.analytics.publisher.reporter.CounterMetric;
import org.wso2.apk.enforcer.analytics.publisher.reporter.MetricEventBuilder;
import org.wso2.apk.enforcer.analytics.publisher.reporter.MetricReporter;
import org.wso2.apk.enforcer.analytics.publisher.reporter.MetricReporterFactory;
import org.wso2.apk.enforcer.analytics.publisher.reporter.log.LogCounterMetric;
import org.wso2.apk.enforcer.analytics.publisher.util.TestUtils;
import org.wso2.apk.enforcer.analytics.publisher.util.UnitTestAppender;

import java.util.List;

public class LogMetricReporterTestCase {

    private static final Logger log = LogManager.getLogger(LogMetricReporterTestCase.class);

    @Test
    public void testLogMetricReporter() throws MetricCreationException, MetricReportingException {

        log.info("Running log metric test case");
        Logger log = LogManager.getLogger(LogCounterMetric.class);
        LoggerContext context = LoggerContext.getContext(false);
        Configuration config = context.getConfiguration();
        UnitTestAppender appender = config.getAppender("UnitTestAppender");

        MetricReporter metricReporter = MetricReporterFactory.getInstance().createMetricReporter(
                "org.wso2.apk.enforcer.analytics.publisher.reporter.log.LogMetricReporter", null);
        CounterMetric metric = metricReporter.createCounterMetric("testCounter", null);
        MetricEventBuilder builder = metric.getEventBuilder();
        builder.addAttribute("attribute1", "value1").addAttribute("attribute2", "value2").addAttribute("attribute3",
                "value3");
        metric.incrementCount(builder);

        List<String> messages = appender.getMessages();

        Assert.assertTrue(TestUtils.isContains(messages, "testCounter"), "Metric name is not properly logged");
        Assert.assertTrue(TestUtils.isContains(messages, "attribute1=value1"), "Metric attribute is not properly "
                + "logged");
        Assert.assertTrue(TestUtils.isContains(messages, "attribute2=value2"),
                "Metric attribute is not properly logged");
        Assert.assertTrue(TestUtils.isContains(messages, "attribute3=value3"),
                "Metric attribute is not properly logged");
    }

    @Test(expectedExceptions = MetricReportingException.class, dependsOnMethods = {"testLogMetricReporter"})
    public void testLogMetricReporterWithInvalidAttributes() throws MetricCreationException, MetricReportingException {

        MetricReporter metricReporter = MetricReporterFactory.getInstance().createMetricReporter(
                "org.wso2.apk.enforcer.analytics.publisher.reporter.log.LogMetricReporter", null);
        CounterMetric metric = metricReporter.createCounterMetric("testCounter", null);
        MetricEventBuilder builder = metric.getEventBuilder();
        builder.addAttribute("attribute1", 123).addAttribute("attribute2", "value2").addAttribute("attribute3",
                "value3");
        metric.incrementCount(builder);
    }
}
