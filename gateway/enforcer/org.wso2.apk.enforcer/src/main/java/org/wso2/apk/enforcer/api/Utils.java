/*
 * Copyright (c) 2021, WSO2 LLC. (http://www.wso2.org).
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
package org.wso2.apk.enforcer.api;

import org.wso2.apk.enforcer.commons.model.APIKeyAuthenticationConfig;
import org.wso2.apk.enforcer.commons.model.AuthenticationConfig;
import org.wso2.apk.enforcer.commons.model.JWTAuthenticationConfig;
import org.wso2.apk.enforcer.commons.model.Oauth2AuthenticationConfig;
import org.wso2.apk.enforcer.config.EnforcerConfig;
import org.wso2.apk.enforcer.discovery.api.APIKey;
import org.wso2.apk.enforcer.discovery.api.EndpointClusterConfig;
import org.wso2.apk.enforcer.discovery.api.Operation;
import org.wso2.apk.enforcer.discovery.api.OperationPolicies;
import org.wso2.apk.enforcer.commons.model.EndpointCluster;
import org.wso2.apk.enforcer.commons.model.EndpointSecurity;
import org.wso2.apk.enforcer.commons.model.Policy;
import org.wso2.apk.enforcer.commons.model.PolicyConfig;
import org.wso2.apk.enforcer.commons.model.RequestContext;
import org.wso2.apk.enforcer.commons.model.ResourceConfig;
import org.wso2.apk.enforcer.commons.model.RetryConfig;
import org.wso2.apk.enforcer.config.ConfigHolder;
import org.wso2.apk.enforcer.constants.APIConstants;
import org.wso2.apk.enforcer.constants.AdapterConstants;

import java.util.ArrayList;
import java.util.List;

/**
 * Utility Methods used across different APIs.
 */
public class Utils {

    /**
     * Construct Enforcer Specific DTOs for endpoints from XDS specific DTOs.
     *
     * @param rpcEndpointCluster XDS specific EndpointCluster Instance
     * @return Enforcer Specific XDS cluster instance
     */
    public static EndpointCluster processEndpoints(org.wso2.apk.enforcer.discovery.api.EndpointCluster
                                                           rpcEndpointCluster) {
        if (rpcEndpointCluster == null || rpcEndpointCluster.getUrlsCount() == 0) {
            return null;
        }
        List<String> urls = new ArrayList<>(1);
        rpcEndpointCluster.getUrlsList().forEach(endpoint -> {
            String url = endpoint.getURLType().toLowerCase() + "://" +
                    endpoint.getHost() + ":" + endpoint.getPort() + endpoint.getBasepath();
            urls.add(url);
        });
        EndpointCluster endpointCluster = new EndpointCluster();
        endpointCluster.setUrls(urls);

        // Getting the first endpoint's basepath would be enough,
        // as all endpoints in the cluster should have the same basepath as of now
        endpointCluster.setBasePath(rpcEndpointCluster.getUrlsList().get(0).getBasepath());
        if (rpcEndpointCluster.hasConfig()) {
            EndpointClusterConfig endpointClusterConfig = rpcEndpointCluster.getConfig();
            if (endpointClusterConfig.hasRetryConfig()) {
                org.wso2.apk.enforcer.discovery.api.RetryConfig rpcRetryConfig
                        = endpointClusterConfig.getRetryConfig();
                RetryConfig retryConfig = new RetryConfig(rpcRetryConfig.getCount(),
                        rpcRetryConfig.getStatusCodesList().toArray(new Integer[0]));
                endpointCluster.setRetryConfig(retryConfig);
            }
            if (endpointClusterConfig.hasTimeoutConfig()) {
                org.wso2.apk.enforcer.discovery.api.TimeoutConfig timeoutConfig
                        = endpointClusterConfig.getTimeoutConfig();
                endpointCluster.setRouteTimeoutInMillis(timeoutConfig.getRouteTimeoutInMillis());
            }
        }
        return endpointCluster;
    }

