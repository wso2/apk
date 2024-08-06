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
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNodeDetailsAddressStore(t *testing.T) {
	store := newProviderStore()
	testCases := []struct {
		name              string
		nodeObject        *corev1.Node
		expectedAddresses []string
	}{
		{
			name: "No node addresses",
			nodeObject: &corev1.Node{
				ObjectMeta: v1.ObjectMeta{Name: "node1"},
				Status:     corev1.NodeStatus{Addresses: []corev1.NodeAddress{{}}},
			},
			expectedAddresses: []string{},
		},
		{
			name: "only external address",
			nodeObject: &corev1.Node{
				ObjectMeta: v1.ObjectMeta{Name: "node1"},
				Status: corev1.NodeStatus{Addresses: []corev1.NodeAddress{{
					Address: "1.1.1.1",
					Type:    corev1.NodeExternalIP,
				}}},
			},
			expectedAddresses: []string{"1.1.1.1"},
		},
		{
			name: "only internal address",
			nodeObject: &corev1.Node{
				ObjectMeta: v1.ObjectMeta{Name: "node1"},
				Status: corev1.NodeStatus{Addresses: []corev1.NodeAddress{{
					Address: "1.1.1.1",
					Type:    corev1.NodeInternalIP,
				}}},
			},
			expectedAddresses: []string{"1.1.1.1"},
		},
		{
			name: "prefer external address",
			nodeObject: &corev1.Node{
				ObjectMeta: v1.ObjectMeta{Name: "node1"},
				Status: corev1.NodeStatus{Addresses: []corev1.NodeAddress{
					{
						Address: "1.1.1.1",
						Type:    corev1.NodeExternalIP,
					},
					{
						Address: "2.2.2.2",
						Type:    corev1.NodeInternalIP,
					},
				}},
			},
			expectedAddresses: []string{"1.1.1.1"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			store.addNode(tc.nodeObject)
			assert.Equal(t, tc.expectedAddresses, store.listNodeAddresses())
			store.removeNode(tc.nodeObject)
		})
	}
}

func TestRace(t *testing.T) {
	s := newProviderStore()

	go func() {
		for {
			s.addNode(&corev1.Node{
				ObjectMeta: v1.ObjectMeta{Name: "node1"},
				Status:     corev1.NodeStatus{Addresses: []corev1.NodeAddress{{}}},
			})
		}
	}()

	_ = s.listNodeAddresses()
}
