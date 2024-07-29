/*
 *  Copyright (c) 2024, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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
 *  This file contains code derived from Envoy Gateway,
 *  https://github.com/envoyproxy/gateway
 *  and is provided here subject to the following:
 *  Copyright Project Envoy Gateway Authors
 *
 */

package translator

import (
	"embed"
	"flag"
	"path/filepath"
	"strings"
	"testing"

	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	resourcev3 "github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"github.com/stretchr/testify/require"
	"github.com/wso2/apk/adapter/internal/operator/gateway-api/ir"
	"github.com/wso2/apk/adapter/internal/operator/gateway-api/xds/utils"
	"github.com/wso2/apk/adapter/pkg/utils/file"
	"sigs.k8s.io/yaml"
)

var (
	//go:embed testdata/out/*
	outFiles embed.FS
	//go:embed testdata/in/*
	inFiles embed.FS

	overrideTestData = flag.Bool("override-testdata", false, "if override the test output data.")
)

func TestTranslateXds(t *testing.T) {
	testCases := []struct {
		name                      string
		dnsDomain                 string
		requireSecrets            bool
		requireEnvoyPatchPolicies bool
	}{
		{
			name: "empty",
		},
		{
			name: "http-route",
		},
		{
			name: "http-route-regex",
		},
		{
			name: "http-route-redirect",
		},
		{
			name: "http-route-mirror",
		},
		{
			name: "http-route-multiple-mirrors",
		},
		{
			name: "http-route-multiple-matches",
		},
		{
			name: "http-route-direct-response",
		},
		{
			name: "http-route-request-headers",
		},
		{
			name: "http-route-response-add-headers",
		},
		{
			name: "http-route-response-remove-headers",
		},
		{
			name: "http-route-response-add-remove-headers",
		},
		{
			name: "http-route-weighted-invalid-backend",
		},
		{
			name: "http-route-dns-cluster",
		},
		{
			name:           "http-route-with-tls-system-truststore",
			requireSecrets: true,
		},
		{
			name:           "http-route-with-tlsbundle",
			requireSecrets: true,
		},
		{
			name:           "http-route-with-tlsbundle-multiple-certs",
			requireSecrets: true,
		},
		{
			name:           "simple-tls",
			requireSecrets: true,
		},
		{
			name:           "mutual-tls",
			requireSecrets: true,
		},
		{
			name:           "http3",
			requireSecrets: true,
		},
		{
			name: "tls-route-passthrough",
		},
		{
			name: "tcp-route-simple",
		},
		{
			name: "tcp-route-complex",
		},
		{
			name: "tcp-route-tls-terminate",
		},
		{
			name: "multiple-simple-tcp-route-same-port",
		},
		{
			name: "http-route-weighted-backend",
		},
		{
			name: "tcp-route-weighted-backend",
		},
		{
			name:           "multiple-listeners-same-port",
			requireSecrets: true,
		},
		{
			name: "udp-route",
		},
		{
			name: "http2-route",
		},
		{
			name: "http-route-rewrite-url-prefix",
		},
		{
			name: "http-route-rewrite-root-path-url-prefix",
		},
		{
			name: "http-route-rewrite-url-fullpath",
		},
		{
			name: "http-route-rewrite-url-host",
		},
		{
			name: "http-route-timeout",
		},
		{
			name: "accesslog",
		},
		{
			name: "tracing",
		},
		{
			name: "metrics-virtual-host",
		},
		{
			name:                      "jsonpatch-missing-resource",
			requireEnvoyPatchPolicies: true,
		},
		{
			name:                      "jsonpatch-invalid-patch",
			requireEnvoyPatchPolicies: true,
		},
		{
			name:                      "jsonpatch-add-op-without-value",
			requireEnvoyPatchPolicies: true,
		},
		{
			name:                      "jsonpatch-move-op-with-value",
			requireEnvoyPatchPolicies: true,
		},
		{
			name: "listener-tcp-keepalive",
		},
		{
			name: "load-balancer",
		},
		{
			name: "cors",
		},
		{
			name: "oidc",
		},
		{
			name: "http-route-partial-invalid",
		},
		{
			name: "listener-proxy-protocol",
		},
		{
			name: "proxy-protocol-upstream",
		},
		{
			name: "basic-auth",
		},
		{
			name: "health-check",
		},
		{
			name: "circuit-breaker",
		},
		{
			name: "suppress-envoy-headers",
		},
		{
			name: "fault-injection",
		},
		{
			name: "tls-with-ciphers-versions-alpn",
		},
		{
			name: "path-settings",
		},
		{
			name: "http1-trailers",
		},
		{
			name: "http1-preserve-case",
		},
		{
			name: "timeout",
		},
		{
			name: "ext-auth",
		},
		{
			name: "http10",
		},
		{
			name: "upstream-tcpkeepalive",
		},
		{
			name: "client-timeout",
		},
		{
			name: "retry-partial-invalid",
		},
		{
			name: "multiple-listeners-same-port-with-different-filters",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			dnsDomain := tc.dnsDomain
			if dnsDomain == "" {
				dnsDomain = "cluster.local"
			}
			ir := requireXdsIRFromInputTestData(t, "xds-ir", tc.name+".yaml")
			tr := &Translator{
				// GlobalRateLimit: &GlobalRateLimitSettings{
				// 	ServiceURL: ratelimit.GetServiceURL("envoy-gateway-system", dnsDomain),
				// },
			}

			tCtx, err := tr.Translate(ir)
			if !strings.HasSuffix(tc.name, "partial-invalid") {
				require.NoError(t, err)
			}

			listeners := tCtx.XdsResources[resourcev3.ListenerType]
			routes := tCtx.XdsResources[resourcev3.RouteType]
			clusters := tCtx.XdsResources[resourcev3.ClusterType]
			endpoints := tCtx.XdsResources[resourcev3.EndpointType]
			if *overrideTestData {
				require.NoError(t, file.Write(requireResourcesToYAMLString(t, listeners), filepath.Join("testdata", "out", "xds-ir", tc.name+".listeners.yaml")))
				require.NoError(t, file.Write(requireResourcesToYAMLString(t, routes), filepath.Join("testdata", "out", "xds-ir", tc.name+".routes.yaml")))
				require.NoError(t, file.Write(requireResourcesToYAMLString(t, clusters), filepath.Join("testdata", "out", "xds-ir", tc.name+".clusters.yaml")))
				require.NoError(t, file.Write(requireResourcesToYAMLString(t, endpoints), filepath.Join("testdata", "out", "xds-ir", tc.name+".endpoints.yaml")))
			}
			require.Equal(t, requireTestDataOutFile(t, "xds-ir", tc.name+".listeners.yaml"), requireResourcesToYAMLString(t, listeners))
			require.Equal(t, requireTestDataOutFile(t, "xds-ir", tc.name+".routes.yaml"), requireResourcesToYAMLString(t, routes))
			require.Equal(t, requireTestDataOutFile(t, "xds-ir", tc.name+".clusters.yaml"), requireResourcesToYAMLString(t, clusters))
			require.Equal(t, requireTestDataOutFile(t, "xds-ir", tc.name+".endpoints.yaml"), requireResourcesToYAMLString(t, endpoints))
			if tc.requireSecrets {
				secrets := tCtx.XdsResources[resourcev3.SecretType]
				if *overrideTestData {
					require.NoError(t, file.Write(requireResourcesToYAMLString(t, secrets), filepath.Join("testdata", "out", "xds-ir", tc.name+".secrets.yaml")))
				}
				require.Equal(t, requireTestDataOutFile(t, "xds-ir", tc.name+".secrets.yaml"), requireResourcesToYAMLString(t, secrets))
			}
		})
	}
}

