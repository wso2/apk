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
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	eg "github.com/envoyproxy/gateway/api/v1alpha1"

	"github.com/wso2/apk/adapter/pkg/logging"
	cache "github.com/wso2/apk/common-controller/internal/cache"
	"github.com/wso2/apk/common-controller/internal/config"
	controlplane "github.com/wso2/apk/common-controller/internal/controlplane"
	loggers "github.com/wso2/apk/common-controller/internal/loggers"
	"github.com/wso2/apk/common-controller/internal/operator/synchronizer"
	syncronizer "github.com/wso2/apk/common-controller/internal/operator/synchronizer"
	"github.com/wso2/apk/common-controller/internal/utils"
	"github.com/wso2/apk/common-go-libs/constants"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	k8client "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/wso2/apk/common-go-libs/apis/dp/v1alpha4"
	dpV2alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v2alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RouteMetadataReconciler reconciles a RouteMetadata object
type RouteMetadataReconciler struct {
	client client.Client
	Scheme *runtime.Scheme
	Store  *cache.RouteMetadataDataStore
}

const (
	// configMapIndex is the index for ConfigMap resources in the RouteMetadata controller
	configMapIndexRouteMetadata = "ConfigMapIndexRouteMetadata"
)

// NewRouteMetadataController creates a new controller for RouteMetadata.
func NewRouteMetadataController(mgr manager.Manager, store *cache.RouteMetadataDataStore) error {
	reconciler := &RouteMetadataReconciler{
		client: mgr.GetClient(),
		Store:  store,
	}

	ctx := context.Background()
	if err := reconciler.addRouteMetadataIndexes(ctx, mgr); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2612, logging.BLOCKER, "Error adding indexes: %v", err))
		return err
	}

	c, err := controller.New(constants.RouteMetadataController, mgr, controller.Options{Reconciler: reconciler})
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2664, logging.BLOCKER,
			"Error creating RouteMetadata controller: %v", err.Error()))
		return err
	}

	conf := config.ReadConfigs()
	predicateRouteMetadata := []predicate.TypedPredicate[*dpV2alpha1.RouteMetadata]{
		predicate.NewTypedPredicateFuncs[*dpV2alpha1.RouteMetadata](
			utils.FilterRouteMetadataByNamespaces(conf.CommonController.Operator.Namespaces),
		),
	}

	if err := c.Watch(source.Kind(mgr.GetCache(), &dpV2alpha1.RouteMetadata{}, &handler.TypedEnqueueRequestForObject[*dpV2alpha1.RouteMetadata]{}, predicateRouteMetadata...)); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2638, logging.BLOCKER,
			"Error watching RouteMetadata resources: %v", err.Error()))
		return err
	}

	predicateConfigMap := []predicate.TypedPredicate[*corev1.ConfigMap]{predicate.NewTypedPredicateFuncs(utils.FilterConfigMapByNamespaces(conf.CommonController.Operator.Namespaces))}
	if err := c.Watch(source.Kind(mgr.GetCache(), &corev1.ConfigMap{},
		handler.TypedEnqueueRequestsFromMapFunc(reconciler.getRouteMetadataForConfigMap), predicateConfigMap...)); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2613, logging.BLOCKER,
			"Error watching ConfigMap resources: %v", err))
		return err
	}

	loggers.LoggerAPKOperator.Debug("RouteMetadata Controller successfully started. Watching RouteMetadata Objects...")
	return nil
}

//+kubebuilder:rbac:groups=dp.wso2.com,resources=routepolicies,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=dp.wso2.com,resources=routepolicies/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=dp.wso2.com,resources=routepolicies/finalizers,verbs=update

