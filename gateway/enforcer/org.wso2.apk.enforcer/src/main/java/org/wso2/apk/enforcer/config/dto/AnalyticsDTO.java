/*
 * Copyright (c) 2021, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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

package org.wso2.apk.enforcer.config.dto;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * Holds the analytics data publishing related Configuration.
 */
public class AnalyticsDTO {

    private boolean enabled = false;

    private Map<String,Object> properties = new HashMap<>();
    private List<AnalyticsPublisherConfigDTO> analyticsPublisherConfigDTOList = new ArrayList<>();
    private AnalyticsReceiverConfigDTO serverConfig;

    public AnalyticsReceiverConfigDTO getServerConfig() {

        return serverConfig;
    }

    public void setServerConfig(AnalyticsReceiverConfigDTO serverConfig) {

        this.serverConfig = serverConfig;
    }

    public List<AnalyticsPublisherConfigDTO> getAnalyticsPublisherConfigDTOList() {

        return analyticsPublisherConfigDTOList;
    }

    public void addAnalyticsPublisherConfig(AnalyticsPublisherConfigDTO analyticsPublisherConfigDTO) {

        this.analyticsPublisherConfigDTOList.add(analyticsPublisherConfigDTO);
    }

    public boolean isEnabled() {

        return enabled;
    }

    public void setEnabled(boolean enabled) {

        this.enabled = enabled;
    }

    public Map<String, Object> getProperties() {

        return properties;
    }

    public void setProperties(Map<String, Object> properties) {

        this.properties = properties;
    }
}
