/*
 *  Copyright (c) 2025, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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
 */
 
package logging

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	egv1a1 "github.com/envoyproxy/gateway/api/v1alpha1"
)

func TestZapLogLevel(t *testing.T) {
	level, err := zapcore.ParseLevel("warn")
	if err != nil {
		t.Errorf("ParseLevel error %v", err)
	}
	zc := zap.NewDevelopmentConfig()
	core := zapcore.NewCore(zapcore.NewConsoleEncoder(zc.EncoderConfig), zapcore.AddSync(os.Stdout), zap.NewAtomicLevelAt(level))
	zapLogger := zap.New(core, zap.AddCaller())
	log := zapLogger.Sugar()
	log.Info("ok", "k1", "v1")
	log.Error(errors.New("new error"), "error")
}

func TestLogger(t *testing.T) {
	logger := NewLogger(egv1a1.DefaultEnvoyGatewayLogging())
	logger.Info("kv msg", "key", "value")
	logger.Sugar().Infof("template %s %d", "string", 123)

	logger.WithName(string(egv1a1.LogComponentGlobalRateLimitRunner)).WithValues("runner", egv1a1.LogComponentGlobalRateLimitRunner).Info("msg", "k", "v")

	defaultLogger := DefaultLogger(egv1a1.LogLevelInfo)
	assert.NotNil(t, defaultLogger.logging)
	assert.NotNil(t, defaultLogger.sugaredLogger)

	fileLogger := FileLogger("/dev/stderr", "fl-test", egv1a1.LogLevelInfo)
	assert.NotNil(t, fileLogger.logging)
	assert.NotNil(t, fileLogger.sugaredLogger)
}

func TestLoggerWithName(t *testing.T) {
	originalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	defer func() {
		// Restore the original stdout and close the pipe
		os.Stdout = originalStdout
		err := w.Close()
		require.NoError(t, err)
	}()

	config := egv1a1.DefaultEnvoyGatewayLogging()
	config.Level[egv1a1.LogComponentInfrastructureRunner] = egv1a1.LogLevelDebug

	logger := NewLogger(config).WithName(string(egv1a1.LogComponentInfrastructureRunner))
	logger.Info("info message")
	logger.Sugar().Debugf("debug message")

	// Read from the pipe (captured stdout)
	outputBytes := make([]byte, 200)
	_, err := r.Read(outputBytes)
	require.NoError(t, err)
	capturedOutput := string(outputBytes)
	assert.Contains(t, capturedOutput, string(egv1a1.LogComponentInfrastructureRunner))
	assert.Contains(t, capturedOutput, "info message")
	assert.Contains(t, capturedOutput, "debug message")
}
