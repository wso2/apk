/*
 * Copyright (c) 2021, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.wso2.apk.apimgt.devportal.utils;

import org.wso2.apk.apimgt.api.model.API;
import org.wso2.apk.apimgt.api.model.APIProduct;
import org.wso2.apk.apimgt.devportal.dto.APIEndpointURLsInnerDefaultVersionURLsDTO;
import org.wso2.apk.apimgt.devportal.dto.APIEndpointURLsInnerDTO;
import org.wso2.apk.apimgt.devportal.dto.APIEndpointURLsInnerURLsDTO;
import org.wso2.apk.apimgt.api.APIConsumer;
import org.wso2.apk.apimgt.api.APIManagementException;
import org.wso2.apk.apimgt.api.model.Environment;
import org.wso2.apk.apimgt.impl.APIConstants;
import org.wso2.apk.apimgt.impl.utils.APIUtil;
import org.wso2.apk.apimgt.rest.api.util.utils.RestApiCommonUtil;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.HashMap;
import java.util.HashSet;
import java.util.List;
import java.util.Map;
import java.util.Set;

/**
 *  This is the util implementation class for store rest API
 */
public class APIUtils {
    /**
     * Extracts the API environment details with access url for each endpoint
     *
     * @param api          API object
     * @param tenantDomain Tenant domain of the API
     * @return the API environment details
     * @throws APIManagementException error while extracting the information
     */
    public static List<APIEndpointURLsInnerDTO> extractEndpointURLs(API api, String tenantDomain)
            throws APIManagementException {
        List<APIEndpointURLsInnerDTO> apiEndpointsList = new ArrayList<>();
        String organization = api.getOrganization();
        Map<String, Environment> environments = APIUtil.getEnvironments(organization);

        Set<String> environmentsPublishedByAPI = new HashSet<>(api.getEnvironments());
        environmentsPublishedByAPI.remove("none");

        Set<String> apiTransports = new HashSet<>(Arrays.asList(api.getTransports().split(",")));
        APIConsumer apiConsumer = RestApiCommonUtil.getLoggedInUserConsumer();

        for (String environmentName : environmentsPublishedByAPI) {
            Environment environment = environments.get(environmentName);
            if (environment != null) {
                APIEndpointURLsInnerURLsDTO APIEndpointURLsInnerURLsDTO = new APIEndpointURLsInnerURLsDTO();
                APIEndpointURLsInnerDefaultVersionURLsDTO APIEndpointURLsInnerDefaultVersionURLsDTO = new APIEndpointURLsInnerDefaultVersionURLsDTO();
                String[] gwEndpoints = null;
                if ("WS".equalsIgnoreCase(api.getType())) {
                    gwEndpoints = environment.getWebsocketGatewayEndpoint().split(",");
                } else {
                    gwEndpoints = environment.getApiGatewayEndpoint().split(",");
                }
                Map<String, String> domains = new HashMap<>();
                if (tenantDomain != null) {
                    domains = apiConsumer.getTenantDomainMappings(tenantDomain,
                            APIConstants.API_DOMAIN_MAPPINGS_GATEWAY);
                }

                String customGatewayUrl = null;
                if (domains != null) {
                    customGatewayUrl = domains.get(APIConstants.CUSTOM_URL);
                }

                for (String gwEndpoint : gwEndpoints) {
                    StringBuilder endpointBuilder = new StringBuilder(gwEndpoint);

                    if (customGatewayUrl != null) {
                        int index = endpointBuilder.indexOf("//");
                        endpointBuilder.replace(index + 2, endpointBuilder.length(), customGatewayUrl);
                        endpointBuilder.append(api.getContext().replace("/t/" + tenantDomain, ""));
                    } else {
                        endpointBuilder.append(api.getContext());
                    }

                    if (gwEndpoint.contains("http:") && apiTransports.contains("http")) {
                        APIEndpointURLsInnerURLsDTO.setHttp(endpointBuilder.toString());
                    } else if (gwEndpoint.contains("https:") && apiTransports.contains("https")) {
                        APIEndpointURLsInnerURLsDTO.setHttps(endpointBuilder.toString());
                    } else if (gwEndpoint.contains("ws:")) {
                        APIEndpointURLsInnerURLsDTO.setWs(endpointBuilder.toString());
                    } else if (gwEndpoint.contains("wss:")) {
                        APIEndpointURLsInnerURLsDTO.setWss(endpointBuilder.toString());
                    }

                    if (api.isDefaultVersion()) {
                        int index = endpointBuilder.lastIndexOf(api.getId().getVersion());
                        endpointBuilder.replace(index, endpointBuilder.length(), "");
                        if (gwEndpoint.contains("http:") && apiTransports.contains("http")) {
                            APIEndpointURLsInnerDefaultVersionURLsDTO.setHttp(endpointBuilder.toString());
                        } else if (gwEndpoint.contains("https:") && apiTransports.contains("https")) {
                            APIEndpointURLsInnerDefaultVersionURLsDTO.setHttps(endpointBuilder.toString());
                        } else if (gwEndpoint.contains("ws:")) {
                            APIEndpointURLsInnerDefaultVersionURLsDTO.setWs(endpointBuilder.toString());
                        } else if (gwEndpoint.contains("wss:")) {
                            APIEndpointURLsInnerDefaultVersionURLsDTO.setWss(endpointBuilder.toString());
                        }
                    }
                }

                APIEndpointURLsInnerDTO APIEndpointURLsInnerDTO = new APIEndpointURLsInnerDTO();
                APIEndpointURLsInnerDTO.setDefaultVersionURLs(APIEndpointURLsInnerDefaultVersionURLsDTO);
                APIEndpointURLsInnerDTO.setUrLs(APIEndpointURLsInnerURLsDTO);

                APIEndpointURLsInnerDTO.setEnvironmentName(environment.getName());
                APIEndpointURLsInnerDTO.setEnvironmentType(environment.getType());

                apiEndpointsList.add(APIEndpointURLsInnerDTO);
            }
        }

        return apiEndpointsList;
    }


