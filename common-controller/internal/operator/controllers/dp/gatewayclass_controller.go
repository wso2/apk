/*
 *  Copyright (c) 2024, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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
	"time"

	"github.com/wso2/apk/adapter/pkg/logging"
	loggers "github.com/wso2/apk/common-controller/internal/loggers"
	constants "github.com/wso2/apk/common-controller/internal/operator/constant"
	"github.com/wso2/apk/common-controller/internal/operator/status"
	k8error "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	k8client "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/source"
	gwapiv1b1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

const (
	gatewayRateLimitPolicyIndex = "gatewayRateLimitPolicyIndex"
	gatewayAPIPolicyIndex       = "gatewayAPIPolicyIndex"
)

var (
	setReadiness sync.Once
)

// GatewayClassReconciler reconciles a Gateway object
type GatewayClassReconciler struct {
	client        k8client.Client
	statusUpdater *status.UpdateHandler
	mgr           manager.Manager
}

// NewGatewayClassController creates a new GatewayClass controller instance. GatewayClass Controllers watches for gwapiv1b1.GatewayClass.
func NewGatewayClassController(mgr manager.Manager, statusUpdater *status.UpdateHandler) error {
	r := &GatewayClassReconciler{
		client:        mgr.GetClient(),
		statusUpdater: statusUpdater,
		mgr:           mgr,
	}

	c, err := controller.New(constants.GatewayClassController, mgr, controller.Options{Reconciler: r})
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2663, logging.BLOCKER,
			"Error creating GatewayClass controller: %v", err.Error()))
		return err
	}

	if err := c.Watch(source.Kind(mgr.GetCache(), &gwapiv1b1.GatewayClass{}), &handler.EnqueueRequestForObject{}); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2639, logging.BLOCKER,
			"Error watching GatewayClass resources: %v", err.Error()))
		return err
	}

	loggers.LoggerAPKOperator.Info("GatwayClasses Controller successfully started. Watching GatewayClass Objects...")
	return nil
}

//+kubebuilder:rbac:groups=dp.wso2.com,resources=gatewayclasses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=dp.wso2.com,resources=gatewayclasses/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=dp.wso2.com,resources=gatewayclasses/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (gatewayClassReconciler *GatewayClassReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// Check whether the Gateway CR exist, if not consider as a DELETE event.
	loggers.LoggerAPKOperator.Info("Reconciling gateway class...")
	var gatewayClassDef gwapiv1b1.GatewayClass
	if err := gatewayClassReconciler.client.Get(ctx, req.NamespacedName, &gatewayClassDef); err != nil {
		if k8error.IsNotFound(err) {
			// Gateway deleted. We dont have to handle this
			return ctrl.Result{}, nil
		}
		return ctrl.Result{
			RequeueAfter: time.Duration(1 * time.Second),
		}, nil
	}
	gatewayClassReconciler.handleGatewayClassStatus(req.NamespacedName, constants.Create, []string{})
	return ctrl.Result{}, nil
}

// handleGatewayClassStatus updates the Gateway CR update
func (gatewayClassReconciler *GatewayClassReconciler) handleGatewayClassStatus(gatewayKey types.NamespacedName, state string, events []string) {
	accept := false
	message := ""

	switch state {
	case constants.Create:
		accept = true
		message = "GatewayClass is deployed successfully"
	case constants.Update:
		accept = true
		message = fmt.Sprintf("GatewayClass update is deployed successfully. %v Updated", events)
	}
	timeNow := metav1.Now()
	gatewayClassReconciler.statusUpdater.Send(status.Update{
		NamespacedName: gatewayKey,
		Resource:       new(gwapiv1b1.GatewayClass),
		UpdateStatus: func(obj k8client.Object) k8client.Object {
			h, ok := obj.(*gwapiv1b1.GatewayClass)
			if !ok {
				loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error3109, logging.BLOCKER, "Error while updating GatewayClass status %v", obj))
			}
			hCopy := h.DeepCopy()
			var gwConditions []metav1.Condition = hCopy.Status.Conditions
			gwConditions[0].Status = "Unknown"
			if accept {
				gwConditions[0].Status = "True"
			} else {
				gwConditions[0].Status = "False"
			}
			gwConditions[0].Message = message
			gwConditions[0].LastTransitionTime = timeNow
			gwConditions[0].Reason = "Reconciled"
			gwConditions[0].Type = "Accepted"
			generation := hCopy.ObjectMeta.Generation
			for i := range gwConditions {
				// Assign generation to ObservedGeneration
				gwConditions[i].ObservedGeneration = generation
			}
			hCopy.Status.Conditions = gwConditions
			return hCopy
		},
	})
}
