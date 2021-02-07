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
	"time"
)

type TestResults struct {
	mode string

	serverName        string
	serverSponsor     string
	duration          time.Duration
	latency           float64
	downloadSpeedMbps float64
	uploadSpeedMbps   float64
}

type TestRunner interface {
	Run(context.Context) (*TestResults, error)
}
