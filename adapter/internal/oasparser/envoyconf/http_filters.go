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

// Package envoyconf generates the envoyconfiguration for listeners, virtual hosts,
// routes, clusters, and endpoints.
package envoyconf

import (
	"time"

	"strings"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_config_ratelimit_v3 "github.com/envoyproxy/go-control-plane/envoy/config/ratelimit/v3"
	cors_filter_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/cors/v3"
	ext_authv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/ext_authz/v3"
	ext_process "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/ext_proc/v3"
	grpc_stats_filter_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/grpc_stats/v3"
	luav3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/lua/v3"
	ratelimit "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/ratelimit/v3"
	routerv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/router/v3"
	wasm_filter_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/wasm/v3"
	hcmv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	wasmv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/wasm/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/wso2/apk/adapter/config"
	logger "github.com/wso2/apk/adapter/internal/loggers"
	logging "github.com/wso2/apk/adapter/internal/logging"
	"google.golang.org/protobuf/proto"

	"github.com/golang/protobuf/ptypes/any"
)

// HTTPExternalProcessor HTTP filter
const HTTPExternalProcessor = "envoy.filters.http.ext_proc"

// RatelimitFilterName Ratelimit filter name
const RatelimitFilterName = "envoy.filters.http.ratelimit"

// getHTTPFilters generates httpFilter configuration
func getHTTPFilters(globalLuaScript string) []*hcmv3.HttpFilter {
	// extAuth := getExtAuthzHTTPFilter()
	extProcessor := getExtProcessHTTPFilter()
	router := getRouterHTTPFilter()
	luaLocal := getLuaFilter(LuaLocal, `
function envoy_on_request(request_handle)
end
function envoy_on_response(response_handle)
end`)
	luaGlobal := getLuaFilter(LuaGlobal, globalLuaScript)
	cors := getCorsHTTPFilter()

	httpFilters := []*hcmv3.HttpFilter{
		cors,
		// extAuth,
		luaLocal,
		luaGlobal,
		extProcessor,
	}
	conf := config.ReadConfigs()
	if conf.Envoy.RateLimit.Enabled {
		rateLimit := getRateLimitFilter()
		httpFilters = append(httpFilters, rateLimit)
	}
	if conf.Envoy.Filters.Compression.Enabled {
		compressionFilter, err := getCompressorFilter()
		if err != nil {
			logger.LoggerXds.ErrorC(logging.PrintError(logging.Error2234, logging.MINOR, "Error occurred while creating the compression filter: %v", err.Error()))
			return httpFilters
		}
		httpFilters = append(httpFilters, compressionFilter)
	}
	httpFilters = append(httpFilters, router)
	return httpFilters
}

// getRouterHTTPFilter gets router http filter.
func getRouterHTTPFilter() *hcmv3.HttpFilter {

	routeFilterConf := routerv3.Router{
		DynamicStats:             nil,
		StartChildSpan:           false,
		UpstreamLog:              nil,
		SuppressEnvoyHeaders:     true,
		StrictCheckHeaders:       nil,
		RespectExpectedRqTimeout: false,
	}

	routeFilterTypedConf, err := anypb.New(&routeFilterConf)
	if err != nil {
		logger.LoggerOasparser.Error("Error marshaling route filter configs. ", err)
	}
	filter := hcmv3.HttpFilter{
		Name:       wellknown.Router,
		ConfigType: &hcmv3.HttpFilter_TypedConfig{TypedConfig: routeFilterTypedConf},
	}
	return &filter
}

