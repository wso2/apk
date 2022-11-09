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

import com.zaxxer.hikari.HikariDataSource;

import java.sql.Connection;
import java.sql.SQLException;

public class DBDataSource {
    static HikariDataSource basicDataSource = new HikariDataSource();
    static String databaseName = "apimdb";

    DBDataSource() throws Exception {
        String ipAddress = "localhost";
        String port = "5432";
        basicDataSource.setDriverClassName("org.postgresql.Driver");
        basicDataSource.setJdbcUrl("jdbc:postgresql://" + ipAddress + ":" + port + "/" + databaseName);
        basicDataSource.setUsername("apimtest");
        basicDataSource.setPassword("apimtest");
        basicDataSource.setAutoCommit(true);
        basicDataSource.setMaximumPoolSize(20);
    }

    /**
     * Get a {@link Connection} object
     *
     * @return {@link Connection} from given DataSource
     */
    public Connection getConnection() throws SQLException {
        return basicDataSource.getConnection();
    }

    /**
     * Return javax.sql.DataSource object
     *
     * @return {@link javax.sql.DataSource} object
     */
    public HikariDataSource getDatasource() throws SQLException {
        return basicDataSource;
    }
}
