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
	"github.com/wso2/apk/common-controller/internal/config"
	loggers "github.com/wso2/apk/common-controller/internal/loggers"
	"github.com/wso2/apk/common-controller/internal/server"
	"github.com/wso2/apk/common-go-libs/pkg/server/model"
	"github.com/wso2/apk/common-controller/internal/utils"
	"github.com/wso2/apk/common-go-libs/constants"
	k8error "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	k8client "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	cpv1alpha3 "github.com/wso2/apk/common-go-libs/apis/cp/v1alpha3"
	dpv1alpha3 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha3"
)

// SubscriptionReconciler reconciles a Subscription object
type SubscriptionReconciler struct {
	client client.Client
	Scheme *runtime.Scheme
	ods    *cache.SubscriptionDataStore
	rlODS  *cache.RatelimitDataStore
}

const (
	subscriptionRatelimitIndex     = "subscriptionRatelimitIndex"
	subscriptionToAIRatelimitIndex = "subscriptionToAIRatelimitIndex"
)

// NewSubscriptionController creates a new Subscription controller instance.
func NewSubscriptionController(mgr manager.Manager, subscriptionStore *cache.SubscriptionDataStore) error {
	r := &SubscriptionReconciler{
		client: mgr.GetClient(),
		ods:    subscriptionStore,
	}
	ctx := context.Background()
	conf := config.ReadConfigs()
	if err := addSubscriptionControllerIndexes(ctx, mgr); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2658, logging.CRITICAL, "Error adding indexes: %v", err))
		return err
	}
	c, err := controller.New(constants.SubscriptionController, mgr, controller.Options{Reconciler: r})
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2608, logging.BLOCKER, "Error creating Subscription controller: %v", err.Error()))
		return err
	}

	if err := c.Watch(source.Kind(mgr.GetCache(), &cpv1alpha3.Subscription{}, &handler.TypedEnqueueRequestForObject[*cpv1alpha3.Subscription]{},
		predicate.NewTypedPredicateFuncs(utils.FilterSubsByNamespaces([]string{utils.GetOperatorPodNamespace()})))); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2609, logging.BLOCKER, "Error watching Subscription resources: %v", err.Error()))
		return err
	}

	predicateRateLimitPolicy := []predicate.TypedPredicate[*dpv1alpha3.RateLimitPolicy]{predicate.NewTypedPredicateFuncs[*dpv1alpha3.RateLimitPolicy](utils.FilterRateLimitPolicyByNamespaces(conf.CommonController.Operator.Namespaces))}
	if err := c.Watch(source.Kind(mgr.GetCache(), &dpv1alpha3.RateLimitPolicy{}, handler.TypedEnqueueRequestsFromMapFunc(r.getSubscriptionForRatelimit),
		predicateRateLimitPolicy...)); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2639, logging.BLOCKER, "Error watching Ratelimit resources: %v", err))
		return err
	}

	predicateAIRatelimitPolicy := []predicate.TypedPredicate[*dpv1alpha3.AIRateLimitPolicy]{predicate.NewTypedPredicateFuncs[*dpv1alpha3.AIRateLimitPolicy](utils.FilterAIRatelimitPolicyByNamespaces(conf.CommonController.Operator.Namespaces))}
	if err := c.Watch(source.Kind(mgr.GetCache(), &dpv1alpha3.AIRateLimitPolicy{}, handler.TypedEnqueueRequestsFromMapFunc(r.getSubscriptionForAIRatelimit),
		predicateAIRatelimitPolicy...)); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2645, logging.BLOCKER, "Error watching AIRatelimitPolicy resources: %v", err))
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
	var subscription cpv1alpha3.Subscription
	if err := subscriptionReconciler.client.Get(ctx, req.NamespacedName, &subscription); err != nil {
		if k8error.IsNotFound(err) {
			_, state := subscriptionReconciler.ods.GetSubscriptionFromStore(subscriptionKey)
			if state {
				// Subscription in cache
				loggers.LoggerAPKOperator.Debugf("Subscription %s/%s not found. Ignoring since object must be deleted", subscriptionKey.Namespace, subscriptionKey.Name)
				utils.SendDeleteSubscriptionEvent(subscriptionKey.Name, subscription)
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

func sendSubUpdates(subscription cpv1alpha3.Subscription) {
	subList := marshalSubscription(subscription)
	server.AddSubscription(subList)
}

func marshalSubscription(subscription cpv1alpha3.Subscription) model.Subscription {
	subscribedAPI := &model.SubscribedAPI{}
	sub := model.Subscription{
		UUID:         subscription.Name,
		SubStatus:    subscription.Spec.SubscriptionStatus,
		Organization: subscription.Spec.Organization,
	}
	sub.RatelimitTier = subscription.Spec.RatelimitRef.Name
	if subscription.Spec.API.Name != "" && subscription.Spec.API.Version != "" {
		subscribedAPI.Name = subscription.Spec.API.Name
		subscribedAPI.Version = subscription.Spec.API.Version
	}
	sub.SubscribedAPI = subscribedAPI
	return sub
}

// addSubscriptionControllerIndexes adds indexes to the Subscription controller
func addSubscriptionControllerIndexes(ctx context.Context, mgr manager.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(ctx, &cpv1alpha3.Subscription{}, subscriptionRatelimitIndex,
		func(rawObj k8client.Object) []string {
			subscription := rawObj.(*cpv1alpha3.Subscription)
			var subscriptionRatelimit []string
			subscriptionRatelimit = append(subscriptionRatelimit,
				types.NamespacedName{
					Name:      string(subscription.Spec.RatelimitRef.Name),
					Namespace: subscription.Namespace,
				}.String())
			return subscriptionRatelimit
		}); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2610, logging.CRITICAL, "Error adding indexes: %v", err))
		return err
	}
	if err := mgr.GetFieldIndexer().IndexField(ctx, &cpv1alpha3.Subscription{}, subscriptionToAIRatelimitIndex,
		func(rawObj k8client.Object) []string {
			subscription := rawObj.(*cpv1alpha3.Subscription)
			var aiRatelimits []string
			aiRatelimits = append(aiRatelimits,
				types.NamespacedName{
					Name:      string(subscription.Spec.RatelimitRef.Name),
					Namespace: subscription.Namespace,
				}.String())
			return aiRatelimits
		}); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2610, logging.CRITICAL, "Error adding indexes: %v", err))
		return err
	}
	return nil
}

