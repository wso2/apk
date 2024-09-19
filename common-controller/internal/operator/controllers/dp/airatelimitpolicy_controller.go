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

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	dpv1alpha3 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha3"
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
	Scheme *runtime.Scheme
	ods    *cache.RatelimitDataStore
}

// NewAIRatelimitController creates a new ratelimitcontroller instance.
func NewAIRatelimitController(mgr manager.Manager, ratelimitStore *cache.RatelimitDataStore) error {
	aiRateLimitPolicyReconciler := &AIRateLimitPolicyReconciler{
		Client: mgr.GetClient(),
		ods: ratelimitStore,
	}

	c, err := controller.New(constants.AIRatelimitController, mgr, controller.Options{Reconciler: aiRateLimitPolicyReconciler})
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2663, logging.BLOCKER,
			"Error creating Ratelimit controller: %v", err.Error()))
		return err
	}

	conf := config.ReadConfigs()
	predicates := []predicate.Predicate{predicate.NewPredicateFuncs(utils.FilterByNamespaces(conf.CommonController.Operator.Namespaces))}

	if err := c.Watch(source.Kind(mgr.GetCache(), &dpv1alpha3.AIRateLimitPolicy{}), &handler.EnqueueRequestForObject{}, predicates...); err != nil {
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
		xds.UpdateRateLimitXDSCacheForAIRatelimitPolicies(r.ods.GetAIRatelimitPolicySpecs())
		xds.UpdateRateLimiterPolicies(conf.CommonController.Server.Label)
	} else {
		loggers.LoggerAPKOperator.Infof("ratelimits found")
		if ratelimitPolicy.Spec.Override == nil {
			ratelimitPolicy.Spec.Override = ratelimitPolicy.Spec.Default
		}
		if ratelimitPolicy.Spec.TargetRef.Name != "" {
			r.ods.AddorUpdateAIRatelimitToStore(ratelimitKey, ratelimitPolicy.Spec)
			xds.UpdateRateLimitXDSCacheForAIRatelimitPolicies(r.ods.GetAIRatelimitPolicySpecs())
			xds.UpdateRateLimiterPolicies(conf.CommonController.Server.Label)
		} else {
			r.ods.DeleteAIRatelimitPolicySpec(ratelimitKey)
			xds.UpdateRateLimitXDSCacheForAIRatelimitPolicies(r.ods.GetAIRatelimitPolicySpecs())
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
