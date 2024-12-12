package util

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func CreateGRPCConnection(ctx context.Context, host, port string, tlsConfig *tls.Config) (*grpc.ClientConn, error) {
	address := fmt.Sprintf("%s:%s", host, port)
	
	kacp := keepalive.ClientParameters{
		Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
		PermitWithoutStream: true,             // send pings even without active streams
	}
	return conn, nil
}
