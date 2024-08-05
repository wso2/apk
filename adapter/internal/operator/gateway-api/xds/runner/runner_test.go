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

package runner

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/wso2/apk/adapter/config"
)

// func TestTLSConfig(t *testing.T) {
// 	// Create trusted CA, server and client certs.
// 	trustedCACert := certyaml.Certificate{
// 		Subject: "cn=trusted-ca",
// 	}
// 	egCertBeforeRotation := certyaml.Certificate{
// 		Subject:         "cn=eg-before-rotation",
// 		SubjectAltNames: []string{"DNS:localhost"},
// 		Issuer:          &trustedCACert,
// 	}
// 	egCertAfterRotation := certyaml.Certificate{
// 		Subject:         "cn=eg-after-rotation",
// 		SubjectAltNames: []string{"DNS:localhost"},
// 		Issuer:          &trustedCACert,
// 	}
// 	trustedEnvoyCert := certyaml.Certificate{
// 		Subject: "cn=trusted-envoy",
// 		Issuer:  &trustedCACert,
// 	}

// 	// Create another CA and a client cert to test that untrusted clients are denied.
// 	untrustedCACert := certyaml.Certificate{
// 		Subject: "cn=untrusted-ca",
// 	}
// 	untrustedClientCert := certyaml.Certificate{
// 		Subject: "cn=untrusted-client",
// 		Issuer:  &untrustedCACert,
// 	}

// 	caCertPool := x509.NewCertPool()
// 	ca, err := trustedCACert.X509Certificate()
// 	require.NoError(t, err)
// 	caCertPool.AddCert(&ca)

// 	tests := map[string]struct {
// 		serverCredentials *certyaml.Certificate
// 		clientCredentials *certyaml.Certificate
// 		expectError       bool
// 	}{
// 		"successful TLS connection established": {
// 			serverCredentials: &egCertBeforeRotation,
// 			clientCredentials: &trustedEnvoyCert,
// 			expectError:       false,
// 		},
// 		"rotating server credentials returns new server cert": {
// 			serverCredentials: &egCertAfterRotation,
// 			clientCredentials: &trustedEnvoyCert,
// 			expectError:       false,
// 		},
// 		"rotating server credentials again to ensure rotation can be repeated": {
// 			serverCredentials: &egCertBeforeRotation,
// 			clientCredentials: &trustedEnvoyCert,
// 			expectError:       false,
// 		},
// 		"fail to connect with client certificate which is not signed by correct CA": {
// 			serverCredentials: &egCertBeforeRotation,
// 			clientCredentials: &untrustedClientCert,
// 			expectError:       true,
// 		},
// 	}

// 	// Create temporary directory to store certificates and key for the server.
// 	configDir, err := os.MkdirTemp("", "eg-testdata-")
// 	require.NoError(t, err)
// 	defer os.RemoveAll(configDir)

// 	caFile := filepath.Join(configDir, "ca.crt")
// 	certFile := filepath.Join(configDir, "tls.crt")
// 	keyFile := filepath.Join(configDir, "tls.key")

// 	// Initial set of credentials must be written into temp directory before
// 	// starting the tests to avoid error at server startup.
// 	err = trustedCACert.WritePEM(caFile, keyFile)
// 	require.NoError(t, err)
// 	err = egCertBeforeRotation.WritePEM(certFile, keyFile)
// 	require.NoError(t, err)

// 	r := New(&Config{})
// 	g := grpc.NewServer(grpc.Creds(credentials.NewTLS(r.tlsConfig(certFile, keyFile, caFile))))
// 	if g == nil {
// 		t.Error("failed to create server")
// 	}

// 	address := "localhost:8001"
// 	l, err := net.Listen("tcp", address)
// 	require.NoError(t, err)

// 	go func() {
// 		err := g.Serve(l)
// 		require.NoError(t, err)
// 	}()
// 	defer g.GracefulStop()

// 	for name, tc := range tests {
// 		tc := tc
// 		t.Run(name, func(t *testing.T) {
// 			// Store certificate and key to temp dir used by serveContext.
// 			err = tc.serverCredentials.WritePEM(certFile, keyFile)
// 			require.NoError(t, err)
// 			clientCert, _ := tc.clientCredentials.TLSCertificate()
// 			receivedCert, err := tryConnect(address, clientCert, caCertPool)
// 			gotError := err != nil
// 			if gotError != tc.expectError {
// 				t.Errorf("Unexpected result when connecting to the server: %s", err)
// 			}
// 			if err == nil {
// 				expectedCert, _ := tc.serverCredentials.X509Certificate()
// 				assert.Equal(t, &expectedCert, receivedCert)
// 			}
// 		})
// 	}
// }

// tryConnect tries to establish TLS connection to the server.
// If successful, return the server certificate.
func tryConnect(address string, clientCert tls.Certificate, caCertPool *x509.CertPool) (*x509.Certificate, error) {
	clientConfig := &tls.Config{
		ServerName:   "localhost",
		MinVersion:   tls.VersionTLS13,
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      caCertPool,
	}
	conn, err := tls.Dial("tcp", address, clientConfig)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	err = peekError(conn)
	if err != nil {
		return nil, err
	}

	return conn.ConnectionState().PeerCertificates[0], nil
}

// peekError is a workaround for TLS 1.3: due to shortened handshake, TLS alert
// from server is received at first read from the socket. To receive alert for
// bad certificate, this function tries to read one byte.
// Adapted from https://golang.org/src/crypto/tls/handshake_client_test.go
func peekError(conn net.Conn) error {
	_ = conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
	_, err := conn.Read(make([]byte, 1))
	if err != nil {
		var netErr net.Error
		if !errors.As(netErr, &netErr) || !netErr.Timeout() {
			return err
		}
	}
	return nil
}

func TestServeXdsServerListenFailed(t *testing.T) {
	conf := config.ReadConfigs()
	// Occupy the address to make listening failed
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", conf.Deployment.Gateway.AdapterXDSPort))
	require.NoError(t, err)
	defer l.Close()

	r := New(&Config{})
	// Don't crash in this function
	r.serveXdsServer(context.Background())
}
