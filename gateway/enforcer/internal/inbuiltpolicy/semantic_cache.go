/*
 *  Copyright (c) 2025, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */

package inbuiltpolicy

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_service_proc_v3 "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"
	v32 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/google/uuid"
	"github.com/wso2/apk/gateway/enforcer/internal/dto"
	"github.com/wso2/apk/gateway/enforcer/internal/logging"
	"github.com/wso2/apk/gateway/enforcer/internal/semanticcache"
)

var (
	globalEmbeddingProvider   semanticcache.EmbeddingProvider
	initEmbeddingProviderOnce sync.Once

	globalVectorStoreProvider   semanticcache.VectorDBProvider
	initVectorStoreProviderOnce sync.Once
)

// SemanticCachePolicy is a struct that represents a semantic cache policy.
type SemanticCachePolicy struct {
	dto.BaseInBuiltPolicy
	embeddingConfig     semanticcache.EmbeddingProviderConfig
	vectorStoreConfig   semanticcache.VectorDBProviderConfig
	embeddingProvider   semanticcache.EmbeddingProvider
	vectorStoreProvider semanticcache.VectorDBProvider
}

// HandleRequestBody is a method that implements the mediation logic for the Semantic Caching policy on request.
func (s *SemanticCachePolicy) HandleRequestBody(logger *logging.Logger, req *envoy_service_proc_v3.ProcessingRequest, resp *envoy_service_proc_v3.ProcessingResponse, props map[string]interface{}) *envoy_service_proc_v3.ProcessingResponse {
	logger.Sugar().Debugf("Beginning request body validation for Semantic Caching policy: %s", s.PolicyID)
	ctx := props["ctx"].(context.Context)
	if len(req.GetRequestBody().Body) == 0 {
		logger.Sugar().Debug("Request body is empty, skipping semantic caching")
		return nil
	}

	embedding, err := s.embeddingProvider.GetEmbedding(logger, string(req.GetRequestBody().Body))
	if err != nil {
		logger.Error(err, "Error in embedding generation.")
		return nil
	}
	logger.Sugar().Debugf("Request Body: %s", string(req.GetRequestBody().Body))
	logger.Sugar().Debugf("Request Body Embedding Length: %d", len(embedding))
	logger.Sugar().Debugf("Request Body Embedding: %f", embedding[:4])

	cacheRetrieveConfig := map[string]interface{}{
		"threshold": s.vectorStoreConfig.Threshold,
		"api_id":    props["matchedAPIUUID"].(string),
		"ctx":       ctx,
	}
	logger.Sugar().Debug("Checking for a cached response in Vector Store")
	cacheResponse, err := s.vectorStoreProvider.Retrieve(logger, embedding, cacheRetrieveConfig)
	if err != nil {
		logger.Error(err, "Error in retrieving cached response from VectorDB.")
	}
	if len(cacheResponse.ResponsePayload) == 0 {
		logger.Sugar().Debug("Cache Miss. Sending Request to the LLM Backend.")
		embeddingBytes, err := json.Marshal(embedding)
		if err != nil {
			logger.Error(err, "failed to marshal embedding")
		} else {
			dynamicMetadataKeyValuePairs, ok := props["dynamicMetadataMap"].(map[string]interface{})
			if ok {
				dynamicMetadataKeyValuePairs[semanticCacheEmbeddingKey] = string(embeddingBytes)
			}
			logger.Sugar().Debugf("Embedding stored in metadata: %s", dynamicMetadataKeyValuePairs[semanticCacheEmbeddingKey])
		}
	} else {
		logger.Sugar().Debug("Semantic Cache Hit")
		logger.Sugar().Debugf("Cached Response: : %+v", cacheResponse.ResponsePayload)
		responseBodyBytes, err := json.Marshal(cacheResponse.ResponsePayload)
		if err != nil {
			logger.Error(err, "failed to marshal cached response payload")
			return nil
		}
		resp := s.buildImmediateResponse(logger, responseBodyBytes)
		return resp
	}

	return nil
}

