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
import org.wso2.apk.apimgt.impl.dto.TierPermissionDTO;

import java.util.Set;

public interface TierDAO {

    /**
     * Update Tier Permissions.
     *
     * @param tierName       Tier Name
     * @param permissionType Permission Type
     * @param roles          Roles
     * @param organization   Organization
     * @throws APIManagementException if fails to update Tier permissions
     */
    void updateTierPermissions(String tierName, String permissionType, String roles, String organization)
            throws APIManagementException;

    /**
     * Delete Tier Permissions.
     *
     * @param tierName     Tier Name
     * @param organization Organization
     * @throws APIManagementException if fails to delete Tier permissions
     */
    void deleteThrottlingPermissions(String tierName, String organization) throws APIManagementException;

    /**
     * Retrieve Tier Permissions.
     *
     * @param organization Organization
     * @throws APIManagementException if fails to retrieve Tier permissions
     */
    Set<TierPermissionDTO> getTierPermissions(String organization) throws APIManagementException;

    /**
     * Retrieve Tier Permission by Tier Name.
     *
     * @param tierName     Tier Name
     * @param organization Organization
     * @throws APIManagementException if fails to retrieve Tier permission
     */
    TierPermissionDTO getThrottleTierPermission(String tierName, String organization) throws APIManagementException;

    /**
     * Update Tier Permissions.
     *
     * @param tierName       Tier Name
     * @param permissionType Permission Type
     * @param roles          Roles
     * @param organization   Organization
     * @throws APIManagementException if fails to update Tier permissions
     */
    void updateThrottleTierPermissions(String tierName, String permissionType, String roles, String organization)
            throws APIManagementException;

    /**
     * Retrieve Tier Permissions by Organization.
     *
     * @param organization Organization
     * @throws APIManagementException if fails to retrieve Tier permissions
     */
    Set<TierPermissionDTO> getThrottleTierPermissions(String organization) throws APIManagementException;

}
