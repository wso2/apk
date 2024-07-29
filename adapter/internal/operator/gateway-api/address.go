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
	v1 "sigs.k8s.io/gateway-api/apis/v1"
)

var _ AddressesTranslator = (*Translator)(nil)

type AddressesTranslator interface {
	ProcessAddresses(gateways []*GatewayContext, infraIR InfraIRMap)
}

func (t *Translator) ProcessAddresses(gateways []*GatewayContext, infraIR InfraIRMap) {
	for _, gateway := range gateways {
		// Infra IR already exist
		irKey := t.getIRKey(gateway.Gateway)
		gwInfraIR := infraIR[irKey]

		var ips []string
		for _, gwadr := range gateway.Spec.Addresses {
			if *gwadr.Type == v1.IPAddressType {
				ips = append(ips, gwadr.Value)
			}
		}
		gwInfraIR.Proxy.Addresses = ips
	}
}
