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

package org.wso2.apk.apimgt.impl.dao.constants;

import org.wso2.apk.apimgt.api.model.policy.PolicyConstants;
import org.wso2.apk.apimgt.impl.APIConstants;

public class SQLConstants {
    public static final String GET_API_FOR_CONTEXT_TEMPLATE_SQL =
            " SELECT " +
            "   API.API_NAME," +
            "   API.API_PROVIDER" +
            " FROM " +
            "   API API" +
            " WHERE " +
            "   API.CONTEXT_TEMPLATE = ? ";

    public static final String GET_VERSIONS_MATCHES_API_NAME_AND_ORGANIZATION_SQL=
            "SELECT API_VERSION FROM API WHERE API_NAME = ? AND API_PROVIDER = ? AND ORGANIZATION = ?";

    public static final String GET_APPLICATION_REGISTRATION_SQL =
            " SELECT REG_ID FROM APPLICATION_REGISTRATION WHERE SUBSCRIBER_ID = ? AND APP_ID = ? AND TOKEN_TYPE = " +
                    "? AND KEY_MANAGER = ?";

    public static final String ADD_APPLICATION_REGISTRATION_SQL =
            " INSERT INTO " +
            "   APPLICATION_REGISTRATION (SUBSCRIBER_ID,WF_REF,APP_ID,TOKEN_TYPE,ALLOWED_DOMAINS," +
            "VALIDITY_PERIOD,TOKEN_SCOPE,INPUTS,KEY_MANAGER) " +
            " VALUES(?,?,?,?,?,?,?,?,?)";

    public static final String ADD_APPLICATION_KEY_MAPPING_SQL =
            " INSERT INTO " +
            "   APPLICATION_KEY_MAPPING (APPLICATION_ID,KEY_TYPE,STATE,KEY_MANAGER,UUID,CREATE_MODE) " +
            " VALUES(?,?,?,?,?,?)";

    public static final String GET_SUBSCRIBED_APIS_OF_USER_SQL =
            " SELECT " +
            "   API.API_PROVIDER AS API_PROVIDER," +
            "   API.API_NAME AS API_NAME," +
            "   API.CONTEXT AS API_CONTEXT, " +
            "   API.API_VERSION AS API_VERSION, " +
            "   SP.TIER_ID AS SP_TIER_ID " +
            " FROM " +
            "   SUBSCRIPTION SP, " +
            "   API API," +
            "   SUBSCRIBER SB, " +
            "   APPLICATION APP " +
            " WHERE " +
            "   SB.USER_ID = ? " +
            "   AND SB.ORGANIZATION = ? " +
            "   AND SB.SUBSCRIBER_ID = APP.SUBSCRIBER_ID " +
            "   AND APP.APPLICATION_ID=SP.APPLICATION_ID " +
            "   AND API.API_ID = SP.API_ID" +
            "   AND SP.SUBS_CREATE_STATE = '" + APIConstants.SubscriptionCreatedStatus.SUBSCRIBE + "'";

    public static final String GET_SUBSCRIBED_APIS_OF_USER_CASE_INSENSITIVE_SQL =
            " SELECT " +
            "   API.API_PROVIDER AS API_PROVIDER," +
            "   API.API_NAME AS API_NAME," +
            "   API.CONTEXT AS API_CONTEXT," +
            "   API.API_VERSION AS API_VERSION, " +
            "   SP.TIER_ID AS SP_TIER_ID " +
            " FROM " +
            "   SUBSCRIPTION SP, " +
            "   API API," +
            "   SUBSCRIBER SB, " +
            "   APPLICATION APP " +
            " WHERE " +
            "   LOWER(SB.USER_ID) = LOWER(?) " +
            "   AND SB.ORGANIZATION = ? " +
            "   AND SB.SUBSCRIBER_ID = APP.SUBSCRIBER_ID " +
            "   AND APP.APPLICATION_ID=SP.APPLICATION_ID " +
            "   AND API.API_ID = SP.API_ID" +
            "   AND SP.SUBS_CREATE_STATE = '" +APIConstants.SubscriptionCreatedStatus.SUBSCRIBE + "'";

    public static final String GET_SUBSCRIBED_API_IDs_BY_APP_ID_SQL =
            " SELECT " +
                    "   API.API_ID " +
                    " FROM " +
                    "   SUBSCRIBER SUB," +
                    "   APPLICATION APP, " +
                    "   SUBSCRIPTION SUBS, " +
                    "   API API " +
                    " WHERE " +
                    "   SUB.ORGANIZATION = ? " +
                    "   AND SUB.SUBSCRIBER_ID=APP.SUBSCRIBER_ID " +
                    "   AND APP.APPLICATION_ID=SUBS.APPLICATION_ID " +
                    "   AND API.API_ID=SUBS.API_ID" +
                    "   AND APP.APPLICATION_ID= ? " +
                    "   AND SUBS.SUBS_CREATE_STATE = '" + APIConstants.SubscriptionCreatedStatus.SUBSCRIBE + "'";

    public static final String GET_INCLUDED_APIS_IN_PRODUCT_SQL =
            "SELECT "
                    + "DISTINCT API_ID "
                    + "FROM API_URL_MAPPING "
                    + "WHERE URL_MAPPING_ID IN "
                    + "(SELECT URL_MAPPING_ID FROM API_PRODUCT_MAPPING WHERE API_ID = ? AND REVISION_UUID = 'Current API')";

    public static final String GET_SUBSCRIBED_APIS_OF_USER_BY_APP_SQL =
            " SELECT " +
                    "   API.API_PROVIDER AS API_PROVIDER," +
                    "   API.API_NAME AS API_NAME," +
                    "   API.CONTEXT AS API_CONTEXT, " +
                    "   API.API_VERSION AS API_VERSION, " +
                    "   SP.TIER_ID AS SP_TIER_ID " +
                    " FROM " +
                    "   SUBSCRIPTION SP, " +
                    "   API API," +
                    "   SUBSCRIBER SB, " +
                    "   APPLICATION APP " +
                    " WHERE " +
                    "   SB.ORGANIZATION = ? " +
                    "   AND SB.SUBSCRIBER_ID = APP.SUBSCRIBER_ID " +
                    "   AND APP.APPLICATION_ID=SP.APPLICATION_ID " +
                    "   AND API.API_ID = SP.API_ID" +
                    "   AND (SP.SUB_STATUS = '" + APIConstants.SubscriptionStatus.UNBLOCKED +
                    "' OR SP.SUB_STATUS = '" + APIConstants.SubscriptionStatus.TIER_UPDATE_PENDING +
                    "' OR SP.SUB_STATUS = '" + APIConstants.SubscriptionStatus.PROD_ONLY_BLOCKED + "')" +
                    "   AND SP.SUBS_CREATE_STATE = '" + APIConstants.SubscriptionCreatedStatus.SUBSCRIBE + "'" +
                    "   AND APP.APPLICATION_ID = ?";

    public static final String GET_SUBSCRIBED_APIS_OF_USER_BY_APP_CASE_INSENSITIVE_SQL =
            " SELECT " +
                    "   API.API_PROVIDER AS API_PROVIDER," +
                    "   API.API_NAME AS API_NAME," +
                    "   API.CONTEXT AS API_CONTEXT," +
                    "   API.API_VERSION AS API_VERSION, " +
                    "   SP.TIER_ID AS SP_TIER_ID " +
                    " FROM " +
                    "   SUBSCRIPTION SP, " +
                    "   API API," +
                    "   SUBSCRIBER SB, " +
                    "   APPLICATION APP " +
                    " WHERE " +
                    "   SB.ORGANIZATION = ? " +
                    "   AND SB.SUBSCRIBER_ID = APP.SUBSCRIBER_ID " +
                    "   AND APP.APPLICATION_ID=SP.APPLICATION_ID " +
                    "   AND API.API_ID = SP.API_ID" +
                    "   AND (SP.SUB_STATUS = '" + APIConstants.SubscriptionStatus.UNBLOCKED +
                    "' OR SP.SUB_STATUS = '" + APIConstants.SubscriptionStatus.TIER_UPDATE_PENDING +
                    "' OR SP.SUB_STATUS = '" + APIConstants.SubscriptionStatus.PROD_ONLY_BLOCKED + "')" +
                    "   AND SP.SUBS_CREATE_STATE = '" + APIConstants.SubscriptionCreatedStatus.SUBSCRIBE + "'" +
                    "   AND APP.APPLICATION_ID = ?";

    public static final String ADD_API_SERVICE_MAPPING_SQL = "INSERT INTO API_SERVICE_MAPPING (API_ID, SERVICE_KEY" +
            ", MD5, TENANT_ID) VALUES (?,?,?,?)";

    public static final String GET_SERVICE_KEY_BY_API_ID_SQL = "SELECT SERVICE_KEY FROM API_SERVICE_MAPPING WHERE " +
            " API_ID = ? AND TENANT_ID = ?";

    public static final String GET_SERVICE_KEY_BY_API_ID_SQL_WITHOUT_TENANT_ID = "SELECT SERVICE_KEY FROM " +
            " API_SERVICE_MAPPING WHERE API_ID = ?";

    public static final String UPDATE_API_SERVICE_MAPPING_SQL = "UPDATE API_SERVICE_MAPPING SET " +
            "   SERVICE_KEY = ?, " +
            "   MD5 = ? " +
            "   WHERE API_ID = ?";

    public static final String GET_MD5_VALUE_OF_SERVICE_BY_API_ID_SQL = "SELECT " +
            "   SERVICE_CATALOG.MD5 AS SERVICE_MD5, " +
            "   SERVICE_CATALOG.SERVICE_NAME, " +
            "   SERVICE_CATALOG.SERVICE_VERSION, " +
            "   API_SERVICE_MAPPING.MD5 AS API_SERVICE_MD5, " +
            "   API_SERVICE_MAPPING.SERVICE_KEY " +
            "   FROM SERVICE_CATALOG " +
            "   INNER JOIN API_SERVICE_MAPPING" +
            "   ON " +
            "   API_SERVICE_MAPPING.SERVICE_KEY = SERVICE_CATALOG.SERVICE_KEY " +
            "   WHERE " +
            "   API_SERVICE_MAPPING.API_ID = ?";

    public static final String GET_SUBSCRIBED_USERS_FOR_API_SQL =
            " SELECT " +
            "   SB.USER_ID, " +
            " FROM " +
            "   SUBSCRIBER SB, " +
            "   APPLICATION APP, " +
            "   SUBSCRIPTION SP, " +
            "   API API " +
            " WHERE " +
            "   API.API_PROVIDER = ? " +
            "   AND API.API_NAME = ?" +
            "   AND API.API_VERSION = ?" +
            "   AND SP.APPLICATION_ID = APP.APPLICATION_ID " +
            "   AND APP.SUBSCRIBER_ID=SB.SUBSCRIBER_ID " +
            "   AND API.API_ID = SP.API_ID" +
            "   AND SP.SUBS_CREATE_STATE = '" + APIConstants.SubscriptionCreatedStatus.SUBSCRIBE + "'";

    public static final String ADD_SUBSCRIBER_SQL =
            " INSERT" +
            "   INTO SUBSCRIBER (USER_ID, ORGANIZATION, EMAIL_ADDRESS, DATE_SUBSCRIBED, CREATED_BY, CREATED_TIME, " +
                    "UPDATED_TIME) " +
            " VALUES (?,?,?,?,?,?,?)";

    public static final String ADD_MONETIZATION_USAGE_PUBLISH_INFO =
            " INSERT INTO MONETIZATION_USAGE (ID, STATE, STATUS, STARTED_TIME, PUBLISHED_TIME) VALUES (?,?,?,?,?)";

    public static final String UPDATE_MONETIZATION_USAGE_PUBLISH_INFO =
            " UPDATE MONETIZATION_USAGE SET STATE = ?, STATUS = ?, STARTED_TIME = ?, PUBLISHED_TIME = ?" +
                    " WHERE ID = ?";

    public static final String GET_MONETIZATION_USAGE_PUBLISH_INFO =
            " SELECT ID, STATE, STATUS, STARTED_TIME, PUBLISHED_TIME FROM MONETIZATION_USAGE";

    public static final String UPDATE_SUBSCRIBER_SQL =
            " UPDATE SUBSCRIBER " +
            " SET" +
            "   USER_ID=?," +
            "   ORGANIZATION=?," +
            "   EMAIL_ADDRESS=?," +
            "   DATE_SUBSCRIBED=?," +
            "   UPDATED_BY=?," +
            "   UPDATED_TIME=? " +
            " WHERE" +
            "   SUBSCRIBER_ID=?";

    public static final String GET_SUBSCRIBER_SQL =
            " SELECT " +
            "   USER_ID, ORGANIZATION, EMAIL_ADDRESS, DATE_SUBSCRIBED " +
            " FROM " +
            "   SUBSCRIBER " +
            " WHERE " +
            "   SUBSCRIBER_ID=?";

    public static final String CHECK_EXISTING_SUBSCRIPTION_API_SQL =
            " SELECT " +
            "   SUB_STATUS, SUBS_CREATE_STATE " +
            " FROM " +
            "   SUBSCRIPTION " +
            " WHERE " +
            "   API_ID = ? " +
            "   AND APPLICATION_ID = ?";

    public static final String RETRIEVE_SUBSCRIPTION_ID_SQL =
            " SELECT " +
            "   SUBSCRIPTION_ID " +
            " FROM " +
            "   SUBSCRIPTION " +
            " WHERE " +
            "   UUID = ? ";

    public static final String ADD_SUBSCRIPTION_SQL =
            " INSERT INTO " +
            "   SUBSCRIPTION (TIER_ID,API_ID,APPLICATION_ID,SUB_STATUS,SUBS_CREATE_STATE,CREATED_BY,CREATED_TIME, " +
            "   UPDATED_TIME, UUID, TIER_ID_PENDING) " +
            " VALUES (?,?,?,?,?,?,?,?,?,?)";

    public static final String UPDATE_SINGLE_SUBSCRIPTION_SQL =
            " UPDATE SUBSCRIPTION " +
            " SET TIER_ID_PENDING = ? " +
            " , SUB_STATUS = ? " +
            " WHERE UUID = ?";

    public static final String GET_SUBSCRIPTION_UUID_SQL =
            " SELECT UUID " +
            " FROM SUBSCRIPTION " +
            " WHERE " +
            "   API_ID = ? " +
            "   AND APPLICATION_ID = ?";

    public static final String GET_SUBSCRIBER_ID_BY_SUBSCRIPTION_UUID_SQL =
            " SELECT APPS.SUBSCRIBER_ID  AS SUBSCRIBER_ID " +
            " FROM " +
            " APPLICATION APPS, " +
            " SUBSCRIPTION SUBS" +
            " WHERE " +
            " SUBS.APPLICATION_ID = APPS.APPLICATION_ID " +
            " AND SUBS.UUID = ?";

    public static final String GET_SUBSCRIPTION_STATUS_BY_UUID_SQL =
            " SELECT SUB_STATUS " +
            " FROM SUBSCRIPTION " +
            " WHERE UUID = ?";

    public static final String UPDATE_SUBSCRIPTION_SQL =
            " UPDATE SUBSCRIPTION " +
            " SET SUBS_CREATE_STATE = '" + APIConstants.SubscriptionCreatedStatus.UN_SUBSCRIBE + "' " +
            " WHERE UUID = ?";

    public static final String REMOVE_SUBSCRIPTION_SQL =
            " DELETE FROM SUBSCRIPTION WHERE UUID = ?";

    public static final String REMOVE_SUBSCRIPTION_BY_ID_SQL =
            " DELETE FROM SUBSCRIPTION WHERE SUBSCRIPTION_ID = ?";

    public static final String REMOVE_ALL_SUBSCRIPTIONS_SQL =
            " DELETE FROM SUBSCRIPTION WHERE API_ID = ?";

    public static final String GET_SUBSCRIPTION_STATUS_BY_ID_SQL =
            " SELECT SUB_STATUS FROM SUBSCRIPTION WHERE SUBSCRIPTION_ID = ?";

    public static final String GET_SUBSCRIPTION_BY_ID_SQL =
            " SELECT " +
            "   SUBS.SUBSCRIPTION_ID AS SUBSCRIPTION_ID, " +
            "   API.API_PROVIDER AS API_PROVIDER, " +
            "   API.API_NAME AS API_NAME, " +
            "   API.API_VERSION AS API_VERSION, " +
            "   API.API_TYPE AS API_TYPE, " +
            "   API.ORGANIZATION AS ORGANIZATION, " +
            "   SUBS.APPLICATION_ID AS APPLICATION_ID, " +
            "   SUBS.TIER_ID AS TIER_ID, " +
            "   SUBS.TIER_ID_PENDING AS TIER_ID_PENDING, " +
            "   SUBS.SUB_STATUS AS SUB_STATUS, " +
            "   SUBS.SUBS_CREATE_STATE AS SUBS_CREATE_STATE, " +
            "   SUBS.UUID AS UUID, " +
             "   API.API_ID AS API_ID," +
             "   API.API_UUID AS API_UUID " +
            " FROM " +
            "   SUBSCRIPTION SUBS," +
            "   API API " +
            " WHERE " +
            "   API.API_ID = SUBS.API_ID " +
            "   AND SUBSCRIPTION_ID = ?";

    public static final String GET_SUBSCRIPTION_BY_UUID_SQL =
            " SELECT " +
            "   SUBS.SUBSCRIPTION_ID AS SUBSCRIPTION_ID, " +
            "   API.API_PROVIDER AS API_PROVIDER, " +
            "   API.API_NAME AS API_NAME, " +
            "   API.API_VERSION AS API_VERSION, " +
            "   API.API_TYPE AS API_TYPE, " +
            "   API.ORGANIZATION AS ORGANIZATION, " +
            "   SUBS.APPLICATION_ID AS APPLICATION_ID, " +
            "   SUBS.TIER_ID AS TIER_ID, " +
            "   SUBS.TIER_ID_PENDING AS TIER_ID_PENDING, " +
            "   SUBS.SUB_STATUS AS SUB_STATUS, " +
            "   SUBS.SUBS_CREATE_STATE AS SUBS_CREATE_STATE, " +
            "   SUBS.UUID AS UUID, " +
            "   SUBS.CREATED_TIME AS CREATED_TIME, " +
            "   SUBS.UPDATED_TIME AS UPDATED_TIME, " +
            "   API.API_UUID AS API_UUID, " +
            "   API.API_ID AS API_ID " +
            " FROM " +
            "   SUBSCRIPTION SUBS," +
            "   API API " +
            " WHERE " +
            "   API.API_ID = SUBS.API_ID " +
            "   AND UUID = ?";

    public static final String GET_TENANT_SUBSCRIBER_SQL =
            " SELECT " +
            "   SUBSCRIBER_ID, " +
            "   USER_ID, " +
            "   ORGANIZATION, " +
            "   EMAIL_ADDRESS, " +
            "   DATE_SUBSCRIBED " +
            " FROM " +
            "   SUBSCRIBER " +
            " WHERE " +
            "   USER_ID = ? " +
            "   AND ORGANIZATION = ?";

    public static final String GET_TENANT_SUBSCRIBER_CASE_INSENSITIVE_SQL =
            " SELECT " +
            "   SUBSCRIBER_ID, " +
            "   USER_ID, " +
            "   ORGANIZATION, " +
            "   EMAIL_ADDRESS, " +
            "   DATE_SUBSCRIBED " +
            " FROM " +
            "   SUBSCRIBER " +
            " WHERE " +
            "   LOWER(USER_ID) = LOWER(?) " +
            "   AND ORGANIZATION = ?";

    public static final String GET_SUBSCRIBED_APIS_SQL =
            " SELECT " +
            "   SUBS.SUBSCRIPTION_ID AS SUBS_ID, " +
            "   API.API_UUID AS API_UUID, " +
            "   API.API_PROVIDER AS API_PROVIDER, " +
            "   API.API_NAME AS API_NAME, " +
            "   API.API_VERSION AS API_VERSION, " +
            "   SUBS.TIER_ID AS TIER_ID, " +
            "   APP.APPLICATION_ID AS APP_ID, " +
            "   SUBS.SUB_STATUS AS SUB_STATUS, " +
            "   SUBS.SUBS_CREATE_STATE AS SUBS_CREATE_STATE, " +
            "   APP.NAME AS APP_NAME, " +
            "   APP.CALLBACK_URL AS CALLBACK_URL, " +
            "   SUBS.UUID AS SUB_UUID, " +
            "   APP.UUID AS APP_UUID " +
            " FROM " +
            "   SUBSCRIBER SUB," +
            "   APPLICATION APP, " +
            "   SUBSCRIPTION SUBS, " +
            "   API API " +
            " WHERE " +
            "   SUB.ORGANIZATION = ? " +
            "   AND APP.APPLICATION_ID=SUBS.APPLICATION_ID " +
            "   AND SUB.SUBSCRIBER_ID=APP.SUBSCRIBER_ID " +
            "   AND API.API_ID=SUBS.API_ID" +
            "   AND APP.NAME= ? " +
            "   AND SUBS.SUBS_CREATE_STATE = '" + APIConstants.SubscriptionCreatedStatus.SUBSCRIBE + "'";

    public static final String GET_SUBSCRIBED_APIS_BY_ID_SQL =
            " SELECT " +
                    "   SUBS.SUBSCRIPTION_ID AS SUBS_ID, " +
                    "   API.API_PROVIDER AS API_PROVIDER, " +
                    "   API.API_NAME AS API_NAME, " +
                    "   API.API_VERSION AS API_VERSION, " +
                    "   SUBS.TIER_ID AS TIER_ID, " +
                    "   APP.APPLICATION_ID AS APP_ID, " +
                    "   SUBS.SUB_STATUS AS SUB_STATUS, " +
                    "   SUBS.SUBS_CREATE_STATE AS SUBS_CREATE_STATE, " +
                    "   APP.NAME AS APP_NAME, " +
                    "   APP.CALLBACK_URL AS CALLBACK_URL, " +
                    "   SUBS.UUID AS SUB_UUID, " +
                    "   APP.UUID AS APP_UUID, " +
                    "   APP.CREATED_BY AS OWNER" +
                    " FROM " +
                    "   SUBSCRIBER SUB," +
                    "   APPLICATION APP, " +
                    "   SUBSCRIPTION SUBS, " +
                    "   API API " +
                    " WHERE " +
                    "   SUB.ORGANIZATION = ? " +
                    "   AND APP.APPLICATION_ID=SUBS.APPLICATION_ID " +
                    "   AND SUB.SUBSCRIBER_ID=APP.SUBSCRIBER_ID " +
                    "   AND API.API_ID=SUBS.API_ID" +
                    "   AND APP.APPLICATION_ID= ? " +
                    "   AND SUBS.SUBS_CREATE_STATE = '" + APIConstants.SubscriptionCreatedStatus.SUBSCRIBE + "'";

    public static final String GET_SUBSCRIPTION_COUNT_SQL =
            " SELECT COUNT(*) AS SUB_COUNT " +
            " FROM " +
            "   SUBSCRIPTION SUBS, APPLICATION APP, SUBSCRIBER SUB " +
            " WHERE SUBS.SUBS_CREATE_STATE ='" + APIConstants.SubscriptionCreatedStatus.SUBSCRIBE + "'" +
            "   AND SUBS.APPLICATION_ID = APP.APPLICATION_ID" +
            "   AND APP.NAME=?" +
            "   AND APP.SUBSCRIBER_ID= SUB.SUBSCRIBER_ID" +
            "   AND SUB.ORGANIZATION=?";

    public static final String GET_SUBSCRIPTION_COUNT_BY_APP_ID_SQL =
            " SELECT COUNT(*) AS SUB_COUNT " +
                    " FROM " +
                    "   SUBSCRIPTION SUBS, APPLICATION APP,API API " +
                    " WHERE SUBS.SUBS_CREATE_STATE ='" + APIConstants.SubscriptionCreatedStatus.SUBSCRIBE + "'" +
                    "   AND SUBS.APPLICATION_ID = APP.APPLICATION_ID" +
                    "   AND API.API_ID = SUBS.API_ID" +
                    "   AND APP.APPLICATION_ID = ?" +
                    "   AND API.ORGANIZATION = ?";

    public static final String GET_SUBSCRIPTION_COUNT_CASE_INSENSITIVE_SQL =
            " SELECT COUNT(*) AS SUB_COUNT " +
            " FROM " +
            "   SUBSCRIPTION SUBS,APPLICATION APP,SUBSCRIBER SUB " +
            " WHERE " +
            "   SUBS.SUBS_CREATE_STATE ='" + APIConstants.SubscriptionCreatedStatus.SUBSCRIBE + "'" +
            "   AND SUBS.APPLICATION_ID = APP.APPLICATION_ID" +
            "   AND APP.NAME=?" +
            "   AND APP.SUBSCRIBER_ID= SUB.SUBSCRIBER_ID" +
            "   AND SUB.ORGANIZATION=?";

    public static final String GET_SUBSCRIPTION_COUNT_BY_APP_ID_CASE_INSENSITIVE_SQL =
            " SELECT COUNT(*) AS SUB_COUNT " +
                    " FROM " +
                    "   SUBSCRIPTION SUBS,APPLICATION APP,SUBSCRIBER SUB " +
                    " WHERE " +
                    "   SUBS.SUBS_CREATE_STATE ='" + APIConstants.SubscriptionCreatedStatus.SUBSCRIBE + "'" +
                    "   AND SUBS.APPLICATION_ID = APP.APPLICATION_ID" +
                    "   AND APP.APPLICATION_ID=?" +
                    "   AND APP.SUBSCRIBER_ID= SUB.SUBSCRIBER_ID" +
                    "   AND SUB.ORGANIZATION=?";

    public static final String GET_PAGINATED_SUBSCRIBED_APIS_SQL =
            " SELECT " +
            "   API.API_TYPE AS TYPE, " +
            "   SUBS.UUID AS SUB_UUID, " +
            "   SUBS.SUBSCRIPTION_ID, " +
            "   API.API_PROVIDER AS API_PROVIDER, " +
            "   API.API_NAME AS API_NAME, " +
            "   API.API_VERSION AS API_VERSION, " +
            "   SUBS.TIER_ID AS TIER_ID, " +
            "   SUBS.TIER_ID_PENDING AS TIER_ID_PENDING, " +
            "   APP.APPLICATION_ID AS APP_ID, " +
            "   APP.UUID AS APP_UUID, " +
            "   SUBS.SUB_STATUS AS SUB_STATUS, " +
            "   SUBS.SUBS_CREATE_STATE AS SUBS_CREATE_STATE, " +
            "   APP.NAME AS APP_NAME, " +
            "   APP.CALLBACK_URL AS CALLBACK_URL " +
            " FROM " +
            "   SUBSCRIBER SUB," +
            "   APPLICATION APP, " +
            "   SUBSCRIPTION SUBS, " +
            "   API API " +
            " WHERE " +
            "   SUB.ORGANIZATION = ? " +
            "   AND SUB.SUBSCRIBER_ID=APP.SUBSCRIBER_ID " +
            "   AND APP.APPLICATION_ID=SUBS.APPLICATION_ID " +
            "   AND API.API_ID=SUBS.API_ID" +
            "   AND APP.NAME= ? " +
            "   AND SUBS.SUBS_CREATE_STATE = '" + APIConstants.SubscriptionCreatedStatus.SUBSCRIBE + "'";


    public static final String GET_PAGINATED_SUBSCRIBED_APIS_BY_APP_ID_SQL =
            " SELECT " +
                    "   SUBS.SUBSCRIPTION_ID, " +
                    "   API.API_PROVIDER AS API_PROVIDER, " +
                    "   API.API_UUID AS API_UUID, " +
                    "   API.API_NAME AS API_NAME, " +
                    "   API.API_TYPE AS TYPE, " +
                    "   API.API_VERSION AS API_VERSION, " +
                    "   SUBS.TIER_ID AS TIER_ID, " +
                    "   SUBS.TIER_ID_PENDING AS TIER_ID_PENDING, " +
                    "   APP.APPLICATION_ID AS APP_ID, " +
                    "   SUBS.SUB_STATUS AS SUB_STATUS, " +
                    "   SUBS.UUID AS SUB_UUID, " +
                    "   SUBS.SUBS_CREATE_STATE AS SUBS_CREATE_STATE, " +
                    "   APP.NAME AS APP_NAME " +
                    " FROM " +
                    "   APPLICATION APP, " +
                    "   SUBSCRIPTION SUBS, " +
                    "   API API " +
                    " WHERE " +
                    "   APP.APPLICATION_ID=SUBS.APPLICATION_ID " +
                    "   AND API.API_ID=SUBS.API_ID " +
                    "   AND APP.APPLICATION_ID = ? " +
                    "   AND API.ORGANIZATION = ?"+
                    "   AND SUBS.SUBS_CREATE_STATE = '" + APIConstants.SubscriptionCreatedStatus.SUBSCRIBE + "'" +
                    " ORDER BY API_NAME ASC";

    public static final String GET_SUBSCRIBED_APIS_BY_APP_ID_SQL =
            " SELECT " +
                    "   SUBS.SUBSCRIPTION_ID, " +
                    "   API.API_PROVIDER AS API_PROVIDER, " +
                    "   API.API_UUID AS API_UUID, " +
                    "   API.API_NAME AS API_NAME, " +
                    "   API.API_TYPE AS TYPE, " +
                    "   API.API_VERSION AS API_VERSION, " +
                    "   SUBS.TIER_ID AS TIER_ID, " +
                    "   SUBS.TIER_ID_PENDING AS TIER_ID_PENDING, " +
                    "   APP.APPLICATION_ID AS APP_ID, " +
                    "   SUBS.SUB_STATUS AS SUB_STATUS, " +
                    "   SUBS.UUID AS SUB_UUID, " +
                    "   SUBS.SUBS_CREATE_STATE AS SUBS_CREATE_STATE, " +
                    "   APP.NAME AS APP_NAME " +
                    " FROM " +
                    "   APPLICATION APP, " +
                    "   SUBSCRIPTION SUBS, " +
                    "   API API " +
                    " WHERE " +
                    "   APP.APPLICATION_ID=SUBS.APPLICATION_ID " +
                    "   AND API.API_ID=SUBS.API_ID " +
                    "   AND APP.APPLICATION_ID = ? " +
                    "   AND SUBS.SUBS_CREATE_STATE = '" + APIConstants.SubscriptionCreatedStatus.SUBSCRIBE + "'" +
                    " ORDER BY API_NAME ASC";

    public static final String GET_SUBSCRIBED_APIS_OF_SUBSCRIBER_SQL =
            " SELECT " +
            "   API.API_TYPE AS TYPE, " +
            "   SUBS.SUBSCRIPTION_ID AS SUBS_ID, " +
            "   API.API_PROVIDER AS API_PROVIDER, " +
            "   API.API_NAME AS API_NAME, " +
            "   API.API_UUID AS API_UUID, " +
            "   API.API_VERSION AS API_VERSION, " +
            "   SUBS.TIER_ID AS TIER_ID, " +
            "   SUBS.TIER_ID_PENDING AS TIER_ID_PENDING, " +
            "   APP.APPLICATION_ID AS APP_ID, " +
            "   SUBS.SUB_STATUS AS SUB_STATUS, " +
            "   SUBS.SUBS_CREATE_STATE AS SUBS_CREATE_STATE, " +
            "   APP.NAME AS APP_NAME, " +
            "   APP.TOKEN_TYPE AS APP_TOKEN_TYPE, " +
            "   APP.CALLBACK_URL AS CALLBACK_URL, " +
            "   SUBS.UUID AS SUB_UUID, " +
            "   APP.UUID AS APP_UUID, " +
            "   APP.CREATED_BY AS OWNER" +
            " FROM " +
            "   SUBSCRIBER SUB," +
            "   APPLICATION APP, " +
            "   SUBSCRIPTION SUBS, " +
            "   API API " +
            " WHERE " +
            "   SUB.ORGANIZATION = ? " +
            "   AND SUB.SUBSCRIBER_ID=APP.SUBSCRIBER_ID " +
            "   AND APP.APPLICATION_ID=SUBS.APPLICATION_ID" +
            "   AND API.API_ID=SUBS.API_ID " +
            "   AND API.ORGANIZATION = ? " +
            "   AND SUBS.SUBS_CREATE_STATE = '" + APIConstants.SubscriptionCreatedStatus.SUBSCRIBE + "'";

    public static final String GET_CLIENT_OF_APPLICATION_SQL =
            " SELECT  CONSUMER_KEY,KEY_MANAGER " +
            " FROM APPLICATION_KEY_MAPPING " +
            " WHERE APPLICATION_ID = ? AND KEY_TYPE = ?";

