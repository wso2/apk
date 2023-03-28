/*
 *  Copyright (c) 2023, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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
	"github.com/wso2/apk/adapter/pkg/logging"
	k8error "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/types"

	dpv1alpha1 "github.com/wso2/apk/adapter/pkg/operator/apis/dp/v1alpha1"
	"github.com/wso2/apk/adapter/pkg/operator/constants"
	"github.com/wso2/apk/adapter/pkg/operator/status"
	"github.com/wso2/apk/adapter/pkg/operator/synchronizer"
	"github.com/wso2/apk/adapter/pkg/operator/utils"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	gwapiv1b1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	ctrl "sigs.k8s.io/controller-runtime"
	k8client "sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	gatewayIndex = "gatewayIndex"
)

// GatewayReconciler reconciles a Gateway object
type GatewayReconciler struct {
	client        k8client.Client
	ods           *synchronizer.OperatorDataStore
	ch            *chan synchronizer.GatewayEvent
	statusUpdater *status.UpdateHandler
	mgr           manager.Manager
}

// NewGatewayController creates a new Gateway controller instance. Gateway Controllers watches for gwapiv1b1.Gateway.
func NewGatewayController(mgr manager.Manager, operatorDataStore *synchronizer.OperatorDataStore, statusUpdater *status.UpdateHandler,
	ch *chan synchronizer.GatewayEvent) error {
	r := &GatewayReconciler{
		client:        mgr.GetClient(),
		ods:           operatorDataStore,
		ch:            ch,
		statusUpdater: statusUpdater,
		mgr:           mgr,
	}
	c, err := controller.New(constants.GatewayController, mgr, controller.Options{Reconciler: r})
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2610, err))
		return err
	}

	conf := config.ReadConfigs()
	ctx := context.Background()
	predicates := []predicate.Predicate{predicate.NewPredicateFuncs(utils.FilterByNamespaces(conf.Adapter.Operator.Namespaces))}

	if err := addGatewayIndexes(ctx, mgr); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2612, err))
		return err
	}

	if err := c.Watch(&source.Kind{Type: &gwapiv1b1.Gateway{}}, &handler.EnqueueRequestForObject{},
		predicates...); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2611, err))
		return err
	}

	predicates = append(predicates, predicate.NewPredicateFuncs(func(object k8client.Object) bool {
		rlPolicy := object.(*dpv1alpha1.RateLimitPolicy)
		if rlPolicy.Spec.TargetRef.Kind == constants.KindGateway {
			return true
		}
		return false
	}))

	if err := c.Watch(&source.Kind{Type: &dpv1alpha1.RateLimitPolicy{}}, handler.EnqueueRequestsFromMapFunc(r.handleCustomRateLimitPolicies), predicates...); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2611, err))
		return err
	}

	loggers.LoggerAPKOperator.Info("Gateway Controller successfully started. Watching Gateway Objects....")
	return nil
}

//+kubebuilder:rbac:groups=dp.wso2.com,resources=gateways,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=dp.wso2.com,resources=gateways/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=dp.wso2.com,resources=gateways/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Gateway object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (gatewayReconciler *GatewayReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// Check whether the Gateway CR exist, if not consider as a DELETE event.
	var gatewayDef gwapiv1b1.Gateway
	if err := gatewayReconciler.client.Get(ctx, req.NamespacedName, &gatewayDef); err != nil {
		gatewayState, found := gatewayReconciler.ods.GetCachedGateway(req.NamespacedName)
		if found && k8error.IsNotFound(err) {
			// The gateway doesn't exist in the gateway Cache, remove it
			gatewayReconciler.ods.DeleteCachedGateway(req.NamespacedName)
			loggers.LoggerAPKOperator.Infof("Delete event has received for Gateway : %s, hence deleted from Gateway cache", req.NamespacedName.String())
			*gatewayReconciler.ch <- synchronizer.GatewayEvent{EventType: constants.Delete, Event: gatewayState}
			return ctrl.Result{}, nil
		}
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2619, req.NamespacedName.String(), err))
		return ctrl.Result{}, nil
	}
	var gwCondition []metav1.Condition = gatewayDef.Status.Conditions
	customRateLimitPolicies, err := gatewayReconciler.getCustomRateLimitPoliciesForGateway(utils.NamespacedName(&gatewayDef))
	if err != nil {
		loggers.LoggerAPKOperator.Infof("XXXXXX Error: %v", err)
	}
	loggers.LoggerAPKOperator.Infof("XXXXXX customRateLimitPolicies: %v", len(customRateLimitPolicies))
	if gwCondition[0].Type != "Accepted" {
		gatewayState := gatewayReconciler.ods.AddGatewayState(gatewayDef, customRateLimitPolicies)
		*gatewayReconciler.ch <- synchronizer.GatewayEvent{EventType: constants.Create, Event: gatewayState}
		gatewayReconciler.handleGatewayStatus(req.NamespacedName, constants.DeployedState, []string{})
	} else if cachedGateway, events, updated :=
		gatewayReconciler.ods.UpdateGatewayState(&gatewayDef, customRateLimitPolicies); updated {
		*gatewayReconciler.ch <- synchronizer.GatewayEvent{EventType: constants.Update, Event: cachedGateway}
		gatewayReconciler.handleGatewayStatus(req.NamespacedName, constants.UpdatedState, events)
	}
	return ctrl.Result{}, nil
}

// handleStatus updates the Gateway CR update
func (gatewayReconciler *GatewayReconciler) handleGatewayStatus(gatewayKey types.NamespacedName, state string, events []string) {
	accept := false
	message := ""
	//event := ""

	switch state {
	case constants.DeployedState:
		accept = true
		message = "Gateway is deployed successfully"
	case constants.UpdatedState:
		accept = true
		message = fmt.Sprintf("Gateway update is deployed successfully. %v Updated", events)
	}
	timeNow := metav1.Now()
	//event = fmt.Sprintf("[%s] %s", timeNow.String(), message)

	gatewayReconciler.statusUpdater.Send(status.Update{
		NamespacedName: gatewayKey,
		Resource:       new(gwapiv1b1.Gateway),
		UpdateStatus: func(obj k8client.Object) k8client.Object {
			h, ok := obj.(*gwapiv1b1.Gateway)
			if !ok {
				loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2626, obj))
			}
			hCopy := h.DeepCopy()
			var gwCondition []metav1.Condition = hCopy.Status.Conditions
			gwCondition[0].Status = "Unknown"
			if accept {
				gwCondition[0].Status = "True"
			} else {
				gwCondition[0].Status = "False"
			}
			gwCondition[0].Message = message
			gwCondition[0].LastTransitionTime = timeNow
			// gwCondition[0].Reason = append(gwCondition[0].Reason, event)
			gwCondition[0].Reason = "Reconciled"
			gwCondition[0].Type = state
			hCopy.Status.Conditions = gwCondition
			return hCopy
		},
	})
}

func (gatewayReconciler *GatewayReconciler) handleCustomRateLimitPolicies(obj k8client.Object) []reconcile.Request {
	// ctx := context.Background()
	ratelimitPolicy, ok := obj.(*dpv1alpha1.RateLimitPolicy)
	if !ok {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2622, ratelimitPolicy))
		return []reconcile.Request{}
	}
	// utils.SelectPolicy(&ratelimitPolicy.Spec.Override, &ratelimitPolicy.Spec.Default, nil, nil)
	return []reconcile.Request{{
		NamespacedName: types.NamespacedName{
			Namespace: ratelimitPolicy.Namespace,
			Name:      string(ratelimitPolicy.Spec.TargetRef.Name),
		},},
	}
}

func (gatewayReconciler *GatewayReconciler) getCustomRateLimitPoliciesForGateway(gatewayName types.NamespacedName) ([]*dpv1alpha1.RateLimitPolicy, error) {
	ctx := context.Background()
	var ratelimitPolicyList dpv1alpha1.RateLimitPolicyList
	var rateLimitPolicies []*dpv1alpha1.RateLimitPolicy
	if err := gatewayReconciler.client.List(ctx, &ratelimitPolicyList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(gatewayIndex, gatewayName.String()),
	}); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2623, err))
		return nil, err
	}
	loggers.LoggerAPKOperator.Infof("XXXXXX ratelimitPolicyList: %v", len(ratelimitPolicyList.Items))
	for _, item := range ratelimitPolicyList.Items {
		rateLimitPolicy := item
		rateLimitPolicies = append(rateLimitPolicies, &rateLimitPolicy)
	}
	return rateLimitPolicies, nil
}


func addGatewayIndexes(ctx context.Context, mgr manager.Manager) error { 
	return mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha1.RateLimitPolicy{}, gatewayIndex,
		func(rawObj k8client.Object) []string {
			ratelimitPolicy := rawObj.(*dpv1alpha1.RateLimitPolicy)
			var gateways []string
			if ratelimitPolicy.Spec.TargetRef.Kind == constants.KindGateway {
				gateways = append(gateways,
					types.NamespacedName{
						Namespace: utils.GetNamespace(
							(*gwapiv1b1.Namespace)(ratelimitPolicy.Spec.TargetRef.Namespace),
							ratelimitPolicy.Namespace),
						Name: string(ratelimitPolicy.Spec.TargetRef.Name),
					}.String())
			}
			return gateways
		})
}
