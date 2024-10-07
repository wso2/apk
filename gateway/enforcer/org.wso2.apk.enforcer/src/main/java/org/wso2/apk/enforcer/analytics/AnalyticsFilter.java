/*
 * Copyright (c) 2021, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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

package org.wso2.apk.enforcer.analytics;

import io.envoyproxy.envoy.service.accesslog.v3.StreamAccessLogsMessage;
import io.opentelemetry.context.Scope;
import org.apache.commons.lang3.StringUtils;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.apache.logging.log4j.ThreadContext;
import org.wso2.apk.enforcer.commons.analytics.collectors.AnalyticsCustomDataProvider;
import org.wso2.apk.enforcer.commons.analytics.collectors.impl.GenericRequestDataCollector;
import org.wso2.apk.enforcer.commons.analytics.exceptions.AnalyticsException;
import org.wso2.apk.enforcer.commons.logging.ErrorDetails;
import org.wso2.apk.enforcer.commons.logging.LoggingConstants;
import org.wso2.apk.enforcer.commons.model.AuthenticationContext;
import org.wso2.apk.enforcer.commons.model.RequestContext;
import org.wso2.apk.enforcer.commons.model.ResourceConfig;
import org.wso2.apk.enforcer.config.ConfigHolder;
import org.wso2.apk.enforcer.config.dto.AnalyticsDTO;
import org.wso2.apk.enforcer.config.dto.AnalyticsPublisherConfigDTO;
import org.wso2.apk.enforcer.constants.APIConstants;
import org.wso2.apk.enforcer.constants.AnalyticsConstants;
import org.wso2.apk.enforcer.constants.MetadataConstants;
import org.wso2.apk.enforcer.tracing.TracingConstants;
import org.wso2.apk.enforcer.tracing.TracingSpan;
import org.wso2.apk.enforcer.tracing.TracingTracer;
import org.wso2.apk.enforcer.tracing.Utils;
import org.wso2.apk.enforcer.util.FilterUtils;

import java.util.ArrayList;
import java.util.List;
import java.util.Map;

import static org.wso2.apk.enforcer.analytics.AnalyticsConstants.CHOREO_FAULT_SCHEMA;
import static org.wso2.apk.enforcer.analytics.AnalyticsConstants.CHOREO_RESPONSE_SCHEMA;
import static org.wso2.apk.enforcer.analytics.AnalyticsConstants.IS_CHOREO_DEPLOYMENT_CONFIG_KEY;

/**
 * This is the filter is for Analytics.
 * If the request is failed at enforcer (due to throttling, authentication failures) the analytics event is
 * published by the filter itself.
 * If the request is allowed to proceed, the dynamic metadata will be populated so that the analytics event can be
 * populated from grpc access logs within AccessLoggingService.
 */
public class AnalyticsFilter {

    private static final Logger logger = LogManager.getLogger(AnalyticsFilter.class);
    private static AnalyticsFilter analyticsFilter;
    private static AnalyticsEventPublisher publisher;
    private static AnalyticsCustomDataProvider analyticsDataProvider;

    private AnalyticsFilter() {

        AnalyticsDTO analyticsConfig = ConfigHolder.getInstance().getConfig().getAnalyticsConfig();
        Map<String, Object> properties = analyticsConfig.getProperties();
        publisher = new DefaultAnalyticsEventPublisher();
        boolean choreoDeployment = false;
        if (properties != null){
            choreoDeployment = (boolean) properties.getOrDefault(IS_CHOREO_DEPLOYMENT_CONFIG_KEY,false);
        }
        if (choreoDeployment){
            publisher = new DefaultAnalyticsEventPublisher(CHOREO_RESPONSE_SCHEMA, CHOREO_FAULT_SCHEMA);
        }
        List<AnalyticsPublisherConfigDTO> analyticsPublisherConfigDTOList =
                ConfigHolder.getInstance().getConfig().getAnalyticsConfig().getAnalyticsPublisherConfigDTOList();
        publisher.init(analyticsPublisherConfigDTOList);
    }

    public static AnalyticsFilter getInstance() {

        if (analyticsFilter == null) {
            synchronized (new Object()) {
                if (analyticsFilter == null) {
                    analyticsFilter = new AnalyticsFilter();
                }
            }
        }
        return analyticsFilter;
    }

