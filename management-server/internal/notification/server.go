/*
 *  Copyright (c) 2022, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */

package notification

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"github.com/wso2/apk/adapter/pkg/logging"
	"github.com/wso2/apk/adapter/pkg/utils/tlsutils"
	"github.com/wso2/apk/management-server/internal/config"
	"github.com/wso2/apk/management-server/internal/logger"
	"github.com/wso2/apk/management-server/internal/synchronizer"
	internal_types "github.com/wso2/apk/management-server/internal/types"
	"github.com/wso2/apk/management-server/internal/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

type notificationService struct {
	notificationService UnimplementedNotificationServiceServer
}

func newnotificationService() *notificationService {
	return &notificationService{}
}

// CreateApplication sends an application create event
func (s *notificationService) CreateApplication(ctx context.Context, application *Application) (*NotificationResponse, error) {
	logger.LoggerNotificationServer.Infof("Message received : %q", &application)
	var event = internal_types.ApplicationEvent{
		Label:         config.ReadConfigs().ManagementServer.NodeLabels[0],
		UUID:          application.Uuid,
		IsRemoveEvent: false,
	}
	synchronizer.AddApplicationEventsToChannel(event)
	return &NotificationResponse{Code: NotificationResponse_OK}, nil
}

// UpdateApplication sends an application update event
func (s *notificationService) UpdateApplication(ctx context.Context, application *Application) (*NotificationResponse, error) {
	logger.LoggerNotificationServer.Infof("Message received : %q", &application)
	var event = internal_types.ApplicationEvent{
		Label:         config.ReadConfigs().ManagementServer.NodeLabels[0],
		UUID:          application.Uuid,
		IsRemoveEvent: false,
	}
	synchronizer.AddApplicationEventsToChannel(event)
	return &NotificationResponse{Code: NotificationResponse_OK}, nil
}

// DeleteApplication sends an application delete event
func (s *notificationService) DeleteApplication(ctx context.Context, application *Application) (*NotificationResponse, error) {
	logger.LoggerNotificationServer.Infof("Message received : %q", &application)
	var event = internal_types.ApplicationEvent{
		Label:         config.ReadConfigs().ManagementServer.NodeLabels[0],
		UUID:          application.Uuid,
		IsRemoveEvent: true,
	}
	synchronizer.AddApplicationEventsToChannel(event)
	return &NotificationResponse{Code: NotificationResponse_OK}, nil
}

// CreateSubscription sends a subscription create event
func (s *notificationService) CreateSubscription(ctx context.Context, subscription *Subscription) (*NotificationResponse, error) {
	logger.LoggerNotificationServer.Infof("Message received : %q", &subscription)
	var event = internal_types.SubscriptionEvent{
		Label:         config.ReadConfigs().ManagementServer.NodeLabels[0],
		UUID:          subscription.Uuid,
		IsRemoveEvent: false,
	}
	synchronizer.AddSubscriptionEventsToChannel(event)
	return &NotificationResponse{Code: NotificationResponse_OK}, nil
}

// UpdateSubscription sends a subscription update event
func (s *notificationService) UpdateSubscription(ctx context.Context, subscription *Subscription) (*NotificationResponse, error) {
	logger.LoggerNotificationServer.Infof("Message received : %q", &subscription)
	var event = internal_types.SubscriptionEvent{
		Label:         config.ReadConfigs().ManagementServer.NodeLabels[0],
		UUID:          subscription.Uuid,
		IsRemoveEvent: false,
	}
	synchronizer.AddSubscriptionEventsToChannel(event)
	return &NotificationResponse{Code: NotificationResponse_OK}, nil
}

// DeleteSubscription sends a subscription delete event
func (s *notificationService) DeleteSubscription(ctx context.Context, subscription *Subscription) (*NotificationResponse, error) {
	logger.LoggerNotificationServer.Infof("Message received : %q", &subscription)
	var event = internal_types.SubscriptionEvent{
		Label:         config.ReadConfigs().ManagementServer.NodeLabels[0],
		UUID:          subscription.Uuid,
		IsRemoveEvent: true,
	}
	synchronizer.AddSubscriptionEventsToChannel(event)
	return &NotificationResponse{Code: NotificationResponse_OK}, nil
}

// StartGRPCServer starts the GRPC server for notifications
func StartGRPCServer() {
	var grpcOptions []grpc.ServerOption
	publicKeyLocation, privateKeyLocation, truststoreLocation := utils.GetKeyLocations()
	cert, err := tlsutils.GetServerCertificate(publicKeyLocation, privateKeyLocation)
	caCertPool := tlsutils.GetTrustedCertPool(truststoreLocation)
	if err == nil {
		grpcOptions = append(grpcOptions, grpc.Creds(
			credentials.NewTLS(&tls.Config{
				Certificates: []tls.Certificate{cert},
				ClientAuth:   tls.RequireAndVerifyClientCert,
				ClientCAs:    caCertPool,
			}),
		))
	} else {
		logger.LoggerNotificationServer.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprintf("Failed to initiate the ssl context, error: %v", err.Error()),
			Severity:  logging.BLOCKER,
			ErrorCode: 1200,
		})
	}
	grpcOptions = append(grpcOptions, grpc.KeepaliveParams(
		keepalive.ServerParameters{
			Time:    time.Duration(5 * time.Minute),
			Timeout: time.Duration(20 * time.Second),
		}),
	)
	grpcServer := grpc.NewServer(grpcOptions...)
	conf := config.ReadConfigs()
	port := conf.ManagementServer.NotificationPort
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logger.LoggerNotificationServer.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprintf("Failed to listen for Notifications on port: %v, error: %v", port, err.Error()),
			Severity:  logging.BLOCKER,
			ErrorCode: 1201,
		})
	}
	// register services
	notificationService := newnotificationService()
	RegisterNotificationServiceServer(grpcServer, notificationService)
	logger.LoggerNotificationServer.Infof("Management server is listening for GRPC connections on port: %v.", port)
	grpcServer.Serve(lis)
}
