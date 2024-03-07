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

import (
	"github.com/jackc/pgx/v5"
	"github.com/wso2/apk/common-controller/internal/loggers"
	"github.com/wso2/apk/common-controller/internal/server"
)

// deployApplicationAttributes deploys application attributes
func deployApplicationwithAttributes(tx pgx.Tx, application server.Application) error {
	PrepareQueries(tx, insertApplication, insertApplicationAttributes)
	err := AddApplication(tx, application.UUID, application.Name, application.Owner, application.OrganizationID)
	if err != nil {
		loggers.LoggerAPI.Error("Error while adding application ", err)
		return err
	}
	for attributeKey, attributeValue := range application.Attributes {
		err = AddApplicationAttributes(tx, application.UUID, attributeKey, attributeValue)
		if err != nil {
			loggers.LoggerAPI.Error("Error while adding application attributes ", err)
			return err
		}
	}
	return nil
}

func updateApplicationAttributes(tx pgx.Tx, application server.Application) error {
	PrepareQueries(tx, insertApplicationAttributes, deleteAllAppAttributes)
	err := DeleteApplicationAttributes(tx, application.UUID)
	if err != nil {
		loggers.LoggerAPI.Error("Error while deleting application attributes ", err)
		return err
	}
	for attributeKey, attributeValue := range application.Attributes {
		err = AddApplicationAttributes(tx, application.UUID, attributeKey, attributeValue)
		if err != nil {
			loggers.LoggerAPI.Error("Error while adding application attributes ", err)
			return err
		}
	}
	return nil
}

// GetAllApplications gets all applications from the database
func GetAllApplications(tx pgx.Tx) ([]server.Application, error) {
	rows, err := ExecDBQueryRows(tx, getAllApplicationAttributes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	appAttributes := make(map[string]map[string]string)
	for rows.Next() {
		var uuid string
		var name string
		var attrib string
		err := rows.Scan(&uuid, &name, &attrib)
		if err != nil {
			return nil, err
		}
		if _, ok := appAttributes[uuid]; !ok {
			appAttributes[uuid] = map[string]string{}
		}
		appAttributes[uuid][name] = attrib
	}

	rows, err = ExecDBQueryRows(tx, getAllApplications)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var applications []server.Application
	for rows.Next() {
		var app server.Application
		err := rows.Scan(&app.UUID, &app.Name, &app.Owner, &app.OrganizationID)
		if err != nil {
			return nil, err
		}
		app.Attributes = appAttributes[app.UUID]
		applications = append(applications, app)
	}
	return applications, nil
}

// GetAllSubscription gets all subscriptions from the database
func GetAllSubscription(tx pgx.Tx) ([]server.Subscription, error) {
	rows, err := ExecDBQueryRows(tx, getAllSubscriptions)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var subscriptions []server.Subscription
	for rows.Next() {
		sub := server.Subscription{
			SubscribedAPI: &server.SubscribedAPI{},
		}

		err := rows.Scan(&sub.UUID, &sub.SubscribedAPI.Name, &sub.SubscribedAPI.Version, &sub.SubStatus, &sub.Organization)
		if err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, sub)
	}
	return subscriptions, nil
}

// GetAllApplicationKeyMappings gets all application key mappings from the database
func GetAllApplicationKeyMappings(tx pgx.Tx) ([]server.ApplicationKeyMapping, error) {
	rows, err := ExecDBQueryRows(tx, getAllApplicationKeyMappings)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var appKeyMappings []server.ApplicationKeyMapping
	for rows.Next() {
		var appKeyMapping server.ApplicationKeyMapping
		err := rows.Scan(&appKeyMapping.ApplicationUUID, &appKeyMapping.SecurityScheme, &appKeyMapping.ApplicationIdentifier,
			&appKeyMapping.KeyType, &appKeyMapping.EnvID, &appKeyMapping.OrganizationID)
		if err != nil {
			return nil, err
		}
		appKeyMappings = append(appKeyMappings, appKeyMapping)
	}
	return appKeyMappings, nil
}

