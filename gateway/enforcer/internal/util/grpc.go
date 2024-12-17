package util

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

func CreateGRPCConnection(ctx context.Context, host, port string, tlsConfig *tls.Config) (*grpc.ClientConn, error) {
	address := fmt.Sprintf("%s:%s", host, port)
	
	kacp := keepalive.ClientParameters{
		Time:                300 * time.Second,
		PermitWithoutStream: true,             
	}

	dialOptions := []grpc.DialOption{
		grpc.WithKeepaliveParams(kacp),
	}

	dialOptions = append(dialOptions, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	conn, err := grpc.NewClient(
		address,
		dialOptions...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection: %v", err)
	}
	return conn, nil
}

func CreateGRPCConnectionWithRetry(ctx context.Context, host, port string, tlsConfig *tls.Config, maxRetries int, retryInterval time.Duration) (*grpc.ClientConn, error) {
	for retries := 0; maxRetries == -1 || retries < maxRetries; retries++ {
		conn, err := CreateGRPCConnection(ctx, host, port, tlsConfig)
		if err == nil {
			return conn, nil
		}
		time.Sleep(retryInterval)
	}
	return nil, fmt.Errorf("failed to create gRPC connection after %d retries", maxRetries)
}

func CreateGRPCConnectionWithRetryAndPanic(ctx context.Context, host, port string, tlsConfig *tls.Config, maxRetries int, retryInterval time.Duration) *grpc.ClientConn {
	conn, err := CreateGRPCConnectionWithRetry(ctx, host, port, tlsConfig, maxRetries, retryInterval)
	if err != nil {
		panic(err)
	}
	return conn
}

func CreateGRPCServer(publicKeyPath, privateKeyPath string) (*grpc.Server, error) {
	// Load TLS credentials
	cert, err := LoadCertificates(publicKeyPath, privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load TLS credentials: %v", err)
	}

	creds := credentials.NewTLS(&tls.Config{Certificates: []tls.Certificate{cert}})
	// Create and return a new gRPC server with the loaded credentials
	return grpc.NewServer(grpc.Creds(creds)), nil
}