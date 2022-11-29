/*
 * Copyright (c) 2022, WSO2 LLC. (http://www.wso2.com).
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

package org.wso2.apk.apimgt.rest.api.util.utils;

import com.fasterxml.jackson.annotation.JsonInclude;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.commons.lang3.exception.ExceptionUtils;
import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.json.simple.JSONObject;
import org.wso2.apk.apimgt.api.*;
import org.wso2.apk.apimgt.api.model.DuplicateAPIException;
import org.wso2.apk.apimgt.impl.utils.APIUtil;

import java.util.*;

public class RestApiUtil {

    public static final Log log = LogFactory.getLog(RestApiUtil.class);

    /**
     * This method is used to get the scope list from the yaml file
     *
     * @return MAP of scope list for all portal
     */
    public static  Map<String, List<String>> getScopesInfoFromAPIYamlDefinitions() throws APIManagementException {

        Map<String, List<String>>   portalScopeList = new HashMap<>();
        //TODO: APK
        return portalScopeList;
    }

    /**
     * Check if the specified throwable e is due to an authorization failure
     * @param e throwable to check
     * @return true if the specified throwable e is due to an authorization failure, false otherwise
     */
    @SuppressWarnings("ThrowableResultOfMethodCallIgnored")
    public static boolean isDueToAuthorizationFailure(Throwable e) {
        Throwable rootCause = getPossibleErrorCause(e);
        return rootCause instanceof APIMgtAuthorizationFailedException;
    }

    /**
     * Attempts to find the actual cause of the throwable 'e'
     *
     * @param e throwable
     * @return the root cause of 'e' if the root cause exists, otherwise returns 'e' itself
     */
    private static Throwable getPossibleErrorCause (Throwable e) {
        Throwable rootCause = ExceptionUtils.getRootCause(e);
        rootCause = rootCause == null ? e : rootCause;
        return rootCause;
    }

    /**
     * To check whether the DevPortal Anonymous Mode is enabled. It can be either enabled globally or tenant vice.
     *
     * @param tenantDomain Tenant domain
     * @return whether devportal anonymous mode is enabled or not
     */
    public static boolean isDevPortalAnonymousEnabled(String tenantDomain) {
        try {
            org.json.simple.JSONObject tenantConfig = APIUtil.getTenantConfig(tenantDomain);
            Object value = tenantConfig.get(APIConstants.API_TENANT_CONF_ENABLE_ANONYMOUS_MODE);
            if (value != null) {
                return Boolean.parseBoolean(value.toString());
            } else {
                return APIUtil.isDevPortalAnonymous();
            }
        } catch (APIManagementException e) {
            log.error("Error while retrieving Anonymous config from registry", e);
        }
        return true;
    }

    public static <T> String getJsonFromDTO(T dto) throws APIManagementException {
        ObjectMapper mapper = new ObjectMapper();
        mapper.setSerializationInclusion(JsonInclude.Include.NON_NULL);
        try {
            return mapper.writeValueAsString(dto);
        } catch (JsonProcessingException e) {
            throw new APIManagementException("Error");
        }
    }

    public static <T> T getDTOFromJson(String json, Class<T> clazz)
            throws APIManagementException{
        ObjectMapper mapper = new ObjectMapper();
        try {
            return mapper.readValue(json, clazz);
        } catch (JsonProcessingException e) {
            throw new APIManagementException("Error");
        }
    }

    /**
     * Check if the specified throwable e is happened as the updated/new resource conflicting with an already existing
     * resource
     *
     * @param e throwable to check
     * @return true if the specified throwable e is happened as the updated/new resource conflicting with an already
     *   existing resource, false otherwise
     */
    @SuppressWarnings("ThrowableResultOfMethodCallIgnored")
    public static boolean isDueToResourceAlreadyExists(Throwable e) {
        Throwable rootCause = getPossibleErrorCause(e);
        return rootCause instanceof APIMgtResourceAlreadyExistsException || rootCause instanceof DuplicateAPIException;
    }

    /**
     * Check if the message of the root cause message of 'e' matches with the specified message
     *
     * @param e throwable to check
     * @param message error message
     * @return true if the message of the root cause of 'e' matches with 'message'
     */
    @SuppressWarnings("ThrowableResultOfMethodCallIgnored")
    public static boolean rootCauseMessageMatches (Throwable e, String message) {
        Throwable rootCause = getPossibleErrorCause(e);
        return rootCause.getMessage().contains(message);
    }

    /**
     * Returns the current logged in consumer's group id
     * @return group id of the current logged in user.
     */
    @SuppressWarnings("unchecked")
    public static String getLoggedInUserGroupId() throws APIManagementException {
        String username = RestApiCommonUtil.getLoggedInUsername();
        String tenantDomain = RestApiCommonUtil.getLoggedInUserTenantDomain();
        JSONObject loginInfoJsonObj = new JSONObject();
        try {
            loginInfoJsonObj.put("user", username);
            loginInfoJsonObj.put("isSuperTenant", tenantDomain.equals("carbon.super"));
            String loginInfoString = loginInfoJsonObj.toJSONString();
            String[] groupIdArr = getGroupIds(loginInfoString);
            String groupId = "";
            if (groupIdArr != null) {
                for (int i = 0; i < groupIdArr.length; i++) {
                    if (groupIdArr[i] != null) {
                        if (i == groupIdArr.length - 1) {
                            groupId = groupId + groupIdArr[i];
                        } else {
                            groupId = groupId + groupIdArr[i] + ",";
                        }
                    }
                }
            }
            return groupId;
        } catch (APIManagementException e) {
            String errorMsg = "Unable to get groupIds of user " + username;
            throw  new APIManagementException(errorMsg, ExceptionCodes.INTERNAL_ERROR);
        }
    }

    private static String[] getGroupIds(String loginInfoString) throws APIManagementException {
        String groupingExtractorClass = APIUtil.getRESTApiGroupingExtractorImplementation();
        return APIUtil.getGroupIdsFromExtractor(loginInfoString, groupingExtractorClass);
    }
}
