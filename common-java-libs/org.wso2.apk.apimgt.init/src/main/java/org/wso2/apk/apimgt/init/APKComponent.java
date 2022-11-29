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

package org.wso2.apk.apimgt.init;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.wso2.apk.apimgt.api.APIManagementException;
import org.wso2.apk.apimgt.impl.APIManagerConfigurationServiceImpl;
import org.wso2.apk.apimgt.impl.ConfigurationHolder;
import org.wso2.apk.apimgt.impl.factory.SQLConstantManagerFactory;
import org.wso2.apk.apimgt.impl.internal.ServiceReferenceHolder;
import org.wso2.apk.apimgt.impl.utils.APIMgtDBUtil;


public class APKComponent {

    private static final Log log = LogFactory.getLog(APKComponent.class);

    public static void activate(String configuration) throws APIManagementException {

        // Set configurations
        ConfigurationHolder config = new ConfigurationHolder();
        ObjectMapper objectMapper = new ObjectMapper();
        objectMapper.disable(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES);
        try {
            config = objectMapper.readValue(configuration, ConfigurationHolder.class);
            APIManagerConfigurationServiceImpl configurationService = new APIManagerConfigurationServiceImpl(config);
            ServiceReferenceHolder.getInstance().setAPIManagerConfigurationService(configurationService);
            log.info("Successfully set the APK configurations...");
        } catch (JsonProcessingException e) {
            throw new APIManagementException("Error while reading configurations");
        }

        try {
            // Initialize database
            APIMgtDBUtil.initialize();
            log.info("Datasource initialized successfully!");
        } catch (Exception e) {
            throw new APIManagementException("Error while initializing database connection", e);
        }

        // Initialise SQL constant manager
        log.debug("Initializing SQL constant manager");
        SQLConstantManagerFactory.initializeSQLConstantManager();

        //TODO: APK
//        // initialize API-M Caches
//        CacheProvider.createTenantConfigCache();
    }
}
