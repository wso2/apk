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

import org.wso2.apk.enforcer.analytics.publisher.exception.MetricReportingException;
import org.wso2.apk.enforcer.analytics.publisher.reporter.AbstractMetricEventBuilder;
import org.wso2.apk.enforcer.analytics.publisher.reporter.GenericInputValidator;
import org.wso2.apk.enforcer.analytics.publisher.reporter.MetricEventBuilder;
import org.wso2.apk.enforcer.analytics.publisher.reporter.MetricSchema;
import org.wso2.apk.enforcer.analytics.publisher.util.Constants;
import org.wso2.apk.enforcer.analytics.publisher.util.EventMapAttributeFilter;
import org.wso2.apk.enforcer.analytics.publisher.util.UserAgentParser;
import ua_parser.Client;

import java.util.HashMap;
import java.util.Map;

/**
 * Event builder for Prometheus Metric Reporter.
 */
public class PrometheusMetricEventBuilder extends AbstractMetricEventBuilder {

    protected Map<String, Class> requiredAttributes;
    private Map<String, Object> eventMap;
    private Boolean isBuilt = false;

    public PrometheusMetricEventBuilder() {
        requiredAttributes = GenericInputValidator.getInstance().getEventProperties(MetricSchema.RESPONSE);
        eventMap = new HashMap<>();
    }

    public PrometheusMetricEventBuilder(Map<String, Class> requiredAttributes) {
        this.requiredAttributes = requiredAttributes;
        eventMap = new HashMap<>();
    }

    @Override
    protected Map<String, Object> buildEvent() {
        if (!isBuilt) {
            // util function to filter required attributes
            eventMap = EventMapAttributeFilter.getInstance().filter(eventMap, requiredAttributes);

            // userAgent raw string is not required and removing
            String userAgentHeader = (String) eventMap.remove(Constants.USER_AGENT_HEADER);
            if (userAgentHeader != null) {
                setUserAgentProperties(userAgentHeader);
            }
            isBuilt = true;
        }
        return eventMap;
    }

    @Override
    public boolean validate() throws MetricReportingException {
        if (!isBuilt) {
            Map<String, String> propertyMap = (Map<String, String>) eventMap.get(Constants.PROPERTIES);
            if (propertyMap != null) {
                copyDefaultPropertiesToRootLevel(propertyMap);
            }
            for (Map.Entry<String, Class> entry : requiredAttributes.entrySet()) {
                Object attribute = eventMap.get(entry.getKey());
                if (attribute == null) {
                    throw new MetricReportingException(entry.getKey() + " is missing in metric data. This metric event "
                            + "will not be processed further.");
                } else if (!attribute.getClass().equals(entry.getValue())) {
                    throw new MetricReportingException(entry.getKey() + " is expecting a " + entry.getValue() + " type "
                            + "attribute while attribute of type "
                            + attribute.getClass() + " is present.");
                }
            }
        }
        return true;
    }

    @Override
    public MetricEventBuilder addAttribute(String key, Object value) throws MetricReportingException {
        eventMap.put(key, value);
        return this;
    }

    private void setUserAgentProperties(String userAgentHeader) {
        String browser = null;
        String platform = null;
        Client client = UserAgentParser.getInstance().parseUserAgent(userAgentHeader);
        if (client != null) {
            browser = client.userAgent.family;
            platform = client.os.family;
        }

        if (browser == null || browser.isEmpty()) {
            browser = Constants.UNKNOWN_VALUE;
        }
        if (platform == null || platform.isEmpty()) {
            platform = Constants.UNKNOWN_VALUE;
        }
        eventMap.put(Constants.USER_AGENT, browser);
        eventMap.put(Constants.PLATFORM, platform);
    }

    private void copyDefaultPropertiesToRootLevel(Map<String, String> properties) {

        if (properties.get(Constants.API_CONTEXT) != null) {
            eventMap.put(Constants.API_CONTEXT, properties.get(Constants.API_CONTEXT));
        }
        if (properties.get(Constants.USER_NAME) != null) {
            eventMap.put(Constants.USER_NAME, properties.get(Constants.USER_NAME));
        }
    }
}
