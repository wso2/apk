package grpcutils

import "github.com/wso2/apk/test/integration/integration/utils/generatedcode/student"

type StudentResponseSatisfier struct{}

// IsSatisfactory checks if the given response is satisfactory based on the expected response.
func (srs StudentResponseSatisfier) IsSatisfactory(response interface{}, expectedResponse ExpectedResponse) bool {
	// Type assert the response to *student.StudentResponse
	resp, ok := response.(*student.StudentResponse)
	if !ok {
		return false // or panic, or handle the error according to your error handling policy
	}

	if resp.Name == expectedResponse.Out.Name && resp.Age == expectedResponse.Out.Age {
		return true
	}
	return false
}