// getGRPCStatsHTTPFilter gets grpc_stats http filter.
func getGRPCStatsHTTPFilter() *hcmv3.HttpFilter {

	gprcStatsFilterConf := grpc_stats_filter_v3.FilterConfig{
		EnableUpstreamStats: true,
		EmitFilterState:     true,
	}
	gprcStatsFilterTypedConf, err := anypb.New(&gprcStatsFilterConf)

	if err != nil {
		logger.LoggerOasparser.Error("Error marshaling grpc stats filter configs. ", err)
	}

	filter := hcmv3.HttpFilter{
		Name:       "grpc_stats",
		ConfigType: &hcmv3.HttpFilter_TypedConfig{TypedConfig: gprcStatsFilterTypedConf},
	}

	return &filter
}

// getCorsHTTPFilter gets cors http filter.
func getCorsHTTPFilter() *hcmv3.HttpFilter {

	corsFilterConf := cors_filter_v3.CorsPolicy{}
	corsFilterTypedConf, err := anypb.New(&corsFilterConf)

	if err != nil {
		logger.LoggerOasparser.Error("Error marshaling cors filter configs. ", err)
	}

	filter := hcmv3.HttpFilter{
		Name:       wellknown.CORS,
		ConfigType: &hcmv3.HttpFilter_TypedConfig{TypedConfig: corsFilterTypedConf},
	}

	return &filter
}

// UpgradeFilters that are applied in websocket upgrade mode
func getUpgradeFilters() []*hcmv3.HttpFilter {

	cors := getCorsHTTPFilter()
	grpcStats := getGRPCStatsHTTPFilter()
	// extAauth := getExtAuthzHTTPFilter()
	apkWebSocketWASM := getAPKWebSocketWASMFilter()
	router := getRouterHTTPFilter()
	upgradeFilters := []*hcmv3.HttpFilter{
		cors,
		grpcStats,
		// extAauth,
		apkWebSocketWASM,
		router,
	}
	return upgradeFilters
}

// getRateLimitFilter configures the ratelimit filter
func getRateLimitFilter() *hcmv3.HttpFilter {
	conf := config.ReadConfigs()

	// X-RateLimit Headers
	var enableXRatelimitHeaders ratelimit.RateLimit_XRateLimitHeadersRFCVersion
	if conf.Envoy.RateLimit.XRateLimitHeaders.Enabled {
		switch strings.ToUpper(conf.Envoy.RateLimit.XRateLimitHeaders.RFCVersion) {
		case ratelimit.RateLimit_DRAFT_VERSION_03.String():
			enableXRatelimitHeaders = ratelimit.RateLimit_DRAFT_VERSION_03
		default:
			defaultType := ratelimit.RateLimit_DRAFT_VERSION_03
			logger.LoggerOasparser.ErrorC(logging.PrintError(logging.Error2240, logging.MAJOR, "Invalid XRatelimitHeaders type, continue with default type %s", defaultType))
			enableXRatelimitHeaders = defaultType
		}
	} else {
		enableXRatelimitHeaders = ratelimit.RateLimit_OFF
	}

	rateLimit := &ratelimit.RateLimit{
		Domain:                  RateLimiterDomain,
		FailureModeDeny:         conf.Envoy.RateLimit.FailureModeDeny,
		EnableXRatelimitHeaders: enableXRatelimitHeaders,
		Timeout: &durationpb.Duration{
			Nanos:   (int32(conf.Envoy.RateLimit.RequestTimeoutInMillis) % 1000) * 1000000,
			Seconds: conf.Envoy.RateLimit.RequestTimeoutInMillis / 1000,
		},
		RateLimitService: &envoy_config_ratelimit_v3.RateLimitServiceConfig{
			TransportApiVersion: corev3.ApiVersion_V3,
			GrpcService: &corev3.GrpcService{
				TargetSpecifier: &corev3.GrpcService_EnvoyGrpc_{
					EnvoyGrpc: &corev3.GrpcService_EnvoyGrpc{
						ClusterName: rateLimitClusterName,
					},
				},
				Timeout: &durationpb.Duration{
					Nanos:   (int32(conf.Envoy.RateLimit.RequestTimeoutInMillis) % 1000) * 1000000,
					Seconds: conf.Envoy.RateLimit.RequestTimeoutInMillis / 1000,
				},
			},
		},
	}
	ext, err2 := anypb.New(rateLimit)
	if err2 != nil {
		logger.LoggerOasparser.ErrorC(logging.PrintError(logging.Error2241, logging.MAJOR, "Error occurred while parsing ratelimit filter config. Error: %s", err2.Error()))
	}
	rlFilter := hcmv3.HttpFilter{
		Name: wellknown.HTTPRateLimit,
		ConfigType: &hcmv3.HttpFilter_TypedConfig{
			TypedConfig: ext,
		},
	}
	return &rlFilter
}

