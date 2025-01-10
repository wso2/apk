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
	"sync"

	api "github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/api"
)

// APIStore is a thread-safe store for APIs.
type APIStore struct {
	apis []*api.Api
	mu   sync.RWMutex
}

// NewAPIStore creates a new instance of APIStore.
func NewAPIStore() *APIStore {
	return &APIStore{
		apis: make([]*api.Api, 0),
	}
}

// AddAPIs adds a list of APIs to the store.
// This method is thread-safe.
func (s *APIStore) AddAPIs(apis []*api.Api) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.apis = apis
}

// GetAPIs retrieves the list of APIs from the store.
// This method is thread-safe.
func (s *APIStore) GetAPIs() []*api.Api {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.apis
}
