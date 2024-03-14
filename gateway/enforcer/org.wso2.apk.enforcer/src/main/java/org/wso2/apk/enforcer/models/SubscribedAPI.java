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

import org.wso2.apk.enforcer.subscription.SubscribedAPIDto;

import java.util.ArrayList;
import java.util.List;
import java.util.regex.Pattern;

/**
 * Model class for Subscription related API details
 */
public class SubscribedAPI {

    private String name;
    private String version;
    private Pattern versionRegexPattern;

    public SubscribedAPI(org.wso2.apk.enforcer.discovery.subscription.SubscribedAPI subscribedApi) {

        this.name = subscribedApi.getName();
        this.version = subscribedApi.getVersion();
        this.versionRegexPattern = Pattern.compile(this.version);
    }

    public SubscribedAPI(SubscribedAPIDto subscribedApi) {

        this.name = subscribedApi.getName();
        this.version = subscribedApi.getVersion();
        this.versionRegexPattern = Pattern.compile(this.version);

    }

    public String getName() {

        return name;
    }

    public void setName(String name) {

        this.name = name;
    }

    public String getVersion() {

        return version;
    }

    public void setVersion(String version) {

        this.version = version;
    }

    public Pattern getVersionRegexPattern() {

        return versionRegexPattern;
    }

    public void setVersionRegexPattern(Pattern versionRegexPattern) {

        this.versionRegexPattern = versionRegexPattern;
    }
}
