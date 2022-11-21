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

package org.wso2.apk.apimgt.backoffice.utils.mappings;

import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.commons.io.FileUtils;
import org.apache.commons.io.IOUtils;
import org.apache.commons.lang3.ArrayUtils;
import org.apache.commons.lang3.RandomStringUtils;
import org.apache.commons.lang3.StringUtils;
import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.json.simple.JSONObject;
import org.json.simple.parser.JSONParser;
import org.json.simple.parser.ParseException;
import org.wso2.apk.apimgt.api.APIManagementException;
import org.wso2.apk.apimgt.api.APIProvider;
import org.wso2.apk.apimgt.api.ExceptionCodes;
import org.wso2.apk.apimgt.api.FaultGatewaysException;
import org.wso2.apk.apimgt.api.model.API;
import org.wso2.apk.apimgt.api.model.APICategory;
import org.wso2.apk.apimgt.api.model.APIIdentifier;
import org.wso2.apk.apimgt.api.model.APIProductIdentifier;
import org.wso2.apk.apimgt.api.model.APIStateChangeResponse;
import org.wso2.apk.apimgt.api.model.ApiTypeWrapper;
import org.wso2.apk.apimgt.api.model.Documentation;
import org.wso2.apk.apimgt.api.model.DocumentationContent;
import org.wso2.apk.apimgt.api.model.Identifier;
import org.wso2.apk.apimgt.api.model.LifeCycleEvent;
import org.wso2.apk.apimgt.api.model.OperationPolicyData;
import org.wso2.apk.apimgt.api.model.ResourceFile;
import org.wso2.apk.apimgt.api.model.Tier;
import org.wso2.apk.apimgt.backoffice.dto.LifecycleHistoryDTO;
import org.wso2.apk.apimgt.backoffice.dto.LifecycleStateDTO;
import org.wso2.apk.apimgt.impl.APIConstants;
import org.wso2.apk.apimgt.impl.importexport.APIImportExportException;
import org.wso2.apk.apimgt.impl.importexport.ExportFormat;
import org.wso2.apk.apimgt.impl.importexport.ImportExportConstants;
import org.wso2.apk.apimgt.impl.importexport.utils.CommonUtil;
import org.wso2.apk.apimgt.impl.utils.APIUtil;
import org.wso2.apk.apimgt.impl.utils.APIVersionStringComparator;
import org.wso2.apk.apimgt.backoffice.utils.crypto.CryptoTool;
import org.wso2.apk.apimgt.backoffice.utils.crypto.CryptoToolException;
import org.wso2.apk.apimgt.backoffice.dto.APIDTO;
import org.wso2.apk.apimgt.backoffice.dto.DocumentDTO;
import org.wso2.apk.apimgt.rest.api.util.utils.RestApiCommonUtil;
import org.wso2.apk.apimgt.rest.api.util.RestApiConstants;
import org.wso2.apk.apimgt.rest.api.util.annotations.Scope;

import java.io.ByteArrayInputStream;
import java.io.ByteArrayOutputStream;
import java.io.File;
import java.io.FileInputStream;
import java.io.FileNotFoundException;
import java.io.FileOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.lang.reflect.Field;
import java.nio.charset.StandardCharsets;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.HashMap;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.Set;

/**
 * This is a publisher rest api utility class.
 */
public class PublisherCommonUtils {

    private static final Log log = LogFactory.getLog(PublisherCommonUtils.class);

