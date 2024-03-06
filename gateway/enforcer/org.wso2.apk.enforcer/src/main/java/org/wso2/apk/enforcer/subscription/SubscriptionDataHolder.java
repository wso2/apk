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

package org.wso2.apk.enforcer.subscription;

import java.util.HashMap;
import java.util.Map;

/**
 * This class holds tenant wise subscription data stores.
 */
public class SubscriptionDataHolder {

    private static final SubscriptionDataHolder instance = new SubscriptionDataHolder();
    Map<String, SubscriptionDataStore> subscriptionDataStoreMap = new HashMap<>();

    public static SubscriptionDataHolder getInstance() {
        return instance;
    }

    public SubscriptionDataStore getSubscriptionDataStore(String organization) {
        return subscriptionDataStoreMap.get(organization);
    }

    public SubscriptionDataStore initializeSubscriptionDataStore(String organization) {
        SubscriptionDataStore subscriptionDataStore = new SubscriptionDataStoreImpl();
        subscriptionDataStoreMap.put(organization, subscriptionDataStore);
        return subscriptionDataStore;
    }

    public int getTotalSubscriptionCount() {
        int totalSubCount = 0;
        for (SubscriptionDataStore store : subscriptionDataStoreMap.values()) {
                totalSubCount += store.getSubscriptionCount();
        }
        return totalSubCount;
    }

    public int getTotalJWTIssuerCount() {
        int totalJWTIssuerCount = 0;
        for (SubscriptionDataStore store : subscriptionDataStoreMap.values()) {
                totalJWTIssuerCount += store.getJWTIssuerCount();
            }
        return totalJWTIssuerCount;
    }

}
