/*
 *  Copyright (c) 2022, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
 *
 *  WSO2 LLC. licenses this file to you under the Apache License,
 *  Version 2.0 (the "License"); you may not use this file except
 *  in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package synchronizer

import (
	"fmt"

	"github.com/wso2/apk/adapter/pkg/logging"
	"github.com/wso2/apk/management-server/internal/database"
	"github.com/wso2/apk/management-server/internal/logger"
	"github.com/wso2/apk/management-server/internal/types"
	"github.com/wso2/apk/management-server/internal/xds"
)

// ApplicationCreateAndDeleteEventChannel represents the channel to send/receive Application events
var ApplicationCreateAndDeleteEventChannel chan types.ApplicationEvent

// SubscriptionCreateAndDeleteEventChannel represents the channel to send/receive Subscription events
var SubscriptionCreateAndDeleteEventChannel chan types.SubscriptionEvent

func init() {
	// Channel to handle application create and delete events
	ApplicationCreateAndDeleteEventChannel = make(chan types.ApplicationEvent)
	// Channel to handle subscription create and delete events
	SubscriptionCreateAndDeleteEventChannel = make(chan types.SubscriptionEvent)
}

// AddApplicationEventsToChannel adds the application event to the channel
func AddApplicationEventsToChannel(event types.ApplicationEvent) {
	ApplicationCreateAndDeleteEventChannel <- event
}

// AddSubscriptionEventsToChannel adds the subscription event to the channel
func AddSubscriptionEventsToChannel(event types.SubscriptionEvent) {
	SubscriptionCreateAndDeleteEventChannel <- event
}

// ProcessApplicationEvents processes the application event
func ProcessApplicationEvents() {
	for e := range ApplicationCreateAndDeleteEventChannel {
		if e.IsRemoveEvent {
			database.DbCache.Delete(e.UUID)
			xds.RemoveApplication(e.Label, e.UUID)
		} else {
			app, err := database.GetApplicationByUUID(e.UUID)
			if err != nil {
				logger.LoggerDatabase.ErrorC(logging.ErrorDetails{
					Message: fmt.Sprintf("Error retrieving application for uuid : %s from database error: %v, "+
						"hence skipping add to xdx cache", e.UUID, err),
					Severity:  logging.MINOR,
					ErrorCode: 1101,
				})
				return
			}
			xds.AddSingleApplication(e.Label, app)
		}
	}
}

// ProcessSubscriptionEvents processes the subscription event
func ProcessSubscriptionEvents() {
	for e := range SubscriptionCreateAndDeleteEventChannel {
		if e.IsRemoveEvent {
			err := database.DbCache.DeleteSubscriptionFromApplication(e.AppUUID, e.UUID)
			if err == database.ErrApplicationNotInCache {
				database.GetApplicationByUUID(e.AppUUID)
			}
			app, _ := database.DbCache.Read(e.AppUUID)
			xds.AddSingleApplication(e.Label, app)
		} else if e.IsUpdateEvent {
			sub, _ := database.GetSubscriptionByUUID(e.UUID)
			err := database.DbCache.UpdateSubscriptionInApplication(e.AppUUID, sub)
			if err == database.ErrApplicationNotInCache {
				database.GetApplicationByUUID(e.AppUUID)
			}
			app, _ := database.DbCache.Read(e.AppUUID)
			xds.AddSingleApplication(e.Label, app)
		} else {
			sub, _ := database.GetSubscriptionByUUID(e.UUID)
			err := database.DbCache.AddSubscriptionForApplication(e.AppUUID, sub)
			if err == database.ErrApplicationNotInCache {
				database.GetApplicationByUUID(e.AppUUID)
			}
			app, _ := database.DbCache.Read(e.AppUUID)
			xds.AddSingleApplication(e.Label, app)
		}
	}
}
