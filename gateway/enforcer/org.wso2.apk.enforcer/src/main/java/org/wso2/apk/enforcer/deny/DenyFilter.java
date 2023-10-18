/*
 * Copyright (c) 2023, WSO2 LLC. (http://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.wso2.apk.enforcer.deny;

import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.wso2.apk.enforcer.commons.Filter;
import org.wso2.apk.enforcer.commons.model.APIConfig;
import org.wso2.apk.enforcer.commons.model.RequestContext;
import org.wso2.apk.enforcer.constants.APIConstants;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.Map;

/**
 * enum for storing types of deny policies
 */
enum DenyPolicyType {
    APPLICATION, SUBSCRIPTION, USER
}

public class DenyFilter implements Filter {

    private static final Logger logger = LogManager.getLogger(DenyFilter.class);

    // Hashmap to keep track of the blocked subs, apps and users.
    private final HashMap<DenyPolicyType, ArrayList<String>> deniedDetailsMap = new HashMap<>();

    @Override
    public void init(APIConfig apiConfig, Map<String, String> configProperties) {
        Filter.super.init(apiConfig, configProperties);
        loadDeniedList();
    }

    /**
     *
     * @param requestContext {@code RequestContext} object
     * @return boolean
     */
    @Override
    public boolean handleRequest(RequestContext requestContext) {
        String username = requestContext.getAuthenticationContext().getUsername();
        String applicationId = requestContext.getAuthenticationContext().getApplicationUUID();
        String subscriptionId = requestContext.getAuthenticationContext().getSubscriber();

        if (isInDeniedList(username, DenyPolicyType.USER) || isInDeniedList(applicationId, DenyPolicyType.APPLICATION) ||
                isInDeniedList(subscriptionId, DenyPolicyType.SUBSCRIPTION)) {
            logger.debug("Request blocked due to deny policy (enforcer).");
            requestContext.getProperties()
                    .put(APIConstants.MessageFormat.STATUS_CODE, APIConstants.StatusCodes.UNAUTHORIZED.getCode());
            requestContext.getProperties().put(APIConstants.MessageFormat.ERROR_CODE, APIConstants.StatusCodes.UNAUTHORIZED.getValue());
            requestContext.getProperties().put(APIConstants.MessageFormat.ERROR_MESSAGE, APIConstants.REQUEST_DENIED_MESSAGE);
            requestContext.getProperties().put(APIConstants.MessageFormat.ERROR_DESCRIPTION, APIConstants.REQUEST_DENIED_DESCRIPTION);
            return false;
        }
        return true;
    }

    private void loadDeniedList() {
        // TODO (Gayangi): implement loading data from the database once database is implemented
        deniedDetailsMap.put(DenyPolicyType.USER, new ArrayList<>());
        deniedDetailsMap.put(DenyPolicyType.APPLICATION, new ArrayList<>());
        deniedDetailsMap.put(DenyPolicyType.SUBSCRIPTION, new ArrayList<>());
    }

    /**
     *
     * @param value Represents the value to be checked to see if it has been blocked
     * @param denyPolicyType
     * @return true if value is in the relevant denied list
     */
    private boolean isInDeniedList(String value, DenyPolicyType denyPolicyType) {
        return deniedDetailsMap.get(denyPolicyType).contains(value);
    }
}
