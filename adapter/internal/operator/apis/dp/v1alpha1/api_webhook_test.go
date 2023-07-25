/*
 *  Copyright (c) 2022, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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

package v1alpha1

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateAPIContext(t *testing.T) {
	type getAPITestItem struct {
		pass    bool
		message string
		context string
	}
	dataItems := []getAPITestItem{
		{
			context: "/base/v1",
			pass:    true,
			message: "",
		},
		{
			context: "/ERROR-Hello%20W/v1",
			pass:    false,
			message: "invalid API context. Does not start with / or includes invalid characters.",
		},
		{
			context: "base/v1",
			pass:    false,
			message: "invalid API context. Does not start with / or includes invalid characters.",
		},
		{
			context: "/" + strings.Repeat("e", 228) + "/v1",
			pass:    true,
			message: "",
		},
		{
			context: "/" + strings.Repeat("e", 229) + "/v1",
			pass:    false,
			message: "API context character length should not exceed 232.",
		},
		{
			context: "/base",
			pass:    false,
			message: "API context value should contain the /{APIVersion} at end.",
		},
	}
	for _, item := range dataItems {
		err := validateAPIContextFormat(item.context, "v1")
		assert.Equal(t, item.pass, err == "", item.message)
	}
}

func TestAPIDisplayNameFormat(t *testing.T) {
	type getAPITestItem struct {
		pass    bool
		message string
		context string
	}
	dataItems := []getAPITestItem{
		{
			context: "My API 1",
			pass:    true,
			message: "",
		},
		{
			context: "My API $1",
			pass:    false,
			message: "invalid API display name. Includes invalid characters.",
		},
		{
			context: strings.Repeat("e", 60),
			pass:    true,
			message: "",
		},
		{
			context: strings.Repeat("e", 61),
			pass:    false,
			message: "API display name character length should not exceed 60.",
		},
	}
	for _, item := range dataItems {
		err := validateAPIDisplayNameFormat(item.context)
		assert.Equal(t, item.pass, err == "", item.message)
	}
}

func TestAPIVersionFormat(t *testing.T) {
	type getAPITestItem struct {
		pass    bool
		message string
		context string
	}
	dataItems := []getAPITestItem{
		{
			context: "v1",
			pass:    true,
			message: "",
		},
		{
			context: "version 1",
			pass:    false,
			message: "invalid API version. Includes invalid characters spaces.",
		},
		{
			context: "v1&v2",
			pass:    false,
			message: "invalid API version. Includes invalid characters &.",
		},
		{
			context: strings.Repeat("v", 30),
			pass:    true,
			message: "",
		},
		{
			context: strings.Repeat("v", 31),
			pass:    false,
			message: "API version length should not exceed 30.",
		},
	}
	for _, item := range dataItems {
		err := validateAPIVersionFormat(item.context)
		assert.Equal(t, item.pass, err == "", item.message)
	}
}
