package grpc

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	api_ads "github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/service/api"
	subscription_service "github.com/wso2/apk/common-go-libs/pkg/discovery/api/wso2/discovery/service/apkmgt"
	subscription_proto_model "github.com/wso2/apk/common-go-libs/pkg/discovery/api/wso2/discovery/subscription"
	rest_server_model "github.com/wso2/apk/common-go-libs/pkg/server/model"
	"github.com/wso2/apk/gateway/enforcer/internal/config"
	data_store "github.com/wso2/apk/gateway/enforcer/internal/datastore"
	"github.com/wso2/apk/gateway/enforcer/internal/logging"
	"github.com/wso2/apk/gateway/enforcer/internal/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// EventingGRPCClient is a client for managing gRPC connections to an eventing service.
// It includes configuration for retries, TLS, and logging.
type EventingGRPCClient struct {
	Host            string
	Port            string
	maxRetries      int
	retryInterval   time.Duration
	tlsConfig       *tls.Config
	grpcConn        *grpc.ClientConn
	ctx             context.Context
	cancel          context.CancelFunc
	client          api_ads.ApiDiscoveryServiceClient
	log             logging.Logger
	subAppDataStore *data_store.SubscriptionApplicationDataStore
}

// NewEventingGRPCClient creates a new instance of EventingGRPCClient.
// It initializes the client with the given host, port, retry parameters, TLS configuration, and logger.
func NewEventingGRPCClient(host string, port string, maxRetries int, retryInterval time.Duration, tlsConfig *tls.Config, cfg *config.Server, dataStore *data_store.SubscriptionApplicationDataStore) *EventingGRPCClient {
	// Create a new APIClient object
	return &EventingGRPCClient{
		Host:            host,
		Port:            port,
		maxRetries:      maxRetries,
		retryInterval:   retryInterval,
		tlsConfig:       tlsConfig,
		grpcConn:        nil,
		log:             cfg.Logger,
		subAppDataStore: dataStore,
	}
}

// InitiateEventingGRPCConnection establishes and maintains a gRPC connection to the eventing service.
// It also handles reconnection logic on errors and listens for incoming event streams.
func (c *EventingGRPCClient) InitiateEventingGRPCConnection() {
	// Generate a unique connection ID
	connectionID := uuid.New().String()

	// Create metadata with the enforcer-uuid
	md := metadata.New(map[string]string{"enforcer-uuid": connectionID})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	// Log the connection ID for debugging
	c.log.Info(fmt.Sprintf("Sending request with metadata: enforcer-uuid=%s", connectionID))

	// Create a gRPC connection
	grpcConn := util.CreateGRPCConnectionWithRetryAndPanic(nil, c.Host, c.Port, c.tlsConfig, c.maxRetries, c.retryInterval)
	c.grpcConn = grpcConn
	client := subscription_service.NewEventStreamServiceClient(grpcConn)

	stream, err := client.StreamEvents(ctx, &subscription_service.Request{Event: "ALL_EVENTS"})
	if err != nil {
		c.cancel()
		c.grpcConn.Close()
		panic(fmt.Errorf("Failed to initiate GRPC connection with CommonController subscription grpc server: %v", err))
	}

	// Handle incoming messages in a separate goroutine
	c.log.Info("Connected to the gRPC stream")
	go func() {
		for {
			resp, err := stream.Recv()
			if err != nil {
				c.log.Error(err, "Failed to receive API stream data")
				c.cancel()
				c.grpcConn.Close()
				go c.InitiateEventingGRPCConnection()
				break
			}
			c.log.Info(fmt.Sprintf("Received config: %v", resp))
		}
	}()
}

// handleNotificationEvent translates the Java method to Go
func (c *EventingGRPCClient) handleNotificationEvent(event *subscription_proto_model.Event) {
	switch event.Type {
	case "ALL_EVENTS":
		log.Println("Received all events from the server")
		c.subAppDataStore.LoadStartupData()
	case "SUBSCRIPTION_CREATED", "SUBSCRIPTION_UPDATED":
		log.Println("Subscription created or updated")
		c.subAppDataStore.AddSubscription(convertProtoSubscriptionToRestSubscription(event.Subscription))
	case "APPLICATION_CREATED", "APPLICATION_UPDATED":
		c.subAppDataStore.AddApplication(convertProtoApplicationToRestApplication(event.Application))
	case "APPLICATION_MAPPING_CREATED", "APPLICATION_MAPPING_UPDATED":
		c.subAppDataStore.AddApplicationMapping(convertProtoApplicationMappingToRestApplicationMapping(event.ApplicationMapping))
	case "APPLICATION_KEY_MAPPING_CREATED", "APPLICATION_KEY_MAPPING_UPDATED":
		c.subAppDataStore.AddApplicationKeyMapping(convertProtoApplicationKeyMappingToRestApplicationKeyMapping(event.ApplicationKeyMapping))
	case "SUBSCRIPTION_DELETED":
		c.subAppDataStore.DeleteSubscription(event.Subscription.Uuid)
	case "APPLICATION_MAPPING_DELETED":
		c.subAppDataStore.DeleteApplicationMapping(event.ApplicationMapping.Uuid)
	case "APPLICATION_KEY_MAPPING_DELETED":
		c.subAppDataStore.DeleteApplicationKeyMapping(event.ApplicationKeyMapping.ApplicationIdentifier)
	case "APPLICATION_DELETED":
		c.subAppDataStore.DeleteApplication(event.Application.Uuid)
	default:
		log.Println("Unknown event type received from the server")
	}
}

func convertProtoApplicationToRestApplication(appSource *subscription_proto_model.Application) *rest_server_model.Application {
	return &rest_server_model.Application{
		UUID:           appSource.Uuid,
		Name:           appSource.Name,
		Owner:          appSource.Owner,
		Attributes:     appSource.Attributes,
		OrganizationID: appSource.Organization,
		TimeStamp:      time.Now().Unix(),
	}
}

func convertProtoSubscriptionToRestSubscription(subSource *subscription_proto_model.Subscription) *rest_server_model.Subscription {
	return &rest_server_model.Subscription{
		UUID:          subSource.Uuid,
		SubStatus:     subSource.SubStatus,
		Organization:  subSource.Organization,
		RatelimitTier: subSource.RatelimitTier,
		SubscribedAPI: &rest_server_model.SubscribedAPI{
			Name:    subSource.SubscribedApi.Name,
			Version: subSource.SubscribedApi.Version,
		},
	}
}

func convertProtoApplicationMappingToRestApplicationMapping(appMapSource *subscription_proto_model.ApplicationMapping) *rest_server_model.ApplicationMapping {
	return &rest_server_model.ApplicationMapping{
		UUID:            appMapSource.Uuid,
		ApplicationRef:  appMapSource.ApplicationRef,
		SubscriptionRef: appMapSource.SubscriptionRef,
		OrganizationID:  appMapSource.Organization,
	}
}

func convertProtoApplicationKeyMappingToRestApplicationKeyMapping(appKeyMapSource *subscription_proto_model.ApplicationKeyMapping) *rest_server_model.ApplicationKeyMapping {
	return &rest_server_model.ApplicationKeyMapping{
		ApplicationUUID:       appKeyMapSource.ApplicationUUID,
		ApplicationIdentifier: appKeyMapSource.ApplicationIdentifier,
		OrganizationID:        appKeyMapSource.Organization,
		SecurityScheme:        appKeyMapSource.SecurityScheme,
		KeyType:               appKeyMapSource.KeyType,
		EnvID:                 appKeyMapSource.EnvID,
	}
}
