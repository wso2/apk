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
package org.wso2.apk.enforcer.analytics.publisher.reporter.moesif.util;

import com.google.gson.annotations.SerializedName;

/**
 * POJO to parse the JSON.
 */
public class MoesifKeyEntry {
    private String uuid;
    @SerializedName("organization_id")
    private String organizationID;
    @SerializedName("moesif_key")
    private String moesifKey;
    private String env;

    public MoesifKeyEntry() {
    }

    public String getMoesif_key() {
        return moesifKey;
    }

    public String getOrganization_id() {
        return organizationID;
    }

    public String getUuid() {
        return uuid;
    }

    public String getEnv() {
        return env;
    }
}
