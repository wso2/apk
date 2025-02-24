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
	apiStore, configStore, jwtIssuerDatastore, modelBasedRoundRobinTracker := xds.CreateXDSClients(cfg)
	// NewJWTTransformer creates a new instance of JWTTransformer.
	jwtTransformer := transformer.NewJWTTransformer(cfg, jwtIssuerDatastore)
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
	go extproc.StartExternalProcessingServer(cfg, apiStore, subAppDatastore, jwtTransformer, modelBasedRoundRobinTracker, revokedJTIStore)
	go jwtbackend.StartJWKSServer(cfg)
	// Wait for the config to be loaded
	cfg.Logger.Sugar().Debug("Waiting for the config to be loaded")
	<-configStore.Notify
	cfg.Logger.Info("Config loaded successfully")
	if len(configStore.GetConfigs()) > 0 && configStore.GetConfigs()[0].Analytics != nil && configStore.GetConfigs()[0].Analytics.Enabled {
		// Start the access log service server
		go grpc.StartAccessLogServiceServer(cfg, configStore)
	}
	// Start the metrics server
	if cfg.Metrics.Enabled && strings.EqualFold(cfg.Metrics.Type, "prometheus") {
		metrics.RegisterDataSources(jwtTransformer, subAppDatastore)
		go metrics.StartPrometheusMetricsServer(cfg.Metrics.Port)
	}
	// Wait forever
	select {}
}
