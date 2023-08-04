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

package cp

import (
	"context"
	"fmt"

	"github.com/wso2/apk/adapter/internal/discovery/xds"
	"github.com/wso2/apk/adapter/internal/loggers"
	cpv1alpha1 "github.com/wso2/apk/adapter/internal/operator/apis/cp/v1alpha1"
	"github.com/wso2/apk/adapter/internal/operator/constants"
	"github.com/wso2/apk/adapter/internal/operator/utils"
	"github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/subscription"
	"github.com/wso2/apk/adapter/pkg/logging"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// SubscriptionReconciler reconciles a Subscription object
type SubscriptionReconciler struct {
	client client.Client
	Scheme *runtime.Scheme
}

// NewSubscriptionController creates a new Subscription controller instance.
func NewSubscriptionController(mgr manager.Manager) error {
	r := &SubscriptionReconciler{
		client: mgr.GetClient(),
	}
	c, err := controller.New(constants.SubscriptionController, mgr, controller.Options{Reconciler: r})
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2608, logging.BLOCKER, "Error creating Subscription controller: %v", err.Error()))
		return err
	}

	if err := c.Watch(&source.Kind{Type: &cpv1alpha1.Subscription{}}, &handler.EnqueueRequestForObject{},
		predicate.NewPredicateFuncs(utils.FilterByNamespaces([]string{utils.GetOperatorPodNamespace()}))); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2609, logging.BLOCKER, "Error watching Subscription resources: %v", err.Error()))
		return err
	}

	loggers.LoggerAPKOperator.Debug("Subscription Controller successfully started. Watching Subscription Objects...")
	return nil
}

//+kubebuilder:rbac:groups=cp.wso2.com,resources=subscriptions,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cp.wso2.com,resources=subscriptions/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cp.wso2.com,resources=subscriptions/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Subscription object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (subscriptionReconciler *SubscriptionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	loggers.LoggerAPKOperator.Debugf("Reconciling subscription: %v", req.NamespacedName.String())

	subscriptionKey := req.NamespacedName
	var subscriptionList = new(cpv1alpha1.SubscriptionList)
	if err := subscriptionReconciler.client.List(ctx, subscriptionList); err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to get subscriptions %s/%s",
			subscriptionKey.Namespace, subscriptionKey.Name)
	}
	sendSubUpdates(subscriptionList)
	return ctrl.Result{}, nil
}

func sendSubUpdates(subscriptionsList *cpv1alpha1.SubscriptionList) {
	subList := marshalSubscriptionList(subscriptionsList.Items)
	xds.UpdateEnforcerSubscriptions(subList)
}

func marshalSubscriptionList(subscriptionList []cpv1alpha1.Subscription) *subscription.SubscriptionList {
	subscriptions := []*subscription.Subscription{}
	for _, subInternal := range subscriptionList {
		sub := &subscription.Subscription{
			Uuid:           string(subInternal.UID),
			ApiRef:         subInternal.Spec.APIRef,
			PolicyId:       subInternal.Spec.PolicyID,
			SubStatus:      subInternal.Spec.SubscriptionStatus,
			ApplicationRef: subInternal.Spec.ApplicationRef,
		}
		subscriptions = append(subscriptions, sub)
	}
	return &subscription.SubscriptionList{
		List: subscriptions,
	}
}
