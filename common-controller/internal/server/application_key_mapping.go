package server

// ApplicationKeyMapping defines the desired state of ApplicationKeyMapping
type ApplicationKeyMapping struct {
	ApplicationUUID       string `json:"applicationUUID,omitempty"`
	SecurityScheme        string `json:"securityScheme,omitempty"`
	ApplicationIdentifier string `json:"applicationIdentifier,omitempty"`
	KeyType               string `json:"keyType,omitempty"`
	EnvID                 string `json:"envID,omitempty"`
}

// ApplicationKeyMappingList contains a list of ApplicationKeyMapping
type ApplicationKeyMappingList struct {
	List []ApplicationKeyMapping `json:"list"`
}
