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
	"github.com/wso2/apk/adapter/internal/discovery/xds/common"
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
	setReadiness   sync.Once
	supportedKinds = []gwapiv1.Kind{gwapiv1.Kind("HTTPRoute")}
	controllerName = "wso2.com/apk-envoy"
)

// GetControllerName returns the controller name that supported by APK
func GetControllerName() string {
	return controllerName
}

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

	if err := c.Watch(source.Kind(mgr.GetCache(), &dpv1alpha1.RateLimitPolicy{}),
		handler.EnqueueRequestsFromMapFunc(r.handleCustomRateLimitPolicies), predicates...); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error3121, logging.BLOCKER, "Error watching Ratelimit resources: %v", err))
		return err
	}

	if err := c.Watch(source.Kind(mgr.GetCache(), &dpv1alpha2.APIPolicy{}), handler.EnqueueRequestsFromMapFunc(r.getGatewaysForAPIPolicy),
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

	if err := c.Watch(source.Kind(mgr.GetCache(), &dpv1alpha1.Backend{}), handler.EnqueueRequestsFromMapFunc(r.getGatewaysForBackend),
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

	if err := c.Watch(source.Kind(mgr.GetCache(), &gwapiv1.HTTPRoute{}),
		handler.EnqueueRequestsFromMapFunc(r.getHTTPRoutes), predicates...); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error3121, logging.BLOCKER, "Error watching HttpRoutes resources: %v", err))
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
	loggers.LoggerAPKOperator.Infof("Reconciling gateway... %s", req.NamespacedName.String())
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

	gatewayStateData, listenerStatuses, err := gatewayReconciler.resolveGatewayState(ctx, gatewayDef)
	// Check whether the status change is needed for gateway
	statusChanged := isStatusChanged(gatewayDef, listenerStatuses)
	loggers.LoggerAPKOperator.Infof("Status changed ? %+v", statusChanged)

	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error3122, logging.BLOCKER, "Error resolving Gateway State %s: %v", req.NamespacedName.String(), err))
		return ctrl.Result{}, err
	}
	state := constants.Update
	var (
		events        = make([]string, 0)
		updated       = false
		cachedGateway synchronizer.GatewayState
	)

	if gwCondition[0].Status != metav1.ConditionTrue {
		gatewayState := gatewayReconciler.ods.AddGatewayState(gatewayDef, gatewayStateData)
		*gatewayReconciler.ch <- synchronizer.GatewayEvent{EventType: constants.Create, Event: gatewayState}
		state = constants.Create
	} else if cachedGateway, events, updated =
		gatewayReconciler.ods.UpdateGatewayState(&gatewayDef, gatewayStateData); updated {
		*gatewayReconciler.ch <- synchronizer.GatewayEvent{EventType: constants.Update, Event: cachedGateway}
		state = constants.Update
	}
	if statusChanged || updated {
		loggers.LoggerAPKOperator.Infof("Updating gateway status. Gateway: %s", utils.NamespacedName(&gatewayDef))
		gatewayReconciler.handleGatewayStatus(req.NamespacedName, state, events, listenerStatuses)
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
	gateway gwapiv1.Gateway) (*synchronizer.GatewayStateData, []gwapiv1.ListenerStatus, error) {
	gatewayState := &synchronizer.GatewayStateData{}
	var err error
	resolvedListenerCerts := make(map[string]map[string][]byte)
	namespace := gwapiv1.Namespace(gateway.Namespace)
	listenerstatuses := make([]gwapiv1.ListenerStatus, 0)
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error3124, logging.MAJOR, "Error while getting http routes: %s", err))
	}
	// Retireve listener Certificates
	for _, listener := range gateway.Spec.Listeners {
		accepted := true
		attachedRouteCount, err := getAttachedRoutesCountForListener(ctx, gatewayReconciler.client, gateway, string(listener.Name))
		if err != nil {
			attachedRouteCount = 0
		}
		listenerStatus := gwapiv1.ListenerStatus{
			Name:           listener.Name,
			SupportedKinds: []gwapiv1.RouteGroupKind{},
			Conditions:     []metav1.Condition{},
			AttachedRoutes: attachedRouteCount,
		}

		listenerDefinedAllowedKinds := listener.AllowedRoutes.Kinds
		actualKinds := make([]gwapiv1.Kind, len(listenerDefinedAllowedKinds))
		for i, obj := range listenerDefinedAllowedKinds {
			actualKinds[i] = obj.Kind
		}
		intersectionKinds := findIntersectionKinds(supportedKinds, actualKinds)
		if len(intersectionKinds) == 0 {
			// If listener does not define any supported kinds then we need to support all of the default supported kinds by the implementation
			intersectionKinds = supportedKinds
		}
		for _, kind := range intersectionKinds {
			listenerStatus.SupportedKinds = append(listenerStatus.SupportedKinds, gwapiv1.RouteGroupKind{
				Group: (*gwapiv1.Group)(&gwapiv1.GroupVersion.Group),
				Kind:  kind,
			})
		}

		if len(intersectionKinds) < len(listenerDefinedAllowedKinds) {
			accepted = false
			listenerStatus.Conditions = append(listenerStatus.Conditions, metav1.Condition{
				Type:               string(gwapiv1.ListenerConditionResolvedRefs),
				Status:             metav1.ConditionFalse,
				Reason:             string(gwapiv1.ListenerReasonInvalidRouteKinds),
				LastTransitionTime: metav1.Now(),
				ObservedGeneration: gateway.Generation,
			})
		}

		if listener.Protocol == gwapiv1.HTTPProtocolType {
			continue
		}
		data, err := gatewayReconciler.resolveListenerSecretRefs(ctx, &listener.TLS.CertificateRefs[0], string(namespace))
		if err != nil {
			accepted = false
			listenerStatus.Conditions = append(listenerStatus.Conditions, metav1.Condition{
				Type:               string(gwapiv1.ListenerConditionResolvedRefs),
				Status:             metav1.ConditionFalse,
				Reason:             string(gwapiv1.ListenerReasonInvalidCertificateRef),
				LastTransitionTime: metav1.Now(),
				ObservedGeneration: gateway.Generation,
			})
			loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error3105, logging.BLOCKER, "Error resolving listener certificates: %v", err))
			return nil, listenerstatuses, err
		}
		resolvedListenerCerts[string(listener.Name)] = data
		if accepted {
			listenerStatus.Conditions = append(listenerStatus.Conditions, metav1.Condition{
				Type:               string(gwapiv1.ListenerConditionAccepted),
				Status:             metav1.ConditionTrue,
				Reason:             string(gwapiv1.ListenerReasonResolvedRefs),
				LastTransitionTime: metav1.Now(),
				ObservedGeneration: gateway.Generation,
			})
			listenerStatus.Conditions = append(listenerStatus.Conditions, metav1.Condition{
				Type:               string(gwapiv1.ListenerConditionProgrammed),
				Status:             metav1.ConditionTrue,
				Reason:             string(gwapiv1.ListenerReasonResolvedRefs),
				LastTransitionTime: metav1.Now(),
				ObservedGeneration: gateway.Generation,
			})

		}
		listenerstatuses = append(listenerstatuses, listenerStatus)
		loggers.LoggerAPKOperator.Debugf("A listener status is added for listener:  %s", string(listenerStatus.Name))
	}
	gatewayState.GatewayResolvedListenerCerts = resolvedListenerCerts
	if gatewayState.GatewayAPIPolicies, err = gatewayReconciler.getAPIPoliciesForGateway(ctx, &gateway); err != nil {
		return nil, listenerstatuses, fmt.Errorf("error while getting gateway apipolicy for gateway: %s, %s", utils.NamespacedName(&gateway).String(), err.Error())
	}
	if gatewayState.GatewayInterceptorServiceMapping, err = gatewayReconciler.getInterceptorServicesForGateway(ctx, gatewayState.GatewayAPIPolicies); err != nil {
		return nil, listenerstatuses, fmt.Errorf("error while getting interceptor service for gateway: %s, %s", utils.NamespacedName(&gateway).String(), err.Error())
	}
	customRateLimitPolicies, err := gatewayReconciler.getCustomRateLimitPoliciesForGateway(utils.NamespacedName(&gateway))
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error3124, logging.MAJOR, "Error while getting custom rate limit policies: %s", err))
	}
	gatewayState.GatewayCustomRateLimitPolicies = customRateLimitPolicies
	gatewayState.GatewayBackendMapping = gatewayReconciler.getResolvedBackendsMapping(ctx, gatewayState)
	return gatewayState, listenerstatuses, nil
}

