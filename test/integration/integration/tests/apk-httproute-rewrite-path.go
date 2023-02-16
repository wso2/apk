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

	"github.com/wso2/apk/test/integration/integration/utils/kubernetes"
	"github.com/wso2/apk/test/integration/integration/utils/suite"
	"sigs.k8s.io/gateway-api/conformance/utils/http"
	gwapisuite "sigs.k8s.io/gateway-api/conformance/utils/suite"
)

func init() {
	IntegrationTests = append(IntegrationTests, HTTPRouteRewritePath)
}

// HTTPRouteRewritePath test
var HTTPRouteRewritePath = suite.IntegrationTest{
	ShortName:   "HTTPRouteRewritePath",
	Description: "An HTTPRoute with path rewrite filter",
	Manifests:   []string{"tests/httproute-rewrite-path.yaml"},
	Features:    []gwapisuite.SupportedFeature{gwapisuite.SupportHTTPRoutePathRewrite},
	Test: func(t *testing.T, suite *suite.IntegrationTestSuite) {
		ns := "gateway-conformance-infra"
		gwAddr, _ := kubernetes.WaitForGatewayAddress(t, suite.Client, suite.TimeoutConfig)

		testCases := []http.ExpectedResponse{
			{
				Request: http.Request{
					Host: "urlrewrite.gw.wso2.com",
					Path: "/rewrite-path-api/1.0.0/prefix/one/two",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path: "/one/two",
					},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
			{
				Request: http.Request{
					Host: "urlrewrite.gw.wso2.com",
					Path: "/rewrite-path-api/1.0.0/full/one/two",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path: "/one",
					},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
		}
		for i := range testCases {
			// Declare tc here to avoid loop variable
			// reuse issues across parallel tests.
			tc := testCases[i]
			t.Run(tc.GetTestCaseName(i), func(t *testing.T) {
				t.Parallel()
				http.MakeRequestAndExpectEventuallyConsistentResponse(t, suite.RoundTripper, suite.TimeoutConfig, gwAddr, tc)
			})
		}
	},
}
