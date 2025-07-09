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
	"fmt"
	"sync"

	"github.com/wso2/apk/adapter/config"
	"github.com/wso2/apk/adapter/internal/discovery/xds"
	"github.com/wso2/apk/adapter/internal/loggers"
	"github.com/wso2/apk/adapter/pkg/logging"
	"golang.org/x/exp/maps"
	corev1 "k8s.io/api/core/v1"
	k8error "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/types"

	"github.com/wso2/apk/adapter/internal/clients/kvresolver"
	"github.com/wso2/apk/adapter/internal/operator/constants"
	"github.com/wso2/apk/adapter/internal/operator/status"
	"github.com/wso2/apk/adapter/internal/operator/synchronizer"
	"github.com/wso2/apk/adapter/internal/operator/utils"
	dpv1alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha1"
	dpv1alpha3 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha3"
	dpv1alpha5 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha5"
	ctrl "sigs.k8s.io/controller-runtime"
	k8client "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
)

const (
	gatewayRateLimitPolicyIndex = "gatewayRateLimitPolicyIndex"
	gatewayAPIPolicyIndex       = "gatewayAPIPolicyIndex"
	tokenIssuerIndex            = "tokenIssuerIndex"
	secretTokenIssuerIndex      = "secretTokenIssuerIndex"
	configmapIssuerIndex        = "configmapIssuerIndex"
	defaultAllEnvironments      = "*"
)

var (
	setReadiness sync.Once
)

// GatewayReconciler reconciles a Gateway object
type GatewayReconciler struct {
	client        k8client.Client
	kvClient      *kvresolver.KVResolverClientImpl
	ods           *synchronizer.OperatorDataStore
	ch            *chan synchronizer.GatewayEvent
	statusUpdater *status.UpdateHandler
	mgr           manager.Manager
}

