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

package dto

import "time"

// Event represents analytics event data.
type Event struct {
	API               *ExtendedAPI                   `json:"api,omitempty" bson:"api"`
	Operation         *Operation             `json:"operation,omitempty" bson:"operation"`
	Target            *Target                `json:"target,omitempty" bson:"target"`
	Application       *Application           `json:"application,omitempty" bson:"application"`
	Latencies         *Latencies             `json:"latencies,omitempty" bson:"latencies"`
	MetaInfo          *MetaInfo              `json:"metaInfo,omitempty" bson:"meta_info"`
	Error             *Error                 `json:"error,omitempty" bson:"error"`
	ProxyResponseCode int                    `json:"proxyResponseCode,omitempty" bson:"proxy_response_code"`
	RequestTimestamp  time.Time                 `json:"requestTimestamp,omitempty" bson:"request_timestamp"`
	UserAgentHeader   string                 `json:"userAgentHeader,omitempty" bson:"user_agent_header"`
	UserName          string                 `json:"userName,omitempty" bson:"user_name"`
	UserIP            string                 `json:"userIp,omitempty" bson:"user_ip"`
	ErrorType         string                 `json:"errorType,omitempty" bson:"error_type"`
	Properties        map[string]interface{} `json:"properties,omitempty" bson:"properties"`
}
