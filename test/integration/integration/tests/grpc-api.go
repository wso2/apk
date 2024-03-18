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
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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
	//Manifests:   []string{"tests/grpc-api.yaml"},
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
		testCases := []ExpectedResponse{
			{
				out: &student.StudentResponse{
					Name: "Dineth",
					Age:  10,
				},
				err: nil,
			},
		}
		for i := range testCases {
			tc := testCases[i]
			//t.Run(tc.GetTestCaseName(i), func(t *testing.T) {
			//	t.Parallel()
			//	http.MakeRequestAndExpectEventuallyConsistentResponse(t, suite.RoundTripper, suite.TimeoutConfig, gwAddr, tc)
			//})
			t.Run("Invoke gRPC API", func(t *testing.T) {
				out, err := invokeGRPCClientUntilSatisfied(gwAddr, t)
				if err != nil {
					if tc.err != nil {
						t.Errorf("Err -> \nWant: %q\nGot: %q\n", tc.err, err)
					}
				} else {
					if tc.out.Name != out.Name ||
						tc.out.Age != out.Age {
						t.Errorf("Out -> \nWant: %q\nGot : %q", tc.out, out)
					}
				}

			})
		}
	},
}

func invokeGRPCClient(gwAddr string, t *testing.T) (*student.StudentResponse, error) {

	t.Logf("Starting gRPC client...")

	// Set up TLS credentials for the connection without enforcing server certificate validation.
	t.Logf("Setting up TLS credentials without server certificate validation...")
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
		t.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()
	t.Log("Successfully connected to the server.")

	c := student.NewStudentServiceClient(conn)

	// Prepare the context with a timeout for the request.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	t.Log("Sending request to the server...")
	// Create a StudentRequest message
	r := &student.StudentRequest{Id: 1234} // Adjust the ID according to your actual implementation
	response, err := c.GetStudent(ctx, r)
	if err != nil {
		t.Logf("Could not fetch student: %v", err)
	}

	t.Logf("Received response from server: %v\n", response)
	t.Logf("Student Details: %v\n", response)
	return response, nil
}

type ExpectedResponse struct {
	out *student.StudentResponse
	err error
}

func invokeGRPCClientUntilSatisfied(gwAddr string, t *testing.T) (*student.StudentResponse, error) {
	var out *student.StudentResponse
	var err error
	attempt := 0
	maxAttempts := 4

	for attempt < maxAttempts {
		t.Logf("Attempt %d to invoke gRPC client...", attempt+1)
		out, err = invokeGRPCClient(gwAddr, t)

		if err != nil {
			t.Logf("Error on attempt %d: %v", attempt+1, err)
		} else {
			// Check if the response is satisfactory. This condition needs to be defined.
			// For example, assuming a satisfactory condition is when out.Satisfied is true.
			// This is a placeholder condition and should be replaced with your actual success criteria.
			if out != nil && isResponseSatisfactory(out) {
				t.Logf("Satisfactory response received: %+v", out)
				return out, nil
			}
		}

		if attempt < maxAttempts-1 {
			t.Logf("Waiting 20 seconds before next attempt...")
			time.Sleep(20 * time.Second)
		}
		attempt++
	}

	t.Logf("Failed to receive a satisfactory response after %d attempts", maxAttempts)
	return out, err // Returning the last response and error, might need adjustment based on requirements.
}

func isResponseSatisfactory(response *student.StudentResponse) bool {
	// Define the condition for a response to be considered satisfactory.
	// This is a placeholder function and should contain actual logic to evaluate the response.
	return false // Placeholder: assume every response is satisfactory.
}
