package grpcutils

import "github.com/wso2/apk/test/integration/integration/utils/generatedcode/student"

type ExpectedResponse struct {
	Out *student.StudentResponse
	Err error
}
type Request struct {
	Host    string
	Headers map[string]string
}
type GRPCTestCase struct {
	Request          Request
	ExpectedResponse ExpectedResponse
	ActualResponse   *student.StudentResponse
	Name             string
}
