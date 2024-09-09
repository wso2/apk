/*
 *  Copyright (c) 2023 WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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

package cache

import (
	"sync"

	logger "github.com/sirupsen/logrus"
	cpv1alpha2 "github.com/wso2/apk/common-go-libs/apis/cp/v1alpha2"
	cpv1alpha3 "github.com/wso2/apk/common-go-libs/apis/cp/v1alpha3"
	"k8s.io/apimachinery/pkg/types"
)

// SubscriptionDataStore is a cache subscription data.
type SubscriptionDataStore struct {
	applicationStore        map[types.NamespacedName]*cpv1alpha2.ApplicationSpec
	subscriptionStore       map[types.NamespacedName]*cpv1alpha3.SubscriptionSpec
	applicationMappingStore map[types.NamespacedName]*cpv1alpha2.ApplicationMappingSpec
	mu                      sync.Mutex
}

// CreateNewSubscriptionDataStore creates a new SubscriptionDataStore.
func CreateNewSubscriptionDataStore() *SubscriptionDataStore {
	return &SubscriptionDataStore{
		applicationStore:        map[types.NamespacedName]*cpv1alpha2.ApplicationSpec{},
		subscriptionStore:       map[types.NamespacedName]*cpv1alpha3.SubscriptionSpec{},
		applicationMappingStore: map[types.NamespacedName]*cpv1alpha2.ApplicationMappingSpec{},
	}
}

// AddorUpdateApplicationToStore adds a new application to the DataStore.
func (ods *SubscriptionDataStore) AddorUpdateApplicationToStore(name types.NamespacedName, application cpv1alpha2.ApplicationSpec) {
	ods.mu.Lock()
	defer ods.mu.Unlock()
	logger.Debug("Adding/Updating application to cache")
	ods.applicationStore[name] = &application
}

// AddorUpdateSubscriptionToStore adds a new subscription to the DataStore.
func (ods *SubscriptionDataStore) AddorUpdateSubscriptionToStore(name types.NamespacedName, subscription cpv1alpha3.SubscriptionSpec) {
	ods.mu.Lock()
	defer ods.mu.Unlock()
	logger.Debug("Adding/Updating subscription to cache")
	ods.subscriptionStore[name] = &subscription
}

// AddorUpdateApplicationMappingToStore adds a new application mapping to the DataStore.
func (ods *SubscriptionDataStore) AddorUpdateApplicationMappingToStore(name types.NamespacedName, applicationMapping cpv1alpha2.ApplicationMappingSpec) {
	ods.mu.Lock()
	defer ods.mu.Unlock()
	logger.Debug("Adding/Updating application mapping to cache")
	ods.applicationMappingStore[name] = &applicationMapping
}

// GetApplicationFromStore get cached application
func (ods *SubscriptionDataStore) GetApplicationFromStore(name types.NamespacedName) (cpv1alpha2.ApplicationSpec, bool) {
	var application cpv1alpha2.ApplicationSpec
	if cachedApplication, found := ods.applicationStore[name]; found {
		logger.Debug("Found cached application")
		return *cachedApplication, true
	}
	return application, false
}

// GetSubscriptionFromStore get cached subscription
func (ods *SubscriptionDataStore) GetSubscriptionFromStore(name types.NamespacedName) (cpv1alpha3.SubscriptionSpec, bool) {
	var subscription cpv1alpha3.SubscriptionSpec
	if cachedSubscription, found := ods.subscriptionStore[name]; found {
		logger.Debug("Found cached subscription")
		return *cachedSubscription, true
	}
	return subscription, false
}

// GetApplicationMappingFromStore get cached application mapping
func (ods *SubscriptionDataStore) GetApplicationMappingFromStore(name types.NamespacedName) (cpv1alpha2.ApplicationMappingSpec, bool) {
	var applicationMapping cpv1alpha2.ApplicationMappingSpec
	if cachedApplicationMapping, found := ods.applicationMappingStore[name]; found {
		logger.Debug("Found cached application mapping")
		return *cachedApplicationMapping, true
	}
	return applicationMapping, false
}

// DeleteApplicationFromStore delete from application cache
func (ods *SubscriptionDataStore) DeleteApplicationFromStore(name types.NamespacedName) {
	ods.mu.Lock()
	defer ods.mu.Unlock()
	logger.Info("Deleting application from cache")
	delete(ods.applicationStore, name)
}

// DeleteSubscriptionFromStore delete from subscription cache
func (ods *SubscriptionDataStore) DeleteSubscriptionFromStore(name types.NamespacedName) {
	ods.mu.Lock()
	defer ods.mu.Unlock()
	logger.Info("Deleting subscription from cache")
	delete(ods.subscriptionStore, name)
}

// DeleteApplicationMappingFromStore delete from application mapping cache
func (ods *SubscriptionDataStore) DeleteApplicationMappingFromStore(name types.NamespacedName) {
	ods.mu.Lock()
	defer ods.mu.Unlock()
	logger.Info("Deleting application mapping from cache")
	delete(ods.applicationMappingStore, name)
}