func TestTranslateXdsNegative(t *testing.T) {
	testCases := []struct {
		name           string
		dnsDomain      string
		requireSecrets bool
	}{
		{
			name: "http-route-invalid",
		},
		{
			name: "tcp-route-invalid",
		},
		{
			name: "tcp-route-invalid-endpoint",
		},
		{
			name: "udp-route-invalid",
		},
		{
			name: "jsonpatch-invalid",
		},
		{
			name: "accesslog-invalid",
		},
		{
			name: "tracing-invalid",
		},
		{
			name: "jsonpatch-invalid-listener",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			dnsDomain := tc.dnsDomain
			if dnsDomain == "" {
				dnsDomain = "cluster.local"
			}
			ir := requireXdsIRFromInputTestData(t, "xds-ir", tc.name+".yaml")
			tr := &Translator{
				// GlobalRateLimit: &GlobalRateLimitSettings{
				// 	ServiceURL: ratelimit.GetServiceURL("envoy-gateway-system", dnsDomain),
				// },
			}

			_, err := tr.Translate(ir)
			require.Error(t, err)
			if tc.name != "jsonpatch-invalid" {
				require.Contains(t, err.Error(), "validation failed for xds resource")
			}
		})
	}
}

// func TestTranslateRateLimitConfig(t *testing.T) {
// 	testCases := []struct {
// 		name string
// 	}{
// 		{
// 			name: "empty-header-matches",
// 		},
// 		{
// 			name: "distinct-match",
// 		},
// 		{
// 			name: "distinct-remote-address-match",
// 		},
// 		{
// 			name: "value-match",
// 		},
// 		{
// 			name: "multiple-matches",
// 		},
// 		{
// 			name: "multiple-rules",
// 		},
// 		{
// 			name: "multiple-routes",
// 		},
// 		{
// 			name: "masked-remote-address-match",
// 		},
// 		{
// 			name: "multiple-masked-remote-address-match-with-same-cidr",
// 		},
// 	}

