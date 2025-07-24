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
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/wso2/apk/test/integration/integration/utils/http"
	"github.com/wso2/apk/test/integration/integration/utils/suite"
	testtypes "github.com/wso2/apk/test/integration/integration/utils/testtypes"
	"gopkg.in/yaml.v2"
)

var apiPolicy testtypes.APIPolicy
var isBackendJWTEnabled bool
var filePathToResource string

func init() {

	// Get the file path for the resource file
	path, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	paths := strings.Split(path, string(os.PathSeparator))
	if paths[len(paths)-1] == "integration" && paths[len(paths)-2] == "integration" {
		filePathToResource = path + "/tests/resources/tests/api-policy-with-jwt-generator.yaml"
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else if paths[len(paths)-1] == "integration" {
		filePathToResource = path + "/integration/tests/resources/tests/api-policy-with-jwt-generator.yaml"
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	//IntegrationTests = append(IntegrationTests, BackendJWTGenerationPolicy)
}

// BackendJWTGenerationPolicy test
var BackendJWTGenerationPolicy = suite.IntegrationTest{
	ShortName:   "BackendJWTGenerationPolicy",
	Description: "Tests API with backend JWT generation policy",
	Manifests:   []string{"tests/api-policy-with-jwt-generator.yaml"},
	Test: func(t *testing.T, suite *suite.IntegrationTestSuite) {
		ns := "gateway-integration-test-infra"
		gwAddr := "api-policy-with-jwt-generator.test.gw.wso2.com:9095"
		token := http.GetTestToken(t)

		yamlFile, err := ioutil.ReadFile(filepath.Clean(filePathToResource))

		if err != nil {
			t.Error(err)
		}
		err = yaml.Unmarshal(yamlFile, &apiPolicy)
		if err != nil {
			t.Error(err)
		}

		if apiPolicy.Spec.Default != nil && apiPolicy.Spec.Default.BackendJWTToken != nil {
			isBackendJWTEnabled = apiPolicy.Spec.Default.BackendJWTToken.Enabled
		}

		if apiPolicy.Spec.Override != nil && apiPolicy.Spec.Override.BackendJWTToken != nil {
			isBackendJWTEnabled = apiPolicy.Spec.Override.BackendJWTToken.Enabled
		}

		var headers map[string]string = nil
		if isBackendJWTEnabled {
			headers = map[string]string{"X-JWT-Assertion": ""}
		}

		testCases := []http.ExpectedResponse{
			{
				Request: http.Request{
					Host: "api-policy-with-jwt-generator.test.gw.wso2.com",
					Path: "/api-policy-with-jwt-generator/v1.0.0/v2/echo-full",
					Headers: map[string]string{
						"content-type": "application/json",
					},
					Method: "GET",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path:    "/v2/echo-full",
						Headers: headers,
					},
				},
				Backend:   "infra-backend-v1",
				Namespace: ns,
			},
			{
				Request: http.Request{
					Host: "api-policy-with-jwt-generator.test.gw.wso2.com",
					Path: "/api-policy-with-jwt-generator/v2/echo-full",
					Headers: map[string]string{
						"content-type": "application/json",
					},
					Method: "GET",
				},
				ExpectedRequest: &http.ExpectedRequest{
					Request: http.Request{
						Path:    "/v2/echo-full",
						Headers: headers,
					},
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
