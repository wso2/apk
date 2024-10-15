/*
 * Copyright (c) 2020, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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

import org.apache.commons.lang3.StringUtils;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.wso2.apk.enforcer.analytics.AnalyticsFilter;
import org.wso2.apk.enforcer.commons.Filter;
import org.wso2.apk.enforcer.commons.dto.ClaimValueDTO;
import org.wso2.apk.enforcer.commons.dto.JWTConfigurationDto;
import org.wso2.apk.enforcer.commons.model.APIConfig;
import org.wso2.apk.enforcer.commons.model.RequestContext;
import org.wso2.apk.enforcer.commons.model.ResourceConfig;
import org.wso2.apk.enforcer.config.ConfigHolder;
import org.wso2.apk.enforcer.config.EnforcerConfig;
import org.wso2.apk.enforcer.constants.APIConstants;
import org.wso2.apk.enforcer.constants.HttpConstants;
import org.wso2.apk.enforcer.cors.CorsFilter;
import org.wso2.apk.enforcer.discovery.api.*;
import org.wso2.apk.enforcer.interceptor.MediationPolicyFilter;
import org.wso2.apk.enforcer.security.AuthFilter;
import org.wso2.apk.enforcer.security.mtls.MtlsUtils;
import org.wso2.apk.enforcer.util.EndpointUtils;
import org.wso2.apk.enforcer.util.FilterUtils;
import org.wso2.apk.enforcer.util.MockImplUtils;

import java.security.KeyStore;
import java.security.KeyStoreException;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * Specific implementation for a Rest API type APIs.
 */
public class RestAPI implements API {

    private static final Logger logger = LogManager.getLogger(RestAPI.class);
    private final List<Filter> filters = new ArrayList<>();
    private APIConfig apiConfig;

    @Override
    public List<Filter> getFilters() {

        return filters;
    }

