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
	"context"
	"crypto/tls"
	"github.com/wso2/apk/test/integration/integration/utils/generatedcode/student"
	"github.com/wso2/apk/test/integration/integration/utils/grpcutils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"sigs.k8s.io/gateway-api/conformance/utils/config"
	"testing"
	"time"

	"github.com/wso2/apk/test/integration/integration/utils/suite"
)

func init() {
	IntegrationTests = append(IntegrationTests, GRPCAPI)
}

// GRPCAPI tests gRPC API
var GRPCAPI = suite.IntegrationTest{
	ShortName:   "GRPCAPI",
	Description: "Tests gRPC API",
	Manifests:   []string{"tests/grpc-api.yaml"},
	Test: func(t *testing.T, suite *suite.IntegrationTestSuite) {
		gwAddr := "grpc.test.gw.wso2.com:9095"
		//token := http.GetTestToken(t)

		//testCases := []http.ExpectedResponse{
		//{
		//	Request: http.Request{
		//		Host:   "gql.test.gw.wso2.com",
		//		Path:   "/gql/v1",
		//		Method: "POST",
		//		Headers: map[string]string{
		//			"Content-Type": "application/json",
		//		},
		//		Body: `{"query":"query{\n    human(id:1000){\n        id\n        name\n    }\n}","variables":{}}`,
		//	},
		//	ExpectedRequest: &http.ExpectedRequest{
		//		Request: http.Request{
		//			Method: ""},
		//	},
		//	Response: http.Response{StatusCode: 200},
		//},
		//}

		testCases := []grpcutils.GRPCTestCase{
			{
				ExpectedResponse: grpcutils.ExpectedResponse{
					Out: &student.StudentResponse{
						Name: "Dineth",
						Age:  10,
					},
					Err: nil,
				},
				ActualResponse: &student.StudentResponse{},
				Name:           "Get Student Details",
			},
		}
		for i := range testCases {
			tc := testCases[i]
			//t.Run(tc.GetTestCaseName(i), func(t *testing.T) {
			//	t.Parallel()
			//	http.MakeRequestAndExpectEventuallyConsistentResponse(t, suite.RoundTripper, suite.TimeoutConfig, gwAddr, tc)
			//})
			t.Run("Invoke gRPC API", func(t *testing.T) {
				t.Parallel()
				invokeGRPCClientUntilSatisfied(gwAddr, t, tc, suite.TimeoutConfig)

			})
		}
	},
}

func invokeGRPCClient(gwAddr string, t *testing.T) (*student.StudentResponse, error) {

	t.Logf("Starting gRPC client...")

	config := &tls.Config{
		InsecureSkipVerify: true, // CAUTION: This disables SSL certificate verification.
	}
	creds := credentials.NewTLS(config)

	// Dial the server with the TLS credentials and a dial timeout.
	dialCtx, dialCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer dialCancel()
	t.Logf("Dialing to server at %s with timeout...", gwAddr)
	conn, err := grpc.DialContext(dialCtx, gwAddr, grpc.WithTransportCredentials(creds), grpc.WithBlock())
	if err != nil {
		t.Logf("Failed to connect: %v", err)
	}
	defer conn.Close()
	//t.Log("Successfully connected to the server.")

	c := student.NewStudentServiceClient(conn)

	// Prepare the context with a timeout for the request.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	t.Log("Sending request to the server...")
	// Create a StudentRequest message
	r := &student.StudentRequest{Id: 1234}
	response, err := c.GetStudent(ctx, r)
	if err != nil {
		t.Logf("Could not fetch student: %v", err)
	}
	t.Logf("Received response from server: %v\n", response)
	return response, nil
}

func invokeGRPCClientUntilSatisfied(gwAddr string, t *testing.T, testCase grpcutils.GRPCTestCase, timeout config.TimeoutConfig) {
	var out *student.StudentResponse
	var err error
	attempt := 0
	maxAttempts := 4
	expected := testCase.ExpectedResponse
	//timeoutDuration := timeout.RequestTimeout * time.Second
	timeoutDuration := 10 * time.Second
	for attempt < maxAttempts {
		t.Logf("Attempt %d to invoke gRPC client...", attempt+1)
		out, err = invokeGRPCClient(gwAddr, t)

		if err != nil {
			t.Logf("Error on attempt %d: %v", attempt+1, err)
		} else {
			if out != nil && isResponseSatisfactory(out, expected) {
				t.Logf("Satisfactory response received: %+v", out)
				return
			}
		}

		if attempt < maxAttempts-1 {
			t.Logf("Waiting %s seconds before next attempt...", timeoutDuration)
			time.Sleep(timeoutDuration)
		}
		attempt++
	}

	t.Logf("Failed to receive a satisfactory response after %d attempts", maxAttempts)
	t.Fail()
}

func isResponseSatisfactory(response *student.StudentResponse, expectedResponse grpcutils.ExpectedResponse) bool {
	if response.Name == expectedResponse.Out.Name && response.Age == expectedResponse.Out.Age {
		return true
	}
	return false
}