// NewGatewayController creates a new Gateway controller instance. Gateway Controllers watches for gwapiv1.Gateway.
func NewGatewayController(mgr manager.Manager, operatorDataStore *synchronizer.OperatorDataStore, statusUpdater *status.UpdateHandler,
	ch *chan synchronizer.GatewayEvent) error {
	r := &GatewayReconciler{
		client:        mgr.GetClient(),
		kvClient:      kvresolver.InitializeKVResolverClient(),
		ods:           operatorDataStore,
		ch:            ch,
		statusUpdater: statusUpdater,
		mgr:           mgr,
	}
	c, err := controller.New(constants.GatewayController, mgr, controller.Options{Reconciler: r})
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error3119, logging.BLOCKER, "Error creating API controller, error: %v", err))
		return err
	}

	ctx := context.Background()
	conf := config.ReadConfigs()

	if err := addGatewayIndexes(ctx, mgr); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error3120, logging.BLOCKER, "Error adding indexes: %v", err))
		return err
	}

	predicateGateway := []predicate.TypedPredicate[*gwapiv1.Gateway]{predicate.NewTypedPredicateFuncs[*gwapiv1.Gateway](utils.FilterGatewayByNamespaces(conf.Adapter.Operator.Namespaces))}
	if err := c.Watch(source.Kind(mgr.GetCache(), &gwapiv1.Gateway{}, &handler.TypedEnqueueRequestForObject[*gwapiv1.Gateway]{},
		predicateGateway...)); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2611, logging.BLOCKER, "Error watching Gateway resources: %v", err))
		return err
	}

	predicateRateLimitPolicy := []predicate.TypedPredicate[*dpv1alpha3.RateLimitPolicy]{predicate.NewTypedPredicateFuncs[*dpv1alpha3.RateLimitPolicy](utils.FilterRateLimitPolicyByNamespaces(conf.Adapter.Operator.Namespaces))}
	if err := c.Watch(source.Kind(mgr.GetCache(), &dpv1alpha3.RateLimitPolicy{}, handler.TypedEnqueueRequestsFromMapFunc(r.handleCustomRateLimitPolicies),
		predicateRateLimitPolicy...)); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2639, logging.BLOCKER, "Error watching Ratelimit resources: %v", err))
		return err
	}

	predicateAPIPolicy := []predicate.TypedPredicate[*dpv1alpha5.APIPolicy]{predicate.NewTypedPredicateFuncs[*dpv1alpha5.APIPolicy](utils.FilterAPIPolicyByNamespaces(conf.Adapter.Operator.Namespaces))}
	if err := c.Watch(source.Kind(mgr.GetCache(), &dpv1alpha5.APIPolicy{}, handler.TypedEnqueueRequestsFromMapFunc(r.getGatewaysForAPIPolicy),
		predicateAPIPolicy...)); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2617, logging.BLOCKER, "Error watching APIPolicy resources: %v", err))
		return err
	}

	predicateInterceptorService := []predicate.TypedPredicate[*dpv1alpha1.InterceptorService]{predicate.NewTypedPredicateFuncs[*dpv1alpha1.InterceptorService](utils.FilterInterceptorServiceByNamespaces(conf.Adapter.Operator.Namespaces))}
	if err := c.Watch(source.Kind(mgr.GetCache(), &dpv1alpha1.InterceptorService{}, handler.TypedEnqueueRequestsFromMapFunc(r.getAPIsForInterceptorService),
		predicateInterceptorService...)); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2640, logging.BLOCKER, "Error watching InterceptorService resources: %v", err))
		return err
	}

	predicateBackendJWT := []predicate.TypedPredicate[*dpv1alpha1.BackendJWT]{predicate.NewTypedPredicateFuncs[*dpv1alpha1.BackendJWT](utils.FilterBackendJWTByNamespaces(conf.Adapter.Operator.Namespaces))}
	if err := c.Watch(source.Kind(mgr.GetCache(), &dpv1alpha1.BackendJWT{}, handler.TypedEnqueueRequestsFromMapFunc(r.getAPIsForBackendJWT),
		predicateBackendJWT...)); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2661, logging.BLOCKER, "Error watching BackendJWT resources: %v", err))
		return err
	}

	predicateBackend := []predicate.TypedPredicate[*dpv1alpha5.Backend]{predicate.NewTypedPredicateFuncs[*dpv1alpha5.Backend](utils.FilterBackendByNamespaces(conf.Adapter.Operator.Namespaces))}
	if err := c.Watch(source.Kind(mgr.GetCache(), &dpv1alpha5.Backend{}, handler.TypedEnqueueRequestsFromMapFunc(r.getGatewaysForBackend),
		predicateBackend...)); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2615, logging.BLOCKER, "Error watching Backend resources: %v", err))
		return err
	}

	predicateConfigMap := []predicate.TypedPredicate[*corev1.ConfigMap]{predicate.NewTypedPredicateFuncs[*corev1.ConfigMap](utils.FilterConfigMapByNamespaces(conf.Adapter.Operator.Namespaces))}
	if err := c.Watch(source.Kind(mgr.GetCache(), &corev1.ConfigMap{}, handler.TypedEnqueueRequestsFromMapFunc(r.getGatewaysForConfigMap),
		predicateConfigMap...)); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2644, logging.BLOCKER, "Error watching ConfigMap resources: %v", err))
		return err
	}

	predicateSecret := []predicate.TypedPredicate[*corev1.Secret]{predicate.NewTypedPredicateFuncs[*corev1.Secret](utils.FilterSecretByNamespaces(conf.Adapter.Operator.Namespaces))}
	if err := c.Watch(source.Kind(mgr.GetCache(), &corev1.Secret{}, handler.TypedEnqueueRequestsFromMapFunc(r.getGatewaysForSecret),
		predicateSecret...)); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2645, logging.BLOCKER, "Error watching Secret resources: %v", err))
		return err
	}
	predicateTokenIssuer := []predicate.TypedPredicate[*dpv1alpha1.TokenIssuer]{predicate.NewTypedPredicateFuncs[*dpv1alpha1.TokenIssuer](utils.FilterTokenIssuerByNamespaces(conf.Adapter.Operator.Namespaces))}
	if err := c.Watch(source.Kind(mgr.GetCache(), &dpv1alpha1.TokenIssuer{}, handler.TypedEnqueueRequestsFromMapFunc(r.getGatewaysForTokenIssuer),
		predicateTokenIssuer...)); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2646, logging.BLOCKER, "Error watching TokenIssuer resources: %v", err))
		return err
	}

	loggers.LoggerAPKOperator.Info("Gateway Controller successfully started. Watching Gateway Objects....")
	return nil
}

