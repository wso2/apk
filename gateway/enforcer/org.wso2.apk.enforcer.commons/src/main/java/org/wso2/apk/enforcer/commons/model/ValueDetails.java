package org.wso2.apk.enforcer.commons.model;

// ValueDetails is used to provide the AI model to the enforcer
public class ValueDetails {
    private String in;
    private String value;

    public ValueDetails(String in, String value) {
        this.in = in;
        this.value = value;
    }

    public String getIn() {
        return in;
    }

    public String getValue() {
        return value;
    }

    public void setIn(String in) {
        this.in = in;
    }

    public void setValue(String value) {
        this.value = value;
    }
}
