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

	"github.com/wso2/apk/adapter/config"
	"github.com/wso2/apk/adapter/internal/loggers"
	"github.com/wso2/apk/adapter/pkg/logging"
	"golang.org/x/exp/maps"
	corev1 "k8s.io/api/core/v1"
	k8error "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/types"

	dpv1alpha1 "github.com/wso2/apk/adapter/internal/operator/apis/dp/v1alpha1"
	"github.com/wso2/apk/adapter/internal/operator/constants"
	"github.com/wso2/apk/adapter/internal/operator/status"
	"github.com/wso2/apk/adapter/internal/operator/synchronizer"
	"github.com/wso2/apk/adapter/internal/operator/utils"
	ctrl "sigs.k8s.io/controller-runtime"
	k8client "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	gwapiv1b1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

const (
	gatewayRateLimitPolicyIndex = "gatewayRateLimitPolicyIndex"
	gatewayAPIPolicyIndex       = "gatewayAPIPolicyIndex"
)

// GatewayReconciler reconciles a Gateway object
type GatewayReconciler struct {
	client        k8client.Client
	ods           *synchronizer.OperatorDataStore
	ch            *chan synchronizer.GatewayEvent
	statusUpdater *status.UpdateHandler
	mgr           manager.Manager
}

// NewGatewayController creates a new Gateway controller instance. Gateway Controllers watches for gwapiv1b1.Gateway.
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
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2610, err))
		return err
	}

	ctx := context.Background()
	conf := config.ReadConfigs()
	predicates := []predicate.Predicate{predicate.NewPredicateFuncs(utils.FilterByNamespaces(conf.Adapter.Operator.Namespaces))}

	if err := addGatewayIndexes(ctx, mgr); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2612, err))
		return err
	}

	if err := c.Watch(&source.Kind{Type: &gwapiv1b1.Gateway{}}, &handler.EnqueueRequestForObject{},
		predicates...); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(3100, err))
		return err
	}

	if err := c.Watch(&source.Kind{Type: &dpv1alpha1.RateLimitPolicy{}},
		handler.EnqueueRequestsFromMapFunc(r.handleCustomRateLimitPolicies), predicates...); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2611, err))
		return err
	}

	if err := c.Watch(&source.Kind{Type: &dpv1alpha1.APIPolicy{}}, handler.EnqueueRequestsFromMapFunc(r.getGatewaysForAPIPolicy),
		predicates...); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(3101, err))
		return err
	}

	if err := c.Watch(&source.Kind{Type: &dpv1alpha1.InterceptorService{}}, handler.EnqueueRequestsFromMapFunc(r.getAPIsForInterceptorService),
		predicates...); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(3110, err))
		return err
	}

	if err := c.Watch(&source.Kind{Type: &dpv1alpha1.Backend{}}, handler.EnqueueRequestsFromMapFunc(r.getGatewaysForBackend),
		predicates...); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(3102, err))
		return err
	}

	if err := c.Watch(&source.Kind{Type: &corev1.ConfigMap{}}, handler.EnqueueRequestsFromMapFunc(r.getGatewaysForConfigMap),
		predicates...); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(3103, err))
		return err
	}

	if err := c.Watch(&source.Kind{Type: &corev1.Secret{}}, handler.EnqueueRequestsFromMapFunc(r.getGatewaysForSecret),
		predicates...); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(3104, err))
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
	var gatewayDef gwapiv1b1.Gateway
	if err := gatewayReconciler.client.Get(ctx, req.NamespacedName, &gatewayDef); err != nil {
		gatewayState, found := gatewayReconciler.ods.GetCachedGateway(req.NamespacedName)
		if found && k8error.IsNotFound(err) {
			// The gateway doesn't exist in the gateway Cache, remove it
			gatewayReconciler.ods.DeleteCachedGateway(req.NamespacedName)
			loggers.LoggerAPKOperator.Infof("Delete event has received for Gateway : %s, hence deleted from Gateway cache", req.NamespacedName.String())
			*gatewayReconciler.ch <- synchronizer.GatewayEvent{EventType: constants.Delete, Event: gatewayState}
			return ctrl.Result{}, nil
		}
		loggers.LoggerAPKOperator.Warnf("Gateway CR related to the reconcile request with key: %s returned error. Assuming Gateway is already deleted, hence ignoring the error : %v", err)
		return ctrl.Result{}, nil
	}
	var gwCondition []metav1.Condition = gatewayDef.Status.Conditions

	gatewayStateData, err := gatewayReconciler.resolveGatewayState(ctx, gatewayDef)
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2620, req.NamespacedName.String(), err))
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
	return ctrl.Result{}, nil
}

