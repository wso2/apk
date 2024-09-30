/*
 *  Copyright (c) 2024, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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
 *  This file contains code derived from Envoy Gateway,
 *  https://github.com/envoyproxy/gateway
 *  and is provided here subject to the following:
 *  Copyright Project Envoy Gateway Authors
 *
 */

package gatewayapi

import (
	"reflect"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
	"sigs.k8s.io/gateway-api/apis/v1alpha2"
)

// GatewayContext wraps a Gateway and provides helper methods for setting conditions, accessing Listeners, etc.
type GatewayContext struct {
	*gwapiv1.Gateway
	listeners []*ListenerContext
}

// ResetListeners resets the listener statuses and re-generates the GatewayContext
// ListenerContexts from the Gateway spec.
func (g *GatewayContext) ResetListeners() {
	lenListener := len(g.Spec.Listeners)
	g.Status.Listeners = make([]gwapiv1.ListenerStatus, lenListener)
	g.listeners = make([]*ListenerContext, lenListener)
	for i := range g.Spec.Listeners {
		listener := &g.Spec.Listeners[i]
		g.Status.Listeners[i] = gwapiv1.ListenerStatus{Name: listener.Name}
		g.listeners[i] = &ListenerContext{
			listenerStatusID: i,
			Listener:         listener,
			gateway:          g.Gateway,
		}
	}
}

func (g *GatewayContext) attachEnvoyProxy(resources *Resources) {
	// if g.Spec.Infrastructure != nil && g.Spec.Infrastructure.ParametersRef != nil && !IsMergeGatewaysEnabled(resources) {
	// 	ref := g.Spec.Infrastructure.ParametersRef
	// 	if string(ref.Group) == egv1a1.GroupVersion.Group && ref.Kind == egv1a1.KindEnvoyProxy {
	// 		ep := resources.GetEnvoyProxy(g.Namespace, ref.Name)
	// 		if ep != nil {
	// 			g.envoyProxy = ep
	// 			return
	// 		}
	// 	}
	// 	// not found, fallthrough to use envoyProxy attached to gatewayclass
	// }

	// g.envoyProxy = resources.EnvoyProxyForGatewayClass
}

// ListenerContext wraps a Listener and provides helper methods for
// setting conditions and other status information on the associated
// Gateway, etc.
type ListenerContext struct {
	*gwapiv1.Listener

	gateway           *gwapiv1.Gateway
	listenerStatusID  int
	namespaceSelector labels.Selector
	tlsSecrets        []*v1.Secret
}

func (l *ListenerContext) SetCondition(conditionType gwapiv1.ListenerConditionType, status metav1.ConditionStatus, reason gwapiv1.ListenerConditionReason, message string) {
	cond := metav1.Condition{
		Type:               string(conditionType),
		Status:             status,
		Reason:             string(reason),
		Message:            message,
		ObservedGeneration: l.gateway.Generation,
		LastTransitionTime: metav1.NewTime(time.Now()),
	}

	index := -1
	for i, existing := range l.gateway.Status.Listeners[l.listenerStatusID].Conditions {
		if existing.Type == cond.Type {
			// return early if the condition is unchanged
			if existing.Status == cond.Status &&
				existing.Reason == cond.Reason &&
				existing.Message == cond.Message &&
				existing.ObservedGeneration == cond.ObservedGeneration {
				return
			}
			index = i
			break
		}
	}

	if index > -1 {
		l.gateway.Status.Listeners[l.listenerStatusID].Conditions[index] = cond
	} else {
		l.gateway.Status.Listeners[l.listenerStatusID].Conditions = append(l.gateway.Status.Listeners[l.listenerStatusID].Conditions, cond)
	}
}

func (l *ListenerContext) SetSupportedKinds(kinds ...gwapiv1.RouteGroupKind) {
	l.gateway.Status.Listeners[l.listenerStatusID].SupportedKinds = kinds
}

func (l *ListenerContext) IncrementAttachedRoutes() {
	l.gateway.Status.Listeners[l.listenerStatusID].AttachedRoutes++
}

func (l *ListenerContext) AttachedRoutes() int32 {
	return l.gateway.Status.Listeners[l.listenerStatusID].AttachedRoutes
}

