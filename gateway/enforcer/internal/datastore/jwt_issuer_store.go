package datastore

import (
	subscription "github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/subscription"
	"sync"
)

// JWTIssuerStore is a thread-safe store for APIs.
type JWTIssuerStore struct {
	jwtIssuers []*subscription.JWTIssuer
	mu   sync.RWMutex
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