    public static final String GET_CONSUMER_KEYS_OF_APPLICATION_SQL =
            " SELECT CONSUMER_KEY " +
            " FROM APPLICATION_KEY_MAPPING " +
            " WHERE APPLICATION_ID = ?";


    public static final String GET_THROTTLE_TIER_PERMISSION_ID_SQL =
            " SELECT THROTTLE_TIER_PERMISSIONS_ID " +
            " FROM THROTTLE_TIER_PERMISSIONS " +
            " WHERE TIER = ? AND " + "ORGANIZATION = ?";

    public static final String ADD_THROTTLE_TIER_PERMISSION_SQL =
            " INSERT INTO" +
            "   THROTTLE_TIER_PERMISSIONS (TIER, PERMISSIONS_TYPE, ROLES, ORGANIZATION)" +
            " VALUES(?, ?, ?, ?)";

    public static final String UPDATE_THROTTLE_TIER_PERMISSION_SQL =
            " UPDATE" +
            "   THROTTLE_TIER_PERMISSIONS " +
            " SET " +
            "   TIER = ?, " +
            "   PERMISSIONS_TYPE = ?," +
            "   ROLES = ? " +
            " WHERE " +
            "   THROTTLE_TIER_PERMISSIONS_ID = ? " +
            "   AND ORGANIZATION = ?";

    public static final String DELETE_THROTTLE_TIER_PERMISSION_SQL = "DELETE FROM THROTTLE_TIER_PERMISSIONS WHERE " +
            "THROTTLE_TIER_PERMISSIONS_ID = ? AND ORGANIZATION = ?";

    public static final String GET_THROTTLE_TIER_PERMISSIONS_SQL =
            " SELECT TIER,PERMISSIONS_TYPE, ROLES " +
            " FROM THROTTLE_TIER_PERMISSIONS " +
            " WHERE ORGANIZATION = ?";

    public static final String GET_THROTTLE_TIER_PERMISSION_SQL =
            " SELECT PERMISSIONS_TYPE, ROLES " +
                    " FROM THROTTLE_TIER_PERMISSIONS " +
                    " WHERE TIER = ? AND ORGANIZATION = ?";

    public static final String DELETE_THROTTLE_TIER_BY_NAME_PERMISSION_SQL =
            " DELETE FROM " +
            " THROTTLE_TIER_PERMISSIONS " +
            " WHERE TIER = ? AND ORGANIZATION = ?";

  //--------------------

    public static final String GET_TIER_PERMISSION_ID_SQL =
            " SELECT TIER_PERMISSIONS_ID " +
            " FROM TIER_PERMISSIONS " +
            " WHERE TIER = ? AND " + "ORGANIZATION = ?";

    public static final String ADD_TIER_PERMISSION_SQL =
            " INSERT INTO" +
            "   TIER_PERMISSIONS (TIER, PERMISSIONS_TYPE, ROLES, ORGANIZATION)" +
            " VALUES(?, ?, ?, ?)";

    public static final String UPDATE_TIER_PERMISSION_SQL =
            " UPDATE" +
            "   TIER_PERMISSIONS " +
            " SET " +
            "   TIER = ?, " +
            "   PERMISSIONS_TYPE = ?," +
            "   ROLES = ? " +
            " WHERE " +
            "   TIER_PERMISSIONS_ID = ? " +
            "   AND ORGANIZATION = ?";

    public static final String GET_TIER_PERMISSIONS_SQL =
            " SELECT TIER , PERMISSIONS_TYPE , ROLES " +
            " FROM TIER_PERMISSIONS " +
            " WHERE ORGANIZATION = ?";

    public static final String GET_PERMISSION_OF_TIER_SQL =
            " SELECT PERMISSIONS_TYPE, ROLES " +
            " FROM TIER_PERMISSIONS " +
            " WHERE TIER = ? AND ORGANIZATION = ?";

    public static final String GET_SUBSCRIBERS_OF_PROVIDER_SQL =
            " SELECT " +
            "   SUBS.USER_ID AS USER_ID," +
            "   SUBS.EMAIL_ADDRESS AS EMAIL_ADDRESS, " +
            "   SUBS.DATE_SUBSCRIBED AS DATE_SUBSCRIBED " +
            " FROM " +
            "   SUBSCRIBER  SUBS," +
            "   APPLICATION  APP, " +
            "   SUBSCRIPTION SUB, " +
            "   API API " +
            " WHERE  " +
            "   SUB.APPLICATION_ID = APP.APPLICATION_ID " +
            "   AND SUBS. SUBSCRIBER_ID = APP.SUBSCRIBER_ID " +
            "   AND API.API_ID = SUB.API_ID " +
            "   AND API.API_PROVIDER = ?" +
            "   AND SUB.SUBS_CREATE_STATE = '" + APIConstants.SubscriptionCreatedStatus.SUBSCRIBE + "'";

    public static final String GET_SUBSCRIBERS_OF_API_SQL =
            " SELECT DISTINCT SB.USER_ID, SB.DATE_SUBSCRIBED " +
            " FROM " +
            "   SUBSCRIBER SB, " +
            "   SUBSCRIPTION SP," +
            "   APPLICATION APP," +
            "   API API " +
            " WHERE " +
            "   API.API_PROVIDER=? " +
            "   AND API.API_NAME=? " +
            "   AND API.API_VERSION=? " +
            "   AND SP.APPLICATION_ID=APP.APPLICATION_ID" +
            "   AND APP.SUBSCRIBER_ID=SB.SUBSCRIBER_ID " +
            "   AND API.API_ID = SP.API_ID" +
            "   AND SP.SUB_STATUS != '" + APIConstants.SubscriptionStatus.REJECTED + "'" +
            "   AND SP.SUBS_CREATE_STATE = '" + APIConstants.SubscriptionCreatedStatus.SUBSCRIBE + "'";

    public static final String GET_API_SUBSCRIPTION_COUNT_BY_API_SQL =
            " SELECT" +
            "   COUNT(SUB.SUBSCRIPTION_ID) AS SUB_ID " +
            " FROM " +
            "   SUBSCRIPTION SUB, " +
            "   API API " +
            " WHERE API.API_PROVIDER=? " +
            "   AND API.API_NAME=?" +
            "   AND API.API_VERSION=?" +
            "   AND API.API_ID=SUB.API_ID" +
            "   AND SUB.SUBS_CREATE_STATE = '" + APIConstants.SubscriptionCreatedStatus.SUBSCRIBE + "'";

    public static final String UPDATE_SUBSCRIPTION_OF_APPLICATION_SQL =
            " UPDATE SUBSCRIPTION " +
            " SET " +
            "   SUB_STATUS = ?, " +
            "   UPDATED_BY = ?, " +
            "   UPDATED_TIME = ? " +
            " WHERE " +
            "   API_ID = ? " +
            "   AND APPLICATION_ID = ?";

    public static final String UPDATE_SUBSCRIPTION_OF_UUID_SQL =
            " UPDATE SUBSCRIPTION " +
            " SET " +
            "   SUB_STATUS = ?, " +
            "   UPDATED_BY = ?, " +
            "   UPDATED_TIME = ? " +
            " WHERE " +
            "   UUID = ?";

    public static final String UPDATE_SUBSCRIPTION_STATUS_SQL =
            " UPDATE SUBSCRIPTION " +
            " SET SUB_STATUS = ? " +
            " WHERE SUBSCRIPTION_ID = ?";

    public static final String UPDATE_SUBSCRIPTION_STATUS_AND_TIER_SQL =
            " UPDATE SUBSCRIPTION " +
                    " SET TIER_ID_PENDING = ? " +
                    " , TIER_ID = ? " +
                    " , SUB_STATUS = ? " +
                    " WHERE SUBSCRIPTION_ID = ?";

    public static final String GET_REGISTRATION_APPROVAL_STATUS_SQL =
            " SELECT KEY_MANAGER,STATE FROM APPLICATION_KEY_MAPPING WHERE APPLICATION_ID = ? AND KEY_TYPE =?";

    public static final String UPDATE_APPLICAITON_KEY_TYPE_MAPPINGS_SQL =
            " UPDATE APPLICATION_KEY_MAPPING SET CONSUMER_KEY = ? , APP_INFO = ? WHERE APPLICATION_ID = ? AND " +
                    "KEY_TYPE = ? AND KEY_MANAGER = ?";

    public static final String UPDATE_APPLICATION_KEY_TYPE_MAPPINGS_METADATA_SQL =
            " UPDATE APPLICATION_KEY_MAPPING SET APP_INFO = ? WHERE APPLICATION_ID = ? AND " +
                    "KEY_TYPE = ? AND KEY_MANAGER = ?";

    public static final String ADD_APPLICATION_KEY_TYPE_MAPPING_SQL =
            " INSERT INTO " +
            " APPLICATION_KEY_MAPPING (APPLICATION_ID,CONSUMER_KEY,KEY_TYPE,STATE,CREATE_MODE,KEY_MANAGER,UUID) " +
            " VALUES (?,?,?,?,?,?,?)";

    public static final String UPDATE_APPLICATION_KEY_MAPPING_SQL =
            " UPDATE APPLICATION_KEY_MAPPING SET STATE = ? WHERE APPLICATION_ID = ? AND KEY_TYPE = ? AND " +
                    "KEY_MANAGER=?";

    public static final String GET_SUBSCRIPTION_SQL =
            " SELECT " +
            "   SUBS.TIER_ID ," +
            "   API.API_PROVIDER ," +
            "   API.API_NAME ," +
            "   API.API_VERSION ," +
            "   SUBS.APPLICATION_ID " +
            " FROM " +
            "   SUBSCRIPTION SUBS," +
            "   SUBSCRIBER SUB, " +
            "   APPLICATION  APP, " +
            "   API API " +
            " WHERE " +
            "   API.API_PROVIDER  = ?" +
            "   AND API.API_NAME = ?" +
            "   AND API.API_VERSION = ?" +
            "   AND SUB.USER_ID = ?" +
            "   AND SUB.ORGANIZATION = ? " +
            "   AND APP.SUBSCRIBER_ID = SUB.SUBSCRIBER_ID" +
            "   AND API.API_ID = SUBS.API_ID" +
            "   AND SUBS.SUBS_CREATE_STATE = '" + APIConstants.SubscriptionCreatedStatus.SUBSCRIBE + "'";

    public static final String GET_SUBSCRIPTION_CASE_INSENSITIVE_SQL =
            " SELECT " +
            "   SUBS.TIER_ID ," +
            "   API.API_PROVIDER ," +
            "   API.API_NAME ," +
            "   API.API_VERSION ," +
            "   SUBS.APPLICATION_ID " +
            " FROM " +
            "   SUBSCRIPTION SUBS," +
            "   SUBSCRIBER SUB, " +
            "   APPLICATION  APP, " +
            "   API API " +
            " WHERE " +
            "   API.API_PROVIDER  = ?" +
            "   AND API.API_NAME = ?" +
            "   AND API.API_VERSION = ?" +
            "   AND LOWER(SUB.USER_ID) = LOWER(?)" +
            "   AND SUB.ORGANIZATION = ? " +
            "   AND APP.SUBSCRIBER_ID = SUB.SUBSCRIBER_ID" +
            "   AND API.API_ID = SUBS.API_ID" +
            "   AND SUBS.SUBS_CREATE_STATE = '" + APIConstants.SubscriptionCreatedStatus.SUBSCRIBE + "'";

    public static final String GET_APP_SUBSCRIPTION_TO_API_SQL =
            " SELECT " +
                    "   SUBS.TIER_ID ," +
                    "   API.API_PROVIDER ," +
                    "   API.API_NAME ," +
                    "   API.API_VERSION ," +
                    "   SUBS.APPLICATION_ID " +
                    " FROM " +
                    "   SUBSCRIPTION SUBS," +
                    "   SUBSCRIBER SUB, " +
                    "   APPLICATION  APP, " +
                    "   API API " +
                    " WHERE " +
                    "   API.API_PROVIDER  = ?" +
                    "   AND API.API_NAME = ?" +
                    "   AND API.API_VERSION = ?" +
                    "   AND SUB.USER_ID = ?" +
                    "   AND SUB.ORGANIZATION = ? " +
                    "   AND SUBS.APPLICATION_ID = ? " +
                    "   AND APP.SUBSCRIBER_ID = SUB.SUBSCRIBER_ID" +
                    "   AND API.API_ID = SUBS.API_ID" +
                    "   AND SUBS.SUBS_CREATE_STATE = '" + APIConstants.SubscriptionCreatedStatus.SUBSCRIBE + "'";

    public static final String GET_APP_SUBSCRIPTION_TO_API_CASE_INSENSITIVE_SQL =
            " SELECT " +
                    "   SUBS.TIER_ID ," +
                    "   API.API_PROVIDER ," +
                    "   API.API_NAME ," +
                    "   API.API_VERSION ," +
                    "   SUBS.APPLICATION_ID " +
                    " FROM " +
                    "   SUBSCRIPTION SUBS," +
                    "   SUBSCRIBER SUB, " +
                    "   APPLICATION  APP, " +
                    "   API API " +
                    " WHERE " +
                    "   API.API_PROVIDER  = ?" +
                    "   AND API.API_NAME = ?" +
                    "   AND API.API_VERSION = ?" +
                    "   AND LOWER(SUB.USER_ID) = LOWER(?)" +
                    "   AND SUB.ORGANIZATION = ? " +
                    "   AND SUBS.APPLICATION_ID = ? " +
                    "   AND APP.SUBSCRIBER_ID = SUB.SUBSCRIBER_ID" +
                    "   AND API.API_ID = SUBS.API_ID" +
                    "   AND SUBS.SUBS_CREATE_STATE = '" + APIConstants.SubscriptionCreatedStatus.SUBSCRIBE + "'";

    public static final String GET_APP_API_USAGE_BY_PROVIDER_SQL =
            " SELECT " +
            "   SUBS.SUBSCRIPTION_ID AS SUBSCRIPTION_ID, " +
            "   SUBS.APPLICATION_ID AS APPLICATION_ID, " +
            "   SUBS.SUB_STATUS AS SUB_STATUS, " +
            "   SUBS.TIER_ID AS TIER_ID, " +
            "   API.API_PROVIDER AS API_PROVIDER, " +
            "   API.API_NAME AS API_NAME, " +
            "   API.API_VERSION AS API_VERSION, " +
            "   SUB.USER_ID AS USER_ID, " +
            "   APP.NAME AS APPNAME, " +
            "   SUBS.UUID AS SUB_UUID, " +
            "   SUBS.TIER_ID AS SUB_TIER_ID, " +
            "   APP.UUID AS APP_UUID, " +
            "   SUBS.SUBS_CREATE_STATE AS SUBS_CREATE_STATE " +
            " FROM " +
            "   SUBSCRIPTION SUBS, " +
            "   APPLICATION APP, " +
            "   SUBSCRIBER SUB, " +
            "   API API " +
            " WHERE " +
            "   SUBS.APPLICATION_ID = APP.APPLICATION_ID " +
            "   AND APP.SUBSCRIBER_ID = SUB.SUBSCRIBER_ID " +
            "   AND API.API_PROVIDER = ? " +
            "   AND API.API_ID = SUBS.API_ID " +
            "   AND SUBS.SUB_STATUS != '" + APIConstants.SubscriptionStatus.REJECTED + "'" +
            " ORDER BY " +
            "   APP.NAME";

    public static final String GET_APP_API_USAGE_BY_UUID_SQL =
            " SELECT " +
            "   SUBS.SUBSCRIPTION_ID AS SUBSCRIPTION_ID, " +
            "   SUBS.APPLICATION_ID AS APPLICATION_ID, " +
            "   SUBS.SUB_STATUS AS SUB_STATUS, " +
            "   SUBS.TIER_ID AS TIER_ID, " +
            "   API.API_PROVIDER AS API_PROVIDER, " +
            "   API.API_NAME AS API_NAME, " +
            "   API.API_VERSION AS API_VERSION, " +
            "   SUB.USER_ID AS USER_ID, " +
            "   APP.NAME AS APPNAME, " +
            "   SUBS.UUID AS SUB_UUID, " +
            "   SUBS.TIER_ID AS SUB_TIER_ID, " +
            "   APP.UUID AS APP_UUID, " +
            "   SUBS.SUBS_CREATE_STATE AS SUBS_CREATE_STATE " +
            " FROM " +
            "   SUBSCRIPTION SUBS, " +
            "   APPLICATION APP, " +
            "   SUBSCRIBER SUB, " +
            "   API API " +
            " WHERE " +
            "   SUBS.APPLICATION_ID = APP.APPLICATION_ID " +
            "   AND APP.SUBSCRIBER_ID = SUB.SUBSCRIBER_ID " +
            "   AND API.API_ID = SUBS.API_ID " +
            "   AND API.API_UUID = ? " +
            "   AND API.ORGANIZATION= ? " +
            "   AND SUBS.SUB_STATUS != '" + APIConstants.SubscriptionStatus.REJECTED + "'" +
            " ORDER BY " +
            "   APP.NAME";

    public static final String GET_SUBSCRIPTIONS_OF_API_SQL =
            " SELECT " +
                    "   SUBS.SUBSCRIPTION_ID AS SUBSCRIPTION_ID, " +
                    "   SUBS.APPLICATION_ID AS APPLICATION_ID, " +
                    "   SUBS.SUB_STATUS AS SUB_STATUS, " +
                    "   SUBS.TIER_ID AS TIER_ID, " +
                    "   API.API_PROVIDER AS API_PROVIDER, " +
                    "   SUB.USER_ID AS USER_ID, " +
                    "   APP.NAME AS APPNAME, " +
                    "   SUBS.UUID AS SUB_UUID, " +
                    "   SUBS.CREATED_TIME AS SUB_CREATED_TIME, " +
                    "   SUBS.TIER_ID AS SUB_TIER_ID, " +
                    "   APP.UUID AS APP_UUID, " +
                    "   SUBS.SUBS_CREATE_STATE AS SUBS_CREATE_STATE " +
                    " FROM " +
                    "   SUBSCRIPTION SUBS, " +
                    "   APPLICATION APP, " +
                    "   SUBSCRIBER SUB, " +
                    "   API API " +
                    " WHERE " +
                    "   SUBS.APPLICATION_ID = APP.APPLICATION_ID " +
                    "   AND APP.SUBSCRIBER_ID = SUB.SUBSCRIBER_ID " +
                    "   AND API.API_NAME = ? " +
                    "   AND API.API_VERSION = ? " +
                    "   AND API.API_PROVIDER = ? " +
                    "   AND API.API_ID = SUBS.API_ID " +
                    "   AND SUBS.SUB_STATUS != '" + APIConstants.SubscriptionStatus.REJECTED + "'" +
                    " ORDER BY " +
                    "   APP.NAME";

    public static final String GET_API_RATING_SQL =
            "SELECT RATING FROM API_RATINGS WHERE API_ID= ? AND SUBSCRIBER_ID=? ";

    public static final String ADD_API_RATING_SQL =
            "INSERT INTO API_RATINGS (RATING_ID, RATING, API_ID, SUBSCRIBER_ID)  VALUES (?,?,?,?)";

    public static final String UPDATE_API_RATING_SQL =
            "UPDATE API_RATINGS SET RATING=? WHERE API_ID= ? AND SUBSCRIBER_ID=?";

    public static final String GET_API_RATING_ID_SQL =
            "SELECT RATING_ID FROM API_RATINGS WHERE API_ID= ? AND SUBSCRIBER_ID=? ";

    public static final String REMOVE_RATING_SQL =
            "DELETE FROM API_RATINGS WHERE RATING_ID =? ";

    public static final String GET_API_RATING_INFO_SQL =
            "SELECT RATING_ID, API_ID, RATING, SUBSCRIBER_ID FROM API_RATINGS WHERE SUBSCRIBER_ID  = ? "
                    + "AND API_ID= ? ";

    public static final String GET_API_ALL_RATINGS_SQL =
            "SELECT RATING_ID, API_ID, RATING, SUBSCRIBER_ID FROM API_RATINGS WHERE API_ID= ? ";

    public static final String GET_SUBSCRIBER_NAME_FROM_ID_SQL =
            "SELECT USER_ID FROM SUBSCRIBER WHERE SUBSCRIBER_ID = ? ";

    public static final String GET_RATING_INFO_BY_ID_SQL =
            "SELECT RATING_ID, API_ID, RATING, SUBSCRIBER_ID FROM API_RATINGS WHERE RATING_ID = ? "
                    + "AND API_ID= ? ";

    public static final String REMOVE_FROM_API_RATING_SQL =
            "DELETE FROM API_RATINGS WHERE API_ID=? ";

    public static final String GET_API_AVERAGE_RATING_SQL =
            " SELECT " +
            "   CAST( SUM(RATING) AS DECIMAL)/COUNT(RATING) AS RATING " +
            " FROM " +
            "   API_RATINGS " +
            " WHERE " +
            "   API_ID =? " +
            " GROUP BY " +
            "   API_ID ";

    public static final String APP_APPLICATION_SQL =
            " INSERT INTO APPLICATION (NAME, SUBSCRIBER_ID, APPLICATION_TIER, " +
            "   CALLBACK_URL, DESCRIPTION, APPLICATION_STATUS, GROUP_ID, CREATED_BY, CREATED_TIME, UPDATED_TIME, " +
                    "UUID, TOKEN_TYPE, ORGANIZATION)" +
            " VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)";

    public static final String UPDATE_APPLICATION_SQL =
            " UPDATE " +
            "   APPLICATION" +
            " SET " +
            "   NAME = ?," +
            "   APPLICATION_TIER = ?, " +
            "   CALLBACK_URL = ?, " +
            "   DESCRIPTION = ?, " +
            "   UPDATED_BY = ?, " +
            "   UPDATED_TIME = ?, " +
            "   TOKEN_TYPE = ? " +
            " WHERE" +
            "   APPLICATION_ID = ?";

    public static final String ADD_APPLICATION_ATTRIBUTES_SQL =
            " INSERT INTO APPLICATION_ATTRIBUTES (APPLICATION_ID, NAME, APP_ATTRIBUTE, ORGANIZATION) VALUES (?,?,?,?)";

    public static final String REMOVE_APPLICATION_ATTRIBUTES_SQL =
            " DELETE FROM " +
                    "   APPLICATION_ATTRIBUTES" +
                    " WHERE" +
                    "   APPLICATION_ID = ?";

    public static final String REMOVE_APPLICATION_ATTRIBUTES_BY_ATTRIBUTE_NAME_SQL =
            " DELETE FROM " +
                    "   APPLICATION_ATTRIBUTES" +
                    " WHERE" +
                    "   NAME = ? AND APPLICATION_ID = ?";

    public static final String UPDATE_APPLICATION_STATUS_SQL =
            " UPDATE APPLICATION SET APPLICATION_STATUS = ? WHERE APPLICATION_ID = ?";

    public static final String GET_APPLICATION_STATUS_BY_ID_SQL =
            "SELECT APPLICATION_STATUS FROM APPLICATION WHERE APPLICATION_ID= ?";

    public static final String GET_APPLICATION_ID_PREFIX =
            " SELECT " +
            "   APP.APPLICATION_ID " +
            " FROM " + "   " +
            "   APPLICATION APP," +
            "   SUBSCRIBER SUB " +
            " WHERE " +
            "   LOWER(APP.NAME) = LOWER(?)" + "   " +
            "   AND APP.ORGANIZATION = ? " +
            "   AND APP.SUBSCRIBER_ID = SUB.SUBSCRIBER_ID";

    public static final String GET_APPLICATION_ID_PREFIX_FOR_GROUP_COMPARISON = " SELECT APP.APPLICATION_ID FROM "
            + "APPLICATION APP, SUBSCRIBER SUB WHERE LOWER(APP.NAME) = LOWER(?) "
            + "AND APP.SUBSCRIBER_ID = SUB.SUBSCRIBER_ID";

    public static final String GET_APPLICATION_ID_SQL =
            "SELECT APPLICATION_ID FROM APPLICATION WHERE  SUBSCRIBER_ID  = ? AND NAME= ?";

    public static final String GET_APPLICATION_NAME_FROM_ID_SQL =
            "SELECT NAME FROM APPLICATION WHERE APPLICATION_ID = ?";

    public static final String GET_BASIC_APPLICATION_DETAILS_PREFIX =
            " SELECT " +
            "   APPLICATION_ID, " +
            "   NAME, " +
            "   APPLICATION_TIER, " +
            "   APP.SUBSCRIBER_ID,  " +
            "   CALLBACK_URL,  " +
            "   DESCRIPTION,  " +
            "   APPLICATION_STATUS,  " +
            "   USER_ID  " +
            " FROM " +
            "   APPLICATION APP, " +
            "   SUBSCRIBER SUB " +
            " WHERE " +
            "   SUB.SUBSCRIBER_ID = APP.SUBSCRIBER_ID ";

    public static final String GET_APPLICATIONS_PREFIX =
            "SELECT " +
            "   APPLICATION_ID, " +
            "   NAME," +
            "   APPLICATION_TIER," +
            "   APP.SUBSCRIBER_ID,  " +
            "   APP.TOKEN_TYPE,  " +
            "   CALLBACK_URL,  " +
            "   DESCRIPTION, " +
            "   APPLICATION_STATUS, " +
            "   USER_ID, " +
            "   GROUP_ID, " +
            "   UUID, " +
            "   APP.CREATED_BY " +
            " FROM" +
            "   APPLICATION APP, " +
            "   SUBSCRIBER SUB  " +
            " WHERE " +
            "   SUB.SUBSCRIBER_ID = APP.SUBSCRIBER_ID ";

    public static final String GET_APPLICATIONS_COUNT =
            "SELECT " +
            "   count(*) count " +
            " FROM" +
            "   APPLICATION APP, " +
            "   SUBSCRIBER SUB  " +
            " WHERE " +
            "   SUB.SUBSCRIBER_ID = APP.SUBSCRIBER_ID " +
            " AND " +
            "   SUB.ORGANIZATION=?" +
            " And "+
            "    ( SUB.CREATED_BY like ?" +
            " AND APP.NAME like ? )";

    public static final String GET_APPLICATION_BY_SUBSCRIBERID_AND_NAME_SQL =
            " SELECT " +
                    "   APP.APPLICATION_ID," +
                    "   APP.NAME," +
                    "   APP.SUBSCRIBER_ID," +
                    "   APP.APPLICATION_TIER," +
                    "   APP.CALLBACK_URL," +
                    "   APP.DESCRIPTION, " +
                    "   APP.SUBSCRIBER_ID," +
                    "   APP.APPLICATION_STATUS, " +
                    "   APP.GROUP_ID, " +
                    "   APP.UPDATED_TIME, "+
                    "   APP.CREATED_TIME, "+
                    "   APP.UUID," +
                    "   APP.TOKEN_TYPE," +
                    "   SUB.USER_ID " +
                    " FROM " +
                    "   SUBSCRIBER SUB," +
                    "   APPLICATION APP " +
                    " WHERE " +
                    "    APP.SUBSCRIBER_ID = ? " +
                    "  AND APP.NAME = ? " +
                    "   AND APP.SUBSCRIBER_ID = SUB.SUBSCRIBER_ID";

    public static final String GET_SIMPLE_APPLICATIONS =
            " SELECT " +
            "   APPLICATION_ID, " +
            "   NAME," +
            "   USER_ID, " +
            "   APP.CREATED_BY " +
            " FROM" +
            "   APPLICATION APP, " +
            "   SUBSCRIBER SUB  " +
            " WHERE " +
            "   SUB.SUBSCRIBER_ID = APP.SUBSCRIBER_ID ";

    public static final String GET_APPLICATIONS_BY_OWNER =
            "SELECT " +
            "   UUID, " +
            "   APPLICATION_ID, " +
            "   NAME," +
            "   CREATED_BY, " +
            "   APPLICATION_STATUS, " +
            "   GROUP_ID  " +
            " FROM" +
            "   APPLICATION " +
            " WHERE " +
            "   CREATED_BY = ?" +
            "   LIMIT ?, ? ";

    public static final String UPDATE_APPLICATION_OWNER =
            "UPDATE APPLICATION " +
            " SET " +
                "CREATED_BY = ? , " +
                "SUBSCRIBER_ID = ? " +
            " WHERE " +
            "   UUID = ? ";

        public static final String GET_APPLICATIONS_COUNNT_CASESENSITVE_WITHGROUPID = "SELECT " +
                "   count(*) count " +
                " FROM" +
                "   APPLICATION APP, " +
                "   SUBSCRIBER SUB  " +
                " WHERE " +
                "   SUB.SUBSCRIBER_ID = APP.SUBSCRIBER_ID " +
                " AND " +
                "   (GROUP_ID= ?  OR  ((GROUP_ID = '' OR GROUP_ID IS NULL) AND LOWER (SUB.USER_ID) = LOWER(?)))"+
                " And "+
                "    NAME like ?";


        public static final String GET_APPLICATIONS_COUNNT_NONE_CASESENSITVE_WITHGROUPID = "SELECT " +
                "   count(*) count " +
                " FROM" +
                "   APPLICATION APP, " +
                "   SUBSCRIBER SUB  " +
                " WHERE " +
                "   SUB.SUBSCRIBER_ID = APP.SUBSCRIBER_ID " +
                " AND " +
                "   (GROUP_ID= ?  OR ((GROUP_ID = '' OR GROUP_ID IS NULL) AND SUB.USER_ID=?))" +
                " And "+
                "    NAME like ?";

    public static final String GET_APPLICATIONS_COUNNT_CASESENSITVE_WITH_MULTIGROUPID = "SELECT " +
            "   count(*) count " +
            " FROM" +
            "   APPLICATION APP, " +
            "   SUBSCRIBER SUB  " +
            " WHERE " +
            "   SUB.SUBSCRIBER_ID = APP.SUBSCRIBER_ID " +
            " AND (" +
            "   (APPLICATION_ID IN ( SELECT APPLICATION_ID FROM APPLICATION_GROUP_MAPPING WHERE GROUP_ID IN ($params) AND TENANT = ?)) " +
            "           OR   " +
            "   LOWER (SUB.USER_ID) = LOWER(?) )"+
            " And "+
            "    NAME like ?";


    public static final String GET_APPLICATIONS_COUNNT_NONE_CASESENSITVE_WITH_MULTIGROUPID = "SELECT " +
            "   count(*) count " +
            " FROM" +
            "   APPLICATION APP, " +
            "   SUBSCRIBER SUB  " +
            " WHERE " +
            "   SUB.SUBSCRIBER_ID = APP.SUBSCRIBER_ID " +
            " AND (" +
            "    (APPLICATION_ID IN ( SELECT APPLICATION_ID FROM APPLICATION_GROUP_MAPPING WHERE GROUP_ID IN ($params) AND TENANT = ?)) " +
            "        OR " +
            "    SUB.USER_ID = ? )" +
            " And "+
            "    NAME like ?";

        public static final String GET_APPLICATIONS_COUNNT_CASESENSITVE = "SELECT " +
                "   count(*) count " +
                " FROM" +
                "   APPLICATION APP, " +
                "   SUBSCRIBER SUB  " +
                " WHERE " +
                "   SUB.SUBSCRIBER_ID = APP.SUBSCRIBER_ID " +
                " AND " +
                "    LOWER(SUB.USER_ID) = LOWER(?)"+
                " And "+
                "    NAME like ?";

        public static final String GET_APPLICATIONS_COUNNT_NONE_CASESENSITVE = "SELECT " +
                "   count(*) count " +
                " FROM" +
                "   APPLICATION APP, " +
                "   SUBSCRIBER SUB  " +
                " WHERE " +
                "   SUB.SUBSCRIBER_ID = APP.SUBSCRIBER_ID " +
                " AND " +
                "   SUB.USER_ID=?" +
                " And "+
                "    NAME like ?";




        public static final String GET_CONSUMER_KEYS_SQL =
            " SELECT " +
            "   MAP.CONSUMER_KEY " +
            " FROM " +
            "   SUBSCRIPTION SUB, " +
            "   APPLICATION_KEY_MAPPING MAP " +
            " WHERE " +
            "   SUB.APPLICATION_ID = MAP.APPLICATION_ID " +
            "   AND SUB.API_ID = ?";

    public static final String GET_SUBSCRIPTION_ID_OF_APPLICATION_SQL =
            "SELECT SUBSCRIPTION_ID FROM SUBSCRIPTION WHERE APPLICATION_ID = ?";

    public static final String GET_CONSUMER_KEY_OF_APPLICATION_SQL =
            "SELECT MAP.CONSUMER_KEY, MAP.CREATE_MODE, KM.NAME, KM.ORGANIZATION FROM APPLICATION_KEY_MAPPING MAP,"
                    + " KEY_MANAGER KM WHERE MAP.APPLICATION_ID = ? AND MAP.KEY_MANAGER = KM.UUID";

