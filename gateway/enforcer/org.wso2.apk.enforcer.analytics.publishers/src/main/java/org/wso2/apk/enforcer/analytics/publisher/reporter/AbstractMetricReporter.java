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

package org.wso2.apk.enforcer.analytics.publisher.reporter;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.wso2.apk.enforcer.analytics.publisher.exception.MetricCreationException;

import java.util.HashMap;
import java.util.Map;

/**
 * Implementation of {@link MetricReporter}. All concrete implementations should extend this class. Validations and
 * metric type creation is enforced using this abstract class
 */
public abstract class AbstractMetricReporter implements MetricReporter {
    private static final Logger log = LoggerFactory.getLogger(AbstractMetricReporter.class);
    private final Map<String, String> properties;
    private Map<String, Metric> metricRegistry;

    protected AbstractMetricReporter(Map<String, String> properties) throws MetricCreationException {
        this.properties = properties;
        metricRegistry = new HashMap<>();
        validateConfigProperties(properties);
    }

    /**
     * Method to validate the configuration properties. Config properties are accepted as a map to increase
     * extendability. Hence this method is responsible to sanitize the map
     *
     * @param properties Configuration properties needed by the implementation
     * @throws MetricCreationException Exception will be throw is any of the required fields are missing or no in the
     *                                 expected format
     */
    protected abstract void validateConfigProperties(Map<String, String> properties) throws MetricCreationException;

    public Map<String, String> getConfiguration() {
        return properties;
    }

    @Override
    public CounterMetric createCounterMetric(String name, MetricSchema schema) throws MetricCreationException {
        Metric metric = metricRegistry.get(name);
        if (metric == null) {
            synchronized (this) {
                if (metricRegistry.get(name) == null) {
                    metric = createCounter(name, schema);
                    metricRegistry.put(name, metric);
                } else {
                    metric = metricRegistry.get(name);
                }
            }
        } else if (!(metric instanceof CounterMetric)) {
            throw new MetricCreationException("Timer Metric with the same name already exists. Please use a different"
                                                      + " name");
        } else if (metric.getSchema() != schema) {
            throw new MetricCreationException("Counter Metric with the same name but different schema already exists."
                                                      + " Please use a different name");
        }
        return (CounterMetric) metric;
    }

    protected abstract CounterMetric createCounter(String name, MetricSchema schema) throws MetricCreationException;

    @Override
    public TimerMetric createTimerMetric(String name) {
        Metric metric = metricRegistry.get(name);
        if (metric == null) {
            synchronized (this) {
                if (metricRegistry.get(name) == null) {
                    metric = createTimer(name);
                    metricRegistry.put(name, metric);
                } else {
                    metric = metricRegistry.get(name);
                }
            }
        }
        if (!(metric instanceof TimerMetric)) {
            log.error("Counter Metric with the same name already exists. Please use a different name");
            return null;
        } else {
            return (TimerMetric) metric;
        }
    }

    protected abstract TimerMetric createTimer(String name);

}
