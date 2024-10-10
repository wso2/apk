package org.wso2.apk.integration.utils;

import io.grpc.ClientInterceptor;
import io.grpc.ForwardingClientCall;
import io.grpc.ForwardingClientCall.SimpleForwardingClientCall;
import io.grpc.ForwardingClientCallListener.SimpleForwardingClientCallListener;
import io.grpc.Metadata;
import io.grpc.MethodDescriptor;
import io.grpc.CallOptions;
import io.grpc.ClientCall;
import io.grpc.Channel;
import java.util.Map;

public class GenericClientInterceptor implements ClientInterceptor {

    private Map<String, String> headers;
    private Metadata responseHeaders;

    public GenericClientInterceptor(Map<String, String> headers) {
        this.headers = headers;
    }

    public void setResponseHeaders(Metadata responseHeaders) {
        this.responseHeaders = responseHeaders;
    }

    public Metadata getResponseHeaders() {
        return this.responseHeaders;
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

                super.start(new SimpleForwardingClientCallListener<RespT>(responseListener) {
                    @Override
                    public void onHeaders(Metadata headers) {
                        /**
                         * if you don't need receive header from server,
                         * you can use {@link io.grpc.stub.MetadataUtils#attachHeaders}
                         * directly to send header
                         */
                        setResponseHeaders(headers);
                        super.onHeaders(headers);
                    }
                }, headersMetadata);
            }
        };
    }
}