    @Override
    public String init(Api api) {

        String vhost = api.getVhost();
        String basePath = api.getBasePath();
        String name = api.getTitle();
        String version = api.getVersion();
        String apiType = api.getApiType();
        List<ResourceConfig> resources = new ArrayList<>();
        String mutualSSL = api.getMutualSSL();
        Map<String, Boolean> applicationSecurity = api.getApplicationSecurityMap();

        for (Resource res : api.getResourcesList()) {
            for (Operation operation : res.getMethodsList()) {
                ResourceConfig resConfig = Utils.buildResource(operation, res.getPath(),
                        APIProcessUtils.convertProtoEndpointSecurity(res.getEndpointSecurityList()));
                resConfig.setPolicyConfig(Utils.genPolicyConfig(operation.getPolicies()));
                resConfig.setEndpoints(Utils.processEndpoints(res.getEndpoints()));
                resources.add(resConfig);
            }
        }

        KeyStore trustStore;
        try {
            trustStore = MtlsUtils.createTrustStore(api.getClientCertificatesList());
        } catch (KeyStoreException e) {
            throw new SecurityException(e);
        }

        BackendJWTTokenInfo backendJWTTokenInfo = api.getBackendJWTTokenInfo();
        JWTConfigurationDto jwtConfigurationDto = new JWTConfigurationDto();

        // If backendJWTTokeInfo is available
        if (api.hasBackendJWTTokenInfo()) {
            Map<String, Claim> claims = backendJWTTokenInfo.getCustomClaimsMap();
            Map<String, ClaimValueDTO> claimsMap = new HashMap<>();
            for (Map.Entry<String, Claim> claimEntry : claims.entrySet()) {
                Claim claim = claimEntry.getValue();
                ClaimValueDTO claimVal = new ClaimValueDTO(claim.getValue(), claim.getType());
                claimsMap.put(claimEntry.getKey(), claimVal);
            }
            EnforcerConfig enforcerConfig = ConfigHolder.getInstance().getConfig();
            jwtConfigurationDto.populateConfigValues(backendJWTTokenInfo.getEnabled(),
                    backendJWTTokenInfo.getHeader(), backendJWTTokenInfo.getSigningAlgorithm(),
                    backendJWTTokenInfo.getEncoding(), enforcerConfig.getJwtConfigurationDto().getPublicCert(),
                    enforcerConfig.getJwtConfigurationDto().getPrivateKey(), backendJWTTokenInfo.getTokenTTL(),
                    claimsMap, enforcerConfig.getJwtConfigurationDto().useKid(),
                    enforcerConfig.getJwtConfigurationDto().getKidValue());
        }

        org.wso2.apk.enforcer.commons.model.AIProvider aiProvider = new org.wso2.apk.enforcer.commons.model.AIProvider();
        if (api.getAiprovider().getEnabled()) {
            aiProvider.setProviderName(api.getAiprovider().getProviderName());
            aiProvider.setProviderAPIVersion(api.getAiprovider().getProviderAPIVersion());
            aiProvider.setOrganization(api.getAiprovider().getOrganization());
            aiProvider.setEnabled(api.getAiprovider().getEnabled());
            aiProvider.setModel(new org.wso2.apk.enforcer.commons.model.ValueDetails(
                    api.getAiprovider().getModel().getIn(), api.getAiprovider().getModel().getValue()));
            aiProvider.setPromptTokens(new org.wso2.apk.enforcer.commons.model.ValueDetails(
                    api.getAiprovider().getPromptTokens().getIn(), api.getAiprovider().getPromptTokens().getValue()));
            aiProvider.setCompletionToken(new org.wso2.apk.enforcer.commons.model.ValueDetails(
                    api.getAiprovider().getCompletionToken().getIn(), api.getAiprovider().getCompletionToken().getValue()));
            aiProvider.setTotalToken(new org.wso2.apk.enforcer.commons.model.ValueDetails(
                    api.getAiprovider().getTotalToken().getIn(), api.getAiprovider().getTotalToken().getValue()));
        }
        byte[] apiDefinition = api.getApiDefinitionFile().toByteArray();

        this.apiConfig = new APIConfig.Builder(name).uuid(api.getId()).vhost(vhost).basePath(basePath).version(version)
                .resources(resources).apiType(apiType).tier(api.getTier())
                .envType(api.getEnvType()).disableAuthentication(api.getDisableAuthentications())
                .disableScopes(api.getDisableScopes()).trustStore(trustStore).organizationId(api.getOrganizationId())
                .mutualSSL(mutualSSL).systemAPI(api.getSystemAPI())
                .applicationSecurity(applicationSecurity).jwtConfigurationDto(jwtConfigurationDto)
                .apiDefinition(apiDefinition).environment(api.getEnvironment())
                .subscriptionValidation(api.getSubscriptionValidation()).transportSecurity(api.getTransportSecurity())
                .aiProvider(aiProvider)
                .build();

        initFilters();
        return basePath;
    }

