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

package envoyconf

import (
	"errors"
	"fmt"
	"net"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	clusterv3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	endpointv3 "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	routev3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	cors_filter_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/cors/v3"
	extAuthService "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/ext_authz/v3"
	lua "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/lua/v3"
	tlsv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
	upstreams_http_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/upstreams/http/v3"
	envoy_type_matcherv3 "github.com/envoyproxy/go-control-plane/envoy/type/matcher/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"

	"github.com/wso2/apk/adapter/config"
	"github.com/wso2/apk/adapter/internal/interceptor"
	logger "github.com/wso2/apk/adapter/internal/loggers"
	logging "github.com/wso2/apk/adapter/internal/logging"
	"github.com/wso2/apk/adapter/internal/oasparser/constants"
	"github.com/wso2/apk/adapter/internal/oasparser/model"
	"github.com/wso2/apk/adapter/internal/svcdiscovery"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/wrappers"
	gwapiv1b1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

// WireLogValues holds debug logging related template values
type WireLogValues struct {
	LogConfig *config.WireLogConfig
}

// CombinedTemplateValues holds combined values for both WireLogValues properties and Interceptor properties in the same level
type CombinedTemplateValues struct {
	WireLogValues
	interceptor.Interceptor
}

// CreateRoutesWithClusters creates envoy routes along with clusters and endpoint instances.
// This creates routes for all the swagger resources and link to clusters.
// Create clusters for endpoints.
// If a resource has resource level endpoint, it create another cluster and
// link it. If resources doesn't has resource level endpoints, those clusters are linked
// to the api level clusters.
func CreateRoutesWithClusters(mgwSwagger model.MgwSwagger, interceptorCerts map[string][]byte, vHost string, organizationID string) (routesP []*routev3.Route,
	clustersP []*clusterv3.Cluster, addressesP []*corev3.Address, err error) {
	var (
		routes    []*routev3.Route
		clusters  []*clusterv3.Cluster
		endpoints []*corev3.Address
	)

	apiTitle := mgwSwagger.GetTitle()
	apiVersion := mgwSwagger.GetVersion()

	conf := config.ReadConfigs()
	timeout := conf.Envoy.ClusterTimeoutInSeconds

	// Create API level interceptor clusters if required
	clustersI, endpointsI, apiRequestInterceptor, apiResponseInterceptor := createInterceptorAPIClusters(mgwSwagger,
		interceptorCerts, vHost, organizationID)
	clusters = append(clusters, clustersI...)
	endpoints = append(endpoints, endpointsI...)

	// Maintain a clusterName-EndpointCluster mapping to prevent duplicate
	// creation of clusters.
	processedEndpoints := map[string]model.EndpointCluster{}

	for _, resource := range mgwSwagger.GetResources() {
		var clusterName string
		resourcePath := resource.GetPath()
		endpoint := resource.GetEndpoints()
		basePath := strings.TrimSuffix(endpoint.Endpoints[0].Basepath, "/")
		existingClusterName := getExistingClusterName(*endpoint, processedEndpoints)
		
		if existingClusterName == "" {
			clusterName = getClusterName(endpoint.EndpointPrefix, organizationID, vHost, mgwSwagger.GetTitle(), apiVersion, resource.GetID())
			cluster, address, err := processEndpoints(clusterName, endpoint, timeout, basePath)
			if err != nil {
				logger.LoggerOasparser.ErrorC(logging.GetErrorByCode(2239, apiTitle, apiVersion, resourcePath, err.Error()))
			} else {
				clusters = append(clusters, cluster)
				endpoints = append(endpoints, address...)
				processedEndpoints[clusterName] = *endpoint
			}
		} else {
			clusterName = existingClusterName
		}
		// Create resource level interceptor clusters if required
		clustersI, endpointsI, operationalReqInterceptors, operationalRespInterceptorVal := createInterceptorResourceClusters(mgwSwagger,
			interceptorCerts, vHost, organizationID, apiRequestInterceptor, apiResponseInterceptor, resource)
		clusters = append(clusters, clustersI...)
		endpoints = append(endpoints, endpointsI...)
		routeP, err := createRoutes(genRouteCreateParams(&mgwSwagger, resource, vHost, basePath, clusterName, *operationalReqInterceptors, *operationalRespInterceptorVal, organizationID, false))
		if err != nil {
			logger.LoggerXds.ErrorC(logging.GetErrorByCode(2231, mgwSwagger.GetTitle(), mgwSwagger.GetVersion(), resource.GetPath(), err.Error()))
			return nil, nil, nil, fmt.Errorf("error while creating routes. %v", err)
		}
		routes = append(routes, routeP...)
	}

	return routes, clusters, endpoints, nil
}

func getClusterName(epPrefix string, organizationID string, vHost string, swaggerTitle string, swaggerVersion string,
	resourceID string) string {
	if resourceID != "" {
		return strings.TrimSpace(organizationID+"_"+epPrefix+"_"+vHost+"_"+strings.Replace(swaggerTitle, " ", "", -1)+swaggerVersion) +
			"_" + strings.Replace(resourceID, " ", "", -1) + "0"
	}
	return strings.TrimSpace(organizationID + "_" + epPrefix + "_" + vHost + "_" + strings.Replace(swaggerTitle, " ", "", -1) +
		swaggerVersion)
}

func getExistingClusterName(endpoint model.EndpointCluster, clusterEndpointMapping map[string]model.EndpointCluster) string {
	for clusterName, endpointValue := range clusterEndpointMapping {
		if reflect.DeepEqual(endpoint, endpointValue) {
			return clusterName
		}
	}
	return ""
}

// CreateLuaCluster creates lua cluster configuration.
func CreateLuaCluster(interceptorCerts map[string][]byte, endpoint model.InterceptEndpoint) (*clusterv3.Cluster, []*corev3.Address, error) {
	logger.LoggerOasparser.Debug("creating a lua cluster ", endpoint.ClusterName)
	return processEndpoints(endpoint.ClusterName, &endpoint.EndpointCluster, endpoint.ClusterTimeout, endpoint.EndpointCluster.Endpoints[0].Basepath)
}

// CreateTracingCluster creates a cluster definition for router's tracing server.
func CreateTracingCluster(conf *config.Config) (*clusterv3.Cluster, []*corev3.Address, error) {
	var epHost string
	var epPort uint32
	var epPath string
	epTimeout := conf.Envoy.ClusterTimeoutInSeconds
	epCluster := &model.EndpointCluster{
		Endpoints: []model.Endpoint{
			{
				Host:    "",
				URLType: "http",
				Port:    uint32(9411),
			},
		},
	}

	if epHost = conf.Tracing.ConfigProperties[tracerHost]; len(epHost) <= 0 {
		return nil, nil, errors.New("invalid host provided for tracing endpoint")
	}
	if epPath = conf.Tracing.ConfigProperties[tracerEndpoint]; len(epPath) <= 0 {
		return nil, nil, errors.New("invalid endpoint path provided for tracing endpoint")
	}
	if port, err := strconv.ParseUint(conf.Tracing.ConfigProperties[tracerPort], 10, 32); err == nil {
		epPort = uint32(port)
	} else {
		return nil, nil, errors.New("invalid port provided for tracing endpoint")
	}

	epCluster.Endpoints[0].Host = epHost
	epCluster.Endpoints[0].Port = epPort
	epCluster.Endpoints[0].Basepath = epPath

	return processEndpoints(tracingClusterName, epCluster, epTimeout, epPath)
}

// processEndpoints creates cluster configuration. AddressConfiguration, cluster name and
// urlType (http or https) is required to be provided.
// timeout cluster timeout
func processEndpoints(clusterName string, clusterDetails *model.EndpointCluster,
	timeout time.Duration, basePath string) (*clusterv3.Cluster, []*corev3.Address, error) {
	// tls configs
	var transportSocketMatches []*clusterv3.Cluster_TransportSocketMatch
	// create loadbalanced/failover endpoints
	var lbEPs []*endpointv3.LocalityLbEndpoints
	// failover priority
	priority := 0
	// epType {loadbalance, failover}
	epType := clusterDetails.EndpointType

	addresses := []*corev3.Address{}

	for i, ep := range clusterDetails.Endpoints {
		// validating the basepath to be same for all upstreams of an api
		if strings.TrimSuffix(ep.Basepath, "/") != basePath {
			return nil, nil, errors.New("endpoint basepath mismatched for " + ep.RawURL + ". expected : " + basePath + " but found : " + ep.Basepath)
		}
		// create addresses for endpoints
		address := createAddress(ep.Host, ep.Port)
		addresses = append(addresses, address)

		// create loadbalance / failover endpoints
		localityLbEndpoints := &endpointv3.LocalityLbEndpoints{
			Priority: uint32(priority),
			LbEndpoints: []*endpointv3.LbEndpoint{
				{
					HostIdentifier: &endpointv3.LbEndpoint_Endpoint{
						Endpoint: &endpointv3.Endpoint{
							Address: address,
						},
					},
				},
			},
		}

		// create tls configs
		if strings.HasPrefix(ep.URLType, httpsURLType) || strings.HasPrefix(ep.URLType, wssURLType) {
			upstreamtlsContext := createUpstreamTLSContext(ep.Certificate, ep.AllowedSANs, address, clusterDetails.HTTP2BackendEnabled)
			marshalledTLSContext, err := anypb.New(upstreamtlsContext)
			if err != nil {
				return nil, nil, errors.New("internal Error while marshalling the upstream TLS Context")
			}
			transportSocketMatch := &clusterv3.Cluster_TransportSocketMatch{
				Name: "ts" + strconv.Itoa(i),
				Match: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"lb_id": structpb.NewStringValue(strconv.Itoa(i)),
					},
				},
				TransportSocket: &corev3.TransportSocket{
					Name: transportSocketName,
					ConfigType: &corev3.TransportSocket_TypedConfig{
						TypedConfig: marshalledTLSContext,
					},
				},
			}
			transportSocketMatches = append(transportSocketMatches, transportSocketMatch)
			localityLbEndpoints.LbEndpoints[0].Metadata = &corev3.Metadata{
				FilterMetadata: map[string]*structpb.Struct{
					"envoy.transport_socket_match": {
						Fields: map[string]*structpb.Value{
							"lb_id": structpb.NewStringValue(strconv.Itoa(i)),
						},
					},
				},
			}
		}
		lbEPs = append(lbEPs, localityLbEndpoints)

		// set priority for next endpoint
		if strings.HasPrefix(epType, "failover") {
			priority = priority + 1
		}
	}
	conf := config.ReadConfigs()

	httpProtocolOptions := &upstreams_http_v3.HttpProtocolOptions{
		UpstreamProtocolOptions: &upstreams_http_v3.HttpProtocolOptions_ExplicitHttpConfig_{
			ExplicitHttpConfig: &upstreams_http_v3.HttpProtocolOptions_ExplicitHttpConfig{
				ProtocolConfig: &upstreams_http_v3.HttpProtocolOptions_ExplicitHttpConfig_HttpProtocolOptions{
					HttpProtocolOptions: &corev3.Http1ProtocolOptions{
						EnableTrailers: config.GetWireLogConfig().LogTrailersEnabled,
					},
				},
			},
		},
	}

	if clusterDetails.HTTP2BackendEnabled {
		httpProtocolOptions.UpstreamProtocolOptions = &upstreams_http_v3.HttpProtocolOptions_ExplicitHttpConfig_{
			ExplicitHttpConfig: &upstreams_http_v3.HttpProtocolOptions_ExplicitHttpConfig{
				ProtocolConfig: &upstreams_http_v3.HttpProtocolOptions_ExplicitHttpConfig_Http2ProtocolOptions{
					Http2ProtocolOptions: &corev3.Http2ProtocolOptions{
						HpackTableSize: &wrapperspb.UInt32Value{
							Value: conf.Envoy.Upstream.HTTP2.HpackTableSize,
						},
						MaxConcurrentStreams: &wrapperspb.UInt32Value{
							Value: conf.Envoy.Upstream.HTTP2.MaxConcurrentStreams,
						},
					},
				},
			},
		}
	}

	ext, err2 := proto.Marshal(httpProtocolOptions)
	if err2 != nil {
		logger.LoggerOasparser.Error(err2)
	}

	cluster := clusterv3.Cluster{
		Name:                 clusterName,
		ConnectTimeout:       durationpb.New(timeout * time.Second),
		ClusterDiscoveryType: &clusterv3.Cluster_Type{Type: clusterv3.Cluster_STRICT_DNS},
		DnsLookupFamily:      clusterv3.Cluster_V4_ONLY,
		LbPolicy:             clusterv3.Cluster_ROUND_ROBIN,
		LoadAssignment: &endpointv3.ClusterLoadAssignment{
			ClusterName: clusterName,
			Endpoints:   lbEPs,
		},
		TransportSocketMatches: transportSocketMatches,
		DnsRefreshRate:         durationpb.New(time.Duration(conf.Envoy.Upstream.DNS.DNSRefreshRate) * time.Millisecond),
		RespectDnsTtl:          conf.Envoy.Upstream.DNS.RespectDNSTtl,
		TypedExtensionProtocolOptions: map[string]*anypb.Any{
			"envoy.extensions.upstreams.http.v3.HttpProtocolOptions": &any.Any{
				TypeUrl: "type.googleapis.com/envoy.extensions.upstreams.http.v3.HttpProtocolOptions",
				Value:   ext,
			},
		},
	}

	if len(clusterDetails.Endpoints) > 1 {
		cluster.HealthChecks = createHealthCheck()
	}

	if clusterDetails.Config != nil && clusterDetails.Config.CircuitBreakers != nil {
		config := clusterDetails.Config.CircuitBreakers
		thresholds := &clusterv3.CircuitBreakers_Thresholds{}
		if config.MaxConnections > 0 {
			thresholds.MaxConnections = wrapperspb.UInt32(uint32(config.MaxConnections))
		}
		if config.MaxConnectionPools > 0 {
			thresholds.MaxConnectionPools = wrapperspb.UInt32(uint32(config.MaxConnectionPools))
		}
		if config.MaxPendingRequests > 0 {
			thresholds.MaxPendingRequests = wrapperspb.UInt32(uint32(config.MaxPendingRequests))
		}
		if config.MaxRequests > 0 {
			thresholds.MaxRequests = wrapperspb.UInt32(uint32(config.MaxRequests))
		}
		if config.MaxRetries > 0 {
			thresholds.MaxRetries = wrapperspb.UInt32(uint32(config.MaxRetries))
		}
		cluster.CircuitBreakers = &clusterv3.CircuitBreakers{
			Thresholds: []*clusterv3.CircuitBreakers_Thresholds{
				thresholds,
			},
		}
	}

	// service discovery itself will be handling loadbancing etc.
	// Therefore mutiple endpoint support is not needed, hence consider only.
	serviceDiscoveryString := clusterDetails.Endpoints[0].ServiceDiscoveryString
	if serviceDiscoveryString != "" {
		//add the api level cluster name to the ClusterConsulKeyMap
		svcdiscovery.ClusterConsulKeyMap[clusterName] = serviceDiscoveryString
		logger.LoggerOasparser.Debugln("Consul cluster added for x-wso2-endpoints: ", clusterName, " ",
			serviceDiscoveryString)
	}

	return &cluster, addresses, nil
}

