package org.wso2.apk.integration.utils.clients;

import io.grpc.StatusRuntimeException;
import io.grpc.netty.GrpcSslContexts;
import io.grpc.netty.NettyChannelBuilder;
import io.netty.handler.ssl.SslContext;
import io.netty.handler.ssl.util.InsecureTrustManagerFactory;
import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;

import io.grpc.ManagedChannel;
import io.grpc.Metadata;

import org.wso2.apk.integration.utils.GenericClientInterceptor;
import org.wso2.apk.integration.utils.clients.student_service.StudentRequest;
import org.wso2.apk.integration.utils.clients.student_service.StudentResponse;
import org.wso2.apk.integration.utils.clients.student_service.StudentServiceDefaultVersionGrpc;
import org.wso2.apk.integration.utils.clients.student_service.StudentServiceGrpc;

import javax.net.ssl.SSLException;
import java.util.Map;
import java.util.concurrent.TimeUnit;

public class SimpleGRPCStudentClient {
    protected Log log = LogFactory.getLog(SimpleGRPCStudentClient.class);
    private static final int EVENTUAL_SUCCESS_RESPONSE_TIMEOUT_IN_SECONDS = 10;
    private final String host;
    private final int port;
    private static Metadata responseHeaders;

    public SimpleGRPCStudentClient(String host, int port) {
        this.host = host;
        this.port = port;
    }

    public StudentResponse GetStudent(Map<String, String> headers) throws StatusRuntimeException {
        ManagedChannel managedChannel = null;
        try {
            SslContext sslContext = GrpcSslContexts.forClient()
                    .trustManager(InsecureTrustManagerFactory.INSTANCE)
                    .build();

            GenericClientInterceptor interceptor = new GenericClientInterceptor(headers);
            managedChannel = NettyChannelBuilder.forAddress(host, port)
                    .sslContext(sslContext)
                    .intercept(interceptor)
                    .build();
            StudentServiceGrpc.StudentServiceBlockingStub blockingStub = StudentServiceGrpc
                    .newBlockingStub(managedChannel);
            if (blockingStub == null) {
                log.error("Failed to create blocking stub");
                throw new RuntimeException("Failed to create blocking stub");
            }
            StudentResponse response = blockingStub.getStudent(StudentRequest.newBuilder().setId(1).build());
            if (response == null) {
                log.error("Failed to get student");
                throw new RuntimeException("Failed to get student");
            }
            setResponseHeaders(interceptor.getResponseHeaders());
            return response;
        } catch (SSLException e) {
            throw new RuntimeException("Failed to create SSL context", e);
        } finally {
            // Shut down the channel to release resources
            if (managedChannel != null) {
                managedChannel.shutdown(); // Initiates a graceful shutdown
                try {
                    // Wait at most 5 seconds for the channel to terminate
                    if (!managedChannel.awaitTermination(5, TimeUnit.SECONDS)) {
                        managedChannel.shutdownNow(); // Force shutdown if it does not complete within the timeout
                    }
                } catch (InterruptedException ie) {
                    managedChannel.shutdownNow(); // Force shutdown if the thread is interrupted
                }
            }
        }
    }

    public static String getResponseHeader(String headerName) {
        Metadata.Key<String> headerValue = Metadata.Key.of(headerName, Metadata.ASCII_STRING_MARSHALLER);
        if (responseHeaders == null) {
            return "";
        }
        return responseHeaders.get(headerValue);
    }

    public void setResponseHeaders(Metadata metadata) {
        SimpleGRPCStudentClient.responseHeaders = metadata;
    }

    public StudentResponse GetStudentDefaultVersion(Map<String, String> headers) throws StatusRuntimeException {
        ManagedChannel managedChannel = null;
        try {
            SslContext sslContext = GrpcSslContexts.forClient()
                    .trustManager(InsecureTrustManagerFactory.INSTANCE)
                    .build();

            GenericClientInterceptor interceptor = new GenericClientInterceptor(headers);
            managedChannel = NettyChannelBuilder.forAddress(host, port)
                    .sslContext(sslContext)
                    .intercept(interceptor)
                    .build();
            StudentServiceDefaultVersionGrpc.StudentServiceBlockingStub blockingStub = StudentServiceDefaultVersionGrpc
                    .newBlockingStub(managedChannel);
            if (blockingStub == null) {
                log.error("Failed to create blocking stub");
                throw new RuntimeException("Failed to create blocking stub");
            }
            StudentResponse response = blockingStub.getStudent(StudentRequest.newBuilder().setId(1).build());
            if (response == null) {
                log.error("Failed to get student");
                throw new RuntimeException("Failed to get student");
            }
            return response;
        } catch (SSLException e) {
            throw new RuntimeException("Failed to create SSL context", e);
        } finally {
            // Shut down the channel to release resources
            if (managedChannel != null) {
                managedChannel.shutdown(); // Initiates a graceful shutdown
                try {
                    // Wait at most 5 seconds for the channel to terminate
                    if (!managedChannel.awaitTermination(5, TimeUnit.SECONDS)) {
                        managedChannel.shutdownNow(); // Force shutdown if it does not complete within the timeout
                    }
                } catch (InterruptedException ie) {
                    managedChannel.shutdownNow(); // Force shutdown if the thread is interrupted
                }
            }
        }
    }

}
