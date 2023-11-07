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

import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

/**
 * Entity for keeping Application related information. Represents an Application in APIM.
 */
public class Application implements CacheableEntity<String> {

    private String uuid;
    private String name = null;
    private String owner = null;
    private Map<String, String> attributes = new ConcurrentHashMap<>();

    public String getUUID() {

        return uuid;
    }

    public void setUUID(String uuid) {

        this.uuid = uuid;
    }

    public String getName() {

        return name;
    }

    public void setName(String name) {

        this.name = name;
    }

    public String getOwner() {

        return owner;
    }

    public void setOwner(String owner) {

        this.owner = owner;
    }

    public Map<String, String> getAttributes() {

        return attributes;
    }

    public void addAttribute(String key, String value) {

        this.attributes.put(key, value);
    }

    public void removeAttribute(String key) {

        this.attributes.remove(key);
    }

    public String getCacheKey() {

        return uuid;
    }

    @Override
    public String toString() {

        return "Application [uuid=" + uuid + ", name=" + name + ", owner=" + owner + ", attributes=" + attributes
                + "]";
    }

}

