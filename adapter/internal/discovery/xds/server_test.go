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
	"github.com/wso2/apk/adapter/internal/oasparser/model"
)

func TestUpdateAPICache(t *testing.T) {
	tests := []struct {
		name               string
		vHosts             []string
		labels             []string
		listeners          []string
		adapterInternalAPI model.AdapterInternalAPI
		EnvType            string
		action             string
		deletedvHosts      []string
	}{
		{
			name:      "Test creating first prod api",
			vHosts:    []string{"prod1.gw.abc.com", "prod2.gw.abc.com"},
			labels:    []string{"default"},
			listeners: []string{"httpslistener"},
			adapterInternalAPI: model.AdapterInternalAPI{
				UUID:           "api-1-uuid",
				EnvType:        "prod",
				OrganizationID: "org-1",
			},
			EnvType: "prod",
			action:  "CREATE",
		},
		{
			name:      "Test creating first sand api",
			vHosts:    []string{"sand3.gw.abc.com", "sand4.gw.abc.com"},
			labels:    []string{"default"},
			listeners: []string{"httpslistener"},
			adapterInternalAPI: model.AdapterInternalAPI{
				UUID:           "app-1-uuid",
				EnvType:        "sand",
				OrganizationID: "org-1",
			},
			EnvType: "sand",
			action:  "CREATE",
		},
		{
			name:      "Test creating second prod api",
			vHosts:    []string{"prod1.gw.pqr.com", "prod2.gw.pqr.com"},
			labels:    []string{"default"},
			listeners: []string{"httpslistener"},
			adapterInternalAPI: model.AdapterInternalAPI{
				UUID:           "api-2-uuid",
				EnvType:        "prod",
				OrganizationID: "org-2",
			},
			EnvType: "prod",
			action:  "CREATE",
		},
		{
			name:      "Test updating first prod api 1 with new vhosts",
			vHosts:    []string{"prod1.gw.abc.com", "prod2.gw.abc.com"},
			labels:    []string{"default"},
			listeners: []string{"httpslistener"},
			adapterInternalAPI: model.AdapterInternalAPI{
				UUID:           "api-1-uuid",
				EnvType:        "prod",
				OrganizationID: "org-1",
			},
			action: "UPDATE",
		},
		{
			name:      "Test deleting api 1 both prod and sand",
			labels:    []string{"default"},
			listeners: []string{"httpslistener"},
			adapterInternalAPI: model.AdapterInternalAPI{
				UUID:           "app-1-uuid",
				OrganizationID: "org-1",
			},
			action: "DELETE",
			deletedvHosts: []string{"prod1.gw.abc.com", "prod2.gw.abc.com",
				"sand3.gw.abc.com", "sand4.gw.abc.com"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			switch test.action {
			case "CREATE":
			case "UPDATE":
				for _, label := range test.labels {
					SanitizeGateway(label, true)
				}
				UpdateAPICache(test.vHosts, test.labels, test.listeners[0], "httpslistener", test.adapterInternalAPI)
				identifier := GetvHostsIdentifier(test.adapterInternalAPI.UUID, "prod")
				actualvHosts, ok := orgIDAPIvHostsMap[test.adapterInternalAPI.OrganizationID][identifier]
				if !ok {
					t.Errorf("orgIDAPIvHostsMap has not updated with new entry with the key: %s, %v",
						identifier, orgIDAPIvHostsMap)
				}
				assert.Equal(t, actualvHosts, test.vHosts, "Not expected vHosts found, expected: %v but found: %v",
					test.vHosts, actualvHosts)
				for _, vhsot := range actualvHosts {
					testExistsInMapping(t, orgAPIMap[test.adapterInternalAPI.OrganizationID],
						GenerateIdentifierForAPIWithUUID(vhsot, test.adapterInternalAPI.UUID), true)
				}
			case "DELETE":
				DeleteAPICREvent(test.labels, test.adapterInternalAPI.UUID, test.adapterInternalAPI.OrganizationID)
				prodIdentifier := GetvHostsIdentifier(test.adapterInternalAPI.UUID, "prod")
				sandIdentifier := GetvHostsIdentifier(test.adapterInternalAPI.UUID, "sand")
				_, prodExists := orgIDAPIvHostsMap[test.adapterInternalAPI.OrganizationID][prodIdentifier]
				_, sandExists := orgIDAPIvHostsMap[test.adapterInternalAPI.OrganizationID][sandIdentifier]
				if prodExists {
					t.Errorf("orgIDAPIvHostsMap has a mapping for prod after api deletion")
				}
				if sandExists {
					t.Errorf("orgIDAPIvHostsMap has a mapping for sand after api deletion")
				}
				for _, vhsot := range test.deletedvHosts {
					testExistsInMapping(t, orgAPIMap[test.adapterInternalAPI.OrganizationID],
						GenerateIdentifierForAPIWithUUID(vhsot, test.adapterInternalAPI.UUID), false)
				}
			}
		})
	}
}

func testExistsInMapping[V any, M map[string]V](t *testing.T, mapping M, key string, checkExists bool) {
	_, ok := mapping[key]
	if checkExists {
		if !ok {
			t.Errorf("Not found mapping for key %s in map %v", key, mapping)
		}
	} else {
		if ok {
			t.Errorf("Found mapping for key %s in map %v", key, mapping)
		}
	}
}
