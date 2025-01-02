package xds

import (
	"github.com/wso2/apk/gateway/enforcer/internal/config"
	"github.com/wso2/apk/gateway/enforcer/internal/util"
	"time"
)

// CreateXDSClients initializes and establishes connections for multiple XDS clients, 
// including API XDS, Config XDS, and JWT Issuer XDS clients. 
// It handles TLS configuration, certificate loading, and connection setup.
func CreateXDSClients(cfg *config.Server) {
	clientCert, err := util.LoadCertificates(cfg.EnforcerPublicKeyPath, cfg.EnforcerPrivateKeyPath)
	if err != nil {
		panic(err)
	}

	// Load the trusted CA certificates
	certPool, err := util.LoadCACertificates(cfg.TrustedAdapterCertsPath)
	if err != nil {
		panic(err)
	}
	
	// Create the TLS configuration
	tlsConfig := util.CreateTLSConfig(clientCert, certPool)
	apiXDSClient := NewAPIXDSClient(cfg.AdapterHost, cfg.AdapterXdsPort, cfg.XdsMaxRetries,time.Duration(cfg.XdsRetryPeriod) * time.Second, tlsConfig, cfg)
	configXDSClient := NewXDSConfigClient(cfg.AdapterHost, cfg.AdapterXdsPort, cfg.XdsMaxRetries,time.Duration(cfg.XdsRetryPeriod) * time.Second, tlsConfig, cfg)
	jwtIssuerXDSClient := NewJWTIssuerXDSClient(cfg.AdapterHost, cfg.AdapterXdsPort, cfg.XdsMaxRetries,time.Duration(cfg.XdsRetryPeriod) * time.Second, tlsConfig, cfg)

	apiXDSClient.InitiateAPIXDSConnection()
	configXDSClient.InitiateConfigXDSConnection()
	jwtIssuerXDSClient.InitiateSubscriptionXDSConnection()
	cfg.Logger.Info("XDS clients initiated successfully")
}