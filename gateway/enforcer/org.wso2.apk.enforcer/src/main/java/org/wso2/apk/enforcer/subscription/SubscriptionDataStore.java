/*
 * Copyright (c) 2020, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 *
 * WSO2 Inc. licenses this file to you under the Apache License,
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

    void addSubscriptions(List<org.wso2.apk.enforcer.discovery.subscription.Subscription> subscriptionList);

    void addApplications(List<org.wso2.apk.enforcer.discovery.subscription.Application> applicationList);

    void addApis(List<APIs> apisList);

    void addApplicationPolicies(
            List<org.wso2.apk.enforcer.discovery.subscription.ApplicationPolicy> applicationPolicyList);

    void addSubscriptionPolicies(
            List<org.wso2.apk.enforcer.discovery.subscription.SubscriptionPolicy> subscriptionPolicyList);

    void addApplicationKeyMappings(
            List<org.wso2.apk.enforcer.discovery.subscription.ApplicationKeyMapping> applicationKeyMappingList);

    /**
     * Filter the API map according to the provided parameters
     * @param name API Name
     * @param context API Context
     * @param version API Version
     * @param uuid API UUID
     * @return Matching list of apis.
     */
    List<API> getMatchingAPIs(String name, String context, String version, String uuid);

    /**
     * Filter the API map according to the provided parameters
     *
     * @param context API Context
     * @param version API Version
     * @return Matching list of apis.
     */
    API getMatchingAPI(String context, String version);

    /**
     * Filter the applications map based on the criteria.
     * @param name Application Name
     * @param organizationID Application tenant domain/ organization id
     * @param uuid Application uuid.
     * @return List of applications which match the given parameters
     */
    List<Application> getMatchingApplications(String name, String organizationID, String uuid);


    /**
     * Filter the application key mapping map based on provided parameters
     * @param applicationUUID  Application uuid
     * @param consumerKey The application consumer key
     * @return List of key mappings which match the given parameters
     */
    List<ApplicationKeyMapping> getMatchingKeyMapping(String applicationUUID, String consumerKey);


    /**
     * Filter the subscriptions map based on the provided parameters
     * @param applicationUUID UUID of the application
     * @param apiUUID UUID of the api
     * @param state Subscription state
     * @return A List of subscriptions which matches the given parameters
     */
    List<Subscription> getMatchingSubscriptions(String applicationUUID, String apiUUID, String state);

    void addJWTIssuers(List<JWTIssuer> jwtIssuers);

    /**
     * Returns the JWTValidator based on Issuer
     * @param issuer issuer in JWT
     * @return JWTValidator Implementation
     */
    JWTValidator getJWTValidatorByIssuer(String issuer);
}
