package org.wso2.apk.enforcer.commons.model;

public class JWTAuthenticationConfig {
    private String Header;
    private boolean sendTokenToUpstream;

    public String getHeader() {
        return Header;
    }

    public void setHeader(String header) {
        Header = header;
    }

    public boolean isSendTokenToUpstream() {
        return sendTokenToUpstream;
    }

    public void setSendTokenToUpstream(boolean sendTokenToUpstream) {
        this.sendTokenToUpstream = sendTokenToUpstream;
    }
}
