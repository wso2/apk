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
	"encoding/json"
	"errors"
	"strconv"

	"github.com/wso2/apk/gateway/enforcer/internal/logging"
)

// AssessmentResult holds the result of payload validation for assessment reporting
type AssessmentResult struct {
	InspectedContent   string
	CategoriesAnalysis []map[string]interface{}
	CategoryMap        map[string]int
	Error              string
	IsResponse         bool
	GuardrailOutput    interface{} // For storing AWS Bedrock Guardrail output or other specific outputs
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
