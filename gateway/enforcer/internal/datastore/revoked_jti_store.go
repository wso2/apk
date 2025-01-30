package datastore

import "time"

type RevokedJTIStore struct {
	revokedJTIs map[string]time.Time // JTI -> expiry time
}

// NewRevokedJTIStore creates a new instance of RevokedJTIStore.
func NewRevokedJTIStore() *RevokedJTIStore {
	return &RevokedJTIStore{
		revokedJTIs: make(map[string]time.Time),
	}
}

// AddJTI adds a JTI to the store with the given expiry time.
func (r *RevokedJTIStore) AddJTI(jti string, expiry time.Time) {
	r.revokedJTIs[jti] = expiry
}

// IsJTIRevoked checks if the given JTI is revoked.
func (r *RevokedJTIStore) IsJTIRevoked(jti string) bool {
	_, ok := r.revokedJTIs[jti]
	return ok
}

// removeExpiredJTIs removes all expired JTIs from the store.
func (r *RevokedJTIStore) removeExpiredJTIs() {
	for jti, expiry := range r.revokedJTIs {
		if time.Now().After(expiry) {
			delete(r.revokedJTIs, jti)
		}
	}
}

// StartRevokedJTIStoreCleanup starts a goroutine to remove expired JTIs from the store periodically.
func (r *RevokedJTIStore) StartRevokedJTIStoreCleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			r.removeExpiredJTIs()
		}
	}()
}

