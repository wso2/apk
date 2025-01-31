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

// Latencies represents latency attributes in an analytics event.
type Latencies struct {
	ResponseLatency          int64 `json:"responseLatency"`
	BackendLatency           int64 `json:"backendLatency"`
	RequestMediationLatency  int64 `json:"requestMediationLatency"`
	ResponseMediationLatency int64 `json:"responseMediationLatency"`
}

// GetResponseLatency returns the response latency.
func (l *Latencies) GetResponseLatency() int64 {
	return l.ResponseLatency
}

// SetResponseLatency sets the response latency.
func (l *Latencies) SetResponseLatency(responseLatency int64) {
	l.ResponseLatency = responseLatency
}

// GetBackendLatency returns the backend latency.
func (l *Latencies) GetBackendLatency() int64 {
	return l.BackendLatency
}

// SetBackendLatency sets the backend latency.
func (l *Latencies) SetBackendLatency(backendLatency int64) {
	l.BackendLatency = backendLatency
}

// GetRequestMediationLatency returns the request mediation latency.
func (l *Latencies) GetRequestMediationLatency() int64 {
	return l.RequestMediationLatency
}

// SetRequestMediationLatency sets the request mediation latency.
func (l *Latencies) SetRequestMediationLatency(requestMediationLatency int64) {
	l.RequestMediationLatency = requestMediationLatency
}

// GetResponseMediationLatency returns the response mediation latency.
func (l *Latencies) GetResponseMediationLatency() int64 {
	return l.ResponseMediationLatency
}

// SetResponseMediationLatency sets the response mediation latency.
func (l *Latencies) SetResponseMediationLatency(responseMediationLatency int64) {
	l.ResponseMediationLatency = responseMediationLatency
}
