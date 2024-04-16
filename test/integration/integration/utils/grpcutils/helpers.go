package grpcutils

import "github.com/wso2/apk/test/integration/integration/utils/generatedcode/student"

type Request struct {
	Host    string
	Headers map[string]string
}

type ExpectedResponse struct {
	Out *student.StudentResponse
	Err error
}

type GRPCTestCase struct {
	Request          Request
	ExpectedResponse ExpectedResponse
	ActualResponse   *student.StudentResponse
	Name             string
}
type ResponseSatisfier interface {
	IsSatisfactory(response interface{}, expectedResponse ExpectedResponse) bool
}
