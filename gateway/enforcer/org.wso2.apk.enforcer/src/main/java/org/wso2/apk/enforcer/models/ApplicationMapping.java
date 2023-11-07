/*
 * Copyright (c) 2023, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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
 * Entity for keeping mapping between Application and Subscription.
 */
public class ApplicationMapping implements CacheableEntity<String> {

    private String uuid = null;
    private String applicationRef = null;
    private String subscriptionRef = null;

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

    public String getCacheKey() {
        return uuid;
    }

    @Override
    public String toString() {
        return "ApplicationMapping [uuid=" + uuid + ", applicationRef=" + applicationRef + ", subscriptionRef="
                + subscriptionRef + "]";
    }
}
