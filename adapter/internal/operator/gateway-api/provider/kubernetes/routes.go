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
	"errors"

	"github.com/wso2/apk/adapter/internal/loggers"
	gatewayapi "github.com/wso2/apk/adapter/internal/operator/gateway-api"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
)

// processHTTPRoutes finds HTTPRoutes corresponding to a gatewayNamespaceName, further checks for
// the backend references and pushes the HTTPRoutes to the resourceTree.
func (r *gatewayReconcilerNew) processHTTPRoutes(ctx context.Context, gatewayNamespaceName string,
	resourceMap *resourceMappings, resourceTree *gatewayapi.Resources) error {
	httpRouteList := &gwapiv1.HTTPRouteList{}

	// extensionRefFilters, err := r.getExtensionRefFilters(ctx)
	// if err != nil {
	// 	return err
	// }
	// for i := range extensionRefFilters {
	// 	filter := extensionRefFilters[i]
	// 	resourceMap.extensionRefFilters[utils.NamespacedName(&filter)] = filter
	// }

	if err := r.client.List(ctx, httpRouteList, &client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(gatewayHTTPRouteIndex, gatewayNamespaceName),
	}); err != nil {
		loggers.LoggerAPKOperator.Error("Failed to list HTTPRoutes", err)
		return err
	}

	for _, httpRoute := range httpRouteList.Items {
		loggers.LoggerAPKOperator.Infof("Processing HTTPRoute %s:%s", httpRoute.Namespace, httpRoute.Name)

		for _, rule := range httpRoute.Spec.Rules {
			for _, backendRef := range rule.BackendRefs {
				backendRef := backendRef
				if err := validateBackendRef(&backendRef.BackendRef); err != nil {
					loggers.LoggerAPKOperator.Error("Invalid backendRef ", err)
					continue
				}

				backendNamespace := gatewayapi.NamespaceDerefOr(backendRef.Namespace, httpRoute.Namespace)
				resourceMap.allAssociatedBackendRefs[gwapiv1.BackendObjectReference{
					Group:     backendRef.BackendObjectReference.Group,
					Kind:      backendRef.BackendObjectReference.Kind,
					Namespace: gatewayapi.NamespacePtrV1Alpha2(backendNamespace),
					Name:      backendRef.Name,
				}] = struct{}{}

				if backendNamespace != httpRoute.Namespace {
					from := ObjectKindNamespacedName{
						kind:      gatewayapi.KindHTTPRoute,
						namespace: httpRoute.Namespace,
						name:      httpRoute.Name,
					}
					to := ObjectKindNamespacedName{
						kind:      gatewayapi.KindDerefOr(backendRef.Kind, gatewayapi.KindService),
						namespace: backendNamespace,
						name:      string(backendRef.Name),
					}
					refGrant, err := r.findReferenceGrant(ctx, from, to)
					switch {
					case err != nil:
						loggers.LoggerAPKOperator.Error("Failed to find ReferenceGrant", err)
					case refGrant == nil:
						loggers.LoggerAPKOperator.Infof("No matching ReferenceGrants found for kind %s:%s to %s:%s",
							from.kind, from.namespace, to.kind, to.namespace)
					default:
						resourceTree.ReferenceGrants = append(resourceTree.ReferenceGrants, refGrant)
						loggers.LoggerAPKOperator.Infof("Added ReferenceGrant to resource map %s:%s",
							refGrant.Namespace, refGrant.Name)
					}
				}
			}

			for i := range rule.Filters {
				filter := rule.Filters[i]

				if err := gatewayapi.ValidateHTTPRouteFilter(&filter); err != nil {
					loggers.LoggerAPKOperator.Errorf("Bypassing filter rule for index %v, %v", i, err)
					continue
				}

				// Load in the backendRefs from any requestMirrorFilters on the HTTPRoute
				if filter.Type == gwapiv1.HTTPRouteFilterRequestMirror {
					// Make sure the config actually exists
					mirrorFilter := filter.RequestMirror
					if mirrorFilter == nil {
						loggers.LoggerAPKOperator.Error(errors.New("invalid requestMirror filter"), "bypassing filter rule", "index", i)
						continue
					}

					mirrorBackendObj := mirrorFilter.BackendRef
					// Wrap the filter's BackendObjectReference into a BackendRef so we can use existing tooling to check it
					weight := int32(1)
					mirrorBackendRef := gwapiv1.BackendRef{
						BackendObjectReference: mirrorBackendObj,
						Weight:                 &weight,
					}

					if err := validateBackendRef(&mirrorBackendRef); err != nil {
						loggers.LoggerAPKOperator.Error("Invalid backendRef ", err)
						continue
					}

					backendNamespace := gatewayapi.NamespaceDerefOr(mirrorBackendRef.Namespace,
						httpRoute.Namespace)
					resourceMap.allAssociatedBackendRefs[gwapiv1.BackendObjectReference{
						Group:     mirrorBackendRef.BackendObjectReference.Group,
						Kind:      mirrorBackendRef.BackendObjectReference.Kind,
						Namespace: gatewayapi.NamespacePtrV1Alpha2(backendNamespace),
						Name:      mirrorBackendRef.Name,
					}] = struct{}{}

					if backendNamespace != httpRoute.Namespace {
						from := ObjectKindNamespacedName{
							kind:      gatewayapi.KindHTTPRoute,
							namespace: httpRoute.Namespace,
							name:      httpRoute.Name,
						}
						to := ObjectKindNamespacedName{
							kind:      gatewayapi.KindDerefOr(mirrorBackendRef.Kind, gatewayapi.KindService),
							namespace: backendNamespace,
							name:      string(mirrorBackendRef.Name),
						}
						refGrant, err := r.findReferenceGrant(ctx, from, to)
						switch {
						case err != nil:
							loggers.LoggerAPKOperator.Error("Failed to find ReferenceGrant", err)
						case refGrant == nil:
							loggers.LoggerAPKOperator.Infof("No matching ReferenceGrants found from %s:%s target %s:%s",
								from.kind, from.namespace, to.kind, to.namespace)
						default:
							resourceTree.ReferenceGrants = append(resourceTree.ReferenceGrants, refGrant)
							loggers.LoggerAPKOperator.Infof("Added ReferenceGrant to resource map %s:%s",
								refGrant.Namespace, refGrant.Name)
						}
					}
				} else if filter.Type == gwapiv1.HTTPRouteFilterExtensionRef {
					// NOTE: filters must be in the same namespace as the HTTPRoute
					// Check if it's a Kind managed by an extension and add to resourceTree
					key := types.NamespacedName{
						Namespace: httpRoute.Namespace,
						Name:      string(filter.ExtensionRef.Name),
					}
					extRefFilter, ok := resourceMap.extensionRefFilters[key]
					if !ok {
						loggers.LoggerAPKOperator.Errorf("Filter not found; bypassing rule name %v for index %v",
							filter.ExtensionRef.Name, i)
						continue
					}

					resourceTree.ExtensionRefFilters = append(resourceTree.ExtensionRefFilters, extRefFilter)
				}
			}
		}

		resourceMap.allAssociatedNamespaces[httpRoute.Namespace] = struct{}{}
		// Discard Status to reduce memory consumption in watchable
		// It will be recomputed by the gateway-api layer
		httpRoute.Status = gwapiv1.HTTPRouteStatus{}
		resourceTree.HTTPRoutes = append(resourceTree.HTTPRoutes, &httpRoute)
	}

	return nil
}
