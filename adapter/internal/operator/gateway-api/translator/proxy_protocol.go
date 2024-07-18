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

package translator

import (
	listenerv3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	proxyprotocolv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/listener/proxy_protocol/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	"github.com/wso2/apk/adapter/internal/operator/gateway-api/ir"
	"google.golang.org/protobuf/types/known/anypb"
)

// patchProxyProtocolFilter builds and appends the Proxy Protocol Filter to the
// HTTP Listener's Listener Filters if applicable.
func patchProxyProtocolFilter(xdsListener *listenerv3.Listener, irListener *ir.HTTPListener) {
	// Return early if unset
	if xdsListener == nil || irListener == nil || !irListener.EnableProxyProtocol {
		return
	}

	// Return early if filter already exists.
	for _, filter := range xdsListener.ListenerFilters {
		if filter.Name == wellknown.ProxyProtocol {
			return
		}
	}

	proxyProtocolFilter := buildProxyProtocolFilter()

	if proxyProtocolFilter != nil {
		xdsListener.ListenerFilters = append(xdsListener.ListenerFilters, proxyProtocolFilter)
	}
}

// buildProxypProtocolFilter returns a Proxy Protocol listener filter from the provided IR listener.
func buildProxyProtocolFilter() *listenerv3.ListenerFilter {
	pp := &proxyprotocolv3.ProxyProtocol{}

	ppAny, err := anypb.New(pp)
	if err != nil {
		return nil
	}

	return &listenerv3.ListenerFilter{
		Name: wellknown.ProxyProtocol,
		ConfigType: &listenerv3.ListenerFilter_TypedConfig{
			TypedConfig: ppAny,
		},
	}
}