// 	for _, tc := range testCases {
// 		tc := tc
// 		t.Run(tc.name, func(t *testing.T) {
// 			in := requireXdsIRListenerFromInputTestData(t, "ratelimit-config", tc.name+".yaml")
// 			out := BuildRateLimitServiceConfig(in)
// 			if *overrideTestData {
// 				require.NoError(t, file.Write(requireYamlRootToYAMLString(t, out), filepath.Join("testdata", "out", "ratelimit-config", tc.name+".yaml")))
// 			}
// 			require.Equal(t, requireTestDataOutFile(t, "ratelimit-config", tc.name+".yaml"), requireYamlRootToYAMLString(t, out))
// 		})
// 	}
// }

// func TestTranslateXdsWithExtension(t *testing.T) {
// 	testCases := []struct {
// 		name           string
// 		requireSecrets bool
// 		err            string
// 	}{
// 		// Require secrets for all the tests since the extension for testing always injects one
// 		{
// 			name:           "empty",
// 			requireSecrets: true,
// 			err:            "",
// 		},
// 		{
// 			name:           "http-route",
// 			requireSecrets: true,
// 			err:            "",
// 		},
// 		{
// 			name:           "http-route-extension-filter",
// 			requireSecrets: true,
// 			err:            "",
// 		},
// 		{
// 			name:           "http-route-extension-route-error",
// 			requireSecrets: true,
// 			err:            "route hook resource error",
// 		},
// 		{
// 			name:           "http-route-extension-virtualhost-error",
// 			requireSecrets: true,
// 			err:            "extension post xds virtual host hook error",
// 		},
// 		{
// 			name:           "http-route-extension-listener-error",
// 			requireSecrets: true,
// 			err:            "extension post xds listener hook error",
// 		},
// 	}

// 	for _, tc := range testCases {
// 		tc := tc
// 		t.Run(tc.name, func(t *testing.T) {
// 			// Testdata for the extension tests is similar to the ir test dat
// 			// New directory is just to keep them separate and easy to understand
// 			ir := requireXdsIRFromInputTestData(t, "extension-xds-ir", tc.name+".yaml")
// 			tr := &Translator{
// 				// GlobalRateLimit: &GlobalRateLimitSettings{
// 				// 	ServiceURL: ratelimit.GetServiceURL("envoy-gateway-system", "cluster.local"),
// 				// },
// 			}
// 			// ext := v1alpha1.ExtensionManager{
// 			// 	Resources: []v1alpha1.GroupVersionKind{
// 			// 		{
// 			// 			Group:   "foo.example.io",
// 			// 			Version: "v1alpha1",
// 			// 			Kind:    "examplefilter",
// 			// 		},
// 			// 	},
// 			// 	Hooks: &v1alpha1.ExtensionHooks{
// 			// 		XDSTranslator: &v1alpha1.XDSTranslatorHooks{
// 			// 			Post: []v1alpha1.XDSTranslatorHook{
// 			// 				v1alpha1.XDSRoute,
// 			// 				v1alpha1.XDSVirtualHost,
// 			// 				v1alpha1.XDSHTTPListener,
// 			// 				v1alpha1.XDSTranslation,
// 			// 			},
// 			// 		},
// 			// 	},
// 			// }
// 			// extMgr := testutils.NewManager(ext)
// 			// tr.ExtensionManager = &extMgr

