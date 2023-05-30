package org.wso2.apk.enforcer.server.swagger;

import io.grpc.netty.shaded.io.netty.buffer.Unpooled;
import io.grpc.netty.shaded.io.netty.channel.ChannelHandlerContext;
import io.grpc.netty.shaded.io.netty.channel.ChannelInboundHandlerAdapter;
import io.grpc.netty.shaded.io.netty.handler.codec.http.DefaultFullHttpResponse;
import io.grpc.netty.shaded.io.netty.handler.codec.http.FullHttpResponse;
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

import java.util.Arrays;
import java.util.HashMap;
import java.util.Map;

public class SwaggerServerHandler extends ChannelInboundHandlerAdapter {

    final SubscriptionDataStore dataStore = SubscriptionDataHolder.getInstance().getTenantSubscriptionStore();
    private final static APIFactory apiFactory = APIFactory.getInstance();

    @Override
    public void channelRead(ChannelHandlerContext ctx, Object msg) throws Exception {
        FullHttpRequest request;
        ResponsePayload responsePayload;

        if(msg instanceof FullHttpRequest) {
            request = (FullHttpRequest) msg;
        } else {
            String error = AdminConstants.ErrorMessages.INTERNAL_SERVER_ERROR;
            responsePayload = new ResponsePayload();
            responsePayload.setError(true);
            responsePayload.setStatus(HttpResponseStatus.INTERNAL_SERVER_ERROR);
            responsePayload.setContent(error);
            buildAndSendResponse(ctx, responsePayload);
            return;
        }
        // Check request uri and invoke correct handler
        String[] uriSections = request.uri().split("\\?");
        String[] params = null;
        String baseURI = uriSections[0];
        if (uriSections.length > 1) {
            params = uriSections[1].split("&");
        }
        boolean isSwagger = Arrays.stream(params).anyMatch(param -> {
            if(APIDefinitionConstants.SWAGGER_DEFINITION.equalsIgnoreCase(param)) {
                return true;
            }
            return false;
        });
        System.out.println("API Base Path: " + params[0]);
        if(isSwagger){
            // load the corresponding swagger definition from the API name
            byte[] apiDefinition = apiFactory.getAPIDefinition(params[0], params[1], params[2]);
            System.out.println("API Base Path: " + params[0]);
            System.out.println("API Definition: " + new String(apiDefinition));
            Map<String, byte[]> map = new HashMap<>();
            map.put("swagger.json", apiDefinition);
            responsePayload = APIDefinitionUtils.buildResponsePayload(map, HttpResponseStatus.OK, false);
            buildAndSendResponse(ctx, responsePayload);
        }
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

    private void getSwaggerDefinition() {

    }
}
