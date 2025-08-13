package graphql

// Operation represents a GraphQL operation with its target, verb, and required scopes.
type Operation struct {
	// Target is the target of the operation, e.g., "query", "mutation", etc.
	Target string   `yaml:"target" json:"target"`
	// Verb is the HTTP verb associated with the operation, e.g., "GET", "POST", etc.
	Verb   string   `yaml:"verb" json:"verb"`
	// Scopes are the required scopes for the operation.
	Scopes []string `yaml:"scopes" json:"scopes"`
}