//+kubebuilder:rbac:groups=dp.wso2.com,resources=gateways,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=dp.wso2.com,resources=gateways/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=dp.wso2.com,resources=gateways/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Gateway object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (gatewayReconciler *GatewayReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// Check whether the Gateway CR exist, if not consider as a DELETE event.
	loggers.LoggerAPKOperator.Infof("Reconciling gateway...")
	var gatewayDef gwapiv1.Gateway
	if err := gatewayReconciler.client.Get(ctx, req.NamespacedName, &gatewayDef); err != nil {
		gatewayState, found := gatewayReconciler.ods.GetCachedGateway(req.NamespacedName)
		if found && k8error.IsNotFound(err) {
			// The gateway doesn't exist in the gateway Cache, remove it
			gatewayReconciler.ods.DeleteCachedGateway(req.NamespacedName)
			loggers.LoggerAPKOperator.Infof("Delete event has received for Gateway : %s, hence deleted from Gateway cache", req.NamespacedName.String())
			*gatewayReconciler.ch <- synchronizer.GatewayEvent{EventType: constants.Delete, Event: gatewayState}
			return ctrl.Result{}, nil
		}
		loggers.LoggerAPKOperator.Warnf("Gateway CR related to the reconcile request with key: returned error. Assuming Gateway is already deleted, hence ignoring the error : %v", err)
		return ctrl.Result{}, nil
	}
	var gwCondition []metav1.Condition = gatewayDef.Status.Conditions

	gatewayStateData, err := gatewayReconciler.resolveGatewayState(ctx, gatewayDef)
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error3122, logging.BLOCKER, "Error resolving Gateway State %s: %v", req.NamespacedName.String(), err))
		return ctrl.Result{}, err
	}

	if gwCondition[0].Type != "Accepted" {
		gatewayState := gatewayReconciler.ods.AddGatewayState(gatewayDef, gatewayStateData)
		*gatewayReconciler.ch <- synchronizer.GatewayEvent{EventType: constants.Create, Event: gatewayState}
		gatewayReconciler.handleGatewayStatus(req.NamespacedName, constants.Create, []string{})
	} else if cachedGateway, events, updated :=
		gatewayReconciler.ods.UpdateGatewayState(&gatewayDef, gatewayStateData); updated {
		*gatewayReconciler.ch <- synchronizer.GatewayEvent{EventType: constants.Update, Event: cachedGateway}
		gatewayReconciler.handleGatewayStatus(req.NamespacedName, constants.Update, events)
	}
	setReadiness.Do(gatewayReconciler.setGatewayReadiness)
	return ctrl.Result{}, nil
}

// setGatewayReadiness sets the gateway readiness status
func (gatewayReconciler *GatewayReconciler) setGatewayReadiness() {
	apisList, _ := utils.RetrieveAPIList(gatewayReconciler.client)
	if len(apisList) == 0 {
		xds.SetReady()
	}
}

// resolveListenerSecretRefs resolves the certificate secret references in the related listeners in Gateway CR
func (gatewayReconciler *GatewayReconciler) resolveListenerSecretRefs(ctx context.Context, secretRef *gwapiv1.SecretObjectReference, gatewayNamespace string) (map[string][]byte, error) {
	var secret corev1.Secret
	namespace := gwapiv1.Namespace(string(*secretRef.Namespace))
	if err := gatewayReconciler.client.Get(ctx, types.NamespacedName{Name: string(secretRef.Name),
		Namespace: utils.GetNamespace(&namespace, gatewayNamespace)}, &secret); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error3123, logging.BLOCKER, "Unable to find associated secret %s in %s: %v", secretRef.Name, string(*secretRef.Namespace), err))
		return nil, err
	}
	return secret.Data, nil
}

// resolveGatewayState resolves the GatewayState struct using gwapiv1.Gateway and resource indexes
func (gatewayReconciler *GatewayReconciler) resolveGatewayState(ctx context.Context,
	gateway gwapiv1.Gateway) (*synchronizer.GatewayStateData, error) {
	gatewayState := &synchronizer.GatewayStateData{}
	var err error
	resolvedListenerCerts := make(map[string]map[string][]byte)
	namespace := gwapiv1.Namespace(gateway.Namespace)
	// Retireve listener Certificates
	for _, listener := range gateway.Spec.Listeners {
		if listener.Protocol == gwapiv1.HTTPProtocolType {
			continue
		}
		data, err := gatewayReconciler.resolveListenerSecretRefs(ctx, &listener.TLS.CertificateRefs[0], string(namespace))
		if err != nil {
			loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error3105, logging.BLOCKER, "Error resolving listener certificates: %v", err))
			return nil, err
		}
		resolvedListenerCerts[string(listener.Name)] = data
	}
	gatewayState.GatewayResolvedListenerCerts = resolvedListenerCerts
	if gatewayState.GatewayAPIPolicies, err = gatewayReconciler.getAPIPoliciesForGateway(ctx, &gateway); err != nil {
		return nil, fmt.Errorf("error while getting gateway apipolicy for gateway: %s, %s", utils.NamespacedName(&gateway).String(), err.Error())
	}
	if gatewayState.GatewayInterceptorServiceMapping, err = gatewayReconciler.getInterceptorServicesForGateway(ctx, gatewayState.GatewayAPIPolicies); err != nil {
		return nil, fmt.Errorf("error while getting interceptor service for gateway: %s, %s", utils.NamespacedName(&gateway).String(), err.Error())
	}
	customRateLimitPolicies, err := gatewayReconciler.getCustomRateLimitPoliciesForGateway(utils.NamespacedName(&gateway))
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error3124, logging.MAJOR, "Error while getting custom rate limit policies: %s", err))
	}
	tokenIssuers, err := GetJWTIssuers(ctx, gatewayReconciler.client, gateway)
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error3126, logging.MAJOR, "Error while getting token issuers: %s", err))
	}
	gatewayState.TokenIssuers = tokenIssuers
	gatewayState.GatewayCustomRateLimitPolicies = customRateLimitPolicies
	gatewayState.GatewayBackendMapping = gatewayReconciler.getResolvedBackendsMapping(ctx, gatewayState)
	return gatewayState, nil
}

