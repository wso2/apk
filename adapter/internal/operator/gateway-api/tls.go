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

package gatewayapi

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	v1 "sigs.k8s.io/gateway-api/apis/v1"
)

// validateTLSSecretData ensures the cert and key provided in a secret
// is not malformed and can be properly parsed
func validateTLSSecretsData(secrets []*corev1.Secret, host *v1.Hostname) error {
	var publicKeyAlgorithm string
	var parseErr error

	pkaSecretSet := make(map[string][]string)
	for _, secret := range secrets {
		certData := secret.Data[corev1.TLSCertKey]

		if err := validateCertificate(certData); err != nil {
			return fmt.Errorf("%s/%s must contain valid %s and %s, unable to validate certificate in %s: %w", secret.Namespace, secret.Name, corev1.TLSCertKey, corev1.TLSPrivateKeyKey, corev1.TLSCertKey, err)
		}

		certBlock, _ := pem.Decode(certData)
		if certBlock == nil {
			return fmt.Errorf("%s/%s must contain valid %s and %s, unable to decode pem data in %s", secret.Namespace, secret.Name, corev1.TLSCertKey, corev1.TLSPrivateKeyKey, corev1.TLSCertKey)
		}

		cert, err := x509.ParseCertificate(certBlock.Bytes)
		if err != nil {
			return fmt.Errorf("%s/%s must contain valid %s and %s, unable to parse certificate in %s: %w", secret.Namespace, secret.Name, corev1.TLSCertKey, corev1.TLSPrivateKeyKey, corev1.TLSCertKey, err)
		}
		publicKeyAlgorithm = cert.PublicKeyAlgorithm.String()

		keyData := secret.Data[corev1.TLSPrivateKeyKey]

		keyBlock, _ := pem.Decode(keyData)
		if keyBlock == nil {
			return fmt.Errorf("%s/%s must contain valid %s and %s, unable to decode pem data in %s", secret.Namespace, secret.Name, corev1.TLSCertKey, corev1.TLSPrivateKeyKey, corev1.TLSPrivateKeyKey)
		}

		matchedFQDN, err := verifyHostname(cert, host)
		if err != nil {
			return fmt.Errorf("%s/%s must contain valid %s and %s, hostname %s does not match Common Name or DNS Names in the certificate %s", secret.Namespace, secret.Name, corev1.TLSCertKey, corev1.TLSPrivateKeyKey, string(*host), corev1.TLSCertKey)
		}
		pkaSecretKey := fmt.Sprintf("%s/%s", publicKeyAlgorithm, matchedFQDN)

		// Check whether the public key algorithm and matched certificate FQDN in the referenced secrets are unique.
		if matchedFQDN, ok := pkaSecretSet[pkaSecretKey]; ok {
			return fmt.Errorf("%s/%s public key algorithm must be unique, matched certificate FQDN %s has a conflicting algorithm [%s]",
				secret.Namespace, secret.Name, matchedFQDN, publicKeyAlgorithm)

		}
		pkaSecretSet[pkaSecretKey] = matchedFQDN

		switch keyBlock.Type {
		case "PRIVATE KEY":
			_, err := x509.ParsePKCS8PrivateKey(keyBlock.Bytes)
			if err != nil {
				parseErr = fmt.Errorf("%s/%s must contain valid %s and %s, unable to parse PKCS8 formatted private key in %s", secret.Namespace, secret.Name, corev1.TLSCertKey, corev1.TLSPrivateKeyKey, corev1.TLSPrivateKeyKey)
			}
		case "RSA PRIVATE KEY":
			_, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
			if err != nil {
				parseErr = fmt.Errorf("%s/%s must contain valid %s and %s, unable to parse PKCS1 formatted private key in %s", secret.Namespace, secret.Name, corev1.TLSCertKey, corev1.TLSPrivateKeyKey, corev1.TLSPrivateKeyKey)
			}
		case "EC PRIVATE KEY":
			_, err := x509.ParseECPrivateKey(keyBlock.Bytes)
			if err != nil {
				parseErr = fmt.Errorf("%s/%s must contain valid %s and %s, unable to parse EC formatted private key in %s", secret.Namespace, secret.Name, corev1.TLSCertKey, corev1.TLSPrivateKeyKey, corev1.TLSPrivateKeyKey)
			}
		default:
			return fmt.Errorf("%s/%s must contain valid %s and %s, %s key format found in %s, supported formats are PKCS1, PKCS8 or EC", secret.Namespace, secret.Name, corev1.TLSCertKey, corev1.TLSPrivateKeyKey, keyBlock.Type, corev1.TLSPrivateKeyKey)
		}
	}

	return parseErr
}

// verifyHostname checks if the listener Hostname matches any domain in the certificate, returns a list of matched hosts.
func verifyHostname(cert *x509.Certificate, host *v1.Hostname) ([]string, error) {
	var matchedHosts []string

	if len(cert.DNSNames) > 0 {
		matchedHosts = computeHosts(cert.DNSNames, host)
	} else {
		matchedHosts = computeHosts([]string{cert.Subject.CommonName}, host)
	}

	if len(matchedHosts) > 0 {
		return matchedHosts, nil
	}

	return nil, x509.HostnameError{Certificate: cert, Host: string(*host)}
}

func validateCertificate(data []byte) error {
	block, _ := pem.Decode(data)
	if block == nil {
		return fmt.Errorf("pem decode failed")
	}
	certs, err := x509.ParseCertificates(block.Bytes)
	if err != nil {
		return err
	}
	now := time.Now()
	roundedTime := now.Truncate(time.Minute)

	// TODO(amali) remove this round logic if any issue happens.
	// it's added due to issue in gateway conformance test certificate
	// Only add a minute if the current time is not exactly on the minute
	if !now.Equal(roundedTime) {
		roundedTime = roundedTime.Add(time.Minute)
	}
	for _, cert := range certs {
		if roundedTime.After(cert.NotAfter) {
			return fmt.Errorf("certificate is expired %v, now: %v", cert.NotAfter, roundedTime)
		}
		if roundedTime.Before(cert.NotBefore) {
			return fmt.Errorf("certificate is not yet valid %v, now %v", cert.NotBefore, roundedTime)
		}
	}
	return nil
}