    /**
     * Update an API.
     *
     * @param originalAPI    Existing API
     * @param apiDtoToUpdate New API DTO to update
     * @param apiProvider    API Provider
     * @param tokenScopes    Scopes of the token
     * @throws ParseException         If an error occurs while parsing the endpoint configuration
     * @throws APIManagementException If an error occurs while updating the API
     * @throws FaultGatewaysException If an error occurs while updating manage of an existing API
     */
    public static API updateApi(API originalAPI, APIDTO apiDtoToUpdate, APIProvider apiProvider, String[] tokenScopes)
            throws APIManagementException, FaultGatewaysException {

        APIIdentifier apiIdentifier = originalAPI.getId();
        // Validate if the USER_REST_API_SCOPES is not set in WebAppAuthenticator when scopes are validated
        if (tokenScopes == null) {
            throw new APIManagementException("Error occurred while updating the  API " + originalAPI.getUUID()
                    + " as the token information hasn't been correctly set internally",
                    ExceptionCodes.TOKEN_SCOPES_NOT_SET);
        }

        Scope[] apiDtoClassAnnotatedScopes = APIDTO.class.getAnnotationsByType(Scope.class);
        boolean hasClassLevelScope = checkClassScopeAnnotation(apiDtoClassAnnotatedScopes, tokenScopes);

        JSONParser parser = new JSONParser();
        String oldEndpointConfigString = originalAPI.getEndpointConfig();
        JSONObject oldEndpointConfig = null;
        if (StringUtils.isNotBlank(oldEndpointConfigString)) {
            try {
                oldEndpointConfig = (JSONObject) parser.parse(oldEndpointConfigString);
            } catch (ParseException e) {
                throw new APIManagementException("Error while parsing endpoint config",
                        ExceptionCodes.JSON_PARSE_ERROR);
            }
        }
        String oldProductionApiSecret = null;
        String oldSandboxApiSecret = null;

        if (oldEndpointConfig != null) {
            if ((oldEndpointConfig.containsKey(APIConstants.ENDPOINT_SECURITY))) {
                JSONObject oldEndpointSecurity = (JSONObject) oldEndpointConfig.get(APIConstants.ENDPOINT_SECURITY);
                if (oldEndpointSecurity.containsKey(APIConstants.OAuthConstants.ENDPOINT_SECURITY_PRODUCTION)) {
                    JSONObject oldEndpointSecurityProduction = (JSONObject) oldEndpointSecurity
                            .get(APIConstants.OAuthConstants.ENDPOINT_SECURITY_PRODUCTION);

                    if (oldEndpointSecurityProduction.get(APIConstants.OAuthConstants.OAUTH_CLIENT_ID) != null
                            && oldEndpointSecurityProduction.get(APIConstants.OAuthConstants.OAUTH_CLIENT_SECRET)
                            != null) {
                        oldProductionApiSecret = oldEndpointSecurityProduction
                                .get(APIConstants.OAuthConstants.OAUTH_CLIENT_SECRET).toString();
                    }
                }
                if (oldEndpointSecurity.containsKey(APIConstants.OAuthConstants.ENDPOINT_SECURITY_SANDBOX)) {
                    JSONObject oldEndpointSecuritySandbox = (JSONObject) oldEndpointSecurity
                            .get(APIConstants.OAuthConstants.ENDPOINT_SECURITY_SANDBOX);

                    if (oldEndpointSecuritySandbox.get(APIConstants.OAuthConstants.OAUTH_CLIENT_ID) != null
                            && oldEndpointSecuritySandbox.get(APIConstants.OAuthConstants.OAUTH_CLIENT_SECRET)
                            != null) {
                        oldSandboxApiSecret = oldEndpointSecuritySandbox
                                .get(APIConstants.OAuthConstants.OAUTH_CLIENT_SECRET).toString();
                    }
                }
            }
        }

//        Map endpointConfig = (Map) apiDtoToUpdate.getEndpointConfig();
//        CryptoTool cryptoTool = CryptoToolUtil.getDefaultCryptoTool();
//
//        // OAuth 2.0 backend protection: API Key and API Secret encryption
//        encryptEndpointSecurityOAuthCredentials(endpointConfig, cryptoTool, oldProductionApiSecret, oldSandboxApiSecret,
//                apiDtoToUpdate);
//
//        // AWS Lambda: secret key encryption while updating the API
//        if (apiDtoToUpdate.getEndpointConfig() != null) {
//            if (endpointConfig.containsKey(APIConstants.AMZN_SECRET_KEY)) {
//                String secretKey = (String) endpointConfig.get(APIConstants.AMZN_SECRET_KEY);
//                if (!StringUtils.isEmpty(secretKey)) {
//                    if (!APIConstants.AWS_SECRET_KEY.equals(secretKey)) {
//                        try {
//                            String encryptedSecretKey = cryptoTool.encryptAndBase64Encode(secretKey.getBytes());
//                            endpointConfig.put(APIConstants.AMZN_SECRET_KEY, encryptedSecretKey);
//                            apiDtoToUpdate.setEndpointConfig(endpointConfig);
//                        } catch (CryptoToolException e) {
//                            throw new APIManagementException(ExceptionCodes.from(ExceptionCodes.ENDPOINT_CRYPTO_ERROR,
//                                    "Error while encrypting AWS secret key"));
//                        }
//
//                    } else {
//                        try {
//                            JSONParser jsonParser = new JSONParser();
//                            JSONObject originalEndpointConfig = (JSONObject) jsonParser
//                                    .parse(originalAPI.getEndpointConfig());
//                            String encryptedSecretKey = (String) originalEndpointConfig
//                                    .get(APIConstants.AMZN_SECRET_KEY);
//                            endpointConfig.put(APIConstants.AMZN_SECRET_KEY, encryptedSecretKey);
//                            apiDtoToUpdate.setEndpointConfig(endpointConfig);
//                        } catch (ParseException e) {
//                            throw new APIManagementException("Error while parsing endpoint config",
//                                    ExceptionCodes.JSON_PARSE_ERROR);
//                        }
//                    }
//                }
//            }
//        }

        if (!hasClassLevelScope) {
            // Validate per-field scopes
            apiDtoToUpdate = getFieldOverriddenAPIDTO(apiDtoToUpdate, originalAPI, tokenScopes);
        }
        //Overriding some properties:
        //API Name change not allowed if OnPrem
//        if (APIUtil.isOnPremResolver()) {
//            apiDtoToUpdate.setName(apiIdentifier.getApiName());
//        }
        apiDtoToUpdate.setVersion(apiIdentifier.getVersion());
        //apiDtoToUpdate.setProvider(apiIdentifier.getProviderName());
        apiDtoToUpdate.setContext(originalAPI.getContextTemplate());
        //apiDtoToUpdate.setLifeCycleStatus(originalAPI.getStatus());
        apiDtoToUpdate.setType(APIDTO.TypeEnum.fromValue(originalAPI.getType()));

//        // Validate API Security
//        List<String> apiSecurity = apiDtoToUpdate.getSecurityScheme();
//        //validation for tiers
//        List<String> tiersFromDTO = apiDtoToUpdate.getPolicies();
//        String originalStatus = originalAPI.getStatus();
//        if (apiSecurity.contains(APIConstants.DEFAULT_API_SECURITY_OAUTH2) || apiSecurity
//                .contains(APIConstants.API_SECURITY_API_KEY)) {
//            if ((tiersFromDTO == null || tiersFromDTO.isEmpty() && !(APIConstants.CREATED.equals(originalStatus)
//                    || APIConstants.PROTOTYPED.equals(originalStatus)))
//                    && !apiDtoToUpdate.getAdvertiseInfo().isAdvertised()) {
//                throw new APIManagementException(
//                        "A tier should be defined if the API is not in CREATED or PROTOTYPED state",
//                        ExceptionCodes.TIER_CANNOT_BE_NULL);
//            }
//        }

//        if (tiersFromDTO != null && !tiersFromDTO.isEmpty()) {
//            //check whether the added API's tiers are all valid
//            Set<Tier> definedTiers = apiProvider.getTiers();
//            List<String> invalidTiers = getInvalidTierNames(definedTiers, tiersFromDTO);
//            if (invalidTiers.size() > 0) {
//                throw new APIManagementException(
//                        "Specified tier(s) " + Arrays.toString(invalidTiers.toArray()) + " are invalid",
//                        ExceptionCodes.TIER_NAME_INVALID);
//            }
//        }
//        if (apiDtoToUpdate.getAccessControlRoles() != null) {
//            String errorMessage = validateUserRoles(apiDtoToUpdate.getAccessControlRoles());
//            if (!errorMessage.isEmpty()) {
//                throw new APIManagementException(errorMessage, ExceptionCodes.INVALID_USER_ROLES);
//            }
//        }
//        if (apiDtoToUpdate.getVisibleRoles() != null) {
//            String errorMessage = validateRoles(apiDtoToUpdate.getVisibleRoles());
//            if (!errorMessage.isEmpty()) {
//                throw new APIManagementException(errorMessage, ExceptionCodes.INVALID_USER_ROLES);
//            }
//        }
//        if (apiDtoToUpdate.getAdditionalProperties() != null) {
//            String errorMessage = validateAdditionalProperties(apiDtoToUpdate.getAdditionalProperties());
//            if (!errorMessage.isEmpty()) {
//                throw new APIManagementException(errorMessage, ExceptionCodes
//                        .from(ExceptionCodes.INVALID_ADDITIONAL_PROPERTIES, apiDtoToUpdate.getName(),
//                                apiDtoToUpdate.getVersion()));
//            }
//        }
        // Validate if resources are empty
        if (apiDtoToUpdate.getOperations() == null || apiDtoToUpdate.getOperations().isEmpty()) {
            throw new APIManagementException(ExceptionCodes.NO_RESOURCES_FOUND);
        }
        API apiToUpdate = APIMappingUtil.fromDTOtoAPI(apiDtoToUpdate, apiIdentifier.getProviderName());
        if (APIConstants.PUBLIC_STORE_VISIBILITY.equals(apiToUpdate.getVisibility())) {
            apiToUpdate.setVisibleRoles(StringUtils.EMPTY);
        }
        apiToUpdate.setUUID(originalAPI.getUUID());
        apiToUpdate.setOrganization(originalAPI.getOrganization());
        //validateScopes(apiToUpdate);
        apiToUpdate.setThumbnailUrl(originalAPI.getThumbnailUrl());

        //preserve monetization status in the update flow
        //apiProvider.configureMonetizationInAPIArtifact(originalAPI); ////////////TODO /////////REG call
        //apiToUpdate.setWsdlUrl(apiDtoToUpdate.getWsdlUrl());
        //apiToUpdate.setGatewayType(apiDtoToUpdate.getGatewayType());

        //validate API categories
        List<APICategory> apiCategories = apiToUpdate.getApiCategories();
        List<APICategory> apiCategoriesList = new ArrayList<>();
        for (APICategory category : apiCategories) {
            category.setOrganization(originalAPI.getOrganization());
            apiCategoriesList.add(category);
        }
        apiToUpdate.setApiCategories(apiCategoriesList);
        if (apiCategoriesList.size() > 0) {
            if (!APIUtil.validateAPICategories(apiCategoriesList, originalAPI.getOrganization())) {
                throw new APIManagementException("Invalid API Category name(s) defined",
                        ExceptionCodes.from(ExceptionCodes.API_CATEGORY_INVALID));
            }
        }

        apiToUpdate.setOrganization(originalAPI.getOrganization());
        apiProvider.updateAPI(apiToUpdate, originalAPI);

        return apiProvider.getAPIbyUUID(originalAPI.getUuid(), originalAPI.getOrganization());
        // TODO use returend api
    }

