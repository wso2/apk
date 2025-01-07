package datastore

import (
	api "github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/api"
	"sync"
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