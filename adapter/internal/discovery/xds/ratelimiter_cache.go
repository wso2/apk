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
	"github.com/wso2/apk/adapter/internal/loggers"
	logging "github.com/wso2/apk/adapter/internal/logging"
	"github.com/wso2/apk/adapter/internal/oasparser/envoyconf"
	"github.com/wso2/apk/adapter/internal/oasparser/model"
	mgw "github.com/wso2/apk/adapter/internal/oasparser/model"
	"github.com/wso2/apk/adapter/pkg/utils/stringutils"
)

var rlsPolicyCache *rateLimitPolicyCache

func getRateLimitUnit(name string) rls_config.RateLimitUnit {
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
		loggers.LoggerXds.ErrorC(logging.GetErrorByCode(1712, name))
		return rls_config.RateLimitUnit_UNKNOWN
	}
}

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
func (r *rateLimitPolicyCache) AddAPILevelRateLimitPolicies(vHosts []string, mgwSwagger *mgw.MgwSwagger) {

	rlsConfigs := []*rls_config.RateLimitDescriptor{}

	// The map apiOperations is used to keep `Pat:HTTPmethod` unique to make sure the Rate Limiter Config to be consistent (not to have duplicate rate limit policies)
	// path -> HTTP method
	apiOperations := make(map[string]map[string]struct{})
	for _, resource := range mgwSwagger.GetResources() {
		path := mgwSwagger.GetXWso2Basepath() + resource.GetPath()
		if _, ok := apiOperations[path]; !ok {
			apiOperations[path] = make(map[string]struct{})
		}
		operationRlsConfigs := []*rls_config.RateLimitDescriptor{}
		for _, operation := range resource.GetMethod() {
			method := operation.GetMethod()
			if _, ok := apiOperations[path][method]; ok {
				// Unreachable if the swagger definition is valid
				loggers.LoggerXds.Warnf("Duplicate API resource HTTP method %q %q in the swagger definition, skipping rate limit policy for the duplicate resource.", path, method)
				continue
			}

			if operation.RateLimitPolicy != nil {
				rlPolicyConfig := parseRateLimitPolicyToXDS(operation.RateLimitPolicy)
				rlConf := &rls_config.RateLimitDescriptor{
					Key:       envoyconf.DescriptorKeyForMethod,
					Value:     method,
					RateLimit: rlPolicyConfig,
				}
				operationRlsConfigs = append(operationRlsConfigs, rlConf)
				apiOperations[path][method] = void
			}
		}
		if len(operationRlsConfigs) > 0 {
			rlsConfig := &rls_config.RateLimitDescriptor{
				Key:         envoyconf.DescriptorKeyForPath,
				Value:       path,
				Descriptors: operationRlsConfigs,
			}
			rlsConfigs = append(rlsConfigs, rlsConfig)
		}
	}

	if mgwSwagger.RateLimitPolicy != nil {
		apiLevelRLPolicyConfig := parseRateLimitPolicyToXDS(mgwSwagger.RateLimitPolicy)
		rlsConfigs = append(rlsConfigs, &rls_config.RateLimitDescriptor{
			Key:   envoyconf.DescriptorKeyForPath,
			Value: mgwSwagger.GetXWso2Basepath(),
			Descriptors: []*rls_config.RateLimitDescriptor{
				{
					Key:       envoyconf.DescriptorKeyForMethod,
					Value:     envoyconf.DescriptorValueForAPIMethod,
					RateLimit: apiLevelRLPolicyConfig,
				},
			},
		},
		)
	}

	if len(rlsConfigs) == 0 {
		return
	}

	org := mgwSwagger.OrganizationID

	r.apiLevelMu.Lock()
	defer r.apiLevelMu.Unlock()
	if _, ok := r.apiLevelRateLimitPolicies[org]; !ok {
		r.apiLevelRateLimitPolicies[org] = make(map[string]map[string][]*rls_config.RateLimitDescriptor)
	}
	for _, vHost := range vHosts {
		if _, ok := r.apiLevelRateLimitPolicies[org][vHost]; !ok {
			r.apiLevelRateLimitPolicies[org][vHost] = make(map[string][]*rls_config.RateLimitDescriptor)
		}
		apiIdentifier := GenerateIdentifierForAPIWithUUID(vHost, mgwSwagger.UUID)
		r.apiLevelRateLimitPolicies[org][vHost][apiIdentifier] = rlsConfigs
	}
}

// DeleteAPILevelRateLimitPolicies deletes inline Rate Limit policies added with the API.
func (r *rateLimitPolicyCache) DeleteAPILevelRateLimitPolicies(org, vHost, apiID string) {
	r.apiLevelMu.Lock()
	defer r.apiLevelMu.Unlock()
	delete(r.apiLevelRateLimitPolicies[org][vHost], apiID)
}

func parseRateLimitPolicyToXDS(policy *mgw.RateLimitPolicy) *rls_config.RateLimitPolicy {
	if policy != nil {
		unit := getRateLimitUnit(policy.SpanUnit)
		return &rls_config.RateLimitPolicy{
			Unit:            unit,
			RequestsPerUnit: uint32(policy.Count),
		}

	}
	return nil
}

