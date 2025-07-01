/* Copyright (c) 2025 WSO2 LLC. (http://www.wso2.com) All Rights Reserved.
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package constants

// API Type related constants.
const (
	API_TYPE_REST     string = "REST"
	API_TYPE_GRAPHQL  string = "GRAPHQL"
	API_TYPE_GRPC     string = "GRPC"
	API_TYPE_ASYNC    string = "ASYNC"
	API_TYPE_SOAP     string = "SOAP"
	API_TYPE_SSE      string = "SSE"
	API_TYPE_WS       string = "WS"
	API_TYPE_WEBSUB   string = "WEBSUB"
)

// ALLOWED_API_TYPES is a list of allowed API types.
var ALLOWED_API_TYPES = []string{
	API_TYPE_REST, 
	API_TYPE_GRAPHQL,
	API_TYPE_GRPC,
}
