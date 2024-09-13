package org.wso2.apk.enforcer.subscription;

import java.io.Serializable;

/**
 * Entity for keeping Application related information. Represents an Application in APK.
 */
public class SubscriptionDto implements Serializable {

    private static final long serialVersionUID = 1L;
    private String uuid;
    private String organization;
    private String subStatus;
    private SubscribedAPIDto subscribedApi;
    private String ratelimitTier;

    public String getUuid() {

        return uuid;
    }

    public void setUuid(String uuid) {

        this.uuid = uuid;
    }

    public String getOrganization() {

        return organization;
    }

    public void setOrganization(String organization) {

        this.organization = organization;
    }

    public String getSubStatus() {

        return subStatus;
    }

    public void setSubStatus(String subStatus) {

        this.subStatus = subStatus;
    }

    public SubscribedAPIDto getSubscribedApi() {

        return subscribedApi;
    }

    public void setSubscribedApi(SubscribedAPIDto subscribedApi) {

        this.subscribedApi = subscribedApi;
    }

    public String getRatelimitTier() {

        return ratelimitTier;
    }

    public void setRatelimitTier(String ratelimitTier) {

        this.ratelimitTier = ratelimitTier;
    }
}
