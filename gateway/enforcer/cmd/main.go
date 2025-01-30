package main

import (
	"time"

	"github.com/wso2/apk/gateway/enforcer/internal/config"
	"github.com/wso2/apk/gateway/enforcer/internal/datastore"
	"github.com/wso2/apk/gateway/enforcer/internal/extproc"
	"github.com/wso2/apk/gateway/enforcer/internal/grpc"
	"github.com/wso2/apk/gateway/enforcer/internal/transformer"
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
	client := grpc.NewEventingGRPCClient(host, port, cfg.XdsMaxRetries, time.Duration(cfg.XdsRetryPeriod)*time.Millisecond, tlsConfig, cfg, subAppDatastore)
	// Start the connection
	client.InitiateEventingGRPCConnection()

	// Create the XDS clients
	apiStore, configStore, jwtIssuerDatastore,modelBasedRoundRobinTracker := xds.CreateXDSClients(cfg)
	// NewJWTTransformer creates a new instance of JWTTransformer.
	jwtTransformer := transformer.NewJWTTransformer(jwtIssuerDatastore)
	// Start the external processing server
	go extproc.StartExternalProcessingServer(cfg, apiStore, subAppDatastore, jwtTransformer,modelBasedRoundRobinTracker)

	// Wait for the config to be loaded
	cfg.Logger.Info("Waiting for the config to be loaded")
	<- configStore.Notify
	cfg.Logger.Info("Config loaded successfully")
	if len(configStore.GetConfigs()) > 0 && configStore.GetConfigs()[0].Analytics != nil && configStore.GetConfigs()[0].Analytics.Enabled {
		// Start the access log service server
		go grpc.StartAccessLogServiceServer(cfg, configStore)
	}

	// Wait forever
	select {}
}
