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
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/wso2/apk/adapter/config"
	"github.com/wso2/apk/adapter/internal/controlplane"
	"github.com/wso2/apk/adapter/internal/discovery/xds"
	"github.com/wso2/apk/adapter/internal/discovery/xds/common"
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
	ctrl "sigs.k8s.io/controller-runtime"
	k8client "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
	"sigs.k8s.io/gateway-api/apis/v1alpha2"

	dpv1alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha1"
	dpv1alpha2 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha2"
	dpv1alpha3 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha3"
	dpv1alpha4 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha4"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	httpRouteAPIIndex = "httpRouteToAPIIndex"
	gqlRouteAPIIndex  = "gqlRouteToAPIIndex"
	grpcRouteAPIIndex = "grpcRouteToAPIIndex"
	// apiAuthenticationIndex Index for API level authentications
	apiAuthenticationIndex = "apiToAuthenticationIndex"
	// apiAuthenticationResourceIndex Index for resource level authentications
	apiAuthenticationResourceIndex = "apiToAuthenticationResourceIndex"
	// apiRateLimitIndex Index for API level ratelimits
	apiRateLimitIndex = "apiToRateLimitIndex"
	// apiRateLimitResourceIndex Index for resource level ratelimits
	apiRateLimitResourceIndex = "apiToRateLimitResourceIndex"
	// gatewayHTTPRouteIndex Index for gateway httproutes
	gatewayHTTPRouteIndex = "gatewayToHTTPRouteIndex"
	// gatewayGRPCRouteIndex Index for gateway grpcroutes
	gatewayGRPCRouteIndex = "gatewayToGRPCRouteIndex"
	// apiAPIPolicyIndex Index for API level apipolicies
	apiAPIPolicyIndex = "apiToAPIPolicyIndex"
	// apiAPIPolicyResourceIndex Index for resource level apipolicies
	apiAPIPolicyResourceIndex        = "apiToAPIPolicyResourceIndex"
	serviceHTTPRouteIndex            = "serviceToHTTPRouteIndex"
	httprouteScopeIndex              = "httprouteToScopeIndex"
	grpcRouteScopeIndex              = "grpcRouteToScopeIndex"
	gqlRouteScopeIndex               = "gqlRouteToScopeIndex"
	configMapBackend                 = "configMapToBackend"
	configMapAPIDefinition           = "configMapToAPIDefinition"
	secretBackend                    = "secretToBackend"
	configMapAuthentication          = "configMapToAuthentication"
	secretAuthentication             = "secretToAuthentication"
	backendHTTPRouteIndex            = "backendToHTTPRouteIndex"
	backendGQLRouteIndex             = "backendToGQLRouteIndex"
	backendGRPCRouteIndex            = "backendToGRPCRouteIndex"
	interceptorServiceAPIPolicyIndex = "interceptorServiceToAPIPolicyIndex"
	backendInterceptorServiceIndex   = "backendToInterceptorServiceIndex"
	backendJWTAPIPolicyIndex         = "backendJWTToAPIPolicyIndex"
	aiRatelimitPolicyToBackendIndex  = "aiRatelimitToPolicyToBackendIndex"
	subscriptionToAPIIndex           = "subscriptionToAPIIndex"
	apiToSubscriptionIndex           = "apiToSubscriptionIndex"
	aiProviderAPIPolicyIndex         = "aiProviderToAPIPolicyIndex"
)

var (
	applyAllAPIsOnce sync.Once
)

// APIReconciler reconciles a API object
type APIReconciler struct {
	client                k8client.Client
	ods                   *synchronizer.OperatorDataStore
	ch                    *chan *synchronizer.APIEvent
	successChannel        *chan synchronizer.SuccessEvent
	statusUpdater         *status.UpdateHandler
	mgr                   manager.Manager
	apiPropagationEnabled bool
}

// NewAPIController creates a new API controller instance. API Controllers watches for dpv1alpha3.API and gwapiv1.HTTPRoute.
func NewAPIController(mgr manager.Manager, operatorDataStore *synchronizer.OperatorDataStore, statusUpdater *status.UpdateHandler,
	ch *chan *synchronizer.APIEvent, successChannel *chan synchronizer.SuccessEvent) error {
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
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2619, logging.BLOCKER, "Error applying startup APIs: %v", err.Error()))
		return err
	}

	conf := config.ReadConfigs()
	apiReconciler.apiPropagationEnabled = conf.Adapter.ControlPlane.EnableAPIPropagation

	predicateAPI := []predicate.TypedPredicate[*dpv1alpha3.API]{predicate.NewTypedPredicateFuncs[*dpv1alpha3.API](utils.FilterAPIByNamespaces(conf.Adapter.Operator.Namespaces))}
	if err := c.Watch(source.Kind(mgr.GetCache(), &dpv1alpha3.API{}, &handler.TypedEnqueueRequestForObject[*dpv1alpha3.API]{}, predicateAPI...)); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2611, logging.BLOCKER, "Error watching API resources: %v", err))
		return err
	}
	if err := addIndexes(ctx, mgr); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2612, logging.BLOCKER, "Error adding indexes: %v", err))
		return err
	}

	predicateHTTPRoute := []predicate.TypedPredicate[*gwapiv1.HTTPRoute]{predicate.NewTypedPredicateFuncs[*gwapiv1.HTTPRoute](utils.FilterHTTPRouteByNamespaces(conf.Adapter.Operator.Namespaces))}
	if err := c.Watch(source.Kind(mgr.GetCache(), &gwapiv1.HTTPRoute{}, handler.TypedEnqueueRequestsFromMapFunc[*gwapiv1.HTTPRoute](apiReconciler.populateAPIReconcileRequestsForHTTPRoute),
		predicateHTTPRoute...)); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2613, logging.BLOCKER, "Error watching HTTPRoute resources: %v", err))
		return err
	}

	predicateGQLRoute := []predicate.TypedPredicate[*dpv1alpha2.GQLRoute]{predicate.NewTypedPredicateFuncs[*dpv1alpha2.GQLRoute](utils.FilterGQLRouteByNamespaces(conf.Adapter.Operator.Namespaces))}
	if err := c.Watch(source.Kind(mgr.GetCache(), &dpv1alpha2.GQLRoute{}, handler.TypedEnqueueRequestsFromMapFunc(apiReconciler.populateAPIReconcileRequestsForGQLRoute),
		predicateGQLRoute...)); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2667, logging.BLOCKER, "Error watching GQLRoute resources: %v", err))
		return err
	}

	predicateGRPCRoute := []predicate.TypedPredicate[*gwapiv1.GRPCRoute]{predicate.NewTypedPredicateFuncs[*gwapiv1.GRPCRoute](utils.FilterGRPCRouteByNamespaces(conf.Adapter.Operator.Namespaces))}
	if err := c.Watch(source.Kind(mgr.GetCache(), &gwapiv1.GRPCRoute{}, handler.TypedEnqueueRequestsFromMapFunc(apiReconciler.populateAPIReconcileRequestsForGRPCRoute),
		predicateGRPCRoute...)); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2668, logging.BLOCKER, "Error watching GRPCRoute resources: %v", err))
		return err
	}

	predicateGateway := []predicate.TypedPredicate[*gwapiv1.Gateway]{predicate.NewTypedPredicateFuncs[*gwapiv1.Gateway](utils.FilterGatewayByNamespaces(conf.Adapter.Operator.Namespaces))}
	if err := c.Watch(source.Kind(mgr.GetCache(), &gwapiv1.Gateway{}, handler.TypedEnqueueRequestsFromMapFunc(apiReconciler.getAPIsForGateway),
		predicateGateway...)); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2611, logging.BLOCKER, "Error watching Gateway resources: %v", err))
		return err
	}

	predicateBackend := []predicate.TypedPredicate[*dpv1alpha2.Backend]{predicate.NewTypedPredicateFuncs[*dpv1alpha2.Backend](utils.FilterBackendByNamespaces(conf.Adapter.Operator.Namespaces))}
	if err := c.Watch(source.Kind(mgr.GetCache(), &dpv1alpha2.Backend{}, handler.TypedEnqueueRequestsFromMapFunc(apiReconciler.populateAPIReconcileRequestsForBackend),
		predicateBackend...)); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2615, logging.BLOCKER, "Error watching Backend resources: %v", err))
		return err
	}

	predicateAuthentication := []predicate.TypedPredicate[*dpv1alpha2.Authentication]{predicate.NewTypedPredicateFuncs[*dpv1alpha2.Authentication](utils.FilterAuthenticationByNamespaces(conf.Adapter.Operator.Namespaces))}
	if err := c.Watch(source.Kind(mgr.GetCache(), &dpv1alpha2.Authentication{}, handler.TypedEnqueueRequestsFromMapFunc(apiReconciler.populateAPIReconcileRequestsForAuthentication),
		predicateAuthentication...)); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2616, logging.BLOCKER, "Error watching Authentication resources: %v", err))
		return err
	}

	predicateInterceptorService := []predicate.TypedPredicate[*dpv1alpha1.InterceptorService]{predicate.NewTypedPredicateFuncs[*dpv1alpha1.InterceptorService](utils.FilterInterceptorServiceByNamespaces(conf.Adapter.Operator.Namespaces))}
	if err := c.Watch(source.Kind(mgr.GetCache(), &dpv1alpha1.InterceptorService{}, handler.TypedEnqueueRequestsFromMapFunc(apiReconciler.populateAPIReconcileRequestsForInterceptorService),
		predicateInterceptorService...)); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2640, logging.BLOCKER, "Error watching InterceptorService resources: %v", err))
		return err
	}

	predicateBackendJWT := []predicate.TypedPredicate[*dpv1alpha1.BackendJWT]{predicate.NewTypedPredicateFuncs[*dpv1alpha1.BackendJWT](utils.FilterBackendJWTByNamespaces(conf.Adapter.Operator.Namespaces))}
	if err := c.Watch(source.Kind(mgr.GetCache(), &dpv1alpha1.BackendJWT{}, handler.TypedEnqueueRequestsFromMapFunc(apiReconciler.populateAPIReconcileRequestsForBackendJWT),
		predicateBackendJWT...)); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2661, logging.BLOCKER, "Error watching BackendJWT resources: %v", err))
		return err
	}

	predicateAPIPolicy := []predicate.TypedPredicate[*dpv1alpha4.APIPolicy]{predicate.NewTypedPredicateFuncs[*dpv1alpha4.APIPolicy](utils.FilterAPIPolicyByNamespaces(conf.Adapter.Operator.Namespaces))}
	if err := c.Watch(source.Kind(mgr.GetCache(), &dpv1alpha4.APIPolicy{}, handler.TypedEnqueueRequestsFromMapFunc(apiReconciler.populateAPIReconcileRequestsForAPIPolicy),
		predicateAPIPolicy...)); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2617, logging.BLOCKER, "Error watching APIPolicy resources: %v", err))
		return err
	}

	predicateRateLimitPolicy := []predicate.TypedPredicate[*dpv1alpha3.RateLimitPolicy]{predicate.NewTypedPredicateFuncs[*dpv1alpha3.RateLimitPolicy](utils.FilterRateLimitPolicyByNamespaces(conf.Adapter.Operator.Namespaces))}
	if err := c.Watch(source.Kind(mgr.GetCache(), &dpv1alpha3.RateLimitPolicy{}, handler.TypedEnqueueRequestsFromMapFunc(apiReconciler.populateAPIReconcileRequestsForRateLimitPolicy),
		predicateRateLimitPolicy...)); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2639, logging.BLOCKER, "Error watching Ratelimit resources: %v", err))
		return err
	}

	predicateScope := []predicate.TypedPredicate[*dpv1alpha1.Scope]{predicate.NewTypedPredicateFuncs[*dpv1alpha1.Scope](utils.FilterScopeByNamespaces(conf.Adapter.Operator.Namespaces))}
	if err := c.Watch(source.Kind(mgr.GetCache(), &dpv1alpha1.Scope{}, handler.TypedEnqueueRequestsFromMapFunc(apiReconciler.populateAPIReconcileRequestsForScope),
		predicateScope...)); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2618, logging.BLOCKER, "Error watching scope resources: %v", err))
		return err
	}

	predicateConfigMap := []predicate.TypedPredicate[*corev1.ConfigMap]{predicate.NewTypedPredicateFuncs[*corev1.ConfigMap](utils.FilterConfigMapByNamespaces(conf.Adapter.Operator.Namespaces))}
	if err := c.Watch(source.Kind(mgr.GetCache(), &corev1.ConfigMap{}, handler.TypedEnqueueRequestsFromMapFunc(apiReconciler.populateAPIReconcileRequestsForConfigMap),
		predicateConfigMap...)); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2644, logging.BLOCKER, "Error watching ConfigMap resources: %v", err))
		return err
	}

	predicateSecret := []predicate.TypedPredicate[*corev1.Secret]{predicate.NewTypedPredicateFuncs[*corev1.Secret](utils.FilterSecretByNamespaces(conf.Adapter.Operator.Namespaces))}
	if err := c.Watch(source.Kind(mgr.GetCache(), &corev1.Secret{}, handler.TypedEnqueueRequestsFromMapFunc(apiReconciler.populateAPIReconcileRequestsForSecret),
		predicateSecret...)); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2645, logging.BLOCKER, "Error watching Secret resources: %v", err))
		return err
	}

	predicateAIProvider := []predicate.TypedPredicate[*dpv1alpha4.AIProvider]{predicate.NewTypedPredicateFuncs[*dpv1alpha4.AIProvider](utils.FilterAIProviderByNamespaces(conf.Adapter.Operator.Namespaces))}
	if err := c.Watch(source.Kind(mgr.GetCache(), &dpv1alpha4.AIProvider{}, handler.TypedEnqueueRequestsFromMapFunc(apiReconciler.populateAPIReconcileRequestsForAIProvider),
		predicateAIProvider...)); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2615, logging.BLOCKER, "Error watching AIPolicy resources: %v", err))
		return err
	}

	predicateAIRatelimitPolicy := []predicate.TypedPredicate[*dpv1alpha3.AIRateLimitPolicy]{predicate.NewTypedPredicateFuncs[*dpv1alpha3.AIRateLimitPolicy](utils.FilterAIRatelimitPolicyByNamespaces(conf.Adapter.Operator.Namespaces))}
	if err := c.Watch(source.Kind(mgr.GetCache(), &dpv1alpha3.AIRateLimitPolicy{}, handler.TypedEnqueueRequestsFromMapFunc(apiReconciler.populateAPIReconcileRequestsForAIRatelimitPolicy),
		predicateAIRatelimitPolicy...)); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2645, logging.BLOCKER, "Error watching AIRatelimitPolicy resources: %v", err))
		return err
	}

	loggers.LoggerAPKOperator.Info("API Controller successfully started. Watching API Objects....")
	go apiReconciler.handleStatus()
	go apiReconciler.handleLabels(ctx)
	return nil
}

// +kubebuilder:rbac:groups=dp.wso2.com,resources=apis,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=dp.wso2.com,resources=apis/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=dp.wso2.com,resources=apis/finalizers,verbs=update
// +kubebuilder:rbac:groups=dp.wso2.com,resources=httproutes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=dp.wso2.com,resources=httproutes/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=dp.wso2.com,resources=httproutes/finalizers,verbs=update
// +kubebuilder:rbac:groups=dp.wso2.com,resources=gqlroutes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=dp.wso2.com,resources=gqlroutes/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=dp.wso2.com,resources=gqlroutes/finalizers,verbs=update
// +kubebuilder:rbac:groups=dp.wso2.com,resources=authentications,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=dp.wso2.com,resources=authentications/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=dp.wso2.com,resources=authentications/finalizers,verbs=update
// +kubebuilder:rbac:groups=dp.wso2.com,resources=grpcroutes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=dp.wso2.com,resources=grpcroutes/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=dp.wso2.com,resources=grpcroutes/finalizers,verbs=update
// +kubebuilder:rbac:groups=dp.wso2.com,resources=apipolicies,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=dp.wso2.com,resources=apipolicies/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=dp.wso2.com,resources=apipolicies/finalizers,verbs=update
// +kubebuilder:rbac:groups=dp.wso2.com,resources=scopes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=dp.wso2.com,resources=scopes/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=dp.wso2.com,resources=scopes/finalizers,verbs=update
// +kubebuilder:rbac:groups=dp.wso2.com,resources=ratelimitpolicies,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=dp.wso2.com,resources=ratelimitpolicies/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=dp.wso2.com,resources=ratelimitpolicies/finalizers,verbs=update
// +kubebuilder:rbac:groups=dp.wso2.com,resources=airatelimitpolicies,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=dp.wso2.com,resources=airatelimitpolicies/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=dp.wso2.com,resources=airatelimitpolicies/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (apiReconciler *APIReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	applyAllAPIsOnce.Do(apiReconciler.applyStartupAPIs)
	loggers.LoggerAPKOperator.Infof("Reconciling for API %s", req.NamespacedName.String())
	// Check whether the API CR exist, if not consider as a DELETE event.
	var apiCR dpv1alpha3.API
	if err := apiReconciler.client.Get(ctx, req.NamespacedName, &apiCR); err != nil {
		apiState, found := apiReconciler.ods.GetCachedAPI(req.NamespacedName)
		if found && k8error.IsNotFound(err) {
			if apiReconciler.apiPropagationEnabled && isAPIPropagatable(&apiState) {
				// Convert api state to api cp data
				loggers.LoggerAPKOperator.Info("Sending API deletion event to agent")
				apiCpData := apiReconciler.convertAPIStateToAPICp(ctx, apiState, "", nil, nil)
				apiCpData.Event = controlplane.EventTypeDelete
				controlplane.AddToEventQueue(apiCpData)
			}
			// The api doesn't exist in the api Cache, remove it
			apiReconciler.ods.DeleteCachedAPI(req.NamespacedName)
			loggers.LoggerAPKOperator.Infof("Delete event received for API : %s with API UUID : %v, hence deleted from API cache",
				req.NamespacedName.String(), string(apiCR.ObjectMeta.UID))
			*apiReconciler.ch <- &synchronizer.APIEvent{EventType: constants.Delete, Events: []synchronizer.APIState{apiState}}
			return ctrl.Result{}, nil
		}
		loggers.LoggerAPKOperator.Warnf("Api CR related to the reconcile request with key: %s returned error. Assuming API with API UUID : %v is already deleted, hence ignoring the error : %v",
			req.NamespacedName.String(), string(apiCR.ObjectMeta.UID), err)
		return ctrl.Result{}, nil
	}

	if apiState, err := apiReconciler.resolveAPIRefs(ctx, apiCR); err != nil {
		loggers.LoggerAPKOperator.Warnf("Error retrieving ref CRs for API in namespace : %s with API UUID : %v, %v",
			req.NamespacedName.String(), string(apiCR.ObjectMeta.UID), err)
		return ctrl.Result{}, nil
	} else if apiState != nil {
		loggers.LoggerAPKOperator.Infof("Ready to deploy CRs for API in namespace : %s with API UUID : %v, %v",
			req.NamespacedName.String(), string(apiCR.ObjectMeta.UID), err)
		*apiReconciler.ch <- apiState
	}
	return ctrl.Result{}, nil
}

// applyStartupAPIs applies the APIs which are already available in the cluster at the startup of the operator.
func (apiReconciler *APIReconciler) applyStartupAPIs() {
	ctx := context.Background()
	apisList, err := utils.RetrieveAPIList(apiReconciler.client)
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2605, logging.CRITICAL, "Unable to list APIs: %v", err))
		return
	}
	combinedapiEvent := &synchronizer.APIEvent{
		EventType: constants.Create,
		Events:    make([]synchronizer.APIState, 0),
	}
	for _, api := range apisList {
		if apiState, err := apiReconciler.resolveAPIRefs(ctx, api); err != nil {
			loggers.LoggerAPKOperator.Warnf("Error retrieving ref CRs for API : %s in namespace : %s with API UUID : %v, %v",
				api.Name, api.Namespace, string(api.ObjectMeta.UID), err)
		} else if apiState != nil {
			combinedapiEvent.Events = append(combinedapiEvent.Events, apiState.Events...)
		}
	}
	// Send all the API events to the channel
	if len(combinedapiEvent.Events) > 0 {
		*apiReconciler.ch <- combinedapiEvent
		loggers.LoggerAPKOperator.Info("Initial APIs were reconciled successfully")
	} else {
		loggers.LoggerAPKOperator.Warn("No startup APIs found")
	}
	xds.SetReady()
}

