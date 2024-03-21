/*
 *  Copyright (c) 2024, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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
	IntegrationTests = append(IntegrationTests, VerifyOldAPIs)
}

// VerifyOldAPIs test
var VerifyOldAPIs = suite.IntegrationTest{
	ShortName:   "VerifyOldAPIs",
	Description: "Verify Old APIs with different CR versions",
	Manifests:   []string{"tests/verify-old-apis.yaml"},
	Test: func(t *testing.T, suite *suite.IntegrationTestSuite) {
		ns := "gateway-integration-test-infra"
		gwAddr1 := "prod-api.test.gw.wso2.com:9095"
		token := http.GetTestToken(t)

		testCases1 := []http.ExpectedResponse{
			// invoke prod api using prod domain name, invokes prod backend
			{
				Request: http.Request{
					Host: "prod-api.test.gw.wso2.com",
					Path: "/verify-old-apis/v1",
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

		for i := range testCases1 {
			tc := testCases1[i]
			tc.Request.Headers = http.AddBearerTokenToHeader(token, tc.Request.Headers)
			t.Run(tc.GetTestCaseName(i), func(t *testing.T) {
				t.Parallel()
				http.MakeRequestAndExpectEventuallyConsistentResponse(t, suite.RoundTripper, suite.TimeoutConfig, gwAddr1, tc)
			})
		}

	},
}
