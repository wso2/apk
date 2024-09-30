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

				// Publish
				r.Xds.Store(key, result)
			}
		},
	)
	loggers.LoggerAPKOperator.Info("subscriber shutting down")
}
