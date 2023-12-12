package org.wso2.apk.enforcer.subscription;

import java.io.Serializable;

public class ApplicationKeyMappingDTO implements Serializable {

    public String getApplicationUUID() {

        return applicationUUID;
    }

    public void setApplicationUUID(String applicationUUID) {

        this.applicationUUID = applicationUUID;
    }

    public String getSecurityScheme() {

        return securityScheme;
    }

    public void setSecurityScheme(String securityScheme) {

        this.securityScheme = securityScheme;
    }

    public String getApplicationIdentifier() {

        return applicationIdentifier;
    }

    public void setApplicationIdentifier(String applicationIdentifier) {

        this.applicationIdentifier = applicationIdentifier;
    }

    public String getKeyType() {

        return keyType;
    }

    public void setKeyType(String keyType) {

        this.keyType = keyType;
    }

    public String getEnvID() {

        return envID;
    }

    public void setEnvID(String envID) {

        this.envID = envID;
    }

    private String applicationUUID;
    private String securityScheme;
    private String applicationIdentifier;
    private String keyType;
    private String envID;

    public String getOrganizationId() {

        return organizationId;
    }

    public void setOrganizationId(String organizationId) {

        this.organizationId = organizationId;
    }

    private String organizationId;

}
