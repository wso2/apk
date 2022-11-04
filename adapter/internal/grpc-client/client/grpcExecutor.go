/*
Copyright 2022.

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

package client

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/wso2/apk/adapter/config"
	logger "github.com/wso2/apk/adapter/internal/loggers"
	"github.com/wso2/apk/adapter/pkg/logging"
	"github.com/wso2/apk/adapter/pkg/tlsutils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	grpcStatus "google.golang.org/grpc/status"
) 



type RetryPolicy struct {
	// Maximum number of time a failed grpc call will be retried. Set negative value to try indefinitely.
	MaxAttempts int;
	// Time delay between retries. (In milli seconds)
	BackOffInMilliSeconds int;
}

func GetConnection() (*grpc.ClientConn, error){
	conf, _ := config.ReadConfigs()
	address := conf.Adapter.GRPCClient.ManagementServerAddress;
	return grpc.Dial(address, []grpc.DialOption{
			grpc.WithTransportCredentials(generateTLSCredentials()),
			grpc.WithBlock()})
}

func generateTLSCredentials() credentials.TransportCredentials {
	conf, _ := config.ReadConfigs()
	certPool := tlsutils.GetTrustedCertPool(conf.Adapter.Truststore.Location)
	// There is a single private-public key pair for XDS server initialization, as well as for XDS client authentication
	certificate, err := tlsutils.GetServerCertificate(conf.Adapter.Keystore.CertPath,
		conf.Adapter.Keystore.KeyPath)
	if err != nil {
		logger.LoggerGRPCClient.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprintf("Error while processing the private-public key pair : %v", err.Error()),
			Severity:  logging.BLOCKER,
			ErrorCode: 1702,
		})
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{certificate},
		RootCAs:      certPool,
	}
	return credentials.NewTLS(tlsConfig)
}


func  ExecuteGRPCCall(connection *grpc.ClientConn, call func() (interface{}, error)) (interface{}, error) {
	conf, _ := config.ReadConfigs()
	maxAttempts := conf.Adapter.GRPCClient.MaxAttempts;
	backOffInMilliSeconds := conf.Adapter.GRPCClient.BackOffInMilliSeconds;
	retries := 0;
	response, err := call();
	for {
		
		if (err != nil) {
			errStatus, _ := grpcStatus.FromError(err)
			logger.LoggerGRPCClient.Errorf("gRPC call failed. errorCode: %s errorMessage: %s", errStatus.Code().String(), errStatus.Message());
			if (maxAttempts < 0) {
				// If max attempts has a negative value, retry indefinitely by setting retry less than max attempts.
				retries = maxAttempts - 1;
			} else {
				retries++;
			}
			if (retries <= maxAttempts) {
				// Retry grpc call after BackOffInMilliSeconds
				time.Sleep(time.Duration(backOffInMilliSeconds) * time.Millisecond)
				response, err = call();
			} else {
				return response, err;
			}
		} else {
			return response, nil;
		}
	}
}
