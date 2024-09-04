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

package dp

import (
	"context"
	"fmt"
	"time"

	logger "github.com/sirupsen/logrus"
	k8error "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/fields"
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

	k8client "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/source"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/wso2/apk/adapter/pkg/logging"
	cache "github.com/wso2/apk/common-controller/internal/cache"
	"github.com/wso2/apk/common-controller/internal/config"
	loggers "github.com/wso2/apk/common-controller/internal/loggers"
	"github.com/wso2/apk/common-controller/internal/utils"
	xds "github.com/wso2/apk/common-controller/internal/xds"
	dpv1alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha1"
	dpv1alpha2 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha2"
	"github.com/wso2/apk/common-go-libs/constants"
)

// RateLimitPolicyReconciler reconciles a RateLimitPolicy object
type RateLimitPolicyReconciler struct {
	client client.Client
	ods    *cache.RatelimitDataStore
	Scheme *runtime.Scheme
}

const (
	// apiRateLimitIndex Index for API level ratelimits
	apiRateLimitIndex = "apiRateLimitIndex"
	// apiRateLimitResourceIndex Index for resource level ratelimits
	httprouteRateLimitIndex = "httprouteRateLimitIndex"
)

// NewratelimitController creates a new ratelimitcontroller instance.
func NewratelimitController(mgr manager.Manager, ratelimitStore *cache.RatelimitDataStore) error {
	ratelimitReconsiler := &RateLimitPolicyReconciler{
		client: mgr.GetClient(),
		ods:    ratelimitStore,
	}

	ctx := context.Background()
	if err := addIndexes(ctx, mgr); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2612, logging.BLOCKER, "Error adding indexes: %v", err))
		return err
	}

	c, err := controller.New(constants.RatelimitController, mgr, controller.Options{Reconciler: ratelimitReconsiler})
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2663, logging.BLOCKER,
			"Error creating Ratelimit controller: %v", err.Error()))
		return err
	}

	conf := config.ReadConfigs()
	predicates := []predicate.Predicate{predicate.NewPredicateFuncs(utils.FilterByNamespaces(conf.CommonController.Operator.Namespaces))}

	if err := c.Watch(source.Kind(mgr.GetCache(), &dpv1alpha2.API{}),
		handler.EnqueueRequestsFromMapFunc(ratelimitReconsiler.getRatelimitForAPI), predicates...); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2611, logging.BLOCKER,
			"Error watching API resources: %v", err))
		return err
	}

	if err := c.Watch(source.Kind(mgr.GetCache(), &gwapiv1.HTTPRoute{}),
		handler.EnqueueRequestsFromMapFunc(ratelimitReconsiler.getRatelimitForHTTPRoute), predicates...); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2613, logging.BLOCKER,
			"Error watching HTTPRoute resources: %v", err))
		return err
	}

	if err := c.Watch(source.Kind(mgr.GetCache(), &dpv1alpha1.RateLimitPolicy{}), &handler.EnqueueRequestForObject{}, predicates...); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2639, logging.BLOCKER,
			"Error watching Ratelimit resources: %v", err.Error()))
		return err
	}

	loggers.LoggerAPKOperator.Debug("RatelimitPolicy Controller successfully started. Watching RatelimitPolicy Objects...")
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
	// Check whether the Ratelimit CR exist, if not consider as a DELETE event.
	loggers.LoggerAPKOperator.Infof("Reconciling ratelimit...")
	conf := config.ReadConfigs()
	ratelimitKey := req.NamespacedName
	var ratelimitPolicy dpv1alpha1.RateLimitPolicy

	// Check k8s RatelimitPolicy Availbility
	if err := ratelimitReconsiler.client.Get(ctx, ratelimitKey, &ratelimitPolicy); err != nil {
		resolveRateLimitAPIPolicyList, found := ratelimitReconsiler.ods.GetResolveRatelimitPolicy(req.NamespacedName)
		// If availble in cache Delete cache and xds
		if found && k8error.IsNotFound(err) {
			ratelimitReconsiler.ods.DeleteResolveRatelimitPolicy(req.NamespacedName)
			xds.DeleteAPILevelRateLimitPolicies(resolveRateLimitAPIPolicyList)
			xds.DeleteResourceLevelRateLimitPolicies(resolveRateLimitAPIPolicyList)
			xds.UpdateRateLimiterPolicies(conf.CommonController.Server.Label)
		}
		resolveCustomRateLimitPolicy, foundCustom := ratelimitReconsiler.ods.GetCachedCustomRatelimitPolicy(req.NamespacedName)
		if foundCustom && k8error.IsNotFound(err) {
			ratelimitReconsiler.ods.DeleteCachedCustomRatelimitPolicy(req.NamespacedName)
			logger.Debug("Deleting CustomRateLimitPolicy : ", resolveCustomRateLimitPolicy)
			xds.DeleteCustomRateLimitPolicies(resolveCustomRateLimitPolicy)
			xds.UpdateRateLimiterPolicies(conf.CommonController.Server.Label)
		}
		resolveSubscriptionRatelimitPolicy, foundSubscription := ratelimitReconsiler.ods.GetResolveSubscriptionRatelimitPolicy(req.NamespacedName)
		if foundSubscription && k8error.IsNotFound(err) {
			ratelimitReconsiler.ods.DeleteSubscriptionRatelimitPolicy(req.NamespacedName)
			logger.Debug("Deleting SubscriptionRateLimitPolicy : ", resolveSubscriptionRatelimitPolicy)
			xds.DeleteSubscriptionRateLimitPolicies(resolveSubscriptionRatelimitPolicy)
			xds.UpdateRateLimiterPolicies(conf.CommonController.Server.Label)
		}

		if k8error.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{
			RequeueAfter: time.Duration(1 * time.Second),
		}, nil
	}

	if ratelimitPolicy.Spec.Override != nil && ratelimitPolicy.Spec.Override.Custom != nil {
		var customRateLimitPolicy = ratelimitReconsiler.marshelCustomRateLimit(ctx, ratelimitKey, ratelimitPolicy)
		ratelimitReconsiler.ods.AddorUpdateCustomRatelimitToStore(ratelimitKey, customRateLimitPolicy)
		xds.UpdateRateLimitXDSCacheForCustomPolicies(customRateLimitPolicy)
		xds.UpdateRateLimiterPolicies(conf.CommonController.Server.Label)
	} else if ratelimitPolicy.Spec.Override != nil && ratelimitPolicy.Spec.Override.Subscription != nil {
		var resolveSubscriptionRatelimitPolicy = ratelimitReconsiler.marshelSubscriptionRateLimit(ratelimitPolicy)
		ratelimitReconsiler.ods.AddorUpdateResolveSubscriptionRatelimitToStore(ratelimitKey, resolveSubscriptionRatelimitPolicy)
		xds.UpdateRateLimitXDSCacheForSubscriptionPolicies(resolveSubscriptionRatelimitPolicy)
		xds.UpdateRateLimiterPolicies(conf.CommonController.Server.Label)
	} else {
		if resolveRatelimitPolicyList, err := ratelimitReconsiler.marshelRateLimit(ctx, ratelimitKey, ratelimitPolicy); err != nil {
			return ctrl.Result{}, err
		} else if len(resolveRatelimitPolicyList) > 0 {
			ratelimitReconsiler.ods.AddorUpdateResolveRatelimitToStore(ratelimitKey, resolveRatelimitPolicyList)
			xds.UpdateRateLimitXDSCache(resolveRatelimitPolicyList)
			xds.UpdateRateLimiterPolicies(conf.CommonController.Server.Label)
		}
	}

	return ctrl.Result{}, nil
}

