package org.wso2.apk.integration.utils.clients;

import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import org.wso2.apk.integration.utils.clients.studentGrpcClient.StudentRequest;
import org.wso2.apk.integration.utils.clients.studentGrpcClient.StudentResponse;
import org.wso2.apk.integration.utils.clients.studentGrpcClient.StudentServiceGrpc;

import java.security.KeyManagementException;
import java.security.KeyStoreException;
import java.security.NoSuchAlgorithmException;

public class SimpleGRPCStudentClient {
    protected Log log = LogFactory.getLog(SimpleGRPCStudentClient.class);
    private static final int EVENTUAL_SUCCESS_RESPONSE_TIMEOUT_IN_SECONDS = 10;

    private ManagedChannel managedChannel;

    public SimpleGRPCStudentClient() throws NoSuchAlgorithmException, KeyStoreException, KeyManagementException {
        this.managedChannel = ManagedChannelBuilder.forAddress("default.gw.wso2.com", 9095).usePlaintext().build();
        log.info("ManagedChannel created");
    }

    public StudentResponse GetStudent() {
        //Create a blocking stub for the StudentService
        StudentServiceGrpc.StudentServiceBlockingStub blockingStub = StudentServiceGrpc.newBlockingStub(managedChannel);

        // Make a synchronous gRPC call to get student details for ID 1
        StudentResponse studentResponse = blockingStub.getStudent(StudentRequest.newBuilder().setId(1).build());

        // Log the response received from the gRPC server
        log.info("response = " + studentResponse.getName() + " " + studentResponse.getAge());
        return studentResponse;
    }


}

