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

package org.wso2.apk.enforcer.commons.model;

import java.util.List;

public class AuthenticationConfig {
    private JWTAuthenticationConfig jwtAuthenticationConfig;
    private List<APIKeyAuthenticationConfig> apiKeyAuthenticationConfigs;
    private InternalKeyConfig internalKeyConfig;
    private Oauth2AuthenticationConfig oauth2AuthenticationConfig;
    private boolean Disabled;

    public JWTAuthenticationConfig getJwtAuthenticationConfig() {
        return jwtAuthenticationConfig;
    }

    public void setJwtAuthenticationConfig(JWTAuthenticationConfig jwtAuthenticationConfig) {
        this.jwtAuthenticationConfig = jwtAuthenticationConfig;
    }

    public List<APIKeyAuthenticationConfig> getApiKeyAuthenticationConfigs() {
        return apiKeyAuthenticationConfigs;
    }

    public void setApiKeyAuthenticationConfigs(List<APIKeyAuthenticationConfig> apiKeyAuthenticationConfigs) {
        this.apiKeyAuthenticationConfigs = apiKeyAuthenticationConfigs;
    }

    public boolean isDisabled() {
        return Disabled;
    }

    public void setDisabled(boolean disabled) {
        Disabled = disabled;
    }

    public InternalKeyConfig getInternalKeyConfig() {
        return internalKeyConfig;
    }

    public void setInternalKeyConfig(InternalKeyConfig internalKeyConfig) {
        this.internalKeyConfig = internalKeyConfig;
    }

    public Oauth2AuthenticationConfig getOauth2AuthenticationConfig() {

        return oauth2AuthenticationConfig;
    }

    public void setOauth2AuthenticationConfig(Oauth2AuthenticationConfig oauth2AuthenticationConfig) {

        this.oauth2AuthenticationConfig = oauth2AuthenticationConfig;
    }
}
