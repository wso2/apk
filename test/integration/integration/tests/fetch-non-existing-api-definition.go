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
	//IntegrationTests = append(IntegrationTests, FetchNonExistingAPIDefinition)
}

// FetchNonExistingAPIDefinition test
var FetchNonExistingAPIDefinition = suite.IntegrationTest{
	ShortName:   "FetchNonExistingAPIDefinition",
	Description: "Tests an invocation on non existing api definition from the api definition route",
	Manifests:   []string{"tests/fetch-non-existing-api-definition.yaml"},
	Test: func(t *testing.T, suite *suite.IntegrationTestSuite) {
		ns := "gateway-integration-test-infra"
		gwAddr := "fetch-non-existing-api-definition.test.gw.wso2.com:9095"
		token := http.GetTestToken(t)

		testCases := []http.ExpectedResponse{
			{
				Request: http.Request{
					Host: "fetch-non-existing-api-definition.gw.wso2.com",
					Path: "/fetch-non-existing-api-definition/v1.0.0/api-definition",
					Headers: map[string]string{
						"content-type": "application/json",
					},
					Method: "GET",
				},
				Response: http.Response{
					StatusCode: 404,
				},
				Backend:      "infra-backend-v1",
				Namespace:    ns,
				TestCaseName: "FetchAPIDefinition",
			},
		}

		for i := range testCases {
			tc := testCases[i]
			tc.Request.Headers = http.AddBearerTokenToHeader(token, tc.Request.Headers)
			t.Run(tc.GetTestCaseName(i), func(t *testing.T) {
				t.Parallel()
				http.MakeRequestAndExpectEventuallyConsistentResponse(t, suite.RoundTripper, suite.TimeoutConfig, gwAddr, tc)
			})
		}
	},
}
