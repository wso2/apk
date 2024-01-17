package org.wso2.apk.enforcer.subscription;

public class SignatureValidationRestDto {

    private ResolvedJWKS jwks;
    private ResolvedCertificate certificate;

    public ResolvedJWKS getJwks() {

        return jwks;
    }

    public void setJwks(ResolvedJWKS jwks) {

        this.jwks = jwks;
    }

    public ResolvedCertificate getCertificate() {

        return certificate;
    }

    public void setCertificate(ResolvedCertificate certificate) {

        this.certificate = certificate;
    }
}
