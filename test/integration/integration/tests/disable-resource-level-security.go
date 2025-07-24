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
	//IntegrationTests = append(IntegrationTests, DisableResourceSecurity)
}

// DisableResourceSecurity test
var DisableResourceSecurity = suite.IntegrationTest{
	ShortName:   "DisableResourceSecurity",
	Description: "Tests API with disabled security",
	Manifests:   []string{"tests/disable-resource-level-security.yaml"},
	Test: func(t *testing.T, suite *suite.IntegrationTestSuite) {
		ns := "gateway-integration-test-infra"
		gwAddr := "disable-resource-security.test.gw.wso2.com:9095"

		testCases := []http.ExpectedResponse{
			{
				Request: http.Request{
					Host: "disable-resource-security.test.gw.wso2.com",
					Path: "/disable-resource-security/v1/users",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path: "/users",
					},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
			{
				Request: http.Request{
					Host: "disable-resource-security.test.gw.wso2.com",
					Path: "/disable-resource-security/v1/orders",
				},
				Response: http.Response{StatusCode: 401},
			},
			{
				Request: http.Request{
					Host: "disable-resource-security.test.gw.wso2.com",
					Path: "/disable-resource-security/users",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path: "/users",
					},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
			{
				Request: http.Request{
					Host: "disable-resource-security.test.gw.wso2.com",
					Path: "/disable-resource-security/orders",
				},
				Response: http.Response{StatusCode: 401},
			},
		}
		for i := range testCases {
			tc := testCases[i]
			// No test token added to the request header
			t.Run(tc.GetTestCaseName(i), func(t *testing.T) {
				t.Parallel()
				http.MakeRequestAndExpectEventuallyConsistentResponse(t, suite.RoundTripper, suite.TimeoutConfig, gwAddr, tc)
			})
		}
	},
}
