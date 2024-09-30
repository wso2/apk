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
package org.wso2.apk.enforcer.commons.model;

import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * Holds the metadata related to the resources/operations of an API.
 */
public class ResourceConfig {

    private String path;
    private String matchID;
    private HttpMethods method;
    private String tier = "Unlimited";
    private EndpointCluster endpoints;
    private EndpointSecurity[] endpointSecurity;
    private PolicyConfig policyConfig;
    private MockedApiConfig mockedApiConfig;
    private AuthenticationConfig authenticationConfig;
    private String[] scopes;

    /**
     * ENUM to hold http operations.
     */
    public enum HttpMethods {
        GET("get"), POST("post"), PUT("put"), DELETE("delete"), HEAD("head"),
        PATCH("patch"), OPTIONS("options"), QUERY("query"), MUTATION("mutation"),
        SUBSCRIPTION("subscription"), GRPC("GRPC");

        private String value;

        private HttpMethods(String value) {
            this.value = value;
        }
    }

    /**
     * Get the matching path Template for the request.
     *
     * @return path Template
     */
    public String getMatchID() {
        return matchID;
    }

    public void setMatchID(String matchID) {
        this.matchID = matchID;
    }

    /**
     * Get the matching path Template for the request.
     *
     * @return path Template
     */
    public String getPath() {
        return path;
    }

    public void setPath(String path) {
        this.path = path;
    }

    /**
     * Get the matching HTTP Method.
     *
     * @return HTTP method
     */
    public HttpMethods getMethod() {
        return method;
    }

    public void setMethod(HttpMethods method) {
        this.method = method;
    }

    /**
     * Get the resource level throttling tier assigned for the corresponding Resource.
     *
     * @return resource level throttling tier
     */
    public String getTier() {
        return tier;
    }

    public void setTier(String tier) {
        this.tier = tier;
    }

    //todo(amali) this don't need to be a map

    /**
     * Get the resource level endpoint cluster map for the corresponding Resource
     * where the map-key is either "PRODUCTION" or "SANDBOX".
     *
     * @return resource level endpoint cluster map
     */
    public EndpointCluster getEndpoints() {
        return endpoints;
    }

    public void setEndpoints(EndpointCluster endpoints) {
        this.endpoints = endpoints;
    }

    /**
     * Provides mock API configurations defined in the JSON.
     *
     * @return MockedApiConfig object with all the configurations defined under operation of a mocked API.
     */
    public MockedApiConfig getMockedApiConfig() {
        return mockedApiConfig;
    }

    public void setMockApiConfig(MockedApiConfig mockedApiConfig) {
        this.mockedApiConfig = mockedApiConfig;
    }

    public AuthenticationConfig getAuthenticationConfig() {
        return authenticationConfig;
    }

    public void setAuthenticationConfig(AuthenticationConfig authenticationConfig) {
        this.authenticationConfig = authenticationConfig;
    }

    public String[] getScopes() {
        return scopes;
    }

    public void setScopes(String[] scopes) {
        this.scopes = scopes;
    }

    public PolicyConfig getPolicyConfig() {
        return policyConfig;
    }

    public void setPolicyConfig(PolicyConfig policyConfig) {
        this.policyConfig = policyConfig;
    }

    public EndpointSecurity[] getEndpointSecurity() {
        return endpointSecurity;
    }

    public void setEndpointSecurity(EndpointSecurity[] endpointSecurity) {
        this.endpointSecurity = endpointSecurity;
    }
}

