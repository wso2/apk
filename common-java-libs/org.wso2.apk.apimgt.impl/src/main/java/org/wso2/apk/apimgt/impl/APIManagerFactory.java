/*
 * Copyright (c) 2022, WSO2 LLC. (http://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.wso2.apk.apimgt.impl;

import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.wso2.apk.apimgt.api.*;
import org.wso2.apk.apimgt.impl.utils.LRUCache;
import org.wso2.apk.apimgt.user.mgt.UserConstants;

import java.util.Map;

public class APIManagerFactory {

    private static final Log log = LogFactory.getLog(APIManagerFactory.class);

    private static final String ANONYMOUS_USER = "__wso2.am.anon__";

    private static final APIManagerFactory instance = new APIManagerFactory();

    private APIManagerCache<APIProvider> providers = new APIManagerCache<APIProvider>(50);

    private APIManagerFactory() {

    }

    private APIProvider newProvider(String username) throws APIManagementException {
        //TODO: APK
//        return new UserAwareAPIProvider(username);
        return new APIProviderImpl(username);
    }

    public APIProvider getAPIProvider(String username) throws APIManagementException {
        return new APIProviderImpl(username);
    }

    public static APIManagerFactory getInstance() {
        return instance;
    }




    public void clearAll() {

        providers.exclusiveLock();
        try {
            for (APIProvider provider : providers.values()) {
                cleanupSilently(provider);
            }
            providers.clear();
        } finally {
            providers.release();
        }
    }

    private void cleanupSilently(APIManager manager) {
        if (manager != null) {
            try {
                manager.cleanup();
            } catch (APIManagementException ignore) {

            }
        }
    }

    private class APIManagerCache<T> extends LRUCache<String,T> {

        public APIManagerCache(int maxEntries) {
            super(maxEntries);
        }

        protected void handleRemovableEntry(Map.Entry<String,T> entry) {
            try {
                ((APIManager) entry.getValue()).cleanup();
            } catch (APIManagementException e) {
                log.warn("Error while cleaning up APIManager instance", e);
            }
        }
    }
}
