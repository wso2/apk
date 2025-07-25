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
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/andybalholm/brotli"
	"github.com/wso2/apk/gateway/enforcer/internal/logging"
	"github.com/wso2/apk/gateway/enforcer/internal/util"
)

// AssessmentResult holds the result of payload validation for assessment reporting
type AssessmentResult struct {
	InspectedContent   string
	CategoriesAnalysis []map[string]interface{}
	CategoryMap        map[string]int
	Error              string
	IsResponse         bool
	GuardrailOutput    interface{} // For storing AWS Bedrock Guardrail output or other specific outputs
	ModifiedPayload    *[]byte     // For storing modified payload content (PII masking/redaction)
}

// ExtractStringValueFromJsonpath extracts a value from a nested JSON structure based on a JSON path.
func ExtractStringValueFromJsonpath(logger *logging.Logger, payload []byte, jsonpath string) (string, error) {
	if jsonpath == "" {
		logger.Sugar().Debugf("No JSONPath provided, returning payload as string")
		return string(payload), nil
	}
	var jsonData map[string]interface{}
	if err := json.Unmarshal(payload, &jsonData); err != nil {
		logger.Error(err, "Error unmarshaling JSON Request Body")
		return "", err
	}
	value, err := extractValueFromJsonpath(jsonData, jsonpath)
	if err != nil {
		logger.Error(err, "Error extracting value from JSON using JSONPath")
		return "", err
	}
	// Convert to string if possible
	switch v := value.(type) {
	case string:
		return v, nil
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case int:
		return strconv.Itoa(v), nil
	default:
		return "", errors.New("value at JSONPath is not a string or number")
	}
}

// DecompressLLMResp will properly decompress the response given from the LLM
func DecompressLLMResp(body []byte) (string, string, error) {
	asString := string(body)
	if util.IsValidJSON(asString) {
		return asString, "", nil
	}

	// Try GZIP first
	gzipReader, err := gzip.NewReader(bytes.NewReader(body))
	if err == nil {
		defer gzipReader.Close()
		unzipped, err := io.ReadAll(gzipReader)
		if err == nil {
			return string(unzipped), "gzip", nil
		}
	}

	// If GZIP failed, try Brotli
	brReader := brotli.NewReader(bytes.NewReader(body))
	unbr, err := io.ReadAll(brReader)
	if err != nil {
		return "", "", fmt.Errorf("failed to decompress response body: %w", err)
	}
	return string(unbr), "br", nil
}

// CompressLLMResp compresses content with a specified compression type
// compressionType can be "gzip", "br" (brotli), or empty for no compression
func CompressLLMResp(content []byte, compressionType string) ([]byte, error) {
	// If no compression requested or content is already valid JSON, return as is
	if compressionType == "" {
		return content, nil
	}

	switch compressionType {
	case "gzip":
		// Compress with GZIP
		var compressedBuffer bytes.Buffer
		gzipWriter := gzip.NewWriter(&compressedBuffer)
		_, err := gzipWriter.Write(content)
		if err != nil {
			return nil, fmt.Errorf("failed to write content to gzip writer: %w", err)
		}
		if err := gzipWriter.Close(); err != nil {
			return nil, fmt.Errorf("failed to close gzip writer: %w", err)
		}
		return compressedBuffer.Bytes(), nil

	case "br":
		// Compress with Brotli
		var compressedBuffer bytes.Buffer
		brWriter := brotli.NewWriter(&compressedBuffer)
		_, err := brWriter.Write(content)
		if err != nil {
			return nil, fmt.Errorf("failed to write content to brotli writer: %w", err)
		}
		if err := brWriter.Close(); err != nil {
			return nil, fmt.Errorf("failed to close brotli writer: %w", err)
		}
		return compressedBuffer.Bytes(), nil

	default:
		return nil, fmt.Errorf("unsupported compression type: %s", compressionType)
	}
}
