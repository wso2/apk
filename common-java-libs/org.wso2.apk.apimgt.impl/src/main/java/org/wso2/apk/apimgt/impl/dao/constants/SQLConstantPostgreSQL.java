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

/**
 * This class holds postgre sql queries.
 */
public class SQLConstantPostgreSQL extends SQLConstants{

    public static final String GET_APPLICATIONS_PREFIX_CASESENSITVE_WITHGROUPID =
            "select distinct x.*,bl.ENABLED from (" +
                    " SELECT " +
                    "   APPLICATION_ID, " +
                    "   NAME," +
                    "   APPLICATION_TIER," +
                    "   APP.SUBSCRIBER_ID,  " +
                    "   APP.CREATED_TIME AS APP_CREATED_TIME, " +
                    "   APP.UPDATED_TIME AS APP_UPDATED_TIME, " +
                    "   CALLBACK_URL,  " +
                    "   DESCRIPTION, " +
                    "   APPLICATION_STATUS, " +
                    "   USER_ID, " +
                    "   GROUP_ID, " +
                    "   UUID, " +
                    "   APP.CREATED_BY AS CREATED_BY " +
                    " FROM" +
                    "   APPLICATION APP, " +
                    "   SUBSCRIBER SUB  " +
                    " WHERE " +
                    "   SUB.SUBSCRIBER_ID = APP.SUBSCRIBER_ID " +
                    " AND " +
                    "   (GROUP_ID= ?  OR  (GROUP_ID='' AND SUB.USER_ID = ?))" +
                    " AND " +
                    "   APP.ORGANIZATION = ? " +
                    " And " +
                    "    NAME like ?" +
                    " ORDER BY $1 $2 " +
                    " offset ? limit  ? "+
                    " )x left join BLOCK_CONDITION bl on  ( bl.TYPE = 'APPLICATION' AND bl.BLOCK_CONDITION = concat(concat(x.USER_ID,':'),x.name)) ";


    public static final String GET_APPLICATIONS_PREFIX_NONE_CASESENSITVE_WITHGROUPID =
            "select distinct x.*,bl.ENABLED from (" +
                    "SELECT " +
                    "   APPLICATION_ID, " +
                    "   NAME," +
                    "   APPLICATION_TIER," +
                    "   APP.SUBSCRIBER_ID,  " +
                    "   APP.CREATED_TIME AS APP_CREATED_TIME, " +
                    "   APP.UPDATED_TIME AS APP_UPDATED_TIME, " +
                    "   CALLBACK_URL,  " +
                    "   DESCRIPTION, " +
                    "   APPLICATION_STATUS, " +
                    "   USER_ID, " +
                    "   GROUP_ID, " +
                    "   UUID, " +
                    "   APP.CREATED_BY AS CREATED_BY " +
                    " FROM" +
                    "   APPLICATION APP, " +
                    "   SUBSCRIBER SUB  " +
                    " WHERE " +
                    "   SUB.SUBSCRIBER_ID = APP.SUBSCRIBER_ID " +
                    " AND " +
                    "   (GROUP_ID= ?  OR (GROUP_ID='' AND LOWER (SUB.USER_ID) =LOWER (?)))"+
                    " AND " +
                    "   APP.ORGANIZATION = ? " +
                    " And "+
                    "    NAME like ?"+
                    " ORDER BY $1 $2 " +
                    " offset ? limit  ? "+
                    " )x left join BLOCK_CONDITION bl on  ( bl.TYPE = 'APPLICATION' AND bl.BLOCK_CONDITION = concat(concat(x.USER_ID,':'),x.name)) ";

