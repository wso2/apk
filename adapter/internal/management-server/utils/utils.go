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

package utils

import (
	"crypto/tls"

	"github.com/wso2/apk/adapter/config"
	"github.com/wso2/apk/adapter/internal/loggers"
	logging "github.com/wso2/apk/adapter/internal/logging"
	"github.com/wso2/apk/adapter/pkg/utils/tlsutils"
	"google.golang.org/grpc/credentials"
)

// GenerateTLSCredentials generate tls credentials
func GenerateTLSCredentials() (credentials.TransportCredentials, error) {
	conf := config.ReadConfigs()
	certPool := tlsutils.GetTrustedCertPool(conf.Adapter.Truststore.Location)
	certificate, err := tlsutils.GetServerCertificate(conf.Adapter.Keystore.CertPath,
		conf.Adapter.Keystore.KeyPath)
	if err != nil {
		loggers.LoggerGRPCClient.ErrorC(logging.PrintError(logging.Error2700, logging.BLOCKER, "Error while processing the private-public key pair : %v", err.Error()))
		return nil, err
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{certificate},
		RootCAs:      certPool,
	}
	return credentials.NewTLS(tlsConfig), nil
}
