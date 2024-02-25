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

package org.wso2.apk.enforcer.discovery;

import com.google.protobuf.Any;
import com.google.rpc.Status;
import io.envoyproxy.envoy.config.core.v3.Node;
import io.envoyproxy.envoy.service.discovery.v3.DiscoveryRequest;
import io.envoyproxy.envoy.service.discovery.v3.DiscoveryResponse;
import io.grpc.ConnectivityState;
import io.grpc.ManagedChannel;
import io.grpc.stub.StreamObserver;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.wso2.apk.enforcer.config.ConfigHolder;
import org.wso2.apk.enforcer.constants.AdapterConstants;
import org.wso2.apk.enforcer.constants.Constants;
import org.wso2.apk.enforcer.discovery.common.XDSCommonUtils;
import org.wso2.apk.enforcer.discovery.scheduler.XdsSchedulerManager;
import org.wso2.apk.enforcer.discovery.service.subscription.JWTIssuerDiscoveryServiceGrpc;
import org.wso2.apk.enforcer.discovery.subscription.JWTIssuer;
import org.wso2.apk.enforcer.discovery.subscription.JWTIssuerList;
import org.wso2.apk.enforcer.jmx.JMXUtils;
import org.wso2.apk.enforcer.metrics.jmx.impl.ExtAuthMetrics;
import org.wso2.apk.enforcer.subscription.SubscriptionDataHolder;
import org.wso2.apk.enforcer.subscription.SubscriptionDataStore;
import org.wso2.apk.enforcer.util.GRPCUtils;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.concurrent.TimeUnit;

/**
 * Client to communicate with JWTIssuer discovery service at the adapter.
 */
public class JWTIssuerDiscoveryClient implements Runnable {

    private static final Logger logger = LogManager.getLogger(JWTIssuerDiscoveryClient.class);
    private static JWTIssuerDiscoveryClient instance;
    private ManagedChannel channel;
    private JWTIssuerDiscoveryServiceGrpc.JWTIssuerDiscoveryServiceStub stub;
    private StreamObserver<DiscoveryRequest> reqObserver;
    private final String host;
    private final String hostname;
    private final int port;

    /**
     * This is a reference to the latest received response from the ADS.
     * <p>
     * Usage: When ack/nack a DiscoveryResponse this value is used to identify the latest received DiscoveryResponse
     * which may not have been acked/nacked so far.
     * </p>
     */

    private DiscoveryResponse latestReceived;
    /**
     * This is a reference to the latest acked response from the ADS.
     * <p>
     * Usage: When nack a DiscoveryResponse this value is used to find the latest successfully processed
     * DiscoveryResponse. Information sent in the nack request will contain information about this response value.
     * </p>
     */
    private DiscoveryResponse latestACKed;

    /**
     * Node struct for the discovery client
     */
    private final Node node;

    private JWTIssuerDiscoveryClient(String host, String hostname, int port) {

        this.host = host;
        this.hostname = hostname;
        this.port = port;
        initConnection();
        this.node = XDSCommonUtils.generateXDSNode(AdapterConstants.COMMON_ENFORCER_LABEL);
        this.latestACKed = DiscoveryResponse.getDefaultInstance();
    }

    private void initConnection() {

        if (GRPCUtils.isReInitRequired(channel)) {
            if (channel != null && !channel.isShutdown()) {
                channel.shutdownNow();
                do {
                    try {
                        channel.awaitTermination(100, TimeUnit.MILLISECONDS);
                    } catch (InterruptedException e) {
                        logger.error("JWTISsuer discovery channel shutdown wait was interrupted", e);
                    }
                } while (!channel.isShutdown());
            }
            this.channel = GRPCUtils.createSecuredChannel(logger, host, port, hostname);
            this.stub = JWTIssuerDiscoveryServiceGrpc.newStub(channel);
        } else if (channel.getState(true) == ConnectivityState.READY) {
            XdsSchedulerManager.getInstance().stopJWTIssuerDiscoveryScheduling();
        }
    }

    public static JWTIssuerDiscoveryClient getInstance() {

        if (instance == null) {
            String sdsHost = ConfigHolder.getInstance().getEnvVarConfig().getAdapterHost();
            String sdsHostname = ConfigHolder.getInstance().getEnvVarConfig().getAdapterHostname();
            int sdsPort = Integer.parseInt(ConfigHolder.getInstance().getEnvVarConfig().getAdapterXdsPort());
            instance = new JWTIssuerDiscoveryClient(sdsHost, sdsHostname, sdsPort);
        }
        return instance;
    }

