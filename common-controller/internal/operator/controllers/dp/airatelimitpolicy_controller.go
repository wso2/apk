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

package dp

import (
	"context"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	dpv1alpha3 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha3"
	dpv1alpha5 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha5"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/wso2/apk/adapter/pkg/logging"
	cache "github.com/wso2/apk/common-controller/internal/cache"
	"github.com/wso2/apk/common-controller/internal/config"
	loggers "github.com/wso2/apk/common-controller/internal/loggers"
	"github.com/wso2/apk/common-controller/internal/utils"
	xds "github.com/wso2/apk/common-controller/internal/xds"
	"github.com/wso2/apk/common-go-libs/constants"
)

// AIRateLimitPolicyReconciler reconciles a AIRateLimitPolicy object
type AIRateLimitPolicyReconciler struct {
	client.Client
	APIReader client.Reader
	Scheme    *runtime.Scheme
	ods       *cache.RatelimitDataStore
}

// NewAIRatelimitController creates a new ratelimitcontroller instance.
func NewAIRatelimitController(mgr manager.Manager, ratelimitStore *cache.RatelimitDataStore) error {
	aiRateLimitPolicyReconciler := &AIRateLimitPolicyReconciler{
		Client:    mgr.GetClient(),
		APIReader: mgr.GetAPIReader(),
		ods:       ratelimitStore,
	}

	c, err := controller.New(constants.AIRatelimitController, mgr, controller.Options{Reconciler: aiRateLimitPolicyReconciler})
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2663, logging.BLOCKER,
			"Error creating Ratelimit controller: %v", err.Error()))
		return err
	}

	conf := config.ReadConfigs()
	predicateAIRatelimitPolicy := []predicate.TypedPredicate[*dpv1alpha3.AIRateLimitPolicy]{predicate.NewTypedPredicateFuncs[*dpv1alpha3.AIRateLimitPolicy](utils.FilterAIRatelimitPolicyByNamespaces(conf.CommonController.Operator.Namespaces))}

	if err := c.Watch(source.Kind(mgr.GetCache(), &dpv1alpha3.AIRateLimitPolicy{}, &handler.TypedEnqueueRequestForObject[*dpv1alpha3.AIRateLimitPolicy]{}, predicateAIRatelimitPolicy...)); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2639, logging.BLOCKER,
			"Error watching Ratelimit resources: %v", err.Error()))
		return err
	}

	loggers.LoggerAPKOperator.Debug("RatelimitPolicy Controller successfully started. Watching RatelimitPolicy Objects...")
	return nil
}

