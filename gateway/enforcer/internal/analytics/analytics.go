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
	"github.com/wso2/apk/gateway/enforcer/internal/analytics/dto"
	"github.com/wso2/apk/gateway/enforcer/internal/config"
	"github.com/wso2/apk/gateway/enforcer/internal/datastore"
)

// EventCategory represents the category of an event.
type EventCategory string

// FaultCategory represents the category of a fault.
type FaultCategory string

const (
	// EventCategorySuccess represents a successful event.
	EventCategorySuccess EventCategory = "SUCCESS"
	// EventCategoryFault represents a fault event.
	EventCategoryFault EventCategory = "FAULT"
	// EventCategoryInvalid represents an invalid event.
	EventCategoryInvalid EventCategory = "INVALID"
	// FaultCategoryTargetConnectivity represents a target connectivity fault.
	FaultCategoryTargetConnectivity FaultCategory = "TARGET_CONNECTIVITY"
	// FaultCategoryOther represents other faults.
	FaultCategoryOther FaultCategory = "OTHER"
)

// Analytics represents Choreo analytics.
type Analytics struct {
	// Cfg represents the server configuration.
	Cfg         *config.Server
	// ConfigStore represents the configuration store.
	ConfigStore *datastore.ConfigStore
}

// Process processes event and publishes the data
func (c *Analytics) Process(event *v3.HTTPAccessLogEntry) {
	if c.GetEventCategory(event) == EventCategoryFault && c.GetFaultType() == FaultCategoryOther {
		return
	}

	// Add logic to publish the event
	_ = c.prepareAnalyticEvent(event)
}

// GetEventCategory returns the event category.
func (c *Analytics) GetEventCategory(logEntry *v3.HTTPAccessLogEntry) EventCategory {
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
func (c *Analytics) GetFaultType() FaultCategory {
	if c.isTargetFaultRequest() {
		return FaultCategoryTargetConnectivity
	}
	return FaultCategoryOther
}

// isTargetFaultRequest checks if the request is a target fault request.
func (c *Analytics) isTargetFaultRequest() bool {
	// Implement the logic to determine if the request is a target fault request.
	return false
}

func (c *Analytics) prepareAnalyticEvent(logEntry *v3.HTTPAccessLogEntry) *dto.Event {
	keyValuePairsFromMetadata := make(map[string]string)
	c.Cfg.Logger.Info(fmt.Sprintf("log entry metadata, %+v", logEntry.CommonProperties))
	if logEntry.CommonProperties != nil && logEntry.CommonProperties.Metadata != nil && logEntry.CommonProperties.Metadata.FilterMetadata != nil {
		if sv, exists := logEntry.CommonProperties.Metadata.FilterMetadata[ExtProcMetadataContextKey]; exists {
			if sv.Fields != nil {
				for key, value := range sv.Fields {
					if value != nil {
						keyValuePairsFromMetadata[key] = value.GetStringValue()
					}
				}
			}
		}
	}
	// Prepare extended API
	extendedAPI := dto.ExtendedAPI{}
	extendedAPI.APIType = keyValuePairsFromMetadata[APITypeKey]
	extendedAPI.APIID = keyValuePairsFromMetadata[APIIDKey]
	extendedAPI.APICreator = keyValuePairsFromMetadata[APICreatorKey]
	extendedAPI.APIName = keyValuePairsFromMetadata[APINameKey]
	extendedAPI.APIVersion = keyValuePairsFromMetadata[APIVersionKey]
	extendedAPI.APICreatorTenantDomain = keyValuePairsFromMetadata[APICreatorTenantDomainKey]
	extendedAPI.OrganizationID = keyValuePairsFromMetadata[APIOrganizationIDKey]
	extendedAPI.APIContext = keyValuePairsFromMetadata[APIContextKey]
	extendedAPI.EnvironmentID = keyValuePairsFromMetadata[APIEnvironmentKey]

	// Prepare operation
	operation := dto.Operation{}
	operation.APIResourceTemplate = keyValuePairsFromMetadata[APIResourceTemplateKey]
	operation.APIMethod = keyValuePairsFromMetadata[logEntry.Request.GetRequestMethod().String()]

	// Prepare target
	target := dto.Target{}
	target.ResponseCacheHit = false
	target.TargetResponseCode = int(logEntry.GetResponse().GetResponseCode().Value)
	target.Destination = keyValuePairsFromMetadata[DestinationKey]

	// Prepare Application
	application := &dto.Application{}
	if keyValuePairsFromMetadata[AppIDKey] == Unknown {
		application = c.getAnonymousApp()
	} else {
		application.ApplicationID = keyValuePairsFromMetadata[AppIDKey]
		application.KeyType = keyValuePairsFromMetadata[AppKeyTypeKey]
		application.ApplicationName = keyValuePairsFromMetadata[AppNameKey]
		application.ApplicationOwner = keyValuePairsFromMetadata[AppOwnerKey]
	}

	properties := logEntry.GetCommonProperties()
	backendResponseRecvTimestamp :=
		(properties.TimeToLastUpstreamRxByte.Seconds * 1000) +
			(int64(properties.TimeToLastUpstreamRxByte.Nanos) / 1_000_000)

	backendRequestSendTimestamp :=
		(properties.TimeToFirstUpstreamTxByte.Seconds * 1000) +
			(int64(properties.TimeToFirstUpstreamTxByte.Nanos) / 1_000_000)

	downstreamResponseSendTimestamp :=
		(properties.TimeToLastDownstreamTxByte.Seconds * 1000) +
			(int64(properties.TimeToLastDownstreamTxByte.Nanos) / 1_000_000)

	// Prepare Latencies
	latencies := dto.Latencies{}
	latencies.BackendLatency = backendResponseRecvTimestamp - backendRequestSendTimestamp
	latencies.RequestMediationLatency = backendRequestSendTimestamp
	latencies.ResponseLatency = downstreamResponseSendTimestamp
	latencies.ResponseMediationLatency = downstreamResponseSendTimestamp - backendResponseRecvTimestamp

	// prepare metaInfo
	metaInfo := dto.MetaInfo{}
	metaInfo.CorrelationID = keyValuePairsFromMetadata[CorrelationIDKey]
	metaInfo.RegionID = keyValuePairsFromMetadata[RegionKey]

	userAgent := logEntry.GetRequest().GetUserAgent()
	userName := keyValuePairsFromMetadata[APIUserNameKey]
	userIP := logEntry.GetCommonProperties().GetDownstreamRemoteAddress().GetSocketAddress().GetAddress()
	if userIP == "" {
		userIP = Unknown
	}
	if userAgent == "" {
		userAgent = Unknown
	}

	event := &dto.Event{}
	event.MetaInfo = &metaInfo
	event.API = &extendedAPI
	event.Operation = &operation
	event.Target = &target
	event.Application = application
	event.Latencies = &latencies
	event.UserAgentHeader = userAgent
	event.UserName = userName
	event.UserIP = userIP
	event.ProxyResponseCode = int(logEntry.GetResponse().GetResponseCode().Value)
	event.RequestTimestamp = logEntry.GetCommonProperties().GetStartTime().String()

	return event
}

func (c *Analytics) getAnonymousApp() *dto.Application {
	application := &dto.Application{}
	application.ApplicationID = anonymousValye
	application.ApplicationName = anonymousValye
	application.KeyType = anonymousValye
	application.ApplicationOwner = anonymousValye
	return application
}
