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
package org.wso2.apk.apimgt.api;

import org.wso2.apk.apimgt.api.dto.KeyManagerConfigurationDTO;
import org.wso2.apk.apimgt.api.model.*;
import org.wso2.apk.apimgt.api.model.policy.*;

import java.io.InputStream;
import java.util.List;
import java.util.Map;

/**
 * APIAdmin responsible for providing helper functionality
 */
public interface APIAdmin {
    /**
     * Returns environments of a given tenant
     *
     * @param tenantDomain tenant domain
     * @return List of environments related to the given tenant
     */
    List<Environment> getAllEnvironments(String tenantDomain) throws APIManagementException;

    /**
     * Returns environment of a given uuid
     *
     * @param tenantDomain tenant domain
     * @return List of environments related to the given tenant
     */
    Environment getEnvironment(String tenantDomain, String uuid) throws APIManagementException;

    /**
     * Creates a new environment in the tenant
     *
     * @param tenantDomain tenant domain
     * @param environment  content to add
     * @return Added environment
     * @throws APIManagementException if failed add environment
     */
    Environment addEnvironment(String tenantDomain, Environment environment) throws APIManagementException;

    /**
     * Delete an existing environment
     *
     * @param tenantDomain tenant domain
     * @param uuid         Environment identifier
     * @throws APIManagementException If failed to delete environment
     */
    void deleteEnvironment(String tenantDomain, String uuid) throws APIManagementException;

    /**
     * Updates the details of the given Environment.
     *
     * @param tenantDomain tenant domain
     * @param environment  content to update
     * @return updated environment
     * @throws APIManagementException if failed to update environment
     */
    Environment updateEnvironment(String tenantDomain, Environment environment) throws APIManagementException;


    Application[] getAllApplicationsOfTenantForMigration(String appTenantDomain) throws APIManagementException;

    /**
     * Returns List of Applications
     *
     * @param user            Logged-in user
     * @param owner           Owner of the application
     * @param organization    Organization
     * @param limit           The limit
     * @param offset          The offset
     * @param applicationName The application name
     * @param sortBy          The sortBy column
     * @param sortOrder       The sort order
     * @return List of applications match to the search conditions
     * @throws APIManagementException
     */
    Application[] getApplicationsWithPagination(String user, String owner, String organization, int limit, int offset,
                                                String applicationName, String sortBy, String sortOrder)
            throws APIManagementException;

    /**
     * Get count of the applications for the tenantId.
     *
     * @param organization      Organization
     * @param searchOwner       content to search applications based on owners
     * @param searchApplication content to search applications based on application
     * @throws APIManagementException if failed to get application
     * @return
     */

    public int getApplicationsCount(String organization, String searchOwner, String searchApplication)
            throws APIManagementException;

    /**
     * This methods loads the monetization implementation class
     *
     * @return monetization implementation class
     * @throws APIManagementException if failed to load monetization implementation class
     */
    Monetization getMonetizationImplClass() throws APIManagementException;

    /**
     * Get the info about the monetization usage publish job
     *
     * @throws APIManagementException if failed to get monetization usage publish info
     */
    MonetizationUsagePublishInfo getMonetizationUsagePublishInfo() throws APIManagementException;

    /**
     * Add the info about the monetization usage publish job
     *
     * @throws APIManagementException if failed to update monetization usage publish info
     */
    void addMonetizationUsagePublishInfo(MonetizationUsagePublishInfo monetizationUsagePublishInfo)
            throws APIManagementException;

    /**
     * Update the info about the monetization usage publish job
     *
     * @throws APIManagementException if failed to update monetization usage publish info
     */
    void updateMonetizationUsagePublishInfo(MonetizationUsagePublishInfo monetizationUsagePublishInfo)
            throws APIManagementException;

    /**
     * Adds a new category for the tenant
     *
     * @param userName     logged in user name
     * @param category     category to add
     * @param organization organization
     * @throws APIManagementException if failed add category
     */
    APICategory addCategory(APICategory category, String userName, String organization) throws APIManagementException;

