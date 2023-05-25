package org.wso2.apk.enforcer.commons.model;

public class APIKeyAuthenticationConfig {
    private String In;
    private String Name;
    private boolean sendTokenToUpstream;

    public String getIn() {
        return In;
    }

    public void setIn(String in) {
        In = in;
    }

    public String getName() {
        return Name;
    }

    public void setName(String name) {
        Name = name;
    }

    public boolean isSendTokenToUpstream() {
        return sendTokenToUpstream;
    }

    public void setSendTokenToUpstream(boolean sendTokenToUpstream) {
        this.sendTokenToUpstream = sendTokenToUpstream;
    }
}
