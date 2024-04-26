package org.wso2.apk.integration.utils.clients;

import io.grpc.StatusRuntimeException;
import io.grpc.netty.GrpcSslContexts;
import io.grpc.netty.NettyChannelBuilder;
import io.netty.handler.ssl.SslContext;
import io.netty.handler.ssl.util.InsecureTrustManagerFactory;
import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import io.grpc.ManagedChannel;
import org.wso2.apk.integration.utils.GenericClientInterceptor;
import org.wso2.apk.integration.utils.clients.studentGrpcClient.StudentRequest;
import org.wso2.apk.integration.utils.clients.studentGrpcClient.StudentResponse;
import org.wso2.apk.integration.utils.clients.studentGrpcClient.StudentServiceGrpc;

import javax.net.ssl.SSLException;
import java.util.Map;


public class SimpleGRPCStudentClient {
    protected Log log = LogFactory.getLog(SimpleGRPCStudentClient.class);
    private static final int EVENTUAL_SUCCESS_RESPONSE_TIMEOUT_IN_SECONDS = 10;
    private final String host;
    private final int port;

    public SimpleGRPCStudentClient(String host, int port) {
        this.host = host;
        this.port = port;
    }

    public StudentResponse GetStudent(Map<String,String> headers) {
        try {
            SslContext sslContext = GrpcSslContexts.forClient()
                    .trustManager(InsecureTrustManagerFactory.INSTANCE)
                    .build();
            GenericClientInterceptor interceptor = new GenericClientInterceptor(headers);
            ManagedChannel managedChannel = NettyChannelBuilder.forAddress(host, port)
                    .sslContext(sslContext)
                    .intercept(interceptor)
                    .build();
            StudentServiceGrpc.StudentServiceBlockingStub blockingStub = StudentServiceGrpc.newBlockingStub(managedChannel);

            return blockingStub.getStudent(StudentRequest.newBuilder().setId(1).build());
        } catch (StatusRuntimeException e) {
            log.error("Failed to retrieve student: " + e.getStatus().getDescription());
            throw e;
        } catch (SSLException e) {
            throw new RuntimeException(e);
        }
    }


}

