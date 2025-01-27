package main

import (
	"time"

	"github.com/wso2/apk/gateway/enforcer/internal/config"
	"github.com/wso2/apk/gateway/enforcer/internal/datastore"
	"github.com/wso2/apk/gateway/enforcer/internal/extproc"
	"github.com/wso2/apk/gateway/enforcer/internal/grpc"
	"github.com/wso2/apk/gateway/enforcer/internal/util"
	"github.com/wso2/apk/gateway/enforcer/internal/xds"
)

func main() {
	cfg := config.GetConfig()
	port := cfg.CommonControllerXdsPort
	host := cfg.CommonControllerHostname
	clientCert, err := util.LoadCertificates(cfg.EnforcerPublicKeyPath, cfg.EnforcerPrivateKeyPath)
	if err != nil {
		panic(err)
	}

	// Load the trusted CA certificates
	certPool, err := util.LoadCACertificates(cfg.TrustedAdapterCertsPath)
	if err != nil {
		panic(err)
	}

	//Create the TLS configuration
	tlsConfig := util.CreateTLSConfig(clientCert, certPool)
	subAppDatastore := datastore.NewSubAppDataStore(cfg)
	client := grpc.NewEventingGRPCClient(host, port, cfg.XdsMaxRetries, time.Duration(cfg.XdsRetryPeriod)*time.Second, tlsConfig, cfg, subAppDatastore)
	// Start the connection
	client.InitiateEventingGRPCConnection()

	// Create the XDS clients
	apiStore, _, _, modelBasedRoundRobinTracker := xds.CreateXDSClients(cfg)

	// Start the external processing server
	go extproc.StartExternalProcessingServer(cfg, apiStore, subAppDatastore, modelBasedRoundRobinTracker)

	// Start the access log service server
	go grpc.StartAccessLogServiceServer(cfg)

	// Wait forever
	select {}
}
