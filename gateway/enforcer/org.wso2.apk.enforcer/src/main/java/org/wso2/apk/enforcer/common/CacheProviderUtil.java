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
package org.wso2.apk.enforcer.common;

import org.wso2.apk.enforcer.api.API;

import java.util.HashSet;
import java.util.Map;
import java.util.Set;
import java.util.concurrent.ConcurrentHashMap;

/**
 * This Class used to keep Organization level cache maps.
 */
public class CacheProviderUtil {

    private final static Map<String, CacheProvider> organizationCache = new ConcurrentHashMap<>();

    private CacheProviderUtil() {

    }

    private static void initializeOrgLevelCache(String organization) {

        String syncKey = organization + "-InitializeCache";
        if (!organizationCache.containsKey(organization)) {
            synchronized (syncKey.intern()) {
                if (!organizationCache.containsKey(organization)) {
                    CacheProvider cacheProvider = new CacheProvider();
                    cacheProvider.init();
                    organizationCache.put(organization, cacheProvider);
                }
            }
        }
    }

    /**
     * initialize Cache map from APIS
     * @param apis APIS available in cluster
     */
    public static void initializeCacheHolder(Map<String, API> apis) {

        Set<String> organizations = new HashSet<>();
        for (Map.Entry<String, API> api : apis.entrySet()) {
            organizations.add(api.getValue().getAPIConfig().getOrganizationId());
        }
        for (String organization : organizations) {
            initializeOrgLevelCache(organization);
        }
        organizationCache.keySet().removeIf(organization -> !organizations.contains(organization));
    }

    /**
     * This method used to get Organization level cache.
     * @param organization organization Name
     * @return CacheProvider if exists
     */
    public static CacheProvider getOrganizationCache(String organization) {

        return organizationCache.get(organization);
    }
}
