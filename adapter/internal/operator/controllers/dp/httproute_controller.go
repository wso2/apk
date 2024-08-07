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
	"fmt"

	"github.com/wso2/apk/adapter/config"
	"github.com/wso2/apk/adapter/internal/discovery/xds/common"
	"github.com/wso2/apk/adapter/internal/loggers"
	"github.com/wso2/apk/adapter/internal/operator/constants"
	"github.com/wso2/apk/adapter/internal/operator/status"
	"github.com/wso2/apk/adapter/internal/operator/synchronizer"
	"github.com/wso2/apk/adapter/internal/operator/utils"
	"github.com/wso2/apk/adapter/pkg/logging"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	k8client "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
)

var httprouteParentKind = "Gateway"

// HTTPRouteReconciler reconciles a HttpRoute object
type HTTPRouteReconciler struct {
	client        k8client.Client
	ods           *synchronizer.OperatorDataStore
	statusUpdater *status.UpdateHandler
	mgr           manager.Manager
}

const (
	conditionTypeProgrammed                                         = "Programmed"
	conditionReasonProgrammedUnknown   gwapiv1.RouteConditionReason = "Unknown"
	conditionReasonConfiguredInGateway gwapiv1.RouteConditionReason = "ConfiguredInGateway"
	conditionReasonTranslationError    gwapiv1.RouteConditionReason = "TranslationError"
)

type referencedGatewaysAndCondition struct {
	gateway      *gwapiv1.Gateway
	condition    metav1.Condition
	listenerName string
}

// NewHTTPRouteController creates a new HttpRoute controller instance. Httproute Controllers watches for gwapiv1.HTTPRoute.
func NewHTTPRouteController(mgr manager.Manager, operatorDataStore *synchronizer.OperatorDataStore, statusUpdater *status.UpdateHandler) error {
	r := &HTTPRouteReconciler{
		client:        mgr.GetClient(),
		ods:           operatorDataStore,
		statusUpdater: statusUpdater,
		mgr:           mgr,
	}
	c, err := controller.New(constants.HTTPRouteController, mgr, controller.Options{Reconciler: r})
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error3119, logging.BLOCKER, "Error creating HttpRoute controller, error: %v", err))
		return err
	}
	conf := config.ReadConfigs()
	predicates := []predicate.Predicate{predicate.NewPredicateFuncs(utils.FilterByNamespaces(conf.Adapter.Operator.Namespaces))}

	if err := c.Watch(source.Kind(mgr.GetCache(), &gwapiv1.HTTPRoute{}), &handler.EnqueueRequestForObject{},
		predicates...); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error3100, logging.BLOCKER, "Error watching HttpRoute resources: %v", err))
		return err
	}

	if err := c.Watch(source.Kind(mgr.GetCache(), &gwapiv1.Gateway{}),
		handler.EnqueueRequestsFromMapFunc(r.handleGateway), predicates...); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error3121, logging.BLOCKER, "Error watching Gateway resources: %v", err))
		return err
	}

	loggers.LoggerAPKOperator.Info("HttpRoute Controller successfully started. Watching Httproute Objects....")
	return nil
}

