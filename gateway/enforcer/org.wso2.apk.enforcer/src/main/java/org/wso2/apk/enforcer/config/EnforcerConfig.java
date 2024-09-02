/*
 * Copyright (c) 2021, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.wso2.apk.enforcer.config;

import org.wso2.apk.enforcer.commons.dto.JWTConfigurationDto;
import org.wso2.apk.enforcer.commons.jwttransformer.DefaultJWTTransformer;
import org.wso2.apk.enforcer.commons.jwttransformer.JWTTransformer;
import org.wso2.apk.enforcer.config.dto.APIKeyIssuerDto;
import org.wso2.apk.enforcer.config.dto.AnalyticsDTO;
import org.wso2.apk.enforcer.config.dto.AuthServiceConfigurationDto;
import org.wso2.apk.enforcer.config.dto.CacheDto;
import org.wso2.apk.enforcer.config.dto.ClientConfigDto;
import org.wso2.apk.enforcer.config.dto.FilterDTO;
import org.wso2.apk.enforcer.config.dto.ManagementCredentialsDto;
import org.wso2.apk.enforcer.config.dto.MetricsDTO;
import org.wso2.apk.enforcer.config.dto.MutualSSLDto;
import org.wso2.apk.enforcer.config.dto.SoapErrorResponseConfigDto;
import org.wso2.apk.enforcer.config.dto.TracingDTO;
import org.wso2.apk.enforcer.jwks.BackendJWKSDto;

import java.util.HashMap;
import java.util.Map;

/**
 * Configuration holder class for Microgateway.
 */
public class EnforcerConfig {

    private AuthServiceConfigurationDto authService;
    private TracingDTO tracingConfig;
    private MetricsDTO metricsConfig;
    private JWTConfigurationDto jwtConfigurationDto;
    private BackendJWKSDto backendJWKSDto;

    private APIKeyIssuerDto apiKeyIssuerDto;
    private APIKeyIssuerDto runtimeTokenIssuerDto;
    private CacheDto cacheDto;
    private String publicCertificatePath = "";
    private String privateKeyPath = "";
    private AnalyticsDTO analyticsConfig;
    private final Map<String, JWTTransformer> jwtTransformerMap = new HashMap<>();
    private MutualSSLDto mtlsInfo;
    private ManagementCredentialsDto management;
    private FilterDTO[] customFilters;

    private SoapErrorResponseConfigDto soapErrorResponseConfigDto;
    private boolean mandateSubscriptionValidation;
    private boolean mandateInternalKeyValidation;
    private ClientConfigDto httpClientConfigDto;
    private boolean enableGatewayClassController;

    public void setEnableGatewayClassController(Boolean enableGatewayClassController) {

        this.enableGatewayClassController = enableGatewayClassController;
    }

    public Boolean  getEnableGatewayClassController() {

        return enableGatewayClassController;
    }
    public ClientConfigDto getHttpClientConfigDto() {

        return httpClientConfigDto;
    }

    public void setHttpClientConfigDto(ClientConfigDto httpClientConfigDto) {

        this.httpClientConfigDto = httpClientConfigDto;
    }

    public AuthServiceConfigurationDto getAuthService() {
        return authService;
    }

    public void setAuthService(AuthServiceConfigurationDto authService) {
        this.authService = authService;
    }



    public void setTracingConfig(TracingDTO tracingConfig) {
        this.tracingConfig = tracingConfig;
    }

    public TracingDTO getTracingConfig() {
        return tracingConfig;
    }

    public MetricsDTO getMetricsConfig() {
        return metricsConfig;
    }

    public void setMetricsConfig(MetricsDTO metricsConfig) {
        this.metricsConfig = metricsConfig;
    }

    public void setJwtConfigurationDto(JWTConfigurationDto jwtConfigurationDto) {
        this.jwtConfigurationDto = jwtConfigurationDto;
    }

    public JWTConfigurationDto getJwtConfigurationDto() {
        return jwtConfigurationDto;
    }

