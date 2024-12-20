package xds

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	api_ads "github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/service/api"
	"github.com/wso2/apk/gateway/enforcer/internal/config"
	"github.com/wso2/apk/gateway/enforcer/internal/logging"
	"github.com/wso2/apk/gateway/enforcer/internal/util"
	"google.golang.org/grpc"
)

// JWTIssuerXDSClient is a client for managing gRPC connections to the API Discovery Service (XDS) 
// for JWT issuer-related configurations. It includes retry logic, TLS configuration, and logging.
type JWTIssuerXDSClient struct {
	Host          string
	Port          string
	maxRetries    int
	retryInterval time.Duration
	tlsConfig     *tls.Config
	grpcConn      *grpc.ClientConn
	ctx           context.Context
	cancel        context.CancelFunc
	client        api_ads.ApiDiscoveryServiceClient
	log           logging.Logger
}

// NewJWTIssuerXDSClient creates a new instance of JWTIssuerXDSClient.
// It initializes the client with the given host, port, retry parameters, TLS configuration, and logger.
func NewJWTIssuerXDSClient(host string, port string, maxRetries int, retryInterval time.Duration, tlsConfig *tls.Config, cfg *config.Server) *JWTIssuerXDSClient {
	// Create a new APIClient object
	return &JWTIssuerXDSClient{
		Host:          host,
		Port:          port,
		maxRetries:    maxRetries,
		retryInterval: retryInterval,
		tlsConfig:     tlsConfig,
		grpcConn:      nil,
		log:           cfg.Logger,
	}
}

// InitiateSubscriptionXDSConnection establishes and maintains a gRPC connection to the API Discovery Service.
// It also handles reconnection logic on errors and listens for incoming JWT issuer configuration streams.
func (c *JWTIssuerXDSClient) InitiateSubscriptionXDSConnection() {
	grpcConn := util.CreateGRPCConnectionWithRetryAndPanic(nil, c.Host, c.Port, c.tlsConfig, c.maxRetries, c.retryInterval)
	c.grpcConn = grpcConn
	client := api_ads.NewApiDiscoveryServiceClient(grpcConn)

	ctx, cancel := context.WithCancel(context.Background())
	c.ctx = ctx
	c.cancel = cancel

	stream, err := client.StreamApis(ctx)
	if err != nil {
		cancel()
		c.grpcConn.Close()
		panic(fmt.Errorf("Failed to initiate XDS connection with API Discovery Service: %v", err))
	}

	go func() {
		for {
			resp, err := stream.Recv()
			if err != nil {
				c.log.Error(err, "Failed to receive jwt issuer")
				cancel()
				c.grpcConn.Close()
				go c.InitiateSubscriptionXDSConnection()
				break
			}
			c.log.Info(fmt.Sprintf("Received config: %v", resp))
		}
	}()
}
