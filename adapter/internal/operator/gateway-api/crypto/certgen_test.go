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

package crypto

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGenerateCerts(t *testing.T) {
	type testcase struct {
		certConfig              *Configuration
		wantEnvoyGatewayDNSName string
		wantEnvoyDNSName        string
	}

	run := func(t *testing.T, name string, tc testcase) {
		t.Helper()

		t.Run(name, func(t *testing.T) {
			t.Helper()

			got, err := GenerateCerts()
			require.NoError(t, err)

			roots := x509.NewCertPool()
			ok := roots.AppendCertsFromPEM(got.CACertificate)
			require.Truef(t, ok, "Failed to set up CA cert for testing, maybe it's an invalid PEM")

			currentTime := time.Now()

			err = verifyCert(got.EnvoyGatewayCertificate, roots, tc.wantEnvoyGatewayDNSName, currentTime)
			require.NoErrorf(t, err, "Validating %s failed", name)

			err = verifyCert(got.EnvoyCertificate, roots, tc.wantEnvoyDNSName, currentTime)
			require.NoErrorf(t, err, "Validating %s failed", name)
		})
	}

	run(t, "no configuration - use defaults", testcase{
		certConfig:              &Configuration{},
		wantEnvoyGatewayDNSName: DefaultEnvoyGatewayDNSPrefix,
		wantEnvoyDNSName:        fmt.Sprintf("*.%s", "apk"),
	})
}

func TestGeneratedValidKubeCerts(t *testing.T) {
	now := time.Now()
	expiry := now.Add(24 * 365 * time.Hour)

	caCert, caKey, err := newCA("envoy-gateway", expiry)
	require.NoErrorf(t, err, "Failed to generate CA cert")

	egCertReq := &certificateRequest{
		caCertPEM:  caCert,
		caKeyPEM:   caKey,
		expiry:     expiry,
		commonName: "envoy-gateway",
		altNames:   kubeServiceNames("envoy-gateway", "apk", "cluster.local"),
	}
	egCert, _, err := newCert(egCertReq)
	require.NoErrorf(t, err, "Failed to generate Envoy Gateway cert")

	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM(caCert)
	require.Truef(t, ok, "Failed to set up CA cert for testing, maybe it's an invalid PEM")

	envoyCertReq := &certificateRequest{
		caCertPEM:  caCert,
		caKeyPEM:   caKey,
		expiry:     expiry,
		commonName: "envoy",
		altNames:   kubeServiceNames("envoy", "apk", "cluster.local"),
	}
	envoyCert, _, err := newCert(envoyCertReq)
	require.NoErrorf(t, err, "Failed to generate Envoy cert")

	envoyRateLimitCertReq := &certificateRequest{
		caCertPEM:  caCert,
		caKeyPEM:   caKey,
		expiry:     expiry,
		commonName: "envoy",
		altNames:   kubeServiceNames("envoy", "apk", "cluster.local"),
	}

	envoyRateLimitCert, _, err := newCert(envoyRateLimitCertReq)
	require.NoErrorf(t, err, "Failed to generate Envoy Rate Limit Client cert")

	tests := []struct {
		name    string
		cert    []byte
		dnsName string
	}{
		{
			name:    "envoy gateway cert",
			cert:    egCert,
			dnsName: "envoy-gateway",
		},
		{
			name:    "envoy cert",
			cert:    envoyCert,
			dnsName: "envoy",
		},
		{
			name:    "envoy rate limit client cert",
			cert:    envoyRateLimitCert,
			dnsName: "envoy",
		},
	}

	for i := range tests {
		tc := tests[i]
		t.Run(tc.name, func(t *testing.T) {
			err := verifyCert(tc.cert, roots, tc.dnsName, now)
			require.NoErrorf(t, err, "Validating %s failed", tc.name)
		})
	}

}

func verifyCert(certPEM []byte, roots *x509.CertPool, dnsname string, currentTime time.Time) error {
	block, _ := pem.Decode(certPEM)
	if block == nil {
		return fmt.Errorf("failed to decode %s certificate from PEM form", dnsname)
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return err
	}

	opts := x509.VerifyOptions{
		DNSName:     dnsname,
		Roots:       roots,
		CurrentTime: currentTime,
	}
	if _, err = cert.Verify(opts); err != nil {
		return fmt.Errorf("certificate verification failed: %w", err)
	}

	return nil
}

func TestGenerateHMACSecret(t *testing.T) {
	bytes, _ := generateHMACSecret()
	encodedSecret := base64.StdEncoding.EncodeToString(bytes)
	fmt.Println("Base64 encoded secret:", encodedSecret)
}