// GetAllAppSubs gets all application subscription mappings from the database
func GetAllAppSubs(tx pgx.Tx) ([]server.ApplicationMapping, error) {
	rows, err := ExecDBQueryRows(tx, getAllAppSubs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var appSubs []server.ApplicationMapping
	for rows.Next() {
		var appSub server.ApplicationMapping
		err := rows.Scan(&appSub.UUID, &appSub.ApplicationRef, &appSub.SubscriptionRef, &appSub.OrganizationID)
		if err != nil {
			return nil, err
		}
		appSubs = append(appSubs, appSub)
	}
	return appSubs, nil
}

// AddApplication adds an application to the database
func AddApplication(tx pgx.Tx, uuid, name, owner, org string) error {
	return ExecDBQuery(tx, insertApplication, uuid, name, owner, org)
}

// UpdateApplication updates an application in the database
func UpdateApplication(tx pgx.Tx, uuid, name, owner, org string) error {
	return ExecDBQuery(tx, updateApplication, uuid, name, owner, org)
}

// DeleteApplication deletes an application from the database
func DeleteApplication(tx pgx.Tx, uuid string) error {
	return ExecDBQuery(tx, deleteApplication, uuid)
}

// DeleteAllApplications deletes all applications from the database
func DeleteAllApplications(tx pgx.Tx) error {
	return ExecDBQuery(tx, deleteAllApplications)
}

// AddApplicationAttributes adds attributes to an application in the database
func AddApplicationAttributes(tx pgx.Tx, appUUID, name, appAttribute string) error {
	return ExecDBQuery(tx, insertApplicationAttributes, appUUID, name, appAttribute)
}

// DeleteApplicationAttributes deletes attributes of an application from the database
func DeleteApplicationAttributes(tx pgx.Tx, appUUID string) error {
	return ExecDBQuery(tx, deleteApplicationAttributes, appUUID)
}

// DeleteAllAppAttributes deletes all attributes of all applications from the database
func DeleteAllAppAttributes(tx pgx.Tx) error {
	return ExecDBQuery(tx, deleteAllAppAttributes)
}

// AddSubscription adds a subscription to the database
func AddSubscription(tx pgx.Tx, uuid, apiName, apiVersion, subStatus, organization string) error {
	return ExecDBQuery(tx, insertSubscription, uuid, apiName, apiVersion, subStatus, organization)
}

// UpdateSubscription updates a subscription in the database
func UpdateSubscription(tx pgx.Tx, uuid, apiName, apiVersion, subStatus, organization string) error {
	return ExecDBQuery(tx, updateSubscription, uuid, apiName, apiVersion, subStatus, organization)
}

// DeleteSubscription deletes a subscription from the database
func DeleteSubscription(tx pgx.Tx, uuid string) error {
	return ExecDBQuery(tx, deleteSubscription, uuid)
}

// DeleteAllSubscriptions deletes all subscriptions from the database
func DeleteAllSubscriptions(tx pgx.Tx) error {
	return ExecDBQuery(tx, deleteAllSubscriptions)
}

// AddApplicationKeyMapping adds a key mapping to the database
func AddApplicationKeyMapping(tx pgx.Tx, applicationUUID, securityScheme, applicationIdentifier, keyType, env,
	organization string) error {
	return ExecDBQuery(tx, insertApplicationKeyMapping, applicationUUID, securityScheme, applicationIdentifier, keyType,
		env, organization)
}

// DeleteApplicationKeyMapping deletes a key mapping from the database
func DeleteApplicationKeyMapping(tx pgx.Tx, applicationUUID, securityScheme, env string) error {
	return ExecDBQuery(tx, deleteApplicationKeyMapping, applicationUUID, securityScheme, env)
}

// UpdateApplicationKeyMapping updates a key mapping in the database
func UpdateApplicationKeyMapping(tx pgx.Tx, applicationUUID, securityScheme, applicationIdentifier, keyType, env,
	organization string) error {
	return ExecDBQuery(tx, updateApplicationKeyMapping, applicationUUID, securityScheme, applicationIdentifier, keyType,
		env, organization)
}

// DeleteAllApplicationKeyMappings deletes all key mappings from the database
func DeleteAllApplicationKeyMappings(tx pgx.Tx) error {
	return ExecDBQuery(tx, deleteAllApplicationKeyMappings)
}

// AddAppSub adds an application subscription mapping to the database
func AddAppSub(tx pgx.Tx, uuid, appUUID, subUUID, organization string) error {
	return ExecDBQuery(tx, insertAppSub, uuid, appUUID, subUUID, organization)
}

// UpdateAppSub updates an application subscription mapping in the database
func UpdateAppSub(tx pgx.Tx, uuid, appUUID, subUUID, organization string) error {
	return ExecDBQuery(tx, updateAppSub, uuid, appUUID, subUUID, organization)
}

// DeleteAppSub deletes an application subscription mapping from the database
func DeleteAppSub(tx pgx.Tx, uuid string) error {
	return ExecDBQuery(tx, deleteAppSub, uuid)
}

// DeleteAllAppSub deletes all application subscription mappings from the database
func DeleteAllAppSub(tx pgx.Tx) error {
	return ExecDBQuery(tx, deleteAllAppSub)
}
