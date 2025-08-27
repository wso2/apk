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
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"k8s.io/apimachinery/pkg/fields"
	k8client "sigs.k8s.io/controller-runtime/pkg/client"
	"k8s.io/apimachinery/pkg/types"
)

// RoutePolicyReconciler reconciles a RoutePolicy object
type RoutePolicyReconciler struct {
	client client.Client
	Scheme *runtime.Scheme
	Store  *cache.RoutePolicyDataStore
}

const (
	configMapIndex = "configmapIndex"
	secretIndex    = "secretIndex"
)

// NewRoutePolicyController creates a new controller for RoutePolicy.
func NewRoutePolicyController(mgr manager.Manager, store *cache.RoutePolicyDataStore) error {
	reconciler := &RoutePolicyReconciler{
		client: mgr.GetClient(),
		Store:  store,
	}

	ctx := context.Background()
	if err := reconciler.addRoutePolicyIndexes(ctx, mgr); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2612, logging.BLOCKER, "Error adding indexes: %v", err))
		return err
	}

	c, err := controller.New(constants.RoutePolicyController, mgr, controller.Options{Reconciler: reconciler})
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2664, logging.BLOCKER,
			"Error creating RoutePolicy controller: %v", err.Error()))
		return err
	}

	conf := config.ReadConfigs()
	predicateRoutePolicy := []predicate.TypedPredicate[*dpV2alpha1.RoutePolicy]{
		predicate.NewTypedPredicateFuncs[*dpV2alpha1.RoutePolicy](
			utils.FilterRoutePolicyByNamespaces(conf.CommonController.Operator.Namespaces),
		),
	}

	if err := c.Watch(source.Kind(mgr.GetCache(), &dpV2alpha1.RoutePolicy{}, &handler.TypedEnqueueRequestForObject[*dpV2alpha1.RoutePolicy]{}, predicateRoutePolicy...)); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2638, logging.BLOCKER,
			"Error watching RoutePolicy resources: %v", err.Error()))
		return err
	}

	predicateConfigMap := []predicate.TypedPredicate[*corev1.ConfigMap]{predicate.NewTypedPredicateFuncs(utils.FilterConfigMapByNamespaces(conf.CommonController.Operator.Namespaces))}
	if err := c.Watch(source.Kind(mgr.GetCache(), &corev1.ConfigMap{},
		handler.TypedEnqueueRequestsFromMapFunc(reconciler.getRoutePolicyForConfigMap), predicateConfigMap...)); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2613, logging.BLOCKER,
			"Error watching ConfigMap resources: %v", err))
		return err
	}

	predicateSecret := []predicate.TypedPredicate[*corev1.Secret]{predicate.NewTypedPredicateFuncs(utils.FilterSecretByNamespaces(conf.CommonController.Operator.Namespaces))}
	if err := c.Watch(source.Kind(mgr.GetCache(), &corev1.Secret{},
		handler.TypedEnqueueRequestsFromMapFunc(reconciler.getRoutePolicyForSecret), predicateSecret...)); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2613, logging.BLOCKER,
			"Error watching Secret resources: %v", err))
		return err
	}

	loggers.LoggerAPKOperator.Debug("RoutePolicy Controller successfully started. Watching RoutePolicy Objects...")
	return nil
}

//+kubebuilder:rbac:groups=dp.wso2.com,resources=routepolicies,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=dp.wso2.com,resources=routepolicies/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=dp.wso2.com,resources=routepolicies/finalizers,verbs=update

