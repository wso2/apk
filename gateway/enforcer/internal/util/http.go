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

package util

import (
	"bytes"
	"crypto/tls"
	"net/http"
	"time"
)

// MakeGETRequest HTTP client for making GET requests with custom TLS config
func MakeGETRequest(url string, tlsConfig *tls.Config, headers map[string]string) (*http.Response, error) {
	// Create a custom HTTP client with the provided TLS configuration
	tr := &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	client := &http.Client{Transport: tr}

	// Create the HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Set request headers
	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	// Execute the request
	return client.Do(req)
}

// MakeHTTPRequestWithRetry makes an HTTP request with the given method, headers, timeout, and retry logic.
func MakeHTTPRequestWithRetry(method, url string, tlsConfig *tls.Config, headers map[string]string, body []byte, timeoutMs int, retryCount int, retryIntervalMs int) (*http.Response, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Duration(timeoutMs) * time.Millisecond,
	}

	var lastErr error
	for i := 0; i < retryCount; i++ {
		req, err := http.NewRequest(method, url, bytes.NewReader(body))
		if err != nil {
			return nil, err
		}
		for key, value := range headers {
			req.Header.Set(key, value)
		}
		resp, err := client.Do(req)
		if err == nil {
			return resp, nil
		}
		lastErr = err
		if i < retryCount-1 {
			time.Sleep(time.Duration(retryIntervalMs) * time.Millisecond)
		}
	}
	return nil, lastErr
}
