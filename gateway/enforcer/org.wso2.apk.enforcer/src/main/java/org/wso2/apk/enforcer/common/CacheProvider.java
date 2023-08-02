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

package org.wso2.apk.enforcer.common;

import com.google.common.cache.CacheBuilder;
import com.google.common.cache.CacheLoader;
import com.google.common.cache.LoadingCache;
import org.wso2.apk.enforcer.commons.dto.JWTValidationInfo;
import org.wso2.apk.enforcer.config.ConfigHolder;
import org.wso2.apk.enforcer.config.dto.CacheDto;
import org.wso2.apk.enforcer.security.jwt.SignedJWTInfo;
import org.wso2.apk.enforcer.security.jwt.validator.JWTConstants;

import java.util.concurrent.TimeUnit;

/**
 * Class for initiating and returning caches.
 */
public class CacheProvider {
    private LoadingCache<String, SignedJWTInfo> gatewaySignedJWTParseCache;
    private LoadingCache<String, String> gatewayTokenCache;
    private LoadingCache<String, JWTValidationInfo> gatewayKeyCache;
    private LoadingCache<String, Boolean> invalidTokenCache;
    private LoadingCache<String, JWTValidationInfo> gatewayJWTTokenCache;
    private LoadingCache<String, String> getGatewayInternalKeyCache;
    private LoadingCache<String, String> getInvalidGatewayInternalKeyCache;
    private LoadingCache<String, JWTValidationInfo> getGatewayInternalKeyDataCache;
    private LoadingCache<String, String> getGatewayAPIKeyCache;
    private LoadingCache<String, String> getInvalidGatewayAPIKeyCache;
    private LoadingCache<String, JWTValidationInfo> getGatewayAPIKeyDataCache;

    private static boolean cacheEnabled = true;
    public void init() {
        CacheDto cacheDto = ConfigHolder.getInstance().getConfig().getCacheDto();
        cacheEnabled = cacheDto.isEnabled();
        int maxSize = cacheDto.getMaximumSize();
        int expiryTime = cacheDto.getExpiryTime();
        gatewaySignedJWTParseCache = initCache(maxSize, expiryTime);
        gatewayTokenCache = initCache(maxSize, expiryTime);
        gatewayKeyCache = initCache(maxSize, expiryTime);
        invalidTokenCache = initCache(maxSize, expiryTime);
        gatewayJWTTokenCache = initCache(maxSize, expiryTime);
        getGatewayInternalKeyCache = initCache(maxSize, expiryTime);
        getGatewayInternalKeyDataCache = initCache(maxSize, expiryTime);
        getInvalidGatewayInternalKeyCache = initCache(maxSize, expiryTime);
        getGatewayAPIKeyCache = initCache(maxSize, expiryTime);
        getInvalidGatewayAPIKeyCache = initCache(maxSize, expiryTime);
        getGatewayAPIKeyDataCache = initCache(maxSize, expiryTime);

    }

    private static LoadingCache initCache(int maxSize, int expiryTime) {
        return CacheBuilder.newBuilder()
                .maximumSize(maxSize)                                  // maximum 10000 tokens can be cached
                .expireAfterAccess(expiryTime, TimeUnit.MINUTES)      // cache will expire after 15 minutes of access
                .build(new CacheLoader<String, String>() {            // build the cacheloader
                    @Override public String load(String s) throws Exception {
                        return JWTConstants.UNAVAILABLE;
                    }

                });
    }


    /**
     * @return Gateway Internal Key cache
     */
    public LoadingCache getGatewayInternalKeyCache() {
        return getGatewayInternalKeyCache;
    }

    /**
     * @return Gateway Internal Key data cache
     */
    public LoadingCache getGatewayInternalKeyDataCache() {
        return getGatewayInternalKeyDataCache;
    }

    /**
     * @return Gateway Internal Key invalid data cache
     */
    public LoadingCache getInvalidGatewayInternalKeyCache() {
        return getInvalidGatewayInternalKeyCache;
    }

    /**
     *
     * @return SignedJWT ParsedCache
     */
    public LoadingCache getGatewaySignedJWTParseCache() {
        return gatewaySignedJWTParseCache;
    }

    /**
     * @return gateway token cache
     */
    public LoadingCache getGatewayTokenCache() {
        return gatewayTokenCache;
    }

    /**
     * @return gateway key cache
     */
    public LoadingCache getGatewayKeyCache() {
        return gatewayKeyCache;
    }

    /**
     * @return gateway invalid token cache
     */
    public LoadingCache getInvalidTokenCache() {
        return invalidTokenCache;
    }

    /**
     * @return JWT token cache
     */
    public LoadingCache getGatewayJWTTokenCache() {
        return gatewayJWTTokenCache;
    }

    /**
     * @return Gateway API key cache
     */
    public LoadingCache getGatewayAPIKeyCache() {
        return getGatewayAPIKeyCache;
    }

    /**
     * @return Gateway API key data cache
     */
    public LoadingCache getGatewayAPIKeyDataCache() {
        return getGatewayAPIKeyDataCache;
    }

    /**
     * @return Gateway API key invalid data cache
     */
    public LoadingCache getInvalidGatewayAPIKeyCache() {
        return getInvalidGatewayAPIKeyCache;
    }
}
