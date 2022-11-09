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

// Package tlsutils contains the utility functions related to tls communication of the adapter
package tlsutils

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	logger "github.com/wso2/apk/adapter/pkg/loggers"
	"github.com/wso2/apk/adapter/pkg/logging"
)

var (
	onceTrustedCertsRead sync.Once
	onceKeyCertsRead     sync.Once
	certificate          tls.Certificate
	certReadErr          error
	caCertPool           *x509.CertPool
)

const (
	pemExtension  string = ".pem"
	crtExtension  string = ".crt"
	authorization string = "Authorization"
)

// GetServerCertificate returns the certificate (used for the restAPI server and xds server) created based on configuration values.
func GetServerCertificate(tlsCertificate string, tlsCertificateKey string) (tls.Certificate, error) {
	certReadErr = nil
	onceKeyCertsRead.Do(func() {
		cert, err := tls.LoadX509KeyPair(string(tlsCertificate), string(tlsCertificateKey))
		if err != nil {
			logger.LoggerTLSUtils.ErrorC(logging.ErrorDetails{
							Message:   fmt.Sprintf("Error while loading the tls keypair. Error: %v", err),
							Severity:  logging.MINOR,
							ErrorCode: 2700,
						})
			certReadErr = err
		}
		certificate = cert
	})
	return certificate, certReadErr
}

// GetTrustedCertPool returns the trusted certificate (used for the restAPI server and xds server) created based on
// the provided directory/file path.
func GetTrustedCertPool(truststoreLocation string) *x509.CertPool {
	onceTrustedCertsRead.Do(func() {
		caCertPool = x509.NewCertPool()
		filepath.Walk(truststoreLocation, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				logger.LoggerTLSUtils.ErrorC(logging.ErrorDetails{
					Message:   fmt.Sprintf("Error while reading the trusted certificates directory/file. Error: %v", err),
					Severity:  logging.MINOR,
					ErrorCode: 2700,
				})
			} else {
				if !info.IsDir() && (filepath.Ext(info.Name()) == pemExtension ||
					filepath.Ext(info.Name()) == crtExtension) {
					caCert, caCertErr := ioutil.ReadFile(path)
					if caCertErr != nil {
						logger.LoggerTLSUtils.ErrorC(logging.ErrorDetails{
							Message:   fmt.Sprintf("Error while reading the certificate file. %v", info.Name()),
							Severity:  logging.MINOR,
							ErrorCode: 2700,
						})
					}
					if IsPublicCertificate(caCert) {
						caCertPool.AppendCertsFromPEM(caCert)
						logger.LoggerTLSUtils.Debugf("%v : Certificate is added as a trusted certificate.", info.Name())
					}
				}
			}
			return nil
		})
	})
	return caCertPool
}

// IsPublicCertificate checks if the file content represents valid public certificate in PEM format.
func IsPublicCertificate(certContent []byte) bool {
	certContentPattern := `\-\-\-\-\-BEGIN\sCERTIFICATE\-\-\-\-\-((.|\n)*)\-\-\-\-\-END\sCERTIFICATE\-\-\-\-\-`
	regex := regexp.MustCompile(certContentPattern)
	return regex.Match(certContent)
}