// resolveListenerSecretRefs resolves the certificate secret references in the related listeners in Gateway CR
func (gatewayReconciler *GatewayReconciler) resolveListenerSecretRefs(ctx context.Context, secretRef *gwapiv1b1.SecretObjectReference, gatewayNamespace string) (map[string][]byte, error) {
	var secret corev1.Secret
	namespace := gwapiv1b1.Namespace(string(*secretRef.Namespace))
	if err := gatewayReconciler.client.Get(ctx, types.NamespacedName{Name: string(secretRef.Name),
		Namespace: utils.GetNamespace(&namespace, gatewayNamespace)}, &secret); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2612, secretRef.Name, string(*secretRef.Namespace), err))
		return nil, err
	}
	return secret.Data, nil
}

// resolveGatewayState resolves the GatewayState struct using gwapiv1b1.Gateway and resource indexes
func (gatewayReconciler *GatewayReconciler) resolveGatewayState(ctx context.Context,
	gateway gwapiv1b1.Gateway) (*synchronizer.GatewayStateData, error) {
	gatewayState := &synchronizer.GatewayStateData{}
	var err error
	resolvedListenerCerts := make(map[string]map[string][]byte)
	namespace := gwapiv1b1.Namespace(gateway.Namespace)
	// Retireve listener Certificates
	for _, listener := range gateway.Spec.Listeners {
		data, err := gatewayReconciler.resolveListenerSecretRefs(ctx, &listener.TLS.CertificateRefs[0], string(namespace))
		if err != nil {
			loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(3105, err))
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
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2650, err))
	}
	gatewayState.GatewayCustomRateLimitPolicies = customRateLimitPolicies
	gatewayState.GatewayBackendMapping = gatewayReconciler.getResolvedBackendsMapping(ctx, gatewayState)
	return gatewayState, nil
}

func (gatewayReconciler *GatewayReconciler) getAPIPoliciesForGateway(ctx context.Context,
	gateway *gwapiv1b1.Gateway) (map[string]dpv1alpha1.APIPolicy, error) {
	apiPolicies := make(map[string]dpv1alpha1.APIPolicy)
	apiPolicyList := &dpv1alpha1.APIPolicyList{}
	if err := gatewayReconciler.client.List(ctx, apiPolicyList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(gatewayAPIPolicyIndex, utils.NamespacedName(gateway).String()),
	}); err != nil {
		return nil, err
	}
	for _, apipolicy := range apiPolicyList.Items {
		apiPolicies[utils.NamespacedName(&apipolicy).String()] = apipolicy
	}
	return apiPolicies, nil
}

// getInterceptorServicesForGateway returns the list of interceptor services for the given gateway
func (gatewayReconciler *GatewayReconciler) getInterceptorServicesForGateway(ctx context.Context,
	gatewayAPIPolicies map[string]dpv1alpha1.APIPolicy) (map[string]dpv1alpha1.InterceptorService, error) {
	allGatewayAPIPolicies := maps.Values(gatewayAPIPolicies)
	interceptorServices := make(map[string]dpv1alpha1.InterceptorService)
	for _, apiPolicy := range allGatewayAPIPolicies {
		if apiPolicy.Spec.Default != nil && len(apiPolicy.Spec.Default.RequestInterceptors) > 0 {
			interceptorPtr := utils.GetInterceptorService(ctx, gatewayReconciler.client, &apiPolicy.Spec.Default.RequestInterceptors[0], nil, apiPolicy.Namespace)
			if interceptorPtr != nil {
				interceptorServices[utils.NamespacedName(interceptorPtr).String()] = *interceptorPtr
			}
		}
		if apiPolicy.Spec.Default != nil && len(apiPolicy.Spec.Default.ResponseInterceptors) > 0 {
			interceptorPtr := utils.GetInterceptorService(ctx, gatewayReconciler.client, &apiPolicy.Spec.Default.ResponseInterceptors[0], nil, apiPolicy.Namespace)
			if interceptorPtr != nil {
				interceptorServices[utils.NamespacedName(interceptorPtr).String()] = *interceptorPtr
			}
		}
		if apiPolicy.Spec.Override != nil && len(apiPolicy.Spec.Override.RequestInterceptors) > 0 {
			interceptorPtr := utils.GetInterceptorService(ctx, gatewayReconciler.client, &apiPolicy.Spec.Override.RequestInterceptors[0], nil, apiPolicy.Namespace)
			if interceptorPtr != nil {
				interceptorServices[utils.NamespacedName(interceptorPtr).String()] = *interceptorPtr
			}
		}
		if apiPolicy.Spec.Override != nil && len(apiPolicy.Spec.Override.ResponseInterceptors) > 0 {
			interceptorPtr := utils.GetInterceptorService(ctx, gatewayReconciler.client, &apiPolicy.Spec.Override.ResponseInterceptors[0], nil, apiPolicy.Namespace)
			if interceptorPtr != nil {
				interceptorServices[utils.NamespacedName(interceptorPtr).String()] = *interceptorPtr
			}
		}
	}
	return interceptorServices, nil
}

