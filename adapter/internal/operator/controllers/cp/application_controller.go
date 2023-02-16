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

package cp

import (
	"context"
	"fmt"

	"github.com/wso2/apk/adapter/internal/discovery/xds"
	"github.com/wso2/apk/adapter/internal/loggers"
	"github.com/wso2/apk/adapter/pkg/logging"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	cpv1alpha1 "github.com/wso2/apk/adapter/internal/operator/apis/cp/v1alpha1"
	"github.com/wso2/apk/adapter/internal/operator/constants"
	"github.com/wso2/apk/adapter/internal/operator/utils"
	"github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/subscription"
)

// ApplicationReconciler reconciles a Application object
type ApplicationReconciler struct {
	client client.Client
	Scheme *runtime.Scheme
}

// NewApplicationController creates a new Application controller instance.
func NewApplicationController(mgr manager.Manager) error {
	r := &ApplicationReconciler{
		client: mgr.GetClient(),
	}
	c, err := controller.New(constants.ApplicationController, mgr, controller.Options{Reconciler: r})
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprintf("Error creating Application controller: %v", err.Error()),
			Severity:  logging.BLOCKER,
			ErrorCode: 2801,
		})
		return err
	}

	if err := c.Watch(&source.Kind{Type: &cpv1alpha1.Application{}}, &handler.EnqueueRequestForObject{},
		predicate.NewPredicateFuncs(utils.FilterByNamespaces([]string{utils.GetOperatorPodNamespace()}))); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprintf("Error watching Application resources: %v", err.Error()),
			Severity:  logging.BLOCKER,
			ErrorCode: 2802,
		})
		return err
	}

	loggers.LoggerAPKOperator.Debug("Application Controller successfully started. Watching Application Objects...")
	return nil
}

//+kubebuilder:rbac:groups=cp.wso2.com,resources=applications,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cp.wso2.com,resources=applications/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cp.wso2.com,resources=applications/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// Modify the Reconcile function to compare the state specified by
// the Application object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (applicationReconciler *ApplicationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	loggers.LoggerAPKOperator.Debugf("Reconciling application: %v", req.NamespacedName.String())

	applicationKey := req.NamespacedName
	var applicationList = new(cpv1alpha1.ApplicationList)
	if err := applicationReconciler.client.List(ctx, applicationList); err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to get applications %s/%s",
			applicationKey.Namespace, applicationKey.Name)
	}

	sendUpdates(applicationList)
	return ctrl.Result{}, nil
}

func sendUpdates(applicationList *cpv1alpha1.ApplicationList) {
	appList := marshalApplicationList(applicationList.Items)
	xds.UpdateEnforcerApplications(appList)

	appKeyMappingList := marshalApplicationKeyMapping(applicationList.Items)
	xds.UpdateEnforcerApplicationKeyMappings(appKeyMappingList)
}

func marshalApplicationList(applicationList []cpv1alpha1.Application) *subscription.ApplicationList {
	applications := []*subscription.Application{}
	for _, appInternal := range applicationList {
		app := &subscription.Application{
			Uuid:       appInternal.Name,
			Name:       appInternal.Spec.Name,
			Policy:     appInternal.Spec.Policy,
			Attributes: appInternal.Spec.Attributes,
		}
		applications = append(applications, app)
	}
	return &subscription.ApplicationList{
		List: applications,
	}
}

func marshalApplicationKeyMapping(applicationList []cpv1alpha1.Application) *subscription.ApplicationKeyMappingList {
	applicationKeyMappings := []*subscription.ApplicationKeyMapping{}
	for _, appInternal := range applicationList {
		for _, consumerKeyInternal := range appInternal.Spec.Keys {
			consumerKey := &subscription.ApplicationKeyMapping{
				ConsumerKey:     consumerKeyInternal.Key,
				KeyManager:      consumerKeyInternal.KeyManager,
				ApplicationUUID: appInternal.Name,
			}
			applicationKeyMappings = append(applicationKeyMappings, consumerKey)
		}
	}
	return &subscription.ApplicationKeyMappingList{
		List: applicationKeyMappings,
	}
}
