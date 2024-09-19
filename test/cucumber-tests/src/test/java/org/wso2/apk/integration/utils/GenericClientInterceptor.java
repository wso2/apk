package org.wso2.apk.integration.utils;

import io.grpc.ClientInterceptor;
import io.grpc.ForwardingClientCall;
import io.grpc.Metadata;
import io.grpc.MethodDescriptor;
import io.grpc.CallOptions;
import io.grpc.ClientCall;
import io.grpc.Channel;
import java.util.Map;

import java.util.Map;

public class GenericClientInterceptor implements ClientInterceptor {

    private Map<String, String> headers;

    public GenericClientInterceptor(Map<String, String> headers) {
        this.headers = headers;
    }

    @Override
    public <ReqT, RespT> ClientCall<ReqT, RespT> interceptCall(
            MethodDescriptor<ReqT, RespT> method, CallOptions callOptions, Channel next) {
        return new ForwardingClientCall.SimpleForwardingClientCall<ReqT, RespT>(
                next.newCall(method, callOptions)) {

            @Override
            public void start(Listener<RespT> responseListener, Metadata headersMetadata) {
                // Set each header in the map to the Metadata headers
                headers.forEach((key, value) -> headersMetadata.put(
                        Metadata.Key.of(key, Metadata.ASCII_STRING_MARSHALLER), value));

                super.start(responseListener, headersMetadata);
            }
        };
    }
}