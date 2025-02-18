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
	"strings"

	v3 "github.com/envoyproxy/go-control-plane/envoy/data/accesslog/v3"
	"github.com/wso2/apk/gateway/enforcer/internal/analytics/dto"
	analytics_publisher "github.com/wso2/apk/gateway/enforcer/internal/analytics/publishers"
	"github.com/wso2/apk/gateway/enforcer/internal/config"
	"github.com/wso2/apk/gateway/enforcer/internal/datastore"
)

// EventCategory represents the category of an event.
type EventCategory string

// FaultCategory represents the category of a fault.
type FaultCategory string

// RFC3339Millis represents the RFC3339 date format with milliseconds.
const RFC3339Millis = "2006-01-02T15:04:05.000Z07:00"

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
	// DefaultAnalyticsPublisher represents the default analytics publisher.
	DefaultAnalyticsPublisher = "default"
	// MoesifAnalyticsPublisher represents the Moesif analytics publisher.
	MoesifAnalyticsPublisher = "moesif"
	// ELKAnalyticsPublisher represents the ELK analytics publisher.
	ELKAnalyticsPublisher = "elk"

	// PromptTokenCountMetadataKey represents the prompt token count metadata key.
	PromptTokenCountMetadataKey string = "aitoken:prompttokencount"
	// CompletionTokenCountMetadataKey represents the completion token count metadata key.
	CompletionTokenCountMetadataKey string = "aitoken:completiontokencount"
	// TotalTokenCountMetadataKey represents the total token count metadata key.
	TotalTokenCountMetadataKey string = "aitoken:totaltokencount"
	// ModelIDMetadataKey represents the model name metadata key.
	ModelIDMetadataKey string = "aitoken:modelid"

	// AIProviderNameMetadataKey represents the AI provider metadata key.
	AIProviderNameMetadataKey string = "ai:providername"
	// AIProviderAPIVersionMetadataKey represents the AI provider API version metadata key.
	AIProviderAPIVersionMetadataKey string = "ai:providerversion"
)

// Analytics represents Choreo analytics.
type Analytics struct {
	// cfg represents the server configuration.
	cfg *config.Server
	// configStore represents the configuration store.
	configStore *datastore.ConfigStore
	// publishers represents the publishers.
	publishers []analytics_publisher.Publisher
}

// NewAnalytics creates a new instance of Analytics.
func NewAnalytics(cfg *config.Server, configStore *datastore.ConfigStore) *Analytics {
	publishers := make([]analytics_publisher.Publisher, 0)
	if len(configStore.GetConfigs()) != 0 {
		config := configStore.GetConfigs()[0]
		if config.Analytics.Enabled {
			for _, pub := range config.Analytics.AnalyticsPublisher {
				cfg.Logger.Sugar().Debug(fmt.Sprintf("Publisher type: %s", pub.Type))
				switch strings.ToLower(pub.Type) {
				case strings.ToLower(ELKAnalyticsPublisher):
					logLevel := "INFO"
					if level, exists := pub.ConfigProperties["logLevel"]; exists {
						logLevel = level
					}
					publishers = append(publishers, analytics_publisher.NewELK(cfg, logLevel))
					cfg.Logger.Sugar().Debug(fmt.Sprintf("ELK publisher added with log level: %s", logLevel))
				case strings.ToLower(MoesifAnalyticsPublisher):
					// publisher := publishers.NewMoesif(cfg, pub.LogLevel)
				case strings.ToLower(DefaultAnalyticsPublisher):
					publisher := analytics_publisher.NewChoreo(cfg, cfg.ChoreoAnalyticsAuthURL, cfg.ChoreoAnalyticsAuthToken)
					if publisher == nil {
						cfg.Logger.Error(nil, "Error while creating Choreo publisher")
					} else {
						publishers = append(publishers, publisher)
						cfg.Logger.Info("Choreo publisher added")
					}
				}
			}
		}
	}
	if len(publishers) == 0 {
		cfg.Logger.Sugar().Debug("No analytics publishers found. Analytics will not be published.")
	}
	return &Analytics{
		cfg:         cfg,
		configStore: configStore,
		publishers:  publishers,
	}
}

