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
	"strconv"
	"time"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	listenerv3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	routev3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	envoy_config_trace_v3 "github.com/envoyproxy/go-control-plane/envoy/config/trace/v3"
	tlsInspectorv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/listener/tls_inspector/v3"
	hcmv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	tlsv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
	metadatav3 "github.com/envoyproxy/go-control-plane/envoy/type/metadata/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	"github.com/golang/protobuf/ptypes/wrappers"

	"github.com/wso2/apk/adapter/config"
	"github.com/wso2/apk/adapter/internal/discovery/xds/common"
	"github.com/wso2/apk/adapter/internal/loggers"
	"github.com/wso2/apk/adapter/internal/oasparser/model"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
)

// CreateRoutesConfigForRds generates the default RouteConfiguration.
// Only the provided virtual hosts will be assigned inside the configuration.
// This is used to provide the configuration for RDS.
func CreateRoutesConfigForRds(vHosts []*routev3.VirtualHost, httpListeners string) *routev3.RouteConfiguration {
	rdsConfigName := httpListeners
	routeConfiguration := routev3.RouteConfiguration{
		Name:                   rdsConfigName,
		VirtualHosts:           vHosts,
		RequestHeadersToRemove: []string{clusterHeaderName},
	}
	return &routeConfiguration
}