// resolveAPIRefs validates following references related to the API
// - HTTPRoutes
func (apiReconciler *APIReconciler) resolveAPIRefs(ctx context.Context, api dpv1alpha3.API) (*synchronizer.APIEvent, error) {
	var prodRouteRefs, sandRouteRefs []string
	if len(api.Spec.Production) > 0 {
		prodRouteRefs = api.Spec.Production[0].RouteRefs
	}
	if len(api.Spec.Sandbox) > 0 {
		sandRouteRefs = api.Spec.Sandbox[0].RouteRefs
	}

	apiState := &synchronizer.APIState{
		APIDefinition: &api,
	}
	var err error
	apiRef := utils.NamespacedName(&api)
	namespace := api.Namespace
	if apiState.Authentications, err = apiReconciler.getAuthenticationsForAPI(ctx, api); err != nil {
		return nil, fmt.Errorf("error while getting API level auth for API : %s in namespace : %s with API UUID : %v, %s",
			apiRef.String(), namespace, string(api.ObjectMeta.UID), err.Error())
	}
	if apiState.RateLimitPolicies, err = apiReconciler.getRatelimitPoliciesForAPI(ctx, api); err != nil {
		return nil, fmt.Errorf("error while getting API level ratelimit for API : %s in namespace : %s with API UUID : %v, %s",
			apiRef.String(), namespace, string(api.ObjectMeta.UID), err.Error())
	}
	if apiState.APIPolicies, err = apiReconciler.getAPIPoliciesForAPI(ctx, api); err != nil {
		return nil, fmt.Errorf("error while getting API level apipolicy for API : %s in namespace : %s with API UUID : %v, %s",
			apiRef.String(), namespace, string(api.ObjectMeta.UID), err.Error())
	}

	if apiState.ResourceAuthentications, err = apiReconciler.getAuthenticationsForResources(ctx, api); err != nil {
		return nil, fmt.Errorf("error while getting httproute resource auth : %s in namespace : %s with API UUID : %v, %s",
			apiRef.String(), namespace, string(api.ObjectMeta.UID), err.Error())
	}
	if apiState.ResourceRateLimitPolicies, err = apiReconciler.getRatelimitPoliciesForResources(ctx, api); err != nil {
		return nil, fmt.Errorf("error while getting httproute resource ratelimit : %s in namespace : %s with API UUID : %v, %s",
			apiRef.String(), namespace, string(api.ObjectMeta.UID), err.Error())
	}
	if apiState.ResourceAPIPolicies, err = apiReconciler.getAPIPoliciesForResources(ctx, api); err != nil {
		return nil, fmt.Errorf("error while getting httproute resource apipolicy %s in namespace : %s with API UUID : %v, %s",
			apiRef.String(), namespace, string(api.ObjectMeta.UID), err.Error())
	}
	if apiState.InterceptorServiceMapping, apiState.BackendJWTMapping, apiState.SubscriptionValidation, apiState.AIProvider, apiState.ResolvedModelBasedRoundRobin, err =
		apiReconciler.getAPIPolicyChildrenRefs(ctx, apiState.APIPolicies, apiState.ResourceAPIPolicies, api); err != nil {
		return nil, fmt.Errorf("error while getting referenced policies in apipolicy %s in namespace : %s with API UUID : %v, %s",
			apiRef.String(), namespace, string(api.ObjectMeta.UID), err.Error())
	}
	if api.Spec.DefinitionFileRef != "" {
		if apiState.APIDefinitionFile, err = apiReconciler.getAPIDefinitionForAPI(ctx, api.Spec.DefinitionFileRef, namespace, api); err != nil {
			return nil, fmt.Errorf("error while getting api definition file of api %s in namespace : %s with API UUID : %v, %s",
				apiRef.String(), namespace, string(api.ObjectMeta.UID), err.Error())
		} else if apiState.APIDefinitionFile == nil && apiState.APIDefinition.Spec.APIType == "GraphQL" {
			return nil, fmt.Errorf("error while getting api definition file of api %s in namespace : %s with API UUID : %v, %s",
				apiRef.String(), namespace, string(api.ObjectMeta.UID), "api definition file not found")
		}
	}
	if len(apiState.Authentications) > 0 {
		if apiState.MutualSSL, err = apiReconciler.resolveAuthentications(ctx, apiState.Authentications); err != nil {
			return nil, fmt.Errorf("error while resolving authentication %v in namespace: %s was not found. %s",
				apiState.Authentications, namespace, err.Error())
		}
	}
	var prodAirl *dpv1alpha3.AIRateLimitPolicy
	if len(prodRouteRefs) > 0 && apiState.APIDefinition.Spec.APIType == "REST" {
		apiState.ProdHTTPRoute = &synchronizer.HTTPRouteState{}
		if apiState.ProdHTTPRoute, prodAirl, err = apiReconciler.resolveHTTPRouteRefs(ctx, apiState.ProdHTTPRoute, prodRouteRefs,
			namespace, apiState.InterceptorServiceMapping, api); err != nil {
			return nil, fmt.Errorf("error while resolving production httpRouteref %s in namespace :%s has not found. %s",
				prodRouteRefs, namespace, err.Error())
		}
		if !apiReconciler.ods.IsGatewayAvailable(types.NamespacedName{
			Name: string(apiState.ProdHTTPRoute.HTTPRouteCombined.Spec.ParentRefs[0].Name),
			Namespace: utils.GetNamespace(apiState.ProdHTTPRoute.HTTPRouteCombined.Spec.ParentRefs[0].Namespace,
				apiState.ProdHTTPRoute.HTTPRouteCombined.Namespace),
		}) {
			return nil, fmt.Errorf("no gateway available for httpRouteref %s in namespace :%s has not found",
				prodRouteRefs, namespace)
		}
	}
	apiState.ProdAIRL = prodAirl
	var sandAirl *dpv1alpha3.AIRateLimitPolicy
	if len(sandRouteRefs) > 0 && apiState.APIDefinition.Spec.APIType == "REST" {
		apiState.SandHTTPRoute = &synchronizer.HTTPRouteState{}
		if apiState.SandHTTPRoute, sandAirl, err = apiReconciler.resolveHTTPRouteRefs(ctx, apiState.SandHTTPRoute, sandRouteRefs,
			namespace, apiState.InterceptorServiceMapping, api); err != nil {
			return nil, fmt.Errorf("error while resolving sandbox httpRouteref %s in namespace :%s has not found. %s",
				sandRouteRefs, namespace, err.Error())
		}
		if !apiReconciler.ods.IsGatewayAvailable(types.NamespacedName{
			Name: string(apiState.SandHTTPRoute.HTTPRouteCombined.Spec.ParentRefs[0].Name),
			Namespace: utils.GetNamespace(apiState.SandHTTPRoute.HTTPRouteCombined.Spec.ParentRefs[0].Namespace,
				apiState.SandHTTPRoute.HTTPRouteCombined.Namespace),
		}) {
			return nil, fmt.Errorf("no gateway available for httpRouteref %s in namespace :%s has not found",
				sandRouteRefs, namespace)
		}
	}
	apiState.SandAIRL = sandAirl
	// handle gql apis
	if len(prodRouteRefs) > 0 && apiState.APIDefinition.Spec.APIType == "GraphQL" {
		if apiState.ProdGQLRoute, err = apiReconciler.resolveGQLRouteRefs(ctx, prodRouteRefs, namespace,
			api); err != nil {
			return nil, fmt.Errorf("error while resolving production gqlRouteref %s in namespace :%s has not found. %s",
				prodRouteRefs, namespace, err.Error())
		}
	}
	if len(sandRouteRefs) > 0 && apiState.APIDefinition.Spec.APIType == "GraphQL" {
		if apiState.SandGQLRoute, err = apiReconciler.resolveGQLRouteRefs(ctx, sandRouteRefs, namespace,
			api); err != nil {
			return nil, fmt.Errorf("error while resolving sandbox gqlRouteref %s in namespace :%s has not found. %s",
				sandRouteRefs, namespace, err.Error())
		}
	}

	// handle gRPC APIs
	if len(prodRouteRefs) > 0 && apiState.APIDefinition.Spec.APIType == constants.GRPC {
		if apiState.ProdGRPCRoute, err = apiReconciler.resolveGRPCRouteRefs(ctx, prodRouteRefs,
			namespace, apiState.InterceptorServiceMapping, api); err != nil {
			return nil, fmt.Errorf("error while resolving production grpcRouteref %s in namespace :%s was not found. %s",
				prodRouteRefs, namespace, err.Error())
		}
		if !apiReconciler.ods.IsGatewayAvailable(types.NamespacedName{
			Name: string(apiState.ProdGRPCRoute.GRPCRouteCombined.Spec.ParentRefs[0].Name),
			Namespace: utils.GetNamespace(apiState.ProdGRPCRoute.GRPCRouteCombined.Spec.ParentRefs[0].Namespace,
				apiState.ProdGRPCRoute.GRPCRouteCombined.Namespace),
		}) {
			return nil, fmt.Errorf("no gateway available for grpcpRouteref %s in namespace :%s was not found",
				prodRouteRefs, namespace)
		}
	}

	if len(sandRouteRefs) > 0 && apiState.APIDefinition.Spec.APIType == constants.GRPC {
		if apiState.SandGRPCRoute, err = apiReconciler.resolveGRPCRouteRefs(ctx, sandRouteRefs,
			namespace, apiState.InterceptorServiceMapping, api); err != nil {
			return nil, fmt.Errorf("error while resolving sandbox grpcRouteref %s in namespace :%s was not found. %s",
				sandRouteRefs, namespace, err.Error())
		}
		if !apiReconciler.ods.IsGatewayAvailable(types.NamespacedName{
			Name: string(apiState.SandGRPCRoute.GRPCRouteCombined.Spec.ParentRefs[0].Name),
			Namespace: utils.GetNamespace(apiState.SandGRPCRoute.GRPCRouteCombined.Spec.ParentRefs[0].Namespace,
				apiState.SandGRPCRoute.GRPCRouteCombined.Namespace),
		}) {
			return nil, fmt.Errorf("no gateway available for grpcRouteref %s in namespace :%s was not found",
				sandRouteRefs, namespace)
		}
	}

	// Validate resource level extension refs resolved
	extRefValErr := apiReconciler.validateRouteExtRefs(apiState)
	if extRefValErr != nil {
		return nil, extRefValErr
	}

	loggers.LoggerAPKOperator.Debugf("Child references are retrieved successfully for API CR %s", apiRef.String())
	storedHash, hashFound := apiState.APIDefinition.ObjectMeta.Labels["apiHash"]
	if !api.Status.DeploymentStatus.Accepted {
		if apiReconciler.apiPropagationEnabled && isAPIPropagatable(apiState) {
			apiHash := apiReconciler.getAPIHash(apiState)
			push := false
			if !hashFound || storedHash != apiHash {
				// Check whether apiHash in the controlplane queue
				if !controlplane.IsAPIHashQueued(apiHash) {
					push = true
				}
			}
			if push {
				loggers.LoggerAPKOperator.Infof("API hash changed sending the API to agent")
				// Publish the api data to CP
				apiCpData := apiReconciler.convertAPIStateToAPICp(ctx, *apiState, apiHash, prodAirl, sandAirl)
				apiCpData.Event = controlplane.EventTypeCreate
				controlplane.AddToEventQueue(apiCpData)
			}
		}
		apiReconciler.ods.AddAPIState(apiRef, apiState)
		apiReconciler.traverseAPIStateAndUpdateOwnerReferences(ctx, *apiState)
		return &synchronizer.APIEvent{EventType: constants.Create, Events: []synchronizer.APIState{*apiState}, UpdatedEvents: []string{}}, nil
	} else if cachedAPI, events, updated :=
		apiReconciler.ods.UpdateAPIState(apiRef, apiState); updated {
		if apiReconciler.apiPropagationEnabled && isAPIPropagatable(apiState) {
			apiHash := apiReconciler.getAPIHash(apiState)
			push := false
			if !hashFound || storedHash != apiHash {
				// Check whether apiHash in the controlplane queue
				if !controlplane.IsAPIHashQueued(apiHash) {
					push = true
				}
			}
			if push {
				loggers.LoggerAPKOperator.Infof("API hash changed sending the API to agent")
				// Publish the api data to CP
				apiCpData := apiReconciler.convertAPIStateToAPICp(ctx, *apiState, apiHash, prodAirl, sandAirl)
				apiCpData.Event = controlplane.EventTypeUpdate
				controlplane.AddToEventQueue(apiCpData)
			}
		}
		apiReconciler.traverseAPIStateAndUpdateOwnerReferences(ctx, *apiState)
		loggers.LoggerAPKOperator.Infof("API CR %s with API UUID : %v is updated on %v", apiRef.String(),
			string(api.ObjectMeta.UID), events)
		return &synchronizer.APIEvent{EventType: constants.Update, Events: []synchronizer.APIState{cachedAPI}, UpdatedEvents: events}, nil
	}

	return nil, nil
}

func isAPIPropagatable(apiState *synchronizer.APIState) bool {
	validOrgs := []string{"carbon.super"}

	if apiState.APIDefinition.Spec.APIType == constants.GRPC {
		return false
	}

	// System APIs should not be propagated to CP
	if apiState.APIDefinition.Spec.SystemAPI {
		return false
	}
	if apiState.ProdGQLRoute == nil && apiState.ProdHTTPRoute == nil {
		return false
	}
	// Only valid organization's APIs can be propagated to CP
	return utils.ContainsString(validOrgs, apiState.APIDefinition.Spec.Organization)
}

func (apiReconciler *APIReconciler) resolveGQLRouteRefs(ctx context.Context, gqlRouteRefs []string,
	namespace string, api dpv1alpha3.API) (*synchronizer.GQLRouteState, error) {
	gqlRouteState, err := apiReconciler.concatGQLRoutes(ctx, gqlRouteRefs, namespace, api)
	if err != nil {
		return nil, err
	}
	gqlRouteState.Scopes, err = apiReconciler.getScopesForGQLRoute(ctx, gqlRouteState.GQLRouteCombined, api)
	return &gqlRouteState, err
}

// resolveHTTPRouteRefs validates following references related to the API
// - Authentications
func (apiReconciler *APIReconciler) resolveHTTPRouteRefs(ctx context.Context, httpRouteState *synchronizer.HTTPRouteState,
	httpRouteRefs []string, namespace string, interceptorServiceMapping map[string]dpv1alpha1.InterceptorService,
	api dpv1alpha3.API) (*synchronizer.HTTPRouteState, *dpv1alpha3.AIRateLimitPolicy, error) {
	var err error
	httpRouteState.HTTPRouteCombined, httpRouteState.HTTPRoutePartitions, err = apiReconciler.concatHTTPRoutes(ctx, httpRouteRefs, namespace, api)
	if err != nil {
		return nil, nil, err
	}
	var airl *dpv1alpha3.AIRateLimitPolicy
	httpRouteState.BackendMapping, airl, err = apiReconciler.getResolvedBackendsMapping(ctx, httpRouteState, interceptorServiceMapping, api)
	if err != nil {
		return nil, nil, err
	}
	httpRouteState.Scopes, err = apiReconciler.getScopesForHTTPRoute(ctx, httpRouteState.HTTPRouteCombined, api)

	return httpRouteState, airl, err
}

func (apiReconciler *APIReconciler) resolveGRPCRouteRefs(ctx context.Context, grpcRouteRefs []string,
	namespace string, interceptorServiceMapping map[string]dpv1alpha1.InterceptorService, api dpv1alpha3.API) (*synchronizer.GRPCRouteState, error) {
	grpcRouteState, err := apiReconciler.concatGRPCRoutes(ctx, grpcRouteRefs, namespace, api)
	if err != nil {
		return nil, err
	}
	grpcRouteState.BackendMapping, err = apiReconciler.getResolvedBackendsMappingForGRPC(ctx, &grpcRouteState, interceptorServiceMapping, api)
	if err != nil {
		return nil, err
	}
	grpcRouteState.Scopes, err = apiReconciler.getScopesForGRPCRoute(ctx, grpcRouteState.GRPCRouteCombined, api)
	return &grpcRouteState, err
}

func (apiReconciler *APIReconciler) concatGRPCRoutes(ctx context.Context, grpcRouteRefs []string,
	namespace string, api dpv1alpha3.API) (synchronizer.GRPCRouteState, error) {
	grpcRouteState := synchronizer.GRPCRouteState{}
	grpcRoutePartitions := make(map[string]*gwapiv1.GRPCRoute)
	for _, grpcRouteRef := range grpcRouteRefs {
		var grpcRoute gwapiv1.GRPCRoute
		namespacedName := types.NamespacedName{Namespace: namespace, Name: grpcRouteRef}
		if err := utils.ResolveRef(ctx, apiReconciler.client, &api, namespacedName, true, &grpcRoute); err != nil {
			return grpcRouteState, fmt.Errorf("error while getting grpcroute %s in namespace :%s, %s", grpcRouteRef,
				namespace, err.Error())
		}
		grpcRoutePartitions[namespacedName.String()] = &grpcRoute
		if grpcRouteState.GRPCRouteCombined == nil {
			grpcRouteState.GRPCRouteCombined = &grpcRoute
		} else {
			grpcRouteState.GRPCRouteCombined.Spec.Rules = append(grpcRouteState.GRPCRouteCombined.Spec.Rules,
				grpcRoute.Spec.Rules...)
		}
	}
	grpcRouteState.GRPCRoutePartitions = grpcRoutePartitions
	backendNamespacedName := types.NamespacedName{
		Name:      string(grpcRouteState.GRPCRouteCombined.Spec.Rules[0].BackendRefs[0].BackendRef.Name),
		Namespace: namespace,
	}
	resolvedBackend := utils.GetResolvedBackend(ctx, apiReconciler.client, backendNamespacedName, &api)
	if resolvedBackend != nil {
		grpcRouteState.BackendMapping = map[string]*dpv1alpha2.ResolvedBackend{
			backendNamespacedName.String(): resolvedBackend,
		}
		return grpcRouteState, nil
	}
	return grpcRouteState, errors.New("error while resolving backend for grpcroute")
}

func (apiReconciler *APIReconciler) concatGQLRoutes(ctx context.Context, gqlRouteRefs []string,
	namespace string, api dpv1alpha3.API) (synchronizer.GQLRouteState, error) {
	gqlRouteState := synchronizer.GQLRouteState{}
	gqlRoutePartitions := make(map[string]*dpv1alpha2.GQLRoute)
	for _, gqlRouteRef := range gqlRouteRefs {
		var gqlRoute dpv1alpha2.GQLRoute
		namespacedName := types.NamespacedName{Namespace: namespace, Name: gqlRouteRef}
		if err := utils.ResolveRef(ctx, apiReconciler.client, &api, namespacedName, true, &gqlRoute); err != nil {
			return gqlRouteState, fmt.Errorf("error while getting gqlroute %s in namespace :%s, %s", gqlRouteRef,
				namespace, err.Error())
		}
		gqlRoutePartitions[namespacedName.String()] = &gqlRoute
		if gqlRouteState.GQLRouteCombined == nil {
			gqlRouteState.GQLRouteCombined = &gqlRoute
		} else {
			gqlRouteState.GQLRouteCombined.Spec.Rules = append(gqlRouteState.GQLRouteCombined.Spec.Rules,
				gqlRoute.Spec.Rules...)
		}
	}
	gqlRouteState.GQLRoutePartitions = gqlRoutePartitions
	backendNamespacedName := types.NamespacedName{
		Name:      string(gqlRouteState.GQLRouteCombined.Spec.BackendRefs[0].Name),
		Namespace: namespace,
	}
	resolvedBackend := utils.GetResolvedBackend(ctx, apiReconciler.client, backendNamespacedName, &api)
	if resolvedBackend != nil {
		gqlRouteState.BackendMapping = map[string]*dpv1alpha2.ResolvedBackend{
			backendNamespacedName.String(): resolvedBackend,
		}
		return gqlRouteState, nil
	}
	return gqlRouteState, errors.New("error while resolving backend for gqlroute")
}

func (apiReconciler *APIReconciler) getScopesForGRPCRoute(ctx context.Context,
	grpcRoute *gwapiv1.GRPCRoute, api dpv1alpha3.API) (map[string]dpv1alpha1.Scope, error) {
	scopes := make(map[string]dpv1alpha1.Scope)
	for _, rule := range grpcRoute.Spec.Rules {
		for _, filter := range rule.Filters {
			if filter.ExtensionRef != nil && filter.ExtensionRef.Kind == constants.KindScope {
				scope := &dpv1alpha1.Scope{}
				if err := utils.ResolveRef(ctx, apiReconciler.client, &api,
					types.NamespacedName{Namespace: grpcRoute.Namespace, Name: string(filter.ExtensionRef.Name)}, false,
					scope); err != nil {
					return nil, fmt.Errorf("error while getting scope %s in namespace :%s, %s", filter.ExtensionRef.Name,
						grpcRoute.Namespace, err.Error())
				}
				scopes[utils.NamespacedName(scope).String()] = *scope
			}
		}
	}
	return scopes, nil
}

func (apiReconciler *APIReconciler) concatHTTPRoutes(ctx context.Context, httpRouteRefs []string,
	namespace string, api dpv1alpha3.API) (*gwapiv1.HTTPRoute, map[string]*gwapiv1.HTTPRoute, error) {
	var combinedHTTPRoute *gwapiv1.HTTPRoute
	httpRoutePartitions := make(map[string]*gwapiv1.HTTPRoute)
	for _, httpRouteRef := range httpRouteRefs {
		var httpRoute gwapiv1.HTTPRoute
		namespacedName := types.NamespacedName{Namespace: namespace, Name: httpRouteRef}
		if err := utils.ResolveRef(ctx, apiReconciler.client, &api,
			namespacedName, true, &httpRoute); err != nil {
			return nil, httpRoutePartitions, fmt.Errorf("error while getting httproute %s in namespace :%s, %s", httpRouteRef, namespace, err.Error())
		}
		httpRoutePartitions[namespacedName.String()] = &httpRoute
		if combinedHTTPRoute == nil {
			combinedHTTPRoute = &httpRoute
		} else {
			combinedHTTPRoute.Spec.Rules = append(combinedHTTPRoute.Spec.Rules, httpRoute.Spec.Rules...)
		}
	}
	return combinedHTTPRoute, httpRoutePartitions, nil
}

func (apiReconciler *APIReconciler) getAuthenticationsForAPI(ctx context.Context,
	api dpv1alpha3.API) (map[string]dpv1alpha2.Authentication, error) {
	nameSpacedName := utils.NamespacedName(&api).String()
	authentications := make(map[string]dpv1alpha2.Authentication)
	authenticationList := &dpv1alpha2.AuthenticationList{}
	if err := apiReconciler.client.List(ctx, authenticationList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(apiAuthenticationIndex, nameSpacedName),
	}); err != nil {
		return nil, err
	}
	for item := range authenticationList.Items {
		authenticationListItem := authenticationList.Items[item]
		if authenticationListItem.Spec.TargetRef.Name == v1alpha2.ObjectName(api.Name) {
			authentications[utils.NamespacedName(&authenticationListItem).String()] = authenticationListItem
		}
	}
	return authentications, nil
}

