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

import com.google.gson.Gson;
import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.wso2.apk.enforcer.analytics.publisher.exception.MetricReportingException;
import org.wso2.apk.enforcer.analytics.publisher.jmx.JMXUtils;
import org.wso2.apk.enforcer.analytics.publisher.jmx.impl.ExtAuthMetrics;
import org.wso2.apk.enforcer.analytics.publisher.reporter.CounterMetric;
import org.wso2.apk.enforcer.analytics.publisher.reporter.GenericInputValidator;
import org.wso2.apk.enforcer.analytics.publisher.reporter.MetricEventBuilder;
import org.wso2.apk.enforcer.analytics.publisher.reporter.MetricSchema;

import java.io.UnsupportedEncodingException;
import java.net.URLEncoder;
import java.nio.charset.StandardCharsets;
import java.util.Map;

/**
 * Prometheus Counter Metrics class, This class can be used to send analytics event as a prometheus metric event.
 */
public class PrometheusCounterMetric implements CounterMetric {
    private static final Log log = LogFactory.getLog(PrometheusCounterMetric.class);
    private final String name;
    private final Gson gson;
    private MetricSchema schema;

    protected PrometheusCounterMetric(String name, MetricSchema schema) {
        this.name = name;
        this.gson = new Gson();
        this.schema = schema;
    }

    @Override
    public int incrementCount(MetricEventBuilder builder) throws MetricReportingException {

        Map<String, Object> event = builder.build();
        String jsonString = gson.toJson(event);
        log.info("JSON String:"+jsonString);
        // Escape the double quotes
//        String escapedJsonString = jsonString.replace("\"", "\\\"");

        // URL encode the JSON string
        String encodedJsonString = null;
        try {
            encodedJsonString = URLEncoder.encode(jsonString, StandardCharsets.UTF_8.toString());
        } catch (UnsupportedEncodingException e) {
            throw new RuntimeException(e);
        }

        log.info("JSON String: " + encodedJsonString);
        APIInvocationEvent apiInvocationEvent = gson.fromJson(jsonString, APIInvocationEvent.class);
        String apiName = apiInvocationEvent.getApiName();
        int proxyResponseCode = apiInvocationEvent.getProxyResponseCode();
        String destination = apiInvocationEvent.getDestination();
        String apiCreatorTenantDomain = apiInvocationEvent.getApiCreatorTenantDomain();
        String platform = apiInvocationEvent.getPlatform();
        String organizationId = apiInvocationEvent.getOrganizationId();
        String apiMethod = apiInvocationEvent.getApiMethod();
        String apiVersion = apiInvocationEvent.getApiVersion();
        String gatewayType = apiInvocationEvent.getGatewayType();
        String environmentId = apiInvocationEvent.getEnvironmentId();
        String apiCreator = apiInvocationEvent.getApiCreator();
        boolean responseCacheHit = apiInvocationEvent.isResponseCacheHit();
        int backendLatency = apiInvocationEvent.getBackendLatency();
        String correlationId = apiInvocationEvent.getCorrelationId();
        int requestMediationLatency = apiInvocationEvent.getRequestMediationLatency();
        String keyType = apiInvocationEvent.getKeyType();
        String apiId = apiInvocationEvent.getApiId();
        String applicationName = apiInvocationEvent.getApplicationName();
        int targetResponseCode = apiInvocationEvent.getTargetResponseCode();
        String requestTimestamp = apiInvocationEvent.getRequestTimestamp();
        String applicationOwner = apiInvocationEvent.getApplicationOwner();
        String userAgent = apiInvocationEvent.getUserAgent();
        String userName = apiInvocationEvent.getUserName();
        String apiResourceTemplate = apiInvocationEvent.getApiResourceTemplate();
        String regionId = apiInvocationEvent.getRegionId();
        int responseLatency = apiInvocationEvent.getResponseLatency();
        int responseMediationLatency = apiInvocationEvent.getResponseMediationLatency();
        String userIp = apiInvocationEvent.getUserIp();
        String apiContext = apiInvocationEvent.getApiContext();
        String applicationId = apiInvocationEvent.getApplicationId();
        String apiType = apiInvocationEvent.getApiType();
        Map<String, String> properties = apiInvocationEvent.getProperties();

        log.info("PrometheusMetrics: " + name.replaceAll("[\r\n]", "") + ", properties :" +
                jsonString.replaceAll("[\r\n]", ""));

        String propertiesString = encodedJsonString;
        log.info("Properties String:"+propertiesString);
        if (JMXUtils.isJMXMetricsEnabled()) {
            ExtAuthMetrics.getInstance(apiInvocationEvent).recordApiMessages();
            log.info("apiName: " + apiName + ", applicationId: " + applicationId);
        }
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
                return new PrometheusMetricEventBuilder(
                        GenericInputValidator.getInstance().getEventProperties(MetricSchema.PROMETHEUS_RESPONSE));
            case ERROR:
                return new PrometheusMetricEventBuilder(
                        GenericInputValidator.getInstance().getEventProperties(MetricSchema.PROMETHEUS_ERROR));
        }
    }
}
