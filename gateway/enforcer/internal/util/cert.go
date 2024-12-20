package util

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
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

// LoadCACertificates loads the CA certificates from the provided file path.
// It reads the certificate file and appends it to a new CertPool. 
// If any error occurs during reading or appending, it returns an error.
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

// CreateTLSConfig creates and returns a new TLS configuration using the provided client certificate
// and CA certificate pool. It sets up the certificates for secure communication.
func CreateTLSConfig(cert tls.Certificate, certPool *x509.CertPool) *tls.Config {
	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      certPool,
	}
}