func createHealthCheck() []*corev3.HealthCheck {
	conf := config.ReadConfigs()
	return []*corev3.HealthCheck{
		{
			Timeout:            durationpb.New(time.Duration(conf.Envoy.Upstream.Health.Timeout) * time.Second),
			Interval:           durationpb.New(time.Duration(conf.Envoy.Upstream.Health.Interval) * time.Second),
			UnhealthyThreshold: wrapperspb.UInt32(uint32(conf.Envoy.Upstream.Health.UnhealthyThreshold)),
			HealthyThreshold:   wrapperspb.UInt32(uint32(conf.Envoy.Upstream.Health.HealthyThreshold)),
			// we only support tcp default healthcheck
			HealthChecker: &corev3.HealthCheck_TcpHealthCheck_{},
		},
	}
}

func createUpstreamTLSContext(upstreamCerts []byte, allowedSANs []string, address *corev3.Address, hTTP2BackendEnabled bool) *tlsv3.UpstreamTlsContext {
	conf := config.ReadConfigs()
	tlsCert := generateTLSCert(conf.Envoy.KeyStore.KeyPath, conf.Envoy.KeyStore.CertPath)
	// Convert the cipher string to a string array
	ciphersArray := strings.Split(conf.Envoy.Upstream.TLS.Ciphers, ",")
	for i := range ciphersArray {
		ciphersArray[i] = strings.TrimSpace(ciphersArray[i])
	}

	upstreamTLSContext := &tlsv3.UpstreamTlsContext{
		CommonTlsContext: &tlsv3.CommonTlsContext{
			TlsParams: &tlsv3.TlsParameters{
				TlsMinimumProtocolVersion: createTLSProtocolVersion(conf.Envoy.Upstream.TLS.MinimumProtocolVersion),
				TlsMaximumProtocolVersion: createTLSProtocolVersion(conf.Envoy.Upstream.TLS.MaximumProtocolVersion),
				CipherSuites:              ciphersArray,
			},
			TlsCertificates: []*tlsv3.TlsCertificate{tlsCert},
		},
	}

	if hTTP2BackendEnabled {
		upstreamTLSContext.CommonTlsContext.AlpnProtocols = []string{"h2", "http/1.1"}
	}

	sanType := tlsv3.SubjectAltNameMatcher_IP_ADDRESS
	// Sni should be assigned when there is a hostname
	if net.ParseIP(address.GetSocketAddress().GetAddress()) == nil {
		upstreamTLSContext.Sni = address.GetSocketAddress().GetAddress()
		// If the address is an IP, then the SAN type should be changed accordingly.
		sanType = tlsv3.SubjectAltNameMatcher_DNS
	}

	if !conf.Envoy.Upstream.TLS.DisableSslVerification {
		var trustedCASrc *corev3.DataSource

		if len(upstreamCerts) > 0 {
			trustedCASrc = &corev3.DataSource{
				Specifier: &corev3.DataSource_InlineBytes{
					InlineBytes: upstreamCerts,
				},
			}
		} else {
			trustedCASrc = &corev3.DataSource{
				Specifier: &corev3.DataSource_Filename{
					Filename: conf.Envoy.Upstream.TLS.TrustedCertPath,
				},
			}
		}

		upstreamTLSContext.CommonTlsContext.ValidationContextType = &tlsv3.CommonTlsContext_ValidationContext{
			ValidationContext: &tlsv3.CertificateValidationContext{
				TrustedCa: trustedCASrc,
			},
		}
	}

	if conf.Envoy.Upstream.TLS.VerifyHostName && !conf.Envoy.Upstream.TLS.DisableSslVerification {
		addressString := address.GetSocketAddress().GetAddress()
		subjectAltNames := []*tlsv3.SubjectAltNameMatcher{
			{
				SanType: sanType,
				Matcher: &envoy_type_matcherv3.StringMatcher{
					MatchPattern: &envoy_type_matcherv3.StringMatcher_Exact{
						Exact: addressString,
					},
				},
			},
		}
		for _, san := range allowedSANs {
			subjectAltNames = append(subjectAltNames, &tlsv3.SubjectAltNameMatcher{
				SanType: sanType,
				Matcher: &envoy_type_matcherv3.StringMatcher{
					MatchPattern: &envoy_type_matcherv3.StringMatcher_SafeRegex{
						SafeRegex: &envoy_type_matcherv3.RegexMatcher{
							Regex: san,
						},
					},
				},
			})
		}
		upstreamTLSContext.CommonTlsContext.GetValidationContext().MatchTypedSubjectAltNames = subjectAltNames
	}
	return upstreamTLSContext
}

