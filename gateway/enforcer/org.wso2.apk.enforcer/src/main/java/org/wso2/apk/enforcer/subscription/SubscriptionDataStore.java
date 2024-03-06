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

import org.wso2.apk.enforcer.discovery.subscription.JWTIssuer;
import org.wso2.apk.enforcer.models.*;
import org.wso2.apk.enforcer.security.jwt.validator.JWTValidator;

import java.util.List;
import java.util.Set;

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
     * Gets Subscription by ID.
     *
     * @param appUUID Application associated with the Subscription (uuid)
     * @param apiUUID Api associated with the Subscription (uuid)
     * @return {@link Subscription}
     */
    Subscription getSubscriptionById(String appUUID, String apiUUID);

    void addSubscriptions(List<SubscriptionDto> subscriptionList);

    void addApplications(List<ApplicationDto> applicationList);

    void addApplicationKeyMappings(
            List<ApplicationKeyMappingDTO> applicationKeyMappingList);

    /**
     * Filter the applicationMapping map based on the provided application UUID.
     *
     * @param uuid Application UUID
     * @return ApplicationMapping which match the given UUID
     */
    Set<ApplicationMapping> getMatchingApplicationMappings(String uuid);

    /**
     * Filter the application key mapping map based on provided parameters
     *
     * @param applicationIdentifier Application identifier
     * @param keyType               Key type, i.e. PRODUCTION or SANDBOX
     * @param securityScheme        Security scheme
     * @param envType
     * @return ApplicationKeyMapping which match the given parameters
     */
    ApplicationKeyMapping getMatchingApplicationKeyMapping(String applicationIdentifier, String keyType,
            String securityScheme, String envType);

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
     *
     * @param issuer      issuer in JWT
     * @param environment environment of the Issuer
     * @return JWTValidator Implementation
     */
    JWTValidator getJWTValidatorByIssuer(String issuer, String environment);

    void addApplication(org.wso2.apk.enforcer.discovery.subscription.Application application);

    void addSubscription(org.wso2.apk.enforcer.discovery.subscription.Subscription subscription);

    void addApplicationMapping(org.wso2.apk.enforcer.discovery.subscription.ApplicationMapping applicationMapping);

    void addApplicationKeyMapping(
            org.wso2.apk.enforcer.discovery.subscription.ApplicationKeyMapping applicationKeyMapping);

    void removeApplicationMapping(org.wso2.apk.enforcer.discovery.subscription.ApplicationMapping applicationMapping);

    void removeApplicationKeyMapping(
            org.wso2.apk.enforcer.discovery.subscription.ApplicationKeyMapping applicationKeyMapping);

    void removeSubscription(org.wso2.apk.enforcer.discovery.subscription.Subscription subscription);

    void removeApplication(org.wso2.apk.enforcer.discovery.subscription.Application application);

    public void addApplicationMappings(List<ApplicationMappingDto> applicationMappingList);

     int getSubscriptionCount();

     int getJWTIssuerCount();
}
