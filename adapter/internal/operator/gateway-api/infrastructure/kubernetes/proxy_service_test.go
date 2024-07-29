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
	"testing"

	"github.com/stretchr/testify/require"
	gatewayapi "github.com/wso2/apk/adapter/internal/operator/gateway-api"
	"github.com/wso2/apk/adapter/internal/operator/gateway-api/infrastructure/kubernetes/proxy"
	"github.com/wso2/apk/adapter/internal/operator/gateway-api/ir"
)

func TestDeleteProxyService(t *testing.T) {
	testCases := []struct {
		name string
	}{
		{
			name: "delete service",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			kube := newTestInfra(t)
			infra := ir.NewInfra()

			infra.Proxy.GetProxyMetadata().Labels[gatewayapi.OwningGatewayNamespaceLabel] = "default"
			infra.Proxy.GetProxyMetadata().Labels[gatewayapi.OwningGatewayNameLabel] = infra.Proxy.Name
			r := proxy.NewResourceRender(kube.Namespace, infra.GetProxyInfra())
			err := kube.createOrUpdateService(context.Background(), r)
			require.NoError(t, err)

			err = kube.deleteService(context.Background(), r)
			require.NoError(t, err)
		})
	}
}