    public static final String REMOVE_APPLICATION_FROM_SUBSCRIPTIONS_SQL =
            "DELETE FROM SUBSCRIPTION WHERE APPLICATION_ID = ?";

    public static final String REMOVE_APPLICATION_FROM_APPLICATION_KEY_MAPPINGS_SQL =
            "DELETE FROM APPLICATION_KEY_MAPPING WHERE APPLICATION_ID = ?";

    public static final String REMOVE_APPLICATION_FROM_DOMAIN_MAPPINGS_SQL =
            "DELETE FROM APP_KEY_DOMAIN_MAPPING WHERE CONSUMER_KEY = ?";

    public static final String REMOVE_APPLICATION_FROM_APPLICATIONS_SQL =
            "DELETE FROM APPLICATION WHERE APPLICATION_ID = ?";

    public static final String REMOVE_APPLICATION_FROM_APPLICATION_REGISTRATIONS_SQL =
            "DELETE FROM APPLICATION_REGISTRATION WHERE APP_ID = ?";

    public static final String GET_CONSUMER_KEY_WITH_MODE_SLQ =
            " SELECT" +
            "   CONSUMER_KEY, " +
            "   KEY_TYPE" +
            " FROM" +
            "   APPLICATION_KEY_MAPPING " +
            " WHERE" +
            "   APPLICATION_ID = ? AND " +
            "   CREATE_MODE = ?";

    public static final String GET_CONSUMER_KEY_FOR_APPLICATION_KEY_TYPE_BY_APP_ID_SQL =
            " SELECT " +
                    "   AKM.CONSUMER_KEY " +
                    " FROM " +
                    "   APPLICATION APP," +
                    "   APPLICATION_KEY_MAPPING AKM," +
                    "   SUBSCRIBER SUB " +
                    " WHERE " +
                    "   SUB.SUBSCRIBER_ID=APP.SUBSCRIBER_ID " +
                    "   AND APP.APPLICATION_ID = AKM.APPLICATION_ID " +
                    "   AND APP.APPLICATION_ID = ? AND AKM.KEY_TYPE=?  ";

    public static final String GET_APPLICATION_ID_BY_CONSUMER_KEY_SQL =
            " SELECT " +
            "   MAP.APPLICATION_ID, " +
            "   MAP.KEY_TYPE " +
            " FROM " +
            "   APPLICATION_KEY_MAPPING MAP " +
            " WHERE " +
            "   MAP.CONSUMER_KEY = ? ";

    public static final String DELETE_APPLICATION_KEY_MAPPING_BY_CONSUMER_KEY_SQL =
            "DELETE FROM APPLICATION_KEY_MAPPING WHERE CONSUMER_KEY = ?";

    public static final String DELETE_APPLICATION_KEY_MAPPING_BY_UUID_SQL =
            "DELETE FROM APPLICATION_KEY_MAPPING WHERE UUID = ?";

    public static final String DELETE_APPLICATION_KEY_MAPPING_BY_APPLICATION_ID_SQL =
            "DELETE FROM APPLICATION_KEY_MAPPING WHERE APPLICATION_ID = ? AND KEY_TYPE = ?";

    public static final String REMOVE_FROM_APPLICATION_REGISTRANTS_SQL =
            "DELETE FROM APPLICATION_REGISTRATION WHERE APP_ID = ? AND TOKEN_TYPE = ? AND KEY_MANAGER = ?";

    public static final String GET_SUBSCRIBER_CASE_INSENSITIVE_SQL =
            " SELECT " +
            "   SUB.SUBSCRIBER_ID AS SUBSCRIBER_ID," +
            "   SUB.USER_ID AS USER_ID, " +
            "   SUB.ORGANIZATION AS ORGANIZATION," +
            "   SUB.EMAIL_ADDRESS AS EMAIL_ADDRESS," +
            "   SUB.DATE_SUBSCRIBED AS DATE_SUBSCRIBED " +
            " FROM " +
            "   SUBSCRIBER SUB " +
            " WHERE " +
            "   LOWER(SUB.USER_ID) = LOWER(?) " +
            "   AND SUB.ORGANIZATION = ?";

    public static final String GET_SUBSCRIBER_DETAILS_SQL =
            " SELECT " +
            "   SUB.SUBSCRIBER_ID AS SUBSCRIBER_ID," +
            "   SUB.USER_ID AS USER_ID, " +
            "   SUB.ORGANIZATION AS ORGANIZATION," +
            "   SUB.EMAIL_ADDRESS AS EMAIL_ADDRESS," +
            "   SUB.DATE_SUBSCRIBED AS DATE_SUBSCRIBED " +
            " FROM " +
            "   SUBSCRIBER SUB " +
            " WHERE " +
            "   SUB.USER_ID = ? " +
            "   AND SUB.ORGANIZATION = ?";

    public static final String GET_API_ID_SQL =
            "SELECT API.API_ID FROM API API WHERE API.API_PROVIDER = ? AND API.API_NAME = ? AND API.API_VERSION = ? ";
    public static final String GET_API_ID_SQL_BY_UUID =
            "SELECT API.API_ID FROM API API WHERE API.API_UUID = ?";
    public static final String GET_LIGHT_WEIGHT_API_INFO_BY_API_IDENTIFIER = "SELECT API_ID,API_UUID,API_PROVIDER," +
            "API_NAME,API_VERSION,CONTEXT,API_TYPE,STATUS FROM API WHERE API_PROVIDER = ? AND API_NAME = ? AND " +
            "API_VERSION = ? AND ORGANIZATION = ?";

    public static final String GET_API_PRODUCT_ID_SQL =
            "SELECT API_ID FROM API WHERE API_PROVIDER = ? AND API_NAME = ? "
                    + "AND API_VERSION = ? AND API_TYPE = '" + APIConstants.API_PRODUCT + "'";

    public static final String GET_API_PRODUCT_SQL =
            "SELECT API_ID, API_TIER FROM API WHERE API_UUID = ? AND API_TYPE = '" + APIConstants.API_PRODUCT + "'";

    public static final String GET_AUDIT_UUID_SQL =
            "SELECT MAP.AUDIT_UUID FROM SECURITY_AUDIT_UUID_MAPPING MAP WHERE MAP.API_ID = ?";

    public static final String ADD_SECURITY_AUDIT_MAP_SQL =
            "INSERT INTO SECURITY_AUDIT_UUID_MAPPING (API_ID, AUDIT_UUID) VALUES (?,?)";

    public static final String REMOVE_SECURITY_AUDIT_MAP_SQL =
            "DELETE FROM SECURITY_AUDIT_UUID_MAPPING WHERE API_ID = ?";

    public static final String ADD_CUSTOM_COMPLEXITY_DETAILS_SQL =
            "INSERT INTO GRAPHQL_COMPLEXITY (UUID, API_ID, TYPE, FIELD, COMPLEXITY_VALUE) VALUES (?,?,?,?,?)";

    public static final String GET_CUSTOM_COMPLEXITY_DETAILS_SQL =
            " SELECT" +
            "   TYPE," +
            "   FIELD," +
            "   COMPLEXITY_VALUE" +
            " FROM" +
            "   GRAPHQL_COMPLEXITY " +
            " WHERE" +
            "   API_ID = ? AND REVISION_UUID IS NULL";


    public static final String UPDATE_CUSTOM_COMPLEXITY_DETAILS_SQL =
            " UPDATE GRAPHQL_COMPLEXITY " +
            " SET " +
            "   COMPLEXITY_VALUE = ? " +
            " WHERE " +
            "    API_ID = ?" +
            "    AND TYPE = ? " +
            "    AND FIELD = ? AND REVISION_UUID IS NULL";

    public static final String REMOVE_FROM_GRAPHQL_COMPLEXITY_SQL =
            "DELETE FROM GRAPHQL_COMPLEXITY WHERE API_ID = ?";

    public static final String ADD_API_LIFECYCLE_EVENT_SQL =
            " INSERT INTO API_LC_EVENT (API_ID, PREVIOUS_STATE, NEW_STATE, USER_ID, TENANT_ID, EVENT_DATE)" +
            " VALUES (?,?,?,?,?,?)";

    public static final String GET_LIFECYCLE_EVENT_SQL =
            " SELECT" +
            "   LC.API_ID AS API_ID," +
            "   LC.PREVIOUS_STATE AS PREVIOUS_STATE," +
            "   LC.NEW_STATE AS NEW_STATE," +
            "   LC.USER_ID AS USER_ID," +
            "   LC.EVENT_DATE AS EVENT_DATE " +
            " FROM" +
            "   API_LC_EVENT LC " +
            " WHERE" +
            "   LC.API_ID = ?";

    public static final String GET_SUBSCRIPTION_DATA_SQL =
            " SELECT" +
            "   SUB.SUBSCRIPTION_ID AS SUBSCRIPTION_ID," +
            "   SUB.TIER_ID AS TIER_ID," +
            "   SUB.APPLICATION_ID AS APPLICATION_ID," +
            "   SUB.SUB_STATUS AS SUB_STATUS," +
            "   API.CONTEXT AS CONTEXT," +
            "   API.API_VERSION AS VERSION" +
            " FROM" +
            "   SUBSCRIPTION SUB," +
            "   API API " +
            " WHERE" +
            "   API.API_PROVIDER = ?" +
            "   AND API.API_NAME = ?" +
            "   AND API.API_VERSION IN (_API_VERSION_LIST_)" +
            "   AND API.API_ID = SUB.API_ID";

    public static final String GET_APPLICATION_DATA_SQL =
            " SELECT" +
            "   SUB.SUBSCRIPTION_ID AS SUBSCRIPTION_ID," +
            "   SUB.TIER_ID AS TIER_ID," +
            "   SUB.SUB_STATUS AS SUB_STATUS," +
            "   APP.APPLICATION_ID AS APPLICATION_ID," +
            "   API.CONTEXT AS CONTEXT " +
            " FROM" +
            "   SUBSCRIPTION SUB," +
            "   APPLICATION APP," +
            "   API API " +
            " WHERE" +
            "   API.API_PROVIDER = ?" +
            "   AND API.API_NAME = ?" +
            "   AND API.API_VERSION = ?" +
            "   AND SUB.APPLICATION_ID = APP.APPLICATION_ID" +
            "   AND API.API_ID = SUB.API_ID" +
            "   AND SUB.SUB_STATUS != '" + APIConstants.SubscriptionStatus.ON_HOLD + "'";

    public static final String ADD_API_SQL =
            " INSERT INTO API (API_PROVIDER,API_NAME,API_VERSION,CONTEXT,CONTEXT_TEMPLATE,CREATED_BY," +
                    "CREATED_TIME,API_TIER,API_TYPE,API_UUID,STATUS,ORGANIZATION,GATEWAY_VENDOR,VERSION_COMPARABLE)" +
                    " VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)";

    public static final String GET_DEFAULT_VERSION_SQL =
            "SELECT DEFAULT_API_VERSION FROM API_DEFAULT_VERSION WHERE API_NAME= ? AND API_PROVIDER= ? ";

    public static final String ADD_WORKFLOW_ENTRY_SQL =
            " INSERT INTO WORKFLOWS (WF_REFERENCE,WF_TYPE,WF_STATUS,WF_CREATED_TIME,WF_STATUS_DESC,TENANT_ID," +
            "TENANT_DOMAIN,WF_EXTERNAL_REFERENCE,WF_METADATA,WF_PROPERTIES)" +
            " VALUES (?,?,?,?,?,?,?,?,?,?)";

    public static final String UPDATE_WORKFLOW_ENTRY_SQL =
            " UPDATE WORKFLOWS " +
            " SET " +
            "   WF_STATUS = ?, " +
            "   WF_STATUS_DESC = ? " +
            " WHERE " +
           "    WF_EXTERNAL_REFERENCE = ?";

    public static final String GET_ALL_WORKFLOW_ENTRY_SQL =
            "SELECT * FROM WORKFLOWS WHERE WF_EXTERNAL_REFERENCE=?";

    public static final String GET_ALL_WORKFLOW_ENTRY_FROM_INTERNAL_REF_SQL =
            "SELECT * FROM WORKFLOWS WHERE WF_REFERENCE=? AND WF_TYPE=?";

    public static final String ADD_PAYLOAD_SQL =
            " UPDATE WORKFLOWS " +
                    " SET " +
                    "   WF_METADATA = ?, " +
                    "   WF_PROPERTIES = ?, " +
                    "   WF_STATUS_DESC = ? " +
                    " WHERE " +
                    "    WF_EXTERNAL_REFERENCE = ?";

    public static final String DELETE_WORKFLOW_REQUEST_SQL=
            " DELETE FROM WORKFLOWS WHERE WF_EXTERNAL_REFERENCE = ?";

    public static final String GET_ALL_WORKFLOW_DETAILS_BY_EXTERNALWORKFLOWREF =
            " SELECT  * FROM WORKFLOWS WHERE WF_EXTERNAL_REFERENCE = ?";

    public static final String GET_ALL_WORKFLOW_DETAILS_BY_WORKFLOW_TYPE =
            " SELECT  * FROM WORKFLOWS WHERE WF_TYPE = ? AND  WF_STATUS = ? AND TENANT_DOMAIN = ?";

    public static final String GET_ALL_WORKFLOW_DETAILS =
            " SELECT  * FROM WORKFLOWS WHERE WF_STATUS = ? AND TENANT_DOMAIN = ?";

    public static final String GET_ALL_WORKFLOW_DETAILS_BY_EXTERNAL_WORKFLOW_REFERENCE =
            " SELECT  * FROM WORKFLOWS " +
            " WHERE WF_EXTERNAL_REFERENCE = ? " +
            " AND WF_STATUS = ? " +
            " AND TENANT_DOMAIN = ?";

    public static final String UPDATE_PUBLISHED_DEFAULT_VERSION_SQL =
            " UPDATE API_DEFAULT_VERSION " +
            " SET " +
            "   PUBLISHED_DEFAULT_API_VERSION = ? " +
            " WHERE" +
            "   API_NAME = ? " +
            "   AND API_PROVIDER = ?";

    public static final String REMOVE_API_DEFAULT_VERSION_SQL =
            "DELETE FROM API_DEFAULT_VERSION WHERE API_NAME = ? AND API_PROVIDER = ?";

    public static final String GET_PUBLISHED_DEFAULT_VERSION_SQL =
            "SELECT PUBLISHED_DEFAULT_API_VERSION FROM API_DEFAULT_VERSION WHERE API_NAME= ? AND API_PROVIDER= ? ";

    public static final String ADD_API_DEFAULT_VERSION_SQL =
            " INSERT INTO " +
            " API_DEFAULT_VERSION(API_NAME,API_PROVIDER,DEFAULT_API_VERSION,PUBLISHED_DEFAULT_API_VERSION,"
                    + "ORGANIZATION) VALUES (?,?,?,?,?)";

    public static final String ADD_URL_MAPPING_SQL =
            " INSERT INTO " +
            " API_URL_MAPPING (API_ID,HTTP_METHOD,AUTH_SCHEME,URL_PATTERN,THROTTLING_TIER,MEDIATION_SCRIPT)" +
            " VALUES (?,?,?,?,?,?)";

    public static final String GET_APPLICATION_BY_NAME_PREFIX =
            " SELECT " +
            "   APP.APPLICATION_ID," +
            "   APP.NAME," +
            "   APP.SUBSCRIBER_ID," +
            "   APP.APPLICATION_TIER," +
            "   APP.CALLBACK_URL," +
            "   APP.DESCRIPTION, " +
            "   APP.SUBSCRIBER_ID," +
            "   APP.APPLICATION_STATUS," +
            "   APP.GROUP_ID," +
            "   APP.UUID," +
            "   APP.CREATED_BY," +
            "   APP.TOKEN_TYPE," +
            "   APP.UPDATED_TIME," +
            "   APP.CREATED_TIME," +
            "   SUB.USER_ID," +
            "   APP.CREATED_BY" +
            " FROM " +
            "   SUBSCRIBER SUB," +
            "   APPLICATION APP";

    public static final String GET_APPLICATION_ATTRIBUTES_BY_APPLICATION_ID =
            " SELECT " +
                    "   APP.APPLICATION_ID," +
                    "   APP.NAME," +
                    "   APP.APP_ATTRIBUTE" +
                    " FROM " +
                    "   APPLICATION_ATTRIBUTES APP WHERE APPLICATION_ID = ?";

    public static final String GET_APPLICATION_BY_ID_SQL =
            " SELECT " +
            "   APP.APPLICATION_ID," +
            "   APP.NAME," +
            "   APP.SUBSCRIBER_ID," +
            "   APP.APPLICATION_TIER," +
            "   APP.CALLBACK_URL," +
            "   APP.DESCRIPTION, " +
            "   APP.SUBSCRIBER_ID," +
            "   APP.APPLICATION_STATUS, " +
            "   APP.UPDATED_TIME," +
            "   APP.CREATED_TIME," +
            "   SUB.USER_ID, " +
            "   APP.GROUP_ID," +
            "   APP.CREATED_BY," +
            "   APP.UUID, " +
            "   APP.ORGANIZATION, " +
            "   APP.TOKEN_TYPE " +
            " FROM " +
            "   SUBSCRIBER SUB," +
            "   APPLICATION APP " +
            " WHERE " +
            "   APPLICATION_ID = ? " +
            "   AND APP.SUBSCRIBER_ID = SUB.SUBSCRIBER_ID";

    public static final String GET_APPLICATION_BY_UUID_SQL =
            " SELECT " +
            "   APP.APPLICATION_ID," +
            "   APP.NAME," +
            "   APP.SUBSCRIBER_ID," +
            "   APP.APPLICATION_TIER," +
            "   APP.CALLBACK_URL," +
            "   APP.DESCRIPTION, " +
            "   APP.SUBSCRIBER_ID," +
            "   APP.APPLICATION_STATUS, " +
            "   APP.GROUP_ID, " +
            "   APP.UPDATED_TIME, "+
            "   APP.CREATED_TIME, "+
            "   APP.UUID," +
            "   APP.TOKEN_TYPE," +
            "   APP.CREATED_BY," +
            "   APP.ORGANIZATION ORGANIZATION,"+
            "   SUB.USER_ID " +
            " FROM " +
            "   SUBSCRIBER SUB," +
            "   APPLICATION APP " +
            " WHERE " +
            "   APP.UUID = ? " +
            "   AND APP.SUBSCRIBER_ID = SUB.SUBSCRIBER_ID";

    public static final String GET_APPLICATION_BY_CLIENT_ID_SQL =
            " SELECT " +
                    "   APP.APPLICATION_ID," +
                    "   APP.NAME," +
                    "   APP.SUBSCRIBER_ID," +
                    "   APP.APPLICATION_TIER," +
                    "   APP.CALLBACK_URL," +
                    "   APP.DESCRIPTION, " +
                    "   APP.SUBSCRIBER_ID," +
                    "   APP.APPLICATION_STATUS, " +
                    "   APP.GROUP_ID, " +
                    "   APP.UPDATED_TIME, "+
                    "   APP.CREATED_TIME, "+
                    "   APP.UUID," +
                    "   APP.CREATED_BY," +
                    "   APP.TOKEN_TYPE," +
                    "   APP.ORGANIZATION," +
                    "   APP_MAP.KEY_TYPE" +
                    " FROM " +
                    "   APPLICATION_KEY_MAPPING APP_MAP," +
                    "   APPLICATION APP " +
                    " WHERE " +
                    "   APP_MAP.CONSUMER_KEY = ? " +
                    "   AND APP.APPLICATION_ID = APP_MAP.APPLICATION_ID";

    public static final String GET_APPLICATION_INFO_BY_CK =
            "SELECT APP.NAME as NAME, APP.UUID as UUID, APP.ORGANIZATION as ORGANIZATION, SUB.USER_ID as OWNER"
            + " FROM"
            + "     APPLICATION APP,"
            + "     APPLICATION_KEY_MAPPING APP_MAP,"
            + "     SUBSCRIBER SUB"
            + " WHERE APP_MAP.CONSUMER_KEY = ?"
            + "  AND APP.APPLICATION_ID = APP_MAP.APPLICATION_ID"
            + "  AND APP.SUBSCRIBER_ID = SUB.SUBSCRIBER_ID";

    public static final String REMOVE_FROM_URI_TEMPLATES_SQL =
            "DELETE FROM API_URL_MAPPING WHERE API_ID = ? AND REVISION_UUID IS NULL";

    //Product ID is recorded under revision column in the API_URL_MAPPING table
    public static final String REMOVE_FROM_URI_TEMPLATES__FOR_PRODUCTS_SQL =
            "DELETE FROM API_URL_MAPPING WHERE REVISION_UUID = ?";

    public static final String GET_ALL_URL_TEMPLATES_SQL =
            " SELECT    " +
                    "   AUM.HTTP_METHOD," +
                    "   AUM.AUTH_SCHEME," +
                    "   AUM.URL_PATTERN," +
                    "   AUM.THROTTLING_TIER," +
                    "   AUM.MEDIATION_SCRIPT, " +
                    "   AUM.URL_MAPPING_ID " +
                    " FROM " +
                    "   API_URL_MAPPING AUM, " +
                    "   API API " +
                    " WHERE" +
                    "   API.CONTEXT= ? " +
                    "   AND API.API_VERSION = ? " +
                    "   AND AUM.API_ID = API.API_ID AND AUM.REVISION_UUID IS NULL " +
                    " ORDER BY URL_MAPPING_ID";

    public static final String UPDATE_API_SQL =
            "UPDATE API " +
                    "SET " +
                    "   CONTEXT = ?, " +
                    "   API_NAME = ?, " +
                    "   CONTEXT_TEMPLATE = ?, " +
                    "   UPDATED_BY = ?," +
                    "   UPDATED_TIME = ?, " +
                    "   API_TIER = ?, " +
                    "   API_TYPE = ?, " +
                    "   GATEWAY_VENDOR = ? " +
                    " WHERE " +
                    "   API_UUID = ? ";

    public static final String GET_ORGANIZATION_BY_API_ID = "SELECT ORGANIZATION FROM API WHERE API_UUID = ?";

    public static final String GET_GATEWAY_VENDOR_BY_API_ID = "SELECT GATEWAY_VENDOR FROM API WHERE API_UUID = ?";

    public static final String FIX_NULL_THROTTLING_TIERS =
            "UPDATE API_URL_MAPPING SET THROTTLING_TIER = 'Unlimited' WHERE " +
                     " THROTTLING_TIER IS NULL";

    public static final String REMOVE_APPLICATION_MAPPINGS_BY_CONSUMER_KEY_SQL =
            "DELETE FROM APPLICATION_KEY_MAPPING WHERE CONSUMER_KEY = ?";

    public static final String REMOVE_FROM_API_LIFECYCLE_SQL =
            "DELETE FROM API_LC_EVENT WHERE API_ID=? ";

    public static final String REMOVE_FROM_API_COMMENT_SQL =
            "DELETE FROM API_COMMENTS WHERE API_ID=? ";

    public static final String REMOVE_FROM_API_SUBSCRIPTION_SQL =
            "DELETE FROM SUBSCRIPTION WHERE API_ID=?";

    public static final String REMOVE_FROM_EXTERNAL_STORES_SQL =
            "DELETE FROM EXTERNAL_STORES WHERE API_ID=?";

    public static final String REMOVE_FROM_API_SQL =
            "DELETE FROM API WHERE API_PROVIDER=? AND API_NAME=? AND API_VERSION=? ";

    public static final String REMOVE_FROM_API_SQL_BY_UUID =
            "DELETE FROM API WHERE API_UUID=?";

    public static final String REMOVE_FROM_API_URL_MAPPINGS_SQL =
            "DELETE FROM API_URL_MAPPING WHERE API_ID = ?";

    public static final String GET_API_LIST_SQL_BY_ORG = "SELECT API.API_ID, API.API_UUID,API.API_NAME," +
            "API.API_VERSION,API.API_PROVIDER FROM API API WHERE API.ORGANIZATION = ?";

    public static final String REMOVE_BULK_APIS_DATA_FROM_API_SQL = "DELETE FROM API WHERE API_UUID IN (_API_UUIDS_)";

    public static final String DELETE_BULK_API_WORKFLOWS_REQUEST_SQL = "DELETE FROM WORKFLOWS WHERE " +
            "WF_TYPE=\"API_STATE\" AND WF_REFERENCE IN (SELECT CONVERT(API.API_ID, CHAR) FROM API API " +
            "WHERE API.API_UUID IN (_API_UUIDS_))";

    public static final String DELETE_BULK_GW_PUBLISHED_API_DETAILS = "DELETE FROM GW_PUBLISHED_API_DETAILS WHERE " +
            "TENANT_DOMAIN = ?";

    public static final String GET_APPLICATION_BY_TIER_SQL =
            " SELECT DISTINCT AMS.APPLICATION_ID,NAME,SUBSCRIBER_ID " +
            " FROM " +
            "   SUBSCRIPTION AMS," +
            "   APPLICATION AMA " +
            "WHERE " +
            "   AMS.TIER_ID=? " +
            "   AND AMS.APPLICATION_ID=AMA.APPLICATION_ID";

    public static final String GET_URL_TEMPLATES_SQL =
            " SELECT " +
            "   URL_PATTERN," +
            "   HTTP_METHOD," +
            "   AUTH_SCHEME," +
            "   THROTTLING_TIER, " +
            "   MEDIATION_SCRIPT " +
            " FROM " +
            "   API_URL_MAPPING " +
            " WHERE " +
            "   API_ID = ? AND REVISION_UUID IS NULL" +
            " ORDER BY " +
            "   URL_MAPPING_ID ASC ";

    public static final String GET_URL_TEMPLATES_OF_API_SQL =
            " SELECT " +
            "  AUM.URL_MAPPING_ID," +
            "   AUM.URL_PATTERN," +
            "   AUM.HTTP_METHOD," +
            "   AUM.AUTH_SCHEME," +
            "   AUM.THROTTLING_TIER," +
            "   AUM.MEDIATION_SCRIPT," +
            "   ARSM.SCOPE_NAME " +
            " FROM " +
            "   API_URL_MAPPING AUM " +
            " INNER JOIN API API ON AUM.API_ID = API.API_ID " +
            " LEFT OUTER JOIN API_RESOURCE_SCOPE_MAPPING ARSM ON AUM.URL_MAPPING_ID = ARSM.URL_MAPPING_ID" +
            " WHERE " +
            "  API.API_UUID = ? AND " +
            "  AUM.REVISION_UUID IS NULL " +
            " ORDER BY AUM.URL_MAPPING_ID ASC ";

    public static final String GET_URL_MAPPING_IDS_OF_API_PRODUCT_SQL =
            " SELECT " +
            "  AUM.URL_MAPPING_ID" +
            " FROM " +
            "   API_URL_MAPPING AUM " +
            "  WHERE AUM.REVISION_UUID = ? ";

    public static final String GET_URL_TEMPLATES_OF_API_WITH_PRODUCT_MAPPINGS_SQL =
            " SELECT " +
                    "  AUM.URL_MAPPING_ID," +
                    "   AUM.URL_PATTERN," +
                    "   AUM.HTTP_METHOD," +
                    "   AUM.AUTH_SCHEME," +
                    "   AUM.THROTTLING_TIER," +
                    "   AUM.MEDIATION_SCRIPT," +
                    "   ARSM.SCOPE_NAME " +
                    " FROM " +
                    "   API_URL_MAPPING AUM " +
                    " INNER JOIN API API ON AUM.API_ID = API.API_ID " +
                    " LEFT OUTER JOIN API_RESOURCE_SCOPE_MAPPING ARSM ON AUM.URL_MAPPING_ID = ARSM.URL_MAPPING_ID" +
                    " WHERE " +
                    "  API.API_UUID = ? " +
                    " ORDER BY AUM.URL_MAPPING_ID ASC ";

    public static final String GET_URL_TEMPLATES_OF_API_REVISION_SQL =
            " SELECT " +
                    "  AUM.URL_MAPPING_ID," +
                    "   AUM.URL_PATTERN," +
                    "   AUM.HTTP_METHOD," +
                    "   AUM.AUTH_SCHEME," +
                    "   AUM.THROTTLING_TIER," +
                    "   AUM.MEDIATION_SCRIPT," +
                    "   ARSM.SCOPE_NAME " +
                    " FROM " +
                    "   API_URL_MAPPING AUM " +
                    " INNER JOIN API API ON AUM.API_ID = API.API_ID " +
                    " LEFT OUTER JOIN API_RESOURCE_SCOPE_MAPPING ARSM ON AUM.URL_MAPPING_ID = ARSM.URL_MAPPING_ID" +
                    " WHERE " +
                    "  API.API_UUID = ? AND " +
                    "  AUM.REVISION_UUID = ? " +
                    " ORDER BY AUM.URL_MAPPING_ID ASC ";

    public static final String GET_API_PRODUCT_URI_TEMPLATE_ASSOCIATION_SQL =
            " SELECT " +
            "  API.API_PROVIDER," +
            "  API.API_NAME," +
            "  API.API_VERSION," +
            "  APM.URL_MAPPING_ID  " +
            "  FROM " +
            "  API API " +
            "  INNER JOIN API_PRODUCT_MAPPING APM ON API.API_ID = APM.API_ID " +
            "  WHERE APM.URL_MAPPING_ID IN " +
            "   (SELECT AUM.URL_MAPPING_ID " +
            "   FROM API_URL_MAPPING AUM " +
            "   INNER JOIN API API ON AUM.API_ID = API.API_ID " +
            "   WHERE API.API_UUID = ? AND " +
            "   AUM.REVISION_UUID IS NULL AND APM.REVISION_UUID = 'Current API')";

    public static final String GET_ASSOCIATED_API_PRODUCT_URL_TEMPLATES_SQL =
            " SELECT " +
                    "  API.API_PROVIDER," +
                    "  API.API_NAME," +
                    "  API.API_VERSION," +
                    "  APM.URL_MAPPING_ID  " +
                    "  FROM " +
                    "  API API " +
                    "  INNER JOIN API_PRODUCT_MAPPING APM ON API.API_ID = APM.API_ID " +
                    "  WHERE APM.URL_MAPPING_ID IN " +
                    "   (SELECT AUM.URL_MAPPING_ID " +
                    "   FROM API_URL_MAPPING AUM " +
                    "   INNER JOIN API API ON AUM.API_ID = API.API_ID " +
                    "   WHERE API.API_UUID = ? AND APM.REVISION_UUID = 'Current API')";

    public static final String ADD_COMMENT_SQL =
            " INSERT INTO " +
                    "API_COMMENTS " +
                    "(COMMENT_ID," +
                    "COMMENT_TEXT," +
                    "CREATED_BY," +
                    "CREATED_TIME," +
                    "API_ID," +
                    "PARENT_COMMENT_ID," +
                    "ENTRY_POINT," +
                    "CATEGORY)" +
            " VALUES (?,?,?,?,?,?,?,?)";

    public static final String GET_COMMENT_SQL =
            "SELECT " +
                "API_COMMENTS.COMMENT_ID, " +
                "API_COMMENTS.COMMENT_TEXT, " +
                "API_COMMENTS.CREATED_BY, " +
                "API_COMMENTS.CREATED_TIME, " +
                "API_COMMENTS.UPDATED_TIME, " +
                "API_COMMENTS.API_ID, " +
                "API_COMMENTS.PARENT_COMMENT_ID, " +
                "API_COMMENTS.ENTRY_POINT, " +
                "API_COMMENTS.CATEGORY " +
            "FROM " +
                "API_COMMENTS, " +
                "API API " +
            "WHERE " +
                "API.API_UUID = ? " +
                "AND API.API_ID = API_COMMENTS.API_ID " +
                "AND API_COMMENTS.COMMENT_ID = ?";

    public static final String DELETE_API_CHILD_COMMENTS =
            "DELETE FROM API_COMMENTS WHERE API_ID = ? AND PARENT_COMMENT_ID IS NOT NULL";

    public static final String DELETE_API_PARENT_COMMENTS =
            "DELETE FROM API_COMMENTS WHERE API_ID = ? AND PARENT_COMMENT_ID IS NULL";

    public static final String GET_IDS_OF_REPLIES_SQL =
            "SELECT " +
                "API_COMMENTS.COMMENT_ID " +
            "FROM " +
                "API_COMMENTS, " +
                "API API " +
            "WHERE " +
                "API.API_UUID = ? " +
                "AND API.API_ID = API_COMMENTS.API_ID " +
                "AND PARENT_COMMENT_ID = ?";

    public static final String GET_REPLIES_COUNT_SQL =
            "SELECT " +
                "COUNT(API_COMMENTS.COMMENT_ID) AS COMMENT_COUNT " +
            "FROM " +
                "API_COMMENTS, " +
                "API API " +
            "WHERE " +
                "API.API_UUID = ? " +
                "AND API.API_ID = API_COMMENTS.API_ID " +
                "AND PARENT_COMMENT_ID = ?";

