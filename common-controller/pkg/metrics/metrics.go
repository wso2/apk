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
	"github.com/wso2/apk/adapter/pkg/logging"
	logger "github.com/wso2/apk/common-controller/internal/loggers"
	metrics "github.com/wso2/apk/common-go-libs/pkg/metrics"
)

// StartPrometheusMetricsServer initializes and starts the metrics server to expose metrics to prometheus.
func StartPrometheusMetricsServer(port int32) {

	collector := metrics.CustomMetricsCollector()
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