// 			tCtx, err := tr.Translate(ir)

// 			if tc.err != "" {
// 				require.EqualError(t, err, tc.err)
// 			} else {
// 				require.NoError(t, err)
// 				listeners := tCtx.XdsResources[resourcev3.ListenerType]
// 				routes := tCtx.XdsResources[resourcev3.RouteType]
// 				clusters := tCtx.XdsResources[resourcev3.ClusterType]
// 				endpoints := tCtx.XdsResources[resourcev3.EndpointType]
// 				if *overrideTestData {
// 					require.NoError(t, file.Write(requireResourcesToYAMLString(t, listeners), filepath.Join("testdata", "out", "extension-xds-ir", tc.name+".listeners.yaml")))
// 					require.NoError(t, file.Write(requireResourcesToYAMLString(t, routes), filepath.Join("testdata", "out", "extension-xds-ir", tc.name+".routes.yaml")))
// 					require.NoError(t, file.Write(requireResourcesToYAMLString(t, clusters), filepath.Join("testdata", "out", "extension-xds-ir", tc.name+".clusters.yaml")))
// 					require.NoError(t, file.Write(requireResourcesToYAMLString(t, endpoints), filepath.Join("testdata", "out", "extension-xds-ir", tc.name+".endpoints.yaml")))
// 				}
// 				require.Equal(t, requireTestDataOutFile(t, "extension-xds-ir", tc.name+".listeners.yaml"), requireResourcesToYAMLString(t, listeners))
// 				require.Equal(t, requireTestDataOutFile(t, "extension-xds-ir", tc.name+".routes.yaml"), requireResourcesToYAMLString(t, routes))
// 				require.Equal(t, requireTestDataOutFile(t, "extension-xds-ir", tc.name+".clusters.yaml"), requireResourcesToYAMLString(t, clusters))
// 				require.Equal(t, requireTestDataOutFile(t, "extension-xds-ir", tc.name+".endpoints.yaml"), requireResourcesToYAMLString(t, endpoints))
// 				if tc.requireSecrets {
// 					secrets := tCtx.XdsResources[resourcev3.SecretType]
// 					if *overrideTestData {
// 						require.NoError(t, file.Write(requireResourcesToYAMLString(t, secrets), filepath.Join("testdata", "out", "extension-xds-ir", tc.name+".secrets.yaml")))
// 					}
// 					require.Equal(t, requireTestDataOutFile(t, "extension-xds-ir", tc.name+".secrets.yaml"), requireResourcesToYAMLString(t, secrets))
// 				}
// 			}
// 		})
// 	}
// }

func requireXdsIRFromInputTestData(t *testing.T, name ...string) *ir.Xds {
	t.Helper()
	elems := append([]string{"testdata", "in"}, name...)
	content, err := inFiles.ReadFile(filepath.Join(elems...))
	require.NoError(t, err)
	ir := &ir.Xds{}
	err = yaml.Unmarshal(content, ir)
	require.NoError(t, err)
	return ir
}

func requireXdsIRListenerFromInputTestData(t *testing.T, name ...string) *ir.HTTPListener {
	t.Helper()
	elems := append([]string{"testdata", "in"}, name...)
	content, err := inFiles.ReadFile(filepath.Join(elems...))
	require.NoError(t, err)
	listener := &ir.HTTPListener{}
	err = yaml.Unmarshal(content, listener)
	require.NoError(t, err)
	return listener
}

func requireTestDataOutFile(t *testing.T, name ...string) string {
	t.Helper()
	elems := append([]string{"testdata", "out"}, name...)
	content, err := outFiles.ReadFile(filepath.Join(elems...))
	require.NoError(t, err)
	return string(content)
}

// func requireYamlRootToYAMLString(t *testing.T, pbRoot *ratelimitv3.RateLimitConfig) string {
// 	str, err := GetRateLimitServiceConfigStr(pbRoot)
// 	require.NoError(t, err)
// 	return str
// }

func requireResourcesToYAMLString(t *testing.T, resources []types.Resource) string {
	jsonBytes, err := utils.MarshalResourcesToJSON(resources)
	require.NoError(t, err)
	data, err := yaml.JSONToYAML(jsonBytes)
	require.NoError(t, err)
	return string(data)
}
