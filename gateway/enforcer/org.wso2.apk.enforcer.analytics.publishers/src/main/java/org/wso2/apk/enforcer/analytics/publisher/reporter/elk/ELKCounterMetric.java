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

package org.wso2.apk.enforcer.analytics.publisher.reporter.elk;

import com.google.gson.Gson;
import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.wso2.apk.enforcer.analytics.publisher.exception.MetricReportingException;
import org.wso2.apk.enforcer.analytics.publisher.reporter.CounterMetric;
import org.wso2.apk.enforcer.analytics.publisher.reporter.GenericInputValidator;
import org.wso2.apk.enforcer.analytics.publisher.reporter.MetricEventBuilder;
import org.wso2.apk.enforcer.analytics.publisher.reporter.MetricSchema;

import java.util.Map;

/**
 * Log Counter Metrics class, This class can be used to log analytics event to a separate log file.
 */
public class ELKCounterMetric implements CounterMetric {
    private static final Log log = LogFactory.getLog(ELKCounterMetric.class);
    private final String name;
    private final Gson gson;
    private MetricSchema schema;

    protected ELKCounterMetric(String name, MetricSchema schema) {
        this.name = name;
        this.gson = new Gson();
        this.schema = schema;
    }

    @Override
    public int incrementCount(MetricEventBuilder builder) throws MetricReportingException {
        Map<String, Object> event = builder.build();
        String jsonString = gson.toJson(event);

        log.info("apimMetrics: " + name.replaceAll("[\r\n]", "") + ", properties :" +
                jsonString.replaceAll("[\r\n]", ""));
        return 0;
    }

    @Override
    public String getName() {
        return name;
    }

    @Override
    public MetricSchema getSchema() {
        return schema;
    }

    @Override
    public MetricEventBuilder getEventBuilder() {

        switch (schema) {
            case RESPONSE:
            default:
                return new ELKMetricEventBuilder(
                        GenericInputValidator.getInstance().getEventProperties(MetricSchema.ELK_RESPONSE));
            case ERROR:
                return new ELKMetricEventBuilder(
                        GenericInputValidator.getInstance().getEventProperties(MetricSchema.ELK_ERROR));
        }
    }
}