func (gatewayReconciler *GatewayReconciler) getAPIPoliciesForGateway(ctx context.Context,
	gateway *gwapiv1.Gateway) (map[string]dpv1alpha5.APIPolicy, error) {
	apiPolicies := make(map[string]dpv1alpha5.APIPolicy)
	apiPolicyList := &dpv1alpha5.APIPolicyList{}
	if err := gatewayReconciler.client.List(ctx, apiPolicyList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(gatewayAPIPolicyIndex, utils.NamespacedName(gateway).String()),
	}); err != nil {
		return nil, err
	}
	for item := range apiPolicyList.Items {
		apiPolicy := apiPolicyList.Items[item]
		apiPolicies[utils.NamespacedName(&apiPolicy).String()] = apiPolicy
	}
	return apiPolicies, nil
}

// getInterceptorServicesForGateway returns the list of interceptor services for the given gateway
func (gatewayReconciler *GatewayReconciler) getInterceptorServicesForGateway(ctx context.Context,
	gatewayAPIPolicies map[string]dpv1alpha5.APIPolicy) (map[string]dpv1alpha1.InterceptorService, error) {
	allGatewayAPIPolicies := maps.Values(gatewayAPIPolicies)
	interceptorServices := make(map[string]dpv1alpha1.InterceptorService)
	for _, apiPolicy := range allGatewayAPIPolicies {
		if apiPolicy.Spec.Default != nil && len(apiPolicy.Spec.Default.RequestInterceptors) > 0 {
			interceptorPtr := utils.GetInterceptorService(ctx, gatewayReconciler.client, apiPolicy.Namespace,
				&apiPolicy.Spec.Default.RequestInterceptors[0], nil)
			if interceptorPtr != nil {
				interceptorServices[utils.NamespacedName(interceptorPtr).String()] = *interceptorPtr
			}
		}
		if apiPolicy.Spec.Default != nil && len(apiPolicy.Spec.Default.ResponseInterceptors) > 0 {
			interceptorPtr := utils.GetInterceptorService(ctx, gatewayReconciler.client, apiPolicy.Namespace,
				&apiPolicy.Spec.Default.ResponseInterceptors[0], nil)
			if interceptorPtr != nil {
				interceptorServices[utils.NamespacedName(interceptorPtr).String()] = *interceptorPtr
			}
		}
		if apiPolicy.Spec.Override != nil && len(apiPolicy.Spec.Override.RequestInterceptors) > 0 {
			interceptorPtr := utils.GetInterceptorService(ctx, gatewayReconciler.client, apiPolicy.Namespace,
				&apiPolicy.Spec.Override.RequestInterceptors[0], nil)
			if interceptorPtr != nil {
				interceptorServices[utils.NamespacedName(interceptorPtr).String()] = *interceptorPtr
			}
		}
		if apiPolicy.Spec.Override != nil && len(apiPolicy.Spec.Override.ResponseInterceptors) > 0 {
			interceptorPtr := utils.GetInterceptorService(ctx, gatewayReconciler.client, apiPolicy.Namespace,
				&apiPolicy.Spec.Override.ResponseInterceptors[0], nil)
			if interceptorPtr != nil {
				interceptorServices[utils.NamespacedName(interceptorPtr).String()] = *interceptorPtr
			}
		}
	}
	return interceptorServices, nil
}

func (gatewayReconciler *GatewayReconciler) getResolvedBackendsMapping(ctx context.Context,
	gatewayStateData *synchronizer.GatewayStateData) map[string]*dpv1alpha5.ResolvedBackend {
	backendMapping := make(map[string]*dpv1alpha5.ResolvedBackend)
	if gatewayStateData.GatewayInterceptorServiceMapping != nil {
		interceptorServices := maps.Values(gatewayStateData.GatewayInterceptorServiceMapping)
		for _, interceptorService := range interceptorServices {
			utils.ResolveAndAddBackendToMapping(ctx, gatewayReconciler.client, gatewayReconciler.kvClient, backendMapping,
				interceptorService.Spec.BackendRef, interceptorService.Namespace, nil)
		}
	}
	return backendMapping
}

