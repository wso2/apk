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

package org.wso2.apk.apimgt.impl.utils;

import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.wso2.apk.apimgt.api.APIManagementException;
import org.wso2.apk.apimgt.api.APIManagerDatabaseException;

import javax.naming.Context;
import javax.naming.InitialContext;
import javax.naming.NamingException;
import javax.sql.DataSource;
import java.sql.Connection;
import java.sql.ResultSet;
import java.sql.SQLException;

public class GatewayArtifactsMgtDBUtil {

    private static final Log log = LogFactory.getLog(GatewayArtifactsMgtDBUtil.class);
    private static volatile DataSource artifactSynchronizerDataSource = null;

    /**
     * Initializes the data source
     *
     * @throws APIManagementException if an error occurs while loading DB configuration
     */
    public static void initialize() throws APIManagerDatabaseException {
        if (artifactSynchronizerDataSource != null) {
            return;
        }
        initDatasource();
    }

    private static synchronized void initDatasource() throws APIManagerDatabaseException {
        if (artifactSynchronizerDataSource == null) {
            if (log.isDebugEnabled()) {
                log.debug("Initializing data source");
            }

            // TODO Read from Config
            String artifactSynchronizerDataSourceName = "";
//            GatewayArtifactSynchronizerProperties gatewayArtifactSynchronizerProperties =
//                    ServiceReferenceHolder.getInstance().getAPIManagerConfigurationService()
//                            .getAPIManagerConfiguration().getGatewayArtifactSynchronizerProperties();
//            String artifactSynchronizerDataSourceName =
//                    gatewayArtifactSynchronizerProperties.getArtifactSynchronizerDataSource();

            if (artifactSynchronizerDataSourceName != null) {
                try {
                    Context ctx = new InitialContext();
                    artifactSynchronizerDataSource = (DataSource) ctx.lookup(artifactSynchronizerDataSourceName);
                } catch (NamingException e) {
                    throw new APIManagerDatabaseException("Error while looking up the data " +
                            "source: " + artifactSynchronizerDataSourceName, e);
                }
            } else {
                log.error(artifactSynchronizerDataSourceName + " not defined in api-manager.xml.");
            }
        }
    }

    /**
     * Utility method to get a new database connection for gatewayRuntime artifacts
     *
     * @return Connection
     * @throws SQLException if failed to get Connection
     */
    public static Connection getArtifactSynchronizerConnection() throws SQLException {
        if (artifactSynchronizerDataSource != null) {
            return artifactSynchronizerDataSource.getConnection();
        }
        throw new SQLException("Data source is not configured properly.");
    }

    /**
     * Close ResultSet
     *
     * @param resultSet ResultSet
     */
    public static void closeResultSet(ResultSet resultSet) {
        if (resultSet != null) {
            try {
                resultSet.close();
            } catch (SQLException e) {
                log.warn("Database error. Could not close ResultSet  - " + e.getMessage(), e);
            }
        }
    }
}
