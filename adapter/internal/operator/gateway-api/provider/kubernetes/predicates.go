/*
 *  Copyright (c) 2024, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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

package kubernetes

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"

	"github.com/wso2/apk/adapter/internal/loggers"
	gatewayapi "github.com/wso2/apk/adapter/internal/operator/gateway-api"
	"github.com/wso2/apk/adapter/internal/operator/utils"
	corev1 "k8s.io/api/core/v1"
	k8errors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
)

// hasMatchingController returns true if the provided object is a GatewayClass
// with a Spec.Controller string matching this Envoy Gateway's controller string,
// or false otherwise.
func (r *gatewayReconcilerNew) hasMatchingController(obj client.Object) bool {
	gc, ok := obj.(*gwapiv1.GatewayClass)
	if !ok {
		loggers.LoggerAPKOperator.Info("bypassing reconciliation due to unexpected object type ", obj)
		return false
	}

	if gc.Spec.ControllerName == gatewayClassControllerName {
		loggers.LoggerAPKOperator.Infof("gatewayclass has matching controller name %s, hence processing",
			gc.Name)
		return true
	}

	loggers.LoggerAPKOperator.Info("bypassing reconciliation due to controller name",
		gc.Spec.ControllerName)
	return false
}

// validateGatewayForReconcile returns true if the provided object is a Gateway
// using a GatewayClass matching the configured gatewayclass controller name.
func (r *gatewayReconcilerNew) validateGatewayForReconcile(obj client.Object) bool {
	gw, ok := obj.(*gwapiv1.Gateway)
	if !ok {
		loggers.LoggerAPKOperator.Info("unexpected object type, bypassing reconciliation object", obj)
		return false
	}

	gatewayClass := &gwapiv1.GatewayClass{}
	key := types.NamespacedName{Name: string(gw.Spec.GatewayClassName)}
	if err := r.client.Get(context.Background(), key, gatewayClass); err != nil {
		loggers.LoggerAPKOperator.Errorf("failed to get gatewayclass name %s, %+v", gw.Spec.GatewayClassName, err)
		return false
	}

	if gatewayClass.Spec.ControllerName != gatewayClassControllerName {
		loggers.LoggerAPKOperator.Infof("gatewayclass name %s for gateway doesn't match configured name %s in %s/%s ",
			string(gatewayClass.Spec.ControllerName), gatewayClassControllerName, gw.Namespace, gw.Name)
		return false
	}
	loggers.LoggerAPKOperator.Info("Gateway CR change is detected for ", gw.Name)
	return true
}

// validateServiceForReconcile tries finding the owning Gateway of the Service
// if it exists, finds the Gateway's Deployment, and further updates the Gateway
// status Ready condition. All Services are pushed for reconciliation.
func (r *gatewayReconcilerNew) validateServiceForReconcile(obj client.Object) bool {
	ctx := context.Background()
	svc, ok := obj.(*corev1.Service)
	if !ok {
		loggers.LoggerAPKOperator.Info("unexpected object type, bypassing reconciliation for object", obj)
		return false
	}
	labels := svc.GetLabels()

	// Check if the Service belongs to a Gateway, if so, update the Gateway status.
	gtw := r.findOwningGateway(ctx, labels)
	if gtw != nil {
		r.updateStatusForGateway(ctx, gtw)
		return false
	}

	nsName := utils.NamespacedName(svc)
	return r.isRouteReferencingBackend(&nsName)

}

// findOwningGateway attempts finds a Gateway using "labels".
func (r *gatewayReconcilerNew) findOwningGateway(ctx context.Context, labels map[string]string) *gwapiv1.Gateway {
	gwName, ok := labels[gatewayapi.OwningGatewayNameLabel]
	if !ok {
		return nil
	}

	gwNamespace, ok := labels[gatewayapi.OwningGatewayNamespaceLabel]
	if !ok {
		return nil
	}

	gatewayKey := types.NamespacedName{Namespace: gwNamespace, Name: gwName}
	gtw := new(gwapiv1.Gateway)
	if err := r.client.Get(ctx, gatewayKey, gtw); err != nil {
		loggers.LoggerAPKOperator.Infof("Gateway %+v not found : %v", gatewayKey, err)
		return nil
	}
	return gtw
}

// isRouteReferencingBackend returns true if the backend(service and serviceImport) is referenced by any of the xRoutes
// in the system, else returns false.
func (r *gatewayReconcilerNew) isRouteReferencingBackend(nsName *types.NamespacedName) bool {
	ctx := context.Background()
	httpRouteList := &gwapiv1.HTTPRouteList{}
	if err := r.client.List(ctx, httpRouteList, &client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(serviceHTTPRouteIndex, nsName.String()),
	}); err != nil {
		loggers.LoggerAPKOperator.Error("unable to find associated HTTPRoutes for the service ", err)
		return false
	}

	// Check how many Route objects refer this Backend
	allAssociatedRoutes := len(httpRouteList.Items)

	return allAssociatedRoutes != 0
}

// envoyServiceForGateway returns the Envoy service, returning nil if the service doesn't exist.
func (r *gatewayReconcilerNew) envoyServiceForGateway(ctx context.Context, gateway *gwapiv1.Gateway) (*corev1.Service, error) {
	key := types.NamespacedName{
		Namespace: r.namespace,
		Name:      infraName(gateway),
	}
	svc := new(corev1.Service)
	if err := r.client.Get(ctx, key, svc); err != nil {
		if k8errors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return svc, nil
}

// envoyDeploymentForGateway returns the Envoy Deployment, returning nil if the Deployment doesn't exist.
func (r *gatewayReconcilerNew) envoyDeploymentForGateway(ctx context.Context, gateway *gwapiv1.Gateway) (*appsv1.Deployment, error) {
	key := types.NamespacedName{
		Namespace: r.namespace,
		Name:      infraName(gateway),
	}
	deployment := new(appsv1.Deployment)
	if err := r.client.Get(ctx, key, deployment); err != nil {
		if k8errors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return deployment, nil
}

func (r *gatewayReconcilerNew) handleNode(obj client.Object) bool {
	ctx := context.Background()
	node, ok := obj.(*corev1.Node)
	if !ok {
		loggers.LoggerAPKOperator.Info("unexpected object type, bypassing reconciliation", obj)
		return false
	}

	key := types.NamespacedName{Name: node.Name}
	if err := r.client.Get(ctx, key, node); err != nil {
		if k8errors.IsNotFound(err) {
			r.store.removeNode(node)
			return true
		}
		loggers.LoggerAPKOperator.Error(err, "unable to find node ", node.Name)
		return false
	}

	r.store.addNode(node)
	return true
}

// validateSecretForReconcile checks whether the Secret belongs to a valid Gateway.
func (r *gatewayReconcilerNew) validateSecretForReconcile(obj client.Object) bool {
	secret, ok := obj.(*corev1.Secret)
	if !ok {
		loggers.LoggerAPKOperator.Info("unexpected object type, bypassing reconciliation", obj)
		return false
	}

	nsName := utils.NamespacedName(secret)

	if r.isGatewayReferencingSecret(&nsName) {
		return true
	}

	return false
}

func (r *gatewayReconcilerNew) isGatewayReferencingSecret(nsName *types.NamespacedName) bool {
	gwList := &gwapiv1.GatewayList{}
	if err := r.client.List(context.Background(), gwList, &client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(secretGatewayIndex, nsName.String()),
	}); err != nil {
		loggers.LoggerAPKOperator.Error(err, "unable to find associated Gateways")
		return false
	}

	if len(gwList.Items) == 0 {
		return false
	}

	for _, gw := range gwList.Items {
		gw := gw
		if !r.validateGatewayForReconcile(&gw) {
			return false
		}
	}
	return true
}

// validateDeploymentForReconcile tries finding the owning Gateway of the Deployment
// if it exists, finds the Gateway's Service, and further updates the Gateway
// status Ready condition. No Deployments are pushed for reconciliation.
func (r *gatewayReconcilerNew) validateDeploymentForReconcile(obj client.Object) bool {
	ctx := context.Background()
	deployment, ok := obj.(*appsv1.Deployment)
	if !ok {
		loggers.LoggerAPKOperator.Info("unexpected object type, bypassing reconciliation", obj)
		return false
	}
	labels := deployment.GetLabels()

	// Only deployments in the configured namespace should be reconciled.
	if deployment.Namespace == r.namespace {
		// Check if the deployment belongs to a Gateway, if so, update the Gateway status.
		gtw := r.findOwningGateway(ctx, labels)
		if gtw != nil {
			r.updateStatusForGateway(ctx, gtw)
			return false
		}
	}

	// There is no need to reconcile the Deployment any further.
	return false
}