// Reconcile reconciles the RouteMetadata CR
func (routeMetadataReconciler *RouteMetadataReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	loggers.LoggerAPKOperator.Infof("Reconciling RouteMetadata: %s", req.NamespacedName)
	routeMetadataKey := req.NamespacedName

	var routeMetadata dpV2alpha1.RouteMetadata
	if err := routeMetadataReconciler.client.Get(ctx, routeMetadataKey, &routeMetadata); err != nil {
		loggers.LoggerAPKOperator.Warnf("RouteMetadata %s not found, might be deleted", routeMetadataKey)
		routeMetadataReconciler.Store.DeleteRouteMetadata(routeMetadataKey.Namespace, routeMetadataKey.Name)
		routeMetadata.ObjectMeta = metav1.ObjectMeta{
			Namespace: routeMetadataKey.Namespace,
			Name:      routeMetadataKey.Name,
		}
		routePolicyString, err := utils.ToJSONString(routeMetadata)
		if err != nil {
			loggers.LoggerAPKOperator.Errorf("Error converting RouteMetadata to JSON: %v", err)
		} else {
			utils.SendRouteMetadataDeletedEvent(routePolicyString)
			loggers.LoggerAPKOperator.Debugf("Deleted RouteMetadata JSON: %s", routePolicyString)
		}
		// utils.SendRouteMetadataDeletedEvent()
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if routeMetadata.Spec.API.DefinitionFileRef == nil {
		namespacedName := types.NamespacedName{
			Namespace: req.Namespace,
			Name:      string(routeMetadata.Spec.API.DefinitionFileRef.Name),
		}
		var cm corev1.ConfigMap
		if err := routeMetadataReconciler.client.Get(ctx, namespacedName, &cm); err != nil {
			loggers.LoggerAPKOperator.Errorf("failed to fetch ConfigMap %s: %v", namespacedName.String(), err)
		}
		val, ok := cm.Data["Definition"]
		if !ok {
			loggers.LoggerAPKOperator.Warnf("key %s not found in ConfigMap %s", "Definition", namespacedName.String())
		}
		routeMetadata.Spec.API.Definition = val
	}

	// Add or update the RouteMetadata in the store
	routeMetadataReconciler.Store.AddOrUpdateRouteMetadata(routeMetadata)

	if val, ok := routeMetadata.Labels["initiatedfromCP"]; !ok || val == "false" {
		loggers.LoggerAPKOperator.Infof("Processing RouteMetadata %s as it's initiated from DP", req.NamespacedName)
		// your processing logic here
		state := synchronizer.APIState{
			APIDefinition: &dpV2alpha1.RouteMetadata{
				Spec: dpV2alpha1.RouteMetadataSpec{
					API: dpV2alpha1.API{}, // assuming API is a value type, not a pointer
				},
			},
		}
		state.APIDefinition.Name = routeMetadata.ObjectMeta.Name
		state.APIDefinition.Namespace = routeMetadata.ObjectMeta.Namespace
		state.APIDefinition.Spec.API.Name = routeMetadata.Spec.API.Name
		state.APIDefinition.Spec.API.Version = routeMetadata.Spec.API.Version
		state.APIDefinition.Spec.API.Type = routeMetadata.Spec.API.Type
		state.APIDefinition.Spec.API.Context = routeMetadata.Spec.API.Context
		state.APIDefinition.Spec.API.Organization = routeMetadata.Spec.API.Organization
		state.APIDefinition.Spec.API.Environment = routeMetadata.Spec.API.Environment
		state.APIDefinition.Spec.API.DefinitionFileRef = routeMetadata.Spec.API.DefinitionFileRef
		loggers.LoggerAPKOperator.Info("Sending API Creation event to agent")
		routes, err := CollectHTTPRoutes(ctx, routeMetadataReconciler.client, &routeMetadata)
		if err != nil {
			return ctrl.Result{}, err
		}

		// 3) classify + build states
		var prodState, sandState *syncronizer.HTTPRouteState

		for _, hr := range routes {
			env := strings.ToLower(hr.GetAnnotations()["gateway.envoyproxy.io/kgw-envtype"]) // "prod" or "sand"
			st, err := buildHTTPRouteStateFromObj(ctx, routeMetadataReconciler.client, hr)
			if err != nil {
				loggers.LoggerAPI.Error(err, "buildHTTPRouteStateFromObj failed", "route", hr.Name)
				return ctrl.Result{}, err
			}

			switch env {
			case "production":
				// If multiple prod routes exist, you can choose a policy (first wins, last wins, merge, etc.)
				if prodState == nil {
					prodState = st
				} else {
					// merge strategy if needed; for now keep the first
					loggers.LoggerAPI.Info("multiple prod HTTPRoutes detected; keeping the first", "existing", prodState.HTTPRouteCombined.Name, "ignored", hr.Name)
				}
			case "sandbox":
				if sandState == nil {
					sandState = st
				} else {
					loggers.LoggerAPI.Info("multiple sandbox HTTPRoutes detected; keeping the first", "existing", sandState.HTTPRouteCombined.Name, "ignored", hr.Name)
				}
			default:
				// If missing/unknown, default to prod unless you prefer to skip
				if prodState == nil {
					loggers.LoggerAPI.Info("HTTPRoute without gateway.envoyproxy.io/kgw-envtype (prod|sand) — defaulting to prod", "route", hr.Name)
					prodState = st
				} else {
					loggers.LoggerAPI.Info("extra HTTPRoute without env annotation ignored", "route", hr.Name)
				}
			}
			state.ProdHTTPRoute = prodState
			state.SandHTTPRoute = sandState
		}
		apiCpData := routeMetadataReconciler.convertAPIStateToAPICp(ctx, state)
		apiCpData.Event = controlplane.EventTypeCreate
		controlplane.AddToEventQueue(apiCpData)
	}

	routePolicyString, err := utils.ToJSONString(routeMetadata)
	if err != nil {
		loggers.LoggerAPKOperator.Errorf("Error converting RouteMetadata to JSON: %v", err)
	} else {
		utils.SendRouteMetadataCreatedOrUpdatedEvent(routePolicyString)
		loggers.LoggerAPKOperator.Debugf("Deleted RouteMetadata JSON: %s", routePolicyString)
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (routeMetadataReconciler *RouteMetadataReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dpV2alpha1.RouteMetadata{}).
		Complete(routeMetadataReconciler)
}

// getRouteMetadataForConfigMap returns a list of reconcile requests for RouteMetadata objects
func (routeMetadataReconciler *RouteMetadataReconciler) getRouteMetadataForConfigMap(ctx context.Context, obj *corev1.ConfigMap) []reconcile.Request {
	configMap := obj

	requests := []reconcile.Request{}

	routeMetadataList := &dpV2alpha1.RouteMetadataList{}
	if err := routeMetadataReconciler.client.List(ctx, routeMetadataList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(configMapIndex, utils.NamespacedName(configMap).String()),
	}); err != nil {
		return []reconcile.Request{}
	}

	for item := range routeMetadataList.Items {
		routePolicy := routeMetadataList.Items[item]
		requests = append(requests, routeMetadataReconciler.AddRouteMetadataRequest(&routePolicy)...)
	}

	return requests
}

// AddRouteMetadataRequest adds a reconcile request for the given RouteMetadata
func (routeMetadataReconciler *RouteMetadataReconciler) AddRouteMetadataRequest(routePolicy *dpV2alpha1.RouteMetadata) []reconcile.Request {
	routePolicyKey := client.ObjectKey{
		Namespace: routePolicy.Namespace,
		Name:      routePolicy.Name,
	}

	loggers.LoggerAPKOperator.Debugf("Adding RouteMetadata request for %s", routePolicyKey.String())
	return []reconcile.Request{{NamespacedName: routePolicyKey}}
}

func (routeMetadataReconciler *RouteMetadataReconciler) addRouteMetadataIndexes(ctx context.Context, mgr manager.Manager) error {
	// Index by referenced ConfigMaps
	if err := mgr.GetFieldIndexer().IndexField(ctx, &dpV2alpha1.RouteMetadata{}, configMapIndex,
		func(rawObj k8client.Object) []string {
			routeMetadata := rawObj.(*dpV2alpha1.RouteMetadata)
			namespacedConfigMaps := make([]string, 0)
			if routeMetadata.Spec.API.DefinitionFileRef != nil {
				namespacedConfigMaps = append(namespacedConfigMaps, types.NamespacedName{
					Namespace: routeMetadata.Namespace,
					Name:      string(routeMetadata.Spec.API.DefinitionFileRef.Name),
				}.String())
			}

			return namespacedConfigMaps
		}); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2614, logging.BLOCKER,
			"Error adding index for RouteMetadata: %v", err))
		return err
	}

	return nil
}

