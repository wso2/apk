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

import org.wso2.apk.enforcer.analytics.publisher.exception.MetricReportingException;
import org.wso2.apk.enforcer.analytics.publisher.reporter.AbstractMetricEventBuilder;
import org.wso2.apk.enforcer.analytics.publisher.reporter.MetricEventBuilder;
import org.wso2.apk.enforcer.analytics.publisher.reporter.MetricSchema;
import org.wso2.apk.enforcer.analytics.publisher.util.Constants;

import java.util.HashMap;
import java.util.Map;

/**
 * Builder class for fault events.
 */
public class DefaultFaultMetricEventBuilder extends AbstractMetricEventBuilder {
    protected final Map<String, Class> requiredAttributes;
    protected final Map<String, Object> eventMap;

    public DefaultFaultMetricEventBuilder() {
        requiredAttributes = DefaultInputValidator.getInstance().getEventProperties(MetricSchema.ERROR);
        eventMap = new HashMap<>();
    }

    protected DefaultFaultMetricEventBuilder(Map<String, Class> requiredAttributes) {
        this.requiredAttributes = requiredAttributes;
        eventMap = new HashMap<>();
    }

    @Override
    public boolean validate() throws MetricReportingException {
        for (Map.Entry<String, Class> entry : requiredAttributes.entrySet()) {
            Object attribute = eventMap.get(entry.getKey());
            if (attribute == null) {
                throw new MetricReportingException(entry.getKey() + " is missing in metric data. This metric event "
                                                           + "will not be processed further.");
            } else if (!attribute.getClass().equals(entry.getValue())) {
                throw new MetricReportingException(entry.getKey() + " is expecting a " + entry.getValue() + " type "
                                                           + "attribute while attribute of type " + attribute.getClass()
                                                           + " is present");
            }
        }
        return true;
    }

    @Override
    public MetricEventBuilder addAttribute(String key, Object value) throws MetricReportingException {
        //all validation is moved to validate method to reduce analytics data processing latency
        eventMap.put(key, value);
        return this;
    }

    @Override
    protected Map<String, Object> buildEvent() {
        eventMap.put(Constants.EVENT_TYPE, Constants.FAULT_EVENT_TYPE);
        // properties object is not required and removing
        eventMap.remove(Constants.PROPERTIES);
        return eventMap;
    }
}
