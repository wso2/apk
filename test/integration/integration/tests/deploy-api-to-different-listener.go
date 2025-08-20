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
	//IntegrationTests = append(IntegrationTests, APIDifferentListener)
}

// APIDifferentListener test
var APIDifferentListener = suite.IntegrationTest{
	ShortName:   "APIDifferentListener",
	Description: "An API is deployed to a different listener other than default gateway listener",
	Manifests:   []string{"tests/api-different-listener.yaml"},
	Test: func(t *testing.T, suite *suite.IntegrationTestSuite) {
		ns := "gateway-integration-test-infra"
		gwAddr := "api.am.wso2.com:9095"
		token := http.GetTestToken(t)

		testCases := []http.ExpectedResponse{
			{
				Request: http.Request{
					Host: "api.am.wso2.com",
					Path: "/test-api-with-different-listener/v1.0.0/user/user123/playlist/watch-later",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path: "/user/user123/playlist/watch-later",
					},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
			{
				Request: http.Request{
					Host: "api.am.wso2.com",
					Path: "/test-api-with-different-listener/v1.0.0/user/user123/other-path",
				},
				Response: http.Response{StatusCode: 404},
			},
		}

		negativeTestCases := []http.ExpectedResponse{
			{
				Request: http.Request{
					Host: "gateway-integration-test-infra.test.gw.wso2.com",
					Path: "/test-api-with-different-listener/v1.0.0/user/user123/playlist/watch-later",
				},
				Response:  http.Response{StatusCode: 404},
				Backend:   "infra-backend-v1",
				Namespace: ns,
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
		for i := range negativeTestCases {
			tc := negativeTestCases[i]
			tc.Request.Headers = http.AddBearerTokenToHeader(token, tc.Request.Headers)
			t.Run(tc.GetTestCaseName(i), func(t *testing.T) {
				t.Parallel()
				http.MakeRequestAndExpectEventuallyConsistentResponse(t, suite.RoundTripper, suite.TimeoutConfig, gwAddr, tc)
			})
		}
	},
}
