/*
 *  Copyright (c) 2022, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/wso2/apk/adapter/config"
	"github.com/wso2/apk/adapter/internal/loggers"
	"github.com/wso2/apk/adapter/pkg/logging"
	"github.com/wso2/apk/adapter/pkg/operator/constants"
	"github.com/wso2/apk/adapter/pkg/operator/status"
	"github.com/wso2/apk/adapter/pkg/operator/synchronizer"
	"github.com/wso2/apk/adapter/pkg/operator/utils"
	k8error "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	gwapiv1b1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	ctrl "sigs.k8s.io/controller-runtime"
	k8client "sigs.k8s.io/controller-runtime/pkg/client"

	dpv1alpha1 "github.com/wso2/apk/adapter/pkg/operator/apis/dp/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	httpRouteAPIIndex = "httpRouteAPIIndex"
	// Index for API level authentications
	httpRouteAuthenticationIndex = "httpRouteAuthenticationIndex"
	// Index for resource level authentications
	httpRouteAuthenticationResourceIndex = "httpRouteAuthenticationResourceIndex"
	httpRouteRateLimitIndex              = "httpRouteRateLimitIndex"
	httpRouteRateLimitResourceIndex      = "httpRouteRateLimitResourceIndex"
	// Index for API level apipolicies
	httpRouteAPIPolicyIndex = "httpRouteAPIPolicyIndex"
	// Index for resource level apipolicies
	httpRouteAPIPolicyResourceIndex = "httpRouteAPIPolicyResourceIndex"
	serviceHTTPRouteIndex           = "serviceHTTPRouteIndex"
	apiScopeIndex                   = "apiScopeIndex"
	configMapBackend                = "configMapBackend"
	secretBackend                   = "secretBackend"
	backendHTTPRouteIndex           = "backendHTTPRouteIndex"
)

// APIReconciler reconciles a API object
type APIReconciler struct {
	client        k8client.Client
	ods           *synchronizer.OperatorDataStore
	ch            *chan synchronizer.APIEvent
	statusUpdater *status.UpdateHandler
	mgr           manager.Manager
}

// NewAPIController creates a new API controller instance. API Controllers watches for dpv1alpha1.API and gwapiv1b1.HTTPRoute.
func NewAPIController(mgr manager.Manager, operatorDataStore *synchronizer.OperatorDataStore, statusUpdater *status.UpdateHandler,
	ch *chan synchronizer.APIEvent) error {
	r := &APIReconciler{
		client:        mgr.GetClient(),
		ods:           operatorDataStore,
		ch:            ch,
		statusUpdater: statusUpdater,
		mgr:           mgr,
	}
	c, err := controller.New(constants.APIController, mgr, controller.Options{Reconciler: r})
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2610, err))
		return err
	}
	ctx := context.Background()

	conf := config.ReadConfigs()
	predicates := []predicate.Predicate{predicate.NewPredicateFuncs(utils.FilterByNamespaces(conf.Adapter.Operator.Namespaces))}

	if err := c.Watch(&source.Kind{Type: &dpv1alpha1.API{}}, &handler.EnqueueRequestForObject{},
		predicates...); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2611, err))
		return err
	}
	if err := addIndexes(ctx, mgr); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2612, err))
		return err
	}

	if err := c.Watch(&source.Kind{Type: &gwapiv1b1.HTTPRoute{}}, handler.EnqueueRequestsFromMapFunc(r.getAPIForHTTPRoute),
		predicates...); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2613, err))
		return err
	}

	if err := c.Watch(&source.Kind{Type: &dpv1alpha1.Backend{}}, handler.EnqueueRequestsFromMapFunc(r.getAPIsForBackend),
		predicates...); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2615, err))
		return err
	}

	if err := c.Watch(&source.Kind{Type: &dpv1alpha1.Authentication{}}, handler.EnqueueRequestsFromMapFunc(r.getAPIsForAuthentication),
		predicates...); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2616, err))
		return err
	}

	if err := c.Watch(&source.Kind{Type: &dpv1alpha1.APIPolicy{}}, handler.EnqueueRequestsFromMapFunc(r.getAPIsForAPIPolicy),
		predicates...); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2617, err))
		return err
	}

	if err := c.Watch(&source.Kind{Type: &dpv1alpha1.Scope{}}, handler.EnqueueRequestsFromMapFunc(r.getAPIsForScope),
		predicates...); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2618, err))
		return err
	}

	if err := c.Watch(&source.Kind{Type: &corev1.ConfigMap{}}, handler.EnqueueRequestsFromMapFunc(r.getAPIsForConfigMap),
		predicates...); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2644, err))
		return err
	}

	if err := c.Watch(&source.Kind{Type: &corev1.Secret{}}, handler.EnqueueRequestsFromMapFunc(r.getAPIsForSecret),
		predicates...); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2645, err))
		return err
	}

	loggers.LoggerAPKOperator.Info("API Controller successfully started. Watching API Objects....")
	return nil
}

// +kubebuilder:rbac:groups=dp.wso2.com,resources=apis,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=dp.wso2.com,resources=apis/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=dp.wso2.com,resources=apis/finalizers,verbs=update
// +kubebuilder:rbac:groups=dp.wso2.com,resources=httproutes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=dp.wso2.com,resources=httproutes/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=dp.wso2.com,resources=httproutes/finalizers,verbs=update
// +kubebuilder:rbac:groups=dp.wso2.com,resources=authentications,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=dp.wso2.com,resources=authentications/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=dp.wso2.com,resources=authentications/finalizers,verbs=update
// +kubebuilder:rbac:groups=dp.wso2.com,resources=apipolicies,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=dp.wso2.com,resources=apipolicies/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=dp.wso2.com,resources=apipolicies/finalizers,verbs=update
// +kubebuilder:rbac:groups=dp.wso2.com,resources=scopes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=dp.wso2.com,resources=scopes/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=dp.wso2.com,resources=scopes/finalizers,verbs=update
// +kubebuilder:rbac:groups=dp.wso2.com,resources=ratelimitpolicies,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=dp.wso2.com,resources=ratelimitpolicies/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=dp.wso2.com,resources=ratelimitpolicies/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (apiReconciler *APIReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	loggers.LoggerAPKOperator.Infof("Reconciling for API %s", req.NamespacedName.String())

	// Check whether the API CR exist, if not consider as a DELETE event.
	var apiDef dpv1alpha1.API
	if err := apiReconciler.client.Get(ctx, req.NamespacedName, &apiDef); err != nil {
		apiState, found := apiReconciler.ods.GetCachedAPI(req.NamespacedName)
		if found && k8error.IsNotFound(err) {
			// The api doesn't exist in the api Cache, remove it
			apiReconciler.ods.DeleteCachedAPI(req.NamespacedName)
			loggers.LoggerAPKOperator.Infof("Delete event has received for API : %s, hence deleted from API cache", req.NamespacedName.String())
			*apiReconciler.ch <- synchronizer.APIEvent{EventType: constants.Delete, Event: apiState}
			return ctrl.Result{}, nil
		}
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2619, req.NamespacedName.String(), err))
		return ctrl.Result{}, nil
	}

	// Retrieve HTTPRoutes
	prodHTTPRoute, sandHTTPRoute, err := apiReconciler.resolveAPIRefs(ctx, apiDef.Spec.ProdHTTPRouteRefs,
		apiDef.Spec.SandHTTPRouteRefs, req.NamespacedName.String(), req.Namespace)
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2620, req.NamespacedName.String(), err))
		return ctrl.Result{}, err
	}
	loggers.LoggerAPKOperator.Debugf("HTTPRoutes are retrieved successfully for API CR %s", req.NamespacedName.String())

	if !apiDef.Status.Accepted {
		apiState := apiReconciler.ods.AddAPIState(apiDef, prodHTTPRoute, sandHTTPRoute)
		*apiReconciler.ch <- synchronizer.APIEvent{EventType: constants.Create, Event: apiState}
		//TODO(amali) update status only after deployed without errors
		apiReconciler.handleStatus(req.NamespacedName, constants.DeployedState, []string{})
	} else if cachedAPI, events, updated :=
		apiReconciler.ods.UpdateAPIState(&apiDef, prodHTTPRoute, sandHTTPRoute); updated {
		*apiReconciler.ch <- synchronizer.APIEvent{EventType: constants.Update, Event: cachedAPI}
		apiReconciler.handleStatus(req.NamespacedName, constants.UpdatedState, events)
	}
	return ctrl.Result{}, nil
}

// resolveAPIRefs validates following references related to the API
// - HTTPRoutes
func (apiReconciler *APIReconciler) resolveAPIRefs(ctx context.Context, prodHTTPRouteRef, sandHTTPRouteRef []string,
	api, namespace string) (*synchronizer.HTTPRouteState, *synchronizer.HTTPRouteState, error) {
	prodHTTPRoute := &synchronizer.HTTPRouteState{
		HTTPRoute: &gwapiv1b1.HTTPRoute{},
	}
	sandHTTPRoute := &synchronizer.HTTPRouteState{
		HTTPRoute: &gwapiv1b1.HTTPRoute{},
	}
	var err error
	if len(prodHTTPRouteRef) > 0 {
		if prodHTTPRoute, err = apiReconciler.resolveHTTPRouteRefs(ctx, prodHTTPRouteRef, namespace, api); err != nil {
			return nil, nil, fmt.Errorf("error while resolving production httpRouteref %s in namespace :%s has not found. %s",
				prodHTTPRouteRef, namespace, err.Error())
		}
	}

	if len(sandHTTPRouteRef) > 0 {
		if sandHTTPRoute, err = apiReconciler.resolveHTTPRouteRefs(ctx, sandHTTPRouteRef, namespace, api); err != nil {
			return nil, nil, fmt.Errorf("error while resolving sandbox httpRouteref %s in namespace :%s has not found. %s",
				sandHTTPRouteRef, namespace, err.Error())
		}
	}
	return prodHTTPRoute, sandHTTPRoute, nil
}

// resolveHTTPRouteRefs validates following references related to the API
// - Authentications
func (apiReconciler *APIReconciler) resolveHTTPRouteRefs(ctx context.Context, httpRouteRef []string, namespace, api string) (*synchronizer.HTTPRouteState, error) {
	httpRouteState := &synchronizer.HTTPRouteState{
		HTTPRoute: &gwapiv1b1.HTTPRoute{},
	}
	var err error
	httpRouteState.HTTPRoute, err = apiReconciler.concatHTTPRoutes(ctx, httpRouteRef, namespace)
	if err != nil {
		return nil, err
	}
	if httpRouteState.Authentications, err = apiReconciler.getAuthenticationsForHTTPRoute(ctx, httpRouteState.HTTPRoute); err != nil {
		return nil, fmt.Errorf("error while getting httproute auth : %s in namespace :%s, %s", httpRouteRef,
			namespace, err.Error())
	}
	if httpRouteState.ResourceAuthentications, err = apiReconciler.getAuthenticationsForResources(ctx, httpRouteState.HTTPRoute); err != nil {
		return nil, fmt.Errorf("error while getting httproute resource auth : %s in namespace :%s, %s", httpRouteRef,
			namespace, err.Error())
	}
	if httpRouteState.RateLimitPolicies, err = apiReconciler.getRatelimitPoliciesForHTTPRoute(ctx, httpRouteState.HTTPRoute); err != nil {
		return nil, fmt.Errorf("error while getting httproute ratelimit : %s in namespace :%s, %s", httpRouteRef,
			namespace, err.Error())
	}
	if httpRouteState.ResourceRateLimitPolicies, err = apiReconciler.getRatelimitPoliciesForResources(ctx, httpRouteState.HTTPRoute); err != nil {
		return nil, fmt.Errorf("error while getting httproute resource ratelimit : %s in namespace :%s, %s", httpRouteRef,
			namespace, err.Error())
	}
	if httpRouteState.APIPolicies, err = apiReconciler.getAPIPoliciesForHTTPRoute(ctx, httpRouteState.HTTPRoute); err != nil {
		return nil, fmt.Errorf("error while getting httproute apipolicy : %s in namespace :%s, %s", httpRouteRef,
			namespace, err.Error())
	}
	if httpRouteState.ResourceAPIPolicies, err = apiReconciler.getAPIPoliciesForResources(ctx, httpRouteState.HTTPRoute); err != nil {
		return nil, fmt.Errorf("error while getting httproute resource apipolicy %s in namespace :%s, %s", httpRouteRef,
			namespace, err.Error())
	}
	httpRouteState.BackendMapping = apiReconciler.getResolvedBackendsMapping(ctx, httpRouteState.HTTPRoute)
	httpRouteState.Scopes, err = apiReconciler.getScopesForHTTPRoute(ctx, httpRouteState.HTTPRoute, api)

	return httpRouteState, err
}

func (apiReconciler *APIReconciler) concatHTTPRoutes(ctx context.Context, httpRouteRefs []string,
	namespace string) (*gwapiv1b1.HTTPRoute, error) {
	var combinedHTTPRoute *gwapiv1b1.HTTPRoute
	for _, httpRouteRef := range httpRouteRefs {
		var httpRoute gwapiv1b1.HTTPRoute
		if err := apiReconciler.client.Get(ctx, types.NamespacedName{Namespace: namespace, Name: httpRouteRef},
			&httpRoute); err != nil {
			return nil, fmt.Errorf("error while getting httproute %s in namespace :%s, %s", httpRouteRef, namespace, err.Error())
		}
		if combinedHTTPRoute == nil {
			combinedHTTPRoute = &httpRoute
		} else {
			combinedHTTPRoute.Spec.Rules = append(combinedHTTPRoute.Spec.Rules, httpRoute.Spec.Rules...)
		}
	}
	return combinedHTTPRoute, nil
}

func (apiReconciler *APIReconciler) getAuthenticationsForHTTPRoute(ctx context.Context,
	httpRoute *gwapiv1b1.HTTPRoute) (map[string]dpv1alpha1.Authentication, error) {
	authentications := make(map[string]dpv1alpha1.Authentication)
	authenticationList := &dpv1alpha1.AuthenticationList{}
	if err := apiReconciler.client.List(ctx, authenticationList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(httpRouteAuthenticationIndex, utils.NamespacedName(httpRoute).String()),
	}); err != nil {
		return nil, err
	}
	for _, item := range authenticationList.Items {
		authentications[utils.NamespacedName(&item).String()] = item
	}
	return authentications, nil
}
func (apiReconciler *APIReconciler) getRatelimitPoliciesForHTTPRoute(ctx context.Context,
	httpRoute *gwapiv1b1.HTTPRoute) (map[string]dpv1alpha1.RateLimitPolicy, error) {
	ratelimitPolicies := make(map[string]dpv1alpha1.RateLimitPolicy)
	ratelimitPolicyList := &dpv1alpha1.RateLimitPolicyList{}
	if err := apiReconciler.client.List(ctx, ratelimitPolicyList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(httpRouteRateLimitIndex, utils.NamespacedName(httpRoute).String()),
	}); err != nil {
		return nil, err
	}
	for _, item := range ratelimitPolicyList.Items {
		ratelimitPolicies[utils.NamespacedName(&item).String()] = item
	}
	return ratelimitPolicies, nil
}

func (apiReconciler *APIReconciler) getScopesForHTTPRoute(ctx context.Context,
	httpRoute *gwapiv1b1.HTTPRoute, api string) (map[string]dpv1alpha1.Scope, error) {
	scopes := make(map[string]dpv1alpha1.Scope)
	for _, rule := range httpRoute.Spec.Rules {
		for _, filter := range rule.Filters {
			if filter.Type == gwapiv1b1.HTTPRouteFilterExtensionRef && filter.ExtensionRef != nil && filter.ExtensionRef.Kind == constants.KindScope {
				scope := &dpv1alpha1.Scope{}
				if err := apiReconciler.client.Get(ctx, types.NamespacedName{Namespace: httpRoute.Namespace, Name: string(filter.ExtensionRef.Name)},
					scope); err != nil {
					return nil, fmt.Errorf("error while getting scope %s in namespace :%s, %s", filter.ExtensionRef.Name,
						httpRoute.Namespace, err.Error())
				}
				scopes[utils.NamespacedName(scope).String()] = *scope
				setIndexScopeForAPI(ctx, apiReconciler.mgr, scope, api)
			}
		}
	}

	return scopes, nil
}

func (apiReconciler *APIReconciler) getAuthenticationsForResources(ctx context.Context,
	httpRoute *gwapiv1b1.HTTPRoute) (map[string]dpv1alpha1.Authentication, error) {
	authentications := make(map[string]dpv1alpha1.Authentication)
	authenticationList := &dpv1alpha1.AuthenticationList{}
	if err := apiReconciler.client.List(ctx, authenticationList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(httpRouteAuthenticationResourceIndex, utils.NamespacedName(httpRoute).String()),
	}); err != nil {
		return nil, err
	}
	for _, item := range authenticationList.Items {
		authentications[utils.NamespacedName(&item).String()] = item
	}
	return authentications, nil
}

func (apiReconciler *APIReconciler) getRatelimitPoliciesForResources(ctx context.Context,
	httpRoute *gwapiv1b1.HTTPRoute) (map[string]dpv1alpha1.RateLimitPolicy, error) {
	ratelimitpolicies := make(map[string]dpv1alpha1.RateLimitPolicy)
	ratelimitPolicyList := &dpv1alpha1.RateLimitPolicyList{}
	if err := apiReconciler.client.List(ctx, ratelimitPolicyList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(httpRouteRateLimitResourceIndex, utils.NamespacedName(httpRoute).String()),
	}); err != nil {
		return nil, err
	}
	for _, item := range ratelimitPolicyList.Items {
		ratelimitpolicies[utils.NamespacedName(&item).String()] = item
	}
	return ratelimitpolicies, nil
}

func (apiReconciler *APIReconciler) getAPIPoliciesForHTTPRoute(ctx context.Context,
	httpRoute *gwapiv1b1.HTTPRoute) (map[string]dpv1alpha1.APIPolicy, error) {
	apiPolicies := make(map[string]dpv1alpha1.APIPolicy)
	apiPolicyList := &dpv1alpha1.APIPolicyList{}
	if err := apiReconciler.client.List(ctx, apiPolicyList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(httpRouteAPIPolicyIndex, utils.NamespacedName(httpRoute).String()),
	}); err != nil {
		return nil, err
	}
	for _, item := range apiPolicyList.Items {
		apiPolicies[utils.NamespacedName(&item).String()] = item
	}
	return apiPolicies, nil
}

func (apiReconciler *APIReconciler) getAPIPoliciesForResources(ctx context.Context,
	httpRoute *gwapiv1b1.HTTPRoute) (map[string]dpv1alpha1.APIPolicy, error) {
	apiPolicies := make(map[string]dpv1alpha1.APIPolicy)
	apiPolicyList := &dpv1alpha1.APIPolicyList{}
	if err := apiReconciler.client.List(ctx, apiPolicyList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(httpRouteAPIPolicyResourceIndex, utils.NamespacedName(httpRoute).String()),
	}); err != nil {
		return nil, err
	}
	for _, item := range apiPolicyList.Items {
		apiPolicies[utils.NamespacedName(&item).String()] = item
	}
	return apiPolicies, nil
}

func (apiReconciler *APIReconciler) getResolvedBackendsMapping(ctx context.Context,
	httpRoute *gwapiv1b1.HTTPRoute) dpv1alpha1.BackendMapping {
	backendMapping := make(dpv1alpha1.BackendMapping)
	for _, rule := range httpRoute.Spec.Rules {
		for _, backend := range rule.BackendRefs {
			backendNamespacedName := types.NamespacedName{
				Name:      string(backend.Name),
				Namespace: utils.GetNamespace(backend.Namespace, httpRoute.Namespace),
			}
			resolvedBackend := apiReconciler.getResolvedBackend(ctx, backendNamespacedName)
			if resolvedBackend != nil {
				backendMapping[backendNamespacedName] = resolvedBackend
			}
		}
	}
	loggers.LoggerAPKOperator.Debugf("Generated backendMapping: %v", backendMapping)
	return backendMapping
}

// getHostNameForService resolves the backed hostname for services.
// When service type is ExternalName then ExternalName property is used as the hostname.
// Otherwise defaulted to service name as <namespace>.<service>
func (apiReconciler *APIReconciler) getHostNameForBackend(ctx context.Context, backend gwapiv1b1.HTTPBackendRef,
	defaultNamespace string) string {
	var service = new(corev1.Service)
	err := apiReconciler.client.Get(context.Background(), types.NamespacedName{
		Name:      string(backend.Name),
		Namespace: utils.GetNamespace(backend.Namespace, defaultNamespace)}, service)
	if err == nil {
		switch service.Spec.Type {
		case corev1.ServiceTypeExternalName:
			return service.Spec.ExternalName
		}
	}
	return utils.GetDefaultHostNameForBackend(backend, defaultNamespace)
}

// getTLSConfigForBackend resolves backend TLS configurations.
func (apiReconciler *APIReconciler) getResolvedBackend(ctx context.Context,
	backendNamespacedName types.NamespacedName) *dpv1alpha1.ResolvedBackend {
	resolvedBackend := dpv1alpha1.ResolvedBackend{}
	resolvedTLSConfig := dpv1alpha1.ResolvedTLSConfig{}

	var backend = new(dpv1alpha1.Backend)
	err := apiReconciler.client.Get(context.Background(), backendNamespacedName, backend)

	if err != nil {
		if !apierrors.IsNotFound(err) {
			loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2637, backendNamespacedName, err.Error()))
		}
		return nil
	}
	resolvedBackend.Services = backend.Spec.Services
	resolvedBackend.Protocol = backend.Spec.Protocol
	if backend.Spec.TLS != nil {
		resolvedTLSConfig.ResolvedCertificate = resolveCertificate(ctx, apiReconciler.client,
			backend.Namespace, *backend.Spec.TLS)
		resolvedTLSConfig.AllowedSANs = backend.Spec.TLS.AllowedSANs
		resolvedBackend.TLS = resolvedTLSConfig
	}
	if backend.Spec.Security != nil {
		resolvedBackend.Security = getResolvedBackendSecurity(ctx, apiReconciler.client,
			backend.Namespace, backend.Spec.Security)
	}
	return &resolvedBackend
}

// getResolvedBackendSecurity resolves backend security configurations.
func getResolvedBackendSecurity(ctx context.Context, client k8client.Client,
	namespace string, security []dpv1alpha1.SecurityConfig) []dpv1alpha1.ResolvedSecurityConfig {
	resolvedSecurity := make([]dpv1alpha1.ResolvedSecurityConfig, len(security))
	for _, sec := range security {
		switch sec.Type {
		case "Basic":
			var err error
			var username string
			var password string
			username, err = utils.GetSecretValue(ctx, client,
				namespace, sec.Basic.SecretRef.Name, sec.Basic.SecretRef.UsernameKey)
			password, err = utils.GetSecretValue(ctx, client,
				namespace, sec.Basic.SecretRef.Name, sec.Basic.SecretRef.PasswordKey)
			if err != nil {
				loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2648, sec.Basic.SecretRef))
			}
			resolvedSecurity = append(resolvedSecurity, dpv1alpha1.ResolvedSecurityConfig{
				Type: "Basic",
				Basic: dpv1alpha1.ResolvedBasicSecurityConfig{
					Username: username,
					Password: password,
				},
			})
		}
	}
	return resolvedSecurity
}

// resolveCertificate reads the certificate from TLSConfig, first checks the certificateInline field,
// if no value then load the certificate from secretRef using util function called GetSecretValue
func resolveCertificate(ctx context.Context, client k8client.Client, namespace string, tlsConfig dpv1alpha1.TLSConfig) string {
	var certificate string
	var err error
	if len(tlsConfig.CertificateInline) > 0 {
		certificate = tlsConfig.CertificateInline
	} else if tlsConfig.SecretRef != nil {
		if certificate, err = utils.GetSecretValue(ctx, client,
			namespace, tlsConfig.SecretRef.Name, tlsConfig.SecretRef.Key); err != nil {
			loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2642, tlsConfig.SecretRef))
		}
	} else if tlsConfig.ConfigMapRef != nil {
		if certificate, err = utils.GetConfigMapValue(ctx, client,
			namespace, tlsConfig.ConfigMapRef.Name, tlsConfig.ConfigMapRef.Key); err != nil {
			loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2643, tlsConfig.ConfigMapRef))
		}
	}
	if len(certificate) > 0 {
		block, _ := pem.Decode([]byte(certificate))
		if block == nil {
			loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2627))
			return ""
		}
		_, err = x509.ParseCertificate(block.Bytes)
		if block == nil {
			loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2641, err.Error()))
			return ""
		}
	}
	return certificate
}

// getAPIForHTTPRoute triggers the API controller reconcile method based on the changes detected
// from HTTPRoute objects. If the changes are done for an API stored in the Operator Data store,
// a new reconcile event will be created and added to the reconcile event queue.
func (apiReconciler *APIReconciler) getAPIForHTTPRoute(obj k8client.Object) []reconcile.Request {
	ctx := context.Background()
	httpRoute, ok := obj.(*gwapiv1b1.HTTPRoute)
	if !ok {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2622, httpRoute))
		return []reconcile.Request{}
	}

	apiList := &dpv1alpha1.APIList{}
	if err := apiReconciler.client.List(ctx, apiList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(httpRouteAPIIndex, utils.NamespacedName(httpRoute).String()),
	}); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2623, utils.NamespacedName(httpRoute).String()))
		return []reconcile.Request{}
	}

	if len(apiList.Items) == 0 {
		loggers.LoggerAPKOperator.Debugf("APIs for HTTPRoute not found: %s", utils.NamespacedName(httpRoute).String())
		return []reconcile.Request{}
	}

	requests := []reconcile.Request{}
	for _, api := range apiList.Items {
		req := reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      api.Name,
				Namespace: api.Namespace},
		}
		requests = append(requests, req)
		loggers.LoggerAPKOperator.Infof("Adding reconcile request for API: %s/%s", api.Namespace, api.Name)
	}
	return requests
}

// getAPIsForConfigMap triggers the API controller reconcile method based on the changes detected
// in configMap resources.
func (apiReconciler *APIReconciler) getAPIsForConfigMap(obj k8client.Object) []reconcile.Request {
	ctx := context.Background()
	configMap, ok := obj.(*corev1.ConfigMap)
	if !ok {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2622, configMap))
		return []reconcile.Request{}
	}

	backendList := &dpv1alpha1.BackendList{}
	if err := apiReconciler.client.List(ctx, backendList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(configMapBackend, utils.NamespacedName(configMap).String()),
	}); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2638, utils.NamespacedName(configMap).String()))
		return []reconcile.Request{}
	}

	requests := []reconcile.Request{}
	for _, backend := range backendList.Items {
		requests = append(requests, apiReconciler.getAPIsForBackend(&backend)...)
	}
	return requests
}

// getAPIsForSecret triggers the API controller reconcile method based on the changes detected
// in secret resources.
func (apiReconciler *APIReconciler) getAPIsForSecret(obj k8client.Object) []reconcile.Request {
	ctx := context.Background()
	secret, ok := obj.(*corev1.Secret)
	if !ok {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2622, secret))
		return []reconcile.Request{}
	}

	backendList := &dpv1alpha1.BackendList{}
	if err := apiReconciler.client.List(ctx, backendList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(secretBackend, utils.NamespacedName(secret).String()),
	}); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2621, utils.NamespacedName(secret).String()))
		return []reconcile.Request{}
	}

	requests := []reconcile.Request{}
	for _, backend := range backendList.Items {
		requests = append(requests, apiReconciler.getAPIsForBackend(&backend)...)
	}
	return requests
}

// getAPIForAuthentication triggers the API controller reconcile method based on the changes detected
// from Authentication objects. If the changes are done for an API stored in the Operator Data store,
// a new reconcile event will be created and added to the reconcile event queue.
func (apiReconciler *APIReconciler) getAPIsForAuthentication(obj k8client.Object) []reconcile.Request {
	authentication, ok := obj.(*dpv1alpha1.Authentication)
	ctx := context.Background()
	if !ok {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2624, authentication))
		return []reconcile.Request{}
	}

	requests := []reconcile.Request{}

	// todo(amali) move this validation to validation hook
	if !(authentication.Spec.TargetRef.Kind == constants.KindHTTPRoute || authentication.Spec.TargetRef.Kind == constants.KindResource) {
		loggers.LoggerAPKOperator.Errorf("Unsupported target ref kind : %s was given for authentication: %s",
			authentication.Spec.TargetRef.Kind, authentication.Name)
		return requests
	}
	loggers.LoggerAPKOperator.Debugf("Finding reconcile API requests for httpRoute: %s in namespace : %s",
		authentication.Spec.TargetRef.Name, authentication.Namespace)

	apiList := &dpv1alpha1.APIList{}

	namespacedName := types.NamespacedName{
		Name: string(authentication.Spec.TargetRef.Name),
		Namespace: utils.GetNamespace(
			(*gwapiv1b1.Namespace)(authentication.Spec.TargetRef.Namespace),
			authentication.Namespace)}.String()

	if err := apiReconciler.client.List(ctx, apiList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(httpRouteAPIIndex, namespacedName)}); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2623, namespacedName))
		return []reconcile.Request{}
	}

	if len(apiList.Items) == 0 {
		loggers.LoggerAPKOperator.Debugf("APIs for HTTPRoute not found: %s", namespacedName)
		return []reconcile.Request{}
	}

	for _, api := range apiList.Items {
		req := reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      api.Name,
				Namespace: api.Namespace},
		}
		requests = append(requests, req)
		loggers.LoggerAPKOperator.Infof("Adding reconcile request for API: %s/%s", api.Namespace, api.Name)
	}
	return requests
}

// getAPIForAuthentication triggers the API controller reconcile method based on the changes detected
// from Authentication objects. If the changes are done for an API stored in the Operator Data store,
// a new reconcile event will be created and added to the reconcile event queue.
func (apiReconciler *APIReconciler) getAPIsForAPIPolicy(obj k8client.Object) []reconcile.Request {
	apiPolicy, ok := obj.(*dpv1alpha1.APIPolicy)
	ctx := context.Background()
	if !ok {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2624, apiPolicy))
		return []reconcile.Request{}
	}

	requests := []reconcile.Request{}

	// todo(amali) move this validation to validation hook
	if !(apiPolicy.Spec.TargetRef.Kind == constants.KindHTTPRoute || apiPolicy.Spec.TargetRef.Kind == constants.KindResource) {
		loggers.LoggerAPKOperator.Errorf("Unsupported target ref kind : %s was given for authentication: %s",
			apiPolicy.Spec.TargetRef.Kind, apiPolicy.Name)
		return requests
	}
	loggers.LoggerAPKOperator.Debugf("Finding reconcile API requests for httpRoute: %s in namespace : %s",
		apiPolicy.Spec.TargetRef.Name, apiPolicy.Namespace)

	apiList := &dpv1alpha1.APIList{}

	namespacedName := types.NamespacedName{
		Name: string(apiPolicy.Spec.TargetRef.Name),
		Namespace: utils.GetNamespace(
			(*gwapiv1b1.Namespace)(apiPolicy.Spec.TargetRef.Namespace),
			apiPolicy.Namespace)}.String()

	if err := apiReconciler.client.List(ctx, apiList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(httpRouteAPIIndex, namespacedName)}); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2623, namespacedName))
		return []reconcile.Request{}
	}

	if len(apiList.Items) == 0 {
		loggers.LoggerAPKOperator.Debugf("APIs for HTTPRoute not found: %s", namespacedName)
		return []reconcile.Request{}
	}

	for _, api := range apiList.Items {
		req := reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      api.Name,
				Namespace: api.Namespace},
		}
		requests = append(requests, req)
		loggers.LoggerAPKOperator.Infof("Adding reconcile request for API: %s/%s", api.Namespace, api.Name)
	}
	return requests
}

// getAPIsForScope triggers the API controller reconcile method based on the changes detected
// from scope objects. If the changes are done for an API stored in the Operator Data store,
// a new reconcile event will be created and added to the reconcile event queue.
func (apiReconciler *APIReconciler) getAPIsForScope(obj k8client.Object) []reconcile.Request {
	scope, ok := obj.(*dpv1alpha1.Scope)
	ctx := context.Background()
	if !ok {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2624, scope))
		return []reconcile.Request{}
	}

	requests := []reconcile.Request{}

	apiList := &dpv1alpha1.APIList{}

	if err := apiReconciler.client.List(ctx, apiList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(httpRouteAPIIndex, utils.NamespacedName(scope).String())}); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2623, utils.NamespacedName(scope).String()))
		return []reconcile.Request{}
	}

	if len(apiList.Items) == 0 {
		loggers.LoggerAPKOperator.Debugf("APIs for HTTPRoute not found: %s", utils.NamespacedName(scope).String())
		return []reconcile.Request{}
	}

	for _, api := range apiList.Items {
		req := reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      api.Name,
				Namespace: api.Namespace},
		}
		requests = append(requests, req)
		loggers.LoggerAPKOperator.Infof("Adding reconcile request for API: %s/%s", api.Namespace, api.Name)
	}
	return requests
}

// getAPIsForBackend triggers the API controller reconcile method based on the changes detected
// in backend resources.
func (apiReconciler *APIReconciler) getAPIsForBackend(obj k8client.Object) []reconcile.Request {
	ctx := context.Background()
	backend, ok := obj.(*dpv1alpha1.Backend)
	if !ok {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2624, backend))
		return []reconcile.Request{}
	}

	httpRouteList := &gwapiv1b1.HTTPRouteList{}
	if err := apiReconciler.client.List(ctx, httpRouteList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(backendHTTPRouteIndex, utils.NamespacedName(backend).String()),
	}); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2625, utils.NamespacedName(backend).String()))
		return []reconcile.Request{}
	}

	if len(httpRouteList.Items) == 0 {
		loggers.LoggerAPKOperator.Debugf("HTTPRoutes for Backend not found: %s", utils.NamespacedName(backend).String())
		return []reconcile.Request{}
	}

	requests := []reconcile.Request{}
	for _, httpRoute := range httpRouteList.Items {
		requests = append(requests, apiReconciler.getAPIForHTTPRoute(&httpRoute)...)
	}
	return requests
}

// addIndexes adds indexing on API, for
//   - prodution and sandbox HTTPRoutes
//     referenced in API objects via `.spec.prodHTTPRouteRef` and `.spec.sandHTTPRouteRef`
//     This helps to find APIs that are affected by a HTTPRoute CRUD operation.
//   - authentications
//     authentication schemes related to httproutes
//     This helps to find authentication schemes binded to HTTPRoute.
//   - apiPolicies
//     apiPolicy schemes related to httproutes
//     This helps to find apiPolicy schemes binded to HTTPRoute.
func addIndexes(ctx context.Context, mgr manager.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha1.API{}, httpRouteAPIIndex,
		func(rawObj k8client.Object) []string {
			api := rawObj.(*dpv1alpha1.API)
			var httpRoutes []string
			for _, ref := range api.Spec.ProdHTTPRouteRefs {
				if ref != "" {
					httpRoutes = append(httpRoutes,
						types.NamespacedName{
							Namespace: api.Namespace,
							Name:      ref,
						}.String())
				}
			}
			for _, ref := range api.Spec.SandHTTPRouteRefs {
				if ref != "" {
					httpRoutes = append(httpRoutes,
						types.NamespacedName{
							Namespace: api.Namespace,
							Name:      ref,
						}.String())
				}
			}
			return httpRoutes
		}); err != nil {
		return err
	}

	// Backend to HTTPRoute indexer
	if err := mgr.GetFieldIndexer().IndexField(ctx, &gwapiv1b1.HTTPRoute{}, backendHTTPRouteIndex,
		func(rawObj k8client.Object) []string {
			httpRoute := rawObj.(*gwapiv1b1.HTTPRoute)
			var backends []string
			for _, rule := range httpRoute.Spec.Rules {
				for _, backendRef := range rule.BackendRefs {
					if backendRef.Kind != nil && *backendRef.Kind == constants.KindBackend {
						backends = append(backends, types.NamespacedName{
							Namespace: utils.GetNamespace(backendRef.Namespace,
								httpRoute.ObjectMeta.Namespace),
							Name: string(backendRef.Name),
						}.String())
					}
				}
			}
			return backends
		}); err != nil {
		return err
	}

	// ConfigMap to Backend indexer
	if err := mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha1.Backend{}, configMapBackend,
		func(rawObj k8client.Object) []string {
			backend := rawObj.(*dpv1alpha1.Backend)
			var configMaps []string
			if backend.Spec.TLS != nil && backend.Spec.TLS.ConfigMapRef != nil && len(backend.Spec.TLS.ConfigMapRef.Name) > 0 {
				configMaps = append(configMaps,
					types.NamespacedName{
						Name:      string(backend.Spec.TLS.ConfigMapRef.Name),
						Namespace: backend.Namespace,
					}.String())
			}
			return configMaps
		}); err != nil {
		return err
	}

	// Secret to Backend indexer
	if err := mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha1.Backend{}, secretBackend,
		func(rawObj k8client.Object) []string {
			backend := rawObj.(*dpv1alpha1.Backend)
			var secrets []string
			if backend.Spec.TLS != nil && backend.Spec.TLS.SecretRef != nil && len(backend.Spec.TLS.SecretRef.Name) > 0 {
				secrets = append(secrets,
					types.NamespacedName{
						Name:      string(backend.Spec.TLS.SecretRef.Name),
						Namespace: backend.Namespace,
					}.String())
			}
			if backend.Spec.Security != nil {
				for _, security := range backend.Spec.Security {
					if security.Type == "Basic" {
						secrets = append(secrets,
							types.NamespacedName{
								Name:      string(security.Basic.SecretRef.Name),
								Namespace: backend.Namespace,
							}.String())
					}
				}
			}
			return secrets
		}); err != nil {
		return err
	}

	// authentication to HTTPRoute indexer
	if err := mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha1.Authentication{}, httpRouteAuthenticationIndex,
		func(rawObj k8client.Object) []string {
			authentication := rawObj.(*dpv1alpha1.Authentication)
			var httpRoutes []string
			if authentication.Spec.TargetRef.Kind == constants.KindHTTPRoute {
				httpRoutes = append(httpRoutes,
					types.NamespacedName{
						Namespace: utils.GetNamespace(
							(*gwapiv1b1.Namespace)(authentication.Spec.TargetRef.Namespace),
							authentication.Namespace),
						Name: string(authentication.Spec.TargetRef.Name),
					}.String())
			}
			return httpRoutes
		}); err != nil {
		return err
	}

	// Till the below is httproute rule name and targetref sectionname is supported,
	// https://gateway-api.sigs.k8s.io/geps/gep-713/?h=multiple+targetrefs#apply-policies-to-sections-of-a-resource-future-extension
	// we will use a temporary kindName called Resource for policy attachments
	// TODO(amali) Fix after the official support is available
	if err := mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha1.Authentication{}, httpRouteAuthenticationResourceIndex,
		func(rawObj k8client.Object) []string {
			authentication := rawObj.(*dpv1alpha1.Authentication)
			var httpRoutes []string
			if authentication.Spec.TargetRef.Kind == constants.KindResource {
				httpRoutes = append(httpRoutes,
					types.NamespacedName{
						Namespace: utils.GetNamespace(
							(*gwapiv1b1.Namespace)(authentication.Spec.TargetRef.Namespace),
							authentication.Namespace),
						Name: string(authentication.Spec.TargetRef.Name),
					}.String())
			}
			return httpRoutes
		}); err != nil {
		return err
	}

	// ratelimite policy to HTTPRoute indexer
	if err := mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha1.RateLimitPolicy{}, httpRouteRateLimitIndex,
		func(rawObj k8client.Object) []string {
			ratelimitPolicy := rawObj.(*dpv1alpha1.RateLimitPolicy)
			var httpRoutes []string
			if ratelimitPolicy.Spec.TargetRef.Kind == constants.KindHTTPRoute {
				httpRoutes = append(httpRoutes,
					types.NamespacedName{
						Namespace: utils.GetNamespace(
							(*gwapiv1b1.Namespace)(ratelimitPolicy.Spec.TargetRef.Namespace),
							ratelimitPolicy.Namespace),
						Name: string(ratelimitPolicy.Spec.TargetRef.Name),
					}.String())
			}
			return httpRoutes
		}); err != nil {
		return err
	}

	// Till the below is httproute rule name and targetref sectionname is supported,
	// https://gateway-api.sigs.k8s.io/geps/gep-713/?h=multiple+targetrefs#apply-policies-to-sections-of-a-resource-future-extension
	// we will use a temporary kindName called Resource for policy attachments
	// TODO(amali) Fix after the official support is available
	if err := mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha1.RateLimitPolicy{}, httpRouteRateLimitResourceIndex,
		func(rawObj k8client.Object) []string {
			ratelimitPolicy := rawObj.(*dpv1alpha1.RateLimitPolicy)
			var httpRoutes []string
			if ratelimitPolicy.Spec.TargetRef.Kind == constants.KindResource {
				httpRoutes = append(httpRoutes,
					types.NamespacedName{
						Namespace: utils.GetNamespace(
							(*gwapiv1b1.Namespace)(ratelimitPolicy.Spec.TargetRef.Namespace),
							ratelimitPolicy.Namespace),
						Name: string(ratelimitPolicy.Spec.TargetRef.Name),
					}.String())
			}
			return httpRoutes
		}); err != nil {
		return err
	}

	// APIPolicy to HTTPRoute indexer
	if err := mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha1.APIPolicy{}, httpRouteAPIPolicyIndex,
		func(rawObj k8client.Object) []string {
			apiPolicy := rawObj.(*dpv1alpha1.APIPolicy)
			var httpRoutes []string
			if apiPolicy.Spec.TargetRef.Kind == constants.KindHTTPRoute {
				httpRoutes = append(httpRoutes,
					types.NamespacedName{
						Namespace: utils.GetNamespace(
							(*gwapiv1b1.Namespace)(apiPolicy.Spec.TargetRef.Namespace), apiPolicy.Namespace),
						Name: string(apiPolicy.Spec.TargetRef.Name),
					}.String())
			}
			return httpRoutes
		}); err != nil {
		return err
	}

	// Till the below is httproute rule name and targetref sectionname is supported,
	// https://gateway-api.sigs.k8s.io/geps/gep-713/?h=multiple+targetrefs#apply-policies-to-sections-of-a-resource-future-extension
	// we will use a temporary kindName called Resource for policy attachments
	// TODO(amali) Fix after the official support is available
	err := mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha1.APIPolicy{}, httpRouteAPIPolicyResourceIndex,
		func(rawObj k8client.Object) []string {
			apiPolicy := rawObj.(*dpv1alpha1.APIPolicy)
			var httpRoutes []string
			if apiPolicy.Spec.TargetRef.Kind == constants.KindResource {
				httpRoutes = append(httpRoutes,
					types.NamespacedName{
						Namespace: utils.GetNamespace(
							(*gwapiv1b1.Namespace)(apiPolicy.Spec.TargetRef.Namespace),
							apiPolicy.Namespace),
						Name: string(apiPolicy.Spec.TargetRef.Name),
					}.String())
			}
			return httpRoutes
		})
	return err
}

// setIndexScopeForAPI sets the index for the Scope CR to API CR
func setIndexScopeForAPI(ctx context.Context, mgr manager.Manager, scope *dpv1alpha1.Scope, api string) error {
	return mgr.GetFieldIndexer().IndexField(ctx, scope, apiScopeIndex,
		func(rawObj k8client.Object) []string {
			return []string{api}
		})
}

// handleStatus updates the API CR update
func (apiReconciler *APIReconciler) handleStatus(apiKey types.NamespacedName, state string, events []string) {
	accept := false
	message := ""
	event := ""

	switch state {
	case constants.DeployedState:
		accept = true
		message = "API is deployed to the gateway."
	case constants.UpdatedState:
		accept = true
		message = fmt.Sprintf("API update is deployed to the gateway. %v Updated", events)
	}
	timeNow := metav1.Now()
	event = fmt.Sprintf("[%s] %s", timeNow.String(), message)

	apiReconciler.statusUpdater.Send(status.Update{
		NamespacedName: apiKey,
		Resource:       new(dpv1alpha1.API),
		UpdateStatus: func(obj k8client.Object) k8client.Object {
			h, ok := obj.(*dpv1alpha1.API)
			if !ok {
				loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2626, obj))
			}
			hCopy := h.DeepCopy()
			hCopy.Status.Status = state
			hCopy.Status.Accepted = accept
			hCopy.Status.Message = message
			hCopy.Status.Events = append(hCopy.Status.Events, event)
			hCopy.Status.TransitionTime = &timeNow
			return hCopy
		},
	})
}
