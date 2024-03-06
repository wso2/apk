/*
 * Copyright (c) 2024, WSO2 LLC. (https://www.wso2.com)
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

// Package metrics holds the implementation for exposing adapter metrics to prometheus
package metrics

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	xds "github.com/wso2/apk/adapter/internal/discovery/xds"
	logger "github.com/wso2/apk/adapter/internal/loggers"
	"github.com/wso2/apk/adapter/pkg/logging"
	commonmetrics "github.com/wso2/apk/common-go-libs/pkg/metrics"
)

var (
	prometheusMetricRegistry = prometheus.NewRegistry()
)

// AdapterCollector contains the descriptions of the custom metrics exposed by the adapter.
// It also uses the metrics defined in common-go-libs
type AdapterCollector struct {
	commonmetrics.Collector
	apis                 *prometheus.Desc
	internalClusterCount *prometheus.Desc
	internalRouteCount   *prometheus.Desc
}

func adapterMetricsCollector() *AdapterCollector {
	return &AdapterCollector{
		Collector: *commonmetrics.CustomMetricsCollector(),
		apis: prometheus.NewDesc(
			"api_count",
			"Number of APIs created.",
			nil, nil,
		),
		internalClusterCount: prometheus.NewDesc(
			"internal_cluster_count",
			"Number of internal clusters created.",
			nil, nil,
		),
		internalRouteCount: prometheus.NewDesc(
			"internal_route_count",
			"Number of internal routes created.",
			nil, nil,
		),
	}
}

// Describe sends all the descriptors of the metrics collected by this Collector
// to the provided channel.
func (collector *AdapterCollector) Describe(ch chan<- *prometheus.Desc) {
	collector.Collector.Describe(ch)
	ch <- collector.apis
	ch <- collector.internalClusterCount
	ch <- collector.internalRouteCount
}

// Collect collects all the relevant Prometheus metrics.
func (collector *AdapterCollector) Collect(ch chan<- prometheus.Metric) {
	collector.Collector.Collect(ch)
	var apisCount float64
	var internalClusterCount float64
	var internalRouteCount float64

	apiCount := xds.GetEnvoyInternalAPICount()
	apisCount = float64(apiCount)

	internalRouteCount = float64(xds.GetEnvoyInternalAPIRoutes())
	internalClusterCount = float64(xds.GetEnvoyInternalAPIClusters())

	ch <- prometheus.MustNewConstMetric(collector.apis, prometheus.GaugeValue, apisCount)
	ch <- prometheus.MustNewConstMetric(collector.internalRouteCount, prometheus.GaugeValue, internalRouteCount)
	ch <- prometheus.MustNewConstMetric(collector.internalClusterCount, prometheus.GaugeValue, internalClusterCount)
}

// StartPrometheusMetricsServer initializes and starts the metrics server to expose metrics to prometheus.
func StartPrometheusMetricsServer(port int32) {

	collector := adapterMetricsCollector()
	prometheus.MustRegister(collector)
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":"+strconv.Itoa(int(port)), nil)
	if err != nil {
		logger.LoggerAPK.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprintln("Prometheus metrics server error:", err),
			Severity:  logging.MAJOR,
			ErrorCode: 1110,
		})
	}
}
