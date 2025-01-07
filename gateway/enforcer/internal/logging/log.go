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
	"io"
	"os"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	egv1a1 "github.com/envoyproxy/gateway/api/v1alpha1"
)

const (
	// ErrorLevel is the logr verbosity level for errors.
	ErrorLevel = 0

	// WarnLevel is the logr verbosity level for warnings.
	WarnLevel = 0

	// InfoLevel is the logr verbosity level for info logs.
	InfoLevel = 0

	// DebugLevel is the logr verbosity level for debug logs.
	DebugLevel = 1

	// TraceLevel is the logr verbosity level for trace logs.
	TraceLevel = 2
)

// Logger is a struct that embeds logr.Logger and provides additional logging capabilities.
// It includes a reference to EnvoyGatewayLogging configuration and a SugaredLogger for
// structured logging.
//
// Fields:
// - logging: A pointer to EnvoyGatewayLogging configuration.
// - sugaredLogger: A SugaredLogger instance for structured logging.
type Logger struct {
	logr.Logger
	logging       *egv1a1.EnvoyGatewayLogging
	sugaredLogger *zap.SugaredLogger
}

// NewLogger creates a new Logger instance that logs to stdout.
// It uses the provided EnvoyGatewayLogging configuration and initializes
// the logger with the default log level for the Gateway component.
func NewLogger(logging *egv1a1.EnvoyGatewayLogging) Logger {
	logger := initZapLogger(os.Stdout, logging, logging.Level[egv1a1.LogComponentGatewayDefault])

	return Logger{
		Logger:        zapr.NewLogger(logger),
		logging:       logging,
		sugaredLogger: logger.Sugar(),
	}
}

// FileLogger creates a Logger instance that logs to the specified file.
// The log level is configured using the provided level parameter.
// If the file cannot be opened, it panics.
func FileLogger(file string, name string, level egv1a1.LogLevel) Logger {
	writer, err := os.OpenFile(file, os.O_WRONLY, 0o666)
	if err != nil {
		panic(err)
	}

	logging := egv1a1.DefaultEnvoyGatewayLogging()
	logger := initZapLogger(writer, logging, level)

	return Logger{
		Logger:        zapr.NewLogger(logger).WithName(name),
		logging:       logging,
		sugaredLogger: logger.Sugar(),
	}
}

// DefaultLogger creates a Logger instance with default logging settings.
// It logs to stdout and uses the specified log level for all components.
func DefaultLogger(level egv1a1.LogLevel) Logger {
	logging := egv1a1.DefaultEnvoyGatewayLogging()
	logger := initZapLogger(os.Stdout, logging, level)

	return Logger{
		Logger:        zapr.NewLogger(logger),
		logging:       logging,
		sugaredLogger: logger.Sugar(),
	}
}

// WithName returns a new Logger instance with the specified name element added
// to the Logger's name.  Successive calls with WithName append additional
// suffixes to the Logger's name.  It's strongly recommended that name segments
// contain only letters, digits, and hyphens (see the package documentation for
// more information).
func (l Logger) WithName(name string) Logger {
	logLevel := l.logging.Level[egv1a1.EnvoyGatewayLogComponent(name)]
	logger := initZapLogger(os.Stdout, l.logging, logLevel)

	return Logger{
		Logger:        zapr.NewLogger(logger).WithName(name),
		logging:       l.logging,
		sugaredLogger: logger.Sugar(),
	}
}

// WithValues returns a new Logger instance with additional key/value pairs.
// See Info for documentation on how key/value pairs work.
func (l Logger) WithValues(keysAndValues ...interface{}) Logger {
	l.Logger = l.Logger.WithValues(keysAndValues...)
	return l
}

// Sugar wraps the base Logger functionality in a slower, but less
// verbose, API. Any Logger can be converted to a SugaredLogger with its Sugar
// method.
//
// Unlike the Logger, the SugaredLogger doesn't insist on structured logging.
// For each log level, it exposes four methods:
//
//   - methods named after the log level for log.Print-style logging
//   - methods ending in "w" for loosely-typed structured logging
//   - methods ending in "f" for log.Printf-style logging
//   - methods ending in "ln" for log.Println-style logging
//
// For example, the methods for InfoLevel are:
//
//	Info(...any)           Print-style logging
//	Infow(...any)          Structured logging (read as "info with")
//	Infof(string, ...any)  Printf-style logging
//	Infoln(...any)         Println-style logging
func (l Logger) Sugar() *zap.SugaredLogger {
	return l.sugaredLogger
}

func initZapLogger(w io.Writer, logging *egv1a1.EnvoyGatewayLogging, level egv1a1.LogLevel) *zap.Logger {
	parseLevel, _ := zapcore.ParseLevel(string(logging.DefaultEnvoyGatewayLoggingLevel(level)))
	core := zapcore.NewCore(zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()), zapcore.AddSync(w), zap.NewAtomicLevelAt(parseLevel))

	return zap.New(core, zap.AddCaller())
}


// Debug logs a debug level message using the provided logger.
// The log level is set to DebugLevel and the message is logged with Info method.
//
// Parameters:
//   log (logr.Logger): The logger instance to use for logging.
//   msg (string): The debug message to log.
func Debug(log logr.Logger, msg string) {
	log.V(DebugLevel).Info(msg)
}

// Info logs an informational message using the provided logger.
// The log level is set to InfoLevel and the message is logged with the Info method.
//
// Parameters:
//   log (logr.Logger): The logger instance to use for logging.
//   msg (string): The informational message to log.
func Info(log logr.Logger, msg string) {
	log.V(InfoLevel).Info(msg)
}