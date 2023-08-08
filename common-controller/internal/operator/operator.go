/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package operator

import (
	"flag"
	"os"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.

	"github.com/wso2/apk/adapter/pkg/logging"
	cache "github.com/wso2/apk/common-controller/internal/cache"
	"github.com/wso2/apk/common-controller/internal/loggers"
	dpv1alpha1 "github.com/wso2/apk/common-controller/internal/operator/api/v1alpha1"
	dpcontrollers "github.com/wso2/apk/common-controller/internal/operator/controller"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	gwapiv1b1 "sigs.k8s.io/gateway-api/apis/v1beta1"
	//+kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(gwapiv1b1.AddToScheme(scheme))
	utilruntime.Must(dpv1alpha1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

// InitOperator initializes the operator
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

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))
	ratelimitStore := cache.CreateNewOperatorDataStore()
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
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
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

	if err = (&dpv1alpha1.InterceptorService{}).SetupWebhookWithManager(mgr); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2652, err))
	}

	if err = (&dpv1alpha1.Backend{}).SetupWebhookWithManager(mgr); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2655, err))
	}
	if err := dpcontrollers.NewratelimitController(mgr, ratelimitStore); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(3114, err))
	}
	//+kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