// CreateListenerByGateway create listeners by provided gateway object with the Route Configuration
// stated as RDS. (routes are not assigned directly to the listener.) RouteConfiguration name
// is assigned using its default value. Route Configuration would be resolved via ADS.
//
// If SecuredListenerPort and ListenerPort both are mentioned, two listeners would be added.
// If neither of the two properies are assigned with non-zero values, adapter would panic.
//
// HTTPConnectionManager with HTTP Filters, Accesslog configuration, TransportSocket
// Configuration is included within the implementation.
//
// Listener Address, ListenerPort Value, SecuredListener Address, and  SecuredListenerPort Values are
// fetched from the configuration accordingly.
//
// The relevant private keys and certificates (for securedListener) are fetched from the filepath
// mentioned in the adapter configuration. These certificate, key values are added
// as inline records (base64 encoded).
func CreateListenerByGateway(gateway *gwapiv1.Gateway, resolvedListenerCerts map[string]map[string][]byte, gwLuaScript string) []*listenerv3.Listener {
	conf := config.ReadConfigs()
	// Prepare a map that contains all the listerners identified in all of the gateways that reconciled so far.
	// This map contains port - listeners per protocol with port
	protocolListenerMap := make(map[gwapiv1.ProtocolType]map[uint32][]gwapiv1.Listener)
	for _, listener := range gateway.Spec.Listeners {
		port := uint32(listener.Port)
		if protocolListenerMap[listener.Protocol] == nil {
			protocolListenerMap[listener.Protocol] = make(map[uint32][]gwapiv1.Listener)
		}
		if protocolListenerMap[listener.Protocol][port] == nil {
			protocolListenerMap[listener.Protocol][port] = []gwapiv1.Listener{}
		}
		protocolListenerMap[listener.Protocol][port] = append(protocolListenerMap[listener.Protocol][port], listener)
	}
	loggers.LoggerAPKOperator.Infof("CreateListenerByGateway is called. ProtocolListenerMap: %+v", protocolListenerMap)
	listenerList := make([]*listenerv3.Listener, 0)
	for protocol, protocolPort := range protocolListenerMap {
		for port, listeners := range protocolPort {
			var httpFilters []*hcmv3.HttpFilter
			upgradeFilters := getUpgradeFilters()
			accessLogs := getAccessLogs()
			var listener *listenerv3.Listener
			var listenerName string
			var filterChains []*listenerv3.FilterChain

			for _, listenerObj := range listeners {
				if listenerObj.Name == "httpslistener" {
					httpFilters = getHTTPFilters(gwLuaScript)
				} else {
					httpFilters = getHTTPFilters(`
		function envoy_on_request(request_handle)
		end
		function envoy_on_response(response_handle)
		end`)
				}
				listenerName = common.GetEnvoyListenerName(string(protocol), port)
				filterChainMatch := &listenerv3.FilterChainMatch{
					ServerNames: []string{string(*listenerObj.Hostname)},
				}
				var transportSocket *corev3.TransportSocket
				if protocol == gwapiv1.HTTPSProtocolType {
					publicCertData := resolvedListenerCerts[string(listenerObj.Name)]["tls.crt"]
					privateKeyData := resolvedListenerCerts[string(listenerObj.Name)]["tls.key"]
					var tlsFilter *tlsv3.DownstreamTlsContext
					tlsCert := generateTLSCertWithStr(string(privateKeyData), string(publicCertData))
					//TODO: Make this configurable using config map from listener object
					if conf.Envoy.Downstream.TLS.MTLSAPIsEnabled {
						tlsFilter = &tlsv3.DownstreamTlsContext{
							// This is false since the authentication will be done at the enforcer
							RequireClientCertificate: &wrappers.BoolValue{
								Value: false,
							},
							CommonTlsContext: &tlsv3.CommonTlsContext{
								//TlsCertificateSdsSecretConfigs
								TlsCertificates: []*tlsv3.TlsCertificate{tlsCert},
								//For the purpose of including peer certificate into the request context
								ValidationContextType: &tlsv3.CommonTlsContext_ValidationContext{
									ValidationContext: &tlsv3.CertificateValidationContext{
										TrustedCa: &corev3.DataSource{
											Specifier: &corev3.DataSource_Filename{
												Filename: conf.Envoy.Downstream.TLS.TrustedCertPath,
											},
										},
									},
								},
							},
						}
					} else {
						tlsFilter = &tlsv3.DownstreamTlsContext{
							CommonTlsContext: &tlsv3.CommonTlsContext{
								//TlsCertificateSdsSecretConfigs
								TlsCertificates: []*tlsv3.TlsCertificate{tlsCert},
								AlpnProtocols:   []string{"h2", "http/1.1"},
							},
						}
					}

					marshalledTLSFilter, err := anypb.New(tlsFilter)
					if err != nil {
						loggers.LoggerOasparser.Fatal("Error while Marshalling the downstream TLS Context for the configuration.")
					}

					transportSocket = &corev3.TransportSocket{
						Name: wellknown.TransportSocketTLS,
						ConfigType: &corev3.TransportSocket_TypedConfig{
							TypedConfig: marshalledTLSFilter,
						},
					}
				}

				var filters []*listenerv3.Filter
				manager := &hcmv3.HttpConnectionManager{
					CodecType:  getListenerCodecType(conf.Envoy.ListenerCodecType),
					StatPrefix: httpConManagerStartPrefix,
					// WebSocket upgrades enabled from the HCM
					UpgradeConfigs: []*hcmv3.HttpConnectionManager_UpgradeConfig{{
						UpgradeType: "websocket",
						Enabled:     &wrappers.BoolValue{Value: true},
						Filters:     upgradeFilters,
					}},
					RouteSpecifier: &hcmv3.HttpConnectionManager_Rds{
						Rds: &hcmv3.Rds{
							RouteConfigName: common.GetEnvoyRouteConfigName(listenerName, string(listenerObj.Name)),
							ConfigSource: &corev3.ConfigSource{
								ConfigSourceSpecifier: &corev3.ConfigSource_Ads{
									Ads: &corev3.AggregatedConfigSource{},
								},
								ResourceApiVersion: corev3.ApiVersion_V3,
							},
						},
					},
					HttpFilters: httpFilters,
					LocalReplyConfig: &hcmv3.LocalReplyConfig{
						Mappers: getErrorResponseMappers(),
					},
					RequestTimeout:        durationpb.New(conf.Envoy.Connection.Timeouts.RequestTimeoutInSeconds * time.Second),        // default disabled
					RequestHeadersTimeout: durationpb.New(conf.Envoy.Connection.Timeouts.RequestHeadersTimeoutInSeconds * time.Second), // default disabled
					StreamIdleTimeout:     durationpb.New(conf.Envoy.Connection.Timeouts.StreamIdleTimeoutInSeconds * time.Second),     // Default 5 mins
					CommonHttpProtocolOptions: &corev3.HttpProtocolOptions{
						IdleTimeout: durationpb.New(conf.Envoy.Connection.Timeouts.IdleTimeoutInSeconds * time.Second), // Default 1 hr
					},
					HttpProtocolOptions: &corev3.Http1ProtocolOptions{
						EnableTrailers: config.GetWireLogConfig().LogTrailersEnabled,
					},
					UseRemoteAddress:     &wrappers.BoolValue{Value: conf.Envoy.UseRemoteAddress},
					AppendXForwardedPort: true,
				}

				if len(accessLogs) > 0 {
					manager.AccessLog = accessLogs
				}

				if conf.Tracing.Enabled {
					if conf.Tracing.Type == TracerTypeOtlp {
						if tracing, err := getTracingOTLP(conf); err == nil {
							manager.Tracing = tracing
							manager.GenerateRequestId = &wrappers.BoolValue{Value: conf.Tracing.Enabled}
						} else {
							loggers.LoggerOasparser.Errorf("Failed to initialize tracing for %s. Router tracing will be disabled. Error: %s",
								TracerTypeOtlp, err)
							conf.Tracing.Enabled = false
						}
					} else if conf.Tracing.Type != TracerTypeAzure {
						if tracing, err := getZipkinTracing(conf); err == nil {
							manager.Tracing = tracing
							manager.GenerateRequestId = &wrappers.BoolValue{Value: conf.Tracing.Enabled}
						} else {
							loggers.LoggerOasparser.Errorf("Failed to initialize tracing for %s. Router tracing will be disabled. Error: %s",
								conf.Tracing.Type, err)
							conf.Tracing.Enabled = false
						}
					}
				}

				pbst, err := anypb.New(manager)
				if err != nil {
					loggers.LoggerOasparser.Fatal(err)
				}
				connectionManagerFilterP := listenerv3.Filter{
					Name: wellknown.HTTPConnectionManager,
					ConfigType: &listenerv3.Filter_TypedConfig{
						TypedConfig: pbst,
					},
				}

				// add filters
				filters = append(filters, &connectionManagerFilterP)
				if protocol == gwapiv1.HTTPSProtocolType {
					filterChains = append(filterChains, &listenerv3.FilterChain{
						FilterChainMatch: filterChainMatch,
						Filters:          filters,
						TransportSocket:  transportSocket,
					})
				} else {
					filterChains = append(filterChains, &listenerv3.FilterChain{
						Filters: filters,
					})
				}

			}

			if protocol == gwapiv1.HTTPSProtocolType {
				listenerHostAddress := defaultListenerHostAddress
				securedListenerAddress := &corev3.Address_SocketAddress{
					SocketAddress: &corev3.SocketAddress{
						Protocol: corev3.SocketAddress_TCP,
						Address:  listenerHostAddress,
						PortSpecifier: &corev3.SocketAddress_PortValue{
							PortValue: port,
						},
					},
				}

				//var tlsInspector *tlsInspectorv3.TlsInspector

				tlsInspector := &tlsInspectorv3.TlsInspector{}
				marshalledListenerFilter, err := anypb.New(tlsInspector)
				if err != nil {
					loggers.LoggerOasparser.Fatal("Error while Marshalling the TlsInspector for the configuration.")
				}

				listenerFilters := []*listenerv3.ListenerFilter{
					{ // TLS Inspector
						Name: wellknown.TlsInspector,
						ConfigType: &listenerv3.ListenerFilter_TypedConfig{
							TypedConfig: marshalledListenerFilter,
						},
					},
				}

				listener = &listenerv3.Listener{
					Name: string(listenerName),
					Address: &corev3.Address{
						Address: securedListenerAddress,
					},
					ListenerFilters: listenerFilters,
					FilterChains:    filterChains,
				}
				loggers.LoggerOasparser.Infof("Secured Listener is added. %s : %d", listenerHostAddress, port)
			} else {
				loggers.LoggerOasparser.Info("No SecuredListenerPort is included.")
			}

			if protocol == gwapiv1.HTTPProtocolType {
				listenerHostAddress := defaultListenerHostAddress
				listenerAddress := &corev3.Address_SocketAddress{
					SocketAddress: &corev3.SocketAddress{
						Protocol: corev3.SocketAddress_TCP,
						Address:  listenerHostAddress,
						PortSpecifier: &corev3.SocketAddress_PortValue{
							PortValue: port,
						},
					},
				}

				listener = &listenerv3.Listener{
					Name: string(listenerName),
					Address: &corev3.Address{
						Address: listenerAddress,
					},
					FilterChains: filterChains,
				}
				loggers.LoggerOasparser.Infof("Non-secured Listener is added. %s : %d", listenerHostAddress, port)
			} else {
				loggers.LoggerOasparser.Info("No Non-securedListenerPort is included.")
			}

			if listeners == nil {
				err := errors.New("no Listeners are configured as no port value is mentioned under securedListenerPort or ListenerPort")
				loggers.LoggerOasparser.Fatal(err)
			}
			listenerList = append(listenerList, listener)
		}
	}
	loggers.LoggerOasparser.Infof("Listener list size. %+v", len(listenerList))
	return listenerList
}

