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

package org.wso2.apk.apimgt.rest.api.util.utils;

import org.apache.commons.io.IOUtils;
import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.wso2.apk.apimgt.api.*;
import org.wso2.apk.apimgt.api.model.Scope;
import org.wso2.apk.apimgt.impl.definitions.OASParserUtil;

import java.io.IOException;
import java.util.*;

public class RestApiUtil {

    public static final Log log = LogFactory.getLog(RestApiUtil.class);

    /**
     * This method is used to get the scope list from the yaml file
     *
     * @return MAP of scope list for all portal
     */
    public static  Map<String, List<String>> getScopesInfoFromAPIYamlDefinitions() throws APIManagementException {

        Map<String, List<String>>   portalScopeList = new HashMap<>();
        String [] fileNameArray = {"/admin-api.yaml", "/publisher-api.yaml", "/devportal-api.yaml", "/service-catalog-api.yaml"};
        for (String fileName : fileNameArray) {
            String definition = null;
            try {
                definition = IOUtils
                        .toString(RestApiUtil.class.getResourceAsStream(fileName), "UTF-8");
            } catch (IOException  e) {
                throw new APIManagementException("Error while reading the swagger definition ,",
                        ExceptionCodes.DEFINITION_EXCEPTION);
            }
            APIDefinition oasParser = OASParserUtil.getOASParser(definition);
            Set<Scope> scopeSet = oasParser.getScopes(definition);
            for (Scope entry : scopeSet) {
                List<String> list = new ArrayList<>();
                list.add(entry.getDescription());
                list.add((fileName.replaceAll("-api.yaml", "").replace("/", "")));
                if (("/service-catalog-api.yaml".equals(fileName))) {
                    if (!entry.getKey().contains("apim:api_view")) {
                        portalScopeList.put(entry.getName(), list);
                    }
                } else {
                    portalScopeList.put(entry.getName(), list);
                }
            }
        }
        return portalScopeList;
    }
}
