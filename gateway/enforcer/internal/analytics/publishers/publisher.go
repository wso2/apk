package publishers

import "github.com/wso2/apk/gateway/enforcer/internal/analytics/dto"

// Publisher represents an analytics publisher.
type Publisher interface {
	Publish(event *dto.Event)
}
