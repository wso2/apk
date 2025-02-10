package publishers

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/Azure/azure-amqp-common-go/v4/auth"
	eventhub "github.com/Azure/azure-event-hubs-go/v3"
	"github.com/wso2/apk/gateway/enforcer/internal/analytics/dto"
	"github.com/wso2/apk/gateway/enforcer/internal/config"
	"github.com/wso2/apk/gateway/enforcer/internal/util"
)

// Choreo represents the ELK publisher
type Choreo struct {
	logLevel    string
	cfg         *config.Server
	hub         *eventhub.Hub
	hashedToken string
}

type choreoTokenProvider struct {
	authURL string
	token   string
	cfg     *config.Server
}

type tokenResponse struct {
	Token string `json:"token"`
}

func (c *choreoTokenProvider) GetToken(uri string) (*auth.Token, error) {
	// clientCert, err := util.LoadCertificates(c.cfg.EnforcerPublicKeyPath, c.cfg.EnforcerPrivateKeyPath)
	// if err != nil {
	// 	panic(err)
	// }

	// // Load the trusted CA certificates
	// certPool, err := util.LoadCACertificates(c.cfg.TrustedAdapterCertsPath)
	// if err != nil {
	// 	panic(err)
	// }

	//Create the TLS configuration
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true, // WARNING: This disables certificate verification
	}
	headers := map[string]string{
		"Authorization": "Bearer " + c.token,
	}
	response, err := util.MakeGETRequest(fmt.Sprintf("%s/%s", c.authURL, "token"), tlsConfig, headers)
	if err != nil {
		return nil, err
	}
	var result tokenResponse
	body, _ := ioutil.ReadAll(response.Body)
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	log.Println("Applications: ", result)
	return &auth.Token{
		Token: result.Token,
	}, nil
}

func getResourceURI(sasToken string) string {
	sasAttributes := strings.Split(sasToken, "&")
	resource := strings.Split(sasAttributes[0], "=")
	resourceURI := ""
	if decodedResourceURI, err := url.QueryUnescape(resource[1]); err == nil {
		resourceURI = decodedResourceURI
	}
	return strings.Replace(resourceURI, "sb://", "", 1)
}

func getNamespace(resourceURI string) string {
	ns := strings.Split(resourceURI, "/")[0]
	if strings.Contains(ns, ".") {
		return strings.Split(ns, ".")[0]
	}
	return ns
}

func getEventHubName(resourceURI string) string {
	parts := strings.SplitN(resourceURI, "/", 2)
	if len(parts) > 1 {
		return parts[1]
	}
	return ""
}

// NewChoreo creates a new ELK publisher
func NewChoreo(cfg *config.Server, authURL, token string) *Choreo {
	ctp := &choreoTokenProvider{
		authURL: authURL,
		token:   token,
		cfg:     cfg,
	}
	tokenFromChoreo, err := ctp.GetToken(authURL)

	if err != nil {
		cfg.Logger.Error(err, "Error while getting token from Choreo. Retrying in 5 seconds")
		// Retry after 5 seconds
		time.Sleep(5 * time.Second)
		return NewChoreo(cfg, authURL, token)
	}
	resourceURI := getResourceURI(tokenFromChoreo.Token)
	ns := getNamespace(resourceURI)
	eventHubName := getEventHubName(resourceURI)

	cfg.Logger.Info(fmt.Sprintf("Resource URI: %s", resourceURI))
	cfg.Logger.Info(fmt.Sprintf("Namespace: %s", ns))
	cfg.Logger.Info(fmt.Sprintf("Event Hub Name: %s", eventHubName))

	hub, err := eventhub.NewHub(ns, eventHubName, ctp)
	if err != nil {
		cfg.Logger.Error(err, "Error while creating event hub")
		return nil
	}
	return &Choreo{
		cfg:         cfg,
		hub:         hub,
		hashedToken: util.ComputeSHA256Hash(token),
	}
}

