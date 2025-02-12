package dto

// ELKResponseEvent represents the ELK response event data.
type ELKResponseEvent struct {
	APIID                    string                 `json:"apiId"`
	APIType                  string                 `json:"apiType"`
	APIName                  string                 `json:"apiName"`
	APIVersion               string                 `json:"apiVersion"`
	APICreator               string                 `json:"apiCreation"` // Added
	APICreatorTenantDomain   string                 `json:"apiCreatorTenantDomain"`
	APIMethod                string                 `json:"apiMethod"`
	APIContext               string                 `json:"apiContext"` // Added
	APIResourceTemplate      string                 `json:"apiResourceTemplate"`
	OrganizationID           string                 `json:"organizationID"`
	EnvironmentID            string                 `json:"environmentID"`
	TargetResponseCode       int                    `json:"targetResponseCode"`
	ProxyResponseCode        int                    `json:"proxyResponseCode"` // Added
	ResponseCacheHit         bool                   `json:"responseCacheHit"`
	Destination              string                 `json:"destination"`
	CorrelationID            string                 `json:"correlationID"`
	RegionID                 string                 `json:"regionID"`
	GatewayType              string                 `json:"gatewayType"`
	ResponseLatency          int64                  `json:"responseLatency"`
	BackendLatency           int64                  `json:"backendLatency"`
	RequestMediationLatency  int64                  `json:"requestMediationLatency"`
	ResponseMediationLatency int64                  `json:"responseMediationLatency"`
	KeyType                  string                 `json:"keyType"`
	ApplicationID            string                 `json:"applicationID"`
	ApplicationName          string                 `json:"applicationName"`
	ApplicationOwner         string                 `json:"applicationOwner"`
	UserAgentHeader          string                 `json:"userAgentHeader"` // Added
	UserName                 string                 `json:"userName"`        // Added
	UserIP                   string                 `json:"userIP"`          // Added
	RequestTimestamp         string                 `json:"requestTimestamp"`
	Properties               map[string]interface{} `json:"properties"`
}

// DefaultFaultEvent represents the default fault event data.
type DefaultFaultEvent struct {
	RequestTimestamp       string                 `json:"requestTimestamp"`
	CorrelationID          string                 `json:"correlationID"`
	KeyType                string                 `json:"keyType"`
	ErrorType              string                 `json:"errorType"`
	ErrorCode              int                    `json:"errorCode"`
	ErrorMessage           string                 `json:"errorMessage"`
	APIID                  string                 `json:"apiId"`
	APIType                string                 `json:"apiType"`
	APIName                string                 `json:"apiName"`
	APIVersion             string                 `json:"apiVersion"`
	APIMethod              string                 `json:"apiMethod"`
	APICreation            string                 `json:"apiCreation"`
	APICreatorTenantDomain string                 `json:"apiCreatorTenantDomain"`
	ApplicationID          string                 `json:"applicationID"`
	ApplicationName        string                 `json:"applicationName"`
	ApplicationOwner       string                 `json:"applicationOwner"`
	RegionID               string                 `json:"regionID"`
	GatewayType            string                 `json:"gatewayType"`
	OrganizationID         string                 `json:"organizationID"`
	EnvironmentID          string                 `json:"environmentID"`
	ProxyResponseCode      int                    `json:"proxyResponseCode"`
	TargetResponseCode     int                    `json:"targetResponseCode"`
	ResponseLatency        int64                  `json:"responseLatency"`
	UserIP                 string                 `json:"userIP"`
	UserAgentHeader        string                 `json:"userAgentHeader"`
	Properties             map[string]interface{} `json:"properties"`
}

// ELKFaultEvent represents the ELK fault event data.
type ELKFaultEvent struct {
	RequestTimestamp       string                 `json:"requestTimestamp"`
	CorrelationID          string                 `json:"correlationID"`
	KeyType                string                 `json:"keyType"`
	ErrorType              string                 `json:"errorType"`
	ErrorCode              int                    `json:"errorCode"`
	ErrorMessage           string                 `json:"errorMessage"`
	APIID                  string                 `json:"apiId"`
	APIType                string                 `json:"apiType"`
	APIName                string                 `json:"apiName"`
	APIVersion             string                 `json:"apiVersion"`
	APIMethod              string                 `json:"apiMethod"`
	APICreation            string                 `json:"apiCreation"`
	APICreatorTenantDomain string                 `json:"apiCreatorTenantDomain"`
	ApplicationID          string                 `json:"applicationID"`
	ApplicationName        string                 `json:"applicationName"`
	ApplicationOwner       string                 `json:"applicationOwner"`
	RegionID               string                 `json:"regionID"`
	GatewayType            string                 `json:"gatewayType"`
	ProxyResponseCode      int                    `json:"proxyResponseCode"`
	TargetResponseCode     int                    `json:"targetResponseCode"`
	UserIP                 string                 `json:"userIP"`
	UserAgentHeader        string                 `json:"userAgentHeader"`
	Properties             map[string]interface{} `json:"properties"`
}

// DefaultResponseEvent represents the default response event data.
type DefaultResponseEvent struct {
	RequestTimestamp         string `json:"requestTimestamp"`
	CorrelationID            string `json:"correlationId"`
	KeyType                  string `json:"keyType"`
	APIID                    string `json:"apiId"`
	APIType                  string `json:"apiType"`
	APIName                  string `json:"apiName"`
	APIVersion               string `json:"apiVersion"`
	APICreator               string `json:"apiCreator"`
	APIMethod                string `json:"apiMethod"`
	APIContext               string `json:"apiContext"`
	APIResourceTemplate      string `json:"apiResourceTemplate"`
	APICreatorTenantDomain   string `json:"apiCreatorTenantDomain"`
	Destination              string `json:"destination"`
	ApplicationID            string `json:"applicationId"`
	ApplicationName          string `json:"applicationName"`
	ApplicationOwner         string `json:"applicationOwner"`
	OrganizationID           string `json:"organizationId"`
	EnvironmentID            string `json:"environmentId"`
	RegionID                 string `json:"regionId"`
	GatewayType              string `json:"gatewayType"`
	UserAgentHeader          string `json:"userAgent"`
	UserName                 string `json:"userName"`
	ProxyResponseCode        int    `json:"proxyResponseCode"`
	TargetResponseCode       int    `json:"targetResponseCode"`
	ResponseCacheHit         bool   `json:"responseCacheHit"`
	ResponseLatency          int64  `json:"responseLatency"`
	BackendLatency           int64  `json:"backendLatency"`
	RequestMediationLatency  int64  `json:"requestMediationLatency"`
	ResponseMediationLatency int64  `json:"responseMediationLatency"`
	UserIP                   string `json:"userIP"`
	EventType                string `json:"eventType"`
	Platform 				string `json:"platform"`

	// Properties               map[string]interface{} `json:"properties"`
}