func (r *rateLimitPolicyCache) generateRateLimitConfig(label string) *rls_config.RateLimitConfig {
	var orgDescriptors []*rls_config.RateLimitDescriptor

	r.apiLevelMu.RLock()
	defer r.apiLevelMu.RUnlock()

	// Generate API level rate limit configurations
	for org, orgPolicies := range r.apiLevelRateLimitPolicies {
		var vHostDescriptors []*rls_config.RateLimitDescriptor
		for vHost, vHostPolicies := range orgPolicies {
			var apiDescriptors []*rls_config.RateLimitDescriptor
			for apiID, apiPolicies := range vHostPolicies {
				// Configure API Level rate limit policies only if, the API is deployed to the gateway label
				// Check API deployed to the gateway label
				if stringutils.StringInSlice(label, orgIDOpenAPIEnvoyMap[org][apiID]) {
					apiDescriptors = append(apiDescriptors, apiPolicies...)
				}
			}
			vHostDescriptor := &rls_config.RateLimitDescriptor{
				Key:         envoyconf.DescriptorKeyForVhost,
				Value:       vHost,
				Descriptors: apiDescriptors,
			}
			vHostDescriptors = append(vHostDescriptors, vHostDescriptor)
		}
		orgDescriptor := &rls_config.RateLimitDescriptor{
			Key:         envoyconf.DescriptorKeyForOrg,
			Value:       org,
			Descriptors: vHostDescriptors,
		}
		orgDescriptors = append(orgDescriptors, orgDescriptor)
	}
	customRateLimitDescriptors := r.generateCustomPolicyRateLimitConfig()
	orgDescriptors = append(orgDescriptors, customRateLimitDescriptors...)
	return &rls_config.RateLimitConfig{
		Name:        envoyconf.RateLimiterDomain,
		Domain:      envoyconf.RateLimiterDomain,
		Descriptors: orgDescriptors,
	}
}

func (r *rateLimitPolicyCache) updateXdsCache(label string) bool {
	rlsConf := r.generateRateLimitConfig(label)

	version := fmt.Sprint(rand.Intn(maxRandomInt))
	snap, err := gcp_cache.NewSnapshot(version, map[gcp_resource.Type][]gcp_types.Resource{
		gcp_resource.RateLimitConfigType: {
			rlsConf,
		},
	})
	if err != nil {
		loggers.LoggerXds.ErrorC(logging.GetErrorByCode(1714, err.Error()))
		return false
	}
	if err := snap.Consistent(); err != nil {
		loggers.LoggerXds.ErrorC(logging.GetErrorByCode(1715, err.Error()))
		return false
	}

	if err := r.xdsCache.SetSnapshot(context.Background(), label, snap); err != nil {
		loggers.LoggerXds.ErrorC(logging.GetErrorByCode(1716, err.Error()))
		return false
	}
	loggers.LoggerXds.Infof("New rate limit cache updated for the label: %q version: %q", label, version)
	loggers.LoggerXds.Debug("Updated rate limit config: ", rlsConf)
	return true
}

func (r *rateLimitPolicyCache) AddCustomRateLimitPolicies(customRateLimitPolicies []*model.CustomRateLimitPolicy) {
	r.customRateLimitPolicies = make(map[string]map[string]*rls_config.RateLimitDescriptor) 
	for _, customRateLimitPolicy := range customRateLimitPolicies {
		if r.customRateLimitPolicies[customRateLimitPolicy.Organization] == nil {
			r.customRateLimitPolicies[customRateLimitPolicy.Organization] = make(map[string]*rls_config.RateLimitDescriptor)
			r.customRateLimitPolicies[customRateLimitPolicy.Organization][customRateLimitPolicy.Key+"_"+customRateLimitPolicy.Value] = &rls_config.RateLimitDescriptor{
				Key:  customRateLimitPolicy.Key,
				Value: customRateLimitPolicy.Value,
				RateLimit: &rls_config.RateLimitPolicy{
					Unit:            getRateLimitUnit(customRateLimitPolicy.RateLimit.Unit),
					RequestsPerUnit: uint32(customRateLimitPolicy.RateLimit.RequestsPerUnit),
				},
			}
		} else {
			r.customRateLimitPolicies[customRateLimitPolicy.Organization][customRateLimitPolicy.Key+"_"+customRateLimitPolicy.Value] = &rls_config.RateLimitDescriptor{
				Key:  customRateLimitPolicy.Key,
				Value: customRateLimitPolicy.Value,
				RateLimit: &rls_config.RateLimitPolicy{
					Unit:            getRateLimitUnit(customRateLimitPolicy.RateLimit.Unit),
					RequestsPerUnit: uint32(customRateLimitPolicy.RateLimit.RequestsPerUnit),
				},
			}
		}
	}
}

func (r *rateLimitPolicyCache) generateCustomPolicyRateLimitConfig() []*rls_config.RateLimitDescriptor{
	var orgDescriptors []*rls_config.RateLimitDescriptor
	for org, customRateLimitPolicies := range r.customRateLimitPolicies {
		descriptors := []*rls_config.RateLimitDescriptor{}
		for _, customRateLimitPolicy := range customRateLimitPolicies {
			descriptors = append(descriptors, customRateLimitPolicy)
		}
		orgDescriptors = append(orgDescriptors, &rls_config.RateLimitDescriptor{
			Key:         envoyconf.DescriptorKeyForOrg,
			Value:       org,
			Descriptors: descriptors,
		})
	}
	return orgDescriptors
}

func init() {
	rlsPolicyCache = &rateLimitPolicyCache{
		xdsCache:                  gcp_cache.NewSnapshotCache(false, IDHash{}, nil),
		apiLevelRateLimitPolicies: make(map[string]map[string]map[string][]*rls_config.RateLimitDescriptor),
		customRateLimitPolicies: make(map[string]map[string]*rls_config.RateLimitDescriptor),
	}
}
