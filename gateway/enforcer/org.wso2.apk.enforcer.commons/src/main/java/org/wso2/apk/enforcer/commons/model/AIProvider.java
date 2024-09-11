package org.wso2.apk.enforcer.commons.model;

// AIProvider is used to provide the AI model to the enforcer
public class AIProvider {
    private String providerName;
    private String providerAPIVersion;
    private String organization;

    private Boolean enabled = false;

    private ValueDetails model;

    private ValueDetails promptTokens;

    private ValueDetails completionToken;

    private ValueDetails totalToken;

    // Get provider name
    public String getProviderName() {
        return providerName;
    }

    // Get provider API version
    public String getProviderAPIVersion() {
        return providerAPIVersion;
    }

    // Get enabled
    public Boolean getEnabled() {
        return enabled;
    }

    // Get organization
    public String getOrganization() {
        return organization;
    }

    // Get model
    public ValueDetails getModel() {
        return model;
    }

    // Get prompt tokens
    public ValueDetails getPromptTokens() {
        return promptTokens;
    }

    // Get completion token
    public ValueDetails getCompletionToken() {
        return completionToken;
    }

    // Get total token
    public ValueDetails getTotalToken() {
        return totalToken;
    }

    // Set provider name
    public void setProviderName(String providerName) {
        this.providerName = providerName;
    }

    // Set provider API version
    public void setProviderAPIVersion(String providerAPIVersion) {
        this.providerAPIVersion = providerAPIVersion;
    }

    // Set enabled
    public void setEnabled(Boolean enabled) {
        this.enabled = enabled;
    }

    // Set organization
    public void setOrganization(String organization) {
        this.organization = organization;
    }

    // Set model
    public void setModel(ValueDetails model) {
        this.model = model;
    }

    // Set prompt tokens
    public void setPromptTokens(ValueDetails promptTokens) {
        this.promptTokens = promptTokens;
    }

    // Set completion token
    public void setCompletionToken(ValueDetails completionToken) {
        this.completionToken = completionToken;
    }

    // Set total token
    public void setTotalToken(ValueDetails totalToken) {
        this.totalToken = totalToken;
    }

}