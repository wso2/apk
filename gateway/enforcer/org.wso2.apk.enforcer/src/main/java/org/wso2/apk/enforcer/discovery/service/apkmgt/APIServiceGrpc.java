package org.wso2.apk.enforcer.discovery.service.apkmgt;

import static io.grpc.MethodDescriptor.generateFullMethodName;
import static io.grpc.stub.ClientCalls.asyncBidiStreamingCall;
import static io.grpc.stub.ClientCalls.asyncClientStreamingCall;
import static io.grpc.stub.ClientCalls.asyncServerStreamingCall;
import static io.grpc.stub.ClientCalls.asyncUnaryCall;
import static io.grpc.stub.ClientCalls.blockingServerStreamingCall;
import static io.grpc.stub.ClientCalls.blockingUnaryCall;
import static io.grpc.stub.ClientCalls.futureUnaryCall;
import static io.grpc.stub.ServerCalls.asyncBidiStreamingCall;
import static io.grpc.stub.ServerCalls.asyncClientStreamingCall;
import static io.grpc.stub.ServerCalls.asyncServerStreamingCall;
import static io.grpc.stub.ServerCalls.asyncUnaryCall;
import static io.grpc.stub.ServerCalls.asyncUnimplementedStreamingCall;
import static io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall;

/**
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler",
    comments = "Source: wso2/discovery/service/apkmgt/apids.proto")
public final class APIServiceGrpc {

  private APIServiceGrpc() {}

  public static final String SERVICE_NAME = "discovery.service.apkmgt.APIService";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.API,
      org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.Response> getCreateAPIMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "createAPI",
      requestType = org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.API.class,
      responseType = org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.Response.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.API,
      org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.Response> getCreateAPIMethod() {
    io.grpc.MethodDescriptor<org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.API, org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.Response> getCreateAPIMethod;
    if ((getCreateAPIMethod = APIServiceGrpc.getCreateAPIMethod) == null) {
      synchronized (APIServiceGrpc.class) {
        if ((getCreateAPIMethod = APIServiceGrpc.getCreateAPIMethod) == null) {
          APIServiceGrpc.getCreateAPIMethod = getCreateAPIMethod =
              io.grpc.MethodDescriptor.<org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.API, org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.Response>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "createAPI"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.API.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.Response.getDefaultInstance()))
              .setSchemaDescriptor(new APIServiceMethodDescriptorSupplier("createAPI"))
              .build();
        }
      }
    }
    return getCreateAPIMethod;
  }

  private static volatile io.grpc.MethodDescriptor<org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.API,
      org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.Response> getUpdateAPIMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "updateAPI",
      requestType = org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.API.class,
      responseType = org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.Response.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.API,
      org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.Response> getUpdateAPIMethod() {
    io.grpc.MethodDescriptor<org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.API, org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.Response> getUpdateAPIMethod;
    if ((getUpdateAPIMethod = APIServiceGrpc.getUpdateAPIMethod) == null) {
      synchronized (APIServiceGrpc.class) {
        if ((getUpdateAPIMethod = APIServiceGrpc.getUpdateAPIMethod) == null) {
          APIServiceGrpc.getUpdateAPIMethod = getUpdateAPIMethod =
              io.grpc.MethodDescriptor.<org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.API, org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.Response>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "updateAPI"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.API.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.Response.getDefaultInstance()))
              .setSchemaDescriptor(new APIServiceMethodDescriptorSupplier("updateAPI"))
              .build();
        }
      }
    }
    return getUpdateAPIMethod;
  }

  private static volatile io.grpc.MethodDescriptor<org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.API,
      org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.Response> getDeleteAPIMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "deleteAPI",
      requestType = org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.API.class,
      responseType = org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.Response.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.API,
      org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.Response> getDeleteAPIMethod() {
    io.grpc.MethodDescriptor<org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.API, org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.Response> getDeleteAPIMethod;
    if ((getDeleteAPIMethod = APIServiceGrpc.getDeleteAPIMethod) == null) {
      synchronized (APIServiceGrpc.class) {
        if ((getDeleteAPIMethod = APIServiceGrpc.getDeleteAPIMethod) == null) {
          APIServiceGrpc.getDeleteAPIMethod = getDeleteAPIMethod =
              io.grpc.MethodDescriptor.<org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.API, org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.Response>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "deleteAPI"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.API.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.Response.getDefaultInstance()))
              .setSchemaDescriptor(new APIServiceMethodDescriptorSupplier("deleteAPI"))
              .build();
        }
      }
    }
    return getDeleteAPIMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static APIServiceStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<APIServiceStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<APIServiceStub>() {
        @java.lang.Override
        public APIServiceStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new APIServiceStub(channel, callOptions);
        }
      };
    return APIServiceStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static APIServiceBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<APIServiceBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<APIServiceBlockingStub>() {
        @java.lang.Override
        public APIServiceBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new APIServiceBlockingStub(channel, callOptions);
        }
      };
    return APIServiceBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static APIServiceFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<APIServiceFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<APIServiceFutureStub>() {
        @java.lang.Override
        public APIServiceFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new APIServiceFutureStub(channel, callOptions);
        }
      };
    return APIServiceFutureStub.newStub(factory, channel);
  }

  /**
   */
  public static abstract class APIServiceImplBase implements io.grpc.BindableService {

    /**
     */
    public void createAPI(org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.API request,
        io.grpc.stub.StreamObserver<org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.Response> responseObserver) {
      asyncUnimplementedUnaryCall(getCreateAPIMethod(), responseObserver);
    }

    /**
     */
    public void updateAPI(org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.API request,
        io.grpc.stub.StreamObserver<org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.Response> responseObserver) {
      asyncUnimplementedUnaryCall(getUpdateAPIMethod(), responseObserver);
    }

    /**
     */
    public void deleteAPI(org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.API request,
        io.grpc.stub.StreamObserver<org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.Response> responseObserver) {
      asyncUnimplementedUnaryCall(getDeleteAPIMethod(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getCreateAPIMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.API,
                org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.Response>(
                  this, METHODID_CREATE_API)))
          .addMethod(
            getUpdateAPIMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.API,
                org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.Response>(
                  this, METHODID_UPDATE_API)))
          .addMethod(
            getDeleteAPIMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.API,
                org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.Response>(
                  this, METHODID_DELETE_API)))
          .build();
    }
  }

  /**
   */
  public static final class APIServiceStub extends io.grpc.stub.AbstractAsyncStub<APIServiceStub> {
    private APIServiceStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected APIServiceStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new APIServiceStub(channel, callOptions);
    }

    /**
     */
    public void createAPI(org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.API request,
        io.grpc.stub.StreamObserver<org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.Response> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getCreateAPIMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void updateAPI(org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.API request,
        io.grpc.stub.StreamObserver<org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.Response> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getUpdateAPIMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void deleteAPI(org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.API request,
        io.grpc.stub.StreamObserver<org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.Response> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getDeleteAPIMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   */
  public static final class APIServiceBlockingStub extends io.grpc.stub.AbstractBlockingStub<APIServiceBlockingStub> {
    private APIServiceBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected APIServiceBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new APIServiceBlockingStub(channel, callOptions);
    }

    /**
     */
    public org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.Response createAPI(org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.API request) {
      return blockingUnaryCall(
          getChannel(), getCreateAPIMethod(), getCallOptions(), request);
    }

    /**
     */
    public org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.Response updateAPI(org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.API request) {
      return blockingUnaryCall(
          getChannel(), getUpdateAPIMethod(), getCallOptions(), request);
    }

    /**
     */
    public org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.Response deleteAPI(org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.API request) {
      return blockingUnaryCall(
          getChannel(), getDeleteAPIMethod(), getCallOptions(), request);
    }
  }

  /**
   */
  public static final class APIServiceFutureStub extends io.grpc.stub.AbstractFutureStub<APIServiceFutureStub> {
    private APIServiceFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected APIServiceFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new APIServiceFutureStub(channel, callOptions);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.Response> createAPI(
        org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.API request) {
      return futureUnaryCall(
          getChannel().newCall(getCreateAPIMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.Response> updateAPI(
        org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.API request) {
      return futureUnaryCall(
          getChannel().newCall(getUpdateAPIMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.Response> deleteAPI(
        org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.API request) {
      return futureUnaryCall(
          getChannel().newCall(getDeleteAPIMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_CREATE_API = 0;
  private static final int METHODID_UPDATE_API = 1;
  private static final int METHODID_DELETE_API = 2;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final APIServiceImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(APIServiceImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_CREATE_API:
          serviceImpl.createAPI((org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.API) request,
              (io.grpc.stub.StreamObserver<org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.Response>) responseObserver);
          break;
        case METHODID_UPDATE_API:
          serviceImpl.updateAPI((org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.API) request,
              (io.grpc.stub.StreamObserver<org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.Response>) responseObserver);
          break;
        case METHODID_DELETE_API:
          serviceImpl.deleteAPI((org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.API) request,
              (io.grpc.stub.StreamObserver<org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.Response>) responseObserver);
          break;
        default:
          throw new AssertionError();
      }
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public io.grpc.stub.StreamObserver<Req> invoke(
        io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        default:
          throw new AssertionError();
      }
    }
  }

  private static abstract class APIServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    APIServiceBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return org.wso2.apk.enforcer.discovery.service.apkmgt.ApiDsProto.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("APIService");
    }
  }

  private static final class APIServiceFileDescriptorSupplier
      extends APIServiceBaseDescriptorSupplier {
    APIServiceFileDescriptorSupplier() {}
  }

  private static final class APIServiceMethodDescriptorSupplier
      extends APIServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    APIServiceMethodDescriptorSupplier(String methodName) {
      this.methodName = methodName;
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.MethodDescriptor getMethodDescriptor() {
      return getServiceDescriptor().findMethodByName(methodName);
    }
  }

  private static volatile io.grpc.ServiceDescriptor serviceDescriptor;

  public static io.grpc.ServiceDescriptor getServiceDescriptor() {
    io.grpc.ServiceDescriptor result = serviceDescriptor;
    if (result == null) {
      synchronized (APIServiceGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new APIServiceFileDescriptorSupplier())
              .addMethod(getCreateAPIMethod())
              .addMethod(getUpdateAPIMethod())
              .addMethod(getDeleteAPIMethod())
              .build();
        }
      }
    }
    return result;
  }
}
