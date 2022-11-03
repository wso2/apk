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
import org.wso2.apk.apimgt.api.model.Environment;

import java.util.List;

public interface EnvironmentDAO {

    /**
     * Returns the Environment for the uuid in the tenant domain.
     *
     * @param organization the organization to look environment
     * @param uuid         UUID of the environment
     * @return Gateway environment with given UUID
     */
    Environment getEnvironment(String organization, String uuid) throws APIManagementException;

    /**
     * Returns the Environments List for the Organization.
     *
     * @param organization Organization name.
     * @return List of Environments.
     */
    List<Environment> getAllEnvironments(String organization) throws APIManagementException;

    /**
     * Add an Environment
     *
     * @param organization Organization
     * @param environment  Environment
     * @return added Environment
     * @throws APIManagementException if failed to add environment
     */
    Environment addEnvironment(String organization, Environment environment) throws APIManagementException;

    /**
     * Delete an Environment
     *
     * @param uuid UUID of the environment
     * @throws APIManagementException if failed to delete environment
     */
    void deleteEnvironment(String uuid) throws APIManagementException;

    /**
     * Update Gateway Environment
     *
     * @param environment Environment to be updated
     * @return Updated Environment
     * @throws APIManagementException if failed to updated Environment
     */
    Environment updateEnvironment(Environment environment) throws APIManagementException;



}