func (routeMetadataReconciler *RouteMetadataReconciler) convertAPIStateToAPICp(ctx context.Context, apiState synchronizer.APIState) controlplane.APICPEvent {
	apiCPEvent := controlplane.APICPEvent{}
	spec := apiState.APIDefinition.Spec.API
	configMap := &corev1.ConfigMap{}
	apiDef := ""
	if spec.DefinitionFileRef != nil {
		err := routeMetadataReconciler.client.Get(ctx, types.NamespacedName{Namespace: apiState.APIDefinition.Namespace, Name: string(apiState.APIDefinition.Spec.API.DefinitionFileRef.Name)}, configMap)
		if err != nil {
			loggers.LoggerAPKOperator.Errorf("Error while loading config map for the api definition: %+v, Error: %v", types.NamespacedName{Namespace: apiState.APIDefinition.Namespace, Name: string(apiState.APIDefinition.Spec.API.DefinitionFileRef.Name)}, err)
		} else {
			for _, val := range configMap.BinaryData {
				buf := bytes.NewReader(val)
				r, err := gzip.NewReader(buf)
				if err != nil {
					loggers.LoggerAPKOperator.Errorf("Error creating gzip reader. Error: %+v", err)
					continue
				}
				defer r.Close()
				decompressed, err := ioutil.ReadAll(r)
				if err != nil {
					loggers.LoggerAPKOperator.Errorf("Error reading decompressed data. Error: %+v", err)
					continue
				}
				apiDef = string(decompressed)
			}
		}
	}
	apiUUID, apiUUIDExists := apiState.APIDefinition.ObjectMeta.Labels["apiUUID"]
	if !apiUUIDExists {
		apiUUID = spec.Name
	}
	revisionID, revisionIDExists := apiState.APIDefinition.ObjectMeta.Labels["revisionID"]
	if !revisionIDExists {
		revisionID = "0"
	}
	loggers.LoggerAPI.Info("dwwdwd", apiDef)
	api := controlplane.API{
		APIName:      spec.Name,
		APIVersion:   spec.Version,
		APIType:      spec.Type,
		BasePath:     spec.Context,
		Organization: spec.Organization,
		Environment:  spec.Environment,
		Definition:   apiDef,
		APIUUID:      apiUUID,
		RevisionID:   revisionID,
		ProdEndpoint: SelectProdEndpointString(apiState.ProdHTTPRoute),
	}
	apiCPEvent.API = api
	fmt.Println(api)
	apiCPEvent.CRName = apiState.APIDefinition.ObjectMeta.Name
	apiCPEvent.CRNamespace = apiState.APIDefinition.ObjectMeta.Namespace
	return apiCPEvent

}

