/*
 *  Copyright (c) 2024, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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
 *  This file contains code derived from Envoy Gateway,
 *  https://github.com/envoyproxy/gateway
 *  and is provided here subject to the following:
 *  Copyright Project Envoy Gateway Authors
 *
 */

package runner

import (
	"context"

	"github.com/wso2/apk/adapter/internal/loggers"
	gatewayapi "github.com/wso2/apk/adapter/internal/operator/gateway-api"
	"github.com/wso2/apk/adapter/internal/operator/message"
	"github.com/wso2/apk/adapter/internal/operator/utils"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
	v1 "sigs.k8s.io/gateway-api/apis/v1"
)

type Config struct {
	ProviderResources *message.ProviderResources
	XdsIR             *message.XdsIR
	InfraIR           *message.InfraIR
}

type Runner struct {
	Config
}

func New(cfg *Config) *Runner {
	return &Runner{Config: *cfg}
}

func (r *Runner) Name() string {
	return "gateway-api"
}

// Start starts the gateway-api translator runner
func (r *Runner) Start(ctx context.Context) (err error) {
	go r.subscribeAndTranslate(ctx)
	return
}

func (r *Runner) subscribeAndTranslate(ctx context.Context) {
	message.HandleSubscription(message.Metadata{Runner: "gateway-api", Message: "provider-resources"}, r.ProviderResources.GatewayAPIResources.Subscribe(ctx),
		func(update message.Update[string, *gatewayapi.ControllerResources], errChan chan error) {
			val := update.Value
			loggers.LoggerAPKOperator.Info("Received an update in provider resources ...")
			// There is only 1 key which is the controller name
			// so when a delete is triggered, delete all IR keys
			if update.Delete || val == nil {
				r.deleteAllIRKeys()
				r.deleteAllStatusKeys()
				return
			}

			// IR keys for watchable
			var curIRKeys, newIRKeys []string

			// Get current IR keys
			for key := range r.InfraIR.LoadAll() {
				curIRKeys = append(curIRKeys, key)
			}

			// Get all status keys from watchable and save them in this StatusesToDelete structure.
			// Iterating through the controller resources, any valid keys will be removed from statusesToDelete.
			// Remaining keys will be deleted from watchable before we exit this function.
			statusesToDelete := r.getAllStatuses()

			for _, resources := range *val {
				// Translate and publish IRs.
				t := &gatewayapi.Translator{
					GatewayControllerName: "wso2.com/apk-gateway-default",
					GatewayClassName:      v1.ObjectName(resources.GatewayClass.Name),
				}

				// Translate to IR
				result := t.Translate(resources)

				// Publish the IRs.
				// Also validate the ir before sending it.
				for key, val := range result.InfraIR {
					if err := val.Validate(); err != nil {
						loggers.LoggerAPKOperator.Error(err, "unable to validate infra ir, skipped sending it")
						errChan <- err
					} else {
						r.InfraIR.Store(key, val)
						newIRKeys = append(newIRKeys, key)
					}
				}

				for key, val := range result.XdsIR {
					if err := val.Validate(); err != nil {
						loggers.LoggerAPKOperator.Error(err, "unable to validate xds ir, skipped sending it")
						errChan <- err
					} else {
						r.XdsIR.Store(key, val)
					}
				}

				// Update Status
				for _, gateway := range result.Gateways {
					gateway := gateway
					key := utils.NamespacedName(gateway)
					r.ProviderResources.GatewayStatuses.Store(key, &gateway.Status)
					delete(statusesToDelete.GatewayStatusKeys, key)
				}
				for _, httpRoute := range result.HTTPRoutes {
					httpRoute := httpRoute
					key := utils.NamespacedName(httpRoute)
					r.ProviderResources.HTTPRouteStatuses.Store(key, &httpRoute.Status)
					delete(statusesToDelete.HTTPRouteStatusKeys, key)
				}
				// 	for _, grpcRoute := range result.GRPCRoutes {
				// 		grpcRoute := grpcRoute
				// 		key := utils.NamespacedName(grpcRoute)
				// 		r.ProviderResources.GRPCRouteStatuses.Store(key, &grpcRoute.Status)
				// 		delete(statusesToDelete.GRPCRouteStatusKeys, key)
				// 	}
				// 	for _, tlsRoute := range result.TLSRoutes {
				// 		tlsRoute := tlsRoute
				// 		key := utils.NamespacedName(tlsRoute)
				// 		r.ProviderResources.TLSRouteStatuses.Store(key, &tlsRoute.Status)
				// 		delete(statusesToDelete.TLSRouteStatusKeys, key)
				// 	}
				// 	for _, tcpRoute := range result.TCPRoutes {
				// 		tcpRoute := tcpRoute
				// 		key := utils.NamespacedName(tcpRoute)
				// 		r.ProviderResources.TCPRouteStatuses.Store(key, &tcpRoute.Status)
				// 		delete(statusesToDelete.TCPRouteStatusKeys, key)
				// 	}
				// 	for _, udpRoute := range result.UDPRoutes {
				// 		udpRoute := udpRoute
				// 		key := utils.NamespacedName(udpRoute)
				// 		r.ProviderResources.UDPRouteStatuses.Store(key, &udpRoute.Status)
				// 		delete(statusesToDelete.UDPRouteStatusKeys, key)
				// 	}

				// 	// Skip updating status for policies with empty status
				// 	// They may have been skipped in this translation because
				// 	// their target is not found (not relevant)

				// 	for _, backendTLSPolicy := range result.BackendTLSPolicies {
				// 		backendTLSPolicy := backendTLSPolicy
				// 		key := utils.NamespacedName(backendTLSPolicy)
				// 		if !(reflect.ValueOf(backendTLSPolicy.Status).IsZero()) {
				// 			r.ProviderResources.BackendTLSPolicyStatuses.Store(key, &backendTLSPolicy.Status)
				// 		}
				// 		delete(statusesToDelete.BackendTLSPolicyStatusKeys, key)
				// 	}

			}

			// Delete IR keys
			// There is a 1:1 mapping between infra and xds IR keys
			delKeys := getIRKeysToDelete(curIRKeys, newIRKeys)
			for _, key := range delKeys {
				r.InfraIR.Delete(key)
				r.XdsIR.Delete(key)
			}

			// Delete status keys
			r.deleteStatusKeys(statusesToDelete)
		},
	)
	loggers.LoggerAPKOperator.Info("shutting down")
}