// getGatewaysForBackend triggers the Gateway controller reconcile method based on the changes detected
// in backend resources.
func (gatewayReconciler *GatewayReconciler) getGatewaysForBackend(ctx context.Context, obj *dpv1alpha5.Backend) []reconcile.Request {
	backend := obj

	requests := []reconcile.Request{}

	interceptorServiceList := &dpv1alpha1.InterceptorServiceList{}
	if err := gatewayReconciler.client.List(ctx, interceptorServiceList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(backendInterceptorServiceIndex, utils.NamespacedName(backend).String()),
	}); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error3118, logging.CRITICAL, "Unable to find associated interceptorServices: %v", utils.NamespacedName(backend).String()))
		return []reconcile.Request{}
	}

	for service := range interceptorServiceList.Items {
		interceptorService := interceptorServiceList.Items[service]
		requests = append(requests, gatewayReconciler.getAPIsForInterceptorService(ctx, &interceptorService)...)
	}

	return requests
}

// getAPIsForInterceptorService triggers the Gateway controller reconcile method based on the changes detected
// in InterceptorService resources.
func (gatewayReconciler *GatewayReconciler) getAPIsForInterceptorService(ctx context.Context, obj *dpv1alpha1.InterceptorService) []reconcile.Request {
	interceptorService := obj

	requests := []reconcile.Request{}

	apiPolicyList := &dpv1alpha5.APIPolicyList{}
	if err := gatewayReconciler.client.List(ctx, apiPolicyList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(interceptorServiceAPIPolicyIndex, utils.NamespacedName(interceptorService).String()),
	}); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error3125, logging.CRITICAL, "Unable to find associated APIPolicies: %s for Interceptor Service", utils.NamespacedName(interceptorService).String()))
		return []reconcile.Request{}
	}

	for item := range apiPolicyList.Items {
		apiPolicy := apiPolicyList.Items[item]
		requests = append(requests, gatewayReconciler.getGatewaysForAPIPolicy(ctx, &apiPolicy)...)
	}

	return requests
}

// getAPIsForBackendJWT triggers the Gateway controller reconcile method based on the changes detected
// in BackendJWT resources.
func (gatewayReconciler *GatewayReconciler) getAPIsForBackendJWT(ctx context.Context, obj *dpv1alpha1.BackendJWT) []reconcile.Request {
	backendJWT := obj
	requests := []reconcile.Request{}

	apiPolicyList := &dpv1alpha5.APIPolicyList{}
	if err := gatewayReconciler.client.List(ctx, apiPolicyList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(backendJWTAPIPolicyIndex, utils.NamespacedName(backendJWT).String()),
	}); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error3125, logging.CRITICAL, "Unable to find associated APIPolicies for BackendJWT: %s", utils.NamespacedName(backendJWT).String()))
		return []reconcile.Request{}
	}

	for item := range apiPolicyList.Items {
		apiPolicy := apiPolicyList.Items[item]
		requests = append(requests, gatewayReconciler.getGatewaysForAPIPolicy(ctx, &apiPolicy)...)
	}

	return requests
}

// getGatewaysForSecret triggers the Gateway controller reconcile method based on the changes detected
// in secret resources.
func (gatewayReconciler *GatewayReconciler) getGatewaysForSecret(ctx context.Context, obj *corev1.Secret) []reconcile.Request {
	secret := obj

	backendList := &dpv1alpha5.BackendList{}
	if err := gatewayReconciler.client.List(ctx, backendList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(secretBackend, utils.NamespacedName(secret).String()),
	}); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error3106, logging.CRITICAL, "Unable to find associated Backends for Secret: %s", utils.NamespacedName(secret).String()))
		return []reconcile.Request{}
	}

	requests := []reconcile.Request{}
	for item := range backendList.Items {
		backend := backendList.Items[item]
		requests = append(requests, gatewayReconciler.getGatewaysForBackend(ctx, &backend)...)
	}
	tokenissuerList := &dpv1alpha1.TokenIssuerList{}
	if err := gatewayReconciler.client.List(ctx, tokenissuerList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(secretTokenIssuerIndex, utils.NamespacedName(secret).String()),
	}); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error3107, logging.CRITICAL, "Unable to find associated TokenIssuers for Secret: %s", utils.NamespacedName(secret).String()))
		return []reconcile.Request{}
	}
	for item := range tokenissuerList.Items {
		tokenissuer := tokenissuerList.Items[item]
		if tokenissuer.Spec.TargetRef.Kind == constants.KindGateway {
			requests = append(requests, gatewayReconciler.getGatewaysForTokenIssuer(ctx, &tokenissuer)...)
		}
	}
	return requests
}