    public void handleGRPCLogMsg(StreamAccessLogsMessage message) {

        if (publisher != null) {
            publisher.handleGRPCLogMsg(message);
        } else {
            logger.error("Cannot publish the analytics event as analytics publisher is null.",
                    ErrorDetails.errorLog(LoggingConstants.Severity.CRITICAL, 5102));
        }
    }

//    public void handleWebsocketFrameRequest(WebSocketFrameRequest frameRequest) {
//        if (publisher != null) {
//            publisher.handleWebsocketFrameRequest(frameRequest);
//        } else {
//            logger.error("Cannot publish the analytics event as analytics publisher is null.",
//                    ErrorDetails.errorLog(LoggingConstants.Severity.CRITICAL, 5102));
//        }
//    }

    public static AnalyticsCustomDataProvider getAnalyticsCustomDataProvider() {

        return analyticsDataProvider;
    }

    public void handleSuccessRequest(RequestContext requestContext) {

        TracingSpan analyticsSpan = null;
        Scope analyticsSpanScope = null;
        try {
            if (Utils.tracingEnabled()) {
                TracingTracer tracer = Utils.getGlobalTracer();
                analyticsSpan = Utils.startSpan(TracingConstants.ANALYTICS_SPAN, tracer);
                analyticsSpanScope = analyticsSpan.getSpan().makeCurrent();
                Utils.setTag(analyticsSpan, APIConstants.LOG_TRACE_ID,
                        ThreadContext.get(APIConstants.LOG_TRACE_ID));
            }
            String apiName = requestContext.getMatchedAPI().getName();
            String apiVersion = requestContext.getMatchedAPI().getVersion();
            String apiType = requestContext.getMatchedAPI().getApiType();
            boolean isMockAPI = requestContext.getMatchedAPI().isMockedApi();
            AuthenticationContext authContext = AnalyticsUtils.getAuthenticationContext(requestContext);
            requestContext.addMetadataToMap(MetadataConstants.API_ID_KEY, AnalyticsUtils.getAPIId(requestContext));
            requestContext.addMetadataToMap(MetadataConstants.API_CREATOR_KEY,
                    AnalyticsUtils.setDefaultIfNull(authContext.getApiPublisher()));
            requestContext.addMetadataToMap(MetadataConstants.API_NAME_KEY, apiName);
            requestContext.addMetadataToMap(MetadataConstants.API_VERSION_KEY, apiVersion);
            requestContext.addMetadataToMap(MetadataConstants.API_TYPE_KEY, apiType);
            requestContext.addMetadataToMap(MetadataConstants.IS_MOCK_API, String.valueOf(isMockAPI));

            String tenantDomain = requestContext.getMatchedAPI().getOrganizationId();
            requestContext.addMetadataToMap(MetadataConstants.API_CREATOR_TENANT_DOMAIN_KEY,
                    tenantDomain == null ? APIConstants.SUPER_TENANT_DOMAIN_NAME : tenantDomain);

            // Default Value would be PRODUCTION
            requestContext.addMetadataToMap(MetadataConstants.APP_KEY_TYPE_KEY,
                    requestContext.getMatchedAPI().getEnvType());
            requestContext.addMetadataToMap(MetadataConstants.APP_UUID_KEY,
                    AnalyticsUtils.setDefaultIfNull(authContext.getApplicationUUID()));
            requestContext.addMetadataToMap(MetadataConstants.APP_NAME_KEY,
                    AnalyticsUtils.setDefaultIfNull(authContext.getApplicationName()));
            requestContext.addMetadataToMap(MetadataConstants.APP_OWNER_KEY,
                    AnalyticsUtils.setDefaultIfNull(authContext.getSubscriber()));

            requestContext.addMetadataToMap(MetadataConstants.CORRELATION_ID_KEY, requestContext.getRequestID());
            requestContext.addMetadataToMap(MetadataConstants.REGION_KEY,
                    ConfigHolder.getInstance().getEnvVarConfig().getEnforcerRegionId());

            // As in the matched API, only the resources under the matched resource template are selected.
            ArrayList<String> resourceTemplate = new ArrayList<>();
            for (ResourceConfig resourceConfig : requestContext.getMatchedResourcePaths()) {
                resourceTemplate.add(resourceConfig.getPath());
            }
            requestContext.addMetadataToMap(MetadataConstants.API_RESOURCE_TEMPLATE_KEY,
                    String.join(",", resourceTemplate));

            requestContext.addMetadataToMap(MetadataConstants.DESTINATION, resolveEndpoint(requestContext));

            requestContext.addMetadataToMap(MetadataConstants.API_ORGANIZATION_ID,
                    requestContext.getMatchedAPI().getOrganizationId());
            requestContext.addMetadataToMap(MetadataConstants.CLIENT_IP_KEY, requestContext.getClientIp());
            requestContext.addMetadataToMap(MetadataConstants.USER_AGENT_KEY,
                    AnalyticsUtils.setDefaultIfNull(requestContext.getHeaders().get("user-agent")));

            // Adding UserName and the APIContext
            String endUserName = authContext.getUsername();
            requestContext.addMetadataToMap(MetadataConstants.API_USER_NAME_KEY,
                    endUserName == null ? APIConstants.END_USER_UNKNOWN : endUserName);
            requestContext.addMetadataToMap(MetadataConstants.API_CONTEXT_KEY,
                    requestContext.getMatchedAPI().getBasePath());
            requestContext.addMetadataToMap(MetadataConstants.API_ENVIRONMENT,
                    requestContext.getMatchedAPI().getEnvironment() == null
                            ? APIConstants.DEFAULT_ENVIRONMENT_NAME
                            : requestContext.getMatchedAPI().getEnvironment());
            // Adding Gateway URL
            String gatewayUrl = requestContext.getHeaders().get(AnalyticsConstants.GATEWAY_URL);
            if (!StringUtils.isNotEmpty(gatewayUrl)) {
                String protocol = requestContext.getHeaders().getOrDefault(AnalyticsConstants.X_FORWARD_PROTO_HEADER,
                        APIConstants.HTTPS_PROTOCOL);
                String port = requestContext.getHeaders().getOrDefault(AnalyticsConstants.X_FORWARD_PORT_HEADER, "");
                String host = requestContext.getMatchedAPI().getVhost();
                String path = requestContext.getRequestPath();
                gatewayUrl = protocol.concat("://").concat(host).concat(":").concat(port).concat(path);
            }
            requestContext.addMetadataToMap(MetadataConstants.GATEWAY_URL, gatewayUrl);
        } finally {
            if (Utils.tracingEnabled()) {
                analyticsSpanScope.close();
                Utils.finishSpan(analyticsSpan);
            }
        }
    }

