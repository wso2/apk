/*
 *  Copyright (c) 2022, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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

package operator

import (
	"flag"
	"fmt"
	"strings"

	"github.com/wso2/apk/adapter/config"
	"github.com/wso2/apk/adapter/internal/loggers"

	"github.com/wso2/apk/adapter/pkg/logging"
	"github.com/wso2/apk/adapter/pkg/metrics"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
	gwapiv1a2 "sigs.k8s.io/gateway-api/apis/v1alpha2"

	dpcontrollers "github.com/wso2/apk/adapter/internal/operator/controllers/dp"
	"github.com/wso2/apk/adapter/internal/operator/status"
	"github.com/wso2/apk/adapter/internal/operator/synchronizer"
	"github.com/wso2/apk/adapter/internal/operator/utils"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	cpv1alpha2 "github.com/wso2/apk/common-go-libs/apis/cp/v1alpha2"
	cpv1alpha3 "github.com/wso2/apk/common-go-libs/apis/cp/v1alpha3"
	dpv1alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha1"
	dpv1alpha2 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha2"
	dpv1alpha3 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha3"
	dpv1alpha4 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha4"
	//+kubebuilder:scaffold:imports
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(dpv1alpha1.AddToScheme(scheme))

	utilruntime.Must(gwapiv1.AddToScheme(scheme))

	utilruntime.Must(gwapiv1a2.AddToScheme(scheme))

	utilruntime.Must(dpv1alpha2.AddToScheme(scheme))

	utilruntime.Must(dpv1alpha3.AddToScheme(scheme))

	utilruntime.Must(dpv1alpha4.AddToScheme(scheme))

	utilruntime.Must(cpv1alpha2.AddToScheme(scheme))
	utilruntime.Must(cpv1alpha3.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

// InitOperator starts the Kubernetes gateway operator
func InitOperator(metricsConfig config.Metrics) {
	var enableLeaderElection bool
	var probeAddr string
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()
	log.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	operatorDataStore := synchronizer.GetOperatorDataStore()

	config := config.ReadConfigs()
	options := ctrl.Options{
		Scheme:                  scheme,
		HealthProbeBindAddress:  probeAddr,
		LeaderElection:          true,
		LeaderElectionID:        "operator-lease.apk.wso2.com",
		LeaderElectionNamespace: utils.GetOperatorPodNamespace(),
		// LeaderElectionReleaseOnCancel defines if the leader should step down voluntarily
		// when the Manager ends. This requires the binary to immediately end when the
		// Manager is stopped, otherwise, this setting is unsafe. Setting this significantly
		// speeds up voluntary leader transitions as the new leader don't have to wait
		// LeaseDuration time first.
		//
		// In the default scaffold provided, the program ends immediately after
		// the manager stops, so would be fine to enable this option. However,
		// if you are doing or is intended to do any operation such as perform cleanups
		// after the manager stops then its usage might be unsafe.
		// LeaderElectionReleaseOnCancel: true,
	}
	if !config.Adapter.DeployResourcesWithClusterRoleBindings {
		defaultNamespaces := config.Adapter.Operator.Namespaces
		defaultNamespaces = append(defaultNamespaces, utils.GetOperatorPodNamespace())
		defaultNSMap := make(map[string]cache.Config)
		for _, ns := range defaultNamespaces {
			defaultNSMap[ns] = cache.Config{}
		}
		options.Cache = cache.Options{
			DefaultNamespaces: defaultNSMap,
		}
	}

	if metricsConfig.Enabled {
		options.Metrics.BindAddress = fmt.Sprintf(":%d", metricsConfig.Port)
		// Register the metrics collector
		if strings.EqualFold(metricsConfig.Type, metrics.PrometheusMetricType) {
			loggers.LoggerAPKOperator.Info("Registering Prometheus metrics collector.")
			metrics.RegisterPrometheusCollector()
		}
	} else {
		options.Metrics.BindAddress = "0"
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), options)

	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2600, logging.BLOCKER, "Unable to start manager: %v", err))
	}

	// TODO: Decide on a buffer size and add to config.
	ch := make(chan *synchronizer.APIEvent, 10)
	successChannel := make(chan synchronizer.SuccessEvent, 10)

	gatewaych := make(chan synchronizer.GatewayEvent, 10)
	updateHandler := status.NewUpdateHandler(mgr.GetClient())
	if err := mgr.Add(updateHandler); err != nil {
		loggers.LoggerAPKOperator.Errorf("Failed to add status update handler %v", err)
	}

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
	if config.PartitionServer.Enabled {
		go synchronizer.SendEventToPartitionServer()
	}

	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2604, logging.BLOCKER, "Problem running manager: %v", err))
	}
}
