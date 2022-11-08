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

import org.apache.commons.lang3.StringUtils;
import org.wso2.apk.apimgt.rest.api.util.RestApiConstants;

import java.util.HashMap;
import java.util.Map;

public class RestApiCommonUtil {
    public static String getLoggedInUsername() {

//        return UserContext.getThreadLocalUserContext().getUsername();
        return "apkuser";
    }

    public static String getLoggedInUserTenantDomain() {

//        return UserContext.getThreadLocalUserContext().getOrganization();
        return "carbon.super";
    }

    /**
     * Returns the next/previous offset/limit parameters properly when current offset, limit and size parameters are
     * specified
     *
     * @param offset current starting index
     * @param limit  current max records
     * @param size   maximum index possible
     * @return the next/previous offset/limit parameters as a hash-map
     */
    public static Map<String, Integer> getPaginationParams(Integer offset, Integer limit, Integer size) {

        Map<String, Integer> result = new HashMap<>();
        if (offset >= size || offset < 0)
            return result;

        int start = offset;
        int end = offset + limit - 1;

        int nextStart = end + 1;
        if (nextStart < size) {
            result.put(RestApiConstants.PAGINATION_NEXT_OFFSET, nextStart);
            result.put(RestApiConstants.PAGINATION_NEXT_LIMIT, limit);
        }

        int previousEnd = start - 1;
        int previousStart = previousEnd - limit + 1;

        if (previousEnd >= 0) {
            if (previousStart < 0) {
                result.put(RestApiConstants.PAGINATION_PREVIOUS_OFFSET, 0);
                result.put(RestApiConstants.PAGINATION_PREVIOUS_LIMIT, limit);
            } else {
                result.put(RestApiConstants.PAGINATION_PREVIOUS_OFFSET, previousStart);
                result.put(RestApiConstants.PAGINATION_PREVIOUS_LIMIT, limit);
            }
        }
        return result;
    }

    /**
     * Returns the paginated url for Applications API when it comes to sortOrder and sortBy
     *
     * @param offset    starting index
     * @param limit     max number of objects returned
     * @param groupId   group ID of the application
     * @param sortOrder specified sorting order ex: ASC
     * @param sortBy    specified parameter for the sort ex: name
     * @return constructed paginated url
     */
    public static String getApplicationPaginatedURLWithSortParams(Integer offset, Integer limit, String groupId,
                                                                  String sortOrder, String sortBy) {

        groupId = groupId == null ? "" : groupId;
        String paginatedURL = RestApiConstants.APPLICATIONS_GET_PAGINATION_URL;
        if (StringUtils.isNoneBlank(sortBy) || StringUtils.isNotBlank(sortOrder)) {
            sortOrder = sortOrder == null ? "" : sortOrder;
            sortBy = sortBy == null ? "" : sortBy;
            paginatedURL = RestApiConstants.APPLICATIONS_GET_PAGINATION_URL_WITH_SORTBY_SORTORDER;
            paginatedURL = paginatedURL.replace(RestApiConstants.SORTBY_PARAM, sortBy);
            paginatedURL = paginatedURL.replace(RestApiConstants.SORTORDER_PARAM, sortOrder);
        }
        paginatedURL = paginatedURL.replace(RestApiConstants.LIMIT_PARAM, String.valueOf(limit));
        paginatedURL = paginatedURL.replace(RestApiConstants.OFFSET_PARAM, String.valueOf(offset));
        paginatedURL = paginatedURL.replace(RestApiConstants.GROUPID_PARAM, groupId);
        return paginatedURL;
    }

    /**
     * Returns the paginated url for Applications API
     *
     * @param offset  starting index
     * @param limit   max number of objects returned
     * @param groupId groupId of the Application
     * @return constructed paginated url
     */
    public static String getApplicationPaginatedURL(Integer offset, Integer limit, String groupId) {

        return getApplicationPaginatedURLWithSortParams(offset, limit, groupId, null, null);
    }

    /**
     * Returns the paginated url for admin  /Applications API
     *
     * @param offset starting index
     * @param limit  max number of objects returned
     * @return constructed paginated url
     */
    public static String getApplicationPaginatedURL(Integer offset, Integer limit) {

        String paginatedURL = RestApiConstants.APPLICATIONS_GET_PAGINATION_URL;
        paginatedURL = paginatedURL.replace(RestApiConstants.LIMIT_PARAM, String.valueOf(limit));
        paginatedURL = paginatedURL.replace(RestApiConstants.OFFSET_PARAM, String.valueOf(offset));
        return paginatedURL;
    }
}
