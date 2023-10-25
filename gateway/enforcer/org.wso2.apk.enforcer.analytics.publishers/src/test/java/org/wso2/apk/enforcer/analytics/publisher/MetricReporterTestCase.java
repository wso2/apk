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
import org.testng.annotations.Test;
import org.wso2.apk.enforcer.analytics.publisher.exception.MetricCreationException;
import org.wso2.apk.enforcer.analytics.publisher.exception.MetricReportingException;
import org.wso2.apk.enforcer.analytics.publisher.reporter.CounterMetric;
import org.wso2.apk.enforcer.analytics.publisher.reporter.MetricEventBuilder;
import org.wso2.apk.enforcer.analytics.publisher.reporter.MetricReporter;
import org.wso2.apk.enforcer.analytics.publisher.reporter.MetricReporterFactory;
import org.wso2.apk.enforcer.analytics.publisher.reporter.MetricSchema;
import org.wso2.apk.enforcer.analytics.publisher.util.Constants;
import org.wso2.apk.enforcer.analytics.publisher.util.TestUtils;

import java.util.HashMap;
import java.util.Map;

public class MetricReporterTestCase {

    private static final Logger log = LogManager.getLogger(MetricReporterTestCase.class);

    @Test(enabled = false, expectedExceptions = MetricCreationException.class)
    public void testMetricReporterCreationWithoutConfigs() throws MetricCreationException, MetricReportingException {

        try {
            createAndPublish(null, null);
        } catch (MetricCreationException e) {
            log.error(e);
            throw e;
        }
    }

    @Test(enabled = false, expectedExceptions = MetricCreationException.class, dependsOnMethods =
            {"testMetricReporterCreationWithoutConfigs"})
    public void testMetricReporterCreationWithMissingConfigs() throws MetricCreationException,
            MetricReportingException {

        Map<String, String> configs = new HashMap<>();
        configs.put("token.api.url", "localhost/token-api");
        configs.put("auth.api.url", "localhost/auth-api");
        configs.put("consumer.secret", "some_secret");
        createAndPublish(configs, null);
    }

    @Test(enabled = false, expectedExceptions = MetricCreationException.class, dependsOnMethods =
            {"testMetricReporterCreationWithMissingConfigs"})
    public void testMetricReporterCreationWithNullConfigs() throws MetricCreationException,
            MetricReportingException {

        Map<String, String> configs = new HashMap<>();
        configs.put("token.api.url", "localhost/token-api");
        configs.put("auth.api.url", "localhost/auth-api");
        configs.put("consumer.secret", "some_secret");
        configs.put("sas.token", "some_token");
        configs.put("consumer.key", null);
        createAndPublish(configs, null);
    }

    @Test(enabled = false, expectedExceptions = MetricCreationException.class, dependsOnMethods =
            {"testMetricReporterCreationWithNullConfigs"})
    public void testMetricReporterCreationWithEmptyConfigs() throws MetricCreationException,
            MetricReportingException {

        Map<String, String> configs = new HashMap<>();
        configs.put("token.api.url", "localhost/token-api");
        configs.put("auth.api.url", "localhost/auth-api");
        configs.put("sas.token", "some_token");
        configs.put("consumer.secret", "some_secret");
        configs.put("consumer.key", "");
        createAndPublish(configs, null);
    }

    @Test
    public void testCompleteFlow() throws MetricCreationException, MetricReportingException, InterruptedException {

        Map<String, String> configMap = new HashMap<>();
        String authURL = System.getenv(Constants.AUTH_API_URL);
        String authToken = System.getenv(Constants.AUTH_API_TOKEN);
        if (authToken != null && authURL != null) {
            configMap.put(Constants.AUTH_API_URL, authURL);
            configMap.put(Constants.AUTH_API_TOKEN, authToken);
        } else {
            return;
        }
        MetricReporter metricReporter = MetricReporterFactory.getInstance().createMetricReporter(configMap);
        CounterMetric metric = metricReporter.createCounterMetric("test-connection-counter", MetricSchema.RESPONSE);
        for (int i = 0; i < 5; i++) {
            MetricEventBuilder builder = metric.getEventBuilder();
            TestUtils.populateBuilder(builder);
            metric.incrementCount(builder);
        }
        Thread.sleep(2000);
        //Assertions will be done after mocking eventhub client
    }

    /**
     * Helper method to create and publish metrics for testing purposes
     *
     * @param configs Config map
     * @param event   Event map
     * @throws MetricCreationException  Thrown if config properties are missing
     * @throws MetricReportingException Thrown if event properties are missing
     */
    private void createAndPublish(Map<String, String> configs, Map<String, String> event)
            throws MetricCreationException, MetricReportingException {

        MetricReporter metricReporter = MetricReporterFactory.getInstance().createMetricReporter(configs);
        CounterMetric counterMetric = metricReporter.createCounterMetric("apim.response", MetricSchema.RESPONSE);
        MetricEventBuilder builder = counterMetric.getEventBuilder();

        counterMetric.incrementCount(builder);
    }
}