    public static final String GET_APPLICATIONS_PREFIX_CASESENSITVE_WITH_MULTIGROUPID =
            "select distinct x.*,bl.ENABLED from (" +
                    " SELECT " +
                    "   APPLICATION_ID, " +
                    "   NAME," +
                    "   APPLICATION_TIER," +
                    "   APP.SUBSCRIBER_ID,  " +
                    "   APP.CREATED_TIME AS APP_CREATED_TIME, " +
                    "   APP.UPDATED_TIME AS APP_UPDATED_TIME, " +
                    "   CALLBACK_URL,  " +
                    "   DESCRIPTION, " +
                    "   APPLICATION_STATUS, " +
                    "   USER_ID, " +
                    "   GROUP_ID, " +
                    "   UUID, " +
                    "   APP.CREATED_BY AS CREATED_BY " +
                    " FROM" +
                    "   APPLICATION APP, " +
                    "   SUBSCRIBER SUB  " +
                    " WHERE " +
                    "   SUB.SUBSCRIBER_ID = APP.SUBSCRIBER_ID " +
                    " AND (" +
                    "    (APPLICATION_ID IN ( SELECT APPLICATION_ID FROM APPLICATION_GROUP_MAPPING WHERE GROUP_ID IN ($params) AND TENANT = ?)) " +
                    "           OR " +
                    "    SUB.USER_ID = ?" +
                    "           OR " +
                    "    (APP.APPLICATION_ID IN (SELECT APPLICATION_ID FROM APPLICATION WHERE GROUP_ID = ?))" +
                    " )" +
                    " AND " +
                    "   APP.ORGANIZATION = ? " +
                    " And "+
                    "    NAME like ?" +
                    " ORDER BY $1 $2 " +
                    " offset ? limit  ? "+
                    " )x left join BLOCK_CONDITION bl on  ( bl.TYPE = 'APPLICATION' AND bl.BLOCK_CONDITION = concat(concat(x.USER_ID,':'),x.name)) ";


    public static final String GET_APPLICATIONS_PREFIX_NONE_CASESENSITVE_WITH_MULTIGROUPID =
            "select distinct x.*,bl.ENABLED from (" +
                    "SELECT " +
                    "   APPLICATION_ID, " +
                    "   NAME," +
                    "   APPLICATION_TIER," +
                    "   APP.SUBSCRIBER_ID,  " +
                    "   APP.CREATED_TIME AS APP_CREATED_TIME, " +
                    "   APP.UPDATED_TIME AS APP_UPDATED_TIME, " +
                    "   CALLBACK_URL,  " +
                    "   DESCRIPTION, " +
                    "   APPLICATION_STATUS, " +
                    "   USER_ID, " +
                    "   GROUP_ID, " +
                    "   UUID, " +
                    "   APP.CREATED_BY AS CREATED_BY " +
                    " FROM" +
                    "   APPLICATION APP, " +
                    "   SUBSCRIBER SUB  " +
                    " WHERE " +
                    "   SUB.SUBSCRIBER_ID = APP.SUBSCRIBER_ID " +
                    " AND (" +
                    "    (APPLICATION_ID IN ( SELECT APPLICATION_ID FROM APPLICATION_GROUP_MAPPING WHERE GROUP_ID " +
                    " IN ($params) AND TENANT = ? ))" +
                    "           OR " +
                    "    (LOWER (SUB.USER_ID) = LOWER (?))" +
                    "           OR " +
                    "    (APP.APPLICATION_ID IN (SELECT APPLICATION_ID FROM APPLICATION WHERE GROUP_ID = ?))" +
                    " )" +
                    " AND " +
                    "   APP.ORGANIZATION = ? " +
                    " And " +
                    "    NAME like ?"+
                    " ORDER BY $1 $2 " +
                    " offset ? limit  ? "+
                    " )x left join BLOCK_CONDITION bl on  ( bl.TYPE = 'APPLICATION' AND bl.BLOCK_CONDITION = concat(concat(x.USER_ID,':'),x.name)) ";