// getGatewaysForConfigMap triggers the API controller reconcile method based on the changes detected
// in configMap resources.
func (gatewayReconciler *GatewayReconciler) getGatewaysForConfigMap(ctx context.Context, obj *corev1.ConfigMap) []reconcile.Request {
	configMap := obj

	backendList := &dpv1alpha5.BackendList{}
	if err := gatewayReconciler.client.List(ctx, backendList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(configMapBackend, utils.NamespacedName(configMap).String()),
	}); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error3108, logging.CRITICAL, "Unable to find associated Backends for ConfigMap: %s", utils.NamespacedName(configMap).String()))
		return []reconcile.Request{}
	}

	requests := []reconcile.Request{}
	for item := range backendList.Items {
		backend := backendList.Items[item]
		requests = append(requests, gatewayReconciler.getGatewaysForBackend(ctx, &backend)...)
	}
	tokenissuerList := &dpv1alpha1.TokenIssuerList{}
	if err := gatewayReconciler.client.List(ctx, tokenissuerList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(configmapIssuerIndex, utils.NamespacedName(configMap).String()),
	}); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error3107, logging.CRITICAL, "Unable to find associated TokenIssuers for ConfigMap: %s", utils.NamespacedName(configMap).String()))
		return []reconcile.Request{}
	}
	for item := range tokenissuerList.Items {
		tokenissuer := tokenissuerList.Items[item]
		if tokenissuer.Spec.TargetRef.Kind == constants.KindGateway {
			requests = append(requests, gatewayReconciler.getGatewaysForTokenIssuer(ctx, &tokenissuer)...)
		}
	}
	return requests
}

// handleStatus updates the Gateway CR update
func (gatewayReconciler *GatewayReconciler) handleGatewayStatus(gatewayKey types.NamespacedName, state string, events []string) {
	accept := false
	message := ""
	//event := ""

	switch state {
	case constants.Create:
		accept = true
		message = "Gateway is deployed successfully"
	case constants.Update:
		accept = true
		message = fmt.Sprintf("Gateway update is deployed successfully. %v Updated", events)
	}
	timeNow := metav1.Now()
	//event = fmt.Sprintf("[%s] %s", timeNow.String(), message)

	gatewayReconciler.statusUpdater.Send(status.Update{
		NamespacedName: gatewayKey,
		Resource:       new(gwapiv1.Gateway),
		UpdateStatus: func(obj k8client.Object) k8client.Object {
			h, ok := obj.(*gwapiv1.Gateway)
			if !ok {
				loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error3109, logging.BLOCKER, "Error while updating Gateway status %v", obj))
			}
			hCopy := h.DeepCopy()
			var gwCondition []metav1.Condition = hCopy.Status.Conditions
			gwCondition[0].Status = "Unknown"
			if accept {
				gwCondition[0].Status = "True"
			} else {
				gwCondition[0].Status = "False"
			}
			gwCondition[0].Message = message
			gwCondition[0].LastTransitionTime = timeNow
			// gwCondition[0].Reason = append(gwCondition[0].Reason, event)
			gwCondition[0].Reason = "Reconciled"
			gwCondition[0].Type = state
			hCopy.Status.Conditions = gwCondition
			return hCopy
		},
	})
}

// handleCustomRateLimitPolicies returns the list of gateway reconcile requests
func (gatewayReconciler *GatewayReconciler) handleCustomRateLimitPolicies(ctx context.Context, obj *dpv1alpha3.RateLimitPolicy) []reconcile.Request {
	ratelimitPolicy := obj
	requests := []reconcile.Request{}
	if ratelimitPolicy.Spec.TargetRef.Kind == constants.KindGateway {

		namespace, err := utils.ValidateAndRetrieveNamespace((*gwapiv1.Namespace)(ratelimitPolicy.Spec.TargetRef.Namespace), ratelimitPolicy.Namespace)

		if err != nil {
			loggers.LoggerAPKOperator.Errorf("Namespace mismatch. TargetRef %s needs to be in the same namespace as the RatelimitPolicy %s. Expected: %s, Actual: %s",
				string(ratelimitPolicy.Spec.TargetRef.Name), ratelimitPolicy.Name, ratelimitPolicy.Namespace, string(*ratelimitPolicy.Spec.TargetRef.Namespace))
			return requests
		}

		requests = append(requests, reconcile.Request{
			NamespacedName: types.NamespacedName{
				Namespace: namespace,
				Name:      string(ratelimitPolicy.Spec.TargetRef.Name),
			},
		})
	}
	return requests
}