func (ratelimitReconsiler *RateLimitPolicyReconciler) getRatelimitForAPI(ctx context.Context, obj k8client.Object) []reconcile.Request {
	api, ok := obj.(*dpv1alpha2.API)
	if !ok {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2622, logging.TRIVIAL,
			"Unexpected object type, bypassing reconciliation: %v", api))
		return []reconcile.Request{}
	}

	requests := []reconcile.Request{}

	ratelimitPolicyList := &dpv1alpha1.RateLimitPolicyList{}
	if err := ratelimitReconsiler.client.List(ctx, ratelimitPolicyList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(apiRateLimitIndex, NamespacedName(api).String()),
	}); err != nil {
		return []reconcile.Request{}
	}

	for item := range ratelimitPolicyList.Items {
		ratelimitPolicy := ratelimitPolicyList.Items[item]
		requests = append(requests, ratelimitReconsiler.AddRatelimitRequest(&ratelimitPolicy)...)
	}

	return requests
}

// AddRatelimitRequest adds a request to reconcile for the given ratelimit policy
func (ratelimitReconsiler *RateLimitPolicyReconciler) AddRatelimitRequest(obj k8client.Object) []reconcile.Request {
	ratelimitPolicy, ok := obj.(*dpv1alpha1.RateLimitPolicy)
	if !ok {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2622, logging.TRIVIAL,
			"Unexpected object type, bypassing reconciliation: %v", ratelimitPolicy))
		return nil
	}

	return []reconcile.Request{{
		NamespacedName: types.NamespacedName{
			Name:      string(ratelimitPolicy.Name),
			Namespace: ratelimitPolicy.Namespace,
		},
	}}
}

