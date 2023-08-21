/*
 * Copyright (c) 2023, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.wso2.apk.enforcer.jwks;

import com.nimbusds.jose.jwk.JWKSet;
import io.grpc.netty.shaded.io.netty.buffer.Unpooled;
import io.grpc.netty.shaded.io.netty.channel.ChannelFuture;
import io.grpc.netty.shaded.io.netty.channel.ChannelFutureListener;
import io.grpc.netty.shaded.io.netty.channel.ChannelHandlerContext;
import io.grpc.netty.shaded.io.netty.channel.SimpleChannelInboundHandler;
import io.grpc.netty.shaded.io.netty.handler.codec.http.DefaultFullHttpResponse;
import io.grpc.netty.shaded.io.netty.handler.codec.http.FullHttpRequest;
import io.grpc.netty.shaded.io.netty.handler.codec.http.FullHttpResponse;
import io.grpc.netty.shaded.io.netty.handler.codec.http.HttpMethod;
import io.grpc.netty.shaded.io.netty.handler.codec.http.HttpObject;
import io.grpc.netty.shaded.io.netty.handler.codec.http.HttpRequest;
import io.grpc.netty.shaded.io.netty.handler.codec.http.HttpResponseStatus;
import io.grpc.netty.shaded.io.netty.handler.codec.http.HttpUtil;
import io.grpc.netty.shaded.io.netty.handler.codec.http.HttpVersion;
import org.apache.http.protocol.HTTP;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.wso2.apk.enforcer.config.ConfigHolder;
import org.wso2.apk.enforcer.constants.HttpConstants;

/**
 * JWKS Request Handler for Backend JWTs
 */
public class JWKSRequestHandler extends SimpleChannelInboundHandler<HttpObject> {
    private static final Logger logger = LogManager.getLogger(JWKSRequestHandler.class);
    private static final String route = "/jwks";
    private static final String CONNECTION = "Connection";
    private static final String CLOSE = "close";
    private static final String KEEP_ALIVE = "keep-alive";
    @Override
    public void channelReadComplete(ChannelHandlerContext ctx) throws Exception {

        ctx.flush();
    }

    @Override
    protected void channelRead0(ChannelHandlerContext ctx, HttpObject msg) {
        FullHttpResponse res = null;
        FullHttpRequest req = null;
        BackendJWKSDto backendJWKSDto = ConfigHolder.getInstance().getConfig().getBackendJWKSDto();
        JWKSet jwks = backendJWKSDto.getJwks();
        if (msg instanceof HttpRequest) {
            req = (FullHttpRequest) msg;
            boolean keepAlive = HttpUtil.isKeepAlive(req);
            String path = req.uri().split("\\?")[0]; //Get the context without query params

            if (!(HttpMethod.GET.equals(req.method()) && path.equals(route))) {
                ctx.fireChannelRead(msg);
                return;
            }
            res = new DefaultFullHttpResponse(HttpVersion.HTTP_1_1, HttpResponseStatus.OK,
                    Unpooled.wrappedBuffer(jwks.toJSONObject().toString().getBytes()));
            res.headers()
                    .set(HTTP.CONN_DIRECTIVE, HTTP.CONN_KEEP_ALIVE)
                    .set(HTTP.CONTENT_TYPE, HttpConstants.APPLICATION_JSON)
                    .setInt(HTTP.CONTENT_LEN, res.content().readableBytes());
            if (keepAlive) {
                if (!req.protocolVersion().isKeepAliveDefault()) {
                    res.headers().set(CONNECTION, KEEP_ALIVE);
                }
            } else {
                // Tell the client we're going to close the connection.
                res.headers().set(CONNECTION, CLOSE);
            }
            ChannelFuture f = ctx.write(res);
            if (!keepAlive) {
                f.addListener(ChannelFutureListener.CLOSE);
            }
        }
    }
}