    public static final String GET_ROOT_COMMENTS_COUNT_SQL =
            "SELECT " +
                "COUNT(API_COMMENTS.COMMENT_ID) AS COMMENT_COUNT " +
            "FROM " +
                "API_COMMENTS, " +
                "API API " +
            "WHERE " +
                "API.API_UUID = ? " +
                "AND API.API_ID = API_COMMENTS.API_ID " +
                "AND PARENT_COMMENT_ID IS NULL";

    public static final String EDIT_COMMENT =
            "UPDATE " +
                "API_COMMENTS " +
            "SET " +
                "COMMENT_TEXT = ?, " +
                "UPDATED_TIME = ?, " +
                "CATEGORY = ? " +
            "WHERE " +
                "API_COMMENTS.API_ID = ? " +
                "AND " +
                "API_COMMENTS.COMMENT_ID = ?";

    public static final String DELETE_COMMENT_SQL =
            "DELETE " +
            "FROM " +
                "API_COMMENTS " +
            "WHERE " +
                "API_COMMENTS.API_ID = ? " +
            "AND " +
                "API_COMMENTS.COMMENT_ID = ?";

    public static final String GET_API_CONTEXT_SQL =
            "SELECT CONTEXT FROM API WHERE CONTEXT= ? AND ORGANIZATION = ?";

    public static final String GET_API_IDENTIFIER_BY_UUID_SQL =
            "SELECT API_PROVIDER, API_NAME, API_VERSION FROM API WHERE API_UUID = ?";
    public static final String GET_API_OR_API_PRODUCT_IDENTIFIER_BY_UUID_SQL =
            "SELECT API_PROVIDER, API_NAME, API_VERSION, API_TYPE FROM API WHERE API_UUID = ?";
    public static final String GET_UUID_BY_IDENTIFIER_SQL =
            "SELECT API_UUID FROM API WHERE API_PROVIDER = ? AND API_NAME = ? AND API_VERSION = ?";
    public static final String GET_UUID_BY_IDENTIFIER_AND_ORGANIZATION_SQL = "SELECT API_UUID FROM API"
            + " WHERE API_NAME = ? AND API_VERSION = ? AND ORGANIZATION = ?";
    public static final String GET_API_TYPE_BY_UUID =
            "SELECT API_TYPE FROM API WHERE API_UUID = ?";

    public static final String GET_API_CONTEXT_BY_API_UUID_SQL =
            "SELECT CONTEXT FROM API WHERE API_UUID = ?";

    public static final String GET_ALL_CONTEXT_SQL = "SELECT CONTEXT FROM API ";

    public static final String GET_ALL_CONTEXT_AND_UUID_SQL = "SELECT CONTEXT, API_UUID FROM API ";

    public static final String GET_APPLICATION_REGISTRATION_ENTRY_BY_SUBSCRIBER_SQL =
            "SELECT " +
            "   APP.APPLICATION_ID," +
            "   APP.UUID," +
            "   APP.TOKEN_TYPE AS APP_TYPE," +
            "   APP.NAME," +
            "   APP.SUBSCRIBER_ID," +
            "   APP.APPLICATION_TIER," +
            "   REG.TOKEN_TYPE," +
            "   REG.TOKEN_SCOPE," +
            "   APP.CALLBACK_URL," +
            "   APP.DESCRIPTION," +
            "   APP.APPLICATION_STATUS," +
            "   SUB.USER_ID," +
            "   REG.ALLOWED_DOMAINS," +
            "   REG.VALIDITY_PERIOD," +
            "   REG.INPUTS, REG.KEY_MANAGER" +
            " FROM " +
            "   APPLICATION_REGISTRATION REG," +
            "   APPLICATION APP," +
            "   SUBSCRIBER SUB" +
            " WHERE " +
            "   REG.SUBSCRIBER_ID=SUB.SUBSCRIBER_ID " +
            "   AND REG.APP_ID = APP.APPLICATION_ID " +
            "   AND REG.WF_REF=?";

    public static final String GET_APPLICATION_REGISTRATION_ENTRY_SQL =
            "SELECT " +
            "   REG.TOKEN_TYPE," +
            "   REG.ALLOWED_DOMAINS," +
            "   REG.VALIDITY_PERIOD" +
            " FROM " +
            "   APPLICATION_REGISTRATION REG, " +
            "   APPLICATION APP " +
            " WHERE " +
            "   REG.APP_ID = APP.APPLICATION_ID " +
            "   AND APP.APPLICATION_ID=?";

    public static final String GET_APPLICATION_REGISTRATION_ID_SQL =
            "SELECT APP_ID FROM APPLICATION_REGISTRATION WHERE WF_REF=?";

    public static final String GET_WORKFLOW_ENTRY_SQL =
            "SELECT " +
            "   REG.WF_REF" +
            " FROM " +
            "   APPLICATION APP, " +
            "   APPLICATION_REGISTRATION REG, " +
            "   SUBSCRIBER SUB" +
            " WHERE " +
            "   APP.NAME=? " +
            "   AND SUB.USER_ID=? " +
            "   AND SUB.SUBSCRIBER_ID = APP.SUBSCRIBER_ID " +
            "   AND REG.APP_ID=APP.APPLICATION_ID";

    public static final String GET_WORKFLOW_ENTRY_BY_APP_ID_SQL =
            "SELECT " +
                    "   REG.WF_REF" +
                    " FROM " +
                    "   APPLICATION APP, " +
                    "   APPLICATION_REGISTRATION REG, " +
                    "   SUBSCRIBER SUB" +
                    " WHERE " +
                    "   APP.APPLICATION_ID=? " +
                    "   AND SUB.USER_ID=? " +
                    "   AND SUB.SUBSCRIBER_ID = APP.SUBSCRIBER_ID " +
                    "   AND REG.APP_ID=APP.APPLICATION_ID";

    public static final String GET_EXTERNAL_WORKFLOW_REFERENCE_SQL =
            "SELECT WF_EXTERNAL_REFERENCE FROM WORKFLOWS WHERE WF_TYPE=? AND WF_REFERENCE=?";

    public static final String REMOVE_WORKFLOW_ENTRY_SQL =
            "DELETE FROM WORKFLOWS WHERE WF_TYPE=? AND WF_EXTERNAL_REFERENCE=?";

    public static final String GET_EXTERNAL_WORKFLOW_REFERENCE_FOR_SUBSCRIPTION_SQL =
            "SELECT " +
            "   AW.WF_EXTERNAL_REFERENCE " +
            " FROM" +
            "   WORKFLOWS AW, " +
            "   SUBSCRIPTION ASUB " +
            " WHERE" +
            "   ASUB.API_ID=? " +
            "   AND ASUB.APPLICATION_ID=? " +
            "   AND AW.WF_REFERENCE=ASUB.SUBSCRIPTION_ID " +
            "   AND AW.WF_TYPE=?";

    public static final String GET_EXTERNAL_WORKFLOW_REFERENCE_FOR_SUBSCRIPTION_POSTGRE_SQL =
            "SELECT" +
            "   AW.WF_EXTERNAL_REFERENCE" +
            " FROM" +
            "   WORKFLOWS AW, " +
            "   SUBSCRIPTION ASUB  " +
            " WHERE" +
            "   ASUB.API_ID=? " +
            "   AND ASUB.APPLICATION_ID=?" +
            "   AND AW.WF_REFERENCE::integer=ASUB.SUBSCRIPTION_ID " +
            "   AND AW.WF_TYPE=?";

    public static final String GET_EXTERNAL_WORKFLOW_FOR_SUBSCRIPTION_SQL =
            " SELECT " +
            "   WF_EXTERNAL_REFERENCE" +
            " FROM " +
            "   WORKFLOWS" +
            " WHERE " +
            "   WF_REFERENCE=?" +
            "   AND WF_TYPE=?";

    public static final String GET_EXTERNAL_WORKFLOW_FOR_SIGNUP_SQL =
            "SELECT " +
            "   WF_EXTERNAL_REFERENCE" +
            " FROM " +
            "   WORKFLOWS WHERE " +
            "   WF_REFERENCE=? " +
            "   AND WF_TYPE=?";

    public static final String GET_PAGINATED_SUBSCRIPTIONS_BY_APPLICATION_SQL =
            "SELECT" +
            "   SUBSCRIPTION_ID " +
            " FROM " +
            "   SUBSCRIPTION " +
            " WHERE " +
            "   APPLICATION_ID=? " +
            "   AND SUB_STATUS=?";

    public static final String GET_SUBSCRIPTION_ID_STATUS_BY_APPLICATION_SQL =
            "SELECT" +
                    "   SUBSCRIPTION_ID, SUB_STATUS" +
                    " FROM " +
                    "   SUBSCRIPTION " +
                    " WHERE " +
                    "   APPLICATION_ID=? ";

    public static final String GET_SUBSCRIPTIONS_BY_API_SQL =
            "SELECT" +
                    "   SUBSCRIPTION_ID" +
                    " FROM " +
                    "   SUBSCRIPTION SUBS," +
                    "   API API " +
                    " WHERE " +
                    "   API.API_UUID = ? " +
                    "   AND API.API_ID = SUBS.API_ID " +
                    "   AND SUB_STATUS = ?";

    public static final String GET_REGISTRATION_WORKFLOW_SQL =
            "SELECT WF_REF FROM APPLICATION_REGISTRATION WHERE APP_ID = ? AND TOKEN_TYPE = ? AND KEY_MANAGER = ?";

    public static final String GET_SUBSCRIPTION_STATUS_SQL =
            "SELECT SUB_STATUS FROM SUBSCRIPTION WHERE API_ID = ? AND APPLICATION_ID = ?";
    
    public static final String GET_SUBSCRIPTION_ID_SQL =
            "SELECT SUBSCRIPTION_ID FROM SUBSCRIPTION WHERE API_ID = ? AND APPLICATION_ID = ?";

    public static final String GET_SUBSCRIPTION_CREATION_STATUS_SQL =
            "SELECT SUBS_CREATE_STATE FROM SUBSCRIPTION WHERE API_ID = ? AND APPLICATION_ID = ?";

    public static final String ADD_EXTERNAL_API_STORE_SQL =
            " INSERT INTO EXTERNAL_STORES (API_ID,STORE_ID,STORE_DISPLAY_NAME,STORE_ENDPOINT,STORE_TYPE," +
                    "LAST_UPDATED_TIME) VALUES (?,?,?,?,?,?)";

    public static final String REMOVE_EXTERNAL_API_STORE_SQL =
            "DELETE FROM EXTERNAL_STORES WHERE API_ID=? AND STORE_ID=? AND STORE_TYPE=?";

    public static final String UPDATE_EXTERNAL_API_STORE_SQL =
            "UPDATE " +
            "   EXTERNAL_STORES" +
            " SET " +
            "   STORE_ENDPOINT = ?, " +
            "   STORE_TYPE = ?, " +
            "   LAST_UPDATED_TIME = ? " +
            " WHERE " +
            "   API_ID = ? AND STORE_ID = ? ";

    public static final String GET_EXTERNAL_API_STORE_DETAILS_SQL =
            "SELECT " +
            "   ES.STORE_ID, " +
            "   ES.STORE_DISPLAY_NAME, " +
            "   ES.STORE_ENDPOINT, " +
            "   ES.STORE_TYPE, " +
            "   ES.LAST_UPDATED_TIME " +
             "FROM " +
            "   EXTERNAL_STORES ES " +
            " WHERE " +
            "   ES.API_ID = ? ";

    public static final String ADD_PRODUCT_RESOURCE_MAPPING_SQL = "INSERT INTO API_PRODUCT_MAPPING "
            + "(API_ID,URL_MAPPING_ID,REVISION_UUID) " + "VALUES (?, ?, ?)";

    public static final String DELETE_FROM_API_PRODUCT_MAPPING_SQL = "DELETE FROM API_PRODUCT_MAPPING WHERE "
            + "API_ID = ? AND REVISION_UUID = 'Current API' ";

    public static final String GET_SCOPE_BY_SUBSCRIBED_API_PREFIX =
            "SELECT DISTINCT ARSM.SCOPE_NAME " +
                    "FROM API_RESOURCE_SCOPE_MAPPING ARSM INNER JOIN API_URL_MAPPING AUM " +
                    "ON ARSM.URL_MAPPING_ID = AUM.URL_MAPPING_ID " +
                    "WHERE AUM.REVISION_UUID IS NULL AND AUM.API_ID IN (";

    public static final char GET_SCOPE_BY_SUBSCRIBED_ID_SUFFIX = ')';

    public static final String GET_SCOPE_BY_SUBSCRIBED_ID_ORACLE_SQL =
            "SELECT DISTINCT ARSM.SCOPE_NAME " +
                    "FROM API_RESOURCE_SCOPE_MAPPING ARSM INNER JOIN API_URL_MAPPING AUM " +
                    "ON ARSM.URL_MAPPING_ID = AUM.URL_MAPPING_ID " +
                    "WHERE AUM.REVISION_UUID IS NULL AND AUM.API_ID IN (";

    public static final String GET_RESOURCE_TO_SCOPE_MAPPING_SQL =
            "SELECT AUM.URL_MAPPING_ID, ARSM.SCOPE_NAME FROM API_URL_MAPPING AUM " +
                    "LEFT JOIN API_RESOURCE_SCOPE_MAPPING ARSM ON AUM.URL_MAPPING_ID = ARSM.URL_MAPPING_ID " +
                    "WHERE AUM.API_ID = ? AND AUM.REVISION_UUID IS NULL";

    public static final String REMOVE_SUBSCRIPTION_BY_APPLICATION_ID_SQL =
            "DELETE FROM SUBSCRIPTION WHERE API_ID = ? AND APPLICATION_ID = ? ";

    public static final String GET_API_NAME_NOT_MATCHING_CONTEXT_SQL =
            "SELECT COUNT(API_ID) AS API_COUNT FROM API WHERE LOWER(API_NAME) = LOWER(?) AND ORGANIZATION = ? AND CONTEXT NOT LIKE ?";

    public static final String GET_API_NAME_MATCHING_CONTEXT_SQL =
            "SELECT COUNT(API_ID) AS API_COUNT FROM API WHERE LOWER(API_NAME) = LOWER(?) AND AND ORGANIZATION = ? CONTEXT LIKE ?";

    public static final String GET_API_NAME_DIFF_CASE_NOT_MATCHING_CONTEXT_SQL =
            "SELECT COUNT(API_ID) AS API_COUNT FROM API WHERE LOWER(API_NAME) = LOWER(?) AND CONTEXT NOT LIKE ? "
        + "AND NOT (API_NAME = ?) AND ORGANIZATION = ?";

    public static final String GET_API_NAME_DIFF_CASE_MATCHING_CONTEXT_SQL =
            "SELECT COUNT(API_ID) AS API_COUNT FROM API WHERE LOWER(API_NAME) = LOWER(?) AND CONTEXT LIKE ? " +
                    "AND NOT (API_NAME = ?) AND ORGANIZATION = ?";

    public static final String GET_CONTEXT_TEMPLATE_COUNT_SQL_MATCHES_ORGANIZATION =
            "SELECT COUNT(CONTEXT_TEMPLATE) AS CTX_COUNT FROM API WHERE LOWER(CONTEXT_TEMPLATE) = ? " +
                    "AND ORGANIZATION = ?";

    public static final String GET_API_NAMES_MATCHES_CONTEXT=
            "SELECT DISTINCT API_NAME FROM API WHERE CONTEXT_TEMPLATE = ?";

    public static final String GET_VERSIONS_MATCHES_CONTEXT=
            "SELECT API_VERSION FROM API WHERE CONTEXT_TEMPLATE = ? AND API_NAME = ?";

    public static final String GET_APPLICATION_MAPPING_FOR_CONSUMER_KEY_SQL =
            "SELECT APPLICATION_ID FROM APPLICATION_KEY_MAPPING WHERE CONSUMER_KEY = ? AND KEY_MANAGER = ?";

    public static final String IS_KEY_MAPPING_EXISTS_FOR_APP_ID_KEY_TYPE_OR_CONSUMER_KEY =
            "SELECT 1 FROM APPLICATION_KEY_MAPPING WHERE " +
                    "((APPLICATION_ID = ? AND KEY_TYPE = ?) OR (CONSUMER_KEY = ?)) AND KEY_MANAGER IN (?,?)";

    public static final String IS_KEY_MAPPING_EXISTS_FOR_APP_ID_KEY_TYPE =
            "SELECT 1 FROM APPLICATION_KEY_MAPPING WHERE APPLICATION_ID = ? AND KEY_TYPE = ? " +
                    "AND KEY_MANAGER IN (?,?)";

    public static final String GET_CONSUMER_KEY_BY_APPLICATION_AND_KEY_SQL =
            " SELECT " +
            "   CONSUMER_KEY,KEY_MANAGER " +
            " FROM " +
            "   APPLICATION_KEY_MAPPING " +
            " WHERE " +
            "   APPLICATION_ID = ? " +
            "   AND KEY_TYPE = ? ";

    public static final String GET_LAST_PUBLISHED_API_VERSION_SQL =
            "SELECT " +
            "   API.API_VERSION " +
            " FROM " +
            "   API API , " +
            "   EXTERNAL_STORES ES " +
            " WHERE " +
            "   ES.API_ID = API.API_ID " +
            "   AND API.API_PROVIDER = ? " +
            "   AND API.API_NAME=? " +
            "   AND ES.STORE_ID =? " +
            " ORDER By API.CREATED_TIME ASC";

    public static final String GET_ALL_ALERT_TYPES =
            "SELECT " +
            "   AT.ALERT_TYPE_ID, " +
            "   AT.ALERT_TYPE_NAME " +
            " FROM " +
            "   ALERT_TYPES AT  " +
            " WHERE " +
            "   STAKE_HOLDER   = ?";


    public static final String GET_ALL_ALERT_TYPES_FOR_ADMIN =
            "SELECT DISTINCT" +
            "   AT.ALERT_TYPE_ID, " +
            "   AT.ALERT_TYPE_NAME " +
            " FROM " +
            "   ALERT_TYPES AT  ";

    public static final String GET_SAVED_ALERT_TYPES_BY_USERNAME =
            " SELECT " +
            "   ALERT_TYPE_ID " +
            " FROM " +
            "   ALERT_TYPES_VALUES " +
            " WHERE " +
            "   USER_NAME = ? " +
            "   AND STAKE_HOLDER   = ? ";

    public static final String GET_SAVED_ALERT_EMAILS =

            " SELECT " +
            "   EMAIL " +
            " FROM " +
            "   ALERT_EMAILLIST , " +
            "   ALERT_EMAILLIST_DETAILS  " +
            " WHERE " +
            "   ALERT_EMAILLIST.EMAIL_LIST_ID = ALERT_EMAILLIST_DETAILS.EMAIL_LIST_ID" +
            "   AND USER_NAME = ? " +
            "   AND STAKE_HOLDER  = ? ";


    public static final String ADD_ALERT_TYPES_VALUES =
            " INSERT INTO ALERT_TYPES_VALUES (ALERT_TYPE_ID, USER_NAME , STAKE_HOLDER) " +
            " VALUES(?,?,?)";

    public static final String ADD_ALERT_EMAIL_LIST =
            " INSERT INTO ALERT_EMAILLIST  (USER_NAME, STAKE_HOLDER) " +
            " VALUES(?,?)";

    public static final String DELETE_ALERTTYPES_BY_USERNAME_AND_STAKE_HOLDER  =
            "DELETE FROM ALERT_TYPES_VALUES WHERE USER_NAME = ? AND STAKE_HOLDER = ?";

    public static final String DELETE_EMAILLIST_BY_EMAIL_LIST_ID  =
            "DELETE FROM ALERT_EMAILLIST_DETAILS   WHERE EMAIL_LIST_ID= ? ";

    public static final String GET_EMAILLISTID_BY_USERNAME_AND_STAKEHOLDER =
            " SELECT " +
            "   EMAIL_LIST_ID " +
            " FROM " +
            "   ALERT_EMAILLIST " +
            " WHERE " +
            "   USER_NAME = ? " +
            "   AND STAKE_HOLDER  = ? ";

    public static final String SAVE_EMAIL_LIST_DETAILS_QUERY =
            " INSERT INTO ALERT_EMAILLIST_DETAILS  (EMAIL_LIST_ID, EMAIL) " +
            " VALUES(?,?)";

    public static final String DELETE_ALERTTYPES_EMAILLISTS_BY_USERNAME_AND_STAKE_HOLDER  =
            "DELETE FROM ALERT_EMAILLIST  WHERE USER_NAME = ? AND STAKE_HOLDER = ?";



    public static final String INSERT_APPLICATION_POLICY_SQL =
            "INSERT INTO APPLICATION_USAGE_PLAN (NAME, DISPLAY_NAME, ORGANIZATION, DESCRIPTION, QUOTA_TYPE, QUOTA, \n" +
                    " QUOTA_UNIT, UNIT_TIME, TIME_UNIT, IS_DEPLOYED, UUID) \n" +
            "VALUES (?,?,?,?,?,?,?,?,?,?,?)";

    public static final String INSERT_APPLICATION_POLICY_WITH_CUSTOM_ATTRIB_SQL =
            "INSERT INTO APPLICATION_USAGE_PLAN (NAME, DISPLAY_NAME, ORGANIZATION, DESCRIPTION, QUOTA_TYPE, QUOTA, \n" +
                    " QUOTA_UNIT, UNIT_TIME, TIME_UNIT, IS_DEPLOYED, UUID,CUSTOM_ATTRIBUTES) \n" +
            "VALUES (?,?,?,?,?,?,?,?,?,?,?,?)";

    public static final String INSERT_SUBSCRIPTION_POLICY_SQL =
            "INSERT INTO BUSINESS_PLAN (NAME, DISPLAY_NAME, ORGANIZATION, DESCRIPTION, QUOTA_TYPE, QUOTA, \n" +
                    " QUOTA_UNIT, UNIT_TIME, TIME_UNIT, IS_DEPLOYED, UUID, RATE_LIMIT_COUNT, \n" +
                    " RATE_LIMIT_TIME_UNIT,STOP_ON_QUOTA_REACH, MAX_DEPTH, MAX_COMPLEXITY, \n" +
                    " BILLING_PLAN,MONETIZATION_PLAN,FIXED_RATE,BILLING_CYCLE,PRICE_PER_REQUEST,CURRENCY, CONNECTIONS_COUNT) \n" +
                    " VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)";

    public static final String INSERT_SUBSCRIPTION_POLICY_WITH_CUSTOM_ATTRIB_SQL =
            "INSERT INTO BUSINESS_PLAN (NAME, DISPLAY_NAME, ORGANIZATION, DESCRIPTION, QUOTA_TYPE, QUOTA, \n" +
                    " QUOTA_UNIT, UNIT_TIME, TIME_UNIT, IS_DEPLOYED, UUID,  RATE_LIMIT_COUNT, \n" +
                    " RATE_LIMIT_TIME_UNIT, STOP_ON_QUOTA_REACH, MAX_DEPTH, MAX_COMPLEXITY, \n" +
                    " BILLING_PLAN, CUSTOM_ATTRIBUTES, MONETIZATION_PLAN, \n" +
                    " FIXED_RATE, BILLING_CYCLE, PRICE_PER_REQUEST, CURRENCY, CONNECTIONS_COUNT) \n" +
                    " VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)";


    public static final String INSERT_GLOBAL_POLICY_SQL =
            "INSERT INTO POLICY_GLOBAL (NAME ,ORGANIZATION, KEY_TEMPLATE, DESCRIPTION ,SIDDHI_QUERY, "
                    + "IS_DEPLOYED, UUID) \n" +
            "VALUES (?,?,?,?,?,?,?)";

    public static final String GET_APP_POLICY_NAMES =
            " SELECT " +
                    "   NAME " +
                    "FROM " +
                    "   APPLICATION_USAGE_PLAN " +
                    " WHERE" +
                    "   ORGANIZATION =?";

    public static final String GET_SUB_POLICY_NAMES =
            " SELECT " +
                    "   NAME " +
                    "FROM " +
                    "   BUSINESS_PLAN " +
                    " WHERE" +
                    "   ORGANIZATION =?";

    public static final String GET_GLOBAL_POLICY_NAMES =
            " SELECT " +
                    "   NAME " +
                    "FROM " +
                    "   POLICY_GLOBAL " +
                    " WHERE" +
                    "   ORGANIZATION =?";

    public static final String GET_GLOBAL_POLICY_KEY_TEMPLATES =
            " SELECT " +
                    "   KEY_TEMPLATE " +
                    "FROM " +
                    "   POLICY_GLOBAL " +
                    " WHERE" +
                    "   ORGANIZATION =?";

    public static final String GET_GLOBAL_POLICY_KEY_TEMPLATE =
            " SELECT " +
                    "   KEY_TEMPLATE " +
                    "FROM " +
                    "   POLICY_GLOBAL " +
                    " WHERE" +
                    "   ORGANIZATION =? AND" +
                    "   KEY_TEMPLATE =? AND" +
                    "   NAME =?";

    public static final String GET_APP_POLICIES =
            " SELECT "+
                    "   * " +
                    "FROM " +
                    "   APPLICATION_USAGE_PLAN " +
                    " WHERE" +
                    "   ORGANIZATION =?";
    public static final String GET_SUBSCRIPTION_POLICIES =
            " SELECT " +
                    "   * " +
                    "FROM " +
                    "   BUSINESS_PLAN " +
                    " WHERE" +
                    "   ORGANIZATION =?";

    public static final String GET_SUBSCRIPTION_POLICIES_BY_POLICY_NAMES_PREFIX =
            " SELECT " +
                    "   * " +
                    "FROM " +
                    "   BUSINESS_PLAN " +
                    " WHERE" +
                    "  NAME IN (";

    public static final String GET_SUBSCRIPTION_POLICIES_BY_POLICY_NAMES_SUFFIX =
            ") AND TENANT_ID =?";

    public static final String GET_GLOBAL_POLICIES =
            " SELECT " +
                    "   * " +
                    "FROM " +
                    "   POLICY_GLOBAL " +
                    " WHERE" +
                    "   ORGANIZATION =?";

    public static final String GET_GLOBAL_POLICY =
            " SELECT " +
            "   * " +
            "FROM " +
            "   POLICY_GLOBAL " +
            " WHERE" +
            "   NAME =?";

    public static final String GET_GLOBAL_POLICY_BY_UUID =
            "SELECT " +
            "   * " +
            "FROM " +
            "   POLICY_GLOBAL " +
            "WHERE" +
            "   UUID =?";

    public static final String GET_APPLICATION_POLICY_SQL =
            "SELECT "+
                    "* " +
            "FROM " +
                    "APPLICATION_USAGE_PLAN " +
            "WHERE " +
                    "NAME = ? AND " +
                    "ORGANIZATION =?";

    public static final String GET_APPLICATION_POLICY_BY_UUID_SQL =
            "SELECT " +
                "* " +
            "FROM " +
                "APPLICATION_USAGE_PLAN " +
            "WHERE " +
                "UUID = ?";

    public static final String GET_SUBSCRIPTION_POLICY_SQL =
            "SELECT "+
                    "* " +
                    "FROM " +
                    "   BUSINESS_PLAN " +
            "WHERE " +
                    "NAME = ? AND " +
                    "ORGANIZATION =?";

    public static final String GET_API_PROVIDER_WITH_NAME_VERSION_FOR_SUPER_TENANT =
            "SELECT API.API_PROVIDER FROM API API WHERE API.API_NAME = ? AND API.API_VERSION = ? AND "
                    + "CONTEXT NOT LIKE '%" + APIConstants.TENANT_PREFIX + "%' ";

    public static final String GET_API_PROVIDER_WITH_NAME_VERSION_FOR_GIVEN_TENANT =
            "SELECT API.API_PROVIDER FROM API API WHERE API.API_NAME = ? AND "
                    + "API.API_VERSION = ? AND API.CONTEXT LIKE ? ";

    public static final String GET_SUBSCRIPTION_POLICY_BY_UUID_SQL =
            "SELECT "+
                    "* " +
                    "FROM " +
                    "   BUSINESS_PLAN " +
            "WHERE " +
                    "UUID = ?";

    public static final String UPDATE_APPLICATION_POLICY_SQL =
            "UPDATE APPLICATION_USAGE_PLAN " +
            "SET " +
                    "DISPLAY_NAME = ?, " +
                    "DESCRIPTION = ?, " +
                    "QUOTA_TYPE = ?, " +
                    "QUOTA = ?, " +
                    "QUOTA_UNIT = ?, " +
                    "UNIT_TIME = ?, " +
                    "TIME_UNIT = ? " +
            "WHERE NAME = ? AND ORGANIZATION = ?";

    public static final String UPDATE_APPLICATION_POLICY_WITH_CUSTOM_ATTRIBUTES_SQL =
            "UPDATE APPLICATION_USAGE_PLAN " +
            "SET " +
                    "DISPLAY_NAME = ?, " +
                    "DESCRIPTION = ?, " +
                    "QUOTA_TYPE = ?, " +
                    "QUOTA = ?, " +
                    "QUOTA_UNIT = ?, " +
                    "UNIT_TIME = ?, " +
                    "TIME_UNIT = ?, " +
                    " CUSTOM_ATTRIBUTES = ? "+
            "WHERE NAME = ? AND ORGANIZATION = ?";

    public static final String UPDATE_APPLICATION_POLICY_BY_UUID_SQL =
            "UPDATE APPLICATION_USAGE_PLAN " +
                    "SET " +
                    "DISPLAY_NAME = ?, " +
                    "DESCRIPTION = ?, " +
                    "QUOTA_TYPE = ?, " +
                    "QUOTA = ?, " +
                    "QUOTA_UNIT = ?, " +
                    "UNIT_TIME = ?, " +
                    "TIME_UNIT = ? " +
                    "WHERE UUID = ?";

    public static final String UPDATE_APPLICATION_POLICY_WITH_CUSTOM_ATTRIBUTES_BY_UUID_SQL =
            "UPDATE APPLICATION_USAGE_PLAN " +
                    "SET " +
                    "DISPLAY_NAME = ?, " +
                    "DESCRIPTION = ?, " +
                    "QUOTA_TYPE = ?, " +
                    "QUOTA = ?, " +
                    "QUOTA_UNIT = ?, " +
                    "UNIT_TIME = ?, " +
                    "TIME_UNIT = ?, " +
                    "CUSTOM_ATTRIBUTES = ? "+
                    "WHERE UUID = ?";

    public static final String UPDATE_SUBSCRIPTION_POLICY_SQL =
            "UPDATE BUSINESS_PLAN " +
            "SET " +
                    "DISPLAY_NAME = ?, " +
                    "DESCRIPTION = ?, " +
                    "QUOTA_TYPE = ?, " +
                    "QUOTA = ?, " +
                    "QUOTA_UNIT = ?, " +
                    "UNIT_TIME = ?, " +
                    "TIME_UNIT = ?, " +
                    "RATE_LIMIT_COUNT = ?," +
                    "RATE_LIMIT_TIME_UNIT = ?, " +
                    "STOP_ON_QUOTA_REACH = ?, " +
                    "MAX_DEPTH = ?, " +
                    "MAX_COMPLEXITY = ?, " +
                    "BILLING_PLAN = ?, " +
                    "MONETIZATION_PLAN = ?," +
                    "FIXED_RATE = ?," +
                    "BILLING_CYCLE = ?," +
                    "PRICE_PER_REQUEST = ?, " +
                    "CURRENCY = ?, " +
                    "CONNECTIONS_COUNT = ?" +
            "WHERE NAME = ? AND ORGANIZATION = ?";

    public static final String UPDATE_SUBSCRIPTION_POLICY_WITH_CUSTOM_ATTRIBUTES_SQL =
            "UPDATE BUSINESS_PLAN " +
            "SET " +
                    "DISPLAY_NAME = ?, " +
                    "DESCRIPTION = ?, " +
                    "QUOTA_TYPE = ?, " +
                    "QUOTA = ?, " +
                    "QUOTA_UNIT = ?, " +
                    "UNIT_TIME = ?, " +
                    "TIME_UNIT = ?, " +
                    "RATE_LIMIT_COUNT = ?," +
                    "RATE_LIMIT_TIME_UNIT = ?, " +
                    "STOP_ON_QUOTA_REACH = ?, " +
                    "MAX_DEPTH = ?, " +
                    "MAX_COMPLEXITY = ?, " +
                    "BILLING_PLAN = ?, "+
                    "CUSTOM_ATTRIBUTES = ?, "+
                    "MONETIZATION_PLAN = ?," +
                    "FIXED_RATE = ?," +
                    "BILLING_CYCLE = ?," +
                    "PRICE_PER_REQUEST = ?, " +
                    "CURRENCY = ?, " +
                    "CONNECTIONS_COUNT = ? " +
            "WHERE NAME = ? AND ORGANIZATION = ?";