func (apiReconciler *APIReconciler) getRatelimitPoliciesForAPI(ctx context.Context,
	api dpv1alpha3.API) (map[string]dpv1alpha3.RateLimitPolicy, error) {
	nameSpacedName := utils.NamespacedName(&api).String()
	ratelimitPolicies := make(map[string]dpv1alpha3.RateLimitPolicy)
	ratelimitPolicyList := &dpv1alpha3.RateLimitPolicyList{}
	if err := apiReconciler.client.List(ctx, ratelimitPolicyList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(apiRateLimitIndex, nameSpacedName),
	}); err != nil {
		return nil, err
	}
	for item := range ratelimitPolicyList.Items {
		rateLimitPolicy := ratelimitPolicyList.Items[item]
		ratelimitPolicies[utils.NamespacedName(&rateLimitPolicy).String()] = rateLimitPolicy
	}
	return ratelimitPolicies, nil
}

func (apiReconciler *APIReconciler) getScopesForGQLRoute(ctx context.Context,
	gqlRoute *dpv1alpha2.GQLRoute, api dpv1alpha3.API) (map[string]dpv1alpha1.Scope, error) {
	scopes := make(map[string]dpv1alpha1.Scope)
	for _, rule := range gqlRoute.Spec.Rules {
		for _, filter := range rule.Filters {
			if filter.ExtensionRef != nil && filter.ExtensionRef.Kind == constants.KindScope {
				scope := &dpv1alpha1.Scope{}
				if err := utils.ResolveRef(ctx, apiReconciler.client, &api,
					types.NamespacedName{Namespace: gqlRoute.Namespace, Name: string(filter.ExtensionRef.Name)}, false,
					scope); err != nil {
					return nil, fmt.Errorf("error while getting scope %s in namespace :%s, %s", filter.ExtensionRef.Name,
						gqlRoute.Namespace, err.Error())
				}
				scopes[utils.NamespacedName(scope).String()] = *scope
			}
		}
	}
	return scopes, nil
}

func (apiReconciler *APIReconciler) getScopesForHTTPRoute(ctx context.Context,
	httpRoute *gwapiv1.HTTPRoute, api dpv1alpha3.API) (map[string]dpv1alpha1.Scope, error) {
	scopes := make(map[string]dpv1alpha1.Scope)
	for _, rule := range httpRoute.Spec.Rules {
		for _, filter := range rule.Filters {
			if filter.Type == gwapiv1.HTTPRouteFilterExtensionRef && filter.ExtensionRef != nil &&
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
	api dpv1alpha3.API) (map[string]dpv1alpha2.Authentication, error) {
	nameSpacedName := utils.NamespacedName(&api).String()
	authentications := make(map[string]dpv1alpha2.Authentication)
	authenticationList := &dpv1alpha2.AuthenticationList{}
	if err := apiReconciler.client.List(ctx, authenticationList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(apiAuthenticationResourceIndex, nameSpacedName),
	}); err != nil {
		return nil, err
	}
	for item := range authenticationList.Items {
		authenticationListItem := authenticationList.Items[item]
		if authenticationListItem.Spec.TargetRef.Name != v1alpha2.ObjectName(api.Name) {
			authentications[utils.NamespacedName(&authenticationListItem).String()] = authenticationListItem
		}
	}
	return authentications, nil
}

func (apiReconciler *APIReconciler) getRatelimitPoliciesForResources(ctx context.Context,
	api dpv1alpha3.API) (map[string]dpv1alpha3.RateLimitPolicy, error) {
	nameSpacedName := utils.NamespacedName(&api).String()
	ratelimitpolicies := make(map[string]dpv1alpha3.RateLimitPolicy)
	ratelimitPolicyList := &dpv1alpha3.RateLimitPolicyList{}
	if err := apiReconciler.client.List(ctx, ratelimitPolicyList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(apiRateLimitResourceIndex, nameSpacedName),
	}); err != nil {
		return nil, err
	}
	for item := range ratelimitPolicyList.Items {
		rateLimitPolicy := ratelimitPolicyList.Items[item]
		ratelimitpolicies[utils.NamespacedName(&rateLimitPolicy).String()] = rateLimitPolicy
	}
	return ratelimitpolicies, nil
}

func (apiReconciler *APIReconciler) getAPIPoliciesForAPI(ctx context.Context, api dpv1alpha3.API) (map[string]dpv1alpha4.APIPolicy, error) {
	nameSpacedName := utils.NamespacedName(&api).String()
	apiPolicies := make(map[string]dpv1alpha4.APIPolicy)
	apiPolicyList := &dpv1alpha4.APIPolicyList{}
	if err := apiReconciler.client.List(ctx, apiPolicyList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(apiAPIPolicyIndex, nameSpacedName),
	}); err != nil {
		return nil, err
	}
	for item := range apiPolicyList.Items {
		apiPolicy := apiPolicyList.Items[item]
		loggers.LoggerAPKOperator.Debugf("API Policy %+v", &apiPolicy)
		apiPolicies[utils.NamespacedName(&apiPolicy).String()] = apiPolicy
	}
	return apiPolicies, nil
}

func (apiReconciler *APIReconciler) getAPIDefinitionForAPI(ctx context.Context,
	apiDefinitionFile, namespace string, api dpv1alpha3.API) ([]byte, error) {
	configMap := &corev1.ConfigMap{}
	if err := utils.ResolveRef(ctx, apiReconciler.client, &api,
		types.NamespacedName{Namespace: namespace, Name: apiDefinitionFile}, true, configMap); err != nil {
		return nil, fmt.Errorf("error while getting swagger definition %s in namespace :%s, %s", apiDefinitionFile,
			namespace, err.Error())
	}

	var apiDef []byte
	for _, val := range configMap.BinaryData {
		// config map data key is "swagger.yaml"
		apiDef = []byte(val)
	}
	return apiDef, nil
}

func (apiReconciler *APIReconciler) getAPIPoliciesForResources(ctx context.Context,
	api dpv1alpha3.API) (map[string]dpv1alpha4.APIPolicy, error) {
	nameSpacedName := utils.NamespacedName(&api).String()
	apiPolicies := make(map[string]dpv1alpha4.APIPolicy)
	apiPolicyList := &dpv1alpha4.APIPolicyList{}
	if err := apiReconciler.client.List(ctx, apiPolicyList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(apiAPIPolicyResourceIndex, nameSpacedName),
	}); err != nil {
		return nil, err
	}
	for item := range apiPolicyList.Items {
		apiPolicy := apiPolicyList.Items[item]
		apiPolicies[utils.NamespacedName(&apiPolicy).String()] = apiPolicy
	}
	return apiPolicies, nil
}

// getAPIPolicyChildrenRefs gets all the referenced policies in apipolicy for the resolving API
// - interceptor services
// - backend JWTs
// - subscription validation
func (apiReconciler *APIReconciler) getAPIPolicyChildrenRefs(ctx context.Context,
	apiPolicies, resourceAPIPolicies map[string]dpv1alpha4.APIPolicy,
	api dpv1alpha3.API) (map[string]dpv1alpha1.InterceptorService, map[string]dpv1alpha1.BackendJWT, bool, *dpv1alpha4.AIProvider, *synchronizer.ResolvedModelBasedRoundRobin, error) {
	allAPIPolicies := append(maps.Values(apiPolicies), maps.Values(resourceAPIPolicies)...)
	interceptorServices := make(map[string]dpv1alpha1.InterceptorService)
	backendJWTs := make(map[string]dpv1alpha1.BackendJWT)
	aiProvider := &dpv1alpha4.AIProvider{}
	modelBasedRoundRobin := &dpv1alpha4.ModelBasedRoundRobin{}
	resolvedModelBasedRoundRobin := &synchronizer.ResolvedModelBasedRoundRobin{}
	subscriptionValidation := false
	for _, apiPolicy := range allAPIPolicies {
		if apiPolicy.Spec.Default != nil {
			subscriptionValidation = subscriptionValidation || apiPolicy.Spec.Default.SubscriptionValidation
		}
		if apiPolicy.Spec.Override != nil {
			subscriptionValidation = subscriptionValidation || apiPolicy.Spec.Override.SubscriptionValidation
		}
	}
	for _, apiPolicy := range allAPIPolicies {
		if apiPolicy.Spec.Default != nil {
			if len(apiPolicy.Spec.Default.RequestInterceptors) > 0 {
				interceptorPtr := utils.GetInterceptorService(ctx, apiReconciler.client, apiPolicy.Namespace,
					&apiPolicy.Spec.Default.RequestInterceptors[0], &api)
				if interceptorPtr != nil {
					interceptorServices[utils.NamespacedName(interceptorPtr).String()] = *interceptorPtr
				}
			}
			if len(apiPolicy.Spec.Default.ResponseInterceptors) > 0 {
				interceptorPtr := utils.GetInterceptorService(ctx, apiReconciler.client, apiPolicy.Namespace,
					&apiPolicy.Spec.Default.ResponseInterceptors[0], &api)
				if interceptorPtr != nil {
					interceptorServices[utils.NamespacedName(interceptorPtr).String()] = *interceptorPtr
				}
			}
			if apiPolicy.Spec.Default.BackendJWTPolicy != nil {
				backendJWTPtr := utils.GetBackendJWT(ctx, apiReconciler.client, apiPolicy.Namespace,
					apiPolicy.Spec.Default.BackendJWTPolicy.Name, &api)
				if backendJWTPtr != nil {
					backendJWTs[utils.NamespacedName(backendJWTPtr).String()] = *backendJWTPtr
				}
			}
			if apiPolicy.Spec.Default.AIProvider != nil {
				loggers.LoggerAPKOperator.Debugf("AIProvider Default found in API Policy. AI Provider Name %s", apiPolicy.Spec.Default.AIProvider.Name)
				aiProviderPtr := utils.GetAIProvider(ctx, apiReconciler.client, apiPolicy.Namespace,
					apiPolicy.Spec.Default.AIProvider.Name, &api)
				if aiProviderPtr != nil {
					aiProvider = aiProviderPtr
				}
			}
		}
		if apiPolicy.Spec.Override != nil {
			if len(apiPolicy.Spec.Override.RequestInterceptors) > 0 {
				interceptorPtr := utils.GetInterceptorService(ctx, apiReconciler.client, apiPolicy.Namespace,
					&apiPolicy.Spec.Override.RequestInterceptors[0], &api)
				if interceptorPtr != nil {
					interceptorServices[utils.NamespacedName(interceptorPtr).String()] = *interceptorPtr
				}
			}
			if len(apiPolicy.Spec.Override.ResponseInterceptors) > 0 {
				interceptorPtr := utils.GetInterceptorService(ctx, apiReconciler.client, apiPolicy.Namespace,
					&apiPolicy.Spec.Override.ResponseInterceptors[0], &api)
				if interceptorPtr != nil {
					interceptorServices[utils.NamespacedName(interceptorPtr).String()] = *interceptorPtr
				}
			}
			if apiPolicy.Spec.Override.BackendJWTPolicy != nil {
				backendJWTPtr := utils.GetBackendJWT(ctx, apiReconciler.client, apiPolicy.Namespace,
					apiPolicy.Spec.Override.BackendJWTPolicy.Name, &api)
				if backendJWTPtr != nil {
					backendJWTs[utils.NamespacedName(backendJWTPtr).String()] = *backendJWTPtr
				}
			}
			if apiPolicy.Spec.Override.AIProvider != nil {
				loggers.LoggerAPKOperator.Debugf("AIProvider override found in API Policy. AI Provider Name %s", apiPolicy.Spec.Override.AIProvider.Name)
				aiProviderPtr := utils.GetAIProvider(ctx, apiReconciler.client, apiPolicy.Namespace,
					apiPolicy.Spec.Override.AIProvider.Name, &api)
				if aiProviderPtr != nil {
					aiProvider = aiProviderPtr
				}
			}

		}
	}
	for _, apiPolicy := range apiPolicies {
		backendMapping := make(map[string]*dpv1alpha2.ResolvedBackend)
		if apiPolicy.Spec.Default != nil {
			if apiPolicy.Spec.Default.ModelBasedRoundRobin != nil {
				loggers.LoggerAPKOperator.Infof("ModelBasedRoundRobin Default found in API Policy. ModelBasedRoundRobin Model %v", apiPolicy.Spec.Default.ModelBasedRoundRobin)
				modelBasedRoundRobin = apiPolicy.Spec.Default.ModelBasedRoundRobin
				resolvedModelBasedRoundRobin = &synchronizer.ResolvedModelBasedRoundRobin{
					OnQuotaExceedSuspendDuration: modelBasedRoundRobin.OnQuotaExceedSuspendDuration,
				}
				if modelBasedRoundRobin.ProductionModels != nil {
					loggers.LoggerAPKOperator.Infof("ProductionModels Default found in API Policy. ModelBasedRoundRobin Model %v", modelBasedRoundRobin.ProductionModels)
					productionModels := apiPolicy.Spec.Default.ModelBasedRoundRobin.ProductionModels
					for _, model := range productionModels {
						resolvedBackend := &dpv1alpha2.ResolvedBackend{}
						if model.BackendRef.Name != "" {
							backendNamespacedName := types.NamespacedName{
								Name:      string(model.BackendRef.Name),
								Namespace: utils.GetNamespace(model.BackendRef.Namespace, apiPolicy.Namespace),
							}
							if _, exists := backendMapping[backendNamespacedName.String()]; !exists {
								resolvedBackend = utils.GetResolvedBackend(ctx, apiReconciler.client, backendNamespacedName, &api)
								if resolvedBackend != nil {
									backendMapping[backendNamespacedName.String()] = resolvedBackend
								} else {
									return nil, nil, false, nil, nil, fmt.Errorf("unable to find backend %s", backendNamespacedName.String())
								}
							} else {
								resolvedBackend = backendMapping[backendNamespacedName.String()]
							}
						}
						resolvedModelWeight := synchronizer.ResolvedModelWeight{
							Model:           model.Model,
							Weight:          model.Weight,
							ResolvedBackend: resolvedBackend,
						}
						resolvedModelBasedRoundRobin.ProductionModels = append(resolvedModelBasedRoundRobin.ProductionModels, resolvedModelWeight)
					}
				}
				if modelBasedRoundRobin.SandboxModels != nil {
					loggers.LoggerAPKOperator.Infof("SandboxModels Default found in API Policy. ModelBasedRoundRobin Model %v", modelBasedRoundRobin.SandboxModels)
					sandboxModels := apiPolicy.Spec.Default.ModelBasedRoundRobin.SandboxModels
					for _, model := range sandboxModels {
						resolvedBackend := &dpv1alpha2.ResolvedBackend{}
						if model.BackendRef.Name != "" {
							backendNamespacedName := types.NamespacedName{
								Name:      string(model.BackendRef.Name),
								Namespace: utils.GetNamespace(model.BackendRef.Namespace, apiPolicy.Namespace),
							}
							if _, exists := backendMapping[backendNamespacedName.String()]; !exists {
								resolvedBackend = utils.GetResolvedBackend(ctx, apiReconciler.client, backendNamespacedName, &api)
								if resolvedBackend != nil {
									backendMapping[backendNamespacedName.String()] = resolvedBackend
								} else {
									return nil, nil, false, nil, nil, fmt.Errorf("unable to find backend %s", backendNamespacedName.String())
								}
							} else {
								resolvedBackend = backendMapping[backendNamespacedName.String()]
							}
						}
						resolvedModelWeight := synchronizer.ResolvedModelWeight{
							Model:           model.Model,
							Weight:          model.Weight,
							ResolvedBackend: resolvedBackend,
						}
						resolvedModelBasedRoundRobin.SandboxModels = append(resolvedModelBasedRoundRobin.SandboxModels, resolvedModelWeight)
					}
				}
			}
		}
		if apiPolicy.Spec.Override != nil {
			if apiPolicy.Spec.Override.ModelBasedRoundRobin != nil {
				loggers.LoggerAPKOperator.Infof("ModelBasedRoundRobin override found in API Policy. ModelBasedRoundRobin Model %v", apiPolicy.Spec.Override.ModelBasedRoundRobin)
				modelBasedRoundRobin = apiPolicy.Spec.Override.ModelBasedRoundRobin
				resolvedModelBasedRoundRobin = &synchronizer.ResolvedModelBasedRoundRobin{
					OnQuotaExceedSuspendDuration: modelBasedRoundRobin.OnQuotaExceedSuspendDuration,
				}
				if modelBasedRoundRobin.ProductionModels != nil {
					loggers.LoggerAPKOperator.Infof("ProductionModels override found in API Policy. ModelBasedRoundRobin Model %v", modelBasedRoundRobin.ProductionModels)
					productionModels := apiPolicy.Spec.Override.ModelBasedRoundRobin.ProductionModels
					for _, model := range productionModels {
						resolvedBackend := &dpv1alpha2.ResolvedBackend{}
						if model.BackendRef.Name != "" {
							backendNamespacedName := types.NamespacedName{
								Name:      string(model.BackendRef.Name),
								Namespace: utils.GetNamespace(model.BackendRef.Namespace, apiPolicy.Namespace),
							}
							if _, exists := backendMapping[backendNamespacedName.String()]; !exists {
								resolvedBackend = utils.GetResolvedBackend(ctx, apiReconciler.client, backendNamespacedName, &api)
								if resolvedBackend != nil {
									backendMapping[backendNamespacedName.String()] = resolvedBackend
								} else {
									return nil, nil, false, nil, nil, fmt.Errorf("unable to find backend %s", backendNamespacedName.String())
								}
							} else {
								resolvedBackend = backendMapping[backendNamespacedName.String()]
							}
						}
						resolvedModelWeight := synchronizer.ResolvedModelWeight{
							Model:           model.Model,
							Weight:          model.Weight,
							ResolvedBackend: resolvedBackend,
						}
						resolvedModelBasedRoundRobin.ProductionModels = append(resolvedModelBasedRoundRobin.ProductionModels, resolvedModelWeight)
					}
				}
				if modelBasedRoundRobin.SandboxModels != nil {
					loggers.LoggerAPKOperator.Infof("SandboxModels override found in API Policy. ModelBasedRoundRobin Model %v", modelBasedRoundRobin.SandboxModels)
					sandboxModels := apiPolicy.Spec.Override.ModelBasedRoundRobin.SandboxModels
					for _, model := range sandboxModels {
						resolvedBackend := &dpv1alpha2.ResolvedBackend{}
						if model.BackendRef.Name != "" {
							backendNamespacedName := types.NamespacedName{
								Name:      string(model.BackendRef.Name),
								Namespace: utils.GetNamespace(model.BackendRef.Namespace, apiPolicy.Namespace),
							}
							if _, exists := backendMapping[backendNamespacedName.String()]; !exists {
								resolvedBackend = utils.GetResolvedBackend(ctx, apiReconciler.client, backendNamespacedName, &api)
								if resolvedBackend != nil {
									backendMapping[backendNamespacedName.String()] = resolvedBackend
								} else {
									return nil, nil, false, nil, nil, fmt.Errorf("unable to find backend %s", backendNamespacedName.String())
								}
							} else {
								resolvedBackend = backendMapping[backendNamespacedName.String()]
							}
						}
						resolvedModelWeight := synchronizer.ResolvedModelWeight{
							Model:           model.Model,
							Weight:          model.Weight,
							ResolvedBackend: resolvedBackend,
						}
						resolvedModelBasedRoundRobin.SandboxModels = append(resolvedModelBasedRoundRobin.SandboxModels, resolvedModelWeight)
					}
				}
			}
		}
	}
	return interceptorServices, backendJWTs, subscriptionValidation, aiProvider, resolvedModelBasedRoundRobin, nil
}

func (apiReconciler *APIReconciler) resolveAuthentications(ctx context.Context,
	authentications map[string]dpv1alpha2.Authentication) (*dpv1alpha2.MutualSSL, error) {
	resolvedMutualSSL := dpv1alpha2.MutualSSL{}
	for _, authentication := range authentications {
		err := utils.GetResolvedMutualSSL(ctx, apiReconciler.client, authentication, &resolvedMutualSSL)
		if err != nil {
			return nil, err
		}
	}
	return &resolvedMutualSSL, nil
}

