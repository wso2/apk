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

}