func createTLSProtocolVersion(tlsVersion string) tlsv3.TlsParameters_TlsProtocol {
	switch tlsVersion {
	case "TLS1_0":
		return tlsv3.TlsParameters_TLSv1_0
	case "TLS1_1":
		return tlsv3.TlsParameters_TLSv1_1
	case "TLS1_2":
		return tlsv3.TlsParameters_TLSv1_2
	case "TLS1_3":
		return tlsv3.TlsParameters_TLSv1_3
	default:
		return tlsv3.TlsParameters_TLS_AUTO
	}
}

// createRoutes creates route elements for the route configurations. API title, VHost, xWso2Basepath, API version,
// endpoint's basePath, resource Object (Microgateway's internal representation), clusterName needs to be provided.
func createRoutes(params *routeCreateParams) (routes []*routev3.Route, err error) {
	title := params.title
	version := params.version
	vHost := params.vHost
	xWso2Basepath := params.xWSO2BasePath
	apiType := params.apiType
	corsPolicy := getCorsPolicy(params.corsPolicy)
	resource := params.resource
	clusterName := params.clusterName
	routeConfig := params.routeConfig
	endpointBasepath := params.endpointBasePath
	requestInterceptor := params.requestInterceptor
	responseInterceptor := params.responseInterceptor
	isDefaultVersion := params.isDefaultVersion

	logger.LoggerOasparser.Debugf("creating routes for API %s ....", title)
	var (
		// The following are common to all routes and does not get updated per operation
		decorator *routev3.Decorator
	)

	basePath := strings.TrimSuffix(xWso2Basepath, "/")
	if isDefaultVersion {
		basePath = getDefaultVersionBasepath(basePath, version)
	}

	resourcePath := ""
	var resourceMethods []string
	var pathMatchType gwapiv1b1.PathMatchType
	if params.apiType == constants.GRAPHQL {
		resourceMethods = []string{"POST"}
	} else {
		resourcePath = resource.GetPath()
		resourceMethods = resource.GetMethodList()
		pathMatchType = resource.GetPathMatchType()
	}
	routePath := generateRoutePath(resourcePath, pathMatchType)

	// route path could be empty only if there is no basePath for API or the endpoint available,
	// and resourcePath is also an empty string.
	// Empty check is added to run the gateway in failsafe mode, as if the decorator string is
	// empty, the route configuration does not apply.
	if strings.TrimSpace(routePath) != "" {
		decorator = &routev3.Decorator{
			Operation: vHost + ":" + routePath,
		}
	}

	var contextExtensions = make(map[string]string)
	contextExtensions[pathContextExtension] = resourcePath
	contextExtensions[vHostContextExtension] = vHost
	if xWso2Basepath != "" {
		contextExtensions[basePathContextExtension] = xWso2Basepath
	} else {
		contextExtensions[basePathContextExtension] = endpointBasepath
	}
	contextExtensions[methodContextExtension] = strings.Join(resourceMethods, " ")
	contextExtensions[apiVersionContextExtension] = version
	contextExtensions[apiNameContextExtension] = title
	// One of these values will be selected and added as the cluster-header http header
	// from enhancer
	// Even if the routing is based on direct cluster, these properties needs to be populated
	// to validate the key type component in the token.
	contextExtensions[clusterNameContextExtension] = clusterName

	extAuthPerFilterConfig := extAuthService.ExtAuthzPerRoute{
		Override: &extAuthService.ExtAuthzPerRoute_CheckSettings{
			CheckSettings: &extAuthService.CheckSettings{
				ContextExtensions: contextExtensions,
				// negation is performing to match the envoy config name (disable_request_body_buffering)
				DisableRequestBodyBuffering: !params.passRequestPayloadToEnforcer,
			},
		},
	}

	b := proto.NewBuffer(nil)
	b.SetDeterministic(true)
	_ = b.Marshal(&extAuthPerFilterConfig)
	extAuthzFilter := &any.Any{
		TypeUrl: extAuthzPerRouteName,
		Value:   b.Bytes(),
	}

	var luaPerFilterConfig lua.LuaPerRoute
	if len(requestInterceptor) < 1 && len(responseInterceptor) < 1 {

		logConf := config.ReadLogConfigs()

		if logConf.WireLogs.Enable {

			templateString := `
local utils = require 'home.wso2.interceptor.lib.utils'
local wire_log_config = {
	log_body_enabled = {{ .LogConfig.LogBodyEnabled }},
	log_headers_enabled = {{ .LogConfig.LogHeadersEnabled }},
	log_trailers_enabled = {{ .LogConfig.LogTrailersEnabled }}
}
function envoy_on_request(request_handle)
	utils.wire_log(request_handle, " >> request body >> ", " >> request headers >> ", " >> request trailers >> ", wire_log_config)
end

function envoy_on_response(response_handle)
	utils.wire_log(response_handle, " << response body << ", " << response headers << ", " << response trailers << ", wire_log_config)
end`
			templateValues := WireLogValues{
				LogConfig: config.GetWireLogConfig(),
			}
			luaPerFilterConfig = lua.LuaPerRoute{
				Override: &lua.LuaPerRoute_SourceCode{SourceCode: &corev3.DataSource{Specifier: &corev3.DataSource_InlineString{
					InlineString: interceptor.GetInterceptor(templateValues, templateString),
				}}},
			}
		} else {
			luaPerFilterConfig = lua.LuaPerRoute{
				Override: &lua.LuaPerRoute_Disabled{Disabled: true},
			}
		}

	} else {
		// read from contextExtensions map since, it is updated with correct values with conditions
		// so, no need to change two places
		iInvCtx := &interceptor.InvocationContext{
			OrganizationID:   params.organizationID,
			BasePath:         contextExtensions[basePathContextExtension],
			SupportedMethods: contextExtensions[methodContextExtension],
			APIName:          contextExtensions[apiNameContextExtension],
			APIVersion:       contextExtensions[apiVersionContextExtension],
			PathTemplate:     contextExtensions[pathContextExtension],
			Vhost:            contextExtensions[vHostContextExtension],
			ClusterName:      contextExtensions[clusterNameContextExtension],
		}
		luaPerFilterConfig = lua.LuaPerRoute{
			Override: &lua.LuaPerRoute_SourceCode{
				SourceCode: &corev3.DataSource{
					Specifier: &corev3.DataSource_InlineString{
						InlineString: getInlineLuaScript(requestInterceptor, responseInterceptor, iInvCtx),
					},
				},
			},
		}
	}

	luaMarshelled := proto.NewBuffer(nil)
	luaMarshelled.SetDeterministic(true)
	_ = luaMarshelled.Marshal(&luaPerFilterConfig)

	luaFilter := &any.Any{
		TypeUrl: luaPerRouteName,
		Value:   luaMarshelled.Bytes(),
	}

	corsFilter, _ := anypb.New(corsPolicy)

	perRouteFilterConfigs := map[string]*any.Any{
		wellknown.HTTPExternalAuthorization: extAuthzFilter,
		wellknown.Lua:                       luaFilter,
		wellknown.CORS:                      corsFilter,
	}

	logger.LoggerOasparser.Debugf("adding route : %s for API : %s", resourcePath, title)

	if resource != nil && resource.HasPolicies() {
		logger.LoggerOasparser.Debug("Start creating routes for resource with policies")

		// Policies are per operation (HTTP method). Therefore, create route per HTTP method.
		for _, operation := range resource.GetOperations() {
			var requestHeadersToAdd []*corev3.HeaderValueOption
			var requestHeadersToRemove []string
			var responseHeadersToAdd []*corev3.HeaderValueOption
			var responseHeadersToRemove []string
			var pathRewriteConfig *envoy_type_matcherv3.RegexMatchAndSubstitute

			hasMethodRewritePolicy := false
			var newMethod string

			// Policies - for request flow
			for _, requestPolicy := range operation.GetPolicies().Request {
				logger.LoggerOasparser.Debug("Adding request flow policies for ", resourcePath, operation.GetMethod())
				switch requestPolicy.Action {

				case constants.ActionHeaderAdd:
					logger.LoggerOasparser.Debugf("Adding %s policy to request flow for %s %s",
						constants.ActionHeaderAdd, resourcePath, operation.GetMethod())
					requestHeaderToAdd, err := generateHeaderToAddRouteConfig(requestPolicy.Parameters)
					if err != nil {
						return nil, fmt.Errorf("error adding request policy %s to operation %s of resource %s."+
							" %v", requestPolicy.Action, operation.GetMethod(), resourcePath, err)
					}
					requestHeadersToAdd = append(requestHeadersToAdd, requestHeaderToAdd)

				case constants.ActionHeaderRemove:
					logger.LoggerOasparser.Debugf("Adding %s policy to request flow for %s %s",
						constants.ActionHeaderRemove, resourcePath, operation.GetMethod())
					requestHeaderToRemove, err := generateHeaderToRemoveString(requestPolicy.Parameters)
					if err != nil {
						return nil, fmt.Errorf("error adding request policy %s to operation %s of resource %s."+
							" %v", requestPolicy.Action, operation.GetMethod(), resourcePath, err)
					}
					requestHeadersToRemove = append(requestHeadersToRemove, requestHeaderToRemove)

				case constants.ActionRewritePath:
					logger.LoggerOasparser.Debugf("Adding %s policy to request flow for %s %s",
						constants.ActionRewritePath, resourcePath, operation.GetMethod())
					regexRewrite, err := generateRewritePathRouteConfig(routePath, resourcePath, endpointBasepath,
						requestPolicy.Parameters, pathMatchType)
					if err != nil {
						errMsg := fmt.Sprintf("Error adding request policy %s to operation %s of resource %s. %v",
							constants.ActionRewritePath, operation.GetMethod(), resourcePath, err)
						logger.LoggerOasparser.ErrorC(logging.GetErrorByCode(2212, constants.ActionRewritePath, operation.GetMethod(), resourcePath, err))
						return nil, errors.New(errMsg)
					}
					pathRewriteConfig = regexRewrite

				case constants.ActionRewriteMethod:
					logger.LoggerOasparser.Debugf("Adding %s policy to request flow for %s %s",
						constants.ActionRewriteMethod, resourcePath, operation.GetMethod())
					hasMethodRewritePolicy, err = isMethodRewrite(resourcePath, operation.GetMethod(), requestPolicy.Parameters)
					if err != nil {
						return nil, err
					}
					if !hasMethodRewritePolicy {
						continue
					}
					newMethod, err = getRewriteMethod(resourcePath, operation.GetMethod(), requestPolicy.Parameters)
					if err != nil {
						return nil, err
					}
				}
			}

			// Policies - for response flow
			for _, responsePolicy := range operation.GetPolicies().Response {
				logger.LoggerOasparser.Debug("Adding response flow policies for ", resourcePath, operation.GetMethod())
				switch responsePolicy.Action {

				case constants.ActionHeaderAdd:
					logger.LoggerOasparser.Debugf("Adding %s policy to response flow for %s %s",
						constants.ActionHeaderAdd, resourcePath, operation.GetMethod())
					responseHeaderToAdd, err := generateHeaderToAddRouteConfig(responsePolicy.Parameters)
					if err != nil {
						return nil, fmt.Errorf("error adding response policy %s to operation %s of resource %s."+
							" %v", responsePolicy.Action, operation.GetMethod(), resourcePath, err)
					}
					responseHeadersToAdd = append(responseHeadersToAdd, responseHeaderToAdd)

				case constants.ActionHeaderRemove:
					logger.LoggerOasparser.Debugf("Adding %s policy to response flow for %s %s",
						constants.ActionHeaderRemove, resourcePath, operation.GetMethod())
					responseHeaderToRemove, err := generateHeaderToRemoveString(responsePolicy.Parameters)
					if err != nil {
						return nil, fmt.Errorf("error adding response policy %s to operation %s of resource %s."+
							" %v", responsePolicy.Action, operation.GetMethod(), resourcePath, err)
					}
					responseHeadersToRemove = append(responseHeadersToRemove, responseHeaderToRemove)
				}
			}

			// TODO: (suksw) preserve header key case?
			if hasMethodRewritePolicy {
				logger.LoggerOasparser.Debugf("Creating two routes to support method rewrite for %s %s. New method: %s",
					resourcePath, operation.GetMethod(), newMethod)
				match1 := generateRouteMatch(routePath)
				match1.Headers = generateHTTPMethodMatcher(operation.GetMethod(), clusterName)
				match2 := generateRouteMatch(routePath)
				match2.Headers = generateHTTPMethodMatcher(newMethod, clusterName)

				//- external routes only accept requests if metadata "method-rewrite" is null
				//- external routes adds the metadata "method-rewrite"
				//- internal routes only accept requests if metadata "method-rewrite" matches
				//  metadataValue <old_method>_to_<new_method>
				match1.DynamicMetadata = generateMetadataMatcherForExternalRoutes()
				metadataValue := operation.GetMethod() + "_to_" + newMethod
				match2.DynamicMetadata = generateMetadataMatcherForInternalRoutes(metadataValue)

				action1 := generateRouteAction(apiType, routeConfig)
				action2 := generateRouteAction(apiType, routeConfig)

				// Create route1 for current method.
				// Do not add policies to route config. Send via enforcer
				route1 := generateRouteConfig(xWso2Basepath+operation.GetMethod(), match1, action1, nil, decorator, perRouteFilterConfigs,
					nil, nil, nil, nil)

				// Create route2 for new method.
				// Add all policies to route config. Do not send via enforcer.
				if pathRewriteConfig != nil {
					action2.Route.RegexRewrite = pathRewriteConfig
				} else {
					action2.Route.RegexRewrite = generateRegexMatchAndSubstitute(routePath, endpointBasepath, resourcePath, pathMatchType)
				}
				configToSkipEnforcer := generateFilterConfigToSkipEnforcer()
				route2 := generateRouteConfig(xWso2Basepath, match2, action2, nil, decorator, configToSkipEnforcer,
					requestHeadersToAdd, requestHeadersToRemove, responseHeadersToAdd, responseHeadersToRemove)

				routes = append(routes, route1)
				routes = append(routes, route2)
			} else {
				logger.LoggerOasparser.Debug("Creating routes for resource with policies", resourcePath, operation.GetMethod())
				// create route for current method. Add policies to route config. Send via enforcer
				action := generateRouteAction(apiType, routeConfig)
				match := generateRouteMatch(routePath)
				match.Headers = generateHTTPMethodMatcher(operation.GetMethod(), clusterName)
				match.DynamicMetadata = generateMetadataMatcherForExternalRoutes()
				if pathRewriteConfig != nil {
					action.Route.RegexRewrite = pathRewriteConfig
				} else {
					action.Route.RegexRewrite = generateRegexMatchAndSubstitute(routePath, endpointBasepath, resourcePath, pathMatchType)
				}
				route := generateRouteConfig(xWso2Basepath, match, action, nil, decorator, perRouteFilterConfigs,
					requestHeadersToAdd, requestHeadersToRemove, responseHeadersToAdd, responseHeadersToRemove)
				routes = append(routes, route)
			}

		}
	} else {
		logger.LoggerOasparser.Debugf("Creating routes for resource : %s that has no policies", resourcePath)
		// No policies defined for the resource. Therefore, create one route for all operations.
		methodRegex := strings.Join(resourceMethods, "|")
		if !strings.Contains(methodRegex, "OPTIONS") {
			methodRegex = methodRegex + "|OPTIONS"
		}
		match := generateRouteMatch(routePath)
		match.Headers = generateHTTPMethodMatcher(methodRegex, clusterName)
		action := generateRouteAction(apiType, routeConfig)
		rewritePath := generateRoutePathForReWrite(basePath, resourcePath, pathMatchType)
		action.Route.RegexRewrite = generateRegexMatchAndSubstitute(rewritePath, endpointBasepath, resourcePath, pathMatchType)

		route := generateRouteConfig(xWso2Basepath, match, action, nil, decorator, perRouteFilterConfigs,
			nil, nil, nil, nil) // general headers to add and remove are included in this methods
		routes = append(routes, route)
	}
	return routes, nil
}