func (apiReconciler *APIReconciler) getResolvedBackendsMapping(ctx context.Context,
	httpRouteState *synchronizer.HTTPRouteState, interceptorServiceMapping map[string]dpv1alpha1.InterceptorService,
	api dpv1alpha3.API) (map[string]*dpv1alpha2.ResolvedBackend, *dpv1alpha3.AIRateLimitPolicy, error) {
	backendMapping := make(map[string]*dpv1alpha2.ResolvedBackend)
	var airl *dpv1alpha3.AIRateLimitPolicy
	// Resolve backends in HTTPRoute
	httpRoute := httpRouteState.HTTPRouteCombined
	ruleIdxToAiRatelimitPolicyMapping := make(map[int]*dpv1alpha3.AIRateLimitPolicy)
	httpRouteState.RuleIdxToAiRatelimitPolicyMapping = ruleIdxToAiRatelimitPolicyMapping
	for id, rule := range httpRoute.Spec.Rules {
		for _, backend := range rule.BackendRefs {
			backendNamespacedName := types.NamespacedName{
				Name:      string(backend.Name),
				Namespace: utils.GetNamespace(backend.Namespace, httpRoute.Namespace),
			}
			aiRLPolicyList := &dpv1alpha3.AIRateLimitPolicyList{}
			if err := apiReconciler.client.List(ctx, aiRLPolicyList, &k8client.ListOptions{
				FieldSelector: fields.OneTermEqualSelector(aiRatelimitPolicyToBackendIndex, backendNamespacedName.String()),
			}); err != nil {
				loggers.LoggerAPKOperator.Debugf("No associated AI ratelimit policy found for : %s", backendNamespacedName.String())
			} else {
				for _, aiRLPolicy := range aiRLPolicyList.Items {
					loggers.LoggerAPKOperator.Debugf("Adding mapping for ruleid: %d to aiRLPolicy: %s", id, utils.NamespacedName(&aiRLPolicy))
					if aiRLPolicy.Spec.Override == nil {
						aiRLPolicy.Spec.Override = aiRLPolicy.Spec.Default
					}
					ruleIdxToAiRatelimitPolicyMapping[id] = &aiRLPolicy
					airl = &aiRLPolicy
				}
			}
			if _, exists := backendMapping[backendNamespacedName.String()]; !exists {
				resolvedBackend := utils.GetResolvedBackend(ctx, apiReconciler.client, backendNamespacedName, &api)
				if resolvedBackend != nil {
					backendMapping[backendNamespacedName.String()] = resolvedBackend
				} else {
					return nil, nil, fmt.Errorf("unable to find backend %s", backendNamespacedName.String())
				}
			}

		}

		for _, filter := range rule.Filters {
			if filter.RequestMirror != nil {
				mirrorBackend := filter.RequestMirror.BackendRef
				mirrorBackendNamespacedName := types.NamespacedName{
					Name:      string(mirrorBackend.Name),
					Namespace: utils.GetNamespace(mirrorBackend.Namespace, httpRoute.Namespace),
				}
				if string(*mirrorBackend.Kind) == constants.KindBackend {
					if _, exists := backendMapping[mirrorBackendNamespacedName.String()]; !exists {
						resolvedMirrorBackend := utils.GetResolvedBackend(ctx, apiReconciler.client, mirrorBackendNamespacedName, &api)
						if resolvedMirrorBackend != nil {
							backendMapping[mirrorBackendNamespacedName.String()] = resolvedMirrorBackend
						} else {
							return nil, nil, fmt.Errorf("unable to find backend %s", mirrorBackendNamespacedName.String())
						}
					}
				} else if string(*mirrorBackend.Kind) == constants.KindService {
					var err error
					service, err := utils.GetService(ctx, apiReconciler.client, utils.GetNamespace(mirrorBackend.Namespace, httpRoute.Namespace), string(mirrorBackend.Name))
					if err != nil {
						return nil, nil, fmt.Errorf("unable to find service %s", mirrorBackendNamespacedName.String())
					}
					backendMapping[mirrorBackendNamespacedName.String()], err = utils.GetResolvedBackendFromService(service, int(*mirrorBackend.Port))
					if err != nil {
						return nil, nil, fmt.Errorf("error in getting service information %s", service)
					}
				}
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
	return backendMapping, airl, nil
}

func (apiReconciler *APIReconciler) getResolvedBackendsMappingForGRPC(ctx context.Context,
	grpcRouteState *synchronizer.GRPCRouteState, interceptorServiceMapping map[string]dpv1alpha1.InterceptorService,
	api dpv1alpha3.API) (map[string]*dpv1alpha2.ResolvedBackend, error) {
	backendMapping := make(map[string]*dpv1alpha2.ResolvedBackend)
	grpcRoute := grpcRouteState.GRPCRouteCombined

	for _, rule := range grpcRoute.Spec.Rules {
		for _, backend := range rule.BackendRefs {
			backendNamespacedName := types.NamespacedName{
				Name:      string(backend.Name),
				Namespace: utils.GetNamespace(backend.Namespace, grpcRoute.Namespace),
			}
			if _, exists := backendMapping[backendNamespacedName.String()]; !exists {
				resolvedBackend := utils.GetResolvedBackend(ctx, apiReconciler.client, backendNamespacedName, &api)
				if resolvedBackend != nil {
					backendMapping[backendNamespacedName.String()] = resolvedBackend
				} else {
					return nil, fmt.Errorf("unable to find backend %s", backendNamespacedName.String())
				}
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
	return backendMapping, nil
}

// These proxy methods are designed as intermediaries for the getAPIsFor<CR objects> methods.
// Their purpose is to encapsulate the process of updating owner references within the reconciliation watch methods.
// By employing these proxies, we prevent redundant owner reference updates for the same object due to the hierarchical structure of these functions.
func (apiReconciler *APIReconciler) populateAPIReconcileRequestsForGQLRoute(ctx context.Context, obj *dpv1alpha2.GQLRoute) []reconcile.Request {
	requests := apiReconciler.getAPIForGQLRoute(ctx, obj)
	if len(requests) > 0 {
		apiReconciler.handleOwnerReference(ctx, obj, &requests)
	}
	return requests
}

func (apiReconciler *APIReconciler) populateAPIReconcileRequestsForHTTPRoute(ctx context.Context, obj *gwapiv1.HTTPRoute) []reconcile.Request {
	requests := apiReconciler.getAPIForHTTPRoute(ctx, obj)
	if len(requests) > 0 {
		apiReconciler.handleOwnerReference(ctx, obj, &requests)
	}
	return requests
}

func (apiReconciler *APIReconciler) populateAPIReconcileRequestsForGRPCRoute(ctx context.Context, obj *gwapiv1.GRPCRoute) []reconcile.Request {
	requests := apiReconciler.getAPIForGRPCRoute(ctx, obj)
	apiReconciler.handleOwnerReference(ctx, obj, &requests)
	return requests
}

func (apiReconciler *APIReconciler) populateAPIReconcileRequestsForConfigMap(ctx context.Context, obj *corev1.ConfigMap) []reconcile.Request {
	requests := apiReconciler.getAPIsForConfigMap(ctx, obj)
	if len(requests) > 0 {
		apiReconciler.handleOwnerReference(ctx, obj, &requests)
	}
	return requests
}

func (apiReconciler *APIReconciler) populateAPIReconcileRequestsForSecret(ctx context.Context, obj *corev1.Secret) []reconcile.Request {
	requests := apiReconciler.getAPIsForSecret(ctx, obj)
	if len(requests) > 0 {
		apiReconciler.handleOwnerReference(ctx, obj, &requests)
	}
	return requests
}

func (apiReconciler *APIReconciler) populateAPIReconcileRequestsForAIRatelimitPolicy(ctx context.Context, obj *dpv1alpha3.AIRateLimitPolicy) []reconcile.Request {
	requests := apiReconciler.getAPIsForAIRatelimitPolicy(ctx, obj)
	if len(requests) > 0 {
		apiReconciler.handleOwnerReference(ctx, obj, &requests)
	}
	return requests
}

func (apiReconciler *APIReconciler) populateAPIReconcileRequestsForAuthentication(ctx context.Context, obj *dpv1alpha2.Authentication) []reconcile.Request {
	requests := apiReconciler.getAPIsForAuthentication(ctx, obj)
	if len(requests) > 0 {
		apiReconciler.handleOwnerReference(ctx, obj, &requests)
	}
	return requests
}

func (apiReconciler *APIReconciler) populateAPIReconcileRequestsForAPIPolicy(ctx context.Context, obj *dpv1alpha4.APIPolicy) []reconcile.Request {
	requests := apiReconciler.getAPIsForAPIPolicy(ctx, obj)
	if len(requests) > 0 {
		apiReconciler.handleOwnerReference(ctx, obj, &requests)
	}
	return requests
}

func (apiReconciler *APIReconciler) populateAPIReconcileRequestsForInterceptorService(ctx context.Context, obj *dpv1alpha1.InterceptorService) []reconcile.Request {
	requests := apiReconciler.getAPIsForInterceptorService(ctx, obj)
	if len(requests) > 0 {
		apiReconciler.handleOwnerReference(ctx, obj, &requests)
	}
	return requests
}

func (apiReconciler *APIReconciler) populateAPIReconcileRequestsForBackendJWT(ctx context.Context, obj *dpv1alpha1.BackendJWT) []reconcile.Request {
	requests := apiReconciler.getAPIsForBackendJWT(ctx, obj)
	if len(requests) > 0 {
		apiReconciler.handleOwnerReference(ctx, obj, &requests)
	}
	return requests
}

func (apiReconciler *APIReconciler) populateAPIReconcileRequestsForRateLimitPolicy(ctx context.Context, obj *dpv1alpha3.RateLimitPolicy) []reconcile.Request {
	requests := apiReconciler.getAPIsForRateLimitPolicy(ctx, obj)
	if len(requests) > 0 {
		apiReconciler.handleOwnerReference(ctx, obj, &requests)
	}
	return requests
}

func (apiReconciler *APIReconciler) populateAPIReconcileRequestsForScope(ctx context.Context, obj *dpv1alpha1.Scope) []reconcile.Request {
	requests := apiReconciler.getAPIsForScope(ctx, obj)
	if len(requests) > 0 {
		apiReconciler.handleOwnerReference(ctx, obj, &requests)
	}
	return requests
}

func (apiReconciler *APIReconciler) populateAPIReconcileRequestsForBackend(ctx context.Context, obj *dpv1alpha2.Backend) []reconcile.Request {
	requests := apiReconciler.getAPIsForBackend(ctx, obj)
	if len(requests) > 0 {
		apiReconciler.handleOwnerReference(ctx, obj, &requests)
	}
	return requests
}

func (apiReconciler *APIReconciler) populateAPIReconcileRequestsForAIProvider(ctx context.Context, obj *dpv1alpha4.AIProvider) []reconcile.Request {
	requests := apiReconciler.getAPIsForAIProvider(ctx, obj)
	if len(requests) > 0 {
		apiReconciler.handleOwnerReference(ctx, obj, &requests)
	}
	return requests
}

func (apiReconciler *APIReconciler) traverseAPIStateAndUpdateOwnerReferences(ctx context.Context, apiState synchronizer.APIState) {
	// travserse through all the children of this API and trigger update owner reference
	if apiState.ProdHTTPRoute != nil {
		for _, httpRoute := range apiState.ProdHTTPRoute.HTTPRoutePartitions {
			apiReconciler.retrieveParentAPIsAndUpdateOwnerReference(ctx, httpRoute)
		}
	}
	if apiState.SandHTTPRoute != nil {
		for _, httpRoute := range apiState.SandHTTPRoute.HTTPRoutePartitions {
			apiReconciler.retrieveParentAPIsAndUpdateOwnerReference(ctx, httpRoute)
		}
	}
	if apiState.ProdGQLRoute != nil {
		for _, gqlRoute := range apiState.ProdGQLRoute.GQLRoutePartitions {
			apiReconciler.retrieveParentAPIsAndUpdateOwnerReference(ctx, gqlRoute)
		}
	}
	if apiState.SandGQLRoute != nil {
		for _, gqlRoute := range apiState.SandGQLRoute.GQLRoutePartitions {
			apiReconciler.retrieveParentAPIsAndUpdateOwnerReference(ctx, gqlRoute)
		}
	}

	if apiState.ProdGRPCRoute != nil {
		for _, grpcRoute := range apiState.ProdGRPCRoute.GRPCRoutePartitions {
			apiReconciler.retrieveParentAPIsAndUpdateOwnerReference(ctx, grpcRoute)
		}
	}
	if apiState.SandGRPCRoute != nil {
		for _, grpcRoute := range apiState.SandGRPCRoute.GRPCRoutePartitions {
			apiReconciler.retrieveParentAPIsAndUpdateOwnerReference(ctx, grpcRoute)
		}
	}

	for _, auth := range apiState.Authentications {
		apiReconciler.retrieveParentAPIsAndUpdateOwnerReference(ctx, &auth)
	}
	for _, auth := range apiState.ResourceAuthentications {
		apiReconciler.retrieveParentAPIsAndUpdateOwnerReference(ctx, &auth)
	}
	for _, ratelimit := range apiState.RateLimitPolicies {
		apiReconciler.retrieveParentAPIsAndUpdateOwnerReference(ctx, &ratelimit)
	}
	for _, ratelimit := range apiState.ResourceRateLimitPolicies {
		apiReconciler.retrieveParentAPIsAndUpdateOwnerReference(ctx, &ratelimit)
	}
	for _, apiPolicy := range apiState.APIPolicies {
		apiReconciler.retrieveParentAPIsAndUpdateOwnerReference(ctx, &apiPolicy)
	}
	for _, apiPolicy := range apiState.ResourceAPIPolicies {
		apiReconciler.retrieveParentAPIsAndUpdateOwnerReference(ctx, &apiPolicy)
	}
	for _, interceptorService := range apiState.InterceptorServiceMapping {
		apiReconciler.retrieveParentAPIsAndUpdateOwnerReference(ctx, &interceptorService)
	}
	if apiState.ProdHTTPRoute != nil {
		for _, backend := range apiState.ProdHTTPRoute.BackendMapping {
			if backend != nil {
				apiReconciler.retrieveParentAPIsAndUpdateOwnerReference(ctx, &backend.Backend)
			}
		}
	}
	if apiState.SandHTTPRoute != nil {
		for _, backend := range apiState.SandHTTPRoute.BackendMapping {
			if backend != nil {
				apiReconciler.retrieveParentAPIsAndUpdateOwnerReference(ctx, &backend.Backend)
			}
		}
	}
	if apiState.ProdGQLRoute != nil {
		for _, backend := range apiState.ProdGQLRoute.BackendMapping {
			if backend != nil {
				apiReconciler.retrieveParentAPIsAndUpdateOwnerReference(ctx, &backend.Backend)
			}
		}
	}
	if apiState.SandGQLRoute != nil {
		for _, backend := range apiState.SandGQLRoute.BackendMapping {
			if backend != nil {
				apiReconciler.retrieveParentAPIsAndUpdateOwnerReference(ctx, &backend.Backend)
			}
		}
	}

	if apiState.ProdGRPCRoute != nil {
		for _, backend := range apiState.ProdGRPCRoute.BackendMapping {
			if &backend != nil {
				apiReconciler.retrieveParentAPIsAndUpdateOwnerReference(ctx, &backend.Backend)
			}
		}
	}
	if apiState.SandGRPCRoute != nil {
		for _, backend := range apiState.SandGRPCRoute.BackendMapping {
			if &backend != nil {
				apiReconciler.retrieveParentAPIsAndUpdateOwnerReference(ctx, &backend.Backend)
			}
		}
	}
	for _, backendJwt := range apiState.BackendJWTMapping {
		apiReconciler.retrieveParentAPIsAndUpdateOwnerReference(ctx, &backendJwt)
	}

}

func (apiReconciler *APIReconciler) retrieveParentAPIsAndUpdateOwnerReference(ctx context.Context, obj k8client.Object) {
	var requests []reconcile.Request
	switch obj.(type) {
	case *dpv1alpha2.Backend:
		var backend dpv1alpha2.Backend
		namesapcedName := types.NamespacedName{
			Name:      string(obj.GetName()),
			Namespace: string(obj.GetNamespace()),
		}
		if err := apiReconciler.client.Get(ctx, namesapcedName, &backend); err != nil {
			loggers.LoggerAPKOperator.Errorf("Unexpected error occured while loading the cr object from cluster %+v", err)
			return
		}
		requests = apiReconciler.getAPIsForBackend(ctx, &backend)
		apiReconciler.handleOwnerReference(ctx, &backend, &requests)
	case *dpv1alpha1.Scope:
		var scope dpv1alpha1.Scope
		namesapcedName := types.NamespacedName{
			Name:      string(obj.GetName()),
			Namespace: string(obj.GetNamespace()),
		}
		if err := apiReconciler.client.Get(ctx, namesapcedName, &scope); err != nil {
			loggers.LoggerAPKOperator.Errorf("Unexpected error occured while loading the cr object from cluster %+v", err)
			return
		}
		requests = apiReconciler.getAPIsForScope(ctx, &scope)
		apiReconciler.handleOwnerReference(ctx, &scope, &requests)
	case *dpv1alpha3.RateLimitPolicy:
		var rl dpv1alpha3.RateLimitPolicy
		namesapcedName := types.NamespacedName{
			Name:      string(obj.GetName()),
			Namespace: string(obj.GetNamespace()),
		}
		if err := apiReconciler.client.Get(ctx, namesapcedName, &rl); err != nil {
			loggers.LoggerAPKOperator.Errorf("Unexpected error occured while loading the cr object from cluster %+v", err)
			return
		}
		requests = apiReconciler.getAPIsForRateLimitPolicy(ctx, &rl)
		apiReconciler.handleOwnerReference(ctx, &rl, &requests)
	case *dpv1alpha1.BackendJWT:
		var backendJWT dpv1alpha1.BackendJWT
		namesapcedName := types.NamespacedName{
			Name:      string(obj.GetName()),
			Namespace: string(obj.GetNamespace()),
		}
		if err := apiReconciler.client.Get(ctx, namesapcedName, &backendJWT); err != nil {
			loggers.LoggerAPKOperator.Errorf("Unexpected error occured while loading the cr object from cluster %+v", err)
			return
		}
		requests = apiReconciler.getAPIsForBackendJWT(ctx, &backendJWT)
		apiReconciler.handleOwnerReference(ctx, &backendJWT, &requests)
	case *dpv1alpha1.InterceptorService:
		var interceptorService dpv1alpha1.InterceptorService
		namesapcedName := types.NamespacedName{
			Name:      string(obj.GetName()),
			Namespace: string(obj.GetNamespace()),
		}
		if err := apiReconciler.client.Get(ctx, namesapcedName, &interceptorService); err != nil {
			loggers.LoggerAPKOperator.Errorf("Unexpected error occured while loading the cr object from cluster %+v", err)
			return
		}
		requests = apiReconciler.getAPIsForInterceptorService(ctx, &interceptorService)
		apiReconciler.handleOwnerReference(ctx, &interceptorService, &requests)
	case *dpv1alpha4.APIPolicy:
		var apiPolicy dpv1alpha4.APIPolicy
		namesapcedName := types.NamespacedName{
			Name:      string(obj.GetName()),
			Namespace: string(obj.GetNamespace()),
		}
		if err := apiReconciler.client.Get(ctx, namesapcedName, &apiPolicy); err != nil {
			loggers.LoggerAPKOperator.Errorf("Unexpected error occured while loading the cr object from cluster %+v", err)
			return
		}
		requests = apiReconciler.getAPIsForAPIPolicy(ctx, &apiPolicy)
		apiReconciler.handleOwnerReference(ctx, &apiPolicy, &requests)
	case *dpv1alpha2.Authentication:
		var auth dpv1alpha2.Authentication
		namesapcedName := types.NamespacedName{
			Name:      string(obj.GetName()),
			Namespace: string(obj.GetNamespace()),
		}
		if err := apiReconciler.client.Get(ctx, namesapcedName, &auth); err != nil {
			loggers.LoggerAPKOperator.Errorf("Unexpected error occured while loading the cr object from cluster %+v", err)
			return
		}
		requests = apiReconciler.getAPIsForAuthentication(ctx, &auth)
		apiReconciler.handleOwnerReference(ctx, &auth, &requests)
	case *corev1.Secret:
		var secret corev1.Secret
		namesapcedName := types.NamespacedName{
			Name:      string(obj.GetName()),
			Namespace: string(obj.GetNamespace()),
		}
		if err := apiReconciler.client.Get(ctx, namesapcedName, &secret); err != nil {
			loggers.LoggerAPKOperator.Errorf("Unexpected error occured while loading the cr object from cluster %+v", err)
			return
		}
		requests = apiReconciler.getAPIsForSecret(ctx, &secret)
		apiReconciler.handleOwnerReference(ctx, &secret, &requests)
	case *corev1.ConfigMap:
		var cm corev1.ConfigMap
		namesapcedName := types.NamespacedName{
			Name:      string(obj.GetName()),
			Namespace: string(obj.GetNamespace()),
		}
		if err := apiReconciler.client.Get(ctx, namesapcedName, &cm); err != nil {
			loggers.LoggerAPKOperator.Errorf("Unexpected error occured while loading the cr object from cluster %+v", err)
			return
		}
		requests = apiReconciler.getAPIsForConfigMap(ctx, &cm)
		apiReconciler.handleOwnerReference(ctx, &cm, &requests)
	case *gwapiv1.HTTPRoute:
		var httpRoute gwapiv1.HTTPRoute
		namesapcedName := types.NamespacedName{
			Name:      string(obj.GetName()),
			Namespace: string(obj.GetNamespace()),
		}
		if err := apiReconciler.client.Get(ctx, namesapcedName, &httpRoute); err != nil {
			loggers.LoggerAPKOperator.Errorf("Unexpected error occured while loading the cr object from cluster %+v", err)
			return
		}
		requests = apiReconciler.getAPIForHTTPRoute(ctx, &httpRoute)
		apiReconciler.handleOwnerReference(ctx, &httpRoute, &requests)
	case *dpv1alpha2.GQLRoute:
		var gqlRoute dpv1alpha2.GQLRoute
		namesapcedName := types.NamespacedName{
			Name:      string(obj.GetName()),
			Namespace: string(obj.GetNamespace()),
		}
		if err := apiReconciler.client.Get(ctx, namesapcedName, &gqlRoute); err != nil {
			loggers.LoggerAPKOperator.Errorf("Unexpected error occured while loading the cr object from cluster %+v", err)
			return
		}
		requests = apiReconciler.getAPIForGQLRoute(ctx, &gqlRoute)
		apiReconciler.handleOwnerReference(ctx, &gqlRoute, &requests)
	case *dpv1alpha4.AIProvider:
		var aiProvider dpv1alpha4.AIProvider
		namesapcedName := types.NamespacedName{
			Name:      string(obj.GetName()),
			Namespace: string(obj.GetNamespace()),
		}
		if err := apiReconciler.client.Get(ctx, namesapcedName, &aiProvider); err != nil {
			loggers.LoggerAPKOperator.Errorf("Unexpected error occured while loading the cr object from cluster %+v", err)
			return
		}
		requests = apiReconciler.getAPIsForAIProvider(ctx, &aiProvider)
	case *gwapiv1.GRPCRoute:
		var grpcRoute gwapiv1.GRPCRoute
		namespaceName := types.NamespacedName{
			Name:      string(obj.GetName()),
			Namespace: string(obj.GetNamespace()),
		}
		if err := apiReconciler.client.Get(ctx, namespaceName, &grpcRoute); err != nil {
			loggers.LoggerAPKOperator.Errorf("Unexpected error occured while loading the cr object from cluster %+v", err)
			return
		}
	default:
		loggers.LoggerAPKOperator.Errorf("Unexpected type found while processing owner reference %+v", obj)
	}

}

// getAPIForGQLRoute triggers the API controller reconcile method based on the changes detected
// from GQLRoute objects. If the changes are done for an API stored in the Operator Data store,
// a new reconcile event will be created and added to the reconcile event queue.
func (apiReconciler *APIReconciler) getAPIForGQLRoute(ctx context.Context, obj k8client.Object) []reconcile.Request {
	gqlRoute, ok := obj.(*dpv1alpha2.GQLRoute)
	if !ok {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2665, logging.TRIVIAL, "Unexpected object type, bypassing reconciliation: %v", gqlRoute))
		return []reconcile.Request{}
	}
	apiList := &dpv1alpha3.APIList{}
	if err := apiReconciler.client.List(ctx, apiList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(gqlRouteAPIIndex, utils.NamespacedName(gqlRoute).String()),
	}); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2666, logging.CRITICAL, "Unable to find associated APIs: %s", utils.NamespacedName(gqlRoute).String()))
		return []reconcile.Request{}
	}
	if len(apiList.Items) == 0 {
		loggers.LoggerAPKOperator.Debugf("APIs for GQLRoute not found: %s", utils.NamespacedName(gqlRoute).String())
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
		loggers.LoggerAPKOperator.Infof("Adding reconcile request for API: %s/%s with API UUID: %v due to change in GQLRoute: %v", api.Namespace, api.Name,
			string(api.ObjectMeta.UID), utils.NamespacedName(gqlRoute).String())
	}
	return requests
}

