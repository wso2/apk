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
	"flag"
	"fmt"
	"strings"

	"github.com/wso2/apk/adapter/config"
	"github.com/wso2/apk/adapter/internal/loggers"
	dpcontrollers "github.com/wso2/apk/adapter/internal/operator/controllers/dp"
	"github.com/wso2/apk/adapter/internal/operator/gateway-api/provider"
	"github.com/wso2/apk/adapter/internal/operator/message"
	"github.com/wso2/apk/adapter/internal/operator/status"
	"github.com/wso2/apk/adapter/internal/operator/synchronizer"
	"github.com/wso2/apk/adapter/pkg/logging"
	"github.com/wso2/apk/adapter/pkg/metrics"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// Provider is the scaffolding for the Kubernetes provider. It sets up dependencies
// and defines the topology of the provider and its managed components, wiring
// them together.
type Provider struct {
	client  client.Client
	manager manager.Manager
}

// New creates a new Provider from the provided EnvoyGateway.
func New(cfg *rest.Config, resources *message.ProviderResources) (*Provider, error) {
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()
	log.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	// TODO: Decide which mgr opts should be exposed through envoygateway.provider.kubernetes API.
	mgrOpts := manager.Options{
		LeaderElection:         false,
		Scheme:                 provider.GetScheme(),
		HealthProbeBindAddress: ":8081",
		LeaderElectionID:       "operator-lease.apk.wso2.com",
	}

	conf := config.ReadConfigs()
	metricsConfig := conf.Adapter.Metrics
	if metricsConfig.Enabled {
		mgrOpts.Metrics.BindAddress = fmt.Sprintf(":%d", metricsConfig.Port)
		// Register the metrics collector
		if strings.EqualFold(metricsConfig.Type, metrics.PrometheusMetricType) {
			loggers.LoggerAPKOperator.Info("Registering Prometheus metrics collector.")
			metrics.RegisterPrometheusCollector()
		}
	} else {
		mgrOpts.Metrics.BindAddress = "0"
	}

	mgr, err := ctrl.NewManager(cfg, mgrOpts)

	if err != nil {
		return nil, fmt.Errorf("failed to create manager: %w", err)
	}

	updateHandler := status.NewUpdateHandler(mgr.GetClient())
	if err := mgr.Add(updateHandler); err != nil {
		return nil, fmt.Errorf("failed to add status update handler %w", err)
	}

	// Create and register the controllers with the manager.
	if err := InitGatewayController(mgr, resources, updateHandler); err != nil {
		return nil, fmt.Errorf("failted to create gatewayapi controller: %w", err)
	}

	// Add health check health probes.
	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		return nil, fmt.Errorf("unable to set up health check: %w", err)
	}

	// Add ready check health probes.
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		return nil, fmt.Errorf("unable to set up ready check: %w", err)
	}

	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2600, logging.BLOCKER, "Unable to start manager: %v", err))
	}

	// TODO: Decide on a buffer size and add to config.
	ch := make(chan *synchronizer.APIEvent, 10)
	successChannel := make(chan synchronizer.SuccessEvent, 10)

	gatewaych := make(chan synchronizer.GatewayEvent, 10)
	if err := mgr.Add(updateHandler); err != nil {
		loggers.LoggerAPKOperator.Errorf("Failed to add status update handler %v", err)
	}
	operatorDataStore := synchronizer.GetOperatorDataStore()
	if err := dpcontrollers.NewGatewayController(mgr, operatorDataStore, updateHandler, &gatewaych); err != nil {
		loggers.LoggerAPKOperator.Errorf("Error creating Gateway controller: %v", err)
	}

	if err := dpcontrollers.NewAPIController(mgr, operatorDataStore, updateHandler, &ch, &successChannel); err != nil {
		loggers.LoggerAPKOperator.Errorf("Error creating API controller: %v", err)
	}

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2602, logging.BLOCKER, "Unable to set up health check: %v", err))
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2603, logging.BLOCKER, "Unable to set up ready check: %v", err))
	}

	go synchronizer.HandleAPILifeCycleEvents(&ch, &successChannel)
	go synchronizer.HandleGatewayLifeCycleEvents(&gatewaych)
	if config.ReadConfigs().PartitionServer.Enabled {
		go synchronizer.SendEventToPartitionServer()
	}

	// if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
	// 	loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2604, logging.BLOCKER, "Problem running manager: %v", err))
	// }

	return &Provider{
		manager: mgr,
		client:  mgr.GetClient(),
	}, nil
}

// Start starts the Provider synchronously until a message is received from ctx.
func (p *Provider) Start(ctx context.Context) error {
	errChan := make(chan error)
	go func() {
		errChan <- p.manager.Start(ctx)
	}()

	// Wait for the manager to exit or an explicit stop.
	select {
	case <-ctx.Done():
		return nil
	case err := <-errChan:
		return err
	}
}