    public static final String UPDATE_SUBSCRIPTION_POLICY_BY_UUID_SQL =
            "UPDATE BUSINESS_PLAN " +
                    "SET " +
                    "DISPLAY_NAME = ?, " +
                    "DESCRIPTION = ?, " +
                    "QUOTA_TYPE = ?, " +
                    "QUOTA = ?, " +
                    "QUOTA_UNIT = ?, " +
                    "UNIT_TIME = ?, " +
                    "TIME_UNIT = ?, " +
                    "RATE_LIMIT_COUNT = ?," +
                    "RATE_LIMIT_TIME_UNIT = ?, " +
                    "STOP_ON_QUOTA_REACH = ?, " +
                    "MAX_DEPTH = ?, " +
                    "MAX_COMPLEXITY = ?, " +
                    "BILLING_PLAN = ?, "+
                    "MONETIZATION_PLAN = ?," +
                    "FIXED_RATE = ?," +
                    "BILLING_CYCLE = ?," +
                    "PRICE_PER_REQUEST = ?, " +
                    "CURRENCY = ?, " +
                    "CONNECTIONS_COUNT = ? " +
                    "WHERE UUID = ?";

    public static final String UPDATE_SUBSCRIPTION_POLICY_WITH_CUSTOM_ATTRIBUTES_BY_UUID_SQL =
            "UPDATE BUSINESS_PLAN " +
                    "SET " +
                    "DISPLAY_NAME = ?, " +
                    "DESCRIPTION = ?, " +
                    "QUOTA_TYPE = ?, " +
                    "QUOTA = ?, " +
                    "QUOTA_UNIT = ?, " +
                    "UNIT_TIME = ?, " +
                    "TIME_UNIT = ?, " +
                    "RATE_LIMIT_COUNT = ?," +
                    "RATE_LIMIT_TIME_UNIT = ?, " +
                    "STOP_ON_QUOTA_REACH = ?, " +
                    "MAX_DEPTH = ?, " +
                    "MAX_COMPLEXITY = ?, " +
                    "BILLING_PLAN = ?, "+
                    "CUSTOM_ATTRIBUTES = ?, "+
                    "MONETIZATION_PLAN = ?," +
                    "FIXED_RATE = ?," +
                    "BILLING_CYCLE = ?," +
                    "PRICE_PER_REQUEST = ?, " +
                    "CURRENCY = ?, " +
                    "CONNECTIONS_COUNT = ? " +
                    "WHERE UUID = ?";

    public static final String UPDATE_GLOBAL_POLICY_SQL =
            "UPDATE POLICY_GLOBAL " +
            "SET " +
                    "DESCRIPTION = ?, " +
                    "SIDDHI_QUERY = ?, " +
                    "KEY_TEMPLATE = ? " +
            "WHERE NAME = ? AND ORGANIZATION = ?";

    public static final String UPDATE_GLOBAL_POLICY_BY_UUID_SQL =
            "UPDATE POLICY_GLOBAL " +
                    "SET " +
                    "DESCRIPTION = ?, " +
                    "SIDDHI_QUERY = ?, " +
                    "KEY_TEMPLATE = ? " +
                    "WHERE UUID = ?";

    public static final String UPDATE_APPLICATION_POLICY_STATUS_SQL =
            "UPDATE APPLICATION_USAGE_PLAN SET IS_DEPLOYED = ? WHERE NAME = ? AND ORGANIZATION = ?";

    public static final String UPDATE_SUBSCRIPTION_POLICY_STATUS_SQL =
            "UPDATE BUSINESS_PLAN SET IS_DEPLOYED = ? WHERE NAME = ? AND ORGANIZATION = ?";

    public static final String UPDATE_GLOBAL_POLICY_STATUS_SQL =
            "UPDATE POLICY_GLOBAL SET IS_DEPLOYED = ? WHERE NAME = ? AND ORGANIZATION = ?";

    public static final String DELETE_APPLICATION_POLICY_SQL =
            "DELETE FROM APPLICATION_USAGE_PLAN WHERE ORGANIZATION = ? AND NAME = ?";

    public static final String DELETE_SUBSCRIPTION_POLICY_SQL =
            "DELETE FROM BUSINESS_PLAN WHERE ORGANIZATION = ? AND NAME = ?";

    public static final String DELETE_GLOBAL_POLICY_SQL =
            "DELETE FROM POLICY_GLOBAL WHERE ORGANIZATION = ? AND NAME = ?";

    public static final String GET_API_DETAILS_SQL = "SELECT * FROM API ";


    public static final String REMOVE_GROUP_ID_MAPPING_SQL =
            "DELETE FROM APPLICATION_GROUP_MAPPING WHERE APPLICATION_ID = ? ";

    public static final String ADD_GROUP_ID_MAPPING_SQL =
            "INSERT INTO APPLICATION_GROUP_MAPPING (APPLICATION_ID, GROUP_ID, TENANT) VALUES (?,?,?)";

    public static final String GET_GROUP_ID_SQL =
            "SELECT GROUP_ID  FROM APPLICATION_GROUP_MAPPING WHERE APPLICATION_ID = ?";

    public static final String REMOVE_MIGRATED_GROUP_ID_SQL =
            "UPDATE APPLICATION SET GROUP_ID = '' WHERE APPLICATION_ID = ?";


    /** Environment related constants **/

    public static final String GET_ENVIRONMENT_BY_ORGANIZATION_SQL =
            "SELECT ID, UUID, NAME, ORGANIZATION, DISPLAY_NAME, DESCRIPTION, PROVIDER " +
            "FROM GATEWAY_ENVIRONMENT " +
            "WHERE ORGANIZATION = ?";

    public static final String GET_ENVIRONMENT_BY_ORGANIZATION_AND_UUID_SQL =
            "SELECT ID, UUID, NAME, ORGANIZATION, DISPLAY_NAME, DESCRIPTION, PROVIDER " +
            "FROM GATEWAY_ENVIRONMENT " +
            "WHERE ORGANIZATION = ? AND UUID = ?";

    public static final String INSERT_ENVIRONMENT_SQL = "INSERT INTO " +
            "GATEWAY_ENVIRONMENT (UUID, NAME, DISPLAY_NAME, DESCRIPTION, PROVIDER, ORGANIZATION) " +
            "VALUES (?,?,?,?,?,?)";

    public static final String INSERT_GATEWAY_VHOSTS_SQL = "INSERT INTO " +
            "GW_VHOST (GATEWAY_ENV_ID, HOST, HTTP_CONTEXT, HTTP_PORT, HTTPS_PORT, WS_PORT, WSS_PORT) " +
            "VALUES (?,?,?,?,?,?,?)";

    public static final String DELETE_GATEWAY_VHOSTS_SQL = "DELETE FROM GW_VHOST WHERE GATEWAY_ENV_ID = ?";

    public static final String GET_ENVIRONMENT_VHOSTS_BY_ID_SQL =
            "SELECT GATEWAY_ENV_ID, HOST, HTTP_CONTEXT, HTTP_PORT, HTTPS_PORT, WS_PORT, WSS_PORT " +
            "FROM GW_VHOST WHERE GATEWAY_ENV_ID = ?";

    public static final String DELETE_ENVIRONMENT_SQL = "DELETE FROM GATEWAY_ENVIRONMENT WHERE UUID = ?";

    public static final String UPDATE_ENVIRONMENT_SQL = "UPDATE GATEWAY_ENVIRONMENT " +
            "SET DISPLAY_NAME = ?, DESCRIPTION = ? " +
            "WHERE UUID = ?";

    public static final String DELETE_API_PRODUCT_SQL =
            "DELETE FROM API WHERE API_PROVIDER = ? AND API_NAME = ? AND API_VERSION = ? AND API_TYPE = '"
                    + APIConstants.API_PRODUCT + "'";

    public static final String UPDATE_PRODUCT_SQL =
            " UPDATE API " +
            " SET" +
            "   API_TIER=?," +
            "   UPDATED_BY=?," +
            "   UPDATED_TIME=?," +
            "   GATEWAY_VENDOR=?" +
            " WHERE" +
            "   API_NAME=? AND API_PROVIDER=? AND API_VERSION=? AND API_TYPE='" + APIConstants.API_PRODUCT +"'";

    public static final String GET_PRODUCT_ID =
            "SELECT API_ID FROM API WHERE API_NAME = ? AND API_PROVIDER = ? AND "
            + "API_VERSION = ? AND API_TYPE='" + APIConstants.API_PRODUCT +"'";

    public static final String GET_URL_TEMPLATES_FOR_API =
            "SELECT URL_PATTERN , URL_MAPPING_ID, HTTP_METHOD FROM API API , API_URL_MAPPING URL "
            + "WHERE API.API_ID = URL.API_ID AND API.API_NAME =? "
            + "AND API.API_VERSION=? AND API.API_PROVIDER=? AND URL.REVISION_UUID IS NULL";

    public static final String GET_URL_TEMPLATES_FOR_API_WITH_UUID =
            "SELECT URL_PATTERN , URL_MAPPING_ID, HTTP_METHOD FROM API API , API_URL_MAPPING URL "
                    + "WHERE API.API_ID = URL.API_ID AND API.API_UUID =? AND URL.REVISION_UUID IS NULL";

    public static final String ADD_API_PRODUCT =
            "INSERT INTO "
            + "API(API_PROVIDER, API_NAME, API_VERSION, CONTEXT,"
            + "API_TIER, CREATED_BY, CREATED_TIME, API_TYPE, API_UUID, STATUS, ORGANIZATION, GATEWAY_VENDOR, VERSION_COMPARABLE) VALUES (?,?,?,?,?,?,?,?,?"
                    + ",?,?,?,?)";

    public static final String GET_RESOURCES_OF_PRODUCT =
            "SELECT API_UM.URL_MAPPING_ID, API_UM.URL_PATTERN, API_UM.HTTP_METHOD, API_UM.AUTH_SCHEME, " +
                "API_UM.THROTTLING_TIER, API.API_PROVIDER, API.API_NAME, API.API_VERSION, API.CONTEXT, API.API_UUID " +
            "FROM API_URL_MAPPING API_UM " +
            "INNER JOIN API API " +
                "ON API.API_ID = API_UM.API_ID " +
            "INNER JOIN API_PRODUCT_MAPPING PROD_MAP " +
                "ON PROD_MAP.URL_MAPPING_ID = API_UM.URL_MAPPING_ID " +
            "WHERE PROD_MAP.API_ID = ? AND API_UM.REVISION_UUID = ? AND PROD_MAP.REVISION_UUID = 'Current API' ";

    public static final String GET_RESOURCES_OF_PRODUCT_REVISION =
            "SELECT API_UM.URL_MAPPING_ID, API_UM.URL_PATTERN, API_UM.HTTP_METHOD, API_UM.AUTH_SCHEME, " +
                    "API_UM.THROTTLING_TIER, API.API_PROVIDER, API.API_NAME, API.API_VERSION, API.CONTEXT, API.API_UUID " +
                    "FROM API_URL_MAPPING API_UM " +
                    "INNER JOIN API API " +
                    "ON API.API_ID = API_UM.API_ID " +
                    "INNER JOIN API_PRODUCT_MAPPING PROD_MAP " +
                    "ON PROD_MAP.URL_MAPPING_ID = API_UM.URL_MAPPING_ID " +
                    "WHERE PROD_MAP.API_ID = ? AND API_UM.REVISION_UUID = ? AND PROD_MAP.REVISION_UUID = ? ";

    public static final String GET_WH_TOPIC_SUBSCRIPTIONS =
            "SELECT SUBSCRIPTION.HUB_TOPIC, SUBSCRIPTION.API_UUID, " +
            "SUBSCRIPTION.DELIVERED_AT, SUBSCRIPTION.DELIVERY_STATE, SUBSCRIPTION.HUB_CALLBACK_URL, APPLICATION.APPLICATION_ID " +
            "FROM WEBHOOKS_SUBSCRIPTION SUBSCRIPTION INNER JOIN APPLICATION APPLICATION " +
            "ON SUBSCRIPTION.APPLICATION_ID = APPLICATION.APPLICATION_ID WHERE APPLICATION.UUID = ?";

    public static final String GET_WH_TOPIC_SUBSCRIPTIONS_BY_API_KEY =
            "SELECT SUBSCRIPTION.HUB_TOPIC, SUBSCRIPTION.API_UUID, " +
            "SUBSCRIPTION.DELIVERED_AT, SUBSCRIPTION.DELIVERY_STATE, SUBSCRIPTION.HUB_CALLBACK_URL, APPLICATION.APPLICATION_ID " +
            "FROM WEBHOOKS_SUBSCRIPTION SUBSCRIPTION INNER JOIN APPLICATION APPLICATION " +
            "ON SUBSCRIPTION.APPLICATION_ID = APPLICATION.APPLICATION_ID WHERE APPLICATION.UUID = ? " +
            "AND SUBSCRIPTION.API_UUID = ?";

    public static final String GET_ALL_TOPICS_BY_API_ID =
            "SELECT DISTINCT URL.URL_PATTERN, URL.HTTP_METHOD, API.API_ID FROM API_URL_MAPPING URL INNER JOIN " +
            "API API ON URL.API_ID = API.API_ID WHERE API.API_UUID = ? AND URL.REVISION_UUID IS NULL ORDER BY URL.URL_PATTERN ASC";

    public static final String GET_SCOPE_KEYS_BY_URL_MAPPING_ID =
            "SELECT SCOPE_NAME FROM API_RESOURCE_SCOPE_MAPPING WHERE URL_MAPPING_ID = ?" ;

    /** API Categories related constants **/

    public static final String ADD_CATEGORY_SQL = "INSERT INTO API_CATEGORIES "
            + "(UUID, NAME, DESCRIPTION, ORGANIZATION) VALUES (?,?,?,?)";

    public static final String GET_CATEGORIES_BY_ORGANIZATION_SQL = "SELECT UUID, NAME, DESCRIPTION FROM API_CATEGORIES "
            + "WHERE ORGANIZATION = ? ORDER BY NAME";

    public static final String IS_API_CATEGORY_NAME_EXISTS = "SELECT COUNT(UUID) AS API_CATEGORY_COUNT FROM "
            + "API_CATEGORIES WHERE LOWER(NAME) = LOWER(?) AND ORGANIZATION = ?";

    public static final String IS_API_CATEGORY_NAME_EXISTS_FOR_ANOTHER_UUID = "SELECT COUNT(UUID) AS API_CATEGORY_COUNT FROM "
            + "API_CATEGORIES WHERE LOWER(NAME) = LOWER(?) AND ORGANIZATION = ? AND UUID != ?";

    public static final String GET_API_CATEGORY_BY_ID = "SELECT * FROM API_CATEGORIES WHERE UUID = ?";

    public static final String UPDATE_API_CATEGORY = "UPDATE API_CATEGORIES SET DESCRIPTION = ?, NAME = ?, ORGANIZATION = ? WHERE UUID = ?";

    public static final String DELETE_API_CATEGORY = "DELETE FROM API_CATEGORIES WHERE UUID = ?";

    public static final String GET_USER_ID = "SELECT USER_ID FROM USER WHERE USER_NAME=?";

    public static final String ADD_USER_ID = "INSERT INTO USER (USER_ID, USER_NAME) VALUES (?,?)";
    public static final String GET_KEY_MAPPING_ID_FROM_APPLICATION =
            "SELECT UUID FROM APPLICATION_KEY_MAPPING WHERE APPLICATION_ID = ? AND KEY_TYPE = ? AND KEY_MANAGER = ?";
    public static final String GET_CONSUMER_KEY_FOR_APPLICATION_KEY_TYPE_APP_ID_KEY_MANAGER_SQL =
            "SELECT CONSUMER_KEY FROM APPLICATION_KEY_MAPPING WHERE APPLICATION_ID = ? AND KEY_TYPE = ? AND " +
                    "KEY_MANAGER = ?";
    public static final String GET_KEY_MAPPING_INFO_FROM_APP_ID = "SELECT UUID,CONSUMER_KEY,KEY_MANAGER,KEY_TYPE," +
            "STATE,APP_INFO,CREATE_MODE FROM APPLICATION_KEY_MAPPING WHERE APPLICATION_ID = ?";

    public static final String ADD_GW_PUBLISHED_API_DETAILS = "INSERT INTO GW_PUBLISHED_API_DETAILS (API_ID, " +
            "API_NAME, API_VERSION, TENANT_DOMAIN,API_TYPE) VALUES (?,?,?,?,?)";

    public static final String ADD_GW_API_ARTIFACT =
            "INSERT INTO GW_API_ARTIFACTS (ARTIFACT,TIME_STAMP, API_ID,REVISION_ID) VALUES (?,?,?,?)";

    public static final String UPDATE_API_ARTIFACT = "UPDATE GW_API_ARTIFACTS SET ARTIFACT = ?, " +
            "TIME_STAMP = ? WHERE (API_ID = ?) AND (REVISION_ID = ?)";

    public static final String GET_API_ARTIFACT = "SELECT ARTIFACT FROM GW_API_ARTIFACTS WHERE API_ID =? AND " +
            "GATEWAY_LABEL =? AND GATEWAY_INSTRUCTION = ?";

    public static final String GET_API_ARTIFACT_ANY_INSTRUCTION = "SELECT ARTIFACT FROM GW_API_ARTIFACTS WHERE " +
            "API_ID =? AND GATEWAY_LABEL =?";

    public static final String GET_API_ID = "SELECT API_ID  FROM GW_PUBLISHED_API_DETAILS " +
            "WHERE API_NAME =? AND " + "TENANT_DOMAIN =? AND API_VERSION =?";

    public static final String GET_API_LABEL = "SELECT GATEWAY_LABEL  FROM GW_API_ARTIFACTS " +
            "WHERE API_ID =? AND GATEWAY_INSTRUCTION = 'Publish' ";

    public static final String GET_ALL_API_ARTIFACT = "SELECT ARTIFACT FROM GW_API_ARTIFACTS WHERE "
            + "GATEWAY_LABEL =? AND GATEWAY_INSTRUCTION = ?";

    public static final String GET_PUBLISHED_GATEWAYS_FOR_API = "SELECT COUNT(*) AS COUNT FROM GW_API_ARTIFACTS" +
            " WHERE API_ID = ? AND GATEWAY_INSTRUCTION = ?";

    public static final String CHECK_API_EXISTS = "SELECT 1 FROM GW_PUBLISHED_API_DETAILS" +
            " WHERE API_ID = ?";

    public static final String CHECK_ARTIFACT_EXISTS = "SELECT 1 FROM GW_API_ARTIFACTS" +
            " WHERE API_ID = ? AND REVISION_ID = ?";
    public static final String ADD_GW_PUBLISHED_LABELS = "INSERT INTO GW_API_DEPLOYMENTS (API_ID,REVISION_ID," +
            "LABEL,VHOST) VALUES (?,?,?,?)";
    public static final String DELETE_GW_PUBLISHED_LABELS = "DELETE FROM GW_API_DEPLOYMENTS WHERE API_ID = ? AND " +
            "REVISION_ID = ?";
    public static final String DELETE_GW_PUBLISHED_LABELS_BY_API_ID_REVISION_ID_DEPLOYMENT = "DELETE FROM " +
            "GW_API_DEPLOYMENTS WHERE API_ID = ? AND REVISION_ID = ? AND LABEL=?";
    public static final String DELETE_FROM_GW_API_ARTIFACTS_WHERE_API_ID_AND_REVISION_ID =
            "DELETE FROM GW_API_ARTIFACTS WHERE API_ID = ? AND REVISION_ID = ?";
    public static final String DELETE_FROM_GW_API_ARTIFACTS_BY_API_ID =
            "DELETE FROM GW_API_ARTIFACTS WHERE API_ID = ?";
    public static final String DELETE_GW_PUBLISHED_LABELS_BY_API_ID =
            "DELETE FROM GW_API_DEPLOYMENTS WHERE API_ID = ?";
    public static final String DELETE_GW_PUBLISHED_API_DETAILS = "DELETE FROM GW_PUBLISHED_API_DETAILS WHERE " +
            "API_ID = ?";
    public static final String RETRIEVE_API_ARTIFACT_PROPERTY_VALUES =
            "SELECT API.ORGANIZATION AS ORGANIZATION, " +
                    "DEPLOYMENT_REVISION_MAPPING.DEPLOYED_TIME AS DEPLOYED_TIME " +
                    "FROM API, DEPLOYMENT_REVISION_MAPPING " +
                    "WHERE API.API_UUID = ? AND DEPLOYMENT_REVISION_MAPPING.NAME = ? " +
                    "AND DEPLOYMENT_REVISION_MAPPING.REVISION_UUID = ?";
    public static final String RETRIEVE_ARTIFACTS_BY_APIID_AND_LABEL =
            "SELECT GW_API_DEPLOYMENTS.REVISION_ID AS REVISION_ID,GW_PUBLISHED_API_DETAILS" +
                    ".TENANT_DOMAIN AS TENANT_DOMAIN," +
                    "GW_PUBLISHED_API_DETAILS" +
                    ".API_PROVIDER AS API_PROVIDER," +
                    "GW_PUBLISHED_API_DETAILS.API_NAME AS API_NAME,GW_PUBLISHED_API_DETAILS.API_VERSION AS API_VERSION," +
                    "GW_PUBLISHED_API_DETAILS.API_TYPE AS API_TYPE,GW_API_ARTIFACTS.ARTIFACT AS ARTIFACT," +
                    "GW_API_DEPLOYMENTS.LABEL AS LABEL,GW_API_DEPLOYMENTS.VHOST AS VHOST, " +
                    "API.CONTEXT AS CONTEXT FROM " +
                    "GW_PUBLISHED_API_DETAILS,GW_API_ARTIFACTS,GW_API_DEPLOYMENTS,API WHERE " +
                    "GW_API_DEPLOYMENTS.API_ID= ? AND GW_API_DEPLOYMENTS.LABEL IN (_GATEWAY_LABELS_) AND " +
                    "GW_PUBLISHED_API_DETAILS.TENANT_DOMAIN = ? " +
                    "AND GW_PUBLISHED_API_DETAILS.API_ID=GW_API_DEPLOYMENTS.API_ID AND " +
                    "GW_API_ARTIFACTS.API_ID=GW_API_DEPLOYMENTS.API_ID AND " +
                    "API.API_UUID=GW_API_DEPLOYMENTS.API_ID AND" +
                    " GW_API_ARTIFACTS.REVISION_ID=GW_API_DEPLOYMENTS.REVISION_ID";

    public static final String RETRIEVE_ALL_ARTIFACTS_BY_APIID_AND_LABEL=
            "SELECT GW_API_DEPLOYMENTS.REVISION_ID AS REVISION_ID," +
                    "GW_PUBLISHED_API_DETAILS" +
                    ".API_PROVIDER AS API_PROVIDER," +
                    "GW_PUBLISHED_API_DETAILS.API_NAME AS API_NAME,GW_PUBLISHED_API_DETAILS.API_VERSION AS API_VERSION," +
                    "GW_PUBLISHED_API_DETAILS.API_TYPE AS API_TYPE,GW_API_ARTIFACTS.ARTIFACT AS ARTIFACT," +
                    "GW_API_DEPLOYMENTS.LABEL AS LABEL,GW_API_DEPLOYMENTS.VHOST AS VHOST, " +
                    "API.CONTEXT AS CONTEXT FROM " +
                    "GW_PUBLISHED_API_DETAILS,GW_API_ARTIFACTS,GW_API_DEPLOYMENTS,API WHERE " +
                    "GW_API_DEPLOYMENTS.API_ID= ? AND GW_API_DEPLOYMENTS.LABEL IN (_GATEWAY_LABELS_) " +
                    "AND GW_PUBLISHED_API_DETAILS.API_ID=GW_API_DEPLOYMENTS.API_ID AND " +
                    "GW_API_ARTIFACTS.API_ID=GW_API_DEPLOYMENTS.API_ID AND " +
                    "API.API_UUID=GW_API_DEPLOYMENTS.API_ID AND" +
                    " GW_API_ARTIFACTS.REVISION_ID=GW_API_DEPLOYMENTS.REVISION_ID";
    public static final String RETRIEVE_ARTIFACTS_BY_MULTIPLE_APIIDs_AND_LABEL =
            "SELECT GW_API_DEPLOYMENTS.API_ID AS API_ID,GW_API_DEPLOYMENTS.REVISION_ID AS REVISION_ID,GW_PUBLISHED_API_DETAILS" +
                    ".TENANT_DOMAIN AS TENANT_DOMAIN," +
                    "GW_PUBLISHED_API_DETAILS" +
                    ".API_PROVIDER AS API_PROVIDER," +
                    "GW_PUBLISHED_API_DETAILS.API_NAME AS API_NAME,GW_PUBLISHED_API_DETAILS.API_VERSION AS API_VERSION," +
                    "GW_PUBLISHED_API_DETAILS.API_TYPE AS API_TYPE,GW_API_ARTIFACTS.ARTIFACT AS ARTIFACT," +
                    "GW_API_DEPLOYMENTS.LABEL AS LABEL,GW_API_DEPLOYMENTS.VHOST AS VHOST, " +
                    "API.CONTEXT AS CONTEXT FROM " +
                    "GW_PUBLISHED_API_DETAILS,GW_API_ARTIFACTS,GW_API_DEPLOYMENTS,API WHERE " +
                    "GW_API_DEPLOYMENTS.API_ID IN (_API_IDS_) AND GW_API_DEPLOYMENTS.LABEL IN (_GATEWAY_LABELS_) AND " +
                    "GW_PUBLISHED_API_DETAILS.TENANT_DOMAIN = ? " +
                    "AND GW_PUBLISHED_API_DETAILS.API_ID=GW_API_DEPLOYMENTS.API_ID AND " +
                    "GW_API_ARTIFACTS.API_ID=GW_API_DEPLOYMENTS.API_ID AND " +
                    "API.API_UUID=GW_API_DEPLOYMENTS.API_ID AND" +
                    " GW_API_ARTIFACTS.REVISION_ID=GW_API_DEPLOYMENTS.REVISION_ID";
    public static final String RETRIEVE_ARTIFACTS_ONLY_BY_MULTIPLE_APIIDS =
            "SELECT GW_API_DEPLOYMENTS.API_ID AS API_ID," +
                    "GW_API_DEPLOYMENTS.REVISION_ID AS REVISION_ID," +
                    "GW_PUBLISHED_API_DETAILS.TENANT_DOMAIN AS TENANT_DOMAIN," +
                    "GW_PUBLISHED_API_DETAILS.API_PROVIDER AS API_PROVIDER," +
                    "GW_PUBLISHED_API_DETAILS.API_NAME AS API_NAME,GW_PUBLISHED_API_DETAILS.API_VERSION AS API_VERSION," +
                    "GW_PUBLISHED_API_DETAILS.API_TYPE AS API_TYPE,GW_API_ARTIFACTS.ARTIFACT AS ARTIFACT," +
                    "GW_API_DEPLOYMENTS.LABEL AS LABEL,GW_API_DEPLOYMENTS.VHOST AS VHOST, " +
                    "API.CONTEXT AS CONTEXT FROM " +
                    "GW_PUBLISHED_API_DETAILS,GW_API_ARTIFACTS,GW_API_DEPLOYMENTS,API WHERE " +
                    "GW_API_DEPLOYMENTS.API_ID IN (_API_IDS_) AND " +
                    "GW_PUBLISHED_API_DETAILS.TENANT_DOMAIN = ? " +
                    "AND GW_PUBLISHED_API_DETAILS.API_ID=GW_API_DEPLOYMENTS.API_ID AND " +
                    "GW_API_ARTIFACTS.API_ID=GW_API_DEPLOYMENTS.API_ID AND " +
                    "API.API_UUID=GW_API_DEPLOYMENTS.API_ID AND" +
                    " GW_API_ARTIFACTS.REVISION_ID=GW_API_DEPLOYMENTS.REVISION_ID";

    public static final String RETRIEVE_ALL_ARTIFACTS_BY_MULTIPLE_APIIDs_AND_LABEL =
            "SELECT GW_API_DEPLOYMENTS.API_ID AS API_ID,GW_API_DEPLOYMENTS.REVISION_ID AS REVISION_ID," +
                    "GW_PUBLISHED_API_DETAILS" +
                    ".API_PROVIDER AS API_PROVIDER," +
                    "GW_PUBLISHED_API_DETAILS.API_NAME AS API_NAME,GW_PUBLISHED_API_DETAILS.API_VERSION AS API_VERSION," +
                    "GW_PUBLISHED_API_DETAILS.API_TYPE AS API_TYPE,GW_API_ARTIFACTS.ARTIFACT AS ARTIFACT," +
                    "GW_API_DEPLOYMENTS.LABEL AS LABEL,GW_API_DEPLOYMENTS.VHOST AS VHOST, " +
                    "API.CONTEXT AS CONTEXT FROM " +
                    "GW_PUBLISHED_API_DETAILS,GW_API_ARTIFACTS,GW_API_DEPLOYMENTS,API WHERE " +
                    "GW_API_DEPLOYMENTS.API_ID IN (_API_IDS_) AND GW_API_DEPLOYMENTS.LABEL IN (_GATEWAY_LABELS_) " +
                    "AND GW_PUBLISHED_API_DETAILS.API_ID=GW_API_DEPLOYMENTS.API_ID AND " +
                    "GW_API_ARTIFACTS.API_ID=GW_API_DEPLOYMENTS.API_ID AND " +
                    "API.API_UUID=GW_API_DEPLOYMENTS.API_ID AND" +
                    " GW_API_ARTIFACTS.REVISION_ID=GW_API_DEPLOYMENTS.REVISION_ID";
    public static final String RETRIEVE_ARTIFACTS_BY_LABEL =
            "SELECT GW_API_DEPLOYMENTS.API_ID AS API_ID,GW_API_DEPLOYMENTS.REVISION_ID AS REVISION_ID," +
                    "GW_PUBLISHED_API_DETAILS.TENANT_DOMAIN AS TENANT_DOMAIN," +
                    "GW_PUBLISHED_API_DETAILS" +
                    ".API_PROVIDER AS API_PROVIDER," +
                    "GW_PUBLISHED_API_DETAILS.API_NAME AS API_NAME,GW_PUBLISHED_API_DETAILS.API_VERSION AS API_VERSION," +
                    "GW_PUBLISHED_API_DETAILS.API_TYPE AS API_TYPE,GW_API_ARTIFACTS.ARTIFACT AS ARTIFACT," +
                    "GW_API_DEPLOYMENTS.LABEL AS LABEL,GW_API_DEPLOYMENTS.VHOST AS VHOST, " +
                    "API.CONTEXT AS CONTEXT FROM " +
                    "GW_PUBLISHED_API_DETAILS,GW_API_ARTIFACTS,GW_API_DEPLOYMENTS,API WHERE " +
                    "GW_API_DEPLOYMENTS.LABEL IN (_GATEWAY_LABELS_) AND GW_PUBLISHED_API_DETAILS.TENANT_DOMAIN " +
                    "= ? " +
                    "AND GW_PUBLISHED_API_DETAILS.API_ID=GW_API_DEPLOYMENTS.API_ID AND " +
                    "GW_API_ARTIFACTS.API_ID=GW_API_DEPLOYMENTS.API_ID AND " +
                    "API.API_UUID=GW_API_DEPLOYMENTS.API_ID AND " +
                    "GW_API_ARTIFACTS.REVISION_ID=GW_API_DEPLOYMENTS.REVISION_ID";