func getInlineLuaScript(requestInterceptor map[string]model.InterceptEndpoint, responseInterceptor map[string]model.InterceptEndpoint,
	requestContext *interceptor.InvocationContext) string {

	i := &interceptor.Interceptor{
		Context:      requestContext,
		RequestFlow:  make(map[string]interceptor.Config),
		ResponseFlow: make(map[string]interceptor.Config),
	}
	if len(requestInterceptor) > 0 {
		i.IsRequestFlowEnabled = true
		for method, op := range requestInterceptor {
			i.RequestFlow[method] = interceptor.Config{
				ExternalCall: &interceptor.HTTPCallConfig{
					ClusterName: op.ClusterName,
					// multiplying in seconds here because in configs we are directly getting config to time.Duration
					// which is in nano seconds, so multiplying it in seconds here
					Timeout:         strconv.FormatInt((op.RequestTimeout * time.Second).Milliseconds(), 10),
					AuthorityHeader: op.EndpointCluster.Endpoints[0].GetAuthorityHeader(),
				},
				Include: op.Includes,
			}
		}
	}
	if len(responseInterceptor) > 0 {
		i.IsResponseFlowEnabled = true
		for method, op := range responseInterceptor {
			i.ResponseFlow[method] = interceptor.Config{
				ExternalCall: &interceptor.HTTPCallConfig{
					ClusterName: op.ClusterName,
					// multiplying in seconds here because in configs we are directly getting config to time.Duration
					// which is in nano seconds, so multiplying it in seconds here
					Timeout:         strconv.FormatInt((op.RequestTimeout * time.Second).Milliseconds(), 10),
					AuthorityHeader: op.EndpointCluster.Endpoints[0].GetAuthorityHeader(),
				},
				Include: op.Includes,
			}
		}
	}
	templateValues := CombinedTemplateValues{
		WireLogValues{
			LogConfig: config.GetWireLogConfig(),
		},
		*i,
	}

	templateString := interceptor.GetTemplate(i.IsRequestFlowEnabled,
		i.IsResponseFlowEnabled)

	return interceptor.GetInterceptor(templateValues, templateString)
}

