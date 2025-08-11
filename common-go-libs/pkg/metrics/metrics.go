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
	"runtime"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
)

var (
	prometheusMetricRegistry = prometheus.NewRegistry()
)

// Collector contains the descriptions of the custom metrics exposed
type Collector struct {
	hostInfo           *prometheus.Desc
	availableCPUs      *prometheus.Desc
	freePhysicalMemory *prometheus.Desc
	totalVirtualMemory *prometheus.Desc
	usedVirtualMemory  *prometheus.Desc
	systemCPULoad      *prometheus.Desc
	loadAvg            *prometheus.Desc
}

// CustomMetricsCollector contains the descriptions of the custom metrics exposed
func CustomMetricsCollector() *Collector {
	return &Collector{
		hostInfo: prometheus.NewDesc(
			"host_info",
			"Host Info",
			[]string{"os"}, nil,
		),
		availableCPUs: prometheus.NewDesc(
			"os_available_cpu_total",
			"Number of available CPUs.",
			nil, nil,
		),
		freePhysicalMemory: prometheus.NewDesc(
			"os_free_physical_memory_bytes",
			"Amount of free physical memory.",
			nil, nil,
		),
		totalVirtualMemory: prometheus.NewDesc(
			"os_total_virtual_memory_bytes",
			"Amount of total virtual memory.",
			nil, nil,
		),
		usedVirtualMemory: prometheus.NewDesc(
			"os_used_virtual_memory_bytes",
			"Amount of used virtual memory.",
			nil, nil,
		),
		systemCPULoad: prometheus.NewDesc(
			"os_system_cpu_load_percentage",
			"System-wide CPU usage as a percentage.",
			nil, nil,
		),
		loadAvg: prometheus.NewDesc(
			"os_system_load_average",
			"Current load of CPU in the host system for the last {x} minutes",
			[]string{"duration"}, nil,
		),
	}
}

// Describe sends all the descriptors of the metrics collected by this Collector
// to the provided channel.
func (collector *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.hostInfo
	ch <- collector.availableCPUs
	ch <- collector.freePhysicalMemory
	ch <- collector.totalVirtualMemory
	ch <- collector.usedVirtualMemory
	ch <- collector.systemCPULoad
	ch <- collector.loadAvg
}

// Collect collects all the relevant Prometheus metrics.
func (collector *Collector) Collect(ch chan<- prometheus.Metric) {

	var hostInfoValue float64
	var availableCPUs float64
	var freePhysicalMemory float64
	var totalVirtualMemory float64
	var usedVirtualMemory float64
	var systemCPULoad float64

	host, err := host.Info()
	if handleError(err, "Failed to get host info") {
		return
	}
	hostInfoValue = 1
	availableCPUs = float64(runtime.NumCPU())

	v, err := mem.VirtualMemory()
	if handleError(err, "Failed to read virtual memory metrics") {
		return
	}
	freePhysicalMemory = float64(v.Free)
	usedVirtualMemory = float64(v.Used)
	totalVirtualMemory = float64(v.Total)

	percentages, err := cpu.Percent(0, false)
	if handleError(err, "Failed to read cpu usage metrics") || len(percentages) == 0 {
		return
	}
	totalPercentage := 0.0
	for _, p := range percentages {
		totalPercentage += p
	}
	averagePercentage := totalPercentage / float64(len(percentages))
	systemCPULoad = averagePercentage

	avg, err := load.Avg()
	if handleError(err, "Failed to read cpu load averages") {
		return
	}

	ch <- prometheus.MustNewConstMetric(collector.hostInfo, prometheus.GaugeValue, hostInfoValue, host.OS)
	ch <- prometheus.MustNewConstMetric(collector.availableCPUs, prometheus.GaugeValue, availableCPUs)
	ch <- prometheus.MustNewConstMetric(collector.freePhysicalMemory, prometheus.GaugeValue, freePhysicalMemory)
	ch <- prometheus.MustNewConstMetric(collector.usedVirtualMemory, prometheus.GaugeValue, usedVirtualMemory)
	ch <- prometheus.MustNewConstMetric(collector.totalVirtualMemory, prometheus.GaugeValue, totalVirtualMemory)
	ch <- prometheus.MustNewConstMetric(collector.systemCPULoad, prometheus.GaugeValue, systemCPULoad)
	ch <- prometheus.MustNewConstMetric(collector.loadAvg, prometheus.GaugeValue, avg.Load1, "1m")
	ch <- prometheus.MustNewConstMetric(collector.loadAvg, prometheus.GaugeValue, avg.Load5, "5m")
	ch <- prometheus.MustNewConstMetric(collector.loadAvg, prometheus.GaugeValue, avg.Load15, "15m")
}

func init() {
	// Register the Go collector with the registry
	goCollector := collectors.NewGoCollector()
	prometheusMetricRegistry.MustRegister(goCollector)
}

func handleError(err error, message string) bool {
	if err != nil {
		return true
	}
	return false
}
