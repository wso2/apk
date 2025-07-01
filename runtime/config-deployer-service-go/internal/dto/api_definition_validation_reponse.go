package dto

type APIDefinitionValidationResponse struct {
	IsValid      bool   `json:"isValid"`                // true if valid
	Content      string `json:"content,omitempty"`      // Original content
	JSONContent  string `json:"jsonContent,omitempty"`  // JSON representation
	ProtoContent []byte `json:"protoContent,omitempty"` // Proto file content
	Protocol     string `json:"protocol,omitempty"`     // Protocol type (e.g., HTTP/GRPC)
	Info         Info   `json:"info,omitempty"`         // Parsed info
	IsInit       bool   `json:"isInit"`                 // Init status
}

type Info struct {
	OpenAPIVersion string   `json:"openAPIVersion"` // OpenAPI version (e.g., 3.0.1)
	Name           string   `json:"name"`           // API name
	Version        string   `json:"version"`        // API version
	Context        string   `json:"context"`        // API context path
	Description    string   `json:"description"`    // API description
	Endpoints      []string `json:"endpoints"`      // List of endpoint URLs
}
