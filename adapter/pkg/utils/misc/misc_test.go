/*
 *  Copyright (c) 2024, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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
 *  This file contains code derived from Envoy Gateway,
 *  https://github.com/envoyproxy/gateway
 *  and is provided here subject to the following:
 *  Copyright Project Envoy Gateway Authors
 *
 */

package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetHashedName(t *testing.T) {
	testCases := []struct {
		name     string
		nsName   string
		length   int
		expected string
	}{
		{"test default name", "http", 6, "http-e0603c49"},
		{"test removing trailing slash", "namespace/name", 10, "namespace-18a6500f"},
		{"test removing trailing hyphen", "apk/eg/http", 6, "apk-eg-9df93c35"},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			result := GetHashedName(tc.nsName, tc.length)
			require.Equal(t, tc.expected, result, "Result does not match expected string")
		})
	}
}
