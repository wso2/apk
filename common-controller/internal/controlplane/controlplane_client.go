package controlplane

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/wso2/apk/common-controller/internal/config"
	"github.com/wso2/apk/common-controller/internal/loggers"
	"github.com/wso2/apk/common-controller/internal/server"
	"github.com/wso2/apk/common-controller/internal/utils"
	"github.com/wso2/apk/common-go-libs/constants"
	apkmgt "github.com/wso2/apk/common-go-libs/pkg/discovery/api/wso2/discovery/service/apkmgt"
	"github.com/wso2/apk/common-go-libs/pkg/discovery/api/wso2/discovery/subscription"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	grpcStatus "google.golang.org/grpc/status"
)

// Agent is a struct that implements Agent interface
type Agent struct {
	hostname         string
	port             int
	controlPlaneID   string
	artifactDeployer ArtifactDeployer
}

var (
	connectionFaultChannel chan bool
	eventStreamingClient   apkmgt.EventStreamService_StreamEventsClient
)

func init() {
	connectionFaultChannel = make(chan bool, 10)
}

// NewControlPlaneAgent creates a new ControlPlaneAgent
func NewControlPlaneAgent(hostname string, port int, controlPlaneID string, artifactDeployer ArtifactDeployer) *Agent {
	return &Agent{hostname: hostname, port: port, controlPlaneID: controlPlaneID, artifactDeployer: artifactDeployer}
}

func (controlPlaneGrpcClient *Agent) initGrpcConnection() (*grpc.ClientConn, error) {
	config := config.ReadConfigs()
	publicKeyLocation, privateKeyLocation, truststoreLocation := utils.GetKeyLocations()
	cert, err := utils.GetServerCertificate(publicKeyLocation, privateKeyLocation)
	if err != nil {
		return nil, err
	}
	caCertPool := utils.GetTrustedCertPool(truststoreLocation)
	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
		ServerName:   config.CommonController.Server.ServerName,
	})
	hostname := fmt.Sprintf("%s:%d", config.CommonController.ControlPlane.Host, config.CommonController.ControlPlane.EventPort)
	backOff := grpc_retry.BackoffLinearWithJitter(config.CommonController.ControlPlane.RetryInterval*time.Second, 0.5)
	conection, err := grpc.Dial(hostname, grpc.WithTransportCredentials(creds), grpc.WithStreamInterceptor(
		grpc_retry.StreamClientInterceptor(grpc_retry.WithBackoff(backOff))))
	if err != nil {
		loggers.LoggerAPKOperator.Errorf("Error while connecting to the control plane %s", err.Error())
		return nil, err
	}
	md := metadata.Pairs("common-controller-uuid", controlPlaneGrpcClient.controlPlaneID)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	client := apkmgt.NewEventStreamServiceClient(conection)
	eventStreamingClient, err = client.StreamEvents(ctx, &apkmgt.Request{Event: controlPlaneGrpcClient.controlPlaneID})
	if err != nil {
		loggers.LoggerAPKOperator.Errorf("Error while initializing streaming %s", err.Error())
		return nil, err
	}
	return conection, nil
}

// StartEventStreaming starts event streaming
func (controlPlaneGrpcClient *Agent) StartEventStreaming() {
	conn := controlPlaneGrpcClient.initializeGrpcStreaming()
	for retryTrueReceived := range connectionFaultChannel {
		// event is always true
		if !retryTrueReceived {
			continue
		}
		time.Sleep(config.ReadConfigs().CommonController.ControlPlane.RetryInterval * time.Second)
		if conn != nil {
			conn.Close()
		}
		loggers.LoggerAPKOperator.Error("Connection lost. Retrying to connect to the control plane")
		conn = controlPlaneGrpcClient.initializeGrpcStreaming()
	}
}

