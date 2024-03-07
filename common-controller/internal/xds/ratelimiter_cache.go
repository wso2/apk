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
	"crypto/rand"
	"fmt"
	"strings"
	"sync"

	gcp_types "github.com/envoyproxy/go-control-plane/pkg/cache/types"
	gcp_cache "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	gcp_resource "github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	rls_config "github.com/envoyproxy/go-control-plane/ratelimit/config/ratelimit/v3"
	logger "github.com/sirupsen/logrus"
	"github.com/wso2/apk/adapter/pkg/logging"
	"github.com/wso2/apk/common-controller/internal/loggers"
	dpv1alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha1"
	"github.com/wso2/apk/common-go-libs/constants"
)

// Constants relevant to the route related ratelimit configurations
const (
	DescriptorKeyForOrg                = "org"
	OrgMetadataKey                     = "customorg"
	DescriptorKeyForEnvironment        = "environment"
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
	// org -> environment -> API-Identifier (i.e. Environment:API-UUID) -> Rate Limit Configs
	apiLevelRateLimitPolicies map[string]map[string]map[string]map[string]*rls_config.RateLimitDescriptor

	// org -> Custom Rate Limit Configs
	customRateLimitPolicies map[string]map[string]*rls_config.RateLimitDescriptor

	// mutex for API level
	apiLevelMu sync.RWMutex
}

// AddAPILevelRateLimitPolicies adds inline Rate Limit policies in APIs to be updated in the Rate Limiter service.
func (r *rateLimitPolicyCache) AddAPILevelRateLimitPolicies(resolveRatelimit dpv1alpha1.ResolveRateLimitAPIPolicy) {

	rlsConfigs := rls_config.RateLimitDescriptor{}
	httpMethods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}
	// The map apiOperations is used to keep `Pat:HTTPmethod` unique to make sure the Rate Limiter Config to be consistent (not to have duplicate rate limit policies)
	// path -> HTTP method

	if len(resolveRatelimit.Resources) != 0 {
		for _, resource := range resolveRatelimit.Resources {
			var org = resolveRatelimit.Organization

			path := resolveRatelimit.BasePath + resolveRatelimit.BasePath + resource.Path
			logger.Debug("path", path)

			method := resource.Method

			rlPolicyConfig := parseRateLimitPolicyToXDS(resource.ResourceRatelimit)
			if method == constants.All {
				for _, httpMethod := range httpMethods {
					rlConf := &rls_config.RateLimitDescriptor{
						Key:       DescriptorKeyForMethod,
						Value:     httpMethod,
						RateLimit: rlPolicyConfig,
					}

					if _, ok := r.apiLevelRateLimitPolicies[org]; !ok {
						r.apiLevelRateLimitPolicies[org] = make(map[string]map[string]map[string]*rls_config.RateLimitDescriptor)
					}

					environment := resolveRatelimit.Environment
					if _, ok := r.apiLevelRateLimitPolicies[org][environment]; !ok {
						r.apiLevelRateLimitPolicies[org][environment] = make(map[string]map[string]*rls_config.RateLimitDescriptor)
					}

					if _, ok := r.apiLevelRateLimitPolicies[org][environment][resolveRatelimit.BasePath+resolveRatelimit.BasePath+resource.Path]; !ok {
						r.apiLevelRateLimitPolicies[org][environment][resolveRatelimit.BasePath+resolveRatelimit.BasePath+resource.Path] = make(map[string]*rls_config.RateLimitDescriptor)
						r.apiLevelRateLimitPolicies[org][environment][resolveRatelimit.BasePath+resolveRatelimit.BasePath+resource.Path][httpMethod] = rlConf
					} else {
						r.apiLevelRateLimitPolicies[org][environment][resolveRatelimit.BasePath+resolveRatelimit.BasePath+resource.Path][httpMethod] = rlConf
					}
				}
			} else {
				rlConf := &rls_config.RateLimitDescriptor{
					Key:       DescriptorKeyForMethod,
					Value:     method,
					RateLimit: rlPolicyConfig,
				}

				r.apiLevelMu.Lock()
				defer r.apiLevelMu.Unlock()
				if _, ok := r.apiLevelRateLimitPolicies[org]; !ok {
					r.apiLevelRateLimitPolicies[org] = make(map[string]map[string]map[string]*rls_config.RateLimitDescriptor)
				}

				environment := resolveRatelimit.Environment
				if _, ok := r.apiLevelRateLimitPolicies[org][environment]; !ok {
					r.apiLevelRateLimitPolicies[org][environment] = make(map[string]map[string]*rls_config.RateLimitDescriptor)
				}

				if _, ok := r.apiLevelRateLimitPolicies[org][environment][resolveRatelimit.BasePath+resolveRatelimit.BasePath+resource.Path]; !ok {
					r.apiLevelRateLimitPolicies[org][environment][resolveRatelimit.BasePath+resolveRatelimit.BasePath+resource.Path] = make(map[string]*rls_config.RateLimitDescriptor)
					r.apiLevelRateLimitPolicies[org][environment][resolveRatelimit.BasePath+resolveRatelimit.BasePath+resource.Path][method] = rlConf
				} else {
					r.apiLevelRateLimitPolicies[org][environment][resolveRatelimit.BasePath+resolveRatelimit.BasePath+resource.Path][method] = rlConf
				}
			}
		}
	} else {
		logger.Debug("Going to APILevel")
		apiLevelRLPolicyConfig := parseRateLimitPolicyToXDS(resolveRatelimit.API)
		rlsConfigs = rls_config.RateLimitDescriptor{

			Key:       DescriptorKeyForMethod,
			Value:     DescriptorValueForAPIMethod,
			RateLimit: apiLevelRLPolicyConfig,
		}

		var org = resolveRatelimit.Organization

		r.apiLevelMu.Lock()
		defer r.apiLevelMu.Unlock()
		if _, ok := r.apiLevelRateLimitPolicies[org]; !ok {
			r.apiLevelRateLimitPolicies[org] = make(map[string]map[string]map[string]*rls_config.RateLimitDescriptor)
		}

		environment := resolveRatelimit.Environment
		if _, ok := r.apiLevelRateLimitPolicies[org][environment]; !ok {
			r.apiLevelRateLimitPolicies[org][environment] = make(map[string]map[string]*rls_config.RateLimitDescriptor)
		}
		if _, ok := r.apiLevelRateLimitPolicies[org][environment][resolveRatelimit.BasePath]; !ok {
			r.apiLevelRateLimitPolicies[org][environment][resolveRatelimit.BasePath] = make(map[string]*rls_config.RateLimitDescriptor)
		}
		r.apiLevelRateLimitPolicies[org][environment][resolveRatelimit.BasePath][DescriptorValueForAPIMethod] = &rlsConfigs
	}
}

