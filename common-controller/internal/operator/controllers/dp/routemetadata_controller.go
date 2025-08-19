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

package dp

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/wso2/apk/adapter/pkg/logging"
	cache "github.com/wso2/apk/common-controller/internal/cache"
	"github.com/wso2/apk/common-controller/internal/config"
	loggers "github.com/wso2/apk/common-controller/internal/loggers"
	"github.com/wso2/apk/common-controller/internal/utils"
	"github.com/wso2/apk/common-go-libs/constants"

	dpV2alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v2alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RouteMetadataReconciler reconciles a RouteMetadata object
type RouteMetadataReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Store  *cache.RouteMetadataDataStore
}

// NewRouteMetadataController creates a new controller for RouteMetadata.
func NewRouteMetadataController(mgr manager.Manager, store *cache.RouteMetadataDataStore) error {
	reconciler := &RouteMetadataReconciler{
		Client: mgr.GetClient(),
		Store:  store,
	}

	c, err := controller.New(constants.RouteMetadataController, mgr, controller.Options{Reconciler: reconciler})
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2664, logging.BLOCKER,
			"Error creating RouteMetadata controller: %v", err.Error()))
		return err
	}

	conf := config.ReadConfigs()
	predicateRouteMetadata := []predicate.TypedPredicate[*dpV2alpha1.RouteMetadata]{
		predicate.NewTypedPredicateFuncs[*dpV2alpha1.RouteMetadata](
			utils.FilterRouteMetadataByNamespaces(conf.CommonController.Operator.Namespaces),
		),
	}

	if err := c.Watch(source.Kind(mgr.GetCache(), &dpV2alpha1.RouteMetadata{}, &handler.TypedEnqueueRequestForObject[*dpV2alpha1.RouteMetadata]{}, predicateRouteMetadata...)); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2638, logging.BLOCKER,
			"Error watching RouteMetadata resources: %v", err.Error()))
		return err
	}

	loggers.LoggerAPKOperator.Debug("RouteMetadata Controller successfully started. Watching RouteMetadata Objects...")
	return nil
}

//+kubebuilder:rbac:groups=dp.wso2.com,resources=routepolicies,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=dp.wso2.com,resources=routepolicies/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=dp.wso2.com,resources=routepolicies/finalizers,verbs=update

// Reconcile reconciles the RouteMetadata CR
func (r *RouteMetadataReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	loggers.LoggerAPKOperator.Infof("Reconciling RouteMetadata: %s", req.NamespacedName)
	routePolicyKey := req.NamespacedName

	var routePolicy dpV2alpha1.RouteMetadata
	if err := r.Client.Get(ctx, routePolicyKey, &routePolicy); err != nil {
		loggers.LoggerAPKOperator.Warnf("RouteMetadata %s not found, might be deleted", routePolicyKey)
		r.Store.DeleteRouteMetadata(routePolicyKey.Namespace, routePolicyKey.Name)
		routePolicy.ObjectMeta = metav1.ObjectMeta{
			Namespace: routePolicyKey.Namespace,
			Name:      routePolicyKey.Name,
		}
		routePolicyString, err := utils.ToJSONString(routePolicy)
		if err != nil {
			loggers.LoggerAPKOperator.Errorf("Error converting RouteMetadata to JSON: %v", err)
		} else {
			utils.SendRouteMetadataDeletedEvent(routePolicyString)
			loggers.LoggerAPKOperator.Debugf("Deleted RouteMetadata JSON: %s", routePolicyString)
		}
		// utils.SendRouteMetadataDeletedEvent()
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Add or update the RouteMetadata in the store
	r.Store.AddOrUpdateRouteMetadata(routePolicy)

	routePolicyString, err := utils.ToJSONString(routePolicy)
	if err != nil {
		loggers.LoggerAPKOperator.Errorf("Error converting RouteMetadata to JSON: %v", err)
	} else {
		utils.SendRouteMetadataCreatedOrUpdatedEvent(routePolicyString)
		loggers.LoggerAPKOperator.Debugf("Deleted RouteMetadata JSON: %s", routePolicyString)
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RouteMetadataReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dpV2alpha1.RouteMetadata{}).
		Complete(r)
}