// CreateTokenRoute generates a route for the jwt /testkey endpoint
func CreateTokenRoute() *routev3.Route {
	var (
		router    routev3.Route
		action    *routev3.Route_Route
		match     *routev3.RouteMatch
		decorator *routev3.Decorator
	)

	match = &routev3.RouteMatch{
		PathSpecifier: &routev3.RouteMatch_Path{
			Path: testKeyPath,
		},
	}

	hostRewriteSpecifier := &routev3.RouteAction_AutoHostRewrite{
		AutoHostRewrite: &wrapperspb.BoolValue{
			Value: true,
		},
	}

	decorator = &routev3.Decorator{
		Operation: testKeyPath,
	}

	perFilterConfig := extAuthService.ExtAuthzPerRoute{
		Override: &extAuthService.ExtAuthzPerRoute_Disabled{
			Disabled: true,
		},
	}

	b := proto.NewBuffer(nil)
	b.SetDeterministic(true)
	_ = b.Marshal(&perFilterConfig)
	filter := &any.Any{
		TypeUrl: extAuthzPerRouteName,
		Value:   b.Bytes(),
	}

	action = &routev3.Route_Route{
		Route: &routev3.RouteAction{
			HostRewriteSpecifier: hostRewriteSpecifier,
			RegexRewrite: &envoy_type_matcherv3.RegexMatchAndSubstitute{
				Pattern: &envoy_type_matcherv3.RegexMatcher{
					Regex: testKeyPath,
				},
				Substitution: "/",
			},
		},
	}

	directClusterSpecifier := &routev3.RouteAction_Cluster{
		Cluster: "token_cluster",
	}
	action.Route.ClusterSpecifier = directClusterSpecifier

	router = routev3.Route{
		Name:      testKeyPath, //Categorize routes with same base path
		Match:     match,
		Action:    action,
		Metadata:  nil,
		Decorator: decorator,
		TypedPerFilterConfig: map[string]*any.Any{
			wellknown.HTTPExternalAuthorization: filter,
		},
	}
	return &router
}

