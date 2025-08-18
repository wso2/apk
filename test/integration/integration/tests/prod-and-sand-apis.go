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
	// //IntegrationTests = append(IntegrationTests, ProdAndSandAPIs)
}

// ProdAndSandAPIs test
var ProdAndSandAPIs = suite.IntegrationTest{
	ShortName:   "ProdAndSandAPIs",
	Description: "Tests API with disabled security",
	Manifests:   []string{"tests/prod-and-sand-apis.yaml"},
	Test: func(t *testing.T, suite *suite.IntegrationTestSuite) {
		ns := "gateway-integration-test-infra"
		gwAddr1 := "prod-api.test.gw.wso2.com:9095"
		gwAddr2 := "sand-api.test.gw.wso2.com:9095"
		token := http.GetTestToken(t)

		testCases1 := []http.ExpectedResponse{
			// invoke prod api using prod domain name, invokes prod backend
			{
				Request: http.Request{
					Host: "prod-api.test.gw.wso2.com",
					Path: "/prod-sand-test-api/v1",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path: "/",
					},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
		}
		testCases2 := []http.ExpectedResponse{
			// invoke sand api using sand domain name, invokes sand backend
			{
				Request: http.Request{
					Host: "sand-api.test.gw.wso2.com",
					Path: "/prod-sand-test-api/v1",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path: "/",
					},
				},
				Backend:   "infra-backend-v2",
				Namespace: ns,
			},
		}

		for i := range testCases1 {
			tc := testCases1[i]
			tc.Request.Headers = http.AddBearerTokenToHeader(token, tc.Request.Headers)
			t.Run(tc.GetTestCaseName(i), func(t *testing.T) {
				t.Parallel()
				http.MakeRequestAndExpectEventuallyConsistentResponse(t, suite.RoundTripper, suite.TimeoutConfig, gwAddr1, tc)
			})
		}
		for i := range testCases2 {
			tc := testCases2[i]
			tc.Request.Headers = http.AddBearerTokenToHeader(token, tc.Request.Headers)
			t.Run(tc.GetTestCaseName(i), func(t *testing.T) {
				t.Parallel()
				http.MakeRequestAndExpectEventuallyConsistentResponse(t, suite.RoundTripper, suite.TimeoutConfig, gwAddr2, tc)
			})
		}
	},
}