    /**
     * Updates an API Category
     *
     * @param apiCategory
     * @throws APIManagementException
     */
    void updateCategory(APICategory apiCategory) throws APIManagementException;

    /**
     * Delete an API Category
     *
     * @param categoryID
     * @param username
     * @throws APIManagementException
     */
    void deleteCategory(String categoryID, String username) throws APIManagementException;

    /**
     * Checks whether an api category exists by the given name
     * <p>
     * 1. in case uuid is null : checks whether the categoryName is already taken in the tenantDomain (this
     * flow is used when adding a new api category)
     * 2. in case uuid is not null: checks whether the categoryName is already taken by any category other than the one
     * defined by the passed uuid in the given tenant
     *
     * @param categoryName
     * @param organization
     * @return true if an api category exists by the given category name
     * @throws APIManagementException
     */
    boolean isCategoryNameExists(String categoryName, String uuid, String organization) throws APIManagementException;

    /**
     * Returns all api categories of the organization
     *
     * @param organization Organization
     * @return
     * @throws APIManagementException
     */
    List<APICategory> getAllAPICategoriesOfOrganization(String organization) throws APIManagementException;

    /**
     * Returns all api categories of the organization with number of APIs for each category
     *
     * @param organization organization of the API
     * @return
     * @throws APIManagementException
     */
    List<APICategory> getAPICategoriesOfOrganization(String organization) throws APIManagementException;

    /**
     * Get API Category identified by the given uuid
     *
     * @param apiCategoryId api category UUID
     * @return
     * @throws APIManagementException
     */
    APICategory getAPICategoryByID(String apiCategoryId) throws APIManagementException;

    /**
     * The method converts the date into timestamp
     *
     * @param date
     * @return Timestamp in long format
     */
    long getTimestamp(String date);

    /**
     * This method used to retrieve key manager configurations for tenant
     *
     * @param organization organization of the key manager
     * @return KeyManagerConfigurationDTO list
     * @throws APIManagementException if error occurred
     */
    List<KeyManagerConfigurationDTO> getKeyManagerConfigurationsByOrganization(String organization) throws APIManagementException;

    /**
     * This method returns all the key managers registered in all the tenants
     *
     * @return
     * @throws APIManagementException
     */
    Map<String, List<KeyManagerConfigurationDTO>> getAllKeyManagerConfigurations() throws APIManagementException;

    /**
     * This method used to retrieve key manager with Id in respective tenant
     *
     * @param organization organization requested
     * @param id           uuid of key manager
     * @return KeyManagerConfigurationDTO for retrieved data
     * @throws APIManagementException
     */
    KeyManagerConfigurationDTO getKeyManagerConfigurationById(String organization, String id)
            throws APIManagementException;

    /**
     * This method is used to check IDP is in the given organization
     *
     * @param organization organization uuid
     * @param resourceId   IDP resource ID
     * @return boolean indication of it's existence
     * @throws APIManagementException
     */
    boolean isIDPExistInOrg(String organization, String resourceId) throws APIManagementException;

    /**
     * Used to get organization UUID of a application by giving consumer key.
     *
     * @param consumerKey consumer key of the application
     * @return ApplicationInfo details of a application
     * @throws APIManagementException
     */
    ApplicationInfo getLightweightApplicationByConsumerKey(String consumerKey) throws APIManagementException;

    /**
     * This method used to check existence of key manager with Id in respective tenant
     *
     * @param organization organization requested
     * @param id           uuid of key manager
     * @return existence
     * @throws APIManagementException
     */
    boolean isKeyManagerConfigurationExistById(String organization, String id) throws APIManagementException;

    /**
     * This method used to create key Manager
     *
     * @param keyManagerConfigurationDTO key manager data
     * @return created key manager
     * @throws APIManagementException
     */
    KeyManagerConfigurationDTO addKeyManagerConfiguration(KeyManagerConfigurationDTO keyManagerConfigurationDTO)
            throws APIManagementException;

