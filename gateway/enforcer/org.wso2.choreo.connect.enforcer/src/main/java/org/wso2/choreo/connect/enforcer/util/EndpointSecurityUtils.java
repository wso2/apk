/*
 * Copyright (c) 2021, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 *
 * WSO2 Inc. licenses this file to you under the Apache License,
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

package org.wso2.choreo.connect.enforcer.util;

import org.wso2.choreo.connect.enforcer.commons.model.EndpointSecurity;
import org.wso2.choreo.connect.enforcer.commons.model.RequestContext;
import org.wso2.choreo.connect.enforcer.commons.model.ResourceConfig;
import org.wso2.choreo.connect.enforcer.constants.APIConstants;

import java.util.Base64;

/**
 * Util methods related to backend endpoint security.
 */
public class EndpointSecurityUtils {

    /**
     * Adds the backend endpoint security header to the given requestContext.
     *
     * @param requestContext requestContext instance to add the backend endpoint security header
     */
    public static void addEndpointSecurity(RequestContext requestContext) {
        EndpointSecurity[] endpointSecurities = null;
        if (requestContext.getMatchedResourcePaths() != null ) {
            // getting only first element as there could be only one resourcepaths for APIs except for graphQL APIs. 
            // For GQL APIs too would only have one endpoint for all resources.
            ResourceConfig resourceConfig = requestContext.getMatchedResourcePaths().get(0);
            if (resourceConfig.getEndpointSecurity() != null) {
                endpointSecurities = resourceConfig.getEndpointSecurity();
            }
        }
        for (EndpointSecurity securityInfo : endpointSecurities) {
            if (securityInfo != null && securityInfo.isEnabled() && 
            APIConstants.AUTHORIZATION_HEADER_BASIC.
                    equalsIgnoreCase(securityInfo.getSecurityType())) {
                requestContext.getRemoveHeaders().remove(APIConstants.AUTHORIZATION_HEADER_DEFAULT
                .toLowerCase());
                requestContext.addOrModifyHeaders(APIConstants.AUTHORIZATION_HEADER_DEFAULT,
                APIConstants.AUTHORIZATION_HEADER_BASIC + ' ' +
                Base64.getEncoder().encodeToString((securityInfo.getUsername() +
                ':' + String.valueOf(securityInfo.getPassword())).getBytes()));
            }
        }
    }
}