// getAPIForHTTPRoute triggers the API controller reconcile method based on the changes detected
// from HTTPRoute objects. If the changes are done for an API stored in the Operator Data store,
// a new reconcile event will be created and added to the reconcile event queue.
func (apiReconciler *APIReconciler) getAPIForHTTPRoute(ctx context.Context, obj k8client.Object) []reconcile.Request {
	httpRoute, ok := obj.(*gwapiv1.HTTPRoute)
	if !ok {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2622, logging.TRIVIAL, "Unexpected object type, bypassing reconciliation: %v", httpRoute))
		return []reconcile.Request{}
	}

	apiList := &dpv1alpha3.APIList{}
	if err := apiReconciler.client.List(ctx, apiList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(httpRouteAPIIndex, utils.NamespacedName(httpRoute).String()),
	}); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2623, logging.CRITICAL, "Unable to find associated APIs: %s", utils.NamespacedName(httpRoute).String()))
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
		loggers.LoggerAPKOperator.Infof("Adding reconcile request for API: %s/%s with API UUID: %v due to HTTPRoute change: %v",
			api.Namespace, api.Name, string(api.ObjectMeta.UID), utils.NamespacedName(httpRoute).String())
	}
	return requests
}

// getAPIForGRPCRoute triggers the API controller reconcile method based on the changes detected
// from GRPCRoute objects. If the changes are done for an API stored in the Operator Data store,
// a new reconcile event will be created and added to the reconcile event queue.
func (apiReconciler *APIReconciler) getAPIForGRPCRoute(ctx context.Context, obj k8client.Object) []reconcile.Request {
	grpcRoute, ok := obj.(*gwapiv1.GRPCRoute)
	if !ok {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2622, logging.TRIVIAL,
			"Unexpected object type, bypassing reconciliation: %v", grpcRoute))
		return []reconcile.Request{}
	}

	apiList := &dpv1alpha3.APIList{}

	if err := apiReconciler.client.List(ctx, apiList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(grpcRouteAPIIndex, utils.NamespacedName(grpcRoute).String()),
	}); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2623, logging.CRITICAL,
			"Unable to find associated APIs: %s", utils.NamespacedName(grpcRoute).String()))
		return []reconcile.Request{}
	}

	if len(apiList.Items) == 0 {
		loggers.LoggerAPKOperator.Debugf("APIs for GRPCRoute not found: %s", utils.NamespacedName(grpcRoute).String())
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
		loggers.LoggerAPKOperator.Infof("Adding reconcile request for API: %s/%s with API UUID: %v", api.Namespace, api.Name,
			string(api.ObjectMeta.UID))
	}
	return requests
}

// getAPIsForConfigMap triggers the API controller reconcile method based on the changes detected
// in configMap resources.
func (apiReconciler *APIReconciler) getAPIsForConfigMap(ctx context.Context, obj k8client.Object) []reconcile.Request {
	configMap, ok := obj.(*corev1.ConfigMap)
	if !ok {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2622, logging.TRIVIAL,
			"Unexpected object type, bypassing reconciliation: %v", configMap))
		return []reconcile.Request{}
	}

	backendList := &dpv1alpha2.BackendList{}
	err := apiReconciler.client.List(ctx, backendList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(configMapBackend, utils.NamespacedName(configMap).String()),
	})
	if err == nil && len(backendList.Items) > 0 {
		requests := []reconcile.Request{}
		for item := range backendList.Items {
			backend := backendList.Items[item]
			requests = append(requests, apiReconciler.getAPIsForBackend(ctx, &backend)...)
		}
		return requests
	}

	apiList := &dpv1alpha3.APIList{}
	err = apiReconciler.client.List(ctx, apiList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(configMapAPIDefinition, utils.NamespacedName(configMap).String()),
	})
	if err == nil {
		requests := []reconcile.Request{}
		for _, api := range apiList.Items {
			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      api.Name,
					Namespace: api.Namespace},
			}
			requests = append(requests, req)
			loggers.LoggerAPKOperator.Infof("Adding reconcile request for API: %s/%s with API UUID: %v due to configmap change: %v",
				api.Namespace, api.Name, string(api.ObjectMeta.UID), utils.NamespacedName(configMap).String())
		}
		return requests
	}

	loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2647, logging.MINOR,
		"Unable to find associated APIs for ConfigMap: %s", utils.NamespacedName(configMap).String()))
	return []reconcile.Request{}
}

// getAPIsForSecret triggers the API controller reconcile method based on the changes detected
// in secret resources.
func (apiReconciler *APIReconciler) getAPIsForSecret(ctx context.Context, obj k8client.Object) []reconcile.Request {
	secret, ok := obj.(*corev1.Secret)
	if !ok {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2622, logging.TRIVIAL, "Unexpected object type, bypassing reconciliation: %v", secret))
		return []reconcile.Request{}
	}

	backendList := &dpv1alpha2.BackendList{}
	if err := apiReconciler.client.List(ctx, backendList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(secretBackend, utils.NamespacedName(secret).String()),
	}); err != nil {
		loggers.LoggerAPKOperator.Debugf("Unable to find associated Backends for Secret: %s", utils.NamespacedName(secret).String())
		return []reconcile.Request{}
	}

	requests := []reconcile.Request{}
	for item := range backendList.Items {
		backend := backendList.Items[item]
		requests = append(requests, apiReconciler.getAPIsForBackend(ctx, &backend)...)
	}
	return requests
}

// getAPIsForAIRatelimitPolicy triggers the API controller reconcile method based on the changes detected
// in AIRatelimitPolicy resources.
func (apiReconciler *APIReconciler) getAPIsForAIRatelimitPolicy(ctx context.Context, obj k8client.Object) []reconcile.Request {
	aiRatelimitPolicy, ok := obj.(*dpv1alpha3.AIRateLimitPolicy)
	if !ok {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2622, logging.TRIVIAL, "Unexpected object type, bypassing reconciliation: %v", obj))
		return []reconcile.Request{}
	}

	if aiRatelimitPolicy.Spec.TargetRef.Kind == constants.KindBackend {
		backend := &dpv1alpha2.Backend{}
		namespacedName := types.NamespacedName{
			Name:      string(aiRatelimitPolicy.Spec.TargetRef.Name),
			Namespace: utils.GetNamespace(aiRatelimitPolicy.Spec.TargetRef.Namespace, aiRatelimitPolicy.GetNamespace()),
		}

		if err := apiReconciler.client.Get(ctx, namespacedName, backend); err != nil {
			loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2621, logging.MINOR, "Unable to find associated Backend for AIratelimitPolicy targetref: %s", namespacedName.String()))
			return []reconcile.Request{}
		}
		return apiReconciler.getAPIsForBackend(ctx, backend)
	}
	return []reconcile.Request{}
}

// getAPIForAuthentication triggers the API controller reconcile method based on the changes detected
// from Authentication objects. If the changes are done for an API stored in the Operator Data store,
// a new reconcile event will be created and added to the reconcile event queue.
func (apiReconciler *APIReconciler) getAPIsForAuthentication(ctx context.Context, obj k8client.Object) []reconcile.Request {
	authentication, ok := obj.(*dpv1alpha2.Authentication)
	if !ok {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2622, logging.TRIVIAL, "Unexpected object type, bypassing reconciliation: %v", authentication))
		return []reconcile.Request{}
	}

	requests := []reconcile.Request{}

	namespace, err := utils.ValidateAndRetrieveNamespace((*gwapiv1.Namespace)(authentication.Spec.TargetRef.Namespace), authentication.Namespace)

	if err != nil {
		loggers.LoggerAPKOperator.Errorf("Namespace mismatch. TargetRef %s needs to be in the same namespace as the Athentication %s. Expected: %s, Actual: %s",
			string(authentication.Spec.TargetRef.Name), authentication.Name, authentication.Namespace, string(*authentication.Spec.TargetRef.Namespace))
		return requests
	}

	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      string(authentication.Spec.TargetRef.Name),
			Namespace: namespace,
		},
	}
	requests = append(requests, req)
	loggers.LoggerAPKOperator.Infof("Adding reconcile request for API: %s/%s due to Authentication change: %v",
		string(authentication.Spec.TargetRef.Name), namespace, utils.NamespacedName(authentication).String())

	return requests
}

// getAPIsForAPIPolicy triggers the API controller reconcile method based on the changes detected
// from APIPolicy objects. If the changes are done for an API stored in the Operator Data store,
// a new reconcile event will be created and added to the reconcile event queue.
func (apiReconciler *APIReconciler) getAPIsForAPIPolicy(ctx context.Context, obj k8client.Object) []reconcile.Request {
	apiPolicy, ok := obj.(*dpv1alpha4.APIPolicy)
	requests := []reconcile.Request{}
	if !ok {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2622, logging.TRIVIAL, "Unexpected object type, bypassing reconciliation: %v", apiPolicy))
		return requests
	}

	if !(apiPolicy.Spec.TargetRef.Kind == constants.KindAPI || apiPolicy.Spec.TargetRef.Kind == constants.KindResource) {
		return requests
	}

	namespace, err := utils.ValidateAndRetrieveNamespace((*gwapiv1.Namespace)(apiPolicy.Spec.TargetRef.Namespace), apiPolicy.Namespace)

	if err != nil {
		loggers.LoggerAPKOperator.Errorf("Namespace mismatch. TargetRef %s needs to be in the same namespace as the ApiPolicy %s. Expected: %s, Actual: %s",
			string(apiPolicy.Spec.TargetRef.Name), apiPolicy.Name, apiPolicy.Namespace, string(*apiPolicy.Spec.TargetRef.Namespace))
		return requests
	}

	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      string(apiPolicy.Spec.TargetRef.Name),
			Namespace: namespace},
	}
	requests = append(requests, req)
	loggers.LoggerAPKOperator.Infof("Adding reconcile request for API: %s/%s due to APIPolicy change: %v",
		string(apiPolicy.Spec.TargetRef.Name), namespace, utils.NamespacedName(apiPolicy).String())

	return requests
}

// getAPIsForAIProvider triggers the API controller reconcile method based on the changes detected
// from AIProvider objects. If the changes are done for an API stored in the Operator Data store,
// a new reconcile event will be created and added to the reconcile event queue.
func (apiReconciler *APIReconciler) getAPIsForAIProvider(ctx context.Context, obj k8client.Object) []reconcile.Request {
	aiProvider, ok := obj.(*dpv1alpha4.AIProvider)
	if !ok {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2622, logging.TRIVIAL, "Unexpected object type, bypassing reconciliation: %v", aiProvider))
		return []reconcile.Request{}
	}

	apiPolicyList := &dpv1alpha4.APIPolicyList{}
	if err := apiReconciler.client.List(ctx, apiPolicyList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(aiProviderAPIPolicyIndex, utils.NamespacedName(aiProvider).String()),
	}); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2649, logging.CRITICAL, "Unable to find associated APIPolicies: %s, error: %v", utils.NamespacedName(aiProvider).String(), err.Error()))
		return []reconcile.Request{}
	}

	requests := []reconcile.Request{}
	for item := range apiPolicyList.Items {
		apiPolicy := apiPolicyList.Items[item]
		requests = append(requests, apiReconciler.getAPIsForAPIPolicy(ctx, &apiPolicy)...)
	}
	return requests
}

// getAPIPoliciesForInterceptorService returns associated APIPolicies for the InterceptorService
// when the changes detected in InterceptorService resoruces.
func (apiReconciler *APIReconciler) getAPIsForInterceptorService(ctx context.Context, obj k8client.Object) []reconcile.Request {
	interceptorService, ok := obj.(*dpv1alpha1.InterceptorService)
	if !ok {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2622, logging.TRIVIAL, "Unexpected object type, bypassing reconciliation: %v", interceptorService))
		return []reconcile.Request{}
	}

	apiPolicyList := &dpv1alpha4.APIPolicyList{}
	if err := apiReconciler.client.List(ctx, apiPolicyList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(interceptorServiceAPIPolicyIndex, utils.NamespacedName(interceptorService).String()),
	}); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2649, logging.CRITICAL, "Unable to find associated APIPolicies: %s, error: %v", utils.NamespacedName(interceptorService).String(), err.Error()))
		return []reconcile.Request{}
	}

	requests := []reconcile.Request{}
	for item := range apiPolicyList.Items {
		apiPolicy := apiPolicyList.Items[item]
		requests = append(requests, apiReconciler.getAPIsForAPIPolicy(ctx, &apiPolicy)...)
	}
	return requests
}

// getAPIsForBackendJWT returns associated apipolicy for the backendjwt
// when the changes detected in backendjwt resources.
func (apiReconciler *APIReconciler) getAPIsForBackendJWT(ctx context.Context, obj k8client.Object) []reconcile.Request {
	backendJWT, ok := obj.(*dpv1alpha1.BackendJWT)
	if !ok {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2622, logging.TRIVIAL, "Unexpected object type, bypassing reconciliation: %v", backendJWT))
		return []reconcile.Request{}
	}

	apiPolicyList := &dpv1alpha4.APIPolicyList{}
	if err := apiReconciler.client.List(ctx, apiPolicyList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(backendJWTAPIPolicyIndex, utils.NamespacedName(backendJWT).String()),
	}); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2651, logging.CRITICAL, "Error while getting interceptor service %s, %s", utils.NamespacedName(backendJWT).String(), err.Error()))
		return []reconcile.Request{}
	}

	requests := []reconcile.Request{}
	for item := range apiPolicyList.Items {
		apiPolicy := apiPolicyList.Items[item]
		requests = append(requests, apiReconciler.getAPIsForAPIPolicy(ctx, &apiPolicy)...)
	}
	return requests
}

// getAPIsForRateLimitPolicy triggers the API controller reconcile method based on the changes detected
// from RateLimitPolicy objects. If the changes are done for an API stored in the Operator Data store,
// a new reconcile event will be created and added to the reconcile event queue.
func (apiReconciler *APIReconciler) getAPIsForRateLimitPolicy(ctx context.Context, obj k8client.Object) []reconcile.Request {
	ratelimitPolicy, ok := obj.(*dpv1alpha3.RateLimitPolicy)
	requests := []reconcile.Request{}
	if !ok {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2622, logging.TRIVIAL, "Unexpected object type, bypassing reconciliation: %v", ratelimitPolicy))
		return requests
	}

	if !(ratelimitPolicy.Spec.TargetRef.Kind == constants.KindAPI || ratelimitPolicy.Spec.TargetRef.Kind == constants.KindResource) {
		return requests
	}

	namespace, err := utils.ValidateAndRetrieveNamespace((*gwapiv1.Namespace)(ratelimitPolicy.Spec.TargetRef.Namespace), ratelimitPolicy.Namespace)

	if err != nil {
		loggers.LoggerAPKOperator.Errorf("Namespace mismatch. TargetRef %s needs to be in the same namespace as the RatelimitPolicy %s. Expected: %s, Actual: %s",
			string(ratelimitPolicy.Spec.TargetRef.Name), ratelimitPolicy.Name, ratelimitPolicy.Namespace, string(*ratelimitPolicy.Spec.TargetRef.Namespace))
		return requests
	}

	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      string(ratelimitPolicy.Spec.TargetRef.Name),
			Namespace: namespace},
	}
	requests = append(requests, req)
	loggers.LoggerAPKOperator.Infof("Adding reconcile request for API: %s/%s due to RateLimitPolicy change: %v",
		string(ratelimitPolicy.Spec.TargetRef.Name), namespace, utils.NamespacedName(ratelimitPolicy).String())

	return requests
}

// getAPIsForScope triggers the API controller reconcile method based on the changes detected
// from scope objects. If the changes are done for an API stored in the Operator Data store,
// a new reconcile event will be created and added to the reconcile event queue.
func (apiReconciler *APIReconciler) getAPIsForScope(ctx context.Context, obj k8client.Object) []reconcile.Request {
	scope, ok := obj.(*dpv1alpha1.Scope)
	if !ok {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2622, logging.TRIVIAL, "Unexpected object type, bypassing reconciliation: %v", scope))
		return []reconcile.Request{}
	}

	httpRouteList := &gwapiv1.HTTPRouteList{}
	if err := apiReconciler.client.List(ctx, httpRouteList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(httprouteScopeIndex, utils.NamespacedName(scope).String()),
	}); err != nil {
		loggers.LoggerAPKOperator.Warn(logging.PrintError(logging.Error2625, logging.CRITICAL, "Unable to find associated HTTPRoutes: %s", utils.NamespacedName(scope).String()))
	}

	if len(httpRouteList.Items) == 0 {
		loggers.LoggerAPKOperator.Debugf("HTTPRoutes for scope not found: %s", utils.NamespacedName(scope).String())
	}

	requests := []reconcile.Request{}
	for item := range httpRouteList.Items {
		httpRoute := httpRouteList.Items[item]
		requests = append(requests, apiReconciler.getAPIForHTTPRoute(ctx, &httpRoute)...)
	}

	gqlRouteList := &dpv1alpha2.GQLRouteList{}
	if err := apiReconciler.client.List(ctx, gqlRouteList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(gqlRouteScopeIndex, utils.NamespacedName(scope).String()),
	}); err != nil {
		loggers.LoggerAPKOperator.Warn(logging.PrintError(logging.Error2625, logging.CRITICAL, "Unable to find associated GQLRoute: %s", utils.NamespacedName(scope).String()))
	}

	if len(gqlRouteList.Items) == 0 {
		loggers.LoggerAPKOperator.Debugf("GQLRoutes for scope not found: %s", utils.NamespacedName(scope).String())
	}

	for item := range gqlRouteList.Items {
		httpRoute := gqlRouteList.Items[item]
		requests = append(requests, apiReconciler.getAPIForGQLRoute(ctx, &httpRoute)...)
	}

	grpcRouteList := &gwapiv1.GRPCRouteList{}
	if err := apiReconciler.client.List(ctx, grpcRouteList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(grpcRouteScopeIndex, utils.NamespacedName(scope).String()),
	}); err != nil {
		loggers.LoggerAPKOperator.Warn(logging.PrintError(logging.Error2625, logging.CRITICAL, "Unable to find associated GRPCRoutes: %s", utils.NamespacedName(scope).String()))
	}

	if len(grpcRouteList.Items) == 0 {
		loggers.LoggerAPKOperator.Debugf("GRPCRoutes for scope not found: %s", utils.NamespacedName(scope).String())
	}
	for item := range grpcRouteList.Items {
		grpcRoute := grpcRouteList.Items[item]
		requests = append(requests, apiReconciler.getAPIForGRPCRoute(ctx, &grpcRoute)...)
	}

	return requests
}

