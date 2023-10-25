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

package org.wso2.apk.enforcer.analytics.publisher.reporter.cloud;

import com.azure.core.amqp.AmqpRetryMode;
import com.azure.core.amqp.AmqpRetryOptions;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.wso2.apk.enforcer.analytics.publisher.client.EventHubClient;
import org.wso2.apk.enforcer.analytics.publisher.exception.MetricCreationException;
import org.wso2.apk.enforcer.analytics.publisher.reporter.AbstractMetricReporter;
import org.wso2.apk.enforcer.analytics.publisher.reporter.CounterMetric;
import org.wso2.apk.enforcer.analytics.publisher.reporter.MetricSchema;
import org.wso2.apk.enforcer.analytics.publisher.reporter.TimerMetric;
import org.wso2.apk.enforcer.analytics.publisher.util.Constants;

import java.time.Duration;
import java.util.List;
import java.util.Map;

/**
 * Choreo Metric Reporter Implementation. This implementation is responsible for sending analytics data into Choreo
 * cloud in a secure and reliable way.
 */
public class DefaultAnalyticsMetricReporter extends AbstractMetricReporter {

    private static final Logger log = LoggerFactory.getLogger(DefaultAnalyticsMetricReporter.class);
    protected EventQueue eventQueue;

    public DefaultAnalyticsMetricReporter(Map<String, String> properties) throws MetricCreationException {
        super(properties);
        int queueSize = Constants.DEFAULT_QUEUE_SIZE;
        int workerThreads = Constants.DEFAULT_WORKER_THREADS;
        int flushingDelay = Constants.DEFAULT_FLUSHING_DELAY;
        if (properties.get(Constants.QUEUE_SIZE) != null) {
            queueSize = Integer.parseInt(properties.get(Constants.QUEUE_SIZE));
        }
        if (properties.get(Constants.WORKER_THREAD_COUNT) != null) {
            workerThreads = Integer.parseInt(properties.get(Constants.WORKER_THREAD_COUNT));
        }
        if (properties.get(Constants.CLIENT_FLUSHING_DELAY) != null) {
            flushingDelay = Integer.parseInt(properties.get(Constants.CLIENT_FLUSHING_DELAY));
        }
        String authToken = properties.get(Constants.AUTH_API_TOKEN);
        String authEndpoint = properties.get(Constants.AUTH_API_URL);
        AmqpRetryOptions retryOptions = createRetryOptions(properties);
        EventHubClient client = new EventHubClient(authEndpoint, authToken, retryOptions, properties);
        eventQueue = new EventQueue(queueSize, workerThreads, client, flushingDelay);
    }

    private AmqpRetryOptions createRetryOptions(Map<String, String> properties) {
        int maxRetries = Constants.DEFAULT_MAX_RETRIES;
        int delay = Constants.DEFAULT_DELAY;
        int maxDelay = Constants.DEFAULT_MAX_DELAY;
        int tryTimeout = Constants.DEFAULT_TRY_TIMEOUT;
        AmqpRetryMode retryMode = AmqpRetryMode.FIXED;
        if (properties.get(Constants.EVENTHUB_CLIENT_MAX_RETRIES) != null) {
            int tempMaxRetries = Integer.parseInt(properties.get(Constants.EVENTHUB_CLIENT_MAX_RETRIES));
            if (tempMaxRetries > 0) {
                maxRetries = tempMaxRetries;
            } else {
                log.warn("Provided " + Constants.EVENTHUB_CLIENT_MAX_RETRIES + "value is less than 0 and not "
                                 + "acceptable. Hence using the default value.");
            }
        }
        if (properties.get(Constants.EVENTHUB_CLIENT_DELAY) != null) {
            int tempDelay = Integer.parseInt(properties.get(Constants.EVENTHUB_CLIENT_DELAY));
            if (tempDelay > 0) {
                delay = tempDelay;
            } else {
                log.warn("Provided " + Constants.EVENTHUB_CLIENT_DELAY + "value is less than 0 and not acceptable. "
                                 + "Hence using the default value.");
            }
        }
        if (properties.get(Constants.EVENTHUB_CLIENT_MAX_DELAY) != null) {
            int tempMaxDelay = Integer.parseInt(properties.get(Constants.EVENTHUB_CLIENT_MAX_DELAY));
            if (tempMaxDelay > 0) {
                maxDelay = tempMaxDelay;
            } else {
                log.warn("Provided " + Constants.EVENTHUB_CLIENT_MAX_DELAY + "value is less than 0 and not acceptable. "
                                 + "Hence using the default value.");
            }
        }
        if (properties.get(Constants.EVENTHUB_CLIENT_TRY_TIMEOUT) != null) {
            int tempTryTimeout = Integer.parseInt(properties.get(Constants.EVENTHUB_CLIENT_TRY_TIMEOUT));
            if (tempTryTimeout > 0) {
                tryTimeout = tempTryTimeout;
            } else {
                log.warn("Provided " + Constants.EVENTHUB_CLIENT_TRY_TIMEOUT + "value is less than 0 and not "
                                 + "acceptable. Hence using the default value.");
            }
        }
        if (properties.get(Constants.EVENTHUB_CLIENT_RETRY_MODE) != null) {
            String tempRetryMode = properties.get(Constants.EVENTHUB_CLIENT_RETRY_MODE);
            if (tempRetryMode.equals(Constants.FIXED)) {
                //do nothing
            } else if (tempRetryMode.equals(Constants.EXPONENTIAL)) {
                retryMode = AmqpRetryMode.EXPONENTIAL;
            } else {
                log.warn("Provided " + Constants.EVENTHUB_CLIENT_RETRY_MODE + "value is not supported. Hence will "
                                 + "using the default value.");
            }
        }
        return new AmqpRetryOptions()
                .setDelay(Duration.ofSeconds(delay))
                .setMaxRetries(maxRetries)
                .setMaxDelay(Duration.ofSeconds(maxDelay))
                .setTryTimeout(Duration.ofSeconds(tryTimeout))
                .setMode(retryMode);

    }

    @Override
    protected void validateConfigProperties(Map<String, String> properties) throws MetricCreationException {
        if (properties != null) {
            List<String> requiredProperties = DefaultInputValidator.getInstance().getConfigProperties();
            for (String property : requiredProperties) {
                if (properties.get(property) == null || properties.get(property).isEmpty()) {
                    throw new MetricCreationException(property + " is missing in config data");
                }
            }
        } else {
            throw new MetricCreationException("Configuration properties cannot be null");
        }
    }

    @Override
    protected CounterMetric createCounter(String name, MetricSchema schema) throws MetricCreationException {
        DefaultCounterMetric counterMetric = new DefaultCounterMetric(name, eventQueue, schema);
        return counterMetric;
    }

    @Override
    protected TimerMetric createTimer(String name) {
        return null;
    }

}
