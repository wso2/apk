/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	gwapiv1b1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/wso2/apk/adapter/pkg/loggers"
	"github.com/wso2/apk/adapter/pkg/logging"
	"github.com/wso2/apk/adapter/pkg/utils/envutils"
	"github.com/wso2/apk/adapter/pkg/utils/stringutils"
	dpv1alpha1 "github.com/wso2/apk/common-adapter/api/v1alpha1"
	constants "github.com/wso2/apk/common-adapter/internal/constant"
)

const (
	apiRateLimitResourceIndex = "apiRateLimitResourceIndex"
	apiRateLimitIndex         = "apiRateLimitIndex"
)

// RateLimitPolicyReconciler reconciles a RateLimitPolicy object
type RateLimitPolicyReconciler struct {
	client client.Client
	Scheme *runtime.Scheme
}

// NewratelimitController creates a new ratelimitcontroller instance.
func NewratelimitController(mgr manager.Manager) error {
	ratelimitReconsiler := &RateLimitPolicyReconciler{
		client: mgr.GetClient(),
	}

	c, err := controller.New(constants.RatelimitController, mgr, controller.Options{Reconciler: ratelimitReconsiler})
	if err != nil {
		loggers.LoggerAuth.ErrorC(logging.GetErrorByCode(3111, err.Error()))
		return err
	}

	if err := c.Watch(&source.Kind{Type: &dpv1alpha1.RateLimitPolicy{}}, &handler.EnqueueRequestForObject{},
		predicate.NewPredicateFuncs(FilterByNamespaces([]string{GetOperatorPodNamespace()}))); err != nil {
		loggers.LoggerAuth.ErrorC(logging.GetErrorByCode(3112, err.Error()))
		return err
	}

	loggers.LoggerAuth.Debug("RatelimitPolicy Controller successfully started. Watching RatelimitPolicy Objects...")
	return nil
}

//+kubebuilder:rbac:groups=dp.wso2.com,resources=ratelimitpolicies,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=dp.wso2.com,resources=ratelimitpolicies/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=dp.wso2.com,resources=ratelimitpolicies/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the RateLimitPolicy object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.15.0/pkg/reconcile
func (ratelimitReconsiler *RateLimitPolicyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	ratelimitKey := req.NamespacedName
	ratelimitPolicyList := &dpv1alpha1.RateLimitPolicyList{}
	if err := ratelimitReconsiler.client.List(ctx, ratelimitPolicyList); err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to get Ratelimit policy %s/%s",
			ratelimitKey.Namespace, ratelimitKey.Name)
	}

	ratelimitPolicies, err := ratelimitReconsiler.getRatelimitPoliciesForAPI(ctx, ratelimitReconsiler.client, ratelimitKey)
	if err != nil {
		loggers.LoggerAuth.ErrorC(logging.GetErrorByCode(3111, err.Error()))
		return ctrl.Result{}, err
	}
	loggers.LoggerAuth.Debugf(ratelimitPolicies)
	return ctrl.Result{}, nil
}

func (ratelimitReconsiler *RateLimitPolicyReconciler) getRatelimitPoliciesForAPI(ctx context.Context,
	apiRef string) (map[string]dpv1alpha1.RateLimitPolicy, error) {
	ratelimitPolicies := make(map[string]dpv1alpha1.RateLimitPolicy)
	ratelimitPolicyList := &dpv1alpha1.RateLimitPolicyList{}
	if err := r.client.List(ctx, ratelimitPolicyList); err != nil {
		return nil, err
	}
	for _, item := range ratelimitPolicyList.Items {
		ratelimitPolicies[NamespacedName(&item).String()] = item
	}
	return ratelimitPolicies, nil
}

// FilterByNamespaces takes a list of namespaces and returns a filter function
// which return true if the input object is in the given namespaces list,
// and returns false otherwise
func FilterByNamespaces(namespaces []string) func(object client.Object) bool {
	return func(object client.Object) bool {
		if namespaces == nil {
			return true
		}
		return stringutils.StringInSlice(object.GetNamespace(), namespaces)
	}
}

// GetOperatorPodNamespace returns the namesapce of the operator pod
func GetOperatorPodNamespace() string {
	return envutils.GetEnv(constants.OperatorPodNamespace,
		constants.OperatorPodNamespaceDefaultValue)
}

// NamespacedName generates namespaced name for Kubernetes objects
func NamespacedName(obj client.Object) types.NamespacedName {
	return types.NamespacedName{
		Namespace: obj.GetNamespace(),
		Name:      obj.GetName(),
	}
}

// GetNamespace reads namespace with a default value
func GetNamespace(namespace *gwapiv1b1.Namespace, defaultNamespace string) string {
	if namespace != nil && *namespace != "" {
		return string(*namespace)
	}
	return defaultNamespace
}

// SetupWithManager sets up the controller with the Manager.
func (ratelimitReconsiler *RateLimitPolicyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dpv1alpha1.RateLimitPolicy{}).
		Complete(r)
}
