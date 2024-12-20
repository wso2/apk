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


// handleNotificationEvent translates the Java method to Go
func handleNotificationEvent(event *Event) {
	switch event.Type {
	case "ALL_EVENTS":
		log.Println("Received all events from the server")
		SubscriptionDataStoreUtil.loadStartupArtifacts()
	case "APPLICATION_CREATED":
		log.Println("********")
		SubscriptionDataStoreUtil.addApplication(event.Application)
	case "SUBSCRIPTION_CREATED", "SUBSCRIPTION_UPDATED":
		SubscriptionDataStoreUtil.addSubscription(event.Subscription)
	case "APPLICATION_MAPPING_CREATED", "APPLICATION_MAPPING_UPDATED":
		SubscriptionDataStoreUtil.addApplicationMapping(event.ApplicationMapping)
	case "APPLICATION_KEY_MAPPING_CREATED", "APPLICATION_KEY_MAPPING_UPDATED":
		SubscriptionDataStoreUtil.addApplicationKeyMapping(event.ApplicationKeyMapping)
	case "APPLICATION_UPDATED":
		SubscriptionDataStoreUtil.addApplication(event.Application)
	case "APPLICATION_MAPPING_DELETED":
		SubscriptionDataStoreUtil.removeApplicationMapping(event.ApplicationMapping)
	case "APPLICATION_KEY_MAPPING_DELETED":
		SubscriptionDataStoreUtil.removeApplicationKeyMapping(event.ApplicationKeyMapping)
	case "SUBSCRIPTION_DELETED":
		SubscriptionDataStoreUtil.removeSubscription(event.Subscription)
	case "APPLICATION_DELETED":
		SubscriptionDataStoreUtil.removeApplication(event.Application)
	default:
		log.Println("Unknown event type received from the server")
	}
}