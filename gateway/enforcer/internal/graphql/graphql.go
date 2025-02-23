package graphql

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
	"github.com/wso2/apk/gateway/enforcer/internal/authentication/authenticator"
	"github.com/wso2/apk/gateway/enforcer/internal/authorization"
	"github.com/wso2/apk/gateway/enforcer/internal/config"
	"github.com/wso2/apk/gateway/enforcer/internal/datastore"
	"github.com/wso2/apk/gateway/enforcer/internal/dto"
	"github.com/wso2/apk/gateway/enforcer/internal/requestconfig"
	"github.com/wso2/apk/gateway/enforcer/internal/transformer"
)

// GQLRequest is to unmarshal the request body into
type GQLRequest struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName"`
	Variables     map[string]interface{} `json:"variables"`
}

// ValidateGraphQLOperation validates/authenticates the incoming GraphQL request.
func ValidateGraphQLOperation(authenticator *authenticator.Authenticator, requestConfigHolder *requestconfig.Holder, metadata *dto.ExternalProcessingEnvoyMetadata, subAppDataStore *datastore.SubscriptionApplicationDataStore, cfg *config.Server, requestBody string, jwtTransformer *transformer.JWTTransformer, revokedJTIStore *datastore.RevokedJTIStore) *dto.ImmediateResponse {
	matchedapi := requestConfigHolder.MatchedAPI
	schemaBytes := matchedapi.APIDefinition
	var sdl string
	if schemaString, err := unzipGzip(schemaBytes); err != nil {
		fmt.Println("unzip gzip not working")
	} else {
		sdl = schemaString
	}

	if sdl == "" {
		return &dto.ImmediateResponse{
			StatusCode: 500,
			Message:    "error while obtaining the GraphQL SDL",
		}
	}

	// Parse the schema into a graphql object
	schema, err := gqlparser.LoadSchema(&ast.Source{Input: sdl})
	if err != nil {
		fmt.Printf("error while parsing the GraphQL SDL: %v", err)
		return &dto.ImmediateResponse{
			StatusCode: 500,
			Message:    "error while parsing the GraphQL SDL",
		}
	}

	// Decode the json into a graphql req
	var gqlReq GQLRequest
	if err := json.Unmarshal([]byte(requestBody), &gqlReq); err != nil {
		fmt.Printf("failed to parse GraphQL request: %v", err)
	}

	cleanedQuery := strings.TrimSpace(gqlReq.Query)

	// parse the query
	document, errs := parser.ParseQuery(&ast.Source{
		Input: cleanedQuery,
	})

	if errs != nil {
		fmt.Printf("invalid query: %v", errs)
		return &dto.ImmediateResponse{
			StatusCode: 400,
			Message:    "invalid request - error in graphql query",
		}
	}

	// validate query against graphql sdl
	validationErrors := validator.Validate(schema, document)
	if len(validationErrors) > 0 {
		fmt.Println("Validation Errors:")
		for _, err := range validationErrors {
			fmt.Println("-", err.Message)
		}
		return &dto.ImmediateResponse{
			StatusCode: 400,
			Message:    "validation error: query does not fit schema",
		}
	}

	for _, operation := range document.Operations {
		for _, selection := range operation.SelectionSet {
			res := findMatchedResource(matchedapi.Resources, operation, selection)
			if res == nil {
				return &dto.ImmediateResponse{
					StatusCode: 404,
					Message:    "bad request - resource not found in schema",
				}
			}
			requestConfigHolder.MatchedResource = res
			if res.AuthenticationConfig != nil && !res.AuthenticationConfig.Disabled && !matchedapi.DisableAuthentication {
				immediateResponse := authorization.Validate(authenticator, requestConfigHolder, subAppDataStore, cfg)
				if immediateResponse != nil {
					return immediateResponse
				}
			} else {
				cfg.Logger.Sugar().Debug(fmt.Sprintf("Skipping authentication for the resource: %s", requestConfigHolder.MatchedResource.Path))
			}
		}
	}

	return nil
}

func findMatchedResource(resources []*requestconfig.Resource, operation *ast.OperationDefinition, selection ast.Selection) *requestconfig.Resource {
	if field, ok := selection.(*ast.Field); ok {
		operation.Name = field.Name
	}

	for _, res := range resources {
		if strings.EqualFold(res.Path, string(operation.Name)) && strings.EqualFold(string(res.Method), string(operation.Operation)) {
			return res
		}
	}
	return nil
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
