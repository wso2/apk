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

// Package metrics holds the implementation for exposing enforcer metrics to prometheus
package metrics

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	commonmetrics "github.com/wso2/apk/common-go-libs/pkg/metrics"
	"github.com/wso2/apk/gateway/enforcer/internal/datastore"
	"github.com/wso2/apk/gateway/enforcer/internal/transformer"
)

var (
	prometheusMetricRegistry = prometheus.NewRegistry()
)

var jwtTransformer *transformer.JWTTransformer
var subAppDataStore *datastore.SubscriptionApplicationDataStore

// enforcerCollector contains the descriptions of the custom metrics exposed by the adapter.
// It also uses the metrics defined in common-go-libs
type enforcerCollector struct {
	commonmetrics.Collector
	tokenIssuers  *prometheus.Desc
	subscriptions *prometheus.Desc
}

func enforcerMetricsCollector() *enforcerCollector {
	return &enforcerCollector{
		Collector: *commonmetrics.CustomMetricsCollector(),
		tokenIssuers: prometheus.NewDesc(
			"token_issuer_count",
			"Number of token issuers created.",
			nil, nil,
		),
		subscriptions: prometheus.NewDesc(
			"subscription_count",
			"Number of subscriptions created.",
			nil, nil,
		),
	}
}

// Describe sends all the descriptors of the metrics collected by this Collector
// to the provided channel.
func (collector *enforcerCollector) Describe(ch chan<- *prometheus.Desc) {
	collector.Collector.Describe(ch)
	ch <- collector.tokenIssuers
	ch <- collector.subscriptions
}

// Collect collects all the relevant Prometheus metrics when Prometheus requests it
func (collector *enforcerCollector) Collect(ch chan<- prometheus.Metric) {
	collector.Collector.Collect(ch)
	var tokenIssuerCount float64
	var subscriptionCount float64
	if jwtTransformer != nil {
		tokenIssuerCount = float64(jwtTransformer.GetTokenIssuerCount())
	}
	if subAppDataStore != nil {
		subscriptionCount = float64(subAppDataStore.GetTotalSubscriptionCount())
	}

	ch <- prometheus.MustNewConstMetric(collector.tokenIssuers, prometheus.GaugeValue, tokenIssuerCount)
	ch <- prometheus.MustNewConstMetric(collector.subscriptions, prometheus.GaugeValue, subscriptionCount)
}

// RegisterDataSources registers the data sources that the metrics would be scraped from.
func RegisterDataSources(transformer *transformer.JWTTransformer, dataStore *datastore.SubscriptionApplicationDataStore) {
	jwtTransformer = transformer
	subAppDataStore = dataStore
}

// StartPrometheusMetricsServer initializes and starts the metrics server to expose metrics to prometheus.
func StartPrometheusMetricsServer(port int32) {

	collector := enforcerMetricsCollector()
	prometheus.MustRegister(collector)
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":"+strconv.Itoa(int(port)), nil)
	fmt.Println("Metrics server started on port " + strconv.Itoa(int(port)))
	if err != nil {
		fmt.Println("Error starting the metrics server")
	}
}