func (ratelimitReconsiler *RateLimitPolicyReconciler) getRatelimitForHTTPRoute(ctx context.Context, obj k8client.Object) []reconcile.Request {
	httpRoute, ok := obj.(*gwapiv1.HTTPRoute)
	if !ok {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2622, logging.TRIVIAL,
			"Unexpected object type, bypassing reconciliation: %v", httpRoute))
		return []reconcile.Request{}
	}

	requests := []reconcile.Request{}

	ratelimitPolicyList := &dpv1alpha1.RateLimitPolicyList{}
	if err := ratelimitReconsiler.client.List(ctx, ratelimitPolicyList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(httprouteRateLimitIndex, NamespacedName(httpRoute).String()),
	}); err != nil {
		return []reconcile.Request{}
	}
	for item := range ratelimitPolicyList.Items {
		ratelimitPolicy := ratelimitPolicyList.Items[item]
		requests = append(requests, ratelimitReconsiler.AddRatelimitRequest(&ratelimitPolicy)...)
	}

	return requests
}

func (ratelimitReconsiler *RateLimitPolicyReconciler) marshelSubscriptionRateLimit(
	ratelimitPolicy dpv1alpha1.RateLimitPolicy) dpv1alpha1.ResolveSubscriptionRatelimitPolicy {

	var resolveSubscriptionRatelimit dpv1alpha1.ResolveSubscriptionRatelimitPolicy
	resolveSubscriptionRatelimit.Name = ratelimitPolicy.Name
	resolveSubscriptionRatelimit.RequestCount.RequestsPerUnit = ratelimitPolicy.Spec.Override.Subscription.RequestCount.RequestsPerUnit
	resolveSubscriptionRatelimit.RequestCount.Unit = ratelimitPolicy.Spec.Override.Subscription.RequestCount.Unit
	if ratelimitPolicy.Spec.Override.Subscription.BurstControl != nil {
		resolveSubscriptionRatelimit.BurstControl.RequestsPerUnit = ratelimitPolicy.Spec.Override.Subscription.BurstControl.RequestsPerUnit
		resolveSubscriptionRatelimit.BurstControl.Unit = ratelimitPolicy.Spec.Override.Subscription.BurstControl.Unit
	}
	resolveSubscriptionRatelimit.StopOnQuotaReach = ratelimitPolicy.Spec.Override.Subscription.StopOnQuotaReach
	resolveSubscriptionRatelimit.Organization = ratelimitPolicy.Spec.Override.Subscription.Organization
	return resolveSubscriptionRatelimit
}

