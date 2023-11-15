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

package cp

import (
	"context"
	"fmt"

	"github.com/wso2/apk/adapter/pkg/logging"
	"github.com/wso2/apk/common-controller/internal/cache"
	"github.com/wso2/apk/common-controller/internal/loggers"
	cpv1alpha2 "github.com/wso2/apk/common-controller/internal/operator/apis/cp/v1alpha2"
	constants "github.com/wso2/apk/common-controller/internal/operator/constant"
	"github.com/wso2/apk/common-controller/internal/server"
	"github.com/wso2/apk/common-controller/internal/utils"
	k8error "k8s.io/apimachinery/pkg/api/errors"
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
)

// ApplicationReconciler reconciles a Application object
type ApplicationReconciler struct {
	client client.Client
	Scheme *runtime.Scheme
	ods    *cache.SubscriptionDataStore
}

// NewApplicationController creates a new Application controller instance
func NewApplicationController(mgr manager.Manager, subscriptionStore *cache.SubscriptionDataStore) error {
	r := &ApplicationReconciler{
		client: mgr.GetClient(),
		ods:    subscriptionStore,
	}
	c, err := controller.New(constants.ApplicationController, mgr, controller.Options{Reconciler: r})
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2606, logging.BLOCKER, "Error creating Application controller: %v", err.Error()))
		return err
	}

	if err := c.Watch(source.Kind(mgr.GetCache(), &cpv1alpha2.Application{}), &handler.EnqueueRequestForObject{},
		predicate.NewPredicateFuncs(utils.FilterByNamespaces([]string{utils.GetOperatorPodNamespace()}))); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2607, logging.BLOCKER, "Error watching Application resources: %v", err.Error()))
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
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Application object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.4/pkg/reconcile
func (applicationReconciler *ApplicationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	applicationKey := req.NamespacedName
	var applicationList = new(cpv1alpha2.ApplicationList)

	loggers.LoggerAPKOperator.Infof("Reconciling application: %v", applicationKey.String())
	if err := applicationReconciler.client.List(ctx, applicationList); err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to get applications %s/%s",
			applicationKey.Namespace, applicationKey.Name)
	}
	var application cpv1alpha2.Application
	if err := applicationReconciler.client.Get(ctx, req.NamespacedName, &application); err != nil {
		if k8error.IsNotFound(err) {
			applicationSpec, found := applicationReconciler.ods.GetApplicationFromStore(applicationKey)
			loggers.LoggerAPKOperator.Infof("Application cr not available in k8s")
			loggers.LoggerAPKOperator.Infof("cached Application spec: %v,%v", applicationSpec, found)
			if found {
				utils.SendAppDeletionEvent(applicationKey.Name, applicationSpec)
				applicationReconciler.ods.DeleteApplicationFromStore(applicationKey)
			} else {
				loggers.LoggerAPKOperator.Infof("Application %s/%s does not exist in k8s", applicationKey.Namespace, applicationKey.Name)
			}
		}
	} else {
		loggers.LoggerAPKOperator.Infof("Application cr available in k8s")
		applicationSpec, found := applicationReconciler.ods.GetApplicationFromStore(applicationKey)
		if found {
			// update
			loggers.LoggerAPKOperator.Infof("Application in ods")
			utils.SendAppUpdateEvent(applicationKey.Name, applicationSpec, application.Spec)

		} else {
			loggers.LoggerAPKOperator.Infof("Application in ods consider as update")
			utils.SendAddApplicationEvent(application)
		}
		applicationReconciler.ods.AddorUpdateApplicationToStore(applicationKey, application.Spec)
	}
	sendAppUpdates(applicationList)
	return ctrl.Result{}, nil
}

func sendAppUpdates(applicationList *cpv1alpha2.ApplicationList) {
	appList := marshalApplicationList(applicationList.Items)
	server.AddApplication(appList)
	appKeyMappingList := marshalApplicationKeyMapping(applicationList.Items)
	server.AddApplicationKeyMapping(appKeyMappingList)
}

func marshalApplicationList(applicationList []cpv1alpha2.Application) server.ApplicationList {
	applications := []server.Application{}
	for _, appInternal := range applicationList {
		app := server.Application{
			UUID:           appInternal.Name,
			Name:           appInternal.Spec.Name,
			Owner:          appInternal.Spec.Owner,
			OrganizationID: appInternal.Spec.Organization,
			Attributes:     appInternal.Spec.Attributes,
		}
		applications = append(applications, app)
	}
	return server.ApplicationList{
		List: applications,
	}
}

func marshalApplicationKeyMapping(applicationList []cpv1alpha2.Application) server.ApplicationKeyMappingList {
	applicationKeyMappings := []server.ApplicationKeyMapping{}
	for _, appInternal := range applicationList {
		var oauth2SecurityScheme = appInternal.Spec.SecuritySchemes.OAuth2
		if oauth2SecurityScheme != nil {
			for _, env := range oauth2SecurityScheme.Environments {
				appIdentifier := server.ApplicationKeyMapping{
					ApplicationUUID:       appInternal.Name,
					SecurityScheme:        constants.OAuth2,
					ApplicationIdentifier: env.AppID,
					KeyType:               env.KeyType,
					EnvID:                 env.EnvID,
				}
				applicationKeyMappings = append(applicationKeyMappings, appIdentifier)
			}
		}
	}
	return server.ApplicationKeyMappingList{
		List: applicationKeyMappings,
	}
}