// CreateVirtualHosts creates VirtualHost configurations for envoy which serves
// request from the vHost domain. The routes array will be included as the routes
// for the created virtual host.
func CreateVirtualHosts(vhostToRouteArrayMap map[string][]*routev3.Route, customRateLimitPolicies []*model.CustomRateLimitPolicy, vhostToSubscriptionAIRL map[string]bool, vhostToSubscriptionRL map[string]bool) []*routev3.VirtualHost {
	virtualHosts := make([]*routev3.VirtualHost, 0, len(vhostToRouteArrayMap))
	var rateLimits []*routev3.RateLimit
	for _, customRateLimitPolicy := range customRateLimitPolicies {
		var actions []*routev3.RateLimit_Action
		actions = append(actions, &routev3.RateLimit_Action{
			ActionSpecifier: &routev3.RateLimit_Action_Metadata{
				Metadata: &routev3.RateLimit_Action_MetaData{
					DescriptorKey: OrgMetadataKey,
					MetadataKey: &metadatav3.MetadataKey{
						Key: MetadataNamespaceForWSO2Policies,
						Path: []*metadatav3.MetadataKey_PathSegment{
							{
								Segment: &metadatav3.MetadataKey_PathSegment_Key{
									Key: OrgMetadataKey,
								},
							},
						},
					},
					Source: routev3.RateLimit_Action_MetaData_DYNAMIC,
				},
			},
		})
		actions = append(actions, &routev3.RateLimit_Action{
			ActionSpecifier: &routev3.RateLimit_Action_Metadata{
				Metadata: &routev3.RateLimit_Action_MetaData{
					DescriptorKey: customRateLimitPolicy.Key,
					MetadataKey: &metadatav3.MetadataKey{
						Key: MetadataNamespaceForCustomPolicies,
						Path: []*metadatav3.MetadataKey_PathSegment{
							{
								Segment: &metadatav3.MetadataKey_PathSegment_Key{
									Key: customRateLimitPolicy.Key,
								},
							},
						},
					},
					Source: routev3.RateLimit_Action_MetaData_DYNAMIC,
				},
			},
		})
		rateLimits = append(rateLimits, &routev3.RateLimit{
			Actions: actions,
		})
	}

	for vhost, routes := range vhostToRouteArrayMap {
		if flag, exists := vhostToSubscriptionAIRL[vhost]; exists && flag {
			rateLimits = append(rateLimits, generateSubscriptionBasedAIRatelimits()...)
		}
		if flag, exists := vhostToSubscriptionRL[vhost]; exists && flag {
			rateLimits = append(rateLimits, generateSubscriptionBasedRatelimits()...)
		}
		virtualHost := &routev3.VirtualHost{
			Name:       vhost,
			Domains:    []string{vhost, fmt.Sprint(vhost, ":*")},
			Routes:     routes,
			RateLimits: rateLimits,
		}
		virtualHosts = append(virtualHosts, virtualHost)
	}
	return virtualHosts
}

