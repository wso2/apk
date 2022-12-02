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

package org.wso2.apk.apimgt.impl.dao.constants;

/**
 * Keep the constants related to environment specific api property DAO.
 */
public class EnvironmentSpecificAPIPropertyConstants {
    public static final String ADD_ENVIRONMENT_SPECIFIC_API_PROPERTIES_SQL =
            "INSERT INTO API_ENVIRONMENT_KEYS(UUID, API_UUID, ENVIRONMENT_ID, PROPERTY_CONFIG) VALUES(?,?,?,?)";

    public static final String UPDATE_ENVIRONMENT_SPECIFIC_API_PROPERTIES_SQL =
            "UPDATE API_ENVIRONMENT_KEYS SET PROPERTY_CONFIG = ? WHERE API_UUID=? AND ENVIRONMENT_ID=?";

    public static final String GET_ENVIRONMENT_SPECIFIC_API_PROPERTIES_SQL =
            "SELECT PROPERTY_CONFIG FROM API_ENVIRONMENT_KEYS WHERE API_UUID=? AND ENVIRONMENT_ID=?";

    public static final String IS_ENVIRONMENT_SPECIFIC_API_PROPERTIES_EXIST_SQL =
            "SELECT 1 FROM API_ENVIRONMENT_KEYS WHERE API_UUID=? AND ENVIRONMENT_ID=?";

    public static final String GET_ENVIRONMENT_SPECIFIC_API_PROPERTIES_BY_APIS_SQL =
            "SELECT GATEWAY_ENVIRONMENT.UUID ENV_ID,"
                    + "       GATEWAY_ENVIRONMENT.NAME ENV_NAME,"
                    + "       API_ENVIRONMENT_KEYS.API_UUID API_UUID,"
                    + "       API_ENVIRONMENT_KEYS.PROPERTY_CONFIG CONFIG"
                    + " FROM API_ENVIRONMENT_KEYS,GATEWAY_ENVIRONMENT"
                    + " WHERE API_ENVIRONMENT_KEYS.ENVIRONMENT_ID = GATEWAY_ENVIRONMENT.UUID AND"
                    + "        API_ENVIRONMENT_KEYS.API_UUID IN (_API_ID_LIST_)"
                    + " ORDER BY API_UUID, ENV_NAME, ENV_ID";

    public static final String GET_ENVIRONMENT_SPECIFIC_API_PROPERTIES_BY_APIS_ENVS_SQL =
            "SELECT API_ENVIRONMENT_KEYS.ENVIRONMENT_ID ENV_ID,"
                    + "       API_ENVIRONMENT_KEYS.API_UUID API_UUID,"
                    + "       API_ENVIRONMENT_KEYS.PROPERTY_CONFIG CONFIG"
                    + " FROM API_ENVIRONMENT_KEYS"
                    + " WHERE API_ENVIRONMENT_KEYS.ENVIRONMENT_ID IN (_ENV_ID_LIST_) AND"
                    + "        API_ENVIRONMENT_KEYS.API_UUID IN (_API_ID_LIST_)"
                    + " ORDER BY API_UUID, ENV_ID";;
}