// Reconcile reconciles the RoutePolicy CR
func (routePolicyReconciler *RoutePolicyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	loggers.LoggerAPKOperator.Infof("Reconciling RoutePolicy: %s", req.NamespacedName)
	routePolicyKey := req.NamespacedName

	var routePolicy dpV2alpha1.RoutePolicy
	if err := routePolicyReconciler.client.Get(ctx, routePolicyKey, &routePolicy); err != nil {
		loggers.LoggerAPKOperator.Warnf("RoutePolicy %s not found, might be deleted", routePolicyKey)
		routePolicyReconciler.Store.DeleteRoutePolicy(routePolicyKey.Namespace, routePolicyKey.Name)
		routePolicy.ObjectMeta = metav1.ObjectMeta{
			Namespace: routePolicyKey.Namespace,
			Name:      routePolicyKey.Name,
		}
		routePolicyString, err := utils.ToJSONString(routePolicy)
		if err != nil {
			loggers.LoggerAPKOperator.Errorf("Error converting RoutePolicy to JSON: %v", err)
		} else {
			utils.SendRoutePolicyDeletedEvent(routePolicyString)
			loggers.LoggerAPKOperator.Debugf("Deleted RoutePolicy JSON: %s", routePolicyString)
		}
		// utils.SendRoutePolicyDeletedEvent()
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	paramList := make([]*dpV2alpha1.Parameter, 0)
	for i := range routePolicy.Spec.RequestMediation {
		mediation := routePolicy.Spec.RequestMediation[i]
		for j := range mediation.Parameters {
			param := mediation.Parameters[j]
			paramList = append(paramList, param)
		}
	}
	for i := range routePolicy.Spec.ResponseMediation {
		mediation := routePolicy.Spec.ResponseMediation[i]
		for j := range mediation.Parameters {
			param := mediation.Parameters[j]
			paramList = append(paramList, param)
		}
	}
	for _, param := range paramList {
		if param.ValueRef != nil {
			namespacedName := types.NamespacedName{
				Namespace: req.Namespace,
				Name:      string(param.ValueRef.Name),
			}
			value := param.Value
			switch param.ValueRef.Kind {
			case constants.KindConfigMap:
				var cm corev1.ConfigMap
				if err := routePolicyReconciler.client.Get(ctx, namespacedName, &cm); err != nil {
					loggers.LoggerAPKOperator.Errorf("failed to fetch ConfigMap %s: %v", namespacedName.String(), err)
					continue
				}
				if v, ok := getValueFromMap(cm.Data, param.Key, namespacedName.String()); ok {
					value = v
				} else {
					continue
				}
			case constants.KindSecret:
				var secret corev1.Secret
				if err := routePolicyReconciler.client.Get(ctx, namespacedName, &secret); err != nil {
					loggers.LoggerAPKOperator.Errorf("failed to fetch Secret %s: %v", namespacedName.String(), err)
					continue
				}
				if v, ok := getValueFromSecret(secret.Data, param.Key, namespacedName.String()); ok {
					value = v
				} else {
					continue
				}
			default:
				loggers.LoggerAPKOperator.Warnf("unsupported ValueRef Kind: %s", param.ValueRef.Kind)
				continue
			}
			// Set the resolved value
			param.Value = value
		}
	}

	// Add or update the RoutePolicy in the store
	routePolicyReconciler.Store.AddOrUpdateRoutePolicy(routePolicy)

	routePolicyString, err := utils.ToJSONString(routePolicy)
	if err != nil {
		loggers.LoggerAPKOperator.Errorf("Error converting RoutePolicy to JSON: %v", err)
	} else {
		utils.SendRoutePolicyCreatedOrUpdatedEvent(routePolicyString)
		loggers.LoggerAPKOperator.Debugf("Deleted RoutePolicy JSON: %s", routePolicyString)
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (routePolicyReconciler *RoutePolicyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dpV2alpha1.RoutePolicy{}).
		Complete(routePolicyReconciler)
}

// 
func (routePolicyReconciler *RoutePolicyReconciler) getRoutePolicyForConfigMap(ctx context.Context, obj *corev1.ConfigMap) []reconcile.Request {
	configMap := obj

	requests := []reconcile.Request{}

	routePolicyList := &dpV2alpha1.RoutePolicyList{}
	if err := routePolicyReconciler.client.List(ctx, routePolicyList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(configMapIndex, utils.NamespacedName(configMap).String()),
	}); err != nil {
		return []reconcile.Request{}
	}

	for item := range routePolicyList.Items {
		routePolicy := routePolicyList.Items[item]
		requests = append(requests, routePolicyReconciler.AddRoutePolicyRequest(&routePolicy)...)
	}

	return requests
}

func (routePolicyReconciler *RoutePolicyReconciler) getRoutePolicyForSecret(ctx context.Context, obj *corev1.Secret) []reconcile.Request {
	secret := obj

	requests := []reconcile.Request{}

	routePolicyList := &dpV2alpha1.RoutePolicyList{}
	if err := routePolicyReconciler.client.List(ctx, routePolicyList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(secretIndex, utils.NamespacedName(secret).String()),
	}); err != nil {
		return []reconcile.Request{}
	}

	for item := range routePolicyList.Items {
		routePolicy := routePolicyList.Items[item]
		requests = append(requests, routePolicyReconciler.AddRoutePolicyRequest(&routePolicy)...)
	}

	return requests
}