func generateSubscriptionBasedRatelimits() []*routev3.RateLimit {
	return []*routev3.RateLimit{&routev3.RateLimit{
		Actions: []*routev3.RateLimit_Action{
			{
				ActionSpecifier: &routev3.RateLimit_Action_Metadata{
					Metadata: &routev3.RateLimit_Action_MetaData{
						DescriptorKey: DescriptorKeyForOrganization,
						MetadataKey: &metadatav3.MetadataKey{
							Key: extAuthzFilterName,
							Path: []*metadatav3.MetadataKey_PathSegment{
								{
									Segment: &metadatav3.MetadataKey_PathSegment_Key{
										Key: descriptorMetadataKeyForOrganization,
									},
								},
							},
						},
						Source:       routev3.RateLimit_Action_MetaData_DYNAMIC,
						SkipIfAbsent: true,
					},
				},
			},
			{
				ActionSpecifier: &routev3.RateLimit_Action_Metadata{
					Metadata: &routev3.RateLimit_Action_MetaData{
						DescriptorKey: DescriptorKeyForSubscription,
						MetadataKey: &metadatav3.MetadataKey{
							Key: extAuthzFilterName,
							Path: []*metadatav3.MetadataKey_PathSegment{
								{
									Segment: &metadatav3.MetadataKey_PathSegment_Key{
										Key: descriptorMetadataKeyForSubscription,
									},
								},
							},
						},
						Source:       routev3.RateLimit_Action_MetaData_DYNAMIC,
						SkipIfAbsent: true,
					},
				},
			},
			{
				ActionSpecifier: &routev3.RateLimit_Action_Metadata{
					Metadata: &routev3.RateLimit_Action_MetaData{
						DescriptorKey: DescriptorKeyForPolicy,
						MetadataKey: &metadatav3.MetadataKey{
							Key: extAuthzFilterName,
							Path: []*metadatav3.MetadataKey_PathSegment{
								{
									Segment: &metadatav3.MetadataKey_PathSegment_Key{
										Key: descriptorMetadataKeyForUsagePolicy,
									},
								},
							},
						},
						Source:       routev3.RateLimit_Action_MetaData_DYNAMIC,
						SkipIfAbsent: true,
					},
				},
			},
		},
	}, &routev3.RateLimit{
		Actions: []*routev3.RateLimit_Action{
			{
				ActionSpecifier: &routev3.RateLimit_Action_Metadata{
					Metadata: &routev3.RateLimit_Action_MetaData{
						DescriptorKey: DescriptorKeyForOrganization,
						MetadataKey: &metadatav3.MetadataKey{
							Key: extAuthzFilterName,
							Path: []*metadatav3.MetadataKey_PathSegment{
								{
									Segment: &metadatav3.MetadataKey_PathSegment_Key{
										Key: descriptorMetadataKeyForOrganization,
									},
								},
							},
						},
						Source:       routev3.RateLimit_Action_MetaData_DYNAMIC,
						SkipIfAbsent: true,
					},
				},
			},
			{
				ActionSpecifier: &routev3.RateLimit_Action_Metadata{
					Metadata: &routev3.RateLimit_Action_MetaData{
						DescriptorKey: DescriptorKeyForSubscription,
						MetadataKey: &metadatav3.MetadataKey{
							Key: extAuthzFilterName,
							Path: []*metadatav3.MetadataKey_PathSegment{
								{
									Segment: &metadatav3.MetadataKey_PathSegment_Key{
										Key: descriptorMetadataKeyForSubscription,
									},
								},
							},
						},
						Source:       routev3.RateLimit_Action_MetaData_DYNAMIC,
						SkipIfAbsent: true,
					},
				},
			},
			{
				ActionSpecifier: &routev3.RateLimit_Action_Metadata{
					Metadata: &routev3.RateLimit_Action_MetaData{
						DescriptorKey: DescriptorKeyForPolicy,
						MetadataKey: &metadatav3.MetadataKey{
							Key: extAuthzFilterName,
							Path: []*metadatav3.MetadataKey_PathSegment{
								{
									Segment: &metadatav3.MetadataKey_PathSegment_Key{
										Key: descriptorMetadataKeyForUsagePolicy,
									},
								},
							},
						},
						Source:       routev3.RateLimit_Action_MetaData_DYNAMIC,
						SkipIfAbsent: true,
					},
				},
			},
			{
				ActionSpecifier: &routev3.RateLimit_Action_GenericKey_{
					GenericKey: &routev3.RateLimit_Action_GenericKey{
						DescriptorKey:   "burst",
						DescriptorValue: "enabled",
					},
				},
			},
		},
	},
	}
}

