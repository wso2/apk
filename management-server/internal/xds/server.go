/*
 *  Copyright (c) 2022, WSO2 LLC. (http://www.wso2.org).
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
	"context"
	"crypto/tls"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/wso2/apk/management-server/internal/database"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	apkmgt_application "github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/apkmgt"
	apkmgt_service "github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/service/apkmgt"
	wso2_cache "github.com/wso2/apk/adapter/pkg/discovery/protocol/cache/v3"
	wso2_resource "github.com/wso2/apk/adapter/pkg/discovery/protocol/resource/v3"
	wso2_server "github.com/wso2/apk/adapter/pkg/discovery/protocol/server/v3"
	"github.com/wso2/apk/adapter/pkg/logging"
	"github.com/wso2/apk/adapter/pkg/utils/tlsutils"
	"github.com/wso2/apk/management-server/internal/config"
	"github.com/wso2/apk/management-server/internal/logger"
	internal_types "github.com/wso2/apk/management-server/internal/types"
	"github.com/wso2/apk/management-server/internal/utils"
	"github.com/wso2/apk/management-server/internal/xds/callbacks"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	applicationCache      wso2_cache.SnapshotCache
	applicationCacheMutex sync.Mutex
	introducedLabels      map[string]bool
)

const (
	maxRandomInt             int    = 999999999
	typeURL                  string = "wso2.discovery.apkmgt.Application"
	grpcMaxConcurrentStreams        = 1000000
)

// IDHash uses ID field as the node hash.
type IDHash struct{}

// ID uses the node ID field
func (IDHash) ID(node *corev3.Node) string {
	if node == nil {
		return "unknown"
	}
	return node.Id
}

var _ wso2_cache.NodeHash = IDHash{}

func init() {
	applicationCache = wso2_cache.NewSnapshotCache(false, IDHash{}, nil)
	rand.Seed(time.Now().UnixNano())
	introducedLabels = make(map[string]bool, 1)
}

// FeedData mock data
func FeedData() {
	config := config.ReadConfigs()
	logger.LoggerXdsServer.Debug("adding mock data")
	version := rand.Intn(maxRandomInt)
	applications := apkmgt_application.Application{
		Uuid: "apiUUID1",
		Name: "name1",
	}
	newSnapshot, _ := wso2_cache.NewSnapshot(fmt.Sprint(version), map[wso2_resource.Type][]types.Resource{
		wso2_resource.APKMgtApplicationType: {&applications},
	})
	applicationCacheMutex.Lock()
	defer applicationCacheMutex.Unlock()
	applicationCache.SetSnapshot(context.Background(), config.ManagementServer.NodeLabels[0], newSnapshot)
}

// AddSingleApplication will update the Application specified by the UUID to the xds cache
func AddSingleApplication(label string, application *apkmgt_application.Application) {
	var newSnapshot wso2_cache.Snapshot
	version := rand.Intn(maxRandomInt)
	//application, errDb := database.GetApplicationByUUID(appUUID)
	currentSnapshot, err := applicationCache.GetSnapshot(label)

	// error occurs if no snapshot is under the provided label
	if err != nil {
		newSnapshot, _ = wso2_cache.NewSnapshot(fmt.Sprint(version), map[wso2_resource.Type][]types.Resource{
			wso2_resource.APKMgtApplicationType: {application},
		})
	} else {
		resourceMap := currentSnapshot.GetResourcesAndTTL(typeURL)
		resourceMap[application.Uuid] = types.ResourceWithTTL{
			Resource: application,
		}
		applicationResources := convertResourceMapToArray(resourceMap)
		newSnapshot, _ = wso2_cache.NewSnapshot(fmt.Sprint(version), map[wso2_resource.Type][]types.Resource{
			wso2_resource.APKMgtApplicationType: applicationResources,
		})
	}
	applicationCacheMutex.Lock()
	defer applicationCacheMutex.Unlock()
	applicationCache.SetSnapshot(context.Background(), label, newSnapshot)
	introducedLabels[label] = true
	logger.LoggerXds.Infof("Application Snapshot is updated for label %s with the version %d. New snapshot "+
		"size is %d.", label, version, len(newSnapshot.GetResourcesAndTTL(typeURL)))

}

// RemoveApplication removes the Application entry from XDS cache
func RemoveApplication(label, appUUID string) {
	var newSnapshot wso2_cache.Snapshot
	version := rand.Intn(maxRandomInt)
	for l := range introducedLabels {
		// If the label does not match with any introduced labels, don't need to delete the application from cache.
		if !strings.EqualFold(label, l) {
			continue
		}
		currentSnapshot, err := applicationCache.GetSnapshot(label)
		if err != nil {
			continue
		}

		resourceMap := currentSnapshot.GetResourcesAndTTL(typeURL)
		_, apiFound := resourceMap[appUUID]
		// If the Application is found, then the xds cache is updated and returned.
		if apiFound {
			logger.LoggerXds.Debugf("Application : %s is found within snapshot for label %s", appUUID, label)
			delete(resourceMap, appUUID)
			apiResources := convertResourceMapToArray(resourceMap)
			newSnapshot, _ = wso2_cache.NewSnapshot(fmt.Sprint(version), map[wso2_resource.Type][]types.Resource{
				wso2_resource.APKMgtApplicationType: apiResources,
			})
			applicationCacheMutex.Lock()
			defer applicationCacheMutex.Unlock()
			applicationCache.SetSnapshot(context.Background(), label, newSnapshot)
			logger.LoggerXds.Infof("API Snaphsot is updated for label %s with the version %d. New snapshot "+
				"size is %d.", label, version, len(newSnapshot.GetResourcesAndTTL(typeURL)))
			return
		}
	}
	logger.LoggerXds.Errorf("Application : %s is not found within snapshot for label %s", appUUID, label)
}

// AddMultipleApplications adds the applications specified in applicationEventArray to the xds cache
// This will ideally be used to populate all applications in the startup of the mgt server.
func AddMultipleApplications(applicationEventArray []*internal_types.ApplicationEvent) {
	snapshotMap := make(map[string]*wso2_cache.Snapshot)
	version := rand.Intn(maxRandomInt)

	for _, event := range applicationEventArray {
		label := event.Label
		appUUID := event.UUID

		application, err := database.GetApplicationByUUID(appUUID)
		if err != nil {
			logger.LoggerDatabase.ErrorC(logging.ErrorDetails{
				Message: fmt.Sprintf("Error retrieving application for uuid : %s from database error: %v, "+
					"hence skipping add to xdx cache", appUUID, err),
				Severity:  logging.MINOR,
				ErrorCode: 1101,
			})
			continue
		}

		snapshotEntry, snapshotFound := snapshotMap[label]
		var newSnapshot wso2_cache.Snapshot

		if !snapshotFound {
			newSnapshot, _ = wso2_cache.NewSnapshot(fmt.Sprint(version), map[wso2_resource.Type][]types.Resource{
				wso2_resource.APKMgtApplicationType: {application},
			})
			snapshotEntry = &newSnapshot
			snapshotMap[label] = &newSnapshot
		} else {
			// error occurs if no snapshot is under the provided label
			resourceMap := snapshotEntry.GetResourcesAndTTL(typeURL)
			resourceMap[appUUID] = types.ResourceWithTTL{
				Resource: application,
			}
			appResources := convertResourceMapToArray(resourceMap)
			newSnapshot, _ = wso2_cache.NewSnapshot(fmt.Sprint(version), map[wso2_resource.Type][]types.Resource{
				wso2_resource.APKMgtApplicationType: appResources,
			})
			snapshotMap[label] = &newSnapshot
		}
	}
	applicationCacheMutex.Lock()
	defer applicationCacheMutex.Unlock()
	for label, snapshotEntry := range snapshotMap {
		applicationCache.SetSnapshot(context.Background(), label, *snapshotEntry)
		introducedLabels[label] = true
		logger.LoggerXds.Infof("Application Snaphsot is updated for label %s with the version %d.", label, version)
	}
}

func convertResourceMapToArray(resourceMap map[string]types.ResourceWithTTL) []types.Resource {
	var appResources []types.Resource
	for _, res := range resourceMap {
		appResources = append(appResources, res.Resource)
	}
	return appResources
}

// SetEmptySnapshot sets an empty snapshot into the applicationCache for the given label
// this is used to set empty snapshot when there are no Applications available for a label
func SetEmptySnapshot(label string) error {
	version := rand.Intn(maxRandomInt)
	newSnapshot, err := wso2_cache.NewSnapshot(fmt.Sprint(version), map[wso2_resource.Type][]types.Resource{
		wso2_resource.APKMgtApplicationType: {},
	})
	if err != nil {
		logger.LoggerXds.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprintf("Error creating empty snapshot. error: %v", err.Error()),
			Severity:  logging.CRITICAL,
			ErrorCode: 1003,
		})
		return err
	}
	applicationCacheMutex.Lock()
	defer applicationCacheMutex.Unlock()
	//performing null check again to avoid race conditions
	_, errSnap := applicationCache.GetSnapshot(label)
	if errSnap != nil && strings.Contains(errSnap.Error(), "no snapshot found for node") {
		errSetSnap := applicationCache.SetSnapshot(context.Background(), label, newSnapshot)
		if errSetSnap != nil {
			logger.LoggerXds.ErrorC(logging.ErrorDetails{
				Message:   fmt.Sprintf("Error setting empty snapshot to applicationCache. error : %v", errSetSnap.Error()),
				Severity:  logging.CRITICAL,
				ErrorCode: 1004,
			})
			return errSetSnap
		}
	}
	return nil
}

// InitAPKMgtServer initializes the APK management server
func InitAPKMgtServer() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	apkMgtAPIDsSrv := wso2_server.NewServer(ctx, applicationCache, &callbacks.Callbacks{})
	publicKeyLocation, privateKeyLocation, truststoreLocation := utils.GetKeyLocations()
	cert, err := tlsutils.GetServerCertificate(publicKeyLocation, privateKeyLocation)
	if err != nil {
		logger.LoggerMGTServer.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprintf("Failed to initiate the ssl context, error: %v", err.Error()),
			Severity:  logging.BLOCKER,
			ErrorCode: 1200,
		})
	}
	caCertPool := tlsutils.GetTrustedCertPool(truststoreLocation)
	var grpcOptions []grpc.ServerOption
	grpcOptions = append(grpcOptions, grpc.MaxConcurrentStreams(grpcMaxConcurrentStreams), grpc.Creds(
		credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{cert},
			ClientAuth:   tls.RequireAndVerifyClientCert,
			ClientCAs:    caCertPool,
		}),
	))
	grpcServer := grpc.NewServer(grpcOptions...)
	apkmgt_service.RegisterAPKMgtDiscoveryServiceServer(grpcServer, apkMgtAPIDsSrv)
	config := config.ReadConfigs()
	port := config.ManagementServer.XDSPort

	//todo (amaliMatharaarachchi) handle error gracefully
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logger.LoggerServer.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprintf("Error while listening on port: %v. Error: %v", port, err.Error()),
			Severity:  logging.BLOCKER,
			ErrorCode: 1000,
		})
	}

	logger.LoggerServer.Infof("APK Management server XDS is starting on port %v.", port)
	if err = grpcServer.Serve(listener); err != nil {
		logger.LoggerServer.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprint("Error while starting APK Management server XDS server."),
			Severity:  logging.BLOCKER,
			ErrorCode: 1001,
		})
	}
}
