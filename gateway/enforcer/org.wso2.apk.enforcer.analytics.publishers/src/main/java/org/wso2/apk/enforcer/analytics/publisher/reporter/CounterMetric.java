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

import org.wso2.apk.enforcer.analytics.publisher.exception.MetricReportingException;

/**
 * Interface for Counter Metric.
 */
public interface CounterMetric extends Metric {
    /**
     * method to increment the count of the metric. Associated properties should be passed along with the method
     * invocation.
     *
     * @param builder {@link MetricEventBuilder} of the this Metric
     * @return current counter value
     * @throws MetricReportingException Exception will be thrown if expected properties are not present
     */
    public int incrementCount(MetricEventBuilder builder) throws MetricReportingException;
}
