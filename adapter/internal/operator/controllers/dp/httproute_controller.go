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
	"fmt"

	"github.com/wso2/apk/adapter/config"
	"github.com/wso2/apk/adapter/internal/loggers"
	"github.com/wso2/apk/adapter/internal/operator/constants"
	"github.com/wso2/apk/adapter/internal/operator/synchronizer"
	"github.com/wso2/apk/adapter/internal/operator/utils"
	"github.com/wso2/apk/adapter/pkg/logging"
	k8error "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/source"
	gwapiv1b1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

// HTTPRouteReconciler reconciles a HTTPRoute object.
type HTTPRouteReconciler struct {
	client client.Client
	ods    *synchronizer.OperatorDataStore
}

// NewHTTPRouteController creates a new HTTPRoute controller.
func NewHTTPRouteController(mgr manager.Manager, ods *synchronizer.OperatorDataStore) error {
	r := &HTTPRouteReconciler{
		client: mgr.GetClient(),
		ods:    ods,
	}
	c, err := controller.New(constants.HTTPRouteController, mgr, controller.Options{Reconciler: r})
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprintf("error creating HttpRoute controller: %v", err),
			Severity:  logging.BLOCKER,
			ErrorCode: 2609,
		})
		return err
	}

	conf := config.ReadConfigs()
	predicates := []predicate.Predicate{predicate.NewPredicateFuncs(utils.FilterByNamespaces(conf.Adapter.Operator.Namespaces))}

	if err := c.Watch(&source.Kind{Type: &gwapiv1b1.HTTPRoute{}}, &handler.EnqueueRequestForObject{},
		predicates...); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprintf("error watching HttpRoute objects: %v", err),
			Severity:  logging.BLOCKER,
			ErrorCode: 2610,
		})
		return err
	}
	return nil
}

// Reconcile gets triggered when a HTTPRoute object gets changed.
func (r *HTTPRouteReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var httpRoute gwapiv1b1.HTTPRoute
	if err := r.client.Get(ctx, req.NamespacedName, &httpRoute); err != nil {
		apiRef, found := r.ods.HTTPRouteToAPIRefs[req.NamespacedName]
		if found && k8error.IsNotFound(err) {
			// The route doesn't exist in the route Cache, remove it
			delete(r.ods.HTTPRouteToAPIRefs, req.NamespacedName)
			loggers.LoggerAPKOperator.Infof("HTTPRoute : %s for API : %s deleted from HTTPRouteToAPIRefs cache",
				req.NamespacedName.String(), apiRef)
			return ctrl.Result{}, nil
		}
		loggers.LoggerAPKOperator.Errorf("httpRoute : %v reconcile request not found. Error : %v.",
			req.NamespacedName, err)
		return ctrl.Result{}, err
	}
	// TODO: Add validation for backendRefs and HttpRoute status.
	loggers.LoggerAPKOperator.Debugf("Reconciled HTTPRoute : %s successfully", httpRoute.Name)
	return ctrl.Result{}, nil
}