    /**
     * This method will encrypt the OAuth 2.0 API Key and API Secret
     *
     * @param endpointConfig         endpoint configuration of API
     * @param cryptoTool             cryptography util
     * @param oldProductionApiSecret existing production API secret
     * @param oldSandboxApiSecret    existing sandbox API secret
     * @param apidto                 API DTO
     * @throws APIManagementException if an error occurs due to a problem in the endpointConfig payload
     */
    public static void encryptEndpointSecurityOAuthCredentials(Map endpointConfig, CryptoTool cryptoTool,
            String oldProductionApiSecret, String oldSandboxApiSecret, APIDTO apidto)
            throws APIManagementException {
        // OAuth 2.0 backend protection: API Key and API Secret encryption
        String customParametersString;
        if (endpointConfig != null) {
            if ((endpointConfig.get(APIConstants.ENDPOINT_SECURITY) != null)) {
                Map endpointSecurity = (Map) endpointConfig.get(APIConstants.ENDPOINT_SECURITY);
                if (endpointSecurity.get(APIConstants.OAuthConstants.ENDPOINT_SECURITY_PRODUCTION) != null) {
                    Map endpointSecurityProduction = (Map) endpointSecurity
                            .get(APIConstants.OAuthConstants.ENDPOINT_SECURITY_PRODUCTION);
                    String productionEndpointType = (String) endpointSecurityProduction
                            .get(APIConstants.OAuthConstants.ENDPOINT_SECURITY_TYPE);

                    // Change default value of customParameters JSONObject to String
                    if (!(endpointSecurityProduction
                            .get(APIConstants.OAuthConstants.OAUTH_CUSTOM_PARAMETERS) instanceof String)) {
                        LinkedHashMap<String, String> customParametersHashMap = (LinkedHashMap<String, String>)
                                endpointSecurityProduction.get(APIConstants.OAuthConstants.OAUTH_CUSTOM_PARAMETERS);
                        customParametersString = JSONObject.toJSONString(customParametersHashMap);
                    } else if (endpointSecurityProduction.get(APIConstants.OAuthConstants.OAUTH_CUSTOM_PARAMETERS)
                            != null) {
                        customParametersString = (String) endpointSecurityProduction
                                .get(APIConstants.OAuthConstants.OAUTH_CUSTOM_PARAMETERS);
                    } else {
                        customParametersString = "{}";
                    }

                    endpointSecurityProduction
                            .put(APIConstants.OAuthConstants.OAUTH_CUSTOM_PARAMETERS, customParametersString);

                    if (APIConstants.OAuthConstants.OAUTH.equals(productionEndpointType)) {
                        if (endpointSecurityProduction.get(APIConstants.OAuthConstants.OAUTH_CLIENT_SECRET) != null
                                && StringUtils.isNotBlank(
                                endpointSecurityProduction.get(APIConstants.OAuthConstants.OAUTH_CLIENT_SECRET)
                                        .toString())) {
                            String apiSecret = endpointSecurityProduction
                                    .get(APIConstants.OAuthConstants.OAUTH_CLIENT_SECRET).toString();
                            try {
                                String encryptedApiSecret = cryptoTool.encryptAndBase64Encode(apiSecret.getBytes());
                                endpointSecurityProduction
                                        .put(APIConstants.OAuthConstants.OAUTH_CLIENT_SECRET, encryptedApiSecret);
                            } catch (CryptoToolException e) {
                                throw new APIManagementException(ExceptionCodes
                                        .from(ExceptionCodes.ENDPOINT_CRYPTO_ERROR,
                                                "Error while encoding OAuth client secret"));
                            }
                        } else if (StringUtils.isNotBlank(oldProductionApiSecret)) {
                            endpointSecurityProduction
                                    .put(APIConstants.OAuthConstants.OAUTH_CLIENT_SECRET, oldProductionApiSecret);
                        } else {
                            String errorMessage = "Client secret is not provided for production endpoint security";
                            throw new APIManagementException(
                                    ExceptionCodes.from(ExceptionCodes.INVALID_ENDPOINT_CREDENTIALS, errorMessage));
                        }
                    }
                    endpointSecurity
                            .put(APIConstants.OAuthConstants.ENDPOINT_SECURITY_PRODUCTION, endpointSecurityProduction);
                    endpointConfig.put(APIConstants.ENDPOINT_SECURITY, endpointSecurity);
                    //apidto.setEndpointConfig(endpointConfig);
                }
                if (endpointSecurity.get(APIConstants.OAuthConstants.ENDPOINT_SECURITY_SANDBOX) != null) {
                    Map endpointSecuritySandbox = (Map) endpointSecurity
                            .get(APIConstants.OAuthConstants.ENDPOINT_SECURITY_SANDBOX);
                    String sandboxEndpointType = (String) endpointSecuritySandbox
                            .get(APIConstants.OAuthConstants.ENDPOINT_SECURITY_TYPE);

                    // Change default value of customParameters JSONObject to String
                    if (!(endpointSecuritySandbox
                            .get(APIConstants.OAuthConstants.OAUTH_CUSTOM_PARAMETERS) instanceof String)) {
                        Map<String, String> customParametersHashMap = (Map<String, String>) endpointSecuritySandbox
                                .get(APIConstants.OAuthConstants.OAUTH_CUSTOM_PARAMETERS);
                        customParametersString = JSONObject.toJSONString(customParametersHashMap);
                    } else if (endpointSecuritySandbox.get(APIConstants.OAuthConstants.OAUTH_CUSTOM_PARAMETERS)
                            != null) {
                        customParametersString = (String) endpointSecuritySandbox
                                .get(APIConstants.OAuthConstants.OAUTH_CUSTOM_PARAMETERS);
                    } else {
                        customParametersString = "{}";
                    }
                    endpointSecuritySandbox
                            .put(APIConstants.OAuthConstants.OAUTH_CUSTOM_PARAMETERS, customParametersString);

                    if (APIConstants.OAuthConstants.OAUTH.equals(sandboxEndpointType)) {
                        if (endpointSecuritySandbox.get(APIConstants.OAuthConstants.OAUTH_CLIENT_SECRET) != null
                                && StringUtils.isNotBlank(
                                endpointSecuritySandbox.get(APIConstants.OAuthConstants.OAUTH_CLIENT_SECRET)
                                        .toString())) {
                            String apiSecret = endpointSecuritySandbox
                                    .get(APIConstants.OAuthConstants.OAUTH_CLIENT_SECRET).toString();
                            try {
                                String encryptedApiSecret = cryptoTool.encryptAndBase64Encode(apiSecret.getBytes());
                                endpointSecuritySandbox
                                        .put(APIConstants.OAuthConstants.OAUTH_CLIENT_SECRET, encryptedApiSecret);
                            } catch (CryptoToolException e) {
                                throw new APIManagementException(ExceptionCodes
                                        .from(ExceptionCodes.ENDPOINT_CRYPTO_ERROR,
                                                "Error while encoding OAuth client secret"));
                            }
                        } else if (StringUtils.isNotBlank(oldSandboxApiSecret)) {
                            endpointSecuritySandbox
                                    .put(APIConstants.OAuthConstants.OAUTH_CLIENT_SECRET, oldSandboxApiSecret);
                        } else {
                            String errorMessage = "Client secret is not provided for sandbox endpoint security";
                            throw new APIManagementException(
                                    ExceptionCodes.from(ExceptionCodes.INVALID_ENDPOINT_CREDENTIALS, errorMessage));
                        }
                    }
                    endpointSecurity
                            .put(APIConstants.OAuthConstants.ENDPOINT_SECURITY_SANDBOX, endpointSecuritySandbox);
                    endpointConfig.put(APIConstants.ENDPOINT_SECURITY, endpointSecurity);
                    //apidto.setEndpointConfig(endpointConfig);
                }
            }
        }
    }