//+kubebuilder:rbac:groups=dp.wso2.com,resources=httproutes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=dp.wso2.com,resources=httproutes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=dp.wso2.com,resources=httproutes/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the HttpRoute object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (httpRouteReconciler *HTTPRouteReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// Check whether the Httproute CR exist, if not consider as a DELETE event.
	loggers.LoggerAPKOperator.Infof("Reconciling HttpRoute... %s", req.NamespacedName.String())
	var httproute gwapiv1.HTTPRoute
	if err := httpRouteReconciler.client.Get(ctx, req.NamespacedName, &httproute); err != nil {
		return ctrl.Result{}, nil
	}
	supportedGatways := getSupportedGatewaysForRoute(ctx, httpRouteReconciler.client, httproute)
	existingStatuses := httproute.Status.Parents
	// Check whether the route has supported gateways
	if supportedGatways == nil || len(supportedGatways) == 0 {
		loggers.LoggerAPKOperator.Debugf("Could not find any supported gateway for the httproute.")

		// There are no supported gateways found for this http route. We need to cleanup the parent statuses if needed
		newStatuses := make([]gwapiv1.RouteParentStatus, 0)
		for _, status := range existingStatuses {
			if string(status.ControllerName) == GetControllerName() {
				newStatuses = append(newStatuses, status)
			}
		}
		if len(newStatuses) != len(existingStatuses) {
			loggers.LoggerAPKOperator.Debugf("Cleaning up unnecessary statuses from HTTPRoute cr.")
			// Need to change the statuses
			httpRouteReconciler.statusUpdater.Send(status.Update{
				NamespacedName: req.NamespacedName,
				Resource:       new(gwapiv1.HTTPRoute),
				UpdateStatus: func(obj k8client.Object) k8client.Object {
					h, ok := obj.(*gwapiv1.HTTPRoute)
					if !ok {
						loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error3109, logging.BLOCKER, "Error while updating HttpRoute status %v", obj))
					}
					hCopy := h.DeepCopy()
					hCopy.Status.Parents = newStatuses
					return hCopy
				},
			})
		}
		// There is no status change needed
	} else {
		// We found some supported gateways so we need to update parent statuses
		// Create a hashMap to store the already existing statuses by parentref properties as key so that we will not duplicate the status.
		routeParentStatusMap := make(map[string]gwapiv1.RouteParentStatus)
		for _, routeParentStatus := range existingStatuses {
			key := fmt.Sprintf("%s/%s/%s/%s", *routeParentStatus.ParentRef.Group, *routeParentStatus.ParentRef.Kind, routeParentStatus.ParentRef.Name, *routeParentStatus.ParentRef.Namespace)
			routeParentStatusMap[key] = routeParentStatus
		}
		needsStatusUpdate := false
		// Prepare a local copy of status for supported gateways
		for _, gatewayWithCondition := range supportedGatways {
			gatewayParentStatus := &gwapiv1.RouteParentStatus{
				ParentRef: gwapiv1.ParentReference{
					Group:     (*gwapiv1.Group)(&gwapiv1.GroupVersion.Group),
					Kind:      common.PointerCopy(gwapiv1.Kind(httprouteParentKind)),
					Namespace: (*gwapiv1.Namespace)(&gatewayWithCondition.gateway.Namespace),
					Name:      gwapiv1.ObjectName(gatewayWithCondition.gateway.Name),
				},
				ControllerName: gwapiv1.GatewayController(GetControllerName()),
				Conditions: []metav1.Condition{{
					Type:               gatewayWithCondition.condition.Type,
					Status:             gatewayWithCondition.condition.Status,
					ObservedGeneration: httproute.Generation,
					LastTransitionTime: metav1.Now(),
					Reason:             gatewayWithCondition.condition.Reason,
				}},
			}
			key := fmt.Sprintf("%s/%s/%s/%s", *gatewayParentStatus.ParentRef.Group, *gatewayParentStatus.ParentRef.Kind, gatewayParentStatus.ParentRef.Name, *gatewayParentStatus.ParentRef.Namespace)
			foundStatus, exists := routeParentStatusMap[key]
			statusChanged := true
			if exists {
				if foundStatus.ControllerName == gatewayParentStatus.ControllerName &&
					*foundStatus.ParentRef.Kind == *gatewayParentStatus.ParentRef.Kind &&
					foundStatus.ParentRef.Name == gatewayParentStatus.ParentRef.Name &&
					*foundStatus.ParentRef.Namespace == *gatewayParentStatus.ParentRef.Namespace &&
					*foundStatus.ParentRef.Group == *gatewayParentStatus.ParentRef.Group &&
					len(foundStatus.Conditions) > 0 &&
					common.AreConditionsSame(foundStatus.Conditions[0], gatewayParentStatus.Conditions[0]) {
					statusChanged = false
				}
			}
			needsStatusUpdate = needsStatusUpdate || statusChanged
			if statusChanged {
				status, found := common.FindElement(routeParentStatusMap[key].Conditions, func(cond metav1.Condition) bool {
					if cond.Type == conditionTypeProgrammed {
						return true
					}
					return false
				})
				if !found {
					gatewayParentStatus.Conditions = append(gatewayParentStatus.Conditions, metav1.Condition{
						Type:               conditionTypeProgrammed,
						Status:             metav1.ConditionUnknown,
						Reason:             string(conditionReasonProgrammedUnknown),
						ObservedGeneration: httproute.Generation,
						LastTransitionTime: metav1.Now(),
					})
				} else {
					gatewayParentStatus.Conditions = append(gatewayParentStatus.Conditions, status)
				}
				routeParentStatusMap[key] = *gatewayParentStatus
			}
		}

		if needsStatusUpdate || len(httproute.Status.Parents) != len(routeParentStatusMap) {
			httproute.Status.Parents = make([]gwapiv1.RouteParentStatus, 0)
			for _, parentStatus := range routeParentStatusMap {
				httproute.Status.Parents = append(httproute.Status.Parents, parentStatus)
			}

			httpRouteReconciler.statusUpdater.Send(status.Update{
				NamespacedName: req.NamespacedName,
				Resource:       new(gwapiv1.HTTPRoute),
				UpdateStatus: func(obj k8client.Object) k8client.Object {
					h, ok := obj.(*gwapiv1.HTTPRoute)
					if !ok {
						loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error3109, logging.BLOCKER, "Error while updating HttpRoute status %v", obj))
					}
					hCopy := h.DeepCopy()
					hCopy.Status.Parents = httproute.Status.Parents
					return hCopy
				},
			})
		}

	}

	return ctrl.Result{}, nil
}

