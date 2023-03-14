/*
 *  Copyright (c) 2023, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
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

package synchronizer

import (
	"github.com/wso2/apk/adapter/internal/loggers"
	"github.com/wso2/apk/adapter/pkg/logging"
)

// GatewayEvent holds the data structure used for passing Gateway
// events from controller go routine to synchronizer
// go routine.
type GatewayEvent struct {
	EventType string
	Event     GatewayState
}

// HandleGatewayLifeCycleEvents handles the Gateway events generated from OperatorDataStore
func HandleGatewayLifeCycleEvents(ch *chan GatewayEvent) {
	loggers.LoggerAPKOperator.Info("Operator synchronizer listening for Gateway lifecycle events...")
	for event := range *ch {
		if event.Event.GatewayDefinition == nil {
			loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2628))
		}
		loggers.LoggerAPKOperator.Infof("%s event received for %v", event.EventType, event.Event.GatewayDefinition.Name)
		var err error
		// switch event.EventType {
		// // case constants.Delete:
		// // 	err = undeployGateway(event.Event)
		// // case constants.Create:
		// // 	err = deployGateway(event.Event)
		// // case constants.Update:
		// // 	err = deployGateway(event.Event)
		// // }
		if err != nil {
			loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2629, event.EventType, err))
		}
	}
}
