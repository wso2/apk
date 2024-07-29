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

	"github.com/wso2/apk/adapter/internal/loggers"
	"github.com/wso2/apk/adapter/internal/operator/gateway-api/infrastructure"
	"github.com/wso2/apk/adapter/internal/operator/gateway-api/ir"
	"github.com/wso2/apk/adapter/internal/operator/message"
)

type Config struct {
	InfraIR *message.InfraIR
}

type Runner struct {
	Config
	mgr infrastructure.Manager
}

func (r *Runner) Name() string {
	return string("infrastructure")
}

func New(cfg *Config) *Runner {
	return &Runner{Config: *cfg}
}

// Start starts the infrastructure runner
func (r *Runner) Start(ctx context.Context) (err error) {

	r.mgr, err = infrastructure.NewManager()
	if err != nil {
		loggers.LoggerAPKOperator.Error(err, "failed to create new manager")
		return err
	}
	go r.subscribeToProxyInfraIR(ctx)

	// Enable global ratelimit if it has been configured.
	// if r.EnvoyGateway.RateLimit != nil {
	// 	go r.enableRateLimitInfra(ctx)
	// }

	loggers.LoggerAPKOperator.Info("started")
	return
}

func (r *Runner) subscribeToProxyInfraIR(ctx context.Context) {
	// Subscribe to resources
	message.HandleSubscription(message.Metadata{Runner: string("infrastructure"), Message: "infra-ir"}, r.InfraIR.Subscribe(ctx),
		func(update message.Update[string, *ir.Infra], errChan chan error) {
			loggers.LoggerAPKOperator.Info("Received an update in infrastructure provider ...")
			val := update.Value

			if update.Delete {
				if err := r.mgr.DeleteProxyInfra(ctx, val); err != nil {
					loggers.LoggerAPKOperator.Error(err, "failed to delete infra")
					errChan <- err
				}
			} else {
				// Manage the proxy infra.
				if len(val.Proxy.Listeners) == 0 {
					loggers.LoggerAPKOperator.Info("Infra IR was updated, but no listeners were found. Skipping infra creation.")
					return
				}

				if err := r.mgr.CreateOrUpdateProxyInfra(ctx, val); err != nil {
					loggers.LoggerAPKOperator.Error("Failed to create new infra ", err)
					errChan <- err
				}
			}
		},
	)
	loggers.LoggerAPKOperator.Info("infra subscriber shutting down")
}

// func (r *Runner) enableRateLimitInfra(ctx context.Context) {
// 	if err := r.mgr.CreateOrUpdateRateLimitInfra(ctx); err != nil {
// 		loggers.LoggerAPKOperator.Error(err, "failed to create ratelimit infra")
// 	}
// }
