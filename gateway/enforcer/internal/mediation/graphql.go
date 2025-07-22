package mediation

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/parser"
	"github.com/vektah/gqlparser/v2/validator"
	dpv2alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v2alpha1"
	"github.com/wso2/apk/common-go-libs/constants"
	"github.com/wso2/apk/gateway/enforcer/internal/config"
	"github.com/wso2/apk/gateway/enforcer/internal/dto"
	"github.com/wso2/apk/gateway/enforcer/internal/logging"
	"github.com/wso2/apk/gateway/enforcer/internal/requestconfig"
	"gopkg.in/yaml.v2"
)

// GraphQL represents the configuration for GraphQL policy in the API Gateway.
type GraphQL struct {
	PolicyName    string `json:"policyName"`
	PolicyVersion string `json:"policyVersion"`
	PolicyID      string `json:"policyID"`
	Operations    string `json:"schema"`
	logger        *logging.Logger
	cfg           *config.Server
}

// Operation represents a GraphQL operation with its target, verb, and required scopes.
type Operation struct {
	// Target is the target of the operation, e.g., "query", "mutation", etc.
	Target string   `yaml:"target" json:"target"`
	// Verb is the HTTP verb associated with the operation, e.g., "GET", "POST", etc.
	Verb   string   `yaml:"verb" json:"verb"`
	// Scopes are the required scopes for the operation.
	Scopes []string `yaml:"scopes" json:"scopes"`
}

const (
	// GraphQLPolicyKeySchema is the key for specifying the GraphQL schema.
	GraphQLPolicyKeySchema = "Schema"
)

// NewGraphQL creates a new GraphQL instance with default values.
func NewGraphQL(mediation *dpv2alpha1.Mediation) *GraphQL {
	Operations := ""
	if val, ok := extractPolicyValue(mediation.Parameters, GraphQLPolicyKeySchema); ok {
		Operations = val
	}
	cfg := config.GetConfig()
	logger := cfg.Logger
	return &GraphQL{
		PolicyName:    "GraphQL",
		PolicyVersion: "v1",
		PolicyID:      "graphql",
		Operations:    Operations,
		logger:        &logger,
		cfg:           cfg,
	}
}

