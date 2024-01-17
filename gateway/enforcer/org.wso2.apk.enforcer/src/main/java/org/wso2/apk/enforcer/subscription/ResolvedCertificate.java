package org.wso2.apk.enforcer.subscription;

public class ResolvedCertificate {

    private String resolvedCertificate;
    private String[] allowedSANs;

    public String getResolvedCertificate() {

        return resolvedCertificate;
    }

    public void setResolvedCertificate(String resolvedCertificate) {

        this.resolvedCertificate = resolvedCertificate;
    }

    public String[] getAllowedSANs() {

        return allowedSANs;
    }

    public void setAllowedSANs(String[] allowedSANs) {

        this.allowedSANs = allowedSANs;
    }
}