    @Override
    public ResponseObject process(RequestContext requestContext) {

        ResponseObject responseObject = new ResponseObject(requestContext.getRequestID());
        responseObject.setRequestPath(requestContext.getRequestPath());
        boolean analyticsEnabled = ConfigHolder.getInstance().getConfig().getAnalyticsConfig().isEnabled();

        Utils.handleCommonHeaders(requestContext);
        boolean isExistsMatchedResourcePath = requestContext.getMatchedResourcePaths() != null &&
                requestContext.getMatchedResourcePaths().size() > 0;
        // This flag is used to apply CORS filter
        boolean isOptionCall = requestContext.getRequestMethod().contains(HttpConstants.OPTIONS);
        if (!isExistsMatchedResourcePath && !isOptionCall) {
            // handle other not allowed non option calls
            requestContext.getProperties()
                    .put(APIConstants.MessageFormat.STATUS_CODE, APIConstants.StatusCodes.NOTFOUND.getCode());
            requestContext.getProperties().put(APIConstants.MessageFormat.ERROR_CODE,
                    APIConstants.StatusCodes.NOTFOUND.getValue());
            requestContext.getProperties().put(APIConstants.MessageFormat.ERROR_MESSAGE,
                    APIConstants.NOT_FOUND_MESSAGE);
            requestContext.getProperties().put(APIConstants.MessageFormat.ERROR_DESCRIPTION,
                    APIConstants.NOT_FOUND_DESCRIPTION);
        }
        if ((isExistsMatchedResourcePath || isOptionCall) && executeFilterChain(requestContext)) {
            EndpointUtils.updateClusterHeaderAndCheckEnv(requestContext);
            responseObject.setOrganizationId(requestContext.getMatchedAPI().getOrganizationId());
            responseObject.setRemoveHeaderMap(requestContext.getRemoveHeaders());
            responseObject.setQueryParamsToRemove(requestContext.getQueryParamsToRemove());
            responseObject.setRemoveAllQueryParams(requestContext.isRemoveAllQueryParams());
            responseObject.setQueryParamsToAdd(requestContext.getQueryParamsToAdd());
            responseObject.setQueryParamMap(requestContext.getQueryParameters());
            responseObject.setStatusCode(APIConstants.StatusCodes.OK.getCode());
            if (requestContext.getAddHeaders() != null && requestContext.getAddHeaders().size() > 0) {
                responseObject.setHeaderMap(requestContext.getAddHeaders());
            }
            if (analyticsEnabled) {
                AnalyticsFilter.getInstance().handleSuccessRequest(requestContext);
            }
            // set metadata for interceptors
            responseObject.setMetaDataMap(requestContext.getMetadataMap());
            if (requestContext.getMatchedAPI().isMockedApi()) {
                MockImplUtils.processMockedApiCall(requestContext, responseObject);
                return responseObject;
            }
        } else {
            // If enforcer stops with a false, it will be passed directly to the client.
            responseObject.setDirectResponse(true);
            responseObject.setStatusCode(Integer.parseInt(
                    requestContext.getProperties().get(APIConstants.MessageFormat.STATUS_CODE).toString()));
            if (requestContext.getProperties().get(APIConstants.MessageFormat.ERROR_CODE) != null) {
                responseObject.setErrorCode(
                        requestContext.getProperties().get(APIConstants.MessageFormat.ERROR_CODE).toString());
            }
            if (requestContext.getProperties().get(APIConstants.MessageFormat.ERROR_MESSAGE) != null) {
                responseObject.setErrorMessage(requestContext.getProperties()
                        .get(APIConstants.MessageFormat.ERROR_MESSAGE).toString());
            }
            if (requestContext.getProperties().get(APIConstants.MessageFormat.ERROR_DESCRIPTION) != null) {
                responseObject.setErrorDescription(requestContext.getProperties()
                        .get(APIConstants.MessageFormat.ERROR_DESCRIPTION).toString());
            }
            if (requestContext.getAddHeaders() != null && requestContext.getAddHeaders().size() > 0) {
                responseObject.setHeaderMap(requestContext.getAddHeaders());
            }
            if (analyticsEnabled && !FilterUtils.isSkippedAnalyticsFaultEvent(responseObject.getErrorCode())) {
                AnalyticsFilter.getInstance().handleFailureRequest(requestContext);
                responseObject.setMetaDataMap(new HashMap<>(0));
            }
        }

        return responseObject;
    }

    @Override
    public APIConfig getAPIConfig() {

        return this.apiConfig;
    }

    private void initFilters() {

        AuthFilter authFilter = new AuthFilter();
        authFilter.init(apiConfig, null);
        this.filters.add(authFilter);

        if (!apiConfig.isSystemAPI()) {
            MediationPolicyFilter mediationPolicyFilter = new MediationPolicyFilter();
            this.filters.add(mediationPolicyFilter);
        }

        // CORS filter is added as the first filter, and it is not customizable.
        CorsFilter corsFilter = new CorsFilter();
        this.filters.add(0, corsFilter);
    }
}
