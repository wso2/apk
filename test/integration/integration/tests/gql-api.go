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
 *
 */

package tests

import (
	"testing"

	"github.com/wso2/apk/test/integration/integration/utils/http"
	"github.com/wso2/apk/test/integration/integration/utils/suite"
)

func init() {
	//IntegrationTests = append(IntegrationTests, GQLAPI)
}

// DisableAPISecurity test
var GQLAPI = suite.IntegrationTest{
	ShortName:   "GQLAPI",
	Description: "Tests GraphQL API",
	Manifests:   []string{"tests/gql-api.yaml"},
	Test: func(t *testing.T, suite *suite.IntegrationTestSuite) {
		gwAddr := "gql.test.gw.wso2.com:9095"
		token := http.GetTestToken(t)

		testCases := []http.ExpectedResponse{
			{
				Request: http.Request{
					Host:   "gql.test.gw.wso2.com",
					Path:   "/gql/v1",
					Method: "POST",
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
					Body: `{"query":"query{\n    human(id:1000){\n        id\n        name\n    }\n}","variables":{}}`,
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Method: ""},
				},
				Response: http.Response{StatusCode: 200},
			},
			{
				Request: http.Request{
					Host:   "gql.test.gw.wso2.com",
					Path:   "/gql/v1",
					Method: "POST",
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
					Body: `{"query":"query{\n    human(id:1000){\n        id\n        name\n    }\n    droid(id:2000){\n        name\n        friends{\n            name\n            appearsIn\n        }\n    }\n}","variables":{}}`,
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Method: "",
					},
				},
				Response: http.Response{StatusCode: 401},
			}, {
				Request: http.Request{
					Host:   "gql.test.gw.wso2.com",
					Path:   "/gql/v1",
					Method: "POST",
					Headers: map[string]string{
						"Content-Type":  "application/json",
						"Authorization": "Bearer " + token,
					},
					Body: `{"query":"query{\n    human(id:1000){\n        id\n        name\n    }\n    droid(id:2000){\n        name\n        friends{\n            name\n            appearsIn\n        }\n    }\n}","variables":{}}`,
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Method: "",
					},
				},
				Response: http.Response{StatusCode: 200},
			},
		}
		for i := range testCases {
			tc := testCases[i]
			t.Run(tc.GetTestCaseName(i), func(t *testing.T) {
				t.Parallel()
				http.MakeRequestAndExpectEventuallyConsistentResponse(t, suite.RoundTripper, suite.TimeoutConfig, gwAddr, tc)
			})
		}
	},
}
