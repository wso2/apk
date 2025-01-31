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

package util

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/fs"
	"io/ioutil"
	"path/filepath"
)

// LoadCertificates loads a client certificate and private key from the given file paths.
// It returns the loaded tls.Certificate and an error if any issues occur during loading.
func LoadCertificates(publicKeyPath, privateKeyPath string) (tls.Certificate, error) {
	clientCert, err := tls.LoadX509KeyPair(publicKeyPath, privateKeyPath)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("failed to load client certificate and key: %v", err)
	}
	return clientCert, nil
}

// // LoadCACertificates loads the CA certificates from the provided file path.
// // It reads the certificate file and appends it to a new CertPool.
// // If any error occurs during reading or appending, it returns an error.
// func LoadCACertificates(trustedCertsPath string) (*x509.CertPool, error) {
// 	caCert, err := ioutil.ReadFile(trustedCertsPath)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to read CA certificate: %v", err)
// 	}

// 	certPool := x509.NewCertPool()
// 	if !certPool.AppendCertsFromPEM(caCert) {
// 		return nil, fmt.Errorf("failed to append CA certificate")
// 	}

// 	return certPool, nil
// }

// LoadCACertificates loads all CA certificates from the provided folder path.
// It reads all .pem or .crt files in the folder and appends them to a new CertPool.
// If any error occurs during reading or appending, it returns an error.
func LoadCACertificates(folderPath string) (*x509.CertPool, error) {
	certPool := x509.NewCertPool()

	err := filepath.WalkDir(folderPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing file %s: %v", path, err)
		}

		// Only process files with .pem or .crt extensions
		if !d.IsDir() && (filepath.Ext(path) == ".pem" || filepath.Ext(path) == ".crt") {
			caCert, err := ioutil.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read CA certificate from %s: %v", path, err)
			}

			if !certPool.AppendCertsFromPEM(caCert) {
				return fmt.Errorf("failed to append CA certificate from %s", path)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return certPool, nil
}

// CreateTLSConfig creates and returns a new TLS configuration using the provided client certificate
// and CA certificate pool. It sets up the certificates for secure communication.
func CreateTLSConfig(cert tls.Certificate, certPool *x509.CertPool) *tls.Config {
	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      certPool,
	}
}

// LoadCertificate loads an x509 certificate from a file path
func LoadCertificate(path string) (*x509.Certificate, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate file: %w", err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %w", err)
	}

	return cert, nil
}
