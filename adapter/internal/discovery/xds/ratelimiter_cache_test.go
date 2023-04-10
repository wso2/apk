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
	"fmt"
	"testing"

	rls_config "github.com/envoyproxy/go-control-plane/ratelimit/config/ratelimit/v3"
	"github.com/stretchr/testify/assert"
	"github.com/wso2/apk/adapter/internal/oasparser/model"
)

func TestAddDeleteAPILevelRateLimitPolicies(t *testing.T) {
	t.Run("Add API level rate limiting", testAddAPILevelRateLimitPolicies)
	t.Run("Delete API level rate limiting", testDeleteAPILevelRateLimitPolicies)
}

//todo(amali) add api with both api and operation level rate limiting

func testAddAPILevelRateLimitPolicies(t *testing.T) {
	p5000PerMin := &model.RateLimitPolicy{Count: 5000, SpanUnit: "MINUTE"}
	p2000PerMin := &model.RateLimitPolicy{Count: 2000, SpanUnit: "MINUTE"}
	p100000PerHOUR := &model.RateLimitPolicy{Count: 100000, SpanUnit: "HOUR"}

	tests := []struct {
		desc                      string
		adapterInternalAPI        *model.AdapterInternalAPI
		apiLevelRateLimitPolicies map[string]map[string]map[string][]*rls_config.RateLimitDescriptor
	}{
		{
			desc:                      "Add an API with no Rate Limit policies",
			adapterInternalAPI:        getDummyAPISwagger("1", nil, nil, nil, nil, nil),
			apiLevelRateLimitPolicies: map[string]map[string]map[string][]*rls_config.RateLimitDescriptor{},
		},
		{
			// Note: Each test case is depend on the earlier test cases
			desc:               "Add an API with API Level Rate Limit Policy",
			adapterInternalAPI: getDummyAPISwagger("2", p5000PerMin, nil, nil, nil, nil),
			apiLevelRateLimitPolicies: map[string]map[string]map[string][]*rls_config.RateLimitDescriptor{
				"org1": {"vhost1": {
					"vhost1:2": {&rls_config.RateLimitDescriptor{
						Key:   "path",
						Value: "/base-path-2",
						Descriptors: []*rls_config.RateLimitDescriptor{
							{
								Key:   "method",
								Value: "ALL",
								RateLimit: &rls_config.RateLimitPolicy{
									Unit:            rls_config.RateLimitUnit_MINUTE,
									RequestsPerUnit: 5000,
								},
							},
						},
					}},
				}},
			},
		},
		{
			// Note: Each test case is depend on the earlier test cases
			desc:               "Add an API with no Rate Limit policies",
			adapterInternalAPI: getDummyAPISwagger("4", nil, nil, nil, nil, nil),
			apiLevelRateLimitPolicies: map[string]map[string]map[string][]*rls_config.RateLimitDescriptor{
				"org1": {"vhost1": {
					"vhost1:2": {&rls_config.RateLimitDescriptor{
						Key:   "path",
						Value: "/base-path-2",
						Descriptors: []*rls_config.RateLimitDescriptor{
							{
								Key:   "method",
								Value: "ALL",
								RateLimit: &rls_config.RateLimitPolicy{
									Unit:            rls_config.RateLimitUnit_MINUTE,
									RequestsPerUnit: 5000,
								},
							},
						},
					}},
				}},
			},
		},
		{
			// Note: Each test case is depend on the earlier test cases
			desc:               "Add an API with Operation Level Rate Limit policies",
			adapterInternalAPI: getDummyAPISwagger("5", nil, p100000PerHOUR, nil, nil, p2000PerMin),
			apiLevelRateLimitPolicies: map[string]map[string]map[string][]*rls_config.RateLimitDescriptor{
				"org1": {"vhost1": {
					"vhost1:2": {&rls_config.RateLimitDescriptor{
						Key:   "path",
						Value: "/base-path-2",
						Descriptors: []*rls_config.RateLimitDescriptor{
							{
								Key:   "method",
								Value: "ALL",
								RateLimit: &rls_config.RateLimitPolicy{
									Unit:            rls_config.RateLimitUnit_MINUTE,
									RequestsPerUnit: 5000,
								},
							},
						},
					}},
					"vhost1:5": {
						&rls_config.RateLimitDescriptor{
							Key:   "path",
							Value: "/base-path-5/res1",
							Descriptors: []*rls_config.RateLimitDescriptor{
								{
									Key:   "method",
									Value: "GET",
									RateLimit: &rls_config.RateLimitPolicy{
										Unit:            rls_config.RateLimitUnit_HOUR,
										RequestsPerUnit: 100000,
									},
								},
							},
						},
						&rls_config.RateLimitDescriptor{
							Key:   "path",
							Value: "/base-path-5/res2",
							Descriptors: []*rls_config.RateLimitDescriptor{
								{
									Key:   "method",
									Value: "POST",
									RateLimit: &rls_config.RateLimitPolicy{
										Unit:            rls_config.RateLimitUnit_MINUTE,
										RequestsPerUnit: 2000,
									},
								},
							},
						},
					},
				}},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			rlsPolicyCache.AddAPILevelRateLimitPolicies([]string{"vhost1"}, test.adapterInternalAPI)
			assert.Equal(t, test.apiLevelRateLimitPolicies, rlsPolicyCache.apiLevelRateLimitPolicies, test.desc)
		})
	}
}