// Process processes the request configuration for GraphQL.
func (g *GraphQL) Process(requestConfig *requestconfig.Holder) *Result {
	// Implement the logic to process the requestConfig for GraphQL
	// This is a placeholder implementation
	result := NewResult()
	operations, err := parseOperations([]byte(g.Operations))
	if err != nil {
		g.logger.Sugar().Errorf("failed to parse operations: %v", err)
		result.ImmediateResponse = true
		result.ImmediateResponseBody = "failed to parse operations: " + err.Error()
		result.ImmediateResponseCode = 500
		result.ImmediateResponseDetail = "failed to parse operations: " + err.Error()
		return result
	}

	schemaBytes := []byte(requestConfig.RouteMetadata.Spec.API.Definition)
	var sdl string
	schemaString := ""
	if schemaString, err = unzipGzip(schemaBytes); err != nil {
		g.logger.Sugar().Errorf("error while unzipping the GraphQL SDL: %v, Sending internal server error", err)
		result.ImmediateResponse = true
		result.ImmediateResponseBody = "Related API definition is not configured properly"
		result.ImmediateResponseCode = 500
		result.ImmediateResponseDetail = "error while unzipping the GraphQL SDL"
		return result
	}
	sdl = schemaString

	if sdl == "" {
		g.logger.Sugar().Error("GraphQL SDL is empty, Sending internal server error")
		result.ImmediateResponse = true
		result.ImmediateResponseBody = "Related API definition is not configured properly"
		result.ImmediateResponseCode = 500
		result.ImmediateResponseDetail = "GraphQL SDL is empty"
	}

	// Parse the schema into a graphql object
	schema, err := gqlparser.LoadSchema(&ast.Source{Input: sdl})
	if err != nil {
		g.logger.Sugar().Errorf("error while parsing the GraphQL SDL: %v", err)
		result.ImmediateResponse = true
		result.ImmediateResponseBody = "Related API definition is not configured properly"
		result.ImmediateResponseCode = 500
		result.ImmediateResponseDetail = "error while parsing the GraphQL SDL"
		return result
	}

	// Decode the json into a graphql req
	var gqlReq GQLRequest
	requestBody := string(requestConfig.RequestBody.Body)
	if err := json.Unmarshal([]byte(requestBody), &gqlReq); err != nil {
		g.logger.Sugar().Errorf("failed to parse GraphQL request: %v", err)
	}

	cleanedQuery := strings.TrimSpace(gqlReq.Query)

	// parse the query
	document, err := parser.ParseQuery(&ast.Source{
		Input: cleanedQuery,
	})

	if err != nil {
		g.logger.Sugar().Errorf("invalid query: %v", err)
		result.ImmediateResponse = true
		result.ImmediateResponseBody = "invalid query: " + err.Error()
		result.ImmediateResponseCode = 400
		result.ImmediateResponseDetail = "invalid query: " + err.Error()
		return result
	}

	// validate query against graphql sdl
	validationErrors := validator.Validate(schema, document)
	if len(validationErrors) > 0 {
		g.logger.Sugar().Error("Validation Errors:")
		validationErrorsString := ""
		for _, err := range validationErrors {
			g.logger.Sugar().Error("-", err.Message)
			validationErrorsString += err.Message + "; "
		}
		result.ImmediateResponse = true
		result.ImmediateResponseBody = "validation error: " + validationErrorsString
		result.ImmediateResponseCode = 400
		result.ImmediateResponseDetail = "validation error: " + validationErrorsString
		return result
	}

	for _, operation := range document.Operations {
		for _, selection := range operation.SelectionSet {
			res := findMatchedOperation(operations, operation, selection)
			if res == nil {
				g.logger.Sugar().Errorf("no matching operation found for selection: %+v", selection)
				result.ImmediateResponse = true
				result.ImmediateResponseBody = "bad request - resource not found in schema"
				result.ImmediateResponseCode = 404
				result.ImmediateResponseDetail = "bad request - resource not found in schema"
				return result
			}
			remoteScopes := requestConfig.JWTAuthnPayloaClaims[constants.ScopesHeaderKey]
			scopes := make([]string, 0)
			if remoteScopes != nil {
				switch remoteScopes := remoteScopes.(type) {
				case string:
					scopes = strings.Split(remoteScopes, " ")
				case []string:
					scopes = remoteScopes
				}
			}
			for _, requiredScope := range res.Scopes {
				scopeMatched := false
				for _, scopeFromJWT := range scopes {
					if strings.EqualFold(requiredScope, scopeFromJWT) {
						g.logger.Sugar().Debugf("Matched scope: %s for operation: %s", requiredScope, res.Target)
						scopeMatched = true
						break
					}
				}
				scopeValidationErrorMessage := dto.ErrorResponse{Code: 900910, ErrorMessage: "The access token does not allow you to access the requested resource", ErrorDescription: "User is NOT authorized to access the Resource: " + operation.Name + ". Scope validation failed."}
				forbiddenJSONMessage, _ := json.MarshalIndent(scopeValidationErrorMessage, "", "  ")

				if !scopeMatched {
					g.logger.Sugar().Errorf("scope %s not found in JWT for operation: %s", requiredScope, res.Target)
					result.ImmediateResponse = true
					result.ImmediateResponseBody = string(forbiddenJSONMessage)
					result.ImmediateResponseCode = 403
					result.ImmediateResponseDetail = string(forbiddenJSONMessage)
					return result
				}
				g.logger.Sugar().Debugf("Scope %s matched for operation: %s", requiredScope, res.Target)
			}
		}
	}

	return result
}

func unzipGzip(compressedData []byte) (string, error) {
	reader, err := gzip.NewReader(bytes.NewBuffer(compressedData))
	if err != nil {
		return "", fmt.Errorf("error creating gzip reader: %v", err)
	}
	defer reader.Close()

	// Read the decompressed data
	schemaString, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("error reading decompressed data of the apiDefinition: %v", err)
	}
	return string(schemaString), nil
}

// GQLRequest is to unmarshal the request body into
type GQLRequest struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName"`
	Variables     map[string]interface{} `json:"variables"`
}

func findMatchedOperation(operations []*Operation, operation *ast.OperationDefinition, selection ast.Selection) *Operation {
	if field, ok := selection.(*ast.Field); ok {
		operation.Name = field.Name
	}

	for _, op := range operations {
		if strings.EqualFold(op.Target, string(operation.Name)) && strings.EqualFold(string(op.Verb), string(operation.Operation)) {
			return op
		}
	}
	return nil
}

func parseOperations(data []byte) ([]*Operation, error) {
	var ops []*Operation

	trimmed := strings.TrimSpace(string(data))

	if strings.HasPrefix(trimmed, "{") || strings.HasPrefix(trimmed, "[") {
		// Try JSON
		if err := json.Unmarshal(data, &ops); err != nil {
			return nil, fmt.Errorf("error parsing JSON: %w", err)
		}
	} else {
		// Try YAML
		if err := yaml.Unmarshal(data, &ops); err != nil {
			return nil, fmt.Errorf("error parsing YAML: %w", err)
		}
	}

	return ops, nil
}
