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
	metrics "github.com/wso2/apk/common-go-libs/pkg/metrics"
	k8smetrics "sigs.k8s.io/controller-runtime/pkg/metrics"
)

// RegisterPrometheusCollector registers the Prometheus collector for metrics.
func RegisterPrometheusCollector() {

	collector := metrics.CustomMetricsCollector()
	k8smetrics.Registry.MustRegister(collector)
}
