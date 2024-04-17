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
	"github.com/wso2/apk/test/integration/integration/utils/grpc-code/student"
	"github.com/wso2/apk/test/integration/integration/utils/grpcutils"
	"github.com/wso2/apk/test/integration/integration/utils/suite"
	"testing"
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
				Method:         student.GetStudent,
			},
			//{
			//	ExpectedResponse: grpcutils.ExpectedResponse{
			//		Out: &student_default_version.StudentResponse{
			//			Name: "Dineth",
			//			Age:  10,
			//		},
			//		Err: nil,
			//	},
			//	ActualResponse: &student_default_version.StudentResponse{},
			//	Name:           "Get Student Details (Default API Version)",
			//	Method:         student_default_version.GetStudent,
			//},
		}
		for i := range testCases {
			tc := testCases[i]
			t.Run("Invoke gRPC API", func(t *testing.T) {
				t.Parallel()
				grpcutils.InvokeGRPCClientUntilSatisfied(gwAddr, t, tc, student.StudentResponseSatisfier{}, tc.Method)
			})
		}
	},
}
