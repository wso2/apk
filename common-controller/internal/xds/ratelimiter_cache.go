/*
 *  Copyright (c) 2023, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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
 */

package xds

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"sync"

	gcp_types "github.com/envoyproxy/go-control-plane/pkg/cache/types"
	gcp_cache "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	gcp_resource "github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	rls_config "github.com/envoyproxy/go-control-plane/ratelimit/config/ratelimit/v3"
	logger "github.com/sirupsen/logrus"
	"github.com/wso2/apk/common-controller/internal/loggers"
	logging "github.com/wso2/apk/common-controller/internal/logging"
	dpv1alpha1 "github.com/wso2/apk/common-controller/internal/operator/api/v1alpha1"
)

// Constants relevant to the route related ratelimit configurations
const (
	DescriptorKeyForOrg                = "org"
	OrgMetadataKey                     = "customorg"
	DescriptorKeyForVhost              = "vhost"
	DescriptorKeyForPath               = "path"
	DescriptorKeyForMethod             = "method"
	DescriptorValueForAPIMethod        = "ALL"
	DescriptorValueForOperationMethod  = ":method"
	MetadataNamespaceForCustomPolicies = "apk.ratelimit.metadata"
	MetadataNamespaceForWSO2Policies   = "envoy.filters.http.ext_authz"
	apiDefinitionClusterName           = "api_definition_cluster"
)

// Constants relevant to the rate limit service
const (
	RateLimiterDomain                    = "Default"
	RateLimitPolicyOperationLevel string = "OPERATION"
	RateLimitPolicyAPILevel       string = "API"
)

var void struct{}

var rlsPolicyCache *rateLimitPolicyCache

type rateLimitPolicyCache struct {
	// xdsCache is the snapshot cache for the rate limiter service
	xdsCache gcp_cache.SnapshotCache

	// TODO: (renuka) move both 'apiLevelRateLimitPolicies' and 'apiLevelMu' to a new struct when doing the App level rate limiting
	// So app level rate limits are in a new struct and refer in this struct.
	// org -> vhost -> API-Identifier (i.e. Vhost:API-UUID) -> Rate Limit Configs
	apiLevelRateLimitPolicies map[string]map[string]map[string][]*rls_config.RateLimitDescriptor

	// org -> Custom Rate Limit Configs
	customRateLimitPolicies map[string]map[string]*rls_config.RateLimitDescriptor

	// mutex for API level
	apiLevelMu sync.RWMutex
}

