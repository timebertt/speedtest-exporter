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
package speedtest

import (
	"context"
	"log"
	"time"

	"github.com/timebertt/speedtest-exporter/pkg/metrics"
)

func Run(ctx context.Context, interval time.Duration, logger *log.Logger, runner TestRunner) error {
	// don't wait `interval`, start test right away
	timer := time.NewTimer(time.Second)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-timer.C:
			timer.Reset(interval)

			logger.Println("starting speedtest")
			results, err := runner.Run(ctx)
			if err != nil {
				return err
			}
			logger.Printf("finished speedtest with results latency: %.0f, download: %.2f, upload: %.2f, duration: %.0f", results.latency, results.downloadSpeedMbps, results.uploadSpeedMbps, results.duration.Seconds())
			results.Record()
		}
	}
}

func (r *TestResults) Record() {
	// counter metrics
	metrics.TestsTotal.WithLabelValues(r.mode, r.serverName, r.serverSponsor).Inc()

	// gauge metrics
	metrics.TestDurationSeconds.WithLabelValues(r.mode).Set(r.duration.Seconds())
	metrics.LatencyMilliseconds.WithLabelValues(r.mode).Set(r.latency)
	metrics.DownloadSpeedMbps.WithLabelValues(r.mode).Set(r.downloadSpeedMbps)
	metrics.UploadSpeedMbps.WithLabelValues(r.mode).Set(r.uploadSpeedMbps)

	// histogram metrics
	metrics.TestDurationSecondsHistogram.WithLabelValues(r.mode).Observe(r.duration.Seconds())
	metrics.LatencyMillisecondsHistogram.WithLabelValues(r.mode).Observe(r.latency)
	metrics.DownloadSpeedMbpsHistogram.WithLabelValues(r.mode).Observe(r.downloadSpeedMbps)
	metrics.UploadSpeedMbpsHistogram.WithLabelValues(r.mode).Observe(r.uploadSpeedMbps)
}
