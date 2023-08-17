/*
 *  Copyright (c) 2023, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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

package http

import (
	"fmt"
)

// AddBearerTokenToHeader adds a bearer token to the request.
func AddBearerTokenToHeader(token string, headers map[string]string) map[string]string {
	return AddCustomBearerTokenHeader("Authorization", token, headers)
}

// AddCustomBearerTokenHeader adds a bearer token to the request with specified auth header name.
func AddCustomBearerTokenHeader(headerName string, token string, headers map[string]string) map[string]string {
	if headers == nil {
		headers = make(map[string]string)
	}
	headers[headerName] = fmt.Sprintf("Bearer %s", token)
	return headers
}