// handleGateway handles the gateway changes and create reconcile reuqest for HTTPRoute reconciler
func (httpRouteReconciler *HTTPRouteReconciler) handleGateway(ctx context.Context, obj k8client.Object) []reconcile.Request {
	gateway, ok := obj.(*gwapiv1.Gateway)
	if !ok {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2622, logging.TRIVIAL, "Unexpected object type, bypassing reconciliation: %v", gateway))
		return []reconcile.Request{}
	}

	httpRouteList := &gwapiv1.HTTPRouteList{}
	if err := httpRouteReconciler.client.List(ctx, httpRouteList, &k8client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(gatewayHTTPRouteIndex, utils.NamespacedName(gateway).String()),
	}); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2625, logging.CRITICAL, "Unable to find associated HTTPRoutes: %s", utils.NamespacedName(gateway).String()))
		return []reconcile.Request{}
	}

	if len(httpRouteList.Items) == 0 {
		loggers.LoggerAPKOperator.Debugf("HTTPRoutes for Gateway not found: %s", utils.NamespacedName(gateway).String())
		return []reconcile.Request{}
	}

	requests := []reconcile.Request{}
	for item := range httpRouteList.Items {
		httpRoute := httpRouteList.Items[item]
		if refersGateway(httpRoute, *gateway) {
			requests = append(requests, reconcile.Request{
				NamespacedName: utils.NamespacedName(&httpRoute),
			})
		} else {
			loggers.LoggerAPKOperator.Debugf("Gateway change observed but HttpRoute: %s does not belongs to this gateway: %s hence not reconciling.",
				utils.NamespacedName(&httpRoute).String(),
				utils.NamespacedName(gateway).String())
		}
	}
	return requests
}