// Publish publishes the event to ELK
func (e *Choreo) Publish(event *dto.Event) {
	e.cfg.Logger.Info(fmt.Sprintf("Publishing event to Choreo: %v", event))
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

func (e *Choreo) publishEvent(event *dto.Event) {
	choreoResponseEvent := &dto.DefaultResponseEvent{
		APIName:                  event.API.APIName,
		APIID:                    event.API.APIID,
		APIType:                  event.API.APIType,
		APIVersion:               event.API.APIVersion,
		OrganizationID:           event.API.OrganizationID,
		EnvironmentID:            event.API.EnvironmentID,
		APICreatorTenantDomain:   event.API.APICreatorTenantDomain,
		APIContext:               event.API.APIContext,
		APIMethod:                event.Operation.APIMethod,
		APIResourceTemplate:      event.Operation.APIResourceTemplate,
		TargetResponseCode:       event.Target.TargetResponseCode,
		ProxyResponseCode:        event.ProxyResponseCode,
		ResponseCacheHit:         event.Target.ResponseCacheHit,
		Destination:              event.Target.Destination,
		CorrelationID:            event.MetaInfo.CorrelationID,
		RegionID:                 event.MetaInfo.RegionID,
		GatewayType:              event.MetaInfo.GatewayType,
		ResponseLatency:          event.Latencies.ResponseLatency,
		BackendLatency:           event.Latencies.BackendLatency,
		RequestMediationLatency:  event.Latencies.RequestMediationLatency,
		ResponseMediationLatency: event.Latencies.ResponseMediationLatency,
		KeyType:                  event.Application.KeyType,
		ApplicationID:            event.Application.ApplicationID,
		ApplicationName:          event.Application.ApplicationName,
		ApplicationOwner:         event.Application.ApplicationOwner,
		UserAgentHeader:          event.UserAgentHeader,
		UserName:                 event.UserName,
		UserIP:                   event.UserIP,
		RequestTimestamp:         event.RequestTimestamp,
		Properties:               event.Properties,
	}

	jsonString, err := util.ToJSONString(choreoResponseEvent)
	if err != nil {
		e.cfg.Logger.Error(err, "Error while converting to JSON string")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	eventFromString := eventhub.NewEventFromString(jsonString)
	if eventFromString.Properties == nil {
		eventFromString.Properties = make(map[string]interface{})
	}
	eventFromString.Properties["token-hash"] = e.hashedToken
	e.cfg.Logger.Info(fmt.Sprintf("Event from string: %+v", eventFromString))
	// err = e.hub.Send(ctx, eventFromString)
	// if err != nil {
	// 	e.cfg.Logger.Error(err, "Error while sending event to Choreo")
	// }

	var events []*eventhub.Event
	events = append(events, eventFromString)

	err = e.hub.SendBatch(ctx, eventhub.NewEventBatchIterator(events...))
	if err != nil {
		e.cfg.Logger.Error(err, "Error while sending event to Choreo")
	}
}

func (e *Choreo) publishFault(event *dto.Event) {
	choreoResponseEvent := &dto.DefaultFaultEvent{
		APIName:                event.API.APIName,
		APIID:                  event.API.APIID,
		APIType:                event.API.APIType,
		APIVersion:             event.API.APIVersion,
		APICreatorTenantDomain: event.API.APICreatorTenantDomain,
		APIMethod:              event.Operation.APIMethod,
		TargetResponseCode:     event.Target.TargetResponseCode,
		ProxyResponseCode:      event.ProxyResponseCode,
		CorrelationID:          event.MetaInfo.CorrelationID,
		RegionID:               event.MetaInfo.RegionID,
		GatewayType:            event.MetaInfo.GatewayType,
		KeyType:                event.Application.KeyType,
		ApplicationID:          event.Application.ApplicationID,
		ApplicationName:        event.Application.ApplicationName,
		ApplicationOwner:       event.Application.ApplicationOwner,
		UserAgentHeader:        event.UserAgentHeader,
		UserIP:                 event.UserIP,
		RequestTimestamp:       event.RequestTimestamp,
		Properties:             event.Properties,
		ErrorType:              "",
		ErrorCode:              event.Target.TargetResponseCode,
		ErrorMessage:           event.Target.ResponseCodeDetail,
	}

	jsonString, err := util.ToJSONString(choreoResponseEvent)
	if err != nil {
		e.cfg.Logger.Error(err, "Error while converting to JSON string")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	eventFromString := eventhub.NewEventFromString(jsonString)
	e.cfg.Logger.Info(fmt.Sprintf("Event from string: %+v", eventFromString))
	err = e.hub.Send(ctx, eventFromString)
	if err != nil {
		e.cfg.Logger.Error(err, "Error while sending event to Choreo")
	}
}

func (e *Choreo) isFault(event *dto.Event) bool {
	return event.Target.ResponseCodeDetail != "via_upstream"
}