    /**
     * Check whether the token has APIDTO class level Scope annotation.
     *
     * @return true if the token has APIDTO class level Scope annotation
     */
    private static boolean checkClassScopeAnnotation(Scope[] apiDtoClassAnnotatedScopes, String[] tokenScopes) {

        for (Scope classAnnotation : apiDtoClassAnnotatedScopes) {
            for (String tokenScope : tokenScopes) {
                if (classAnnotation.name().equals(tokenScope)) {
                    return true;
                }
            }
        }
        return false;
    }

    /**
     * Override the API DTO field values with the user passed new values considering the field-wise scopes defined as
     * allowed to update in REST API definition yaml.
     */
    private static JSONObject overrideDTOValues(JSONObject originalApiDtoJson, JSONObject newApiDtoJson, Field field,
                                                String[] tokenScopes, Scope[] fieldAnnotatedScopes)
            throws APIManagementException {

        for (String tokenScope : tokenScopes) {
            for (Scope scopeAnt : fieldAnnotatedScopes) {
                if (scopeAnt.name().equals(tokenScope)) {
                    // do the overriding
                    originalApiDtoJson.put(field.getName(), newApiDtoJson.get(field.getName()));
                    return originalApiDtoJson;
                }
            }
        }
        throw new APIManagementException("User is not authorized to update one or more API fields. None of the "
                + "required scopes found in user token to update the field. So the request will be failed.",
                ExceptionCodes.INVALID_SCOPE);
    }

