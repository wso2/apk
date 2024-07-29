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
	"github.com/wso2/apk/adapter/internal/operator/gateway-api/ir"
	"github.com/wso2/apk/adapter/internal/operator/gateway-api/translator"
	"github.com/wso2/apk/adapter/internal/operator/message"
)

type Config struct {
	XdsIR *message.XdsIR
	Xds   *message.Xds
	// ExtensionManager  extension.Manager
	ProviderResources *message.ProviderResources
}

type Runner struct {
	Config
}

func New(cfg *Config) *Runner {
	return &Runner{Config: *cfg}
}

func (r *Runner) Name() string {
	return string("xds-translator")
}

// Start starts the xds-translator runner
func (r *Runner) Start(ctx context.Context) (err error) {
	go r.subscribeAndTranslate(ctx)
	loggers.LoggerAPKOperator.Info("Started xds translator ...")
	return
}

func (r *Runner) subscribeAndTranslate(ctx context.Context) {
	// Subscribe to resources
	message.HandleSubscription(message.Metadata{Runner: "xds-translator", Message: "xds-ir"}, r.XdsIR.Subscribe(ctx),
		func(update message.Update[string, *ir.Xds], errChan chan error) {
			loggers.LoggerAPKOperator.Info("Received an update in xds translator ...")
			key := update.Key
			val := update.Value

			if update.Delete {
				r.Xds.Delete(key)
			} else {
				// Translate to xds resources
				t := &translator.Translator{}

				// Set the extension manager if an extension is loaded
				// if r.ExtensionManager != nil {
				// 	t.ExtensionManager = &r.ExtensionManager
				// }

				// Set the rate limit service URL if global rate limiting is enabled.
				// if r.EnvoyGateway.RateLimit != nil {
				// 	t.GlobalRateLimit = &translator.GlobalRateLimitSettings{
				// 		ServiceURL: ratelimit.GetServiceURL(r.Namespace, r.DNSDomain),
				// 		FailClosed: r.EnvoyGateway.RateLimit.FailClosed,
				// 	}
				// 	if r.EnvoyGateway.RateLimit.Timeout != nil {
				// 		t.GlobalRateLimit.Timeout = r.EnvoyGateway.RateLimit.Timeout.Duration
				// 	}
				// }

				result, err := t.Translate(val)
				if err != nil {
					loggers.LoggerAPKOperator.Error("Failed to translate xds ir ", err)
					errChan <- err
				}

				// xDS translation is done in a best-effort manner, so the result
				// may contain partial resources even if there are errors.
				if result == nil {
					loggers.LoggerAPKOperator.Info("No xds resources to publish")
					return
				}

				// Get all status keys from watchable and save them in the map statusesToDelete.
				// Iterating through result.EnvoyPatchPolicyStatuses, any valid keys will be removed from statusesToDelete.
				// Remaining keys will be deleted from watchable before we exit this function.
				// statusesToDelete := make(map[ktypes.NamespacedName]bool)
				// for key := range r.ProviderResources.EnvoyPatchPolicyStatuses.LoadAll() {
				// 	statusesToDelete[key] = true
				// }

				// // Publish EnvoyPatchPolicyStatus
				// for _, e := range result.EnvoyPatchPolicyStatuses {
				// 	key := ktypes.NamespacedName{
				// 		Name:      e.Name,
				// 		Namespace: e.Namespace,
				// 	}
				// 	// Skip updating status for policies with empty status
				// 	// They may have been skipped in this translation because
				// 	// their target is not found (not relevant)
				// 	if !(reflect.ValueOf(e.Status).IsZero()) {
				// 		r.ProviderResources.EnvoyPatchPolicyStatuses.Store(key, e.Status)
				// 	}
				// 	delete(statusesToDelete, key)
				// }
				// // Discard the EnvoyPatchPolicyStatuses to reduce memory footprint
				// result.EnvoyPatchPolicyStatuses = nil

				// Publish
				r.Xds.Store(key, result)

				// Delete all the deletable status keys
				// for key := range statusesToDelete {
				// 	r.ProviderResources.EnvoyPatchPolicyStatuses.Delete(key)
				// }
			}
		},
	)
	loggers.LoggerAPKOperator.Info("subscriber shutting down")
}
