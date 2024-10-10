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
	"github.com/wso2/apk/common-controller/internal/config"
	"github.com/wso2/apk/common-controller/internal/loggers"
	"github.com/wso2/apk/common-controller/internal/server"
	k8error "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	k8client "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/wso2/apk/common-controller/internal/utils"
	cpv1alpha2 "github.com/wso2/apk/common-go-libs/apis/cp/v1alpha2"
	cpv1alpha3 "github.com/wso2/apk/common-go-libs/apis/cp/v1alpha3"

	"github.com/wso2/apk/common-go-libs/constants"
)

// ApplicationMappingReconciler reconciles a ApplicationMapping object
type ApplicationMappingReconciler struct {
	client k8client.Client
	Scheme *runtime.Scheme
	ods    *cache.SubscriptionDataStore
}

const (
	applicationIndex  = "applicationIndex"
	subscriptionIndex = "subscriptionIndex"
)

// NewApplicationMappingController creates a new Application and Subscription mapping (i.e. ApplicationMapping) controller instance
func NewApplicationMappingController(mgr manager.Manager, subscriptionStore *cache.SubscriptionDataStore) error {
	r := &ApplicationMappingReconciler{
		client: mgr.GetClient(),
		ods:    subscriptionStore,
	}
	ctx := context.Background()
	conf := config.ReadConfigs()

	if err := addApplicationMappingControllerIndexes(ctx, mgr); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2658, logging.CRITICAL, "Error adding indexes: %v", err))
		return err
	}
	c, err := controller.New(constants.ApplicationMappingController, mgr, controller.Options{Reconciler: r})
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2610, logging.BLOCKER, "Error creating ApplicationMapping controller: %v", err.Error()))
		return err
	}

	if err := c.Watch(source.Kind(mgr.GetCache(), &cpv1alpha2.ApplicationMapping{}, &handler.TypedEnqueueRequestForObject[*cpv1alpha2.ApplicationMapping]{},
		predicate.NewTypedPredicateFuncs(utils.FilterAppMappingByNamespaces([]string{utils.GetOperatorPodNamespace()})))); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2611, logging.BLOCKER, "Error watching ApplicationMapping resources: %v", err.Error()))
		return err
	}

	predicateApp := []predicate.TypedPredicate[*cpv1alpha2.Application]{predicate.NewTypedPredicateFuncs(utils.FilterAppByNamespaces(conf.CommonController.Operator.Namespaces))}
	if err := c.Watch(source.Kind(mgr.GetCache(), &cpv1alpha2.Application{}, handler.TypedEnqueueRequestsFromMapFunc(r.getApplicationMappingsForApplication),
		predicateApp...)); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2613, logging.BLOCKER, "Error watching Application resources: %v", err))
		return err
	}

	predicateSubs := []predicate.TypedPredicate[*cpv1alpha3.Subscription]{predicate.NewTypedPredicateFuncs(utils.FilterSubsByNamespaces(conf.CommonController.Operator.Namespaces))}
	if err := c.Watch(source.Kind(mgr.GetCache(), &cpv1alpha3.Subscription{}, handler.TypedEnqueueRequestsFromMapFunc(r.getApplicationMappingsForSubscription),
		predicateSubs...)); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2613, logging.BLOCKER, "Error watching Subscription resources: %v", err))
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

	loggers.LoggerAPKOperator.Debugf("Reconciling application mapping: %v", applicationMappingKey.String())
	var applicationMapping cpv1alpha2.ApplicationMapping
	if err := r.client.Get(ctx, req.NamespacedName, &applicationMapping); err != nil {
		if k8error.IsNotFound(err) {
			loggers.LoggerAPKOperator.Debugf("Application mapping %s/%s not found in k8s", applicationMappingKey.Namespace, applicationMappingKey.Name)
			applicationMapping, found := r.ods.GetApplicationMappingFromStore(applicationMappingKey)
			if found {
				loggers.LoggerAPKOperator.Debugf("Application mapping %s/%s found in operator data store. Deleting from operator data store and sending delete event to server", applicationMappingKey.Namespace, applicationMappingKey.Name)
				resolvedApplicationMapping := server.GetApplicationMappingFromStore(applicationMappingKey.Name)
				utils.SendDeleteApplicationMappingEvent(applicationMappingKey.Name, applicationMapping, resolvedApplicationMapping.OrganizationID)
				r.ods.DeleteApplicationMappingFromStore(applicationMappingKey)
				server.DeleteApplicationMapping(applicationMappingKey.Name)
			} else {
				loggers.LoggerAPKOperator.Debugf("Application mapping %s/%s not found. Ignoring since object must be deleted", applicationMappingKey.Namespace, applicationMappingKey.Name)
			}
		}
	} else {
		var application cpv1alpha2.Application
		if err := r.client.Get(ctx, types.NamespacedName{Name: string(applicationMapping.Spec.ApplicationRef), Namespace: applicationMapping.Namespace}, &application); err != nil {
			loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2614, logging.CRITICAL, "Error getting Application: %v", err))
			return ctrl.Result{}, nil
		}
		var subscription cpv1alpha3.Subscription
		if err := r.client.Get(ctx, types.NamespacedName{Name: string(applicationMapping.Spec.SubscriptionRef), Namespace: applicationMapping.Namespace}, &subscription); err != nil {
			loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2615, logging.CRITICAL, "Error getting Subscription: %v", err))
			return ctrl.Result{}, nil
		}
		loggers.LoggerAPKOperator.Debugf("Reconsile completed Application mapping :%v,Subscription %v application : %v", applicationMapping, subscription, application)
		sendUpdates(&applicationMapping, application, subscription)
		utils.SendCreateApplicationMappingEvent(applicationMapping, application, subscription)
		r.ods.AddorUpdateApplicationMappingToStore(applicationMappingKey, applicationMapping.Spec)
	}
	return ctrl.Result{}, nil
}

