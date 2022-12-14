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
	"github.com/wso2/apk/adapter/internal/operator/status"
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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const httpRouteAPIIndex = "httpRouteAPIIndex"

// APIReconciler reconciles a API object
type APIReconciler struct {
	client        client.Client
	ods           *synchronizer.OperatorDataStore
	ch            *chan synchronizer.APIEvent
	statusUpdater *status.UpdateHandler
}

// NewAPIController creates a new API controller instance. API Controllers watches for dpv1alpha1.API and gwapiv1b1.HTTPRoute.
func NewAPIController(mgr manager.Manager, operatorDataStore *synchronizer.OperatorDataStore, statusUpdater *status.UpdateHandler,
	ch *chan synchronizer.APIEvent) error {
	r := &APIReconciler{
		client:        mgr.GetClient(),
		ods:           operatorDataStore,
		ch:            ch,
		statusUpdater: statusUpdater,
	}
	c, err := controller.New(constants.APIController, mgr, controller.Options{Reconciler: r})
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprintf("Error creating API controller : %v", err),
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
			Message:   fmt.Sprintf("Error watching API resources: %v", err),
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
			Message:   fmt.Sprintf("Error watching HTTPRoute resources: %v", err),
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

	// Check whether the API CR exist, if not consider as a DELETE event.
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

	// Reject invalid API
	if apiDef.Status.Status == constants.InvalidState {
		loggers.LoggerAPKOperator.Debugf("API CR is rejected as it has been invalidated already.")
		return ctrl.Result{}, nil
	}

	// Validate API CR
	if apiDef.Status.Status == "" {
		if valid, err := validateAPICR(apiDef.Spec); !valid {
			apiReconciler.handleStatus(req.NamespacedName, constants.InvalidState, []string{})
			if err != nil {
				loggers.LoggerAPKOperator.ErrorC(logging.ErrorDetails{
					Message:   fmt.Sprintf("Error validating API CR in namespace : %s, %v", req.NamespacedName.String(), err),
					Severity:  logging.TRIVIAL,
					ErrorCode: 2604,
				})
				return ctrl.Result{}, err
			}
		}
		apiReconciler.handleStatus(req.NamespacedName, constants.ValidatedState, []string{})
	}

	// Retrieve HTTPRoutes
	prodHTTPRoute, sandHTTPRoute, err := getHTTPRoutes(ctx, apiReconciler.client, req.Namespace,
		apiDef.Spec.ProdHTTPRouteRef, apiDef.Spec.SandHTTPRouteRef)
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprintf("Error retrieving HttpRoute CRs in namespace : %s, %v", req.NamespacedName.String(), err),
			Severity:  logging.TRIVIAL,
			ErrorCode: 2604,
		})
		return ctrl.Result{}, err
	}
	loggers.LoggerAPKOperator.Debugf("HTTPRoutes are retrieved successfully for API CR %s", req.NamespacedName.String())

	// Check whether the Operator Data store contains the received API.
	cachedAPI, found := apiReconciler.ods.APIStore[req.NamespacedName]

	if !found {
		apiState := apiReconciler.ods.AddNewAPI(apiDef, prodHTTPRoute, sandHTTPRoute)
		*apiReconciler.ch <- synchronizer.APIEvent{EventType: constants.Create, Event: apiState}
		apiReconciler.handleStatus(req.NamespacedName, constants.DeployedState, []string{})
		return ctrl.Result{}, nil
	}

	if events, updated := updateAPIState(&apiDef, cachedAPI, prodHTTPRoute, sandHTTPRoute); updated {
		*apiReconciler.ch <- synchronizer.APIEvent{EventType: constants.Update, Event: *cachedAPI}
		apiReconciler.handleStatus(req.NamespacedName, constants.UpdatedState, events)
	}
	return ctrl.Result{}, nil
}

// validateAPICR validates fields in API
// TODO(amali) Add validations for all fields
func validateAPICR(apiSpec dpv1alpha1.APISpec) (bool, error) {
	if apiSpec.ProdHTTPRouteRef == "" && apiSpec.SandHTTPRouteRef == "" {
		return false, errors.New("an endpoint should have given for the API")
	}
	return true, nil
}

// getHTTPRoutes validates the HTTPRouteRefs related to a particular API by checking whether the
// HTTPRoutes exists in the controller cache or not.
func getHTTPRoutes(ctx context.Context, client client.Client, namespace string,
	prodHTTPRouteRef string, sandHTTPRouteRef string) (*gwapiv1b1.HTTPRoute, *gwapiv1b1.HTTPRoute, error) {
	var prodHTTPRoute *gwapiv1b1.HTTPRoute
	var sandHTTPRoute *gwapiv1b1.HTTPRoute

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

// updateAPIState update the APIState on ref updates
func updateAPIState(apiDef *dpv1alpha1.API, cachedAPI *synchronizer.APIState, prodHTTPRoute *gwapiv1b1.HTTPRoute,
	sandHTTPRoute *gwapiv1b1.HTTPRoute) ([]string, bool) {
	var updated bool
	events := []string{}
	if apiDef.Generation > cachedAPI.APIDefinition.Generation {
		cachedAPI.APIDefinition = apiDef
		updated = true
		events = append(events, "API Definition")
	}
	if prodHTTPRoute != nil && (prodHTTPRoute.UID != cachedAPI.ProdHTTPRoute.UID ||
		prodHTTPRoute.Generation > cachedAPI.ProdHTTPRoute.Generation) {
		cachedAPI.ProdHTTPRoute = prodHTTPRoute
		updated = true
		events = append(events, "Production Endpoint")
	}
	if sandHTTPRoute != nil && (sandHTTPRoute.UID != cachedAPI.SandHTTPRoute.UID ||
		sandHTTPRoute.Generation > cachedAPI.SandHTTPRoute.Generation) {
		cachedAPI.SandHTTPRoute = sandHTTPRoute
		updated = true
		events = append(events, "Sandbox Endpoint")
	}
	return events, updated
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

// handleStatus updates the API CR update
func (apiReconciler *APIReconciler) handleStatus(apiKey types.NamespacedName, state string, events []string) {
	accept := false
	message := ""
	event := ""

	switch state {
	case constants.DeployedState:
		accept = true
		message = "API is deployed to the gateway."
	case constants.InvalidState:
		accept = false
		message = "Rejected due to invalid data."
	case constants.ValidatedState:
		accept = true
		message = "Successfully validated."
	case constants.UpdatedState:
		accept = true
		message = fmt.Sprintf("API update is deployed to the gateway. %v Updated", events)
	}
	timeNow := metav1.Now()
	event = fmt.Sprintf("[%s] %s", timeNow.String(), message)

	apiReconciler.statusUpdater.Send(status.Update{
		NamespacedName: apiKey,
		Resource:       new(dpv1alpha1.API),
		UpdateStatus: func(obj client.Object) client.Object {
			h, ok := obj.(*dpv1alpha1.API)
			if !ok {
				loggers.LoggerAPKOperator.ErrorC(logging.ErrorDetails{
					Message:   fmt.Sprintf("Unsupported object type %T", obj),
					Severity:  logging.BLOCKER,
					ErrorCode: 2617,
				})
			}
			hCopy := h.DeepCopy()
			hCopy.Status.Status = state
			hCopy.Status.Accepted = accept
			hCopy.Status.Message = message
			hCopy.Status.Events = append(hCopy.Status.Events, event)
			hCopy.Status.TransitionTime = &timeNow
			return hCopy
		},
	})
}
