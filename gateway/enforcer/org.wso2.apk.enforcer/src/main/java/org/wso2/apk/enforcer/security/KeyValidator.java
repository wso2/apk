/*
 * Copyright (c) 2020, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.wso2.apk.enforcer.security;

import com.nimbusds.jwt.JWTClaimsSet;
import net.minidev.json.JSONObject;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.wso2.apk.enforcer.commons.exception.APISecurityException;
import org.wso2.apk.enforcer.commons.logging.ErrorDetails;
import org.wso2.apk.enforcer.commons.logging.LoggingConstants;
import org.wso2.apk.enforcer.commons.model.APIConfig;
import org.wso2.apk.enforcer.commons.model.ResourceConfig;
import org.wso2.apk.enforcer.constants.APIConstants;
import org.wso2.apk.enforcer.constants.APISecurityConstants;
import org.wso2.apk.enforcer.dto.APIKeyValidationInfoDTO;
import org.wso2.apk.enforcer.models.Application;
import org.wso2.apk.enforcer.models.ApplicationKeyMapping;
import org.wso2.apk.enforcer.models.ApplicationMapping;
import org.wso2.apk.enforcer.models.Subscription;
import org.wso2.apk.enforcer.subscription.SubscriptionDataHolder;
import org.wso2.apk.enforcer.subscription.SubscriptionDataStore;

import java.util.List;
import java.util.Set;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * Does the subscription and scope validation.
 */
public class KeyValidator {

    private static final Logger log = LogManager.getLogger(KeyValidator.class);

    /**
     * Validate the scopes related to the given validationContext.
     *
     * @param validationContext the token validation context
     * @return true is the scopes are valid
     * this will indicate the message body for the error response
     */
    public static boolean validateScopes(TokenValidationContext validationContext) throws APISecurityException {

        if (validationContext.isCacheHit()) {
            return true;
        }
        APIKeyValidationInfoDTO apiKeyValidationInfoDTO = validationContext.getValidationInfoDTO();

        if (apiKeyValidationInfoDTO == null) {
            log.error("Error while validating scopes. Key validation information has not been set.",
                    ErrorDetails.errorLog(LoggingConstants.Severity.MINOR, 6603));
            throw new APISecurityException(APIConstants.StatusCodes.UNAUTHENTICATED.getCode(),
                    APISecurityConstants.API_AUTH_GENERAL_ERROR,
                    "Error while validating scopes. Key validation information has not been set");
        }

        Set<String> scopesFromToken = apiKeyValidationInfoDTO.getScopes();

        List<ResourceConfig> matchedResources;
        // when it is a graphQL api multiple matching resources will be returned.
        matchedResources = validationContext.getMatchingResourceConfigs();

        boolean allScopesValidated = true;
        // failedResourcePath - used to identify resource paths with failed scope validation.
        String failedResourcePath = "";
        for (ResourceConfig matchedResource : matchedResources) {
            // scopesValidated - indicate scope has validated
            boolean scopesValidated = false;
            String resourcePath = matchedResource.getPath();
            String[] scopesToValidate = matchedResource.getScopes();
            for (String scope : scopesToValidate) {
                if (scopesFromToken.contains(scope)) {
                    scopesValidated = true;
                    break;
                }
            }
            if (scopesToValidate.length > 0 && !scopesValidated) {
                allScopesValidated = false;
                failedResourcePath = resourcePath;
                break;
            }
        }
        if (!allScopesValidated) {
            apiKeyValidationInfoDTO.setAuthorized(false);
            apiKeyValidationInfoDTO.setValidationStatus(APIConstants.KeyValidationStatus.INVALID_SCOPE);
            String message = "User is NOT authorized to access the Resource: " + failedResourcePath
                    + ". Scope validation failed.";
            throw new APISecurityException(APIConstants.StatusCodes.UNAUTHORIZED.getCode(),
                    APISecurityConstants.INVALID_SCOPE, message);
        }
        return true;
    }

