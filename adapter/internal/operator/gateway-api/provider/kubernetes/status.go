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
	"fmt"

	"github.com/wso2/apk/adapter/internal/loggers"
	"github.com/wso2/apk/adapter/internal/operator/message"
	"github.com/wso2/apk/adapter/internal/operator/status"
	"github.com/wso2/apk/adapter/internal/operator/utils"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
)

// subscribeAndUpdateStatus subscribes to gateway API object status updates and
// writes it into the Kubernetes API Server.
func (r *gatewayReconcilerNew) subscribeAndUpdateStatus(ctx context.Context) {
	// Gateway object status updater
	go func() {
		message.HandleSubscription(
			message.Metadata{Runner: "provider", Message: "gateway-status"},
			r.resources.GatewayStatuses.Subscribe(ctx),
			func(update message.Update[types.NamespacedName, *gwapiv1.GatewayStatus], errChan chan error) {
				// skip delete updates.
				if update.Delete {
					return
				}
				// Get gateway object
				gtw := new(gwapiv1.Gateway)
				if err := r.client.Get(ctx, update.Key, gtw); err != nil {
					loggers.LoggerAPKOperator.Error(err, "gateway not found", "namespace", gtw.Namespace, "name", gtw.Name)
					errChan <- err
					return
				}
				// Set the updated Status and call the status update
				gtw.Status = *update.Value
				r.updateStatusForGateway(ctx, gtw)
			},
		)
		loggers.LoggerAPKOperator.Info("gateway status subscriber shutting down")
	}()

	// HTTPRoute object status updater
	go func() {
		message.HandleSubscription(
			message.Metadata{Runner: "provider", Message: "httproute-status"},
			r.resources.HTTPRouteStatuses.Subscribe(ctx),
			func(update message.Update[types.NamespacedName, *gwapiv1.HTTPRouteStatus], errChan chan error) {
				// skip delete updates.
				if update.Delete {
					return
				}
				key := update.Key
				val := update.Value
				r.statusUpdater.Send(status.Update{
					NamespacedName: key,
					Resource:       new(gwapiv1.HTTPRoute),
					UpdateStatus: func(obj client.Object) client.Object {
						h, ok := obj.(*gwapiv1.HTTPRoute)
						if !ok {
							err := fmt.Errorf("unsupported object type %T", obj)
							errChan <- err
							panic(err)
						}
						hCopy := h.DeepCopy()
						hCopy.Status.Parents = val.Parents
						return hCopy
					},
				})
			},
		)
		loggers.LoggerAPKOperator.Info("HttpRoute status subscriber shutting down")
	}()

}

func (r *gatewayReconcilerNew) updateStatusForGateway(ctx context.Context, gtw *gwapiv1.Gateway) {
	// nil check for unit tests.
	if r.statusUpdater == nil {
		return
	}

	// Get deployment
	deploy, err := r.envoyDeploymentForGateway(ctx, gtw)
	if err != nil || deploy == nil {
		loggers.LoggerAPKOperator.Infof("Failed to get Deployment for gateway %s/%s, %+v",
			gtw.Namespace, gtw.Name, err)
	}

	// Get service
	svc, err := r.envoyServiceForGateway(ctx, gtw)
	if err != nil || svc == nil {
		loggers.LoggerAPKOperator.Infof("Failed to get Service for gateway %s:%s, %+v",
			gtw.Namespace, gtw.Name, err)
	}
	// update accepted condition
	status.UpdateGatewayStatusAcceptedCondition(gtw, true)
	// update address field and programmed condition
	status.UpdateGatewayStatusProgrammedCondition(gtw, svc, deploy, r.store.listNodeAddresses()...)

	key := utils.NamespacedName(gtw)

	// publish status
	r.statusUpdater.Send(status.Update{
		NamespacedName: key,
		Resource:       new(gwapiv1.Gateway),
		UpdateStatus: func(obj client.Object) client.Object {
			g, ok := obj.(*gwapiv1.Gateway)
			if !ok {
				panic(fmt.Sprintf("unsupported object type %T", obj))
			}
			gCopy := g.DeepCopy()
			gCopy.Status.Conditions = gtw.Status.Conditions
			gCopy.Status.Addresses = gtw.Status.Addresses
			gCopy.Status.Listeners = gtw.Status.Listeners
			return gCopy
		},
	})
}