// initializeGrpcStreaming starts event streaming
func (controlPlaneGrpcClient *Agent) initializeGrpcStreaming() *grpc.ClientConn {
	conn, err := controlPlaneGrpcClient.initGrpcConnection()
	if err != nil {
		loggers.LoggerAPKOperator.Errorf("Error while initializing the connection. error: %s", err.Error())
		connectionFaultChannel <- true
		return conn
	}
	go func() {
		for {
			resp, err := eventStreamingClient.Recv()
			if err == io.EOF {
				connectionFaultChannel <- true
				return
			}
			if err != nil {
				errStatus, _ := grpcStatus.FromError(err)
				if errStatus.Code() == codes.Unavailable ||
					errStatus.Code() == codes.DeadlineExceeded ||
					errStatus.Code() == codes.Canceled ||
					errStatus.Code() == codes.ResourceExhausted ||
					errStatus.Code() == codes.Aborted ||
					errStatus.Code() == codes.Internal {
					loggers.LoggerAPKOperator.Errorf("Connection unavailable. errorCode: %s errorMessage: %s",
						errStatus.Code().String(), errStatus.Message())
					connectionFaultChannel <- true
				}
				return
			}
			controlPlaneGrpcClient.handleEvents(resp)
		}
	}()
	return conn
}
func (controlPlaneGrpcClient *Agent) handleEvents(event *subscription.Event) {
	loggers.LoggerAPKOperator.Infof("Received event %s", event.Type)
	if event.Type == constants.AllEvnts {
		retrieveAllData()
	} else if event.Type == constants.ApplicationCreated {
		loggers.LoggerAPKOperator.Infof("Received APPLICATION_CREATED event.")
		if event.Application != nil {
			application := server.Application{UUID: event.Application.Uuid,
				Name:           event.Application.Name,
				Owner:          event.Application.Owner,
				OrganizationID: event.Application.Organization,
				Attributes:     event.Application.Attributes,
			}
			loggers.LoggerAPKOperator.Infof("Received Application %s", application.UUID)
			controlPlaneGrpcClient.artifactDeployer.DeployApplication(application)
		}
	} else if event.Type == constants.ApplicationUpdated {
		if event.Application != nil {
			application := server.Application{UUID: event.Application.Uuid,
				Name:           event.Application.Name,
				Owner:          event.Application.Owner,
				OrganizationID: event.Application.Organization,
				Attributes:     event.Application.Attributes,
			}
			loggers.LoggerAPKOperator.Infof("Received Application %s", application.UUID)
			controlPlaneGrpcClient.artifactDeployer.UpdateApplication(application)
		}
	} else if event.Type == constants.ApplicationDeleted {
		if event.Application != nil {
			application := server.Application{UUID: event.Application.Uuid,
				Name:           event.Application.Name,
				Owner:          event.Application.Owner,
				OrganizationID: event.Application.Organization,
				Attributes:     event.Application.Attributes,
			}
			loggers.LoggerAPKOperator.Infof("Received Application %s", application.UUID)
			controlPlaneGrpcClient.artifactDeployer.DeleteApplication(application.UUID)
		}
	} else if event.Type == constants.SubscriptionCreated {
		loggers.LoggerAPKOperator.Infof("Received SUBSCRIPTION_CREATED event.")
		if event.Subscription != nil {
			subscription := server.Subscription{UUID: event.Subscription.Uuid,
				Organization:  event.Subscription.Organization,
				SubStatus:     event.Subscription.SubStatus,
				SubscribedAPI: &server.SubscribedAPI{Name: event.Subscription.SubscribedApi.Name, Version: event.Subscription.SubscribedApi.Version},
			}
			loggers.LoggerAPKOperator.Infof("Received Subscription %s", subscription.UUID)
			controlPlaneGrpcClient.artifactDeployer.DeploySubscription(subscription)
		}
	} else if event.Type == constants.SubscriptionUpdated {
		loggers.LoggerAPKOperator.Infof("Received SUBSCRIPTION_UPDATED event.")
		if event.Subscription != nil {
			subscription := server.Subscription{UUID: event.Subscription.Uuid,
				Organization:  event.Subscription.Organization,
				SubStatus:     event.Subscription.SubStatus,
				SubscribedAPI: &server.SubscribedAPI{Name: event.Subscription.SubscribedApi.Name, Version: event.Subscription.SubscribedApi.Version},
			}
			loggers.LoggerAPKOperator.Infof("Received Subscription %s", subscription.UUID)
			controlPlaneGrpcClient.artifactDeployer.UpdateSubscription(subscription)
		}
	} else if event.Type == constants.SubscriptionDeleted {
		loggers.LoggerAPKOperator.Infof("Received SUBSCRIPTION_DELETED event.")
		if event.Subscription != nil {
			subscription := server.Subscription{UUID: event.Subscription.Uuid,
				Organization:  event.Subscription.Organization,
				SubStatus:     event.Subscription.SubStatus,
				SubscribedAPI: &server.SubscribedAPI{Name: event.Subscription.SubscribedApi.Name, Version: event.Subscription.SubscribedApi.Version},
			}
			loggers.LoggerAPKOperator.Infof("Received Subscription %s", subscription.UUID)
			controlPlaneGrpcClient.artifactDeployer.DeleteSubscription(subscription.UUID)
		}
	} else if event.Type == constants.ApplicationKeyMappingCreated {
		loggers.LoggerAPKOperator.Infof("Received APPLICATION_KEY_MAPPING_CREATED event.")
		if event.ApplicationKeyMapping != nil {
			applicationKeyMapping := server.ApplicationKeyMapping{ApplicationUUID: event.ApplicationKeyMapping.ApplicationUUID,
				SecurityScheme:        event.ApplicationKeyMapping.SecurityScheme,
				ApplicationIdentifier: event.ApplicationKeyMapping.ApplicationIdentifier,
				KeyType:               event.ApplicationKeyMapping.KeyType,
				EnvID:                 event.ApplicationKeyMapping.EnvID,
				OrganizationID:        event.ApplicationKeyMapping.Organization,
			}
			loggers.LoggerAPKOperator.Infof("Received ApplicationKeyMapping %s", applicationKeyMapping.ApplicationUUID)
			controlPlaneGrpcClient.artifactDeployer.DeployKeyMappings(applicationKeyMapping)
		}
	} else if event.Type == constants.ApplicationKeyMappingDeleted {
		loggers.LoggerAPKOperator.Infof("Received APPLICATION_KEY_MAPPING_DELETED event.")
		if event.ApplicationKeyMapping != nil {
			applicationKeyMapping := server.ApplicationKeyMapping{ApplicationUUID: event.ApplicationKeyMapping.ApplicationUUID,
				SecurityScheme:        event.ApplicationKeyMapping.SecurityScheme,
				ApplicationIdentifier: event.ApplicationKeyMapping.ApplicationIdentifier,
				KeyType:               event.ApplicationKeyMapping.KeyType,
				EnvID:                 event.ApplicationKeyMapping.EnvID,
				OrganizationID:        event.ApplicationKeyMapping.Organization,
			}
			loggers.LoggerAPKOperator.Infof("Received ApplicationKeyMapping %s", applicationKeyMapping.ApplicationUUID)
			controlPlaneGrpcClient.artifactDeployer.DeleteKeyMappings(applicationKeyMapping)
		}
	} else if event.Type == constants.ApplicationMappingCreated {
		loggers.LoggerAPKOperator.Infof("Received APPLICATION_MAPPING_CREATED event.")
		if event.ApplicationMapping != nil {
			applicationMapping := server.ApplicationMapping{UUID: event.ApplicationMapping.Uuid,
				ApplicationRef:  event.ApplicationMapping.ApplicationRef,
				SubscriptionRef: event.ApplicationMapping.SubscriptionRef,
				OrganizationID:  event.ApplicationMapping.Organization,
			}
			loggers.LoggerAPKOperator.Infof("Received ApplicationMapping %s", applicationMapping.UUID)
			controlPlaneGrpcClient.artifactDeployer.DeployApplicationMappings(applicationMapping)
		}
	} else if event.Type == constants.ApplicationMappingDeleted {
		loggers.LoggerAPKOperator.Infof("Received APPLICATION_MAPPING_DELETED event.")
		if event.ApplicationMapping != nil {
			applicationMapping := server.ApplicationMapping{UUID: event.ApplicationMapping.Uuid,
				ApplicationRef:  event.ApplicationMapping.ApplicationRef,
				SubscriptionRef: event.ApplicationMapping.SubscriptionRef,
				OrganizationID:  event.ApplicationMapping.Organization,
			}
			loggers.LoggerAPKOperator.Infof("Received ApplicationMapping %s", applicationMapping.UUID)
			controlPlaneGrpcClient.artifactDeployer.DeleteApplicationMappings(applicationMapping.UUID)
		}
	} else if event.Type == constants.ApplicationMappingUpdated {
		loggers.LoggerAPKOperator.Infof("Received APPLICATION_MAPPING_UPDATED event.")
		if event.ApplicationMapping != nil {
			applicationMapping := server.ApplicationMapping{UUID: event.ApplicationMapping.Uuid,
				ApplicationRef:  event.ApplicationMapping.ApplicationRef,
				SubscriptionRef: event.ApplicationMapping.SubscriptionRef,
				OrganizationID:  event.ApplicationMapping.Organization,
			}
			loggers.LoggerAPKOperator.Infof("Received ApplicationMapping %s", applicationMapping.UUID)
			controlPlaneGrpcClient.artifactDeployer.UpdateApplicationMappings(applicationMapping)
		}

	}
}
func retrieveAllData() {

}