func generateSubscriptionBasedAIRatelimits() []*routev3.RateLimit {
	rateLimitForRequestTokenCount := routev3.RateLimit{
		Actions: []*routev3.RateLimit_Action{
			{
				ActionSpecifier: &routev3.RateLimit_Action_Metadata{
					Metadata: &routev3.RateLimit_Action_MetaData{
						DescriptorKey: DescriptorKeyForAIRequestTokenCountForSubscriptionBasedAIRL,
						MetadataKey: &metadatav3.MetadataKey{
							Key: extAuthzFilterName,
							Path: []*metadatav3.MetadataKey_PathSegment{
								&metadatav3.MetadataKey_PathSegment{
									Segment: &metadatav3.MetadataKey_PathSegment_Key{
										Key: DynamicMetadataKeyForOrganizationAndAIRLPolicy,
									},
								},
							},
						},
						Source:       routev3.RateLimit_Action_MetaData_DYNAMIC,
						SkipIfAbsent: true,
					},
				},
			},
			{
				ActionSpecifier: &routev3.RateLimit_Action_Metadata{
					Metadata: &routev3.RateLimit_Action_MetaData{
						DescriptorKey: DescriptorKeyForAISubscription,
						MetadataKey: &metadatav3.MetadataKey{
							Key: extAuthzFilterName,
							Path: []*metadatav3.MetadataKey_PathSegment{
								&metadatav3.MetadataKey_PathSegment{
									Segment: &metadatav3.MetadataKey_PathSegment_Key{
										Key: DynamicMetadataKeyForSubscription,
									},
								},
							},
						},
						Source:       routev3.RateLimit_Action_MetaData_DYNAMIC,
						SkipIfAbsent: true,
					},
				},
			},
		},
	}
	rateLimitForResponseTokenCount := routev3.RateLimit{
		Actions: []*routev3.RateLimit_Action{
			{
				ActionSpecifier: &routev3.RateLimit_Action_Metadata{
					Metadata: &routev3.RateLimit_Action_MetaData{
						DescriptorKey: DescriptorKeyForAIResponseTokenCountForSubscriptionBasedAIRL,
						MetadataKey: &metadatav3.MetadataKey{
							Key: extAuthzFilterName,
							Path: []*metadatav3.MetadataKey_PathSegment{
								&metadatav3.MetadataKey_PathSegment{
									Segment: &metadatav3.MetadataKey_PathSegment_Key{
										Key: DynamicMetadataKeyForOrganizationAndAIRLPolicy,
									},
								},
							},
						},
						Source:       routev3.RateLimit_Action_MetaData_DYNAMIC,
						SkipIfAbsent: true,
					},
				},
			},
			{
				ActionSpecifier: &routev3.RateLimit_Action_Metadata{
					Metadata: &routev3.RateLimit_Action_MetaData{
						DescriptorKey: DescriptorKeyForAISubscription,
						MetadataKey: &metadatav3.MetadataKey{
							Key: extAuthzFilterName,
							Path: []*metadatav3.MetadataKey_PathSegment{
								&metadatav3.MetadataKey_PathSegment{
									Segment: &metadatav3.MetadataKey_PathSegment_Key{
										Key: DynamicMetadataKeyForSubscription,
									},
								},
							},
						},
						Source:       routev3.RateLimit_Action_MetaData_DYNAMIC,
						SkipIfAbsent: true,
					},
				},
			},
		},
	}
	rateLimitForRequestCount := routev3.RateLimit{
		Actions: []*routev3.RateLimit_Action{
			{
				ActionSpecifier: &routev3.RateLimit_Action_Metadata{
					Metadata: &routev3.RateLimit_Action_MetaData{
						DescriptorKey: DescriptorKeyForAIRequestCountForSubscriptionBasedAIRL,
						MetadataKey: &metadatav3.MetadataKey{
							Key: extAuthzFilterName,
							Path: []*metadatav3.MetadataKey_PathSegment{
								&metadatav3.MetadataKey_PathSegment{
									Segment: &metadatav3.MetadataKey_PathSegment_Key{
										Key: DynamicMetadataKeyForOrganizationAndAIRLPolicy,
									},
								},
							},
						},
						Source:       routev3.RateLimit_Action_MetaData_DYNAMIC,
						SkipIfAbsent: true,
					},
				},
			},
			{
				ActionSpecifier: &routev3.RateLimit_Action_Metadata{
					Metadata: &routev3.RateLimit_Action_MetaData{
						DescriptorKey: DescriptorKeyForAISubscription,
						MetadataKey: &metadatav3.MetadataKey{
							Key: extAuthzFilterName,
							Path: []*metadatav3.MetadataKey_PathSegment{
								&metadatav3.MetadataKey_PathSegment{
									Segment: &metadatav3.MetadataKey_PathSegment_Key{
										Key: DynamicMetadataKeyForSubscription,
									},
								},
							},
						},
						Source:       routev3.RateLimit_Action_MetaData_DYNAMIC,
						SkipIfAbsent: true,
					},
				},
			},
		},
	}
	rateLimitForTotalTokenCount := routev3.RateLimit{
		Actions: []*routev3.RateLimit_Action{
			{
				ActionSpecifier: &routev3.RateLimit_Action_Metadata{
					Metadata: &routev3.RateLimit_Action_MetaData{
						DescriptorKey: DescriptorKeyForAITotalTokenCountForSubscriptionBasedAIRL,
						MetadataKey: &metadatav3.MetadataKey{
							Key: extAuthzFilterName,
							Path: []*metadatav3.MetadataKey_PathSegment{
								&metadatav3.MetadataKey_PathSegment{
									Segment: &metadatav3.MetadataKey_PathSegment_Key{
										Key: DynamicMetadataKeyForOrganizationAndAIRLPolicy,
									},
								},
							},
						},
						Source:       routev3.RateLimit_Action_MetaData_DYNAMIC,
						SkipIfAbsent: true,
					},
				},
			},
			{
				ActionSpecifier: &routev3.RateLimit_Action_Metadata{
					Metadata: &routev3.RateLimit_Action_MetaData{
						DescriptorKey: DescriptorKeyForAISubscription,
						MetadataKey: &metadatav3.MetadataKey{
							Key: extAuthzFilterName,
							Path: []*metadatav3.MetadataKey_PathSegment{
								&metadatav3.MetadataKey_PathSegment{
									Segment: &metadatav3.MetadataKey_PathSegment_Key{
										Key: DynamicMetadataKeyForSubscription,
									},
								},
							},
						},
						Source:       routev3.RateLimit_Action_MetaData_DYNAMIC,
						SkipIfAbsent: true,
					},
				},
			},
		},
	}
	return []*routev3.RateLimit{&rateLimitForRequestTokenCount, &rateLimitForResponseTokenCount, &rateLimitForRequestCount, &rateLimitForTotalTokenCount}
}

