package model

type API struct {
	Name              string   `json:"name"`              // API name
	BasePath          string   `json:"basePath"`          // API base path
	Version           string   `json:"version"`           // API version
	Type              string   `json:"type"`              // API type (e.g., REST, GraphQL)
	Endpoint          string   `json:"endpoint"`          // Endpoint URL
	APISecurity       string   `json:"apiSecurity"`       // Security definition
	Scopes            []string `json:"scopes"`            // Array of scopes
	GraphQLSchema     string   `json:"graphQLSchema"`     // GraphQL schema string
	ProtoDefinition   string   `json:"protoDefinition"`   // gRPC proto content
	SwaggerDefinition string   `json:"swaggerDefinition"` // Swagger/OpenAPI content
	Environment       string   `json:"environment"`       // Deployment environment
}
