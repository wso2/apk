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
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	ext_procv3 "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"
	"github.com/wso2/apk/gateway/enforcer/internal/requestconfig"
	"net/url"
	"sort"
	"strings"
	"time"
)

// Constants for AWS Signature Version 4
const (
	amzDateFormat            = "20060102T150405Z"
	dateFormat               = "20060102"
	aws4Algorithm            = "AWS4-HMAC-SHA256"
	aws4Request              = "aws4_request"
	aws4Credential           = "Credential"
	aws4SignedHeaders        = "SignedHeaders"
	aws4Signature            = "Signature"
	headerHost               = "host"
	headerXAmzDate           = "x-amz-date"
	headerXAmzSecurityToken  = "x-amz-security-token"
	headerXAmzContentSHA256  = "x-amz-content-sha256"
	headerContentType        = "content-type"
	applicationJSONMediaType = "application/json"
	authorizationHeader      = "Authorization"
)

// GenerateAWSSignatureHeaders generates the AWS Signature Version 4 headers for a request.
func GenerateAWSSignatureHeaders(matchedAPI *requestconfig.API, matchedResource *requestconfig.Resource,
	req *ext_procv3.ProcessingRequest) (map[string]string, error) {
	host, method, service, uri, queryString, payload, region, accessKey, secretKey,
		sessionToken, incomingHeaders := extractAWSSignatureParameters(matchedAPI, matchedResource, req)
	if accessKey != "" && secretKey != "" && host != "" {
		canonicalRequest, stringToSign, awsHeaders, err := generateAWSSignature(host, method, service, uri, queryString, payload,
			accessKey, secretKey, region, sessionToken, incomingHeaders)
		fmt.Printf("Canonical Request: \n%s\nString to Sign: \n%s\n", canonicalRequest, stringToSign)
		if err != nil {
			return nil, fmt.Errorf("missing required fields: 'accessKey', 'secretKey', or 'region'")
		}
		return awsHeaders, nil
	}
	return nil, fmt.Errorf("AWS signature generation skipped - missing required credentials or host information")
}

// extractAWSSignatureParameters extracts AWS signature parameters from request context and configuration
func extractAWSSignatureParameters(matchedAPI *requestconfig.API, matchedResource *requestconfig.Resource,
	req *ext_procv3.ProcessingRequest) (host, method, service, uri, queryString, payload,
	region, accessKey, secretKey, sessionToken string, incomingHeaders map[string]string) {
	// Extract host from endpoint URL
	if matchedResource.Endpoints != nil && len(matchedResource.Endpoints.URLs) > 0 {
		endpointURL := matchedResource.Endpoints.URLs[0]
		if strings.HasPrefix(endpointURL, "https://") {
			hostWithPort := strings.TrimPrefix(endpointURL, "https://")
			if colonIndex := strings.Index(hostWithPort, ":"); colonIndex != -1 {
				host = hostWithPort[:colonIndex]
			} else {
				host = hostWithPort
			}
		}
	}

	// Extract method
	method = string(matchedResource.Method)

	// Remove base path from URI if present
	uri = strings.TrimPrefix(matchedResource.RouteMetadataAttributes.URI, matchedAPI.BasePath)

	// Extract query string from the request path if present
	if questionIndex := strings.Index(uri, "?"); questionIndex != -1 {
		queryString = uri[questionIndex+1:]
		uri = uri[:questionIndex]
	}

	uriParts := strings.Split(uri, "/")
	for i, part := range uriParts {
		uriParts[i] = url.QueryEscape(part)
	}
	uri = strings.Join(uriParts, "/")

	// Extract payload from request headers
	if req != nil && req.GetRequestBody() != nil {
		payload = string(req.GetRequestBody().Body)
	}

	// Extract AWS credentials from endpoint security configuration
	service, region, accessKey, secretKey = extractAWSCredentials(matchedAPI, matchedResource)

	// Prepare incoming headers map
	incomingHeaders = make(map[string]string)

	return
}

// extractAWSCredentials extracts AWS credentials from endpoint security configuration
func extractAWSCredentials(matchedAPI *requestconfig.API, matchedResource *requestconfig.Resource) (service, region,
	accessKey, secretKey string) {
	// Check API level endpoint security
	if matchedAPI != nil && matchedAPI.EndpointSecurity != nil {
		for _, es := range matchedAPI.EndpointSecurity {
			if es.Enabled && es.SecurityType == "AWSKey" {
				if es.CustomParameters != nil {
					service = es.CustomParameters["service"]
					region = es.CustomParameters["region"]
					accessKey = es.CustomParameters["accessKey"]
					secretKey = es.CustomParameters["secretKey"]
					return
				}
			}
		}
	}

	// Check resource level endpoint security
	if matchedResource != nil && matchedResource.EndpointSecurity != nil {
		for _, es := range matchedResource.EndpointSecurity {
			if es.Enabled && es.SecurityType == "AWSKey" {
				if es.CustomParameters != nil {
					service = es.CustomParameters["service"]
					region = es.CustomParameters["region"]
					accessKey = es.CustomParameters["accessKey"]
					secretKey = es.CustomParameters["secretKey"]
					return
				}
			}
		}
	}

	return
}

