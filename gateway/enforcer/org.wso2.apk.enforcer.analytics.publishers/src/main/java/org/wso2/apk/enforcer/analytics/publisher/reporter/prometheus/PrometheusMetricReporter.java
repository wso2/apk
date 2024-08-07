/*
 * Copyright (c) 2024, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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

package org.wso2.apk.enforcer.analytics.publisher.reporter.prometheus;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.wso2.apk.enforcer.analytics.publisher.exception.MetricCreationException;
import org.wso2.apk.enforcer.analytics.publisher.reporter.AbstractMetricReporter;
import org.wso2.apk.enforcer.analytics.publisher.reporter.CounterMetric;
import org.wso2.apk.enforcer.analytics.publisher.reporter.MetricSchema;
import org.wso2.apk.enforcer.analytics.publisher.reporter.TimerMetric;

import java.util.Map;

/**
 * Prometheus Metric Reporter class.
 */
public class PrometheusMetricReporter extends AbstractMetricReporter {
    private static final Logger log = LoggerFactory.getLogger(PrometheusMetricReporter.class);

    public PrometheusMetricReporter(Map<String, String> properties) throws MetricCreationException {
        super(properties);
        log.info("LogMetricReporter successfully initialized");
    }

    @Override
    protected void validateConfigProperties(Map<String, String> properties) throws MetricCreationException {
        //nothing to validate
    }

    @Override
    protected CounterMetric createCounter(String name, MetricSchema schema) {
        return new PrometheusCounterMetric(name, schema);
    }

    @Override
    protected TimerMetric createTimer(String name) {
        return null;
    }
}
