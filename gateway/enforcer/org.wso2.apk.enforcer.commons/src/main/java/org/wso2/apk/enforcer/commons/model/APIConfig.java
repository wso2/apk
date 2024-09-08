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
package org.wso2.apk.enforcer.commons.model;

import org.wso2.apk.enforcer.commons.dto.JWTConfigurationDto;

import java.security.KeyStore;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * APIConfig contains the details related to the MatchedAPI for the inbound request.
 */
public class APIConfig {
    private String name;
    private String version;
    private String vhost;
    private String basePath;
    private String apiType;
//    private Map<String, EndpointCluster> endpoints; // "PRODUCTION" OR "SANDBOX" -> endpoint cluster
    private String envType;
    private String apiLifeCycleState;
    private String authorizationHeader;
    private String organizationId;
    private String uuid;
    private String tier;
    private boolean disableAuthentication;
    private boolean disableScopes;
    private List<ResourceConfig> resources = new ArrayList<>();
    private boolean isMockedApi;
    private KeyStore trustStore;
    private String mutualSSL;
    private boolean transportSecurity;
    private boolean applicationSecurity;
    private GraphQLSchemaDTO graphQLSchemaDTO;
    private JWTConfigurationDto jwtConfigurationDto;
    private boolean systemAPI;
    private byte[] apiDefinition;
    private String environment;
    private boolean subscriptionValidation;
    private EndpointSecurity[] endpointSecurity;
    private EndpointCluster endpoints;
    private AIProvider aiProvider;

    public AIProvider getAiProvider() {
        return aiProvider;
    }

    public void setAiProvider(AIProvider aiProvider) {
        this.aiProvider = aiProvider;
    }

    public EndpointCluster getEndpoints() {
        return endpoints;
    }

    public void getEndpointSecurity(EndpointSecurity[] endpointSecurity) {
        this.endpointSecurity = endpointSecurity;
    }

    /**
     * getApiType returns the API type. This could be one of the following.
     * HTTP, WS, WEBHOOK
     *
     * @return the apiType
     */
    public String getApiType() {
        return apiType;
    }

    /**
     * getEnvType returns the API's env type
     * whether the key type is production or sandbox.
     *
     * @return getEnvType returns type of the env. Production or Sandbox
     */
    public String getEnvType() {
        return envType;
    }

    /**
     * Corresponding API's organization UUID (TenantDomain) is returned.
     *
     * @return Organization UUID
     */
    public String getOrganizationId() {
        return organizationId;
    }

    /**
     * Corresponding API's API UUID is returned.
     * @return API UUID
     */
    public String getUuid() {
        return uuid;
    }

    /**
     * Corresponding API's API Name is returned.
     * @return API name
     */
    public String getName() {
        return name;
    }

    /**
     * Corresponding API's API Version is returned.
     * @return API version
     */
    public String getVersion() {
        return version;
    }

    /**
     * Corresponding API's Host is returned.
     * @return API's host
     */
    public String getVhost() {
        return vhost;
    }

    /**
     * Corresponding API's Basepath is returned.
     * @return basePath of the API
     */
    public String getBasePath() {
        return basePath;
    }

    /**
     * Current API Lifecycle state is returned.
     * @return lifecycle state
     */
    public String getApiLifeCycleState() {
        return apiLifeCycleState;
    }

    /**
     * API level Throttling tier assigned for the corresponding API.
     * @return API level throttling tier
     */
    public String getTier() {
        return tier;
    }

    /**
     * If the authentication is disabled for the API .
     *
     * @return true if the authentication is disabled for the API.
     */
    public boolean isDisableAuthentication() {
        return disableAuthentication;
    }

    /**
     * If the scopes are disabled for the API .
     *
     * @return true if the scopes are disabled for the API.
     */
    public boolean isDisableScopes() {
        return disableScopes;
    }

    /**
     * Returns the complete list of resources under the corresponding API.
     * Each operation in the openAPI definition is listed under here.
     * @return Resources of the API.
     */
    public List<ResourceConfig> getResources() {
        return resources;
    }

    /**
     * Returns whether a given API is a mocked API or not.
     *
     * @return boolean value to denote isMockedApi or not.
     */
    public boolean isMockedApi() {
        return isMockedApi;
    }

    /**
     * Returns the truststore for the corresponding API.
     *
     * @return TrustStore
     */
    public KeyStore getTrustStore() {
        return trustStore;
    }


    /**
     * Returns the mTLS optionality for the corresponding API.
     *
     * @return mTLS optionality
     */
    public String getMutualSSL() {
        return mutualSSL;
    }

    /**
     * Returns if transport security (mTLS) is enabled or disabled for the corresponding API.
     *
     * @return transportSecurity enabled
     */
    public boolean isTransportSecurity() {
        return transportSecurity;
    }

    /**
     * Returns the application security optionality for the corresponding API.
     *
     * @return application security optionality
     */
    public boolean getApplicationSecurity() {
        return applicationSecurity;
    }

    /**
     * Returns graphQL definitions and schemes.
     *
     * @return GraphQLSchemaDTO.
     */
    public GraphQLSchemaDTO getGraphQLSchemaDTO() {
        return graphQLSchemaDTO;
    }

    public boolean isSystemAPI() {
        return systemAPI;
    }

    /**
     * Returns the API definition.
     * @return byte array of the API definition.
     */
    public byte[] getApiDefinition() {
        return apiDefinition;
    }