func (l *ListenerContext) AllowsKind(kind gwapiv1.RouteGroupKind) bool {
	for _, allowed := range l.gateway.Status.Listeners[l.listenerStatusID].SupportedKinds {
		if GroupDerefOr(allowed.Group, "") == GroupDerefOr(kind.Group, "") &&
			allowed.Kind == kind.Kind {
			return true
		}
	}

	return false
}

func (l *ListenerContext) AllowsNamespace(namespace *v1.Namespace) bool {
	if namespace == nil {
		return false
	}

	if l.AllowedRoutes == nil || l.AllowedRoutes.Namespaces == nil || l.AllowedRoutes.Namespaces.From == nil {
		return l.gateway.Namespace == namespace.Name
	}

	switch *l.AllowedRoutes.Namespaces.From {
	case gwapiv1.NamespacesFromAll:
		return true
	case gwapiv1.NamespacesFromSelector:
		if l.namespaceSelector == nil {
			return false
		}
		return l.namespaceSelector.Matches(labels.Set(namespace.Labels))
	default:
		// NamespacesFromSame is the default
		return l.gateway.Namespace == namespace.Name
	}
}

func (l *ListenerContext) IsReady() bool {
	for _, cond := range l.gateway.Status.Listeners[l.listenerStatusID].Conditions {
		if cond.Type == string(gwapiv1.ListenerConditionProgrammed) && cond.Status == metav1.ConditionTrue {
			return true
		}
	}

	return false
}

func (l *ListenerContext) GetConditions() []metav1.Condition {
	return l.gateway.Status.Listeners[l.listenerStatusID].Conditions
}

func (l *ListenerContext) SetTLSSecrets(tlsSecrets []*v1.Secret) {
	l.tlsSecrets = tlsSecrets
}

// RouteContext represents a generic Route object (HTTPRoute, TLSRoute, etc.)
// that can reference Gateway objects.
type RouteContext interface {
	client.Object
}

// HTTPRouteContext wraps an HTTPRoute and provides helper methods for
// accessing the route's parents.
type HTTPRouteContext struct {
	// GatewayControllerName is the name of the Gateway API controller.
	GatewayControllerName string

	*gwapiv1.HTTPRoute

	ParentRefs map[gwapiv1.ParentReference]*RouteParentContext
}

// GRPCRouteContext wraps a GRPCRoute and provides helper methods for
// accessing the route's parents.
type GRPCRouteContext struct {
	// GatewayControllerName is the name of the Gateway API controller.
	GatewayControllerName string

	*v1alpha2.GRPCRoute

	ParentRefs map[gwapiv1.ParentReference]*RouteParentContext
}

// TLSRouteContext wraps a TLSRoute and provides helper methods for
// accessing the route's parents.
type TLSRouteContext struct {
	// GatewayControllerName is the name of the Gateway API controller.
	GatewayControllerName string

	*v1alpha2.TLSRoute

	ParentRefs map[gwapiv1.ParentReference]*RouteParentContext
}

// UDPRouteContext wraps a UDPRoute and provides helper methods for
// accessing the route's parents.
type UDPRouteContext struct {
	// GatewayControllerName is the name of the Gateway API controller.
	GatewayControllerName string

	*v1alpha2.UDPRoute

	ParentRefs map[gwapiv1.ParentReference]*RouteParentContext
}

// TCPRouteContext wraps a TCPRoute and provides helper methods for
// accessing the route's parents.
type TCPRouteContext struct {
	// GatewayControllerName is the name of the Gateway API controller.
	GatewayControllerName string

	*v1alpha2.TCPRoute

	ParentRefs map[gwapiv1.ParentReference]*RouteParentContext
}

// GetRouteType returns the Kind of the Route object, HTTPRoute,
// TLSRoute, TCPRoute, UDPRoute etc.
func GetRouteType(route RouteContext) gwapiv1.Kind {
	rv := reflect.ValueOf(route).Elem()
	return gwapiv1.Kind(rv.FieldByName("Kind").String())
}

// TODO: [v1alpha2-gwapiv1] This should not be required once all Route
// objects being implemented are of type gwapiv1.
// GetHostnames returns the hosts targeted by the Route object.
func GetHostnames(route RouteContext) []string {
	rv := reflect.ValueOf(route).Elem()
	kind := rv.FieldByName("Kind").String()
	if kind == KindTCPRoute || kind == KindUDPRoute {
		return nil
	}

	hs := rv.FieldByName("Spec").FieldByName("Hostnames")
	hostnames := make([]string, hs.Len())
	for i := 0; i < len(hostnames); i++ {
		hostnames[i] = hs.Index(i).String()
	}
	return hostnames
}