func getSupportedGatewaysForRoute(ctx context.Context, client k8client.Client, httpRoute gwapiv1.HTTPRoute) []referencedGatewaysAndCondition {
	parentRefs := httpRoute.Spec.ParentRefs
	if parentRefs == nil {
		loggers.LoggerAPKOperator.Errorf("Parent ref not found for HTTPRoute: %s/%s", httpRoute.Namespace, httpRoute.Name)
		return nil
	}
	referencedGatewayList := make([]referencedGatewaysAndCondition, 0)
	for _, parentRef := range parentRefs {
		namespace := httpRoute.GetNamespace()
		if parentRef.Namespace != nil {
			namespace = string(*parentRef.Namespace)
		}
		name := string(parentRef.Name)

		// Try to fetch the referenced gateway
		gateway := gwapiv1.Gateway{}
		if err := client.Get(ctx, types.NamespacedName{
			Namespace: namespace,
			Name:      name,
		}, &gateway); err != nil {
			if apierrors.IsNotFound(err) {
				loggers.LoggerAPKOperator.Infof("Could not find gateway: %s/%s", namespace, name)
				// There might be other gateway in this list that can be found. So keep search.
				continue
			}
			loggers.LoggerAPKOperator.Errorf("Error while fetching gateway: %s/%s. Error: %s", namespace, name, err)
			return nil
		}
		var (
			httpRouteMatched        = false
			hostnameMatched         = false
			portMatched             = false
			satisfiesAllowedRoutes  = false
			satisfiesSupportedKinds = false
			satisfiesListenerName   = false
		)

		for _, listener := range gateway.Spec.Listeners {
			if ok, err := checkIfRouteMatchesAllowedRoutes(ctx, client, httpRoute, listener, gateway.Namespace, parentRef.Namespace); err != nil {
				return nil
			} else if !ok {
				loggers.LoggerAPKOperator.Debugf("HttpRoute(%s/%s) did not match any allowed routes in gateway: %s/%s, listener: %s", httpRoute.Namespace, httpRoute.Name, namespace, name, listener.Name)
				continue
			}
			satisfiesAllowedRoutes = true
			if err := checkIfMatchingReadyListenerExistsInStatus(httpRoute, listener, gateway.Status.Listeners); err != nil {
				loggers.LoggerAPKOperator.Debugf("Gateway(%s/%s) listener: %s does not have a ready listener for this HttpRoute(%s/%s). Error: %s", namespace, name, listener.Name, httpRoute.Namespace, httpRoute.Name, err.Error())
				continue
			}
			satisfiesSupportedKinds = true
			if parentRef.SectionName != nil {
				if *parentRef.SectionName != "" && *parentRef.SectionName != listener.Name {
					loggers.LoggerAPKOperator.Debugf("Gateway(%s/%s) listener: %s does not have a matching listener with section name %s for this HttpRoute(%s/%s)", namespace, name, listener.Name, *parentRef.SectionName, httpRoute.Namespace, httpRoute.Name)
					continue
				}
				satisfiesListenerName = true
			}
			if parentRef.Port != nil {
				if *parentRef.Port != listener.Port {
					loggers.LoggerAPKOperator.Debugf("Gateway(%s/%s) listener: %s does not have a matching port (%d) for this HttpRoute(%s/%s)", namespace, name, listener.Name, int32(*parentRef.Port), httpRoute.Namespace, httpRoute.Name)
					continue
				}
				portMatched = true
			}
			hostnameMatched = isGatewayListenerHostnameMatched(listener, httpRoute.Spec.Hostnames)
			httpRouteMatched = hostnameMatched
		}

		if httpRouteMatched {
			var listenerName string
			if parentRef.SectionName != nil && *parentRef.SectionName != "" {
				listenerName = string(*parentRef.SectionName)
			}

			referencedGatewayList = append(referencedGatewayList, referencedGatewaysAndCondition{
				gateway:      &gateway,
				listenerName: listenerName,
				condition: metav1.Condition{
					Type:               string(gwapiv1.RouteConditionAccepted),
					Status:             metav1.ConditionTrue,
					Reason:             string(gwapiv1.RouteReasonAccepted),
					ObservedGeneration: httpRoute.GetGeneration(),
				},
			})
		} else {
			reason := gwapiv1.RouteReasonNoMatchingParent
			if !hostnameMatched {
				reason = gwapiv1.RouteReasonNoMatchingListenerHostname
			} else if (parentRef.SectionName) != nil && !satisfiesListenerName {
				reason = gwapiv1.RouteReasonNoMatchingParent
			} else if (parentRef.Port != nil) && !portMatched {
				reason = gwapiv1.RouteReasonNoMatchingParent
			} else if !satisfiesAllowedRoutes || !satisfiesSupportedKinds {
				reason = gwapiv1.RouteReasonNotAllowedByListeners
			}
			var listenerName string
			if parentRef.SectionName != nil && *parentRef.SectionName != "" {
				listenerName = string(*parentRef.SectionName)
			}

			referencedGatewayList = append(referencedGatewayList, referencedGatewaysAndCondition{
				gateway:      &gateway,
				listenerName: listenerName,
				condition: metav1.Condition{
					Type:               string(gwapiv1.RouteConditionAccepted),
					Status:             metav1.ConditionFalse,
					Reason:             string(reason),
					ObservedGeneration: httpRoute.GetGeneration(),
				},
			})
		}
	}
	return referencedGatewayList
}