func (gatewayReconciler *GatewayReconciler) getResolvedBackendsMapping(ctx context.Context,
	gatewayStateData *synchronizer.GatewayStateData) dpv1alpha1.BackendMapping {
	backendMapping := make(dpv1alpha1.BackendMapping)
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
func (gatewayReconciler *GatewayReconciler) getGatewaysForBackend(obj k8client.Object) []reconcile.Request {
	ctx := context.Background()
	backend, ok := obj.(*dpv1alpha1.Backend)
	if !ok {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2624, backend))
		return []reconcile.Request{}
	}

	requests := []reconcile.Request{}

	interceptorServiceList := &dpv1alpha1.InterceptorServiceList{}
	if err := gatewayReconciler.client.List(ctx, interceptorServiceList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(backendInterceptorServiceIndex, utils.NamespacedName(backend).String()),
	}); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2649, utils.NamespacedName(backend).String()))
		return []reconcile.Request{}
	}

	for _, interceptorService := range interceptorServiceList.Items {
		requests = append(requests, gatewayReconciler.getAPIsForInterceptorService(&interceptorService)...)
	}

	return requests
}

// getAPIsForInterceptorService triggers the Gateway controller reconcile method based on the changes detected
// in InterceptorService resources.
func (gatewayReconciler *GatewayReconciler) getAPIsForInterceptorService(obj k8client.Object) []reconcile.Request {
	ctx := context.Background()
	interceptorService, ok := obj.(*dpv1alpha1.InterceptorService)
	if !ok {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(3107, interceptorService))
		return []reconcile.Request{}
	}

	requests := []reconcile.Request{}

	apiPolicyList := &dpv1alpha1.APIPolicyList{}
	if err := gatewayReconciler.client.List(ctx, apiPolicyList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(interceptorServiceAPIPolicyIndex, utils.NamespacedName(interceptorService).String()),
	}); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2649, utils.NamespacedName(interceptorService).String()))
		return []reconcile.Request{}
	}

	for _, apiPolicy := range apiPolicyList.Items {
		requests = append(requests, gatewayReconciler.getGatewaysForAPIPolicy(&apiPolicy)...)
	}

	return requests
}

// getGatewaysForSecret triggers the Gateway controller reconcile method based on the changes detected
// in secret resources.
func (gatewayReconciler *GatewayReconciler) getGatewaysForSecret(obj k8client.Object) []reconcile.Request {
	ctx := context.Background()
	secret, ok := obj.(*corev1.Secret)
	if !ok {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(3107, secret))
		return []reconcile.Request{}
	}

	backendList := &dpv1alpha1.BackendList{}
	if err := gatewayReconciler.client.List(ctx, backendList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(secretBackend, utils.NamespacedName(secret).String()),
	}); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(3106, utils.NamespacedName(secret).String()))
		return []reconcile.Request{}
	}

	requests := []reconcile.Request{}
	for _, backend := range backendList.Items {
		requests = append(requests, gatewayReconciler.getGatewaysForBackend(&backend)...)
	}
	return requests
}

