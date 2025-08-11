/* Copyright (c) 2025 WSO2 LLC. (http://www.wso2.com) All Rights Reserved.
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package main

import (

	// "strings"

	"github.com/wso2/apk/config-deployer-service-go/internal/api/routes"
	"github.com/wso2/apk/config-deployer-service-go/internal/config"
	"github.com/wso2/apk/config-deployer-service-go/internal/logging"
	"github.com/wso2/apk/config-deployer-service-go/internal/services/validators"
	"os"
	// "github.com/wso2/apk/gateway/enforcer/internal/datastore"
	// "github.com/wso2/apk/gateway/enforcer/internal/extproc"
	// "github.com/wso2/apk/gateway/enforcer/internal/grpc"
	// "github.com/wso2/apk/gateway/enforcer/internal/jwtbackend"
	// metrics "github.com/wso2/apk/gateway/enforcer/internal/metrics"
	// "github.com/wso2/apk/gateway/enforcer/internal/tokenrevocation"
	// "github.com/wso2/apk/gateway/enforcer/internal/transformer"
	// "github.com/wso2/apk/gateway/enforcer/internal/util"
	// "github.com/wso2/apk/gateway/enforcer/internal/xds"
)

func main() {
	cfg := config.GetConfig()
	logging.LoggerMain.Info("Server starting in main")

	apkSchemaLocation := "resources/conf/apk-schema.json"
	apkConfSchemaContent, err := os.ReadFile(apkSchemaLocation)
	if err != nil {
		logging.LoggerMain.Error("Failed to read APK schema file", err)
		panic(err)
	}
	validators.GlobalAPKConfValidator = validators.NewAPKConfValidator(string(apkConfSchemaContent))
	// port := cfg.CommonControllerXdsPort
	// host := cfg.CommonControllerHostname
	// clientCert, err := util.LoadCertificates(cfg.EnforcerPublicKeyPath, cfg.EnforcerPrivateKeyPath)
	// if err != nil {
	// 	panic(err)
	// }

	// // Load the trusted CA certificates
	// certPool, err := util.LoadCACertificates(cfg.TrustedAdapterCertsPath)
	// if err != nil {
	// 	panic(err)
	// }

	// //Create the TLS configuration
	// tlsConfig := util.CreateTLSConfig(clientCert, certPool)
	// subAppDatastore := datastore.NewSubAppDataStore(cfg)
	// client := grpc.NewEventingGRPCClient(host, port, cfg.XdsMaxRetries, time.Duration(cfg.XdsRetryPeriod)*time.Millisecond, tlsConfig, cfg, subAppDatastore)
	// // Start the connection
	// client.InitiateEventingGRPCConnection()

	// // Create the XDS clients
	// apiStore, configStore, jwtIssuerDatastore, modelBasedRoundRobinTracker := xds.CreateXDSClients(cfg)
	// // NewJWTTransformer creates a new instance of JWTTransformer.
	// jwtTransformer := transformer.NewJWTTransformer(cfg, jwtIssuerDatastore)
	// var revokedJTIStore *datastore.RevokedJTIStore
	// if cfg.TokenRevocationEnabled {
	// 	revokedJTIStore = datastore.NewRevokedJTIStore()
	// 	revokedJTIStore.StartRevokedJTIStoreCleanup(time.Duration(cfg.RevokedTokenCleanupInterval) * time.Second)
	// 	if cfg.IsRedisTLSEnabled {
	// 		tokenrevocation.NewRevokedTokenFetcher(cfg, revokedJTIStore, tlsConfig).Start()
	// 	} else {
	// 		tokenrevocation.NewRevokedTokenFetcher(cfg, revokedJTIStore, nil).Start()
	// 	}
	// }
	// // Start the external processing server
	// jwtbackend.JWKKEy, err = jwtbackend.ReadAndConvertToJwks(cfg)
	// if err != nil {
	// 	cfg.Logger.Sugar().Errorf("Failed to generate JWKS: %v", err)
	// }
	// go extproc.StartExternalProcessingServer(cfg, apiStore, subAppDatastore, jwtTransformer, modelBasedRoundRobinTracker, revokedJTIStore)
	go routes.StartArtifactGeneratorServer(cfg)
	// Wait for the config to be loaded
	logging.LoggerMain.Debug("Waiting for the config to be loaded")
	// <-configStore.Notify
	logging.LoggerMain.Info("Config loaded successfully")
	// if len(configStore.GetConfigs()) > 0 && configStore.GetConfigs()[0].Analytics != nil && configStore.GetConfigs()[0].Analytics.Enabled {
	// 	// Start the access log service server
	// 	go grpc.StartAccessLogServiceServer(cfg, configStore)
	// }
	// // Start the metrics server
	// if cfg.Metrics.Enabled && strings.EqualFold(cfg.Metrics.Type, "prometheus") {
	// 	metrics.RegisterDataSources(jwtTransformer, subAppDatastore)
	// 	go metrics.StartPrometheusMetricsServer(cfg.Metrics.Port)
	// }
	// Wait forever
	select {}
}
