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

// EventCategory represents the category of an event.
type EventCategory string
// FaultCategory represents the category of a fault.
type FaultCategory string

const (
	// EventCategorySuccess represents a successful event.
	EventCategorySuccess            EventCategory = "SUCCESS"
	// EventCategoryFault represents a fault event.
	EventCategoryFault              EventCategory = "FAULT"
	// EventCategoryInvalid represents an invalid event.
	EventCategoryInvalid            EventCategory = "INVALID"
	// FaultCategoryTargetConnectivity represents a target connectivity fault.
	FaultCategoryTargetConnectivity FaultCategory = "TARGET_CONNECTIVITY"
	// FaultCategoryOther represents other faults.
	FaultCategoryOther              FaultCategory = "OTHER"
)
