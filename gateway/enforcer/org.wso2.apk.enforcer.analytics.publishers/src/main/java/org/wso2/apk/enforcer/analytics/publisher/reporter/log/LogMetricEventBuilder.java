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

package org.wso2.apk.enforcer.analytics.publisher.reporter.log;

import org.wso2.apk.enforcer.analytics.publisher.exception.MetricReportingException;
import org.wso2.apk.enforcer.analytics.publisher.reporter.AbstractMetricEventBuilder;
import org.wso2.apk.enforcer.analytics.publisher.reporter.MetricEventBuilder;

import java.util.HashMap;
import java.util.Map;

/**
 * Event builder for log Metric Reporter.
 */
public class LogMetricEventBuilder extends AbstractMetricEventBuilder {
    private Map<String, Object> eventMap = new HashMap<>();

    @Override
    protected Map<String, Object> buildEvent() {
        return eventMap;
    }

    @Override
    public boolean validate() throws MetricReportingException {
        for (Object value : eventMap.values()) {
            if (!(value instanceof String)) {
                throw new MetricReportingException("Only attributes of type String is supported");
            }
        }
        return true;
    }

    @Override public MetricEventBuilder addAttribute(String key, Object value) throws MetricReportingException {
        eventMap.put(key, value);
        return this;
    }
}
