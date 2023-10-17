package org.wso2.apk.enforcer.config.dto;

import java.util.Map;

public class AnalyticsPublisherConfigDTO {

    private boolean isEnabled;
    private String type;
    private Map<String, String> configProperties;

    public boolean isEnabled() {

        return isEnabled;
    }

    public void setEnabled(boolean enabled) {

        isEnabled = enabled;
    }

    public String getType() {

        return type;
    }

    public void setType(String type) {

        this.type = type;
    }

    public Map<String, String> getConfigProperties() {

        return configProperties;
    }

    public void setConfigProperties(Map<String, String> configProperties) {

        this.configProperties = configProperties;
    }

    public AnalyticsPublisherConfigDTO(boolean isEnabled, String type, Map<String, String> configProperties) {

        this.isEnabled = isEnabled;
        this.type = type;
        this.configProperties = configProperties;
    }
}
