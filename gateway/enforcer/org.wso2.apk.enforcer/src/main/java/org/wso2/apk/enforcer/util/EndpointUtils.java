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

package org.wso2.apk.enforcer.util;

import org.apache.commons.lang3.StringUtils;
import org.wso2.apk.enforcer.commons.model.EndpointCluster;
import org.wso2.apk.enforcer.commons.model.EndpointSecurity;
import org.wso2.apk.enforcer.commons.model.RequestContext;
import org.wso2.apk.enforcer.commons.model.ResourceConfig;
import org.wso2.apk.enforcer.commons.model.RetryConfig;
import org.wso2.apk.enforcer.config.ConfigHolder;
import org.wso2.apk.enforcer.config.EnforcerConfig;
import org.wso2.apk.enforcer.constants.APIConstants;

import java.util.Base64;
import org.wso2.apk.enforcer.constants.AdapterConstants;

/**
 * Util methods related to backend endpoint security.
 */
public class EndpointUtils {

    /**
     * Adds the backend endpoint security header to the given requestContext.
     *
     * @param requestContext requestContext instance to add the backend endpoint security header
     */
    public static void addEndpointSecurity(RequestContext requestContext) {
        if (requestContext.getMatchedResourcePaths() != null) {
            // getting only first element as there could be only one resourcepaths for APIs except for graphQL APIs. 
            // For GQL APIs too would only have one endpoint for all resources.
            ResourceConfig resourceConfig = requestContext.getMatchedResourcePaths().get(0);
            if (resourceConfig.getEndpointSecurity() != null) {
                for (EndpointSecurity securityInfo : resourceConfig.getEndpointSecurity()) {
                    if (securityInfo != null && securityInfo.isEnabled() &&
                            APIConstants.AUTHORIZATION_HEADER_BASIC.
                                    equalsIgnoreCase(securityInfo.getSecurityType())) {
                        requestContext.getRemoveHeaders().remove(APIConstants.AUTHORIZATION_HEADER_DEFAULT
                                .toLowerCase());
                        requestContext.addOrModifyHeaders(APIConstants.AUTHORIZATION_HEADER_DEFAULT,
                                APIConstants.AUTHORIZATION_HEADER_BASIC + ' ' +
                                        Base64.getEncoder().encodeToString((securityInfo.getUsername() +
                                                ':' + String.valueOf(securityInfo.getPassword())).getBytes()));
                    }
                }
            }
        }
    }

    /**
     * Update the cluster header based on the keyType and authenticate the token against its respective endpoint
     * environment.
     *
     * @param requestContext request Context
     */
    public static void updateClusterHeaderAndCheckEnv(RequestContext requestContext) {
        EnforcerConfig enforcerConfig = ConfigHolder.getInstance().getConfig();
        if (!enforcerConfig.getEnableGatewayClassController()) {
            requestContext.addOrModifyHeaders(AdapterConstants.CLUSTER_HEADER, requestContext.getClusterHeader());
        }
        requestContext.getRemoveHeaders().remove(AdapterConstants.CLUSTER_HEADER);
        addRouterHttpHeaders(requestContext);
        addEndpointSecurity(requestContext);
    }

    private static void addRouterHttpHeaders(RequestContext requestContext) {
        // requestContext.getMatchedResourcePaths() will only have one element for non GraphQL APIs.
        // Also, GraphQL APIs doesn't have resource level endpoint configs
        ResourceConfig resourceConfig = requestContext.getMatchedResourcePaths().get(0);
        // In websockets case, the endpoints object becomes null. Hence it would result
        // in a NPE, if it is not checked.
        if (resourceConfig.getEndpoints() != null) {
            EndpointCluster endpointCluster = resourceConfig.getEndpoints();
            addRetryAndTimeoutConfigHeaders(requestContext, endpointCluster);
            handleEmptyPathHeader(requestContext, endpointCluster.getBasePath());
        }
    }

    private static void addRetryAndTimeoutConfigHeaders(RequestContext requestContext, EndpointCluster endpointCluster) {
        RetryConfig retryConfig = endpointCluster.getRetryConfig();
        if (retryConfig != null) {
            addRetryConfigHeaders(requestContext, retryConfig);
        }
        Integer timeout = endpointCluster.getRouteTimeoutInMillis();
        if (timeout != null) {
            addTimeoutHeaders(requestContext, timeout);
        }
    }

    /**
     * This will fix sending upstream an empty path header issue.
     *
     * @param requestContext request context
     * @param basePath       endpoint basepath
     */
    private static void handleEmptyPathHeader(RequestContext requestContext, String basePath) {
        if (StringUtils.isNotBlank(basePath)) {
            return;
        }
        // remaining path after removing the context and the version from the invoked path.
        String remainingPath = StringUtils.removeStartIgnoreCase(requestContext.getHeaders()
                .get(APIConstants.PATH_HEADER).split("\\?")[0], requestContext.getMatchedAPI().getBasePath());
        // if the :path will be empty after applying the route's substitution, then we have to add a "/" forcefully
        // to avoid :path being empty.
        if (StringUtils.isBlank(remainingPath)) {
            String[] splittedPath = requestContext.getHeaders().get(APIConstants.PATH_HEADER).split("\\?");
            String newPath = splittedPath.length > 1 ? splittedPath[0] + "/?" + splittedPath[1] : splittedPath[0] + "/";
            requestContext.addOrModifyHeaders(APIConstants.PATH_HEADER, newPath);
        }
    }

    private static void addRetryConfigHeaders(RequestContext requestContext, RetryConfig retryConfig) {
        requestContext.addOrModifyHeaders(AdapterConstants.HttpRouterHeaders.RETRY_ON,
                AdapterConstants.HttpRouterHeaderValues.RETRIABLE_STATUS_CODES);
        requestContext.addOrModifyHeaders(AdapterConstants.HttpRouterHeaders.MAX_RETRIES,
                Integer.toString(retryConfig.getCount()));
        requestContext.addOrModifyHeaders(AdapterConstants.HttpRouterHeaders.RETRIABLE_STATUS_CODES,
                StringUtils.join(retryConfig.getStatusCodes(), ","));
    }

    private static void addTimeoutHeaders(RequestContext requestContext, Integer routeTimeoutInMillis) {
        requestContext.addOrModifyHeaders(AdapterConstants.HttpRouterHeaders.UPSTREAM_REQ_TIMEOUT_MS,
                Integer.toString(routeTimeoutInMillis));
    }
}