// TODO: (VirajSalaka) Still the following method is not utilized as Sds is not implement. Keeping the Implementation for future reference
// func generateDefaultSdsSecretFromConfigfile(privateKeyPath string, pulicKeyPath string) (*tlsv3.Secret, error) {
// 	var secret tlsv3.Secret
// 	tlsCert := generateTLSCert(privateKeyPath, pulicKeyPath)
// 	secret = tlsv3.Secret{
// 		Name: defaultListenerSecretConfigName,
// 		Type: &tlsv3.Secret_TlsCertificate{
// 			TlsCertificate: tlsCert,
// 		},
// 	}
// 	return &secret, nil
// }

// generateTLSCert generates the TLS Certiificate with given private key filepath and the corresponding public Key filepath.
// The files should be mounted to the router container unless the default cert is used.
func generateTLSCert(privateKeyPath string, publicKeyPath string) *tlsv3.TlsCertificate {
	tlsCert := tlsv3.TlsCertificate{
		PrivateKey: &corev3.DataSource{
			Specifier: &corev3.DataSource_Filename{
				Filename: privateKeyPath,
			},
		},
		CertificateChain: &corev3.DataSource{
			Specifier: &corev3.DataSource_Filename{
				Filename: publicKeyPath,
			},
		},
	}
	return &tlsCert
}