    public static final String RETRIEVE_ALL_ARTIFACTS_BY_LABEL =
            "SELECT GW_API_DEPLOYMENTS.API_ID AS API_ID,GW_API_DEPLOYMENTS.REVISION_ID AS REVISION_ID," +
                    "GW_PUBLISHED_API_DETAILS" +
                    ".API_PROVIDER AS API_PROVIDER," +
                    "GW_PUBLISHED_API_DETAILS.API_NAME AS API_NAME,GW_PUBLISHED_API_DETAILS.API_VERSION AS API_VERSION," +
                    "GW_PUBLISHED_API_DETAILS.API_TYPE AS API_TYPE,GW_API_ARTIFACTS.ARTIFACT AS ARTIFACT," +
                    "GW_API_DEPLOYMENTS.LABEL AS LABEL,GW_API_DEPLOYMENTS.VHOST AS VHOST, " +
                    "API.CONTEXT AS CONTEXT FROM " +
                    "GW_PUBLISHED_API_DETAILS,GW_API_ARTIFACTS,GW_API_DEPLOYMENTS,API WHERE " +
                    "GW_API_DEPLOYMENTS.LABEL IN (_GATEWAY_LABELS_) " +
                    "AND GW_PUBLISHED_API_DETAILS.API_ID=GW_API_DEPLOYMENTS.API_ID AND " +
                    "GW_API_ARTIFACTS.API_ID=GW_API_DEPLOYMENTS.API_ID AND " +
                    "API.API_UUID=GW_API_DEPLOYMENTS.API_ID AND " +
                    "GW_API_ARTIFACTS.REVISION_ID=GW_API_DEPLOYMENTS.REVISION_ID";
    public static final String RETRIEVE_ARTIFACTS =
            "SELECT GW_API_DEPLOYMENTS.API_ID AS API_ID,GW_API_DEPLOYMENTS.REVISION_ID AS REVISION_ID," +
                    "GW_PUBLISHED_API_DETAILS.TENANT_DOMAIN AS TENANT_DOMAIN,GW_PUBLISHED_API_DETAILS.API_PROVIDER AS " +
                    "API_PROVIDER,GW_PUBLISHED_API_DETAILS.API_NAME AS API_NAME,GW_PUBLISHED_API_DETAILS.API_VERSION AS " +
                    "API_VERSION," +
                    "GW_PUBLISHED_API_DETAILS.API_TYPE AS API_TYPE,GW_API_ARTIFACTS.ARTIFACT AS ARTIFACT," +
                    "GW_API_DEPLOYMENTS.LABEL AS LABEL,GW_API_DEPLOYMENTS.VHOST AS VHOST, " +
                    "API.CONTEXT AS CONTEXT FROM " +
                    "GW_PUBLISHED_API_DETAILS,GW_API_ARTIFACTS,GW_API_DEPLOYMENTS,API WHERE " +
                    "GW_PUBLISHED_API_DETAILS.API_ID=GW_API_DEPLOYMENTS.API_ID AND " +
                    "GW_API_ARTIFACTS.API_ID=GW_API_DEPLOYMENTS.API_ID AND " +
                    "API.API_UUID=GW_API_DEPLOYMENTS.API_ID AND " +
                    "GW_API_ARTIFACTS.REVISION_ID=GW_API_DEPLOYMENTS.REVISION_ID AND " +
                    "GW_PUBLISHED_API_DETAILS.TENANT_DOMAIN = ?";

    public static final String RETRIEVE_ALL_ARTIFACTS =
            "SELECT GW_API_DEPLOYMENTS.API_ID AS API_ID,GW_API_DEPLOYMENTS.REVISION_ID AS REVISION_ID," +
                    "GW_PUBLISHED_API_DETAILS.API_PROVIDER AS " +
                    "API_PROVIDER, GW_PUBLISHED_API_DETAILS.API_NAME AS API_NAME,GW_PUBLISHED_API_DETAILS.API_VERSION AS " +
                    "API_VERSION," +
                    "GW_PUBLISHED_API_DETAILS.API_TYPE AS API_TYPE,GW_API_ARTIFACTS.ARTIFACT AS ARTIFACT," +
                    "GW_API_DEPLOYMENTS.LABEL AS LABEL,GW_API_DEPLOYMENTS.VHOST AS VHOST, " +
                    "API.CONTEXT AS CONTEXT FROM " +
                    "GW_PUBLISHED_API_DETAILS,GW_API_ARTIFACTS,GW_API_DEPLOYMENTS,API WHERE " +
                    "GW_PUBLISHED_API_DETAILS.API_ID=GW_API_DEPLOYMENTS.API_ID AND " +
                    "GW_API_ARTIFACTS.API_ID=GW_API_DEPLOYMENTS.API_ID AND " +
                    "API.API_UUID=GW_API_DEPLOYMENTS.API_ID AND " +
                    "GW_API_ARTIFACTS.REVISION_ID=GW_API_DEPLOYMENTS.REVISION_ID";

    public static final String UPDATE_API_STATUS = "UPDATE API SET STATUS = ? WHERE API_ID = ?";
    public static final String RETRIEVE_API_STATUS_FROM_UUID = "SELECT STATUS FROM API WHERE API_UUID = ?";
    public static final String RETRIEVE_API_INFO_FROM_UUID = "SELECT API_UUID, API_PROVIDER, API_NAME, API_VERSION, " +
            "CONTEXT, CONTEXT_TEMPLATE, API_TIER, API_TYPE, CREATED_BY, CREATED_TIME, UPDATED_BY, UPDATED_TIME, " +
            " ORGANIZATION, REVISIONS_CREATED, STATUS FROM API WHERE API_UUID = ?";
    public static final String RETRIEVE_DEFAULT_VERSION = "SELECT DEFAULT_API_VERSION " +
            "FROM API WHERE API_NAME = ? AND API_PROVIDER =?";
    public static final String UPDATE_REVISION_CREATED_BY_API_SQL = "UPDATE API SET REVISIONS_CREATED = ? WHERE " +
            "API_UUID = ?";
    public static final String ADD_API_REVISION_METADATA = "INSERT INTO API_REVISION_METADATA (API_UUID," +
            "REVISION_UUID,API_TIER) VALUES(?,?,(SELECT API_TIER FROM API WHERE API_UUID = ? ))";
    public static final String DELETE_API_REVISION_METADATA = "DELETE FROM API_REVISION_METADATA WHERE API_UUID = " +
            "? AND REVISION_UUID = ?";
    public static final String GET_REVISIONED_API_TIER_SQL = "SELECT API_TIER FROM API_REVISION_METADATA WHERE " +
            "API_UUID = ? AND REVISION_UUID = ?";
    public static final String RESTORE_API_REVISION_METADATA = "UPDATE API SET API_TIER = (SELECT API_TIER FROM " +
            "API_REVISION_METADATA WHERE API_UUID = ? AND REVISION_UUID = ?) WHERE API_UUID = ?";
    public static final String ADD_PER_API_LOGGING_SQL =
            "UPDATE API SET LOG_LEVEL=? WHERE API_UUID=? AND ORGANIZATION=?";
    public static final String RETRIEVE_PER_API_LOGGING_OFF_SQL =
            "SELECT API.API_UUID, API.LOG_LEVEL, API.API_NAME, API.CONTEXT, API.API_VERSION " +
                    "FROM API WHERE API.LOG_LEVEL = 'OFF' AND API.ORGANIZATION=?";
    public static final String RETRIEVE_PER_API_LOGGING_BASIC_SQL =
            "SELECT API.API_UUID, API.LOG_LEVEL, API.API_NAME, API.CONTEXT, API.API_VERSION " +
                    "FROM API WHERE API.LOG_LEVEL = 'BASIC' AND API.ORGANIZATION=?";
    public static final String RETRIEVE_PER_API_LOGGING_STANDARD_SQL =
            "SELECT API.API_UUID, API.LOG_LEVEL, API.API_NAME, API.CONTEXT, API.API_VERSION " +
                    "FROM API WHERE API.LOG_LEVEL = 'STANDARD' AND API.ORGANIZATION=?";
    public static final String RETRIEVE_PER_API_LOGGING_FULL_SQL =
            "SELECT API.API_UUID, API.LOG_LEVEL, API.API_NAME, API.CONTEXT, API.API_VERSION " +
                    "FROM API WHERE API.LOG_LEVEL = 'FULL' AND API.ORGANIZATION=?";
    public static final String RETRIEVE_ALL_PER_API_LOGGING_SQL =
            "SELECT API.API_UUID, API.LOG_LEVEL, API.API_NAME, API.CONTEXT, API.API_VERSION " +
            "FROM API WHERE API.LOG_LEVEL <> 'OFF'";
    public static final String RETRIEVE_PER_API_LOGGING_ALL_SQL =
            "SELECT API.API_UUID, API.LOG_LEVEL, API.API_NAME, API.CONTEXT, API.API_VERSION " +
            "FROM API WHERE API.ORGANIZATION=?";
    public static final String RETRIEVE_PER_API_LOGGING_BY_UUID_SQL =
            "SELECT API.API_UUID, API.LOG_LEVEL, API.API_NAME, API.CONTEXT, API.API_VERSION " +
            "FROM API WHERE API.API_UUID = ? AND API.ORGANIZATION=?";
    public static final String GATEWAY_LABEL_REGEX = "_GATEWAY_LABELS_";
    public static final String API_ID_REGEX = "_API_IDS_";
    public static final String API_UUID_REGEX = "_API_UUIDS_";
    public static final int API_ID_CHUNK_SIZE = 25;

    /** Throttle related constants**/

    public static class ThrottleSQLConstants{

    	public static final String QUOTA_TYPE_BANDWIDTH = PolicyConstants.BANDWIDTH_TYPE;

    	public static final String QUOTA_TYPE_REQUESTCOUNT = PolicyConstants.REQUEST_COUNT_TYPE;

		public static final String GET_POLICY_NAMES = " SELECT " + "   NAME " + "FROM " + "   API_THROTTLE_POLICY"
				+ " WHERE" + "   TYPE = ?" + "   AND ORGANIZATION =?";

		public static final String GET_EXISTING_POLICY_SQL = "SELECT POLICY_ID FROM API_THROTTLE_POLICY WHERE NAME = ? AND ORGANIZATION = ? ";

		public static final String INSERT_API_POLICY_SQL = "INSERT INTO API_THROTTLE_POLICY (NAME, DISPLAY_NAME, ORGANIZATION, DESCRIPTION, DEFAULT_QUOTA_TYPE, \n"
				+ "  DEFAULT_QUOTA, DEFAULT_QUOTA_UNIT, DEFAULT_UNIT_TIME, DEFAULT_TIME_UNIT , IS_DEPLOYED, UUID, APPLICABLE_LEVEL) \n"
				+ " VALUES (?,?,?,?,? ,?,?,?,?,? ,?,?)";

		public static final String INSERT_API_POLICY_WITH_ID_SQL = "INSERT INTO API_THROTTLE_POLICY (NAME, DISPLAY_NAME, ORGANIZATION, DESCRIPTION, DEFAULT_QUOTA_TYPE, \n"
				+ " DEFAULT_QUOTA, DEFAULT_QUOTA_UNIT, DEFAULT_UNIT_TIME, DEFAULT_TIME_UNIT, \n"
				+ " IS_DEPLOYED, UUID, APPLICABLE_LEVEL, POLICY_ID) \n" + "VALUES (?,?,?,?,?, ?,?,?,?,? ,?,?,?)";


        public static final String UPDATE_API_POLICY_BY_UUID_SQL = "UPDATE API_THROTTLE_POLICY SET DISPLAY_NAME = ?, "
                + "DESCRIPTION = ?, DEFAULT_QUOTA_TYPE = ?, DEFAULT_QUOTA = ?, DEFAULT_QUOTA_UNIT = ?, "
                + "DEFAULT_UNIT_TIME = ?, DEFAULT_TIME_UNIT = ? WHERE UUID = ?";

		public static final String UPDATE_API_POLICY_SQL = "UPDATE API_THROTTLE_POLICY SET DISPLAY_NAME = ?,"
                + "DESCRIPTION = ?, DEFAULT_QUOTA_TYPE = ?, DEFAULT_QUOTA = ?, DEFAULT_QUOTA_UNIT = ?,"
                + "DEFAULT_UNIT_TIME = ?, DEFAULT_TIME_UNIT = ? WHERE NAME = ? AND ORGANIZATION = ?";

		public static final String GET_API_POLICY_NAMES = " SELECT " + "   NAME " + "FROM "
				+ "   API_THROTTLE_POLICY " + " WHERE" + "   ORGANIZATION =?";

		public static final String GET_API_POLICIES = " SELECT " + "   * " + "FROM " + "   API_THROTTLE_POLICY "
				+ " WHERE" + "   ORGANIZATION =? ORDER BY NAME";

		public static final String GET_API_POLICY_ID_SQL = "SELECT " + "POLICY_ID, UUID " + "FROM "
				+ " API_THROTTLE_POLICY " + "WHERE " + "NAME = ? AND " + "ORGANIZATION = ?";

        public static final String GET_API_POLICY_ID_BY_UUID_SQL = "SELECT " + "POLICY_ID, UUID " + "FROM "
                + " API_THROTTLE_POLICY " + "WHERE " + "UUID = ?";

		public static final String GET_API_POLICY_SQL = "SELECT " + "* " + "FROM " + "API_THROTTLE_POLICY "
				+ " WHERE " + "NAME = ? AND " + "ORGANIZATION =?";

        public static final String GET_API_POLICY_BY_UUID_SQL = "SELECT " + "* " + "FROM " + "API_THROTTLE_POLICY "
                + " WHERE " + "UUID = ?";

		public static final String UPDATE_API_POLICY_STATUS_SQL = "UPDATE API_THROTTLE_POLICY SET IS_DEPLOYED = ? WHERE NAME = ? AND ORGANIZATION = ?";

		public static final String DELETE_API_POLICY_SQL = "DELETE FROM API_THROTTLE_POLICY WHERE ORGANIZATION = ? AND NAME = ?";

		public static final String INSERT_CONDITION_GROUP_SQL = "INSERT INTO CONDITION_GROUP(POLICY_ID, QUOTA_TYPE,QUOTA,QUOTA_UNIT,UNIT_TIME,TIME_UNIT,DESCRIPTION) \n"
															+ " VALUES (?,?,?,?,?,?,?)";
        public static final String DELETE_CONDITION_GROUP_SQL = "DELETE FROM CONDITION_GROUP WHERE POLICY_ID = ?";

        public static final String UPDATE_CONDITION_GROUP_SQL =  "UPDATE CONDITION_GROUP SET QUOTA_TYPE = ?, "
                + "QUOTA = ?, QUOTA_UNIT = ?, UNIT_TIME = ?, TIME_UNIT = ?, DESCRIPTION = ? WHERE POLICY_ID = ?";

        public static final String GET_PIPELINES_SQL = "SELECT " + "CONDITION_GROUP_ID, " + "QUOTA_TYPE, " + "QUOTA, "
				+ " QUOTA_UNIT, " + "UNIT_TIME, " + "TIME_UNIT, "+ "DESCRIPTION " + "FROM " + "CONDITION_GROUP " + "WHERE " + "POLICY_ID =?";

		public static final String GET_IP_CONDITIONS_SQL = "SELECT " + "STARTING_IP, " + "ENDING_IP, " + "SPECIFIC_IP, "
				+ " WITHIN_IP_RANGE " + "FROM " + "IP_CONDITION " + "WHERE " + "CONDITION_GROUP_ID = ? ";

		public static final String GET_HEADER_CONDITIONS_SQL = "SELECT " + "HEADER_FIELD_NAME, " + "HEADER_FIELD_VALUE , IS_HEADER_FIELD_MAPPING "
				+ " FROM " + "HEADER_FIELD_CONDITION " + "WHERE " + "CONDITION_GROUP_ID =?";

		public static final String GET_JWT_CLAIM_CONDITIONS_SQL = "SELECT " + "CLAIM_URI, " + "CLAIM_ATTRIB , IS_CLAIM_MAPPING " + "FROM "
				+ " JWT_CLAIM_CONDITION " + "WHERE " + "CONDITION_GROUP_ID =?";

		public static final String GET_QUERY_PARAMETER_CONDITIONS_SQL = "SELECT " + "PARAMETER_NAME, "
				+ " PARAMETER_VALUE , IS_PARMAPPING " + "FROM " + "QUERY_PARAMETER_CONDITION " + "WHERE " + "CONDITION_GROUP_ID =?";

		public static final String INSERT_QUERY_PARAMETER_CONDITION_SQL = "INSERT INTO QUERY_PARAMETER_CONDITION(CONDITION_GROUP_ID,PARAMETER_NAME,PARAMETER_VALUE, IS_PARMAPPING) \n"
				+ " VALUES (?,?,?,?)";

		public static final String INSERT_HEADER_FIELD_CONDITION_SQL = "INSERT INTO HEADER_FIELD_CONDITION(CONDITION_GROUP_ID,HEADER_FIELD_NAME,HEADER_FIELD_VALUE, IS_HEADER_FIELD_MAPPING) \n"
				+ " VALUES (?,?,?,?)";

		public static final String INSERT_JWT_CLAIM_CONDITION_SQL = "INSERT INTO JWT_CLAIM_CONDITION(CONDITION_GROUP_ID,CLAIM_URI,CLAIM_ATTRIB,IS_CLAIM_MAPPING) \n"
				+ " VALUES (?,?,?,?)";

		public static final String INSERT_IP_CONDITION_SQL =
	            " INSERT INTO IP_CONDITION(STARTING_IP,ENDING_IP,SPECIFIC_IP,WITHIN_IP_RANGE,CONDITION_GROUP_ID ) \n" +
	            " VALUES (?,?,?,?,?)";

		public static final String IS_ANY_POLICY_CONTENT_AWARE_WITHOUT_API_POLICY_SQL = "SELECT APPPOLICY.ORGANIZATION, APPPOLICY.QUOTA_TYPE "
				+ " FROM APPLICATION_USAGE_PLAN APPPOLICY," + "BUSINESS_PLAN SUBPOLICY "
				+ " WHERE APPPOLICY.ORGANIZATION =? AND " + "APPPOLICY.NAME =? AND " + "SUBPOLICY.NAME=? ";

		public static final String IS_ANY_POLICY_CONTENT_AWARE_SQL = "select sum(c) as c from("
				+ " (SELECT count(*) as c"
				+ " FROM API_THROTTLE_POLICY APIPOLICY where APIPOLICY.NAME =?  AND APIPOLICY.ORGANIZATION =? AND APIPOLICY.DEFAULT_QUOTA_TYPE = 'bandwidthVolume')"
				+ " union "
				+ " (SELECT count(*) as c"
				+ " FROM API_THROTTLE_POLICY APIPOLICY , CONDITION_GROUP cg where APIPOLICY.NAME =?  AND APIPOLICY.ORGANIZATION =? AND cg.POLICY_ID = APIPOLICY.POLICY_ID AND cg.QUOTA_TYPE = 'bandwidthVolume')"
				+ " union "
				+ " (SELECT count(*) as c"
				+ " FROM API_THROTTLE_POLICY APIPOLICY, API_URL_MAPPING RS, CONDITION_GROUP cg where"
				+ " RS.API_ID = ? AND RS.REVISION_UUID IS NULL AND APIPOLICY.NAME = RS.THROTTLING_TIER AND APIPOLICY.ORGANIZATION =? AND cg.POLICY_ID = APIPOLICY.POLICY_ID AND cg.QUOTA_TYPE = 'bandwidthVolume' "
				+ " ) "
				+ " union "
				+ "  (SELECT count(*) as c"
				+ " FROM API_THROTTLE_POLICY APIPOLICY, API_URL_MAPPING RS where "
				+ " RS.API_ID = ? AND RS.REVISION_UUID IS NULL AND APIPOLICY.NAME = RS.THROTTLING_TIER AND APIPOLICY.ORGANIZATION =? AND APIPOLICY.DEFAULT_QUOTA_TYPE = 'bandwidthVolume') "
				+ " union "
				+ " (SELECT count(*) as c FROM BUSINESS_PLAN SUBPOLICY WHERE SUBPOLICY.NAME= ? AND SUBPOLICY.ORGANIZATION = ? AND SUBPOLICY.QUOTA_TYPE = 'bandwidthVolume')"
				+ " union "
				+ " (SELECT count(*) as c FROM APPLICATION_USAGE_PLAN APPPOLICY where APPPOLICY.NAME = ? AND APPPOLICY.ORGANIZATION = ? AND APPPOLICY.QUOTA_TYPE = 'bandwidthVolume')"
				+ " ) x";

        public static final String GET_CONDITION_GROUPS_FOR_POLICIES_SQL = "SELECT grp.CONDITION_GROUP_ID ,AUM.HTTP_METHOD,AUM.AUTH_SCHEME, pol.APPLICABLE_LEVEL, "
                + " AUM.URL_PATTERN,AUM.THROTTLING_TIER,AUM.MEDIATION_SCRIPT,AUM.URL_MAPPING_ID, pol.DEFAULT_QUOTA_TYPE  "
                + " FROM API_URL_MAPPING AUM"
                + " INNER JOIN  API API ON AUM.API_ID = API.API_ID"
                + " LEFT OUTER JOIN API_THROTTLE_POLICY pol ON AUM.THROTTLING_TIER = pol.NAME "
                + " LEFT OUTER JOIN CONDITION_GROUP grp ON pol.POLICY_ID  = grp.POLICY_ID"
                + " where API.CONTEXT= ? AND API.API_VERSION = ? AND pol.ORGANIZATION = ? AND AUM.REVISION_UUID IS NULL"
                /*+ " GROUP BY AUM.HTTP_METHOD,AUM.URL_PATTERN, AUM.URL_MAPPING_ID"*/
                + " ORDER BY AUM.URL_MAPPING_ID";

        public static final String GET_CONDITION_GROUPS_FOR_POLICIES_IN_PRODUCTS_SQL = "SELECT AUM.HTTP_METHOD, AUM.AUTH_SCHEME, AUM.URL_PATTERN, AUM.THROTTLING_TIER, " +
                "AUM.MEDIATION_SCRIPT, AUM.URL_MAPPING_ID, POL.APPLICABLE_LEVEL, GRP.CONDITION_GROUP_ID, POL.DEFAULT_QUOTA_TYPE " +
                "FROM API_URL_MAPPING AUM, API_PRODUCT_MAPPING APM, API API, API_THROTTLE_POLICY POL " +
                "LEFT OUTER JOIN CONDITION_GROUP GRP ON POL.POLICY_ID  = GRP.POLICY_ID " +
                "WHERE APM.API_ID = API.API_ID " +
                "AND API.CONTEXT = ? AND API.API_VERSION = ? AND POL.ORGANIZATION = ? " +
                "AND APM.URL_MAPPING_ID = AUM.URL_MAPPING_ID AND AUM.THROTTLING_TIER = POL.NAME AND AUM.REVISION_UUID IS NULL " +
                "ORDER BY AUM.URL_MAPPING_ID";

        public static final String ADD_BLOCK_CONDITION_SQL =
                "INSERT INTO BLOCK_CONDITION (TYPE,BLOCK_CONDITION,ENABLED,ORGANIZATION,UUID) VALUES (?,?,?,?,?)";
        public static final String GET_BLOCK_CONDITION_BY_ORG_SQL =
                "SELECT CONDITION_ID,TYPE,BLOCK_CONDITION,ENABLED,ORGANIZATION,UUID FROM BLOCK_CONDITION WHERE ORGANIZATION =?";
        public static final String GET_BLOCK_CONDITION_SQL =
                "SELECT TYPE,BLOCK_CONDITION,ENABLED,ORGANIZATION,UUID FROM BLOCK_CONDITION WHERE CONDITION_ID =?";
        public static final String GET_BLOCK_CONDITION_BY_UUID_SQL =
                "SELECT CONDITION_ID,TYPE,BLOCK_CONDITION,ENABLED,ORGANIZATION,UUID FROM BLOCK_CONDITION WHERE UUID =?";
        public static final String UPDATE_BLOCK_CONDITION_STATE_SQL =
                "UPDATE BLOCK_CONDITION SET ENABLED = ? WHERE CONDITION_ID = ?";
        public static final String UPDATE_BLOCK_CONDITION_STATE_BY_UUID_SQL =
                "UPDATE BLOCK_CONDITION SET ENABLED = ? WHERE UUID = ?";
        public static final String DELETE_BLOCK_CONDITION_SQL =
                "DELETE FROM BLOCK_CONDITION WHERE CONDITION_ID=?";
        public static final String DELETE_BLOCK_CONDITION_BY_UUID_SQL =
                "DELETE FROM BLOCK_CONDITION WHERE UUID=?";
        public static final String BLOCK_CONDITION_EXIST_SQL =
                "SELECT CONDITION_ID,TYPE,BLOCK_CONDITION,ENABLED,ORGANIZATION,UUID FROM BLOCK_CONDITION WHERE ORGANIZATION =? "
                        + "AND TYPE =? AND BLOCK_CONDITION =?";
        public static final String GET_SUBSCRIPTION_BLOCK_CONDITION_BY_VALUE_AND_DOMAIN_SQL =
                "SELECT CONDITION_ID,TYPE,BLOCK_CONDITION,ENABLED,ORGANIZATION,UUID FROM BLOCK_CONDITION WHERE "
                        + "BLOCK_CONDITION = ? AND ORGANIZATION = ? ";

        public static final String TIER_HAS_SUBSCRIPTION = " select count(sub.TIER_ID) as c from SUBSCRIPTION sub, API api "
        		+ " where sub.TIER_ID = ? and api.API_PROVIDER like ? and sub.API_ID = api.API_ID ";

        public static final String TIER_ATTACHED_TO_RESOURCES_API = " select sum(c) as c from("
        		+ " (select count(api.API_TIER) as c from  API api where api.API_TIER = ? and api.API_PROVIDER like ? )"
        		+ "		union "
        		+ " (select count(map.THROTTLING_TIER) as c from API_URL_MAPPING map, API api"
        		+ "  where map.THROTTLING_TIER = ? and api.API_PROVIDER like ?  and map.API_ID = api.API_ID and map.REVISION_UUID IS NULL)) x ";

        public static final String TIER_ATTACHED_TO_APPLICATION = " SELECT count(APPLICATION_TIER) as c FROM APPLICATION where APPLICATION_TIER = ? ";

        public static final String GET_TIERS_WITH_BANDWIDTH_QUOTA_TYPE_SQL = "SELECT NAME "
                + "FROM API_THROTTLE_POLICY LEFT JOIN CONDITION_GROUP "
                + "ON API_THROTTLE_POLICY.POLICY_ID = CONDITION_GROUP.POLICY_ID "
                + "WHERE "
                + "(DEFAULT_QUOTA_TYPE  = '" + QUOTA_TYPE_BANDWIDTH + "' OR QUOTA_TYPE  = '"+ QUOTA_TYPE_BANDWIDTH + "') "
                + "AND ORGANIZATION = ?";
        public static final String TIER_HAS_ATTACHED_TO_APPLICATION = "SELECT 1 FROM APPLICATION WHERE " +
                "ORGANIZATION = ? " +
                "AND APPLICATION.APPLICATION_TIER = ?";

        public static final String TIER_HAS_ATTACHED_TO_SUBSCRIPTION_SUPER_TENANT = "(SELECT 1 from SUBSCRIPTION " +
                "WHERE API_ID IN (SELECT API_ID FROM API WHERE ORGANIZATION = ? AND TIER_ID_PENDING = ?)) " +
                "UNION " +
                "(SELECT 1 FROM SUBSCRIPTION WHERE  API_ID IN (SELECT API_ID FROM API WHERE ORGANIZATION = ? AND" +
                " TIER_ID = ?))";

        public static final String TIER_HAS_ATTACHED_TO_API_RESOURCE_TENANT = "(SELECT 1 FROM API WHERE " +
                "ORGANIZATION = ? AND API_TIER = ?) UNION (SELECT 1 FROM API_URL_MAPPING WHERE API_ID IN " +
                "(SELECT API_ID FROM API WHERE ORGANIZATION = ?) AND THROTTLING_TIER = ?)";
    }

    public static class CertificateConstants {
        public static final String INSERT_CERTIFICATE = "INSERT INTO CERTIFICATE_METADATA " +
                "(TENANT_ID, END_POINT, ALIAS,CERTIFICATE) VALUES(?, ?, ?,?)";

        public static final String GET_CERTIFICATES = "SELECT * FROM CERTIFICATE_METADATA WHERE TENANT_ID=?";

        public static final String GET_CERTIFICATE_ALL_TENANTS = "SELECT * FROM CERTIFICATE_METADATA WHERE " +
                "(ALIAS=?)";
        public static final String GET_CERTIFICATE_TENANT = "SELECT * FROM CERTIFICATE_METADATA WHERE TENANT_ID=? " +
                "AND (ALIAS=? OR END_POINT=?)";
        public static final String GET_CERTIFICATE_TENANT_ALIAS_ENDPOINT = "SELECT * FROM CERTIFICATE_METADATA " +
                       "WHERE TENANT_ID=? AND ALIAS=? AND END_POINT=?";

        public static final String DELETE_CERTIFICATES = "DELETE FROM CERTIFICATE_METADATA WHERE TENANT_ID=? " +
                "AND ALIAS=?";

        public static final String CERTIFICATE_COUNT_QUERY = "SELECT COUNT(*) AS count FROM CERTIFICATE_METADATA " +
                "WHERE TENANT_ID=?";

        public static final String SELECT_CERTIFICATE_FOR_ALIAS = "SELECT * FROM CERTIFICATE_METADATA "
                + "WHERE ALIAS=?";
        public static final String CERTIFICATE_EXIST =
                "SELECT 1 FROM CERTIFICATE_METADATA WHERE ALIAS=? AND TENANT_ID=?";
    }

    public static class ClientCertificateConstants{
        public static final String INSERT_CERTIFICATE = "INSERT INTO API_CLIENT_CERTIFICATE " +
                "(CERTIFICATE, TENANT_ID, ALIAS, API_ID, TIER_NAME) VALUES(?, ?, ?, (SELECT API_ID FROM API WHERE " +
                "API_PROVIDER = ? AND API_NAME = ? AND API_VERSION = ? AND ORGANIZATION = ? ), ?)";

        public static final String GET_CERTIFICATES_FOR_API = "SELECT ALIAS FROM API_CLIENT_CERTIFICATE WHERE "
                + "TENANT_ID=? and API_ID=(SELECT API_ID FROM API WHERE API_PROVIDER = ? AND API_NAME = ? AND " +
                "API_VERSION = ? ) and REMOVED=? and REVISION_UUID ='Current API'";

        public static final String DELETE_CERTIFICATES_FOR_API = "DELETE FROM API_CLIENT_CERTIFICATE "
                + "WHERE TENANT_ID=? and API_ID=(SELECT API_ID FROM API WHERE API_PROVIDER = ? AND API_NAME = ? " +
                "AND API_VERSION = ? ) and REMOVED=? and REVISION_UUID ='Current API'";

        public static final String SELECT_CERTIFICATE_FOR_ALIAS = "SELECT ALIAS FROM API_CLIENT_CERTIFICATE "
                + "WHERE ALIAS=? AND REMOVED=? AND TENANT_ID =? and REVISION_UUID ='Current API'";

        public static final String SELECT_CERTIFICATE_FOR_TENANT =
                "SELECT AC.CERTIFICATE, AC.ALIAS, AC.TIER_NAME, AA.API_PROVIDER, AA.API_NAME, "
                        + "AA.API_VERSION FROM API_CLIENT_CERTIFICATE AC, API AA "
                        + "WHERE AC.REMOVED=? AND AC.TENANT_ID=? AND AA.API_ID=AC.API_ID AND AC.REVISION_UUID ='Current API'";

        public static final String SELECT_CERTIFICATE_FOR_TENANT_ALIAS =
                "SELECT AC.CERTIFICATE, AC.ALIAS, AC.TIER_NAME, AA.API_PROVIDER, AA.API_NAME, AA.API_VERSION "
                        + "FROM API_CLIENT_CERTIFICATE AC, API AA "
                        + "WHERE AC.REMOVED=? AND AC.TENANT_ID=? AND AC.ALIAS=? AND AA.API_ID=AC.API_ID AND AC.REVISION_UUID ='Current API'";

        public static final String SELECT_CERTIFICATE_FOR_TENANT_ALIAS_APIID =
                "SELECT AC.CERTIFICATE, AC.ALIAS, AC.TIER_NAME FROM API_CLIENT_CERTIFICATE AC "
                        + "WHERE AC.REMOVED=? AND AC.TENANT_ID=? AND AC.ALIAS=? AND AC.API_ID = ? AND AC.REVISION_UUID ='Current API'";

        public static final String SELECT_CERTIFICATE_FOR_TENANT_APIID =
                "SELECT AC.CERTIFICATE, AC.ALIAS, AC.TIER_NAME FROM API_CLIENT_CERTIFICATE AC "
                        + "WHERE AC.REMOVED=? AND AC.TENANT_ID=? AND AC.API_ID=? AND AC.REVISION_UUID ='Current API'";

        public static final String PRE_DELETE_CERTIFICATES = "DELETE FROM API_CLIENT_CERTIFICATE "
                + "WHERE TENANT_ID=? AND REMOVED=? AND REVISION_UUID ='Current API' AND ALIAS=? AND API_ID=(SELECT API_ID FROM API WHERE " +
                "API_PROVIDER = ? AND API_NAME = ? AND API_VERSION = ? )";

        public static final String PRE_DELETE_CERTIFICATES_WITHOUT_APIID = "DELETE FROM API_CLIENT_CERTIFICATE "
                + "WHERE TENANT_ID=? AND REMOVED=? and ALIAS=? AND REVISION_UUID ='Current API'";

