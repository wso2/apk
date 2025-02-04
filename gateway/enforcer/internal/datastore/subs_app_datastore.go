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
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
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
	applications                map[string]map[string]*subscription_model.Application                   // OrganozationID -> ApplicationUUID -> Application
	applicationMappings         map[string]map[string]map[string]*subscription_model.ApplicationMapping // OrganizationID -> ApplicationRef -> ApplicationMappingUUID -> ApplicationMapping
	applicationKeyMappings      map[string]map[string]*subscription_model.ApplicationKeyMapping         // OrganizationID -> ApplicationKeyMappingCacheKey -> ApplicationKeyMapping
	subscriptions               map[string]map[string]*subscription_model.Subscription                  // OrganizationID -> SubscriptionUUID -> Subscription
	mu                          sync.RWMutex
	commonControllerRestBaseURL string
}

// NewSubAppDataStore Initialize the datastore
func NewSubAppDataStore(cfg *config.Server) *SubscriptionApplicationDataStore {
	return &SubscriptionApplicationDataStore{
		applications:                make(map[string]map[string]*subscription_model.Application),
		applicationMappings:         make(map[string]map[string]map[string]*subscription_model.ApplicationMapping),
		applicationKeyMappings:      make(map[string]map[string]*subscription_model.ApplicationKeyMapping),
		subscriptions:               make(map[string]map[string]*subscription_model.Subscription),
		commonControllerRestBaseURL: "https://" + cfg.CommonControllerHost + ":" + cfg.CommonControllerRestPort,
	}
}

// AddApplication adds a new Application
func (ds *SubscriptionApplicationDataStore) AddApplication(app *subscription_model.Application) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	if _, exists := ds.applications[app.OrganizationID]; !exists {
		ds.applications[app.OrganizationID] = make(map[string]*subscription_model.Application)
	}
	ds.applications[app.OrganizationID][app.UUID] = app
}

// AddApplicationMapping add a new ApplicationMapping
func (ds *SubscriptionApplicationDataStore) AddApplicationMapping(mapping *subscription_model.ApplicationMapping) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	if _, exists := ds.applicationMappings[mapping.OrganizationID]; !exists {
		ds.applicationMappings[mapping.OrganizationID] = make(map[string]map[string]*subscription_model.ApplicationMapping)
		ds.applicationMappings[mapping.OrganizationID][mapping.ApplicationRef] = make(map[string]*subscription_model.ApplicationMapping)
	} else if _, exists := ds.applicationMappings[mapping.OrganizationID][mapping.ApplicationRef]; !exists {
		ds.applicationMappings[mapping.OrganizationID][mapping.ApplicationRef] = make(map[string]*subscription_model.ApplicationMapping)
	}
	ds.applicationMappings[mapping.OrganizationID][mapping.ApplicationRef][mapping.UUID] = mapping
}

// AddApplicationKeyMapping adds a new ApplicationKeyMapping
func (ds *SubscriptionApplicationDataStore) AddApplicationKeyMapping(keyMapping *subscription_model.ApplicationKeyMapping) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	if _, exists := ds.applicationKeyMappings[keyMapping.OrganizationID]; !exists {
		ds.applicationKeyMappings[keyMapping.OrganizationID] = make(map[string]*subscription_model.ApplicationKeyMapping)
	}

	ds.applicationKeyMappings[keyMapping.OrganizationID][util.PrepareApplicationKeyMappingCacheKey(keyMapping.ApplicationIdentifier, keyMapping.KeyType, keyMapping.SecurityScheme, keyMapping.EnvID)] = keyMapping
}

// AddSubscription Add a new Subscription
func (ds *SubscriptionApplicationDataStore) AddSubscription(subscription *subscription_model.Subscription) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	if _, exists := ds.subscriptions[subscription.Organization]; !exists {
		ds.subscriptions[subscription.Organization] = make(map[string]*subscription_model.Subscription)
	}
	ds.subscriptions[subscription.Organization][subscription.UUID] = subscription
}

// DeleteApplication Delete an Application by UUID
func (ds *SubscriptionApplicationDataStore) DeleteApplication(application *subscription_model.Application) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	if _, exists := ds.applications[application.OrganizationID]; exists {
		delete(ds.applications[application.OrganizationID], application.UUID)
		return nil
	}
	return errors.New("application not found")
}

