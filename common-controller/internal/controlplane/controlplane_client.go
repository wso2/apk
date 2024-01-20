package controlplane

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"

	"github.com/wso2/apk/common-controller/internal/config"
	"github.com/wso2/apk/common-controller/internal/loggers"
	"github.com/wso2/apk/common-controller/internal/server"
	"github.com/wso2/apk/common-controller/internal/utils"
	"github.com/wso2/apk/common-go-libs/constants"
	apkmgt "github.com/wso2/apk/common-go-libs/pkg/discovery/api/wso2/discovery/service/apkmgt"
	"github.com/wso2/apk/common-go-libs/pkg/discovery/api/wso2/discovery/subscription"
)

// Agent is a struct that implements Agent interface
type Agent struct {
	hostname         string
	port             int
	controlPlaneID   string
	artifactDeployer ArtifactDeployer
}

// NewControlPlaneAgent creates a new ControlPlaneAgent
func NewControlPlaneAgent(hostname string, port int, controlPlaneID string, artifactDeployer ArtifactDeployer) *Agent {
	return &Agent{hostname: hostname, port: port, controlPlaneID: controlPlaneID, artifactDeployer: artifactDeployer}
}

// StartEventStreaming starts event streaming
func (controlPlaneGrpcClient *Agent) StartEventStreaming() error {
	config := config.ReadConfigs()
	publicKeyLocation, privateKeyLocation, truststoreLocation := utils.GetKeyLocations()
	cert, err := utils.GetServerCertificate(publicKeyLocation, privateKeyLocation)
	if err != nil {
		return err
	}
	caCertPool := utils.GetTrustedCertPool(truststoreLocation)
	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
		ServerName:   config.CommonController.Server.ServerName,
	})
	hostname := fmt.Sprintf("%s%d", controlPlaneGrpcClient.hostname, controlPlaneGrpcClient.port)
	conn, err := grpc.Dial(hostname, grpc.WithTransportCredentials(creds))
	if err != nil {
		return err
	}
	md := metadata.Pairs("common-controller-uuid", controlPlaneGrpcClient.controlPlaneID)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	client := apkmgt.NewEventStreamServiceClient(conn)
	streamClient, err := client.StreamEvents(ctx, &apkmgt.Request{Event: controlPlaneGrpcClient.controlPlaneID})
	if err != nil {
		return err
	}
	done := make(chan bool)
	go func() {
		for {
			resp, err := streamClient.Recv()
			if err == io.EOF {
				done <- true
				return
			}
			if err != nil {
				log.Fatalf("can not receive %v", err)
			}
			controlPlaneGrpcClient.handleEvents(resp)
		}
	}()

	<-done
	log.Printf("finished")
	return nil
}
func (controlPlaneGrpcClient *Agent) handleEvents(event *subscription.Event) {
	loggers.LoggerAPKOperator.Infof("Received event %s", event.Type)
	if event.Type == constants.AllEvnts {
		loggers.LoggerAPKOperator.Infof("Received event %s", event.Type)
	} else if event.Type == constants.ApplicationCreated {
		if event.Application != nil {
			application := server.Application{UUID: event.Application.Uuid,
				Name:           event.Application.Name,
				Owner:          event.Application.Owner,
				OrganizationID: event.Application.Organization,
				Attributes:     event.Application.Attributes,
			}
			controlPlaneGrpcClient.artifactDeployer.DeployApplication(application)
		}
	}
}
