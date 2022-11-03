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

package org.wso2.apk.apimgt.impl.dao;

import org.wso2.apk.apimgt.api.APIManagementException;
import org.wso2.apk.apimgt.api.model.APICategory;
import org.wso2.apk.apimgt.api.model.MonetizationUsagePublishInfo;
import org.wso2.apk.apimgt.api.model.botDataAPI.BotDetectionData;

import java.io.InputStream;
import java.util.List;

public interface AdminDAO {

    /**
     * Derives info about monetization usage publish job
     *
     * @return info about the monetization usage publish job
     * @throws APIManagementException
     */
    MonetizationUsagePublishInfo getMonetizationUsagePublishInfo() throws APIManagementException;

    /**
     * Updates info about monetization usage publish job
     *
     * @throws APIManagementException
     */
    void updateUsagePublishInfo(MonetizationUsagePublishInfo monetizationUsagePublishInfo)
            throws APIManagementException;

    /**
     * Add info about monetization usage publish job
     *
     * @throws APIManagementException
     */
    void addMonetizationUsagePublishInfo(MonetizationUsagePublishInfo monetizationUsagePublishInfo)
            throws APIManagementException;

    /**
     * Adds an API category
     *
     * @param category     Category
     * @param organization Organization
     * @return Category
     */
    APICategory addCategory(APICategory category, String organization) throws APIManagementException;

    /**
     * Update API Category
     *
     * @param apiCategory API category object with updated details
     * @throws APIManagementException
     */
    void updateCategory(APICategory apiCategory) throws APIManagementException;

    /**
     * Delete API Category
     *
     * @param categoryID API category ID
     * @throws APIManagementException
     */
    void deleteCategory(String categoryID) throws APIManagementException;

    /**
     * Get all available API categories of the organization
     *
     * @param organization
     * @return
     * @throws APIManagementException
     */
    List<APICategory> getAllCategories(String organization) throws APIManagementException;

    /**
     * Checks whether the given category name is already available under given tenant domain with any UUID other than
     * the given UUID
     *
     * @param categoryName
     * @param uuid
     * @param organization
     * @return
     */
    boolean isAPICategoryNameExists(String categoryName, String uuid, String organization) throws APIManagementException;

    /**
     * Get API category by ID
     *
     * @param apiCategoryID Category ID
     * @return
     * @throws APIManagementException
     */
    APICategory getAPICategoryByID(String apiCategoryID) throws APIManagementException;

    /**
     * Adds a tenant theme to the database
     *
     * @param organization tenant ID of user
     * @param themeContent content of the tenant theme
     * @throws APIManagementException if an error occurs when adding a tenant theme to the database
     */
    void addTenantTheme(String organization, InputStream themeContent) throws APIManagementException;

    /**
     * Updates an existing tenant theme in the database
     *
     * @param organization tenant ID of user
     * @param themeContent content of the tenant theme
     * @throws APIManagementException if an error occurs when updating an existing tenant theme in the database
     */
    void updateTenantTheme(String organization, InputStream themeContent) throws APIManagementException;

    /**
     * Retrieves a tenant theme from the database
     *
     * @param tenantId tenant ID of user
     * @return content of the tenant theme
     * @throws APIManagementException if an error occurs when retrieving a tenant theme from the database
     */
    InputStream getTenantTheme(String tenantId) throws APIManagementException;

    /**
     * Checks whether a tenant theme exist for a particular tenant
     *
     * @param organization tenant ID of user
     * @return true if a tenant theme exist for a particular tenant ID, false otherwise
     * @throws APIManagementException if an error occurs when determining whether a tenant theme exists for a given
     *                                tenant ID
     */
    boolean isTenantThemeExist(String organization) throws APIManagementException;

    /**
     * Deletes a tenant theme from the database
     *
     * @param organization tenant ID of user
     * @throws APIManagementException if an error occurs when deleting a tenant theme from the database
     */
    void deleteTenantTheme(String organization) throws APIManagementException;

}