// generate TLS certs as inline strings
func generateTLSCertWithStr(privateKey string, publicKey string) *tlsv3.TlsCertificate {
	tlsCert := tlsv3.TlsCertificate{
		PrivateKey: &corev3.DataSource{
			Specifier: &corev3.DataSource_InlineString{
				InlineString: privateKey,
			},
		},
		CertificateChain: &corev3.DataSource{
			Specifier: &corev3.DataSource_InlineString{
				InlineString: publicKey,
			},
		},
	}
	return &tlsCert
}

func getZipkinTracing(conf *config.Config) (*hcmv3.HttpConnectionManager_Tracing, error) {
	var endpoint string
	var maxPathLength uint32

	if endpoint = conf.Tracing.ConfigProperties[tracerEndpoint]; len(endpoint) <= 0 {
		return nil, errors.New("invalid endpoint path provided for tracing endpoint")
	}
	if length, err := strconv.ParseUint(conf.Tracing.ConfigProperties[tracerMaxPathLength], 10, 32); err == nil {
		maxPathLength = uint32(length)
	} else {
		return nil, errors.New("invalid max path length provided for tracing endpoint")
	}

	providerConf := &envoy_config_trace_v3.ZipkinConfig{
		CollectorCluster:         tracingClusterName,
		CollectorEndpoint:        endpoint,
		CollectorEndpointVersion: envoy_config_trace_v3.ZipkinConfig_HTTP_JSON,
	}

	typedConf, err := anypb.New(providerConf)
	if err != nil {
		return nil, err
	}

	tracing := &hcmv3.HttpConnectionManager_Tracing{
		Provider: &envoy_config_trace_v3.Tracing_Http{
			Name: tracerNameZipkin,
			ConfigType: &envoy_config_trace_v3.Tracing_Http_TypedConfig{
				TypedConfig: typedConf,
			},
		},
		MaxPathTagLength: &wrappers.UInt32Value{Value: maxPathLength},
	}

	return tracing, nil
}

