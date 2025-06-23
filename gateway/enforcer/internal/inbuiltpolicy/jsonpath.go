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
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var arrayIndexRegex = regexp.MustCompile(`^([a-zA-Z0-9_]+)\[(-?\d+)\]$`)

// extractValueFromJsonpath extracts a value from a nested JSON structure based on a JSON path.
func extractValueFromJsonpath(data map[string]interface{}, jsonpath string) (interface{}, error) {
	keys := strings.Split(jsonpath, ".")
	if len(keys) > 0 && keys[0] == "$" {
		keys = keys[1:]
	}

	return extractRecursive(data, keys)
}

func extractRecursive(current interface{}, keys []string) (interface{}, error) {
	if len(keys) == 0 {
		return current, nil
	}
	key := keys[0]
	remaining := keys[1:]

	if key == "*" {
		var results []interface{}
		switch node := current.(type) {
		case map[string]interface{}:
			for _, v := range node {
				res, err := extractRecursive(v, remaining)
				if err == nil {
					results = append(results, res)
				}
			}
		case []interface{}:
			for _, v := range node {
				res, err := extractRecursive(v, remaining)
				if err == nil {
					results = append(results, res)
				}
			}
		default:
			return nil, errors.New("wildcard used on non-iterable node")
		}
		return results, nil
	}

	if matches := arrayIndexRegex.FindStringSubmatch(key); len(matches) == 3 {
		arrayName := matches[1]
		idxStr := matches[2]
		idx, err := strconv.Atoi(idxStr)
		if err != nil {
			return nil, errors.New("invalid array index: " + idxStr)
		}
		if node, ok := current.(map[string]interface{}); ok {
			if arrVal, exists := node[arrayName]; exists {
				if arr, ok := arrVal.([]interface{}); ok {
					if idx < 0 {
						idx = len(arr) + idx
					}
					if idx < 0 || idx >= len(arr) {
						return nil, errors.New("array index out of range: " + idxStr)
					}
					return extractRecursive(arr[idx], remaining)
				}
				return nil, errors.New("not an array: " + arrayName)
			}
			return nil, errors.New("key not found: " + arrayName)
		}
		return nil, errors.New("invalid structure for key: " + arrayName)
	}

	if node, ok := current.(map[string]interface{}); ok {
		if val, exists := node[key]; exists {
			return extractRecursive(val, remaining)
		}
		return nil, errors.New("key not found: " + key)
	}
	return nil, errors.New("invalid structure for key: " + key)
}