// AddRoutePolicyRequest adds a request to reconcile for the given route policy
func (routePolicyReconciler *RoutePolicyReconciler) AddRoutePolicyRequest(obj k8client.Object) []reconcile.Request {
	routePolicy, ok := obj.(*dpV2alpha1.RoutePolicy)
	if !ok {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2622, logging.TRIVIAL,
			"Unexpected object type, bypassing reconciliation: %v", routePolicy))
		return nil
	}

	return []reconcile.Request{{
		NamespacedName: types.NamespacedName{
			Name:      string(routePolicy.Name),
			Namespace: routePolicy.Namespace,
		},
	}}
}

func (routePolicyReconciler *RoutePolicyReconciler) addRoutePolicyIndexes(ctx context.Context, mgr manager.Manager) error {
	// Index by referenced ConfigMaps
	if err := mgr.GetFieldIndexer().IndexField(ctx, &dpV2alpha1.RoutePolicy{}, configMapIndex,
		func(rawObj k8client.Object) []string {
			routePolicy := rawObj.(*dpV2alpha1.RoutePolicy)
			namespacedConfigMaps := make([]string, 0)
			if routePolicy.Spec.RequestMediation != nil {
				for _, mediation := range routePolicy.Spec.RequestMediation {
					if mediation.Parameters != nil {
						for _, param := range mediation.Parameters {
							if param.ValueRef != nil {
								if param.ValueRef.Kind == constants.KindConfigMap {
									// If ValueRef is set, use the namespaced name of the referenced object
									namespacedConfigMaps = append(namespacedConfigMaps, types.NamespacedName{
										Namespace: routePolicy.GetNamespace(),
										Name:      string(param.ValueRef.Name),
									}.String())
								}
							}
						}
					}
				}
			}

			return namespacedConfigMaps
		}); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2614, logging.BLOCKER,
			"Error adding index for RoutePolicy: %v", err))
		return err
	}

	// Index by referenced Secrets
	if err := mgr.GetFieldIndexer().IndexField(ctx, &dpV2alpha1.RoutePolicy{}, secretIndex,
		func(rawObj k8client.Object) []string {
			routePolicy := rawObj.(*dpV2alpha1.RoutePolicy)
			namespacedSecrets := make([]string, 0)
			if routePolicy.Spec.RequestMediation != nil {
				for _, mediation := range routePolicy.Spec.RequestMediation {
					if mediation.Parameters != nil {
						for _, param := range mediation.Parameters {
							if param.ValueRef != nil && param.ValueRef.Kind == constants.KindSecret {
								namespacedSecrets = append(namespacedSecrets, types.NamespacedName{
									Namespace: routePolicy.GetNamespace(),
									Name:      string(param.ValueRef.Name),
								}.String())
							}
						}
					}
				}
			}
			return namespacedSecrets
		}); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2614, logging.BLOCKER,
			"Error adding Secret index for RoutePolicy: %v", err))
		return err
	}

	return nil
}

func getValueFromMap(data map[string]string, key, name string) (string, bool) {
	if val, ok := data[key]; ok {
		return val, true
	}
	loggers.LoggerAPKOperator.Warnf("key %s not found in ConfigMap %s", key, name)
	for k, v := range data {
		loggers.LoggerAPKOperator.Warnf("key %s not found in ConfigMap %s, falling back to first key %s", key, name, k)
		return v, true
	}
	return "", false
}

func getValueFromSecret(data map[string][]byte, key, name string) (string, bool) {
	if val, ok := data[key]; ok {
		return string(val), true
	}
	loggers.LoggerAPKOperator.Warnf("key %s not found in Secret %s", key, name)
	for k, v := range data {
		loggers.LoggerAPKOperator.Warnf("key %s not found in Secret %s, falling back to first key %s", key, name, k)
		return string(v), true
	}
	return "", false
}
