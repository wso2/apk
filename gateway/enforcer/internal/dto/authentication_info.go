package dto

// APIKeyAuthenticationInfo represents the API key authentication info.
type APIKeyAuthenticationInfo struct {
	Valid             bool             `json:"valid"` // Valid
	IssuedTime        int64            // Issued time
	ExpiryTime        int64            // Expiry time
	Keytype           string           // Key type
	PermittedReferer  []string         // Permitted referer
	PermittedIP       []string         // Permitted IP
	Application       *ApplicationInfo `json:"application"`       // Application info
	ValidationCode    int              `json:"validationCode"`    // Validation code
	ValidationMessage string           `json:"validationMessage"` // Validation message
}

// ApplicationInfo represents the application info.
type ApplicationInfo struct {
	//ID application id
	ID float64
	//UUID application uuid
	UUID string
}
