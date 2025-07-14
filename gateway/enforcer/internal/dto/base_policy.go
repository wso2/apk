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

import (
	envoy_service_proc_v3 "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"
	"github.com/wso2/apk/gateway/enforcer/internal/logging"
)

// InBuiltPolicy is an interface that defines the methods required for a policy in the enforcer.
type InBuiltPolicy interface {
	GetPolicyName() string
	GetPolicyID() string
	GetPolicyVersion() string
	GetParameters() map[string]string
	GetPolicyOrder() int
	HandleRequestHeaders(cfg *logging.Logger, req *envoy_service_proc_v3.ProcessingRequest, resp *envoy_service_proc_v3.ProcessingResponse, props map[string]interface{}) *envoy_service_proc_v3.ProcessingResponse
	HandleRequestBody(cfg *logging.Logger, req *envoy_service_proc_v3.ProcessingRequest, resp *envoy_service_proc_v3.ProcessingResponse, props map[string]interface{}) *envoy_service_proc_v3.ProcessingResponse
	HandleResponseHeaders(cfg *logging.Logger, req *envoy_service_proc_v3.ProcessingRequest, resp *envoy_service_proc_v3.ProcessingResponse, props map[string]interface{}) *envoy_service_proc_v3.ProcessingResponse
	HandleResponseBody(cfg *logging.Logger, req *envoy_service_proc_v3.ProcessingRequest, resp *envoy_service_proc_v3.ProcessingResponse, props map[string]interface{}) *envoy_service_proc_v3.ProcessingResponse
}

// BaseInBuiltPolicy is a struct that implements the Policy interface.
// It serves as a base class for all policies, providing default implementations for the methods.
type BaseInBuiltPolicy struct {
	PolicyName    string
	PolicyID      string
	PolicyVersion string
	Parameters    map[string]string
	PolicyOrder   int
}

// GetPolicyName returns the name of the policy.
func (p *BaseInBuiltPolicy) GetPolicyName() string {
	return p.PolicyName
}

// GetPolicyID returns the ID of the policy.
func (p *BaseInBuiltPolicy) GetPolicyID() string {
	return p.PolicyID
}

// GetPolicyVersion returns the version of the policy.
func (p *BaseInBuiltPolicy) GetPolicyVersion() string {
	return p.PolicyVersion
}

// GetParameters returns the parameters of the policy.
func (p *BaseInBuiltPolicy) GetParameters() map[string]string {
	if p.Parameters == nil {
		p.Parameters = make(map[string]string)
	}
	return p.Parameters
}

// GetPolicyOrder returns the order of the policy.
func (p *BaseInBuiltPolicy) GetPolicyOrder() int {
	return p.PolicyOrder
}

// HandleRequestHeaders is a method that implements the mediation logic for the policy on request headers.
func (p *BaseInBuiltPolicy) HandleRequestHeaders(cfg *logging.Logger, req *envoy_service_proc_v3.ProcessingRequest, resp *envoy_service_proc_v3.ProcessingResponse, props map[string]interface{}) *envoy_service_proc_v3.ProcessingResponse {
	cfg.Sugar().Debugf("BaseInBuiltPolicy HandleRequestHeaders called for policy: %s", p.PolicyName)
	return nil // Default implementation does nothing
}

// HandleRequestBody is a method that implements the mediation logic for the policy on request body.
func (p *BaseInBuiltPolicy) HandleRequestBody(cfg *logging.Logger, req *envoy_service_proc_v3.ProcessingRequest, resp *envoy_service_proc_v3.ProcessingResponse, props map[string]interface{}) *envoy_service_proc_v3.ProcessingResponse {
	cfg.Sugar().Debugf("BaseInBuiltPolicy HandleRequestBody called for policy: %s", p.PolicyName)
	return nil // Default implementation does nothing
}

// HandleResponseHeaders is a method that implements the mediation logic for the policy on response headers.
func (p *BaseInBuiltPolicy) HandleResponseHeaders(cfg *logging.Logger, req *envoy_service_proc_v3.ProcessingRequest, resp *envoy_service_proc_v3.ProcessingResponse, props map[string]interface{}) *envoy_service_proc_v3.ProcessingResponse {
	cfg.Sugar().Debugf("BaseInBuiltPolicy HandleResponseHeaders called for policy: %s", p.PolicyName)
	return nil // Default implementation does nothing
}

// HandleResponseBody is a method that implements the mediation logic for the policy on response body.
func (p *BaseInBuiltPolicy) HandleResponseBody(cfg *logging.Logger, req *envoy_service_proc_v3.ProcessingRequest, resp *envoy_service_proc_v3.ProcessingResponse, props map[string]interface{}) *envoy_service_proc_v3.ProcessingResponse {
	cfg.Sugar().Debugf("BaseInBuiltPolicy HandleResponseBody called for policy: %s", p.PolicyName)
	return nil // Default implementation does nothing
}