    /**
     * This method used to update key Manager
     *
     * @param keyManagerConfigurationDTO key manager data
     * @return updated key manager
     * @throws APIManagementException
     */
    KeyManagerConfigurationDTO updateKeyManagerConfiguration(KeyManagerConfigurationDTO keyManagerConfigurationDTO)
            throws APIManagementException;

    /**
     * hTis method used to delete IDP mapped with key manager
     *
     * @param organization               organization requested
     * @param keyManagerConfigurationDTO key manager data
     * @throws APIManagementException
     */
    void deleteIdentityProvider(String organization, KeyManagerConfigurationDTO keyManagerConfigurationDTO)
            throws APIManagementException;

    /**
     * This method used to delete key manager
     *
     * @param organization               organization requested
     * @param keyManagerConfigurationDTO key manager data
     * @throws APIManagementException
     */
    void deleteKeyManagerConfigurationById(String organization, KeyManagerConfigurationDTO keyManagerConfigurationDTO)
            throws APIManagementException;

    /**
     * This method used to retrieve key manager from name
     *
     * @param organization organization requested
     * @param name         name requested
     * @return keyManager data
     * @throws APIManagementException
     */
    KeyManagerConfigurationDTO getKeyManagerConfigurationByName(String organization, String name)
            throws APIManagementException;

    /**
     * This method used to check the existence of the scope name for the particular user
     *
     * @param username  user to be validated
     * @param scopeName scope name to be checked
     * @return true if a scope exists by the given username
     * @throws APIManagementException
     */
    boolean isScopeExistsForUser(String username, String scopeName)
            throws APIManagementException;

    /**
     * This method used to check the existence of the scope name
     *
     * @param username  logged in username to get the tenantDomain
     * @param scopeName scope name to be checked
     * @return true if a scope exists
     * @throws APIManagementException
     */
    boolean isScopeExists(String username, String scopeName)
            throws APIManagementException;

    /**
     * Adds a tenant theme to the database
     *
     * @param organization     tenant ID of user
     * @param themeContent content of the tenant theme
     * @throws APIManagementException if an error occurs when adding a tenant theme to the database
     */
    void addTenantTheme(String organization, InputStream themeContent) throws APIManagementException;

    /**
     * Updates an existing tenant theme in the database
     *
     * @param organization     tenant ID of user
     * @param themeContent content of the tenant theme
     * @throws APIManagementException if an error occurs when updating an existing tenant theme in the database
     */
    void updateTenantTheme(String organization, InputStream themeContent) throws APIManagementException;

    /**
     * Retrieves a tenant theme from the database
     *
     * @param organization tenant ID of user
     * @return content of the tenant theme
     * @throws APIManagementException if an error occurs when retrieving a tenant theme from the database
     */
    InputStream getTenantTheme(String organization) throws APIManagementException;

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

    String getTenantConfig(String organization) throws APIManagementException;

    void updateTenantConfig(String organization, String config) throws APIManagementException;

    String getTenantConfigSchema(String organization);

    /**
     * Get policy list for given level
     *
     * @param organization organization of user
     * @param level        policy level
     * @return
     * @throws APIManagementException
     */
    Policy[] getPolicies(String organization, String level) throws APIManagementException;

    Policy getPolicyByNameAndType(String organization, String level, String name) throws APIManagementException;

    /**
     * Get api throttling policy by name
     *
     * @param username   name of the user
     * @param policyName name of the policy
     * @throws APIManagementException
     */
    APIPolicy getAPIPolicy(String username, String policyName) throws APIManagementException;

    /**
     * Get application throttling policy by name
     *
     * @param username   name of the user
     * @param policyName name of the policy
     * @throws APIManagementException
     */
    ApplicationPolicy getApplicationPolicy(String username, String policyName) throws APIManagementException;

