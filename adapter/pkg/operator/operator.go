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

	"github.com/wso2/apk/adapter/config"
	"github.com/wso2/apk/adapter/internal/loggers"
	"github.com/wso2/apk/adapter/internal/management-server/xds"
	"github.com/wso2/apk/adapter/pkg/logging"
	gwapiv1a2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gwapiv1b1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	cpcontrollers "github.com/wso2/apk/adapter/pkg/operator/controllers/cp"
	dpcontrollers "github.com/wso2/apk/adapter/pkg/operator/controllers/dp"
	"github.com/wso2/apk/adapter/pkg/operator/status"
	"github.com/wso2/apk/adapter/pkg/operator/synchronizer"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	cpv1alpha1 "github.com/wso2/apk/adapter/pkg/operator/apis/cp/v1alpha1"
	dpv1alpha1 "github.com/wso2/apk/adapter/pkg/operator/apis/dp/v1alpha1"
	//+kubebuilder:scaffold:imports
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(dpv1alpha1.AddToScheme(scheme))

	utilruntime.Must(gwapiv1b1.AddToScheme(scheme))

	utilruntime.Must(gwapiv1a2.AddToScheme(scheme))

	utilruntime.Must(cpv1alpha1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

// InitOperator starts the Kubernetes gateway operator
func InitOperator() {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	operatorDataStore := synchronizer.CreateNewOperatorDataStore()

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "73c5c496.wso2.com",
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
	})
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2600, err))
	}

	// TODO: Decide on a buffer size and add to config.
	ch := make(chan synchronizer.APIEvent, 10)

	updateHandler := status.NewUpdateHandler(mgr.GetClient())
	if err := mgr.Add(updateHandler); err != nil {
		loggers.LoggerAPKOperator.Errorf("Failed to add status update handler %v", err)
	}

	if err := dpcontrollers.NewAPIController(mgr, operatorDataStore, updateHandler, &ch); err != nil {
		loggers.LoggerAPKOperator.Errorf("Error creating API controller: %v", err)
	}

	if err = (&dpv1alpha1.API{}).SetupWebhookWithManager(mgr); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2601, err))
	}

	if err = (&dpv1alpha1.RateLimitPolicy{}).SetupWebhookWithManager(mgr); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2637, err))
	}

	if err = (&dpv1alpha1.APIPolicy{}).SetupWebhookWithManager(mgr); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2638, err))
	}

	if err := cpcontrollers.NewApplicationController(mgr); err != nil {
		loggers.LoggerAPKOperator.Errorf("Error creating Application controller: %v", err)
	}

	if err := cpcontrollers.NewSubscriptionController(mgr); err != nil {
		loggers.LoggerAPKOperator.Errorf("Error creating Subscription controller: %v", err)
	}

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2602, err))
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2603, err))
	}

	go synchronizer.HandleAPILifeCycleEvents(&ch)
	if config.ReadConfigs().ManagementServer.Enabled {
		go xds.InitApkMgtXDSClient()
		go xds.HandleApplicationEventsFromMgtServer(mgr.GetClient(), mgr.GetAPIReader())
		go synchronizer.SendAPIToAPKMgtServer()
	}

	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2604, err))
	}
}
