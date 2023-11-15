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
	"github.com/wso2/apk/common-controller/internal/server"
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

	cpv1alpha2 "github.com/wso2/apk/common-controller/internal/operator/apis/cp/v1alpha2"
	constants "github.com/wso2/apk/common-controller/internal/operator/constant"
	"github.com/wso2/apk/common-controller/internal/utils"
)

// ApplicationMappingReconciler reconciles a ApplicationMapping object
type ApplicationMappingReconciler struct {
	client client.Client
	Scheme *runtime.Scheme
	ods    *cache.SubscriptionDataStore
}

// NewApplicationMappingController creates a new Application and Subscription mapping (i.e. ApplicationMapping) controller instance
func NewApplicationMappingController(mgr manager.Manager, subscriptionStore *cache.SubscriptionDataStore) error {
	r := &ApplicationMappingReconciler{
		client: mgr.GetClient(),
		ods:    subscriptionStore,
	}
	c, err := controller.New(constants.ApplicationMappingController, mgr, controller.Options{Reconciler: r})
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2610, logging.BLOCKER, "Error creating ApplicationMapping controller: %v", err.Error()))
		return err
	}

	if err := c.Watch(source.Kind(mgr.GetCache(), &cpv1alpha2.ApplicationMapping{}), &handler.EnqueueRequestForObject{},
		predicate.NewPredicateFuncs(utils.FilterByNamespaces([]string{utils.GetOperatorPodNamespace()}))); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2611, logging.BLOCKER, "Error watching ApplicationMapping resources: %v", err.Error()))
		return err
	}

	loggers.LoggerAPKOperator.Debug("ApplicationMapping Controller successfully started. Watching ApplicationMapping Objects...")
	return nil
}

//+kubebuilder:rbac:groups=cp.wso2.com,resources=applicationmappings,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cp.wso2.com,resources=applicationmappings/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cp.wso2.com,resources=applicationmappings/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ApplicationMapping object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.4/pkg/reconcile
func (r *ApplicationMappingReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	applicationMappingKey := req.NamespacedName
	var applicationMappingList = new(cpv1alpha2.ApplicationMappingList)

	loggers.LoggerAPKOperator.Debugf("Reconciling application mapping: %v", applicationMappingKey.String())
	if err := r.client.List(ctx, applicationMappingList); err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to get application mappings %s/%s",
			applicationMappingKey.Namespace, applicationMappingKey.Name)
	}
	sendUpdates(applicationMappingList)
	var applicationMapping cpv1alpha2.ApplicationMapping
	if err := r.client.Get(ctx, req.NamespacedName, &applicationMapping); err != nil {
		if k8error.IsNotFound(err) {
			applicationMapping, found := r.ods.GetApplicationMappingFromStore(applicationMappingKey)
			if found {
				utils.SendDeleteApplicationMappingEvent(applicationMappingKey.Name, applicationMapping)
				r.ods.DeleteApplicationMappingFromStore(applicationMappingKey)
			} else {
				loggers.LoggerAPKOperator.Debugf("Application mapping %s/%s not found. Ignoring since object must be deleted", applicationMappingKey.Namespace, applicationMappingKey.Name)
			}
		}
	} else {
		utils.SendCreateApplicationMappingEvent(applicationMapping)
		r.ods.AddorUpdateApplicationMappingToStore(applicationMappingKey, applicationMapping.Spec)
	}
	return ctrl.Result{}, nil
}

func sendUpdates(applicationMappingList *cpv1alpha2.ApplicationMappingList) {
	appMappingList := marshalApplicationMappingList(applicationMappingList.Items)
	server.AddApplicationMapping(appMappingList)
}

func marshalApplicationMappingList(applicationMappingList []cpv1alpha2.ApplicationMapping) server.ApplicationMappingList {
	applicationMappings := []server.ApplicationMapping{}
	for _, appMappingInternal := range applicationMappingList {
		appMapping := server.ApplicationMapping{
			UUID:            appMappingInternal.Name,
			ApplicationRef:  appMappingInternal.Spec.ApplicationRef,
			SubscriptionRef: appMappingInternal.Spec.SubscriptionRef,
		}
		applicationMappings = append(applicationMappings, appMapping)
	}
	return server.ApplicationMappingList{
		List: applicationMappings,
	}
}
