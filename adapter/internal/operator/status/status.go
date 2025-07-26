/*
 *  Copyright (c) 2022, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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
	"sync"
	"time"

	"github.com/wso2/apk/adapter/internal/loggers"
	dpv1alpha3 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha3"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Update contains information for a status update event.
type Update struct {
	NamespacedName types.NamespacedName
	Resource       client.Object
	UpdateStatus   func(client.Object) client.Object
}

// UpdateHandler handles status updates.
type UpdateHandler struct {
	client        client.Client
	updateChannel chan Update
}

// DedupingUpdateHandler wraps UpdateHandler and deduplicates updates for the same resource.
type DedupingUpdateHandler struct {
	handler   *UpdateHandler
	pending   map[types.NamespacedName]Update
	mutex     sync.Mutex
	flushTick time.Duration
	stopChan  chan struct{}
}

// NewUpdateHandler creates a new status update handler.
func NewUpdateHandler(client client.Client) *UpdateHandler {
	return &UpdateHandler{
		client:        client,
		updateChannel: make(chan Update, 50), // Smaller buffer to prevent OOM
	}
}

// NewDedupingUpdateHandler creates a deduping wrapper for UpdateHandler.
func NewDedupingUpdateHandler(handler *UpdateHandler, flushTick time.Duration) *DedupingUpdateHandler {
	if flushTick == 0 {
		flushTick = 2 * time.Second
	}
	return &DedupingUpdateHandler{
		handler:   handler,
		pending:   make(map[types.NamespacedName]Update),
		flushTick: flushTick,
		stopChan:  make(chan struct{}),
	}
}

// applyUpdate performs the status patch on a CR.
func (updateHandler *UpdateHandler) applyUpdate(update Update) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := retry.OnError(retry.DefaultBackoff, func(err error) bool {
		return errors.IsConflict(err) || errors.IsServerTimeout(err)
	}, func() error {
		latest := update.Resource.DeepCopyObject().(client.Object)
		if err := updateHandler.client.Get(ctx, update.NamespacedName, latest); err != nil {
			if errors.IsNotFound(err) {
				loggers.LoggerAPKOperator.Warnf("API CR %s not found, skipping status update",
					update.NamespacedName.String())
				return nil
			}
			return err
		}

		updatedObj := update.UpdateStatus(latest)
		if isStatusEqual(latest, updatedObj) {
			loggers.LoggerAPKOperator.Debugf("Status unchanged for %s, skipping update",
				update.NamespacedName.String())
			return nil
		}

		loggers.LoggerAPKOperator.Debugf("Patching status for %s ...", update.NamespacedName.String())
		return updateHandler.client.Status().Patch(ctx, updatedObj, client.MergeFrom(latest))
	})

	if err != nil {
		loggers.LoggerAPKOperator.Errorf("Unable to patch status for %s, Kind: %+v, Error: %v",
			update.NamespacedName.String(), update.Resource.GetObjectKind(), err)
	}
}

// Start starts the status update handler goroutine.
func (updateHandler *UpdateHandler) Start(ctx context.Context) error {
	loggers.LoggerAPKOperator.Info("Started status update handler")
	defer loggers.LoggerAPKOperator.Info("Stopped status update handler")

	for {
		select {
		case update := <-updateHandler.updateChannel:
			loggers.LoggerAPKOperator.Debugf("Received a status update event for %s",
				update.NamespacedName.String())
			updateHandler.applyUpdate(update)
		case <-ctx.Done():
			return nil
		}
	}
}

// Send adds a status update event to the update channel.
func (updateHandler *UpdateHandler) Send(update Update) {
	select {
	case updateHandler.updateChannel <- update:
	default:
		loggers.LoggerAPKOperator.Warnf("Dropping status update for %s - queue is full",
			update.NamespacedName.String())
	}
}

// Send deduplicates updates for the same resource and flushes later.
func (d *DedupingUpdateHandler) Send(update Update) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.pending[update.NamespacedName] = update
}

// StartDeduper runs the deduper and forwards unique updates to UpdateHandler.
func (d *DedupingUpdateHandler) StartDeduper(ctx context.Context) {
	ticker := time.NewTicker(d.flushTick)
	defer ticker.Stop()

	loggers.LoggerAPKOperator.Infof("Started deduping status update handler with flush interval %v", d.flushTick)

	for {
		select {
		case <-ticker.C:
			d.flushPending()
		case <-ctx.Done():
			loggers.LoggerAPKOperator.Info("Stopped deduping status update handler")
			return
		}
	}
}

// flushPending pushes all pending unique updates to the underlying UpdateHandler.
func (d *DedupingUpdateHandler) flushPending() {
	d.mutex.Lock()
	updates := make([]Update, 0, len(d.pending))
	for _, u := range d.pending {
		updates = append(updates, u)
	}
	d.pending = make(map[types.NamespacedName]Update)
	d.mutex.Unlock()

	for _, u := range updates {
		loggers.LoggerAPKOperator.Debugf("Flushing deduped status update for %s", u.NamespacedName)
		d.handler.Send(u)
	}
}

// isStatusEqual checks if two objects have equivalent status.
func isStatusEqual(objA, objB interface{}) bool {
	switch a := objA.(type) {
	case *dpv1alpha3.API:
		if b, ok := objB.(*dpv1alpha3.API); ok {
			return compareAPIs(a, b)
		}
	}
	return false
}

// compareAPIs compares the status of API CRs.
func compareAPIs(api1 *dpv1alpha3.API, api2 *dpv1alpha3.API) bool {
	return api1.Status.DeploymentStatus.Message == api2.Status.DeploymentStatus.Message
}
