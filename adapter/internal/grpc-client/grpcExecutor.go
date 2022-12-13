/*
 *  Copyright (c) 2022, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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

package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/wso2/apk/adapter/config"
	logger "github.com/wso2/apk/adapter/internal/loggers"
	"github.com/wso2/apk/adapter/pkg/logging"
	"github.com/wso2/apk/adapter/pkg/utils/tlsutils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	grpcStatus "google.golang.org/grpc/status"
)

// RetryPolicy holds configuration for grpc connection retries
type RetryPolicy struct {
	// Maximum number of time a failed grpc call will be retried. Set negative value to try indefinitely.
	MaxAttempts int
	// Time delay between retries. (In milli seconds)
	BackOffInMilliSeconds int
}

// GetConnection creates and returns a grpc client connection
func GetConnection(address string) (*grpc.ClientConn, error) {
	transportCredentials, err := generateTLSCredentials()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return grpc.DialContext(ctx, address, grpc.WithTransportCredentials(transportCredentials), grpc.WithBlock())
}

func generateTLSCredentials() (credentials.TransportCredentials, error) {
	conf := config.ReadConfigs()
	certPool := tlsutils.GetTrustedCertPool(conf.Adapter.Truststore.Location)
	certificate, err := tlsutils.GetServerCertificate(conf.Adapter.Keystore.CertPath,
		conf.Adapter.Keystore.KeyPath)
	if err != nil {
		logger.LoggerGRPCClient.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprintf("Error while processing the private-public key pair : %v", err.Error()),
			Severity:  logging.BLOCKER,
			ErrorCode: 2700,
		})
		return nil, err
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{certificate},
		RootCAs:      certPool,
	}
	return credentials.NewTLS(tlsConfig), nil
}

// ExecuteGRPCCall executes a grpc call
func ExecuteGRPCCall(connection *grpc.ClientConn, call func() (interface{}, error)) (interface{}, error) {
	conf := config.ReadConfigs()
	maxAttempts := conf.Adapter.GRPCClient.MaxAttempts
	backOffInMilliSeconds := conf.Adapter.GRPCClient.BackOffInMilliSeconds
	retries := 0
	response, err := call()
	for {
		if err != nil {
			errStatus, _ := grpcStatus.FromError(err)
			logger.LoggerGRPCClient.ErrorC(logging.ErrorDetails{
				Message:   fmt.Sprintf("GRPC call failed. errorCode: %v errorMessage: %v", errStatus.Code().String(), errStatus.Message()),
				Severity:  logging.CRITICAL,
				ErrorCode: 2701,
			})
			if maxAttempts < 0 {
				// If max attempts has a negative value, retry indefinitely by setting retry less than max attempts.
				retries = maxAttempts - 1
			} else {
				retries++
			}
			if retries <= maxAttempts {
				// Retry grpc call after BackOffInMilliSeconds
				time.Sleep(time.Duration(backOffInMilliSeconds) * time.Millisecond)
				response, err = call()
			} else {
				return response, err
			}
		} else {
			return response, nil
		}
	}
}