func (ratelimitReconsiler *RateLimitPolicyReconciler) marshelRateLimit(ctx context.Context, ratelimitKey types.NamespacedName,
	ratelimitPolicy dpv1alpha1.RateLimitPolicy) ([]dpv1alpha1.ResolveRateLimitAPIPolicy, error) {

	policyList := []dpv1alpha1.ResolveRateLimitAPIPolicy{}
	var api dpv1alpha2.API

	if err := ratelimitReconsiler.client.Get(ctx, types.NamespacedName{
		Namespace: ratelimitKey.Namespace,
		Name:      string(ratelimitPolicy.Spec.TargetRef.Name)},
		&api); err != nil {
		return nil, fmt.Errorf("error while getting API : %v, %s", string(ratelimitPolicy.Spec.TargetRef.Name), err.Error())
	}

	organization := api.Spec.Organization
	basePath := api.Spec.BasePath
	environment := utils.GetEnvironment(api.Spec.Environment)

	// API Level Rate limit policy
	if ratelimitPolicy.Spec.TargetRef.Kind == constants.KindAPI {

		var resolveRatelimit dpv1alpha1.ResolveRateLimitAPIPolicy
		if ratelimitPolicy.Spec.Override != nil {
			resolveRatelimit.API.RequestsPerUnit = ratelimitPolicy.Spec.Override.API.RequestsPerUnit
			resolveRatelimit.API.Unit = ratelimitPolicy.Spec.Override.API.Unit
		} else {
			resolveRatelimit.API.RequestsPerUnit = ratelimitPolicy.Spec.Default.API.RequestsPerUnit
			resolveRatelimit.API.Unit = ratelimitPolicy.Spec.Default.API.Unit
		}

		resolveRatelimit.Environment = environment
		resolveRatelimit.Organization = organization
		resolveRatelimit.BasePath = basePath
		resolveRatelimit.UUID = string(api.ObjectMeta.UID)
		policyList = append(policyList, resolveRatelimit)
	}

	// Resource Level Rate limit policy
	if ratelimitPolicy.Spec.TargetRef.Kind == constants.KindResource {

		var resolveRatelimit dpv1alpha1.ResolveRateLimitAPIPolicy
		resolveRatelimit.Organization = organization
		resolveRatelimit.BasePath = basePath
		resolveRatelimit.UUID = string(api.ObjectMeta.UID)
		resolveRatelimit.Environment = environment

		if len(api.Spec.Production) > 0 && api.Spec.APIType == "REST" {
			resolveResourceList, err := ratelimitReconsiler.getHTTPRouteResourceList(ctx, ratelimitKey, ratelimitPolicy,
				api.Spec.Production[0].RouteRefs)
			if err != nil {
				return nil, err
			}
			if len(resolveResourceList) > 0 {
				resolveRatelimit.Resources = resolveResourceList
				policyList = append(policyList, resolveRatelimit)
			}
		}

		if len(api.Spec.Sandbox) > 0 && api.Spec.APIType == "REST" {
			resolveResourceList, err := ratelimitReconsiler.getHTTPRouteResourceList(ctx, ratelimitKey, ratelimitPolicy,
				api.Spec.Sandbox[0].RouteRefs)
			if err != nil {
				return nil, err
			}
			if len(resolveResourceList) > 0 {
				resolveRatelimit.Resources = resolveResourceList
				resolveRatelimit.Environment += "_sandbox"
				policyList = append(policyList, resolveRatelimit)
			}
		}
	}

	return policyList, nil
}

func (ratelimitReconsiler *RateLimitPolicyReconciler) getHTTPRouteResourceList(ctx context.Context, ratelimitKey types.NamespacedName,
	ratelimitPolicy dpv1alpha1.RateLimitPolicy, httpRefs []string) ([]dpv1alpha1.ResolveResource, error) {

	var resolveResourceList []dpv1alpha1.ResolveResource
	var httpRoute gwapiv1.HTTPRoute

	for _, ref := range httpRefs {
		if ref != "" {
			if err := ratelimitReconsiler.client.Get(ctx, types.NamespacedName{
				Namespace: ratelimitKey.Namespace,
				Name:      ref},
				&httpRoute); err != nil {
				return nil, fmt.Errorf("error while getting HTTPRoute : %v for API : %v, %s", string(ref),
					string(ratelimitPolicy.Spec.TargetRef.Name), err.Error())
			}
			for _, rule := range httpRoute.Spec.Rules {
				for _, filter := range rule.Filters {
					if filter.ExtensionRef != nil {
						if filter.ExtensionRef.Kind == constants.KindRateLimitPolicy && string(filter.ExtensionRef.Name) == ratelimitPolicy.Name {
							var resolveResource dpv1alpha1.ResolveResource
							resolveResource.Path = *rule.Matches[0].Path.Value
							if rule.Matches[0].Method != nil {
								resolveResource.Method = string(*rule.Matches[0].Method)
							} else {
								resolveResource.Method = constants.All
							}
							resolveResource.PathMatchType = *rule.Matches[0].Path.Type
							if ratelimitPolicy.Spec.Override != nil {
								resolveResource.ResourceRatelimit.RequestsPerUnit = ratelimitPolicy.Spec.Override.API.RequestsPerUnit
								resolveResource.ResourceRatelimit.Unit = ratelimitPolicy.Spec.Override.API.Unit
							} else {
								resolveResource.ResourceRatelimit.RequestsPerUnit = ratelimitPolicy.Spec.Default.API.RequestsPerUnit
								resolveResource.ResourceRatelimit.Unit = ratelimitPolicy.Spec.Default.API.Unit
							}
							resolveResourceList = append(resolveResourceList, resolveResource)
						}
					}
				}

			}
		}
	}

	return resolveResourceList, nil
}

