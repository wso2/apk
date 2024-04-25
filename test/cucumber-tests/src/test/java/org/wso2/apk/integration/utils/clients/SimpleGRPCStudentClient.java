package org.wso2.apk.integration.utils.clients;

import io.grpc.StatusRuntimeException;
import io.grpc.netty.GrpcSslContexts;
import io.grpc.netty.NettyChannelBuilder;
import io.netty.handler.ssl.SslContext;
import io.netty.handler.ssl.util.InsecureTrustManagerFactory;
import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import io.grpc.ManagedChannel;
import org.wso2.apk.integration.utils.JWTClientInterceptor;
import org.wso2.apk.integration.utils.clients.studentGrpcClient.StudentRequest;
import org.wso2.apk.integration.utils.clients.studentGrpcClient.StudentResponse;
import org.wso2.apk.integration.utils.clients.studentGrpcClient.StudentServiceGrpc;

import javax.net.ssl.SSLException;
import java.security.KeyManagementException;
import java.security.KeyStoreException;
import java.security.NoSuchAlgorithmException;

public class SimpleGRPCStudentClient {
    protected Log log = LogFactory.getLog(SimpleGRPCStudentClient.class);
    private static final int EVENTUAL_SUCCESS_RESPONSE_TIMEOUT_IN_SECONDS = 10;

    public SimpleGRPCStudentClient() throws NoSuchAlgorithmException, KeyStoreException, KeyManagementException {
        this.GetStudent();
    }

    public StudentResponse GetStudent() {
        try {
            //sleep for 5 seconds
            try {
                log.info("Sleeping for 5 seconds");
                Thread.sleep(5000);
                log.info("Woke up after 5 seconds");
            } catch (InterruptedException e) {
                log.error("Thread sleep interrupted");
            }
            SslContext sslContext = GrpcSslContexts.forClient()
                    .trustManager(InsecureTrustManagerFactory.INSTANCE)
                    .build();

            ManagedChannel managedChannel = NettyChannelBuilder.forAddress("default.gw.wso2.com", 9095)
                    .sslContext(sslContext)
                    .intercept(new JWTClientInterceptor("eyJhbGciOiJSUzI1NiIsICJ0eXAiOiJKV1QiLCAia2lkIjoiZ2F0ZXdheV9jZXJ0aWZpY2F0ZV9hbGlhcyJ9.eyJpc3MiOiJodHRwczovL2lkcC5hbS53c28yLmNvbS90b2tlbiIsICJzdWIiOiI0NWYxYzVjOC1hOTJlLTExZWQtYWZhMS0wMjQyYWMxMjAwMDIiLCAiYXVkIjoiYXVkMSIsICJleHAiOjE3MTQwMzc1NDMsICJuYmYiOjE3MTQwMzM5NDMsICJpYXQiOjE3MTQwMzM5NDMsICJqdGkiOiIwMWVmMDJkZS01NjhmLTE0NTgtYTJlMS0wYTk3N2NmMTA5MGMiLCAiY2xpZW50SWQiOiI0NWYxYzVjOC1hOTJlLTExZWQtYWZhMS0wMjQyYWMxMjAwMDIiLCAic2NvcGUiOiJhcGs6YXBpX2NyZWF0ZSJ9.NbvlpbNL0Q3Op7I36nWSMb9R6zCtUeI7las0vOYMNQxMgLrTOBLkXmd9EfSg46fqOD7a9YqoGmKgn5UhXcQSFhtUwAKvbYDvnTyYfT6X3fBqFWl59xt74yJ8f6cSRSBb88Is0qDCWoTVgM-5eTqb93uU8KC0LG1YwU3OoxoiuM1_ix1qbugb-X7gYVvfqttnHl_0e-4jgNN5YLQl8xo8DBs9D-yDWDkDQpj_NonCY1AXqlrynmKbf7kRPR3abHJiF07BQoXJzOXUv4lyHJ1K6DnHj9l2w-KNAWDTQ-kffkVtkQ3hNjbl0Q5ieHsVXLQ9HB1AkOGrho6W8CIO4qlb9A")) // replace "your-jwt-token" with your actual JWT token
                    .build();
            //Create a blocking stub for the StudentService
            StudentServiceGrpc.StudentServiceBlockingStub blockingStub = StudentServiceGrpc.newBlockingStub(managedChannel);
            // Make a synchronous gRPC call to get student details for ID 1
            StudentResponse studentResponse = blockingStub.getStudent(StudentRequest.newBuilder().setId(1).build());
            // Log the response received from the gRPC server
            log.info("response = " + studentResponse.getName() + " " + studentResponse.getAge());
            return studentResponse;
        } catch (StatusRuntimeException e) {
            log.error("Failed to retrieve student: " + e.getStatus().getDescription());
            throw e;
        } catch (SSLException e) {
            throw new RuntimeException(e);
        }
    }


}

