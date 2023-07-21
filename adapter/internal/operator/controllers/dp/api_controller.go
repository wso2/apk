/*
 *  Copyright (c) 2022, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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

package dp

import (
	"context"
	"fmt"
	"sync"

	"github.com/wso2/apk/adapter/config"
	"github.com/wso2/apk/adapter/internal/discovery/xds"
	"github.com/wso2/apk/adapter/internal/loggers"
	"github.com/wso2/apk/adapter/internal/operator/constants"
	"github.com/wso2/apk/adapter/internal/operator/status"
	"github.com/wso2/apk/adapter/internal/operator/synchronizer"
	"github.com/wso2/apk/adapter/internal/operator/utils"
	"github.com/wso2/apk/adapter/pkg/logging"
	"golang.org/x/exp/maps"
	k8error "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	gwapiv1b1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	ctrl "sigs.k8s.io/controller-runtime"
	k8client "sigs.k8s.io/controller-runtime/pkg/client"

	dpv1alpha1 "github.com/wso2/apk/adapter/internal/operator/apis/dp/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	httpRouteAPIIndex = "httpRouteAPIIndex"
	// apiAuthenticationIndex Index for API level authentications
	apiAuthenticationIndex = "apiAuthenticationIndex"
	// apiAuthenticationResourceIndex Index for resource level authentications
	apiAuthenticationResourceIndex = "apiAuthenticationResourceIndex"
	// apiRateLimitIndex Index for API level ratelimits
	apiRateLimitIndex = "apiRateLimitIndex"
	// apiRateLimitResourceIndex Index for resource level ratelimits
	apiRateLimitResourceIndex = "apiRateLimitResourceIndex"
	// gatewayHTTPRouteIndex Index for gateway httproutes
	gatewayHTTPRouteIndex = "gatewayHTTPRouteIndex"
	// apiAPIPolicyIndex Index for API level apipolicies
	apiAPIPolicyIndex = "apiAPIPolicyIndex"
	// apiAPIPolicyResourceIndex Index for resource level apipolicies
	apiAPIPolicyResourceIndex        = "apiAPIPolicyResourceIndex"
	serviceHTTPRouteIndex            = "serviceHTTPRouteIndex"
	httprouteScopeIndex              = "httprouteScopeIndex"
	configMapBackend                 = "configMapBackend"
	secretBackend                    = "secretBackend"
	backendHTTPRouteIndex            = "backendHTTPRouteIndex"
	interceptorServiceAPIPolicyIndex = "interceptorServiceAPIPolicyIndex"
	backendInterceptorServiceIndex   = "backendInterceptorServiceIndex"
)

var (
	applyAllAPIsOnce sync.Once
)

// APIReconciler reconciles a API object
type APIReconciler struct {
	client         k8client.Client
	ods            *synchronizer.OperatorDataStore
	ch             *chan synchronizer.APIEvent
	successChannel *chan synchronizer.SuccessEvent
	statusUpdater  *status.UpdateHandler
	mgr            manager.Manager
}

// NewAPIController creates a new API controller instance. API Controllers watches for dpv1alpha1.API and gwapiv1b1.HTTPRoute.
func NewAPIController(mgr manager.Manager, operatorDataStore *synchronizer.OperatorDataStore, statusUpdater *status.UpdateHandler,
	ch *chan synchronizer.APIEvent, successChannel *chan synchronizer.SuccessEvent) error {
	apiReconciler := &APIReconciler{
		client:         mgr.GetClient(),
		ods:            operatorDataStore,
		ch:             ch,
		successChannel: successChannel,
		statusUpdater:  statusUpdater,
		mgr:            mgr,
	}
	ctx := context.Background()

	c, err := controller.New(constants.APIController, mgr, controller.Options{Reconciler: apiReconciler})
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2619, err))
		return err
	}

	conf := config.ReadConfigs()
	predicates := []predicate.Predicate{predicate.NewPredicateFuncs(utils.FilterByNamespaces(conf.Adapter.Operator.Namespaces))}

	if err := c.Watch(&source.Kind{Type: &dpv1alpha1.API{}}, &handler.EnqueueRequestForObject{},
		predicates...); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2611, err))
		return err
	}
	if err := addIndexes(ctx, mgr); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2612, err))
		return err
	}

	if err := c.Watch(&source.Kind{Type: &gwapiv1b1.HTTPRoute{}}, handler.EnqueueRequestsFromMapFunc(apiReconciler.getAPIForHTTPRoute),
		predicates...); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2613, err))
		return err
	}

	if err := c.Watch(&source.Kind{Type: &gwapiv1b1.Gateway{}}, handler.EnqueueRequestsFromMapFunc(apiReconciler.getAPIsForGateway),
		predicates...); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2611, err))
		return err
	}

	if err := c.Watch(&source.Kind{Type: &dpv1alpha1.Backend{}}, handler.EnqueueRequestsFromMapFunc(apiReconciler.getAPIsForBackend),
		predicates...); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2615, err))
		return err
	}

	if err := c.Watch(&source.Kind{Type: &dpv1alpha1.Authentication{}}, handler.EnqueueRequestsFromMapFunc(apiReconciler.getAPIsForAuthentication),
		predicates...); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2616, err))
		return err
	}

	if err := c.Watch(&source.Kind{Type: &dpv1alpha1.InterceptorService{}}, handler.EnqueueRequestsFromMapFunc(apiReconciler.getAPIsForInterceptorService),
		predicates...); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2640, err))
		return err
	}

	if err := c.Watch(&source.Kind{Type: &dpv1alpha1.APIPolicy{}}, handler.EnqueueRequestsFromMapFunc(apiReconciler.getAPIsForAPIPolicy),
		predicates...); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2617, err))
		return err
	}

	if err := c.Watch(&source.Kind{Type: &dpv1alpha1.RateLimitPolicy{}}, handler.EnqueueRequestsFromMapFunc(apiReconciler.getAPIsForRateLimitPolicy),
		predicates...); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2639, err))
		return err
	}

	if err := c.Watch(&source.Kind{Type: &dpv1alpha1.Scope{}}, handler.EnqueueRequestsFromMapFunc(apiReconciler.getAPIsForScope),
		predicates...); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2618, err))
		return err
	}

	if err := c.Watch(&source.Kind{Type: &corev1.ConfigMap{}}, handler.EnqueueRequestsFromMapFunc(apiReconciler.getAPIsForConfigMap),
		predicates...); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2644, err))
		return err
	}

	if err := c.Watch(&source.Kind{Type: &corev1.Secret{}}, handler.EnqueueRequestsFromMapFunc(apiReconciler.getAPIsForSecret),
		predicates...); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2645, err))
		return err
	}

	loggers.LoggerAPKOperator.Info("API Controller successfully started. Watching API Objects....")
	go apiReconciler.handleStatus()
	return nil
}

// +kubebuilder:rbac:groups=dp.wso2.com,resources=apis,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=dp.wso2.com,resources=apis/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=dp.wso2.com,resources=apis/finalizers,verbs=update
// +kubebuilder:rbac:groups=dp.wso2.com,resources=httproutes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=dp.wso2.com,resources=httproutes/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=dp.wso2.com,resources=httproutes/finalizers,verbs=update
// +kubebuilder:rbac:groups=dp.wso2.com,resources=authentications,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=dp.wso2.com,resources=authentications/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=dp.wso2.com,resources=authentications/finalizers,verbs=update
// +kubebuilder:rbac:groups=dp.wso2.com,resources=apipolicies,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=dp.wso2.com,resources=apipolicies/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=dp.wso2.com,resources=apipolicies/finalizers,verbs=update
// +kubebuilder:rbac:groups=dp.wso2.com,resources=scopes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=dp.wso2.com,resources=scopes/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=dp.wso2.com,resources=scopes/finalizers,verbs=update
// +kubebuilder:rbac:groups=dp.wso2.com,resources=ratelimitpolicies,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=dp.wso2.com,resources=ratelimitpolicies/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=dp.wso2.com,resources=ratelimitpolicies/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (apiReconciler *APIReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	applyAllAPIsOnce.Do(apiReconciler.applyStartupAPIs)
	loggers.LoggerAPKOperator.Infof("Reconciling for API %s", req.NamespacedName.String())
	// Check whether the API CR exist, if not consider as a DELETE event.
	var apiCR dpv1alpha1.API
	if err := apiReconciler.client.Get(ctx, req.NamespacedName, &apiCR); err != nil {
		apiState, found := apiReconciler.ods.GetCachedAPI(req.NamespacedName)
		if found && k8error.IsNotFound(err) {
			// The api doesn't exist in the api Cache, remove it
			apiReconciler.ods.DeleteCachedAPI(req.NamespacedName)
			loggers.LoggerAPKOperator.Infof("Delete event has received for API : %s, hence deleted from API cache", req.NamespacedName.String())
			*apiReconciler.ch <- synchronizer.APIEvent{EventType: constants.Delete, Event: apiState}
			return ctrl.Result{}, nil
		}
		loggers.LoggerAPKOperator.Warnf("Api CR related to the reconcile request with key: %s returned error. Assuming API is already deleted, hence ignoring the error : %v", req.NamespacedName.String(), err)
		return ctrl.Result{}, nil
	}

	if apiState, err := apiReconciler.resolveAPIRefs(ctx, apiCR); err != nil {
		loggers.LoggerAPKOperator.Warnf("Error retrieving ref CRs for API in namespace : %s, %v", req.NamespacedName.String(), err)
		return ctrl.Result{}, err
	} else if apiState != nil {
		loggers.LoggerAPKOperator.Infof("Ready to deploy CRs for API in namespace : %s, %v", req.NamespacedName.String(), err)
		*apiReconciler.ch <- *apiState
	}
	return ctrl.Result{}, nil
}

// applyStartupAPIs applies the APIs which are already available in the cluster at the startup of the operator.
func (apiReconciler *APIReconciler) applyStartupAPIs() {
	ctx := context.Background()
	apisList, err := retrieveAPIList(apiReconciler.client)
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2601, err))
		return
	}
	for _, api := range apisList {
		if apiState, err := apiReconciler.resolveAPIRefs(ctx, api); err != nil {
			loggers.LoggerAPKOperator.Warnf("Error retrieving ref CRs for API : %s in namespace : %s, %v", api.Name, api.Namespace, err)
		} else if apiState != nil {
			*apiReconciler.ch <- *apiState
		}
	}
	xds.SetReady()
	loggers.LoggerAPKOperator.Info("Initial APIs were deployed successfully")
}

func retrieveAPIList(k8sclient k8client.Client) ([]dpv1alpha1.API, error) {
	ctx := context.Background()
	conf := config.ReadConfigs()
	namespaces := conf.Adapter.Operator.Namespaces
	var apis []dpv1alpha1.API
	if namespaces == nil {
		apiList := &dpv1alpha1.APIList{}
		if err := k8sclient.List(ctx, apiList, &k8client.ListOptions{}); err != nil {
			return nil, err
		}
		apis = make([]dpv1alpha1.API, len(apiList.Items))
		copy(apis[:], apiList.Items[:])
	} else {
		for _, namespace := range namespaces {
			apiList := &dpv1alpha1.APIList{}
			if err := k8sclient.List(ctx, apiList, &k8client.ListOptions{Namespace: namespace}); err != nil {
				return nil, err
			}
			apis = append(apis, apiList.Items...)
		}
	}
	return apis, nil
}

// resolveAPIRefs validates following references related to the API
// - HTTPRoutes
func (apiReconciler *APIReconciler) resolveAPIRefs(ctx context.Context, api dpv1alpha1.API) (*synchronizer.APIEvent, error) {
	var prodHTTPRouteRefs, sandHTTPRouteRefs []string
	if len(api.Spec.Production) > 0 {
		prodHTTPRouteRefs = api.Spec.Production[0].HTTPRouteRefs
	}
	if len(api.Spec.Sandbox) > 0 {
		sandHTTPRouteRefs = api.Spec.Sandbox[0].HTTPRouteRefs
	}

	apiState := &synchronizer.APIState{
		APIDefinition: &api,
	}
	var err error
	apiRef := utils.NamespacedName(&api)
	namespace := api.Namespace
	if apiState.Authentications, err = apiReconciler.getAuthenticationsForAPI(ctx, api); err != nil {
		return nil, fmt.Errorf("error while getting API level auth for API : %s in namespace :%s, %s", apiRef.String(),
			namespace, err.Error())
	}
	if apiState.RateLimitPolicies, err = apiReconciler.getRatelimitPoliciesForAPI(ctx, api); err != nil {
		return nil, fmt.Errorf("error while getting API level ratelimit for API : %s in namespace :%s, %s", apiRef.String(),
			namespace, err.Error())
	}
	if apiState.APIPolicies, err = apiReconciler.getAPIPoliciesForAPI(ctx, api); err != nil {
		return nil, fmt.Errorf("error while getting API level apipolicy for API : %s in namespace :%s, %s", apiRef.String(),
			namespace, err.Error())
	}

	if apiState.ResourceAuthentications, err = apiReconciler.getAuthenticationsForResources(ctx, api); err != nil {
		return nil, fmt.Errorf("error while getting httproute resource auth : %s in namespace :%s, %s", apiRef.String(),
			namespace, err.Error())
	}
	if apiState.ResourceRateLimitPolicies, err = apiReconciler.getRatelimitPoliciesForResources(ctx, api); err != nil {
		return nil, fmt.Errorf("error while getting httproute resource ratelimit : %s in namespace :%s, %s", apiRef.String(),
			namespace, err.Error())
	}
	if apiState.ResourceAPIPolicies, err = apiReconciler.getAPIPoliciesForResources(ctx, api); err != nil {
		return nil, fmt.Errorf("error while getting httproute resource apipolicy %s in namespace :%s, %s", apiRef.String(),
			namespace, err.Error())
	}
	if apiState.InterceptorServiceMapping, err = apiReconciler.getInterceptorServices(ctx, apiState.APIPolicies,
		apiState.ResourceAPIPolicies, api); err != nil {
		return nil, fmt.Errorf("error while getting interceptor services %s in namespace :%s, %s", apiRef.String(),
			namespace, err.Error())
	}
	if api.Spec.DefinitionFileRef != "" {
		if apiState.APIDefinitionFile, err = apiReconciler.getAPIDefinitionForAPI(ctx, api.Spec.DefinitionFileRef, namespace, api); err != nil {
			return nil, fmt.Errorf("error while getting api definition file of api %s in namespace :%s, %s", apiRef.String(),
				namespace, err.Error())
		}
	}

	if len(prodHTTPRouteRefs) > 0 {
		apiState.ProdHTTPRoute = &synchronizer.HTTPRouteState{}
		if apiState.ProdHTTPRoute, err = apiReconciler.resolveHTTPRouteRefs(ctx, apiState.ProdHTTPRoute, prodHTTPRouteRefs,
			namespace, apiState.InterceptorServiceMapping, api); err != nil {
			return nil, fmt.Errorf("error while resolving production httpRouteref %s in namespace :%s has not found. %s",
				prodHTTPRouteRefs, namespace, err.Error())
		}
		if !apiReconciler.ods.IsGatewayAvailable(types.NamespacedName{
			Name: string(apiState.ProdHTTPRoute.HTTPRoute.Spec.ParentRefs[0].Name),
			Namespace: utils.GetNamespace(apiState.ProdHTTPRoute.HTTPRoute.Spec.ParentRefs[0].Namespace,
				apiState.ProdHTTPRoute.HTTPRoute.Namespace),
		}) {
			return nil, fmt.Errorf("no gateway available for httpRouteref %s in namespace :%s has not found",
				prodHTTPRouteRefs, namespace)
		}
	}

	if len(sandHTTPRouteRefs) > 0 {
		apiState.SandHTTPRoute = &synchronizer.HTTPRouteState{}
		if apiState.SandHTTPRoute, err = apiReconciler.resolveHTTPRouteRefs(ctx, apiState.SandHTTPRoute, sandHTTPRouteRefs,
			namespace, apiState.InterceptorServiceMapping, api); err != nil {
			return nil, fmt.Errorf("error while resolving sandbox httpRouteref %s in namespace :%s has not found. %s",
				sandHTTPRouteRefs, namespace, err.Error())
		}
		if !apiReconciler.ods.IsGatewayAvailable(types.NamespacedName{
			Name:      string(apiState.SandHTTPRoute.HTTPRoute.Spec.ParentRefs[0].Name),
			Namespace: utils.GetNamespace(apiState.SandHTTPRoute.HTTPRoute.Spec.ParentRefs[0].Namespace, apiState.SandHTTPRoute.HTTPRoute.Namespace),
		}) {
			return nil, fmt.Errorf("no gateway available for httpRouteref %s in namespace :%s has not found",
				sandHTTPRouteRefs, namespace)
		}
	}

	loggers.LoggerAPKOperator.Debugf("HTTPRoutes are retrieved successfully for API CR %s", apiRef.String())

	if !api.APIDeploymentEvent.Accepted {
		apiReconciler.ods.AddAPIState(apiRef, apiState)
		return &synchronizer.APIEvent{EventType: constants.Create, Event: *apiState, UpdatedEvents: []string{}}, nil
	} else if cachedAPI, events, updated :=
		apiReconciler.ods.UpdateAPIState(apiRef, apiState); updated {
		apiReconciler.removeOldOwnerRefs(ctx, cachedAPI)
		return &synchronizer.APIEvent{EventType: constants.Update, Event: cachedAPI, UpdatedEvents: events}, nil
	}

	return nil, nil
}

func (apiReconciler *APIReconciler) removeOldOwnerRefs(ctx context.Context, apiState synchronizer.APIState) {
	api := apiState.APIDefinition
	if apiState.ProdHTTPRoute != nil {
		apiReconciler.removeOldOwnerRefsFromHTTPRoute(ctx, apiState.ProdHTTPRoute, api.Name, api.Namespace)
	}
	if apiState.SandHTTPRoute != nil {
		apiReconciler.removeOldOwnerRefsFromHTTPRoute(ctx, apiState.SandHTTPRoute, api.Name, api.Namespace)
	}

	// remove old owner refs from interceptor services
	interceptorServiceList := &dpv1alpha1.InterceptorServiceList{}
	if err := apiReconciler.client.List(ctx, interceptorServiceList, &k8client.ListOptions{
		Namespace: api.Namespace,
	}); err != nil {
		loggers.LoggerAPKOperator.Errorf("error while listing CRs for API CR %s, %s",
			api.Name, err.Error())
	}
	for _, interceptorService := range interceptorServiceList.Items {
		// check interceptorService has similar item inside the apiState.InterceptorServiceMapping
		interceptorServiceFound := false
		for _, attachedInterceptorService := range apiState.InterceptorServiceMapping {
			if attachedInterceptorService.Name == interceptorService.Name {
				interceptorServiceFound = true
				break
			}
		}
		if !interceptorServiceFound {
			// remove owner reference
			apiReconciler.removeOldOwnerRefsFromChild(ctx, &interceptorService, api.Name, api.Namespace)
		}
	}
}

func (apiReconciler *APIReconciler) removeOldOwnerRefsFromHTTPRoute(ctx context.Context,
	httpRouteState *synchronizer.HTTPRouteState, apiName, apiNamespace string) {
	httpRoute := httpRouteState.HTTPRoute
	if httpRoute != nil {
		// scope CRs
		scopeList := &dpv1alpha1.ScopeList{}
		if err := apiReconciler.client.List(ctx, scopeList, &k8client.ListOptions{
			Namespace: apiNamespace,
		}); err != nil {
			loggers.LoggerAPKOperator.Errorf("error while listing authentication CRs for API CR %s, %s",
				apiName, err.Error())
		}
		for _, scope := range scopeList.Items {
			// check scope has similar item inside the apiState.ProdHTTPRoute.Scopes
			scopeFound := false
			for _, attachedScope := range httpRouteState.Scopes {
				if scope.GetName() == attachedScope.GetName() {
					scopeFound = true
					break
				}
			}
			if !scopeFound {
				apiReconciler.removeOldOwnerRefsFromChild(ctx, &scope, apiName, apiNamespace)
			}
		}

		// backend CRs
		backendList := &dpv1alpha1.BackendList{}
		if err := apiReconciler.client.List(ctx, backendList, &k8client.ListOptions{
			Namespace: apiNamespace,
		}); err != nil {
			loggers.LoggerAPKOperator.Errorf("error while listing authentication CRs for API CR %s, %s",
				apiName, err.Error())
		}
		for _, backend := range backendList.Items {
			// check backend has similar item inside the apiState.ProdHTTPRoute.Backends
			backendFound := false
			for _, attachedBackend := range httpRouteState.BackendMapping {
				if backend.GetName() == attachedBackend.Backend.Name {
					backendFound = true
					break
				}
			}
			if !backendFound {
				apiReconciler.removeOldOwnerRefsFromChild(ctx, &backend, apiName, apiNamespace)
			}
		}
	}
}

func (apiReconciler *APIReconciler) removeOldOwnerRefsFromChild(ctx context.Context, child k8client.Object,
	apiName, apiNamespace string) {
	ownerReferences := child.GetOwnerReferences()
	for i, ownerRef := range ownerReferences {
		if ownerRef.Kind == "API" && ownerRef.Name == apiName {
			// delete the element from ownerReferences list
			ownerReferences = append(ownerReferences[:i], ownerReferences[i+1:]...)
			child.SetOwnerReferences(ownerReferences)
			if err := utils.UpdateCR(ctx, apiReconciler.client, child); err != nil {
				loggers.LoggerAPKOperator.Errorf("error while updating CR %s, %s",
					child.GetName(), err.Error())
			}
			break
		}
	}
}

// resolveHTTPRouteRefs validates following references related to the API
// - Authentications
func (apiReconciler *APIReconciler) resolveHTTPRouteRefs(ctx context.Context, httpRouteState *synchronizer.HTTPRouteState,
	httpRouteRefs []string, namespace string, interceptorServiceMapping map[string]dpv1alpha1.InterceptorService,
	api dpv1alpha1.API) (*synchronizer.HTTPRouteState, error) {
	var err error
	httpRouteState.HTTPRoute, err = apiReconciler.concatHTTPRoutes(ctx, httpRouteRefs, namespace, api)
	if err != nil {
		return nil, err
	}
	httpRouteState.BackendMapping = apiReconciler.getResolvedBackendsMapping(ctx, httpRouteState, interceptorServiceMapping, api)
	httpRouteState.Scopes, err = apiReconciler.getScopesForHTTPRoute(ctx, httpRouteState.HTTPRoute, api)
	return httpRouteState, err
}

func (apiReconciler *APIReconciler) concatHTTPRoutes(ctx context.Context, httpRouteRefs []string,
	namespace string, api dpv1alpha1.API) (*gwapiv1b1.HTTPRoute, error) {
	var combinedHTTPRoute *gwapiv1b1.HTTPRoute
	for _, httpRouteRef := range httpRouteRefs {
		var httpRoute gwapiv1b1.HTTPRoute
		if err := utils.ResolveRef(ctx, apiReconciler.client, &api,
			types.NamespacedName{Namespace: namespace, Name: httpRouteRef}, true, &httpRoute); err != nil {
			return nil, fmt.Errorf("error while getting httproute %s in namespace :%s, %s", httpRouteRef, namespace, err.Error())
		}
		if combinedHTTPRoute == nil {
			combinedHTTPRoute = &httpRoute
		} else {
			combinedHTTPRoute.Spec.Rules = append(combinedHTTPRoute.Spec.Rules, httpRoute.Spec.Rules...)
		}
	}
	return combinedHTTPRoute, nil
}

func (apiReconciler *APIReconciler) getAuthenticationsForAPI(ctx context.Context,
	api dpv1alpha1.API) (map[string]dpv1alpha1.Authentication, error) {
	nameSpacedName := utils.NamespacedName(&api).String()
	authentications := make(map[string]dpv1alpha1.Authentication)
	authenticationList := &dpv1alpha1.AuthenticationList{}
	if err := apiReconciler.client.List(ctx, authenticationList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(apiAuthenticationIndex, nameSpacedName),
	}); err != nil {
		return nil, err
	}
	for _, item := range authenticationList.Items {
		if err := utils.UpdateOwnerReference(ctx, apiReconciler.client, &item, api, true); err != nil {
			return nil, err
		}
		authentications[utils.NamespacedName(&item).String()] = item
	}
	return authentications, nil
}

func (apiReconciler *APIReconciler) getRatelimitPoliciesForAPI(ctx context.Context,
	api dpv1alpha1.API) (map[string]dpv1alpha1.RateLimitPolicy, error) {
	nameSpacedName := utils.NamespacedName(&api).String()
	ratelimitPolicies := make(map[string]dpv1alpha1.RateLimitPolicy)
	ratelimitPolicyList := &dpv1alpha1.RateLimitPolicyList{}
	if err := apiReconciler.client.List(ctx, ratelimitPolicyList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(apiRateLimitIndex, nameSpacedName),
	}); err != nil {
		return nil, err
	}
	for _, item := range ratelimitPolicyList.Items {
		if err := utils.UpdateOwnerReference(ctx, apiReconciler.client, &item, api, true); err != nil {
			return nil, err
		}
		ratelimitPolicies[utils.NamespacedName(&item).String()] = item
	}
	return ratelimitPolicies, nil
}

func (apiReconciler *APIReconciler) getScopesForHTTPRoute(ctx context.Context,
	httpRoute *gwapiv1b1.HTTPRoute, api dpv1alpha1.API) (map[string]dpv1alpha1.Scope, error) {
	scopes := make(map[string]dpv1alpha1.Scope)
	for _, rule := range httpRoute.Spec.Rules {
		for _, filter := range rule.Filters {
			if filter.Type == gwapiv1b1.HTTPRouteFilterExtensionRef && filter.ExtensionRef != nil &&
				filter.ExtensionRef.Kind == constants.KindScope {
				scope := &dpv1alpha1.Scope{}
				if err := utils.ResolveRef(ctx, apiReconciler.client, &api,
					types.NamespacedName{Namespace: httpRoute.Namespace, Name: string(filter.ExtensionRef.Name)}, false,
					scope); err != nil {
					return nil, fmt.Errorf("error while getting scope %s in namespace :%s, %s", filter.ExtensionRef.Name,
						httpRoute.Namespace, err.Error())
				}
				scopes[utils.NamespacedName(scope).String()] = *scope
			}
		}
	}

	return scopes, nil
}

func (apiReconciler *APIReconciler) getAuthenticationsForResources(ctx context.Context,
	api dpv1alpha1.API) (map[string]dpv1alpha1.Authentication, error) {
	nameSpacedName := utils.NamespacedName(&api).String()
	authentications := make(map[string]dpv1alpha1.Authentication)
	authenticationList := &dpv1alpha1.AuthenticationList{}
	if err := apiReconciler.client.List(ctx, authenticationList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(apiAuthenticationResourceIndex, nameSpacedName),
	}); err != nil {
		return nil, err
	}
	for _, item := range authenticationList.Items {
		if err := utils.UpdateOwnerReference(ctx, apiReconciler.client, &item, api, true); err != nil {
			return nil, err
		}
		authentications[utils.NamespacedName(&item).String()] = item
	}
	return authentications, nil
}

func (apiReconciler *APIReconciler) getRatelimitPoliciesForResources(ctx context.Context,
	api dpv1alpha1.API) (map[string]dpv1alpha1.RateLimitPolicy, error) {
	nameSpacedName := utils.NamespacedName(&api).String()
	ratelimitpolicies := make(map[string]dpv1alpha1.RateLimitPolicy)
	ratelimitPolicyList := &dpv1alpha1.RateLimitPolicyList{}
	if err := apiReconciler.client.List(ctx, ratelimitPolicyList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(apiRateLimitResourceIndex, nameSpacedName),
	}); err != nil {
		return nil, err
	}
	for _, item := range ratelimitPolicyList.Items {
		if err := utils.UpdateOwnerReference(ctx, apiReconciler.client, &item, api, true); err != nil {
			return nil, err
		}
		ratelimitpolicies[utils.NamespacedName(&item).String()] = item
	}
	return ratelimitpolicies, nil
}

func (apiReconciler *APIReconciler) getAPIPoliciesForAPI(ctx context.Context,
	api dpv1alpha1.API) (map[string]dpv1alpha1.APIPolicy, error) {
	nameSpacedName := utils.NamespacedName(&api).String()
	apiPolicies := make(map[string]dpv1alpha1.APIPolicy)
	apiPolicyList := &dpv1alpha1.APIPolicyList{}
	if err := apiReconciler.client.List(ctx, apiPolicyList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(apiAPIPolicyIndex, nameSpacedName),
	}); err != nil {
		return nil, err
	}
	for _, item := range apiPolicyList.Items {
		if err := utils.UpdateOwnerReference(ctx, apiReconciler.client, &item, api, true); err != nil {
			return nil, err
		}
		apiPolicies[utils.NamespacedName(&item).String()] = item
	}
	return apiPolicies, nil
}

func (apiReconciler *APIReconciler) getAPIDefinitionForAPI(ctx context.Context,
	apiDefinitionFile, namespace string, api dpv1alpha1.API) ([]byte, error) {
	configMap := &corev1.ConfigMap{}
	if err := utils.ResolveRef(ctx, apiReconciler.client, &api,
		types.NamespacedName{Namespace: namespace, Name: apiDefinitionFile}, true, configMap); err != nil {
		return nil, fmt.Errorf("error while getting swagger definition %s in namespace :%s, %s", apiDefinitionFile,
			namespace, err.Error())
	}
	apiDef := make(map[string][]byte)
	for _, val := range configMap.BinaryData {
		// config map data key is "swagger.yaml"
		apiDef["apiDef"] = []byte(val)
	}
	return apiDef["apiDef"], nil
}

func (apiReconciler *APIReconciler) getAPIPoliciesForResources(ctx context.Context,
	api dpv1alpha1.API) (map[string]dpv1alpha1.APIPolicy, error) {
	nameSpacedName := utils.NamespacedName(&api).String()
	apiPolicies := make(map[string]dpv1alpha1.APIPolicy)
	apiPolicyList := &dpv1alpha1.APIPolicyList{}
	if err := apiReconciler.client.List(ctx, apiPolicyList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(apiAPIPolicyResourceIndex, nameSpacedName),
	}); err != nil {
		return nil, err
	}
	for _, item := range apiPolicyList.Items {
		if err := utils.UpdateOwnerReference(ctx, apiReconciler.client, &item, api, true); err != nil {
			return nil, err
		}
		apiPolicies[utils.NamespacedName(&item).String()] = item
	}
	return apiPolicies, nil
}

// getInterceptorServices gets all the interceptor services for the resolving API
func (apiReconciler *APIReconciler) getInterceptorServices(ctx context.Context,
	apiPolicies, resourceAPIPolicies map[string]dpv1alpha1.APIPolicy, api dpv1alpha1.API) (map[string]dpv1alpha1.InterceptorService, error) {
	allAPIPolicies := append(maps.Values(apiPolicies), maps.Values(resourceAPIPolicies)...)
	interceptorServices := make(map[string]dpv1alpha1.InterceptorService)
	for _, apiPolicy := range allAPIPolicies {
		if apiPolicy.Spec.Default != nil && len(apiPolicy.Spec.Default.RequestInterceptors) > 0 {
			interceptorPtr := utils.GetInterceptorService(ctx, apiReconciler.client, &apiPolicy.Spec.Default.RequestInterceptors[0], &api)
			if interceptorPtr != nil {
				interceptorServices[utils.NamespacedName(interceptorPtr).String()] = *interceptorPtr
			}
		}
		if apiPolicy.Spec.Default != nil && len(apiPolicy.Spec.Default.ResponseInterceptors) > 0 {
			interceptorPtr := utils.GetInterceptorService(ctx, apiReconciler.client, &apiPolicy.Spec.Default.ResponseInterceptors[0], &api)
			if interceptorPtr != nil {
				interceptorServices[utils.NamespacedName(interceptorPtr).String()] = *interceptorPtr
			}
		}
		if apiPolicy.Spec.Override != nil && len(apiPolicy.Spec.Override.RequestInterceptors) > 0 {
			interceptorPtr := utils.GetInterceptorService(ctx, apiReconciler.client, &apiPolicy.Spec.Override.RequestInterceptors[0], &api)
			if interceptorPtr != nil {
				interceptorServices[utils.NamespacedName(interceptorPtr).String()] = *interceptorPtr
			}
		}
		if apiPolicy.Spec.Override != nil && len(apiPolicy.Spec.Override.ResponseInterceptors) > 0 {
			interceptorPtr := utils.GetInterceptorService(ctx, apiReconciler.client, &apiPolicy.Spec.Override.ResponseInterceptors[0], &api)
			if interceptorPtr != nil {
				interceptorServices[utils.NamespacedName(interceptorPtr).String()] = *interceptorPtr
			}
		}
	}
	return interceptorServices, nil
}

func (apiReconciler *APIReconciler) getResolvedBackendsMapping(ctx context.Context,
	httpRouteState *synchronizer.HTTPRouteState, interceptorServiceMapping map[string]dpv1alpha1.InterceptorService,
	api dpv1alpha1.API) dpv1alpha1.BackendMapping {
	backendMapping := make(dpv1alpha1.BackendMapping)

	// Resolve backends in HTTPRoute
	httpRoute := httpRouteState.HTTPRoute
	for _, rule := range httpRoute.Spec.Rules {
		for _, backend := range rule.BackendRefs {
			backendNamespacedName := types.NamespacedName{
				Name:      string(backend.Name),
				Namespace: utils.GetNamespace(backend.Namespace, httpRoute.Namespace),
			}
			resolvedBackend := utils.GetResolvedBackend(ctx, apiReconciler.client, backendNamespacedName, &api)
			if resolvedBackend != nil {
				backendMapping[backendNamespacedName] = resolvedBackend
			}
		}
	}

	// Resolve backends in InterceptorServices
	interceptorServices := maps.Values(interceptorServiceMapping)
	for _, interceptorService := range interceptorServices {
		utils.ResolveAndAddBackendToMapping(ctx, apiReconciler.client, backendMapping,
			interceptorService.Spec.BackendRef, interceptorService.Namespace, &api)
	}

	loggers.LoggerAPKOperator.Debugf("Generated backendMapping: %v", backendMapping)
	return backendMapping
}

// getAPIForHTTPRoute triggers the API controller reconcile method based on the changes detected
// from HTTPRoute objects. If the changes are done for an API stored in the Operator Data store,
// a new reconcile event will be created and added to the reconcile event queue.
func (apiReconciler *APIReconciler) getAPIForHTTPRoute(obj k8client.Object) []reconcile.Request {
	ctx := context.Background()
	httpRoute, ok := obj.(*gwapiv1b1.HTTPRoute)
	if !ok {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2622, httpRoute))
		return []reconcile.Request{}
	}

	apiList := &dpv1alpha1.APIList{}
	if err := apiReconciler.client.List(ctx, apiList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(httpRouteAPIIndex, utils.NamespacedName(httpRoute).String()),
	}); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2623, utils.NamespacedName(httpRoute).String()))
		return []reconcile.Request{}
	}

	if len(apiList.Items) == 0 {
		loggers.LoggerAPKOperator.Debugf("APIs for HTTPRoute not found: %s", utils.NamespacedName(httpRoute).String())
		return []reconcile.Request{}
	}

	requests := []reconcile.Request{}
	for _, api := range apiList.Items {
		req := reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      api.Name,
				Namespace: api.Namespace},
		}
		requests = append(requests, req)
		loggers.LoggerAPKOperator.Infof("Adding reconcile request for API: %s/%s", api.Namespace, api.Name)
	}
	return requests
}

// getAPIsForConfigMap triggers the API controller reconcile method based on the changes detected
// in configMap resources.
func (apiReconciler *APIReconciler) getAPIsForConfigMap(obj k8client.Object) []reconcile.Request {
	ctx := context.Background()
	configMap, ok := obj.(*corev1.ConfigMap)
	if !ok {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2622, configMap))
		return []reconcile.Request{}
	}

	backendList := &dpv1alpha1.BackendList{}
	if err := apiReconciler.client.List(ctx, backendList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(configMapBackend, utils.NamespacedName(configMap).String()),
	}); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2638, utils.NamespacedName(configMap).String()))
		return []reconcile.Request{}
	}

	requests := []reconcile.Request{}
	for _, backend := range backendList.Items {
		requests = append(requests, apiReconciler.getAPIsForBackend(&backend)...)
	}
	return requests
}

// getAPIsForSecret triggers the API controller reconcile method based on the changes detected
// in secret resources.
func (apiReconciler *APIReconciler) getAPIsForSecret(obj k8client.Object) []reconcile.Request {
	ctx := context.Background()
	secret, ok := obj.(*corev1.Secret)
	if !ok {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2622, secret))
		return []reconcile.Request{}
	}

	backendList := &dpv1alpha1.BackendList{}
	if err := apiReconciler.client.List(ctx, backendList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(secretBackend, utils.NamespacedName(secret).String()),
	}); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2621, utils.NamespacedName(secret).String()))
		return []reconcile.Request{}
	}

	requests := []reconcile.Request{}
	for _, backend := range backendList.Items {
		requests = append(requests, apiReconciler.getAPIsForBackend(&backend)...)
	}
	return requests
}

// getAPIForAuthentication triggers the API controller reconcile method based on the changes detected
// from Authentication objects. If the changes are done for an API stored in the Operator Data store,
// a new reconcile event will be created and added to the reconcile event queue.
func (apiReconciler *APIReconciler) getAPIsForAuthentication(obj k8client.Object) []reconcile.Request {
	authentication, ok := obj.(*dpv1alpha1.Authentication)
	if !ok {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2624, authentication))
		return []reconcile.Request{}
	}

	requests := []reconcile.Request{}

	// todo(amali) move this validation to validation hook
	if !(authentication.Spec.TargetRef.Kind == constants.KindAPI || authentication.Spec.TargetRef.Kind == constants.KindResource) {
		loggers.LoggerAPKOperator.Errorf("Unsupported target ref kind : %s was given for authentication: %s",
			authentication.Spec.TargetRef.Kind, authentication.Name)
		return requests
	}

	namespace := utils.GetNamespace((*gwapiv1b1.Namespace)(authentication.Spec.TargetRef.Namespace), authentication.Namespace)

	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      string(authentication.Spec.TargetRef.Name),
			Namespace: namespace,
		},
	}
	requests = append(requests, req)
	loggers.LoggerAPKOperator.Infof("Adding reconcile request for API: %s/%s", string(authentication.Spec.TargetRef.Name), namespace)

	return requests
}

// getAPIForAuthentication triggers the API controller reconcile method based on the changes detected
// from Authentication objects. If the changes are done for an API stored in the Operator Data store,
// a new reconcile event will be created and added to the reconcile event queue.
func (apiReconciler *APIReconciler) getAPIsForAPIPolicy(obj k8client.Object) []reconcile.Request {
	apiPolicy, ok := obj.(*dpv1alpha1.APIPolicy)
	if !ok {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2624, apiPolicy))
		return []reconcile.Request{}
	}
	requests := []reconcile.Request{}

	if apiPolicy.Spec.TargetRef.Kind == constants.KindGateway {
		return []reconcile.Request{}
	}

	// todo(amali) move this validation to validation hook
	if !(apiPolicy.Spec.TargetRef.Kind == constants.KindAPI || apiPolicy.Spec.TargetRef.Kind == constants.KindResource) {
		loggers.LoggerAPKOperator.Errorf("Unsupported target ref kind : %s was given for authentication: %s",
			apiPolicy.Spec.TargetRef.Kind, apiPolicy.Name)
		return requests
	}

	namespace := utils.GetNamespace((*gwapiv1b1.Namespace)(apiPolicy.Spec.TargetRef.Namespace), apiPolicy.Namespace)

	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      string(apiPolicy.Spec.TargetRef.Name),
			Namespace: namespace},
	}
	requests = append(requests, req)
	loggers.LoggerAPKOperator.Infof("Adding reconcile request for API: %s/%s", string(apiPolicy.Spec.TargetRef.Name), namespace)

	return requests
}

// getAPIPoliciesForInterceptorService returns associated APIPolicies for the InterceptorService
// when the changes detected in InterceptorService resoruces.
func (apiReconciler *APIReconciler) getAPIsForInterceptorService(obj k8client.Object) []reconcile.Request {
	interceptorService, ok := obj.(*dpv1alpha1.InterceptorService)
	if !ok {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2624, interceptorService))
		return []reconcile.Request{}
	}

	ctx := context.Background()
	apiPolicyList := &dpv1alpha1.APIPolicyList{}
	if err := apiReconciler.client.List(ctx, apiPolicyList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(interceptorServiceAPIPolicyIndex, utils.NamespacedName(interceptorService).String()),
	}); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2651, utils.NamespacedName(interceptorService).String(), err.Error()))
		return []reconcile.Request{}
	}

	requests := []reconcile.Request{}
	for _, apiPolicy := range apiPolicyList.Items {
		requests = append(requests, apiReconciler.getAPIsForAPIPolicy(&apiPolicy)...)
	}
	return requests
}

// getAPIForAuthentication triggers the API controller reconcile method based on the changes detected
// from Authentication objects. If the changes are done for an API stored in the Operator Data store,
// a new reconcile event will be created and added to the reconcile event queue.
func (apiReconciler *APIReconciler) getAPIsForRateLimitPolicy(obj k8client.Object) []reconcile.Request {
	ratelimitPolicy, ok := obj.(*dpv1alpha1.RateLimitPolicy)
	if !ok {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2624, ratelimitPolicy))
		return []reconcile.Request{}
	}
	requests := []reconcile.Request{}

	// todo(amali) move this validation to validation hook
	if !(ratelimitPolicy.Spec.TargetRef.Kind == constants.KindAPI || ratelimitPolicy.Spec.TargetRef.Kind == constants.KindResource) {
		loggers.LoggerAPKOperator.Errorf("Unsupported target ref kind : %s was given for authentication: %s",
			ratelimitPolicy.Spec.TargetRef.Kind, ratelimitPolicy.Name)
		return requests
	}

	namespace := utils.GetNamespace((*gwapiv1b1.Namespace)(ratelimitPolicy.Spec.TargetRef.Namespace), ratelimitPolicy.Namespace)

	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      string(ratelimitPolicy.Spec.TargetRef.Name),
			Namespace: namespace},
	}
	requests = append(requests, req)
	loggers.LoggerAPKOperator.Infof("Adding reconcile request for API: %s/%s", string(ratelimitPolicy.Spec.TargetRef.Name), namespace)

	return requests
}

// getAPIsForScope triggers the API controller reconcile method based on the changes detected
// from scope objects. If the changes are done for an API stored in the Operator Data store,
// a new reconcile event will be created and added to the reconcile event queue.
func (apiReconciler *APIReconciler) getAPIsForScope(obj k8client.Object) []reconcile.Request {
	scope, ok := obj.(*dpv1alpha1.Scope)
	if !ok {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2624, scope))
		return []reconcile.Request{}
	}
	ctx := context.Background()

	httpRouteList := &gwapiv1b1.HTTPRouteList{}
	if err := apiReconciler.client.List(ctx, httpRouteList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(httprouteScopeIndex, utils.NamespacedName(scope).String()),
	}); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2625, utils.NamespacedName(scope).String()))
		return []reconcile.Request{}
	}

	if len(httpRouteList.Items) == 0 {
		loggers.LoggerAPKOperator.Debugf("HTTPRoutes for scope not found: %s", utils.NamespacedName(scope).String())
	}
	requests := []reconcile.Request{}
	for _, httpRoute := range httpRouteList.Items {
		requests = append(requests, apiReconciler.getAPIForHTTPRoute(&httpRoute)...)
	}
	return requests
}

// getAPIsForBackend triggers the API controller reconcile method based on the changes detected
// in backend resources.
func (apiReconciler *APIReconciler) getAPIsForBackend(obj k8client.Object) []reconcile.Request {
	ctx := context.Background()
	backend, ok := obj.(*dpv1alpha1.Backend)
	if !ok {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2624, backend))
		return []reconcile.Request{}
	}

	httpRouteList := &gwapiv1b1.HTTPRouteList{}
	if err := apiReconciler.client.List(ctx, httpRouteList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(backendHTTPRouteIndex, utils.NamespacedName(backend).String()),
	}); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2625, utils.NamespacedName(backend).String()))
		return []reconcile.Request{}
	}

	if len(httpRouteList.Items) == 0 {
		loggers.LoggerAPKOperator.Debugf("HTTPRoutes for Backend not found: %s", utils.NamespacedName(backend).String())
	}

	requests := []reconcile.Request{}
	for _, httpRoute := range httpRouteList.Items {
		requests = append(requests, apiReconciler.getAPIForHTTPRoute(&httpRoute)...)
	}

	// Create API reconcile events when Backend reffered from InterceptorService
	interceptorServiceList := &dpv1alpha1.InterceptorServiceList{}
	if err := apiReconciler.client.List(ctx, interceptorServiceList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(backendInterceptorServiceIndex, utils.NamespacedName(backend).String()),
	}); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2649, utils.NamespacedName(backend).String()))
		return []reconcile.Request{}
	}

	if len(interceptorServiceList.Items) == 0 {
		loggers.LoggerAPKOperator.Debugf("InterceptorService for Backend not found: %s", utils.NamespacedName(backend).String())
	}

	for _, interceptorService := range interceptorServiceList.Items {
		requests = append(requests, apiReconciler.getAPIsForInterceptorService(&interceptorService)...)
	}

	return requests
}

// getAPIsForGateway triggers the API controller reconcile method based on the changes detected
// in gateway resources.
func (apiReconciler *APIReconciler) getAPIsForGateway(obj k8client.Object) []reconcile.Request {
	ctx := context.Background()
	gateway, ok := obj.(*gwapiv1b1.Gateway)
	if !ok {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2624, gateway))
		return []reconcile.Request{}
	}

	httpRouteList := &gwapiv1b1.HTTPRouteList{}
	if err := apiReconciler.client.List(ctx, httpRouteList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(gatewayHTTPRouteIndex, utils.NamespacedName(gateway).String()),
	}); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2625, utils.NamespacedName(gateway).String()))
		return []reconcile.Request{}
	}

	if len(httpRouteList.Items) == 0 {
		loggers.LoggerAPKOperator.Debugf("HTTPRoutes for Gateway not found: %s", utils.NamespacedName(gateway).String())
		return []reconcile.Request{}
	}

	requests := []reconcile.Request{}
	for _, httpRoute := range httpRouteList.Items {
		requests = append(requests, apiReconciler.getAPIForHTTPRoute(&httpRoute)...)
	}
	return requests
}

// addIndexes adds indexing on API, for
//   - prodution and sandbox HTTPRoutes
//     referenced in API objects via `.spec.prodHTTPRouteRef` and `.spec.sandHTTPRouteRef`
//     This helps to find APIs that are affected by a HTTPRoute CRUD operation.
//   - authentications
//     authentication schemes related to httproutes
//     This helps to find authentication schemes binded to HTTPRoute.
//   - apiPolicies
//     apiPolicy schemes related to httproutes
//     This helps to find apiPolicy schemes binded to HTTPRoute.
func addIndexes(ctx context.Context, mgr manager.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha1.API{}, httpRouteAPIIndex,
		func(rawObj k8client.Object) []string {
			api := rawObj.(*dpv1alpha1.API)
			var httpRoutes []string
			if len(api.Spec.Production) > 0 {
				for _, ref := range api.Spec.Production[0].HTTPRouteRefs {
					if ref != "" {
						httpRoutes = append(httpRoutes,
							types.NamespacedName{
								Namespace: api.Namespace,
								Name:      ref,
							}.String())
					}
				}
			}
			if len(api.Spec.Sandbox) > 0 {
				for _, ref := range api.Spec.Sandbox[0].HTTPRouteRefs {
					if ref != "" {
						httpRoutes = append(httpRoutes,
							types.NamespacedName{
								Namespace: api.Namespace,
								Name:      ref,
							}.String())
					}
				}
			}
			return httpRoutes
		}); err != nil {
		return err
	}

	if err := mgr.GetFieldIndexer().IndexField(ctx, &gwapiv1b1.HTTPRoute{}, httprouteScopeIndex,
		func(rawObj k8client.Object) []string {
			httpRoute := rawObj.(*gwapiv1b1.HTTPRoute)
			var scopes []string
			for _, rule := range httpRoute.Spec.Rules {
				for _, filter := range rule.Filters {
					if filter.Type == gwapiv1b1.HTTPRouteFilterExtensionRef {
						if filter.ExtensionRef != nil && filter.ExtensionRef.Kind == constants.KindScope {
							scopes = append(scopes, types.NamespacedName{
								Namespace: httpRoute.Namespace,
								Name:      string(filter.ExtensionRef.Name),
							}.String())
						}
					}
				}
			}
			return scopes
		}); err != nil {
		return err
	}

	// Backend to HTTPRoute indexer
	if err := mgr.GetFieldIndexer().IndexField(ctx, &gwapiv1b1.HTTPRoute{}, backendHTTPRouteIndex,
		func(rawObj k8client.Object) []string {
			httpRoute := rawObj.(*gwapiv1b1.HTTPRoute)
			var backends []string
			for _, rule := range httpRoute.Spec.Rules {
				for _, backendRef := range rule.BackendRefs {
					if backendRef.Kind != nil && *backendRef.Kind == constants.KindBackend {
						backends = append(backends, types.NamespacedName{
							Namespace: utils.GetNamespace(backendRef.Namespace,
								httpRoute.ObjectMeta.Namespace),
							Name: string(backendRef.Name),
						}.String())
					}
				}
			}
			return backends
		}); err != nil {
		return err
	}

	// Gateway to HTTPRoute indexer
	if err := mgr.GetFieldIndexer().IndexField(ctx, &gwapiv1b1.HTTPRoute{}, gatewayHTTPRouteIndex,
		func(rawObj k8client.Object) []string {
			httpRoute := rawObj.(*gwapiv1b1.HTTPRoute)
			var gateways []string
			for _, parentRef := range httpRoute.Spec.ParentRefs {
				gateways = append(gateways, types.NamespacedName{
					Namespace: utils.GetNamespace(parentRef.Namespace,
						httpRoute.Namespace),
					Name: string(parentRef.Name),
				}.String())
			}
			return gateways
		}); err != nil {
		return err
	}

	// ConfigMap to Backend indexer
	if err := mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha1.Backend{}, configMapBackend,
		func(rawObj k8client.Object) []string {
			backend := rawObj.(*dpv1alpha1.Backend)
			var configMaps []string
			if backend.Spec.TLS != nil && backend.Spec.TLS.ConfigMapRef != nil && len(backend.Spec.TLS.ConfigMapRef.Name) > 0 {
				configMaps = append(configMaps,
					types.NamespacedName{
						Name:      string(backend.Spec.TLS.ConfigMapRef.Name),
						Namespace: backend.Namespace,
					}.String())
			}
			return configMaps
		}); err != nil {
		return err
	}

	// Secret to Backend indexer
	if err := mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha1.Backend{}, secretBackend,
		func(rawObj k8client.Object) []string {
			backend := rawObj.(*dpv1alpha1.Backend)
			var secrets []string
			if backend.Spec.TLS != nil && backend.Spec.TLS.SecretRef != nil && len(backend.Spec.TLS.SecretRef.Name) > 0 {
				secrets = append(secrets,
					types.NamespacedName{
						Name:      string(backend.Spec.TLS.SecretRef.Name),
						Namespace: backend.Namespace,
					}.String())
			}
			if backend.Spec.Security != nil {
				if backend.Spec.Security.Basic != nil {
					secrets = append(secrets,
						types.NamespacedName{
							Name:      string(backend.Spec.Security.Basic.SecretRef.Name),
							Namespace: backend.Namespace,
						}.String())
				}
			}
			return secrets
		}); err != nil {
		return err
	}

	// authentication to API indexer
	if err := mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha1.Authentication{}, apiAuthenticationIndex,
		func(rawObj k8client.Object) []string {
			authentication := rawObj.(*dpv1alpha1.Authentication)
			var apis []string
			if authentication.Spec.TargetRef.Kind == constants.KindAPI {
				apis = append(apis,
					types.NamespacedName{
						Namespace: utils.GetNamespace(
							(*gwapiv1b1.Namespace)(authentication.Spec.TargetRef.Namespace),
							authentication.Namespace),
						Name: string(authentication.Spec.TargetRef.Name),
					}.String())
			}
			return apis
		}); err != nil {
		return err
	}

	// Till the below is httproute rule name and targetref sectionname is supported,
	// https://gateway-api.sigs.k8s.io/geps/gep-713/?h=multiple+targetrefs#apply-policies-to-sections-of-a-resource-future-extension
	// we will use a temporary kindName called Resource for policy attachments
	// TODO(amali) Fix after the official support is available
	if err := mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha1.Authentication{}, apiAuthenticationResourceIndex,
		func(rawObj k8client.Object) []string {
			authentication := rawObj.(*dpv1alpha1.Authentication)
			var apis []string
			if authentication.Spec.TargetRef.Kind == constants.KindResource {
				apis = append(apis,
					types.NamespacedName{
						Namespace: utils.GetNamespace(
							(*gwapiv1b1.Namespace)(authentication.Spec.TargetRef.Namespace),
							authentication.Namespace),
						Name: string(authentication.Spec.TargetRef.Name),
					}.String())
			}
			return apis
		}); err != nil {
		return err
	}

	// ratelimite policy to API indexer
	if err := mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha1.RateLimitPolicy{}, apiRateLimitIndex,
		func(rawObj k8client.Object) []string {
			ratelimitPolicy := rawObj.(*dpv1alpha1.RateLimitPolicy)
			var apis []string
			if ratelimitPolicy.Spec.TargetRef.Kind == constants.KindAPI {
				apis = append(apis,
					types.NamespacedName{
						Namespace: utils.GetNamespace(
							(*gwapiv1b1.Namespace)(ratelimitPolicy.Spec.TargetRef.Namespace),
							ratelimitPolicy.Namespace),
						Name: string(ratelimitPolicy.Spec.TargetRef.Name),
					}.String())
			}
			return apis
		}); err != nil {
		return err
	}

	// Till the below is httproute rule name and targetref sectionname is supported,
	// https://gateway-api.sigs.k8s.io/geps/gep-713/?h=multiple+targetrefs#apply-policies-to-sections-of-a-resource-future-extension
	// we will use a temporary kindName called Resource for policy attachments
	// TODO(amali) Fix after the official support is available
	if err := mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha1.RateLimitPolicy{}, apiRateLimitResourceIndex,
		func(rawObj k8client.Object) []string {
			ratelimitPolicy := rawObj.(*dpv1alpha1.RateLimitPolicy)
			var apis []string
			if ratelimitPolicy.Spec.TargetRef.Kind == constants.KindResource {
				apis = append(apis,
					types.NamespacedName{
						Namespace: utils.GetNamespace(
							(*gwapiv1b1.Namespace)(ratelimitPolicy.Spec.TargetRef.Namespace),
							ratelimitPolicy.Namespace),
						Name: string(ratelimitPolicy.Spec.TargetRef.Name),
					}.String())
			}
			return apis
		}); err != nil {
		return err
	}

	// backend to InterceptorService indexer
	if err := mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha1.InterceptorService{}, backendInterceptorServiceIndex,
		func(rawObj k8client.Object) []string {
			interceptorService := rawObj.(*dpv1alpha1.InterceptorService)
			var backends []string
			backends = append(backends,
				types.NamespacedName{
					Namespace: utils.GetNamespace(
						(*gwapiv1b1.Namespace)(&interceptorService.Spec.BackendRef.Namespace), interceptorService.Namespace),
					Name: string(interceptorService.Spec.BackendRef.Name),
				}.String())
			return backends
		}); err != nil {
		return err
	}

	// interceptorService to APIPolicy indexer
	if err := mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha1.APIPolicy{}, interceptorServiceAPIPolicyIndex,
		func(rawObj k8client.Object) []string {
			apiPolicy := rawObj.(*dpv1alpha1.APIPolicy)
			var interceptorServices []string
			if apiPolicy.Spec.Default != nil && len(apiPolicy.Spec.Default.RequestInterceptors) > 0 {
				interceptorServices = append(interceptorServices,
					types.NamespacedName{
						Namespace: utils.GetNamespace(
							(*gwapiv1b1.Namespace)(&apiPolicy.Spec.Default.RequestInterceptors[0].Namespace), apiPolicy.Namespace),
						Name: string(apiPolicy.Spec.Default.RequestInterceptors[0].Name),
					}.String())
			}
			if apiPolicy.Spec.Override != nil && len(apiPolicy.Spec.Override.RequestInterceptors) > 0 {
				interceptorServices = append(interceptorServices,
					types.NamespacedName{
						Namespace: utils.GetNamespace(
							(*gwapiv1b1.Namespace)(&apiPolicy.Spec.Override.RequestInterceptors[0].Namespace), apiPolicy.Namespace),
						Name: string(apiPolicy.Spec.Override.RequestInterceptors[0].Name),
					}.String())
			}
			if apiPolicy.Spec.Default != nil && len(apiPolicy.Spec.Default.ResponseInterceptors) > 0 {
				interceptorServices = append(interceptorServices,
					types.NamespacedName{
						Namespace: utils.GetNamespace(
							(*gwapiv1b1.Namespace)(&apiPolicy.Spec.Default.ResponseInterceptors[0].Namespace), apiPolicy.Namespace),
						Name: string(apiPolicy.Spec.Default.ResponseInterceptors[0].Name),
					}.String())
			}
			if apiPolicy.Spec.Override != nil && len(apiPolicy.Spec.Override.ResponseInterceptors) > 0 {
				interceptorServices = append(interceptorServices,
					types.NamespacedName{
						Namespace: utils.GetNamespace(
							(*gwapiv1b1.Namespace)(&apiPolicy.Spec.Override.ResponseInterceptors[0].Namespace), apiPolicy.Namespace),
						Name: string(apiPolicy.Spec.Override.ResponseInterceptors[0].Name),
					}.String())
			}
			return interceptorServices
		}); err != nil {
		return err
	}

	// httpRoute to APIPolicy indexer
	if err := mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha1.APIPolicy{}, apiAPIPolicyIndex,
		func(rawObj k8client.Object) []string {
			apiPolicy := rawObj.(*dpv1alpha1.APIPolicy)
			var apis []string
			if apiPolicy.Spec.TargetRef.Kind == constants.KindAPI {
				apis = append(apis,
					types.NamespacedName{
						Namespace: utils.GetNamespace(
							(*gwapiv1b1.Namespace)(apiPolicy.Spec.TargetRef.Namespace), apiPolicy.Namespace),
						Name: string(apiPolicy.Spec.TargetRef.Name),
					}.String())
			}
			return apis
		}); err != nil {
		return err
	}

	// api to APIPolicy in resource level indexer
	// Till the below is httproute rule name and targetref sectionname is supported,
	// https://gateway-api.sigs.k8s.io/geps/gep-713/?h=multiple+targetrefs#apply-policies-to-sections-of-a-resource-future-extension
	// we will use a temporary kindName called Resource for policy attachments
	// TODO(amali) Fix after the official support is available
	err := mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha1.APIPolicy{}, apiAPIPolicyResourceIndex,
		func(rawObj k8client.Object) []string {
			apiPolicy := rawObj.(*dpv1alpha1.APIPolicy)
			var apis []string
			if apiPolicy.Spec.TargetRef.Kind == constants.KindResource {
				apis = append(apis,
					types.NamespacedName{
						Namespace: utils.GetNamespace(
							(*gwapiv1b1.Namespace)(apiPolicy.Spec.TargetRef.Namespace),
							apiPolicy.Namespace),
						Name: string(apiPolicy.Spec.TargetRef.Name),
					}.String())
			}
			return apis
		})
	return err
}

// handleStatus updates the API CR update
func (apiReconciler *APIReconciler) handleStatus() {
	for successEvent := range *apiReconciler.successChannel {
		accept := false
		message := ""
		event := ""

		switch successEvent.State {
		case constants.Create:
			accept = true
			message = "API is deployed to the gateway."
		case constants.Update:
			accept = true
			message = fmt.Sprintf("API update is deployed to the gateway. %v were Updated", successEvent.Events)
		}
		timeNow := metav1.Now()
		event = fmt.Sprintf("[%s] %s", timeNow.String(), message)

		apiReconciler.statusUpdater.Send(status.Update{
			NamespacedName: successEvent.APINamespacedName,
			Resource:       new(dpv1alpha1.API),
			UpdateStatus: func(obj k8client.Object) k8client.Object {
				h, ok := obj.(*dpv1alpha1.API)
				if !ok {
					loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2626, obj))
				}
				hCopy := h.DeepCopy()
				hCopy.APIDeploymentEvent.Status = successEvent.State
				hCopy.APIDeploymentEvent.Accepted = accept
				hCopy.APIDeploymentEvent.Message = message
				hCopy.APIDeploymentEvent.Events = append(hCopy.APIDeploymentEvent.Events, event)
				hCopy.APIDeploymentEvent.TransitionTime = &timeNow
				return hCopy
			},
		})
	}
}