    /**
     * Get the API DTO object in which the API field values are overridden with the user passed new values.
     *
     * @throws APIManagementException
     */
    private static APIDTO getFieldOverriddenAPIDTO(APIDTO apidto, API originalAPI, String[] tokenScopes)
            throws APIManagementException {

        APIDTO originalApiDTO;
        APIDTO updatedAPIDTO;

        try {
            originalApiDTO = APIMappingUtil.fromAPItoDTO(originalAPI);

            Field[] fields = APIDTO.class.getDeclaredFields();
            ObjectMapper mapper = new ObjectMapper();
            String newApiDtoJsonString = mapper.writeValueAsString(apidto);
            JSONParser parser = new JSONParser();
            JSONObject newApiDtoJson = (JSONObject) parser.parse(newApiDtoJsonString);

            String originalApiDtoJsonString = mapper.writeValueAsString(originalApiDTO);
            JSONObject originalApiDtoJson = (JSONObject) parser.parse(originalApiDtoJsonString);

            for (Field field : fields) {
                Scope[] fieldAnnotatedScopes = field.getAnnotationsByType(Scope.class);
                String originalElementValue = mapper.writeValueAsString(originalApiDtoJson.get(field.getName()));
                String newElementValue = mapper.writeValueAsString(newApiDtoJson.get(field.getName()));

                if (!StringUtils.equals(originalElementValue, newElementValue)) {
                    originalApiDtoJson = overrideDTOValues(originalApiDtoJson, newApiDtoJson, field, tokenScopes,
                            fieldAnnotatedScopes);
                }
            }

            updatedAPIDTO = mapper.readValue(originalApiDtoJson.toJSONString(), APIDTO.class);

        } catch (IOException | ParseException e) {
            String msg = "Error while processing API DTO json strings";
            throw new APIManagementException(msg, e, ExceptionCodes.JSON_PARSE_ERROR);
        }
        return updatedAPIDTO;
    }

    /**
     * Update thumbnail of an API/API Product
     *
     * @param fileInputStream Input stream
     * @param fileContentType The content type of the image
     * @param apiProvider     API Provider
     * @param apiId           API/API Product UUID
     * @param tenantDomain    Tenant domain of the API
     * @throws APIManagementException If an error occurs while updating the thumbnail
     */
    public static void updateThumbnail(InputStream fileInputStream, String fileContentType, APIProvider apiProvider,
                                       String apiId, String tenantDomain) throws APIManagementException {
        ResourceFile apiImage = new ResourceFile(fileInputStream, fileContentType);
        apiProvider.setThumbnailToAPI(apiId, apiImage, tenantDomain);
    }