// AddAPILevelRateLimitPolicies adds inline Rate Limit policies in APIs to be updated in the Rate Limiter service.
func (r *rateLimitPolicyCache) AddAPILevelRateLimitPolicies(vHosts []string, resolveRatelimit dpv1alpha1.ResolveRateLimitAPIPolicy) {

	rlsConfigs := []*rls_config.RateLimitDescriptor{}

	// The map apiOperations is used to keep `Pat:HTTPmethod` unique to make sure the Rate Limiter Config to be consistent (not to have duplicate rate limit policies)
	// path -> HTTP method
	apiOperations := make(map[string]map[string]struct{})
	if len(resolveRatelimit.Resources) != 0 {
		logger.Info("Going to resorcelevel")
		for _, resource := range resolveRatelimit.Resources {

			path := resolveRatelimit.Context + resolveRatelimit.Context + resource.Path
			logger.Info("path", path)
			if _, ok := apiOperations[path]; !ok {
				apiOperations[path] = make(map[string]struct{})
			}

			method := resource.Method

			rlPolicyConfig := parseRateLimitPolicyToXDS(resource.ResourceRatelimit)
			rlConf := &rls_config.RateLimitDescriptor{
				Key:       DescriptorKeyForMethod,
				Value:     method,
				RateLimit: rlPolicyConfig,
			}

			rlsConfig := &rls_config.RateLimitDescriptor{
				Key:   DescriptorKeyForPath,
				Value: path,
				Descriptors: []*rls_config.RateLimitDescriptor{
					rlConf,
				},
			}
			rlsConfigs = append(rlsConfigs, rlsConfig)
			var org = resolveRatelimit.Organization

			r.apiLevelMu.Lock()
			defer r.apiLevelMu.Unlock()
			if _, ok := r.apiLevelRateLimitPolicies[org]; !ok {
				r.apiLevelRateLimitPolicies[org] = make(map[string]map[string][]*rls_config.RateLimitDescriptor)
				logger.Info("org", org)
			}
			for _, vHost := range vHosts {
				logger.Info("vHost", vHost)
				if _, ok := r.apiLevelRateLimitPolicies[org][vHost]; !ok {
					r.apiLevelRateLimitPolicies[org][vHost] = make(map[string][]*rls_config.RateLimitDescriptor)
				}
				apiIdentifier := GenerateIdentifierForAPIWithUUID(vHost, resolveRatelimit.UUID) + resource.Path + resource.Method
				r.apiLevelRateLimitPolicies[org][vHost][apiIdentifier] = rlsConfigs
			}
		}
	} else {
		logger.Info("Going to APILevel")
		apiLevelRLPolicyConfig := parseRateLimitPolicyToXDS(resolveRatelimit.API)
		rlsConfigs = append(rlsConfigs, &rls_config.RateLimitDescriptor{
			Key:   DescriptorKeyForPath,
			Value: resolveRatelimit.Context,
			Descriptors: []*rls_config.RateLimitDescriptor{
				{
					Key:       DescriptorKeyForMethod,
					Value:     DescriptorValueForAPIMethod,
					RateLimit: apiLevelRLPolicyConfig,
				},
			},
		},
		)
		if len(rlsConfigs) == 0 {
			return
		}

		var org = resolveRatelimit.Organization

		r.apiLevelMu.Lock()
		defer r.apiLevelMu.Unlock()
		if _, ok := r.apiLevelRateLimitPolicies[org]; !ok {
			r.apiLevelRateLimitPolicies[org] = make(map[string]map[string][]*rls_config.RateLimitDescriptor)
			logger.Info("org", org)
		}
		for _, vHost := range vHosts {
			logger.Info("vHost", vHost)
			if _, ok := r.apiLevelRateLimitPolicies[org][vHost]; !ok {
				r.apiLevelRateLimitPolicies[org][vHost] = make(map[string][]*rls_config.RateLimitDescriptor)
			}
			apiIdentifier := GenerateIdentifierForAPIWithUUID(vHost, resolveRatelimit.UUID)
			r.apiLevelRateLimitPolicies[org][vHost][apiIdentifier] = rlsConfigs
			logger.Info("apiIdentifier", apiIdentifier)
		}
	}
}

// DeleteAPILevelRateLimitPolicies deletes inline Rate Limit policies added with the API.
func (r *rateLimitPolicyCache) DeleteAPILevelRateLimitPolicies(org string, vHosts []string, apiID string) {
	r.apiLevelMu.Lock()
	defer r.apiLevelMu.Unlock()
	for _, vHost := range vHosts {
		apiIdentifier := GenerateIdentifierForAPIWithUUID(vHost, apiID)
		logger.Info("apiIdentifier", apiIdentifier)
		logger.Info("r.apiLevelRateLimitPolicies[org][vHost]", r.apiLevelRateLimitPolicies[org][vHost])
		delete(r.apiLevelRateLimitPolicies[org][vHost], apiIdentifier)
	}
}

// DeleteAPILevelRateLimitPolicies deletes inline Rate Limit policies added with the API.
func (r *rateLimitPolicyCache) DeleteResourceLevelRateLimitPolicies(org string, vHosts []string, apiID string, path string, method string) {
	r.apiLevelMu.Lock()
	defer r.apiLevelMu.Unlock()
	for _, vHost := range vHosts {
		apiIdentifier := GenerateIdentifierForAPIWithUUID(vHost, apiID) + path + method
		logger.Info("apiIdentifier", apiIdentifier)
		logger.Info("r.apiLevelRateLimitPolicies[org][vHost]", r.apiLevelRateLimitPolicies[org][vHost])
		delete(r.apiLevelRateLimitPolicies[org][vHost], apiIdentifier)
	}
}

// DeleteCustomRateLimitPolicies deletes Custom Rate Limit policies added.
func (r *rateLimitPolicyCache) DeleteCustomRateLimitPolicies(customRateLimitPolicy dpv1alpha1.CustomRateLimitPolicyDef) {
	r.apiLevelMu.Lock()
	defer r.apiLevelMu.Unlock()
	delete(r.customRateLimitPolicies[customRateLimitPolicy.Organization], customRateLimitPolicy.Key+"_"+customRateLimitPolicy.Value)
}