// getExtProcessHTTPFilter gets ExtAauthz http filter.
func getExtProcessHTTPFilter() *hcmv3.HttpFilter {
	conf := config.ReadConfigs()
	externalProcessor := &ext_process.ExternalProcessor{
		GrpcService: &corev3.GrpcService{
			TargetSpecifier: &corev3.GrpcService_EnvoyGrpc_{
				EnvoyGrpc: &corev3.GrpcService_EnvoyGrpc{
					ClusterName: extAuthzClusterName,
				},
			},
			Timeout: durationpb.New(conf.Envoy.EnforcerResponseTimeoutInSeconds * time.Second),
		},
		FailureModeAllow: true,
		ProcessingMode: &ext_process.ProcessingMode{
			ResponseBodyMode:   ext_process.ProcessingMode_BUFFERED,
			RequestHeaderMode:  ext_process.ProcessingMode_SEND,
			ResponseHeaderMode: ext_process.ProcessingMode_SEND,
			// RequestHeaderMode:  ext_process.ProcessingMode_SKIP,
			// ResponseHeaderMode: ext_process.ProcessingMode_SKIP,
			RequestBodyMode:   ext_process.ProcessingMode_BUFFERED,
		},
		MetadataOptions: &ext_process.MetadataOptions{
			ForwardingNamespaces: &ext_process.MetadataOptions_MetadataNamespaces{
				Untyped: []string{"envoy.filters.http.ext_authz", "envoy.filters.http.ext_proc"},
			},
			ReceivingNamespaces: &ext_process.MetadataOptions_MetadataNamespaces{
				Untyped: []string{"envoy.filters.http.ext_proc"},
			},
		},
		RequestAttributes:  []string{"xds.route_metadata", "request.method"},
		ResponseAttributes: []string{"xds.route_metadata"},
		MessageTimeout:     durationpb.New(conf.Envoy.EnforcerResponseTimeoutInSeconds * time.Second),
	}
	ext, err2 := anypb.New(externalProcessor)
	if err2 != nil {
		logger.LoggerOasparser.Error(err2)
	}
	extProcessFilter := hcmv3.HttpFilter{
		Name: HTTPExternalProcessor,
		ConfigType: &hcmv3.HttpFilter_TypedConfig{
			TypedConfig: ext,
		},
	}
	return &extProcessFilter
}

