package tests

import (
	"fmt"
	"testing"

	"github.com/wso2/apk/test/integration/integration/utils/suite"
	"sigs.k8s.io/gateway-api/conformance/utils/http"
	gwapisuite "sigs.k8s.io/gateway-api/conformance/utils/suite"
)

func init() {
	IntegrationTests = append(IntegrationTests, APKHTTPRouteRewritePath)
}

// APKHTTPRouteRewritePath test
var APKHTTPRouteRewritePath = suite.IntegrationTest{
	ShortName:   "APKHTTPRouteRewritePath",
	Description: "An HTTPRoute with path rewrite filter",
	Manifests:   []string{"tests/apk-httproute-rewrite-path.yaml"},
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
		token := "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6IlpUYzJOV1UyTXprM01XRXhNRE0zWVRjeE1HSTFNVGcxWlRCaVl6YzJNakpoWm1Sak1XWTFaQT09In0.eyJhdWQiOiJodHRwOlwvXC9vcmcud3NvMi5hcGltZ3RcL2dhdGV3YXkiLCJzdWIiOiJhZG1pbiIsImlzcyI6Imh0dHBzOlwvXC9ndy53c28yLmNvbVwvdGVzdGtleSIsImtleXR5cGUiOiJQUk9EVUNUSU9OIiwiZXhwIjoxNjc2NDY0Mzk5LCJpYXQiOjE2NzY0NjA3OTksImp0aSI6IjA0ZDU5ZTY0LTcxMjItNGQ1ZC05N2ZjLTYzMTU0OTY3ZWMzOSJ9.z96YjfIw7lLXmhKgBPEN0kYWqjVBDBDW1p7nnV6O5OJ7ta97kwAe29FcimcKY1wBqm45yK2Yi68ANvVIFS36Og49K_cOT94mKdFkOCLzXM12jIo9qDvrT6ao-Na6yd0cG3-jpVMG8xzbYTCkTZd-FNoJKY_xo7CxsjnQnhPL-MW53j6UpoG17v5xL_O6Y6LS7FBu_1vSqKOuee1aGMXpzyDH-uxc5DuD750EfFD3tiaE9nuFzsJE4f4LvjHzApgHWiIsIb_c0KU4I79MSd-F-9-KuGIQaCpZ_4Qom9aS-uxf6uTN8G2bFkLoRC8ZOquJL3x8DMRn3cPYWXHlQQtuJg"

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
