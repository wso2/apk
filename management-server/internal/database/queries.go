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
	QueryGetApplicationByUUID string = " SELECT " +
		"   APP.UUID," +
		"   APP.NAME," +
		"   APP.SUBSCRIBER_ID," +
		"   APP.ORGANIZATION ORGANIZATION," +
		"   SUB.USER_ID " +
		" FROM " +
		"   AM_SUBSCRIBER SUB," +
		"   AM_APPLICATION APP " +
		" WHERE " +
		"   APP.UUID = $1 " +
		"   AND APP.SUBSCRIBER_ID = SUB.SUBSCRIBER_ID"

	QueryGetAllSubscriptionsForApplication string = "select " +
		"	SUB.uuid as UUID, " +
		"	API.api_uuid as API_UUID, " +
		"	API.api_version as API_VERSION, " +
		"	SUB.sub_status as SUB_STATUS, " +
		"	APP.organization as ORGANIZATION, " +
		"	SUB.created_by as CREATED_BY " +
		" FROM " +
		" AM_APPLICATION APP, AM_SUBSCRIPTION SUB, AM_API API " +
		" where 1 = 1 " +
		"	AND APP.application_id = SUB.application_id " +
		"	AND SUB.api_id = API.api_id " +
		"	AND APP.uuid = $1"
)
