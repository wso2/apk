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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateBasePath(t *testing.T) {
	type getXWso2BasepathTestItem struct {
		errorNil bool
		message  string
		context  string
	}
	dataItems := []getXWso2BasepathTestItem{
		{
			context:  "/v1/base",
			errorNil: true,
			message:  "valid basepath",
		},
		{
			context:  "/ERROR-Hello%20W",
			errorNil: false,
			message:  "basepath must not include invalid characters",
		},
		{
			context:  "base",
			errorNil: false,
			message:  "basepath must start with /",
		},
		{
			context:  "",
			errorNil: false,
			message:  "basepath must not be empty",
		},
	}
	for _, item := range dataItems {
		err := validateContext(item.context)
		assert.Equal(t, item.errorNil, err == "", item.message)
	}
}
