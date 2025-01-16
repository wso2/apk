/*
 *  Copyright (c) 2025, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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
 */

package xds

import (
	"time"

	core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	v3 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"github.com/wso2/apk/gateway/enforcer/internal/config"
	"github.com/wso2/apk/gateway/enforcer/internal/datastore"
	"github.com/wso2/apk/gateway/enforcer/internal/util"
	status "google.golang.org/genproto/googleapis/rpc/status"
	structpb "google.golang.org/protobuf/types/known/structpb"
)

// CreateXDSClients initializes and establishes connections for multiple XDS clients,
// including API XDS, Config XDS, and JWT Issuer XDS clients.
// It handles TLS configuration, certificate loading, and connection setup.
func CreateXDSClients(cfg *config.Server) (*datastore.APIStore, *datastore.ConfigStore, *datastore.JWTIssuerStore) {
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
	apiDatastore := datastore.NewAPIStore()
	configDatastore := datastore.NewConfigStore()
	jwtIssuerDatastore := datastore.NewJWTIssuerStore()
	apiXDSClient := NewAPIXDSClient(cfg.AdapterHost, cfg.AdapterXdsPort, cfg.XdsMaxRetries, time.Duration(cfg.XdsRetryPeriod)*time.Millisecond, tlsConfig, cfg, apiDatastore)
	configXDSClient := NewXDSConfigClient(cfg.AdapterHost, cfg.AdapterXdsPort, cfg.XdsMaxRetries, time.Duration(cfg.XdsRetryPeriod)*time.Millisecond, tlsConfig, cfg, configDatastore)
	jwtIssuerXDSClient := NewJWTIssuerXDSClient(cfg.AdapterHost, cfg.AdapterXdsPort, cfg.XdsMaxRetries, time.Duration(cfg.XdsRetryPeriod)*time.Millisecond, tlsConfig, cfg, jwtIssuerDatastore)

	apiXDSClient.InitiateAPIXDSConnection()
	configXDSClient.InitiateConfigXDSConnection()
	jwtIssuerXDSClient.InitiateSubscriptionXDSConnection()
	cfg.Logger.Info("XDS clients initiated successfully")
	return apiDatastore, configDatastore, jwtIssuerDatastore
}

// CreateNode creates a new Node object with the given node ID and instance identifier.
func CreateNode(nodeID string, instanceIdentifier string) *core_v3.Node {
	fields := map[string]*structpb.Value{
		"instanceIdentifier": structpb.NewStringValue(instanceIdentifier),
	}
	return &core_v3.Node{
		Id:       nodeID,
		Metadata: &structpb.Struct{Fields: fields},
	}
}

// DiscoveryRequestForNode creates a new DiscoveryRequest for the given parameters.
func DiscoveryRequestForNode(node *core_v3.Node, versionInfo, nonce string, errorDetail *status.Status, typedURL string) *v3.DiscoveryRequest {
	return &v3.DiscoveryRequest{
		Node:          node,
		TypeUrl:       typedURL,
		VersionInfo:   versionInfo,
		ResponseNonce: nonce,
		ErrorDetail:   errorDetail,
	}
}