// DeleteApplicationMapping Delete an ApplicationMapping by UUID
func (ds *SubscriptionApplicationDataStore) DeleteApplicationMapping(appMap *subscription_model.ApplicationMapping) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	if _, exists := ds.applicationMappings[appMap.OrganizationID]; exists {
		if _, exists := ds.applicationMappings[appMap.OrganizationID][appMap.ApplicationRef]; exists {
			delete(ds.applicationMappings[appMap.OrganizationID][appMap.ApplicationRef], appMap.UUID)
			return nil
		}
		return errors.New("applicationMapping not found")
	}
	return errors.New("applicationMapping not found")
}

// DeleteApplicationKeyMapping Delete an ApplicationKeyMapping by UUID
func (ds *SubscriptionApplicationDataStore) DeleteApplicationKeyMapping(appKeyMap *subscription_model.ApplicationKeyMapping) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	if _, exists := ds.applicationKeyMappings[appKeyMap.OrganizationID]; exists {
		delete(ds.applicationKeyMappings[appKeyMap.OrganizationID], util.PrepareApplicationKeyMappingCacheKey(appKeyMap.ApplicationIdentifier, appKeyMap.KeyType, appKeyMap.SecurityScheme, appKeyMap.EnvID))
		return nil
	}
	return errors.New("applicationKeyMapping not found")
}

// DeleteSubscription Delete a Subscription by UUID
func (ds *SubscriptionApplicationDataStore) DeleteSubscription(subscription *subscription_model.Subscription) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	if _, exists := ds.subscriptions[subscription.Organization]; exists {
		delete(ds.subscriptions[subscription.Organization], subscription.UUID)
		return nil
	}
	return errors.New("subscription not found")
}

// GetApplicationMappings Get an ApplicationMapping by UUID
func (ds *SubscriptionApplicationDataStore) GetApplicationMappings(org string, appID string) map[string]*subscription_model.ApplicationMapping {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	if _, exists := ds.applicationMappings[org]; exists {
		if _, exists := ds.applicationMappings[org][appID]; exists {
			return ds.applicationMappings[org][appID]
		}
	}
	return nil
}

// GetSubscriptions Get an Subscription by UUID
func (ds *SubscriptionApplicationDataStore) GetSubscriptions(org string, subscriptionID string) map[string]*subscription_model.Subscription {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	if _, exists := ds.subscriptions[org]; exists {
		if _, exists := ds.subscriptions[org][subscriptionID]; exists {
			return ds.subscriptions[org]
		}
	}
	return nil
}

// GetSubscription Get an Subscription by UUID
func (ds *SubscriptionApplicationDataStore) GetSubscription(org string, subscriptionID string) *subscription_model.Subscription {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	if _, exists := ds.subscriptions[org]; exists {
		if _, exists := ds.subscriptions[org][subscriptionID]; exists {
			return ds.subscriptions[org][subscriptionID]
		}
	}
	return nil
}

// GetApplicationKeyMapping Get an ApplicationKeyMapping by UUID
func (ds *SubscriptionApplicationDataStore) GetApplicationKeyMapping(org string, appKeyMapKey string) *subscription_model.ApplicationKeyMapping {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	if _, exists := ds.applicationKeyMappings[org]; exists {
		if _, exists := ds.applicationKeyMappings[org][appKeyMapKey]; exists {
			return ds.applicationKeyMappings[org][appKeyMapKey]
		}
	}
	return nil
}

// GetApplication Get an Application by UUID
func (ds *SubscriptionApplicationDataStore) GetApplication(org string, appID string) *subscription_model.Application {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	if _, exists := ds.applications[org]; exists {
		if _, exists := ds.applications[org][appID]; exists {
			return ds.applications[org][appID]
		}
	}
	return nil
}

// // UpdateApplication Update an Application
// func (ds *SubscriptionApplicationDataStore) UpdateApplication(app *subscription_model.Application) error {
// 	ds.mu.Lock()
// 	defer ds.mu.Unlock()
// 	if _, exists := ds.applications[app.UUID]; exists {
// 		ds.applications[app.UUID] = app
// 		return nil
// 	}
// 	return errors.New("application not found")
// }

// // UpdateApplicationMapping Update an ApplicationMapping
// func (ds *SubscriptionApplicationDataStore) UpdateApplicationMapping(mapping *subscription_model.ApplicationMapping) error {
// 	ds.mu.Lock()
// 	defer ds.mu.Unlock()
// 	if _, exists := ds.applicationMappings[mapping.UUID]; exists {
// 		ds.applicationMappings[mapping.UUID] = mapping
// 		return nil
// 	}
// 	return errors.New("applicationMapping not found")
// }

