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
import org.wso2.apk.apimgt.api.dto.KeyManagerConfigurationDTO;

import java.util.List;

public interface KeyManagerDAO {

    /**
     * Get Key Manager Configuration By Organization.
     *
     * @param organization Organization Name
     * @return List of Key Manager Configurations
     * @throws APIManagementException if failed to Configurations
     */
    List<KeyManagerConfigurationDTO> getKeyManagerConfigurationsByOrganization(String organization)
            throws APIManagementException;

    /**
     * Get All Key Manager Configurations.
     *
     * @return List of Key Manager Configurations
     * @throws APIManagementException if failed to Configurations
     */
    List<KeyManagerConfigurationDTO> getKeyManagerConfigurations() throws APIManagementException;

    /**
     * Get Key Manager Configuration By Organization and UUID.
     *
     * @param organization Organization Name
     * @param uuid KeyManager UUID
     * @return List of Key Manager Configurations
     * @throws APIManagementException if failed to Configurations
     */
    KeyManagerConfigurationDTO getKeyManagerConfigurationByID(String organization, String uuid)
            throws APIManagementException;

    /**
     * Checks Key Manager exists in given organization.
     *
     * @param organization Organization Name
     * @param resourceId KeyManager UUID
     * @return boolean
     * @throws APIManagementException
     */
    boolean isIDPExistInOrg(String organization, String resourceId) throws APIManagementException;

    /**
     * Checks Key Manager Configuration exists in given organization by UUID.
     *
     * @param organization Organization Name
     * @param resourceId KeyManager UUID
     * @return boolean
     * @throws APIManagementException
     */
    boolean isKeyManagerConfigurationExistById(String organization, String resourceId) throws APIManagementException;

    /**
     * Checks Key Manager Configuration exists in given organization by UUID.
     *
     * @param organization Organization Name
     * @param name KeyManager Name
     * @return boolean
     * @throws APIManagementException
     */
    boolean isKeyManagerConfigurationExistByName(String organization, String name)
            throws APIManagementException;

    /**
     * Add Key Manager Configuration
     *
     * @param keyManagerConfigurationDTO Key Manager Configuration DTO
     * @throws APIManagementException if error
     */
    void addKeyManagerConfiguration(KeyManagerConfigurationDTO keyManagerConfigurationDTO)
            throws APIManagementException;

    /**
     * Update Key Manager Configuration
     *
     * @param keyManagerConfigurationDTO Key Manager Configuration DTO
     * @throws APIManagementException if error
     */
    void updateKeyManagerConfiguration(KeyManagerConfigurationDTO keyManagerConfigurationDTO)
            throws APIManagementException;

    /**
     * Delete Key Manager Configuration
     *
     * @param id Key Manager UUID
     * @param organization Organization Name
     * @throws APIManagementException if error
     */
    void deleteKeyManagerConfigurationById(String id, String organization) throws APIManagementException;

    /**
     * Retrieve Key Manager Configuration by Name
     *
     * @param name Key Manager Name
     * @param organization Organization Name
     * @throws APIManagementException if error
     */
    KeyManagerConfigurationDTO getKeyManagerConfigurationByName(String organization, String name)
            throws APIManagementException;

    /**
     * Get Key Manager Configuration By UUID.
     *
     * @param uuid KeyManager UUID
     * @return Key Manager Configurations DTO
     * @throws APIManagementException if failed to retrieve Configuration
     */
    KeyManagerConfigurationDTO getKeyManagerConfigurationByUUID(String uuid) throws APIManagementException;
}