// getExtAuthzHTTPFilter gets ExtAauthz http filter.
func getExtAuthzHTTPFilter() *hcmv3.HttpFilter {
	conf := config.ReadConfigs()
	extAuthzConfig := &ext_authv3.ExtAuthz{
		// This would clear the route cache only if there is a header added/removed or changed
		// within ext-authz filter.
		ClearRouteCache:        true,
		IncludePeerCertificate: true,
		TransportApiVersion:    corev3.ApiVersion_V3,
		Services: &ext_authv3.ExtAuthz_GrpcService{
			GrpcService: &corev3.GrpcService{
				TargetSpecifier: &corev3.GrpcService_EnvoyGrpc_{
					EnvoyGrpc: &corev3.GrpcService_EnvoyGrpc{
						ClusterName: extAuthzClusterName,
					},
				},
				Timeout: durationpb.New(conf.Envoy.EnforcerResponseTimeoutInSeconds * time.Second),
				InitialMetadata: []*corev3.HeaderValue{
					{
						Key:   "x-request-id",
						Value: "%REQ(x-request-id)%",
					},
				},
			},
		},
	}

	// configures envoy to handle request body and GraphQL APIs require below configs to pass request
	// payload to the enforcer.
	extAuthzConfig.WithRequestBody = &ext_authv3.BufferSettings{
		MaxRequestBytes:     conf.Envoy.PayloadPassingToEnforcer.MaxRequestBytes,
		AllowPartialMessage: conf.Envoy.PayloadPassingToEnforcer.AllowPartialMessage,
		PackAsBytes:         conf.Envoy.PayloadPassingToEnforcer.PackAsBytes,
	}

	ext, err2 := anypb.New(extAuthzConfig)
	if err2 != nil {
		logger.LoggerOasparser.Error(err2)
	}
	extAuthzFilter := hcmv3.HttpFilter{
		Name: wellknown.HTTPExternalAuthorization,
		ConfigType: &hcmv3.HttpFilter_TypedConfig{
			TypedConfig: ext,
		},
	}
	return &extAuthzFilter
}

// getLuaFilter gets Lua http filter.
func getLuaFilter(filterName, defaultScript string) *hcmv3.HttpFilter {
	luaConfig := &luav3.Lua{
		DefaultSourceCode: &corev3.DataSource{
			Specifier: &corev3.DataSource_InlineString{
				InlineString: defaultScript,
			},
		},
	}
	ext, err2 := anypb.New(luaConfig)
	if err2 != nil {
		logger.LoggerOasparser.Error(err2)
	}
	luaFilter := hcmv3.HttpFilter{
		Name: filterName,
		ConfigType: &hcmv3.HttpFilter_TypedConfig{
			TypedConfig: ext,
		},
	}
	return &luaFilter
}

func getAPKWebSocketWASMFilter() *hcmv3.HttpFilter {
	config := &wrappers.StringValue{
		Value: `{
			"node_id": "apk_node_1",
			"rate_limit_service": "ext-authz",
			"timeout": "20s",
			"failure_mode_deny": "true"
		}`,
	}
	a, err := proto.Marshal(config)
	if err != nil {
		logger.LoggerOasparser.Error(err)
	}
	apkWebsocketWASMConfig := wasmv3.PluginConfig{
		Name:   apkWebSocketWASMFilterName,
		RootId: apkWebSocketWASMFilterRoot,
		Vm: &wasmv3.PluginConfig_VmConfig{
			VmConfig: &wasmv3.VmConfig{
				VmId:             apkWASMVmID,
				Runtime:          apkWASMVmRuntime,
				AllowPrecompiled: true,
				Code: &corev3.AsyncDataSource{
					Specifier: &corev3.AsyncDataSource_Local{
						Local: &corev3.DataSource{
							Specifier: &corev3.DataSource_Filename{
								Filename: apkWebSocketWASM,
							},
						},
					},
				},
			},
		},
		Configuration: &any.Any{
			TypeUrl: "type.googleapis.com/google.protobuf.StringValue",
			Value:   a,
		},
	}

	apkWebSocketWASMFilterConfig := &wasm_filter_v3.Wasm{
		Config: &apkWebsocketWASMConfig,
	}

	ext, err2 := proto.Marshal(apkWebSocketWASMFilterConfig)
	if err2 != nil {
		logger.LoggerOasparser.Error(err2)
	}
	apkWebSocketFilter := hcmv3.HttpFilter{
		Name: apkWebSocketWASMFilterName,
		ConfigType: &hcmv3.HttpFilter_TypedConfig{
			TypedConfig: &any.Any{
				TypeUrl: "type.googleapis.com/envoy.extensions.filters.http.wasm.v3.Wasm",
				Value:   ext,
			},
		},
	}
	return &apkWebSocketFilter

}