// DeleteAPILevelRateLimitPolicies deletes inline Rate Limit policies added with the API.
func (r *rateLimitPolicyCache) DeleteAPILevelRateLimitPolicies(org string, environment string, basePath string) {
	r.apiLevelMu.Lock()
	defer r.apiLevelMu.Unlock()
	delete(r.apiLevelRateLimitPolicies[org][environment][basePath], DescriptorValueForAPIMethod)
}

// DeleteAPILevelRateLimitPolicies deletes inline Rate Limit policies added with the API.
func (r *rateLimitPolicyCache) DeleteResourceLevelRateLimitPolicies(org string, environment string, basePath string, path string, method string) {
	httpMethods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}
	r.apiLevelMu.Lock()
	defer r.apiLevelMu.Unlock()
	if method == constants.All {
		for _, httpMethod := range httpMethods {
			delete(r.apiLevelRateLimitPolicies[org][environment][basePath+basePath+path], httpMethod)
		}
	} else {
		delete(r.apiLevelRateLimitPolicies[org][environment][basePath+basePath+path], method)
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
	// Generate API level rate limit configurations
	for org, orgPolicies := range r.apiLevelRateLimitPolicies {
		var envDescriptors []*rls_config.RateLimitDescriptor
		for env, envPolicies := range orgPolicies {
			var apiPathDiscriptors []*rls_config.RateLimitDescriptor
			for path, apiPathPolicies := range envPolicies {
				// Configure API Level rate limit policies only if, the API is deployed to the gateway label
				// Check API deployed to the gateway label
				var methodDescriptors []*rls_config.RateLimitDescriptor
				for _, methodPolicies := range apiPathPolicies {
					methodDescriptors = append(methodDescriptors, methodPolicies)

				}
				apiPathDiscriptor := &rls_config.RateLimitDescriptor{
					Key:         DescriptorKeyForPath,
					Value:       path,
					Descriptors: methodDescriptors,
				}
				apiPathDiscriptors = append(apiPathDiscriptors, apiPathDiscriptor)

			}
			envDescriptor := &rls_config.RateLimitDescriptor{
				Key:         DescriptorKeyForEnvironment,
				Value:       env,
				Descriptors: apiPathDiscriptors,
			}
			envDescriptors = append(envDescriptors, envDescriptor)
		}

		orgDescriptor := &rls_config.RateLimitDescriptor{
			Key:         DescriptorKeyForOrg,
			Value:       org,
			Descriptors: envDescriptors,
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
				Unit:            getRateLimitUnit(customRateLimitPolicy.Unit),
				RequestsPerUnit: uint32(customRateLimitPolicy.RequestsPerUnit),
			},
		}
	} else {
		r.customRateLimitPolicies[customRateLimitPolicy.Organization][customRateLimitPolicy.Key+"_"+customRateLimitPolicy.Value] = &rls_config.RateLimitDescriptor{
			Key:   customRateLimitPolicy.Key,
			Value: customRateLimitPolicy.Value,
			RateLimit: &rls_config.RateLimitPolicy{
				Unit:            getRateLimitUnit(customRateLimitPolicy.Unit),
				RequestsPerUnit: uint32(customRateLimitPolicy.RequestsPerUnit),
			},
		}
	}
}

