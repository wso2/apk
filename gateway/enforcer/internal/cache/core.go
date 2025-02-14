package cache

import (
	"encoding/json"
	"fmt"
	"time"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_service_proc_v3 "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"
	v32 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/wso2/apk/gateway/enforcer/internal/datastore"
	"github.com/wso2/apk/gateway/enforcer/internal/dto"
)

// CheckCacheForKey checks if the key is in the cache
func CheckCacheForKey(key string, cacheStore datastore.CacheStore) (string, error) {

	return cacheStore.Get(key)
	// TODO: check vector similarity search if redis cache miss

}

// Caches the response value
func cacheResponse(key string, value string, cacheStore datastore.CacheStore) {
	err := cacheStore.Set(key, value)
	if err != nil {
		fmt.Printf("[AI-CACHE] cache set failed, key: %s, error: %v", key, err)
		return
	}
	fmt.Printf("[AI-CACHE] cache set success, key: %s, length of value: %d", key, len(value))

}

// SendCachedHTTPResponse makes ext_proc response for cached value
func SendCachedHTTPResponse(cachedResponse string, resp *envoy_service_proc_v3.ProcessingResponse) {

	llmResponse := dto.LLMResponse{
		ID:      "chatcmpl-123", // You may want to generate a unique ID
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   "gpt-3.5-turbo",
		Usage: dto.Usage{
			PromptTokens:     0,
			CompletionTokens: 0,
			TotalTokens:      0,
		},
		Choices: []dto.Choice{
			{
				Index: 0,
				Message: dto.Message{
					Role:    "assistant",
					Content: cachedResponse,
				},
				Delta:        []any{nil},
				FinishReason: "stop",
			},
		},
	}

	httpBody, _ := json.Marshal(llmResponse)
	httpBodyLength := len(httpBody)

	headers := &envoy_service_proc_v3.HeaderMutation{
		SetHeaders: []*corev3.HeaderValueOption{
			{
				Header: &corev3.HeaderValue{
					Key:      "Content-Length",
					RawValue: []byte(fmt.Sprintf("%d", httpBodyLength)),
				},
			},
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

	rbq := &envoy_service_proc_v3.ImmediateResponse{
		Status: &v32.HttpStatus{
			Code: v32.StatusCode_OK,
		},
		Headers: headers,
		Body:    httpBody,
	}

	resp.Response = &envoy_service_proc_v3.ProcessingResponse_ImmediateResponse{
		ImmediateResponse: rbq,
	}

}
