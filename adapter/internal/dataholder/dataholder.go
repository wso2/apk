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

package dataholder

import (
	k8types "k8s.io/apimachinery/pkg/types"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
)

// The following variables will be used to store the state of the apk.
// This data should not be utilized by operator thread as its not designed for parallel access.
var (
	// This variable in the structure of gateway's namespaced name -> gateway
	gatewayMap map[string]gwapiv1.Gateway
)

func init() {
	gatewayMap = make(map[string]gwapiv1.Gateway)
}

// GetGatewayMap returns list of cached gateways
func GetGatewayMap() map[string]gwapiv1.Gateway {
	return gatewayMap
}

// UpdateGateway caches the gateway
func UpdateGateway(gateway gwapiv1.Gateway) {
	gatewayMap[k8types.NamespacedName{Name: gateway.Name, Namespace: gateway.Namespace}.String()] = gateway
}

// RemoveGateway removes the gateway from the cache
func RemoveGateway(gateway gwapiv1.Gateway) {
	delete(gatewayMap, k8types.NamespacedName{Name: gateway.Name, Namespace: gateway.Namespace}.String())
}

// GetAllGatewayListenerSections return the list of all the listeners that stored in the gateway cache
func GetAllGatewayListenerSections() []gwapiv1.Listener {
	listeners := make([]gwapiv1.Listener, 0)
	for _, gateway := range gatewayMap {
		listeners = append(listeners, gateway.Spec.Listeners...)
	}
	return listeners
}

// GetListenersAsPortalPortMap returns a map that have a structure protocol -> port -> list of listeners for that port and protocol combination
// Data is derived based on the current status of the gatwayMap cache
func GetListenersAsPortalPortMap() map[string]map[uint32][]gwapiv1.Listener {
	listenersAsPortalPortMap := make(map[string]map[uint32][]gwapiv1.Listener)
	for _, gateway := range gatewayMap {
		for _, listener := range gateway.Spec.Listeners {
			protocol := string(listener.Protocol)
			port := uint32(listener.Port)
			if portMap, portFound := listenersAsPortalPortMap[protocol]; portFound {
				if listenersList, listenerListFound := portMap[port]; listenerListFound {
					if listenersList == nil {
						listenersList = []gwapiv1.Listener{listener}
					} else {
						listenersList = append(listenersList, listener)
					}
					listenersAsPortalPortMap[protocol][port] = listenersList
				} else {
					listenerList := []gwapiv1.Listener{listener}
					listenersAsPortalPortMap[protocol][port] = listenerList
				}
			} else {
				listenersAsPortalPortMap[protocol] = make(map[uint32][]gwapiv1.Listener)
				listenerList := []gwapiv1.Listener{listener}
				listenersAsPortalPortMap[protocol][port] = listenerList
			}
		}
	}
	return listenersAsPortalPortMap
}
