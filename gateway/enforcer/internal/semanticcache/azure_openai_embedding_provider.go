package semanticcache

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/wso2/apk/gateway/enforcer/internal/logging"
)

// AzureOpenAIEmbeddingProvider implements the EmbeddingProvider interface for Azure OpenAI
type AzureOpenAIEmbeddingProvider struct {
	authHeaderName string
	azureAPIKey    string
	endpointURL    string
	client         *http.Client
}

// Init initializes the Azure OpenAI embedding provider with configuration
func (a *AzureOpenAIEmbeddingProvider) Init(logger *logging.Logger, config EmbeddingProviderConfig) error {
	err := ValidateEmbeddingProviderConfigProps(config)
	if err != nil {
		return fmt.Errorf("invalid embedding provider config properties: %v", err)
	}
	a.azureAPIKey = config.APIKey
	a.endpointURL = config.EmbeddingEndpoint
	a.authHeaderName = config.AuthHeaderName
	timeout := DefaultTimeout
	if v, err := strconv.Atoi(config.TimeOut); err == nil {
		timeout = v
	}

	a.client = &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}
	return nil
}

// GetType returns the type of the embedding provider
func (a *AzureOpenAIEmbeddingProvider) GetType() string {
	return "AZURE_OPENAI"
}

// GetEmbedding generates an embedding vector for a single input text
func (a *AzureOpenAIEmbeddingProvider) GetEmbedding(logger *logging.Logger, input string) ([]float32, error) {
	requestBody := map[string]interface{}{
		"input": input,
	}
	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", a.endpointURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set(a.authHeaderName, a.azureAPIKey) //Header should be "api-key"
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response map[string]interface{}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, err
	}

	data := response["data"].([]interface{})[0].(map[string]interface{})
	embedding := data["embedding"].([]interface{})
	embeddingResult := make([]float32, len(embedding))
	for i, value := range embedding {
		embeddingResult[i] = float32(value.(float64))
	}

	return embeddingResult, nil
}

// GetEmbeddings generates embedding vectors for multiple input texts
func (a *AzureOpenAIEmbeddingProvider) GetEmbeddings(logger *logging.Logger, inputs []string) ([][]float32, error) {
	requestBody := map[string]interface{}{
		"input": inputs,
	}
	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", a.endpointURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set(a.authHeaderName, a.azureAPIKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response map[string]interface{}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, err
	}

	data := response["data"].([]interface{})
	var embeddings [][]float32
	for _, dataNode := range data {
		dataMap := dataNode.(map[string]interface{})
		embedding := dataMap["embedding"].([]interface{})
		embeddingResult := make([]float32, len(embedding))
		for i, value := range embedding {
			embeddingResult[i] = float32(value.(float64))
		}
		embeddings = append(embeddings, embeddingResult)
	}

	return embeddings, nil
}
