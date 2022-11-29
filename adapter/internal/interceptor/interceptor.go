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

package interceptor

import (
	"bytes"
	"fmt"
	"text/template"

	logger "github.com/wso2/apk/adapter/internal/loggers"
	"github.com/wso2/apk/adapter/pkg/logging"
)

//Interceptor hold values used for interceptor
type Interceptor struct {
	Context               *InvocationContext
	IsRequestFlowEnabled  bool
	IsResponseFlowEnabled bool
	RequestFlow           map[string]Config // key:operation method -> value:config
	ResponseFlow          map[string]Config // key:operation method -> value:config
}

//Config hold config values used for request/response interceptors
type Config struct {
	Enable       bool
	ExternalCall *HTTPCallConfig
	Include      *RequestInclusions
}

//HTTPCallConfig hold values used for external interceptor engine
type HTTPCallConfig struct {
	ClusterName     string
	Timeout         string // in milli seconds
	AuthorityHeader string
}

// RequestInclusions represents which should be included in the request payload to the interceptor service
type RequestInclusions struct {
	InvocationContext bool
	RequestHeaders    bool
	RequestBody       bool
	RequestTrailer    bool
	ResponseHeaders   bool
	ResponseBody      bool
	ResponseTrailers  bool
}

// InvocationContext represents static details of the invocation context of a request for the resource path
// runtime details such as actual path will be populated from the lua script and set in the invocation context
type InvocationContext struct {
	OrganizationID   string
	BasePath         string
	SupportedMethods string
	APIName          string
	APIVersion       string
	PathTemplate     string
	Vhost            string
	ProdClusterName  string
	SandClusterName  string
}

