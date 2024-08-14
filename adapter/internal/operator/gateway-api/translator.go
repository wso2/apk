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
	"github.com/wso2/apk/adapter/internal/operator/gateway-api/ir"
	"golang.org/x/exp/maps"
	"k8s.io/apimachinery/pkg/runtime/schema"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
)

const (
	KindConfigMap        = "ConfigMap"
	KindBackendTLSPolicy = "BackendTLSPolicy"
	KindGateway          = "Gateway"
	KindGatewayClass     = "GatewayClass"
	KindGRPCRoute        = "GRPCRoute"
	KindHTTPRoute        = "HTTPRoute"
	KindNamespace        = "Namespace"
	KindTLSRoute         = "TLSRoute"
	KindTCPRoute         = "TCPRoute"
	KindUDPRoute         = "UDPRoute"
	KindService          = "Service"
	KindServiceImport    = "ServiceImport"
	KindSecret           = "Secret"
	KindBackend          = "Backend"

	GroupMultiClusterService = "multicluster.x-k8s.io"
	// OwningGatewayNamespaceLabel is the owner reference label used for managed infra.
	// The value should be the namespace of the accepted Envoy Gateway.
	OwningGatewayNamespaceLabel = "apk.wso2.com/owning-gateway-namespace"

	OwningGatewayClassLabel = "apk.wso2.com/owning-gatewayclass"
	// OwningGatewayNameLabel is the owner reference label used for managed infra.
	// The value should be the name of the accepted Envoy Gateway.
	OwningGatewayNameLabel = "apk.wso2.com/owning-gateway-name"

	// minEphemeralPort is the first port in the ephemeral port range.
	minEphemeralPort = 1024
	// wellKnownPortShift is the constant added to the well known port (1-1023)
	// to convert it into an ephemeral port.
	wellKnownPortShift = 10000
)

var _ TranslatorManager = (*Translator)(nil)

type TranslatorManager interface {
	Translate(resources *Resources) *TranslateResult
	GetRelevantGateways(gateways []*gwapiv1.Gateway) []*GatewayContext

	RoutesTranslator
	ListenersTranslator
	AddressesTranslator
	FiltersTranslator
}

// Translator translates Gateway API resources to IRs and computes status
// for Gateway API resources.
type Translator struct {
	// GatewayControllerName is the name of the Gateway API controller
	GatewayControllerName string

	// GatewayClassName is the name of the GatewayClass
	// to process Gateways for.
	GatewayClassName gwapiv1.ObjectName

	// GlobalRateLimitEnabled is true when global
	// ratelimiting has been configured by the admin.
	GlobalRateLimitEnabled bool

	// EndpointRoutingDisabled can be set to true to use
	// the Service Cluster IP for routing to the backend
	// instead.
	EndpointRoutingDisabled bool

	// MergeGateways is true when all Gateway Listeners
	// should be merged under the parent GatewayClass.
	MergeGateways bool

	// EnvoyPatchPolicyEnabled when the EnvoyPatchPolicy
	// feature is enabled.
	EnvoyPatchPolicyEnabled bool

	// ExtensionGroupKinds stores the group/kind for all resources
	// introduced by an Extension so that the translator can
	// store referenced resources in the IR for later use.
	ExtensionGroupKinds []schema.GroupKind

	// Namespace is the namespace that Envoy Gateway runs in.
	Namespace string
}

type TranslateResult struct {
	Resources
	XdsIR   XdsIRMap   `json:"xdsIR" yaml:"xdsIR"`
	InfraIR InfraIRMap `json:"infraIR" yaml:"infraIR"`
}

func newTranslateResult(gateways []*GatewayContext,
	httpRoutes []*HTTPRouteContext,
	grpcRoutes []*GRPCRouteContext,
	tlsRoutes []*TLSRouteContext,
	tcpRoutes []*TCPRouteContext,
	udpRoutes []*UDPRouteContext,
	xdsIR XdsIRMap, infraIR InfraIRMap) *TranslateResult {
	translateResult := &TranslateResult{
		XdsIR:   xdsIR,
		InfraIR: infraIR,
	}

	for _, gateway := range gateways {
		translateResult.Gateways = append(translateResult.Gateways, gateway.Gateway)
	}
	for _, httpRoute := range httpRoutes {
		translateResult.HTTPRoutes = append(translateResult.HTTPRoutes, httpRoute.HTTPRoute)
	}
	for _, grpcRoute := range grpcRoutes {
		translateResult.GRPCRoutes = append(translateResult.GRPCRoutes, grpcRoute.GRPCRoute)
	}
	for _, tlsRoute := range tlsRoutes {
		translateResult.TLSRoutes = append(translateResult.TLSRoutes, tlsRoute.TLSRoute)
	}
	for _, tcpRoute := range tcpRoutes {
		translateResult.TCPRoutes = append(translateResult.TCPRoutes, tcpRoute.TCPRoute)
	}
	for _, udpRoute := range udpRoutes {
		translateResult.UDPRoutes = append(translateResult.UDPRoutes, udpRoute.UDPRoute)
	}
	return translateResult
}