    /**
     * Add document DTO.
     *
     * @param documentDto Document DTO
     * @param apiId       API UUID
     * @return Added documentation
     * @param organization  Identifier of an Organization
     * @throws APIManagementException If an error occurs when retrieving API Identifier,
     *                                when checking whether the documentation exists and when adding the documentation
     */
    public static Documentation addDocumentationToAPI(DocumentDTO documentDto, String apiId, String organization)
            throws APIManagementException {

        APIProvider apiProvider = RestApiCommonUtil.getLoggedInUserProvider();
        Documentation documentation = DocumentationMappingUtil.fromDTOtoDocumentation(documentDto);
        String documentName = documentDto.getName();
        if (documentDto.getType() == null) {
            throw new APIManagementException("Documentation type cannot be empty",
                    ExceptionCodes.PARAMETER_NOT_PROVIDED);
        }
        if (documentDto.getType() == DocumentDTO.TypeEnum.OTHER && StringUtils
                .isBlank(documentDto.getOtherTypeName())) {
            //check otherTypeName for not null if doc type is OTHER
            throw new APIManagementException("otherTypeName cannot be empty if type is OTHER.",
                    ExceptionCodes.PARAMETER_NOT_PROVIDED);
        }
        String sourceUrl = documentDto.getSourceUrl();
        if (documentDto.getSourceType() == DocumentDTO.SourceTypeEnum.URL && (
                StringUtils.isBlank(sourceUrl) || !RestApiCommonUtil.isURL(sourceUrl))) {
            throw new APIManagementException("Invalid document sourceUrl Format",
                    ExceptionCodes.PARAMETER_NOT_PROVIDED);
        }

        if (apiProvider.isDocumentationExist(apiId, documentName, organization)) {
            throw new APIManagementException("Requested document '" + documentName + "' already exists",
                    ExceptionCodes.DOCUMENT_ALREADY_EXISTS);
        }
        documentation = apiProvider.addDocumentation(apiId, documentation, organization);

        return documentation;
    }

    /**
     * Add documentation content of inline and markdown documents.
     *
     * @param documentation Documentation
     * @param apiProvider   API Provider
     * @param apiId         API/API Product UUID
     * @param documentId    Document ID
     * @param organization  Identifier of the organization
     * @param inlineContent Inline content string
     * @throws APIManagementException If an error occurs while adding the documentation content
     */
    public static void addDocumentationContent(Documentation documentation, APIProvider apiProvider, String apiId,
                                               String documentId, String organization, String inlineContent)
            throws APIManagementException {
        DocumentationContent content = new DocumentationContent();
        content.setSourceType(DocumentationContent.ContentSourceType.valueOf(documentation.getSourceType().toString()));
        content.setTextContent(inlineContent);
        apiProvider.addDocumentationContent(apiId, documentId, organization, content);
    }

    /**
     * Add documentation content of files.
     *
     * @param inputStream  Input Stream
     * @param mediaType    Media type of the document
     * @param filename     File name
     * @param apiProvider  API Provider
     * @param apiId        API/API Product UUID
     * @param documentId   Document ID
     * @param organization organization of the API
     * @throws APIManagementException If an error occurs while adding the documentation file
     */
    public static void addDocumentationContentForFile(InputStream inputStream, String mediaType, String filename,
                                                      APIProvider apiProvider, String apiId,
                                                      String documentId, String organization)
            throws APIManagementException {
        DocumentationContent content = new DocumentationContent();
        ResourceFile resourceFile = new ResourceFile(inputStream, mediaType);
        resourceFile.setName(filename);
        content.setResourceFile(resourceFile);
        content.setSourceType(DocumentationContent.ContentSourceType.FILE);
        apiProvider.addDocumentationContent(apiId, documentId, organization, content);
    }

    /**
     * Checks whether the list of tiers are valid given the all valid tiers.
     *
     * @param allTiers     All defined tiers
     * @param currentTiers tiers to check if they are a subset of defined tiers
     * @return null if there are no invalid tiers or returns the set of invalid tiers if there are any
     */
    public static List<String> getInvalidTierNames(Set<Tier> allTiers, List<String> currentTiers) {

        List<String> invalidTiers = new ArrayList<>();
        for (String tierName : currentTiers) {
            boolean isTierValid = false;
            for (Tier definedTier : allTiers) {
                if (tierName.equals(definedTier.getName())) {
                    isTierValid = true;
                    break;
                }
            }
            if (!isTierValid) {
                invalidTiers.add(tierName);
            }
        }
        return invalidTiers;
    }

    /**
     * Change the lifecycle state of an API or API Product identified by UUID
     *
     * @param action       LC state change action
     * @param apiTypeWrapper API Type Wrapper (API or API Product)
     * @param lcChecklist  LC state change check list
     * @param organization Organization of logged-in user
     * @return APIStateChangeResponse
     * @throws APIManagementException Exception if there is an error when changing the LC state of API or API Product
     */
    public static APIStateChangeResponse changeApiOrApiProductLifecycle(String action, ApiTypeWrapper apiTypeWrapper,
                                                                        String lcChecklist, String organization)
            throws APIManagementException {

        String[] checkListItems = lcChecklist != null ? lcChecklist.split(APIConstants.DELEM_COMMA) : new String[0];
        APIProvider apiProvider = RestApiCommonUtil.getLoggedInUserProvider();

        Map<String, Object> apiLCData = apiProvider.getAPILifeCycleData(apiTypeWrapper.getUuid(), organization);

        String[] nextAllowedStates = (String[]) apiLCData.get(APIConstants.LC_NEXT_STATES);
        if (!ArrayUtils.contains(nextAllowedStates, action)) {
            throw new APIManagementException("Action '" + action + "' is not allowed. Allowed actions are "
                    + Arrays.toString(nextAllowedStates), ExceptionCodes.from(ExceptionCodes
                    .UNSUPPORTED_LIFECYCLE_ACTION, action));
        }

        //check and set lifecycle check list items including "Deprecate Old Versions" and "Require Re-Subscription".
        Map<String, Boolean> lcMap = new HashMap<>();
        for (String checkListItem : checkListItems) {
            String[] attributeValPair = checkListItem.split(APIConstants.DELEM_COLON);
            if (attributeValPair.length == 2) {
                String checkListItemName = attributeValPair[0].trim();
                boolean checkListItemValue = Boolean.parseBoolean(attributeValPair[1].trim());
                lcMap.put(checkListItemName, checkListItemValue);
            }
        }

        return apiProvider.changeLifeCycleStatus(organization, apiTypeWrapper, action, lcMap);
    }