func testDeleteAPILevelRateLimitPolicies(t *testing.T) {
	tests := []struct {
		desc                      string
		org                       string
		vHost                     string
		apiID                     string
		apiLevelRateLimitPolicies map[string]map[string]map[string][]*rls_config.RateLimitDescriptor
	}{
		{
			desc:  "Delete API with API level rate limits: vhost1:2",
			org:   "org1",
			vHost: "vhost1",
			apiID: "vhost1:2",
			apiLevelRateLimitPolicies: map[string]map[string]map[string][]*rls_config.RateLimitDescriptor{
				"org1": {"vhost1": {
					"vhost1:5": {
						&rls_config.RateLimitDescriptor{
							Key:   "path",
							Value: "/base-path-5/res1",
							Descriptors: []*rls_config.RateLimitDescriptor{
								{
									Key:   "method",
									Value: "GET",
									RateLimit: &rls_config.RateLimitPolicy{
										Unit:            rls_config.RateLimitUnit_HOUR,
										RequestsPerUnit: 100000,
									},
								},
							},
						},
						&rls_config.RateLimitDescriptor{
							Key:   "path",
							Value: "/base-path-5/res2",
							Descriptors: []*rls_config.RateLimitDescriptor{
								{
									Key:   "method",
									Value: "POST",
									RateLimit: &rls_config.RateLimitPolicy{
										Unit:            rls_config.RateLimitUnit_MINUTE,
										RequestsPerUnit: 2000,
									},
								},
							},
						},
					},
				}},
			},
		},
		{
			desc:  "Delete API with no API level rate limits: vhost1:1",
			org:   "org1",
			vHost: "vhost1",
			apiID: "vhost1:3",
			apiLevelRateLimitPolicies: map[string]map[string]map[string][]*rls_config.RateLimitDescriptor{
				"org1": {"vhost1": {
					"vhost1:5": {
						&rls_config.RateLimitDescriptor{
							Key:   "path",
							Value: "/base-path-5/res1",
							Descriptors: []*rls_config.RateLimitDescriptor{
								{
									Key:   "method",
									Value: "GET",
									RateLimit: &rls_config.RateLimitPolicy{
										Unit:            rls_config.RateLimitUnit_HOUR,
										RequestsPerUnit: 100000,
									},
								},
							},
						},
						&rls_config.RateLimitDescriptor{
							Key:   "path",
							Value: "/base-path-5/res2",
							Descriptors: []*rls_config.RateLimitDescriptor{
								{
									Key:   "method",
									Value: "POST",
									RateLimit: &rls_config.RateLimitPolicy{
										Unit:            rls_config.RateLimitUnit_MINUTE,
										RequestsPerUnit: 2000,
									},
								},
							},
						},
					},
				}},
			},
		},
		{
			desc:  "Delete API with operation level rate limits: vhost1:5",
			org:   "org1",
			vHost: "vhost1",
			apiID: "vhost1:5",
			apiLevelRateLimitPolicies: map[string]map[string]map[string][]*rls_config.RateLimitDescriptor{
				"org1": {"vhost1": {}},
			},
		},
		{
			desc:  "Delete API in an Org that has no APIs with rate limits",
			org:   "org4",
			vHost: "vhost1",
			apiID: "vhost1:5",
			apiLevelRateLimitPolicies: map[string]map[string]map[string][]*rls_config.RateLimitDescriptor{
				"org1": {"vhost1": {}},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			rlsPolicyCache.DeleteAPILevelRateLimitPolicies(test.org, test.vHost, test.apiID)
			assert.Equal(t, test.apiLevelRateLimitPolicies, rlsPolicyCache.apiLevelRateLimitPolicies)
		})
	}
}

