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

package operator

import (
	"flag"
	"os"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.

	"github.com/google/uuid"
	"github.com/wso2/apk/adapter/pkg/logging"
	cache "github.com/wso2/apk/common-controller/internal/cache"
	"github.com/wso2/apk/common-controller/internal/config"
	"github.com/wso2/apk/common-controller/internal/controlplane"
	"github.com/wso2/apk/common-controller/internal/loggers"
	cpcontrollers "github.com/wso2/apk/common-controller/internal/operator/controllers/cp"
	dpcontrollers "github.com/wso2/apk/common-controller/internal/operator/controllers/dp"
	cpv1alpha2 "github.com/wso2/apk/common-go-libs/apis/cp/v1alpha2"
	dpv1alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha1"
	dpv1alpha2 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha2"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	gwapiv1b1 "sigs.k8s.io/gateway-api/apis/v1beta1"
	//+kubebuilder:scaffold:imports
	"github.com/wso2/apk/common-controller/internal/operator/status"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(gwapiv1b1.AddToScheme(scheme))
	utilruntime.Must(dpv1alpha1.AddToScheme(scheme))
	utilruntime.Must(dpv1alpha2.AddToScheme(scheme))
	utilruntime.Must(cpv1alpha2.AddToScheme(scheme))
	utilruntime.Must(cpv1alpha2.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

// InitOperator initializes the operator
func InitOperator() {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	controlPlaneID := uuid.New().String()
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

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))
	ratelimitStore := cache.CreateNewOperatorDataStore()
	subscriptionStore := cache.CreateNewSubscriptionDataStore()
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: probeAddr,
		// LeaderElection:         true,
		// LeaderElectionID:       "operator-lease.apk.wso2.com",
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
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2600, logging.BLOCKER, "Unable to start manager: %v", err))
		os.Exit(1)
	}

	if err = (&dpv1alpha1.API{}).SetupWebhookWithManager(mgr); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2601, logging.MAJOR,
			"Unable to create webhook API, error: %v", err))
	}

	if err = (&dpv1alpha2.API{}).SetupWebhookWithManager(mgr); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2601, logging.MAJOR,
			"Unable to create webhook API, error: %v", err))
	}

	if err = (&dpv1alpha1.RateLimitPolicy{}).SetupWebhookWithManager(mgr); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2637, logging.MAJOR,
			"Unable to create webhook for Ratelimit, error: %v", err))
	}

	if err = (&dpv1alpha2.APIPolicy{}).SetupWebhookWithManager(mgr); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2638, logging.MAJOR,
			"Unable to create webhook for APIPolicy, error: %v", err))
	}

	if err = (&dpv1alpha2.Authentication{}).SetupWebhookWithManager(mgr); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2638, logging.MAJOR,
			"Unable to create webhook for Authentication, error: %v", err))
	}

	if err = (&dpv1alpha1.InterceptorService{}).SetupWebhookWithManager(mgr); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2652, logging.MAJOR,
			"Unable to create webhook for InterceptorService, error: %v", err))
	}

	if err = (&dpv1alpha1.Backend{}).SetupWebhookWithManager(mgr); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2655, logging.MAJOR,
			"Unable to create webhook for Backend, error: %v", err))
	}
	if err := dpcontrollers.NewratelimitController(mgr, ratelimitStore); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error3114, logging.MAJOR,
			"Error creating JWT Issuer controller, error: %v", err))
	}
	if err := cpcontrollers.NewApplicationController(mgr, subscriptionStore); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error3115, logging.MAJOR,
			"Error creating Application controller, error: %v", err))
	}
	if err := cpcontrollers.NewSubscriptionController(mgr, subscriptionStore); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error3116, logging.MAJOR,
			"Error creating Subscription controller, error: %v", err))
	}
	if err := cpcontrollers.NewApplicationMappingController(mgr, subscriptionStore); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error3117, logging.MAJOR,
			"Error creating Application Mapping controller, error: %v", err))
	}

	updateHandler := status.NewUpdateHandler(mgr.GetClient())
	if err := mgr.Add(updateHandler); err != nil {
		loggers.LoggerAPKOperator.Errorf("Failed to add status update handler %v", err)
	}
	if err := dpcontrollers.NewGatewayClassController(mgr, updateHandler); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error3114, logging.MAJOR,
			"Error creating GatewayClass controller, error: %v", err))
	}

	//+kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2602, logging.BLOCKER, "Unable to set up health check: %v", err))
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2603, logging.BLOCKER, "Unable to set up ready check: %v", err))
		os.Exit(1)
	}
	config := config.ReadConfigs()
	if config.CommonController.ControlPlane.Enabled {
		go func() {
			var controlPlane controlplane.ArtifactDeployer
			if config.CommonController.ControlPlane.Persistence.Type == "K8s" {
				controlPlane = controlplane.NewK8sArtifactDeployer(mgr)

			}
			grpcClient := controlplane.NewControlPlaneAgent(config.CommonController.ControlPlane.Host, config.CommonController.ControlPlane.EventPort, controlPlaneID, controlPlane)
			if grpcClient != nil {
				grpcClient.StartEventStreaming()
			}
		}()
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2604, logging.BLOCKER, "Problem running manager: %v", err))
		os.Exit(1)
	}
}