    /**
     * Extracts the API environment details with access url for each endpoint
     *
     * @param apiProduct   API object
     * @param tenantDomain Tenant domain of the API
     * @return the API environment details
     * @throws APIManagementException error while extracting the information
     */
    public static List<APIEndpointURLsInnerDTO> extractEndpointURLs(APIProduct apiProduct, String tenantDomain)
            throws APIManagementException {
        List<APIEndpointURLsInnerDTO> apiEndpointsList = new ArrayList<>();
        String organization = apiProduct.getOrganization();
        Map<String, Environment> environments = APIUtil.getEnvironments(organization);

        Set<String> environmentsPublishedByAPI = new HashSet<>(apiProduct.getEnvironments());
        environmentsPublishedByAPI.remove("none");

        Set<String> apiTransports = new HashSet<>(Arrays.asList(apiProduct.getTransports().split(",")));
        APIConsumer apiConsumer = RestApiCommonUtil.getLoggedInUserConsumer();

        for (String environmentName : environmentsPublishedByAPI) {
            Environment environment = environments.get(environmentName);
            if (environment != null) {
                APIEndpointURLsInnerURLsDTO APIEndpointURLsInnerURLsDTO = new APIEndpointURLsInnerURLsDTO();
                String[] gwEndpoints = null;
                gwEndpoints = environment.getApiGatewayEndpoint().split(",");

                Map<String, String> domains = new HashMap<>();
                if (tenantDomain != null) {
                    domains = apiConsumer.getTenantDomainMappings(tenantDomain,
                            APIConstants.API_DOMAIN_MAPPINGS_GATEWAY);
                }

                String customGatewayUrl = null;
                if (domains != null) {
                    customGatewayUrl = domains.get(APIConstants.CUSTOM_URL);
                }

                for (String gwEndpoint : gwEndpoints) {
                    StringBuilder endpointBuilder = new StringBuilder(gwEndpoint);

                    if (customGatewayUrl != null) {
                        int index = endpointBuilder.indexOf("//");
                        endpointBuilder.replace(index + 2, endpointBuilder.length(), customGatewayUrl);
                        endpointBuilder.append(apiProduct.getContext().replace("/t/" + tenantDomain, ""));
                    } else {
                        endpointBuilder.append(apiProduct.getContext());
                    }

                    if (gwEndpoint.contains("http:") && apiTransports.contains("http")) {
                        APIEndpointURLsInnerURLsDTO.setHttp(endpointBuilder.toString());
                    } else if (gwEndpoint.contains("https:") && apiTransports.contains("https")) {
                        APIEndpointURLsInnerURLsDTO.setHttps(endpointBuilder.toString());
                    }
                }

                APIEndpointURLsInnerDTO APIEndpointURLsInnerDTO = new APIEndpointURLsInnerDTO();
                APIEndpointURLsInnerDTO.setUrLs(APIEndpointURLsInnerURLsDTO);
                APIEndpointURLsInnerDTO.setEnvironmentName(environment.getName());
                APIEndpointURLsInnerDTO.setEnvironmentType(environment.getType());

                apiEndpointsList.add(APIEndpointURLsInnerDTO);
            }
        }

        return apiEndpointsList;
    }
}
