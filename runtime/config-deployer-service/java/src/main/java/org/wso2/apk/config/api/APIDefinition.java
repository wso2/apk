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
package org.wso2.apk.config.api;

import org.wso2.apk.config.model.API;
import org.wso2.apk.config.model.SwaggerData;
import org.wso2.apk.config.model.URITemplate;

import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * APIDefinition is responsible for providing uri templates, scopes and
 * save the api definition according to the permission and visibility
 */

@SuppressWarnings("unused")
public abstract class APIDefinition {

    private static final Pattern CURLY_BRACES_PATTERN = Pattern.compile("(?<=\\{)(?!\\s*\\{)[^{}]+");

    /**
     * This method extracts the URI templates from the API definition
     *
     * @return URI templates
     */
    public abstract Set<URITemplate> getURITemplates(String resourceConfigsJSON) throws APIManagementException;

    /**
     * This method extracts the scopes from the API definition
     *
     * @param resourceConfigsJSON resource json
     * @return scopes
     */
    public abstract String[] getScopes(String resourceConfigsJSON) throws APIManagementException;

    public abstract String generateAPIDefinition(API api) throws APIManagementException;

    /**
     * This method generates API definition using the given api's URI templates and the swagger.
     * It will alter the provided swagger definition based on the URI templates. For example: if there is a new
     * URI template which is not included in the swagger, it will be added to the swagger as a basic resource. Any
     * additional resources inside the swagger will be removed from the swagger. Changes to scopes, throtting policies,
     * on the resource will be updated on the swagger
     *
     * @param swagger     swagger definition
     * @return API definition in string format
     * @throws APIManagementException if error occurred when generating API Definition
     */
    public abstract String generateAPIDefinition(API api, String swagger) throws APIManagementException;

    /**
     * Extract and return path parameters in the given URI template
     *
     * @param uriTemplate URI Template value
     * @return path parameters in the given URI template
     */
    public List<String> getPathParamNames(String uriTemplate) {
        List<String> params = new ArrayList<>();

        Matcher bracesMatcher = CURLY_BRACES_PATTERN.matcher(uriTemplate);
        while (bracesMatcher.find()) {
            params.add(bracesMatcher.group());
        }
        return params;
    }

    /**
     * Creates a helper resource path map using provided swagger data.
     * Creates map in below format:
     * /order      -> [post -> resource1]
     * /order/{id} -> [get -> resource2, put -> resource3, ..]
     *
     * @return a structured uri template map using provided Swagger Data Resource Paths
     */
    public Map<String, Map<String, SwaggerData.Resource>> getResourceMap(API api) {
        SwaggerData swaggerData = new SwaggerData(api);
        Map<String, Map<String, SwaggerData.Resource>> uriTemplateMap = new LinkedHashMap<>();
        for (SwaggerData.Resource resource : swaggerData.getResources()) {
            Map<String, SwaggerData.Resource> resources = uriTemplateMap.computeIfAbsent(resource.getPath(), k -> new LinkedHashMap<>());
            resources.put(resource.getVerb().toUpperCase(), resource);
        }
        return uriTemplateMap;
    }

    /**
     * This method validates the given OpenAPI definition by content
     *
     * @param apiDefinition     OpenAPI Definition content
     * @param returnJsonContent whether to return the converted json form of the OpenAPI definition
     * @return APIDefinitionValidationResponse object with validation information
     */
    public abstract APIDefinitionValidationResponse validateAPIDefinition(String apiDefinition,
                                                                          boolean returnJsonContent) throws APIManagementException;

    public abstract API getAPIFromDefinition(String content) throws APIManagementException;

    /**
     * This method changes the URI templates from the API definition as it support different schemes
     * @param resourceConfigsJSON json String of oasDefinition
     * @throws APIManagementException throws if an error occurred
     * @return String
     */
    public abstract String processOtherSchemeScopes(String resourceConfigsJSON)
            throws APIManagementException;

    /**
     * Get parser Type
     *
     * @return String parserType
     */
    public abstract String getType();
    public abstract boolean canHandleDefinition(String definition);
}