// HandleResponseBody is a method that implements the mediation logic for the Semantic Caching policy on response.
func (s *SemanticCachePolicy) HandleResponseBody(logger *logging.Logger, req *envoy_service_proc_v3.ProcessingRequest, resp *envoy_service_proc_v3.ProcessingResponse, props map[string]interface{}) *envoy_service_proc_v3.ProcessingResponse {
	logger.Sugar().Debugf("Beginning response body validation for Semantic Caching policy: %s", s.PolicyID)
	ctx := props["ctx"].(context.Context)

	if props["responseHeaders"] == "200" {
		logger.Sugar().Debug("Semantic Cache RespBody logic gets hit.")
		httpBody := req.GetResponseBody().Body
		// Unmarshal the JSON data into the map
		var responseData map[string]interface{}
		err := json.Unmarshal(httpBody, &responseData)
		if err != nil {
			logger.Error(err, "Error unmarshaling JSON Response Body")
		}

		serializedEmbedding := props["embedding"].(string)
		if serializedEmbedding == "" {
			logger.Sugar().Debug("No semantic embedding found in metadata. Skipping cache storage.")
		} else {
			logger.Sugar().Debug("Found embedding in metadata. Storing new response in cache.")

			// Deserialize the embedding string back into a vector.
			var embedding []float32
			err := json.Unmarshal([]byte(serializedEmbedding), &embedding)
			if err != nil {
				logger.Error(err, "failed to unmarshal embedding from metadata")
				return nil
			}
			// Store the embedding and response body in the Vector DB.
			cr := semanticcache.CacheResponse{
				ResponsePayload:     responseData,
				RequestHash:         uuid.New().String(),
				ResponseFetchedTime: time.Now(),
			}
			err = s.vectorStoreProvider.Store(logger, embedding, cr, map[string]interface{}{
				"api_id": props["matchedAPIUUID"].(string),
				"ctx":    ctx,
			})
			if err != nil {
				logger.Error(err, "Failed to store response in vector DB")
				return nil
			}
			logger.Sugar().Debug("Response stored in the cahce store.")
		}
	}

	return nil
}

// NewSemanticCachingPolicy initializes the NewSemanticCachingPolicy policy from the given InBuiltPolicy.
func NewSemanticCachingPolicy(logger *logging.Logger, inBuiltPolicy dto.InBuiltPolicy) *SemanticCachePolicy {
	logger.Sugar().Debugf("Initializing Semantic Caching policy: %s", inBuiltPolicy.GetPolicyID())
	semanticCachePolicy := &SemanticCachePolicy{
		BaseInBuiltPolicy: dto.BaseInBuiltPolicy{
			PolicyName:    inBuiltPolicy.GetPolicyName(),
			PolicyID:      inBuiltPolicy.GetPolicyID(),
			PolicyVersion: inBuiltPolicy.GetPolicyVersion(),
			Parameters:    inBuiltPolicy.GetParameters(),
			PolicyOrder:   inBuiltPolicy.GetPolicyOrder(),
		},
	}

	for key, value := range inBuiltPolicy.GetParameters() {
		switch key {
		case "header_name":
			semanticCachePolicy.embeddingConfig.AuthHeaderName = value
		case "api_key":
			semanticCachePolicy.embeddingConfig.APIKey = value
		case "embedding_endpoint":
			semanticCachePolicy.embeddingConfig.EmbeddingEndpoint = value
		case "embedding_provider":
			semanticCachePolicy.embeddingConfig.EmbeddingProvider = value
		case "embedding_model":
			semanticCachePolicy.embeddingConfig.EmbeddingModel = value
		case "vector_store_provider":
			semanticCachePolicy.vectorStoreConfig.VectorStoreProvider = value
		case "embedding_dimention":
			semanticCachePolicy.vectorStoreConfig.EmbeddingDimention = value
		case "threshold":
			semanticCachePolicy.vectorStoreConfig.Threshold = value
		case "db_host":
			semanticCachePolicy.vectorStoreConfig.DBHost = value
		case "db_port":
			port, err := strconv.Atoi(value)
			if err == nil {
				semanticCachePolicy.vectorStoreConfig.DBPort = port
			}
		case "username":
			semanticCachePolicy.vectorStoreConfig.Username = value
		case "password":
			semanticCachePolicy.vectorStoreConfig.Password = value
		case "database":
			semanticCachePolicy.vectorStoreConfig.DatabaseName = value
		}
	}

	initEmbeddingProviderOnce.Do(func() {
		logger.Sugar().Debug("First semantic cache policy loaded. Initializing shared Embedding Provider...")
		var err error
		globalEmbeddingProvider, err = initializeEmbeddingProvider(logger, semanticCachePolicy.embeddingConfig)
		if err != nil {
			logger.Error(err, "Failed to initialize the shared Embedding Provider")
		}
	})

	initVectorStoreProviderOnce.Do(func() {
		logger.Sugar().Debug("First semantic cache policy loaded. Initializing shared Vector DB Provider...")
		var err error
		globalVectorStoreProvider, err = initializeVectorDBProvider(logger, semanticCachePolicy.vectorStoreConfig)
		if err != nil {
			logger.Error(err, "Failed to initialize the shared Vector DB Provider")
		}
	})

	if globalEmbeddingProvider == nil || globalVectorStoreProvider == nil {
		return nil
	}
	semanticCachePolicy.embeddingProvider = globalEmbeddingProvider
	semanticCachePolicy.vectorStoreProvider = globalVectorStoreProvider
	return semanticCachePolicy
}