type StatusesToDelete struct {
	GatewayStatusKeys   map[types.NamespacedName]bool
	HTTPRouteStatusKeys map[types.NamespacedName]bool
}

func (r *Runner) getAllStatuses() *StatusesToDelete {
	// Maps storing status keys to be deleted
	ds := &StatusesToDelete{
		GatewayStatusKeys:   make(map[types.NamespacedName]bool),
		HTTPRouteStatusKeys: make(map[types.NamespacedName]bool),
	}

	// Get current status keys
	for key := range r.ProviderResources.GatewayStatuses.LoadAll() {
		ds.GatewayStatusKeys[key] = true
	}
	for key := range r.ProviderResources.HTTPRouteStatuses.LoadAll() {
		ds.HTTPRouteStatusKeys[key] = true
	}

	return ds
}

func (r *Runner) deleteStatusKeys(ds *StatusesToDelete) {
	for key := range ds.GatewayStatusKeys {
		r.ProviderResources.GatewayStatuses.Delete(key)
		delete(ds.GatewayStatusKeys, key)
	}
	for key := range ds.HTTPRouteStatusKeys {
		r.ProviderResources.HTTPRouteStatuses.Delete(key)
		delete(ds.HTTPRouteStatusKeys, key)
	}
}

// deleteAllIRKeys deletes all XdsIR and InfraIR
func (r *Runner) deleteAllIRKeys() {
	for key := range r.InfraIR.LoadAll() {
		r.InfraIR.Delete(key)
		r.XdsIR.Delete(key)
	}
}

// deleteAllStatusKeys deletes all status keys stored by the subscriber.
func (r *Runner) deleteAllStatusKeys() {
	// Fields of GatewayAPIStatuses
	for key := range r.ProviderResources.GatewayStatuses.LoadAll() {
		r.ProviderResources.GatewayStatuses.Delete(key)
	}
	for key := range r.ProviderResources.HTTPRouteStatuses.LoadAll() {
		r.ProviderResources.HTTPRouteStatuses.Delete(key)
	}
}

// getIRKeysToDelete returns the list of IR keys to delete
// based on the difference between the current keys and the
// new keys parameters passed to the function.
func getIRKeysToDelete(curKeys, newKeys []string) []string {
	curSet := sets.NewString(curKeys...)
	newSet := sets.NewString(newKeys...)

	delSet := curSet.Difference(newSet)

	return delSet.List()
}