// // UpdateApplicationKeyMapping Update an ApplicationKeyMapping
// func (ds *SubscriptionApplicationDataStore) UpdateApplicationKeyMapping(keyMapping *subscription_model.ApplicationKeyMapping) error {
// 	ds.mu.Lock()
// 	defer ds.mu.Unlock()
// 	if _, exists := ds.applicationKeyMappings[keyMapping.ApplicationIdentifier]; exists {
// 		ds.applicationKeyMappings[keyMapping.ApplicationIdentifier] = keyMapping
// 		return nil
// 	}
// 	return errors.New("applicationKeyMapping not found")
// }

// // UpdateSubscription Update a Subscription
// func (ds *SubscriptionApplicationDataStore) UpdateSubscription(subscription *subscription_model.Subscription) error {
// 	ds.mu.Lock()
// 	defer ds.mu.Unlock()
// 	if _, exists := ds.subscriptions[subscription.UUID]; exists {
// 		ds.subscriptions[subscription.UUID] = subscription
// 		return nil
// 	}
// 	return errors.New("subscription not found")
// }

// LoadStartupData loads all the necessary startup data into the SubscriptionApplicationDataStore.
// It retrieves all subscriptions, applications, application mappings, and application key mappings,
// and adds them to the data store. If any error occurs during the retrieval process, it returns the error.
func (ds *SubscriptionApplicationDataStore) LoadStartupData() error {
	subsList, err := ds.loadAllSubscriptions()
	if err != nil {
		return err
	}
	for _, sub := range subsList.List {
		ds.AddSubscription(&sub)
	}
	appList, err := ds.loadAllApplications()
	if err != nil {
		return err
	}
	for _, app := range appList.List {
		ds.AddApplication(&app)
	}
	appMappingList, err := ds.loadAllApplicationMappings()
	if err != nil {
		return err
	}
	for _, appMapping := range appMappingList.List {
		ds.AddApplicationMapping(&appMapping)
	}
	appKeyMappingList, err := ds.loadAllApplicationKeyMappings()
	if err != nil {
		return err
	}
	for _, appKeyMapping := range appKeyMappingList.List {
		ds.AddApplicationKeyMapping(&appKeyMapping)
	}
	return nil
}

// Get all applications
func (ds *SubscriptionApplicationDataStore) loadAllApplications() (*subscription_model.ApplicationList, error) {
	url := fmt.Sprintf("%s/applications", ds.commonControllerRestBaseURL)
	// Get the TLS configuration
	tlsConfig, err := GetTLSConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get TLS config: %w", err)
	}
	resp, err := util.MakeGETRequest(url, tlsConfig)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result subscription_model.ApplicationList
	body, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	log.Println("Applications: ", result)
	return &result, nil
}

// Get all subscriptions
func (ds *SubscriptionApplicationDataStore) loadAllSubscriptions() (*subscription_model.SubscriptionList, error) {
	url := fmt.Sprintf("%s/subscriptions", ds.commonControllerRestBaseURL)
	// Get the TLS configuration
	tlsConfig, err := GetTLSConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get TLS config: %w", err)
	}
	resp, err := util.MakeGETRequest(url, tlsConfig)
	log.Println("Response: ", resp)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result subscription_model.SubscriptionList
	body, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	log.Println("Subscription: ", result)
	return &result, nil
}

// Get all application mappings
func (ds *SubscriptionApplicationDataStore) loadAllApplicationMappings() (*subscription_model.ApplicationMappingList, error) {
	url := fmt.Sprintf("%s/applicationmappings", ds.commonControllerRestBaseURL)
	tlsConfig, err := GetTLSConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get TLS config: %w", err)
	}
	resp, err := util.MakeGETRequest(url, tlsConfig)
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
func (ds *SubscriptionApplicationDataStore) loadAllApplicationKeyMappings() (*subscription_model.ApplicationKeyMappingList, error) {
	url := fmt.Sprintf("%s/applicationkeymappings", ds.commonControllerRestBaseURL)
	tlsConfig, err := GetTLSConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get TLS config: %w", err)
	}
	resp, err := util.MakeGETRequest(url, tlsConfig)
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

// GetTLSConfig loads and returns a TLS configuration
func GetTLSConfig() (*tls.Config, error) {
	cfg := config.GetConfig()

	// Load the client certificate and private key
	clientCert, err := util.LoadCertificates(cfg.EnforcerPublicKeyPath, cfg.EnforcerPrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load client certificate and private key: %w", err)
	}

	// Load the trusted CA certificates
	certPool, err := util.LoadCACertificates(cfg.TrustedAdapterCertsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load trusted CA certificates: %w", err)
	}

	// Create and return the TLS configuration
	tlsConfig := util.CreateTLSConfig(clientCert, certPool)
	return tlsConfig, nil
}
