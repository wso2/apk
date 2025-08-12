/* Copyright (c) 2025 WSO2 LLC. (http://www.wso2.com) All Rights Reserved.
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package config

import (
	"fmt"
	"log"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"

	egv1a1 "github.com/envoyproxy/gateway/api/v1alpha1"
	dpv2alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v2alpha1"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
	gwapiv1alpha3 "sigs.k8s.io/gateway-api/apis/v1alpha3"
)

var GlobalManager ctrl.Manager

// InitializeManager initializes the Kubernetes controller manager
func InitializeManager() error {
	var probeAddr string
	var scheme = runtime.NewScheme()

	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(dpv2alpha1.AddToScheme(scheme))
	utilruntime.Must(egv1a1.AddToScheme(scheme))
	utilruntime.Must(corev1.AddToScheme(scheme))
	utilruntime.Must(gwapiv1alpha3.Install(scheme))
	utilruntime.Must(gwapiv1.Install(scheme))

	options := ctrl.Options{
		Scheme:                 scheme,
		HealthProbeBindAddress: probeAddr,
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), options)
	if err != nil {
		return fmt.Errorf("unable to start kubernetes controller manager: %w", err)
	}

	GlobalManager = mgr
	return nil
}

// GetManager returns the global manager instance
func GetManager() ctrl.Manager {
	if GlobalManager == nil {
		log.Fatal("Manager not initialized. Call InitializeManager() first.")
	}
	return GlobalManager
}
