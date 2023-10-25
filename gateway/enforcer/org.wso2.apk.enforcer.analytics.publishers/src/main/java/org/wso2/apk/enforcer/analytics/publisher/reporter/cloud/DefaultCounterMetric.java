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

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.wso2.apk.enforcer.analytics.publisher.client.ClientStatus;
import org.wso2.apk.enforcer.analytics.publisher.exception.MetricCreationException;
import org.wso2.apk.enforcer.analytics.publisher.reporter.CounterMetric;
import org.wso2.apk.enforcer.analytics.publisher.reporter.MetricEventBuilder;
import org.wso2.apk.enforcer.analytics.publisher.reporter.MetricSchema;

import java.util.concurrent.atomic.AtomicInteger;

/**
 * Implementation of {@link CounterMetric} for Choroe Metric Reporter.
 */
public class DefaultCounterMetric implements CounterMetric {
    private static final Logger log = LoggerFactory.getLogger(DefaultCounterMetric.class);
    private String name;
    private EventQueue queue;
    private MetricSchema schema;
    private ClientStatus status;
    private final AtomicInteger failureCount;

    public DefaultCounterMetric(String name, EventQueue queue, MetricSchema schema) throws MetricCreationException {
        //Constructor should be made protected. Keeping public till testing plan is finalized
        this.name = name;
        this.queue = queue;
        if (schema == MetricSchema.ERROR || schema == MetricSchema.RESPONSE
                || schema == MetricSchema.CHOREO_ERROR || schema == MetricSchema.CHOREO_RESPONSE) {
            this.schema = schema;
        } else {
            throw new MetricCreationException("Default Counter Metric only supports " + MetricSchema.RESPONSE + ", "
                    + ", " + MetricSchema.ERROR + ", " + MetricSchema.CHOREO_RESPONSE + " and "
                    + MetricSchema.CHOREO_RESPONSE + " types.");
        }
        this.status = queue.getClient().getStatus();
        this.failureCount = new AtomicInteger(0);
    }

    @Override
    public String getName() {
        return name;
    }

    @Override public MetricSchema getSchema() {
        return schema;
    }

    @Override
    public int incrementCount(MetricEventBuilder builder) {
        if (!(status == ClientStatus.NOT_CONNECTED)) {
            queue.put(builder);
        } else {
            if (failureCount.incrementAndGet() % 1000 == 0) {
                log.error("Eventhub client is not connected. " + failureCount.incrementAndGet() + " events dropped so "
                                  + "far. Please correct your configuration and restart the instance.");
            }
        }
        return 0;
    }

    /**
     * Returns Event Builder used for this CounterMetric. Depending on the schema different types of builders will be
     * returned.
     *
     * @return {@link MetricEventBuilder} for this {@link CounterMetric}
     */
    @Override
    public MetricEventBuilder getEventBuilder() {
        switch (schema) {
            case RESPONSE:
                return new DefaultResponseMetricEventBuilder();
            case ERROR:
                return new DefaultFaultMetricEventBuilder();
            case CHOREO_RESPONSE:
                return new DefaultChoreoResponseMetricEventBuilder();
            case CHOREO_ERROR:
                return new DefaultChoreoFaultMetricEventBuilder();
            default:
                // will not happen
                return null;
        }
    }
}
