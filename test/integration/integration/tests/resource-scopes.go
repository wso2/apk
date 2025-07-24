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
	//IntegrationTests = append(IntegrationTests, ResourceScopes)
}

// ResourceScopes test
var ResourceScopes = suite.IntegrationTest{
	ShortName:   "ResourceScopes",
	Description: "Tests resource with scopes",
	Manifests:   []string{"tests/resource-scopes.yaml"},
	Test: func(t *testing.T, suite *suite.IntegrationTestSuite) {
		ns := "gateway-integration-test-infra"
		gwAddr := "resource-scopes.test.gw.wso2.com:9095"

		tokenWithoutScopes := http.GetTestToken(t)
		tokenWithMatchingScopes := http.GetTestToken(t, "read:pets", "write:pets")
		tokenWithNotMatchingScopes := http.GetTestToken(t, "no:pets")
		testTokens := []string{
			tokenWithoutScopes,
			tokenWithMatchingScopes,
			tokenWithNotMatchingScopes,
			tokenWithoutScopes,
		}

		testCases := []http.ExpectedResponse{
			// Without scopes in both api and token test
			{
				Request: http.Request{
					Host: "resource-scopes.test.gw.wso2.com",
					Path: "/resource-scopes/v1/pet/123",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path: "/pet/123",
					},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
			// With scopes in api and matching scopes in token test
			{
				Request: http.Request{
					Host: "resource-scopes.test.gw.wso2.com",
					Path: "/resource-scopes/v1/pets/findByTags",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path: "/pets/findByTags",
					},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
			// With scopes in api and but not matching scopes in token test
			{
				Request: http.Request{
					Host: "resource-scopes.test.gw.wso2.com",
					Path: "/resource-scopes/v1/pets/findByTags",
				},
				Response: http.Response{StatusCode: 403},
			},
			// With scopes in api but no scopes in token test
			{
				Request: http.Request{
					Host: "resource-scopes.test.gw.wso2.com",
					Path: "/resource-scopes/v1/pets/findByTags",
				},
				Response: http.Response{StatusCode: 403},
			},
		}
		for i := range testCases {
			tc := testCases[i]
			tc.Request.Headers = http.AddBearerTokenToHeader(testTokens[i], tc.Request.Headers)
			t.Run(tc.GetTestCaseName(i), func(t *testing.T) {
				t.Parallel()
				http.MakeRequestAndExpectEventuallyConsistentResponse(t, suite.RoundTripper, suite.TimeoutConfig, gwAddr, tc)
			})
		}
	},
}