    public static final String GET_APPLICATIONS_PREFIX_CASESENSITVE =
            "select distinct x.*,bl.ENABLED from (" +
                    "SELECT " +
                    "   APPLICATION_ID, " +
                    "   NAME," +
                    "   APPLICATION_TIER," +
                    "   APP.SUBSCRIBER_ID,  " +
                    "   APP.CREATED_TIME AS APP_CREATED_TIME, " +
                    "   APP.UPDATED_TIME AS APP_UPDATED_TIME, " +
                    "   CALLBACK_URL,  " +
                    "   DESCRIPTION, " +
                    "   APPLICATION_STATUS, " +
                    "   USER_ID, " +
                    "   GROUP_ID, " +
                    "   UUID, " +
                    "   APP.CREATED_BY AS CREATED_BY " +
                    " FROM" +
                    "   APPLICATION APP, " +
                    "   SUBSCRIBER SUB  " +
                    " WHERE " +
                    "   SUB.SUBSCRIBER_ID = APP.SUBSCRIBER_ID " +
                    " AND " +
                    "    SUB.USER_ID = ?"+
                    " AND " +
                    "   APP.ORGANIZATION = ? " +
                    " And "+
                    "    NAME like ?"+
                    " ORDER BY $1 $2 " +
                    " offset ? limit  ? "+
                    " )x left join BLOCK_CONDITION bl on  ( bl.TYPE = 'APPLICATION' AND bl.BLOCK_CONDITION = concat(concat(x.USER_ID,':'),x.name)) ";

    public static final String GET_APPLICATIONS_PREFIX_NONE_CASESENSITVE =
            "select distinct x.*,bl.ENABLED from (" +
                    "SELECT " +
                    "   APPLICATION_ID, " +
                    "   NAME," +
                    "   APPLICATION_TIER," +
                    "   APP.SUBSCRIBER_ID,  " +
                    "   APP.CREATED_TIME AS APP_CREATED_TIME, " +
                    "   APP.UPDATED_TIME AS APP_UPDATED_TIME, " +
                    "   CALLBACK_URL,  " +
                    "   DESCRIPTION, " +
                    "   APPLICATION_STATUS, " +
                    "   USER_ID, " +
                    "   GROUP_ID, " +
                    "   UUID, " +
                    "   APP.CREATED_BY AS CREATED_BY " +
                    " FROM" +
                    "   APPLICATION APP, " +
                    "   SUBSCRIBER SUB  " +
                    " WHERE " +
                    "   SUB.SUBSCRIBER_ID = APP.SUBSCRIBER_ID " +
                    " AND " +
                    "   LOWER (SUB.USER_ID) =LOWER (?)" +
                    " AND " +
                    "   APP.ORGANIZATION = ? " +
                    " And "+
                    "    NAME like ?"+
                    " ORDER BY $1 $2 " +
                    " offset ? limit  ? "+
                    " )x left join BLOCK_CONDITION bl on  ( bl.TYPE = 'APPLICATION' AND bl.BLOCK_CONDITION = concat(concat(x.USER_ID,':'),x.name)) ";

    public static final String GET_APPLICATIONS_BY_ORGANIZATION =
            "   SELECT " +
                    "   APP.APPLICATION_ID as APPLICATION_ID, " +
                    "   SUB.CREATED_BY AS CREATED_BY," +
                    "   APP.GROUP_ID AS GROUP_ID, " +
                    "   APP.CREATED_TIME AS APP_CREATED_TIME, " +
                    "   APP.UPDATED_TIME AS APP_UPDATED_TIME, " +
                    "   SUB.ORGANIZATION AS ORGANIZATION, " +
                    "   SUB.SUBSCRIBER_ID AS SUBSCRIBER_ID, " +
                    "   APP.UUID AS UUID," +
                    "   APP.NAME AS NAME," +
                    "   APP.APPLICATION_STATUS as APPLICATION_STATUS  " +
                    " FROM" +
                    "   APPLICATION APP, " +
                    "   SUBSCRIBER SUB  " +
                    " WHERE " +
                    "   SUB.SUBSCRIBER_ID = APP.SUBSCRIBER_ID " +
                    " AND " +
                    "    SUB.ORGANIZATION = ? "+
                    " And "+
                    "    ( SUB.CREATED_BY like ?"+
                    " AND APP.NAME like ?"+
                    " ) ORDER BY $1 $2 " +
                    " offset ? limit  ? ";

    public static final String GET_REPLIES_SQL =
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
                "AND PARENT_COMMENT_ID = ? " +
                "ORDER BY API_COMMENTS.CREATED_TIME ASC OFFSET ? LIMIT ?";

    public static final String GET_ROOT_COMMENTS_SQL =
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
                "AND PARENT_COMMENT_ID IS NULL " +
                "ORDER BY API_COMMENTS.CREATED_TIME DESC OFFSET ? LIMIT ?";
}
