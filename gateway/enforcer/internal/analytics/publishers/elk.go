package publishers

import (
	"github.com/wso2/apk/gateway/enforcer/internal/analytics/dto"
	"github.com/wso2/apk/gateway/enforcer/internal/config"
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
	// Implement the ELK publish logic
	if e.isFault(event) {
		e.publishFault(event)
	} else {
		e.publishEvent(event)
	}
}

func (e *ELK) publishEvent(event *dto.Event) {
	// Implement the ELK publish event logic
}

func (e *ELK) publishFault(event *dto.Event) {
	// Implement the ELK publish fault logic
}

func (e *ELK) isFault(event *dto.Event) bool {
	return event.Target.ResponseCodeDetail != "via_upstream"
}