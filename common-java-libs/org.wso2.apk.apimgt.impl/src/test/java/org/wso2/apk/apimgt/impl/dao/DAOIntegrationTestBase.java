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

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.testng.annotations.BeforeTest;
import org.wso2.apk.apimgt.impl.utils.APIMgtDBUtil;

import java.sql.Connection;

public class DAOIntegrationTestBase {
    protected DBDataSource dataSource;
    private static final Logger log = LoggerFactory.getLogger(DAOIntegrationTestBase.class);

    public DAOIntegrationTestBase() {
    }

    @BeforeTest
    public void init() throws Exception {
        dataSource = new DBDataSource();
        int maxRetries = 5;
        long maxWait = 5000;
        while (maxRetries > 0) {
            try (Connection connection = dataSource.getConnection()) {
                log.info("Database Connection Successful");
                APIMgtDBUtil.initialize(dataSource.getDatasource());
                break;
            } catch (Exception e) {
                if (maxRetries > 0) {
                    log.warn("Couldn't connect into database retrying after next 5 seconds");
                    maxRetries--;
                    try {
                        Thread.sleep(maxWait);
                    } catch (InterruptedException e1) {
                    }
                } else {
                    log.error("Max tries 5 exceed to connect");
                    throw e;
                }
            }
        }
    }
}