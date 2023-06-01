package org.wso2.apk.enforcer.server.swagger;

import io.grpc.netty.shaded.io.netty.buffer.Unpooled;
import io.grpc.netty.shaded.io.netty.channel.ChannelHandlerContext;
import io.grpc.netty.shaded.io.netty.channel.ChannelInboundHandlerAdapter;
import io.grpc.netty.shaded.io.netty.channel.SimpleChannelInboundHandler;
import io.grpc.netty.shaded.io.netty.handler.codec.http.DefaultFullHttpResponse;
import io.grpc.netty.shaded.io.netty.handler.codec.http.FullHttpMessage;
import io.grpc.netty.shaded.io.netty.handler.codec.http.FullHttpResponse;
import io.grpc.netty.shaded.io.netty.handler.codec.http.HttpObject;
import io.grpc.netty.shaded.io.netty.handler.codec.http.HttpRequest;
import io.grpc.netty.shaded.io.netty.handler.codec.http.HttpResponseStatus;
import io.grpc.netty.shaded.io.netty.handler.codec.http.HttpVersion;
import io.grpc.netty.shaded.io.netty.util.CharsetUtil;
import io.netty.handler.codec.http.FullHttpRequest;
import org.apache.http.protocol.HTTP;
import org.wso2.apk.enforcer.api.APIFactory;
import org.wso2.apk.enforcer.constants.APIDefinitionConstants;
import org.wso2.apk.enforcer.constants.AdminConstants;
import org.wso2.apk.enforcer.constants.HttpConstants;
import org.wso2.apk.enforcer.models.ResponsePayload;
import org.wso2.apk.enforcer.subscription.SubscriptionDataHolder;
import org.wso2.apk.enforcer.subscription.SubscriptionDataStore;

import java.nio.charset.StandardCharsets;
import java.util.Arrays;
import java.util.HashMap;
import java.util.Map;

public class SwaggerServerHandler extends SimpleChannelInboundHandler<HttpObject> {

//    final SubscriptionDataStore dataStore = SubscriptionDataHolder.getInstance().getTenantSubscriptionStore();
    private final static APIFactory apiFactory = APIFactory.getInstance();

    @Override
    public void channelRead0(ChannelHandlerContext ctx, HttpObject msg) throws Exception {
        HttpRequest request;
        ResponsePayload responsePayload;

        if(msg instanceof HttpRequest) {
            request = (HttpRequest) msg;
        } else {
            String error = AdminConstants.ErrorMessages.INTERNAL_SERVER_ERROR;
            responsePayload = new ResponsePayload();
            responsePayload.setError(true);
            responsePayload.setStatus(HttpResponseStatus.INTERNAL_SERVER_ERROR);
            responsePayload.setContent(error);
            buildAndSendResponse(ctx, responsePayload);
            return;
        }

        String [] params = request.uri().split("/");
        boolean isSwagger = Arrays.stream(params).anyMatch(param -> {
            if(APIDefinitionConstants.SWAGGER_DEFINITION.equalsIgnoreCase(param)) {
                return true;
            }
            return false;
        });
        final String basePath = "/" + params[1] + "/" + params[2];
        if(isSwagger){
            // load the corresponding swagger definition from the API name
            byte[] apiDefinition = apiFactory.getAPIDefinition(basePath, params[2], params[3]);
            Map<String,String> map = new HashMap<>();
            map.put("swagger.json", new String(apiDefinition, StandardCharsets.UTF_8));
            responsePayload = APIDefinitionUtils.buildResponsePayload(map, HttpResponseStatus.OK, false);
            buildAndSendResponse(ctx, responsePayload);
        }

        // If the request is not for swagger definition, then send a bad request
        String error = AdminConstants.ErrorMessages.BAD_REQUEST;
        responsePayload = new ResponsePayload();
        responsePayload.setError(true);
        responsePayload.setStatus(HttpResponseStatus.BAD_REQUEST);
        responsePayload.setContent(error);
        buildAndSendResponse(ctx, responsePayload);
    }

    @Override
    public void channelReadComplete(ChannelHandlerContext ctx) {
        ctx.flush();
    }

    private void buildAndSendResponse(ChannelHandlerContext ctx, ResponsePayload response) {
        FullHttpResponse httpResponse = new DefaultFullHttpResponse(HttpVersion.HTTP_1_1,
                response.getStatus(),
                Unpooled.copiedBuffer(response.getContent(), CharsetUtil.UTF_8));
        httpResponse.headers().set(HTTP.CONTENT_TYPE, HttpConstants.APPLICATION_JSON);
        httpResponse.headers().set(HTTP.CONTENT_LEN, httpResponse.content().readableBytes());
        ctx.writeAndFlush(httpResponse);
    }
}
