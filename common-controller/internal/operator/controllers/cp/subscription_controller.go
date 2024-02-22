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

	"github.com/wso2/apk/adapter/pkg/logging"
	"github.com/wso2/apk/common-controller/internal/cache"
	loggers "github.com/wso2/apk/common-controller/internal/loggers"
	"github.com/wso2/apk/common-controller/internal/server"
	"github.com/wso2/apk/common-controller/internal/utils"
	"github.com/wso2/apk/common-go-libs/constants"
	k8error "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/source"

	cpv1alpha2 "github.com/wso2/apk/common-go-libs/apis/cp/v1alpha2"
)

// SubscriptionReconciler reconciles a Subscription object
type SubscriptionReconciler struct {
	client client.Client
	Scheme *runtime.Scheme
	ods    *cache.SubscriptionDataStore
}

// NewSubscriptionController creates a new Subscription controller instance.
func NewSubscriptionController(mgr manager.Manager, subscriptionStore *cache.SubscriptionDataStore) error {
	r := &SubscriptionReconciler{
		client: mgr.GetClient(),
		ods:    subscriptionStore,
	}
	c, err := controller.New(constants.SubscriptionController, mgr, controller.Options{Reconciler: r})
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2608, logging.BLOCKER, "Error creating Subscription controller: %v", err.Error()))
		return err
	}

	if err := c.Watch(source.Kind(mgr.GetCache(), &cpv1alpha2.Subscription{}), &handler.EnqueueRequestForObject{},
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
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.4/pkg/reconcile
func (subscriptionReconciler *SubscriptionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	loggers.LoggerAPKOperator.Debugf("Reconciling subscription: %v", req.NamespacedName.String())

	subscriptionKey := req.NamespacedName
	var subscription cpv1alpha2.Subscription
	if err := subscriptionReconciler.client.Get(ctx, req.NamespacedName, &subscription); err != nil {
		if k8error.IsNotFound(err) {
			subscriptionSpec, state := subscriptionReconciler.ods.GetSubscriptionFromStore(subscriptionKey)
			if state {
				// Subscription in cache
				loggers.LoggerAPKOperator.Debugf("Subscription %s/%s not found. Ignoring since object must be deleted", subscriptionKey.Namespace, subscriptionKey.Name)
				utils.SendDeleteSubscriptionEvent(subscriptionKey.Name, subscriptionSpec)
				subscriptionReconciler.ods.DeleteSubscriptionFromStore(subscriptionKey)
				server.DeleteSubscription(subscriptionKey.Name)
				return ctrl.Result{}, nil
			}
		}
	} else {
		sendSubUpdates(subscription)
		utils.SendAddSubscriptionEvent(subscription)
		subscriptionReconciler.ods.AddorUpdateSubscriptionToStore(subscriptionKey, subscription.Spec)
	}
	return ctrl.Result{}, nil
}

func sendSubUpdates(subscription cpv1alpha2.Subscription) {
	subList := marshalSubscription(subscription)
	server.AddSubscription(subList)
}

func marshalSubscription(subscription cpv1alpha2.Subscription) server.Subscription {
	subscribedAPI := &server.SubscribedAPI{}
	sub := server.Subscription{
		UUID:         subscription.Name,
		SubStatus:    subscription.Spec.SubscriptionStatus,
		Organization: subscription.Spec.Organization,
	}
	if subscription.Spec.API.Name != "" && subscription.Spec.API.Version != "" {
		subscribedAPI.Name = subscription.Spec.API.Name
		subscribedAPI.Version = subscription.Spec.API.Version
	}
	sub.SubscribedAPI = subscribedAPI
	return sub
}