func (ratelimitReconsiler *RateLimitPolicyReconciler) marshelCustomRateLimit(ctx context.Context, ratelimitKey types.NamespacedName,
	ratelimitPolicy dpv1alpha1.RateLimitPolicy) dpv1alpha1.CustomRateLimitPolicyDef {
	var customRateLimitPolicy dpv1alpha1.CustomRateLimitPolicyDef
	// Custom Rate limit policy
	if ratelimitPolicy.Spec.TargetRef.Kind == constants.KindGateway {
		customRateLimitPolicy = getCustomRateLimitPolicy(&ratelimitPolicy)
		logger.Debug("CustomRateLimitPolicy : ", customRateLimitPolicy)
	}
	return customRateLimitPolicy
}

// getCustomRateLimitPolicy returns the custom rate limit policy.
func getCustomRateLimitPolicy(customRateLimitPolicy *dpv1alpha1.RateLimitPolicy) dpv1alpha1.CustomRateLimitPolicyDef {
	customRLPolicy := *dpv1alpha1.ParseCustomRateLimitPolicy(*customRateLimitPolicy)
	logger.Debug("customRLPolicy:", customRLPolicy)
	return customRLPolicy
}

func addIndexes(ctx context.Context, mgr manager.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(ctx, &gwapiv1.HTTPRoute{}, httprouteRateLimitIndex,
		func(rawObj k8client.Object) []string {
			httpRoute := rawObj.(*gwapiv1.HTTPRoute)
			var ratelimitPolicy []string
			for _, rule := range httpRoute.Spec.Rules {
				for _, filter := range rule.Filters {
					if filter.ExtensionRef != nil {
						if filter.ExtensionRef.Kind == constants.KindRateLimitPolicy {
							ratelimitPolicy = append(ratelimitPolicy,
								types.NamespacedName{
									Namespace: httpRoute.Namespace,
									Name:      string(filter.ExtensionRef.Name),
								}.String())
						}
					}
				}
			}
			return ratelimitPolicy
		}); err != nil {
		return err
	}

	// ratelimite policy to API indexer
	err := mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha1.RateLimitPolicy{}, apiRateLimitIndex,
		func(rawObj k8client.Object) []string {
			ratelimitPolicy := rawObj.(*dpv1alpha1.RateLimitPolicy)
			var apis []string
			apis = append(apis,
				types.NamespacedName{
					Namespace: ratelimitPolicy.Namespace,
					Name:      string(ratelimitPolicy.Spec.TargetRef.Name),
				}.String())
			return apis
		})
	return err
}

// NamespacedName generates namespaced name for Kubernetes objects
func NamespacedName(obj client.Object) types.NamespacedName {
	return types.NamespacedName{
		Namespace: obj.GetNamespace(),
		Name:      obj.GetName(),
	}
}

// GetNamespace reads namespace with a default value
func GetNamespace(namespace *gwapiv1.Namespace, defaultNamespace string) string {
	if namespace != nil && *namespace != "" {
		return string(*namespace)
	}
	return defaultNamespace
}

// SetupWithManager sets up the controller with the Manager.
func (ratelimitReconsiler *RateLimitPolicyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dpv1alpha1.RateLimitPolicy{}).
		Complete(ratelimitReconsiler)
}