// buildResponse is a method that builds the response body for the WordCountGuardrail policy.
func (s *SemanticCachePolicy) buildImmediateResponse(logger *logging.Logger, cachedResponseBytes []byte) *envoy_service_proc_v3.ProcessingResponse {
	logger.Sugar().Debugf("Building immediate response for Semantic Caching policy: %s", s.PolicyID)
	logger.Sugar().Debugf("Cached response content: %s", cachedResponseBytes)
	headers := &envoy_service_proc_v3.HeaderMutation{
		SetHeaders: []*corev3.HeaderValueOption{
			{
				Header: &corev3.HeaderValue{
					Key:      "Content-Type",
					RawValue: []byte("application/json"),
				},
			},
			{
				Header: &corev3.HeaderValue{
					Key:      "X-Cache-Status",
					RawValue: []byte("HIT"),
				},
			},
		},
	}

	resp := &envoy_service_proc_v3.ProcessingResponse{
		Response: &envoy_service_proc_v3.ProcessingResponse_ImmediateResponse{
			ImmediateResponse: &envoy_service_proc_v3.ImmediateResponse{
				Status: &v32.HttpStatus{
					Code: v32.StatusCode_OK,
				},
				Headers: headers,
				Body:    cachedResponseBytes,
			},
		},
	}
	return resp
}

// initializeEmbeddingProvider initializes the embedding provider based on the passed configuration.
func initializeEmbeddingProvider(logger *logging.Logger, embeddingProviderConfig semanticcache.EmbeddingProviderConfig) (semanticcache.EmbeddingProvider, error) {
	logger.Sugar().Debugf("Initializing embedding provider with config: %+v", embeddingProviderConfig)
	var embeddingProvider semanticcache.EmbeddingProvider
	switch embeddingProviderConfig.EmbeddingProvider {
	case "MISTRAL":
		logger.Info("Initializing Mistral Embedding Provider...")
		embeddingProvider = &semanticcache.MistralEmbeddingProvider{}
	case "OPENAI":
		logger.Info("Initializing Redis Embedding Provider...")
		embeddingProvider = &semanticcache.OpenAIEmbeddingProvider{}
	case "AZURE_OPENAI":
		logger.Info("Initializing Azure OpenAI Embedding Provider...")
		embeddingProvider = &semanticcache.AzureOpenAIEmbeddingProvider{}
	default:
		return nil, fmt.Errorf("unsupported embedding provider: %s", embeddingProviderConfig.EmbeddingProvider)
	}

	err := embeddingProvider.Init(logger, embeddingProviderConfig)
	if err != nil {
		logger.Sugar().Errorf("Failed to initialize embedding provider: %s", err.Error())
		return nil, err
	}
	logger.Sugar().Infof("Successfully initialized %s embedding provider", embeddingProvider.GetType())
	return embeddingProvider, nil
}

// initializeVectorDBProvider initializes the vector database provider based on the passed configuration.
func initializeVectorDBProvider(logger *logging.Logger, vectorStoreProviderConfig semanticcache.VectorDBProviderConfig) (semanticcache.VectorDBProvider, error) {
	logger.Sugar().Debugf("Initializing vectorDB store with config: %+v", vectorStoreProviderConfig)
	var vectorStoreProvider semanticcache.VectorDBProvider
	switch vectorStoreProviderConfig.VectorStoreProvider {
	case "REDIS":
		logger.Info("Initializing Redis Vector DB Provider...")
		vectorStoreProvider = &semanticcache.RedisVectorDBProvider{}
	case "MILVUS":
		logger.Info("Initializing Milvus Vector DB Provider...")
		vectorStoreProvider = &semanticcache.MilvusVectorDBProvider{}
	default:
		return nil, fmt.Errorf("unsupported vector store provider: %s", vectorStoreProviderConfig.VectorStoreProvider)
	}
	err := vectorStoreProvider.Init(logger, vectorStoreProviderConfig)
	if err != nil {
		logger.Sugar().Errorf("Failed to initialize %s vector DB provider: %s", vectorStoreProvider.GetType(), err)
		return nil, err
	}
	logger.Sugar().Infof("Successfully initialized %s vector DB provider", vectorStoreProvider.GetType())
	// Creating the index in the vector store
	indexCreationErr := vectorStoreProvider.CreateIndex(logger)
	if indexCreationErr != nil {
		logger.Error(indexCreationErr, "Failed to create index in the vector DB")
		return nil, indexCreationErr
	}
	logger.Sugar().Infof("Successfully created the index")
	return vectorStoreProvider, nil
}