// getGatewaysForConfigMap triggers the API controller reconcile method based on the changes detected
// in configMap resources.
func (gatewayReconciler *GatewayReconciler) getGatewaysForConfigMap(obj k8client.Object) []reconcile.Request {
	ctx := context.Background()
	configMap, ok := obj.(*corev1.ConfigMap)
	if !ok {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(3107, configMap))
		return []reconcile.Request{}
	}

	backendList := &dpv1alpha1.BackendList{}
	if err := gatewayReconciler.client.List(ctx, backendList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(configMapBackend, utils.NamespacedName(configMap).String()),
	}); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(3108, utils.NamespacedName(configMap).String()))
		return []reconcile.Request{}
	}

	requests := []reconcile.Request{}
	for _, backend := range backendList.Items {
		requests = append(requests, gatewayReconciler.getGatewaysForBackend(&backend)...)
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
		Resource:       new(gwapiv1b1.Gateway),
		UpdateStatus: func(obj k8client.Object) k8client.Object {
			h, ok := obj.(*gwapiv1b1.Gateway)
			if !ok {
				loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(3109, obj))
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
func (gatewayReconciler *GatewayReconciler) handleCustomRateLimitPolicies(obj k8client.Object) []reconcile.Request {
	ratelimitPolicy, ok := obj.(*dpv1alpha1.RateLimitPolicy)
	if !ok {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2622, ratelimitPolicy))
		return []reconcile.Request{}
	}
	requests := []reconcile.Request{}
	if ratelimitPolicy.Spec.TargetRef.Kind == constants.KindGateway {
		requests = append(requests, reconcile.Request{
			NamespacedName: types.NamespacedName{
				Namespace: ratelimitPolicy.Namespace,
				Name:      string(ratelimitPolicy.Spec.TargetRef.Name),
			},
		})
	}
	return requests
}

// getCustomRateLimitPoliciesForGateway returns the list of custom rate limit policies for a gateway
func (gatewayReconciler *GatewayReconciler) getCustomRateLimitPoliciesForGateway(gatewayName types.NamespacedName) ([]*dpv1alpha1.RateLimitPolicy, error) {
	ctx := context.Background()
	var ratelimitPolicyList dpv1alpha1.RateLimitPolicyList
	var rateLimitPolicies []*dpv1alpha1.RateLimitPolicy
	if err := gatewayReconciler.client.List(ctx, &ratelimitPolicyList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(gatewayRateLimitPolicyIndex, gatewayName.String()),
	}); err != nil {
		return nil, err
	}
	for _, item := range ratelimitPolicyList.Items {
		rateLimitPolicy := item
		rateLimitPolicies = append(rateLimitPolicies, &rateLimitPolicy)
	}
	return rateLimitPolicies, nil
}

// getGatewaysForAPIPolicy triggers the Gateway controller reconcile method
// based on the changes detected from APIPolicy objects.
func (gatewayReconciler *GatewayReconciler) getGatewaysForAPIPolicy(obj k8client.Object) []reconcile.Request {
	apiPolicy, ok := obj.(*dpv1alpha1.APIPolicy)
	if !ok {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2624, apiPolicy))
		return nil
	}

	if !(apiPolicy.Spec.TargetRef.Kind == constants.KindGateway) {
		return nil
	}

	return []reconcile.Request{{
		NamespacedName: types.NamespacedName{
			Name: string(apiPolicy.Spec.TargetRef.Name),
			Namespace: utils.GetNamespace(
				(*gwapiv1b1.Namespace)(apiPolicy.Spec.TargetRef.Namespace), apiPolicy.Namespace),
		},
	}}
}

// addGatewayIndexes adds indexers related to Gateways
func addGatewayIndexes(ctx context.Context, mgr manager.Manager) error {
	// Gateway to RateLimitPolicy indexer
	if err := mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha1.RateLimitPolicy{}, gatewayRateLimitPolicyIndex,
		func(rawObj k8client.Object) []string {
			ratelimitPolicy := rawObj.(*dpv1alpha1.RateLimitPolicy)
			var gateways []string
			if ratelimitPolicy.Spec.TargetRef.Kind == constants.KindGateway {
				gateways = append(gateways,
					types.NamespacedName{
						Namespace: utils.GetNamespace(
							(*gwapiv1b1.Namespace)(ratelimitPolicy.Spec.TargetRef.Namespace),
							ratelimitPolicy.Namespace),
						Name: string(ratelimitPolicy.Spec.TargetRef.Name),
					}.String())
			}
			return gateways
		}); err != nil {
		return err
	}

	// Gateway to APIPolicy indexer
	err := mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha1.APIPolicy{}, gatewayAPIPolicyIndex,
		func(rawObj k8client.Object) []string {
			apiPolicy := rawObj.(*dpv1alpha1.APIPolicy)
			var httpRoutes []string
			if apiPolicy.Spec.TargetRef.Kind == constants.KindGateway {
				httpRoutes = append(httpRoutes,
					types.NamespacedName{
						Namespace: utils.GetNamespace(
							(*gwapiv1b1.Namespace)(apiPolicy.Spec.TargetRef.Namespace), apiPolicy.Namespace),
						Name: string(apiPolicy.Spec.TargetRef.Name),
					}.String())
			}
			return httpRoutes
		})
	return err
}
