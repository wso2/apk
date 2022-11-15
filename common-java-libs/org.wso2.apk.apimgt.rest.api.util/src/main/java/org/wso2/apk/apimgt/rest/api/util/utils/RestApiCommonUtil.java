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
import org.wso2.apk.apimgt.api.APIConsumer;
import org.wso2.apk.apimgt.api.APIManagementException;
import org.wso2.apk.apimgt.api.APIProvider;
import org.wso2.apk.apimgt.api.model.API;
import org.wso2.apk.apimgt.impl.APIManagerFactory;
import org.wso2.apk.apimgt.rest.api.util.RestApiConstants;
import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.joda.time.DateTime;
import org.joda.time.format.DateTimeFormatter;
import org.joda.time.format.ISODateTimeFormat;
import org.wso2.apk.apimgt.api.APIManagementException;
import org.wso2.apk.apimgt.api.APIProvider;
import org.wso2.apk.apimgt.api.ExceptionCodes;
import org.wso2.apk.apimgt.api.model.APIInfo;
import org.wso2.apk.apimgt.impl.APIManagerFactory;

import java.util.Date;
import java.util.HashMap;
import java.util.Map;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class RestApiCommonUtil {

    private static final Log log = LogFactory.getLog(RestApiCommonUtil.class);

    public static APIProvider getLoggedInUserProvider() throws APIManagementException {

        return APIManagerFactory.getInstance().getAPIProvider(getLoggedInUsername());
    }
    public static String getLoggedInUsername() {

//        return UserContext.getThreadLocalUserContext().getUsername();
        return "apkuser";
    }

    public static String getLoggedInUserTenantDomain() {

//        return UserContext.getThreadLocalUserContext().getOrganization();
        return "carbon.super";
    }

    /**
     * Returns date in RFC3339 format.
     * Example: 2008-11-13T12:23:30-08:00
     *
     * @param date Date object
     * @return date string in RFC3339 format.
     */
    public static String getRFC3339Date(Date date) {

        DateTimeFormatter jodaDateTimeFormatter = ISODateTimeFormat.dateTime();
        DateTime dateTime = new DateTime(date);
        return jodaDateTimeFormatter.print(dateTime);
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
     * Returns the paginated url for Endpoint Certificate Usage API
     *
     * @param alias  alias of certificate
     * @param offset starting index
     * @param limit  max number of objects returned
     * @return constructed paginated url
     */
    public static String getCertificateUsagePaginatedURL(String alias, Integer offset, Integer limit) {
        String paginatedURL = RestApiConstants.ENDPOINT_CERTIFICATE_USAGE_GET_PAGINATION_URL;
        paginatedURL = paginatedURL.replace(RestApiConstants.ALIAS_PARAM, alias);
        paginatedURL = paginatedURL.replace(RestApiConstants.LIMIT_PARAM, String.valueOf(limit));
        paginatedURL = paginatedURL.replace(RestApiConstants.OFFSET_PARAM, String.valueOf(offset));
        return paginatedURL;
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

    /**
     * Returns an APIConsumer which is corresponding to the current logged in user taken from the carbon context
     *
     * @return an APIConsumer which is corresponding to the current logged in user
     * @throws APIManagementException
     */
    public static APIConsumer getLoggedInUserConsumer() throws APIManagementException {

        return APIManagerFactory.getInstance().getAPIConsumer(getLoggedInUsername());
    }

    public static APIConsumer getConsumer(String subscriberName) throws APIManagementException {

        return APIManagerFactory.getInstance().getAPIConsumer(subscriberName);
    }

    public static APIConsumer getConsumer(String subscriberName, String organization) throws APIManagementException {

        return APIManagerFactory.getInstance().getAPIConsumer(subscriberName, organization);
    }

    public static String getSubscriptionPaginatedURLForAPIId(Integer offset, Integer limit, String apiId,
                                                             String groupId) {

        groupId = groupId == null ? "" : groupId;
        String paginatedURL = RestApiConstants.SUBSCRIPTIONS_GET_PAGINATION_URL_APIID;
        paginatedURL = paginatedURL.replace(RestApiConstants.LIMIT_PARAM, String.valueOf(limit));
        paginatedURL = paginatedURL.replace(RestApiConstants.OFFSET_PARAM, String.valueOf(offset));
        paginatedURL = paginatedURL.replace(RestApiConstants.APIID_PARAM, apiId);
        paginatedURL = paginatedURL.replace(RestApiConstants.GROUPID_PARAM, groupId);
        return paginatedURL;
    }

    public static String getDocumentationPaginatedURL(Integer offset, Integer limit, String apiId) {

        String paginatedURL = RestApiConstants.DOCUMENTS_GET_PAGINATION_URL;
        paginatedURL = paginatedURL.replace(RestApiConstants.LIMIT_PARAM, String.valueOf(limit));
        paginatedURL = paginatedURL.replace(RestApiConstants.OFFSET_PARAM, String.valueOf(offset));
        paginatedURL = paginatedURL.replace(RestApiConstants.APIID_PARAM, apiId);
        return paginatedURL;
    }

    public static String getTagsPaginatedURL(Integer offset, Integer limit) {

        String paginatedURL = RestApiConstants.TAGS_GET_PAGINATION_URL;
        paginatedURL = paginatedURL.replace(RestApiConstants.LIMIT_PARAM, String.valueOf(limit));
        paginatedURL = paginatedURL.replace(RestApiConstants.OFFSET_PARAM, String.valueOf(offset));
        return paginatedURL;
    }

    public static String getAPIPaginatedURL(Integer offset, Integer limit, String query) {

        String paginatedURL = RestApiConstants.APIS_GET_PAGINATION_URL;
        paginatedURL = paginatedURL.replace(RestApiConstants.LIMIT_PARAM, String.valueOf(limit));
        paginatedURL = paginatedURL.replace(RestApiConstants.OFFSET_PARAM, String.valueOf(offset));
        paginatedURL = paginatedURL.replace(RestApiConstants.QUERY_PARAM, query);
        return paginatedURL;
    }

    /**
     * Returns the paginated url for API ratings
     *
     * @param offset starting index
     * @param limit  max number of objects returned
     * @return constructed paginated url
     */
    public static String getRatingPaginatedURL(Integer offset, Integer limit, String apiId) {

        String paginatedURL = RestApiConstants.RATINGS_GET_PAGINATION_URL;
        paginatedURL = paginatedURL.replace(RestApiConstants.LIMIT_PARAM, String.valueOf(limit));
        paginatedURL = paginatedURL.replace(RestApiConstants.OFFSET_PARAM, String.valueOf(offset));
        paginatedURL = paginatedURL.replace(RestApiConstants.APIID_PARAM, apiId);
        return paginatedURL;
    }

    /**
     * Returns the paginated url for tiers
     *
     * @param tierLevel tier level (api/application or resource)
     * @param offset    starting index
     * @param limit     max number of objects returned
     * @return constructed paginated url
     */
    public static String getTiersPaginatedURL(String tierLevel, Integer offset, Integer limit) {

        String paginatedURL = RestApiConstants.TIERS_GET_PAGINATION_URL;
        paginatedURL = paginatedURL.replace(RestApiConstants.TIER_LEVEL_PARAM, tierLevel);
        paginatedURL = paginatedURL.replace(RestApiConstants.LIMIT_PARAM, String.valueOf(limit));
        paginatedURL = paginatedURL.replace(RestApiConstants.OFFSET_PARAM, String.valueOf(offset));
        return paginatedURL;
    }


    /**
     * Returns the paginated url for tags
     *
     * @param offset starting index
     * @param limit  max number of objects returned
     * @return constructed paginated url
     */
    public static String getResourcePathPaginatedURL(Integer offset, Integer limit) {

        String paginatedURL = RestApiConstants.RESOURCE_PATH_PAGINATION_URL;
        paginatedURL = paginatedURL.replace(RestApiConstants.LIMIT_PARAM, String.valueOf(limit));
        paginatedURL = paginatedURL.replace(RestApiConstants.OFFSET_PARAM, String.valueOf(offset));
        return paginatedURL;
    }

    /**
     * Returns the paginated url for product documentations
     *
     * @param offset starting index
     * @param limit  max number of objects returned
     * @return constructed paginated url
     */
    public static String getProductDocumentPaginatedURL(Integer offset, Integer limit, String apiId) {

        String paginatedURL = RestApiConstants.PRODUCT_DOCUMENTS_GET_PAGINATION_URL;
        paginatedURL = paginatedURL.replace(RestApiConstants.LIMIT_PARAM, String.valueOf(limit));
        paginatedURL = paginatedURL.replace(RestApiConstants.OFFSET_PARAM, String.valueOf(offset));
        paginatedURL = paginatedURL.replace(RestApiConstants.APIID_PARAM, apiId);
        return paginatedURL;
    }

    public static APIProvider getProvider(String username) throws APIManagementException {

        return APIManagerFactory.getInstance().getAPIProvider(username);
    }

    /**
     * Url validator, Allow any url with https and http.
     * Allow any url without fully qualified domain
     *
     * @param url Url as string
     * @return boolean type stating validated or not
     */
    public static boolean isURL(String url) {

        Pattern pattern = Pattern.compile("^(http|https)://(.)+", Pattern.CASE_INSENSITIVE);
        Matcher matcher = pattern.matcher(url);
        return matcher.matches();

    }

    /**
     * @param apiId UUID of the API
     * @return API details
     * @throws APIManagementException when API does not exist in the DB
     */
    public static APIInfo validateAPIExistence(String apiId) throws APIManagementException {

        APIProvider apiProvider = getLoggedInUserProvider();
        APIInfo apiInfo = apiProvider.getAPIInfoByUUID(apiId);
        if (apiInfo == null) {
            throw new APIManagementException("Couldn't retrieve existing API with API UUID: "
                    + apiId, ExceptionCodes.from(ExceptionCodes.API_NOT_FOUND,
                    apiId));
        }
        return apiInfo;
    }

}
