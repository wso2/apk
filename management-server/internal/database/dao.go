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

import (
	"encoding/json"
	"fmt"
	"time"

	apkmgt "github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/apkmgt"
	apiProtos "github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/service/apkmgt"
	"github.com/wso2/apk/adapter/pkg/logging"
	"github.com/wso2/apk/management-server/internal/logger"
)

// DbCache is a pointer to an ApplicationLocalCache
var DbCache *ApplicationLocalCache

func init() {
	DbCache = NewApplicationLocalCache(cleanupInterval)
}

type artifact struct {
	APIName      string `json:"apiName"`
	ID           string `json:"id"`
	Context      string `json:"context"`
	Version      string `json:"version"`
	ProviderName string `json:"providerName"`
	Status       string `json:"status"`
}

// GetApplicationByUUID retrives an application using uuid and returns it
func GetApplicationByUUID(uuid string) (*apkmgt.Application, error) {
	rows, _ := ExecDBQuery(queryGetApplicationByUUID, uuid)
	rows.Next()
	values, err := rows.Values()
	if err != nil {
		return nil, err
	}
	subs, _ := getSubscriptionsForApplication(uuid)
	keys, _ := getConsumerKeysForApplication(uuid)
	application := &apkmgt.Application{
		Uuid:          values[0].(string),
		Name:          values[1].(string),
		Owner:         "",  //ToDo : Check how to get Owner from db
		Attributes:    nil, //ToDo : check the values for Attributes
		Subscriber:    "",
		Organization:  values[3].(string),
		Subscriptions: subs,
		ConsumerKeys:  keys,
	}
	DbCache.Update(application, time.Now().Unix()+ttl.Microseconds())
	return application, nil
}

// GetCachedApplicationByUUID returns the Application details from the cache.
// If the application is not available in the cache, it will fetch the application from DB.
func GetCachedApplicationByUUID(uuid string) (*apkmgt.Application, error) {
	if app, ok := DbCache.Read(uuid); ok == nil {
		return app, nil
	}
	return GetApplicationByUUID(uuid)
}

// getSubscriptionsForApplication returns all subscriptions from DB, for a given application.
func getSubscriptionsForApplication(appUUID string) ([]*apkmgt.Subscription, error) {
	rows, _ := ExecDBQuery(queryGetAllSubscriptionsForApplication, appUUID)
	var subs []*apkmgt.Subscription
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return nil, err
		}
		subs = append(subs, &apkmgt.Subscription{
			Uuid:               values[0].(string),
			ApiUuid:            values[1].(string),
			PolicyId:           "",
			SubscriptionStatus: values[3].(string),
			Organization:       values[4].(string),
			CreatedBy:          values[5].(string),
		})
	}
	return subs, nil
}

// getConsumerKeysForApplication returns all Consumer Keys from DB, for a given application.
func getConsumerKeysForApplication(appUUID string) ([]*apkmgt.ConsumerKey, error) {
	rows, _ := ExecDBQuery(queryConsumerKeysForApplication, appUUID)
	var keys []*apkmgt.ConsumerKey
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return nil, err
		}
		keys = append(keys, &apkmgt.ConsumerKey{
			Key:        values[0].(string),
			KeyManager: values[1].(string),
		})
	}
	return keys, nil
}

// GetSubscriptionByUUID returns the Application details from the DB for a given subscription UUID.
func GetSubscriptionByUUID(subUUID string) (*apkmgt.Subscription, error) {
	rows, _ := ExecDBQuery(querySubscriptionByUUID, subUUID)
	rows.Next()
	values, err := rows.Values()
	if err != nil {
		return nil, err
	}
	return &apkmgt.Subscription{
		Uuid:               values[0].(string),
		ApiUuid:            values[1].(string),
		PolicyId:           "",
		SubscriptionStatus: values[2].(string),
		Organization:       values[3].(string),
		CreatedBy:          values[4].(string),
	}, nil
}

// CreateAPI creates an API in the DB
func CreateAPI(api *apiProtos.API) error {
	_, err := ExecDBQuery(queryCreateAPI, &api.Uuid, &api.Name, &api.Provider,
		&api.Version, &api.Context, &api.OrganizationId, &api.CreatedBy, time.Now(), &api.Type, marshalArtifact(api), "PUBLISHED")

	if err != nil {
		logger.LoggerDatabase.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprintf("Error creating API %q, Error: %v", api.Uuid, err.Error()),
			Severity:  logging.CRITICAL,
			ErrorCode: 1201,
		})
		return err
	}
	return nil
}

// UpdateAPI updates the given API in the DB
func UpdateAPI(api *apiProtos.API) error {
	_, err := ExecDBQuery(queryUpdateAPI, &api.Uuid, &api.Name, &api.Provider,
		&api.Version, &api.Context, &api.OrganizationId, &api.UpdatedBy, time.Now(), &api.Type, marshalArtifact(api), "PUBLISHED")
	if err != nil {
		logger.LoggerDatabase.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprintf("Error updating API %q, Error: %v", api.Uuid, err.Error()),
			Severity:  logging.CRITICAL,
			ErrorCode: 1202,
		})
		return err
	}
	return nil
}

// DeleteAPI deletes the given API in the DB
func DeleteAPI(api *apiProtos.API) error {
	_, err := ExecDBQuery(queryDeleteAPI, api.Uuid)
	if err != nil {
		logger.LoggerDatabase.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprintf("Error deleting API %q, Error: %v", api.Uuid, err.Error()),
			Severity:  logging.CRITICAL,
			ErrorCode: 1203,
		})
		return err
	}
	return nil
}

func marshalArtifact(api *apiProtos.API) string {
	artifact := &artifact{APIName: api.Name,
		ID:           api.Uuid,
		Context:      api.Context,
		Version:      api.Version,
		ProviderName: api.Provider,
		Status:       "PUBLISHED",
	}
	jsonString, err := json.Marshal(artifact)
	if err != nil {
		return "{}"
	}
	return string(jsonString)
}