func TestGenerateRateLimitConfig(t *testing.T) {
	tests := []struct {
		desc                      string
		orgIDOpenAPIEnvoyMap      map[string]map[string][]string
		apiLevelRateLimitPolicies map[string]map[string]map[string][]*rls_config.RateLimitDescriptor
		rlsConfig                 *rls_config.RateLimitConfig
	}{
		{
			desc: "Test config with multiple labels",
			orgIDOpenAPIEnvoyMap: map[string]map[string][]string{
				"org1": {
					"vhost1:2": []string{"Default"},
					"vhost1:3": []string{"Dev"},
					"vhost1:5": []string{"Default"},
				},
			},
			apiLevelRateLimitPolicies: map[string]map[string]map[string][]*rls_config.RateLimitDescriptor{
				"org1": {"vhost1": {
					"vhost1:2": {&rls_config.RateLimitDescriptor{
						Key:   "path",
						Value: "/base-path-2",
						Descriptors: []*rls_config.RateLimitDescriptor{
							{
								Key:   "method",
								Value: "ALL",
								RateLimit: &rls_config.RateLimitPolicy{
									Unit:            rls_config.RateLimitUnit_MINUTE,
									RequestsPerUnit: 5000,
								},
							},
						},
					}},
					"vhost1:3": {&rls_config.RateLimitDescriptor{
						Key:   "path",
						Value: "/base-path-2",
						Descriptors: []*rls_config.RateLimitDescriptor{
							{
								Key:   "method",
								Value: "ALL",
								RateLimit: &rls_config.RateLimitPolicy{
									Unit:            rls_config.RateLimitUnit_MINUTE,
									RequestsPerUnit: 5000,
								},
							},
						},
					}},
					"vhost1:5": {
						&rls_config.RateLimitDescriptor{
							Key:   "path",
							Value: "/base-path-5/res1",
							Descriptors: []*rls_config.RateLimitDescriptor{
								{
									Key:   "method",
									Value: "GET",
									RateLimit: &rls_config.RateLimitPolicy{
										Unit:            rls_config.RateLimitUnit_HOUR,
										RequestsPerUnit: 100000,
									},
								},
								{
									Key:   "method",
									Value: "POST",
									RateLimit: &rls_config.RateLimitPolicy{
										Unit:            rls_config.RateLimitUnit_MINUTE,
										RequestsPerUnit: 1000,
									},
								},
							},
						},
						&rls_config.RateLimitDescriptor{
							Key:   "path",
							Value: "/base-path-5/res2",
							Descriptors: []*rls_config.RateLimitDescriptor{
								{
									Key:   "method",
									Value: "POST",
									RateLimit: &rls_config.RateLimitPolicy{
										Unit:            rls_config.RateLimitUnit_MINUTE,
										RequestsPerUnit: 2000,
									},
								},
							},
						},
					},
				}},
			},
			rlsConfig: &rls_config.RateLimitConfig{
				Name:   "Default",
				Domain: "Default",
				Descriptors: []*rls_config.RateLimitDescriptor{
					{
						Key:   "org",
						Value: "org1",
						Descriptors: []*rls_config.RateLimitDescriptor{
							{
								Key:   "vhost",
								Value: "vhost1",
								Descriptors: []*rls_config.RateLimitDescriptor{
									{
										Key:   "path",
										Value: "/base-path-2",
										Descriptors: []*rls_config.RateLimitDescriptor{
											{
												Key:   "method",
												Value: "ALL",
												RateLimit: &rls_config.RateLimitPolicy{
													Unit:            rls_config.RateLimitUnit_MINUTE,
													RequestsPerUnit: 5000,
												},
											},
										},
									},
									{
										Key:   "path",
										Value: "/base-path-5/res1",
										Descriptors: []*rls_config.RateLimitDescriptor{
											{
												Key:   "method",
												Value: "GET",
												RateLimit: &rls_config.RateLimitPolicy{
													Unit:            rls_config.RateLimitUnit_HOUR,
													RequestsPerUnit: 100000,
												},
											},
											{
												Key:   "method",
												Value: "POST",
												RateLimit: &rls_config.RateLimitPolicy{
													Unit:            rls_config.RateLimitUnit_MINUTE,
													RequestsPerUnit: 1000,
												},
											},
										},
									},
									{
										Key:   "path",
										Value: "/base-path-5/res2",
										Descriptors: []*rls_config.RateLimitDescriptor{
											{
												Key:   "method",
												Value: "POST",
												RateLimit: &rls_config.RateLimitPolicy{
													Unit:            rls_config.RateLimitUnit_MINUTE,
													RequestsPerUnit: 2000,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			orgIDOpenAPIEnvoyMap = test.orgIDOpenAPIEnvoyMap
			c := &rateLimitPolicyCache{
				apiLevelRateLimitPolicies: test.apiLevelRateLimitPolicies,
			}
			actualConf := c.generateRateLimitConfig("Default")
			// Construct "expected" and "actual" here, since the diff gen by assert is bit difficult to read.
			valuesAsStr := fmt.Sprintf("expected: %v\nactual: %v", test.rlsConfig, actualConf)

			// Test descriptors inside Org1, vHost1 (because the order of the elements can not be guaranteed)
			assert.ElementsMatch(t, test.rlsConfig.Descriptors[0].Descriptors[0].Descriptors,
				actualConf.Descriptors[0].Descriptors[0].Descriptors, valuesAsStr)

			// Test other parts of the config
			test.rlsConfig.Descriptors[0].Descriptors[0] = nil
			actualConf.Descriptors[0].Descriptors[0] = nil
			assert.Equal(t, test.rlsConfig, actualConf)
		})
	}
}

func getDummyAPISwagger(apiID string, apiPolicy, res1GetPolicy, res1PostPolicy, res2GetPolicy,
	res2PostPolicy *model.RateLimitPolicy) *model.AdapterInternalAPI {

	res1GetOp := model.NewOperation("GET", nil, nil)
	res1GetOp.RateLimitPolicy = res1GetPolicy
	res1PostOp := model.NewOperation("POST", nil, nil)
	res1PostOp.RateLimitPolicy = res1PostPolicy
	res2GetOp := model.NewOperation("GET", nil, nil)
	res2GetOp.RateLimitPolicy = res2GetPolicy
	res2PostOp := model.NewOperation("POST", nil, nil)
	res2PostOp.RateLimitPolicy = res2PostPolicy

	res1 := model.CreateMinimalDummyResourceForTests("/res1", []*model.Operation{res1GetOp, res1PostOp}, "id1", nil, false)
	res2 := model.CreateMinimalDummyResourceForTests("/res2", []*model.Operation{res2GetOp, res2PostOp}, "id2", nil, false)

	adapterInternalAPI := model.CreateDummyAdapterInternalAPIForTests(fmt.Sprintf("API-%s", apiID), "v1.0.0", fmt.Sprintf("/base-path-%s", apiID), []*model.Resource{
		&res1, &res2,
	})
	adapterInternalAPI.UUID = apiID
	adapterInternalAPI.RateLimitPolicy = apiPolicy
	adapterInternalAPI.OrganizationID = "org1"
	return adapterInternalAPI
}
