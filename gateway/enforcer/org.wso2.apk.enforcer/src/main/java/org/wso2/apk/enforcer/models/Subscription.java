/*
 * Copyright (c) 2020, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.wso2.apk.enforcer.models;

import org.wso2.apk.enforcer.common.CacheableEntity;

/**
 * Entity for representing a SubscriptionDTO in APIM.
 */
public class Subscription implements CacheableEntity<String> {

    private String subscriptionId = null;
    private String subscriptionStatus = null;
    private String organization = null;
    private SubscribedAPI subscribedApi = null;
    private String ratelimitTier = null;
    private long timeStamp;

    public String getSubscriptionId() {

        return subscriptionId;
    }

    public void setSubscriptionId(String subscriptionId) {

        this.subscriptionId = subscriptionId;
    }

    public String getSubscriptionStatus() {

        return subscriptionStatus;
    }

    public void setSubscriptionStatus(String subscriptionStatus) {

        this.subscriptionStatus = subscriptionStatus;
    }

    public String getOrganization() {

        return organization;
    }

    public void setOrganization(String organization) {

        this.organization = organization;
    }

    public SubscribedAPI getSubscribedApi() {

        return subscribedApi;
    }

    public void setSubscribedApi(SubscribedAPI subscribedApi) {

        this.subscribedApi = subscribedApi;
    }

    public long getTimeStamp() {

        return timeStamp;
    }

    public void setTimeStamp(long timeStamp) {

        this.timeStamp = timeStamp;
    }

    @Override
    public String getCacheKey() {

        return subscriptionId;
    }

    public String getRatelimitTier() {
        return ratelimitTier;
    }

    public void setRatelimitTier(String ratelimitTier) {
        this.ratelimitTier = ratelimitTier;
    }

    @Override
    public String toString() {

        return "Subscription{" +
                "subscriptionId='" + subscriptionId + '\'' +
                ", subscriptionStatus='" + subscriptionStatus + '\'' +
                ", organization='" + organization + '\'' +
                ", subscribedApi=" + subscribedApi +
                ", timeStamp=" + timeStamp +
                '}';
    }
}