// getCustomRateLimitPoliciesForGateway returns the list of custom rate limit policies for a gateway
func (gatewayReconciler *GatewayReconciler) getCustomRateLimitPoliciesForGateway(gatewayName types.NamespacedName) (map[string]*dpv1alpha3.RateLimitPolicy, error) {
	ctx := context.Background()
	var ratelimitPolicyList dpv1alpha3.RateLimitPolicyList
	rateLimitPolicies := make(map[string]*dpv1alpha3.RateLimitPolicy)
	if err := gatewayReconciler.client.List(ctx, &ratelimitPolicyList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(gatewayRateLimitPolicyIndex, gatewayName.String()),
	}); err != nil {
		return nil, err
	}
	for _, item := range ratelimitPolicyList.Items {
		rateLimitPolicy := item
		rateLimitPolicies[utils.NamespacedName(&rateLimitPolicy).String()] = &rateLimitPolicy
	}
	return rateLimitPolicies, nil
}

// getGatewaysForAPIPolicy triggers the Gateway controller reconcile method
// based on the changes detected from APIPolicy objects.
func (gatewayReconciler *GatewayReconciler) getGatewaysForAPIPolicy(ctx context.Context, obj *dpv1alpha5.APIPolicy) []reconcile.Request {
	apiPolicy := obj

	if !(apiPolicy.Spec.TargetRef.Kind == constants.KindGateway) {
		return nil
	}

	namespace, err := utils.ValidateAndRetrieveNamespace((*gwapiv1.Namespace)(apiPolicy.Spec.TargetRef.Namespace), apiPolicy.Namespace)

	if err != nil {
		loggers.LoggerAPKOperator.Errorf("Namespace mismatch. TargetRef %s needs to be in the same namespace as the ApiPolicy %s. Expected: %s, Actual: %s",
			string(apiPolicy.Spec.TargetRef.Name), apiPolicy.Name, apiPolicy.Namespace, string(*apiPolicy.Spec.TargetRef.Namespace))
		return nil
	}

	return []reconcile.Request{{
		NamespacedName: types.NamespacedName{
			Name:      string(apiPolicy.Spec.TargetRef.Name),
			Namespace: namespace,
		},
	}}
}

// getGatewaysForTokenIssuer triggers the Gateway controller reconcile method
// based on the changes detected from TokenIssuer objects.
func (gatewayReconciler *GatewayReconciler) getGatewaysForTokenIssuer(ctx context.Context, obj *dpv1alpha1.TokenIssuer) []reconcile.Request {
	tokenIssuer := obj

	if !(tokenIssuer.Spec.TargetRef.Kind == constants.KindGateway) {
		return nil
	}

	namespace, err := utils.ValidateAndRetrieveNamespace((*gwapiv1.Namespace)(tokenIssuer.Spec.TargetRef.Namespace), tokenIssuer.Namespace)

	if err != nil {
		loggers.LoggerAPKOperator.Errorf("Namespace mismatch. TargetRef %s needs to be in the same namespace as the TokenIssuer %s. Expected: %s, Actual: %s",
			string(tokenIssuer.Spec.TargetRef.Name), tokenIssuer.Name, tokenIssuer.Namespace, string(*tokenIssuer.Spec.TargetRef.Namespace))
		return nil
	}

	return []reconcile.Request{{
		NamespacedName: types.NamespacedName{
			Name:      string(tokenIssuer.Spec.TargetRef.Name),
			Namespace: namespace,
		},
	}}
}

