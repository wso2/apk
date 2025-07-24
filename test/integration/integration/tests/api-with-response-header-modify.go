/*
Copyright 2022 The Kubernetes Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package tests

import (
	"testing"

	"github.com/wso2/apk/test/integration/integration/utils/http"
	"github.com/wso2/apk/test/integration/integration/utils/suite"
)

func init() {
	//IntegrationTests = append(IntegrationTests, APIWithResoponseHeaderModify)
}

// APIWithResoponseHeaderModify test
var APIWithResoponseHeaderModify = suite.IntegrationTest{
	ShortName:   "APIWithResoponseHeaderModify",
	Description: "An API with response header modify",
	Manifests:   []string{"tests/api-with-response-header-modify.yaml"},
	Test: func(t *testing.T, suite *suite.IntegrationTestSuite) {
		ns := "gateway-integration-test-infra"
		gwAddr := "gateway-integration-test-infra.test.gw.wso2.com:9095"
		token := http.GetTestToken(t)

		testCases := []http.ExpectedResponse{
			{
				Request: http.Request{
					Host: "gateway-integration-test-infra.test.gw.wso2.com",
					Path: "/response-header-modify/1.0.0/set",
				},
				BackendSetResponseHeaders: map[string]string{
					"Some-Other-Header": "val",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path: "/set",
					},
				},
				Response: http.Response{
					Headers: map[string]string{
						"Some-Other-Header": "val",
						"X-Header-Set":      "set-overwrites-values",
					},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
			{
				Request: http.Request{
					Host: "gateway-integration-test-infra.test.gw.wso2.com",
					Path: "/response-header-modify/1.0.0/set",
				},
				BackendSetResponseHeaders: map[string]string{
					"Some-Other-Header": "val",
					"X-Header-Set":      "some-other-value",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path: "/set",
					},
				},
				Response: http.Response{
					Headers: map[string]string{
						"Some-Other-Header": "val",
						"X-Header-Set":      "set-overwrites-values",
					},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
			{
				Request: http.Request{
					Host: "gateway-integration-test-infra.test.gw.wso2.com",
					Path: "/response-header-modify/1.0.0/add",
				},
				BackendSetResponseHeaders: map[string]string{
					"Some-Other-Header": "val",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path: "/add",
					},
				},
				Response: http.Response{
					Headers: map[string]string{
						"Some-Other-Header": "val",
						"X-Header-Add":      "add-appends-values",
					},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
			{
				Request: http.Request{
					Host: "gateway-integration-test-infra.test.gw.wso2.com",
					Path: "/response-header-modify/1.0.0/add",
				},
				BackendSetResponseHeaders: map[string]string{
					"Some-Other-Header": "val",
					"X-Header-Add":      "some-other-value",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path: "/add",
					},
				},
				Response: http.Response{
					Headers: map[string]string{
						"Some-Other-Header": "val",
						"X-Header-Add":      "add-appends-values",
					},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			}, {
				Request: http.Request{
					Host: "gateway-integration-test-infra.test.gw.wso2.com",
					Path: "/response-header-modify/1.0.0/remove",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path: "/remove",
					},
				},
				BackendSetResponseHeaders: map[string]string{
					"X-Header-Remove": "val",
				},
				Response: http.Response{
					AbsentHeaders: []string{"X-Header-Remove"},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},

			// Check default api path
			{
				Request: http.Request{
					Host: "gateway-integration-test-infra.test.gw.wso2.com",
					Path: "/response-header-modify/set",
				},
				BackendSetResponseHeaders: map[string]string{
					"Some-Other-Header": "val",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path: "/set",
					},
				},
				Response: http.Response{
					Headers: map[string]string{
						"Some-Other-Header": "val",
						"X-Header-Set":      "set-overwrites-values",
					},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
			{
				Request: http.Request{
					Host: "gateway-integration-test-infra.test.gw.wso2.com",
					Path: "/response-header-modify/set",
				},
				BackendSetResponseHeaders: map[string]string{
					"Some-Other-Header": "val",
					"X-Header-Set":      "some-other-value",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path: "/set",
					},
				},
				Response: http.Response{
					Headers: map[string]string{
						"Some-Other-Header": "val",
						"X-Header-Set":      "set-overwrites-values",
					},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
			{
				Request: http.Request{
					Host: "gateway-integration-test-infra.test.gw.wso2.com",
					Path: "/response-header-modify/add",
				},
				BackendSetResponseHeaders: map[string]string{
					"Some-Other-Header": "val",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path: "/add",
					},
				},
				Response: http.Response{
					Headers: map[string]string{
						"Some-Other-Header": "val",
						"X-Header-Add":      "add-appends-values",
					},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
			{
				Request: http.Request{
					Host: "gateway-integration-test-infra.test.gw.wso2.com",
					Path: "/response-header-modify/add",
				},
				BackendSetResponseHeaders: map[string]string{
					"Some-Other-Header": "val",
					"X-Header-Add":      "some-other-value",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path: "/add",
					},
				},
				Response: http.Response{
					Headers: map[string]string{
						"Some-Other-Header": "val",
						"X-Header-Add":      "add-appends-values",
					},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
			{
				Request: http.Request{
					Host: "gateway-integration-test-infra.test.gw.wso2.com",
					Path: "/response-header-modify/remove",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path: "/remove",
					},
				},
				BackendSetResponseHeaders: map[string]string{
					"X-Header-Remove": "val",
				},
				Response: http.Response{
					AbsentHeaders: []string{"X-Header-Remove"},
				},
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
	},
}
