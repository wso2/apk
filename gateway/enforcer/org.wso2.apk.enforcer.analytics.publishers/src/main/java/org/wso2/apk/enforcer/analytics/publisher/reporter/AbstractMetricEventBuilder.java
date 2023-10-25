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

import java.util.Map;

/**
 * Abstract implementation of {@link MetricEventBuilder}. Ensures that subclasses perform validations when adding
 * elements.
 */
public abstract class AbstractMetricEventBuilder implements MetricEventBuilder {
    @Override
    public Map<String, Object> build() throws MetricReportingException {
        if (validate()) {
            return buildEvent();
        }
        throw new MetricReportingException("Validation failure occurred when building the event");
    }

    /**
     * Process the added data and return as a flat {@link Map}.
     *
     * @return Map representing attributes of Metric Event
     */
    protected abstract Map<String, Object> buildEvent();
}
