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

import org.wso2.apk.enforcer.analytics.publisher.exception.MetricCreationException;

import java.util.Map;

/**
 * Base interface for all Metric Reporter implementations. {@link AbstractMetricReporter} will implement this
 * interface and any concrete implementations should extend {@link AbstractMetricReporter}
 */
public interface MetricReporter {

    /**
     * Create and return {@link CounterMetric} to instrument collected metrics.
     *
     * @param name   Name of the metric
     * @param schema Metric schema
     * @return {@link CounterMetric}
     * @throws MetricCreationException if error occurred when creating CounterMetric
     */
    CounterMetric createCounterMetric(String name, MetricSchema schema) throws MetricCreationException;

    /**
     * Create and return {@link TimerMetric} to instrument collected metrics.
     *
     * @param name Name of the metric
     * @return {@link TimerMetric}
     */
    TimerMetric createTimerMetric(String name);

    /**
     * Returns the currently set configurations. Setting configurations will be mandated through constructor
     *
     * @return {@link Map} representing current configurations
     */
    Map<String, String> getConfiguration();

}