func (r *rateLimitPolicyCache) updateXdsCache(label string) bool {
	rlsConf := r.generateRateLimitConfig()
	version := fmt.Sprint(rand.Int(rand.Reader, maxRandomBigInt()))
	snap, err := gcp_cache.NewSnapshot(version, map[gcp_resource.Type][]gcp_types.Resource{
		gcp_resource.RateLimitConfigType: {
			rlsConf,
		},
	})
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error1714, logging.MAJOR,
			"Error while creating the rate limit snapshot: %v", err.Error()))
		return false
	}
	if err := snap.Consistent(); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error1715, logging.MAJOR,
			"Inconsistent rate limiter snapshot: %v", err.Error()))
		return false
	}

	if err := r.xdsCache.SetSnapshot(context.Background(), label, snap); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error1716, logging.MAJOR,
			"Error while updating the rate limit snapshot: %v", err.Error()))
		return false
	}
	loggers.LoggerAPKOperator.Infof("New rate limit cache updated for the label: %q version: %q", label, version)
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
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error1712, logging.MAJOR,
			"Unknown rate limit unit %q, defaulting to UNKNOWN", name))
		return rls_config.RateLimitUnit_UNKNOWN
	}
}

func init() {
	rlsPolicyCache = &rateLimitPolicyCache{
		xdsCache:                  gcp_cache.NewSnapshotCache(false, IDHash{}, nil),
		apiLevelRateLimitPolicies: make(map[string]map[string]map[string]map[string]*rls_config.RateLimitDescriptor),
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
	version := fmt.Sprint(rand.Int(rand.Reader, maxRandomBigInt()))
	snap, err := gcp_cache.NewSnapshot(version, map[gcp_resource.Type][]gcp_types.Resource{
		gcp_resource.RateLimitConfigType: {
			rls,
		},
	})
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error1714, logging.MAJOR,
			"Error while creating the rate limit snapshot: %v", err.Error()))
		return false
	}
	if err := snap.Consistent(); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error1715, logging.MAJOR,
			"Inconsistent rate limiter snapshot: %v", err.Error()))
		return false
	}

	if err := r.xdsCache.SetSnapshot(context.Background(), label, snap); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error1716, logging.MAJOR,
			"Error while updating the rate limit snapshot: %v", err.Error()))
		return false
	}
	loggers.LoggerAPKOperator.Infof("New rate limit cache updated for the label: %q version: %q", label, version)
	loggers.LoggerAPKOperator.Debug("Updated rate limit config: ", rls)
	return true
}
