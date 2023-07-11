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

	logger "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	gwapiv1b1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/wso2/apk/adapter/pkg/loggers"
	"github.com/wso2/apk/adapter/pkg/logging"
	"github.com/wso2/apk/adapter/pkg/utils/envutils"
	"github.com/wso2/apk/adapter/pkg/utils/stringutils"
	dpv1alpha1 "github.com/wso2/apk/common-adapter/internal/operator/api/v1alpha1"
	constants "github.com/wso2/apk/common-adapter/internal/operator/constant"
	xds "github.com/wso2/apk/common-adapter/internal/xds"
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
	logger.Info("Incoming 1")
	c, err := controller.New(constants.RatelimitController, mgr, controller.Options{Reconciler: ratelimitReconsiler})
	if err != nil {
		loggers.LoggerAuth.ErrorC(logging.GetErrorByCode(3111, err.Error()))
		return err
	}

	if err := c.Watch(&source.Kind{Type: &dpv1alpha1.RateLimitPolicy{}}, &handler.EnqueueRequestForObject{}); err != nil {
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
	logger.Info("Incoming ")
	ratelimitKey := req.NamespacedName
	var api dpv1alpha1.API
	var vhost []string
	var resolveResourceList []dpv1alpha1.ResolveResource
	var resolveRatelimit dpv1alpha1.ResolveRateLimitAPIPolicy
	var ratelimitPolicy dpv1alpha1.RateLimitPolicy
	if err := ratelimitReconsiler.client.Get(ctx, ratelimitKey, &ratelimitPolicy); err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to get Ratelimit policy %s/%s",
			ratelimitKey.Namespace, ratelimitKey.Name)
	}
	logger.Info(" xxxxxxxxxxxxxx", ratelimitPolicy.Spec.TargetRef.Kind)
	// API Level Rate limit policy
	if ratelimitPolicy.Spec.TargetRef.Kind == constants.KindAPI {
		logger.Info("API Kind ")
		ratelimitReconsiler.client.Get(ctx, types.NamespacedName{
			Namespace: ratelimitKey.Namespace,
			Name:      string(ratelimitPolicy.Spec.TargetRef.Name)},
			&api)

		var organization = api.Spec.Organization
		var context = api.Spec.Context
		var httpRoute gwapiv1b1.HTTPRoute
		logger.Info("context ", context)
		if len(api.Spec.Production) > 0 {
			for _, ref := range api.Spec.Production[0].HTTPRouteRefs {
				if ref != "" {
					ratelimitReconsiler.client.Get(ctx, types.NamespacedName{
						Namespace: ratelimitKey.Namespace,
						Name:      ref},
						&httpRoute)
					for _, hostName := range httpRoute.Spec.Hostnames {
						vhost = append(vhost, string(hostName))
					}
				}
			}
		}
		if len(api.Spec.Sandbox) > 0 {
			for _, ref := range api.Spec.Sandbox[0].HTTPRouteRefs {
				if ref != "" {
					ratelimitReconsiler.client.Get(ctx, types.NamespacedName{
						Namespace: ratelimitKey.Namespace,
						Name:      ref},
						&httpRoute)
					for _, hostName := range httpRoute.Spec.Hostnames {
						vhost = append(vhost, string(hostName))
					}
				}
			}
		}
		resolveRatelimit.API.RequestsPerUnit = ratelimitPolicy.Spec.Default.API.RateLimit.RequestsPerUnit
		resolveRatelimit.API.Unit = ratelimitPolicy.Spec.Default.API.RateLimit.Unit
		resolveRatelimit.Vhost = vhost
		resolveRatelimit.Organization = organization
		resolveRatelimit.Context = context
		resolveRatelimit.UUID = string(api.ObjectMeta.UID)
	}

	// Resource Level Rate limit policy
	if ratelimitPolicy.Spec.TargetRef.Kind == constants.KindResource {
		ratelimitReconsiler.client.Get(ctx, types.NamespacedName{
			Namespace: ratelimitKey.Namespace,
			Name:      string(ratelimitPolicy.Spec.TargetRef.Name)},
			&api)
		var organization = api.Spec.Organization
		var context = api.Spec.Context
		var httpRoute gwapiv1b1.HTTPRoute
		if len(api.Spec.Production) > 0 {
			for _, ref := range api.Spec.Production[0].HTTPRouteRefs {
				if ref != "" {
					ratelimitReconsiler.client.Get(ctx, types.NamespacedName{
						Namespace: ratelimitKey.Namespace,
						Name:      ref},
						&httpRoute)
					for _, rule := range httpRoute.Spec.Rules {
						for _, filter := range rule.Filters {
							if filter.ExtensionRef != nil {
								if filter.ExtensionRef.Kind == constants.KindRateLimitPolicy {
									var resolveResource dpv1alpha1.ResolveResource
									resolveResource.Path = *rule.Matches[0].Path.Value
									resolveResource.Method = string(*rule.Matches[0].Method)
									resolveResource.PathMatchType = *rule.Matches[0].Path.Type
									resolveResourceList = append(resolveResourceList, resolveResource)
								}
							}
						}

					}
				}
			}
		}
		if len(api.Spec.Sandbox) > 0 {
			for _, ref := range api.Spec.Sandbox[0].HTTPRouteRefs {
				if ref != "" {
					ratelimitReconsiler.client.Get(ctx, types.NamespacedName{
						Namespace: ratelimitKey.Namespace,
						Name:      ref},
						&httpRoute)
					for _, rule := range httpRoute.Spec.Rules {
						for _, filter := range rule.Filters {
							if filter.ExtensionRef != nil {
								if filter.ExtensionRef.Kind == constants.KindRateLimitPolicy {
									var resolveResource dpv1alpha1.ResolveResource
									resolveResource.Path = *rule.Matches[0].Path.Value
									resolveResource.Method = string(*rule.Matches[0].Method)
									resolveResource.PathMatchType = *rule.Matches[0].Path.Type
									resolveResourceList = append(resolveResourceList, resolveResource)
								}
							}
						}

					}
				}
			}
		}
		resolveRatelimit.Organization = organization
		resolveRatelimit.Context = context
		resolveRatelimit.UUID = string(api.ObjectMeta.UID)
		resolveRatelimit.Resources = resolveResourceList
		resolveRatelimit.API.RequestsPerUnit = ratelimitPolicy.Spec.Default.API.RateLimit.RequestsPerUnit
		resolveRatelimit.API.Unit = ratelimitPolicy.Spec.Default.API.RateLimit.Unit
	}

	xds.UpdateRateLimitXDSCache(vhost, resolveRatelimit)
	xds.UpdateRateLimiterPolicies("default")

	return ctrl.Result{}, nil
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

// SetupWithManager sets up the controller with the Manager.
func (ratelimitReconsiler *RateLimitPolicyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dpv1alpha1.RateLimitPolicy{}).
		Complete(ratelimitReconsiler)
}
