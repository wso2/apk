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

package controllers

import (
	"context"
	"errors"
	"fmt"

	"github.com/wso2/apk/adapter/config"
	"github.com/wso2/apk/adapter/internal/loggers"
	"github.com/wso2/apk/adapter/internal/operator/constants"
	"github.com/wso2/apk/adapter/internal/operator/synchronizer"
	"github.com/wso2/apk/adapter/internal/operator/utils"
	"github.com/wso2/apk/adapter/pkg/logging"
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
	"sigs.k8s.io/controller-runtime/pkg/client"

	dpv1alpha1 "github.com/wso2/apk/adapter/internal/operator/apis/dp/v1alpha1"
)

const httpRouteAPIIndex = "httpRouteAPIIndex"

// APIReconciler reconciles a API object
type APIReconciler struct {
	client client.Client
	ods    *synchronizer.OperatorDataStore
	ch     *chan synchronizer.APIEvent
}

// NewAPIController creates a new API controller instance. API Controllers watches for dpv1alpha1.API and gwapiv1b1.HTTPRoute.
func NewAPIController(mgr manager.Manager, operatorDataStore *synchronizer.OperatorDataStore, ch *chan synchronizer.APIEvent) error {
	r := &APIReconciler{
		client: mgr.GetClient(),
		ods:    operatorDataStore,
		ch:     ch,
	}
	c, err := controller.New(constants.APIController, mgr, controller.Options{Reconciler: r})
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprintf("error creating API controller:%v", err),
			Severity:  logging.BLOCKER,
			ErrorCode: 2600,
		})
		return err
	}
	ctx := context.Background()

	conf := config.ReadConfigs()
	predicates := []predicate.Predicate{predicate.NewPredicateFuncs(utils.FilterByNamespaces(conf.Adapter.Operator.Namespaces))}

	if err := c.Watch(&source.Kind{Type: &dpv1alpha1.API{}}, &handler.EnqueueRequestForObject{},
		predicates...); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprintf("error watching API resources: %v", err),
			Severity:  logging.BLOCKER,
			ErrorCode: 2601,
		})
		return err
	}
	if err := addAPIIndexers(ctx, mgr); err != nil {
		return err
	}

	if err := c.Watch(&source.Kind{Type: &gwapiv1b1.HTTPRoute{}}, handler.EnqueueRequestsFromMapFunc(r.getAPIForHTTPRoute),
		predicates...); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprintf("error watching HTTPRoute resources: %v", err),
			Severity:  logging.BLOCKER,
			ErrorCode: 2602,
		})
		return err
	}
	loggers.LoggerAPKOperator.Info("API Controller successfully started. Watching API Objects....")
	return nil
}

//+kubebuilder:rbac:groups=*,resources=apis,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=*,resources=apis/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=*,resources=apis/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (apiReconciler *APIReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	// 1. Check whether the API CR exist, if not consider as a DELETE event.
	var apiDef dpv1alpha1.API
	if err := apiReconciler.client.Get(ctx, req.NamespacedName, &apiDef); err != nil {
		apiState, found := apiReconciler.ods.APIStore[req.NamespacedName]
		if found && k8error.IsNotFound(err) {
			event := *apiState
			// The api doesn't exist in the api Cache, remove it
			delete(apiReconciler.ods.APIStore, req.NamespacedName)
			loggers.LoggerAPKOperator.Infof("Delete event has received for API : %s, hence deleted from API cache", req.NamespacedName.String())
			*apiReconciler.ch <- synchronizer.APIEvent{EventType: constants.Delete, Event: event}
			return ctrl.Result{}, nil
		}
		loggers.LoggerAPKOperator.ErrorC(logging.ErrorDetails{
			Message: fmt.Sprintf("Api CR related to the reconcile request with key: %s returned error."+
				" Assuming API is already deleted, hence ignoring the error : %v",
				req.NamespacedName.String(), err),
			Severity:  logging.TRIVIAL,
			ErrorCode: 2604,
		})
		return ctrl.Result{}, nil
	}

	// 2. Handle HTTPRoute validation
	prodHTTPRoute, sandHTTPRoute, err := validateHTTPRouteRefs(ctx, apiReconciler.client, req.Namespace,
		apiDef.Spec.ProdHTTPRouteRef, apiDef.Spec.SandHTTPRouteRef)
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprintf("error validating HttpRoute CRs: %v", err),
			Severity:  logging.TRIVIAL,
			ErrorCode: 2604,
		})
		return ctrl.Result{}, err
	}
	loggers.LoggerAPKOperator.Debugf("HTTPRoute validation has passed for API CR %s", req.NamespacedName.String())

	// 3. Check whether the Operator Data store contains the received API.
	cachedAPI, found := apiReconciler.ods.GetAPI(utils.NamespacedName(&apiDef))

	if !found {
		apiState := apiReconciler.ods.AddNewAPI(apiDef, prodHTTPRoute, sandHTTPRoute)
		*apiReconciler.ch <- synchronizer.APIEvent{EventType: constants.Create, Event: apiState}
		return ctrl.Result{}, nil
	}
	apiState := synchronizer.APIState{}
	if apiDef.Generation > cachedAPI.APIDefinition.Generation {
		apiStateUpdate, err := apiReconciler.ods.UpdateAPIDef(apiDef)
		if err != nil {
			loggers.LoggerAPKOperator.ErrorC(logging.ErrorDetails{
				Message:   fmt.Sprintf("error updating API CR in operator data store: %v", err),
				Severity:  logging.TRIVIAL,
				ErrorCode: 2606,
			})
			return ctrl.Result{}, err
		}
		apiState = apiStateUpdate
	}
	if prodHTTPRoute != nil && (cachedAPI.ProdHTTPRoute == nil || prodHTTPRoute.UID != cachedAPI.ProdHTTPRoute.UID || prodHTTPRoute.Generation > cachedAPI.ProdHTTPRoute.Generation) {
		apiStateUpdate, err := apiReconciler.ods.UpdateHTTPRoute(utils.NamespacedName(&apiDef), prodHTTPRoute, true)
		if err != nil {
			loggers.LoggerAPKOperator.ErrorC(logging.ErrorDetails{
				Message:   fmt.Sprintf("error updating prod HTTPRoute CR in operator data store: %v", err),
				Severity:  logging.TRIVIAL,
				ErrorCode: 2607,
			})
			return ctrl.Result{}, err
		}
		apiState = apiStateUpdate
	}
	if sandHTTPRoute != nil && (cachedAPI.SandHTTPRoute == nil || sandHTTPRoute.UID != cachedAPI.SandHTTPRoute.UID || sandHTTPRoute.Generation > cachedAPI.SandHTTPRoute.Generation) {
		apiStateUpdate, err := apiReconciler.ods.UpdateHTTPRoute(utils.NamespacedName(&apiDef), sandHTTPRoute, false)
		if err != nil {
			loggers.LoggerAPKOperator.ErrorC(logging.ErrorDetails{
				Message:   fmt.Sprintf("error updating sand HTTPRoute CR in operator data store: %v", err),
				Severity:  logging.TRIVIAL,
				ErrorCode: 2617,
			})
			return ctrl.Result{}, err
		}
		apiState = apiStateUpdate
	}
	if (synchronizer.APIState{}) != apiState {
		*apiReconciler.ch <- synchronizer.APIEvent{EventType: constants.Update, Event: apiState}
	}
	return ctrl.Result{}, nil
}