    public static ResourceConfig buildResource(Operation operation, String resPath, EndpointSecurity[] endpointSecurity) {
        ResourceConfig resource = new ResourceConfig();
        resource.setMatchID(operation.getMatchID());
        resource.setPath(resPath);
        resource.setMethod(ResourceConfig.HttpMethods.valueOf(operation.getMethod().toUpperCase()));
        resource.setTier(operation.getTier());
        resource.setEndpointSecurity(endpointSecurity);
        AuthenticationConfig authenticationConfig = new AuthenticationConfig();

        if (ConfigHolder.getInstance().getConfig()
                .getMandateInternalKeyValidation()) {
            JWTAuthenticationConfig jwtAuthenticationConfig = getDefaultJwtAuthenticationConfig();
            authenticationConfig.setJwtAuthenticationConfig(jwtAuthenticationConfig);
        }

        if (operation.hasApiAuthentication()) {
            authenticationConfig.setDisabled(operation.getApiAuthentication().getDisabled());
            if (operation.getApiAuthentication().hasOauth2()) {
                Oauth2AuthenticationConfig oAuth2AuthenticationConfig = new Oauth2AuthenticationConfig();
                oAuth2AuthenticationConfig.setHeader(operation.getApiAuthentication().getOauth2().getHeader());
                oAuth2AuthenticationConfig.setSendTokenToUpstream(operation.getApiAuthentication().getOauth2()
                        .getSendTokenToUpstream());
                authenticationConfig.setOauth2AuthenticationConfig(oAuth2AuthenticationConfig);
            }
            if (operation.getApiAuthentication().hasJwt()) {
                JWTAuthenticationConfig jwtAuthenticationConfig = getJwtAuthenticationConfig(operation);
                authenticationConfig.setJwtAuthenticationConfig(jwtAuthenticationConfig);
            }
            List<APIKeyAuthenticationConfig> apiKeyAuthenticationConfigs = getApiKeyAuthenticationConfigs(operation);
            authenticationConfig.setApiKeyAuthenticationConfigs(apiKeyAuthenticationConfigs);
        }
        resource.setAuthenticationConfig(authenticationConfig);
        resource.setScopes(operation.getScopesList().toArray(new String[0]));
        return resource;
    }

    private static List<APIKeyAuthenticationConfig> getApiKeyAuthenticationConfigs(Operation operation) {
        List<APIKeyAuthenticationConfig> apiKeyAuthenticationConfigs = new ArrayList<>();
        for (APIKey apiKey : operation.getApiAuthentication().getApikeyList()) {
            APIKeyAuthenticationConfig apiKeyAuthenticationConfig = new APIKeyAuthenticationConfig();
            apiKeyAuthenticationConfig.setIn(apiKey.getIn());
            apiKeyAuthenticationConfig.setName(apiKey.getName());
            apiKeyAuthenticationConfig.setSendTokenToUpstream(apiKey.getSendTokenToUpstream());
            apiKeyAuthenticationConfigs.add(apiKeyAuthenticationConfig);
        }
        return apiKeyAuthenticationConfigs;
    }

    private static JWTAuthenticationConfig getJwtAuthenticationConfig(Operation operation) {
        JWTAuthenticationConfig jwtAuthenticationConfig = new JWTAuthenticationConfig();
        jwtAuthenticationConfig.setHeader(operation.getApiAuthentication().getJwt().getHeader());
        jwtAuthenticationConfig.setSendTokenToUpstream(operation.getApiAuthentication().getJwt()
                .getSendTokenToUpstream());
        List<String> audience = new ArrayList<>();
        for (int i = 0; i < operation.getApiAuthentication().getJwt().getAudienceCount(); i++) {
            audience.add(operation.getApiAuthentication().getJwt().getAudience(i));
        }
        jwtAuthenticationConfig.setAudience(audience);
        return jwtAuthenticationConfig;
    }

    private static JWTAuthenticationConfig getDefaultJwtAuthenticationConfig() {
        JWTAuthenticationConfig jwtAuthenticationConfig = new JWTAuthenticationConfig();
        jwtAuthenticationConfig.setHeader(APIConstants.TEST_CONSOLE_KEY_HEADER);
        jwtAuthenticationConfig.setSendTokenToUpstream(false);
        return jwtAuthenticationConfig;
    }

    public static PolicyConfig genPolicyConfig(OperationPolicies operationPolicies) {
        PolicyConfig policyConfig = new PolicyConfig();
        if (operationPolicies.getRequestCount() > 0) {
            policyConfig.setRequest(genPolicyList(operationPolicies.getRequestList()));
        }
        if (operationPolicies.getResponseCount() > 0) {
            policyConfig.setResponse(genPolicyList(operationPolicies.getResponseList()));
        }
        if (operationPolicies.getFaultCount() > 0) {
            policyConfig.setFault(genPolicyList(operationPolicies.getFaultList()));
        }
        return policyConfig;
    }

    private static ArrayList<Policy> genPolicyList
            (List<org.wso2.apk.enforcer.discovery.api.Policy> operationPoliciesList) {
        ArrayList<Policy> policyList = new ArrayList<>();
        for (org.wso2.apk.enforcer.discovery.api.Policy policy : operationPoliciesList) {
            policyList.add(new Policy(policy.getAction(), policy.getParametersMap()));
        }
        return policyList;
    }

    /**
     * Set common authentication headers as headers to be removed and protected headers,
     * so that they will be removed by the router before being sent to the backend.
     *
     * @param requestContext requestContext
     */
    public static void handleCommonHeaders(RequestContext requestContext) {
        // not allow clients to set cluster header manually
        requestContext.getRemoveHeaders().add(AdapterConstants.CLUSTER_HEADER);
    }
}
