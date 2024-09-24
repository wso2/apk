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

	"github.com/wso2/apk/adapter/internal/operator/constants"
	"github.com/wso2/apk/adapter/internal/operator/status"
	"github.com/wso2/apk/adapter/internal/operator/synchronizer"
	"github.com/wso2/apk/adapter/internal/operator/utils"
	dpv1alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha1"
	dpv1alpha2 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha2"
	dpv1alpha3 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha3"
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
)

var (
	setReadiness sync.Once
)

// GatewayReconciler reconciles a Gateway object
type GatewayReconciler struct {
	client        k8client.Client
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
	predicates := []predicate.Predicate{predicate.NewPredicateFuncs(utils.FilterByNamespaces(conf.Adapter.Operator.Namespaces))}

	if err := addGatewayIndexes(ctx, mgr); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error3120, logging.BLOCKER, "Error adding indexes: %v", err))
		return err
	}

	if err := c.Watch(source.Kind(mgr.GetCache(), &gwapiv1.Gateway{}), &handler.EnqueueRequestForObject{},
		predicates...); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error3100, logging.BLOCKER, "Error watching Gateway resources: %v", err))
		return err
	}

	if err := c.Watch(source.Kind(mgr.GetCache(), &dpv1alpha3.RateLimitPolicy{}),
		handler.EnqueueRequestsFromMapFunc(r.handleCustomRateLimitPolicies), predicates...); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error3121, logging.BLOCKER, "Error watching Ratelimit resources: %v", err))
		return err
	}

	if err := c.Watch(source.Kind(mgr.GetCache(), &dpv1alpha3.APIPolicy{}), handler.EnqueueRequestsFromMapFunc(r.getGatewaysForAPIPolicy),
		predicates...); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error3101, logging.BLOCKER, "Error watching APIPolicy resources: %v", err))
		return err
	}

	if err := c.Watch(source.Kind(mgr.GetCache(), &dpv1alpha1.InterceptorService{}), handler.EnqueueRequestsFromMapFunc(r.getAPIsForInterceptorService),
		predicates...); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error3110, logging.BLOCKER, "Error watching InterceptorService resources: %v", err))
		return err
	}

	if err := c.Watch(source.Kind(mgr.GetCache(), &dpv1alpha1.BackendJWT{}), handler.EnqueueRequestsFromMapFunc(r.getAPIsForBackendJWT),
		predicates...); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error3126, logging.BLOCKER, "Error watching BackendJWT resources: %v", err))
		return err
	}

	if err := c.Watch(source.Kind(mgr.GetCache(), &dpv1alpha2.Backend{}), handler.EnqueueRequestsFromMapFunc(r.getGatewaysForBackend),
		predicates...); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error3102, logging.BLOCKER, "Error watching Backend resources: %v", err))
		return err
	}

	if err := c.Watch(source.Kind(mgr.GetCache(), &corev1.ConfigMap{}), handler.EnqueueRequestsFromMapFunc(r.getGatewaysForConfigMap),
		predicates...); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error3103, logging.BLOCKER, "Error watching ConfigMap resources: %v", err))
		return err
	}

	if err := c.Watch(source.Kind(mgr.GetCache(), &corev1.Secret{}), handler.EnqueueRequestsFromMapFunc(r.getGatewaysForSecret),
		predicates...); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error3104, logging.BLOCKER, "Error watching Secret resources: %v", err))
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
	gatewayState.GatewayCustomRateLimitPolicies = customRateLimitPolicies
	gatewayState.GatewayBackendMapping = gatewayReconciler.getResolvedBackendsMapping(ctx, gatewayState)
	return gatewayState, nil
}

