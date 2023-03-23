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

package kubernetes

import (
	"context"
	"net"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/wso2/apk/test/integration/integration/constants"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/gateway-api/conformance/utils/config"
)

// NamespacesMustBeAccepted waits until all Pods are marked ready.
func NamespacesMustBeAccepted(t *testing.T, c client.Client, timeoutConfig config.TimeoutConfig, namespaces []string) {
	t.Helper()

	waitErr := wait.PollImmediate(1*time.Second, timeoutConfig.NamespacesMustBeReady, func() (bool, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		for _, ns := range namespaces {
			podList := &v1.PodList{}
			err := c.List(ctx, podList, client.InNamespace(ns))
			if err != nil {
				t.Errorf("Error listing Pods: %v", err)
			}
			for _, pod := range podList.Items {
				if !findPodConditionInList(t, pod.Status.Conditions, "Ready", "True") &&
					pod.Status.Phase != v1.PodSucceeded {
					t.Logf("%s/%s Pod not ready yet", ns, pod.Name)
					return false, nil
				}
			}
		}
		t.Logf("Gateways and Pods in %s namespaces ready", strings.Join(namespaces, ", "))
		return true, nil
	})
	require.NoErrorf(t, waitErr, "error waiting for %s namespaces to be ready", strings.Join(namespaces, ", "))
}

// WaitForGatewayAddress waits until at least one IP Address has been set in the
// Gateway infra exposed service.
func WaitForGatewayAddress(t *testing.T, c client.Client, timeoutConfig config.TimeoutConfig) string {
	// Use http port for now, ideally we should get the port from the Gateway or from a config.
	port := strconv.FormatInt(int64(constants.GatewayServicePort), 10)
	return WaitForIPAddress(t, c, timeoutConfig, port)
}

// WaitForAPIListenerAddress waits until at least one IP Address has been set in the
// Gateway infra exposed service.
func WaitForAPIListenerAddress(t *testing.T, c client.Client, timeoutConfig config.TimeoutConfig) string {
	// Use http port for now, ideally we should get the port from the APIListener or from a config.
	port := strconv.FormatInt(int64(constants.APIListenerServicePort), 10)
	return WaitForIPAddress(t, c, timeoutConfig, port)
}

// WaitForIPAddress waits until at least one IP Address has been set in the
// Gateway infra exposed service.
func WaitForIPAddress(t *testing.T, c client.Client, timeoutConfig config.TimeoutConfig, port string) string {
	t.Helper()

	var ipAddr string
	name := constants.GatewayServiceName
	namespace := constants.GatewayServiceNamespace

	waitErr := wait.PollImmediate(1*time.Second, timeoutConfig.GatewayMustHaveAddress, func() (bool, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		svc := &v1.Service{}
		if err := c.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, svc); err != nil {
			t.Fatalf("failed to get service %s/%s: %v", namespace, name, err)
			return false, nil
		}

		if len(svc.Status.LoadBalancer.Ingress) == 0 {
			t.Fatalf("service %s/%s has no external IP address", namespace, name)
		}

		ipAddr = svc.Status.LoadBalancer.Ingress[0].IP

		return true, nil
	})

	require.NoErrorf(t, waitErr, "error waiting for Gateway service to have an IP address")
	return net.JoinHostPort(ipAddr, port)
}

func findPodConditionInList(t *testing.T, conditions []v1.PodCondition, condName, condValue string) bool {
	t.Helper()

	for _, cond := range conditions {
		if cond.Type == v1.PodConditionType(condName) {
			if cond.Status == v1.ConditionStatus(condValue) {
				return true
			}
			t.Logf("%s condition set to %s, expected %s", condName, cond.Status, condValue)
		}
	}

	t.Logf("%s was not in conditions list", condName)
	return false
}