    /**
     * Validate subscriptions for access tokens by utilizing the consumer key.
     *
     * @param validationInfo Token validation related details. This will be populated based on the available data during
     *                       the subscription validation.
     * @throws APISecurityException throws if subscription validation fails.
     */
    public static void validateSubscriptionUsingConsumerKey(APIKeyValidationInfoDTO validationInfo)
            throws APISecurityException {

        Application app;
        Subscription sub = null;
        ApplicationKeyMapping keyMapping;
        Set<ApplicationMapping> appMappings;
        String apiName = validationInfo.getApiName();
        String apiContext = validationInfo.getApiContext();
        String apiVersion = validationInfo.getApiVersion();
        String consumerKey = validationInfo.getConsumerKey();
        String securityScheme = validationInfo.getSecurityScheme();
        String keyType = validationInfo.getEnvType();

        log.debug("Before validating subscriptions");
        log.debug("Validation Info : { name : {}, context : {}, version : {}, consumerKey : {} }",
                apiName, apiContext, apiVersion, consumerKey);

        SubscriptionDataStore datastore =
                SubscriptionDataHolder.getInstance().getSubscriptionDataStore(validationInfo.getSubscriberOrganization());

        if (datastore != null) {
            // Get application key mapping using the consumer key, key type and security scheme
            keyMapping = datastore.getMatchingApplicationKeyMapping(consumerKey, keyType, securityScheme,
                    validationInfo.getEnvironment());

            if (keyMapping != null) {
                // Get application and application mapping using application UUID
                String applicationUUID = keyMapping.getApplicationUUID();
                app = datastore.getApplicationById(applicationUUID);
                appMappings = datastore.getMatchingApplicationMappings(applicationUUID);

                if (appMappings != null && app != null) {
                    // Get subscription using the subscription UUID
                    for (ApplicationMapping appMapping : appMappings) {
                        String subscriptionUUID = appMapping.getSubscriptionUUID();
                        Subscription subscription = datastore.getMatchingSubscription(subscriptionUUID);

                        if (validationInfo.getApiName().equals(subscription.getSubscribedApi().getName())) {
                            // Validate API version
                            Pattern pattern = subscription.getSubscribedApi().getVersionRegexPattern();
                            String versionToMatch = validationInfo.getApiVersion();
                            Matcher matcher = pattern.matcher(versionToMatch);
                            if (matcher.matches()) {
                                sub = subscription;
                                break;
                            }
                        }
                    }
                    // Validate subscription
                    if (sub != null) {
                        validate(validationInfo, app, sub);
                        if (!validationInfo.isAuthorized() && validationInfo.getValidationStatus() == 0) {
                            // Scenario where validation failed and message is not set
                            validationInfo.setValidationStatus(
                                    APIConstants.KeyValidationStatus.API_AUTH_RESOURCE_FORBIDDEN);
                        }
                        log.debug("After validating subscriptions");
                        return;
                    } else {
                        log.error(
                                "Valid subscription not found for access token. " +
                                        "application: {}, app_UUID: {}, API name: {}, API context: {} API version" +
                                        " : {}",
                                app.getName(), app.getUUID(), apiName, apiContext, apiVersion);
                    }
                } else {
                    log.error(
                            "Valid application and / or application mapping not found for application uuid : " + applicationUUID);
                }
            } else {
                log.error(
                        "Valid application key mapping not found in the data store for access token. " +
                                "Application identifier: {}, key type : {}, security scheme : {}",
                        consumerKey, keyType, securityScheme);
            }
        } else {
            log.error("Subscription data store is null");
            throw new APISecurityException(APIConstants.StatusCodes.UNAUTHENTICATED.getCode(),
                    APISecurityConstants.API_AUTH_GENERAL_ERROR,
                    APISecurityConstants.API_AUTH_GENERAL_ERROR_MESSAGE);
        }
        // If the execution reaches this point, it means that the subscription validation has failed.
        log.error("User is NOT authorized to access the API. Subscription validation failed for consumer key : "
                + consumerKey);
        throw new

                APISecurityException(APIConstants.StatusCodes.UNAUTHORIZED.getCode(),

                APISecurityConstants.API_AUTH_FORBIDDEN, APISecurityConstants.API_AUTH_FORBIDDEN_MESSAGE);
    }

    /**
     * Validate subscriptions for API keys.
     *
     * @param apiUuid    uuid of the API
     * @param apiContext API context, used for logging purposes and to extract the tenant domain
     * @param payload    JWT claims set extracted from the API key
     * @return validation information about the request
     */
    public static APIKeyValidationInfoDTO validateSubscription(String apiUuid, String apiContext, APIConfig api,
                                                               JWTClaimsSet payload) {

        log.debug("Before validating subscriptions with API key. API_uuid: {}, context: {}", apiUuid, apiContext);

        Application app = null;
        Subscription sub = null;

        SubscriptionDataStore datastore =
                SubscriptionDataHolder.getInstance().getSubscriptionDataStore(api.getOrganizationId());
        if (datastore != null) {
            JSONObject appObject = (JSONObject) payload.getClaim(APIConstants.JwtTokenConstants.APPLICATION);
            String appUuid = appObject.getAsString("uuid");
            if (!appObject.isEmpty() && !appUuid.isEmpty()) {
                app = datastore.getApplicationById(appUuid);
                if (app != null) {
                    sub = datastore.getSubscriptionById(app.getUUID(), api.getUuid());
                    if (sub != null) {
                        log.debug("All information is retrieved from the in memory data store.");
                    } else {
                        log.info(
                                "Valid subscription not found for API key. " +
                                        "application: {} app_UUID: {} API_name: {} API_UUID : {}",
                                app.getName(), app.getUUID(), api.getName(), api.getUuid());
                    }
                } else {
                    log.info("Application not found in the data store for uuid {}", appUuid);
                }
            } else {
                log.info("Application claim not found in jwt for uuid");
            }
        } else {
            log.error("Subscription data store is null");
        }

        APIKeyValidationInfoDTO infoDTO = new APIKeyValidationInfoDTO();
        if (app != null && sub != null) {
            validate(infoDTO, app, sub);
        }
        if (!infoDTO.isAuthorized() && infoDTO.getValidationStatus() == 0) {
            //Scenario where validation failed and message is not set
            infoDTO.setValidationStatus(APIConstants.KeyValidationStatus.API_AUTH_RESOURCE_FORBIDDEN);
        }
        log.debug("After validating subscriptions with API key.");
        return infoDTO;
    }

    private static void validate(APIKeyValidationInfoDTO infoDTO, Application app, Subscription sub) {

        // Validate subscription status
        String subscriptionStatus = sub.getSubscriptionStatus();
        if (APIConstants.SubscriptionStatus.INACTIVE.equals(subscriptionStatus)) {
            infoDTO.setValidationStatus(APIConstants.KeyValidationStatus.SUBSCRIPTION_INACTIVE);
            infoDTO.setAuthorized(false);
            return;
        }
        infoDTO.setApplicationUUID(app.getUUID());
        infoDTO.setSubscriber(app.getOwner());
        infoDTO.setApplicationName(app.getName());
        infoDTO.setAppAttributes(app.getAttributes());
        infoDTO.setAuthorized(true);
    }
}
