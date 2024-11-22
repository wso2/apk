/*
 *  Copyright (c) 2024, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 *  This file contains code derived from Envoy Gateway,
 *  https://github.com/envoyproxy/gateway
 *  and is provided here subject to the following:
 *  Copyright Project Envoy Gateway Authors
 *
 */

package envoy

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/wso2/apk/adapter/internal/loggers"
	"github.com/wso2/apk/adapter/internal/operator/gateway-api/bootstrap"
)

const (
	// ShutdownManagerPort is the port Envoy shutdown manager will listen on.
	ShutdownManagerPort = 19002
	// ShutdownManagerHealthCheckPath is the path used for health checks.
	ShutdownManagerHealthCheckPath = "/healthz"
	// ShutdownManagerReadyPath is the path used to indicate shutdown readiness.
	ShutdownManagerReadyPath = "/shutdown/ready"
	// ShutdownReadyFile is the file used to indicate shutdown readiness.
	ShutdownReadyFile = "/tmp/shutdown-ready"
)

// ShutdownManager serves shutdown manager process for Envoy proxies.
func ShutdownManager(readyTimeout time.Duration) error {
	// Setup HTTP handler
	handler := http.NewServeMux()
	handler.HandleFunc(ShutdownManagerHealthCheckPath, func(_ http.ResponseWriter, _ *http.Request) {})
	handler.HandleFunc(ShutdownManagerReadyPath, func(w http.ResponseWriter, _ *http.Request) {
		shutdownReadyHandler(w, readyTimeout, ShutdownReadyFile)
	})

	// Setup HTTP server
	srv := http.Server{
		Handler:           handler,
		Addr:              fmt.Sprintf(":%d", ShutdownManagerPort),
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       15 * time.Second,
	}

	// Setup signal handling
	c := make(chan struct{})
	go func() {
		s := make(chan os.Signal, 1)
		signal.Notify(s, os.Interrupt, syscall.SIGTERM)

		r := <-s
		// loggers.LoggerAPKOperator.Info(fmt.Sprintf("received %s", unix.SignalName(r.(syscall.Signal))))
		loggers.LoggerAPKOperator.Info(fmt.Sprintf("received %s", strconv.Itoa(int(r.(syscall.Signal)))))

		// Shutdown HTTP server without interrupting active connections
		if err := srv.Shutdown(context.Background()); err != nil {
			loggers.LoggerAPKOperator.Error(err, "server shutdown error")
		}
		close(c)
	}()

	// Start HTTP server
	loggers.LoggerAPKOperator.Info("starting shutdown manager")
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		loggers.LoggerAPKOperator.Error(err, "starting shutdown manager failed")
	}

	// Wait until done
	<-c
	return nil
}

// shutdownReadyHandler handles the endpoint used by a preStop hook on the Envoy
// container to block until ready to terminate. After the graceful drain process
// has completed a file will be written to indicate shutdown readiness.
func shutdownReadyHandler(w http.ResponseWriter, readyTimeout time.Duration, readyFile string) {
	var startTime = time.Now()

	loggers.LoggerAPKOperator.Info("received shutdown ready request")

	// Poll for shutdown readiness
	for {
		// Check if ready timeout is exceeded
		elapsedTime := time.Since(startTime)
		if elapsedTime > readyTimeout {
			loggers.LoggerAPKOperator.Info("shutdown readiness timeout exceeded")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err := os.Stat(readyFile)
		switch {
		case os.IsNotExist(err):
			time.Sleep(1 * time.Second)
		case err != nil:
			loggers.LoggerAPKOperator.Error(err, "error checking for shutdown readiness")
		case err == nil:
			loggers.LoggerAPKOperator.Info("shutdown readiness detected")
			return
		}
	}
}

// Shutdown is called from a preStop hook on the shutdown-manager container where
// it will initiate a graceful drain sequence on the Envoy proxy and block until
// connections are drained or a timeout is exceeded.
func Shutdown(drainTimeout time.Duration, minDrainDuration time.Duration, exitAtConnections int) error {
	var startTime = time.Now()
	var allowedToExit = false

	loggers.LoggerAPKOperator.Info(fmt.Sprintf("initiating graceful drain with %.0f second minimum drain period and %.0f second timeout",
		minDrainDuration.Seconds(), drainTimeout.Seconds()))

	// Start failing active health checks
	if err := postEnvoyAdminAPI("healthcheck/fail"); err != nil {
		loggers.LoggerAPKOperator.Error(err, "error failing active health checks")
	}

	// Initiate graceful drain sequence
	if err := postEnvoyAdminAPI("drain_listeners?graceful&skip_exit"); err != nil {
		loggers.LoggerAPKOperator.Error(err, "error initiating graceful drain")
	}

	// Poll total connections from Envoy admin API until minimum drain period has
	// been reached and total connections reaches threshold or timeout is exceeded
	for {
		elapsedTime := time.Since(startTime)

		conn, err := getTotalConnections()
		if err != nil {
			loggers.LoggerAPKOperator.Error(err, "error getting total connections")
		}

		if elapsedTime > minDrainDuration && !allowedToExit {
			loggers.LoggerAPKOperator.Info(fmt.Sprintf("minimum drain period reached; will exit when total connections reaches %d", exitAtConnections))
			allowedToExit = true
		}

		if elapsedTime > drainTimeout {
			loggers.LoggerAPKOperator.Info("graceful drain sequence timeout exceeded")
			break
		} else if allowedToExit && conn != nil && *conn <= exitAtConnections {
			loggers.LoggerAPKOperator.Info("graceful drain sequence completed")
			break
		}

		time.Sleep(1 * time.Second)
	}

	// Signal to shutdownReadyHandler that drain process is complete
	if _, err := os.Create(ShutdownReadyFile); err != nil {
		loggers.LoggerAPKOperator.Error(err, "error creating shutdown ready file")
		return err
	}

	return nil
}

// postEnvoyAdminAPI sends a POST request to the Envoy admin API
func postEnvoyAdminAPI(path string) error {
	resp, err := http.Post(fmt.Sprintf("http://%s:%d/%s",
		bootstrap.EnvoyAdminAddress, bootstrap.EnvoyAdminPort, path), "application/json", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response status: %s", resp.Status)
	}
	return nil

}

// getTotalConnections retrieves the total number of open connections from Envoy's server.total_connections stat
func getTotalConnections() (*int, error) {
	// Send request to Envoy admin API to retrieve server.total_connections stat
	resp, err := http.Get(fmt.Sprintf("http://%s:%d//stats?filter=^server\\.total_connections$&format=json",
		bootstrap.EnvoyAdminAddress, bootstrap.EnvoyAdminPort))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response status: %s", resp.Status)
	}
	// Define struct to decode JSON response into; expecting a single stat in the response in the format:
	// {"stats":[{"name":"server.total_connections","value":123}]}
	var r *struct {
		Stats []struct {
			Name  string
			Value int
		}
	}

	// Decode JSON response into struct
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}

	// Defensive check for empty stats
	if len(r.Stats) == 0 {
		return nil, fmt.Errorf("no stats found")
	}

	// Log and return total connections
	c := r.Stats[0].Value
	loggers.LoggerAPKOperator.Info(fmt.Sprintf("total connections: %d", c))
	return &c, nil

}