// TODO: [v1alpha2-gwapiv1] This should not be required once all Route
// objects being implemented are of type gwapiv1.
// GetParentReferences returns the ParentReference of the Route object.
func GetParentReferences(route RouteContext) []gwapiv1.ParentReference {
	rv := reflect.ValueOf(route).Elem()
	kind := rv.FieldByName("Kind").String()
	pr := rv.FieldByName("Spec").FieldByName("ParentRefs")
	if kind == KindHTTPRoute || kind == KindGRPCRoute {
		return pr.Interface().([]gwapiv1.ParentReference)
	}

	parentReferences := make([]gwapiv1.ParentReference, pr.Len())
	for i := 0; i < len(parentReferences); i++ {
		p := pr.Index(i).Interface().(gwapiv1.ParentReference)
		parentReferences[i] = UpgradeParentReference(p)
	}
	return parentReferences
}

// GetRouteStatus returns the RouteStatus object associated with the Route.
func GetRouteStatus(route RouteContext) *gwapiv1.RouteStatus {
	rv := reflect.ValueOf(route).Elem()
	rs := rv.FieldByName("Status").FieldByName("RouteStatus").Interface().(gwapiv1.RouteStatus)
	return &rs
}

// GetRouteParentContext returns RouteParentContext by using the Route
// objects' ParentReference.
func GetRouteParentContext(route RouteContext, forParentRef gwapiv1.ParentReference) *RouteParentContext {
	rv := reflect.ValueOf(route).Elem()
	pr := rv.FieldByName("ParentRefs")
	if pr.IsNil() {
		mm := reflect.MakeMap(reflect.TypeOf(map[gwapiv1.ParentReference]*RouteParentContext{}))
		pr.Set(mm)
	}

	if p := pr.MapIndex(reflect.ValueOf(forParentRef)); p.IsValid() && !p.IsZero() {
		ctx := p.Interface().(*RouteParentContext)
		return ctx
	}

	isHTTPRoute := false
	if rv.FieldByName("Kind").String() == KindHTTPRoute {
		isHTTPRoute = true
	}

	var parentRef *gwapiv1.ParentReference
	specParentRefs := rv.FieldByName("Spec").FieldByName("ParentRefs")
	for i := 0; i < specParentRefs.Len(); i++ {
		p := specParentRefs.Index(i).Interface().(gwapiv1.ParentReference)
		up := p
		if !isHTTPRoute {
			up = UpgradeParentReference(p)
		}
		if reflect.DeepEqual(up, forParentRef) {
			if isHTTPRoute {
				parentRef = &p
			} else {
				upgraded := UpgradeParentReference(p)
				parentRef = &upgraded
			}
			break
		}
	}
	if parentRef == nil {
		panic("parentRef not found")
	}

	routeParentStatusIdx := -1
	statusParents := rv.FieldByName("Status").FieldByName("Parents")
	for i := 0; i < statusParents.Len(); i++ {
		p := statusParents.Index(i).FieldByName("ParentRef").Interface().(gwapiv1.ParentReference)
		if !isHTTPRoute {
			p = UpgradeParentReference(p)
			defaultNamespace := gwapiv1.Namespace(metav1.NamespaceDefault)
			if forParentRef.Namespace == nil {
				forParentRef.Namespace = &defaultNamespace
			}
			if p.Namespace == nil {
				p.Namespace = &defaultNamespace
			}
		}
		if reflect.DeepEqual(p, forParentRef) {
			routeParentStatusIdx = i
			break
		}
	}
	if routeParentStatusIdx == -1 {
		tmpPR := forParentRef
		if !isHTTPRoute {
			tmpPR = DowngradeParentReference(tmpPR)
		}
		rParentStatus := v1alpha2.RouteParentStatus{
			ControllerName: v1alpha2.GatewayController(rv.FieldByName("GatewayControllerName").String()),
			ParentRef:      tmpPR,
		}
		statusParents.Set(reflect.Append(statusParents, reflect.ValueOf(rParentStatus)))
		routeParentStatusIdx = statusParents.Len() - 1
	}

	ctx := &RouteParentContext{
		ParentReference:      parentRef,
		routeParentStatusIdx: routeParentStatusIdx,
	}
	rctx := reflect.ValueOf(ctx)
	rctx.Elem().FieldByName(string(GetRouteType(route))).Set(rv.Field(1))
	pr.SetMapIndex(reflect.ValueOf(forParentRef), rctx)
	return ctx
}

