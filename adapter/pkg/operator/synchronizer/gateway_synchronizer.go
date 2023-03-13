/*
 *  Copyright (c) 2023, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
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

package synchronizer

import (
	"github.com/wso2/apk/adapter/internal/discovery/xds"
	"github.com/wso2/apk/adapter/internal/loggers"
	"github.com/wso2/apk/adapter/pkg/logging"
	"github.com/wso2/apk/adapter/pkg/operator/constants"
	gwapiv1b1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

// GatewayEvent holds the data structure used for passing Gateway
// events from controller go routine to synchronizer
// go routine.
type GatewayEvent struct {
	EventType string
	Event     GatewayState
}

// HandleGatewayLifeCycleEvents handles the Gateway events generated from OperatorDataStore
func HandleGatewayLifeCycleEvents(ch *chan GatewayEvent) {
	loggers.LoggerAPKOperator.Info("Operator synchronizer listening for Gateway lifecycle events...")
	for event := range *ch {
		if event.Event.GatewayDefinition == nil {
			loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2628))
		}
		loggers.LoggerAPKOperator.Infof("%s event received for %v", event.EventType, event.Event.GatewayDefinition.Name)
		var err error
		switch event.EventType {
		case constants.Delete:
			err = undeployGateway(event.Event)
		case constants.Create:
			err = deployGateway(event.Event)
		case constants.Update:
			err = deployGateway(event.Event)
		}
		if err != nil {
			loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2629, event.EventType, err))
		}
	}
}

// deployGateway deploys the related Gateway in CREATE and UPDATE events.
func deployGateway(gatewayState GatewayState) error {
	var err error
	if gatewayState.GatewayDefinition != nil {
		_, err = AddGateway(gatewayState.GatewayDefinition)
	}
	return err
}

// undeployGateway undeploys the related Gateway in DELETE events.
func undeployGateway(gatewayState GatewayState) error {
	var err error
	if gatewayState.GatewayDefinition != nil {
		_, err = DeleteGateway(gatewayState.GatewayDefinition)
	}
	return err
}

// AddGateway adds a Gateway to the XDS server.
func AddGateway(gateway *gwapiv1b1.Gateway) (string, error) {
	xds.GenerateGlobalClusters(gateway.Name)
	listeners, clusters, routes, endpoints, apis := xds.GenerateEnvoyResoucesForGateway(gateway)
	loggers.LoggerAPKOperator.Infof("listeners: %v", listeners)
	loggers.LoggerAPKOperator.Infof("clusters: %v", clusters)
	loggers.LoggerAPKOperator.Infof("routes: %v", routes)
	loggers.LoggerAPKOperator.Infof("endpoints: %v", endpoints)
	loggers.LoggerAPKOperator.Infof("apis: %v", apis)
	xds.UpdateXdsCacheWithLock(gateway.Name, endpoints, clusters, routes, listeners)
	xds.UpdateEnforcerApis(gateway.Name, apis, "")
	return "", nil
}

// DeleteGateway deletes a Gateway from the XDS server.
func DeleteGateway(gateway *gwapiv1b1.Gateway) (string, error) {
	//return xds.DeleteGateway(gateway)
	return "", nil
}

// func UpdateXdsCacheWithLock(label string, endpoints []types.Resource, clusters []types.Resource, routes []types.Resource,
// 	listeners []types.Resource) bool {
// 	mutexForXdsUpdate.Lock()
// 	defer mutexForXdsUpdate.Unlock()
// 	return updateXdsCache(label, endpoints, clusters, routes, listeners)
// }

// // use UpdateXdsCacheWithLock to avoid race conditions
// func updateXdsCache(label string, endpoints []types.Resource, clusters []types.Resource, routes []types.Resource, listeners []types.Resource) bool {
// 	version := rand.Intn(maxRandomInt)
// 	// TODO: (VirajSalaka) kept same version for all the resources as we are using simple cache implementation.
// 	// Will be updated once decide to move to incremental XDS
// 	snap, errNewSnap := envoy_cachev3.NewSnapshot(fmt.Sprint(version), map[envoy_resource.Type][]types.Resource{
// 		envoy_resource.EndpointType: endpoints,
// 		envoy_resource.ClusterType:  clusters,
// 		envoy_resource.ListenerType: listeners,
// 		envoy_resource.RouteType:    routes,
// 	})
// 	if errNewSnap != nil {
// 		logger.LoggerXds.ErrorC(logging.GetErrorByCode(1413, errNewSnap.Error()))
// 		return false
// 	}
// 	snap.Consistent()
// 	//TODO: (VirajSalaka) check
// 	errSetSnap := cache.SetSnapshot(context.Background(), label, snap)
// 	if errSetSnap != nil {
// 		logger.LoggerXds.ErrorC(logging.GetErrorByCode(1414, errSetSnap.Error()))
// 		return false
// 	}
// 	logger.LoggerXds.Infof("New Router cache updated for the label: " + label + " version: " + fmt.Sprint(version))
// 	return true
// }

// // GenerateEnvoyResoucesForLabel generates envoy resources for a given label
// // This method will list out all APIs mapped to the label. and generate envoy resources for all of these APIs.
// func GenerateEnvoyResoucesForLabel(label string) ([]types.Resource, []types.Resource, []types.Resource,
// 	[]types.Resource, []types.Resource) {
// 	var clusterArray []*clusterv3.Cluster
// 	var vhostToRouteArrayMap = make(map[string][]*routev3.Route)
// 	var endpointArray []*corev3.Address
// 	var apis []types.Resource

// 	for organizationID, entityMap := range orgIDOpenAPIEnvoyMap {
// 		for apiKey, labels := range entityMap {
// 			if stringutils.StringInSlice(label, labels) {
// 				vhost, err := ExtractVhostFromAPIIdentifier(apiKey)
// 				if err != nil {
// 					logger.LoggerXds.ErrorC(logging.GetErrorByCode(1411, err.Error(), organizationID))
// 					continue
// 				}
// 				isDefaultVersion := false
// 				if enforcerAPISwagger, ok := orgIDAPIMgwSwaggerMap[organizationID][apiKey]; ok {
// 					isDefaultVersion = enforcerAPISwagger.IsDefaultVersion
// 				} else {
// 					// If the mgwSwagger is not found, proceed with other APIs. (Unreachable condition at this point)
// 					// If that happens, there is no purpose in processing clusters too.
// 					continue
// 				}
// 				// If it is a default versioned API, the routes are added to the end of the existing array.
// 				// Otherwise the routes would be added to the front.
// 				// /fooContext/2.0.0/* resource path should be matched prior to the /fooContext/* .
// 				if isDefaultVersion {
// 					vhostToRouteArrayMap[vhost] = append(vhostToRouteArrayMap[vhost], orgIDOpenAPIRoutesMap[organizationID][apiKey]...)
// 				} else {
// 					vhostToRouteArrayMap[vhost] = append(orgIDOpenAPIRoutesMap[organizationID][apiKey], vhostToRouteArrayMap[vhost]...)
// 				}
// 				clusterArray = append(clusterArray, orgIDOpenAPIClustersMap[organizationID][apiKey]...)
// 				endpointArray = append(endpointArray, orgIDOpenAPIEndpointsMap[organizationID][apiKey]...)
// 				enfocerAPI, ok := orgIDOpenAPIEnforcerApisMap[organizationID][apiKey]
// 				if ok {
// 					apis = append(apis, enfocerAPI)
// 				}
// 				// listenerArrays = append(listenerArrays, openAPIListenersMap[apiKey])
// 			}
// 		}
// 	}

// 	// If the token endpoint is enabled, the token endpoint also needs to be added.
// 	conf := config.ReadConfigs()
// 	enableJwtIssuer := conf.Enforcer.JwtIssuer.Enabled
// 	systemHost := conf.Envoy.SystemHost
// 	if enableJwtIssuer {
// 		routeToken := envoyconf.CreateTokenRoute()
// 		vhostToRouteArrayMap[systemHost] = append(vhostToRouteArrayMap[systemHost], routeToken)
// 	}

// 	// Add health endpoint
// 	routeHealth := envoyconf.CreateHealthEndpoint()
// 	vhostToRouteArrayMap[systemHost] = append(vhostToRouteArrayMap[systemHost], routeHealth)

// 	// Add the readiness endpoint. isReady flag will be set to true once all the apis are fetched from the control plane
// 	if isReady {
// 		readynessEndpoint := envoyconf.CreateReadyEndpoint()
// 		vhostToRouteArrayMap[systemHost] = append(vhostToRouteArrayMap[systemHost], readynessEndpoint)
// 	}

// 	listenerArray, listenerFound := envoyListenerConfigMap[label]
// 	routesConfig, routesConfigFound := envoyRouteConfigMap[label]
// 	if !listenerFound && !routesConfigFound {
// 		listenerArray, routesConfig = oasParser.GetProductionListenerAndRouteConfig(vhostToRouteArrayMap)
// 		envoyListenerConfigMap[label] = listenerArray
// 		envoyRouteConfigMap[label] = routesConfig
// 	} else {
// 		// If the routesConfig exists, the listener exists too
// 		oasParser.UpdateRoutesConfig(routesConfig, vhostToRouteArrayMap)
// 	}
// 	clusterArray = append(clusterArray, envoyClusterConfigMap[label]...)
// 	endpointArray = append(endpointArray, envoyEndpointConfigMap[label]...)
// 	endpoints, clusters, listeners, routeConfigs := oasParser.GetCacheResources(endpointArray, clusterArray, listenerArray, routesConfig)
// 	return endpoints, clusters, listeners, routeConfigs, apis
// }
