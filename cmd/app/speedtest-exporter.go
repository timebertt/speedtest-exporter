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
package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
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
			return opts.run(cmd.Context(), logger)
		},
	}
	cmd.SilenceErrors = true

	opts.addFlags(cmd.Flags())

	return cmd
}

type options struct {
	interval    time.Duration
	bindAddress string
	port        int
}

func (o *options) addFlags(fs *pflag.FlagSet) {
	fs.DurationVar(&o.interval, "interval", 1*time.Minute, "Interval in which to execute speedtests")
	fs.StringVar(&o.bindAddress, "bind-address", "0.0.0.0", "Address for the metrics endpoint to listen on")
	fs.IntVar(&o.port, "port", 8080, "Port for the metrics endpoint to listen on")
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
		if err := speedtest.Run(ctx, o.interval, logger); err != nil {
			logger.Fatalf("error running speedtests: %v", err)
		}
	}()

	<-ctx.Done()
	logger.Println("shutdown signal received, shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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
