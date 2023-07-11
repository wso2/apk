package org.wso2.apk.enforcer.commons.dto;

public class ClaimValueDTO {
    private Object value;
    private String type;

    public ClaimValueDTO(final Object value, final String type) {
        this.value = value;
        this.type = type;
    }
    public Object getValue() {
        return value;
    }

    public String getType() {
        return type;
    }
}
