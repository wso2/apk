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
package org.wso2.apk.enforcer.api;

import java.util.List;

import org.wso2.apk.enforcer.commons.model.EndpointSecurity;

/**
 * APIProcessUtils is used to convert the Endpoint Security DTO used in proto files into Enforcer specific
 * Object.
 */
public class APIProcessUtils {
    public static EndpointSecurity[] convertProtoEndpointSecurity
        (List<org.wso2.apk.enforcer.discovery.api.SecurityInfo> protoSecurityInfo) {
        EndpointSecurity[] securityInfo = new EndpointSecurity[protoSecurityInfo.size()];
        for (int i = 0; i < protoSecurityInfo.size(); i++) {
            EndpointSecurity security = new EndpointSecurity();
            security.setSecurityType(protoSecurityInfo.get(i).getSecurityType());
            security.setUsername(protoSecurityInfo.get(i).getUsername());
            security.setPassword(protoSecurityInfo.get(i).getPassword().toCharArray());
            security.setCustomParameters(protoSecurityInfo.get(i).getCustomParametersMap());
            security.setEnabled(protoSecurityInfo.get(i).getEnabled());
            securityInfo[i] = security;
        }
        return securityInfo;
    }
}
