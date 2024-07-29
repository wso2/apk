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

package runner

import (
	"context"
	"testing"
	"time"

	resourcev3 "github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"github.com/wso2/apk/adapter/internal/operator/gateway-api/ir"
	"github.com/wso2/apk/adapter/internal/operator/message"

	"github.com/stretchr/testify/require"
)

func TestRunner(t *testing.T) {
	// Setup
	xdsIR := new(message.XdsIR)
	xds := new(message.Xds)
	pResource := new(message.ProviderResources)
	// cfg, err := config.New()
	// require.NoError(t, err)
	r := New(&Config{
		// Server:            *cfg,
		ProviderResources: pResource,
		XdsIR:             xdsIR,
		Xds:               xds,
	})

	ctx := context.Background()
	// Start
	err := r.Start(ctx)
	require.NoError(t, err)

	// xDS is nil at start
	require.Equal(t, map[string]*ir.Xds{}, xdsIR.LoadAll())

	// test translation
	path := "example"
	res := ir.Xds{
		HTTP: []*ir.HTTPListener{
			{
				Name:      "test",
				Address:   "0.0.0.0",
				Port:      80,
				Hostnames: []string{"example.com"},
				Routes: []*ir.HTTPRoute{
					{
						Name: "test-route",
						PathMatch: &ir.StringMatch{
							Exact: &path,
						},
						Destination: &ir.RouteDestination{
							Name: "test-dest",
							Settings: []*ir.DestinationSetting{
								{
									Endpoints: []*ir.DestinationEndpoint{
										{
											Host: "10.11.12.13",
											Port: 8080,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	xdsIR.Store("test", &res)
	require.Eventually(t, func() bool {
		out := xds.LoadAll()
		if out == nil {
			return false
		}
		if out["test"] == nil {
			return false
		}
		// Ensure an xds listener is created
		return len(out["test"].XdsResources[resourcev3.ListenerType]) == 1
	}, time.Second*5, time.Millisecond*50)

	// Delete the IR triggering an xds delete
	xdsIR.Delete("test")
	require.Eventually(t, func() bool {
		out := xds.LoadAll()
		// Ensure that xds has no key, value pairs
		return len(out) == 0
	}, time.Second*5, time.Millisecond*50)

}

func TestRunner_withExtensionManager(t *testing.T) {
	// Setup
	xdsIR := new(message.XdsIR)
	xds := new(message.Xds)
	pResource := new(message.ProviderResources)

	// cfg, err := config.New()
	// require.NoError(t, err)
	r := New(&Config{
		// Server:            *cfg,
		ProviderResources: pResource,
		XdsIR:             xdsIR,
		Xds:               xds,
		// ExtensionManager:  &extManagerMock{},
	})

	ctx := context.Background()
	// Start
	err := r.Start(ctx)
	require.NoError(t, err)

	// xDS is nil at start
	require.Equal(t, map[string]*ir.Xds{}, xdsIR.LoadAll())

	// test translation
	path := "example"
	res := ir.Xds{
		HTTP: []*ir.HTTPListener{
			{
				Name:      "test",
				Address:   "0.0.0.0",
				Port:      80,
				Hostnames: []string{"example.com"},
				Routes: []*ir.HTTPRoute{
					{
						Name: "test-route",
						PathMatch: &ir.StringMatch{
							Exact: &path,
						},
						Destination: &ir.RouteDestination{
							Name: "test-dest",
							Settings: []*ir.DestinationSetting{
								{
									Endpoints: []*ir.DestinationEndpoint{
										{
											Host: "10.11.12.13",
											Port: 8080,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	xdsIR.Store("test", &res)
	require.Eventually(t, func() bool {
		out := xds.LoadAll()
		// xDS translation is done in a best-effort manner, so event the extension
		// manager returns an error, the xDS resources should still be created.
		return len(out) == 1
	}, time.Second*5, time.Millisecond*50)
}

// type extManagerMock struct {
// 	types.Manager
// }

// func (m *extManagerMock) GetPostXDSHookClient(xdsHookType v1alpha1.XDSTranslatorHook) types.XDSHookClient {
// 	if xdsHookType == v1alpha1.XDSHTTPListener {
// 		return &xdsHookClientMock{}
// 	}

// 	return nil
// }

// type xdsHookClientMock struct {
// 	types.XDSHookClient
// }

// func (c *xdsHookClientMock) PostHTTPListenerModifyHook(*listenerv3.Listener) (*listenerv3.Listener, error) {
// 	return nil, fmt.Errorf("assuming a network error during the call")
// }