        public static final String DELETE_CERTIFICATES = "UPDATE API_CLIENT_CERTIFICATE SET REMOVED = ? "
                + "WHERE TENANT_ID=? AND REVISION_UUID ='Current API' AND ALIAS=? AND API_ID=(SELECT API_ID FROM API WHERE API_PROVIDER = ? AND " +
                "API_NAME = ? AND API_VERSION = ? )";

        public static final String DELETE_CERTIFICATES_WITHOUT_APIID = "UPDATE API_CLIENT_CERTIFICATE SET REMOVED=? "
                + "WHERE TENANT_ID=? AND ALIAS=? AND REVISION_UUID ='Current API'";

        public static final String CERTIFICATE_COUNT_QUERY = "SELECT COUNT(*) AS count FROM API_CLIENT_CERTIFICATE " +
                "WHERE TENANT_ID=? AND REMOVED=? AND REVISION_UUID ='Current API'";
    }

    /**
     * Static class to hold database queries related to SYSTEM_APPS table
     */
    public static class SystemApplicationConstants {

        public static final String INSERT_SYSTEM_APPLICATION =
                "INSERT INTO SYSTEM_APPS " + "(NAME,CONSUMER_KEY,CONSUMER_SECRET,TENANT_DOMAIN,CREATED_TIME) " +
                        "VALUES (?,?,?,?,?)";

        public static final String GET_APPLICATIONS =
                "SELECT * FROM " + "SYSTEM_APPS WHERE TENANT_DOMAIN = ?";

        public static final String GET_CLIENT_CREDENTIALS_FOR_APPLICATION =
                "SELECT CONSUMER_KEY,CONSUMER_SECRET FROM " + "SYSTEM_APPS WHERE NAME = ? AND TENANT_DOMAIN = ?";

        public static final String DELETE_SYSTEM_APPLICATION = "DELETE FROM SYSTEM_APPS WHERE NAME = ? AND " +
                "TENANT_DOMAIN = ?";

        public static final String CHECK_CLIENT_CREDENTIALS_EXISTS = "SELECT CONSUMER_KEY,CONSUMER_SECRET " +
                "FROM SYSTEM_APPS WHERE NAME = ? AND TENANT_DOMAIN = ?";
    }

    public static class BotDataConstants {

        public static final String ADD_NOTIFICATION = "INSERT INTO NOTIFICATION_SUBSCRIBER (UUID, CATEGORY," +
        "NOTIFICATION_METHOD, SUBSCRIBER_ADDRESS) VALUES(?,?,?,?)";

        public static final String GET_SAVED_ALERT_EMAILS =
                " SELECT UUID, SUBSCRIBER_ADDRESS FROM NOTIFICATION_SUBSCRIBER";

        public static final String DELETE_EMAIL_BY_UUID =
                "DELETE FROM NOTIFICATION_SUBSCRIBER WHERE UUID= ?";

        public static final String GET_ALERT_SUBSCRIPTION_BY_UUID =
                "SELECT UUID, SUBSCRIBER_ADDRESS FROM NOTIFICATION_SUBSCRIBER WHERE UUID = ?";

        public static final String GET_ALERT_SUBSCRIPTION_BY_EMAIL =
                "SELECT UUID, SUBSCRIBER_ADDRESS FROM NOTIFICATION_SUBSCRIBER WHERE SUBSCRIBER_ADDRESS = ?";

        public static final String GET_BOT_DETECTED_DATA =
                "from BOT_DATA SELECT request_time, message_id, http_method, headers, message_body, client_ip";

    }

    public static class RevokedJWTConstants {

        public static final String ADD_JWT_SIGNATURE = "INSERT INTO REVOKED_JWT (UUID, SIGNATURE," +
                "EXPIRY_TIMESTAMP, TENANT_ID, TOKEN_TYPE) VALUES(?,?,?,?,?)";
        public static final String CHECK_REVOKED_TOKEN_EXIST = "SELECT 1 FROM REVOKED_JWT WHERE UUID = ?";
        public static final String DELETE_REVOKED_JWT = "DELETE FROM REVOKED_JWT WHERE EXPIRY_TIMESTAMP < ?";
    }

    //Shared Scopes related constants
    public static final String ADD_SHARED_SCOPE = "INSERT INTO SHARED_SCOPE (NAME, UUID, TENANT_ID) VALUES (?,?,?)";
    public static final String DELETE_SHARED_SCOPE = "DELETE FROM SHARED_SCOPE WHERE NAME = ? AND TENANT_ID = ?";
    public static final String GET_SHARED_SCOPE_BY_UUID = "SELECT NAME FROM SHARED_SCOPE WHERE UUID = ?";
    public static final String GET_ALL_SHARED_SCOPE_KEYS_BY_TENANT = "SELECT NAME FROM SHARED_SCOPE " +
            "WHERE TENANT_ID = ?";
    public static final String IS_SHARED_SCOPE_NAME_EXISTS = "SELECT 1 FROM SHARED_SCOPE " +
            "WHERE TENANT_ID = ? AND NAME = ?";
    public static final String GET_ALL_SHARED_SCOPES_BY_TENANT = "SELECT UUID, NAME FROM SHARED_SCOPE " +
            "WHERE TENANT_ID = ?";
    public static final String GET_SHARED_SCOPE_USAGE_COUNT_BY_TENANT =
            "SELECT SS.NAME, SS.UUID, "
                    + "(SELECT COUNT(*) FROM API_RESOURCE_SCOPE_MAPPING RSM WHERE RSM.SCOPE_NAME=SS.NAME AND "
                    + "RSM.TENANT_ID = ?) usages "
                    + "FROM SHARED_SCOPE SS "
                    + "WHERE SS.TENANT_ID = ?";
    public static final String GET_SHARED_SCOPE_API_USAGE_BY_TENANT =
            "SELECT AA.API_ID, AA.API_NAME, AA.CONTEXT, AA.API_VERSION, AA.API_PROVIDER "
                    + "FROM SHARED_SCOPE ASSC, API_RESOURCE_SCOPE_MAPPING AARSM, "
                    + "API_URL_MAPPING AAUM, API AA "
                    + "WHERE ASSC.NAME=AARSM.SCOPE_NAME AND "
                    + "AARSM.URL_MAPPING_ID=AAUM.URL_MAPPING_ID AND "
                    + "AAUM.API_ID=AA.API_ID AND AAUM.REVISION_UUID IS NULL AND "
                    + "ASSC.UUID=? AND "
                    + "AARSM.TENANT_ID=? "
                    + "GROUP BY AA.API_ID, AA.API_NAME, AA.CONTEXT, AA.API_VERSION, AA.API_PROVIDER";

    public static final String GET_SHARED_SCOPE_URI_USAGE_BY_TENANT =
            "SELECT AAUM.URL_PATTERN, AAUM.HTTP_METHOD "
                    + "FROM SHARED_SCOPE ASSC, API_RESOURCE_SCOPE_MAPPING AARSM, "
                    + "API_URL_MAPPING AAUM, API AA "
                    + "WHERE ASSC.NAME=AARSM.SCOPE_NAME AND "
                    + "AARSM.URL_MAPPING_ID=AAUM.URL_MAPPING_ID AND "
                    + "AAUM.API_ID=AA.API_ID AND "
                    + "ASSC.UUID=? AND "
                    + "AARSM.TENANT_ID=? AND "
                    + "AA.API_ID=? AND AAUM.REVISION_UUID IS NULL";

    //Resource Scope related constants
    public static final String ADD_API_RESOURCE_SCOPE_MAPPING =
            "INSERT INTO API_RESOURCE_SCOPE_MAPPING (SCOPE_NAME, URL_MAPPING_ID, TENANT_ID) VALUES (?, ?, ?)";
    public static final String IS_SCOPE_ATTACHED_LOCALLY =
            "SELECT API.API_NAME, API.API_PROVIDER "
                    + "FROM API_RESOURCE_SCOPE_MAPPING ARSM, API_URL_MAPPING AUM, API "
                    + "WHERE ARSM.SCOPE_NAME = ? AND "
                    + "API.ORGANIZATION = ? AND "
                    + "ARSM.TENANT_ID = ? AND "
                    + "ARSM.SCOPE_NAME NOT IN (SELECT GS.NAME FROM SHARED_SCOPE GS WHERE GS.TENANT_ID = ?) AND "
                    + "ARSM.URL_MAPPING_ID = AUM.URL_MAPPING_ID AND "
                    + "AUM.API_ID = API.API_ID AND AUM.REVISION_UUID IS NULL";
    public static final String IS_SCOPE_ATTACHED =
            "SELECT 1 FROM API_RESOURCE_SCOPE_MAPPING WHERE SCOPE_NAME = ? AND TENANT_ID = ?";

    public static final String REMOVE_RESOURCE_SCOPE_URL_MAPPING_SQL =
            " DELETE FROM API_RESOURCE_SCOPE_MAPPING "
                    + "WHERE URL_MAPPING_ID IN ( SELECT URL_MAPPING_ID FROM API_URL_MAPPING WHERE API_ID = ? AND REVISION_UUID IS NULL )";

    public static final String GET_UNVERSIONED_LOCAL_SCOPES_FOR_API_SQL =
            "SELECT DISTINCT ARSM.SCOPE_NAME "
                    + "FROM API_RESOURCE_SCOPE_MAPPING ARSM INNER JOIN API_URL_MAPPING AUM "
                    + "ON ARSM.URL_MAPPING_ID = AUM.URL_MAPPING_ID "
                    + "WHERE AUM.API_ID = ? AND AUM.REVISION_UUID IS NULL AND ARSM.TENANT_ID = ? AND "
                    + "ARSM.SCOPE_NAME NOT IN (SELECT GS.NAME FROM SHARED_SCOPE GS WHERE GS.TENANT_ID = ?) AND "
                    + "ARSM.SCOPE_NAME NOT IN ( "
                    + "SELECT ARSM2.SCOPE_NAME FROM API_RESOURCE_SCOPE_MAPPING ARSM2 "
                    + "INNER JOIN API_URL_MAPPING AUM2 ON ARSM2.URL_MAPPING_ID = AUM2.URL_MAPPING_ID "
                    + "WHERE AUM2.API_ID != ? AND ARSM2.TENANT_ID = ?)";

    public static final String GET_VERSIONED_LOCAL_SCOPES_FOR_API_SQL =
            "SELECT DISTINCT ARSM.SCOPE_NAME "
                    + "FROM API_RESOURCE_SCOPE_MAPPING ARSM INNER JOIN API_URL_MAPPING AUM "
                    + "ON ARSM.URL_MAPPING_ID = AUM.URL_MAPPING_ID "
                    + "WHERE AUM.API_ID = ? AND AUM.REVISION_UUID IS NULL AND ARSM.TENANT_ID = ? AND "
                    + "ARSM.SCOPE_NAME NOT IN (SELECT GS.NAME FROM SHARED_SCOPE GS WHERE GS.TENANT_ID = ?) AND "
                    + "ARSM.SCOPE_NAME IN ( "
                    + "SELECT ARSM2.SCOPE_NAME FROM API_RESOURCE_SCOPE_MAPPING ARSM2 "
                    + "INNER JOIN API_URL_MAPPING AUM2 ON ARSM2.URL_MAPPING_ID = AUM2.URL_MAPPING_ID "
                    + "WHERE AUM2.API_ID != ? AND ARSM2.TENANT_ID = ?)";

    public static final String GET_ALL_LOCAL_SCOPES_FOR_API_SQL =
            "SELECT DISTINCT ARSM.SCOPE_NAME "
                    + "FROM API_RESOURCE_SCOPE_MAPPING ARSM INNER JOIN API_URL_MAPPING AUM "
                    + "ON ARSM.URL_MAPPING_ID = AUM.URL_MAPPING_ID "
                    + "WHERE AUM.API_ID = ? AND AUM.REVISION_UUID IS NULL AND ARSM.TENANT_ID = ? AND "
                    + "ARSM.SCOPE_NAME NOT IN (SELECT GS.NAME FROM SHARED_SCOPE GS WHERE GS.TENANT_ID = ?)";

    public static final String GET_URL_TEMPLATES_WITH_SCOPES_FOR_API_SQL =
            " SELECT AUM.URL_MAPPING_ID, "
                    + "AUM.URL_PATTERN, "
                    + "AUM.HTTP_METHOD, "
                    + "AUM.AUTH_SCHEME, "
                    + "AUM.THROTTLING_TIER, "
                    + "AUM.MEDIATION_SCRIPT, "
                    + "ARSM.SCOPE_NAME "
                    + "FROM "
                    + "API_URL_MAPPING AUM "
                    + "INNER JOIN API_RESOURCE_SCOPE_MAPPING ARSM ON AUM.URL_MAPPING_ID = ARSM.URL_MAPPING_ID "
                    + "AND AUM.API_ID = ?";

    public static final String GET_API_SCOPES_SQL =
            " SELECT ARSM.SCOPE_NAME FROM API_RESOURCE_SCOPE_MAPPING ARSM, API_URL_MAPPING AUM "
                    + "WHERE ARSM.URL_MAPPING_ID = AUM.URL_MAPPING_ID AND AUM.API_ID = ? AND AUM.REVISION_UUID IS NULL";

    public static final String INSERT_SCOPE_SQL = "INSERT INTO SCOPE (NAME,DISPLAY_NAME,DESCRIPTION,TENANT_ID," +
            "SCOPE_TYPE) VALUES (?,?,?,?,?)";
    public static final String UPDATE_SCOPE_SQL = "UPDATE SCOPE SET DISPLAY_NAME = ?, DESCRIPTION = ? " +
            "WHERE NAME = ? AND TENANT_ID = ?";
    public static final String DELETE_SCOPE_SQL = "DELETE FROM SCOPE WHERE NAME = ? AND TENANT_ID = ?";
    public static final String GET_SCOPE_SQL = "SELECT NAME AS SCOPE_KEY,DISPLAY_NAME AS DISPLAY_NAME,DESCRIPTION AS " +
            "DESCRIPTION,TENANT_ID AS TENANT_ID FROM SCOPE WHERE NAME = ? AND TENANT_ID = ?";
    public static final String ADD_SCOPE_MAPPING = "INSERT INTO SCOPE_BINDING (SCOPE_ID, SCOPE_BINDING, " +
            "BINDING_TYPE) VALUES((SELECT SCOPE_ID FROM SCOPE WHERE NAME = ? AND TENANT_ID = ?),?,?)";
    public static final String DELETE_SCOPE_MAPPING =
            "DELETE FROM SCOPE_BINDING WHERE SCOPE_ID = (SELECT SCOPE_ID FROM SCOPE WHERE NAME = ? " +
                    "AND TENANT_ID = ?)";
    public static final String RETRIEVE_SCOPE_MAPPING =
            "SELECT SCOPE_BINDING FROM SCOPE_BINDING WHERE SCOPE_ID = (SELECT SCOPE_ID FROM SCOPE " +
                    "WHERE NAME = ? AND TENANT_ID = ?) AND BINDING_TYPE = ?";
    public static final String GET_SCOPES_SQL = "SELECT NAME AS SCOPE_KEY,DISPLAY_NAME AS DISPLAY_NAME,DESCRIPTION AS " +
            "DESCRIPTION,TENANT_ID AS TENANT_ID FROM SCOPE WHERE TENANT_ID = ? AND SCOPE_TYPE = ?";

    public static final String SCOPE_EXIST_SQL = "SELECT 1 FROM SCOPE WHERE NAME = ? AND TENANT_ID = ?";

    /**
     * Static class to hold database queries related to webhooks subscriptions
     */
    public static class WebhooksSqlConstants {
        public static final String FIND_SUBSCRIPTION =
                "SELECT WH_SUBSCRIPTION_ID FROM WEBHOOKS_SUBSCRIPTION WHERE API_UUID = ? AND " +
                        "APPLICATION_ID = ? AND TENANT_DOMAIN = ? AND HUB_CALLBACK_URL = ? AND HUB_TOPIC = ?";
        public static final String ADD_SUBSCRIPTION =
                "INSERT INTO WEBHOOKS_SUBSCRIPTION (API_UUID, APPLICATION_ID, TENANT_DOMAIN, " +
                        "HUB_CALLBACK_URL, HUB_TOPIC, HUB_SECRET, HUB_LEASE_SECONDS, UPDATED_AT, EXPIRY_AT, " +
                        "DELIVERY_STATE) VALUES (?,?,?,?,?,?,?,?,?,?)";
        public static final String UPDATE_EXISTING_SUBSCRIPTION = "UPDATE WEBHOOKS_SUBSCRIPTION SET " +
                "HUB_SECRET = ?, HUB_LEASE_SECONDS = ?, UPDATED_AT = ?, EXPIRY_AT = ?  WHERE WH_SUBSCRIPTION_ID = ?";
        public static final String DELETE_IF_EXISTS_SUBSCRIBER =
                "DELETE FROM WEBHOOKS_SUBSCRIPTION WHERE API_UUID = ? AND APPLICATION_ID = ? AND TENANT_DOMAIN = ? AND " +
                        "HUB_CALLBACK_URL = ? AND HUB_TOPIC = ?";
        public static final String ADD_UNSUBSCRIPTION =
                "INSERT INTO WEBHOOKS_UNSUBSCRIPTION (API_UUID, APPLICATION_ID, TENANT_DOMAIN, " +
                        "HUB_CALLBACK_URL, HUB_TOPIC, HUB_SECRET, HUB_LEASE_SECONDS, ADDED_AT) VALUES (?,?,?,?,?,?,?,?)";
        public static final String GET_ALL_VALID_SUBSCRIPTIONS =
                "SELECT WH.API_UUID AS API_UUID, " +
                        "WH.APPLICATION_ID AS APPLICATION_ID, " +
                        "WH.HUB_CALLBACK_URL AS HUB_CALLBACK_URL, " +
                        "WH.HUB_TOPIC AS HUB_TOPIC, " +
                        "WH.HUB_SECRET AS HUB_SECRET, " +
                        "WH.EXPIRY_AT AS EXPIRY_AT, " +
                        "API.CONTEXT AS API_CONTEXT, "  +
                        "API.API_VERSION AS API_VERSION, "  +
                        "API.API_TIER AS API_TIER, "  +
                        "API.API_ID AS API_ID, "  +
                        "SUB.TIER_ID AS SUB_TIER, " +
                        "APP.APPLICATION_TIER AS APPLICATION_TIER, " +
                        "SUBSCRIBER.USER_ID AS SUBSCRIBER, " +
                        "SUBSCRIBER.ORGANIZATION AS ORGANIZATION " +
                        "FROM WEBHOOKS_SUBSCRIPTION WH, " +
                        "API API, " +
                        "SUBSCRIPTION SUB, " +
                        "APPLICATION APP, " +
                        "SUBSCRIBER SUBSCRIBER " +
                        "WHERE WH.EXPIRY_AT >= ? AND WH.TENANT_DOMAIN = ? " +
                        "AND API.API_ID = SUB.API_ID " +
                        "AND WH.APPLICATION_ID = SUB.APPLICATION_ID " +
                        "AND API.API_UUID = WH.API_UUID " +
                        "AND APP.SUBSCRIBER_ID = SUBSCRIBER.SUBSCRIBER_ID ";

        public static final String GET_ALL_VALID_SUBSCRIPTIONS_POSTGRE_SQL =
                "SELECT WH.API_UUID AS API_UUID, " +
                        "WH.APPLICATION_ID AS APPLICATION_ID, " +
                        "WH.HUB_CALLBACK_URL AS HUB_CALLBACK_URL, " +
                        "WH.HUB_TOPIC AS HUB_TOPIC, " +
                        "WH.HUB_SECRET AS HUB_SECRET, " +
                        "WH.EXPIRY_AT AS EXPIRY_AT, " +
                        "API.CONTEXT AS API_CONTEXT, "  +
                        "API.API_VERSION AS API_VERSION, "  +
                        "API.API_TIER AS API_TIER, "  +
                        "API.API_ID AS API_ID, "  +
                        "SUB.TIER_ID AS SUB_TIER, " +
                        "APP.APPLICATION_TIER AS APPLICATION_TIER, " +
                        "SUBSCRIBER.USER_ID AS SUBSCRIBER, " +
                        "SUBSCRIBER.ORGANIZATION AS ORGANIZATION " +
                        "FROM WEBHOOKS_SUBSCRIPTION WH, " +
                        "API API, " +
                        "SUBSCRIPTION SUB, " +
                        "APPLICATION APP, " +
                        "SUBSCRIBER SUBSCRIBER " +
                        "WHERE WH.EXPIRY_AT >= ? AND WH.TENANT_DOMAIN = ? " +
                        "AND API.API_ID = SUB.API_ID " +
                        "AND WH.APPLICATION_ID::integer = SUB.APPLICATION_ID " +
                        "AND API.API_UUID = WH.API_UUID " +
                        "AND APP.SUBSCRIBER_ID = SUBSCRIBER.SUBSCRIBER_ID ";

        public static final String UPDATE_DELIVERY_STATE =
                "UPDATE WEBHOOKS_SUBSCRIPTION SET DELIVERED_AT = ?, DELIVERY_STATE = ? WHERE API_UUID = ? AND " +
                        "APPLICATION_ID = ? AND TENANT_DOMAIN = ? AND HUB_CALLBACK_URL = ? AND HUB_TOPIC = ?";
        public static final String GET_THROTTLE_LIMIT =
                "SELECT CONNECTIONS_COUNT FROM BUSINESS_PLAN WHERE NAME = ? AND ORGANIZATION = ?";
        public static final String GET_CURRENT_CONNECTIONS_COUNT =
                " SELECT COUNT(*) AS SUB_COUNT " +
                        " FROM " +
                        " WEBHOOKS_SUBSCRIPTION" +
                        " WHERE API_UUID = ? " +
                        " AND APPLICATION_ID = ?" +
                        " AND TENANT_DOMAIN = ?";
    }
    public static class KeyManagerSqlConstants {
        public static final String ADD_KEY_MANAGER =
                " INSERT INTO KEY_MANAGER (UUID,NAME,DESCRIPTION,TYPE,CONFIGURATION,ORGANIZATION,ENABLED," +
                        "DISPLAY_NAME,TOKEN_TYPE,EXTERNAL_REFERENCE_ID) VALUES (?,?,?,?,?,?,?,?,?,?)";
        public static final String UPDATE_KEY_MANAGER =
                "UPDATE KEY_MANAGER SET NAME = ?,DESCRIPTION = ?,TYPE = ?,CONFIGURATION = ?,ORGANIZATION = ?," +
                        "ENABLED = ?,DISPLAY_NAME = ?,TOKEN_TYPE = ?, EXTERNAL_REFERENCE_ID = ? WHERE UUID = ?";

        public static final String DELETE_KEY_MANAGER =
                "DELETE FROM KEY_MANAGER WHERE UUID = ? AND ORGANIZATION = ?";
    }

    /**
     * Static class to hold database queries related to TENANT_THEMES table
     */
    public static class TenantThemeConstants {

        public static final String ADD_TENANT_THEME = "INSERT INTO TENANT_THEMES (ORGANIZATION, THEME) VALUES (?,?)";
        public static final String UPDATE_TENANT_THEME = "UPDATE TENANT_THEMES SET THEME = ? WHERE TENANT_ID = ?";
        public static final String DELETE_TENANT_THEME = "DELETE FROM TENANT_THEMES WHERE ORGANIZATION = ?";
        public static final String GET_TENANT_THEME = "SELECT * FROM TENANT_THEMES WHERE ORGANIZATION = ?";
    }

    public static final String GET_API_VERSIONS =
            "SELECT API.API_VERSION FROM API API WHERE API.API_PROVIDER = ? AND API.API_NAME = ? AND ORGANIZATION = ?";
    public static final String GET_API_VERSIONS_UUID =
            "SELECT API.API_UUID, API.STATUS, API.API_VERSION, API.API_TYPE, API.VERSION_COMPARABLE, API.CONTEXT, "
                    + "API.CONTEXT_TEMPLATE FROM API API WHERE API.API_PROVIDER = ? AND API.API_NAME = ? ";
    public static class APIRevisionSqlConstants {
        public static final String ADD_API_REVISION =
                " INSERT INTO REVISION (ID, API_UUID, REVISION_UUID, DESCRIPTION, CREATED_BY, CREATED_TIME)" +
                        " VALUES (?,?,?,?,?,?)";
        public static final String GET_URL_MAPPINGS_WITH_SCOPE_AND_PRODUCT_ID = "SELECT AUM.HTTP_METHOD, AUM.AUTH_SCHEME, " +
                "AUM.URL_PATTERN, AUM.THROTTLING_TIER, AUM.MEDIATION_SCRIPT, ARSM.SCOPE_NAME, PROD_MAP.API_ID " +
                "FROM API_URL_MAPPING AUM LEFT JOIN API_RESOURCE_SCOPE_MAPPING ARSM ON AUM.URL_MAPPING_ID = ARSM.URL_MAPPING_ID " +
                "LEFT JOIN API_PRODUCT_MAPPING PROD_MAP ON AUM.URL_MAPPING_ID = PROD_MAP.URL_MAPPING_ID " +
                "WHERE AUM.API_ID = ? AND AUM.REVISION_UUID IS NULL";
        public static final String GET_REVISIONED_URL_MAPPINGS_ID = "SELECT URL_MAPPING_ID FROM API_URL_MAPPING " +
                "WHERE API_ID = ? AND REVISION_UUID = ? AND HTTP_METHOD = ? AND AUTH_SCHEME = ? AND URL_PATTERN = ? " +
                "AND THROTTLING_TIER = ? ";
        public static final String GET_URL_MAPPINGS_ID = "SELECT URL_MAPPING_ID FROM API_URL_MAPPING " +
                "WHERE API_ID = ? AND HTTP_METHOD = ? AND AUTH_SCHEME = ? AND URL_PATTERN = ? " +
                "AND THROTTLING_TIER = ? AND REVISION_UUID = ?";
        public static final String INSERT_URL_MAPPINGS = "INSERT INTO API_URL_MAPPING(API_ID, HTTP_METHOD," +
                " AUTH_SCHEME, URL_PATTERN, THROTTLING_TIER, REVISION_UUID) VALUES(?,?,?,?,?,?)";
        public static final String GET_CLIENT_CERTIFICATES = "SELECT ALIAS, CERTIFICATE," +
                " TIER_NAME FROM API_CLIENT_CERTIFICATE WHERE API_ID = ? AND REVISION_UUID='Current API'";
        public static final String INSERT_CLIENT_CERTIFICATES = "INSERT INTO API_CLIENT_CERTIFICATE(TENANT_ID, " +
                "ALIAS, API_ID, CERTIFICATE, REMOVED, TIER_NAME, REVISION_UUID) VALUES(?,?,?,?,?,?,?)";
        public static final String GET_GRAPHQL_COMPLEXITY = "SELECT TYPE, FIELD, COMPLEXITY_VALUE " +
                "FROM GRAPHQL_COMPLEXITY WHERE API_ID = ? AND REVISION_UUID IS NULL";
        public static final String INSERT_GRAPHQL_COMPLEXITY = "INSERT INTO GRAPHQL_COMPLEXITY(UUID, API_ID, TYPE," +
                " FIELD, COMPLEXITY_VALUE, REVISION_UUID) VALUES(?,?,?,?,?,?)";
        public static final String INSERT_SCOPE_RESOURCE_MAPPING = "INSERT INTO API_RESOURCE_SCOPE_MAPPING" +
                "(SCOPE_NAME, URL_MAPPING_ID, TENANT_ID) VALUES (?, ?, ?)";
        public static final String INSERT_PRODUCT_RESOURCE_MAPPING = "INSERT INTO API_PRODUCT_MAPPING" +
                "(API_ID, URL_MAPPING_ID) VALUES (?, ?)";
        public static final String INSERT_PRODUCT_REVISION_RESOURCE_MAPPING = "INSERT INTO API_PRODUCT_MAPPING" +
                "(API_ID, URL_MAPPING_ID, REVISION_UUID) VALUES (?, ?, ?)";
        public static final String DELETE_API_REVISION =
                "DELETE FROM REVISION WHERE REVISION_UUID = ?";
        public static final String GET_REVISION_COUNT_BY_API_UUID = "SELECT COUNT(ID) FROM REVISION WHERE API_UUID = ?";
        public static final String GET_MOST_RECENT_REVISION_ID = "SELECT REVISIONS_CREATED FROM API WHERE API_UUID" +
                " = ?";
        public static final String GET_REVISION_BY_REVISION_UUID = "SELECT * FROM REVISION WHERE REVISION_UUID = ?";
        public static final String GET_REVISION_UUID = "SELECT REVISION_UUID FROM REVISION WHERE API_UUID = ? " +
                "AND ID = ?";

