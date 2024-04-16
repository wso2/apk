package student

import (
	"context"
	"github.com/wso2/apk/test/integration/integration/utils/grpcutils"
	"google.golang.org/grpc"
	"log"
	"time"
)

type StudentResponseSatisfier struct{}

// IsSatisfactory checks if the given response is satisfactory based on the expected response.
func (srs StudentResponseSatisfier) IsSatisfactory(response any, expectedResponse grpcutils.ExpectedResponse) bool {
	// Type assert the response to *student.StudentResponse
	resp, respOk := response.(*StudentResponse)
	if !respOk {
		log.Println("Failed to assert response as *student.StudentResponse")
		return false
	}
	// Type assert the expected output to *student.StudentResponse
	expectedResp, expOk := expectedResponse.Out.(*StudentResponse)
	if !expOk {
		log.Println("Failed to assert expectedResponse.Out as *student.StudentResponse")
		return false
	}

	// Compare the actual response with the expected response
	if resp.Name == expectedResp.Name && resp.Age == expectedResp.Age {
		log.Println("Response is satisfactory")
		return true
	} else {
		log.Println("Response does not match the expected output")
		return false
	}
}
func GetStudent(conn *grpc.ClientConn) (any, error) {
	c := NewStudentServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	r := &StudentRequest{Id: 1234}
	response, err := c.GetStudent(ctx, r)
	if err != nil {
		return nil, err
	}
	return response, nil
}
