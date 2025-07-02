

package cache

import (
	"sync"
	dpv2alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v2alpha1"
	"k8s.io/apimachinery/pkg/types"
)

// RoutePolicyDataStore holds RoutePolicy objects,
type RoutePolicyDataStore struct {
	resolveRoutePolicyStore map[string]dpv2alpha1.RoutePolicy
	mu                     sync.Mutex
}

var routePolicyDataStore *RoutePolicyDataStore

func init() {
	// Initialize the RoutePolicyDataStore
	// This is typically done in the main package or where the application starts.
	// It ensures that the data store is ready to use when the application runs.
	routePolicyDataStore = createNewRoutePolicyDataStore()
}

// createNewRoutePolicyDataStore creates a new RoutePolicyDataStore.
func createNewRoutePolicyDataStore() *RoutePolicyDataStore {
	return &RoutePolicyDataStore{
		resolveRoutePolicyStore: make(map[string]dpv2alpha1.RoutePolicy),
	}
}

// GetRoutePolicyDataStore returns the singleton instance of RoutePolicyDataStore.
func GetRoutePolicyDataStore() *RoutePolicyDataStore {
	if routePolicyDataStore == nil {
		routePolicyDataStore = createNewRoutePolicyDataStore()
	}
	return routePolicyDataStore
}

// AddOrUpdateRoutePolicy adds or updates a route policy in the RoutePolicyDataStore.
func (rds *RoutePolicyDataStore) AddOrUpdateRoutePolicy(routePolicy dpv2alpha1.RoutePolicy) {
	rds.mu.Lock()
	defer rds.mu.Unlock()
	namespacedName := types.NamespacedName{
		Namespace: routePolicy.Namespace,
		Name:      routePolicy.Name,
	}
	key := namespacedName.String()
	rds.resolveRoutePolicyStore[key] = routePolicy
}

// GetRoutePolicy retrieves a route policy from the RoutePolicyDataStore.
func (rds *RoutePolicyDataStore) GetRoutePolicy(namespace, name string) (dpv2alpha1.RoutePolicy, bool) {
	rds.mu.Lock()
	defer rds.mu.Unlock()
	namespacedName := types.NamespacedName{
		Namespace: namespace,
		Name:      name,
	}
	key := namespacedName.String()
	routePolicy, found := rds.resolveRoutePolicyStore[key]
	return routePolicy, found
}

// DeleteRoutePolicy deletes a route policy from the RoutePolicyDataStore.
func (rds *RoutePolicyDataStore) DeleteRoutePolicy(namespace, name string) {
	rds.mu.Lock()
	defer rds.mu.Unlock()
	namespacedName := types.NamespacedName{
		Namespace: namespace,
		Name:      name,
	}
	key := namespacedName.String()
	delete(rds.resolveRoutePolicyStore, key)
}

// GetRoutePolicies retrieves all route policies from the RoutePolicyDataStore.
func (rds *RoutePolicyDataStore) GetRoutePolicies() map[string]dpv2alpha1.RoutePolicy {
	rds.mu.Lock()
	defer rds.mu.Unlock()
	routePolicies := make(map[string]dpv2alpha1.RoutePolicy)
	for key, routePolicy := range rds.resolveRoutePolicyStore {
		routePolicies[key] = routePolicy
	}
	return routePolicies
}
