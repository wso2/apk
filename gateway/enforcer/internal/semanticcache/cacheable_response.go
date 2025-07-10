package semanticcache

import (
	"regexp"
	"time"
)

// CacheResponse holds everything needed to cache an HTTP response.
type CacheResponse struct {
	ResponsePayload     map[string]interface{}
	RequestHash         string
	Timeout             time.Duration
	HeaderProperties    map[string]interface{}
	StatusCode          string
	StatusReason        string
	JSON                bool
	ProtocolType        string
	HTTPMethod          string
	ResponseCodePattern string
	responseCodeRegex   *regexp.Regexp
	ResponseFetchedTime time.Time
	CacheControlEnabled bool
	AddAgeHeaderEnabled bool
}

// Clean resets payload and headers to free up memory.
func (c *CacheResponse) Clean() {
	c.ResponsePayload = nil
	c.HeaderProperties = nil
}
