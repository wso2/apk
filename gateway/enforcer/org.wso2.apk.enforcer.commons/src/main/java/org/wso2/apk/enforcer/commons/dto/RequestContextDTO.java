/*
 * Copyright (c) 2021 WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.wso2.apk.enforcer.commons.dto;

import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;

import java.util.Map;

/**
 * Representation of Request Information.
 */
public class RequestContextDTO {
    private static final Log log = LogFactory.getLog(RequestContextDTO.class);

    // request message information
    private MsgInfoDTO msgInfo;
    // invoked API request information
    private APIRequestInfoDTO apiRequestInfo;
    // custom property map used to populate customProperty key template value
    private Map<String, Object> customProperty;

    public MsgInfoDTO getMsgInfo() {

        return msgInfo;
    }

    public void setMsgInfo(MsgInfoDTO msgInfo) {

        this.msgInfo = msgInfo;
    }

    public APIRequestInfoDTO getApiRequestInfo() {

        return apiRequestInfo;
    }

    public void setApiRequestInfo(APIRequestInfoDTO apiRequestInfo) {

        this.apiRequestInfo = apiRequestInfo;
    }

    public Map<String, Object> getCustomProperty() {

        return customProperty;
    }

    public void setCustomProperty(Map<String, Object> customProperty) {

        this.customProperty = customProperty;
    }
}

