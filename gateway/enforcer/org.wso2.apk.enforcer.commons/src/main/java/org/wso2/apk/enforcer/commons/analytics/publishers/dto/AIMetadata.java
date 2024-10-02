/*
 *  Copyright (c) 2024, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */

package org.wso2.apk.enforcer.commons.analytics.publishers.dto;

import com.fasterxml.jackson.annotation.JsonProperty;

/**
 * AI metadata in analytics event.
 */
public class AIMetadata {
    @JsonProperty("model")
    private String model;

    @JsonProperty("vendor_name")
    private String vendorName;

    @JsonProperty("vendor_version")
    private String vendorVersion;

    public String getModel() {

        return model;
    }

    public void setModel(String model) {

        this.model = model;
    }

    public String getVendorName() {

        return vendorName;
    }

    public void setVendorName(String vendorName) {

        this.vendorName = vendorName;
    }

    public String getVendorVersion() {

        return vendorVersion;
    }

    public void setVendorVersion(String vendorVersion) {

        this.vendorVersion = vendorVersion;
    }
}