func sendUpdates(applicationMapping *cpv1alpha2.ApplicationMapping, application cpv1alpha2.Application, subscription cpv1alpha3.Subscription) {
	resolvedApplication := marshalApplication(application)
	appMapping := marshalApplicationMapping(applicationMapping, resolvedApplication)
	server.AddApplicationMapping(appMapping)
}

func marshalApplicationMapping(applicationMapping *cpv1alpha2.ApplicationMapping, application server.Application) server.ApplicationMapping {
	return server.ApplicationMapping{
		UUID:            applicationMapping.Name,
		ApplicationRef:  applicationMapping.Spec.ApplicationRef,
		SubscriptionRef: applicationMapping.Spec.SubscriptionRef,
		OrganizationID:  application.OrganizationID,
	}
}

// addApplicationMappingControllerIndexes adds indexes to the ApplicationMapping controller
func addApplicationMappingControllerIndexes(ctx context.Context, mgr manager.Manager) error {
	// Secret to TokenIssuer indexer
	if err := mgr.GetFieldIndexer().IndexField(ctx, &cpv1alpha2.ApplicationMapping{}, applicationIndex,
		func(rawObj k8client.Object) []string {
			applicationMapping := rawObj.(*cpv1alpha2.ApplicationMapping)
			var application []string
			application = append(application,
				types.NamespacedName{
					Name:      string(applicationMapping.Spec.ApplicationRef),
					Namespace: applicationMapping.Namespace,
				}.String())
			return application
		}); err != nil {
		return err
	}
	err := mgr.GetFieldIndexer().IndexField(ctx, &cpv1alpha2.ApplicationMapping{}, subscriptionIndex,
		func(rawObj k8client.Object) []string {
			applicationMapping := rawObj.(*cpv1alpha2.ApplicationMapping)
			var subscriptions []string
			subscriptions = append(subscriptions,
				types.NamespacedName{
					Name:      string(applicationMapping.Spec.SubscriptionRef),
					Namespace: applicationMapping.Namespace,
				}.String())
			return subscriptions
		})
	return err
}

// getApplicationMappingsForApplication triggers the ApplicationMapping controller reconcile method based on the changes detected
// from Application objects. If the changes are done for an API stored in the Operator Data store,
// a new reconcile event will be created and added to the reconcile event queue.
func (r *ApplicationMappingReconciler) getApplicationMappingsForApplication(ctx context.Context, obj *cpv1alpha2.Application) []reconcile.Request {
	application := obj
	applicationMappingList := &cpv1alpha2.ApplicationMappingList{}
	if err := r.client.List(ctx, applicationMappingList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(applicationIndex, utils.NamespacedName(application).String()),
	}); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2623, logging.CRITICAL, "Unable to find associated application mappings: %s", utils.NamespacedName(application).String()))
		return []reconcile.Request{}
	}

	if len(applicationMappingList.Items) == 0 {
		loggers.LoggerAPKOperator.Debugf("ApplicationMappings for Application %s/%s not found", application.Namespace, application.Name)
		return []reconcile.Request{}
	}

	requests := []reconcile.Request{}
	for _, applicationMapping := range applicationMappingList.Items {
		req := reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      applicationMapping.Name,
				Namespace: applicationMapping.Namespace},
		}
		requests = append(requests, req)
		loggers.LoggerAPKOperator.Debugf("Adding reconcile request for ApplicationMapping: %s/%s with Application UUID: %v", applicationMapping.Namespace, applicationMapping.Name,
			string(applicationMapping.ObjectMeta.UID))
	}
	return requests
}

// getApplicationMappingsForSubscription triggers the ApplicationMapping controller reconcile method based on the changes detected
// from Subscription objects. If the changes are done for an API stored in the Operator Data store,
func (r *ApplicationMappingReconciler) getApplicationMappingsForSubscription(ctx context.Context, obj *cpv1alpha3.Subscription) []reconcile.Request {
	subscription := obj
	applicationMappingList := &cpv1alpha2.ApplicationMappingList{}
	if err := r.client.List(ctx, applicationMappingList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(subscriptionIndex, utils.NamespacedName(subscription).String()),
	}); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2623, logging.CRITICAL, "Unable to find associated Application mappings: %s", utils.NamespacedName(subscription).String()))
		return []reconcile.Request{}
	}

	if len(applicationMappingList.Items) == 0 {
		loggers.LoggerAPKOperator.Debugf("ApplicationMappings for Subscription %s/%s not found", subscription.Namespace, subscription.Name)
		return []reconcile.Request{}
	}

	requests := []reconcile.Request{}
	for _, applicationMapping := range applicationMappingList.Items {
		req := reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      applicationMapping.Name,
				Namespace: applicationMapping.Namespace},
		}
		requests = append(requests, req)
		loggers.LoggerAPKOperator.Debugf("Adding reconcile request for ApplicationMapping: %s/%s with Subscription UUID: %v", applicationMapping.Namespace, applicationMapping.Name,
			string(applicationMapping.ObjectMeta.UID))
	}
	return requests
}