    /**
     * Returns the subscription validation status.
     * @return true if subscription validation is enabled.
     */
    public boolean isSubscriptionValidation() {
        return subscriptionValidation;
    }

    public JWTConfigurationDto getJwtConfigurationDto() {
        return jwtConfigurationDto;
    }

    /**
     * Returns the environment of the API.
     * @return String.
     */
    public String getEnvironment() {
        return environment;
    }

    /**
     * Implements builder pattern to build an API Config object.
     */
    public static class Builder {

        private String name;
        private String version;
        private String vhost;
        private String basePath;
        private String apiType;
        private String envType;
        private String apiLifeCycleState;
        private String organizationId;
        private String uuid;
        private String tier;
        private boolean disableAuthentication;
        private boolean disableScopes;
        private List<ResourceConfig> resources = new ArrayList<>();
        private boolean isMockedApi;
        private KeyStore trustStore;
        private String mutualSSL;
        private boolean applicationSecurity;
        private GraphQLSchemaDTO graphQLSchemaDTO;
        private boolean systemAPI;
        private byte[] apiDefinition;
        private boolean subscriptionValidation;
        private JWTConfigurationDto jwtConfigurationDto;
        private String environment;
        private boolean transportSecurity;

        private AIProvider aiProvider;

        public Builder(String name) {
            this.name = name;
        }

        public Builder version(String version) {
            this.version = version;
            return this;
        }

        public Builder vhost(String vhost) {
            this.vhost = vhost;
            return this;
        }

        public Builder aiProvider(AIProvider aiProvider) {
            this.aiProvider = aiProvider;
            return this;
        }

        public Builder basePath(String basePath) {
            this.basePath = basePath;
            return this;
        }

        public Builder apiType(String apiType) {
            this.apiType = apiType;
            return this;
        }

        public Builder apiLifeCycleState(String apiLifeCycleState) {
            this.apiLifeCycleState = apiLifeCycleState;
            return this;
        }

        public Builder tier(String tier) {
            this.tier = tier;
            return this;
        }

        public Builder disableAuthentication(boolean enabled) {
            this.disableAuthentication = enabled;
            return this;
        }

        public Builder disableScopes(boolean enabled) {
            this.disableScopes = enabled;
            return this;
        }

        public Builder resources(List<ResourceConfig> resources) {
            this.resources = resources;
            return this;
        }

        public Builder envType(String envType) {
            this.envType = envType;
            return this;
        }

        public Builder organizationId(String organizationId) {
            this.organizationId = organizationId;
            return this;
        }

        public Builder uuid(String uuid) {
            this.uuid = uuid;
            return this;
        }

        public Builder graphQLSchemaDTO(GraphQLSchemaDTO graphQLSchemaDTO) {
            this.graphQLSchemaDTO = graphQLSchemaDTO;
            return this;
        }

        public Builder mockedApi(boolean isMockedApi) {
            this.isMockedApi = isMockedApi;
            return this;
        }

        public Builder trustStore(KeyStore trustStore) {
            this.trustStore = trustStore;
            return this;
        }

        public Builder mutualSSL(String mutualSSL) {
            this.mutualSSL = mutualSSL;
            return this;
        }

        public Builder applicationSecurity(boolean applicationSecurity) {
            this.applicationSecurity = applicationSecurity;
            return this;
        }

        public Builder systemAPI(boolean systemAPI) {
            this.systemAPI = systemAPI;
            return this;
        }

        public Builder jwtConfigurationDto(JWTConfigurationDto jwtConfigurationDto) {
            this.jwtConfigurationDto = jwtConfigurationDto;
            return this;
        }

        public Builder apiDefinition(byte[] apiDefinition) {
            this.apiDefinition = apiDefinition;
            return this;
        }

        public Builder environment(String environment) {
            this.environment = environment;
            return this;
        }

        public Builder subscriptionValidation(boolean subscriptionValidation) {
            this.subscriptionValidation = subscriptionValidation;
            return this;
        }

        public Builder transportSecurity(boolean transportSecurity) {
            this.transportSecurity = transportSecurity;
            return this;
        }

        public APIConfig build() {
            APIConfig apiConfig = new APIConfig();
            apiConfig.name = this.name;
            apiConfig.vhost = this.vhost;
            apiConfig.basePath = this.basePath;
            apiConfig.version = this.version;
            apiConfig.apiLifeCycleState = this.apiLifeCycleState;
            apiConfig.resources = this.resources;
            apiConfig.apiType = this.apiType;
            apiConfig.envType = this.envType;
            apiConfig.tier = this.tier;
            apiConfig.disableAuthentication = this.disableAuthentication;
            apiConfig.disableScopes = this.disableScopes;
            apiConfig.organizationId = this.organizationId;
            apiConfig.uuid = this.uuid;
            apiConfig.isMockedApi = this.isMockedApi;
            apiConfig.trustStore = this.trustStore;
            apiConfig.mutualSSL = this.mutualSSL;
            apiConfig.transportSecurity = this.transportSecurity;
            apiConfig.applicationSecurity = this.applicationSecurity;
            apiConfig.graphQLSchemaDTO = this.graphQLSchemaDTO;
            apiConfig.systemAPI = this.systemAPI;
            apiConfig.jwtConfigurationDto = this.jwtConfigurationDto;
            apiConfig.apiDefinition = this.apiDefinition;
            apiConfig.environment = this.environment;
            apiConfig.subscriptionValidation = this.subscriptionValidation;
            apiConfig.aiProvider = this.aiProvider;
            return apiConfig;
        }
    }

    private APIConfig() {
    }
}
