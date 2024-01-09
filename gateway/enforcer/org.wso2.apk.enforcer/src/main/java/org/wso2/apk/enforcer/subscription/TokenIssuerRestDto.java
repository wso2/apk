package org.wso2.apk.enforcer.subscription;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class TokenIssuerRestDto {
    private String name;
    private String organization;
    private String issuer;
    private String consumerKeyClaim;
    private String scopesClaim;
    private List<String> environments = new ArrayList<>();
    private Map<String,String> claimMappings = new HashMap<>();
    private SignatureValidationRestDto signatureValidation;

    public String getName() {

        return name;
    }

    public void setName(String name) {

        this.name = name;
    }

    public String getOrganization() {

        return organization;
    }

    public void setOrganization(String organization) {

        this.organization = organization;
    }

    public String getIssuer() {

        return issuer;
    }

    public void setIssuer(String issuer) {

        this.issuer = issuer;
    }

    public String getConsumerKeyClaim() {

        return consumerKeyClaim;
    }

    public void setConsumerKeyClaim(String consumerKeyClaim) {

        this.consumerKeyClaim = consumerKeyClaim;
    }

    public String getScopesClaim() {

        return scopesClaim;
    }

    public void setScopesClaim(String scopesClaim) {

        this.scopesClaim = scopesClaim;
    }

    public List<String> getEnvironments() {

        return environments;
    }

    public void setEnvironments(List<String> environments) {

        this.environments = environments;
    }

    public Map<String, String> getClaimMappings() {

        return claimMappings;
    }

    public void setClaimMappings(Map<String, String> claimMappings) {

        this.claimMappings = claimMappings;
    }

    public SignatureValidationRestDto getSignatureValidation() {

        return signatureValidation;
    }

    public void setSignatureValidation(SignatureValidationRestDto signatureValidation) {

        this.signatureValidation = signatureValidation;
    }
}