// CollectHTTPRoutes collects all unique HTTPRoute objects referenced by RouteMetadata annotations.
func CollectHTTPRoutes(ctx context.Context, c client.Client, rm *dpV2alpha1.RouteMetadata) ([]*gwapiv1.HTTPRoute, error) {
	seen := make(map[string]struct{})
	var routes []*gwapiv1.HTTPRoute

	for k, v := range rm.GetAnnotations() {
		if strings.HasPrefix(k, "dp.wso2.com/httproute_") {
			for _, name := range strings.Split(v, ",") {
				name = strings.TrimSpace(name)
				if name == "" {
					continue
				}
				if _, exists := seen[name]; exists {
					continue
				}
				seen[name] = struct{}{}

				// Fetch HTTPRoute from same namespace as RouteMetadata
				hr := &gwapiv1.HTTPRoute{}
				if err := c.Get(ctx, client.ObjectKey{Namespace: rm.Namespace, Name: name}, hr); err != nil {
					return nil, fmt.Errorf("failed to get HTTPRoute %s/%s: %w", rm.Namespace, name, err)
				}
				routes = append(routes, hr)
			}
		}
	}
	return routes, nil
}

// buildHTTPRouteStateFromObj initializes maps and resolves Backend refs into BackendMapping.
func buildHTTPRouteStateFromObj(ctx context.Context, kc client.Client, hr *gwapiv1.HTTPRoute) (*syncronizer.HTTPRouteState, error) {
	if hr == nil {
		return nil, fmt.Errorf("nil HTTPRoute")
	}

	st := &syncronizer.HTTPRouteState{
		HTTPRouteCombined:   hr.DeepCopy(),
		HTTPRoutePartitions: make(map[string]*gwapiv1.HTTPRoute),
		BackendMapping:      make(map[string]*v1alpha4.ResolvedBackend),
	}

	// Resolve Envoy Gateway Backend refs -> populate BackendMapping
	if err := resolveEGBackends(ctx, kc, st); err != nil {
		return nil, err
	}
	return st, nil
}

