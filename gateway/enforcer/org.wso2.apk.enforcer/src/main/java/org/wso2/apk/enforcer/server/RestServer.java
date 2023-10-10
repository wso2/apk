/*
 * Copyright (c) 2021, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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

import io.grpc.netty.shaded.io.netty.bootstrap.ServerBootstrap;
import io.grpc.netty.shaded.io.netty.channel.Channel;
import io.grpc.netty.shaded.io.netty.channel.ChannelOption;
import io.grpc.netty.shaded.io.netty.channel.EventLoopGroup;
import io.grpc.netty.shaded.io.netty.channel.nio.NioEventLoopGroup;
import io.grpc.netty.shaded.io.netty.channel.socket.nio.NioServerSocketChannel;
import io.grpc.netty.shaded.io.netty.handler.logging.LogLevel;
import io.grpc.netty.shaded.io.netty.handler.logging.LoggingHandler;
import io.grpc.netty.shaded.io.netty.handler.ssl.SslContext;
import io.grpc.netty.shaded.io.netty.handler.ssl.SslContextBuilder;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.wso2.apk.enforcer.config.ConfigHolder;
import org.wso2.apk.enforcer.jwks.JWKSRequestHandler;
import org.wso2.apk.enforcer.jwks.JWKSServerInitializer;
import org.wso2.apk.enforcer.server.swagger.SwaggerServerInitializer;

import java.io.File;
import java.nio.file.Paths;

import javax.net.ssl.SSLException;

/**
 * TokenServer to handle JWT /testkey endpoint backend HTTPS service.
 */
public class RestServer {

    private static final Logger logger = LogManager.getLogger(RestServer.class);

    public void initServer() throws SSLException, InterruptedException {

        logger.info("New Rest Server New");
        // Configure SSL
        final SslContext sslCtx;
        final SslContextBuilder ssl;
        File certFile = Paths.get(ConfigHolder.getInstance().getEnvVarConfig().getEnforcerPublicKeyPath()).toFile();
        File keyFile = Paths.get(ConfigHolder.getInstance().getEnvVarConfig().getEnforcerPrivateKeyPath()).toFile();
        ssl = SslContextBuilder.forServer(certFile, keyFile);
        ssl.trustManager(ConfigHolder.getInstance().getTrustManagerFactory());
        sslCtx = ssl.build();

        // Create the multithreaded event loops for the server
        final EventLoopGroup bossGroup = new NioEventLoopGroup(Runtime.getRuntime().availableProcessors());
        final EventLoopGroup workerGroup = new NioEventLoopGroup(Runtime.getRuntime().availableProcessors() * 2);
        try {
            ServerBootstrap swaggerServer = new ServerBootstrap();
            swaggerServer.option(ChannelOption.SO_BACKLOG, 1024);
            swaggerServer.group(bossGroup, workerGroup)
                    .channel(NioServerSocketChannel.class)
                    .handler(new LoggingHandler(LogLevel.INFO))
                    .childHandler(new SwaggerServerInitializer(sslCtx));
            Channel swaggerChannel = swaggerServer.bind(8084).sync().channel();
            logger.info("API Definition endpoint started Listening in port : " + 8084);

            ServerBootstrap jWKSServer = new ServerBootstrap();
            jWKSServer.option(ChannelOption.SO_BACKLOG, 1024);
            jWKSServer.group(bossGroup, workerGroup)
                    .channel(NioServerSocketChannel.class)
                    .handler(new JWKSRequestHandler())
                    .childHandler(new JWKSServerInitializer(sslCtx));
            Channel jwksChannel = jWKSServer.bind(9092).sync().channel();
            logger.info("JWKS endpoint started Listening in port : " + 9092);
            jwksChannel.closeFuture().sync();

            swaggerChannel.closeFuture().sync();

        } finally {
            workerGroup.shutdownGracefully();
            bossGroup.shutdownGracefully();
        }
    }
}
