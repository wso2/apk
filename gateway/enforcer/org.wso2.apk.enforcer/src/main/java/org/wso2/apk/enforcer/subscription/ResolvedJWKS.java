package org.wso2.apk.enforcer.subscription;

public class ResolvedJWKS {

    private String url;
    private ResolvedCertificate tls;

    public String getUrl() {

        return url;
    }

    public void setUrl(String url) {

        this.url = url;
    }

    public ResolvedCertificate getTls() {

        return tls;
    }

    public void setTls(ResolvedCertificate tls) {

        this.tls = tls;
    }
}
