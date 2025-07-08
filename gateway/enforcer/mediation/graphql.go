package mediation

import (
	dpv2alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v2alpha1"
	"github.com/wso2/apk/gateway/enforcer/internal/requestconfig"
)

// GraphQL represents the configuration for GraphQL policy in the API Gateway.
type GraphQL struct {
	PolicyName    string `json:"policyName"`
	PolicyVersion string `json:"policyVersion"`
	PolicyID      string `json:"policyID"`
	Schema        string `json:"schema"`
}

const (
	// GraphQLPolicyKeySchema is the key for specifying the GraphQL schema.
	GraphQLPolicyKeySchema = "Schema"
)

// NewGraphQL creates a new GraphQL instance with default values.
func NewGraphQL(mediation *dpv2alpha1.Mediation) *GraphQL {
	schema := ""
	if val, ok := extractPolicyValue(mediation.Parameters, GraphQLPolicyKeySchema); ok {
		schema = val
	}
	return &GraphQL{
		PolicyName:    "GraphQL",
		PolicyVersion: "v1",
		PolicyID:      "graphql",
		Schema:        schema,
	}
}

// Process processes the request configuration for GraphQL.
func (g *GraphQL) Process(requestConfig *requestconfig.Holder) *Result {
	// Implement the logic to process the requestConfig for GraphQL
	// This is a placeholder implementation
	result := &Result{}

	// Add logic to handle GraphQL schema processing here

	return result
}
