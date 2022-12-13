/*
 *  Copyright (c) 2022, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */

package database

const (
	queryGetApplicationByUUID string = " SELECT " +
		"   APP.UUID," +
		"   APP.NAME," +
		"   APP.SUBSCRIBER_ID," +
		"   APP.ORGANIZATION ORGANIZATION," +
		"   SUB.USER_ID " +
		" FROM " +
		"   SUBSCRIBER SUB," +
		"   APPLICATION APP " +
		" WHERE " +
		"   APP.UUID = $1 " +
		"   AND APP.SUBSCRIBER_ID = SUB.SUBSCRIBER_ID"

	queryGetAllSubscriptionsForApplication string = "select " +
		"	SUB.uuid as UUID, " +
		"	API.api_uuid as API_UUID, " +
		"	API.api_version as API_VERSION, " +
		"	SUB.sub_status as SUB_STATUS, " +
		"	APP.organization as ORGANIZATION, " +
		"	SUB.created_by as CREATED_BY " +
		" FROM " +
		" APPLICATION APP, SUBSCRIPTION SUB, API API " +
		" where 1 = 1 " +
		"	AND APP.application_id = SUB.application_id " +
		"	AND SUB.api_id = API.api_id " +
		"	AND APP.uuid = $1"

	queryConsumerKeysForApplication string = "select " +
		"	APPKEY.consumer_key, " +
		"	APPKEY.key_manager " +
		" from " +
		"	application_key_mapping APPKEY, " +
		"	application APP " +
		" where 1=1 " +
		"	AND APP.application_id = APPKEY.application_id " +
		"	AND APP.UUID = $1"

	querySubscriptionByUUID string = "select " +
		"	SUB.uuid, " +
		"	API.api_uuid, " +
		"	SUB.sub_status, " +
		"	API.organization, " +
		"	SUB.created_by " +
		" from " +
		"	subscription SUB, " +
		"	api API " +
		" where 1=1 " +
		"	AND SUB.api_id = API.api_id " +
		"	AND SUB.uuid = $1"

	queryCreateAPI string = "INSERT INTO API " +
		"(API_UUID, API_NAME, API_PROVIDER, API_VERSION," +
		"CONTEXT, ORGANIZATION, CREATED_BY, CREATED_TIME, API_TYPE, ARTIFACT, STATUS)" +
		" VALUES " + "($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)"

	queryDeleteAPI string = "DELETE FROM API" +
		" WHERE " +
		"API_UUID = $1"

	queryUpdateAPI string = "UPDATE API SET " +
		"API_NAME = $2, " +
		"API_PROVIDER = $3, " +
		"API_VERSION = $4, " +
		"CONTEXT = $5, " +
		"ORGANIZATION = $6, " +
		"UPDATED_BY = $7, " +
		"UPDATED_TIME = $8, " +
		"API_TYPE = $9, " +
		"ARTIFACT = $10, " +
		"STATUS = $11" +
		" WHERE " +
		"API_UUID = $1"
)
