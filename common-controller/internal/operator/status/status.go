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

package status

import (
	"context"

	"github.com/wso2/apk/common-controller/internal/loggers"
	dpv1alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Update contain status update event information
type Update struct {
	NamespacedName types.NamespacedName
	Resource       client.Object
	UpdateStatus   func(client.Object) client.Object
}

// UpdateHandler handles status updates
type UpdateHandler struct {
	client        client.Client
	updateChannel chan Update
}

// NewUpdateHandler get a new status update handler
func NewUpdateHandler(client client.Client) *UpdateHandler {
	return &UpdateHandler{
		client:        client,
		updateChannel: make(chan Update, 100),
	}
}

// applyUpdate perform the status update on CR
func (updateHandler *UpdateHandler) applyUpdate(update Update) {
	err := retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		if err := updateHandler.client.Get(context.Background(), update.NamespacedName, update.Resource); err != nil {
			return err
		}

		resourceCopy := update.UpdateStatus(update.Resource)
		if isStatusEqual(update.Resource, resourceCopy) {
			loggers.LoggerAPKOperator.Debugf("Status unchanged, hence not updating. %s", update.NamespacedName.String())
			return nil
		}
		loggers.LoggerAPKOperator.Debugf("Status is updating for %s ...", update.NamespacedName.String())
		return updateHandler.client.Status().Update(context.Background(), resourceCopy)
	})

	if err != nil {
		loggers.LoggerAPKOperator.Errorf("Unable to update status for %s, Error : %v ", update.NamespacedName.String(), err)
	}
}

// Start starts the status update handler go routine.
func (updateHandler *UpdateHandler) Start(ctx context.Context) error {
	loggers.LoggerAPKOperator.Info("Started status update handler")
	defer loggers.LoggerAPKOperator.Info("Stopped status update handler")

	for {
		select {
		case update := <-updateHandler.updateChannel:
			loggers.LoggerAPKOperator.Debugf("Received a status update in %s", update.NamespacedName.String())
			updateHandler.applyUpdate(update)
		case <-ctx.Done():
			return nil
		}
	}
}

// Send public method to add status update events to the update channel.
func (updateHandler *UpdateHandler) Send(update Update) {
	loggers.LoggerAPKOperator.Debugf("SEND Received a status update in %s", update.NamespacedName.String())
	updateHandler.updateChannel <- update
}

// isStatusEqual checks if two objects have equivalent status.
// Supported:
//   - API
func isStatusEqual(objA, objB interface{}) bool {
	switch a := objA.(type) {
	case *dpv1alpha1.API:
		if b, ok := objB.(*dpv1alpha1.API); ok {
			return compareAPIs(a, b)
		}
	}
	return false
}

// compareAPIs compares status in API CRs.
func compareAPIs(api1 *dpv1alpha1.API, api2 *dpv1alpha1.API) bool {
	return api1.Status.DeploymentStatus.Message == api2.Status.DeploymentStatus.Message
}
