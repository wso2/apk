package xds

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	config_ads "github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/service/config"
	"github.com/wso2/apk/gateway/enforcer/internal/config"
	"github.com/wso2/apk/gateway/enforcer/internal/logging"
	"github.com/wso2/apk/gateway/enforcer/internal/util"
	"google.golang.org/grpc"
)

// ConfigXDSClient is a client for managing gRPC connections to the Config Discovery Service (XDS).
// It handles retry logic, TLS configuration, and logging for configuration data streams.
type ConfigXDSClient struct {
	Host          string
	Port          string
	maxRetries    int
	retryInterval time.Duration
	tlsConfig     *tls.Config
	grpcConn      *grpc.ClientConn
	ctx           context.Context
	cancel        context.CancelFunc
	client        config_ads.ConfigDiscoveryServiceClient
	log           logging.Logger
}

// NewXDSConfigClient creates a new instance of ConfigXDSClient.
// It initializes the client with the given host, port, retry parameters, TLS configuration, and logger.
func NewXDSConfigClient(host string, port string, maxRetries int, retryInterval time.Duration, tlsConfig *tls.Config, cfg *config.Server) *ConfigXDSClient {
	// Create a new APIClient object
	return &ConfigXDSClient{
		Host:          host,
		Port:          port,
		maxRetries:    maxRetries,
		retryInterval: retryInterval,
		tlsConfig:     tlsConfig,
		grpcConn:      nil,
		log:           cfg.Logger,
	}
}

// InitiateConfigXDSConnection establishes and maintains a gRPC connection to the Config Discovery Service.
// It also handles reconnection logic on errors and listens for incoming configuration streams.
func (c *ConfigXDSClient) InitiateConfigXDSConnection() {
	grpcConn := util.CreateGRPCConnectionWithRetryAndPanic(nil, c.Host, c.Port, c.tlsConfig, c.maxRetries, c.retryInterval)
	c.grpcConn = grpcConn
	client := config_ads.NewConfigDiscoveryServiceClient(grpcConn)

	ctx, cancel := context.WithCancel(context.Background())
	c.ctx = ctx
	c.cancel = cancel

	stream, err := client.StreamConfigs(ctx)
	if err != nil {
		cancel()
		c.grpcConn.Close()
		panic(fmt.Errorf("Failed to initiate XDS connection with API Discovery Service: %v", err))
	}

	go func() {
		for {
			resp, err := stream.Recv()
			if err != nil {
				c.log.Error(err, "Failed to receive config data")
				cancel()
				c.grpcConn.Close()
				go c.InitiateConfigXDSConnection()
				break
			}
			c.log.Info(fmt.Sprintf("Received config: %v", resp))
		}
	}()
}
