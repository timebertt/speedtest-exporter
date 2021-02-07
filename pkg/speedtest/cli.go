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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"time"
)

type CLITestRunner struct {
	Logger *log.Logger
}

func (c CLITestRunner) Run(ctx context.Context) (*TestResults, error) {
	timeStart := time.Now()

	outBuffer := &bytes.Buffer{}

	cmd := exec.CommandContext(ctx, "speedtest-cli", "--json")
	cmd.Stdout = outBuffer

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	testDuration := time.Since(timeStart)

	rawResults := &cliRawResults{}
	if err := json.Unmarshal(outBuffer.Bytes(), rawResults); err != nil {
		return nil, fmt.Errorf("error unmarshalling cli results: %w", err)
	}

	return &TestResults{
		mode:              "cli",
		serverName:        rawResults.Server.Name,
		serverSponsor:     rawResults.Server.Sponsor,
		duration:          testDuration,
		latency:           rawResults.Ping,
		downloadSpeedMbps: rawResults.Download / 1e6,
		uploadSpeedMbps:   rawResults.Upload / 1e6,
	}, nil
}

type cliRawResults struct {
	Download float64 `json:"download"`
	Upload   float64 `json:"upload"`
	Ping     float64 `json:"ping"`
	Server   struct {
		Name     string  `json:"name"`
		Sponsor  string  `json:"sponsor"`
		Distance float64 `json:"d"`
	} `json:"server"`
}