    /**
     * Retrieve lifecycle history of API or API Product by Identifier
     *
     * @param uuid    Unique UUID of API or API Product
     * @return LifecycleHistoryDTO object
     * @throws APIManagementException exception if there is an error when retrieving the LC history
     */
    public static LifecycleHistoryDTO getLifecycleHistoryDTO(String uuid, APIProvider apiProvider)
            throws APIManagementException {

        List<LifeCycleEvent> lifeCycleEvents = apiProvider.getLifeCycleEvents(uuid);
        return APIMappingUtil.fromLifecycleHistoryModelToDTO(lifeCycleEvents);
    }

    /**
     * Get lifecycle state information of API or API Product
     *
     * @param identifier   Unique identifier of API or API Product
     * @param organization Organization of logged-in user
     * @return LifecycleStateDTO object
     * @throws APIManagementException if there is en error while retrieving the lifecycle state information
     */
    public static LifecycleStateDTO getLifecycleStateInformation(Identifier identifier, String organization)
            throws APIManagementException {

        APIProvider apiProvider = RestApiCommonUtil.getLoggedInUserProvider();
        Map<String, Object> apiLCData = apiProvider.getAPILifeCycleData(identifier.getUUID(), organization);
        if (apiLCData == null) {
            String type;
            if (identifier instanceof APIProductIdentifier) {
                type = APIConstants.API_PRODUCT;
            } else {
                type = APIConstants.API_IDENTIFIER_TYPE;
            }
            throw new APIManagementException("Error while getting lifecycle state for " + type + " with ID "
                    + identifier, ExceptionCodes.from(ExceptionCodes.LIFECYCLE_STATE_INFORMATION_NOT_FOUND, type,
                    identifier.getUUID()));
        } else {
            boolean apiOlderVersionExist = false;
            // check whether other versions of the current API exists
            APIVersionStringComparator comparator = new APIVersionStringComparator();
            Set<String> versions =
                    apiProvider.getAPIVersions(APIUtil.replaceEmailDomain(identifier.getProviderName()),
                            identifier.getName(), organization);

            for (String tempVersion : versions) {
                if (comparator.compare(tempVersion, identifier.getVersion()) < 0) {
                    apiOlderVersionExist = true;
                    break;
                }
            }
            return APIMappingUtil.fromLifecycleModelToDTO(apiLCData, apiOlderVersionExist);
        }
    }


    /**
     * Attaches a file to the specified document
     *
     * @param apiId         identifier of the API, the document belongs to
     * @param documentation Documentation object
     * @param inputStream   input Stream containing the file
     * @param fileName      File name
     * @param mediaType     Media type
     * @param organization  identifier of an organization
     * @throws APIManagementException if unable to add the file
     */
    public static void attachFileToDocument(String apiId, Documentation documentation, InputStream inputStream,
                                            String fileName, String mediaType, String organization)
            throws APIManagementException {

        APIProvider apiProvider = RestApiCommonUtil.getLoggedInUserProvider();
        String documentId = documentation.getId();
        String randomFolderName = RandomStringUtils.randomAlphanumeric(10);
        String tmpFolder = System.getProperty(RestApiConstants.JAVA_IO_TMPDIR) + File.separator
                + RestApiConstants.DOC_UPLOAD_TMPDIR + File.separator + randomFolderName;
        File docFile = new File(tmpFolder);

        boolean folderCreated = docFile.mkdirs();
        if (!folderCreated) {
            throw new APIManagementException("Failed to add content to the document " + documentId,
                    ExceptionCodes.INTERNAL_ERROR);
        }

        InputStream docInputStream = null;
        try {
            if (StringUtils.isBlank(fileName)) {
                fileName = RestApiConstants.DOC_NAME_DEFAULT + randomFolderName;
                log.warn(
                        "Couldn't find the name of the uploaded file for the document " + documentId + ". Using name '"
                                + fileName + "'");
            }
            //APIIdentifier apiIdentifier = APIMappingUtil
            //        .getAPIIdentifierFromUUID(apiId, tenantDomain);

            transferFile(inputStream, fileName, docFile.getAbsolutePath());
            docInputStream = new FileInputStream(docFile.getAbsolutePath() + File.separator + fileName);
            mediaType = mediaType == null ? RestApiConstants.APPLICATION_OCTET_STREAM : mediaType;
            PublisherCommonUtils
                    .addDocumentationContentForFile(docInputStream, mediaType, fileName, apiProvider, apiId,
                            documentId, organization);
            docFile.deleteOnExit();
        } catch (FileNotFoundException e) {
            throw new APIManagementException("Unable to read the file from path ", e, ExceptionCodes.INTERNAL_ERROR);
        } finally {
            IOUtils.closeQuietly(docInputStream);
        }
    }

