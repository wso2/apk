package xds

import (
	"github.com/wso2/apk/gateway/enforcer/internal/config"
	"github.com/wso2/apk/gateway/enforcer/internal/util"
	"time"
)

func CreateXDSClients(cfg *config.Server) {
	clientCert, err := util.LoadCertificates(cfg.EnforcerPublicKeyPath, cfg.EnforcerPrivateKeyPath)
	if err != nil {
		panic(err)
	}

	// Load the trusted CA certificates
	certPoll, err := util.LoadCACertificates(cfg.TrustedAdapterCertsPath)
	if err != nil {
		panic(err)
	}
	
	// Create the TLS configuration
	tlsConfig := util.CreateTLSConfig(clientCert, certPoll)
	apiXDSClient := NewAPIXDSClient(cfg.AdapterHost, cfg.AdapterXdsPort, cfg.XdsMaxRetries,time.Duration(cfg.XdsRetryPeriod) * time.Second, tlsConfig, cfg)
	configXDSClient := NewXDSConfigClient(cfg.AdapterHost, cfg.AdapterXdsPort, cfg.XdsMaxRetries,time.Duration(cfg.XdsRetryPeriod) * time.Second, tlsConfig, cfg)
	jwtIssuerXDSClient := NewJWTIssuerXDSClient(cfg.AdapterHost, cfg.AdapterXdsPort, cfg.XdsMaxRetries,time.Duration(cfg.XdsRetryPeriod) * time.Second, tlsConfig, cfg)

	apiXDSClient.InitiateAPIXDSConnection()
	configXDSClient.InitiateConfigXDSConnection()
	jwtIssuerXDSClient.InitiateSubscriptionXDSConnection()
}