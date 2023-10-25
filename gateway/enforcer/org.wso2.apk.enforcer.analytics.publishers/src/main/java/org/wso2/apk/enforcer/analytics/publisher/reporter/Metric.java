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

/**
 * Base class for {@link CounterMetric} and {@link TimerMetric}.
 */
public interface Metric {
    /**
     * Returns name of the metric. Name is unique per given Reporter
     *
     * @return Name of the metric
     */
    public String getName();

    /**
     * Method to get schema name.
     *
     * @return Schema name of this Metric
     */
    public MetricSchema getSchema();

    /**
     * Returns the Event Builder which can produce an event conforming to the schema of the {@link Metric}. Users
     * should get a builder instance though here, populate the relevant fields and return back to Metric class.
     * @return MetricEventBuilder conforming to schema of Metric
     */
    public MetricEventBuilder getEventBuilder();
}
