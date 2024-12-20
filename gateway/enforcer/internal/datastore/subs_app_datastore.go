package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"sync"

	subscription_model "github.com/wso2/apk/common-go-libs/pkg/discovery/api/wso2/discovery/subscription"
	"github.com/wso2/apk/gateway/enforcer/internal/util"
	"github.com/wso2/apk/gateway/enforcer/internal/config"
)

// Define the datastore
type DataStore struct {
	applications                map[string]*subscription_model.Application
	applicationMappings         map[string]*subscription_model.ApplicationMapping
	applicationKeyMappings      map[string]*subscription_model.ApplicationKeyMapping
	subscriptions               map[string]*subscription_model.Subscription
	mu                          sync.RWMutex
	commonControllerRestBaseUrl string
}

// Initialize the datastore
func NewDataStore(cfg *config.Server) *DataStore {
	return &DataStore{
		applications:           make(map[string]*subscription_model.Application),
		applicationMappings:    make(map[string]*subscription_model.ApplicationMapping),
		applicationKeyMappings: make(map[string]*subscription_model.ApplicationKeyMapping),
		subscriptions:          make(map[string]*subscription_model.Subscription),
		commonControllerRestBaseUrl: "https://" + cfg.CommonControllerHost + ":" + cfg.CommonControllerRestPort ,
	}
}

// Add a new Application
func (ds *DataStore) AddApplication(app *subscription_model.Application) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.applications[app.Uuid] = app
}

// Add a new ApplicationMapping
func (ds *DataStore) AddApplicationMapping(mapping *subscription_model.ApplicationMapping) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.applicationMappings[mapping.Uuid] = mapping
}

// Add a new ApplicationKeyMapping
func (ds *DataStore) AddApplicationKeyMapping(keyMapping *subscription_model.ApplicationKeyMapping) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.applicationKeyMappings[keyMapping.ApplicationIdentifier] = keyMapping
}

// Add a new Subscription
func (ds *DataStore) AddSubscription(subscription *subscription_model.Subscription) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.subscriptions[subscription.Uuid] = subscription
}

// Delete an Application by UUID
func (ds *DataStore) DeleteApplication(id string) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	if _, exists := ds.applications[id]; exists {
		delete(ds.applications, id)
		return nil
	}
	return errors.New("application not found")
}

// Delete an ApplicationMapping by UUID
func (ds *DataStore) DeleteApplicationMapping(id string) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	if _, exists := ds.applicationMappings[id]; exists {
		delete(ds.applicationMappings, id)
		return nil
	}
	return errors.New("applicationMapping not found")
}

// Delete an ApplicationKeyMapping by UUID
func (ds *DataStore) DeleteApplicationKeyMapping(id string) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	if _, exists := ds.applicationKeyMappings[id]; exists {
		delete(ds.applicationKeyMappings, id)
		return nil
	}
	return errors.New("applicationKeyMapping not found")
}

// Delete a Subscription by UUID
func (ds *DataStore) DeleteSubscription(id string) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	if _, exists := ds.subscriptions[id]; exists {
		delete(ds.subscriptions, id)
		return nil
	}
	return errors.New("subscription not found")
}

// Update an Application
func (ds *DataStore) UpdateApplication(app *subscription_model.Application) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	if _, exists := ds.applications[app.Uuid]; exists {
		ds.applications[app.Uuid] = app
		return nil
	}
	return errors.New("application not found")
}

// Update an ApplicationMapping
func (ds *DataStore) UpdateApplicationMapping(mapping *subscription_model.ApplicationMapping) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	if _, exists := ds.applicationMappings[mapping.Uuid]; exists {
		ds.applicationMappings[mapping.Uuid] = mapping
		return nil
	}
	return errors.New("applicationMapping not found")
}

// Update an ApplicationKeyMapping
func (ds *DataStore) UpdateApplicationKeyMapping(keyMapping *subscription_model.ApplicationKeyMapping) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	if _, exists := ds.applicationKeyMappings[keyMapping.ApplicationIdentifier]; exists {
		ds.applicationKeyMappings[keyMapping.ApplicationIdentifier] = keyMapping
		return nil
	}
	return errors.New("applicationKeyMapping not found")
}

// Update a Subscription
func (ds *DataStore) UpdateSubscription(subscription *subscription_model.Subscription) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	if _, exists := ds.subscriptions[subscription.Uuid]; exists {
		ds.subscriptions[subscription.Uuid] = subscription
		return nil
	}
	return errors.New("subscription not found")
}

// Get all applications
func (ds *DataStore) getAllApplications() (*subscription_model.Application, error) {
	url := fmt.Sprintf("%s/applications", ds.commonControllerRestBaseUrl)
	resp, err := util.MakeGETRequest(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result ApplicationListDto
	body, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Get all subscriptions
func (ds *DataStore) getAllSubscriptions() (*SubscriptionListDto, error) {
	url := fmt.Sprintf("%s/subscriptions", ds.commonControllerRestBaseUrl)
	resp, err := util.MakeGETRequest(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result SubscriptionListDto
	body, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Get all application mappings
func (ds *DataStore) getAllApplicationMappings() (*ApplicationMappingDtoList, error) {
	url := fmt.Sprintf("%s/applicationmappings", ds.commonControllerRestBaseUrl)
	resp, err := util.MakeGETRequest(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result ApplicationMappingDtoList
	body, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Get all application key mappings
func (ds *DataStore) getAllApplicationKeyMappings() (*ApplicationKeyMappingDtoList, error) {
	url := fmt.Sprintf("%s/applicationkeymappings", ds.commonControllerRestBaseUrl)
	resp, err := util.MakeGETRequest(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result ApplicationKeyMappingDtoList
	body, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
