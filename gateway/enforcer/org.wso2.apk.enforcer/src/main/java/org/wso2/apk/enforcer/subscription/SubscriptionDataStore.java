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

package org.wso2.apk.enforcer.subscription;

import org.wso2.apk.enforcer.discovery.subscription.APIs;
import org.wso2.apk.enforcer.discovery.subscription.JWTIssuer;
import org.wso2.apk.enforcer.models.*;
import org.wso2.apk.enforcer.security.jwt.validator.JWTValidator;

import java.util.List;

/**
 * A Facade for obtaining Subscription related Data.
 */
public interface SubscriptionDataStore {

    /**
     * Gets an {@link Application} by appUUID.
     *
     * @param appUUID of the Application
     * @return {@link Application} with the appId
     */
    Application getApplicationById(String appUUID);

    /**
     * Get API by Context and Version.
     *
     * @param uuid UUID of the API
     * @return {@link API} entry represented by Context and Version.
     */
    API getApiByContextAndVersion(String uuid);

    /**
     * Gets Subscription by ID.
     *
     * @param appUUID Application associated with the Subscription (uuid)
     * @param apiUUID Api associated with the Subscription (uuid)
     * @return {@link Subscription}
     */
    Subscription getSubscriptionById(String appUUID, String apiUUID);

    void addSubscriptions(List<SubscriptionDto> subscriptionList);

    void addApplications(List<ApplicationDto> applicationList);

    void addApis(List<APIs> apisList);

    void addApplicationKeyMappings(
            List<ApplicationKeyMappingDTO> applicationKeyMappingList);

    /**
     * Filter the API map according to the provided parameters
     *
     * @param context API Context
     * @param version API Version
     * @return Matching list of apis.
     */
    API getMatchingAPI(String context, String version);

    /**
     * Filter the applicationMapping map based on the provided application UUID.
     *
     * @param uuid Application UUID
     * @return ApplicationMapping which match the given UUID
     */
    ApplicationMapping getMatchingApplicationMapping(String uuid);

    /**
     * Filter the application key mapping map based on provided parameters
     *
     * @param applicationIdentifier Application identifier
     * @param keyType               Key type, i.e. PRODUCTION or SANDBOX
     * @param securityScheme        Security scheme
     * @return ApplicationKeyMapping which match the given parameters
     */
    ApplicationKeyMapping getMatchingApplicationKeyMapping(String applicationIdentifier, String keyType,
            String securityScheme);

    /**
     * Filter the applications map based on the provided parameters.
     *
     * @param uuid UUID of the application
     * @return Application which match the given UUID
     */
    Application getMatchingApplication(String uuid);

    /**
     * Filter the subscriptions map based on the provided parameters.
     *
     * @param uuid UUID of the subscription
     * @return Subscription which matches the given UUID
     */
    Subscription getMatchingSubscription(String uuid);

    void addJWTIssuers(List<JWTIssuer> jwtIssuers);

    /**
     * Returns the JWTValidator based on Issuer
     * @param issuer issuer in JWT
     * @param environment environment of the Issuer
     * @return JWTValidator Implementation
     */
    JWTValidator getJWTValidatorByIssuer(String issuer, String organization, String environment);
}
