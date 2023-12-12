package org.wso2.apk.enforcer.subscription;
/*
 * Copyright (c) 2023, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import io.grpc.ClientInterceptor;
import io.grpc.ConnectivityState;
import io.grpc.ManagedChannel;
import io.grpc.Metadata;
import io.grpc.netty.shaded.io.grpc.netty.NettyChannelBuilder;
import io.grpc.stub.MetadataUtils;
import io.grpc.stub.StreamObserver;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.wso2.apk.enforcer.config.ConfigHolder;
import org.wso2.apk.enforcer.discovery.scheduler.XdsSchedulerManager;
import org.wso2.apk.enforcer.discovery.service.apkmgt.EventStreamServiceGrpc;
import org.wso2.apk.enforcer.discovery.service.apkmgt.Request;
import org.wso2.apk.enforcer.discovery.subscription.Application;
import org.wso2.apk.enforcer.discovery.subscription.Event;
import org.wso2.apk.enforcer.util.GRPCUtils;

import java.util.UUID;
import java.util.concurrent.TimeUnit;

/**
 * Client to communicate with JWTIssuer discovery service at the adapter.
 */
public class EventingGrpcClient implements Runnable {

    private static final Logger logger = LogManager.getLogger(EventingGrpcClient.class);
    private static EventingGrpcClient instance;
    private ManagedChannel channel;
    private EventStreamServiceGrpc.EventStreamServiceStub stub;
    private final String host;
    private final String hostname;
    private final int port;

    private EventingGrpcClient(String host, String hostname, int port) {

        this.host = host;
        this.hostname = hostname;
        this.port = port;
        initConnection();
    }

    private void initConnection() {

        if (GRPCUtils.isReInitRequired(channel)) {
            if (channel != null && !channel.isShutdown()) {
                channel.shutdownNow();
                do {
                    try {
                        channel.awaitTermination(100, TimeUnit.MILLISECONDS);
                    } catch (InterruptedException e) {
                        logger.error("JWTIssuer discovery channel shutdown wait was interrupted", e);
                    }
                } while (!channel.isShutdown());
            }
            Metadata metadata = new Metadata();
            String connectionId = UUID.randomUUID().toString();
            metadata.put(Metadata.Key.of("enforcer-uuid", Metadata.ASCII_STRING_MARSHALLER),
                    connectionId);
            logger.info("Enforcer UUID: " + connectionId);
            this.channel = GRPCUtils.createSecuredChannel(logger, host, port, hostname);
            ClientInterceptor metadataInterceptor = MetadataUtils.newAttachHeadersInterceptor(metadata);

            this.stub = EventStreamServiceGrpc.newStub(channel).withInterceptors(metadataInterceptor);
        } else if (channel.getState(true) == ConnectivityState.READY) {
            XdsSchedulerManager.getInstance().stopEventStreamScheduling();
        }
    }

    public static EventingGrpcClient getInstance() {

        if (instance == null) {
            String sdsHost = ConfigHolder.getInstance().getEnvVarConfig().getCommonControllerHost();
            String sdsHostname = ConfigHolder.getInstance().getEnvVarConfig().getCommonControllerHostname();
            int sdsPort = Integer.parseInt(ConfigHolder.getInstance().getEnvVarConfig().getCommonControllerXdsPort());
            instance = new EventingGrpcClient(sdsHost, sdsHostname, sdsPort);
        }
        return instance;
    }

    public void run() {

        initConnection();
        watchEvents();
    }

    public void watchEvents() {

        Request request = Request.newBuilder().setEvent("event").build();

        stub.streamEvents(request, new StreamObserver<>() {
            @Override
            public void onNext(Event event) {

                handleNotificationEvent(event);
                XdsSchedulerManager.getInstance().stopEventStreamScheduling();
            }

            @Override
            public void onError(Throwable t) {

                logger.error("Event error", t);
                XdsSchedulerManager.getInstance().startEventScheduling();

            }

            @Override
            public void onCompleted() {

                logger.info("Completed====");
            }
        });
    }
    private void handleNotificationEvent(Event event) {

        switch (event.getType()) {
            case "ALL_EVENTS":
                logger.info("Received all events from the server");
                SubscriptionDataStoreUtil.getInstance().loadStartupArtifacts();
                break;
            case "APPLICATION_CREATED":
                Application application = event.getApplication();
                SubscriptionDataStoreUtil.addApplication(application);
                break;
            case "SUBSCRIPTION_CREATED":
            case "SUBSCRIPTION_UPDATED":
                SubscriptionDataStoreUtil.addSubscription(event.getSubscription());

                break;
            case "APPLICATION_MAPPING_CREATED":
            case "APPLICATION_MAPPING_UPDATED":
                SubscriptionDataStoreUtil.addApplicationMapping(event.getApplicationMapping());
                break;
            case "APPLICATION_KEY_MAPPING_CREATED":
            case "APPLICATION_KEY_MAPPING_UPDATED":
                SubscriptionDataStoreUtil.addApplicationKeyMapping(event.getApplicationKeyMapping());
                break;
            case "APPLICATION_UPDATED":
                SubscriptionDataStoreUtil.addApplication(event.getApplication());
                break;
            case "APPLICATION_MAPPING_DELETED":
                SubscriptionDataStoreUtil.removeApplicationMapping(event.getApplicationMapping());
                break;
            case "APPLICATION_KEY_MAPPING_DELETED":
                SubscriptionDataStoreUtil.removeApplicationKeyMapping(event.getApplicationKeyMapping());
                break;
            case "SUBSCRIPTION_DELETED":
                SubscriptionDataStoreUtil.removeSubscription(event.getSubscription());
                break;
            case "APPLICATION_DELETED":
                SubscriptionDataStoreUtil.removeApplication(event.getApplication());
                break;
            default:
                logger.error("Unknown event type received from the server");
                break;
        }
    }
}
