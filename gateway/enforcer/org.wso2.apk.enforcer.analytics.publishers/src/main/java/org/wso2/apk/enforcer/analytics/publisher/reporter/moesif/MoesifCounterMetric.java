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
package org.wso2.apk.enforcer.analytics.publisher.reporter.moesif;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.wso2.apk.enforcer.analytics.publisher.exception.MetricReportingException;
import org.wso2.apk.enforcer.analytics.publisher.reporter.CounterMetric;
import org.wso2.apk.enforcer.analytics.publisher.reporter.GenericInputValidator;
import org.wso2.apk.enforcer.analytics.publisher.reporter.MetricEventBuilder;
import org.wso2.apk.enforcer.analytics.publisher.reporter.MetricSchema;

/**
 * Implementation of {@link CounterMetric} for Moesif Metric Reporter.
 */
public class MoesifCounterMetric implements CounterMetric {
    private static final Logger log = LoggerFactory.getLogger(MoesifCounterMetric.class);
    private EventQueue queue;
    private String name;
    private MetricSchema schema;

    public MoesifCounterMetric(String name, EventQueue queue, MetricSchema schema) {
        this.name = name;
        this.schema = schema;
        this.queue = queue;
    }


    @Override
    public int incrementCount(MetricEventBuilder metricEventBuilder) throws MetricReportingException {
        queue.put(metricEventBuilder);
        return 0;
    }

    @Override
    public String getName() {
        return this.name;
    }

    @Override
    public MetricSchema getSchema() {
        return this.schema;
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
            default:
                return new MoesifMetricEventBuilder(
                        GenericInputValidator.getInstance().getEventProperties(MetricSchema.RESPONSE));
            case ERROR:
                return new MoesifMetricEventBuilder(
                        GenericInputValidator.getInstance().getEventProperties(MetricSchema.ERROR));
            case CHOREO_RESPONSE:
                return new MoesifMetricEventBuilder(
                        GenericInputValidator.getInstance().getEventProperties(MetricSchema.CHOREO_RESPONSE));
            case CHOREO_ERROR:
                return new MoesifMetricEventBuilder(
                        GenericInputValidator.getInstance().getEventProperties(MetricSchema.CHOREO_ERROR));
        }
    }
}
