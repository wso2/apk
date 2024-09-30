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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	fakeclient "sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var (
	namespace = "apk"
)

// func TestCreateOrUpdateProxyConfigMap(t *testing.T) {
// 	infra := ir.NewInfra()
// 	infra.Proxy.Name = "test"
// 	infra.Proxy.GetProxyMetadata().Labels[gatewayapi.OwningGatewayNamespaceLabel] = "default"
// 	infra.Proxy.GetProxyMetadata().Labels[gatewayapi.OwningGatewayNameLabel] = infra.Proxy.Name

// 	testCases := []struct {
// 		name    string
// 		current *corev1.ConfigMap
// 		expect  *corev1.ConfigMap
// 	}{
// 		{
// 			name: "create configmap",
// 			expect: &corev1.ConfigMap{
// 				ObjectMeta: metav1.ObjectMeta{
// 					Namespace: namespace,
// 					Name:      "envoy-test-9f86d081",
// 					Labels: map[string]string{
// 						"app.kubernetes.io/name":               "envoy",
// 						"app.kubernetes.io/component":          "proxy",
// 						"app.kubernetes.io/managed-by":         "apk-gateway",
// 						gatewayapi.OwningGatewayNamespaceLabel: "default",
// 						gatewayapi.OwningGatewayNameLabel:      "test",
// 					},
// 				},
// 				Data: map[string]string{
// 					proxy.SdsCAFilename:   proxy.SdsCAConfigMapData,
// 					proxy.SdsCertFilename: proxy.SdsCertConfigMapData,
// 				},
// 			},
// 		},
// 		{
// 			name: "update configmap",
// 			current: &corev1.ConfigMap{
// 				ObjectMeta: metav1.ObjectMeta{
// 					Namespace: namespace,
// 					Name:      "envoy-test",
// 					Labels: map[string]string{
// 						"app.kubernetes.io/name":               "envoy",
// 						"app.kubernetes.io/component":          "proxy",
// 						"app.kubernetes.io/managed-by":         "apk-gateway",
// 						gatewayapi.OwningGatewayNamespaceLabel: "default",
// 						gatewayapi.OwningGatewayNameLabel:      "test",
// 					},
// 				},
// 				Data: map[string]string{"foo": "bar"},
// 			},
// 			expect: &corev1.ConfigMap{
// 				ObjectMeta: metav1.ObjectMeta{
// 					Namespace: namespace,
// 					Name:      "envoy-test-9f86d081",
// 					Labels: map[string]string{
// 						"app.kubernetes.io/name":               "envoy",
// 						"app.kubernetes.io/component":          "proxy",
// 						"app.kubernetes.io/managed-by":         "apk-gateway",
// 						gatewayapi.OwningGatewayNamespaceLabel: "default",
// 						gatewayapi.OwningGatewayNameLabel:      "test",
// 					},
// 				},
// 				Data: map[string]string{
// 					proxy.SdsCAFilename:   proxy.SdsCAConfigMapData,
// 					proxy.SdsCertFilename: proxy.SdsCertConfigMapData,
// 				},
// 			},
// 		},
// 	}

// 	for _, tc := range testCases {
// 		tc := tc
// 		t.Run(tc.name, func(t *testing.T) {
// 			var cli client.Client
// 			if tc.current != nil {
// 				cli = fakeclient.NewClientBuilder().WithObjects(tc.current).Build()
// 			} else {
// 				cli = fakeclient.NewClientBuilder().Build()
// 			}
// 			kube := NewInfra(cli)
// 			r := proxy.NewResourceRender(kube.Namespace, infra.GetProxyInfra())
// 			err := kube.createOrUpdateConfigMap(context.Background(), r)
// 			require.NoError(t, err)
// 			actual := &corev1.ConfigMap{
// 				ObjectMeta: metav1.ObjectMeta{
// 					Namespace: tc.expect.Namespace,
// 					Name:      tc.expect.Name,
// 				},
// 			}
// 			require.NoError(t, kube.Client.Get(context.Background(), client.ObjectKeyFromObject(actual), actual))
// 			require.Equal(t, tc.expect.Data, actual.Data)
// 			assert.True(t, apiequality.Semantic.DeepEqual(tc.expect.Labels, actual.Labels))
// 		})
// 	}
// }

func TestDeleteConfigProxyMap(t *testing.T) {

	infra := ir.NewInfra()
	infra.Proxy.Name = "test"

	testCases := []struct {
		name    string
		current *corev1.ConfigMap
		expect  bool
	}{
		{
			name: "delete configmap",
			current: &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: namespace,
					Name:      "envoy-test",
				},
			},
			expect: true,
		},
		{
			name: "configmap not found",
			current: &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: namespace,
					Name:      "foo",
				},
			},
			expect: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			cli := fakeclient.NewClientBuilder().WithObjects(tc.current).Build()
			kube := NewInfra(cli)

			infra.Proxy.GetProxyMetadata().Labels[gatewayapi.OwningGatewayNamespaceLabel] = "default"
			infra.Proxy.GetProxyMetadata().Labels[gatewayapi.OwningGatewayNameLabel] = infra.Proxy.Name

			r := proxy.NewResourceRender(kube.Namespace, infra.GetProxyInfra())
			cm := &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: kube.Namespace,
					Name:      r.Name(),
				},
			}
			err := kube.Client.Delete(context.Background(), cm)
			require.NoError(t, err)
		})
	}
}
