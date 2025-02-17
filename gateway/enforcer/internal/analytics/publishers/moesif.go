package publishers

import (
	"fmt"
	"sync"
	"time"

	"github.com/moesif/moesifapi-go"
	"github.com/moesif/moesifapi-go/models"
	"github.com/wso2/apk/gateway/enforcer/internal/analytics/dto"
	"github.com/wso2/apk/gateway/enforcer/internal/config"
)

const (
	anonymous = "anonymous"
)

// Moesif represents a Moesif publisher.
type Moesif struct {
	cfg *config.Server
	api    moesifapi.API
	events []*models.EventModel
	mu     sync.Mutex
}

// NewMoesif creates a new Moesif publisher.
func NewMoesif(cfg *config.Server) *Moesif {

	apiClient := moesifapi.NewAPI(cfg.MoesifToken)
	moesif := &Moesif{
		cfg: cfg,
		events: []*models.EventModel{},
		api:    apiClient,
		mu:     sync.Mutex{},
	}
	go func() {
		for {
			time.Sleep(time.Duration(cfg.MoesifPublishInterval) * time.Second)
			moesif.mu.Lock()
			if len(moesif.events) > 0 {
				moesif.cfg.Logger.Info(fmt.Sprintf("Publishing %d events to Moesif", len(moesif.events)))
				err := moesif.api.QueueEvents(moesif.events)
				if err != nil {
					moesif.cfg.Logger.Error(err, "Error publishing events to Moesif")
				}
				moesif.events = []*models.EventModel{}
			}
			moesif.mu.Unlock()
		}
	}()
	return moesif
}

// Publish publishes an event to Moesif.
func (m *Moesif) Publish(event *dto.Event) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cfg.Logger.Info("Publishing event to Moesif")
	m.cfg.Logger.Info(fmt.Sprintf("Event: %+v", event))
	uri := event.API.APIContext + event.Operation.APIResourceTemplate
	if event.Target.Destination != "" {
		uri = event.Target.Destination
	}
	
	req := models.EventRequestModel{
		Time:       &event.RequestTimestamp,
		Uri:        uri,
		Verb:       "GET",
		ApiVersion: &event.API.APIVersion,
		IpAddress:  &event.UserIP,
		Headers: map[string]interface{}{
			"User-Agent":   event.UserAgentHeader,
			"Content-Type": "application/json",
		},
		Body: nil,
	}
	respTime := event.RequestTimestamp
	if event.Latencies != nil {
		respTime = event.RequestTimestamp.Add(time.Duration(event.Latencies.ResponseLatency) * time.Millisecond)
	}
	
	rspHeaders := map[string]string{
		"Vary":          "Accept-Encoding",
		"Pragma":        "no-cache",
		"Expires":       "-1",
		"Content-Type":  "application/json; charset=utf-8",
		"Cache-Control": "no-cache",
	}

	rsp := models.EventResponseModel{
		Time:    &respTime,
		Status:  event.ProxyResponseCode,
		Headers: rspHeaders,
	}

	userID := anonymous
	eventModel := &models.EventModel{
		Request:  req,
		Response: rsp,
		UserId:   &userID,
	}
	m.events = append(m.events, eventModel)
	m.cfg.Logger.Info(fmt.Sprintf("Event added to the queue. Queue size: %d", len(m.events)))
}
