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
	"testing"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
)

type GroupKindNamespacedName struct {
	Group     gwapiv1.Group
	Kind      gwapiv1.Kind
	Namespace gwapiv1.Namespace
	Name      gwapiv1.ObjectName
}

// TestGatewayClassHasMatchingController tests the hasMatchingController
// predicate function.
func TestGatewayClassHasMatchingController(t *testing.T) {
	testCases := []struct {
		name   string
		obj    *gwapiv1.GatewayClass
		client client.Client
		expect bool
	}{
		{
			name:   "matching controller name",
			obj:    GetGatewayClass("test-gc", gatewayClassControllerName, nil),
			expect: true,
		},
		{
			name:   "non-matching controller name",
			obj:    GetGatewayClass("test-gc", "not.configured/controller", nil),
			expect: false,
		},
	}

	r := gatewayReconcilerNew{}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			res := r.hasMatchingController(tc.obj)
			require.Equal(t, tc.expect, res)
		})
	}
}

// GetGatewayClass returns a sample GatewayClass.
func GetGatewayClass(name string, controller gwapiv1.GatewayController, envoyProxy *GroupKindNamespacedName) *gwapiv1.GatewayClass {
	gwc := &gwapiv1.GatewayClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: gwapiv1.GatewayClassSpec{
			ControllerName: controller,
		},
	}

	if envoyProxy != nil {
		gwc.Spec.ParametersRef = &gwapiv1.ParametersReference{
			Group:     envoyProxy.Group,
			Kind:      envoyProxy.Kind,
			Name:      string(envoyProxy.Name),
			Namespace: &envoyProxy.Namespace,
		}
	}

	return gwc
}