// CreateHealthEndpoint generates a route for the jwt /health endpoint
// Replies with direct response.
func CreateHealthEndpoint() *routev3.Route {
	var (
		router    routev3.Route
		match     *routev3.RouteMatch
		decorator *routev3.Decorator
	)

	match = &routev3.RouteMatch{
		PathSpecifier: &routev3.RouteMatch_Path{
			Path: healthPath,
		},
	}

	decorator = &routev3.Decorator{
		Operation: healthPath,
	}

	perFilterConfig := extAuthService.ExtAuthzPerRoute{
		Override: &extAuthService.ExtAuthzPerRoute_Disabled{
			Disabled: true,
		},
	}

	b := proto.NewBuffer(nil)
	b.SetDeterministic(true)
	_ = b.Marshal(&perFilterConfig)
	filter := &any.Any{
		TypeUrl: extAuthzPerRouteName,
		Value:   b.Bytes(),
	}

	router = routev3.Route{
		Name:  healthPath, //Categorize routes with same base path
		Match: match,
		Action: &routev3.Route_DirectResponse{
			DirectResponse: &routev3.DirectResponseAction{
				Status: 200,
				Body: &corev3.DataSource{
					Specifier: &corev3.DataSource_InlineString{
						InlineString: healthEndpointResponse,
					},
				},
			},
		},
		Metadata:  nil,
		Decorator: decorator,
		TypedPerFilterConfig: map[string]*any.Any{
			wellknown.HTTPExternalAuthorization: filter,
		},
	}
	return &router
}

// CreateReadyEndpoint generates a route for the router /ready endpoint
// Replies with direct response.
func CreateReadyEndpoint() *routev3.Route {
	var (
		router    routev3.Route
		match     *routev3.RouteMatch
		decorator *routev3.Decorator
	)

	match = &routev3.RouteMatch{
		PathSpecifier: &routev3.RouteMatch_Path{
			Path: readyPath,
		},
	}

	decorator = &routev3.Decorator{
		Operation: readyPath,
	}

	perFilterConfig := extAuthService.ExtAuthzPerRoute{
		Override: &extAuthService.ExtAuthzPerRoute_Disabled{
			Disabled: true,
		},
	}

	b := proto.NewBuffer(nil)
	b.SetDeterministic(true)
	_ = b.Marshal(&perFilterConfig)
	filter := &any.Any{
		TypeUrl: extAuthzPerRouteName,
		Value:   b.Bytes(),
	}

	router = routev3.Route{
		Name:  readyPath, //Categorize routes with same base path
		Match: match,
		Action: &routev3.Route_DirectResponse{
			DirectResponse: &routev3.DirectResponseAction{
				Status: 200,
				Body: &corev3.DataSource{
					Specifier: &corev3.DataSource_InlineString{
						InlineString: readyEndpointResponse,
					},
				},
			},
		},
		Metadata:  nil,
		Decorator: decorator,
		TypedPerFilterConfig: map[string]*any.Any{
			wellknown.HTTPExternalAuthorization: filter,
		},
	}
	return &router
}

// generateRoutePath generates route paths for the api resources.
func generateRoutePath(resourcePath string, pathMatchType gwapiv1b1.PathMatchType) string {
	newPath := strings.TrimSuffix(resourcePath, "/")
	switch pathMatchType {
	case gwapiv1b1.PathMatchExact:
		return fmt.Sprintf("^%s([/]{0,1})", regexp.QuoteMeta(newPath))
	case gwapiv1b1.PathMatchRegularExpression:
		return fmt.Sprintf("^%s([/]{0,1})", newPath)
	case gwapiv1b1.PathMatchPathPrefix:
		fallthrough
	default:
		return fmt.Sprintf("^%s((?:/.*)*)", regexp.QuoteMeta(newPath))
	}
}

// generateRoutePath generates route paths for path rewrite matching.
func generateRoutePathForReWrite(basePath, resourcePath string, pathMatchType gwapiv1b1.PathMatchType) string {
	switch pathMatchType {
	case gwapiv1b1.PathMatchExact:
		fallthrough
	case gwapiv1b1.PathMatchPathPrefix:
		fallthrough
	default:
		return generateRoutePath(resourcePath, pathMatchType)
	case gwapiv1b1.PathMatchRegularExpression:
		return fmt.Sprintf("^(%s)", strings.TrimSuffix(resourcePath, "/"))
	}
}

// generateSubstitutionString returns a regex that has indexes to place the path variables extracted by capture groups
func generateSubstitutionString(resourcePath string, pathMatchType gwapiv1b1.PathMatchType) string {
	var resourceRegex string
	switch pathMatchType {
	case gwapiv1b1.PathMatchExact:
		resourceRegex = resourcePath
	case gwapiv1b1.PathMatchPathPrefix:
		resourceRegex = fmt.Sprintf("%s\\1", strings.TrimSuffix(resourcePath, "/"))
	case gwapiv1b1.PathMatchRegularExpression:
		resourceRegex = "\\1"
	}
	return resourceRegex
}

