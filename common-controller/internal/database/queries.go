/*
 *  Copyright (c) 2024, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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
	insertApplication     = "INSERT INTO APPLICATION (UUID, NAME, OWNER, ORGANIZATION) VALUES ($1, $2, $3, $4);"
	getAllApplications    = "SELECT UUID, NAME, OWNER, ORGANIZATION FROM APPLICATION;"
	updateApplication     = "UPDATE APPLICATION SET NAME = $2, OWNER = $3, ORGANIZATION = $4 WHERE UUID = $1;"
	deleteApplication     = "DELETE FROM APPLICATION WHERE UUID = $1;"
	deleteAllApplications = "DELETE FROM APPLICATION;"

	insertApplicationAttributes = "INSERT INTO APPLICATION_ATTRIBUTES (APPLICATION_UUID, NAME, APP_ATTRIBUTE) VALUES ($1, $2, $3);"
	getAllApplicationAttributes = "SELECT APPLICATION_UUID, NAME, APP_ATTRIBUTE FROM APPLICATION_ATTRIBUTES;"
	deleteApplicationAttributes = "DELETE FROM APPLICATION_ATTRIBUTES WHERE APPLICATION_UUID = $1;"
	deleteAllAppAttributes      = "DELETE FROM APPLICATION_ATTRIBUTES;"

	insertSubscription     = "INSERT INTO SUBSCRIPTION (UUID, API_NAME, API_VERSION, SUB_STATUS, ORGANIZATION) VALUES ($1, $2, $3, $4, $5);"
	getAllSubscriptions    = "SELECT UUID, API_NAME, API_VERSION, SUB_STATUS, ORGANIZATION FROM SUBSCRIPTION;"
	updateSubscription     = "UPDATE SUBSCRIPTION SET API_NAME = $2, API_VERSION = $3, SUB_STATUS = $4, ORGANIZATION = $5 WHERE UUID = $1;"
	deleteSubscription     = "DELETE FROM SUBSCRIPTION WHERE UUID = $1;"
	deleteAllSubscriptions = "DELETE FROM SUBSCRIPTION;"

	insertAppSub    = "INSERT INTO APPLICATION_SUBSCRIPTION_MAPPING (UUID, APPLICATION_UUID, SUBSCRIPTION_UUID, ORGANIZATION) VALUES ($1, $2, $3, $4);"
	getAllAppSubs   = "SELECT UUID, APPLICATION_UUID, SUBSCRIPTION_UUID, ORGANIZATION FROM APPLICATION_SUBSCRIPTION_MAPPING;"
	updateAppSub    = "UPDATE APPLICATION_SUBSCRIPTION_MAPPING SET APPLICATION_UUID = $2, SUBSCRIPTION_UUID = $3, ORGANIZATION = $4 WHERE UUID = $1;"
	deleteAppSub    = "DELETE FROM APPLICATION_SUBSCRIPTION_MAPPING WHERE UUID = $1;"
	deleteAllAppSub = "DELETE FROM APPLICATION_SUBSCRIPTION_MAPPING;"

	getAllApplicationKeyMappings    = "SELECT APPLICATION_UUID, SECURITY_SCHEME, APPLICATION_IDENTIFIER, KEY_TYPE, ENVIRONMENT, ORGANIZATION FROM APPLICATION_KEY_MAPPING;"
	insertApplicationKeyMapping     = "INSERT INTO APPLICATION_KEY_MAPPING (APPLICATION_UUID, SECURITY_SCHEME, APPLICATION_IDENTIFIER, KEY_TYPE, ENVIRONMENT, ORGANIZATION) VALUES ($1, $2, $3, $4, $5, $6);"
	updateApplicationKeyMapping     = "UPDATE APPLICATION_KEY_MAPPING SET APPLICATION_IDENTIFIER = $3, KEY_TYPE = $4, ORGANIZATION = $6 WHERE APPLICATION_UUID = $1, SECURITY_SCHEME = $2 AND ENVIRONMENT=$5;"
	deleteApplicationKeyMapping     = "DELETE FROM APPLICATION_KEY_MAPPING WHERE APPLICATION_UUID = $1, SECURITY_SCHEME = $2 AND ENVIRONMENT=$3;"
	deleteAllApplicationKeyMappings = "DELETE FROM APPLICATION_KEY_MAPPING"
)
