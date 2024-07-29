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

package infrastructure

import (
	"context"

	"github.com/wso2/apk/adapter/internal/operator/gateway-api/infrastructure/kubernetes"
	"github.com/wso2/apk/adapter/internal/operator/gateway-api/ir"
	"sigs.k8s.io/controller-runtime/pkg/client"
	clicfg "sigs.k8s.io/controller-runtime/pkg/client/config"
)

var _ Manager = (*kubernetes.Infra)(nil)

// Manager provides the scaffolding for managing infrastructure.
type Manager interface {
	// CreateOrUpdateProxyInfra creates or updates infra.
	CreateOrUpdateProxyInfra(ctx context.Context, infra *ir.Infra) error
	// DeleteProxyInfra deletes infra.
	DeleteProxyInfra(ctx context.Context, infra *ir.Infra) error
	// // CreateOrUpdateRateLimitInfra creates or updates rate limit infra.
	// CreateOrUpdateRateLimitInfra(ctx context.Context) error
	// DeleteRateLimitInfra deletes rate limit infra.
	// DeleteRateLimitInfra(ctx context.Context) error
}

// NewManager returns a new infrastructure Manager.
func NewManager() (Manager, error) {
	var mgr Manager
	cli, err := client.New(clicfg.GetConfigOrDie(), client.Options{})
	if err != nil {
		return nil, err
	}
	mgr = kubernetes.NewInfra(cli)
	return mgr, nil
}
