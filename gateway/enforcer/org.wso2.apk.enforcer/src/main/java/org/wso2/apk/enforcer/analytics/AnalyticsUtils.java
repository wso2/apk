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

import io.envoyproxy.envoy.data.accesslog.v3.HTTPAccessLogEntry;
import org.apache.commons.lang3.StringUtils;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.wso2.apk.enforcer.commons.model.AuthenticationContext;
import org.wso2.apk.enforcer.commons.model.RequestContext;
import org.wso2.apk.enforcer.config.ConfigHolder;
import org.wso2.apk.enforcer.constants.AnalyticsConstants;
import org.wso2.apk.enforcer.constants.MetadataConstants;
import org.wso2.apk.enforcer.models.API;
import org.wso2.apk.enforcer.subscription.SubscriptionDataHolder;

import java.lang.reflect.Constructor;
import java.lang.reflect.InvocationTargetException;
import java.util.Map;

import static org.wso2.apk.enforcer.analytics.AnalyticsConstants.GATEWAY_TYPE_CONFIG_KEY;
import static org.wso2.apk.enforcer.analytics.AnalyticsConstants.DEFAULT_GATEWAY_TYPE;

/**
 * Common Utility functions
 */
public class AnalyticsUtils {
    private static final Logger logger = LogManager.getLogger(AnalyticsUtils.class);

    public static String getAPIId(RequestContext requestContext) {
        return requestContext.getMatchedAPI().getUuid();
    }
    public static String getOrganization(RequestContext requestContext) {
        return requestContext.getMatchedAPI().getOrganizationId();
    }
    public static String setDefaultIfNull(String value) {
        return value == null ? AnalyticsConstants.DEFAULT_FOR_UNASSIGNED : value;
    }

    public static String getGatewayType() {
        Map<String, Object> properties = ConfigHolder.getInstance().getConfig().getAnalyticsConfig().getProperties();
        String gatewayType = DEFAULT_GATEWAY_TYPE;
        if (properties != null) {
            gatewayType = (String) properties.getOrDefault(GATEWAY_TYPE_CONFIG_KEY, DEFAULT_GATEWAY_TYPE);
        }
        return gatewayType;
    }

    /**
     * Extracts Authentication Context from the request Context. If Authentication Context is not available,
     * new Authentication Context object will be created with authenticated property is set to false.
     *
     * @param requestContext {@code RequestContext} object
     * @return {@code AuthenticationContext} object
     */
    public static AuthenticationContext getAuthenticationContext(RequestContext requestContext) {
        AuthenticationContext authContext = requestContext.getAuthenticationContext();
        // When authentication failure happens authContext remains null
        if (authContext == null) {
            authContext = new AuthenticationContext();
            authContext.setAuthenticated(false);
        }
        return authContext;
    }

    /**
     * Decides if the logEntry corresponds to a mock API. The "x-wso2-is-mock-api" is only set when
     * handling mock-api-request.
     *
     * @param logEntry Access Log Entry
     * @return true if the logEntry has the metadata called "x-wso2-is-mock-api" and its value is true
     */
    public static boolean isMockAPISuccessRequest(HTTPAccessLogEntry logEntry) {

        return (!StringUtils.isEmpty(logEntry.getResponse().getResponseCodeDetails())) &&
                logEntry.getResponse().getResponseCodeDetails()
                        .equals(AnalyticsConstants.EXT_AUTH_DENIED_RESPONSE_DETAIL) &&
                logEntry.hasCommonProperties() &&
                logEntry.getCommonProperties().hasMetadata() &&
                logEntry.getCommonProperties().getMetadata().getFilterMetadataMap()
                        .get(MetadataConstants.EXT_AUTH_METADATA_CONTEXT_KEY) != null &&
                logEntry.getCommonProperties().getMetadata()
                        .getFilterMetadataMap().get(MetadataConstants.EXT_AUTH_METADATA_CONTEXT_KEY).getFieldsMap()
                        .containsKey(MetadataConstants.IS_MOCK_API) &&
                Boolean.parseBoolean(logEntry.getCommonProperties().getMetadata()
                        .getFilterMetadataMap().get(MetadataConstants.EXT_AUTH_METADATA_CONTEXT_KEY).getFieldsMap()
                        .get(MetadataConstants.IS_MOCK_API).getStringValue());
    }

}
