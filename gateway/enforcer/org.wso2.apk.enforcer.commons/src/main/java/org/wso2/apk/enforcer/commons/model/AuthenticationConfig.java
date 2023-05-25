package org.wso2.apk.enforcer.commons.model;

import java.util.List;

public class AuthenticationConfig {
    private JWTAuthenticationConfig jwtAuthenticationConfig;
    private List<APIKeyAuthenticationConfig> apiKeyAuthenticationConfigs;
    private InternalKeyConfig internalKeyConfig;
    private boolean Disabled;
    private boolean sendTokenToUpstream;

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

    public boolean isSendTokenToUpstream() {
        return sendTokenToUpstream;
    }

    public void setSendTokenToUpstream(boolean sendTokenToUpstream) {
        this.sendTokenToUpstream = sendTokenToUpstream;
    }

    public InternalKeyConfig getInternalKeyConfig() {
        return internalKeyConfig;
    }

    public void setInternalKeyConfig(InternalKeyConfig internalKeyConfig) {
        this.internalKeyConfig = internalKeyConfig;
    }
}