func getListenerCodecType(codecType string) hcmv3.HttpConnectionManager_CodecType {
	switch codecType {
	case "AUTO":
		return hcmv3.HttpConnectionManager_AUTO
	case "HTTP1":
		return hcmv3.HttpConnectionManager_HTTP1
	case "HTTP2":
		return hcmv3.HttpConnectionManager_HTTP2
	default:
		return hcmv3.HttpConnectionManager_AUTO
	}
}

func getTracingOTLP(conf *config.Config) (*hcmv3.HttpConnectionManager_Tracing, error) {

	var maxPathLength uint32
	var connectionTimeout uint32

	if length, err := strconv.ParseUint(conf.Tracing.ConfigProperties[tracerMaxPathLength], 10, 32); err == nil {
		maxPathLength = uint32(length)
	} else {
		return nil, errors.New("invalid max path length provided for tracing endpoint")
	}

	if timeout, err := strconv.ParseUint(conf.Tracing.ConfigProperties[tracerConnectionTimeout], 10, 32); err == nil {
		connectionTimeout = uint32(timeout)
	} else {
		connectionTimeout = 20
		loggers.LoggerOasparser.Infof("Setting up default connection timeout for tracing endpoint as %d seconds", connectionTimeout)
	}

	providerConf := &envoy_config_trace_v3.OpenTelemetryConfig{
		GrpcService: &corev3.GrpcService{
			TargetSpecifier: &corev3.GrpcService_EnvoyGrpc_{
				EnvoyGrpc: &corev3.GrpcService_EnvoyGrpc{
					ClusterName: tracingClusterName,
				},
			},
			Timeout: durationpb.New(time.Duration(connectionTimeout) * time.Second),
		},
		ServiceName: tracerServiceNameRouter,
	}

	typedConf, err := anypb.New(providerConf)
	if err != nil {
		return nil, err
	}

	tracing := &hcmv3.HttpConnectionManager_Tracing{
		Provider: &envoy_config_trace_v3.Tracing_Http{
			Name: tracerNameOpenTelemetry,
			ConfigType: &envoy_config_trace_v3.Tracing_Http_TypedConfig{
				TypedConfig: typedConf,
			},
		},
		MaxPathTagLength: &wrappers.UInt32Value{Value: maxPathLength},
	}

	return tracing, nil
}