// generateAWSSignature creates AWS Signature Version 4 headers for authenticating requests.
// It constructs the signature based on the provided request parameters and AWS credentials.
func generateAWSSignature(host, method, service, uri, queryString, payload, accessKey, secretKey, region,
	sessionToken string, incomingHeaders map[string]string) (string, string, map[string]string, error) {
	if accessKey == "" || secretKey == "" || region == "" {
		return "", "", nil, fmt.Errorf("missing required fields: 'accessKey', 'secretKey', or 'region'")
	}

	// Step 1: Create date stamps
	now := time.Now().UTC()
	amzDate := now.Format(amzDateFormat)
	dateStamp := now.Format(dateFormat)

	// Step 2: Create canonical headers
	headers := make(map[string]string)
	for k, v := range incomingHeaders {
		headers[strings.ToLower(k)] = v
	}
	headers[headerHost] = host
	headers[headerXAmzDate] = amzDate
	if sessionToken != "" {
		headers[headerXAmzSecurityToken] = sessionToken
	}

	payloadHash := getSha256Digest(payload)
	if payload != "" {
		headers[headerContentType] = applicationJSONMediaType
		headers[headerXAmzContentSHA256] = payloadHash
	}

	// Build canonical headers string and signed headers list
	var canonicalHeaders strings.Builder
	var signedHeadersBuilder strings.Builder
	var headerKeys []string
	for k := range headers {
		headerKeys = append(headerKeys, k)
	}
	sort.Strings(headerKeys)

	for _, k := range headerKeys {
		canonicalHeaders.WriteString(fmt.Sprintf("%s:%s\n", k, headers[k]))
		signedHeadersBuilder.WriteString(k + ";")
	}
	signedHeaders := strings.TrimSuffix(signedHeadersBuilder.String(), ";")

	// Step 3: Create canonical request
	canonicalQueryString := createCanonicalQueryString(queryString)
	canonicalRequest := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s",
		method, uri, canonicalQueryString, canonicalHeaders.String(), signedHeaders, payloadHash)

	// Step 4: Create string to sign
	credentialScope := fmt.Sprintf("%s/%s/%s/%s", dateStamp, region, service, aws4Request)
	stringToSign := fmt.Sprintf("%s\n%s\n%s\n%s",
		aws4Algorithm, amzDate, credentialScope, getSha256Digest(canonicalRequest))

	// Step 5: Calculate signature
	signingKey, err := getSignatureKey(secretKey, dateStamp, region, service)
	if err != nil {
		return "", "", nil, fmt.Errorf("error getting signature key: %w", err)
	}
	signature := hex.EncodeToString(hmacSHA256([]byte(stringToSign), signingKey))

	// Step 6: Create authorization header
	authorizationHeaderValue := fmt.Sprintf("%s %s=%s/%s, %s=%s, %s=%s",
		aws4Algorithm, aws4Credential, accessKey, credentialScope, aws4SignedHeaders, signedHeaders, aws4Signature, signature)

	// Create result map with all required headers
	authHeaders := make(map[string]string)
	for k, v := range headers {
		authHeaders[k] = v
	}
	authHeaders[authorizationHeader] = authorizationHeaderValue

	return canonicalRequest, stringToSign, authHeaders, nil
}

// createCanonicalQueryString sorts the query string parameters into canonical form.
func createCanonicalQueryString(queryString string) string {
	if queryString == "" {
		return ""
	}
	params := strings.Split(queryString, "&")
	sort.Strings(params)
	return strings.Join(params, "&")
}

// getSha256Digest returns the SHA-256 hash of a string.
func getSha256Digest(input string) string {
	hasher := sha256.New()
	hasher.Write([]byte(input))
	return hex.EncodeToString(hasher.Sum(nil))
}

// getSignatureKey calculates the AWS Signature Version 4 signing key.
func getSignatureKey(key, dateStamp, regionName, serviceName string) ([]byte, error) {
	kSecret := []byte("AWS4" + key)
	kDate := hmacSHA256([]byte(dateStamp), kSecret)
	kRegion := hmacSHA256([]byte(regionName), kDate)
	kService := hmacSHA256([]byte(serviceName), kRegion)
	kSigning := hmacSHA256([]byte(aws4Request), kService)
	return kSigning, nil
}

// hmacSHA256 computes the HMAC-SHA256 hash.
func hmacSHA256(data, key []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}
