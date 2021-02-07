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
	"fmt"
	"log"
	"time"

	"github.com/showwin/speedtest-go/speedtest"

	"github.com/timebertt/speedtest-exporter/pkg/metrics"
)

func Run(ctx context.Context, interval time.Duration, logger *log.Logger) error {
	ticker := time.NewTicker(interval)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err := runTest(logger); err != nil {
				return err
			}
		}
	}
}

func runTest(logger *log.Logger) error {
	timeStart := time.Now()

	user, _ := speedtest.FetchUserInfo()

	serverList, _ := speedtest.FetchServerList(user)
	targets, _ := serverList.FindServer([]int{})
	target := targets[0]

	logger.Printf("starting speedtest against server %q (%s)", target.Name, target.Sponsor)
	if err := target.PingTest(); err != nil {
		return fmt.Errorf("error executing ping test: %w", err)
	}
	if err := target.DownloadTest(false); err != nil {
		return fmt.Errorf("error executing download test: %w", err)
	}
	if err := target.UploadTest(false); err != nil {
		return fmt.Errorf("error executing upload test: %w", err)
	}

	testDuration := time.Since(timeStart)
	logger.Printf("finished speedtest with results latency: %s, download: %f, upload: %f, duration: %s", target.Latency, target.DLSpeed, target.ULSpeed, testDuration)
	metrics.TestsTotal.WithLabelValues(target.Name, target.Sponsor).Inc()
	metrics.TestDurationSecondsHistogram.WithLabelValues().Observe(testDuration.Seconds())
	metrics.LatencyMillisecondsHistogram.WithLabelValues().Observe(float64(target.Latency.Milliseconds()))
	metrics.DownloadSpeedMbpsHistogram.WithLabelValues().Observe(target.DLSpeed)
	metrics.UploadSpeedMbpsHistogram.WithLabelValues().Observe(target.ULSpeed)

	return nil
}
