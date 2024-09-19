package grpcutils

import (
	"crypto/tls"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Request struct {
	Host    string
	Headers map[string]string
}
type ClientCreator[T any] func(conn *grpc.ClientConn) T
type ExpectedResponse struct {
	Out any
	Err error
}

type GRPCTestCase struct {
	Request          Request
	ExpectedResponse ExpectedResponse
	ActualResponse   any
	Name             string
	Method           func(conn *grpc.ClientConn) (any, error)
	Satisfier        ResponseSatisfier
}
type ResponseSatisfier interface {
	IsSatisfactory(response interface{}, expectedResponse ExpectedResponse) bool
}

func DialGRPCServer(gwAddr string, t *testing.T) (*grpc.ClientConn, error) {
	// Set up a connection to the server.
	t.Logf("Dialing gRPC server at %s...", gwAddr)
	creds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})
	conn, err := grpc.Dial(gwAddr, grpc.WithTransportCredentials(creds))
	if err != nil {
		t.Fatalf("Could not connect to the server: %v", err)
	}
	return conn, nil
}
func InvokeGRPCClientUntilSatisfied(gwAddr string, t *testing.T, testCase GRPCTestCase, satisfier ResponseSatisfier, fn ExecuteClientCall) {
	//(delay to allow CRs to be applied)
	time.Sleep(5 * time.Second)

	var out any
	var err error
	attempt := 0
	maxAttempts := 4
	expected := testCase.ExpectedResponse
	timeoutDuration := 50 * time.Second
	for attempt < maxAttempts {
		t.Logf("Attempt %d to invoke gRPC client...", attempt+1)
		out, err = InvokeGRPCClient(gwAddr, t, fn)

		if err != nil {
			t.Logf("Error on attempt %d: %v", attempt+1, err)
		} else {
			if satisfier.IsSatisfactory(out, expected) {
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

type ExecuteClientCall func(conn *grpc.ClientConn) (any, error)

func InvokeGRPCClient(gwAddr string, t *testing.T, fn ExecuteClientCall) (any, error) {

	conn, err := DialGRPCServer(gwAddr, t)
	if err != nil {
		t.Fatalf("Could not connect to the server: %v", err)
	}

	response, err := fn(conn)
	if err != nil {
		return nil, err
	}
	return response, nil
}