func (t *Translator) Translate(resources *Resources) *TranslateResult {
	// Get Gateways belonging to our GatewayClass.
	gateways := t.GetRelevantGateways(resources.Gateways)

	// Build IR maps.
	xdsIR, infraIR := t.InitIRs(gateways, resources)

	// Process all Listeners for all relevant Gateways.
	t.ProcessListeners(gateways, xdsIR, infraIR, resources)

	// Process all Addresses for all relevant Gateways.
	t.ProcessAddresses(gateways, infraIR)

	// Process all relevant HTTPRoutes.
	httpRoutes := t.ProcessHTTPRoutes(resources.HTTPRoutes, gateways, resources, xdsIR)

	// Process all relevant GRPCRoutes.
	// grpcRoutes := t.ProcessGRPCRoutes(resources.GRPCRoutes, gateways, resources, xdsIR)

	// // Process all relevant TLSRoutes.
	// tlsRoutes := t.ProcessTLSRoutes(resources.TLSRoutes, gateways, resources, xdsIR)

	// // Process all relevant TCPRoutes.
	// tcpRoutes := t.ProcessTCPRoutes(resources.TCPRoutes, gateways, resources, xdsIR)

	// // Process all relevant UDPRoutes.
	// udpRoutes := t.ProcessUDPRoutes(resources.UDPRoutes, gateways, resources, xdsIR)

	routes := []RouteContext{}
	for _, h := range httpRoutes {
		routes = append(routes, h)
	}
	// for _, g := range grpcRoutes {
	// 	routes = append(routes, g)
	// }
	// for _, t := range tlsRoutes {
	// 	routes = append(routes, t)
	// }
	// for _, t := range tcpRoutes {
	// 	routes = append(routes, t)
	// }
	// for _, u := range udpRoutes {
	// 	routes = append(routes, u)
	// }

	// Process BackendTrafficPolicies
	// backendTrafficPolicies := t.ProcessBackendTrafficPolicies(
	// 	resources.BackendTrafficPolicies, gateways, routes, xdsIR)

	// // Process SecurityPolicies
	// securityPolicies := t.ProcessSecurityPolicies(
	// 	resources.SecurityPolicies, gateways, routes, resources, xdsIR)

	t.ProcessAPIs(resources.APIs, resources.HTTPRoutes, gateways, httpRoutes, resources, xdsIR)

	// Sort xdsIR based on the Gateway API spec
	sortXdsIRMap(xdsIR)

	return newTranslateResult(gateways, httpRoutes, nil, nil, nil, nil, xdsIR, infraIR)

}

// GetRelevantGateways returns GatewayContexts, containing a copy of the original
// Gateway with the Listener statuses reset.
func (t *Translator) GetRelevantGateways(gateways []*gwapiv1.Gateway) []*GatewayContext {
	var relevant []*GatewayContext

	for _, gateway := range gateways {
		if gateway == nil {
			panic("received nil gateway")
		}

		if gateway.Spec.GatewayClassName == t.GatewayClassName {
			gc := &GatewayContext{
				Gateway: gateway.DeepCopy(),
			}
			gc.ResetListeners()

			relevant = append(relevant, gc)
		}
	}

	return relevant
}

// InitIRs checks if mergeGateways is enabled in EnvoyProxy config and initializes XdsIR and InfraIR maps with adequate keys.
func (t *Translator) InitIRs(gateways []*GatewayContext, resources *Resources) (map[string]*ir.Xds, map[string]*ir.Infra) {
	xdsIR := make(XdsIRMap)
	infraIR := make(InfraIRMap)

	var irKey string
	for _, gateway := range gateways {
		gwXdsIR := &ir.Xds{}
		gwInfraIR := ir.NewInfra()
		labels := infrastructureLabels(gateway.Gateway)
		annotations := infrastructureAnnotations(gateway.Gateway)
		gwInfraIR.Proxy.GetProxyMetadata().Annotations = annotations

		irKey = irStringKey(gateway.Gateway.Namespace, gateway.Gateway.Name)

		maps.Copy(labels, GatewayOwnerLabels(gateway.Namespace, gateway.Name))
		gwInfraIR.Proxy.GetProxyMetadata().Labels = labels

		gwInfraIR.Proxy.Name = irKey
		// save the IR references in the map before the translation starts
		xdsIR[irKey] = gwXdsIR
		infraIR[irKey] = gwInfraIR
	}

	return xdsIR, infraIR
}

func infrastructureAnnotations(gtw *gwapiv1.Gateway) map[string]string {
	if gtw.Spec.Infrastructure != nil && len(gtw.Spec.Infrastructure.Annotations) > 0 {
		res := make(map[string]string)
		for k, v := range gtw.Spec.Infrastructure.Annotations {
			res[string(k)] = string(v)
		}
		return res
	}
	return nil
}

func infrastructureLabels(gtw *gwapiv1.Gateway) map[string]string {
	res := make(map[string]string)
	if gtw.Spec.Infrastructure != nil {
		for k, v := range gtw.Spec.Infrastructure.Labels {
			res[string(k)] = string(v)
		}
	}
	return res
}

// XdsIR and InfraIR map keys by default are {GatewayNamespace}/{GatewayName}, but if mergeGateways is set, they are merged under {GatewayClassName} key.
func (t *Translator) getIRKey(gateway *gwapiv1.Gateway) string {
	irKey := irStringKey(gateway.Namespace, gateway.Name)
	if t.MergeGateways {
		return string(t.GatewayClassName)
	}

	return irKey
}