    /**
     * Get subscription throttling policy by name
     *
     * @param username   name of the user
     * @param policyName name of the policy
     * @throws APIManagementException
     */
    SubscriptionPolicy getSubscriptionPolicy(String username, String policyName) throws APIManagementException;

    /**
     * @return List of block Conditions
     * @throws APIManagementException
     */
    List<BlockConditionsDTO> getBlockConditions(String organization) throws APIManagementException;

    /**
     * Get api throttling policy by uuid
     *
     * @param uuid UUID of the policy
     * @throws APIManagementException
     */
    APIPolicy getAPIPolicyByUUID(String uuid) throws APIManagementException;

    /**
     * Get application throttling policy by uuid
     *
     * @param uuid UUID of the policy
     * @throws APIManagementException
     */
    ApplicationPolicy getApplicationPolicyByUUID(String uuid) throws APIManagementException;

    /**
     * Get subscription throttling policy by uuid
     *
     * @param uuid UUID of the policy
     * @throws APIManagementException
     */
    SubscriptionPolicy getSubscriptionPolicyByUUID(String uuid) throws APIManagementException;

    /**
     * Retrieves a block condition by its UUID
     *
     * @param uuid uuid of the block condition
     * @return Retrieve a block Condition
     * @throws APIManagementException
     */
    BlockConditionsDTO getBlockConditionByUUID(String uuid) throws APIManagementException;

    /**
     * Add a block condition with condition status
     *
     * @param conditionType   type of the condition (IP, Context .. )
     * @param conditionValue  value of the condition
     * @param conditionStatus status of the condition
     * @param organization    organization
     * @return UUID of the new Block Condition
     * @throws APIManagementException
     */
    String addBlockCondition(String conditionType, String conditionValue, boolean conditionStatus, String organization)
            throws APIManagementException;

    void addPolicy(Policy policy, String username) throws APIManagementException;

    /**
     * Updates a block condition given its UUID
     *
     * @param uuid  uuid of the block condition
     * @param state state of condition
     * @return state change success or not
     * @throws APIManagementException
     */
    boolean updateBlockConditionByUUID(String uuid, String state) throws APIManagementException;

    /**
     * Updates throttle policy in global CEP, gateway and database.
     * <p>
     * Database transactions and deployements are not rolledback on failiure.
     * A flag will be inserted into the database whether the operation was
     * successfull or not.
     * </p>
     *
     * @param policy updated {@link Policy} object
     * @throws APIManagementException
     */
    void updatePolicy(Policy policy) throws APIManagementException;

    /**
     * Deletes a block condition given its UUID
     *
     * @param uuid uuid of the block condition
     * @return true if successfully deleted
     * @throws APIManagementException
     */
    boolean deleteBlockConditionByUUID(String uuid) throws APIManagementException;

    /**
     * Delete throttling policy
     *
     * @param username    username
     * @param policyLevel policy type
     * @param policyName  policy name
     * @throws APIManagementException
     */
    void deletePolicy(String username, String policyLevel, String policyName) throws APIManagementException;

    boolean hasAttachments(String username, String policyName, String policyLevel, String organization)
            throws APIManagementException;

    /**
     * Get the given Subscription Throttle Policy Permission
     *
     * @return Subscription Throttle Policy
     * @throws APIManagementException If failed to retrieve Subscription Throttle Policy Permission
     */
    Object getThrottleTierPermission(String tierName, String organization) throws APIManagementException;

    /**
     * Update Throttle Tier Permissions
     *
     * @param tierName       Tier Name
     * @param permissionType Permission Type
     * @param roles          Roles
     * @throws APIManagementException If failed to update subscription status
     */
    void updateThrottleTierPermissions(String tierName, String permissionType, String roles, String organization)
            throws APIManagementException;

    /**
     * Delete the Tier Permissions
     *
     * @param tierName     Tier Name
     * @param organization Organization
     * @throws APIManagementException
     */
    void deleteTierPermissions(String tierName, String organization) throws APIManagementException;

}