    public String resolveEndpoint(RequestContext requestContext) {

        // For MockAPIs there is no backend, Hence "MockImplementation" is added as the destination.
        if (requestContext.getMatchedAPI() != null && requestContext.getMatchedAPI().isMockedApi()) {
            return "MockImplementation";
        }
        // This does not cause problems at the moment Since the current microgateway supports only one URL
        try {
            return requestContext.getMatchedResourcePaths().get(0).getEndpoints().getUrls().get(0);
        } catch (Exception e) {
            return "";
        }
    }

    public void handleFailureRequest(RequestContext requestContext) {

        TracingSpan analyticsSpan = null;
        Scope analyticsSpanScope = null;

        try {
            if (Utils.tracingEnabled()) {
                TracingTracer tracer = Utils.getGlobalTracer();
                analyticsSpan = Utils.startSpan(TracingConstants.ANALYTICS_FAILURE_SPAN, tracer);
                analyticsSpanScope = analyticsSpan.getSpan().makeCurrent();
                Utils.setTag(analyticsSpan, APIConstants.LOG_TRACE_ID,
                        ThreadContext.get(APIConstants.LOG_TRACE_ID));
            }
            if (publisher == null) {
                logger.error("Cannot publish the failure event as analytics publisher is null.",
                        ErrorDetails.errorLog(LoggingConstants.Severity.CRITICAL, 5103));
                return;
            }
            ChoreoFaultAnalyticsProvider provider = new ChoreoFaultAnalyticsProvider(requestContext);
            // To avoid incrementing counter for options call
            if (provider.getProxyResponseCode() == 200 || provider.getProxyResponseCode() == 204) {
                return;
            }
            GenericRequestDataCollector dataCollector = new GenericRequestDataCollector(provider);
            try {
                dataCollector.collectData();
                logger.debug("Analytics event for failure event is published.");
            } catch (AnalyticsException e) {
                logger.error("Error while publishing the analytics event. ", e);
            }
        } finally {
            if (Utils.tracingEnabled()) {
                analyticsSpanScope.close();
                Utils.finishSpan(analyticsSpan);
            }
        }
    }

}
