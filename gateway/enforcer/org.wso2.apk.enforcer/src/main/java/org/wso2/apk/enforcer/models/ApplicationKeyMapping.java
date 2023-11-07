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
 * Entity for keeping mapping between Application and Consumer key.
 */
public class ApplicationKeyMapping implements CacheableEntity<String> {

    private String applicationUUID;
    private String securityScheme;
    private String applicationIdentifier;
    private String keyType;
    private String envId;

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

    public String getEnvId() {
        return envId;
    }

    public void setEnvId(String envId) {
        this.envId = envId;
    }

    @Override
    public String getCacheKey() {
        return securityScheme + CacheableEntity.DELEM_PERIOD + applicationIdentifier;
    }

    @Override
    public String toString() {
        return "ApplicationKeyMapping{" +
                "applicationUUID='" + applicationUUID + '\'' +
                ", securityScheme='" + securityScheme + '\'' +
                ", applicationIdentifier='" + applicationIdentifier + '\'' +
                ", keyType='" + keyType + '\'' +
                ", envId='" + envId + '\'' +
                '}';
    }
}
