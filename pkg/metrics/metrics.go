/*
Copyright © 2021 Tim Ebert

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	Namespace = "speedtest"
)

var (
	TestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: Namespace,
		Name:      "tests_total",
		Help:      "Total number of tests executed",
	}, []string{"mode", "server", "sponsor"})
	TestDurationSecondsHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: Namespace,
		Name:      "test_duration_seconds_histogram",
		Help:      "Test durations observed (histogram)",
		Buckets:   []float64{1, 15, 30, 45, 60, 75, 90, 105, 120, 135, 150},
	}, []string{"mode"})
	TestDurationSeconds = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Name:      "test_duration_seconds",
		Help:      "Test durations observed",
	}, []string{"mode"})
	LatencyMillisecondsHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: Namespace,
		Name:      "latency_milliseconds_histogram",
		Help:      "Latency in milliseconds (histogram)",
		Buckets:   []float64{0.25, 0.5, 1, 2, 4, 8, 16, 32, 64, 128, 256, 512},
	}, []string{"mode"})
	LatencyMilliseconds = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Name:      "latency_milliseconds",
		Help:      "Latency in milliseconds",
	}, []string{"mode"})
	DownloadSpeedMbpsHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: Namespace,
		Name:      "download_speed_mbps_histogram",
		Help:      "Download speed in Mb/s (histogram)",
		Buckets:   []float64{0.25, 0.5, 1, 2, 4, 8, 16, 32, 64, 128, 256, 512},
	}, []string{"mode"})
	DownloadSpeedMbps = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Name:      "download_speed_mbps",
		Help:      "Download speed in Mb/s",
	}, []string{"mode"})
	UploadSpeedMbpsHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: Namespace,
		Name:      "upload_speed_mbps_histogram",
		Help:      "Upload speed in Mb/s (histogram)",
		Buckets:   []float64{0.25, 0.5, 1, 2, 4, 8, 16, 32, 64, 128, 256, 512},
	}, []string{"mode"})
	UploadSpeedMbps = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Name:      "upload_speed_mbps",
		Help:      "Upload speed in Mb/s",
	}, []string{"mode"})
)
