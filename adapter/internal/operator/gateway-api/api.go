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
	"strings"

	"github.com/wso2/apk/adapter/internal/operator/gateway-api/ir"
	dpv1alpha2 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha2"
	"k8s.io/apimachinery/pkg/types"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func (t *Translator) ProcessAPIs(apis []*dpv1alpha2.API, httpRoutes []*gwapiv1.HTTPRoute, gateways []*GatewayContext, httpRouteContexts []*HTTPRouteContext, resources *Resources, xdsIR XdsIRMap) {
	for _, api := range apis {
		// For each API find all the related httproutes. If not found any continue.
		routeNamespacedNames := make([]string, 0)
		for _, httpRouteProds := range api.Spec.Production {
			for _, httpRouteRef := range httpRouteProds.RouteRefs {
				routeNamespacedNames = append(routeNamespacedNames, types.NamespacedName{
					Namespace: api.GetNamespace(),
					Name:      httpRouteRef,
				}.String())
			}
		}
		for _, httpRouteSands := range api.Spec.Sandbox {
			for _, httpRouteRef := range httpRouteSands.RouteRefs {
				routeNamespacedNames = append(routeNamespacedNames, types.NamespacedName{
					Namespace: api.GetNamespace(),
					Name:      httpRouteRef,
				}.String())
			}
		}
		for _, httpRouteContext := range httpRouteContexts {
			if !Contains(routeNamespacedNames, GetNamespacedName(httpRouteContext)) {
				continue
			}
			prefix := irRoutePrefix(httpRouteContext)
			parentRefs := GetParentReferences(httpRouteContext)
			for _, p := range parentRefs {
				parentRefCtx := GetRouteParentContext(httpRouteContext, p)
				gtwCtx := parentRefCtx.GetGateway()
				if gtwCtx == nil {
					continue
				}
				irKey := t.getIRKey(gtwCtx.Gateway)
				for _, listener := range parentRefCtx.listeners {
					irListener := xdsIR[irKey].GetHTTPListener(irListenerName(listener))
					if irListener != nil {
						for _, r := range irListener.Routes {
							if strings.HasPrefix(r.Name, prefix) {
								extAuth := buildExtAuth()
								r.ExtAuth = extAuth
							}
						}
					}
				}
			}
		}

	}
}

func buildExtAuth() *ir.ExtAuth {
	grpcExtAuthService := ir.GRPCExtAuthService{
		Destination: ir.RouteDestination{
			Name: ExtAuthClusterName,
		},
	}
	flag := true
	extAuth := &ir.ExtAuth{
		Name:                "common",
		GRPC:                &grpcExtAuthService,
		UseBootstrapCluster: &flag,
	}
	return extAuth
}
