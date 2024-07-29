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

package bootstrap

import (
	"flag"
	"fmt"
	"os"
	"path"
)

var (
	overrideTestData = flag.Bool("override-testdata", false, "if override the test output data.")
)

// func TestGetRenderedBootstrapConfig(t *testing.T) {
// 	cases := []struct {
// 		name         string
// 		proxyMetrics *egv1a1.ProxyMetrics
// 	}{
// 		{
// 			name: "disable-prometheus",
// 			proxyMetrics: &egv1a1.ProxyMetrics{
// 				Prometheus: &egv1a1.ProxyPrometheusProvider{
// 					Disable: true,
// 				},
// 			},
// 		},
// 		{
// 			name: "enable-prometheus",
// 			proxyMetrics: &egv1a1.ProxyMetrics{
// 				Prometheus: &egv1a1.ProxyPrometheusProvider{},
// 			},
// 		},
// 		{
// 			name: "otel-metrics",
// 			proxyMetrics: &egv1a1.ProxyMetrics{
// 				Prometheus: &egv1a1.ProxyPrometheusProvider{
// 					Disable: true,
// 				},
// 				Sinks: []egv1a1.ProxyMetricSink{
// 					{
// 						Type: egv1a1.MetricSinkTypeOpenTelemetry,
// 						OpenTelemetry: &egv1a1.ProxyOpenTelemetrySink{
// 							Host: "otel-collector.monitoring.svc",
// 							Port: 4317,
// 						},
// 					},
// 				},
// 			},
// 		},
// 		{
// 			name: "custom-stats-matcher",
// 			proxyMetrics: &egv1a1.ProxyMetrics{
// 				Matches: []egv1a1.StringMatch{
// 					{
// 						Type:  ptr.To(egv1a1.StringMatchExact),
// 						Value: "http.foo.bar.cluster.upstream_rq",
// 					},
// 					{
// 						Type:  ptr.To(egv1a1.StringMatchPrefix),
// 						Value: "http",
// 					},
// 					{
// 						Type:  ptr.To(egv1a1.StringMatchSuffix),
// 						Value: "upstream_rq",
// 					},
// 					{
// 						Type:  ptr.To(egv1a1.StringMatchRegularExpression),
// 						Value: "virtual.*",
// 					},
// 					{
// 						Type:  ptr.To(egv1a1.StringMatchPrefix),
// 						Value: "cluster",
// 					},
// 				},
// 			},
// 		},
// 	}

// 	for _, tc := range cases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			got, err := GetRenderedBootstrapConfig(tc.proxyMetrics)
// 			require.NoError(t, err)

// 			if *overrideTestData {
// 				// nolint:gosec
// 				err = os.WriteFile(path.Join("testdata", "render", fmt.Sprintf("%s.yaml", tc.name)), []byte(got), 0644)
// 				require.NoError(t, err)
// 				return
// 			}

// 			expected, err := readTestData(tc.name)
// 			require.NoError(t, err)
// 			assert.Equal(t, expected, got)
// 		})
// 	}
// }

func readTestData(caseName string) (string, error) {
	filename := path.Join("testdata", "render", fmt.Sprintf("%s.yaml", caseName))

	b, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
