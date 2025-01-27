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

// Event represents analytics event data.
type Event struct {
	API               *API                   `json:"api,omitempty" bson:"api"`
	Operation         *Operation             `json:"operation,omitempty" bson:"operation"`
	Target            *Target                `json:"target,omitempty" bson:"target"`
	Application       *Application           `json:"application,omitempty" bson:"application"`
	Latencies         *Latencies             `json:"latencies,omitempty" bson:"latencies"`
	MetaInfo          *MetaInfo              `json:"metaInfo,omitempty" bson:"meta_info"`
	Error             *Error                 `json:"error,omitempty" bson:"error"`
	ProxyResponseCode int                    `json:"proxyResponseCode,omitempty" bson:"proxy_response_code"`
	RequestTimestamp  string                 `json:"requestTimestamp,omitempty" bson:"request_timestamp"`
	UserAgentHeader   string                 `json:"userAgentHeader,omitempty" bson:"user_agent_header"`
	UserName          string                 `json:"userName,omitempty" bson:"user_name"`
	UserIP            string                 `json:"userIp,omitempty" bson:"user_ip"`
	ErrorType         string                 `json:"errorType,omitempty" bson:"error_type"`
	Properties        map[string]interface{} `json:"properties,omitempty" bson:"properties"`
}

// Getters and Setters (optional in Go)

// GetAPI returns the API.
func (e *Event) GetAPI() *API {
	return e.API
}

// SetAPI sets the API.
func (e *Event) SetAPI(api *API) {
	e.API = api
}

// GetOperation returns the Operation.
func (e *Event) GetOperation() *Operation {
	return e.Operation
}

// SetOperation sets the Operation.
func (e *Event) SetOperation(operation *Operation) {
	e.Operation = operation
}

// GetTarget returns the Target.
func (e *Event) GetTarget() *Target {
	return e.Target
}

// SetTarget sets the Target.
func (e *Event) SetTarget(target *Target) {
	e.Target = target
}

// GetApplication returns the Application.
func (e *Event) GetApplication() *Application {
	return e.Application
}

// SetApplication sets the Application.
func (e *Event) SetApplication(application *Application) {
	e.Application = application
}

// GetLatencies returns the Latencies.
func (e *Event) GetLatencies() *Latencies {
	return e.Latencies
}

// SetLatencies sets the Latencies.
func (e *Event) SetLatencies(latencies *Latencies) {
	e.Latencies = latencies
}

// GetMetaInfo returns the MetaInfo.
func (e *Event) GetMetaInfo() *MetaInfo {
	return e.MetaInfo
}

// SetMetaInfo sets the MetaInfo.
func (e *Event) SetMetaInfo(metaInfo *MetaInfo) {
	e.MetaInfo = metaInfo
}

// GetError returns the Error.
func (e *Event) GetError() *Error {
	return e.Error
}

// SetError sets the Error.
func (e *Event) SetError(errorL *Error) {
	e.Error = errorL
}

// GetProxyResponseCode returns the ProxyResponseCode.
func (e *Event) GetProxyResponseCode() int {
	return e.ProxyResponseCode
}

// SetProxyResponseCode sets the ProxyResponseCode.
func (e *Event) SetProxyResponseCode(proxyResponseCode int) {
	e.ProxyResponseCode = proxyResponseCode
}

// GetRequestTimestamp returns the RequestTimestamp.
func (e *Event) GetRequestTimestamp() string {
	return e.RequestTimestamp
}

// SetRequestTimestamp sets the RequestTimestamp.
func (e *Event) SetRequestTimestamp(requestTimestamp string) {
	e.RequestTimestamp = requestTimestamp
}

// GetUserAgentHeader returns the UserAgentHeader.
func (e *Event) GetUserAgentHeader() string {
	return e.UserAgentHeader
}

// SetUserAgentHeader sets the UserAgentHeader.
func (e *Event) SetUserAgentHeader(userAgentHeader string) {
	e.UserAgentHeader = userAgentHeader
}

// GetUserName returns the UserName.
func (e *Event) GetUserName() string {
	return e.UserName
}

// SetUserName sets the UserName.
func (e *Event) SetUserName(userName string) {
	e.UserName = userName
}

// GetUserIP returns the UserIP.
func (e *Event) GetUserIP() string {
	return e.UserIP
}

// SetUserIP sets the UserIP.
func (e *Event) SetUserIP(userIP string) {
	e.UserIP = userIP
}

// GetProperties returns the Properties.
func (e *Event) GetProperties() map[string]interface{} {
	return e.Properties
}

// SetProperties sets the Properties.
func (e *Event) SetProperties(properties map[string]interface{}) {
	e.Properties = properties
}
