package cache

import (
	"encoding/json"
	"fmt"

	envoy_service_proc_v3 "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"
	"github.com/wso2/apk/gateway/enforcer/internal/datastore"
	"github.com/wso2/apk/gateway/enforcer/internal/dto"
	"github.com/wso2/apk/gateway/enforcer/internal/util"
)

// HandleHTTPRequestBody handles http request body
func HandleHTTPRequestBody(requestID string, cacheStore datastore.CacheStore, keyStore *datastore.IncomingRequestCacheKeyStore, req *envoy_service_proc_v3.ProcessingRequest, resp *envoy_service_proc_v3.ProcessingResponse) {
	httpBody := req.GetRequestBody().Body

	var llmRequest dto.LLMRequest
	if err := json.Unmarshal(httpBody, &llmRequest); err != nil {
		fmt.Printf("[AI-CACHE] Error unmarshaling JSON Reuqest Body. %v", err)
		return
	}

	key, has := llmRequest.GetKey()
	if !has {
		fmt.Printf("[AI-CACHE] cache key not found in request body.")
		return
	}

	cachedResponse, err := CheckCacheForKey(key, cacheStore)
	if err != nil {
		fmt.Printf("[AI-CACHE] error retrieving key: %s from cache, error: %v", key, err)
		keyStore.Set(requestID, key) // TODO: perform only if cache miss
		return
	}

	SendCachedHTTPResponse(cachedResponse, resp)
}

// HandleHTTPResponseBody handles http response body
func HandleHTTPResponseBody(requestID string, cacheStore datastore.CacheStore, keyStore *datastore.IncomingRequestCacheKeyStore, req *envoy_service_proc_v3.ProcessingRequest, resp *envoy_service_proc_v3.ProcessingResponse) {
	httpBody := req.GetResponseBody().Body

	var llmResponse dto.LLMResponse

	uncompressedBody, err := util.DecompressIfGzip(httpBody)
	if err != nil {
		fmt.Printf("[AI-CACHE] Error decompressing response body, error: %v", err)
		return
	}

	if err := json.Unmarshal(uncompressedBody, &llmResponse); err != nil {
		fmt.Printf("[AI-CACHE] Error unmarshaling JSON Response Body, error: %v", err)
		return
	}

	key, hasKey := keyStore.Pop(requestID)
	if !hasKey {
		fmt.Printf("[AI-CACHE] cache key not found for request ID: %s", requestID)
		return
	}

	responseValue, hasValue := llmResponse.GetValue()
	if !hasValue {
		fmt.Printf("[AI-CACHE] cached value for key %s is missing or empty", key)
		return
	}

	cacheResponse(key, responseValue, cacheStore)
}
