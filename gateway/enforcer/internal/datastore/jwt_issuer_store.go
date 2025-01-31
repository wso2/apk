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
	"fmt"
	"sync"

	subscription "github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/subscription"
)

// JWTIssuerStore is a thread-safe store for APIs.
type JWTIssuerStore struct {
	jwtIssuers map[string]map[string]*subscription.JWTIssuer
	mu         sync.RWMutex
}

// NewJWTIssuerStore creates a new instance of JWTIssuerStore.
func NewJWTIssuerStore() *JWTIssuerStore {
	return &JWTIssuerStore{
		jwtIssuers: make(map[string]map[string]*subscription.JWTIssuer),
	}
}

// AddJWTIssuers adds a list of config to the store.
// This method is thread-safe.
func (s *JWTIssuerStore) AddJWTIssuers(apis []*subscription.JWTIssuer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	orgWizeJWTIssuers := make(map[string]map[string]*subscription.JWTIssuer)
	for _, api := range apis {
		fmt.Printf("Adding JWT Issuer: %v\n", api)
		if _, ok := orgWizeJWTIssuers[api.Organization]; !ok {
			orgWizeJWTIssuers[api.Organization] = make(map[string]*subscription.JWTIssuer)
		}
		orgWizeJWTIssuers[api.Organization][api.Issuer] = api
	}
	fmt.Printf("JWT Issuers: %v\n", orgWizeJWTIssuers)
	s.jwtIssuers = orgWizeJWTIssuers
}

// GetJWTIssuerByOrganizationAndIssuer returns the JWTIssuer for the given organization and issuer.
// This method is thread-safe.
func (s *JWTIssuerStore) GetJWTIssuerByOrganizationAndIssuer(organization, issuer string) *subscription.JWTIssuer {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if orgWiseJWTIssuers, ok := s.jwtIssuers[organization]; ok {
		if jwtIssuer, ok := orgWiseJWTIssuers[issuer]; ok {
			return jwtIssuer
		}
	}
	return nil
}