// validateHTTPRouteRefs validates the HTTPRouteRefs related to a particular API by checking whether the
// HTTPRoutes exists in the controller cache or not.
//
// TODO : Consider HTTPRoute status also when validating.
func validateHTTPRouteRefs(ctx context.Context, client client.Client, namespace string,
	prodHTTPRouteRef string, sandHTTPRouteRef string) (*gwapiv1b1.HTTPRoute, *gwapiv1b1.HTTPRoute, error) {
	var prodHTTPRoute *gwapiv1b1.HTTPRoute
	var sandHTTPRoute *gwapiv1b1.HTTPRoute
	if prodHTTPRouteRef == "" && sandHTTPRouteRef == "" {
		return nil, nil, errors.New("an endpoint should have given for the API")
	}

	if prodHTTPRouteRef != "" {
		httpRoute := gwapiv1b1.HTTPRoute{}
		if err := client.Get(ctx, types.NamespacedName{Namespace: namespace, Name: prodHTTPRouteRef}, &httpRoute); err != nil {
			return nil, nil, fmt.Errorf("production HttpRoute %s in namespace :%s has not found. %s",
				prodHTTPRouteRef, namespace, err.Error())
		}
		prodHTTPRoute = &httpRoute
	}

	if sandHTTPRouteRef != "" {
		httpRoute := gwapiv1b1.HTTPRoute{}
		if err := client.Get(ctx, types.NamespacedName{Namespace: namespace, Name: sandHTTPRouteRef}, &httpRoute); err != nil {
			return nil, nil, fmt.Errorf("error fetching SandHTTPRoute: %s in namespace : %s. %v",
				sandHTTPRouteRef, namespace, err)
		}
		sandHTTPRoute = &httpRoute
	}

	return prodHTTPRoute, sandHTTPRoute, nil
}

// getAPIForHTTPRoute triggers the API controller reconcile method based on the changes detected
// from HTTPRoute objects. If the changes are done for an API stored in the Operator Data store,
// a new reconcile event will be created and added to the reconcile event queue.
func (apiReconciler *APIReconciler) getAPIForHTTPRoute(obj client.Object) []reconcile.Request {
	ctx := context.Background()
	httpRoute, ok := obj.(*gwapiv1b1.HTTPRoute)
	if !ok {
		loggers.LoggerAPKOperator.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprintf("Unexpected object type, bypassing reconciliation: %v", httpRoute),
			Severity:  logging.TRIVIAL,
			ErrorCode: 2608,
		})
		return []reconcile.Request{}
	}

	apiList := &dpv1alpha1.APIList{}
	if err := apiReconciler.client.List(ctx, apiList, &client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(httpRouteAPIIndex, utils.NamespacedName(httpRoute).String()),
	}); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprintf("Unable to find associated APIs: %s", utils.NamespacedName(httpRoute).String()),
			Severity:  logging.CRITICAL,
			ErrorCode: 2610,
		})
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

// addAPIIndexers adds indexing on API, for prodution and sandbox HTTPRoutes
// referenced in API objects via `.spec.prodHTTPRouteRef` and `.spec.sandHTTPRouteRef`.
// This helps to find APIs that are affected by a HTTPRoute CRUD operation.
func addAPIIndexers(ctx context.Context, mgr manager.Manager) error {
	err := mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha1.API{}, httpRouteAPIIndex, func(rawObj client.Object) []string {
		api := rawObj.(*dpv1alpha1.API)
		var httpRoutes []string
		if api.Spec.ProdHTTPRouteRef != "" {
			httpRoutes = append(httpRoutes,
				types.NamespacedName{
					Namespace: api.Namespace,
					Name:      api.Spec.ProdHTTPRouteRef,
				}.String())
		}
		if api.Spec.SandHTTPRouteRef != "" {
			httpRoutes = append(httpRoutes,
				types.NamespacedName{
					Namespace: api.Namespace,
					Name:      api.Spec.SandHTTPRouteRef,
				}.String())
		}
		return httpRoutes
	})
	return err
}
