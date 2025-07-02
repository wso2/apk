package cache

import (
	"sync"

	dpv2alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v2alpha1"
	"k8s.io/apimachinery/pkg/types"
)

// RouteMetadataDataStore holds RouteMetadata objects,
type RouteMetadataDataStore struct {
	resolveRouteMetadataStore map[string]dpv2alpha1.RouteMetadata
	mu                         sync.Mutex
}

var routeMetadataDataStore *RouteMetadataDataStore

func init() {
	// Initialize the RouteMetadataDataStore
	routeMetadataDataStore = createNewRouteMetadataDataStore()
}

// createNewRouteMetadataDataStore creates a new RouteMetadataDataStore.
func createNewRouteMetadataDataStore() *RouteMetadataDataStore {
	return &RouteMetadataDataStore{
		resolveRouteMetadataStore: make(map[string]dpv2alpha1.RouteMetadata),
	}
}

// GetRouteMetadataDataStore returns the singleton instance of RouteMetadataDataStore.
func GetRouteMetadataDataStore() *RouteMetadataDataStore {
	if routeMetadataDataStore == nil {
		routeMetadataDataStore = createNewRouteMetadataDataStore()
	}
	return routeMetadataDataStore
}

// AddOrUpdateRouteMetadata adds or updates a route metadata entry in the store.
func (rmds *RouteMetadataDataStore) AddOrUpdateRouteMetadata(routeMetadata dpv2alpha1.RouteMetadata) {
	rmds.mu.Lock()
	defer rmds.mu.Unlock()
	namespacedName := types.NamespacedName{
		Namespace: routeMetadata.Namespace,
		Name:      routeMetadata.Name,
	}
	key := namespacedName.String()
	rmds.resolveRouteMetadataStore[key] = routeMetadata
}

// GetRouteMetadata retrieves a route metadata entry from the store.
func (rmds *RouteMetadataDataStore) GetRouteMetadata(namespace, name string) (dpv2alpha1.RouteMetadata, bool) {
	rmds.mu.Lock()
	defer rmds.mu.Unlock()
	namespacedName := types.NamespacedName{
		Namespace: namespace,
		Name:      name,
	}
	key := namespacedName.String()
	routeMetadata, found := rmds.resolveRouteMetadataStore[key]
	return routeMetadata, found
}

// DeleteRouteMetadata deletes a route metadata entry from the store.
func (rmds *RouteMetadataDataStore) DeleteRouteMetadata(namespace, name string) {
	rmds.mu.Lock()
	defer rmds.mu.Unlock()
	namespacedName := types.NamespacedName{
		Namespace: namespace,
		Name:      name,
	}
	key := namespacedName.String()
	delete(rmds.resolveRouteMetadataStore, key)
}

// GetRouteMetadatas retrieves all route metadata entries from the store.
func (rmds *RouteMetadataDataStore) GetRouteMetadatas() map[string]dpv2alpha1.RouteMetadata {
	rmds.mu.Lock()
	defer rmds.mu.Unlock()
	metadatas := make(map[string]dpv2alpha1.RouteMetadata)
	for key, metadata := range rmds.resolveRouteMetadataStore {
		metadatas[key] = metadata
	}
	return metadatas
}