    /**
     * This method uploads a given file to specified location
     *
     * @param uploadedInputStream input stream of the file
     * @param newFileName         name of the file to be created
     * @param storageLocation     destination of the new file
     * @throws APIManagementException if the file transfer fails
     */
    public static void transferFile(InputStream uploadedInputStream, String newFileName, String storageLocation)
            throws APIManagementException {
        FileOutputStream outFileStream = null;

        try {
            outFileStream = new FileOutputStream(new File(storageLocation, newFileName));
            int read;
            byte[] bytes = new byte[1024];
            while ((read = uploadedInputStream.read(bytes)) != -1) {
                outFileStream.write(bytes, 0, read);
            }
        } catch (IOException e) {
            String errorMessage = "Error in transferring files.";
            log.error(errorMessage, e);
            throw new APIManagementException(errorMessage, e, ExceptionCodes.INTERNAL_ERROR);
        } finally {
            IOUtils.closeQuietly(outFileStream);
        }
    }

    /**
     * This method validates monetization properties
     *
     * @param monetizationProperties map of monetization properties
     * @throws APIManagementException
     */
    public static void validateMonetizationProperties(Map<String, String> monetizationProperties)
            throws APIManagementException {

        String errorMessage;
        if (monetizationProperties != null) {
            for (Map.Entry<String, String> entry : monetizationProperties.entrySet()) {
                String monetizationPropertyKey = entry.getKey().trim();
                String propertyValue = entry.getValue();
                if (monetizationPropertyKey.contains(" ")) {
                    errorMessage = "Monetization property names should not contain space character. " +
                            "Monetization property '" + monetizationPropertyKey + "' "
                            + "contains space in it.";
                    throw new APIManagementException(errorMessage, ExceptionCodes.INVALID_PARAMETERS_PROVIDED);
                }
                // Maximum allowable characters of registry property name and value is 100 and 1000.
                // Hence we are restricting them to be within 80 and 900.
                if (monetizationPropertyKey.length() > 80) {
                    errorMessage = "Monetization property name can have maximum of 80 characters. " +
                            "Monetization property '" + monetizationPropertyKey + "' + contains "
                            + monetizationPropertyKey.length() + "characters";
                    throw new APIManagementException(errorMessage, ExceptionCodes.INVALID_PARAMETERS_PROVIDED);
                }
                if (propertyValue.length() > 900) {
                    errorMessage = "Monetization property value can have maximum of 900 characters. " +
                            "Property '" + monetizationPropertyKey + "' + "
                            + "contains a value with " + propertyValue.length() + "characters";
                    throw new APIManagementException(errorMessage, ExceptionCodes.INVALID_PARAMETERS_PROVIDED);
                }
            }
        }
    }

    /**
     * This method is used to read input stream of a file and return the string content.
     * @param fileInputStream File input stream
     * @return String
     * @throws APIManagementException*/
    public static String readInputStream(InputStream fileInputStream)
            throws APIManagementException {

        String content = null;
        if (fileInputStream != null) {
            try {
                ByteArrayOutputStream outputStream = new ByteArrayOutputStream();
                IOUtils.copy(fileInputStream, outputStream);
                byte[] sequenceBytes = outputStream.toByteArray();
                InputStream inSequenceStream = new ByteArrayInputStream(sequenceBytes);
                content = IOUtils.toString(inSequenceStream, StandardCharsets.UTF_8.name());
            } catch (IOException e) {
                throw new APIManagementException("Error occurred while reading inputs", e,
                        ExceptionCodes.INTERNAL_ERROR);
            }

        }
        return content;
    }

    public static File exportOperationPolicyData(OperationPolicyData policyData, String format)
            throws APIManagementException {

        File exportFolder = null;
        try {
            exportFolder = CommonUtil.createTempDirectoryFromName(policyData.getSpecification().getName()
                    + "_" + policyData.getSpecification().getVersion());
            String exportAPIBasePath = exportFolder.toString();
            String archivePath =
                    exportAPIBasePath.concat(File.separator + policyData.getSpecification().getName());
            CommonUtil.createDirectory(archivePath);
            String policyName = archivePath + File.separator + policyData.getSpecification().getName();
            if (policyData.getSpecification() != null) {
                if (format.equalsIgnoreCase(ExportFormat.YAML.name())) {
                    CommonUtil.writeDtoToFile(policyName, ExportFormat.YAML,
                            ImportExportConstants.TYPE_POLICY_SPECIFICATION,
                            policyData.getSpecification());
                } else if (format.equalsIgnoreCase(ExportFormat.JSON.name())) {
                    CommonUtil.writeDtoToFile(policyName, ExportFormat.JSON,
                            ImportExportConstants.TYPE_POLICY_SPECIFICATION,
                            policyData.getSpecification());
                }
            }
            if (policyData.getSynapsePolicyDefinition() != null) {
                CommonUtil.writeFile(policyName + APIConstants.SYNAPSE_POLICY_DEFINITION_EXTENSION,
                        policyData.getSynapsePolicyDefinition().getContent());
            }
            if (policyData.getCcPolicyDefinition() != null) {
                CommonUtil.writeFile(policyName + APIConstants.CC_POLICY_DEFINITION_EXTENSION,
                        policyData.getCcPolicyDefinition().getContent());
            }

            CommonUtil.archiveDirectory(exportAPIBasePath);
            FileUtils.deleteQuietly(new File(exportAPIBasePath));
            return new File(exportAPIBasePath + APIConstants.ZIP_FILE_EXTENSION);
        } catch (APIImportExportException | IOException e) {
            throw new APIManagementException("Error while exporting operation policy", e,
                    ExceptionCodes.INTERNAL_ERROR);
        }
    }
}
