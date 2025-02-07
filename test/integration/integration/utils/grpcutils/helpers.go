package grpcutils

import (
	"context"
	"crypto/tls"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
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
	Method           func(conn *grpc.ClientConn, ctx context.Context) (any, error)
	Satisfier        ResponseSatisfier
}
type ResponseSatisfier interface {
	IsSatisfactory(response interface{}, expectedResponse ExpectedResponse) bool
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

	wrappedFn := func(conn *grpc.ClientConn, ctx context.Context) (any, error) {
		authHeader, exists := testCase.Request.Headers["Authorization"]
		if exists {
			md := metadata.New(nil) // Create empty metadata
			md.Append("authorization", authHeader)
			ctx = metadata.NewOutgoingContext(ctx, md)

		}
		return fn(conn, ctx)
	}
	for attempt < maxAttempts {
		t.Logf("Attempt %d to invoke gRPC client...", attempt+1)
		out, err = InvokeGRPCClient(gwAddr, t, wrappedFn, testCase.Request.Headers["Authorization"])

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

type ExecuteClientCall func(conn *grpc.ClientConn, ctx context.Context) (any, error)

func InvokeGRPCClient(gwAddr string, t *testing.T, fn ExecuteClientCall, authHeader string) (any, error) {
	conn, err := DialGRPCServer(gwAddr, t, authHeader)
	if err != nil {
		t.Fatalf("Could not connect to the server: %v", err)
	}
	ctx := context.Background()
	response, err := fn(conn, ctx)

	if err != nil {
		return nil, err
	}
	return response, nil
}

type JWTAuth struct {
	token string
}

// GetRequestMetadata adds the Authorization header
func (j *JWTAuth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": j.token,
	}, nil
}

// RequireTransportSecurity indicates if TLS is required (false if using insecure connection)
func (j *JWTAuth) RequireTransportSecurity() bool {
	return false
}

func DialGRPCServer(gwAddr string, t *testing.T, authHeader string) (*grpc.ClientConn, error) {
	t.Logf("Dialing gRPC server at %s...", gwAddr)
	creds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})
	conn, err := grpc.Dial(gwAddr, grpc.WithTransportCredentials(creds), grpc.WithPerRPCCredentials(&JWTAuth{token: authHeader}))
	if err != nil {
		t.Fatalf("Could not connect to the server: %v", err)
	}
	return conn, nil
}