func (gatewayReconciler *GatewayReconciler) getAPIPoliciesForGateway(ctx context.Context,
	gateway *gwapiv1.Gateway) (map[string]dpv1alpha2.APIPolicy, error) {
	apiPolicies := make(map[string]dpv1alpha2.APIPolicy)
	apiPolicyList := &dpv1alpha2.APIPolicyList{}
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
	gatewayAPIPolicies map[string]dpv1alpha2.APIPolicy) (map[string]dpv1alpha1.InterceptorService, error) {
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
	gatewayStateData *synchronizer.GatewayStateData) map[string]*dpv1alpha1.ResolvedBackend {
	backendMapping := make(map[string]*dpv1alpha1.ResolvedBackend)
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
	backend, ok := obj.(*dpv1alpha1.Backend)
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

// getHTTPRoutes returns the list of gateway reconcile requests
func (gatewayReconciler *GatewayReconciler) getHTTPRoutes(ctx context.Context, obj k8client.Object) []reconcile.Request {
	httpRoute, ok := obj.(*gwapiv1.HTTPRoute)
	if !ok {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error3107, logging.TRIVIAL, "Unexpected object type, bypassing reconciliation: %v", httpRoute))
		return []reconcile.Request{}
	}
	requests := []reconcile.Request{}
	for _, refs := range httpRoute.Spec.ParentRefs {
		if *refs.Kind == constants.KindGateway {
			namespace := ""
			if refs.Namespace != nil {
				namespace = string(*refs.Namespace)
			}
			if namespace == "" {
				namespace = httpRoute.Namespace
			}
			requests = append(requests, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Namespace: namespace,
					Name:      string(refs.Name),
				},
			})
		}
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

	apiPolicyList := &dpv1alpha2.APIPolicyList{}
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

	apiPolicyList := &dpv1alpha2.APIPolicyList{}
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

	backendList := &dpv1alpha1.BackendList{}
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

	backendList := &dpv1alpha1.BackendList{}
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
func (gatewayReconciler *GatewayReconciler) handleGatewayStatus(gatewayKey types.NamespacedName, state string,
	events []string, listeners []gwapiv1.ListenerStatus) {
	message := ""

	switch state {
	case constants.Create:
		message = "Gateway is deployed successfully"
	case constants.Update:
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
			generation := hCopy.ObjectMeta.Generation
			gwCondition[0].Status = "True"
			gwCondition[0].Message = message
			gwCondition[0].LastTransitionTime = timeNow
			// gwCondition[0].Reason = append(gwCondition[0].Reason, event)
			gwCondition[0].Reason = "Reconciled"
			gwCondition[0].Type = constants.Accept
			for i := range gwCondition {
				// Assign generation to ObservedGeneration
				gwCondition[i].ObservedGeneration = generation
			}
			hCopy.Status.Conditions = gwCondition
			for _, listener := range hCopy.Status.Listeners {
				for _, listener1 := range listeners {
					if string(listener.Name) == string(listener1.Name) {
						listener1.AttachedRoutes = listener.AttachedRoutes
					}
				}
			}
			hCopy.Status.Listeners = listeners
			return hCopy
		},
	})
}