var (
	// commonTemplate contains common lua code for request and response intercept
	// Note: this template only applies if request or response interceptor is enabled
	commonTemplate = `
 local interceptor = require 'home.wso2.interceptor.lib.interceptor'
 local utils = require 'home.wso2.interceptor.lib.utils'
 {{if .IsResponseFlowEnabled -}} {{/* resp_flow details are required in req flow if request info needed in resp flow */}}
 local resp_flow_list = {  
	 {{- range $key, $value := .ResponseFlow -}} 
		 {{- $key }} = {invocationContext = {{$value.Include.InvocationContext}}, requestHeaders = {{$value.Include.RequestHeaders}}, requestBody = {{$value.Include.RequestBody}}, requestTrailer = {{$value.Include.RequestTrailer}}, responseHeaders = {{$value.Include.ResponseHeaders}}, responseBody = {{$value.Include.ResponseBody}}, responseTrailers = {{$value.Include.ResponseTrailers}}}, 
	 {{- end -}}}
 {{- else -}}local resp_flow_list = {}{{end}} {{/* if resp_flow disabled no need req info in resp path */}}
 local inv_context = {
	 organizationId = "{{.Context.OrganizationID}}",
	 basePath = "{{.Context.BasePath}}",
	 supportedMethods = "{{.Context.SupportedMethods}}",
	 apiName = "{{.Context.APIName}}",
	 apiVersion = "{{.Context.APIVersion}}",
	 pathTemplate = "{{.Context.PathTemplate}}",
	 vhost = "{{.Context.Vhost}}",
	 prodClusterName = "{{.Context.ProdClusterName}}",
	 sandClusterName = "{{.Context.SandClusterName}}"
 }
 local wire_log_config = {
	 log_body_enabled = {{ .LogConfig.LogBodyEnabled }},
	 log_headers_enabled = {{ .LogConfig.LogHeadersEnabled }},
	 log_trailers_enabled = {{ .LogConfig.LogTrailersEnabled }}
 }
 `
	requestInterceptorTemplate = `
 local req_flow_list = {  
	 {{- range $key, $value := .RequestFlow -}} 
		 {{- $key }} = {invocationContext = {{$value.Include.InvocationContext}}, requestHeaders = {{$value.Include.RequestHeaders}}, requestBody = {{$value.Include.RequestBody}}, requestTrailer = {{$value.Include.RequestTrailer}}}, 
	 {{- end -}}}
 local req_call_config = {  
	 {{- range $key, $value := .RequestFlow -}} 
		 {{- $key }} = {cluster_name = "{{$value.ExternalCall.ClusterName}}", timeout = {{$value.ExternalCall.Timeout}}, authority_header = "{{$value.ExternalCall.AuthorityHeader}}"}, 
	 {{- end -}}}
 function envoy_on_request(request_handle)
	 interceptor.handle_request_interceptor(
		 request_handle, req_call_config, req_flow_list, resp_flow_list, inv_context, false, wire_log_config
	 )
 end
 `

	responseInterceptorTemplate = `
 local res_call_config = {  
	 {{- range $key, $value := .ResponseFlow -}} 
		 {{- $key }} = {cluster_name = "{{$value.ExternalCall.ClusterName}}", timeout={{$value.ExternalCall.Timeout}}, authority_header = "{{$value.ExternalCall.AuthorityHeader}}"}, 
	 {{- end -}}}
 function envoy_on_response(response_handle)
	 interceptor.handle_response_interceptor(
		 response_handle, res_call_config, resp_flow_list, wire_log_config
	 )
 end
 `
	// defaultRequestInterceptorTemplate is the template that is applied when request flow is disabled
	// just updated req flow info with  resp flow without calling interceptor service
	defaultRequestInterceptorTemplate = `
 function envoy_on_request(request_handle)
	 utils.wire_log_headers(request_handle, " >> request headers >> ", {{ .LogConfig.LogHeadersEnabled }})
	 utils.wire_log_body(request_handle, " >> request body >> ", {{ .LogConfig.LogBodyEnabled }})
	 interceptor.handle_request_interceptor(request_handle, {}, {}, resp_flow_list, inv_context, true, { log_body_enabled = false, log_headers_enabled = false, log_trailers_enabled = false })
	 utils.wire_log_trailers(request_handle, " >> request trailers >> ", {{ .LogConfig.LogTrailersEnabled }})
 end
 `
	// defaultResponseInterceptorTemplate is the template that is applied when response flow is disabled
	defaultResponseInterceptorTemplate = `
 function envoy_on_response(response_handle)
	 utils.wire_log(response_handle, " << response body << ", " << response headers << ", " << response trailers << ", wire_log_config)
 end
 `
	// emptyRequestInterceptorTemplate is the template that is applied when request flow and response flow is disabled
	emptyInterceptorTemplate = `
 function envoy_on_request(request_handle)
 end
 function envoy_on_response(response_handle)
 end
 `
)

//GetInterceptor inject values and get request interceptor
// Note: This method is called only if one of request or response interceptor is enabled
func GetInterceptor(templateValues any, templateString string) string {
	t, err := template.New("lua-filter").Parse(templateString)
	if err != nil {
		logger.LoggerInterceptor.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprintf("error while parsing the interceptor template: %v", err.Error()),
			Severity:  logging.CRITICAL,
			ErrorCode: 1800,
		})
		return emptyInterceptorTemplate
	}
	templ := template.Must(t, err)
	var out bytes.Buffer

	err = templ.Execute(&out, templateValues)
	if err != nil {
		logger.LoggerInterceptor.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprintf("executing request interceptor template: %v", err.Error()),
			Severity:  logging.CRITICAL,
			ErrorCode: 1801,
		})
		return emptyInterceptorTemplate
	}
	return out.String()
}

// GetTemplate returns the combined tempalate for the interceptor lua script
func GetTemplate(isReqIntercept bool, isResIntercept bool) string {
	reqT := defaultRequestInterceptorTemplate
	resT := defaultResponseInterceptorTemplate
	if isReqIntercept {
		reqT = requestInterceptorTemplate
	}
	if isResIntercept {
		resT = responseInterceptorTemplate
	}
	return commonTemplate + reqT + resT
}
