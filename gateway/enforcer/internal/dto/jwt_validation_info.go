package dto

// JWTValidationInfo represents the JWT validation info
type JWTValidationInfo struct {
	Issuer    string                 `json:"issuer"`    // Issuer
	ClientID  string                 `json:"clientId"`  // Client ID
	Subject   string                 `json:"subject"`   // Subject
	Audiences *[]string              `json:"audiences"` // Audiences
	Scopes    *[]string              `json:"scopes"`    // Scopes
	Claims    map[string]interface{} `json:"claims"`    // Claims
}
