/*
 *  Copyright (c) 2022, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
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

	"github.com/wso2/apk/adapter/internal/loggers"
	"github.com/wso2/apk/adapter/internal/operator/synchronizer"
	"github.com/wso2/apk/adapter/internal/operator/utils"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	gwapiv1b1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	dpv1alpha1 "github.com/wso2/apk/adapter/internal/operator/api/v1alpha1"
)

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
	c, err := controller.New("API", mgr, controller.Options{Reconciler: r})
	if err != nil {
		loggers.LoggerAPKOperator.Errorf("Error creating API controller: %v", err)
		return err
	}

	if err := c.Watch(&source.Kind{Type: &dpv1alpha1.API{}}, &handler.EnqueueRequestForObject{}); err != nil {
		loggers.LoggerAPKOperator.Errorf("Error watching API resources: %v", err)
		return err
	}
	if err := c.Watch(&source.Kind{Type: &gwapiv1b1.HTTPRoute{}}, handler.EnqueueRequestsFromMapFunc(r.getAPIForHTTPRoute)); err != nil {
		loggers.LoggerAPKOperator.Errorf("Error watching HttpRoute from API Controller: %v", err)
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
// TODO(user): Modify the Reconcile function to compare the state specified by
// the API object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *APIReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	// 1. Check whether the API Def exist, if not DELETE event.
	var apiDef dpv1alpha1.API
	if err := r.client.Get(ctx, req.NamespacedName, &apiDef); err != nil {
		loggers.LoggerAPKOperator.Errorf("apiDef related to reconcile with key: %v not found", req.NamespacedName.String())
	}

	// 2. Handle Http route validation
	prodHTTPRoute, sandHTTPRoute, err := validateHTTPRouteRefs(ctx, r.client, req.Namespace, apiDef.Spec.ProdHTTPRouteRef, apiDef.Spec.SandHTTPRouteRef)
	if err != nil {
		loggers.LoggerAPKOperator.Errorf("Error validating the HttpRouteRefs for API: %v", apiDef.Spec.APIDisplayName)
		return ctrl.Result{}, err
	}

	// 3. Check the Operator data store for the received API event
	cachedAPI, found := r.ods.GetAPI(utils.NamespacedName(&apiDef))

	if !found {
		apiState, err := r.ods.AddNewAPI(apiDef, prodHTTPRoute, sandHTTPRoute)
		if err != nil {
			loggers.LoggerAPKOperator.Errorf("Error storing the new API in the operator data store: %v", err)
			return ctrl.Result{}, err
		}
		*r.ch <- synchronizer.APIEvent{EventType: "CREATE", Event: apiState}
		return ctrl.Result{}, nil
	}
	apiState := synchronizer.APIState{}
	if apiDef.Generation > cachedAPI.APIDefinition.Generation {
		apiStateUpdate, err := r.ods.UpdateAPIDef(apiDef)
		if err != nil {
			return ctrl.Result{}, err
		}
		apiState = apiStateUpdate
	}
	if prodHTTPRoute.Generation > cachedAPI.ProdHTTPRoute.Generation {
		apiStateUpdate, err := r.ods.UpdateHTTPRoute(utils.NamespacedName(&apiDef), prodHTTPRoute, true)
		if err != nil {
			return ctrl.Result{}, err
		}
		apiState = apiStateUpdate
	}
	*r.ch <- synchronizer.APIEvent{EventType: "UPDATE", Event: apiState}
	return ctrl.Result{}, nil

}

func validateHTTPRouteRefs(ctx context.Context, client client.Client, namespace string,
	prodHTTPRouteRef string, sandHTTPRouteRef string) (gwapiv1b1.HTTPRoute, gwapiv1b1.HTTPRoute, error) {
	var prodHTTPRoute gwapiv1b1.HTTPRoute
	if err := client.Get(ctx, types.NamespacedName{Namespace: namespace, Name: prodHTTPRouteRef}, &prodHTTPRoute); err != nil {
		loggers.LoggerAPKOperator.Errorf("Production HttpRoute not found: %v", prodHTTPRouteRef)
		return gwapiv1b1.HTTPRoute{}, gwapiv1b1.HTTPRoute{}, err
	}
	var sandHTTPRoute gwapiv1b1.HTTPRoute
	if sandHTTPRouteRef != "" {
		if err := client.Get(ctx, types.NamespacedName{Namespace: namespace, Name: sandHTTPRouteRef}, &sandHTTPRoute); err != nil {
			loggers.LoggerAPKOperator.Errorf("Error fetching SandHTTPRoute: %v:%v", sandHTTPRouteRef, err)
			return prodHTTPRoute, gwapiv1b1.HTTPRoute{}, err
		}
	}
	return prodHTTPRoute, sandHTTPRoute, nil
}

func (r *APIReconciler) getAPIForHTTPRoute(obj client.Object) []reconcile.Request {
	httpRoute, ok := obj.(*gwapiv1b1.HTTPRoute)
	if !ok {
		loggers.LoggerAPKOperator.Errorf("Unexpected object type, bypassing reconciliation: %v", httpRoute)
	}
	apiRef, found := r.ods.GetAPIForHTTPRoute(utils.NamespacedName(httpRoute))
	if !found {
		loggers.LoggerAPKOperator.Infof("API for HttpRoute not found: %v", httpRoute.Name)
		return []reconcile.Request{}
	}
	requests := []reconcile.Request{}
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      apiRef.Name,
			Namespace: apiRef.Namespace},
	}
	loggers.LoggerAPKOperator.Infof("Adding reconcile request: %v", req.NamespacedName)
	requests = append(requests, req)
	return requests
}