func (r *rateLimitPolicyCache) generateRateLimitConfig() *rls_config.RateLimitConfig {
	var orgDescriptors []*rls_config.RateLimitDescriptor

	r.apiLevelMu.RLock()
	defer r.apiLevelMu.RUnlock()
	logger.Info("r.apiLevelRateLimitPolicies", r.apiLevelRateLimitPolicies)
	// Generate API level rate limit configurations
	for org, orgPolicies := range r.apiLevelRateLimitPolicies {
		logger.Info("org", org)
		var vHostDescriptors []*rls_config.RateLimitDescriptor
		for vHost, vHostPolicies := range orgPolicies {
			logger.Info("vHost", vHost)
			var apiDescriptors []*rls_config.RateLimitDescriptor
			for _, apiPolicies := range vHostPolicies {
				logger.Info("apiPolicies", apiPolicies)
				// Configure API Level rate limit policies only if, the API is deployed to the gateway label
				// Check API deployed to the gateway label
				apiDescriptors = append(apiDescriptors, apiPolicies...)
			}
			vHostDescriptor := &rls_config.RateLimitDescriptor{
				Key:         DescriptorKeyForVhost,
				Value:       vHost,
				Descriptors: apiDescriptors,
			}
			vHostDescriptors = append(vHostDescriptors, vHostDescriptor)
		}
		orgDescriptor := &rls_config.RateLimitDescriptor{
			Key:         DescriptorKeyForOrg,
			Value:       org,
			Descriptors: vHostDescriptors,
		}
		orgDescriptors = append(orgDescriptors, orgDescriptor)
	}

	// Add custom rate limit policies as organization level rate limit policies
	customRateLimitDescriptors := r.generateCustomPolicyRateLimitConfig()
	orgDescriptors = append(orgDescriptors, customRateLimitDescriptors...)

	return &rls_config.RateLimitConfig{
		Name:        RateLimiterDomain,
		Domain:      RateLimiterDomain,
		Descriptors: orgDescriptors,
	}
}

// AddCustomRateLimitPolicies adds custom rate limit policies to the rateLimitPolicyCache.
func (r *rateLimitPolicyCache) AddCustomRateLimitPolicies(customRateLimitPolicy dpv1alpha1.CustomRateLimitPolicyDef) {
	if r.customRateLimitPolicies[customRateLimitPolicy.Organization] == nil {
		r.customRateLimitPolicies[customRateLimitPolicy.Organization] = make(map[string]*rls_config.RateLimitDescriptor)
		r.customRateLimitPolicies[customRateLimitPolicy.Organization][customRateLimitPolicy.Key+"_"+customRateLimitPolicy.Value] = &rls_config.RateLimitDescriptor{
			Key:   customRateLimitPolicy.Key,
			Value: customRateLimitPolicy.Value,
			RateLimit: &rls_config.RateLimitPolicy{
				Unit:            getRateLimitUnit(customRateLimitPolicy.RateLimit.Unit),
				RequestsPerUnit: uint32(customRateLimitPolicy.RateLimit.RequestsPerUnit),
			},
		}
	} else {
		r.customRateLimitPolicies[customRateLimitPolicy.Organization][customRateLimitPolicy.Key+"_"+customRateLimitPolicy.Value] = &rls_config.RateLimitDescriptor{
			Key:   customRateLimitPolicy.Key,
			Value: customRateLimitPolicy.Value,
			RateLimit: &rls_config.RateLimitPolicy{
				Unit:            getRateLimitUnit(customRateLimitPolicy.RateLimit.Unit),
				RequestsPerUnit: uint32(customRateLimitPolicy.RateLimit.RequestsPerUnit),
			},
		}
	}
}

