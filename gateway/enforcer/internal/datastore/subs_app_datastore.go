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

// Define the datastore
type SubscriptionApplicationDataStore struct {
	applications                map[string]*subscription_model.Application
	applicationMappings         map[string]*subscription_model.ApplicationMapping
	applicationKeyMappings      map[string]*subscription_model.ApplicationKeyMapping
	subscriptions               map[string]*subscription_model.Subscription
	mu                          sync.RWMutex
	commonControllerRestBaseUrl string
}

// Initialize the datastore
func NewDataStore(cfg *config.Server) *SubscriptionApplicationDataStore {
	return &SubscriptionApplicationDataStore{
		applications:                make(map[string]*subscription_model.Application),
		applicationMappings:         make(map[string]*subscription_model.ApplicationMapping),
		applicationKeyMappings:      make(map[string]*subscription_model.ApplicationKeyMapping),
		subscriptions:               make(map[string]*subscription_model.Subscription),
		commonControllerRestBaseUrl: "https://" + cfg.CommonControllerHost + ":" + cfg.CommonControllerRestPort,
	}
}

// Add a new Application
func (ds *SubscriptionApplicationDataStore) AddApplication(app *subscription_model.Application) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.applications[app.UUID] = app
}

// Add a new ApplicationMapping
func (ds *SubscriptionApplicationDataStore) AddApplicationMapping(mapping *subscription_model.ApplicationMapping) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.applicationMappings[mapping.UUID] = mapping
}

// Add a new ApplicationKeyMapping
func (ds *SubscriptionApplicationDataStore) AddApplicationKeyMapping(keyMapping *subscription_model.ApplicationKeyMapping) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.applicationKeyMappings[keyMapping.ApplicationIdentifier] = keyMapping
}

// Add a new Subscription
func (ds *SubscriptionApplicationDataStore) AddSubscription(subscription *subscription_model.Subscription) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.subscriptions[subscription.UUID] = subscription
}

// Delete an Application by UUID
func (ds *SubscriptionApplicationDataStore) DeleteApplication(id string) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	if _, exists := ds.applications[id]; exists {
		delete(ds.applications, id)
		return nil
	}
	return errors.New("application not found")
}

// Delete an ApplicationMapping by UUID
func (ds *SubscriptionApplicationDataStore) DeleteApplicationMapping(id string) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	if _, exists := ds.applicationMappings[id]; exists {
		delete(ds.applicationMappings, id)
		return nil
	}
	return errors.New("applicationMapping not found")
}

// Delete an ApplicationKeyMapping by UUID
func (ds *SubscriptionApplicationDataStore) DeleteApplicationKeyMapping(id string) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	if _, exists := ds.applicationKeyMappings[id]; exists {
		delete(ds.applicationKeyMappings, id)
		return nil
	}
	return errors.New("applicationKeyMapping not found")
}

// Delete a Subscription by UUID
func (ds *SubscriptionApplicationDataStore) DeleteSubscription(id string) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	if _, exists := ds.subscriptions[id]; exists {
		delete(ds.subscriptions, id)
		return nil
	}
	return errors.New("subscription not found")
}

// Update an Application
func (ds *SubscriptionApplicationDataStore) UpdateApplication(app *subscription_model.Application) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	if _, exists := ds.applications[app.UUID]; exists {
		ds.applications[app.UUID] = app
		return nil
	}
	return errors.New("application not found")
}

// Update an ApplicationMapping
func (ds *SubscriptionApplicationDataStore) UpdateApplicationMapping(mapping *subscription_model.ApplicationMapping) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	if _, exists := ds.applicationMappings[mapping.UUID]; exists {
		ds.applicationMappings[mapping.UUID] = mapping
		return nil
	}
	return errors.New("applicationMapping not found")
}

// Update an ApplicationKeyMapping
func (ds *SubscriptionApplicationDataStore) UpdateApplicationKeyMapping(keyMapping *subscription_model.ApplicationKeyMapping) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	if _, exists := ds.applicationKeyMappings[keyMapping.ApplicationIdentifier]; exists {
		ds.applicationKeyMappings[keyMapping.ApplicationIdentifier] = keyMapping
		return nil
	}
	return errors.New("applicationKeyMapping not found")
}

// Update a Subscription
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
	} else {
		for _, sub := range subsList.List {
			ds.AddSubscription(&sub)
		}
	}
	appList, err := ds.getAllApplications()
	if err != nil {
		return err
	} else {
		for _, app := range appList.List {
			ds.AddApplication(&app)
		}
	}
	appMappingList, err := ds.getAllApplicationMappings()
	if err != nil {
		return err
	} else {
		for _, appMapping := range appMappingList.List {
			ds.AddApplicationMapping(&appMapping)
		}
	}
	appKeyMappingList, err := ds.getAllApplicationKeyMappings()
	if err != nil {
		return err
	} else {
		for _, appKeyMapping := range appKeyMappingList.List {
			ds.AddApplicationKeyMapping(&appKeyMapping)
		}
	}
	return nil
}

// Get all applications
func (ds *SubscriptionApplicationDataStore) getAllApplications() (*subscription_model.ApplicationList, error) {
	url := fmt.Sprintf("%s/applications", ds.commonControllerRestBaseUrl)
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
	url := fmt.Sprintf("%s/subscriptions", ds.commonControllerRestBaseUrl)
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
	url := fmt.Sprintf("%s/applicationmappings", ds.commonControllerRestBaseUrl)
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
	url := fmt.Sprintf("%s/applicationkeymappings", ds.commonControllerRestBaseUrl)
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
