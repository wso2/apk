package grpc

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	api_ads "github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/service/api"
	subscription_service "github.com/wso2/apk/common-go-libs/pkg/discovery/api/wso2/discovery/service/apkmgt"
	"github.com/wso2/apk/gateway/enforcer/internal/config"
	"github.com/wso2/apk/gateway/enforcer/internal/logging"
	"github.com/wso2/apk/gateway/enforcer/internal/util"
	"google.golang.org/grpc"
)

// EventingGRPCClient is a client for managing gRPC connections to an eventing service.
// It includes configuration for retries, TLS, and logging.
type EventingGRPCClient struct {
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

// NewEventingGRPCClient creates a new instance of EventingGRPCClient.
// It initializes the client with the given host, port, retry parameters, TLS configuration, and logger.
func NewEventingGRPCClient(host string, port string, maxRetries int, retryInterval time.Duration, tlsConfig *tls.Config, cfg *config.Server) *EventingGRPCClient {
	// Create a new APIClient object
	return &EventingGRPCClient{
		Host:          host,
		Port:          port,
		maxRetries:    maxRetries,
		retryInterval: retryInterval,
		tlsConfig:     tlsConfig,
		grpcConn:      nil,
		log:           cfg.Logger,
	}
}

// InitiateEventingGRPCConnection establishes and maintains a gRPC connection to the eventing service.
// It also handles reconnection logic on errors and listens for incoming event streams.
func (c *EventingGRPCClient) InitiateEventingGRPCConnection() {
	grpcConn := util.CreateGRPCConnectionWithRetryAndPanic(nil, c.Host, c.Port, c.tlsConfig, c.maxRetries, c.retryInterval)
	c.grpcConn = grpcConn
	client := subscription_service.NewEventStreamServiceClient(grpcConn)

	ctx, cancel := context.WithCancel(context.Background())
	c.ctx = ctx
	c.cancel = cancel

	stream, err := client.StreamEvents(ctx, &subscription_service.Request{Event: "your-event"})
	if err != nil {
		cancel()
		c.grpcConn.Close()
		panic(fmt.Errorf("Failed to initiate GRPC connection with CommonController subscription grpc server: %v", err))
	}

	go func() {
		for {
			resp, err := stream.Recv()
			if err != nil {
				c.log.Error(err, "Failed to receive API stream data")
				cancel()
				c.grpcConn.Close()
				go c.InitiateEventingGRPCConnection()
				break
			}
			c.log.Info(fmt.Sprintf("Received config: %v", resp))
		}
	}()
}
