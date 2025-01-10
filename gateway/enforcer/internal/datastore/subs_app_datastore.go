/*
 *  Copyright (c) 2025, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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

package datastore

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"sync"

	subscription_model "github.com/wso2/apk/common-go-libs/pkg/server/model"
	"github.com/wso2/apk/gateway/enforcer/internal/config"
	"github.com/wso2/apk/gateway/enforcer/internal/util"
)

// SubscriptionApplicationDataStore is a data store that holds information about applications,
// application mappings, application key mappings, and subscriptions. It provides thread-safe
// access to these data structures using a read-write mutex.
//
// Fields:
// - applications: A map of application IDs to Application objects.
// - applicationMappings: A map of application IDs to ApplicationMapping objects.
// - applicationKeyMappings: A map of application IDs to ApplicationKeyMapping objects.
// - subscriptions: A map of subscription IDs to Subscription objects.
// - mu: A read-write mutex to ensure thread-safe access to the data store.
// - commonControllerRestBaseUrl: The base URL for the common controller REST API.
type SubscriptionApplicationDataStore struct {
	applications                map[string]*subscription_model.Application
	applicationMappings         map[string]*subscription_model.ApplicationMapping
	applicationKeyMappings      map[string]*subscription_model.ApplicationKeyMapping
	subscriptions               map[string]*subscription_model.Subscription
	mu                          sync.RWMutex
	commonControllerRestBaseURL string
}

// NewDataStore Initialize the datastore
func NewDataStore(cfg *config.Server) *SubscriptionApplicationDataStore {
	return &SubscriptionApplicationDataStore{
		applications:                make(map[string]*subscription_model.Application),
		applicationMappings:         make(map[string]*subscription_model.ApplicationMapping),
		applicationKeyMappings:      make(map[string]*subscription_model.ApplicationKeyMapping),
		subscriptions:               make(map[string]*subscription_model.Subscription),
		commonControllerRestBaseURL: "https://" + cfg.CommonControllerHost + ":" + cfg.CommonControllerRestPort,
	}
}

// AddApplication adds a new Application
func (ds *SubscriptionApplicationDataStore) AddApplication(app *subscription_model.Application) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.applications[app.UUID] = app
}

// AddApplicationMapping add a new ApplicationMapping
func (ds *SubscriptionApplicationDataStore) AddApplicationMapping(mapping *subscription_model.ApplicationMapping) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.applicationMappings[mapping.UUID] = mapping
}

// AddApplicationKeyMapping adds a new ApplicationKeyMapping
func (ds *SubscriptionApplicationDataStore) AddApplicationKeyMapping(keyMapping *subscription_model.ApplicationKeyMapping) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.applicationKeyMappings[keyMapping.ApplicationIdentifier] = keyMapping
}

// AddSubscription Add a new Subscription
func (ds *SubscriptionApplicationDataStore) AddSubscription(subscription *subscription_model.Subscription) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.subscriptions[subscription.UUID] = subscription
}

// DeleteApplication Delete an Application by UUID
func (ds *SubscriptionApplicationDataStore) DeleteApplication(id string) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	if _, exists := ds.applications[id]; exists {
		delete(ds.applications, id)
		return nil
	}
	return errors.New("application not found")
}

// DeleteApplicationMapping Delete an ApplicationMapping by UUID
func (ds *SubscriptionApplicationDataStore) DeleteApplicationMapping(id string) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	if _, exists := ds.applicationMappings[id]; exists {
		delete(ds.applicationMappings, id)
		return nil
	}
	return errors.New("applicationMapping not found")
}

// DeleteApplicationKeyMapping Delete an ApplicationKeyMapping by UUID
func (ds *SubscriptionApplicationDataStore) DeleteApplicationKeyMapping(id string) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	if _, exists := ds.applicationKeyMappings[id]; exists {
		delete(ds.applicationKeyMappings, id)
		return nil
	}
	return errors.New("applicationKeyMapping not found")
}

// DeleteSubscription Delete a Subscription by UUID
func (ds *SubscriptionApplicationDataStore) DeleteSubscription(id string) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	if _, exists := ds.subscriptions[id]; exists {
		delete(ds.subscriptions, id)
		return nil
	}
	return errors.New("subscription not found")
}

// UpdateApplication Update an Application
func (ds *SubscriptionApplicationDataStore) UpdateApplication(app *subscription_model.Application) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	if _, exists := ds.applications[app.UUID]; exists {
		ds.applications[app.UUID] = app
		return nil
	}
	return errors.New("application not found")
}

// UpdateApplicationMapping Update an ApplicationMapping
func (ds *SubscriptionApplicationDataStore) UpdateApplicationMapping(mapping *subscription_model.ApplicationMapping) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	if _, exists := ds.applicationMappings[mapping.UUID]; exists {
		ds.applicationMappings[mapping.UUID] = mapping
		return nil
	}
	return errors.New("applicationMapping not found")
}

// UpdateApplicationKeyMapping Update an ApplicationKeyMapping
func (ds *SubscriptionApplicationDataStore) UpdateApplicationKeyMapping(keyMapping *subscription_model.ApplicationKeyMapping) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	if _, exists := ds.applicationKeyMappings[keyMapping.ApplicationIdentifier]; exists {
		ds.applicationKeyMappings[keyMapping.ApplicationIdentifier] = keyMapping
		return nil
	}
	return errors.New("applicationKeyMapping not found")
}

// UpdateSubscription Update a Subscription
func (ds *SubscriptionApplicationDataStore) UpdateSubscription(subscription *subscription_model.Subscription) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	if _, exists := ds.subscriptions[subscription.UUID]; exists {
		ds.subscriptions[subscription.UUID] = subscription
		return nil
	}
	return errors.New("subscription not found")
}

// LoadStartupData loads all the necessary startup data into the SubscriptionApplicationDataStore.
// It retrieves all subscriptions, applications, application mappings, and application key mappings,
// and adds them to the data store. If any error occurs during the retrieval process, it returns the error.
func (ds *SubscriptionApplicationDataStore) LoadStartupData() error {
	subsList, err := ds.getAllSubscriptions()
	if err != nil {
		return err
	}
	for _, sub := range subsList.List {
		ds.AddSubscription(&sub)
	}
	appList, err := ds.getAllApplications()
	if err != nil {
		return err
	}
	for _, app := range appList.List {
		ds.AddApplication(&app)
	}
	appMappingList, err := ds.getAllApplicationMappings()
	if err != nil {
		return err
	}
	for _, appMapping := range appMappingList.List {
		ds.AddApplicationMapping(&appMapping)
	}
	appKeyMappingList, err := ds.getAllApplicationKeyMappings()
	if err != nil {
		return err
	}
	for _, appKeyMapping := range appKeyMappingList.List {
		ds.AddApplicationKeyMapping(&appKeyMapping)
	}
	return nil
}

// Get all applications
func (ds *SubscriptionApplicationDataStore) getAllApplications() (*subscription_model.ApplicationList, error) {
	url := fmt.Sprintf("%s/applications", ds.commonControllerRestBaseURL)
	resp, err := util.MakeGETRequest(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result subscription_model.ApplicationList
	body, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Get all subscriptions
func (ds *SubscriptionApplicationDataStore) getAllSubscriptions() (*subscription_model.SubscriptionList, error) {
	url := fmt.Sprintf("%s/subscriptions", ds.commonControllerRestBaseURL)
	resp, err := util.MakeGETRequest(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result subscription_model.SubscriptionList
	body, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Get all application mappings
func (ds *SubscriptionApplicationDataStore) getAllApplicationMappings() (*subscription_model.ApplicationMappingList, error) {
	url := fmt.Sprintf("%s/applicationmappings", ds.commonControllerRestBaseURL)
	resp, err := util.MakeGETRequest(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result subscription_model.ApplicationMappingList
	body, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Get all application key mappings
func (ds *SubscriptionApplicationDataStore) getAllApplicationKeyMappings() (*subscription_model.ApplicationKeyMappingList, error) {
	url := fmt.Sprintf("%s/applicationkeymappings", ds.commonControllerRestBaseURL)
	resp, err := util.MakeGETRequest(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result subscription_model.ApplicationKeyMappingList
	body, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
