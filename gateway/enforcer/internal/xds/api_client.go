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

type APIXDSClient struct {
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

func NewAPIXDSClient(host string, port string, maxRetries int, retryInterval time.Duration, tlsConfig *tls.Config, cfg *config.Server) *APIXDSClient {
	// Create a new APIClient object
	return &APIXDSClient{
		Host:          host,
		Port:          port,
		maxRetries:    maxRetries,
		retryInterval: retryInterval,
		tlsConfig:     tlsConfig,
		grpcConn:      nil,
		log:           cfg.Logger,
	}
}

func (c *APIXDSClient) InitiateAPIXDSConnection() {
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
				c.log.Error(err, "Failed to receive API stream data")
				cancel()
				c.grpcConn.Close()
				go c.InitiateAPIXDSConnection()
				break
			}
			c.log.Info(fmt.Sprintf("Received config: %v", resp))
		}
	}()
}