        public static final String GET_REVISION_UUID_BY_ORGANIZATION = "SELECT REVISION_UUID FROM REVISION, API WHERE REVISION.API_UUID = ? " +
                "AND REVISION.ID = ? AND API.ORGANIZATION = ?";
        public static final String GET_EARLIEST_REVISION_ID = "SELECT REVISION_UUID FROM REVISION WHERE API_UUID = ? " +
                "ORDER BY ID ASC FETCH NEXT 1 ROWS ONLY";
        public static final String GET_EARLIEST_REVISION_ID_MYSQL = "SELECT REVISION_UUID FROM REVISION WHERE API_UUID = ? " +
                "ORDER BY ID ASC LIMIT 1";
        public static final String GET_EARLIEST_REVISION_ID_MSSQL =
                "SELECT TOP 1 REVISION_UUID FROM REVISION WHERE API_UUID = ? " + "ORDER BY ID ASC";
        public static final String GET_MOST_RECENT_REVISION_UUID = "SELECT REVISION_UUID FROM REVISION WHERE " +
                "API_UUID = ? ORDER BY ID DESC FETCH NEXT 1 ROWS ONLY";
        public static final String GET_MOST_RECENT_REVISION_UUID_MSSQL = "SELECT TOP 1 REVISION_UUID FROM REVISION WHERE " +
                "API_UUID = ? ORDER BY ID DESC";
        public static final String GET_MOST_RECENT_REVISION_UUID_MYSQL = "SELECT REVISION_UUID FROM REVISION WHERE " +
                "API_UUID = ? ORDER BY ID DESC LIMIT 1";
        public static final String GET_REVISION_APIID_BY_REVISION_UUID = "SELECT API_UUID, ID FROM REVISION WHERE REVISION_UUID = ?";
        public static final String GET_REVISIONS_BY_API_UUID = "SELECT ID, REVISION_UUID, DESCRIPTION, CREATED_TIME, " +
                "CREATED_BY FROM REVISION WHERE API_UUID = ? ORDER BY ID";
        public static final String ADD_API_REVISION_DEPLOYMENT_MAPPING =
                " INSERT INTO DEPLOYMENT_REVISION_MAPPING (NAME, VHOST, REVISION_UUID, DISPLAY_ON_DEVPORTAL, DEPLOYED_TIME)" +
                        " VALUES (?,?,?,?,?)";
        public static final String ADD_DEPLOYED_API_REVISION =
                "INSERT INTO DEPLOYED_REVISION (NAME, VHOST, REVISION_UUID, DEPLOYED_TIME)" +
                        " VALUES (?,?,?,?)";
        public static final String DELETE_API_REVISION_DEPLOYMENTS_MAPPING_BY_REVISION_UUID =
                " DELETE FROM DEPLOYMENT_REVISION_MAPPING WHERE REVISION_UUID = ?";
        public static final String GET_API_REVISION_DEPLOYMENT_MAPPING_BY_NAME_AND_REVISION_UUID
                = "SELECT * FROM DEPLOYMENT_REVISION_MAPPING WHERE NAME = ? AND REVISION_UUID = ? ";
        public static final String GET_API_REVISION_DEPLOYMENT_MAPPING_BY_REVISION_UUID
                = "SELECT * FROM DEPLOYMENT_REVISION_MAPPING WHERE REVISION_UUID = ?";
        static final String GET_API_REVISION_DEPLOYMENTS
                = "(SELECT NAME, VHOST, REVISION_UUID, DEPLOYED_TIME, 0 AS DISPLAY_ON_DEVPORTAL, NULL AS DEPLOY_TIME " +
                "FROM DEPLOYED_REVISION DR " +
                "UNION " +
                "SELECT NAME, VHOST, REVISION_UUID, NULL AS DEPLOYED_TIME, DISPLAY_ON_DEVPORTAL, " +
                "DEPLOYED_TIME AS DEPLOY_TIME " +
                "FROM DEPLOYMENT_REVISION_MAPPING DRM) ";
        public static final String GET_API_REVISION_DEPLOYMENTS_BY_API_UUID
                = "SELECT * FROM " + GET_API_REVISION_DEPLOYMENTS + "AD " +
                "WHERE AD.REVISION_UUID " +
                "IN " +
                "(SELECT REVISION_UUID FROM REVISION WHERE API_UUID = ?)";
        static final String GET_API_REVISION_DEPLOYMENTS_POSTGRES
                = "(SELECT NAME, VHOST, REVISION_UUID, DEPLOYED_TIME, false AS DISPLAY_ON_DEVPORTAL, NULL AS DEPLOY_TIME " +
                "FROM DEPLOYED_REVISION DR " +
                "UNION " +
                "SELECT NAME, VHOST, REVISION_UUID, NULL AS DEPLOYED_TIME, DISPLAY_ON_DEVPORTAL, " +
                "DEPLOYED_TIME AS DEPLOY_TIME " +
                "FROM DEPLOYMENT_REVISION_MAPPING DRM) ";
        public static final String GET_API_REVISION_DEPLOYMENTS_BY_API_UUID_POSTGRES
                = "SELECT * FROM " + GET_API_REVISION_DEPLOYMENTS_POSTGRES + "AD " +
                "WHERE AD.REVISION_UUID " +
                "IN " +
                "(SELECT REVISION_UUID FROM REVISION WHERE API_UUID = ?)";
        public static final String GET_API_REVISION_DEPLOYMENT_MAPPING_BY_API_UUID
                = "SELECT * FROM DEPLOYMENT_REVISION_MAPPING ADRM LEFT JOIN REVISION AR ON " +
                "ADRM.REVISION_UUID = AR.REVISION_UUID WHERE AR.API_UUID = ?";
        public static final String GET_DEPLOYED_REVISION_BY_API_UUID
                = "SELECT * FROM DEPLOYED_REVISION ADRM LEFT JOIN REVISION AR ON " +
                "ADRM.REVISION_UUID = AR.REVISION_UUID WHERE AR.API_UUID = ?";
        public static final String REMOVE_API_REVISION_DEPLOYMENT_MAPPING =
                " DELETE FROM DEPLOYMENT_REVISION_MAPPING WHERE NAME = ? AND REVISION_UUID = ?";
        public static final String REMOVE_DEPLOYED_API_REVISION =
                " DELETE FROM DEPLOYED_REVISION WHERE NAME = ? AND REVISION_UUID = ?";
        public static final String SET_UN_DEPLOYED_API_REVISION =
                "UPDATE DEPLOYED_REVISION SET DEPLOYED_TIME = NULL WHERE NAME = ? AND REVISION_UUID = ?";
        public static final String UPDATE_API_REVISION_DEPLOYMENT_MAPPING =
                " UPDATE DEPLOYMENT_REVISION_MAPPING SET DISPLAY_ON_DEVPORTAL = ? WHERE NAME = ? AND REVISION_UUID = ? ";
        public static final String REMOVE_CURRENT_API_ENTRIES_IN_API_URL_MAPPING_BY_API_ID =
                "DELETE FROM API_URL_MAPPING WHERE API_ID = ? AND REVISION_UUID IS NULL";
        public static final String REMOVE_CURRENT_API_PRODUCT_ENTRIES_IN_API_URL_MAPPING =
                "DELETE FROM API_URL_MAPPING WHERE REVISION_UUID = ?";
        public static final String GET_URL_MAPPINGS_WITH_SCOPE_AND_PRODUCT_ID_BY_REVISION_UUID = "SELECT AUM.HTTP_METHOD, AUM.AUTH_SCHEME, " +
                "AUM.URL_PATTERN, AUM.THROTTLING_TIER, AUM.MEDIATION_SCRIPT, ARSM.SCOPE_NAME, PROD_MAP.API_ID " +
                "FROM API_URL_MAPPING AUM LEFT JOIN API_RESOURCE_SCOPE_MAPPING ARSM ON AUM.URL_MAPPING_ID = ARSM.URL_MAPPING_ID " +
                "LEFT JOIN API_PRODUCT_MAPPING PROD_MAP ON AUM.URL_MAPPING_ID = PROD_MAP.URL_MAPPING_ID " +
                "WHERE AUM.API_ID = ? AND AUM.REVISION_UUID = ?";
        public static final String INSERT_URL_MAPPINGS_CURRENT_API = "INSERT INTO API_URL_MAPPING(API_ID, HTTP_METHOD," +
                " AUTH_SCHEME, URL_PATTERN, THROTTLING_TIER) VALUES(?,?,?,?,?)";
        public static final String GET_CURRENT_API_URL_MAPPINGS_ID = "SELECT URL_MAPPING_ID FROM API_URL_MAPPING " +
                "WHERE API_ID = ? AND REVISION_UUID IS NULL AND HTTP_METHOD = ? AND AUTH_SCHEME = ? AND URL_PATTERN = ? " +
                "AND THROTTLING_TIER = ? ";
        public static final String REMOVE_CURRENT_API_ENTRIES_IN_API_CLIENT_CERTIFICATE_BY_API_ID =
                "DELETE FROM API_CLIENT_CERTIFICATE WHERE API_ID = ? AND REVISION_UUID='Current API'";
        public static final String GET_CLIENT_CERTIFICATES_BY_REVISION_UUID = "SELECT ALIAS, CERTIFICATE," +
                " TIER_NAME FROM API_CLIENT_CERTIFICATE WHERE API_ID = ? AND REVISION_UUID = ?";
        public static final String INSERT_CLIENT_CERTIFICATES_AS_CURRENT_API = "INSERT INTO API_CLIENT_CERTIFICATE(TENANT_ID, " +
                "ALIAS, API_ID, CERTIFICATE, REMOVED, TIER_NAME, REVISION_UUID) VALUES(?,?,?,?,?,?,?)";
        public static final String REMOVE_CURRENT_API_ENTRIES_IN_GRAPHQL_COMPLEXITY_BY_API_ID =
                "DELETE FROM GRAPHQL_COMPLEXITY WHERE API_ID = ? AND REVISION_UUID IS NULL";
        public static final String GET_GRAPHQL_COMPLEXITY_BY_REVISION_UUID = "SELECT TYPE, FIELD, COMPLEXITY_VALUE " +
                "FROM GRAPHQL_COMPLEXITY WHERE API_ID = ? AND REVISION_UUID = ?";
        public static final String INSERT_GRAPHQL_COMPLEXITY_AS_CURRENT_API = "INSERT INTO GRAPHQL_COMPLEXITY(UUID, API_ID, TYPE," +
                " FIELD, COMPLEXITY_VALUE) VALUES(?,?,?,?,?)";
        public static final String REMOVE_REVISION_ENTRIES_IN_API_URL_MAPPING_BY_REVISION_UUID =
                "DELETE FROM API_URL_MAPPING WHERE API_ID = ? AND REVISION_UUID = ?";
        public static final String REMOVE_PRODUCT_REVISION_ENTRIES_IN_API_URL_MAPPING_BY_REVISION_UUID =
                "DELETE FROM API_URL_MAPPING WHERE REVISION_UUID = ?";
        public static final String REMOVE_REVISION_ENTRIES_IN_API_CLIENT_CERTIFICATE_BY_REVISION_UUID =
                "DELETE FROM API_CLIENT_CERTIFICATE WHERE API_ID = ? AND REVISION_UUID = ?";
        public static final String REMOVE_REVISION_ENTRIES_IN_GRAPHQL_COMPLEXITY_BY_REVISION_UUID =
                "DELETE FROM GRAPHQL_COMPLEXITY WHERE API_ID = ? AND REVISION_UUID = ?";
        public static final String REMOVE_CURRENT_API_ENTRIES_IN_API_PRODUCT_MAPPING_BY_API_PRODUCT_ID =
                "DELETE FROM API_PRODUCT_MAPPING WHERE API_ID = ? AND REVISION_UUID='Current API'";
        public static final String GET_PRODUCT_RESOURCES_BY_REVISION_UUID = "SELECT URL_MAPPING_ID " +
                "FROM API_PRODUCT_MAPPING WHERE API_ID = ? AND REVISION_UUID = ?";
        public static final String REMOVE_REVISION_ENTRIES_IN_API_PRODUCT_MAPPING_BY_REVISION_UUID =
                "DELETE FROM API_PRODUCT_MAPPING WHERE API_ID = ? AND REVISION_UUID = ?";
        public static final String GET_URL_MAPPINGS_WITH_SCOPE_AND_PRODUCT_ID_BY_PRODUCT_ID = "SELECT AUM.HTTP_METHOD, AUM.AUTH_SCHEME, " +
                "AUM.URL_PATTERN, AUM.THROTTLING_TIER, AUM.MEDIATION_SCRIPT, ARSM.SCOPE_NAME, AUM.API_ID " +
                "FROM API_URL_MAPPING AUM LEFT JOIN API_RESOURCE_SCOPE_MAPPING ARSM ON AUM.URL_MAPPING_ID = ARSM.URL_MAPPING_ID " +
                "LEFT JOIN API_PRODUCT_MAPPING PROD_MAP ON AUM.URL_MAPPING_ID = PROD_MAP.URL_MAPPING_ID " +
                "WHERE PROD_MAP.API_ID = ? AND PROD_MAP.REVISION_UUID = 'Current API'";
        public static final String GET_URL_MAPPINGS_WITH_SCOPE_BY_URL_MAPPING_ID = "SELECT AUM.HTTP_METHOD, AUM.AUTH_SCHEME, " +
                "AUM.URL_PATTERN, AUM.THROTTLING_TIER, AUM.MEDIATION_SCRIPT, ARSM.SCOPE_NAME, AUM.API_ID " +
                "FROM API_URL_MAPPING AUM LEFT JOIN API_RESOURCE_SCOPE_MAPPING ARSM ON AUM.URL_MAPPING_ID = ARSM.URL_MAPPING_ID " +
                "WHERE AUM.URL_MAPPING_ID = ?";
        public static final String GET_CUURENT_API_PRODUCT_RESOURCES = "SELECT URL_MAPPING_ID " +
                "FROM API_PRODUCT_MAPPING WHERE API_ID = ? AND REVISION_UUID = 'Current API'";
        public static final String REMOVE_PRODUCT_ENTRIES_IN_API_URL_MAPPING_BY_URL_MAPPING_ID =
                "DELETE FROM API_URL_MAPPING WHERE URL_MAPPING_ID = ?";
        public static final String CHECK_API_REVISION_DEPLOYMENT_AVAILABILITY_BY_API_UUID = "SELECT 1 FROM " +
                "DEPLOYMENT_REVISION_MAPPING WHERE REVISION_UUID IN  (SELECT REVISION_UUID FROM REVISION WHERE " +
                "API_UUID = ?)";
        public static final String GET_API_PRODUCT_REVISION_URL_MAPPINGS_BY_REVISION_UUID =
                "SELECT AUM.HTTP_METHOD, AUM.AUTH_SCHEME, " +
                "AUM.URL_PATTERN, AUM.THROTTLING_TIER, AUM.MEDIATION_SCRIPT, AUM.API_ID " +
                "FROM API_URL_MAPPING AUM " +
                "WHERE AUM.REVISION_UUID = ? ";

        public static final String GET_API_PRODUCT_REVISION_SCOPE_MAPPINGS_BY_REVISION_UUID =
                "SELECT AUM.HTTP_METHOD, AUM.URL_PATTERN, ARSM.SCOPE_NAME, ARSM.URL_MAPPING_ID " +
                "FROM API_RESOURCE_SCOPE_MAPPING ARSM LEFT JOIN API_URL_MAPPING AUM " +
                "ON ARSM.URL_MAPPING_ID = AUM.URL_MAPPING_ID WHERE AUM.REVISION_UUID = ?";
    }

    /**
     * Static class to hold database queries related to SERVICE_CATALOG table
     */
    public static class ServiceCatalogConstants {

        public static final String ADD_SERVICE = "INSERT INTO SERVICE_CATALOG " +
                "(UUID, SERVICE_KEY, MD5, SERVICE_NAME, SERVICE_VERSION, TENANT_ID, SERVICE_URL, " +
                "DEFINITION_TYPE, DEFINITION_URL, DESCRIPTION, " +
                "SECURITY_TYPE, MUTUAL_SSL_ENABLED, CREATED_TIME, LAST_UPDATED_TIME, CREATED_BY, UPDATED_BY, " +
                "SERVICE_DEFINITION) " +
                "VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)";

        public static final String UPDATE_SERVICE_BY_KEY = "UPDATE SERVICE_CATALOG SET " +
                "MD5 = ?," +
                "SERVICE_NAME = ?," +
                "TENANT_ID = ?," +
                "SERVICE_URL = ?," +
                "DEFINITION_URL = ?," +
                "DESCRIPTION = ?," +
                "SECURITY_TYPE = ?," +
                "MUTUAL_SSL_ENABLED = ?," +
                "LAST_UPDATED_TIME = ?," +
                "UPDATED_BY = ?," +
                "SERVICE_DEFINITION = ? " +
                "WHERE SERVICE_KEY = ? AND TENANT_ID = ?";

        public static final String DELETE_SERVICE_BY_SERVICE_ID = "DELETE FROM SERVICE_CATALOG WHERE UUID = ? " +
                "AND TENANT_ID = ?";

        public static final String GET_SERVICE_BY_SERVICE_KEY = "SELECT UUID, SERVICE_NAME, SERVICE_KEY, MD5, " +
                "   SERVICE_VERSION, SERVICE_URL, DEFINITION_TYPE, DEFINITION_URL, DESCRIPTION, SECURITY_TYPE, " +
                "   MUTUAL_SSL_ENABLED, CREATED_TIME, LAST_UPDATED_TIME, CREATED_BY, UPDATED_BY, SERVICE_DEFINITION " +
                "   FROM SERVICE_CATALOG WHERE SERVICE_KEY = ? AND TENANT_ID = ?";

        public static final String GET_SERVICE_MD5_BY_NAME_AND_VERSION = "SELECT MD5 FROM SERVICE_CATALOG " +
                "WHERE SERVICE_NAME = ? AND SERVICE_VERSION = ? AND TENANT_ID = ?";

        public static final String GET_SERVICE_MD5_BY_SERVICE_KEY = "SELECT MD5 FROM SERVICE_CATALOG " +
                "WHERE SERVICE_KEY = ? AND TENANT_ID = ?";

        public static final String ADD_ENDPOINT_RESOURCES = "INSERT INTO SERVICE_CATALOG (UUID, SERVICE_DEFINITION," +
                " METADATA) VALUES (?,?,?)";

        public static final String GET_SERVICE_BY_NAME_AND_VERSION = "SELECT UUID, SERVICE_NAME, SERVICE_KEY, MD5," +
                " SERVICE_VERSION, SERVICE_URL, DEFINITION_TYPE, DEFINITION_URL, DESCRIPTION, SECURITY_TYPE," +
                " MUTUAL_SSL_ENABLED, CREATED_TIME, LAST_UPDATED_TIME, CREATED_BY, UPDATED_BY, SERVICE_DEFINITION " +
                " FROM SERVICE_CATALOG WHERE SERVICE_NAME = ? AND SERVICE_VERSION = ? AND TENANT_ID = ?";

        public static final String GET_SERVICE_BY_SERVICE_ID = "SELECT " +
                "   UUID, " +
                "   SERVICE_KEY," +
                "   MD5," +
                "   SERVICE_NAME," +
                "   SERVICE_VERSION," +
                "   SERVICE_URL," +
                "   DEFINITION_TYPE," +
                "   DEFINITION_URL," +
                "   DESCRIPTION," +
                "   SECURITY_TYPE," +
                "   MUTUAL_SSL_ENABLED," +
                "   CREATED_TIME," +
                "   LAST_UPDATED_TIME," +
                "   CREATED_BY," +
                "   UPDATED_BY," +
                "   SERVICE_DEFINITION" +
                "   FROM SERVICE_CATALOG WHERE UUID = ? " +
                "   AND TENANT_ID = ?";

        public static final String GET_USAGE_OF_SERVICES_BY_SERVICE_ID = "SELECT " +
                "   API.API_ID, " +
                "   API.API_UUID, " +
                "   API.API_NAME, " +
                "   API.API_VERSION, " +
                "   API.CONTEXT," +
                "   API.API_PROVIDER " +
                "   FROM API INNER JOIN API_SERVICE_MAPPING ON " +
                "   API_SERVICE_MAPPING.API_ID = API.API_ID " +
                "   WHERE SERVICE_KEY = ? AND TENANT_ID = ?";

        public static final String GET_SERVICE_KEY_BY_SERVICE_UUID = "SELECT SERVICE_KEY FROM SERVICE_CATALOG WHERE" +
                "   UUID = ? AND TENANT_ID = ?";

    }

    /**
     * Static class to hold database queries related to API_OPERATION_POLICY_MAPPING table
     */
    public static class OperationPolicyConstants {

        // API policy mapping
        public static final String ADD_API_OPERATION_POLICY_MAPPING =
                "INSERT INTO API_OPERATION_POLICY_MAPPING " +
                        " (URL_MAPPING_ID, POLICY_UUID, DIRECTION, PARAMETERS, POLICY_ORDER) " +
                        " VALUES (?,?,?,?,?)";

        public static final String DELETE_OPERATION_POLICY_BY_POLICY_ID =
                "DELETE FROM API_OPERATION_POLICY WHERE POLICY_UUID = ?";

        public static final String GET_OPERATION_POLICIES_BY_URI_TEMPLATE_ID =
                "SELECT " +
                        " OP.POLICY_NAME, OP.POLICY_VERSION, OPM.DIRECTION, OPM.PARAMETERS, OPM.POLICY_ORDER, OPM.POLICY_UUID " +
                        " FROM " +
                        " API_URL_MAPPING AUM " +
                        " INNER JOIN API_OPERATION_POLICY_MAPPING OPM ON AUM.URL_MAPPING_ID = OPM.URL_MAPPING_ID" +
                        " INNER JOIN OPERATION_POLICY OP ON OPM.POLICY_UUID = OP.POLICY_UUID " +
                        " AND " +
                        " AUM.URL_MAPPING_ID = OPM.URL_MAPPING_ID " +
                        " WHERE " +
                        " AUM.URL_MAPPING_ID = ?";

        public static final String GET_OPERATION_POLICIES_OF_API_SQL =
                " SELECT " +
                        " AUM.URL_MAPPING_ID, AUM.URL_PATTERN, AUM.HTTP_METHOD," +
                        " OP.POLICY_NAME, OP.POLICY_VERSION, OPM.PARAMETERS, OPM.DIRECTION, OPM.POLICY_ORDER, OPM.POLICY_UUID" +
                        " FROM " +
                        " API_URL_MAPPING AUM " +
                        " INNER JOIN API API ON AUM.API_ID = API.API_ID " +
                        " INNER JOIN API_OPERATION_POLICY_MAPPING OPM ON AUM.URL_MAPPING_ID = OPM.URL_MAPPING_ID" +
                        " INNER JOIN OPERATION_POLICY OP ON OPM.POLICY_UUID = OP.POLICY_UUID " +
                        " INNER JOIN API_OPERATION_POLICY AOP ON OPM.POLICY_UUID = AOP.POLICY_UUID " +
                        " WHERE " +
                        " API.API_ID = ? " +
                        " AND " +
                        " AUM.REVISION_UUID IS NULL " +
                        " ORDER BY AUM.URL_MAPPING_ID ASC ";

        public static final String GET_OPERATION_POLICIES_FOR_API_REVISION_SQL =
                " SELECT " +
                        " AUM.URL_MAPPING_ID, AUM.URL_PATTERN, AUM.HTTP_METHOD," +
                        " OP.POLICY_NAME, OP.POLICY_VERSION, OPM.PARAMETERS, OPM.DIRECTION, OPM.POLICY_ORDER, OPM.POLICY_UUID" +
                        " FROM " +
                        " API_URL_MAPPING AUM " +
                        " INNER JOIN API API ON AUM.API_ID = API.API_ID " +
                        " INNER JOIN API_OPERATION_POLICY_MAPPING OPM ON AUM.URL_MAPPING_ID = OPM.URL_MAPPING_ID" +
                        " INNER JOIN OPERATION_POLICY OP ON OPM.POLICY_UUID = OP.POLICY_UUID " +
                        " INNER JOIN API_OPERATION_POLICY AOP ON OPM.POLICY_UUID = AOP.POLICY_UUID " +
                        " WHERE " +
                        " API.API_ID = ? " +
                        " AND " +
                        " AUM.REVISION_UUID = ? " +
                        " ORDER BY AUM.URL_MAPPING_ID ASC ";

        public static final String GET_OPERATION_POLICIES_PER_API_PRODUCT_SQL =
                " SELECT " +
                        " AUM.URL_MAPPING_ID, AUM.URL_PATTERN, AUM.HTTP_METHOD, " +
                        " OP.POLICY_NAME, OP.POLICY_VERSION, OPM.PARAMETERS, OPM.DIRECTION, OPM.POLICY_ORDER, OPM.POLICY_UUID " +
                        " FROM " +
                        " API_URL_MAPPING AUM " +
                        " INNER JOIN API_OPERATION_POLICY_MAPPING OPM ON AUM.URL_MAPPING_ID = OPM.URL_MAPPING_ID " +
                        " INNER JOIN OPERATION_POLICY OP ON OPM.POLICY_UUID = OP.POLICY_UUID " +
                        " WHERE " +
                        " AUM.REVISION_UUID = ? " +
                        " ORDER BY AUM.URL_MAPPING_ID ASC ";

        public static final String ADD_COMMON_OPERATION_POLICY =
                "INSERT INTO COMMON_OPERATION_POLICY (POLICY_UUID) VALUES (?)";

        public static final String ADD_API_SPECIFIC_OPERATION_POLICY =
                "INSERT INTO API_OPERATION_POLICY(POLICY_UUID, API_UUID, CLONED_POLICY_UUID) VALUES (?,?,?)";

        public static final String ADD_API_SPECIFIC_OPERATION_POLICY_WITH_REVISION =
                "INSERT INTO API_OPERATION_POLICY (POLICY_UUID, API_UUID, CLONED_POLICY_UUID, REVISION_UUID) VALUES (?,?,?,?)";

        public static final String ADD_OPERATION_POLICY =
                "INSERT INTO OPERATION_POLICY " +
                        " (POLICY_UUID, POLICY_NAME, POLICY_VERSION, DISPLAY_NAME, POLICY_DESCRIPTION, APPLICABLE_FLOWS, GATEWAY_TYPES, " +
                        " API_TYPES, POLICY_PARAMETERS, ORGANIZATION, POLICY_CATEGORY, " +
                        " POLICY_MD5) " +
                        " VALUES (?,?,?,?,?,?,?,?,?,?,?,?)";

        public static final String UPDATE_API_OPERATION_POLICY_BY_POLICY_ID =
                "UPDATE  API_OPERATION_POLICY SET CLONED_POLICY_UUID = ?  WHERE  POLICY_UUID = ?";

        public static final String UPDATE_OPERATION_POLICY_CONTENT =
                "UPDATE " +
                        " OPERATION_POLICY " +
                        " SET " +
                        " POLICY_NAME = ?, POLICY_VERSION = ?, DISPLAY_NAME = ?, POLICY_DESCRIPTION = ?, APPLICABLE_FLOWS = ?, GATEWAY_TYPES = ?, " +
                        " API_TYPES = ?, POLICY_PARAMETERS = ?, ORGANIZATION = ?, POLICY_CATEGORY = ?, " +
                        " POLICY_MD5 = ? " +
                        " WHERE " +
                        " POLICY_UUID = ?";

        public static final String GET_OPERATION_POLICY_FROM_POLICY_ID =
                "SELECT " +
                        " POLICY_NAME, POLICY_VERSION, DISPLAY_NAME, POLICY_DESCRIPTION, APPLICABLE_FLOWS, GATEWAY_TYPES, API_TYPES, " +
                        " POLICY_PARAMETERS, POLICY_CATEGORY, POLICY_MD5, ORGANIZATION " +
                        " FROM " +
                        " OPERATION_POLICY " +
                        " WHERE " +
                        " POLICY_UUID = ?";

        public static final String GET_API_SPECIFIC_OPERATION_POLICY_FROM_POLICY_ID =
                "SELECT " +
                        " OP.POLICY_UUID, OP.POLICY_NAME, OP.POLICY_VERSION, OP.DISPLAY_NAME, OP.POLICY_DESCRIPTION, OP.APPLICABLE_FLOWS, OP.GATEWAY_TYPES, OP.API_TYPES, " +
                        " OP.POLICY_PARAMETERS, OP.POLICY_CATEGORY, OP.POLICY_MD5, " +
                        " AOP.API_UUID, AOP.REVISION_UUID, AOP.CLONED_POLICY_UUID " +
                        " FROM " +
                        " OPERATION_POLICY OP INNER JOIN API_OPERATION_POLICY AOP ON OP.POLICY_UUID = AOP.POLICY_UUID " +
                        " WHERE " +
                        " OP.POLICY_UUID = ? AND OP.ORGANIZATION = ? AND AOP.API_UUID = ?";

        public static final String GET_COMMON_OPERATION_POLICY_WITH_OUT_DEFINITION_FROM_POLICY_ID =
                "SELECT " +
                        " OP.POLICY_UUID, OP.POLICY_NAME, OP.POLICY_VERSION, OP.DISPLAY_NAME, OP.POLICY_DESCRIPTION, OP.APPLICABLE_FLOWS, OP.GATEWAY_TYPES, OP.API_TYPES, " +
                        " OP.POLICY_PARAMETERS, OP.POLICY_CATEGORY, OP.POLICY_MD5 " +
                        " FROM " +
                        " OPERATION_POLICY OP INNER JOIN COMMON_OPERATION_POLICY COP ON OP.POLICY_UUID = COP.POLICY_UUID " +
                        " WHERE " +
                        " OP.POLICY_UUID = ? AND OP.ORGANIZATION = ?";


        public static final String GET_API_SPECIFIC_OPERATION_POLICY_FROM_POLICY_NAME =
                "SELECT " +
                        " OP.POLICY_UUID, OP.POLICY_NAME, OP.POLICY_VERSION, OP.DISPLAY_NAME, OP.POLICY_DESCRIPTION, OP.APPLICABLE_FLOWS, OP.GATEWAY_TYPES, OP.API_TYPES, " +
                        " OP.POLICY_PARAMETERS, OP.POLICY_CATEGORY, OP.POLICY_MD5, " +
                        " AOP.API_UUID, AOP.REVISION_UUID, AOP.CLONED_POLICY_UUID " +
                        " FROM " +
                        " OPERATION_POLICY OP INNER JOIN API_OPERATION_POLICY AOP ON OP.POLICY_UUID = AOP.POLICY_UUID " +
                        " WHERE " +
                        " OP.POLICY_NAME = ? AND OP.POLICY_VERSION = ? AND OP.ORGANIZATION = ? AND AOP.API_UUID = ? ";

        public static final String GET_COMMON_OPERATION_POLICY_FROM_POLICY_NAME =
                "SELECT " +
                        " OP.POLICY_UUID, OP.POLICY_NAME, OP.POLICY_VERSION, OP.DISPLAY_NAME, OP.POLICY_DESCRIPTION, OP.APPLICABLE_FLOWS, OP.GATEWAY_TYPES, OP.API_TYPES, " +
                        " OP.POLICY_PARAMETERS, OP.POLICY_CATEGORY, OP.POLICY_MD5 " +
                        " FROM " +
                        " OPERATION_POLICY OP INNER JOIN COMMON_OPERATION_POLICY COP ON OP.POLICY_UUID = COP.POLICY_UUID " +
                        " WHERE " +
                        " OP.POLICY_NAME = ? AND OP.POLICY_VERSION = ? AND OP.ORGANIZATION = ?";

        public static final String GET_ALL_COMMON_OPERATION_POLICIES =
                "SELECT " +
                        " OP.POLICY_UUID, OP.POLICY_NAME, OP.POLICY_VERSION, OP.DISPLAY_NAME, OP.POLICY_DESCRIPTION, OP.APPLICABLE_FLOWS, OP.GATEWAY_TYPES, " +
                        " OP.API_TYPES, OP.POLICY_PARAMETERS, OP.POLICY_CATEGORY, OP.POLICY_MD5 " +
                        " FROM " +
                        " OPERATION_POLICY OP INNER JOIN COMMON_OPERATION_POLICY COP ON OP.POLICY_UUID = COP.POLICY_UUID " +
                        " WHERE " +
                        " OP.ORGANIZATION = ?";

        public static final String GET_ALL_API_SPECIFIC_OPERATION_POLICIES_WITHOUT_CLONED_POLICIES =
                "SELECT " +
                        " OP.POLICY_UUID, OP.POLICY_NAME, OP.POLICY_VERSION, OP.DISPLAY_NAME, OP.POLICY_DESCRIPTION, OP.APPLICABLE_FLOWS, OP.GATEWAY_TYPES, " +
                        " OP.API_TYPES, OP.POLICY_PARAMETERS, OP.POLICY_CATEGORY, OP.POLICY_MD5 " +
                        " FROM " +
                        " OPERATION_POLICY OP INNER JOIN API_OPERATION_POLICY AOP ON OP.POLICY_UUID = AOP.POLICY_UUID " +
                        " WHERE " +
                        " OP.ORGANIZATION = ? AND AOP.API_UUID = ? AND AOP.REVISION_UUID IS NULL " +
                        " AND AOP.CLONED_POLICY_UUID IS NULL";

        public static final String GET_EXISTING_POLICY_USAGES_BY_POLICY_UUID =
                "SELECT COUNT(POLICY_UUID) AS POLICY_COUNT FROM API_OPERATION_POLICY_MAPPING " +
                        " WHERE POLICY_UUID = ?";

        public static final String DELETE_OPERATION_POLICY_BY_ID =
                "DELETE FROM OPERATION_POLICY WHERE POLICY_UUID = ?";

        public static final String GET_CLONED_POLICY_ID_FOR_COMMON_POLICY_ID =
                "SELECT POLICY_UUID FROM API_OPERATION_POLICY " +
                        " WHERE " +
                        " CLONED_POLICY_UUID = ?  AND API_UUID = ? AND REVISION_UUID IS NULL";

        public static final String GET_ALL_API_SPECIFIC_POLICIES_FOR_API_ID =
                "SELECT POLICY_UUID FROM API_OPERATION_POLICY WHERE API_UUID = ?";

        public static final String GET_ALL_API_SPECIFIC_POLICIES_FOR_REVISION_UUID =
                "SELECT POLICY_UUID FROM API_OPERATION_POLICY WHERE API_UUID = ? AND REVISION_UUID = ?";

        public static final String GET_ALL_CLONED_POLICIES_FOR_API =
                "SELECT POLICY_UUID FROM API_OPERATION_POLICY WHERE API_UUID = ? AND REVISION_UUID IS NULL " +
                        " AND CLONED_POLICY_UUID IS NOT NULL";

        public static final String GET_OPERATION_POLICY_DEFINITION_FROM_POLICY_ID =
                "SELECT POLICY_DEFINITION, GATEWAY_TYPE, DEFINITION_MD5  FROM OPERATION_POLICY_DEFINITION " +
                        " WHERE POLICY_UUID = ? ";

        public static final String ADD_OPERATION_POLICY_DEFINITION =
                "INSERT INTO OPERATION_POLICY_DEFINITION " +
                        " (POLICY_UUID, GATEWAY_TYPE, DEFINITION_MD5, POLICY_DEFINITION) " +
                        " VALUES (?,?,?,?)";

        public static final String UPDATE_OPERATION_POLICY_DEFINITION =
                "UPDATE OPERATION_POLICY_DEFINITION SET DEFINITION_MD5 = ?, POLICY_DEFINITION = ? " +
                        " WHERE POLICY_UUID = ? AND GATEWAY_TYPE = ?";

        public static final String GET_COMMON_OPERATION_POLICY_NAMES_FOR_ORGANIZATION =
                "SELECT OP.POLICY_NAME, OP.POLICY_VERSION FROM OPERATION_POLICY OP INNER JOIN COMMON_OPERATION_POLICY COP " +
                        " ON OP.POLICY_UUID = COP.POLICY_UUID WHERE OP.ORGANIZATION = ?";
    }

    /**
     * Static class to hold database queries related to SYSTEM_CONFIGS table
     */
    public static class SystemConfigsConstants {
        public static final String ADD_SYSTEM_CONFIG_SQL = "INSERT INTO SYSTEM_CONFIGS "
                + "(ORGANIZATION,CONFIG_TYPE,CONFIGURATION) VALUES (?,?,?)";
        public static final String GET_SYSTEM_CONFIG_SQL = "SELECT CONFIGURATION FROM SYSTEM_CONFIGS "
                + "WHERE ORGANIZATION = ? AND CONFIG_TYPE = ?";
        public static final String UPDATE_SYSTEM_CONFIG_SQL = "UPDATE SYSTEM_CONFIGS "
                + "SET CONFIGURATION = ? WHERE ORGANIZATION = ? AND CONFIG_TYPE = ?";
    }

}
