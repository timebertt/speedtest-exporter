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
package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/timebertt/speedtest-exporter/pkg/speedtest"
)

func NewSpeedTestExporterCommand() *cobra.Command {
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	opts := &options{}

	cmd := &cobra.Command{
		Use:   "speedtest-exporter",
		Short: "A prometheus exporter for speedtest results",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			logger.Println("starting speedtest-exporter...")
			cmd.Flags().VisitAll(func(flag *pflag.Flag) {
				logger.Printf("[FLAG] --%s=%s", flag.Name, flag.Value)
			})

			if err := opts.validate(); err != nil {
				return err
			}
			cmd.SilenceUsage = true

			if err := opts.complete(logger); err != nil {
				return err
			}

			return opts.run(cmd.Context(), logger)
		},
	}
	cmd.SilenceErrors = true

	opts.addFlags(cmd.Flags())

	return cmd
}

type speedtestMode string

const (
	modeGo  speedtestMode = "go"
	modeCLI speedtestMode = "cli"
)

type options struct {
	interval    time.Duration
	bindAddress string
	port        int
	mode        speedtestMode

	runner speedtest.TestRunner
}

type modeFlag struct {
	setFunc func(speedtestMode)
	value   *speedtestMode
}

func (m *modeFlag) String() string {
	return string(*m.value)
}

func (m *modeFlag) Set(s string) error {
	val := strings.ToLower(s)
	switch val {
	case string(modeGo):
		m.setFunc(modeGo)
	case string(modeCLI):
		m.setFunc(modeCLI)
	default:
		return fmt.Errorf("must either be %q or %q", modeGo, modeCLI)
	}
	*m.value = speedtestMode(s)
	return nil
}

func (m *modeFlag) Type() string {
	return "string"
}

func (o *options) addFlags(fs *pflag.FlagSet) {
	fs.DurationVar(&o.interval, "interval", 90*time.Second, "Interval in which to execute speedtests")
	fs.StringVar(&o.bindAddress, "bind-address", "0.0.0.0", "Address for the metrics endpoint to listen on")
	fs.IntVar(&o.port, "port", 8080, "Port for the metrics endpoint to listen on")

	var modeValue modeFlag
	o.mode = modeGo
	modeValue.value = &o.mode
	modeValue.setFunc = func(m speedtestMode) {
		o.mode = m
	}
	fs.Var(&modeValue, "mode", fmt.Sprintf("Speedtest mode, %q or %q (default %q)", modeGo, modeCLI, modeGo))
}

func (o *options) validate() error {
	if o.interval <= 0 {
		return fmt.Errorf("interval must be greater or equal to zero, got: %v", o.interval)
	}
	if o.bindAddress == "" {
		return fmt.Errorf("bind address not set")
	}
	if o.port <= 0 {
		return fmt.Errorf("port not set")
	}
	if o.mode == modeCLI {
		if _, err := exec.LookPath("speedtest-cli"); err != nil {
			return fmt.Errorf("speedtest-cli not installed, please head to https://github.com/sivel/speedtest-cli and download it")
		}
	}

	return nil
}

func (o *options) complete(logger *log.Logger) error {
	switch o.mode {
	case modeGo:
		o.runner = &speedtest.GoTestRunner{Logger: logger}
	case modeCLI:
		o.runner = &speedtest.CLITestRunner{Logger: logger}
	default:
		return fmt.Errorf("invalid mode: %s", o.mode)
	}
	return nil
}

func (o *options) run(ctx context.Context, logger *log.Logger) error {
	listenAddress := fmt.Sprintf("%s:%d", o.bindAddress, o.port)
	logger.Printf("start listening on %s", listenAddress)

	mux := &http.ServeMux{}
	healthHandler := &okHandler{logger}
	mux.Handle("/readyz", healthHandler)
	mux.Handle("/healthz", healthHandler)
	mux.Handle("/metrics", promhttp.Handler())
	server := http.Server{Addr: listenAddress, Handler: mux}

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			logger.Fatalf("error listening: %v", err)
		}
	}()

	done := make(chan struct{})
	go func() {
		defer close(done)
		if err := speedtest.Run(ctx, o.interval, logger, o.runner); err != nil {
			logger.Fatalf("error running speedtest: %v", err)
		}
	}()

	<-ctx.Done()
	logger.Println("shutdown signal received, shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Printf("error shutting down http server: %v", err)
	}

	select {
	case <-shutdownCtx.Done():
	case <-done:
	}

	return nil
}

type okHandler struct {
	logger *log.Logger
}

func (h *okHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("ok"))
	if err != nil {
		h.logger.Printf("error writing response")
	}
}