// RouteParentContext wraps a ParentReference and provides helper methods for
// setting conditions and other status information on the associated
// HTTPRoute, TLSRoute etc.
type RouteParentContext struct {
	*gwapiv1.ParentReference

	// TODO: [v1alpha2-gwapiv1] This can probably be replaced with
	// a single field pointing to *gwapiv1.RouteStatus.
	HTTPRoute *gwapiv1.HTTPRoute
	GRPCRoute *v1alpha2.GRPCRoute
	TLSRoute  *v1alpha2.TLSRoute
	TCPRoute  *v1alpha2.TCPRoute
	UDPRoute  *v1alpha2.UDPRoute

	routeParentStatusIdx int
	listeners            []*ListenerContext
}

// GetGateway returns the GatewayContext if parent resource is a gateway.
func (r *RouteParentContext) GetGateway() *GatewayContext {
	if r == nil || len(r.listeners) == 0 {
		return nil
	}
	gwc := &GatewayContext{
		Gateway: r.listeners[0].gateway,
	}
	return gwc
}

func (r *RouteParentContext) SetListeners(listeners ...*ListenerContext) {
	r.listeners = append(r.listeners, listeners...)
}

func (r *RouteParentContext) SetCondition(route RouteContext, conditionType gwapiv1.RouteConditionType, status metav1.ConditionStatus, reason gwapiv1.RouteConditionReason, message string) {
	cond := metav1.Condition{
		Type:               string(conditionType),
		Status:             status,
		Reason:             string(reason),
		Message:            message,
		ObservedGeneration: route.GetGeneration(),
		LastTransitionTime: metav1.NewTime(time.Now()),
	}

	idx := -1
	routeStatus := GetRouteStatus(route)
	for i, existing := range routeStatus.Parents[r.routeParentStatusIdx].Conditions {
		if existing.Type == cond.Type {
			// return early if the condition is unchanged
			if existing.Status == cond.Status &&
				existing.Reason == cond.Reason &&
				existing.Message == cond.Message &&
				existing.ObservedGeneration == cond.ObservedGeneration {
				return
			}
			idx = i
			break
		}
	}

	if idx > -1 {
		routeStatus.Parents[r.routeParentStatusIdx].Conditions[idx] = cond
	} else {
		routeStatus.Parents[r.routeParentStatusIdx].Conditions = append(routeStatus.Parents[r.routeParentStatusIdx].Conditions, cond)
	}
}

func (r *RouteParentContext) ResetConditions(route RouteContext) {
	routeStatus := GetRouteStatus(route)
	routeStatus.Parents[r.routeParentStatusIdx].Conditions = make([]metav1.Condition, 0)
}

func (r *RouteParentContext) HasCondition(route RouteContext, condType gwapiv1.RouteConditionType, status metav1.ConditionStatus) bool {
	var conditions []metav1.Condition
	routeStatus := GetRouteStatus(route)
	conditions = routeStatus.Parents[r.routeParentStatusIdx].Conditions
	for _, c := range conditions {
		if c.Type == string(condType) && c.Status == status {
			return true
		}
	}
	return false
}

// BackendRefContext represents a generic BackendRef object (HTTPBackendRef, GRPCBackendRef or BackendRef itself)
type BackendRefContext any

func GetBackendRef(b BackendRefContext) *gwapiv1.BackendRef {
	rv := reflect.ValueOf(b)
	br := rv.FieldByName("BackendRef")
	if br.IsValid() {
		backendRef := br.Interface().(gwapiv1.BackendRef)
		return &backendRef

	}
	backendRef := b.(gwapiv1.BackendRef)
	return &backendRef
}

func GetFilters(b BackendRefContext) any {
	return reflect.ValueOf(b).FieldByName("Filters").Interface()
}