func (r *rateLimitPolicyCache) updateXdsCache(label string) bool {
	rlsConf := r.generateRateLimitConfig()
	logger.Info("label", rlsConf)
	version := fmt.Sprint(rand.Intn(maxRandomInt))
	logger.Info("rlsConf", rlsConf)
	snap, err := gcp_cache.NewSnapshot(version, map[gcp_resource.Type][]gcp_types.Resource{
		gcp_resource.RateLimitConfigType: {
			rlsConf,
		},
	})
	logger.Info("snap", snap)
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(1714, err.Error()))
		return false
	}
	if err := snap.Consistent(); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(1715, err.Error()))
		return false
	}

	if err := r.xdsCache.SetSnapshot(context.Background(), label, snap); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(1716, err.Error()))
		return false
	}
	loggers.LoggerAPKOperator.Infof("New rate limit cache updated for the label: %q version: %q, API_UUID: %v", label, version, logging.GetValueFromLogContext("API_UUID"))
	loggers.LoggerAPKOperator.Debug("Updated rate limit config: ", rlsConf)
	return true
}

func parseRateLimitPolicyToXDS(policy dpv1alpha1.ResolveRateLimit) *rls_config.RateLimitPolicy {
	loggers.LoggerAPKOperator.Info("Rate count unit: ", policy.RequestsPerUnit)
	unit := getRateLimitUnit(policy.Unit)
	return &rls_config.RateLimitPolicy{
		Unit:            unit,
		RequestsPerUnit: uint32(policy.RequestsPerUnit),
	}

}

func getRateLimitUnit(name string) rls_config.RateLimitUnit {
	loggers.LoggerAPKOperator.Info("Rate limit unit: ", name)
	switch strings.ToUpper(name) {
	case "SECOND":
		return rls_config.RateLimitUnit_SECOND
	case "MINUTE":
		return rls_config.RateLimitUnit_MINUTE
	case "HOUR":
		return rls_config.RateLimitUnit_HOUR
	case "DAY":
		return rls_config.RateLimitUnit_DAY
	default:
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(1712, name))
		return rls_config.RateLimitUnit_UNKNOWN
	}
}

func init() {
	rlsPolicyCache = &rateLimitPolicyCache{
		xdsCache:                  gcp_cache.NewSnapshotCache(false, IDHash{}, nil),
		apiLevelRateLimitPolicies: make(map[string]map[string]map[string][]*rls_config.RateLimitDescriptor),
		customRateLimitPolicies:   make(map[string]map[string]*rls_config.RateLimitDescriptor),
	}
}

// generateCustomPolicyRateLimitConfig generates rate limit configurations for custom rate limit policies
// based on the policies stored in the rateLimitPolicyCache.
func (r *rateLimitPolicyCache) generateCustomPolicyRateLimitConfig() []*rls_config.RateLimitDescriptor {
	var orgDescriptors []*rls_config.RateLimitDescriptor
	for org, customRateLimitPolicies := range r.customRateLimitPolicies {
		descriptors := []*rls_config.RateLimitDescriptor{}
		for _, customRateLimitPolicy := range customRateLimitPolicies {
			descriptors = append(descriptors, customRateLimitPolicy)
		}
		orgDescriptors = append(orgDescriptors, &rls_config.RateLimitDescriptor{
			Key:         OrgMetadataKey,
			Value:       org,
			Descriptors: descriptors,
		})
	}
	return orgDescriptors
}

// SetEmptySnapshot sets an empty snapshot into the apiCache for the given label
// this is used to set empty snapshot when there are no APIs available for a label
func (r *rateLimitPolicyCache) SetEmptySnapshot(label string) bool {
	var rls = &rls_config.RateLimitConfig{
		Name:        RateLimiterDomain,
		Domain:      RateLimiterDomain,
		Descriptors: []*rls_config.RateLimitDescriptor{},
	}
	version := fmt.Sprint(rand.Intn(maxRandomInt))
	snap, err := gcp_cache.NewSnapshot(version, map[gcp_resource.Type][]gcp_types.Resource{
		gcp_resource.RateLimitConfigType: {
			rls,
		},
	})
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(1714, err.Error()))
		return false
	}
	if err := snap.Consistent(); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(1715, err.Error()))
		return false
	}

	if err := r.xdsCache.SetSnapshot(context.Background(), label, snap); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(1716, err.Error()))
		return false
	}
	loggers.LoggerAPKOperator.Infof("New rate limit cache updated for the label: %q version: %q, API_UUID: %v", label, version, logging.GetValueFromLogContext("API_UUID"))
	loggers.LoggerAPKOperator.Debug("Updated rate limit config: ", rls)
	return true
}
