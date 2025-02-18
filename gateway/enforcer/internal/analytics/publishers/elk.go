package publishers

import (
	"fmt"

	"github.com/wso2/apk/gateway/enforcer/internal/analytics/dto"
	"github.com/wso2/apk/gateway/enforcer/internal/config"
	"github.com/wso2/apk/gateway/enforcer/internal/util"
)

// ELK represents the ELK publisher
type ELK struct {
	logLevel string
	cfg      *config.Server
}

// NewELK creates a new ELK publisher
func NewELK(cfg *config.Server, logLevel string) *ELK {
	return &ELK{
		logLevel: logLevel,
		cfg:      cfg,
	}
}

// Publish publishes the event to ELK
func (e *ELK) Publish(event *dto.Event) {
	e.cfg.Logger.Sugar().Debug(fmt.Sprintf("Publishing event to ELK: %v", event))
	defer func() {
		if r := recover(); r != nil {
			e.cfg.Logger.Error(nil, fmt.Sprintf("Recovered from panic: %v", r))
		}
	}()
	// Implement the ELK publish logic
	if e.isFault(event) {
		e.publishFault(event)
	} else {
		e.publishEvent(event)
	}
}

func (e *ELK) publishEvent(event *dto.Event) {
	elkResponseEvent := &dto.ELKResponseEvent{
		APIName: event.API.APIName,
		APIID:   event.API.APIID,
		APIType: event.API.APIType,
		APIVersion: event.API.APIVersion,
		OrganizationID: event.API.OrganizationID,
		EnvironmentID: event.API.EnvironmentID,
		APICreator: event.API.APICreator,
		APICreatorTenantDomain: event.API.APICreatorTenantDomain,
		APIContext: event.API.APIContext,
		APIMethod: event.Operation.APIMethod,
		APIResourceTemplate: event.Operation.APIResourceTemplate,
		TargetResponseCode: event.Target.TargetResponseCode,
		ProxyResponseCode: event.ProxyResponseCode,
		ResponseCacheHit: event.Target.ResponseCacheHit,
		Destination: event.Target.Destination,
		CorrelationID: event.MetaInfo.CorrelationID,
		RegionID: event.MetaInfo.RegionID,
		GatewayType: event.MetaInfo.GatewayType,
		ResponseLatency: event.Latencies.ResponseLatency,
		BackendLatency: event.Latencies.BackendLatency,
		RequestMediationLatency: event.Latencies.RequestMediationLatency,
		ResponseMediationLatency: event.Latencies.ResponseMediationLatency,
		KeyType: event.Application.KeyType,
		ApplicationID: event.Application.ApplicationID,
		ApplicationName: event.Application.ApplicationName,
		ApplicationOwner: event.Application.ApplicationOwner,
		UserAgentHeader: event.UserAgentHeader,
		UserName: event.UserName,
		UserIP: event.UserIP,
		RequestTimestamp: event.RequestTimestamp,
		Properties: event.Properties,
	}

	jsonString, err := util.ToJSONString(elkResponseEvent)
	if err != nil {
		e.cfg.Logger.Error(err, "Error while converting to JSON string")
		return
	}
	e.cfg.Logger.Info(fmt.Sprintf("apimMetrics: %s, properties: %s", "apim:response", jsonString))
}

func (e *ELK) publishFault(event *dto.Event) {
	elkResponseEvent := &dto.ELKFaultEvent{
		APIName: event.API.APIName,
		APIID:   event.API.APIID,
		APIType: event.API.APIType,
		APIVersion: event.API.APIVersion,
		APICreatorTenantDomain: event.API.APICreatorTenantDomain,
		APIMethod: event.Operation.APIMethod,
		TargetResponseCode: event.Target.TargetResponseCode,
		ProxyResponseCode: event.ProxyResponseCode,
		CorrelationID: event.MetaInfo.CorrelationID,
		RegionID: event.MetaInfo.RegionID,
		GatewayType: event.MetaInfo.GatewayType,
		KeyType: event.Application.KeyType,
		ApplicationID: event.Application.ApplicationID,
		ApplicationName: event.Application.ApplicationName,
		ApplicationOwner: event.Application.ApplicationOwner,
		UserAgentHeader: event.UserAgentHeader,
		UserIP: event.UserIP,
		RequestTimestamp: event.RequestTimestamp,
		Properties: event.Properties,
		ErrorType: "",
		ErrorCode: event.Target.TargetResponseCode,
		ErrorMessage: event.Target.ResponseCodeDetail,
	}

	jsonString, err := util.ToJSONString(elkResponseEvent)
	if err != nil {
		e.cfg.Logger.Error(err, "Error while converting to JSON string")
		return
	}
	e.cfg.Logger.Info(fmt.Sprintf("apimMetrics: %s, properties: %s", "apim:faulty", jsonString))
}

func (e *ELK) isFault(event *dto.Event) bool {
	return event.Target.ResponseCodeDetail != "via_upstream"
}