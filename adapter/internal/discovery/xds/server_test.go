/*
 *  Copyright (c) 2022, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wso2/apk/adapter/config"
	"github.com/wso2/apk/adapter/internal/loggers"
	"github.com/wso2/apk/adapter/internal/oasparser/model"
	semantic_version "github.com/wso2/apk/adapter/pkg/semanticversion"
)

func TestOrgMapUpdates(t *testing.T) {
	orgAPIMap = make(map[string]map[string]*EnvoyInternalAPI)
	orgIDLatestAPIVersionMap = make(map[string]map[string]map[string]semantic_version.SemVersion)
	conf := config.ReadConfigs()
	conf.Envoy.EnableIntelligentRouting = true

	api1uuid := &model.AdapterInternalAPI{
		UUID:           "api-1-uuid",
		EnvType:        "prod",
		OrganizationID: "org-1",
	}
	api1uuid.SetName("api-1")
	api1uuid.SetVersion("v1.0.0")
	api1sanduuid := &model.AdapterInternalAPI{
		UUID:           "api-1-uuid",
		EnvType:        "sand",
		OrganizationID: "org-1",
	}
	api1sanduuid.SetName("api-1")
	api1sanduuid.SetVersion("v1.0.1")
	api2uuid := &model.AdapterInternalAPI{
		UUID:           "api-2-uuid",
		EnvType:        "prod",
		OrganizationID: "org-1",
	}
	api2uuid.SetName("api-2")
	api2uuid.SetVersion("v0.0.1")
	ptrOne := new(int)
	*ptrOne = 1
	tests := []struct {
		name                             string
		vHosts                           []string
		labels                           map[string]struct{}
		listeners                        []string
		adapterInternalAPI               *model.AdapterInternalAPI
		EnvType                          string
		action                           string
		deletedvHosts                    []string
		expectedOrgAPIMap                map[string]map[string]*EnvoyInternalAPI
		expectedOrgIDLatestAPIVersionMap map[string]map[string]map[string]semantic_version.SemVersion
	}{
		{
			name:               "Test creating first prod api",
			vHosts:             []string{"prod1.gw.abc.com", "prod2.gw.abc.com"},
			labels:             map[string]struct{}{"default": {}},
			listeners:          []string{"httpslistener"},
			adapterInternalAPI: api1uuid,
			EnvType:            "prod",
			action:             "CREATE",
			expectedOrgAPIMap: map[string]map[string]*EnvoyInternalAPI{
				"org-1": {
					"prod1.gw.abc.com:api-1-uuid": &EnvoyInternalAPI{
						envoyLabels:        map[string]struct{}{"default": {}},
						adapterInternalAPI: api1uuid,
					},
					"prod2.gw.abc.com:api-1-uuid": &EnvoyInternalAPI{
						envoyLabels:        map[string]struct{}{"default": {}},
						adapterInternalAPI: api1uuid,
					},
				},
			},
			expectedOrgIDLatestAPIVersionMap: map[string]map[string]map[string]semantic_version.SemVersion{
				"org-1": {
					"prod1.gw.abc.com:api-1": {
						"v1": semantic_version.SemVersion{
							Version: "v1.0.0",
							Major:   1,
							Minor:   0,
							Patch:   new(int),
						},
						"v1.0": semantic_version.SemVersion{
							Version: "v1.0.0",
							Major:   1,
							Minor:   0,
							Patch:   new(int),
						},
					},
					"prod2.gw.abc.com:api-1": {
						"v1": semantic_version.SemVersion{
							Version: "v1.0.0",
							Major:   1,
							Minor:   0,
							Patch:   new(int),
						},
						"v1.0": semantic_version.SemVersion{
							Version: "v1.0.0",
							Major:   1,
							Minor:   0,
							Patch:   new(int),
						},
					},
				},
			},
		},
		{
			name:               "Test creating first sand api",
			vHosts:             []string{"sand3.gw.abc.com", "sand4.gw.abc.com"},
			labels:             map[string]struct{}{"default": {}},
			listeners:          []string{"httpslistener"},
			adapterInternalAPI: api1sanduuid,
			EnvType:            "sand",
			action:             "UPDATE",
			expectedOrgAPIMap: map[string]map[string]*EnvoyInternalAPI{
				"org-1": {
					"prod1.gw.abc.com:api-1-uuid": &EnvoyInternalAPI{
						envoyLabels:        map[string]struct{}{"default": {}},
						adapterInternalAPI: api1uuid,
					},
					"prod2.gw.abc.com:api-1-uuid": &EnvoyInternalAPI{
						envoyLabels:        map[string]struct{}{"default": {}},
						adapterInternalAPI: api1uuid,
					},
					"sand3.gw.abc.com:api-1-uuid": &EnvoyInternalAPI{
						envoyLabels:        map[string]struct{}{"default": {}},
						adapterInternalAPI: api1sanduuid,
					},
					"sand4.gw.abc.com:api-1-uuid": &EnvoyInternalAPI{
						envoyLabels:        map[string]struct{}{"default": {}},
						adapterInternalAPI: api1sanduuid,
					},
				},
			},
			expectedOrgIDLatestAPIVersionMap: map[string]map[string]map[string]semantic_version.SemVersion{
				"org-1": {
					"prod1.gw.abc.com:api-1": {
						"v1": semantic_version.SemVersion{
							Version: "v1.0.0",
							Major:   1,
							Minor:   0,
							Patch:   new(int),
						},
						"v1.0": semantic_version.SemVersion{
							Version: "v1.0.0",
							Major:   1,
							Minor:   0,
							Patch:   new(int),
						},
					},
					"prod2.gw.abc.com:api-1": {
						"v1": semantic_version.SemVersion{
							Version: "v1.0.0",
							Major:   1,
							Minor:   0,
							Patch:   new(int),
						},
						"v1.0": semantic_version.SemVersion{
							Version: "v1.0.0",
							Major:   1,
							Minor:   0,
							Patch:   new(int),
						},
					},
					"sand3.gw.abc.com:api-1": {
						"v1": semantic_version.SemVersion{
							Version: "v1.0.1",
							Major:   1,
							Minor:   0,
							Patch:   ptrOne,
						},
						"v1.0": semantic_version.SemVersion{
							Version: "v1.0.1",
							Major:   1,
							Minor:   0,
							Patch:   ptrOne,
						},
					},
					"sand4.gw.abc.com:api-1": {
						"v1": semantic_version.SemVersion{
							Version: "v1.0.1",
							Major:   1,
							Minor:   0,
							Patch:   ptrOne,
						},
						"v1.0": semantic_version.SemVersion{
							Version: "v1.0.1",
							Major:   1,
							Minor:   0,
							Patch:   ptrOne,
						},
					},
				},
			},
		},
		{
			name:               "Test creating second prod api",
			vHosts:             []string{"prod1.gw.pqr.com", "prod2.gw.pqr.com"},
			labels:             map[string]struct{}{"default": {}},
			listeners:          []string{"httpslistener"},
			adapterInternalAPI: api2uuid,
			EnvType:            "prod",
			action:             "CREATE",
			expectedOrgAPIMap: map[string]map[string]*EnvoyInternalAPI{
				"org-1": {
					"prod1.gw.abc.com:api-1-uuid": &EnvoyInternalAPI{
						envoyLabels:        map[string]struct{}{"default": {}},
						adapterInternalAPI: api1uuid,
					},
					"prod2.gw.abc.com:api-1-uuid": &EnvoyInternalAPI{
						envoyLabels:        map[string]struct{}{"default": {}},
						adapterInternalAPI: api1uuid,
					},
					"sand3.gw.abc.com:api-1-uuid": &EnvoyInternalAPI{
						envoyLabels:        map[string]struct{}{"default": {}},
						adapterInternalAPI: api1sanduuid,
					},
					"sand4.gw.abc.com:api-1-uuid": &EnvoyInternalAPI{
						envoyLabels:        map[string]struct{}{"default": {}},
						adapterInternalAPI: api1sanduuid,
					},
					"prod1.gw.pqr.com:api-2-uuid": &EnvoyInternalAPI{
						envoyLabels:        map[string]struct{}{"default": {}},
						adapterInternalAPI: api2uuid,
					},
					"prod2.gw.pqr.com:api-2-uuid": &EnvoyInternalAPI{
						envoyLabels:        map[string]struct{}{"default": {}},
						adapterInternalAPI: api2uuid,
					},
				},
			},
			expectedOrgIDLatestAPIVersionMap: map[string]map[string]map[string]semantic_version.SemVersion{
				"org-1": {
					"prod1.gw.abc.com:api-1": {
						"v1": semantic_version.SemVersion{
							Version: "v1.0.0", // fix it should still be v1.0.0
							Major:   1,
							Minor:   0,
							Patch:   new(int),
						},
						"v1.0": semantic_version.SemVersion{
							Version: "v1.0.0",
							Major:   1,
							Minor:   0,
							Patch:   new(int),
						},
					},
					"prod2.gw.abc.com:api-1": {
						"v1": semantic_version.SemVersion{
							Version: "v1.0.0",
							Major:   1,
							Minor:   0,
							Patch:   new(int),
						},
						"v1.0": semantic_version.SemVersion{
							Version: "v1.0.0",
							Major:   1,
							Minor:   0,
							Patch:   new(int),
						},
					},
					"sand3.gw.abc.com:api-1": {
						"v1": semantic_version.SemVersion{
							Version: "v1.0.1",
							Major:   1,
							Minor:   0,
							Patch:   ptrOne,
						},
						"v1.0": semantic_version.SemVersion{
							Version: "v1.0.1",
							Major:   1,
							Minor:   0,
							Patch:   ptrOne,
						},
					},
					"sand4.gw.abc.com:api-1": {
						"v1": semantic_version.SemVersion{
							Version: "v1.0.1",
							Major:   1,
							Minor:   0,
							Patch:   ptrOne,
						},
						"v1.0": semantic_version.SemVersion{
							Version: "v1.0.1",
							Major:   1,
							Minor:   0,
							Patch:   ptrOne,
						},
					},
					"prod1.gw.pqr.com:api-2": {
						"v0": semantic_version.SemVersion{
							Version: "v0.0.1",
							Major:   0,
							Minor:   0,
							Patch:   ptrOne,
						},
						"v0.0": semantic_version.SemVersion{
							Version: "v0.0.1",
							Major:   0,
							Minor:   0,
							Patch:   ptrOne,
						},
					},
					"prod2.gw.pqr.com:api-2": {
						"v0": semantic_version.SemVersion{
							Version: "v0.0.1",
							Major:   0,
							Minor:   0,
							Patch:   ptrOne,
						},
						"v0.0": semantic_version.SemVersion{
							Version: "v0.0.1",
							Major:   0,
							Minor:   0,
							Patch:   ptrOne,
						},
					},
				},
			},
		},
		// {
		// 	name:               "Test updating first prod api 1 with new vhosts",
		// 	vHosts:             []string{"prod3.gw.abc.com", "prod4.gw.abc.com"},
		// 	labels:             map[string]struct{}{"default": {}},
		// 	listeners:          []string{"httpslistener"},
		// 	adapterInternalAPI: api1uuid,
		// 	action:             "UPDATE",
		// 	expectedOrgAPIMap: map[string]map[string]*EnvoyInternalAPI{
		// 		"org-1": {
		// 			"prod3.gw.abc.com:api-1-uuid": &EnvoyInternalAPI{
		// 				envoyLabels:        map[string]struct{}{"default": {}},
		// 				adapterInternalAPI: api1uuid,
		// 			},
		// 			"prod4.gw.abc.com:api-1-uuid": &EnvoyInternalAPI{
		// 				envoyLabels:        map[string]struct{}{"default": {}},
		// 				adapterInternalAPI: api1uuid,
		// 			},
		// 			"sand3.gw.abc.com:api-1-uuid": &EnvoyInternalAPI{
		// 				envoyLabels:        map[string]struct{}{"default": {}},
		// 				adapterInternalAPI: api1sanduuid,
		// 			},
		// 			"sand4.gw.abc.com:api-1-uuid": &EnvoyInternalAPI{
		// 				envoyLabels:        map[string]struct{}{"default": {}},
		// 				adapterInternalAPI: api1sanduuid,
		// 			},
		// 			"prod1.gw.pqr.com:api-2-uuid": &EnvoyInternalAPI{
		// 				envoyLabels:        map[string]struct{}{"default": {}},
		// 				adapterInternalAPI: api2uuid,
		// 			},
		// 			"prod2.gw.pqr.com:api-2-uuid": &EnvoyInternalAPI{
		// 				envoyLabels:        map[string]struct{}{"default": {}},
		// 				adapterInternalAPI: api2uuid,
		// 			},
		// 		},
		// 	},
		// },
		// {
		// 	name:               "Test deleting api 1 both prod and sand",
		// 	labels:             map[string]struct{}{"default": {}},
		// 	listeners:          []string{"httpslistener"},
		// 	adapterInternalAPI: api1uuid,
		// 	action:             "DELETE",
		// 	expectedOrgAPIMap: map[string]map[string]*EnvoyInternalAPI{
		// 		"org-1": {
		// 			"prod1.gw.pqr.com:api-2-uuid": &EnvoyInternalAPI{
		// 				envoyLabels:        map[string]struct{}{"default": {}},
		// 				adapterInternalAPI: api2uuid,
		// 			},
		// 			"prod2.gw.pqr.com:api-2-uuid": &EnvoyInternalAPI{
		// 				envoyLabels:        map[string]struct{}{"default": {}},
		// 				adapterInternalAPI: api2uuid,
		// 			},
		// 		},
		// 	},
		// },
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			switch test.action {
			case "CREATE":
				loggers.LoggerAPI.Infof("Creating API: %v", test.adapterInternalAPI.UUID)
				for label := range test.labels {
					SanitizeGateway(label, true)
				}
				RemoveAPIFromAllInternalMaps(test.adapterInternalAPI.UUID)
				UpdateAPICache(test.vHosts, test.labels, test.listeners[0], "httpslistener", test.adapterInternalAPI)
			case "UPDATE":
				loggers.LoggerAPI.Infof("Updating API: %v", test.adapterInternalAPI.UUID)
				for label := range test.labels {
					SanitizeGateway(label, true)
				}
				UpdateAPICache(test.vHosts, test.labels, test.listeners[0], "httpslistener", test.adapterInternalAPI)
			case "DELETE":
				loggers.LoggerAPI.Infof("Deleting API: %v", test.adapterInternalAPI.UUID)
				DeleteAPI(test.adapterInternalAPI.UUID, test.labels)
			}
			assert.Equal(t, len(test.expectedOrgAPIMap), len(orgAPIMap), "orgAPIMap length is different, expected: %v but found: %v",
				test.expectedOrgAPIMap, orgAPIMap)
			for orgID, orgAPIs := range test.expectedOrgAPIMap {
				if orgAPI, ok := orgAPIMap[orgID]; !ok {
					t.Errorf("orgAPIMap has no entry with the organization: %s", orgID)
				} else {
					assert.Equal(t, len(test.expectedOrgAPIMap[orgID]), len(orgAPIs), "orgAPI length is different, expected: %v but found: %v",
						len(test.expectedOrgAPIMap[orgID]), len(orgAPIs))
					for vuuid, expectedAPI := range test.expectedOrgAPIMap[orgID] {
						if actualAPI, ok := orgAPI[vuuid]; !ok {
							t.Errorf("orgAPIMap has not updated with new entry with the key: %s, %v",
								vuuid, orgAPIMap)
						} else {
							assert.Equal(t, expectedAPI.adapterInternalAPI.UUID, actualAPI.adapterInternalAPI.UUID, "Not expected API UUID found, expected: %v but found: %v",
								expectedAPI.adapterInternalAPI.UUID, actualAPI.adapterInternalAPI.UUID)
							assert.Equal(t, expectedAPI.adapterInternalAPI.EnvType, actualAPI.adapterInternalAPI.EnvType, "Not expected API EnvType found, expected: %v but found: %v",
								expectedAPI.adapterInternalAPI.EnvType, actualAPI.adapterInternalAPI.EnvType)
							assert.Equal(t, expectedAPI.adapterInternalAPI.OrganizationID, actualAPI.adapterInternalAPI.OrganizationID, "Not expected API OrganizationID found, expected: %v but found: %v",
								expectedAPI.adapterInternalAPI.OrganizationID, actualAPI.adapterInternalAPI.OrganizationID)
						}

					}
				}
			}
			assert.Equal(t, len(test.expectedOrgIDLatestAPIVersionMap), len(orgIDLatestAPIVersionMap), "orgIDLatestAPIVersionMap length is different, expected: %v but found: %v",
				len(test.expectedOrgIDLatestAPIVersionMap), len(orgIDLatestAPIVersionMap))
			for orgID, orgAPIs := range test.expectedOrgIDLatestAPIVersionMap {
				if orgAPI, ok := orgIDLatestAPIVersionMap[orgID]; !ok {
					t.Errorf("orgIDLatestAPIVersionMap has no entry with the organization: %s, %v", orgID, orgIDLatestAPIVersionMap)
				} else {
					assert.Equal(t, len(test.expectedOrgIDLatestAPIVersionMap[orgID]), len(orgAPIs), "orgAPI length is different, expected: %v but found: %v",
						len(test.expectedOrgIDLatestAPIVersionMap[orgID]), len(orgAPIs))
					for vuuid, expectedAPI := range test.expectedOrgIDLatestAPIVersionMap[orgID] {
						if actualAPI, ok := orgAPI[vuuid]; !ok {
							t.Errorf("orgIDLatestAPIVersionMap has not updated with new entry with the key for %s, %v",
								vuuid, orgIDLatestAPIVersionMap)
						} else {
							assert.Equal(t, len(expectedAPI), len(actualAPI), "orgAPI for %v length is different, expected: %v but found: %v",
								vuuid, len(expectedAPI), len(actualAPI))
							for version, expectedVersion := range expectedAPI {
								if actualVersion, ok := actualAPI[version]; !ok {
									t.Errorf("orgIDLatestAPIVersionMap has not updated with new entry with the key for %v: %s, %v",
										vuuid, version, orgIDLatestAPIVersionMap)
								} else {
									assert.Equal(t, expectedVersion.Version, actualVersion.Version, "Not expected API Version found for %v in %v, expected: %v but found: %v",
										vuuid, version, expectedVersion.Version, actualVersion.Version)
									assert.Equal(t, expectedVersion.Major, actualVersion.Major, "Not expected API Major found for %v in %v, expected: %v but found: %v",
										vuuid, version, expectedVersion.Major, actualVersion.Major)
									assert.Equal(t, expectedVersion.Minor, actualVersion.Minor, "Not expected API Minor found for %v in %v, expected: %v but found: %v",
										vuuid, version, expectedVersion.Minor, actualVersion.Minor)
									assert.Equal(t, *expectedVersion.Patch, *actualVersion.Patch, "Not expected API Patch found for %v in %v, expected: %v but found: %v",
										vuuid, version, *expectedVersion.Patch, *actualVersion.Patch)
								}
							}
						}
					}
				}
			}
		})
	}
}