func isMethodRewrite(resourcePath, method string, policyParams interface{}) (isMethodRewrite bool, err error) {
	var paramsToRewriteMethod map[string]interface{}
	var ok bool
	if paramsToRewriteMethod, ok = policyParams.(map[string]interface{}); !ok {
		return false, fmt.Errorf("error while processing policy parameter map for "+
			"request policy %s to operation %s of resource %s. Map: %v",
			constants.ActionRewriteMethod, method, resourcePath, policyParams)
	}

	currentMethod, exists := paramsToRewriteMethod[constants.CurrentMethod]
	if !exists {
		return true, nil
	}
	currentMethodString, _ := currentMethod.(string)

	if currentMethodString == "<no value>" { // the package text/template return this for keys that does not exist
		return true, nil
	}

	if currentMethodString != method {
		return false, nil
	}
	return true, nil // currentMethodString == method
}

func getRewriteMethod(resourcePath, method string, policyParams interface{}) (rewriteMethod string, err error) {
	var paramsToRewriteMethod map[string]interface{}
	var ok bool
	if paramsToRewriteMethod, ok = policyParams.(map[string]interface{}); !ok {
		return "", fmt.Errorf("error while processing policy parameter map for "+
			"request policy %s to operation %s of resource %s. Map: %v",
			constants.ActionRewriteMethod, method, resourcePath, policyParams)
	}

	updatedMethod, exists := paramsToRewriteMethod[constants.UpdatedMethod]
	if !exists {
		return "", fmt.Errorf("error adding request policy %s to operation %s of resource %s."+
			" Policy parameter updatedMethod not found",
			constants.ActionRewriteMethod, method, resourcePath)
	}
	updatedMethodString, isString := updatedMethod.(string)
	if !isString {
		return "", fmt.Errorf("error adding request policy %s to operation %s of resource %s."+
			" Policy parameter updatedMethod is in incorrect format", constants.ActionRewriteMethod,
			method, resourcePath)
	}
	return updatedMethodString, nil
}

func getUpgradeConfig(apiType string) []*routev3.RouteAction_UpgradeConfig {
	var upgradeConfig []*routev3.RouteAction_UpgradeConfig
	if apiType == constants.WS {
		upgradeConfig = []*routev3.RouteAction_UpgradeConfig{{
			UpgradeType: "websocket",
			Enabled:     &wrappers.BoolValue{Value: true},
		}}
	} else {
		upgradeConfig = []*routev3.RouteAction_UpgradeConfig{{
			UpgradeType: "websocket",
			Enabled:     &wrappers.BoolValue{Value: false},
		}}
	}
	return upgradeConfig
}

func getCorsPolicy(corsConfig *model.CorsConfig) *cors_filter_v3.CorsPolicy {

	if corsConfig == nil || !corsConfig.Enabled {
		return nil
	}

	stringMatcherArray := []*envoy_type_matcherv3.StringMatcher{}
	for _, origin := range corsConfig.AccessControlAllowOrigins {

		// * is considered to be the wild card
		formattedString := regexp.QuoteMeta(origin)
		formattedString = strings.ReplaceAll(formattedString, regexp.QuoteMeta("*"), ".*")

		regexMatcher := &envoy_type_matcherv3.StringMatcher{
			MatchPattern: &envoy_type_matcherv3.StringMatcher_SafeRegex{
				SafeRegex: &envoy_type_matcherv3.RegexMatcher{
					Regex: formattedString,
				},
			},
		}
		stringMatcherArray = append(stringMatcherArray, regexMatcher)
	}

	corsPolicy := &cors_filter_v3.CorsPolicy{
		AllowCredentials: &wrapperspb.BoolValue{
			Value: corsConfig.AccessControlAllowCredentials,
		},
	}

	if len(stringMatcherArray) > 0 {
		corsPolicy.AllowOriginStringMatch = stringMatcherArray
	}
	if len(corsConfig.AccessControlAllowMethods) > 0 {
		corsPolicy.AllowMethods = strings.Join(corsConfig.AccessControlAllowMethods, ", ")
	}
	if len(corsConfig.AccessControlAllowHeaders) > 0 {
		corsPolicy.AllowHeaders = strings.Join(corsConfig.AccessControlAllowHeaders, ", ")
	}
	if len(corsConfig.AccessControlExposeHeaders) > 0 {
		corsPolicy.ExposeHeaders = strings.Join(corsConfig.AccessControlExposeHeaders, ", ")
	}
	return corsPolicy
}

func genRouteCreateParams(swagger *model.MgwSwagger, resource *model.Resource, vHost, endpointBasePath string,
	clusterName string, requestInterceptor map[string]model.InterceptEndpoint,
	responseInterceptor map[string]model.InterceptEndpoint, organizationID string, isSandbox bool) *routeCreateParams {
	params := &routeCreateParams{
		organizationID:               organizationID,
		title:                        swagger.GetTitle(),
		apiType:                      swagger.GetAPIType(),
		version:                      swagger.GetVersion(),
		vHost:                        vHost,
		xWSO2BasePath:                swagger.GetXWso2Basepath(),
		AuthHeader:                   swagger.GetXWSO2AuthHeader(),
		clusterName:                  clusterName,
		endpointBasePath:             endpointBasePath,
		corsPolicy:                   swagger.GetCorsConfig(),
		resource:                     resource,
		requestInterceptor:           requestInterceptor,
		responseInterceptor:          responseInterceptor,
		passRequestPayloadToEnforcer: swagger.GetXWso2RequestBodyPass(),
		isDefaultVersion:             swagger.IsDefaultVersion,
	}
	return params
}

// createAddress generates an address from the given host and port
func createAddress(remoteHost string, port uint32) *corev3.Address {
	address := corev3.Address{Address: &corev3.Address_SocketAddress{
		SocketAddress: &corev3.SocketAddress{
			Address:  remoteHost,
			Protocol: corev3.SocketAddress_TCP,
			PortSpecifier: &corev3.SocketAddress_PortValue{
				PortValue: uint32(port),
			},
		},
	}}
	return &address
}

// getMaxStreamDuration configures a maximum duration for a websocket route.
func getMaxStreamDuration(apiType string) *routev3.RouteAction_MaxStreamDuration {
	var maxStreamDuration *routev3.RouteAction_MaxStreamDuration
	if apiType == constants.WS {
		maxStreamDuration = &routev3.RouteAction_MaxStreamDuration{
			MaxStreamDuration: &durationpb.Duration{
				Seconds: 60 * 60 * 24,
			},
		}
	}
	return maxStreamDuration
}

func getDefaultVersionBasepath(basePath string, version string) string {
	// Following is used to replace only the version when basepath = /foo/v2 and version = v2 and context => /foo/v2/v2
	indexOfVersionString := strings.LastIndex(basePath, "/"+version)
	context := strings.Replace(basePath, "/"+version, "", indexOfVersionString)

	// Having ?: in the regex below, avoids this regex acting as a capturing group. Without this the basepath
	// would again be added in the locations of path variables when sending the request to backend.
	return fmt.Sprintf("(?:%s|%s)", basePath, context)
}

