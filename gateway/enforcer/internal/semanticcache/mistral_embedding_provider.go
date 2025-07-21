package semanticcache

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/wso2/apk/gateway/enforcer/internal/logging"
)

// MistralEmbeddingProvider implements the EmbeddingProvider interface for Mistral
type MistralEmbeddingProvider struct {
	authHeaderName string
	mistralAPIKey string
	endpointURL   string
	model         string
	client        *http.Client
}

// Init initializes the Mistral embedding provider with configuration
func (m *MistralEmbeddingProvider) Init(logger *logging.Logger, config EmbeddingProviderConfig) error {
	err := ValidateEmbeddingProviderConfigProps(config)
	if err != nil {
		return fmt.Errorf("invalid embedding provider config properties: %v", err)
	}
	m.mistralAPIKey = config.APIKey
	m.endpointURL = config.EmbeddingEndpoint
	m.model = config.EmbeddingModel
	m.authHeaderName = config.AuthHeaderName
	timeout := DefaultTimeout
	if v, err := strconv.Atoi(config.TimeOut); err == nil {
		timeout = v
	}
	m.client = &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}
	return nil
}

// GetType returns the type of the embedding provider
func (m *MistralEmbeddingProvider) GetType() string {
	return "MISTRAL"
}

// GetEmbedding generates an embedding vector for a single input text
func (m *MistralEmbeddingProvider) GetEmbedding(logger *logging.Logger, input string) ([]float32, error) {
	requestBody := map[string]string{
		"model": m.model,
		"input": input,
	}
	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", m.endpointURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	logger.Sugar().Debugf("Sending request to Mistral API with key: %s", m.mistralAPIKey)
	req.Header.Set(m.authHeaderName, "Bearer "+m.mistralAPIKey) // Header should be "Authorization"
	req.Header.Set("Content-Type", "application/json")
	resp, err := m.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	logger.Sugar().Debugf("Response from Mistral API: %+v", respBody)

	var response map[string]interface{}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, err
	}
	dataArray, ok := response["data"].([]interface{})
	if !ok || len(dataArray) == 0 {
		return nil, errors.New("no data found in embedding response")
	}
	firstItem, ok := dataArray[0].(map[string]interface{})
	if !ok {
		return nil, errors.New("invalid data format")
	}
	rawEmbedding, ok := firstItem["embedding"].([]interface{})
	if !ok {
		return nil, errors.New("embedding field missing or invalid")
	}
	embedding := make([]float32, len(rawEmbedding))
	for i, val := range rawEmbedding {
		switch v := val.(type) {
		case float64:
			embedding[i] = float32(v)
		default:
			return nil, fmt.Errorf("unexpected value type in embedding: %T", v)
		}
	}
	return embedding, nil
}

// GetEmbeddings generates embedding vectors for multiple input texts
func (m *MistralEmbeddingProvider) GetEmbeddings(logger *logging.Logger, inputs []string) ([][]float32, error) {
	requestBody := map[string]interface{}{
		"model": m.model,
		"input": inputs,
	}
	body, err := json.Marshal(requestBody)
	if err != nil {
		logger.Sugar().Debug("Error marshalling request body: %v", err)
		return nil, err
	}

	req, err := http.NewRequest("POST", m.endpointURL, bytes.NewBuffer(body))
	if err != nil {
		logger.Sugar().Debug("Error creating new HTTP request: %v", err)
		return nil, err
	}
	req.Header.Set(m.authHeaderName, "Bearer "+m.mistralAPIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		logger.Sugar().Debug("Error sending request to Mistral API: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Sugar().Debug("Error reading response body: %v", err)
		return nil, err
	}
    
    if resp.StatusCode != http.StatusOK {
        errStr := fmt.Sprintf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
        logger.Sugar().Debug(errStr)
        return nil, errors.New(errStr)
    }

	var response map[string]interface{}
	if err := json.Unmarshal(respBody, &response); err != nil {
		logger.Sugar().Debug("Error unmarshalling response JSON: %v", err)
		return nil, err
	}

	dataArray, ok := response["data"].([]interface{})
	if !ok {
		return nil, errors.New("no 'data' field found in embedding response or it's not an array")
	}
	allEmbeddings := make([][]float32, 0, len(dataArray))
	for i, item := range dataArray {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid data format for item at index %d", i)
		}

		rawEmbedding, ok := itemMap["embedding"].([]interface{})
		if !ok {
			return nil, fmt.Errorf("embedding field missing or invalid for item at index %d", i)
		}

		embedding := make([]float32, len(rawEmbedding))
		for j, val := range rawEmbedding {
			switch v := val.(type) {
			case float64:
				embedding[j] = float32(v)
			default:
				return nil, fmt.Errorf("unexpected value type in embedding at index %d: %T", j, v)
			}
		}
		allEmbeddings = append(allEmbeddings, embedding)
	}

	return allEmbeddings, nil
}
