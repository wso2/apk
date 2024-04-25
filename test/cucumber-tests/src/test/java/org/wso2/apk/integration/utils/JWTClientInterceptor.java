package org.wso2.apk.integration.utils;
import io.grpc.*;

public class JWTClientInterceptor implements ClientInterceptor {

    private String jwtToken;

    public JWTClientInterceptor(String jwtToken) {
        this.jwtToken = jwtToken;
    }

    @Override
    public <ReqT, RespT> ClientCall<ReqT, RespT> interceptCall(
            MethodDescriptor<ReqT, RespT> method, CallOptions callOptions, Channel next) {
        return new ForwardingClientCall.SimpleForwardingClientCall<ReqT, RespT>(
                next.newCall(method, callOptions)) {

            @Override
            public void start(Listener<RespT> responseListener, Metadata headers) {
                headers.put(
                        Metadata.Key.of("Authorization", Metadata.ASCII_STRING_MARSHALLER),
                        "Bearer " + jwtToken);
                super.start(responseListener, headers);
            }
        };
    }
}

