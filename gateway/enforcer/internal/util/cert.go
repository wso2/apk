package util

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
)

func LoadCertificates(publicKeyPath, privateKeyPath string) (tls.Certificate, error) {
	clientCert, err := tls.LoadX509KeyPair(publicKeyPath, privateKeyPath)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("failed to load client certificate and key: %v", err)
	}
	return clientCert, nil
}

func LoadCACertificates(trustedCertsPath string) (*x509.CertPool, error) {
	caCert, err := ioutil.ReadFile(trustedCertsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA certificate: %v", err)
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("failed to append CA certificate")
	}

	return certPool, nil
}

func CreateTLSConfig(cert tls.Certificate, certPool *x509.CertPool) *tls.Config {
	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      certPool,
	}
}