// addGatewayIndexes adds indexers related to Gateways
func addGatewayIndexes(ctx context.Context, mgr manager.Manager) error {
	// Gateway to RateLimitPolicy indexer
	if err := mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha3.RateLimitPolicy{}, gatewayRateLimitPolicyIndex,
		func(rawObj k8client.Object) []string {
			ratelimitPolicy := rawObj.(*dpv1alpha3.RateLimitPolicy)
			var gateways []string
			if ratelimitPolicy.Spec.TargetRef.Kind == constants.KindGateway {

				namespace, err := utils.ValidateAndRetrieveNamespace((*gwapiv1.Namespace)(ratelimitPolicy.Spec.TargetRef.Namespace), ratelimitPolicy.Namespace)

				if err != nil {
					loggers.LoggerAPKOperator.Errorf("Namespace mismatch. TargetRef %s needs to be in the same namespace as the RatelimitPolicy %s. Expected: %s, Actual: %s",
						string(ratelimitPolicy.Spec.TargetRef.Name), ratelimitPolicy.Name, ratelimitPolicy.Namespace, string(*ratelimitPolicy.Spec.TargetRef.Namespace))
					return gateways
				}

				gateways = append(gateways,
					types.NamespacedName{
						Namespace: namespace,
						Name:      string(ratelimitPolicy.Spec.TargetRef.Name),
					}.String())
			}
			return gateways
		}); err != nil {
		return err
	}

	// Gateway to APIPolicy indexer
	if err := mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha5.APIPolicy{}, gatewayAPIPolicyIndex,
		func(rawObj k8client.Object) []string {
			apiPolicy := rawObj.(*dpv1alpha5.APIPolicy)
			var httpRoutes []string
			if apiPolicy.Spec.TargetRef.Kind == constants.KindGateway {

				namespace, err := utils.ValidateAndRetrieveNamespace((*gwapiv1.Namespace)(apiPolicy.Spec.TargetRef.Namespace), apiPolicy.Namespace)

				if err != nil {
					loggers.LoggerAPKOperator.Errorf("Namespace mismatch. TargetRef %s needs to be in the same namespace as the ApiPolicy %s. Expected: %s, Actual: %s",
						string(apiPolicy.Spec.TargetRef.Name), apiPolicy.Name, apiPolicy.Namespace, string(*apiPolicy.Spec.TargetRef.Namespace))
					return httpRoutes
				}

				httpRoutes = append(httpRoutes,
					types.NamespacedName{
						Namespace: namespace,
						Name:      string(apiPolicy.Spec.TargetRef.Name),
					}.String())
			}
			return httpRoutes
		}); err != nil {
		return err
	}
	// Secret to TokenIssuer indexer
	if err := mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha1.TokenIssuer{}, secretTokenIssuerIndex,
		func(rawObj k8client.Object) []string {
			jwtIssuer := rawObj.(*dpv1alpha1.TokenIssuer)
			var secrets []string
			if jwtIssuer.Spec.SignatureValidation.Certificate != nil && jwtIssuer.Spec.SignatureValidation.Certificate.SecretRef != nil && len(jwtIssuer.Spec.SignatureValidation.Certificate.SecretRef.Name) > 0 {
				secrets = append(secrets,
					types.NamespacedName{
						Name:      string(jwtIssuer.Spec.SignatureValidation.Certificate.SecretRef.Name),
						Namespace: jwtIssuer.Namespace,
					}.String())
			}
			if jwtIssuer.Spec.SignatureValidation.JWKS != nil && jwtIssuer.Spec.SignatureValidation.JWKS.TLS != nil && jwtIssuer.Spec.SignatureValidation.JWKS.TLS.SecretRef != nil && len(jwtIssuer.Spec.SignatureValidation.JWKS.TLS.SecretRef.Name) > 0 {
				secrets = append(secrets,
					types.NamespacedName{
						Name:      string(jwtIssuer.Spec.SignatureValidation.JWKS.TLS.SecretRef.Name),
						Namespace: jwtIssuer.Namespace,
					}.String())
			}
			return secrets
		}); err != nil {
		return err
	}
	// Configmap to TokenIssuer indexer
	err := mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha1.TokenIssuer{}, configmapIssuerIndex,
		func(rawObj k8client.Object) []string {
			tokenIssuer := rawObj.(*dpv1alpha1.TokenIssuer)
			var configMaps []string
			if tokenIssuer.Spec.SignatureValidation.Certificate != nil && tokenIssuer.Spec.SignatureValidation.Certificate.ConfigMapRef != nil && len(tokenIssuer.Spec.SignatureValidation.Certificate.ConfigMapRef.Name) > 0 {
				configMaps = append(configMaps,
					types.NamespacedName{
						Name:      string(tokenIssuer.Spec.SignatureValidation.Certificate.ConfigMapRef.Name),
						Namespace: tokenIssuer.Namespace,
					}.String())
			}
			if tokenIssuer.Spec.SignatureValidation.JWKS != nil && tokenIssuer.Spec.SignatureValidation.JWKS.TLS != nil && tokenIssuer.Spec.SignatureValidation.JWKS.TLS.ConfigMapRef != nil && len(tokenIssuer.Spec.SignatureValidation.JWKS.TLS.ConfigMapRef.Name) > 0 {
				configMaps = append(configMaps,
					types.NamespacedName{
						Name:      string(tokenIssuer.Spec.SignatureValidation.JWKS.TLS.ConfigMapRef.Name),
						Namespace: tokenIssuer.Namespace,
					}.String())
			}
			return configMaps
		})
	return err
}
