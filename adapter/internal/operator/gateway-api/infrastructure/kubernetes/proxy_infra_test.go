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
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gatewayapi "github.com/wso2/apk/adapter/internal/operator/gateway-api"
	"github.com/wso2/apk/adapter/internal/operator/gateway-api/infrastructure/kubernetes/proxy"
	"github.com/wso2/apk/adapter/internal/operator/gateway-api/ir"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	fakeclient "sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func newTestInfra(t *testing.T) *Infra {
	cli := fakeclient.NewClientBuilder().Build()
	return newTestInfraWithClient(t, cli)
}

func TestCmpBytes(t *testing.T) {
	m1 := map[string][]byte{}
	m1["a"] = []byte("aaa")
	m2 := map[string][]byte{}
	m2["a"] = []byte("aaa")

	assert.True(t, reflect.DeepEqual(m1, m2))
	assert.False(t, reflect.DeepEqual(nil, m2))
	assert.False(t, reflect.DeepEqual(m1, nil))
}

func newTestInfraWithClient(t *testing.T, cli client.Client) *Infra {
	return NewInfra(cli)
}

func TestCreateProxyInfra(t *testing.T) {
	// Infra with Gateway owner labels.
	infraWithLabels := ir.NewInfra()
	infraWithLabels.GetProxyInfra().GetProxyMetadata().Labels = proxy.EnvoyAppLabel()
	infraWithLabels.GetProxyInfra().GetProxyMetadata().Labels[gatewayapi.OwningGatewayNamespaceLabel] = "default"
	infraWithLabels.GetProxyInfra().GetProxyMetadata().Labels[gatewayapi.OwningGatewayNameLabel] = "test-gw"

	testCases := []struct {
		name   string
		in     *ir.Infra
		expect bool
	}{
		{
			name:   "infra-with-expected-labels",
			in:     infraWithLabels,
			expect: true,
		},
		{
			name:   "default infra without Gateway owner labels",
			in:     ir.NewInfra(),
			expect: false,
		},
		{
			name:   "nil-infra",
			in:     nil,
			expect: false,
		},
		{
			name: "nil-infra-proxy",
			in: &ir.Infra{
				Proxy: nil,
			},
			expect: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			kube := newTestInfra(t)
			// Create or update the proxy infra.
			err := kube.CreateOrUpdateProxyInfra(context.Background(), tc.in)
			if !tc.expect {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				// Verify all resources were created via the fake kube client.
				sa := &corev1.ServiceAccount{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: kube.Namespace,
						Name:      proxy.ExpectedResourceHashedName(tc.in.Proxy.Name),
					},
				}
				require.NoError(t, kube.Client.Get(context.Background(), client.ObjectKeyFromObject(sa), sa))

				// cm := &corev1.ConfigMap{
				// 	ObjectMeta: metav1.ObjectMeta{
				// 		Namespace: kube.Namespace,
				// 		Name:      proxy.ExpectedResourceHashedName(tc.in.Proxy.Name),
				// 	},
				// }
				// require.NoError(t, kube.Client.Get(context.Background(), client.ObjectKeyFromObject(cm), cm))

				deploy := &appsv1.Deployment{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: kube.Namespace,
						Name:      proxy.ExpectedResourceHashedName(tc.in.Proxy.Name),
					},
				}
				require.NoError(t, kube.Client.Get(context.Background(), client.ObjectKeyFromObject(deploy), deploy))

				svc := &corev1.Service{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: kube.Namespace,
						Name:      proxy.ExpectedResourceHashedName(tc.in.Proxy.Name),
					},
				}
				require.NoError(t, kube.Client.Get(context.Background(), client.ObjectKeyFromObject(svc), svc))
			}
		})
	}
}

func TestDeleteProxyInfra(t *testing.T) {

	testCases := []struct {
		name   string
		in     *ir.Infra
		expect bool
	}{
		{
			name:   "nil infra",
			in:     nil,
			expect: false,
		},
		{
			name:   "default infra",
			in:     ir.NewInfra(),
			expect: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			kube := newTestInfra(t)

			err := kube.DeleteProxyInfra(context.Background(), tc.in)
			if !tc.expect {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
