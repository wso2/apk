/*
 * Copyright (c) 2020, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.wso2.apk.enforcer.server;

import io.grpc.Server;
import io.grpc.ServerInterceptors;
import io.grpc.netty.shaded.io.grpc.netty.NettyServerBuilder;
import io.grpc.netty.shaded.io.netty.channel.EventLoopGroup;
import io.grpc.netty.shaded.io.netty.channel.nio.NioEventLoopGroup;
import io.grpc.netty.shaded.io.netty.channel.socket.nio.NioServerSocketChannel;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.wso2.apk.enforcer.analytics.AccessLoggingService;
import org.wso2.apk.enforcer.api.APIFactory;
import org.wso2.apk.enforcer.commons.logging.ErrorDetails;
import org.wso2.apk.enforcer.commons.logging.LoggingConstants;
import org.wso2.apk.enforcer.config.ConfigHolder;
import org.wso2.apk.enforcer.config.EnforcerConfig;
import org.wso2.apk.enforcer.config.dto.AuthServiceConfigurationDto;
import org.wso2.apk.enforcer.config.dto.ThreadPoolConfig;
import org.wso2.apk.enforcer.discovery.ConfigDiscoveryClient;
import org.wso2.apk.enforcer.grpc.ExtAuthService;
 import org.wso2.apk.enforcer.grpc.ExternalProcessorService;
import org.wso2.apk.enforcer.grpc.HealthService;
import org.wso2.apk.enforcer.grpc.interceptors.AccessLogInterceptor;
import org.wso2.apk.enforcer.grpc.interceptors.OpenTelemetryInterceptor;
import org.wso2.apk.enforcer.jmx.JMXAgent;
import org.wso2.apk.enforcer.metrics.MetricsManager;
import org.wso2.apk.enforcer.security.jwt.validator.RevokedJWTDataHolder;
import org.wso2.apk.enforcer.subscription.SubscriptionDataStoreUtil;
import org.wso2.apk.enforcer.tracing.TracerFactory;
import org.wso2.apk.enforcer.tracing.TracingException;
import org.wso2.apk.enforcer.tracing.Utils;
import org.wso2.apk.enforcer.util.TLSUtils;

import java.io.IOException;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.TimeUnit;
import javax.net.ssl.SSLException;

/**
 * gRPC netty based server that handles the incoming requests.
 */
public class AuthServer {

    private static final Logger logger = LogManager.getLogger(AuthServer.class);

    public static void main(String[] args) {
        try {
            // initialize the config holder
            ConfigHolder.getInstance();

            // wait until configurations are fetched from cds
            CountDownLatch latch = new CountDownLatch(1);
            ConfigDiscoveryClient cds = ConfigDiscoveryClient.init(latch);
            cds.requestInitConfig();
            try {
                latch.await();
            } catch (InterruptedException e) {
                logger.error("Error while waiting for configurations from adapter",
                        ErrorDetails.errorLog(LoggingConstants.Severity.BLOCKER, 6700), e);
                System.exit(1);
            }

            EnforcerConfig enforcerConfig = ConfigHolder.getInstance().getConfig();

            // rest api call initialization happens here
            APIFactory.getInstance().init();

            // Initialize tracing objects
            if (enforcerConfig.getTracingConfig().isTracingEnabled()) {
                try {
                    TracerFactory.getInstance().init();
                    Utils.setTracingEnabled(true);
                    logger.info("Tracing is enabled.");
                } catch (TracingException e) {
                    logger.error("Error enabling tracing",
                            ErrorDetails.errorLog(LoggingConstants.Severity.CRITICAL, 6901), e);
                    // prevent further tracing activation by disabling the config
                    Utils.setTracingEnabled(false);
                }
            } else {
                logger.debug("Tracing is disabled.");
            }

            // Create a new server to listen on port 8081
            Server server = initServer();

            // Enable global filters
            if (enforcerConfig.getAnalyticsConfig().isEnabled() ||
                    enforcerConfig.getMetricsConfig().isMetricsEnabled()) {
                AccessLoggingService accessLoggingService = new AccessLoggingService();
                accessLoggingService.init();
                if (enforcerConfig.getMetricsConfig().isMetricsEnabled()) {
                    //Initialize metrics
                    MetricsManager.initializeMetrics(enforcerConfig.getMetricsConfig());
                }
                if (enforcerConfig.getAnalyticsConfig().isEnabled()) {
                    logger.info("analytics filter is enabled.");
                }
            } else {
                logger.debug("analytics filter is disabled.");
            }
            // Start receiving revoked tokens from redis cache
            RevokedTokenRedisClient.retrieveAndSubscribe();

            // Start the server
            server.start();
            logger.info("Sever started Listening in port : " + 8081);

            // Initialize JMX Agent
            JMXAgent.initJMXAgent();
            SubscriptionDataStoreUtil.initializeLoadingTasks();;
            SubscriptionDataStoreUtil.getInstance().loadStartupArtifacts();
            RevokedJWTDataHolder.getInstance().init();

            // Create a new server to listen on port 8082
            RestServer restServer = new RestServer();
            restServer.initServer();

            // Don't exit the main thread. Wait until server is terminated.
            server.awaitTermination();
        } catch (IOException e) {
            logger.error("Error while starting the enforcer gRPC server or http server.",
                    ErrorDetails.errorLog(LoggingConstants.Severity.BLOCKER, 6702), e);
            System.exit(1);
        } catch (InterruptedException e) {
            logger.error("Enforcer server main thread interrupted.",
                    ErrorDetails.errorLog(LoggingConstants.Severity.BLOCKER, 6703), e);
            System.exit(1);
        } catch (Exception ex) {
            // Printing the stack trace in case logger might not have been initialized
            ex.printStackTrace();
            System.exit(1);
        }
    }

    private static Server initServer() throws SSLException {
        final EventLoopGroup bossGroup = new NioEventLoopGroup(Runtime.getRuntime().availableProcessors());
        final EventLoopGroup workerGroup = new NioEventLoopGroup(Runtime.getRuntime().availableProcessors() * 2);
        AuthServiceConfigurationDto authServerConfig = ConfigHolder.getInstance().getConfig().getAuthService();
        ThreadPoolConfig threadPoolConfig = authServerConfig.getThreadPool();
        EnforcerWorkerPool enforcerWorkerPool = new EnforcerWorkerPool(threadPoolConfig.getCoreSize(),
                threadPoolConfig.getMaxSize(), threadPoolConfig.getKeepAliveTime(), threadPoolConfig.getQueueSize(),
                Constants.EXTERNAL_AUTHZ_THREAD_GROUP, Constants.EXTERNAL_AUTHZ_THREAD_ID);
        System.out.println("test");
        return NettyServerBuilder.forPort(authServerConfig.getPort())
                .keepAliveTime(authServerConfig.getKeepAliveTime(), TimeUnit.SECONDS).bossEventLoopGroup(bossGroup)
                .workerEventLoopGroup(workerGroup)
                .addService(ServerInterceptors.intercept(new ExtAuthService(), new OpenTelemetryInterceptor(),
                        new AccessLogInterceptor()))
                .addService(new HealthService())
                .addService(new ExternalProcessorService())
//                .addService(ServerInterceptors.intercept(new WebSocketFrameService(), new AccessLogInterceptor()))
                .maxInboundMessageSize(authServerConfig.getMaxMessageSize())
                .maxInboundMetadataSize(authServerConfig.getMaxHeaderLimit()).channelType(NioServerSocketChannel.class)
                .executor(enforcerWorkerPool.getExecutor())
                .sslContext(TLSUtils.buildGRPCServerSSLContext())
                .build();
    }
}
