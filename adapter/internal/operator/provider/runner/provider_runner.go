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
	"fmt"

	"github.com/wso2/apk/adapter/internal/loggers"
	provider "github.com/wso2/apk/adapter/internal/operator/gateway-api/provider/kubernetes"
	"github.com/wso2/apk/adapter/internal/operator/message"
	ctrl "sigs.k8s.io/controller-runtime"
)

type Config struct {
	ProviderResources *message.ProviderResources
}

type Runner struct {
	Config
}

func New(cfg *Config) *Runner {
	return &Runner{Config: *cfg}
}

func (r *Runner) Name() string {
	return string("provider")
}

// Start the provider runner
func (r *Runner) Start(ctx context.Context) (err error) {
	loggers.LoggerAPKOperator.Info("Started Runner Kubernetes operator...")
	cfg, err := ctrl.GetConfig()
	if err != nil {
		return fmt.Errorf("failed to get kubeconfig: %w", err)
	}
	p, err := provider.New(cfg, r.ProviderResources)
	if err != nil {
		return fmt.Errorf("failed to create provider %s: %w", "Kubernetes", err)
	}
	go func() {
		err := p.Start(ctx)
		if err != nil {
			loggers.LoggerAPKOperator.Error("Unable to start kubernetes operator provider", err)
		}
	}()
	return nil
}
