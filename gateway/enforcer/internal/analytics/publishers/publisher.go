package publishers

import "github.com/wso2/apk/gateway/enforcer/internal/analytics/dto"

type Publisher interface {
	Publish(event *dto.Event)
}
