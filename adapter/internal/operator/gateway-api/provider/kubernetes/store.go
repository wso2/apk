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
	"sync"

	corev1 "k8s.io/api/core/v1"
)

type nodeDetails struct {
	name    string
	address string
}

// kubernetesProviderStore holds cached information for the kubernetes provider.
type kubernetesProviderStore struct {
	// nodes holds information required for updating Gateway status with the Node
	// addresses, in case the Gateway is exposed on every Node of the cluster, using
	// Service of type NodePort.
	nodes map[string]nodeDetails
	mu    sync.Mutex
}

func newProviderStore() *kubernetesProviderStore {
	return &kubernetesProviderStore{
		nodes: make(map[string]nodeDetails),
	}
}

func (p *kubernetesProviderStore) addNode(n *corev1.Node) {
	details := nodeDetails{name: n.Name}

	var internalIP, externalIP string
	for _, addr := range n.Status.Addresses {
		if addr.Type == corev1.NodeExternalIP {
			externalIP = addr.Address
		}
		if addr.Type == corev1.NodeInternalIP {
			internalIP = addr.Address
		}
	}

	// In certain scenarios (like in local KinD clusters), the Node
	// externalIP is not provided, in that case we default back
	// to the internalIP of the Node.
	if externalIP != "" {
		details.address = externalIP
	} else if internalIP != "" {
		details.address = internalIP
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.nodes[n.Name] = details
}

func (p *kubernetesProviderStore) removeNode(n *corev1.Node) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.nodes, n.Name)
}

func (p *kubernetesProviderStore) listNodeAddresses() []string {
	addrs := []string{}
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, n := range p.nodes {
		if n.address != "" {
			addrs = append(addrs, n.address)
		}
	}
	return addrs
}
