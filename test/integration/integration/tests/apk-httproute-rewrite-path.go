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
	"fmt"
	"testing"

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
		// routeNN := types.NamespacedName{Name: "rewrite-path", Namespace: ns}
		// gwNN := types.NamespacedName{Name: "same-namespace", Namespace: ns}
		// TODO: Create a util function to wait until API is deployed
		// TODO: Create a util fnction to get the router service IP address
		gwAddr := "192.168.1.20:9090"
		// kubernetes.GatewayAndHTTPRoutesMustBeAccepted(t, suite.Client, suite.TimeoutConfig, suite.ControllerName, kubernetes.NewGatewayRef(gwNN), routeNN)
		// TODO: Create a util fnction to test key and renew hourly or get a key each time before invoke
		token := "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6IlpUYzJOV1UyTXprM01XRXhNRE0zWVRjeE1HSTFNVGcxWlRCaVl6YzJNakpoWm1Sak1XWTFaQT09In0.eyJhdWQiOiJodHRwOlwvXC9vcmcud3NvMi5hcGltZ3RcL2dhdGV3YXkiLCJzdWIiOiJhZG1pbiIsImlzcyI6Imh0dHBzOlwvXC9ndy53c28yLmNvbVwvdGVzdGtleSIsImtleXR5cGUiOiJQUk9EVUNUSU9OIiwiZXhwIjoxNjc2NTM5OTAzLCJpYXQiOjE2NzY1MzYzMDMsImp0aSI6Ijk4MDdlMTE2LTQ3YTAtNGU1YS1hZmZiLWMyNDU1Y2NhMTA0YSJ9.k0A0Eow_CvbTSXqRn1bxRP-2NOBJ9gon-Scgt4E3qZ9diSioRk68ZR4NdB8qBwG-gcKqQp0hk352P0X6jVij4BiPtpM4UqhUZb82X1tYGmj_UCulcic7Zi4xELXe8lCWQ-ayPfO-n3Tv4LV-jKxef9Opu44qr7m0mVivlmbpckJCK9z_Af2HicjD3oUr5CD5sEcNLGbnd54LNVHB7z1VwDyDiae5N_dmxFbT-1leEqZQPpaBApcRUmZQ99RuaGZYE6p_0HHFwMLyCqVx7_biy2lT5J6F0ery8eQiyO96xTNQ9AbeVwZ8knipyqaCSpeCWqpswmLBxex7xpcljDKRoQ"

		testCases := []http.ExpectedResponse{
			{
				Request: http.Request{
					Host: "urlrewrite.gw.wso2.com",
					Path: "/rewrite-path-api/1.0.0/prefix/one/two",
					Headers: map[string]string{
						"Authorization": fmt.Sprintf("Bearer %s", token),
					},
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
					Headers: map[string]string{
						"Authorization": fmt.Sprintf("Bearer %s", token),
					},
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
