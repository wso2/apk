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

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	k8client "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/wso2/apk/adapter/pkg/logging"
	cache "github.com/wso2/apk/common-controller/internal/cache"
	"github.com/wso2/apk/common-controller/internal/config"
	loggers "github.com/wso2/apk/common-controller/internal/loggers"
	"github.com/wso2/apk/common-controller/internal/utils"
	"github.com/wso2/apk/common-go-libs/constants"

	dpV2alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v2alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RouteMetadataReconciler reconciles a RouteMetadata object
type RouteMetadataReconciler struct {
	client client.Client
	Scheme *runtime.Scheme
	Store  *cache.RouteMetadataDataStore
}

const (
	// configMapIndex is the index for ConfigMap resources in the RouteMetadata controller
	configMapIndexRouteMetadata = "ConfigMapIndexRouteMetadata"
)

// NewRouteMetadataController creates a new controller for RouteMetadata.
func NewRouteMetadataController(mgr manager.Manager, store *cache.RouteMetadataDataStore) error {
	reconciler := &RouteMetadataReconciler{
		client: mgr.GetClient(),
		Store:  store,
	}

	ctx := context.Background()
	if err := reconciler.addRouteMetadataIndexes(ctx, mgr); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2612, logging.BLOCKER, "Error adding indexes: %v", err))
		return err
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

	predicateConfigMap := []predicate.TypedPredicate[*corev1.ConfigMap]{predicate.NewTypedPredicateFuncs(utils.FilterConfigMapByNamespaces(conf.CommonController.Operator.Namespaces))}
	if err := c.Watch(source.Kind(mgr.GetCache(), &corev1.ConfigMap{},
		handler.TypedEnqueueRequestsFromMapFunc(reconciler.getRouteMetadataForConfigMap), predicateConfigMap...)); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2613, logging.BLOCKER,
			"Error watching ConfigMap resources: %v", err))
		return err
	}

	loggers.LoggerAPKOperator.Debug("RouteMetadata Controller successfully started. Watching RouteMetadata Objects...")
	return nil
}

//+kubebuilder:rbac:groups=dp.wso2.com,resources=routepolicies,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=dp.wso2.com,resources=routepolicies/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=dp.wso2.com,resources=routepolicies/finalizers,verbs=update

// Reconcile reconciles the RouteMetadata CR
func (routeMetadataReconciler *RouteMetadataReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	loggers.LoggerAPKOperator.Infof("Reconciling RouteMetadata: %s", req.NamespacedName)
	routeMetadataKey := req.NamespacedName

	var routeMetadata dpV2alpha1.RouteMetadata
	if err := routeMetadataReconciler.client.Get(ctx, routeMetadataKey, &routeMetadata); err != nil {
		loggers.LoggerAPKOperator.Warnf("RouteMetadata %s not found, might be deleted", routeMetadataKey)
		routeMetadataReconciler.Store.DeleteRouteMetadata(routeMetadataKey.Namespace, routeMetadataKey.Name)
		routeMetadata.ObjectMeta = metav1.ObjectMeta{
			Namespace: routeMetadataKey.Namespace,
			Name:      routeMetadataKey.Name,
		}
		routePolicyString, err := utils.ToJSONString(routeMetadata)
		if err != nil {
			loggers.LoggerAPKOperator.Errorf("Error converting RouteMetadata to JSON: %v", err)
		} else {
			utils.SendRouteMetadataDeletedEvent(routePolicyString)
			loggers.LoggerAPKOperator.Debugf("Deleted RouteMetadata JSON: %s", routePolicyString)
		}
		// utils.SendRouteMetadataDeletedEvent()
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if routeMetadata.Spec.API.DefinitionFileRef == nil {
		namespacedName := types.NamespacedName{
			Namespace: req.Namespace,
			Name:      string(routeMetadata.Spec.API.DefinitionFileRef.Name),
		}
		var cm corev1.ConfigMap
		if err := routeMetadataReconciler.client.Get(ctx, namespacedName, &cm); err != nil {
			loggers.LoggerAPKOperator.Errorf("failed to fetch ConfigMap %s: %v", namespacedName.String(), err)
		}
		val, ok := cm.Data["Definition"]
		if !ok {
			loggers.LoggerAPKOperator.Warnf("key %s not found in ConfigMap %s", "Definition", namespacedName.String())
		}
		routeMetadata.Spec.API.Definition = val
	}

	// Add or update the RouteMetadata in the store
	routeMetadataReconciler.Store.AddOrUpdateRouteMetadata(routeMetadata)

	routePolicyString, err := utils.ToJSONString(routeMetadata)
	if err != nil {
		loggers.LoggerAPKOperator.Errorf("Error converting RouteMetadata to JSON: %v", err)
	} else {
		utils.SendRouteMetadataCreatedOrUpdatedEvent(routePolicyString)
		loggers.LoggerAPKOperator.Debugf("Deleted RouteMetadata JSON: %s", routePolicyString)
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (routeMetadataReconciler *RouteMetadataReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dpV2alpha1.RouteMetadata{}).
		Complete(routeMetadataReconciler)
}

// getRouteMetadataForConfigMap returns a list of reconcile requests for RouteMetadata objects
func (routeMetadataReconciler *RouteMetadataReconciler) getRouteMetadataForConfigMap(ctx context.Context, obj *corev1.ConfigMap) []reconcile.Request {
	configMap := obj

	requests := []reconcile.Request{}

	routeMetadataList := &dpV2alpha1.RouteMetadataList{}
	if err := routeMetadataReconciler.client.List(ctx, routeMetadataList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(configMapIndex, NamespacedName(configMap).String()),
	}); err != nil {
		return []reconcile.Request{}
	}

	for item := range routeMetadataList.Items {
		routePolicy := routeMetadataList.Items[item]
		requests = append(requests, routeMetadataReconciler.AddRouteMetadataRequest(&routePolicy)...)
	}

	return requests
}

// AddRouteMetadataRequest adds a reconcile request for the given RouteMetadata
func (routeMetadataReconciler *RouteMetadataReconciler) AddRouteMetadataRequest(routePolicy *dpV2alpha1.RouteMetadata) []reconcile.Request {
	routePolicyKey := client.ObjectKey{
		Namespace: routePolicy.Namespace,
		Name:      routePolicy.Name,
	}

	loggers.LoggerAPKOperator.Debugf("Adding RouteMetadata request for %s", routePolicyKey.String())
	return []reconcile.Request{{NamespacedName: routePolicyKey}}
}

func (routeMetadataReconciler *RouteMetadataReconciler) addRouteMetadataIndexes(ctx context.Context, mgr manager.Manager) error {
	// Index by referenced ConfigMaps
	if err := mgr.GetFieldIndexer().IndexField(ctx, &dpV2alpha1.RouteMetadata{}, configMapIndex,
		func(rawObj k8client.Object) []string {
			routeMetadata := rawObj.(*dpV2alpha1.RouteMetadata)
			namespacedConfigMaps := make([]string, 0)
			if routeMetadata.Spec.API.DefinitionFileRef != nil {
				namespacedConfigMaps = append(namespacedConfigMaps, types.NamespacedName{
					Namespace: routeMetadata.Namespace,
					Name:      string(routeMetadata.Spec.API.DefinitionFileRef.Name),
				}.String())
			}

			return namespacedConfigMaps
		}); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2614, logging.BLOCKER,
			"Error adding index for RouteMetadata: %v", err))
		return err
	}

	return nil
}