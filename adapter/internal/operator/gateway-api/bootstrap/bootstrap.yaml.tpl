admin:
  access_log:
  - name: envoy.access_loggers.file
    typed_config:
      "@type": type.googleapis.com/envoy.extensions.access_loggers.file.v3.FileAccessLog
      path: {{ .AdminServer.AccessLogPath }}
  address:
    socket_address:
      address: {{ .AdminServer.Address }}
      port_value: {{ .AdminServer.Port }}
{{- if .StatsMatcher  }}
stats_config:
  stats_matcher:
    inclusion_list:
      patterns:
      {{- range $_, $item := .StatsMatcher.Exacts }}
      - exact: {{$item}}
      {{- end}}
      {{- range $_, $item := .StatsMatcher.Prefixs }}
      - prefix: {{$item}}
      {{- end}}
      {{- range $_, $item := .StatsMatcher.Suffixs }}
      - suffix: {{$item}}
      {{- end}}
      {{- range $_, $item := .StatsMatcher.RegularExpressions }}
      - safe_regex:
          google_re2: {}
          regex: {{js $item}}
      {{- end}}
{{- end }}
layered_runtime:
  layers:
  - name: global_config
    static_layer:
      envoy.restart_features.use_eds_cache_for_ads: true
      re2.max_program_size.error_level: 4294967295
      re2.max_program_size.warn_level: 1000
dynamic_resources:
  ads_config:
    api_type: DELTA_GRPC
    transport_api_version: V3
    grpc_services:
    - envoy_grpc:
        cluster_name: xds_cluster
    set_node_on_first_message_only: true
  lds_config:
    ads: {}
    resource_api_version: V3
  cds_config:
    ads: {}
    resource_api_version: V3
{{- if .OtelMetricSinks }}
stats_sinks:
{{- range $idx, $sink := .OtelMetricSinks }}
- name: "envoy.stat_sinks.open_telemetry"
  typed_config:
    "@type": type.googleapis.com/envoy.extensions.stat_sinks.open_telemetry.v3.SinkConfig
    grpc_service:
      envoy_grpc:
        cluster_name: otel_metric_sink_{{ $idx }}
{{- end }}
{{- end }}
static_resources:
  listeners:
  - name: envoy-gateway-proxy-ready-{{ .ReadyServer.Address }}-{{ .ReadyServer.Port }}
    address:
      socket_address:
        address: {{ .ReadyServer.Address }}
        port_value: {{ .ReadyServer.Port }}
        protocol: TCP
    filter_chains:
    - filters:
      - name: envoy.filters.network.http_connection_manager
        typed_config:
          "@type": type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager
          stat_prefix: eg-ready-http
          route_config:
            name: local_route
          http_filters:
          - name: envoy.filters.http.health_check
            typed_config:
              "@type": type.googleapis.com/envoy.extensions.filters.http.health_check.v3.HealthCheck
              pass_through_mode: false
              headers:
              - name: ":path"
                string_match:
                  exact: {{ .ReadyServer.ReadinessPath }}
          - name: envoy.filters.http.ext_authz
            typed_config:
              "@type": type.googleapis.com/envoy.extensions.filters.http.ext_authz.v3.ExtAuthz
              clear_route_cache: true
              include_peer_certificate: true
              transport_api_version: 2
              grpc_service: 
                  envoy_grpc: 
                    cluster_name: ext-authz
          - name: envoy.filters.http.router
            typed_config:
              "@type": type.googleapis.com/envoy.extensions.filters.http.router.v3.Router
  clusters:
  {{- range $idx, $sink := .OtelMetricSinks }}
  - name: otel_metric_sink_{{ $idx }}
    connect_timeout: 0.250s
    type: STRICT_DNS
    typed_extension_protocol_options:
      envoy.extensions.upstreams.http.v3.HttpProtocolOptions:
        "@type": "type.googleapis.com/envoy.extensions.upstreams.http.v3.HttpProtocolOptions"
        explicit_http_config:
          http2_protocol_options: {}
    lb_policy: ROUND_ROBIN
    load_assignment:
      cluster_name: otel_metric_sink_{{ $idx }}
      endpoints:
      - lb_endpoints:
        - endpoint:
            address:
              socket_address:
                address: {{ $sink.Address }}
                port_value: {{ $sink.Port }}
  {{- end }}
  - connect_timeout: 10s
    load_assignment:
      cluster_name: xds_cluster
      endpoints:
      - load_balancing_weight: 1
        lb_endpoints:
        - load_balancing_weight: 1
          endpoint:
            address:
              socket_address:
                address: apk-test-wso2-apk-adapter-service.apk.svc
                port_value: {{ .XdsServer.Port }}
    typed_extension_protocol_options:
      envoy.extensions.upstreams.http.v3.HttpProtocolOptions:
        "@type": "type.googleapis.com/envoy.extensions.upstreams.http.v3.HttpProtocolOptions"
        explicit_http_config:
          http2_protocol_options:
            connection_keepalive:
              interval: 30s
              timeout: 5s
    name: xds_cluster
    type: STRICT_DNS
    transport_socket:
      name: envoy.transport_sockets.tls
      typed_config:
        "@type": type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.UpstreamTlsContext
        common_tls_context:
          tls_params:
            tls_maximum_protocol_version: TLSv1_3
          tls_certificates:
            private_key:
              filename: '/home/wso2/security/keystore/router.key'
            certificate_chain:
              filename: '/home/wso2/security/keystore/router.crt'
          validation_context:
            trusted_ca:
              filename: '/home/wso2/security/truststore/adapter.crt'
  - name: ext-authz
    type: STRICT_DNS
    connect_timeout: 20s
    typed_extension_protocol_options:
      envoy.extensions.upstreams.http.v3.HttpProtocolOptions:
        "@type": type.googleapis.com/envoy.extensions.upstreams.http.v3.HttpProtocolOptions
        explicit_http_config:
          http2_protocol_options: {}
    transport_socket:
      name: envoy.transport_sockets.tls
      typed_config:
        "@type": type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.UpstreamTlsContext
        common_tls_context:
          tls_params:
            tls_minimum_protocol_version: TLSv1_2
            tls_maximum_protocol_version: TLSv1_2
          tls_certificates:
            private_key:
              filename: "/home/wso2/security/keystore/router.key"
            certificate_chain:
              filename: "/home/wso2/security/keystore/router.crt"
          validation_context:
            trusted_ca:
              filename: "/home/wso2/security/truststore/enforcer.crt"
    load_assignment:
      cluster_name: ext-authz
      endpoints:
        - lb_endpoints:
            - endpoint:
                address:
                  socket_address:
                    address: enforcer
                    port_value: 8081
  