// handleCustomRateLimitPolicies returns the list of gateway reconcile requests
func (gatewayReconciler *GatewayReconciler) handleCustomRateLimitPolicies(ctx context.Context, obj k8client.Object) []reconcile.Request {
	ratelimitPolicy, ok := obj.(*dpv1alpha1.RateLimitPolicy)
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
func (gatewayReconciler *GatewayReconciler) getCustomRateLimitPoliciesForGateway(gatewayName types.NamespacedName) (map[string]*dpv1alpha1.RateLimitPolicy, error) {
	ctx := context.Background()
	var ratelimitPolicyList dpv1alpha1.RateLimitPolicyList
	rateLimitPolicies := make(map[string]*dpv1alpha1.RateLimitPolicy)
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
	apiPolicy, ok := obj.(*dpv1alpha2.APIPolicy)
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
	if err := mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha1.RateLimitPolicy{}, gatewayRateLimitPolicyIndex,
		func(rawObj k8client.Object) []string {
			ratelimitPolicy := rawObj.(*dpv1alpha1.RateLimitPolicy)
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
	err := mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha2.APIPolicy{}, gatewayAPIPolicyIndex,
		func(rawObj k8client.Object) []string {
			apiPolicy := rawObj.(*dpv1alpha2.APIPolicy)
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

func findIntersectionKinds(list1, list2 []gwapiv1.Kind) []gwapiv1.Kind {
	intersection := []gwapiv1.Kind{}
	set := make(map[string]bool)
	for _, v := range list1 {
		set[string(v)] = true
	}
	for _, v := range list2 {
		if set[string(v)] {
			intersection = append(intersection, v)
		}
	}
	return intersection
}

// findDiffFromSecondListKinds return a list of elements in list2 that are not in the list1
func findDiffFromSecondListKinds(list1, list2 []gwapiv1.Kind) []gwapiv1.Kind {
	diff := []gwapiv1.Kind{}
	set := make(map[string]bool)
	for _, v := range list1 {
		set[string(v)] = true
	}
	for _, v := range list2 {
		if !set[string(v)] {
			diff = append(diff, v)
		}
	}
	return diff
}

// getAttachedRoutesForListener returns the attached route count for a specific listener in a gatway
func getAttachedRoutesCountForListener(ctx context.Context, client k8client.Client, gateway gwapiv1.Gateway, listenerName string) (int32, error) {
	httpRouteList := gwapiv1.HTTPRouteList{}
	if err := client.List(ctx, &httpRouteList); err != nil {
		return 0, err
	}

	var attachedRoutesCount int32
	for _, httpRoute := range httpRouteList.Items {
		_, found := common.FindElement(httpRoute.Status.Parents, func(parentStatus gwapiv1.RouteParentStatus) bool {
			parentNamespacedName := types.NamespacedName{
				Namespace: string(*parentStatus.ParentRef.Namespace),
				Name:      string(parentStatus.ParentRef.Name),
			}.String()
			gatewayNamespacedName := utils.NamespacedName(&gateway).String()
			if parentNamespacedName == gatewayNamespacedName {
				if len(parentStatus.Conditions) >= 1 && parentStatus.Conditions[0].Status == metav1.ConditionTrue {
					// Check whether the listername matches
					_, matched := common.FindElement(httpRoute.Spec.ParentRefs, func(parentRef gwapiv1.ParentReference) bool {
						if string(*parentRef.SectionName) == listenerName {
							return true
						}
						return false
					})
					return matched
				}
			}
			return false
		})
		if found {
			attachedRoutesCount++
		}
	}
	return attachedRoutesCount, nil
}

func isStatusChanged(gateway gwapiv1.Gateway, statuses []gwapiv1.ListenerStatus) bool {
	if len(gateway.Status.Listeners) != len(statuses) {
		return true
	}
	for _, status1 := range gateway.Status.Listeners {
		flag := false
		for _, status2 := range statuses {
			if status1.Name == status2.Name &&
				status1.AttachedRoutes == status2.AttachedRoutes &&
				len(status1.Conditions) == len(status2.Conditions) &&
				len(status1.SupportedKinds) == len(status2.SupportedKinds) {
				flag = common.BothListContainsSameConditions(status1.Conditions, status2.Conditions)
				if flag {
					continue
				}
			}
		}
		if !flag {
			return true
		}
	}
	return false
}