func resolveEGBackends(ctx context.Context, kc client.Client, st *syncronizer.HTTPRouteState) error {
	rt := st.HTTPRouteCombined
	routeNS := rt.Namespace

	seen := make(map[string]struct{}) // de‑dup per route (ns/name:port@host)

	for ruleIdx, rule := range rt.Spec.Rules {
		for _, br := range rule.BackendRefs {
			if br.Group != nil && string(*br.Group) != "gateway.envoyproxy.io" {
				continue
			}
			if br.Kind != nil && string(*br.Kind) != "Backend" {
				continue
			}
			if br.Name == "" {
				continue
			}

			beNS := routeNS
			if br.Namespace != nil && *br.Namespace != "" {
				beNS = string(*br.Namespace)
			}

			var be eg.Backend
			if err := kc.Get(ctx, types.NamespacedName{Namespace: beNS, Name: string(br.Name)}, &be); err != nil {
				return fmt.Errorf("get Envoy Backend %s/%s for route %s/%s: %w",
					beNS, br.Name, routeNS, rt.Name, err)
			}

			eps := extractFQDNs(&be)
			for _, ep := range eps {
				key := beNS + "/" + string(br.Name) + ":" + strconv.Itoa(int(ep.Port)) + "@" + ep.Host
				if _, ok := seen[key]; ok {
					continue
				}
				seen[key] = struct{}{}

				// TODO: Fill your actual fields in ResolvedBackend (Host, Port, TLS, etc.)
				st.BackendMapping[key] = &v1alpha4.ResolvedBackend{}
				// Example if your struct has these fields:
				// st.BackendMapping[key] = &v1alpha4.ResolvedBackend{Host: ep.Host, Port: int(ep.Port)}
			}

			_ = ruleIdx // keep if you later map per-rule policies
		}
	}
	return nil
}

type fqdnEP struct {
	Host string
	Port int32
}

func extractFQDNs(be *eg.Backend) []fqdnEP {
	var out []fqdnEP
	if be == nil {
		return out
	}
	for _, ep := range be.Spec.Endpoints {
		if ep.FQDN != nil && ep.FQDN.Hostname != "" && ep.FQDN.Port != 0 {
			out = append(out, fqdnEP{Host: strings.TrimSpace(ep.FQDN.Hostname), Port: ep.FQDN.Port})
		}
		// (Optionally handle IP/Unix endpoints here)
	}
	return out
}

// SelectProdEndpointString returns the first backend endpoint in the HTTPRouteState's BackendMapping
func SelectProdEndpointString(st *synchronizer.HTTPRouteState) string {
	if st == nil || len(st.BackendMapping) == 0 {
		return ""
	}

	// pick the first key deterministically
	for key := range st.BackendMapping {
		// key format: ns/name:port@host
		at := strings.LastIndexByte(key, '@')
		if at < 0 {
			continue
		}
		left, host := key[:at], key[at+1:]

		col := strings.LastIndexByte(left, ':')
		if col < 0 {
			continue
		}
		portStr := left[col+1:]

		return host + ":" + portStr
	}
	return ""
}