    public void run() {

        initConnection();
        watchJWTIssuers();
    }

    public void watchJWTIssuers() {

        reqObserver = stub.streamJWTIssuers(new StreamObserver<>() {
            @Override
            public void onNext(DiscoveryResponse response) {

                logger.info("JWTIssuer creation event received with version : " + response.getVersionInfo());
                logger.debug("Received JWTIssuer discovery response " + response);
                XdsSchedulerManager.getInstance().stopJWTIssuerDiscoveryScheduling();
                latestReceived = response;
                try {
                    List<JWTIssuer> jwtIssuers = new ArrayList<>();
                    for (Any res : response.getResourcesList()) {
                        jwtIssuers.addAll(res.unpack(JWTIssuerList.class).getListList());
                    }
                    Map<String, List<JWTIssuer>> orgWizeIssuerMap = new HashMap<>();
                    for (JWTIssuer jwtIssuer : jwtIssuers) {
                        List<JWTIssuer> jwtIssuerList = orgWizeIssuerMap.computeIfAbsent(jwtIssuer.getOrganization(),
                                k -> new ArrayList<>());
                        jwtIssuerList.add(jwtIssuer);
                    }
                    orgWizeIssuerMap.forEach((k, v) -> {
                        SubscriptionDataStore subscriptionDataStore =
                                SubscriptionDataHolder.getInstance().getSubscriptionDataStore(k);
                        if (subscriptionDataStore == null) {
                            subscriptionDataStore =
                                    SubscriptionDataHolder.getInstance().initializeSubscriptionDataStore(k);
                        }
                        subscriptionDataStore.addJWTIssuers(v);
                    });

                    if (JMXUtils.isJMXMetricsEnabled()) {
                        ExtAuthMetrics.getInstance().recordJWTIssuerMetrics(jwtIssuers.size());
                    }
                    logger.info("Number of jwt issuers received : " + jwtIssuers.size());
                    ack();
                } catch (Exception e) {
                    // catching generic error here to wrap any grpc communication errors in the runtime
                    onError(e);
                }
            }

            @Override
            public void onError(Throwable throwable) {

                logger.error("Error occurred during JWTIssuer discovery", throwable);
                XdsSchedulerManager.getInstance().startJWTIssuerDiscoveryScheduling();
                nack(throwable);
            }

            @Override
            public void onCompleted() {

                logger.info("Completed receiving JWT Issuer data");
            }
        });

        try {
            DiscoveryRequest req = DiscoveryRequest.newBuilder()
                    .setNode(node)
                    .setVersionInfo(latestACKed.getVersionInfo())
                    .setTypeUrl(Constants.JWT_ISSUER_LIST_TYPE_URL).build();
            reqObserver.onNext(req);
            logger.debug("Sent Discovery request for type url: " + Constants.JWT_ISSUER_LIST_TYPE_URL);

        } catch (Exception e) {
            logger.error("Unexpected error occurred in JWTIssuer discovery service", e);
            reqObserver.onError(e);
        }
    }

    /**
     * Send acknowledgement of successfully processed DiscoveryResponse from the xDS server. This is part of the xDS
     * communication protocol.
     */
    private void ack() {

        DiscoveryRequest req = DiscoveryRequest.newBuilder()
                .setNode(node)
                .setVersionInfo(latestReceived.getVersionInfo())
                .setResponseNonce(latestReceived.getNonce())
                .setTypeUrl(Constants.JWT_ISSUER_LIST_TYPE_URL).build();
        reqObserver.onNext(req);
        latestACKed = latestReceived;
    }

    private void nack(Throwable e) {

        if (latestReceived == null) {
            return;
        }
        DiscoveryRequest req = DiscoveryRequest.newBuilder()
                .setNode(node)
                .setVersionInfo(latestACKed.getVersionInfo())
                .setResponseNonce(latestReceived.getNonce())
                .setTypeUrl(Constants.JWT_ISSUER_LIST_TYPE_URL)
                .setErrorDetail(Status.newBuilder().setMessage(e.getMessage()))
                .build();
        reqObserver.onNext(req);
    }
}
