package mediation

// ExternalInput defines the JSON contract sent to the plugin runner.
// Keep this struct JSON-friendly and stable across components.
type ExternalInput struct {
	Phase           string            `json:"phase"`
	PolicyName      string            `json:"policyName"`
	PolicyID        string            `json:"policyID,omitempty"`
	PolicyVersion   string            `json:"policyVersion,omitempty"`
	Parameters      map[string]string `json:"parameters,omitempty"`
	Attributes      map[string]string `json:"attributes,omitempty"`
	RequestHeaders  map[string]string `json:"requestHeaders,omitempty"`
	ResponseHeaders map[string]string `json:"responseHeaders,omitempty"`
	RequestBody     string            `json:"requestBody,omitempty"`
	ResponseBody    string            `json:"responseBody,omitempty"`
}

// ExternalOutput defines the JSON contract returned by the plugin runner.
// This mirrors the gateway Result fields using plain JSON types.
type ExternalOutput struct {
	AddHeaders                   map[string]string      `json:"addHeaders,omitempty"`
	RemoveHeaders                []string               `json:"removeHeaders,omitempty"`
	ModifyBody                   bool                   `json:"modifyBody,omitempty"`
	Body                         string                 `json:"body,omitempty"`
	ImmediateResponse            bool                   `json:"immediateResponse,omitempty"`
	ImmediateResponseCode        int32                  `json:"immediateResponseCode,omitempty"`
	ImmediateResponseBody        string                 `json:"immediateResponseBody,omitempty"`
	ImmediateResponseDetail      string                 `json:"immediateResponseDetail,omitempty"`
	ImmediateResponseHeaders     map[string]string      `json:"immediateResponseHeaders,omitempty"`
	ImmediateResponseContentType string                 `json:"immediateResponseContentType,omitempty"`
	StopFurtherProcessing        bool                   `json:"stopFurtherProcessing,omitempty"`
	Metadata                     map[string]interface{} `json:"metadata,omitempty"`
}
