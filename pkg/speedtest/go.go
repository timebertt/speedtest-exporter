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
)

type GoTestRunner struct {
	Logger *log.Logger
}

func (g GoTestRunner) Run(_ context.Context) (*TestResults, error) {
	timeStart := time.Now()

	user, _ := speedtest.FetchUserInfo()

	serverList, _ := speedtest.FetchServerList(user)
	targets, _ := serverList.FindServer([]int{})
	target := targets[0]

	g.Logger.Printf("selected server %q (%s)", target.Name, target.Sponsor)
	if err := target.PingTest(); err != nil {
		return nil, fmt.Errorf("error executing ping test: %w", err)
	}
	if err := target.DownloadTest(false); err != nil {
		return nil, fmt.Errorf("error executing download test: %w", err)
	}
	if err := target.UploadTest(false); err != nil {
		return nil, fmt.Errorf("error executing upload test: %w", err)
	}

	testDuration := time.Since(timeStart)

	return &TestResults{
		mode:              "go",
		serverName:        target.Name,
		serverSponsor:     target.Sponsor,
		duration:          testDuration,
		latency:           float64(target.Latency.Milliseconds()),
		downloadSpeedMbps: target.DLSpeed,
		uploadSpeedMbps:   target.ULSpeed,
	}, nil
}
