/*
 *  Copyright (c) 2025, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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
 
package analytics

import (
	"fmt"

	v3 "github.com/envoyproxy/go-control-plane/envoy/data/accesslog/v3"
	"github.com/wso2/apk/gateway/enforcer/internal/config"
)

// ChoreoAnalytics represents Choreo analytics.
type ChoreoAnalytics struct {
	// Event represents the event.
	Cfg *config.Server
}

// Process processes event and publishes the data
func (c *ChoreoAnalytics) Process(event *v3.HTTPAccessLogEntry) {
	if c.GetEventCategory(event) == EventCategoryFault && c.GetFaultType() == FaultCategoryOther {
		return
	}

	// Add logic to publish the event
	c.extractDataFromEvent(event)
}

// GetEventCategory returns the event category.
func (c *ChoreoAnalytics) GetEventCategory(logEntry *v3.HTTPAccessLogEntry) EventCategory {
	if logEntry.GetResponse() != nil && logEntry.GetResponse().GetResponseCodeDetails() == UpstreamSuccessResponseDetail {
		return EventCategorySuccess
	} else if logEntry.GetResponse() != nil &&
		logEntry.GetResponse().GetResponseCode().GetValue() != 200 &&
		logEntry.GetResponse().GetResponseCode().GetValue() != 204 {
		return EventCategoryFault
	}
	return EventCategoryInvalid
}

// GetFaultType returns the fault type.
func (c *ChoreoAnalytics) GetFaultType() FaultCategory {
	if c.isTargetFaultRequest() {
		return FaultCategoryTargetConnectivity
	}
	return FaultCategoryOther
}

// isTargetFaultRequest checks if the request is a target fault request.
func (c *ChoreoAnalytics) isTargetFaultRequest() bool {
	// Implement the logic to determine if the request is a target fault request.
	return false
}

func (c *ChoreoAnalytics) extractDataFromEvent(logEntry *v3.HTTPAccessLogEntry) {
	c.Cfg.Logger.Info(fmt.Sprintf("log entry metadata, %+v", logEntry.CommonProperties))
}
