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
)

func TestGenerateIdentifierForAPIWithUUID(t *testing.T) {
	setupInternalMemoryMapsWithTestSamples()
	tests := []struct {
		name  string
		uuid  string
		vhost string
	}{
		{
			name:  "Get_identifier_from_uuid_and_vhost",
			uuid:  "e2cb0839-700b-4226-8239-eead31353f19",
			vhost: "org2.foo.com",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			identifier := GenerateIdentifierForAPIWithUUID(test.vhost, test.uuid)
			if identifier != test.vhost+":"+test.uuid {
				t.Errorf("expected identifier %v but found %v", test.vhost+":"+test.uuid, identifier)
			}
		})
	}
}

func setupInternalMemoryMapsWithTestSamples() {
	apiToVhostsMap = map[string]map[string]struct{}{
		// The same API uuid is deployed in two org with two gateway environments
		"111-PetStore-org1": {"org1.wso2.com": void, "org2.foo.com": void},
		"333-Pizza-org1":    {"org1.foo.com": void, "org2.foo.com": void, "org2.wso2.com": void},
	}
	apiUUIDToGatewayToVhosts = map[string]map[string]string{
		// PetStore:v1 in Org1
		"111-PetStore-org1": {
			"Default":   "org1.wso2.com",
			"us-region": "org1.wso2.com",
		},
		// PetStore:v1 in Org2
		"222-PetStore-org2": {
			"Default": "org2.foo.com",
		},
		// Pizza:v1 in Org1
		"333-Pizza-org1": {
			"us-region": "org1.foo.com",
		},
		// Pizza:v1 in Org2
		"444-Pizza-org2": {
			"Default":   "org2.foo.com",
			"us-region": "org2.wso2.com",
		},
	}
}
