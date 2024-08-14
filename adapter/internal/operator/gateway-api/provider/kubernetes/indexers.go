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

package kubernetes

import (
	"context"

	"github.com/wso2/apk/adapter/internal/loggers"
	gatewayapi "github.com/wso2/apk/adapter/internal/operator/gateway-api"
	dpv1alpha2 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha2"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
	gwapiv1a2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gwapiv1b1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

const (
	classGatewayIndex          = "classGatewayIndex"
	gatewayTLSRouteIndex       = "gatewayTLSRouteIndex"
	gatewayHTTPRouteIndex      = "gatewayHTTPRouteIndex"
	httpRouteAPIIndex          = "httpRouteAPIIndex"
	gatewayGRPCRouteIndex      = "gatewayGRPCRouteIndex"
	gatewayTCPRouteIndex       = "gatewayTCPRouteIndex"
	gatewayUDPRouteIndex       = "gatewayUDPRouteIndex"
	secretGatewayIndex         = "secretGatewayIndex"
	targetRefGrantRouteIndex   = "targetRefGrantRouteIndex"
	backendHTTPRouteIndex      = "backendHTTPRouteIndex"
	serviceHTTPRouteIndex      = "serviceHTTPRouteIndex"
	backendGRPCRouteIndex      = "backendGRPCRouteIndex"
	backendTLSRouteIndex       = "backendTLSRouteIndex"
	backendTCPRouteIndex       = "backendTCPRouteIndex"
	backendUDPRouteIndex       = "backendUDPRouteIndex"
	secretSecurityPolicyIndex  = "secretSecurityPolicyIndex"
	backendSecurityPolicyIndex = "backendSecurityPolicyIndex"
	configMapCtpIndex          = "configMapCtpIndex"
	secretCtpIndex             = "secretCtpIndex"
	configMapBtlsIndex         = "configMapBtlsIndex"
)

func addReferenceGrantIndexers(ctx context.Context, mgr manager.Manager) error {
	return mgr.GetFieldIndexer().IndexField(ctx, &gwapiv1b1.ReferenceGrant{}, targetRefGrantRouteIndex, func(rawObj client.Object) []string {
		refGrant := rawObj.(*gwapiv1b1.ReferenceGrant)
		var referredServices []string
		for _, target := range refGrant.Spec.To {
			referredServices = append(referredServices, string(target.Kind))
		}
		return referredServices
	})
}

// addHTTPRouteIndexers adds indexing on HTTPRoute.
//   - For Service, ServiceImports objects that are referenced in HTTPRoute objects via `.spec.rules.backendRefs`.
//     This helps in querying for HTTPRoutes that are affected by a particular Service CRUD.
func addHTTPRouteIndexers(ctx context.Context, mgr manager.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(ctx, &gwapiv1.HTTPRoute{}, gatewayHTTPRouteIndex, gatewayHTTPRouteIndexFunc); err != nil {
		return err
	}
	// Backend to HTTPRoute indexer
	if err := mgr.GetFieldIndexer().IndexField(ctx, &gwapiv1.HTTPRoute{}, backendHTTPRouteIndex, backendHTTPRouteIndexFunc); err != nil {
		return err
	}

	return mgr.GetFieldIndexer().IndexField(ctx, &gwapiv1.HTTPRoute{}, serviceHTTPRouteIndex, serviceHTTPRouteIndexFunc)
}

// addAPIIndexers adds indexing on API.
//   - For Service, ServiceImports objects that are referenced in HTTPRoute objects via `.spec.rules.backendRefs`.
//     This helps in querying for HTTPRoutes that are affected by a particular Service CRUD.
func addAPIIndexers(ctx context.Context, mgr manager.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(ctx, &dpv1alpha2.API{}, httpRouteAPIIndex, httpRouteAPIIndexFunc); err != nil {
		return err
	}
	return nil
}

func gatewayHTTPRouteIndexFunc(rawObj client.Object) []string {
	httproute := rawObj.(*gwapiv1.HTTPRoute)
	var gateways []string
	for _, parent := range httproute.Spec.ParentRefs {
		if parent.Kind == nil || string(*parent.Kind) == gatewayapi.KindGateway {
			// If an explicit Gateway namespace is not provided, use the HTTPRoute namespace to
			// lookup the provided Gateway Name.
			gateways = append(gateways,
				types.NamespacedName{
					Namespace: gatewayapi.NamespaceDerefOr(parent.Namespace, httproute.Namespace),
					Name:      string(parent.Name),
				}.String(),
			)
		}
	}
	return gateways
}

func httpRouteAPIIndexFunc(rawObj client.Object) []string {
	api := rawObj.(*dpv1alpha2.API)
	var httpRoutes []string
	if len(api.Spec.Production) > 0 {
		for _, ref := range api.Spec.Production[0].RouteRefs {
			if ref != "" {
				httpRoutes = append(httpRoutes,
					types.NamespacedName{
						Namespace: api.Namespace,
						Name:      ref,
					}.String())
			}
		}
	}
	if len(api.Spec.Sandbox) > 0 {
		for _, ref := range api.Spec.Sandbox[0].RouteRefs {
			if ref != "" {
				httpRoutes = append(httpRoutes,
					types.NamespacedName{
						Namespace: api.Namespace,
						Name:      ref,
					}.String())
			}
		}
	}
	return httpRoutes
}

func serviceHTTPRouteIndexFunc(rawObj client.Object) []string {
	httproute := rawObj.(*gwapiv1.HTTPRoute)
	var backendRefs []string
	for _, rule := range httproute.Spec.Rules {
		for _, backend := range rule.BackendRefs {
			if backend.Kind == nil || string(*backend.Kind) == gatewayapi.KindService {
				// If an explicit Backend namespace is not provided, use the HTTPRoute namespace to
				// lookup the provided Gateway Name.
				backendRefs = append(backendRefs,
					types.NamespacedName{
						Namespace: gatewayapi.NamespaceDerefOr(backend.Namespace, httproute.Namespace),
						Name:      string(backend.Name),
					}.String(),
				)
			}
		}
	}
	return backendRefs
}

func backendHTTPRouteIndexFunc(rawObj client.Object) []string {
	httproute := rawObj.(*gwapiv1.HTTPRoute)
	var backendRefs []string
	for _, rule := range httproute.Spec.Rules {
		for _, backend := range rule.BackendRefs {
			if backend.Kind == nil || string(*backend.Kind) == gatewayapi.KindBackend {
				// If an explicit Backend namespace is not provided, use the HTTPRoute namespace to
				// lookup the provided Gateway Name.
				backendRefs = append(backendRefs,
					types.NamespacedName{
						Namespace: gatewayapi.NamespaceDerefOr(backend.Namespace, httproute.Namespace),
						Name:      string(backend.Name),
					}.String(),
				)
			}
		}
	}
	return backendRefs
}

// addGRPCRouteIndexers adds indexing on GRPCRoute, for Service objects that are
// referenced in GRPCRoute objects via `.spec.rules.backendRefs`. This helps in
// querying for GRPCRoutes that are affected by a particular Service CRUD.
func addGRPCRouteIndexers(ctx context.Context, mgr manager.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(ctx, &gwapiv1a2.GRPCRoute{}, gatewayGRPCRouteIndex, gatewayGRPCRouteIndexFunc); err != nil {
		return err
	}

	return mgr.GetFieldIndexer().IndexField(ctx, &gwapiv1a2.GRPCRoute{}, backendGRPCRouteIndex, backendGRPCRouteIndexFunc)
}

func gatewayGRPCRouteIndexFunc(rawObj client.Object) []string {
	grpcroute := rawObj.(*gwapiv1a2.GRPCRoute)
	var gateways []string
	for _, parent := range grpcroute.Spec.ParentRefs {
		if parent.Kind == nil || string(*parent.Kind) == gatewayapi.KindGateway {
			// If an explicit Gateway namespace is not provided, use the GRPCRoute namespace to
			// lookup the provided Gateway Name.
			gateways = append(gateways,
				types.NamespacedName{
					Namespace: gatewayapi.NamespaceDerefOr(parent.Namespace, grpcroute.Namespace),
					Name:      string(parent.Name),
				}.String(),
			)
		}
	}
	return gateways
}

func backendGRPCRouteIndexFunc(rawObj client.Object) []string {
	grpcroute := rawObj.(*gwapiv1a2.GRPCRoute)
	var backendRefs []string
	for _, rule := range grpcroute.Spec.Rules {
		for _, backend := range rule.BackendRefs {
			if backend.Kind == nil || string(*backend.Kind) == gatewayapi.KindService {
				// If an explicit Backend namespace is not provided, use the GRPCRoute namespace to
				// lookup the provided Gateway Name.
				backendRefs = append(backendRefs,
					types.NamespacedName{
						Namespace: gatewayapi.NamespaceDerefOr(backend.Namespace, grpcroute.Namespace),
						Name:      string(backend.Name),
					}.String(),
				)
			}
		}
	}
	return backendRefs
}

// addTLSRouteIndexers adds indexing on TLSRoute, for Service objects that are
// referenced in TLSRoute objects via `.spec.rules.backendRefs`. This helps in
// querying for TLSRoutes that are affected by a particular Service CRUD.
func addTLSRouteIndexers(ctx context.Context, mgr manager.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(ctx, &gwapiv1a2.TLSRoute{}, gatewayTLSRouteIndex, func(rawObj client.Object) []string {
		tlsRoute := rawObj.(*gwapiv1a2.TLSRoute)
		var gateways []string
		for _, parent := range tlsRoute.Spec.ParentRefs {
			if string(*parent.Kind) == gatewayapi.KindGateway {
				// If an explicit Gateway namespace is not provided, use the TLSRoute namespace to
				// lookup the provided Gateway Name.
				gateways = append(gateways,
					types.NamespacedName{
						Namespace: gatewayapi.NamespaceDerefOrAlpha(parent.Namespace, tlsRoute.Namespace),
						Name:      string(parent.Name),
					}.String(),
				)
			}
		}
		return gateways
	}); err != nil {
		return err
	}

	return mgr.GetFieldIndexer().IndexField(ctx, &gwapiv1a2.TLSRoute{}, backendTLSRouteIndex, backendTLSRouteIndexFunc)
}

func backendTLSRouteIndexFunc(rawObj client.Object) []string {
	tlsroute := rawObj.(*gwapiv1a2.TLSRoute)
	var backendRefs []string
	for _, rule := range tlsroute.Spec.Rules {
		for _, backend := range rule.BackendRefs {
			if backend.Kind == nil || string(*backend.Kind) == gatewayapi.KindService {
				// If an explicit Backend namespace is not provided, use the TLSRoute namespace to
				// lookup the provided Gateway Name.
				backendRefs = append(backendRefs,
					types.NamespacedName{
						Namespace: gatewayapi.NamespaceDerefOrAlpha(backend.Namespace, tlsroute.Namespace),
						Name:      string(backend.Name),
					}.String(),
				)
			}
		}
	}
	return backendRefs
}

// addTCPRouteIndexers adds indexing on TCPRoute, for Service objects that are
// referenced in TCPRoute objects via `.spec.rules.backendRefs`. This helps in
// querying for TCPRoutes that are affected by a particular Service CRUD.
func addTCPRouteIndexers(ctx context.Context, mgr manager.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(ctx, &gwapiv1a2.TCPRoute{}, gatewayTCPRouteIndex, func(rawObj client.Object) []string {
		tcpRoute := rawObj.(*gwapiv1a2.TCPRoute)
		var gateways []string
		for _, parent := range tcpRoute.Spec.ParentRefs {
			if string(*parent.Kind) == gatewayapi.KindGateway {
				// If an explicit Gateway namespace is not provided, use the TCPRoute namespace to
				// lookup the provided Gateway Name.
				gateways = append(gateways,
					types.NamespacedName{
						Namespace: gatewayapi.NamespaceDerefOrAlpha(parent.Namespace, tcpRoute.Namespace),
						Name:      string(parent.Name),
					}.String(),
				)
			}
		}
		return gateways
	}); err != nil {
		return err
	}

	return mgr.GetFieldIndexer().IndexField(ctx, &gwapiv1a2.TCPRoute{}, backendTCPRouteIndex, backendTCPRouteIndexFunc)
}

func backendTCPRouteIndexFunc(rawObj client.Object) []string {
	tcpRoute := rawObj.(*gwapiv1a2.TCPRoute)
	var backendRefs []string
	for _, rule := range tcpRoute.Spec.Rules {
		for _, backend := range rule.BackendRefs {
			if backend.Kind == nil || string(*backend.Kind) == gatewayapi.KindService {
				// If an explicit Backend namespace is not provided, use the TCPRoute namespace to
				// lookup the provided Gateway Name.
				backendRefs = append(backendRefs,
					types.NamespacedName{
						Namespace: gatewayapi.NamespaceDerefOrAlpha(backend.Namespace, tcpRoute.Namespace),
						Name:      string(backend.Name),
					}.String(),
				)
			}
		}
	}
	return backendRefs
}

// addUDPRouteIndexers adds indexing on UDPRoute.
//   - For Gateway objects that are referenced in UDPRoute objects via `.spec.parentRefs`. This helps in
//     querying for UDPRoutes that are affected by a particular Gateway CRUD.
//   - For Service objects that are referenced in UDPRoute objects via `.spec.rules.backendRefs`. This helps in
//     querying for UDPRoutes that are affected by a particular Service CRUD.
func addUDPRouteIndexers(ctx context.Context, mgr manager.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(ctx, &gwapiv1a2.UDPRoute{}, gatewayUDPRouteIndex, func(rawObj client.Object) []string {
		udpRoute := rawObj.(*gwapiv1a2.UDPRoute)
		var gateways []string
		for _, parent := range udpRoute.Spec.ParentRefs {
			if string(*parent.Kind) == gatewayapi.KindGateway {
				// If an explicit Gateway namespace is not provided, use the UDPRoute namespace to
				// lookup the provided Gateway Name.
				gateways = append(gateways,
					types.NamespacedName{
						Namespace: gatewayapi.NamespaceDerefOrAlpha(parent.Namespace, udpRoute.Namespace),
						Name:      string(parent.Name),
					}.String(),
				)
			}
		}
		return gateways
	}); err != nil {
		return err
	}

	return mgr.GetFieldIndexer().IndexField(ctx, &gwapiv1a2.UDPRoute{}, backendUDPRouteIndex, backendUDPRouteIndexFunc)
}

func backendUDPRouteIndexFunc(rawObj client.Object) []string {
	udproute := rawObj.(*gwapiv1a2.UDPRoute)
	var backendRefs []string
	for _, rule := range udproute.Spec.Rules {
		for _, backend := range rule.BackendRefs {
			if backend.Kind == nil || string(*backend.Kind) == gatewayapi.KindService {
				// If an explicit Backend namespace is not provided, use the UDPRoute namespace to
				// lookup the provided Gateway Name.
				backendRefs = append(backendRefs,
					types.NamespacedName{
						Namespace: gatewayapi.NamespaceDerefOrAlpha(backend.Namespace, udproute.Namespace),
						Name:      string(backend.Name),
					}.String(),
				)
			}
		}
	}
	return backendRefs
}

// addGatewayIndexers adds indexing on Gateway, for Secret objects and gatewayclass that are
// referenced in Gateway objects. This helps in querying for Gateways that are
// affected by a particular Secret CRUD.
func addGatewayIndexers(ctx context.Context, mgr manager.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(ctx, &gwapiv1.Gateway{}, secretGatewayIndex, secretGatewayIndexFunc); err != nil {
		loggers.LoggerAPKOperator.Errorf("Error while adding index %s, %v ", secretGatewayIndex, err)
		return err
	}

	return mgr.GetFieldIndexer().IndexField(ctx, &gwapiv1.Gateway{}, classGatewayIndex, func(rawObj client.Object) []string {
		gateway := rawObj.(*gwapiv1.Gateway)
		loggers.LoggerAPKOperator.Infof("Index %s added for %s:%s to %s", classGatewayIndex, gateway.Namespace, gateway.Name,
			gateway.Spec.GatewayClassName)
		return []string{string(gateway.Spec.GatewayClassName)}
	})
}

func secretGatewayIndexFunc(rawObj client.Object) []string {
	gateway := rawObj.(*gwapiv1.Gateway)
	var secretReferences []string
	for _, listener := range gateway.Spec.Listeners {
		if listener.TLS == nil || *listener.TLS.Mode != gwapiv1.TLSModeTerminate {
			continue
		}
		for _, cert := range listener.TLS.CertificateRefs {
			if *cert.Kind == gatewayapi.KindSecret {
				// If an explicit Secret namespace is not provided, use the Gateway namespace to
				// lookup the provided Secret Name.
				secretReferences = append(secretReferences,
					types.NamespacedName{
						Namespace: gatewayapi.NamespaceDerefOr(cert.Namespace, gateway.Namespace),
						Name:      string(cert.Name),
					}.String(),
				)
			}
		}
	}
	return secretReferences
}

// addBtlsIndexers adds indexing on BackendTLSPolicy, for ConfigMap objects that are
// referenced in BackendTLSPolicy objects. This helps in querying for BackendTLSPolicies that are
// affected by a particular ConfigMap CRUD.
func addBtlsIndexers(ctx context.Context, mgr manager.Manager) error {
	return mgr.GetFieldIndexer().IndexField(ctx, &gwapiv1a2.BackendTLSPolicy{}, configMapBtlsIndex, configMapBtlsIndexFunc)
}

func configMapBtlsIndexFunc(rawObj client.Object) []string {
	btls := rawObj.(*gwapiv1a2.BackendTLSPolicy)
	var configMapReferences []string
	if btls.Spec.TLS.CACertRefs != nil {
		for _, caCertRef := range btls.Spec.TLS.CACertRefs {
			if string(caCertRef.Kind) == gatewayapi.KindConfigMap {
				configMapReferences = append(configMapReferences,
					types.NamespacedName{
						Namespace: btls.Namespace,
						Name:      string(caCertRef.Name),
					}.String(),
				)
			}
		}
	}
	return configMapReferences
}