// checkIfRouteMatchesAllowedRoutes determines if the provided route matches the
// criteria defined in the listener's AllowedRoutes field.
func checkIfRouteMatchesAllowedRoutes(
	ctx context.Context,
	client k8client.Client,
	httpRoute gwapiv1.HTTPRoute,
	listener gwapiv1.Listener,
	gatewayNamespace string,
	parentRefNamespace *gwapiv1.Namespace,
) (bool, error) {
	if listener.AllowedRoutes == nil {
		return true, nil
	}

	if len(listener.AllowedRoutes.Kinds) > 0 {
		_, found := common.FindElement(listener.AllowedRoutes.Kinds, func(allowedGroupKind gwapiv1.RouteGroupKind) bool {
			httpRouteVersionKind := httpRoute.GetObjectKind().GroupVersionKind()
			return (allowedGroupKind.Group != nil && string(*allowedGroupKind.Group) == httpRouteVersionKind.Group) && string(allowedGroupKind.Kind) == httpRouteVersionKind.Kind
		})
		if !found {
			return false, nil
		}
	}

	if listener.AllowedRoutes.Namespaces == nil || listener.AllowedRoutes.Namespaces.From == nil {
		return true, nil
	}

	switch *listener.AllowedRoutes.Namespaces.From {
	case gwapiv1.NamespacesFromAll:
		return true, nil

	case gwapiv1.NamespacesFromSame:
		if parentRefNamespace == nil {
			return gatewayNamespace == httpRoute.GetNamespace(), nil
		}
		return httpRoute.GetNamespace() == string(*parentRefNamespace), nil

	case gwapiv1.NamespacesFromSelector:
		namespace := corev1.Namespace{}
		if err := client.Get(ctx, types.NamespacedName{Name: httpRoute.GetNamespace()}, &namespace); err != nil {
			return false, fmt.Errorf("failed to get namespace %s: %w", httpRoute.GetNamespace(), err)
		}

		s, err := metav1.LabelSelectorAsSelector(listener.AllowedRoutes.Namespaces.Selector)
		if err != nil {
			return false, fmt.Errorf(
				"failed to convert AllowedRoutes LabelSelector %s to Selector for listener %s: %w",
				listener.AllowedRoutes.Namespaces.Selector, listener.Name, err,
			)
		}

		ok := s.Matches(labels.Set(namespace.Labels))
		return ok, nil

	default:
		return false, fmt.Errorf(
			"unknown listener.AllowedRoutes.Namespaces.From value: %s for listener %s",
			*listener.AllowedRoutes.Namespaces.From, listener.Name,
		)
	}
}

// checkIfMatchingReadyListenerExistsInStatus determines if there exists a matching ready listener
// in the provided list of listener statuses.
func checkIfMatchingReadyListenerExistsInStatus(route gwapiv1.HTTPRoute, listener gwapiv1.Listener, listenerStatuses []gwapiv1.ListenerStatus) error {

	// Find the relative gateway listener status for the gateway listener
	listenerStatus, listenerFound := common.FindElement(listenerStatuses, func(listenerStatus gwapiv1.ListenerStatus) bool {
		return listenerStatus.Name != listener.Name
	})
	if !listenerFound {
		return fmt.Errorf("Cannot find a listener status for this route in gateway")
	}

	// Check if the programmed status exists
	programmedStatus, foundProgrammedStatus := common.FindElement(listenerStatus.Conditions, func(c metav1.Condition) bool {
		return c.Type == string(gwapiv1.ListenerConditionProgrammed)
	})
	if !foundProgrammedStatus {
		return fmt.Errorf("Cannot find a programmed listener status for this route in gateway")
	}
	if programmedStatus.Status != "True" {
		return fmt.Errorf("Programmed status is not active yet")
	}

	return nil
}

func isGatewayListenerHostnameMatched(listener gwapiv1.Listener, hostnames []gwapiv1.Hostname) bool {
	// If httpRoute does not specify any hostnames then we can assume it accept gateway listener's hostname
	if len(hostnames) == 0 {
		return true
	}

	// If listener hostname is nil or empty it will accept all hostnames from the httpRoute
	if listener.Hostname == nil || *listener.Hostname == "" {
		return true
	}

	for _, hostname := range hostnames {
		if common.MatchesHostname(string(hostname), string(*listener.Hostname)) {
			return true
		}
	}

	return false
}

func refersGateway(httpRoute gwapiv1.HTTPRoute, gateway gwapiv1.Gateway) bool {
	_, found := common.FindElement(httpRoute.Spec.ParentRefs, func(parentRef gwapiv1.ParentReference) bool {
		namespace := ""
		if parentRef.Namespace != nil {
			namespace = string(*parentRef.Namespace)
		}
		if namespace == "" {
			namespace = httpRoute.GetNamespace()
		}
		referingGatewayNamespacedName := types.NamespacedName{
			Namespace: namespace,
			Name:      string(parentRef.Name),
		}
		gatewayNamespace := types.NamespacedName{
			Namespace: gateway.Namespace,
			Name:      gateway.Name,
		}
		if referingGatewayNamespacedName.String() == gatewayNamespace.String() {
			return true
		}

		return false
	})
	return found
}
