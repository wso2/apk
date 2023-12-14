package org.wso2.apk.enforcer.subscription;

import java.io.Serializable;

public class ApplicationMappingDto implements Serializable {
    private String uuid;

    public String getUuid() {

        return uuid;
    }

    public void setUuid(String uuid) {

        this.uuid = uuid;
    }

    public String getApplicationRef() {

        return applicationRef;
    }

    public void setApplicationRef(String applicationRef) {

        this.applicationRef = applicationRef;
    }

    public String getSubscriptionRef() {

        return subscriptionRef;
    }

    public void setSubscriptionRef(String subscriptionRef) {

        this.subscriptionRef = subscriptionRef;
    }

    private String applicationRef;
    private String subscriptionRef;
    private String organizationId;

    public String getOrganizationId() {

        return organizationId;
    }

    public void setOrganizationId(String organizationId) {

        this.organizationId = organizationId;
    }
}
