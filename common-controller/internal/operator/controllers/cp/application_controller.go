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

	"github.com/wso2/apk/adapter/pkg/logging"
	"github.com/wso2/apk/common-controller/internal/cache"
	"github.com/wso2/apk/common-controller/internal/loggers"
	"github.com/wso2/apk/common-controller/internal/server"
	"github.com/wso2/apk/common-go-libs/pkg/server/model"
	"github.com/wso2/apk/common-controller/internal/utils"
	cpv1alpha2 "github.com/wso2/apk/common-go-libs/apis/cp/v1alpha2"
	"github.com/wso2/apk/common-go-libs/constants"
	k8error "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
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

	if err := c.Watch(source.Kind(mgr.GetCache(), &cpv1alpha2.Application{}, &handler.TypedEnqueueRequestForObject[*cpv1alpha2.Application]{},
		predicate.NewTypedPredicateFuncs(utils.FilterAppByNamespaces([]string{utils.GetOperatorPodNamespace()})))); err != nil {
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

	loggers.LoggerAPKOperator.Debugf("Reconciling application: %v", applicationKey.String())
	var application cpv1alpha2.Application
	if err := applicationReconciler.client.Get(ctx, req.NamespacedName, &application); err != nil {
		if k8error.IsNotFound(err) {
			applicationSpec, found := applicationReconciler.ods.GetApplicationFromStore(applicationKey)
			loggers.LoggerAPKOperator.Debugf("Application cr not available in k8s")
			loggers.LoggerAPKOperator.Debugf("cached Application spec: %v,%v", applicationSpec, found)
			if found {
				utils.SendAppDeletionEvent(applicationKey.Name, applicationSpec)
				applicationReconciler.ods.DeleteApplicationFromStore(applicationKey)
				server.DeleteApplication(applicationKey.Name)
			} else {
				loggers.LoggerAPKOperator.Debugf("Application %s/%s does not exist in k8s", applicationKey.Namespace, applicationKey.Name)
			}
		}
	} else {
		loggers.LoggerAPKOperator.Debugf("Application cr available in k8s")
		applicationSpec, found := applicationReconciler.ods.GetApplicationFromStore(applicationKey)
		if found {
			// update
			loggers.LoggerAPKOperator.Debugf("Application in ods")
			utils.SendAppUpdateEvent(applicationKey.Name, applicationSpec, application.Spec)
		} else {
			loggers.LoggerAPKOperator.Debugf("Application in ods consider as update")
			utils.SendAddApplicationEvent(application)
		}
		applicationReconciler.ods.AddorUpdateApplicationToStore(applicationKey, application.Spec)
		applicationReconciler.sendAppUpdates(application, found)
	}
	return ctrl.Result{}, nil
}

func (applicationReconciler *ApplicationReconciler) sendAppUpdates(application cpv1alpha2.Application, update bool) {
	resolvedApplication := marshalApplication(application)
	if update {
		server.DeleteApplication(application.Name)
	}
	server.AddApplication(resolvedApplication)
	if application.Spec.SecuritySchemes != nil {
		appKeyMappingList := marshalApplicationKeyMapping(application)
		for _, applicationKeyMapping := range appKeyMappingList {
			server.AddApplicationKeyMapping(applicationKeyMapping)
		}
	}
}

func marshalApplication(application cpv1alpha2.Application) model.Application {
	return model.Application{
		UUID:           application.Name,
		Name:           application.Spec.Name,
		Owner:          application.Spec.Owner,
		OrganizationID: application.Spec.Organization,
		Attributes:     application.Spec.Attributes,
	}
}

func marshalApplicationKeyMapping(appInternal cpv1alpha2.Application) []model.ApplicationKeyMapping {
	applicationKeyMappings := []model.ApplicationKeyMapping{}
	var oauth2SecurityScheme = appInternal.Spec.SecuritySchemes.OAuth2
	if oauth2SecurityScheme != nil {
		for _, env := range oauth2SecurityScheme.Environments {
			appIdentifier := model.ApplicationKeyMapping{
				ApplicationUUID:       appInternal.Name,
				SecurityScheme:        constants.OAuth2,
				ApplicationIdentifier: env.AppID,
				KeyType:               env.KeyType,
				EnvID:                 env.EnvID,
				OrganizationID:        appInternal.Spec.Organization,
			}
			applicationKeyMappings = append(applicationKeyMappings, appIdentifier)
		}
	}
	return applicationKeyMappings
}
