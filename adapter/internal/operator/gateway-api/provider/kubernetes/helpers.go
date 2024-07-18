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
 */
package kubernetes

import (
	"context"
	"fmt"

	"github.com/wso2/apk/adapter/internal/operator/constants"
	gatewayapi "github.com/wso2/apk/adapter/internal/operator/gateway-api"
	"github.com/wso2/apk/adapter/internal/operator/gateway-api/infrastructure/kubernetes/proxy"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
)

type controlledClasses struct {
	// matchedClasses holds all GatewayClass objects with matching controllerName.
	matchedClasses []*gwapiv1.GatewayClass
}

func (cc *controlledClasses) addMatch(gc *gwapiv1.GatewayClass) {
	cc.matchedClasses = append(cc.matchedClasses, gc)
}

func (cc *controlledClasses) removeMatch(gc *gwapiv1.GatewayClass) {
	// First remove gc from matchedClasses.
	for i, matchedGC := range cc.matchedClasses {
		if matchedGC.Name == gc.Name {
			cc.matchedClasses[i] = cc.matchedClasses[len(cc.matchedClasses)-1]
			cc.matchedClasses = cc.matchedClasses[:len(cc.matchedClasses)-1]
			break
		}
	}
}

// terminatesTLS returns true if the provided gateway contains a listener configured
// for TLS termination.
func terminatesTLS(listener *gwapiv1.Listener) bool {
	if listener.TLS != nil &&
		(listener.Protocol == gwapiv1.HTTPSProtocolType ||
			listener.Protocol == gwapiv1.TLSProtocolType) &&
		listener.TLS.Mode != nil &&
		*listener.TLS.Mode == gwapiv1.TLSModeTerminate {
		return true
	}
	return false
}

// refsSecret returns true if ref refers to a Secret.
func refsSecret(ref *gwapiv1.SecretObjectReference) bool {
	return (ref.Group == nil || *ref.Group == corev1.GroupName) &&
		(ref.Kind == nil || *ref.Kind == gatewayapi.KindSecret)
}

type ObjectKindNamespacedName struct {
	kind      string
	namespace string
	name      string
}

// validateParentRefs validates the provided routeParentReferences, returning the
// referenced Gateways managed by Envoy Gateway. The only supported parentRef
// is a Gateway.
func validateParentRefs(ctx context.Context, client client.Client, namespace string,
	gatewayClassController gwapiv1.GatewayController,
	routeParentReferences []gwapiv1.ParentReference) ([]gwapiv1.Gateway, error) {

	var gateways []gwapiv1.Gateway
	for i := range routeParentReferences {
		ref := routeParentReferences[i]
		if ref.Kind != nil && *ref.Kind != "Gateway" {
			return nil, fmt.Errorf("invalid Kind %q", *ref.Kind)
		}
		if ref.Group != nil && *ref.Group != gwapiv1.GroupName {
			return nil, fmt.Errorf("invalid Group %q", *ref.Group)
		}

		// Ensure the referenced Gateway exists, using the route's namespace unless
		// specified by the parentRef.
		ns := namespace
		if ref.Namespace != nil {
			ns = string(*ref.Namespace)
		}
		gwKey := types.NamespacedName{
			Namespace: ns,
			Name:      string(ref.Name),
		}

		gw := new(gwapiv1.Gateway)
		if err := client.Get(ctx, gwKey, gw); err != nil {
			return nil, fmt.Errorf("failed to get gateway %s/%s: %w", gwKey.Namespace, gwKey.Name, err)
		}

		gcKey := types.NamespacedName{Name: string(gw.Spec.GatewayClassName)}
		gc := new(gwapiv1.GatewayClass)
		if err := client.Get(ctx, gcKey, gc); err != nil {
			return nil, fmt.Errorf("failed to get gatewayclass %s: %w", gcKey.Name, err)
		}
		if gc.Spec.ControllerName == gatewayClassController {
			gateways = append(gateways, *gw)
		}
	}

	return gateways, nil
}

// validateBackendRef validates that ref is a reference to a local Service.
// TODO: Add support for:
//   - Validating weights.
//   - Validating ports.
//   - Referencing HTTPRoutes.
func validateBackendRef(ref *gwapiv1.BackendRef) error {
	switch {
	case ref == nil:
		return nil
	case gatewayapi.GroupDerefOr(ref.Group, corev1.GroupName) != corev1.GroupName && gatewayapi.GroupDerefOr(ref.Group, corev1.GroupName) != constants.GroupName:
		return fmt.Errorf("invalid group; must be nil, empty string or %q, given %q", constants.GroupName, gatewayapi.GroupDerefOr(ref.Group, corev1.GroupName))
	case gatewayapi.KindDerefOr(ref.Kind, gatewayapi.KindService) != gatewayapi.KindService && gatewayapi.KindDerefOr(ref.Kind, gatewayapi.KindService) != constants.KindBackend:
		return fmt.Errorf("invalid kind %q; must be %q or %q, given %q",
			*ref.BackendObjectReference.Kind, gatewayapi.KindService, constants.KindBackend, gatewayapi.GroupDerefOr(ref.Group, corev1.GroupName))
	}

	return nil
}

// infraName returns expected name for the EnvoyProxy infra resources.
// By default it returns hashed string from {GatewayNamespace}/{GatewayName},
func infraName(gateway *gwapiv1.Gateway) string {
	namespace := gateway.Namespace
	if len(namespace) > 5 {
		namespace = namespace[:5]
	}
	infraName := fmt.Sprintf("%s/%s", namespace, gateway.Name)
	return proxy.ExpectedResourceHashedName(infraName)
}
