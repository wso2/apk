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
			t.Run("Invoke gRPC API", func(t *testing.T) {
				t.Parallel()
				invokeGRPCClientUntilSatisfied(gwAddr, t, tc, suite.TimeoutConfig, grpcutils.StudentResponseSatisfier{})

			})
		}
	},
}

func invokeGRPCClient(gwAddr string, t *testing.T, timeout config.TimeoutConfig) (*student.StudentResponse, error) {

	t.Logf("Starting gRPC client...")

	config := &tls.Config{
		InsecureSkipVerify: true,
	}
	creds := credentials.NewTLS(config)

	// Dial the server with the TLS credentials and a dial timeout.
	dialCtx, dialCancel := context.WithTimeout(context.Background(), 60*time.Second)
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
	ctx, cancel := context.WithTimeout(context.Background(), timeout.RequestTimeout)
	defer cancel()

	t.Log("Sending request to the server...")
	// Create a StudentRequest message
	r := &student.StudentRequest{Id: 1234}
	response, err := c.GetStudent(ctx, r)
	if err != nil {
		t.Logf("Could not fetch student: %v", err)
		t.Logf("Error: %v\n", response)
	}

	return response, nil
}

func invokeGRPCClientUntilSatisfied(gwAddr string, t *testing.T, testCase grpcutils.GRPCTestCase, timeout config.TimeoutConfig, satisfier grpcutils.ResponseSatisfier) {
	//(delay to allow CRs to be applied)
	time.Sleep(5 * time.Second)
	var out *student.StudentResponse
	var err error
	attempt := 0
	maxAttempts := 4
	expected := testCase.ExpectedResponse
	timeoutDuration := 50 * time.Second
	for attempt < maxAttempts {
		t.Logf("Attempt %d to invoke gRPC client...", attempt+1)
		out, err = invokeGRPCClient(gwAddr, t, timeout)

		if err != nil {
			t.Logf("Error on attempt %d: %v", attempt+1, err)
		} else {
			if satisfier.IsSatisfactory(out, expected) {
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