func createInterceptorAPIClusters(mgwSwagger model.MgwSwagger, interceptorCerts map[string][]byte, vHost string, organizationID string) (clustersP []*clusterv3.Cluster,
	addressesP []*corev3.Address, apiRequestInterceptorEndpoint *model.InterceptEndpoint, apiResponseInterceptorEndpoint *model.InterceptEndpoint) {
	var (
		clusters  []*clusterv3.Cluster
		endpoints []*corev3.Address

		apiRequestInterceptor  model.InterceptEndpoint
		apiResponseInterceptor model.InterceptEndpoint
	)
	apiTitle := mgwSwagger.GetTitle()
	apiVersion := mgwSwagger.GetVersion()
	apiRequestInterceptor = mgwSwagger.GetInterceptor(mgwSwagger.GetVendorExtensions(), xWso2requestInterceptor, APILevelInterceptor)
	// if lua filter exists on api level, add cluster
	if apiRequestInterceptor.Enable {
		logger.LoggerOasparser.Debugf("API level request interceptors found for %v : %v", apiTitle, apiVersion)
		apiRequestInterceptor.ClusterName = getClusterName(requestInterceptClustersNamePrefix, organizationID, vHost,
			apiTitle, apiVersion, "")
		cluster, addresses, err := CreateLuaCluster(interceptorCerts, apiRequestInterceptor)
		if err != nil {
			apiRequestInterceptor = model.InterceptEndpoint{}
			logger.LoggerOasparser.Errorf("Error while adding api level request intercepter external cluster for %s. %v",
				apiTitle, err.Error())
		} else {
			clusters = append(clusters, cluster)
			endpoints = append(endpoints, addresses...)
		}
	}
	apiResponseInterceptor = mgwSwagger.GetInterceptor(mgwSwagger.GetVendorExtensions(), xWso2responseInterceptor, APILevelInterceptor)
	// if lua filter exists on api level, add cluster
	if apiResponseInterceptor.Enable {
		logger.LoggerOasparser.Debugln("API level response interceptors found for " + mgwSwagger.GetID())
		apiResponseInterceptor.ClusterName = getClusterName(responseInterceptClustersNamePrefix, organizationID, vHost,
			apiTitle, apiVersion, "")
		cluster, addresses, err := CreateLuaCluster(interceptorCerts, apiResponseInterceptor)
		if err != nil {
			apiResponseInterceptor = model.InterceptEndpoint{}
			logger.LoggerOasparser.Errorf("Error while adding api level response intercepter external cluster for %s. %v", apiTitle, err.Error())
		} else {
			clusters = append(clusters, cluster)
			endpoints = append(endpoints, addresses...)
		}
	}
	return clusters, endpoints, &apiRequestInterceptor, &apiResponseInterceptor
}

func createInterceptorResourceClusters(mgwSwagger model.MgwSwagger, interceptorCerts map[string][]byte, vHost string, organizationID string,
	apiRequestInterceptor *model.InterceptEndpoint, apiResponseInterceptor *model.InterceptEndpoint, resource *model.Resource) (clustersP []*clusterv3.Cluster, addressesP []*corev3.Address,
	operationalReqInterceptorsEndpoint *map[string]model.InterceptEndpoint, operationalRespInterceptorValEndpoint *map[string]model.InterceptEndpoint) {
	var (
		clusters  []*clusterv3.Cluster
		endpoints []*corev3.Address
	)
	resourceRequestInterceptor := apiRequestInterceptor
	resourceResponseInterceptor := apiResponseInterceptor
	apiTitle := mgwSwagger.GetTitle()
	apiVersion := mgwSwagger.GetVersion()
	reqInterceptorVal := mgwSwagger.GetInterceptor(resource.GetVendorExtensions(), xWso2requestInterceptor, ResourceLevelInterceptor)
	if reqInterceptorVal.Enable {
		logger.LoggerOasparser.Debugf("Resource level request interceptors found for %v:%v-%v", apiTitle, apiVersion, resource.GetPath())
		reqInterceptorVal.ClusterName = getClusterName(requestInterceptClustersNamePrefix, organizationID, vHost,
			apiTitle, apiVersion, resource.GetID())
		cluster, addresses, err := CreateLuaCluster(interceptorCerts, reqInterceptorVal)
		if err != nil {
			logger.LoggerOasparser.Errorf("Error while adding resource level request intercept external cluster for %s. %v",
				apiTitle, err.Error())
		} else {
			resourceRequestInterceptor = &reqInterceptorVal
			clusters = append(clusters, cluster)
			endpoints = append(endpoints, addresses...)
		}
	}

	// create operational level response interceptor clusters
	operationalReqInterceptors := mgwSwagger.GetOperationInterceptors(*apiRequestInterceptor, *resourceRequestInterceptor, resource.GetMethod(), true)
	for method, opI := range operationalReqInterceptors {
		if opI.Enable && opI.Level == OperationLevelInterceptor {
			logger.LoggerOasparser.Debugf("Operation level request interceptors found for %v:%v-%v-%v", apiTitle, apiVersion, resource.GetPath(),
				opI.ClusterName)
			opID := opI.ClusterName
			opI.ClusterName = getClusterName(requestInterceptClustersNamePrefix, organizationID, vHost, apiTitle, apiVersion, opID)
			operationalReqInterceptors[method] = opI // since cluster name is updated
			cluster, addresses, err := CreateLuaCluster(interceptorCerts, opI)
			if err != nil {
				logger.LoggerOasparser.Errorf("Error while adding operational level request intercept external cluster for %v:%v-%v-%v. %v",
					apiTitle, apiVersion, resource.GetPath(), opID, err.Error())
				// setting resource level interceptor to failed operation level interceptor.
				operationalReqInterceptors[method] = *resourceRequestInterceptor
			} else {
				clusters = append(clusters, cluster)
				endpoints = append(endpoints, addresses...)
			}
		}
	}

	// create resource level response interceptor cluster
	respInterceptorVal := mgwSwagger.GetInterceptor(resource.GetVendorExtensions(), xWso2responseInterceptor, ResourceLevelInterceptor)
	if respInterceptorVal.Enable {
		logger.LoggerOasparser.Debugf("Resource level response interceptors found for %v:%v-%v"+apiTitle, apiVersion, resource.GetPath())
		respInterceptorVal.ClusterName = getClusterName(responseInterceptClustersNamePrefix, organizationID,
			vHost, apiTitle, apiVersion, resource.GetID())
		cluster, addresses, err := CreateLuaCluster(interceptorCerts, respInterceptorVal)
		if err != nil {
			logger.LoggerOasparser.Errorf("Error while adding resource level response intercept external cluster for %s. %v",
				apiTitle, err.Error())
		} else {
			resourceResponseInterceptor = &respInterceptorVal
			clusters = append(clusters, cluster)
			endpoints = append(endpoints, addresses...)
		}
	}

	// create operation level response interceptor clusters
	operationalRespInterceptorVal := mgwSwagger.GetOperationInterceptors(*apiResponseInterceptor, *resourceResponseInterceptor, resource.GetMethod(),
		false)
	for method, opI := range operationalRespInterceptorVal {
		if opI.Enable && opI.Level == OperationLevelInterceptor {
			logger.LoggerOasparser.Debugf("Operational level response interceptors found for %v:%v-%v-%v", apiTitle, apiVersion, resource.GetPath(),
				opI.ClusterName)
			opID := opI.ClusterName
			opI.ClusterName = getClusterName(responseInterceptClustersNamePrefix, organizationID, vHost, apiTitle, apiVersion, opID)
			operationalRespInterceptorVal[method] = opI // since cluster name is updated
			cluster, addresses, err := CreateLuaCluster(interceptorCerts, opI)
			if err != nil {
				logger.LoggerOasparser.Errorf("Error while adding operational level response intercept external cluster for %v:%v-%v-%v. %v",
					apiTitle, apiVersion, resource.GetPath(), opID, err.Error())
				// setting resource level interceptor to failed operation level interceptor.
				operationalRespInterceptorVal[method] = *resourceResponseInterceptor
			} else {
				clusters = append(clusters, cluster)
				endpoints = append(endpoints, addresses...)
			}
		}
	}
	return clusters, endpoints, &operationalReqInterceptors, &operationalRespInterceptorVal
}