// getAPIsForBackend triggers the API controller reconcile method based on the changes detected
// in backend resources.
func (apiReconciler *APIReconciler) getAPIsForBackend(ctx context.Context, obj k8client.Object) []reconcile.Request {
	backend, ok := obj.(*dpv1alpha2.Backend)
	if !ok {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2622, logging.TRIVIAL, "Unexpected object type, bypassing reconciliation: %v", backend))
		return []reconcile.Request{}
	}

	httpRouteList := &gwapiv1.HTTPRouteList{}
	if err := apiReconciler.client.List(ctx, httpRouteList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(backendHTTPRouteIndex, utils.NamespacedName(backend).String()),
	}); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2625, logging.CRITICAL, "Unable to find associated HTTPRoutes: %s", utils.NamespacedName(backend).String()))
		return []reconcile.Request{}
	}

	if len(httpRouteList.Items) == 0 {
		loggers.LoggerAPKOperator.Debugf("HTTPRoutes for Backend not found: %s", utils.NamespacedName(backend).String())
	}

	requests := []reconcile.Request{}
	for item := range httpRouteList.Items {
		httpRoute := httpRouteList.Items[item]
		requests = append(requests, apiReconciler.getAPIForHTTPRoute(ctx, &httpRoute)...)
	}

	gqlRouteList := &dpv1alpha2.GQLRouteList{}
	if err := apiReconciler.client.List(ctx, gqlRouteList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(backendGQLRouteIndex, utils.NamespacedName(backend).String()),
	}); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2625, logging.CRITICAL, "Unable to find associated HTTPRoutes: %s", utils.NamespacedName(backend).String()))
		return []reconcile.Request{}
	}

	if len(gqlRouteList.Items) == 0 {
		loggers.LoggerAPKOperator.Debugf("GQLRoutes for Backend not found: %s", utils.NamespacedName(backend).String())
	}
	for item := range gqlRouteList.Items {
		gqlRoute := gqlRouteList.Items[item]
		requests = append(requests, apiReconciler.getAPIForGQLRoute(ctx, &gqlRoute)...)
	}

	grpcRouteList := &gwapiv1.GRPCRouteList{}
	if err := apiReconciler.client.List(ctx, grpcRouteList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(backendGRPCRouteIndex, utils.NamespacedName(backend).String()),
	}); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2625, logging.CRITICAL, "Unable to find associated GRPCRoutes: %s", utils.NamespacedName(backend).String()))
		return []reconcile.Request{}
	}

	if len(grpcRouteList.Items) == 0 {
		loggers.LoggerAPKOperator.Debugf("GRPCRoutes for Backend not found: %s", utils.NamespacedName(backend).String())
	}

	for item := range grpcRouteList.Items {
		grpcRoute := grpcRouteList.Items[item]
		requests = append(requests, apiReconciler.getAPIForGRPCRoute(ctx, &grpcRoute)...)
	}

	// Create API reconcile events when Backend reffered from InterceptorService
	interceptorServiceList := &dpv1alpha1.InterceptorServiceList{}
	if err := apiReconciler.client.List(ctx, interceptorServiceList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(backendInterceptorServiceIndex, utils.NamespacedName(backend).String()),
	}); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2649, logging.CRITICAL, "Unable to find associated APIPolicies: %s", utils.NamespacedName(backend).String()))
		return []reconcile.Request{}
	}

	if len(interceptorServiceList.Items) == 0 {
		loggers.LoggerAPKOperator.Debugf("InterceptorService for Backend not found: %s", utils.NamespacedName(backend).String())
	}

	for item := range interceptorServiceList.Items {
		interceptorService := interceptorServiceList.Items[item]
		requests = append(requests, apiReconciler.getAPIsForInterceptorService(ctx, &interceptorService)...)
	}

	return requests
}

// getAPIsForGateway triggers the API controller reconcile method based on the changes detected
// in gateway resources.
func (apiReconciler *APIReconciler) getAPIsForGateway(ctx context.Context, obj *gwapiv1.Gateway) []reconcile.Request {
	gateway := obj

	httpRouteList := &gwapiv1.HTTPRouteList{}
	if err := apiReconciler.client.List(ctx, httpRouteList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(gatewayHTTPRouteIndex, utils.NamespacedName(gateway).String()),
	}); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2625, logging.CRITICAL, "Unable to find associated HTTPRoutes: %s", utils.NamespacedName(gateway).String()))
		return []reconcile.Request{}
	}

	grpcRouteList := &gwapiv1.GRPCRouteList{}
	if err := apiReconciler.client.List(ctx, grpcRouteList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(gatewayGRPCRouteIndex, utils.NamespacedName(gateway).String()),
	}); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2625, logging.CRITICAL, "Unable to find associated GRPCRoutes: %s", utils.NamespacedName(gateway).String()))
		return []reconcile.Request{}
	}

	requests := []reconcile.Request{}
	for item := range httpRouteList.Items {
		httpRoute := httpRouteList.Items[item]
		requests = append(requests, apiReconciler.getAPIForHTTPRoute(ctx, &httpRoute)...)
	}
	for item := range grpcRouteList.Items {
		grpcRoute := grpcRouteList.Items[item]
		requests = append(requests, apiReconciler.getAPIForGRPCRoute(ctx, &grpcRoute)...)
	}
	return requests
}

