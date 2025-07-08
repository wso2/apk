package main

import (
	"time"

	"strings"

	"github.com/wso2/apk/gateway/enforcer/internal/config"
	"github.com/wso2/apk/gateway/enforcer/internal/datastore"
	"github.com/wso2/apk/gateway/enforcer/internal/extproc"
	"github.com/wso2/apk/gateway/enforcer/internal/grpc"
	"github.com/wso2/apk/gateway/enforcer/internal/jwtbackend"
	metrics "github.com/wso2/apk/gateway/enforcer/internal/metrics"
	"github.com/wso2/apk/gateway/enforcer/internal/tokenrevocation"
	"github.com/wso2/apk/gateway/enforcer/internal/util"
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
	routePolicyAndMetadataDS := datastore.NewRoutePolicyAndMetadataDataStore(cfg)
	client := grpc.NewEventingGRPCClient(host, port, cfg.XdsMaxRetries, time.Duration(cfg.XdsRetryPeriod)*time.Millisecond, tlsConfig, cfg, subAppDatastore, routePolicyAndMetadataDS)
	// Start the connection
	client.InitiateEventingGRPCConnection()

	
	var revokedJTIStore *datastore.RevokedJTIStore
	if cfg.TokenRevocationEnabled {
		revokedJTIStore = datastore.NewRevokedJTIStore()
		revokedJTIStore.StartRevokedJTIStoreCleanup(time.Duration(cfg.RevokedTokenCleanupInterval) * time.Second)
		if cfg.IsRedisTLSEnabled {
			tokenrevocation.NewRevokedTokenFetcher(cfg, revokedJTIStore, tlsConfig).Start()
		} else {
			tokenrevocation.NewRevokedTokenFetcher(cfg, revokedJTIStore, nil).Start()
		}
	}
	// Start the external processing server
	jwtbackend.JWKKEy, err = jwtbackend.ReadAndConvertToJwks(cfg)
	if err != nil {
		cfg.Logger.Sugar().Errorf("Failed to generate JWKS: %v", err)
	}
	go extproc.StartExternalProcessingServer(cfg, subAppDatastore, routePolicyAndMetadataDS, revokedJTIStore)
	go jwtbackend.StartJWKSServer(cfg)
	
	
	// Start the metrics server
	if cfg.Metrics.Enabled && strings.EqualFold(cfg.Metrics.Type, "prometheus") {
		metrics.RegisterDataSources(subAppDatastore)
		go metrics.StartPrometheusMetricsServer(cfg.Metrics.Port)
	}
	// Wait forever
	select {}
}