//+kubebuilder:rbac:groups=dp.wso2.com,resources=airatelimitpolicies,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=dp.wso2.com,resources=airatelimitpolicies/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=dp.wso2.com,resources=airatelimitpolicies/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the AIRateLimitPolicy object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.4/pkg/reconcile
func (r *AIRateLimitPolicyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	loggers.LoggerAPKOperator.Infof("AIRatelimit reconcile...")
	// TODO(user): your logic here
	ratelimitKey := req.NamespacedName
	var ratelimitPolicy dpv1alpha3.AIRateLimitPolicy
	conf := config.ReadConfigs()

	// Check k8s RatelimitPolicy Availbility
	if err := r.Client.Get(ctx, ratelimitKey, &ratelimitPolicy); err != nil {
		loggers.LoggerAPKOperator.Errorf("Error retrieving AIRatelimit")
		// It could be deletion event. So lets try to delete the related entried from the ods and update xds
		r.ods.DeleteAIRatelimitPolicySpec(ratelimitKey)
		r.ods.DeleteSubscriptionBasedAIRatelimitPolicySpec(ratelimitKey)
		xds.UpdateRateLimitXDSCacheForAIRatelimitPolicies(r.ods.GetAIRatelimitPolicySpecs())
		xds.UpdateRateLimitXDSCacheForSubscriptionBasedAIRatelimitPolicies(r.ods.GetSubscriptionBasedAIRatelimitPolicySpecs())
		xds.UpdateRateLimiterPolicies(conf.CommonController.Server.Label)
	} else {
		loggers.LoggerAPKOperator.Infof("ratelimits found")
		if ratelimitPolicy.Spec.Override == nil {
			ratelimitPolicy.Spec.Override = ratelimitPolicy.Spec.Default
		}
		if ratelimitPolicy.Spec.TargetRef.Kind == "Backend" {
			var backend dpv1alpha5.Backend
			var ns string
			if ratelimitPolicy.Spec.TargetRef.Namespace != nil {
				ns = string(*ratelimitPolicy.Spec.TargetRef.Namespace)
			} else {
				ns = ratelimitPolicy.Namespace
			}
			if err := r.Client.Get(ctx, types.NamespacedName{Namespace: ns, Name: string(ratelimitPolicy.Spec.TargetRef.Name)}, &backend); err != nil {
				// If not found in cache, fallback to direct API call
				loggers.LoggerAPKOperator.Errorf("Error retrieving Backend: %v", err)
				if err := r.APIReader.Get(ctx, types.NamespacedName{Namespace: ns, Name: string(ratelimitPolicy.Spec.TargetRef.Name)}, &backend); err != nil {
					loggers.LoggerAPKOperator.Errorf("Error retrieving Backend directly: %v", err)
				}
			}
			if backend.Name != "" {
				loggers.LoggerAPKOperator.Infof("Backend found: %s", backend.Name)
				// Prepare owner references for the route
				if len(backend.GetOwnerReferences()) == 1 && backend.GetOwnerReferences()[0].Kind == "API" {
					loggers.LoggerAPKOperator.Infof("Owner references found for Backend: %s", backend.Name)
					preparedOwnerReferences := backend.GetOwnerReferences()[0]
					// Decide whether we need an update
					updateRequired := false
					if len(ratelimitPolicy.GetOwnerReferences()) != 1 {
						updateRequired = true
					} else {
						_, found := FindElement(ratelimitPolicy.GetOwnerReferences(), func(refLocal metav1.OwnerReference) bool {
							return refLocal.UID == preparedOwnerReferences.UID && refLocal.Name == preparedOwnerReferences.Name && refLocal.APIVersion == preparedOwnerReferences.APIVersion && refLocal.Kind == preparedOwnerReferences.Kind
						})
						if !found {
							updateRequired = true
						}
					}
					if updateRequired {
						ratelimitPolicy.SetOwnerReferences([]metav1.OwnerReference{preparedOwnerReferences})
						utils.UpdateCR(ctx, r.Client, &ratelimitPolicy)
					}
				}
			}
			r.ods.AddorUpdateAIRatelimitToStore(ratelimitKey, ratelimitPolicy.Spec)
			xds.UpdateRateLimitXDSCacheForAIRatelimitPolicies(r.ods.GetAIRatelimitPolicySpecs())
			xds.UpdateRateLimiterPolicies(conf.CommonController.Server.Label)
		} else if strings.EqualFold(string(ratelimitPolicy.Spec.TargetRef.Kind), "Subscription") {
			r.ods.AddorUpdateSubscriptionBasedAIRatelimitToStore(ratelimitKey, ratelimitPolicy.Spec)
			xds.UpdateRateLimitXDSCacheForSubscriptionBasedAIRatelimitPolicies(r.ods.GetSubscriptionBasedAIRatelimitPolicySpecs())
			xds.UpdateRateLimiterPolicies(conf.CommonController.Server.Label)
		} else {
			r.ods.DeleteAIRatelimitPolicySpec(ratelimitKey)
			r.ods.DeleteSubscriptionBasedAIRatelimitPolicySpec(ratelimitKey)
			xds.UpdateRateLimitXDSCacheForAIRatelimitPolicies(r.ods.GetAIRatelimitPolicySpecs())
			xds.UpdateRateLimitXDSCacheForSubscriptionBasedAIRatelimitPolicies(r.ods.GetSubscriptionBasedAIRatelimitPolicySpecs())
			xds.UpdateRateLimiterPolicies(conf.CommonController.Server.Label)
		}
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AIRateLimitPolicyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dpv1alpha3.AIRateLimitPolicy{}).
		Complete(r)
}

// FindElement searches for an element in a slice based on a given predicate.
// It returns the element and true if the element was found.
func FindElement[T any](collection []T, predicate func(item T) bool) (T, bool) {
	for _, item := range collection {
		if predicate(item) {
			return item, true
		}
	}
	var dummy T
	return dummy, false
}
