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

	subscription "github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/subscription"
)

// JWTIssuerStore is a thread-safe store for APIs.
type JWTIssuerStore struct {
	jwtIssuers []*subscription.JWTIssuer
	mu         sync.RWMutex
}

// NewJWTIssuerStore creates a new instance of JWTIssuerStore.
func NewJWTIssuerStore() *JWTIssuerStore {
	return &JWTIssuerStore{
		jwtIssuers: make([]*subscription.JWTIssuer, 0),
	}
}

// AddJWTIssuers adds a list of config to the store.
// This method is thread-safe.
func (s *JWTIssuerStore) AddJWTIssuers(apis []*subscription.JWTIssuer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.jwtIssuers = apis
}

// GetJWTIssuers retrieves the list of Config from the store.
// This method is thread-safe.
func (s *JWTIssuerStore) GetJWTIssuers() []*subscription.JWTIssuer {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.jwtIssuers
}