    public CacheDto getCacheDto() {
        return cacheDto;
    }

    public void setCacheDto(CacheDto cacheDto) {
        this.cacheDto = cacheDto;
    }

    public void setPublicCertificatePath(String certPath) {
        this.publicCertificatePath = certPath;
    }

    public String getPublicCertificatePath() {
        return publicCertificatePath;
    }

    public void setPrivateKeyPath(String keyPath) {
        this.privateKeyPath = keyPath;
    }

    public String getPrivateKeyPath() {
        return privateKeyPath;
    }

    public AnalyticsDTO getAnalyticsConfig() {
        return analyticsConfig;
    }

    public void setAnalyticsConfig(AnalyticsDTO analyticsConfig) {
        this.analyticsConfig = analyticsConfig;
    }

    public JWTTransformer getJwtTransformer(String issuer) {
        if (jwtTransformerMap.containsKey(issuer)) {
            return jwtTransformerMap.get(issuer);
        }
        synchronized (jwtTransformerMap) {
            // check the map again, if two threads blocks and one add the default one
            // so the next thread also check if the default added by previous one
            if (jwtTransformerMap.containsKey(issuer)) {
                return jwtTransformerMap.get(issuer);
            }
            JWTTransformer defaultJWTTransformer = new DefaultJWTTransformer();
            jwtTransformerMap.put(issuer, defaultJWTTransformer);
            return defaultJWTTransformer;
        }
    }

    public void setJwtTransformers(Map<String, JWTTransformer> jwtTransformerMap) {
        this.jwtTransformerMap.putAll(jwtTransformerMap);
    }
    public MutualSSLDto getMtlsInfo() {
        return mtlsInfo;
    }

    public void setMtlsInfo(MutualSSLDto mtlsInfo) {
        this.mtlsInfo = mtlsInfo;
    }

    public ManagementCredentialsDto getManagement() {
        return management;
    }

    public void setManagement(ManagementCredentialsDto management) {
        this.management = management;
    }

    public SoapErrorResponseConfigDto getSoapErrorResponseConfigDto() {
        return soapErrorResponseConfigDto;
    }

    public void setSoapErrorResponseConfigDto(SoapErrorResponseConfigDto soapErrorResponseConfigDto) {
        this.soapErrorResponseConfigDto = soapErrorResponseConfigDto;
    }

    public FilterDTO[] getCustomFilters() {
        return customFilters;
    }

    public void setCustomFilters(FilterDTO[] customFilters) {
        this.customFilters = customFilters;
    }

    public APIKeyIssuerDto getApiKeyIssuerDto() {

        return apiKeyIssuerDto;
    }

    public void setApiKeyIssuerDto(APIKeyIssuerDto apiKeyIssuerDto) {

        this.apiKeyIssuerDto = apiKeyIssuerDto;
    }

    public APIKeyIssuerDto getRuntimeTokenIssuerDto() {

        return runtimeTokenIssuerDto;
    }

    public void setRuntimeTokenIssuerDto(APIKeyIssuerDto runtimeTokenIssuerDto) {

        this.runtimeTokenIssuerDto = runtimeTokenIssuerDto;
    }
    public BackendJWKSDto getBackendJWKSDto() {
        return backendJWKSDto;
    }

    public void setBackendJWKSDto(BackendJWKSDto backendJWKSDto) {
        this.backendJWKSDto = backendJWKSDto;
    }

    public boolean getMandateSubscriptionValidation() {
        return mandateSubscriptionValidation;
    }

    public void setMandateSubscriptionValidation(boolean mandateSubscriptionValidation) {
        this.mandateSubscriptionValidation = mandateSubscriptionValidation;
    }

    public boolean getMandateInternalKeyValidation() {
        return mandateInternalKeyValidation;
    }

    public void setMandateInternalKeyValidation(boolean mandateInternalKeyValidation) {
        this.mandateInternalKeyValidation = mandateInternalKeyValidation;
    }
}