func (gatewayReconciler *GatewayReconciler) getAPIPoliciesForGateway(ctx context.Context,
	gateway *gwapiv1.Gateway) (map[string]dpv1alpha3.APIPolicy, error) {
	apiPolicies := make(map[string]dpv1alpha3.APIPolicy)
	apiPolicyList := &dpv1alpha3.APIPolicyList{}
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
	gatewayAPIPolicies map[string]dpv1alpha3.APIPolicy) (map[string]dpv1alpha1.InterceptorService, error) {
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
	gatewayStateData *synchronizer.GatewayStateData) map[string]*dpv1alpha2.ResolvedBackend {
	backendMapping := make(map[string]*dpv1alpha2.ResolvedBackend)
	if gatewayStateData.GatewayInterceptorServiceMapping != nil {
		interceptorServices := maps.Values(gatewayStateData.GatewayInterceptorServiceMapping)
		for _, interceptorService := range interceptorServices {
			utils.ResolveAndAddBackendToMapping(ctx, gatewayReconciler.client, backendMapping,
				interceptorService.Spec.BackendRef, interceptorService.Namespace, nil)
		}
	}
	return backendMapping
}

// getGatewaysForBackend triggers the Gateway controller reconcile method based on the changes detected
// in backend resources.
func (gatewayReconciler *GatewayReconciler) getGatewaysForBackend(ctx context.Context, obj k8client.Object) []reconcile.Request {
	backend, ok := obj.(*dpv1alpha2.Backend)
	if !ok {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error3107, logging.TRIVIAL, "Unexpected object type, bypassing reconciliation: %v", backend))
		return []reconcile.Request{}
	}

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
func (gatewayReconciler *GatewayReconciler) getAPIsForInterceptorService(ctx context.Context, obj k8client.Object) []reconcile.Request {
	interceptorService, ok := obj.(*dpv1alpha1.InterceptorService)
	if !ok {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error3107, logging.TRIVIAL, "Unexpected object type, bypassing reconciliation: %v", interceptorService))
		return []reconcile.Request{}
	}

	requests := []reconcile.Request{}

	apiPolicyList := &dpv1alpha3.APIPolicyList{}
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
func (gatewayReconciler *GatewayReconciler) getAPIsForBackendJWT(ctx context.Context, obj k8client.Object) []reconcile.Request {
	backendJWT, ok := obj.(*dpv1alpha1.BackendJWT)
	if !ok {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error3107, logging.TRIVIAL, "Unexpected object type, bypassing reconciliation: %v", backendJWT))
		return []reconcile.Request{}
	}

	requests := []reconcile.Request{}

	apiPolicyList := &dpv1alpha3.APIPolicyList{}
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
func (gatewayReconciler *GatewayReconciler) getGatewaysForSecret(ctx context.Context, obj k8client.Object) []reconcile.Request {
	secret, ok := obj.(*corev1.Secret)
	if !ok {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error3107, logging.TRIVIAL, "Unexpected object type, bypassing reconciliation: %v", secret))
		return []reconcile.Request{}
	}

	backendList := &dpv1alpha2.BackendList{}
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
	return requests
}

// getGatewaysForConfigMap triggers the API controller reconcile method based on the changes detected
// in configMap resources.
func (gatewayReconciler *GatewayReconciler) getGatewaysForConfigMap(ctx context.Context, obj k8client.Object) []reconcile.Request {
	configMap, ok := obj.(*corev1.ConfigMap)
	if !ok {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error3107, logging.TRIVIAL, "Unexpected object type, bypassing reconciliation: %v", configMap))
		return []reconcile.Request{}
	}

	backendList := &dpv1alpha2.BackendList{}
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
func (gatewayReconciler *GatewayReconciler) handleCustomRateLimitPolicies(ctx context.Context, obj k8client.Object) []reconcile.Request {
	ratelimitPolicy, ok := obj.(*dpv1alpha3.RateLimitPolicy)
	if !ok {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error3107, logging.TRIVIAL, "Unexpected object type, bypassing reconciliation: %v", ratelimitPolicy))
		return []reconcile.Request{}
	}
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
func (gatewayReconciler *GatewayReconciler) getGatewaysForAPIPolicy(ctx context.Context, obj k8client.Object) []reconcile.Request {
	apiPolicy, ok := obj.(*dpv1alpha3.APIPolicy)
	if !ok {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error3107, logging.TRIVIAL, "Unexpected object type, bypassing reconciliation: %v", apiPolicy))
		return nil
	}

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
	err := mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha3.APIPolicy{}, gatewayAPIPolicyIndex,
		func(rawObj k8client.Object) []string {
			apiPolicy := rawObj.(*dpv1alpha3.APIPolicy)
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
		})
	return err
}