// getApplicationMappingsForSubscription triggers the ApplicationMapping controller reconcile method based on the changes detected
// from Subscription objects. If the changes are done for an API stored in the Operator Data store,
func (subscriptionReconciler *SubscriptionReconciler) getSubscriptionForRatelimit(ctx context.Context, obj *dpv1alpha3.RateLimitPolicy) []reconcile.Request {
	ratelimit := obj
	subList := &cpv1alpha3.SubscriptionList{}
	if err := subscriptionReconciler.client.List(ctx, subList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(subscriptionIndex, utils.NamespacedName(ratelimit).String()),
	}); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2623, logging.CRITICAL, "Unable to find associated Application mappings: %s", utils.NamespacedName(ratelimit).String()))
		return []reconcile.Request{}
	}

	if len(subList.Items) == 0 {
		loggers.LoggerAPKOperator.Debugf("ApplicationMappings for Subscription %s/%s not found", ratelimit.Namespace, ratelimit.Name)
		return []reconcile.Request{}
	}

	requests := []reconcile.Request{}
	for _, subscription := range subList.Items {
		req := reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      subscription.Name,
				Namespace: subscription.Namespace},
		}
		requests = append(requests, req)
		loggers.LoggerAPKOperator.Debugf("Adding reconcile request for ApplicationMapping: %s/%s with Subscription UUID: %v", subscription.Namespace, subscription.Name,
			string(subscription.ObjectMeta.UID))
	}
	return requests
}

// getSubscriptionForAIRatelimit get the associated subscription reconcile request for a AIRatelimit resource change
func (subscriptionReconciler *SubscriptionReconciler) getSubscriptionForAIRatelimit(ctx context.Context, obj *dpv1alpha3.AIRateLimitPolicy) []reconcile.Request {
	airatelimit := obj
	subList := &cpv1alpha3.SubscriptionList{}
	if err := subscriptionReconciler.client.List(ctx, subList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(subscriptionToAIRatelimitIndex, utils.NamespacedName(airatelimit).String()),
	}); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2623, logging.CRITICAL, "Unable to find associated AI ratelimits: %s", utils.NamespacedName(airatelimit).String()))
		return []reconcile.Request{}
	}

	if len(subList.Items) == 0 {
		loggers.LoggerAPKOperator.Debugf("AIRatelimit for Subscription %s/%s not found", airatelimit.Namespace, airatelimit.Name)
		return []reconcile.Request{}
	}

	requests := []reconcile.Request{}
	for _, subscription := range subList.Items {
		req := reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      subscription.Name,
				Namespace: subscription.Namespace},
		}
		requests = append(requests, req)
		loggers.LoggerAPKOperator.Debugf("Adding reconcile request for AIratelimit policy: %s/%s with Subscription UUID: %v", subscription.Namespace, subscription.Name,
			string(subscription.ObjectMeta.UID))
	}
	return requests
}