// addIndexes adds indexing on API, for
//   - production and sandbox HTTPRoutes
//     referenced in API objects via `.spec.prodHTTPRouteRef` and `.spec.sandHTTPRouteRef`
//     This helps to find APIs that are affected by a HTTPRoute CRUD operation.
//   - authentications
//     authentication schemes related to httproutes
//     This helps to find authentication schemes binded to HTTPRoute.
//   - apiPolicies
//     apiPolicy schemes related to httproutes
//     This helps to find apiPolicy schemes binded to HTTPRoute.
func addIndexes(ctx context.Context, mgr manager.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha3.API{}, httpRouteAPIIndex,
		func(rawObj k8client.Object) []string {
			api := rawObj.(*dpv1alpha3.API)
			var httpRoutes []string
			if len(api.Spec.Production) > 0 {
				for _, ref := range api.Spec.Production[0].RouteRefs {
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
				for _, ref := range api.Spec.Sandbox[0].RouteRefs {
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

	if err := mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha3.API{}, gqlRouteAPIIndex,
		func(rawObj k8client.Object) []string {
			api := rawObj.(*dpv1alpha3.API)
			var gqlRoutes []string
			if len(api.Spec.Production) > 0 {
				for _, ref := range api.Spec.Production[0].RouteRefs {
					if ref != "" {
						gqlRoutes = append(gqlRoutes,
							types.NamespacedName{
								Namespace: api.Namespace,
								Name:      ref,
							}.String())
					}
				}
			}
			if len(api.Spec.Sandbox) > 0 {
				for _, ref := range api.Spec.Sandbox[0].RouteRefs {
					if ref != "" {
						gqlRoutes = append(gqlRoutes,
							types.NamespacedName{
								Namespace: api.Namespace,
								Name:      ref,
							}.String())
					}
				}
			}
			return gqlRoutes
		}); err != nil {
		return err
	}

	if err := mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha3.API{}, grpcRouteAPIIndex,
		func(rawObj k8client.Object) []string {
			api := rawObj.(*dpv1alpha3.API)
			if api.Spec.APIType != constants.GRPC {
				return nil
			}
			var grpcRoutes []string
			if len(api.Spec.Production) > 0 {
				for _, ref := range api.Spec.Production[0].RouteRefs {
					if ref != "" {
						grpcRoutes = append(grpcRoutes,
							types.NamespacedName{
								Namespace: api.Namespace,
								Name:      ref,
							}.String())
					}
				}
			}
			if len(api.Spec.Sandbox) > 0 {
				for _, ref := range api.Spec.Sandbox[0].RouteRefs {
					if ref != "" {
						grpcRoutes = append(grpcRoutes,
							types.NamespacedName{
								Namespace: api.Namespace,
								Name:      ref,
							}.String())
					}
				}
			}
			return grpcRoutes
		}); err != nil {
		return err
	}

	if err := mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha3.API{}, configMapAPIDefinition,
		func(rawObj k8client.Object) []string {
			api := rawObj.(*dpv1alpha3.API)
			var configMaps []string
			if api.Spec.DefinitionFileRef != "" {
				configMaps = append(configMaps,
					types.NamespacedName{
						Name:      string(api.Spec.DefinitionFileRef),
						Namespace: api.Namespace,
					}.String())
			}
			return configMaps
		}); err != nil {
		return err
	}

	if err := mgr.GetFieldIndexer().IndexField(ctx, &gwapiv1.HTTPRoute{}, httprouteScopeIndex,
		func(rawObj k8client.Object) []string {
			httpRoute := rawObj.(*gwapiv1.HTTPRoute)
			var scopes []string
			for _, rule := range httpRoute.Spec.Rules {
				for _, filter := range rule.Filters {
					if filter.Type == gwapiv1.HTTPRouteFilterExtensionRef {
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

	if err := mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha2.GQLRoute{}, gqlRouteAPIIndex,
		func(rawObj k8client.Object) []string {
			gqlRoute := rawObj.(*dpv1alpha2.GQLRoute)
			var scopes []string
			for _, rule := range gqlRoute.Spec.Rules {
				for _, filter := range rule.Filters {
					if filter.ExtensionRef != nil && filter.ExtensionRef.Kind == constants.KindScope {
						scopes = append(scopes, types.NamespacedName{
							Namespace: gqlRoute.Namespace,
							Name:      string(filter.ExtensionRef.Name),
						}.String())
					}
				}
			}
			return scopes
		}); err != nil {
		return err
	}

	// Backend to HTTPRoute indexer
	if err := mgr.GetFieldIndexer().IndexField(ctx, &gwapiv1.HTTPRoute{}, backendHTTPRouteIndex,
		func(rawObj k8client.Object) []string {
			httpRoute := rawObj.(*gwapiv1.HTTPRoute)
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

	// Backend to GQLRoute indexer
	if err := mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha2.GQLRoute{}, backendGQLRouteIndex,
		func(rawObj k8client.Object) []string {
			gqlRoute := rawObj.(*dpv1alpha2.GQLRoute)
			var backends []string
			for _, backendRef := range gqlRoute.Spec.BackendRefs {
				if backendRef.Kind != nil && *backendRef.Kind == constants.KindBackend {
					backends = append(backends, types.NamespacedName{
						Namespace: utils.GetNamespace(backendRef.Namespace,
							gqlRoute.ObjectMeta.Namespace),
						Name: string(backendRef.Name),
					}.String())
				}
			}

			return backends
		}); err != nil {
		return err
	}

	// Backend to GRPCRoute indexer
	if err := mgr.GetFieldIndexer().IndexField(ctx, &gwapiv1.GRPCRoute{}, backendGRPCRouteIndex,
		func(rawObj k8client.Object) []string {
			grpcRoute := rawObj.(*gwapiv1.GRPCRoute)
			var backends []string
			for _, rule := range grpcRoute.Spec.Rules {
				for _, backendRef := range rule.BackendRefs {
					if backendRef.Kind != nil && *backendRef.Kind == constants.KindBackend {
						backends = append(backends, types.NamespacedName{
							Namespace: utils.GetNamespace(backendRef.Namespace,
								grpcRoute.ObjectMeta.Namespace),
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
	if err := mgr.GetFieldIndexer().IndexField(ctx, &gwapiv1.HTTPRoute{}, gatewayHTTPRouteIndex,
		func(rawObj k8client.Object) []string {
			httpRoute := rawObj.(*gwapiv1.HTTPRoute)
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

	//Gateway to GRPCRoute indexer
	if err := mgr.GetFieldIndexer().IndexField(ctx, &gwapiv1.GRPCRoute{}, gatewayGRPCRouteIndex,
		func(rawObj k8client.Object) []string {
			grpcRoute := rawObj.(*gwapiv1.GRPCRoute)
			var gateways []string
			for _, parentRef := range grpcRoute.Spec.ParentRefs {
				gateways = append(gateways, types.NamespacedName{
					Namespace: utils.GetNamespace(parentRef.Namespace,
						grpcRoute.Namespace),
					Name: string(parentRef.Name),
				}.String())
			}
			return gateways
		}); err != nil {
		return err
	}

	// ConfigMap to Backend indexer
	if err := mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha2.Backend{}, configMapBackend,
		func(rawObj k8client.Object) []string {
			backend := rawObj.(*dpv1alpha2.Backend)
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
	if err := mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha2.Backend{}, secretBackend,
		func(rawObj k8client.Object) []string {
			backend := rawObj.(*dpv1alpha2.Backend)
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
				if backend.Spec.Security.APIKey != nil {
					secrets = append(secrets,
						types.NamespacedName{
							Name:      string(backend.Spec.Security.APIKey.ValueFrom.Name),
							Namespace: backend.Namespace,
						}.String())
				}
			}
			return secrets
		}); err != nil {
		return err
	}

	// AIRatelimitPolicy to Backend indexer
	if err := mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha3.AIRateLimitPolicy{}, aiRatelimitPolicyToBackendIndex,
		func(rawObj k8client.Object) []string {
			aiRatelimitPolicy := rawObj.(*dpv1alpha3.AIRateLimitPolicy)
			var backends []string
			namespace := utils.GetNamespace(aiRatelimitPolicy.Spec.TargetRef.Namespace, aiRatelimitPolicy.GetNamespace())
			backends = append(backends, types.NamespacedName{
				Name:      string(aiRatelimitPolicy.Spec.TargetRef.Name),
				Namespace: namespace,
			}.String())
			return backends
		}); err != nil {
		return err
	}

	// API to Subscription indexer
	if err := mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha2.API{}, apiToSubscriptionIndex,
		func(rawObj k8client.Object) []string {
			api := rawObj.(*dpv1alpha2.API)
			var apis []string
			subscriptionIdentifierForIndex := fmt.Sprintf("%s_%s", api.Spec.APIName, api.Spec.APIVersion)
			apis = append(apis, subscriptionIdentifierForIndex)
			return apis
		}); err != nil {
		return err
	}

	// authentication to API indexer
	if err := mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha2.Authentication{}, apiAuthenticationIndex,
		func(rawObj k8client.Object) []string {
			authentication := rawObj.(*dpv1alpha2.Authentication)
			var apis []string
			if authentication.Spec.TargetRef.Kind == constants.KindAPI {
				namespace, err := utils.ValidateAndRetrieveNamespace((*gwapiv1.Namespace)(authentication.Spec.TargetRef.Namespace), authentication.Namespace)
				if err != nil {
					loggers.LoggerAPKOperator.Errorf("Namespace mismatch. TargetRef %s needs to be in the same namespace as the Athentication %s. Expected: %s, Actual: %s",
						string(authentication.Spec.TargetRef.Name), authentication.Name, authentication.Namespace, string(*authentication.Spec.TargetRef.Namespace))
					return apis
				}

				apis = append(apis,
					types.NamespacedName{
						Namespace: namespace,
						Name:      string(authentication.Spec.TargetRef.Name),
					}.String())
			}
			return apis
		}); err != nil {
		return err
	}

	// Secret to Authentication indexer
	if err := mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha2.Authentication{}, secretAuthentication,
		func(rawObj k8client.Object) []string {
			authentication := rawObj.(*dpv1alpha2.Authentication)
			var secrets []string
			if authentication.Spec.Default != nil && authentication.Spec.Default.AuthTypes != nil && authentication.Spec.Default.AuthTypes.MutualSSL != nil && authentication.Spec.Default.AuthTypes.MutualSSL.SecretRefs != nil && len(authentication.Spec.Default.AuthTypes.MutualSSL.SecretRefs) > 0 {
				for _, secret := range authentication.Spec.Default.AuthTypes.MutualSSL.SecretRefs {
					if len(secret.Name) > 0 {
						secrets = append(secrets,
							types.NamespacedName{
								Name:      string(secret.Name),
								Namespace: authentication.Namespace,
							}.String())
					}
				}
			}

			if authentication.Spec.Override != nil && authentication.Spec.Override.AuthTypes != nil && authentication.Spec.Override.AuthTypes.MutualSSL != nil && authentication.Spec.Override.AuthTypes.MutualSSL.SecretRefs != nil && len(authentication.Spec.Override.AuthTypes.MutualSSL.SecretRefs) > 0 {
				for _, secret := range authentication.Spec.Override.AuthTypes.MutualSSL.SecretRefs {
					if len(secret.Name) > 0 {
						secrets = append(secrets,
							types.NamespacedName{
								Name:      string(secret.Name),
								Namespace: authentication.Namespace,
							}.String())
					}
				}

			}
			return secrets
		}); err != nil {
		return err
	}

	// ConfigMap to Authentication indexer
	if err := mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha2.Authentication{}, configMapAuthentication,
		func(rawObj k8client.Object) []string {
			authentication := rawObj.(*dpv1alpha2.Authentication)
			var configMaps []string
			if authentication.Spec.Default != nil && authentication.Spec.Default.AuthTypes != nil && authentication.Spec.Default.AuthTypes.MutualSSL != nil && authentication.Spec.Default.AuthTypes.MutualSSL.ConfigMapRefs != nil && len(authentication.Spec.Default.AuthTypes.MutualSSL.ConfigMapRefs) > 0 {
				for _, configMap := range authentication.Spec.Default.AuthTypes.MutualSSL.ConfigMapRefs {
					if len(configMap.Name) > 0 {
						configMaps = append(configMaps,
							types.NamespacedName{
								Name:      string(configMap.Name),
								Namespace: authentication.Namespace,
							}.String())
					}
				}
			}

			if authentication.Spec.Override != nil && authentication.Spec.Override.AuthTypes != nil && authentication.Spec.Override.AuthTypes.MutualSSL != nil && authentication.Spec.Override.AuthTypes.MutualSSL.ConfigMapRefs != nil && len(authentication.Spec.Override.AuthTypes.MutualSSL.ConfigMapRefs) > 0 {
				for _, configMap := range authentication.Spec.Override.AuthTypes.MutualSSL.ConfigMapRefs {
					if len(configMap.Name) > 0 {
						configMaps = append(configMaps,
							types.NamespacedName{
								Name:      string(configMap.Name),
								Namespace: authentication.Namespace,
							}.String())
					}
				}

			}
			return configMaps
		}); err != nil {
		return err
	}

	// Till the below is httproute rule name and targetref sectionname is supported,
	// https://gateway-api.sigs.k8s.io/geps/gep-713/?h=multiple+targetrefs#apply-policies-to-sections-of-a-resource-future-extension
	// we will use a temporary kindName called Resource for policy attachments
	// TODO(amali) Fix after the official support is available
	if err := mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha2.Authentication{}, apiAuthenticationResourceIndex,
		func(rawObj k8client.Object) []string {
			authentication := rawObj.(*dpv1alpha2.Authentication)
			var apis []string
			if authentication.Spec.TargetRef.Kind == constants.KindResource {

				namespace, err := utils.ValidateAndRetrieveNamespace((*gwapiv1.Namespace)(authentication.Spec.TargetRef.Namespace), authentication.Namespace)

				if err != nil {
					loggers.LoggerAPKOperator.Errorf("Namespace mismatch. TargetRef %s needs to be in the same namespace as the Athentication %s. Expected: %s, Actual: %s",
						string(authentication.Spec.TargetRef.Name), authentication.Name, authentication.Namespace, string(*authentication.Spec.TargetRef.Namespace))
					return apis
				}

				apis = append(apis,
					types.NamespacedName{
						Namespace: namespace,
						Name:      string(authentication.Spec.TargetRef.Name),
					}.String())
			}
			return apis
		}); err != nil {
		return err
	}

	// ratelimit policy to API indexer
	if err := mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha3.RateLimitPolicy{}, apiRateLimitIndex,
		func(rawObj k8client.Object) []string {
			ratelimitPolicy := rawObj.(*dpv1alpha3.RateLimitPolicy)
			var apis []string
			if ratelimitPolicy.Spec.TargetRef.Kind == constants.KindAPI {

				namespace, err := utils.ValidateAndRetrieveNamespace((*gwapiv1.Namespace)(ratelimitPolicy.Spec.TargetRef.Namespace), ratelimitPolicy.Namespace)

				if err != nil {
					loggers.LoggerAPKOperator.Errorf("Namespace mismatch. TargetRef %s needs to be in the same namespace as the RatelimitPolicy %s. Expected: %s, Given: %s",
						string(ratelimitPolicy.Spec.TargetRef.Name), ratelimitPolicy.Name, ratelimitPolicy.Namespace, string(*ratelimitPolicy.Spec.TargetRef.Namespace))
					return apis
				}

				apis = append(apis,
					types.NamespacedName{
						Namespace: namespace,
						Name:      string(ratelimitPolicy.Spec.TargetRef.Name),
					}.String())
			}
			return apis
		}); err != nil {
		return err
	}

	// AIProvider to APIPolicy indexer
	if err := mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha4.APIPolicy{}, aiProviderAPIPolicyIndex,
		func(rawObj k8client.Object) []string {
			apiPolicy := rawObj.(*dpv1alpha4.APIPolicy)
			var aiProviders []string
			if apiPolicy.Spec.Default != nil && apiPolicy.Spec.Default.AIProvider != nil {
				aiProviders = append(aiProviders,
					types.NamespacedName{
						Namespace: apiPolicy.Namespace,
						Name:      string(apiPolicy.Spec.Default.AIProvider.Name),
					}.String())
			}
			if apiPolicy.Spec.Override != nil && apiPolicy.Spec.Override.AIProvider != nil {
				aiProviders = append(aiProviders,
					types.NamespacedName{
						Namespace: apiPolicy.Namespace,
						Name:      string(apiPolicy.Spec.Override.AIProvider.Name),
					}.String())
			}
			return aiProviders
		}); err != nil {
		return err
	}

	// Till the below is httproute rule name and targetref sectionname is supported,
	// https://gateway-api.sigs.k8s.io/geps/gep-713/?h=multiple+targetrefs#apply-policies-to-sections-of-a-resource-future-extension
	// we will use a temporary kindName called Resource for policy attachments
	// TODO(amali) Fix after the official support is available
	if err := mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha3.RateLimitPolicy{}, apiRateLimitResourceIndex,
		func(rawObj k8client.Object) []string {
			ratelimitPolicy := rawObj.(*dpv1alpha3.RateLimitPolicy)
			var apis []string
			if ratelimitPolicy.Spec.TargetRef.Kind == constants.KindResource {

				namespace, err := utils.ValidateAndRetrieveNamespace((*gwapiv1.Namespace)(ratelimitPolicy.Spec.TargetRef.Namespace), ratelimitPolicy.Namespace)

				if err != nil {
					loggers.LoggerAPKOperator.Errorf("Namespace mismatch. TargetRef %s needs to be in the same namespace as the RatelimitPolicy %s. Expected: %s, Given: %s",
						string(ratelimitPolicy.Spec.TargetRef.Name), ratelimitPolicy.Name, ratelimitPolicy.Namespace, string(*ratelimitPolicy.Spec.TargetRef.Namespace))
					return apis
				}

				apis = append(apis,
					types.NamespacedName{
						Namespace: namespace,
						Name:      string(ratelimitPolicy.Spec.TargetRef.Name),
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
					Namespace: interceptorService.Namespace,
					Name:      string(interceptorService.Spec.BackendRef.Name),
				}.String())
			return backends
		}); err != nil {
		return err
	}

	// interceptorService to APIPolicy indexer
	if err := mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha4.APIPolicy{}, interceptorServiceAPIPolicyIndex,
		func(rawObj k8client.Object) []string {
			apiPolicy := rawObj.(*dpv1alpha4.APIPolicy)
			var interceptorServices []string
			if apiPolicy.Spec.Default != nil && len(apiPolicy.Spec.Default.RequestInterceptors) > 0 {
				interceptorServices = append(interceptorServices,
					types.NamespacedName{
						Namespace: apiPolicy.Namespace,
						Name:      string(apiPolicy.Spec.Default.RequestInterceptors[0].Name),
					}.String())
			}
			if apiPolicy.Spec.Override != nil && len(apiPolicy.Spec.Override.RequestInterceptors) > 0 {
				interceptorServices = append(interceptorServices,
					types.NamespacedName{
						Namespace: apiPolicy.Namespace,
						Name:      string(apiPolicy.Spec.Override.RequestInterceptors[0].Name),
					}.String())
			}
			if apiPolicy.Spec.Default != nil && len(apiPolicy.Spec.Default.ResponseInterceptors) > 0 {
				interceptorServices = append(interceptorServices,
					types.NamespacedName{
						Namespace: apiPolicy.Namespace,
						Name:      string(apiPolicy.Spec.Default.ResponseInterceptors[0].Name),
					}.String())
			}
			if apiPolicy.Spec.Override != nil && len(apiPolicy.Spec.Override.ResponseInterceptors) > 0 {
				interceptorServices = append(interceptorServices,
					types.NamespacedName{
						Namespace: apiPolicy.Namespace,
						Name:      string(apiPolicy.Spec.Override.ResponseInterceptors[0].Name),
					}.String())
			}
			return interceptorServices
		}); err != nil {
		return err
	}

	// backendjwt to APIPolicy indexer
	if err := mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha4.APIPolicy{}, backendJWTAPIPolicyIndex,
		func(rawObj k8client.Object) []string {
			apiPolicy := rawObj.(*dpv1alpha4.APIPolicy)
			var backendJWTs []string
			if apiPolicy.Spec.Default != nil && apiPolicy.Spec.Default.BackendJWTPolicy != nil {
				backendJWTs = append(backendJWTs,
					types.NamespacedName{
						Namespace: apiPolicy.Namespace,
						Name:      string(apiPolicy.Spec.Default.BackendJWTPolicy.Name),
					}.String())
			}
			if apiPolicy.Spec.Override != nil && apiPolicy.Spec.Override.BackendJWTPolicy != nil {
				backendJWTs = append(backendJWTs,
					types.NamespacedName{
						Namespace: apiPolicy.Namespace,
						Name:      string(apiPolicy.Spec.Override.BackendJWTPolicy.Name),
					}.String())
			}
			return backendJWTs
		}); err != nil {
		return err
	}

	// httpRoute to APIPolicy indexer
	if err := mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha4.APIPolicy{}, apiAPIPolicyIndex,
		func(rawObj k8client.Object) []string {
			apiPolicy := rawObj.(*dpv1alpha4.APIPolicy)
			var apis []string
			if apiPolicy.Spec.TargetRef.Kind == constants.KindAPI {

				namespace, err := utils.ValidateAndRetrieveNamespace((*gwapiv1.Namespace)(apiPolicy.Spec.TargetRef.Namespace), apiPolicy.Namespace)

				if err != nil {
					loggers.LoggerAPKOperator.Errorf("Namespace mismatch. TargetRef %s needs to be in the same namespace as the ApiPolicy %s. Expected: %s, Actual: %s",
						string(apiPolicy.Spec.TargetRef.Name), apiPolicy.Name, apiPolicy.Namespace, string(*apiPolicy.Spec.TargetRef.Namespace))
					return apis
				}

				apis = append(apis,
					types.NamespacedName{
						Namespace: namespace,
						Name:      string(apiPolicy.Spec.TargetRef.Name),
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
	err := mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha4.APIPolicy{}, apiAPIPolicyResourceIndex,
		func(rawObj k8client.Object) []string {
			apiPolicy := rawObj.(*dpv1alpha4.APIPolicy)
			var apis []string
			if apiPolicy.Spec.TargetRef.Kind == constants.KindResource {

				namespace, err := utils.ValidateAndRetrieveNamespace((*gwapiv1.Namespace)(apiPolicy.Spec.TargetRef.Namespace), apiPolicy.Namespace)

				if err != nil {
					loggers.LoggerAPKOperator.Errorf("Namespace mismatch. TargetRef %s needs to be in the same namespace as the ApiPolicy %s. Expected: %s, Actual: %s",
						string(apiPolicy.Spec.TargetRef.Name), apiPolicy.Name, apiPolicy.Namespace, string(*apiPolicy.Spec.TargetRef.Namespace))
					return apis
				}

				apis = append(apis,
					types.NamespacedName{
						Namespace: namespace,
						Name:      string(apiPolicy.Spec.TargetRef.Name),
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

		for _, apiName := range successEvent.APINamespacedName { // handle startup multiple apis
			apiReconciler.statusUpdater.Send(status.Update{
				NamespacedName: apiName,
				Resource:       new(dpv1alpha3.API),
				UpdateStatus: func(obj k8client.Object) k8client.Object {
					h, ok := obj.(*dpv1alpha3.API)
					if !ok {
						loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2626, logging.BLOCKER, "Unsupported object type %T", obj))
					}
					hCopy := h.DeepCopy()
					hCopy.Status.DeploymentStatus.Status = successEvent.State
					hCopy.Status.DeploymentStatus.Accepted = accept
					hCopy.Status.DeploymentStatus.Message = message
					events := hCopy.Status.DeploymentStatus.Events
					// Keep the first 2 and last 3 events, remove the rest
					if len(events) > 6 {
						// Truncate the events and add the special event
						events = append(events[:2], events[len(events)-3:]...)
					}

					hCopy.Status.DeploymentStatus.Events = append(events, event)
					hCopy.Status.DeploymentStatus.TransitionTime = &timeNow
					return hCopy
				},
			})
		}
	}
}

type patchStringValue struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value string `json:"value"`
}

func (apiReconciler *APIReconciler) handleLabels(ctx context.Context) {
	loggers.LoggerAPKOperator.Info("A thread assigned to handle label updates to API CR.")
	for labelUpdate := range *controlplane.GetLabelQueue() {
		loggers.LoggerAPKOperator.Infof("Starting to process label update for API %s/%s. Labels: %+v", labelUpdate.Namespace, labelUpdate.Name, labelUpdate.Labels)

		patchOps := []patchStringValue{}
		for key, value := range labelUpdate.Labels {
			patchOps = append(patchOps, patchStringValue{
				Op:    "replace",
				Path:  fmt.Sprintf("/metadata/labels/%s", key),
				Value: value,
			})
		}
		payloadBytes, _ := json.Marshal(patchOps)
		apiCR := dpv1alpha3.API{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: labelUpdate.Namespace,
				Name:      labelUpdate.Name,
			},
		}

		err := apiReconciler.client.Patch(ctx, &apiCR, k8client.RawPatch(types.JSONPatchType, payloadBytes))
		if err != nil {
			loggers.LoggerAPKOperator.Errorf("Failed to patch api %s/%s with patch: %+v, error: %+v", labelUpdate.Name, labelUpdate.Namespace, patchOps, err)
			// Patch did not work it could be due to labels field does not exists. Lets try to update the CR with labels field.
			var apiCR dpv1alpha3.API
			if err := apiReconciler.client.Get(ctx, types.NamespacedName{Namespace: labelUpdate.Namespace, Name: labelUpdate.Name}, &apiCR); err == nil {
				if apiCR.ObjectMeta.Labels == nil {
					apiCR.ObjectMeta.Labels = map[string]string{}
				}
				for key, value := range labelUpdate.Labels {
					apiCR.ObjectMeta.Labels[key] = value
				}
				crUpdateError := apiReconciler.client.Update(ctx, &apiCR)
				if crUpdateError != nil {
					loggers.LoggerAPKOperator.Errorf("Error while updating the API CR for api labels. Error: %+v", crUpdateError)
				}
			} else {
				loggers.LoggerAPKOperator.Errorf("Error while loading api: %s/%s, Error: %v", labelUpdate.Name, labelUpdate.Namespace, err)
			}
		}
	}
}

func (apiReconciler *APIReconciler) handleOwnerReference(ctx context.Context, obj k8client.Object, apiRequests *[]reconcile.Request) {
	apis := []dpv1alpha3.API{}
	for _, req := range *apiRequests {
		var apiCR dpv1alpha3.API
		if err := apiReconciler.client.Get(ctx, req.NamespacedName, &apiCR); err == nil {
			apis = append(apis, apiCR)
		} else {
			loggers.LoggerAPKOperator.Errorf("Error while loading api: %+v, Error: %v", req, err)
		}
	}
	// Prepare owner references for the route
	preparedOwnerReferences := prepareOwnerReference(apis)
	// Decide whether we need an update
	updateRequired := false
	if len(obj.GetOwnerReferences()) != len(preparedOwnerReferences) {
		updateRequired = true
	} else {
		for _, ref := range preparedOwnerReferences {
			_, found := common.FindElement(obj.GetOwnerReferences(), func(refLocal metav1.OwnerReference) bool {
				return refLocal.UID == ref.UID && refLocal.Name == ref.Name && refLocal.APIVersion == ref.APIVersion && refLocal.Kind == ref.Kind
			})
			if !found {
				updateRequired = true
				break
			}
		}
	}
	if updateRequired {
		obj.SetOwnerReferences(preparedOwnerReferences)
		utils.UpdateCR(ctx, apiReconciler.client, obj)
	}
}

func prepareOwnerReference(apiItems []dpv1alpha3.API) []metav1.OwnerReference {
	ownerReferences := []metav1.OwnerReference{}
	uidMap := make(map[string]bool)
	for _, ref := range apiItems {
		if ref.Spec.SystemAPI {
			continue
		}
		if _, exists := uidMap[string(ref.UID)]; !exists {
			ownerReferences = append(ownerReferences, metav1.OwnerReference{
				APIVersion: ref.APIVersion,
				Kind:       ref.Kind,
				Name:       ref.Name,
				UID:        ref.UID,
			})
			uidMap[string(ref.UID)] = true
		}
	}
	return ownerReferences
}

func (apiReconciler *APIReconciler) convertAPIStateToAPICp(ctx context.Context, apiState synchronizer.APIState, apiHash string, prodAIRL *dpv1alpha3.AIRateLimitPolicy, sandAIRL *dpv1alpha3.AIRateLimitPolicy) controlplane.APICPEvent {
	apiCPEvent := controlplane.APICPEvent{}
	spec := apiState.APIDefinition.Spec
	configMap := &corev1.ConfigMap{}
	apiDef := ""
	if spec.DefinitionFileRef != "" {
		err := apiReconciler.client.Get(ctx, types.NamespacedName{Namespace: apiState.APIDefinition.Namespace, Name: spec.DefinitionFileRef}, configMap)
		if err != nil {
			loggers.LoggerAPKOperator.Errorf("Error while loading config map for the api definition: %+v, Error: %v", types.NamespacedName{Namespace: apiState.APIDefinition.Namespace, Name: spec.DefinitionFileRef}, err)
		} else {
			for _, val := range configMap.BinaryData {
				buf := bytes.NewReader(val)
				r, err := gzip.NewReader(buf)
				if err != nil {
					loggers.LoggerAPKOperator.Errorf("Error creating gzip reader. Error: %+v", err)
					continue
				}
				defer r.Close()
				decompressed, err := ioutil.ReadAll(r)
				if err != nil {
					loggers.LoggerAPKOperator.Errorf("Error reading decompressed data. Error: %+v", err)
					continue
				}
				apiDef = string(decompressed)
			}
		}
	}
	apiUUID, apiUUIDExists := apiState.APIDefinition.ObjectMeta.Labels["apiUUID"]
	if !apiUUIDExists {
		apiUUID = spec.APIName
	}
	revisionID, revisionIDExists := apiState.APIDefinition.ObjectMeta.Labels["revisionID"]
	if !revisionIDExists {
		revisionID = "0"
	}
	properties := make(map[string]string)
	for _, val := range spec.APIProperties {
		properties[val.Name] = val.Value
	}
	prodEndpoint, sandEndpoint, endpointProtocol, prodAPIKeyName, prodAPIKeyIn, prodAPIKeyValue, prodBasicUsername, prodBasicPassword, sandAPIKeyName, sandAPIKeyIn, sandAPIKeyValue, sandBasicUsername, sandBasicPassword, prodEndpointSecurityType, sandEndpointSecurityType, prodEndpointSecurityEnabled, sandEndpointSecurityEnabled := findProdSandEndpoints(&apiState)
	corsPolicy := pickOneCorsForCP(&apiState)
	vhost := getProdVhost(&apiState)
	sandVhost := geSandVhost(&apiState)
	securityScheme, authHeader, apiKeyHeader := prepareSecuritySchemeForCP(&apiState)
	operations := prepareOperations(&apiState)
	var sandAIRLToAgent *controlplane.AIRL
	var prodAIRLToAgent *controlplane.AIRL
	if prodAIRL != nil {
		var promptTC, completionTC, totalTC, requestC *uint32
		var timeUnit string
		if prodAIRL.Spec.Override.TokenCount != nil {
			promptTC = &prodAIRL.Spec.Override.TokenCount.RequestTokenCount
			completionTC = &prodAIRL.Spec.Override.TokenCount.ResponseTokenCount
			totalTC = &prodAIRL.Spec.Override.TokenCount.TotalTokenCount
			timeUnit = prodAIRL.Spec.Override.TokenCount.Unit
		}
		if prodAIRL.Spec.Override.RequestCount != nil {
			timeUnit = prodAIRL.Spec.Override.RequestCount.Unit
			requestC = &prodAIRL.Spec.Override.RequestCount.RequestsPerUnit
		}
		prodAIRLToAgent = &controlplane.AIRL{
			PromptTokenCount:     promptTC,
			CompletionTokenCount: completionTC,
			TotalTokenCount:      totalTC,
			TimeUnit:             timeUnit,
			RequestCount:         requestC,
		}
	}
	if sandAIRL != nil {
		var promptTC, completionTC, totalTC, requestC *uint32
		var timeUnit string
		if sandAIRL.Spec.Override.TokenCount != nil {
			promptTC = &sandAIRL.Spec.Override.TokenCount.RequestTokenCount
			completionTC = &sandAIRL.Spec.Override.TokenCount.ResponseTokenCount
			totalTC = &sandAIRL.Spec.Override.TokenCount.TotalTokenCount
			timeUnit = sandAIRL.Spec.Override.TokenCount.Unit
		}
		if sandAIRL.Spec.Override.RequestCount != nil {
			timeUnit = sandAIRL.Spec.Override.RequestCount.Unit
			requestC = &sandAIRL.Spec.Override.RequestCount.RequestsPerUnit
		}
		sandAIRLToAgent = &controlplane.AIRL{
			PromptTokenCount:     promptTC,
			CompletionTokenCount: completionTC,
			TotalTokenCount:      totalTC,
			TimeUnit:             timeUnit,
			RequestCount:         requestC,
		}
	}
	prodEndpointSecurity := controlplane.EndpointSecurity{
		Enabled:       prodEndpointSecurityEnabled,
		SecurityType:  prodEndpointSecurityType,
		APIKeyName:    prodAPIKeyName,
		APIKeyIn:      prodAPIKeyIn,
		APIKeyValue:   prodAPIKeyValue,
		BasicUsername: prodBasicUsername,
		BasicPassword: prodBasicPassword,
	}
	sandEndpointSecurity := controlplane.EndpointSecurity{
		Enabled:       sandEndpointSecurityEnabled,
		SecurityType:  sandEndpointSecurityType,
		APIKeyName:    sandAPIKeyName,
		APIKeyIn:      sandAPIKeyIn,
		APIKeyValue:   sandAPIKeyValue,
		BasicUsername: sandBasicUsername,
		BasicPassword: sandBasicPassword,
	}
	subType := "DEFAULT"
	aiConfiguration := controlplane.AIConfiguration{}
	loggers.LoggerAPKOperator.Debugf("AI Provider in state: %+v", apiState.AIProvider)
	if apiState.AIProvider != nil {
		loggers.LoggerAPKOperator.Debugf("AIProvider is found")
		subType = "AIAPI"
		aiConfiguration = controlplane.AIConfiguration{
			LLMProviderID:         apiState.AIProvider.ObjectMeta.Name,
			LLMProviderName:       apiState.AIProvider.Spec.ProviderName,
			LLMProviderAPIVersion: apiState.AIProvider.Spec.ProviderAPIVersion,
		}
	}
	loggers.LoggerAPKOperator.Debugf("Resolved aiConfiguration: %+v", aiConfiguration)

	api := controlplane.API{
		APIName:              spec.APIName,
		APIVersion:           spec.APIVersion,
		IsDefaultVersion:     spec.IsDefaultVersion,
		APIType:              spec.APIType,
		APISubType:           subType,
		BasePath:             spec.BasePath,
		Organization:         spec.Organization,
		Environment:          spec.Environment,
		SystemAPI:            spec.SystemAPI,
		Definition:           apiDef,
		APIUUID:              apiUUID,
		RevisionID:           revisionID,
		APIProperties:        properties,
		ProdEndpoint:         prodEndpoint,
		SandEndpoint:         sandEndpoint,
		ProdEndpointSecurity: prodEndpointSecurity,
		SandEndpointSecurity: sandEndpointSecurity,
		EndpointProtocol:     endpointProtocol,
		CORSPolicy:           corsPolicy,
		Vhost:                vhost,
		SandVhost:            sandVhost,
		SecurityScheme:       securityScheme,
		AuthHeader:           authHeader,
		Operations:           operations,
		APIHash:              apiHash,
		APIKeyHeader:         apiKeyHeader,
		SandAIRL:             sandAIRLToAgent,
		ProdAIRL:             prodAIRLToAgent,
		AIConfiguration:      aiConfiguration,
	}
	apiCPEvent.API = api
	apiCPEvent.CRName = apiState.APIDefinition.ObjectMeta.Name
	apiCPEvent.CRNamespace = apiState.APIDefinition.ObjectMeta.Namespace
	return apiCPEvent

}

func (apiReconciler *APIReconciler) validateRouteExtRefs(apiState *synchronizer.APIState) error {
	extRefs := []*gwapiv1.LocalObjectReference{}
	if apiState.ProdHTTPRoute != nil {
		for _, httpRoute := range apiState.ProdHTTPRoute.HTTPRoutePartitions {
			for _, rule := range httpRoute.Spec.Rules {
				for _, filter := range rule.Filters {
					extRefs = append(extRefs, filter.ExtensionRef)
				}
			}
		}
	}
	if apiState.SandHTTPRoute != nil {
		for _, httpRoute := range apiState.SandHTTPRoute.HTTPRoutePartitions {
			for _, rule := range httpRoute.Spec.Rules {
				for _, filter := range rule.Filters {
					extRefs = append(extRefs, filter.ExtensionRef)
				}
			}
		}
	}
	if apiState.ProdGQLRoute != nil {
		for _, gql := range apiState.ProdGQLRoute.GQLRoutePartitions {
			for _, rule := range gql.Spec.Rules {
				for _, filter := range rule.Filters {
					extRefs = append(extRefs, filter.ExtensionRef)
				}
			}
		}
	}
	if apiState.SandGQLRoute != nil {
		for _, gql := range apiState.SandGQLRoute.GQLRoutePartitions {
			for _, rule := range gql.Spec.Rules {
				for _, filter := range rule.Filters {
					extRefs = append(extRefs, filter.ExtensionRef)
				}
			}
		}
	}
	for _, extRef := range extRefs {
		if extRef != nil {
			extKind := string(extRef.Kind)
			key := types.NamespacedName{Namespace: string(apiState.APIDefinition.Namespace), Name: string(extRef.Name)}.String()
			if extKind == "APIPolicy" {
				_, found := apiState.ResourceAPIPolicies[key]
				if !found {
					return fmt.Errorf("apipolicy not added to the ResourceAPIPolicies map yet. Key: %s", key)
				}
			}
			if extKind == "RateLimitPolicy" {
				_, found := apiState.ResourceRateLimitPolicies[key]
				if !found {
					return fmt.Errorf("ratelimitPolicy not added to the ResourceRateLimitPolicies map yet. Key: %s", key)
				}
			}
			if extKind == "Authentication" {
				_, found := apiState.ResourceAuthentications[key]
				if !found {
					return fmt.Errorf("authentication not added to the resourse Authentication map yet. Key: %s", key)
				}
			}
		}
	}
	return nil
}

func (apiReconciler *APIReconciler) getAPIHash(apiState *synchronizer.APIState) string {
	getUniqueID := func(obj interface{}, fields ...string) string {
		defer func() {
			if r := recover(); r != nil {
				loggers.LoggerAPK.Infof("Error occured while extracting values using reflection. Error: %+v", r)
			}
		}()
		if obj == nil {
			return "nil"
		}
		var sb strings.Builder
		objValue := reflect.ValueOf(obj)
		if objValue.Kind() == reflect.Ptr {
			objValue = objValue.Elem()
		}
		for _, field := range fields {
			fieldNames := strings.Split(field, ".")
			name1 := fieldNames[0]
			name2 := fieldNames[1]
			if objValue.IsValid() && objValue.FieldByName(name1).IsValid() {
				if objValue.FieldByName(name1).FieldByName(name2).IsValid() {
					v := objValue.FieldByName(name1).FieldByName(name2)
					switch v.Kind() {
					case reflect.String:
						sb.WriteString(v.String())
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						sb.WriteString(strconv.FormatInt(v.Int(), 10))
					}
				}
			}
		}
		return sb.String()
	}

	uniqueIDs := make([]string, 0)
	uniqueIDs = append(uniqueIDs, getUniqueID(apiState.APIDefinition, "ObjectMeta.Name", "ObjectMeta.Namespace", "ObjectMeta.Generation"))
	uniqueIDs = append(uniqueIDs, getUniqueID(apiState.SandAIRL, "ObjectMeta.Name", "ObjectMeta.Namespace", "ObjectMeta.Generation"))
	uniqueIDs = append(uniqueIDs, getUniqueID(apiState.ProdAIRL, "ObjectMeta.Name", "ObjectMeta.Namespace", "ObjectMeta.Generation"))
	for _, auth := range apiState.Authentications {
		uniqueIDs = append(uniqueIDs, getUniqueID(auth, "ObjectMeta.Name", "ObjectMeta.Namespace", "ObjectMeta.Generation"))
	}
	for _, arl := range apiState.RateLimitPolicies {
		uniqueIDs = append(uniqueIDs, getUniqueID(arl, "ObjectMeta.Name", "ObjectMeta.Namespace", "ObjectMeta.Generation"))
	}
	for _, ra := range apiState.ResourceAuthentications {
		uniqueIDs = append(uniqueIDs, getUniqueID(ra, "ObjectMeta.Name", "ObjectMeta.Namespace", "ObjectMeta.Generation"))
	}
	for _, rrl := range apiState.ResourceRateLimitPolicies {
		uniqueIDs = append(uniqueIDs, getUniqueID(rrl, "ObjectMeta.Name", "ObjectMeta.Namespace", "ObjectMeta.Generation"))
	}
	for _, ral := range apiState.ResourceAPIPolicies {
		uniqueIDs = append(uniqueIDs, getUniqueID(ral, "ObjectMeta.Name", "ObjectMeta.Namespace", "ObjectMeta.Generation"))
	}
	for _, ap := range apiState.APIPolicies {
		uniqueIDs = append(uniqueIDs, getUniqueID(ap, "ObjectMeta.Name", "ObjectMeta.Namespace", "ObjectMeta.Generation"))
	}
	for _, ism := range apiState.InterceptorServiceMapping {
		uniqueIDs = append(uniqueIDs, getUniqueID(ism, "ObjectMeta.Name", "ObjectMeta.Namespace", "ObjectMeta.Generation"))
	}
	for _, bjm := range apiState.BackendJWTMapping {
		uniqueIDs = append(uniqueIDs, getUniqueID(bjm, "ObjectMeta.Name", "ObjectMeta.Namespace", "ObjectMeta.Generation"))
	}
	if apiState.ProdHTTPRoute != nil {
		for _, phr := range apiState.ProdHTTPRoute.HTTPRoutePartitions {
			uniqueIDs = append(uniqueIDs, getUniqueID(phr, "ObjectMeta.Name", "ObjectMeta.Namespace", "ObjectMeta.Generation"))
		}
		for _, backend := range apiState.ProdHTTPRoute.BackendMapping {
			uniqueIDs = append(uniqueIDs, getUniqueID(backend.Backend, "ObjectMeta.Name", "ObjectMeta.Namespace", "ObjectMeta.Generation"))
		}
	}
	if apiState.SandHTTPRoute != nil {
		for _, shr := range apiState.SandHTTPRoute.HTTPRoutePartitions {
			uniqueIDs = append(uniqueIDs, getUniqueID(shr, "ObjectMeta.Name", "ObjectMeta.Namespace", "ObjectMeta.Generation"))
		}
		for _, backend := range apiState.SandHTTPRoute.BackendMapping {
			uniqueIDs = append(uniqueIDs, getUniqueID(backend.Backend, "ObjectMeta.Name", "ObjectMeta.Namespace", "ObjectMeta.Generation"))
		}
	}
	if apiState.ProdGQLRoute != nil {
		for _, pgqr := range apiState.ProdGQLRoute.GQLRoutePartitions {
			uniqueIDs = append(uniqueIDs, getUniqueID(pgqr, "ObjectMeta.Name", "ObjectMeta.Namespace", "ObjectMeta.Generation"))
		}
		for _, backend := range apiState.ProdGQLRoute.BackendMapping {
			uniqueIDs = append(uniqueIDs, getUniqueID(backend.Backend, "ObjectMeta.Name", "ObjectMeta.Namespace", "ObjectMeta.Generation"))
		}
	}
	if apiState.SandGQLRoute != nil {
		for _, sgqr := range apiState.SandGQLRoute.GQLRoutePartitions {
			uniqueIDs = append(uniqueIDs, getUniqueID(sgqr, "ObjectMeta.Name", "ObjectMeta.Namespace", "ObjectMeta.Generation"))
		}
		for _, backend := range apiState.SandGQLRoute.BackendMapping {
			uniqueIDs = append(uniqueIDs, getUniqueID(backend.Backend, "ObjectMeta.Name", "ObjectMeta.Namespace", "ObjectMeta.Generation"))
		}
	}

	sort.Strings(uniqueIDs)
	joinedUniqueIDs := strings.Join(uniqueIDs, "")
	mutualSSLUniqueID := ""
	if apiState.MutualSSL != nil {
		mutualSSLUniqueID += strconv.FormatBool(apiState.MutualSSL.Disabled) + apiState.MutualSSL.Required + strings.Join(apiState.MutualSSL.ClientCertificates, "")
	}
	joinedUniqueIDs = joinedUniqueIDs + strconv.FormatBool(apiState.SubscriptionValidation) + mutualSSLUniqueID
	hash := sha256.Sum256([]byte(joinedUniqueIDs))
	hashedString := hex.EncodeToString(hash[:])
	truncatedHash := hashedString[:62]
	loggers.LoggerAPK.Debugf("Prepared unique string for api %s/%s: %s, Prepared hash: %s, Truncatd hash to store: %s", apiState.APIDefinition.ObjectMeta.Name,
		apiState.APIDefinition.ObjectMeta.Namespace, joinedUniqueIDs, hashedString, truncatedHash)
	return truncatedHash
}

func findProdSandEndpoints(apiState *synchronizer.APIState) (string, string, string, string, string, string, string, string, string, string, string, string, string, string, string, bool, bool) {
	prodEndpoint := ""
	sandEndpoint := ""
	endpointProtocol := ""
	prodEndpointSecurityType := ""
	prodEndpointSecurityEnabled := false
	prodAPIKeyName := ""
	prodAPIKeyIn := ""
	prodAPIKeyValue := ""
	prodBasicUsername := ""
	prodBasicPassword := ""
	sandEndpointSecurityType := ""
	sandEndpointSecurityEnabled := false
	sandAPIKeyName := ""
	sandAPIKeyIn := ""
	sandAPIKeyValue := ""
	sandBasicUsername := ""
	sandBasicPassword := ""

	if apiState.ProdHTTPRoute != nil {
		for _, backend := range apiState.ProdHTTPRoute.BackendMapping {
			if len(backend.Backend.Spec.Services) > 0 {
				prodEndpoint = fmt.Sprintf("%s:%d", backend.Backend.Spec.Services[0].Host, backend.Backend.Spec.Services[0].Port)
				endpointProtocol = string(backend.Backend.Spec.Protocol)
			}
			if backend.Security.Basic.Username != "" && backend.Security.Basic.Password != "" {
				prodEndpointSecurityEnabled = true
				prodEndpointSecurityType = "basic"
				prodBasicUsername = backend.Security.Basic.Username
				prodBasicPassword = backend.Security.Basic.Password
			}
			if backend.Security.APIKey.Name != "" && backend.Security.APIKey.Value != "" && backend.Security.APIKey.In != "" {
				prodEndpointSecurityEnabled = true
				prodEndpointSecurityType = "apikey"
				prodAPIKeyName = backend.Security.APIKey.Name
				prodAPIKeyIn = string(backend.Security.APIKey.In)
				prodAPIKeyValue = backend.Security.APIKey.Value
			}
		}
	}
	if apiState.SandHTTPRoute != nil {
		for _, backend := range apiState.SandHTTPRoute.BackendMapping {
			if len(backend.Backend.Spec.Services) > 0 {
				sandEndpoint = fmt.Sprintf("%s:%d", backend.Backend.Spec.Services[0].Host, backend.Backend.Spec.Services[0].Port)
				endpointProtocol = string(backend.Backend.Spec.Protocol)
			}
			if backend.Security.Basic.Username != "" && backend.Security.Basic.Password != "" {
				sandEndpointSecurityEnabled = true
				sandEndpointSecurityType = "basic"
				sandBasicUsername = backend.Security.Basic.Username
				sandBasicPassword = backend.Security.Basic.Password
			}
			if backend.Security.APIKey.Name != "" && backend.Security.APIKey.Value != "" && backend.Security.APIKey.In != "" {
				sandEndpointSecurityEnabled = true
				sandEndpointSecurityType = "apikey"
				sandAPIKeyName = backend.Security.APIKey.Name
				sandAPIKeyIn = string(backend.Security.APIKey.In)
				sandAPIKeyValue = backend.Security.APIKey.Value
			}
		}
	}
	if apiState.ProdGQLRoute != nil {
		for _, backend := range apiState.ProdGQLRoute.BackendMapping {
			if len(backend.Backend.Spec.Services) > 0 {
				prodEndpoint = fmt.Sprintf("%s:%d", backend.Backend.Spec.Services[0].Host, backend.Backend.Spec.Services[0].Port)
				endpointProtocol = string(backend.Backend.Spec.Protocol)
			}
			if backend.Security.Basic.Username != "" && backend.Security.Basic.Password != "" {
				prodEndpointSecurityEnabled = true
				prodEndpointSecurityType = "basic"
				prodBasicUsername = backend.Security.Basic.Username
				prodBasicPassword = backend.Security.Basic.Password
			}
			if backend.Security.APIKey.Name != "" && backend.Security.APIKey.Value != "" && backend.Security.APIKey.In != "" {
				prodEndpointSecurityEnabled = true
				prodEndpointSecurityType = "apikey"
				prodAPIKeyName = backend.Security.APIKey.Name
				prodAPIKeyIn = string(backend.Security.APIKey.In)
				prodAPIKeyValue = backend.Security.APIKey.Value
			}
		}
	}
	if apiState.SandGQLRoute != nil {
		for _, backend := range apiState.SandGQLRoute.BackendMapping {
			if len(backend.Backend.Spec.Services) > 0 {
				sandEndpoint = fmt.Sprintf("%s:%d", backend.Backend.Spec.Services[0].Host, backend.Backend.Spec.Services[0].Port)
				endpointProtocol = string(backend.Backend.Spec.Protocol)
			}
			if backend.Security.Basic.Username != "" && backend.Security.Basic.Password != "" {
				sandEndpointSecurityEnabled = true
				sandEndpointSecurityType = "basic"
				sandBasicUsername = backend.Security.Basic.Username
				sandBasicPassword = backend.Security.Basic.Password
			}
			if backend.Security.APIKey.Name != "" && backend.Security.APIKey.Value != "" && backend.Security.APIKey.In != "" {
				sandEndpointSecurityEnabled = true
				sandEndpointSecurityType = "apikey"
				sandAPIKeyName = backend.Security.APIKey.Name
				sandAPIKeyIn = string(backend.Security.APIKey.In)
				sandAPIKeyValue = backend.Security.APIKey.Value
			}
		}
	}
	return prodEndpoint, sandEndpoint, endpointProtocol, prodAPIKeyName, prodAPIKeyIn, prodAPIKeyValue, prodBasicUsername, prodBasicPassword, sandAPIKeyName, sandAPIKeyIn, sandAPIKeyValue, sandBasicUsername, sandBasicPassword, prodEndpointSecurityType, sandEndpointSecurityType, prodEndpointSecurityEnabled, sandEndpointSecurityEnabled
}

func pickOneCorsForCP(apiState *synchronizer.APIState) *controlplane.CORSPolicy {
	apiPolicies := []dpv1alpha4.APIPolicy{}
	for _, apiPolicy := range apiState.APIPolicies {
		apiPolicies = append(apiPolicies, apiPolicy)
	}
	for _, apiPolicy := range apiState.ResourceAPIPolicies {
		apiPolicies = append(apiPolicies, apiPolicy)
	}
	for _, apiPolicy := range apiPolicies {
		corsPolicy := dpv1alpha4.CORSPolicy{}
		found := false
		if apiPolicy.Spec.Override != nil && apiPolicy.Spec.Override.CORSPolicy != nil {
			corsPolicy = *apiPolicy.Spec.Override.CORSPolicy
			found = true
		} else if apiPolicy.Spec.Default != nil && apiPolicy.Spec.Default.CORSPolicy != nil {
			corsPolicy = *apiPolicy.Spec.Default.CORSPolicy
			found = true
		}
		if found {
			modifiedCors := controlplane.CORSPolicy{}
			modifiedCors.AccessControlAllowCredentials = corsPolicy.AccessControlAllowCredentials
			modifiedCors.AccessControlAllowHeaders = corsPolicy.AccessControlAllowHeaders
			modifiedCors.AccessControlAllowOrigins = corsPolicy.AccessControlAllowOrigins
			modifiedCors.AccessControlExposeHeaders = corsPolicy.AccessControlExposeHeaders
			modifiedCors.AccessControlMaxAge = corsPolicy.AccessControlMaxAge
			modifiedCors.AccessControlAllowMethods = corsPolicy.AccessControlAllowMethods
			return &modifiedCors
		}
	}
	return nil
}

func getProdVhost(apiState *synchronizer.APIState) string {
	if apiState.ProdHTTPRoute != nil {
		for _, httpRoute := range apiState.ProdHTTPRoute.HTTPRoutePartitions {
			if len(httpRoute.Spec.Hostnames) > 0 {
				return string(httpRoute.Spec.Hostnames[0])
			}
		}
	}
	if apiState.ProdGQLRoute != nil {
		for _, gql := range apiState.ProdGQLRoute.GQLRoutePartitions {
			if len(gql.Spec.Hostnames) > 0 {
				return string(gql.Spec.Hostnames[0])
			}
		}
	}
	return "default.gw.wso2.com"
}

func geSandVhost(apiState *synchronizer.APIState) string {
	if apiState.SandHTTPRoute != nil {
		for _, httpRoute := range apiState.SandHTTPRoute.HTTPRoutePartitions {
			if len(httpRoute.Spec.Hostnames) > 0 {
				return string(httpRoute.Spec.Hostnames[0])
			}
		}
	}
	if apiState.SandGQLRoute != nil {
		for _, gql := range apiState.SandGQLRoute.GQLRoutePartitions {
			if len(gql.Spec.Hostnames) > 0 {
				return string(gql.Spec.Hostnames[0])
			}
		}
	}
	return "sandbox.default.gw.wso2.com"
}

func prepareSecuritySchemeForCP(apiState *synchronizer.APIState) ([]string, string, string) {
	var pickedAuth *dpv1alpha2.Authentication
	authHeader := "Authorization"
	apiKeyHeader := "ApiKey"
	for _, auth := range apiState.Authentications {
		pickedAuth = &auth
		break
	}
	if pickedAuth != nil {
		var authSpec *dpv1alpha2.AuthSpec
		if pickedAuth.Spec.Override != nil {
			authSpec = pickedAuth.Spec.Override
		} else {
			authSpec = pickedAuth.Spec.Default
		}
		if authSpec != nil {
			if authSpec.AuthTypes != nil {
				authSchemes := []string{}
				isAuthMandatory := false
				isMTLSMandatory := false
				if authSpec.AuthTypes.OAuth2.Required == "mandatory" {
					isAuthMandatory = true
				}
				if !authSpec.AuthTypes.OAuth2.Disabled {
					authSchemes = append(authSchemes, "oauth2")
					if authSpec.AuthTypes.OAuth2.Header != "" {
						authHeader = authSpec.AuthTypes.OAuth2.Header
					}
				}
				if authSpec.AuthTypes.MutualSSL != nil && authSpec.AuthTypes.MutualSSL.Required == "mandatory" {
					isMTLSMandatory = true
				}
				if authSpec.AuthTypes.MutualSSL != nil && !authSpec.AuthTypes.MutualSSL.Disabled {
					authSchemes = append(authSchemes, "mutualssl")
					if isMTLSMandatory {
						authSchemes = append(authSchemes, "mutualssl_mandatory")
					} else {
						authSchemes = append(authSchemes, "mutualssl_optional")
					}
				}
				if len(authSpec.AuthTypes.APIKey.Keys) > 0 {
					authSchemes = append(authSchemes, "api_key")
					for _, apiKey := range authSpec.AuthTypes.APIKey.Keys {
						if apiKey.In == "Header" {
							apiKeyHeader = apiKey.Name
						}
					}
				}
				if isAuthMandatory {
					authSchemes = append(authSchemes, "oauth_basic_auth_api_key_mandatory")
				} else {
					authSchemes = append(authSchemes, "oauth_basic_auth_api_key_optional")
				}
				return authSchemes, authHeader, apiKeyHeader
			}
		}
	}
	return []string{"oauth2", "oauth_basic_auth_api_key_mandatory"}, authHeader, apiKeyHeader
}

func prepareOperations(apiState *synchronizer.APIState) []controlplane.Operation {
	operations := []controlplane.Operation{}
	if apiState.ProdHTTPRoute != nil && apiState.ProdHTTPRoute.HTTPRouteCombined != nil {
		for _, rule := range apiState.ProdHTTPRoute.HTTPRouteCombined.Spec.Rules {
			scopes := []string{}
			requestAddHeaders := []controlplane.Header{}
			responseAddHeaders := []controlplane.Header{}
			requestRemoveHeaders := []string{}
			responseRemoveHeaders := []string{}
			for _, filter := range rule.Filters {
				if filter.ExtensionRef != nil && filter.ExtensionRef.Kind == "Scope" {
					scope, found := apiState.ProdHTTPRoute.Scopes[types.NamespacedName{Namespace: apiState.APIDefinition.ObjectMeta.Namespace, Name: string(filter.ExtensionRef.Name)}.String()]
					if found {
						scopes = append(scopes, scope.Spec.Names...)
					}
				}

				if filter.RequestHeaderModifier != nil {
					requestHeaderModifier := filter.RequestHeaderModifier
					for _, addHeader := range requestHeaderModifier.Add {
						requestAddHeaders = append(requestAddHeaders, controlplane.Header{Name: string(addHeader.Name), Value: string(addHeader.Value)})
					}
					for _, setHeader := range requestHeaderModifier.Set {
						requestAddHeaders = append(requestAddHeaders, controlplane.Header{Name: string(setHeader.Name), Value: string(setHeader.Value)})
					}
					for _, removeHeader := range requestHeaderModifier.Remove {
						requestRemoveHeaders = append(requestRemoveHeaders, removeHeader)
					}
				}

				if filter.ResponseHeaderModifier != nil {
					responseHeaderModifier := filter.ResponseHeaderModifier
					for _, addHeader := range responseHeaderModifier.Add {
						responseAddHeaders = append(responseAddHeaders, controlplane.Header{Name: string(addHeader.Name), Value: string(addHeader.Value)})
					}
					for _, setHeader := range responseHeaderModifier.Set {
						responseAddHeaders = append(responseAddHeaders, controlplane.Header{Name: string(setHeader.Name), Value: string(setHeader.Value)})
					}
					for _, removeHeader := range responseHeaderModifier.Remove {
						responseRemoveHeaders = append(responseRemoveHeaders, removeHeader)
					}
				}
			}

			for _, match := range rule.Matches {
				path := "/"
				verb := "GET"
				if match.Path != nil && match.Path.Value != nil {
					path = *match.Path.Value
				}
				if match.Method != nil {
					verb = string(*match.Method)
				}
				if match.Path.Type == nil || *match.Path.Type == gwapiv1.PathMatchPathPrefix {
					path = path + "*"
				}
				path = "^" + path + "$"
				operations = append(operations, controlplane.Operation{
					Path:   path,
					Verb:   verb,
					Scopes: scopes,
					Headers: controlplane.Headers{
						RequestHeaders: controlplane.HeaderModifier{
							AddHeaders:    requestAddHeaders,
							RemoveHeaders: requestRemoveHeaders,
						},
						ResponseHeaders: controlplane.HeaderModifier{
							AddHeaders:    responseAddHeaders,
							RemoveHeaders: responseRemoveHeaders,
						},
					},
				})
			}
		}
	}
	if apiState.ProdGQLRoute != nil && apiState.ProdGQLRoute.GQLRouteCombined != nil {
		for _, rule := range apiState.ProdGQLRoute.GQLRouteCombined.Spec.Rules {
			scopes := []string{}
			for _, filter := range rule.Filters {
				if filter.ExtensionRef.Kind == "Scope" {
					scope, found := apiState.ProdGQLRoute.Scopes[types.NamespacedName{Namespace: apiState.APIDefinition.ObjectMeta.Namespace, Name: string(filter.ExtensionRef.Name)}.String()]
					if found {
						scopes = append(scopes, scope.Spec.Names...)
					}
				}
			}
			for _, match := range rule.Matches {
				path := ""
				verb := "QUERY"
				if match.Path != nil {
					path = *match.Path
				}
				if match.Type != nil {
					verb = string(*match.Type)
				}
				operations = append(operations, controlplane.Operation{Path: path, Verb: verb, Scopes: scopes})
			}
		}
	}

	return operations
}