// Process processes event and publishes the data
func (c *Analytics) Process(event *v3.HTTPAccessLogEntry) {
	if c.isInvalid(event) {
		c.cfg.Logger.Error(nil, "Invalid event received from the access log service")
		return
	}

	// Add logic to publish the event
	analyticEvent := c.prepareAnalyticEvent(event)
	for _, publisher := range c.publishers {
		publisher.Publish(analyticEvent)
	}

}

// GetEventCategory returns the event category.
func (c *Analytics) isInvalid(logEntry *v3.HTTPAccessLogEntry) bool {
	return logEntry.GetResponse() == nil
}

// GetFaultType returns the fault type.
func (c *Analytics) GetFaultType() FaultCategory {
	return FaultCategoryOther
}

func (c *Analytics) prepareAnalyticEvent(logEntry *v3.HTTPAccessLogEntry) *dto.Event {
	keyValuePairsFromMetadata := make(map[string]string)
	c.cfg.Logger.Sugar().Debug(fmt.Sprintf("log entry, %+v", logEntry))
	if logEntry.CommonProperties != nil && logEntry.CommonProperties.Metadata != nil && logEntry.CommonProperties.Metadata.FilterMetadata != nil {
		if sv, exists := logEntry.CommonProperties.Metadata.FilterMetadata[ExtProcMetadataContextKey]; exists {
			if sv.Fields != nil {
				c.cfg.Logger.Sugar().Debug(fmt.Sprintf("Filter metadata: %+v", sv))
				for key, value := range sv.Fields {
					if value != nil {
						keyValuePairsFromMetadata[key] = value.GetStringValue()
					}
				}
			}
		}
	}
	event := &dto.Event{}
	for key, value := range keyValuePairsFromMetadata {
		c.cfg.Logger.Sugar().Debug(fmt.Sprintf("Metadata key: %s, value: %s", key, value))
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
	operation.APIMethod = logEntry.Request.GetRequestMethod().String()

	// Prepare target
	target := dto.Target{}
	target.ResponseCacheHit = false
	target.TargetResponseCode = int(logEntry.GetResponse().GetResponseCode().Value)
	target.Destination = keyValuePairsFromMetadata[DestinationKey]
	target.ResponseCodeDetail = logEntry.GetResponse().GetResponseCodeDetails()

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
	if properties != nil && properties.TimeToLastUpstreamRxByte != nil && properties.TimeToFirstUpstreamTxByte != nil && properties.TimeToLastDownstreamTxByte != nil {
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
		event.Latencies = &latencies
	}

	// prepare metaInfo
	metaInfo := dto.MetaInfo{}
	metaInfo.CorrelationID = logEntry.GetRequest().RequestId
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

	event.MetaInfo = &metaInfo
	event.API = &extendedAPI
	event.Operation = &operation
	event.Target = &target
	event.Application = application
	event.UserAgentHeader = userAgent
	event.UserName = userName
	event.UserIP = userIP
	event.ProxyResponseCode = int(logEntry.GetResponse().GetResponseCode().Value)
	event.RequestTimestamp = logEntry.GetCommonProperties().GetStartTime().AsTime().Format(RFC3339Millis)
	event.Properties = make(map[string]interface{}, 0)

	aiMetadata := dto.AIMetadata{}
	aiMetadata.VendorName = keyValuePairsFromMetadata[AIProviderNameMetadataKey]
	aiMetadata.VendorVersion = keyValuePairsFromMetadata[AIProviderAPIVersionMetadataKey]
	aiMetadata.Model = keyValuePairsFromMetadata[ModelIDMetadataKey]
	event.Properties["aiMetadata"] = aiMetadata

	aiTokenUsage := dto.AITokenUsage{}
	aiTokenUsage.PromptToken = keyValuePairsFromMetadata[PromptTokenCountMetadataKey]
	aiTokenUsage.CompletionToken = keyValuePairsFromMetadata[CompletionTokenCountMetadataKey]
	aiTokenUsage.TotalToken = keyValuePairsFromMetadata[TotalTokenCountMetadataKey]
	event.Properties["aiTokenUsage"] = aiTokenUsage
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
