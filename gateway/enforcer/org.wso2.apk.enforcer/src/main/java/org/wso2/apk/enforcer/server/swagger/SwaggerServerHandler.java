package org.wso2.apk.enforcer.server.swagger;

import io.grpc.netty.shaded.io.netty.buffer.Unpooled;
import io.grpc.netty.shaded.io.netty.channel.ChannelHandlerContext;
import io.grpc.netty.shaded.io.netty.channel.SimpleChannelInboundHandler;
import io.grpc.netty.shaded.io.netty.handler.codec.http.DefaultFullHttpResponse;
import io.grpc.netty.shaded.io.netty.handler.codec.http.FullHttpResponse;
import io.grpc.netty.shaded.io.netty.handler.codec.http.HttpObject;
import io.grpc.netty.shaded.io.netty.handler.codec.http.HttpRequest;
import io.grpc.netty.shaded.io.netty.handler.codec.http.HttpResponseStatus;
import io.grpc.netty.shaded.io.netty.handler.codec.http.HttpVersion;
import io.grpc.netty.shaded.io.netty.util.CharsetUtil;
import org.apache.http.protocol.HTTP;
import org.wso2.apk.enforcer.api.APIFactory;
import org.wso2.apk.enforcer.constants.APIDefinitionConstants;
import org.wso2.apk.enforcer.constants.AdminConstants;
import org.wso2.apk.enforcer.constants.HttpConstants;
import org.wso2.apk.enforcer.models.ResponsePayload;
import java.util.HashMap;
import java.util.Map;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

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

        boolean isSwagger = false;

        //check if it's GRPC request using the header
        String header = request.headers().get("ApiType");
        String[] params = request.uri().split("/");
        final String basePath;
        final String vHost;
        final String queryParam;
        final String version;
        //if len params is 3, then it's a GRPC request else other
        final String type = params.length == 3 ? "GRPC" : "REST";
        if (type.equals("GRPC")) {
            basePath = "/" + params[1];
            vHost = params[2].split("\\?")[0];
            queryParam = params[2].split("\\?")[1];
            version = extractVersionFromGrpcBasePath(params[1]);
        } else {
            basePath = "/" + params[1] + "/" + params[2];
            vHost = params[3].split("\\?")[0];
            queryParam = params[3].split("\\?")[1];
            version = params[2];
        }

        if (APIDefinitionConstants.SWAGGER_DEFINITION.equalsIgnoreCase(queryParam)) {
            isSwagger = true;
        }
        if(isSwagger){
            // load the corresponding swagger definition from the API name
            byte[] apiDefinition = apiFactory.getAPIDefinition(basePath, version, vHost);
            if (apiDefinition == null || apiDefinition.length == 0) {
                String error = AdminConstants.ErrorMessages.NOT_FOUND;
                responsePayload = new ResponsePayload();
                responsePayload.setError(true);
                responsePayload.setStatus(HttpResponseStatus.NOT_FOUND);
                responsePayload.setContent(error);
                buildAndSendResponse(ctx, responsePayload);
                return;
            }
            Map<String,String> map = new HashMap<>();
            map.put("apiDefinition", APIDefinitionUtils.ReadGzip(apiDefinition));
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

    private static String extractVersionFromGrpcBasePath(String input) {
        // Pattern to match '.v' followed by digits and optional periods followed by more digits
        Pattern pattern = Pattern.compile("\\.v\\d+(\\.\\d+)*");
        Matcher matcher = pattern.matcher(input);

        if (matcher.find()) {
            return matcher.group().substring(1); // Returns the matched version
        }
        return "";
    }
}
