package org.wso2.apk.enforcer.discovery.service.subscription;

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
 * <pre>
 * [#protodoc-title: AppMappingDS]
 * </pre>
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler",
    comments = "Source: wso2/discovery/service/subscription/app_mapping_ds.proto")
public final class ApplicationMappingDiscoveryServiceGrpc {

  private ApplicationMappingDiscoveryServiceGrpc() {}

  public static final String SERVICE_NAME = "discovery.service.subscription.ApplicationMappingDiscoveryService";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<io.envoyproxy.envoy.service.discovery.v3.DiscoveryRequest,
      io.envoyproxy.envoy.service.discovery.v3.DiscoveryResponse> getStreamApplicationMappingsMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "StreamApplicationMappings",
      requestType = io.envoyproxy.envoy.service.discovery.v3.DiscoveryRequest.class,
      responseType = io.envoyproxy.envoy.service.discovery.v3.DiscoveryResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.BIDI_STREAMING)
  public static io.grpc.MethodDescriptor<io.envoyproxy.envoy.service.discovery.v3.DiscoveryRequest,
      io.envoyproxy.envoy.service.discovery.v3.DiscoveryResponse> getStreamApplicationMappingsMethod() {
    io.grpc.MethodDescriptor<io.envoyproxy.envoy.service.discovery.v3.DiscoveryRequest, io.envoyproxy.envoy.service.discovery.v3.DiscoveryResponse> getStreamApplicationMappingsMethod;
    if ((getStreamApplicationMappingsMethod = ApplicationMappingDiscoveryServiceGrpc.getStreamApplicationMappingsMethod) == null) {
      synchronized (ApplicationMappingDiscoveryServiceGrpc.class) {
        if ((getStreamApplicationMappingsMethod = ApplicationMappingDiscoveryServiceGrpc.getStreamApplicationMappingsMethod) == null) {
          ApplicationMappingDiscoveryServiceGrpc.getStreamApplicationMappingsMethod = getStreamApplicationMappingsMethod =
              io.grpc.MethodDescriptor.<io.envoyproxy.envoy.service.discovery.v3.DiscoveryRequest, io.envoyproxy.envoy.service.discovery.v3.DiscoveryResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.BIDI_STREAMING)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "StreamApplicationMappings"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  io.envoyproxy.envoy.service.discovery.v3.DiscoveryRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  io.envoyproxy.envoy.service.discovery.v3.DiscoveryResponse.getDefaultInstance()))
              .setSchemaDescriptor(new ApplicationMappingDiscoveryServiceMethodDescriptorSupplier("StreamApplicationMappings"))
              .build();
        }
      }
    }
    return getStreamApplicationMappingsMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static ApplicationMappingDiscoveryServiceStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<ApplicationMappingDiscoveryServiceStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<ApplicationMappingDiscoveryServiceStub>() {
        @java.lang.Override
        public ApplicationMappingDiscoveryServiceStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new ApplicationMappingDiscoveryServiceStub(channel, callOptions);
        }
      };
    return ApplicationMappingDiscoveryServiceStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static ApplicationMappingDiscoveryServiceBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<ApplicationMappingDiscoveryServiceBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<ApplicationMappingDiscoveryServiceBlockingStub>() {
        @java.lang.Override
        public ApplicationMappingDiscoveryServiceBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new ApplicationMappingDiscoveryServiceBlockingStub(channel, callOptions);
        }
      };
    return ApplicationMappingDiscoveryServiceBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static ApplicationMappingDiscoveryServiceFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<ApplicationMappingDiscoveryServiceFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<ApplicationMappingDiscoveryServiceFutureStub>() {
        @java.lang.Override
        public ApplicationMappingDiscoveryServiceFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new ApplicationMappingDiscoveryServiceFutureStub(channel, callOptions);
        }
      };
    return ApplicationMappingDiscoveryServiceFutureStub.newStub(factory, channel);
  }

  /**
   * <pre>
   * [#protodoc-title: AppMappingDS]
   * </pre>
   */
  public static abstract class ApplicationMappingDiscoveryServiceImplBase implements io.grpc.BindableService {

    /**
     */
    public io.grpc.stub.StreamObserver<io.envoyproxy.envoy.service.discovery.v3.DiscoveryRequest> streamApplicationMappings(
        io.grpc.stub.StreamObserver<io.envoyproxy.envoy.service.discovery.v3.DiscoveryResponse> responseObserver) {
      return asyncUnimplementedStreamingCall(getStreamApplicationMappingsMethod(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getStreamApplicationMappingsMethod(),
            asyncBidiStreamingCall(
              new MethodHandlers<
                io.envoyproxy.envoy.service.discovery.v3.DiscoveryRequest,
                io.envoyproxy.envoy.service.discovery.v3.DiscoveryResponse>(
                  this, METHODID_STREAM_APPLICATION_MAPPINGS)))
          .build();
    }
  }

  /**
   * <pre>
   * [#protodoc-title: AppMappingDS]
   * </pre>
   */
  public static final class ApplicationMappingDiscoveryServiceStub extends io.grpc.stub.AbstractAsyncStub<ApplicationMappingDiscoveryServiceStub> {
    private ApplicationMappingDiscoveryServiceStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected ApplicationMappingDiscoveryServiceStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new ApplicationMappingDiscoveryServiceStub(channel, callOptions);
    }

    /**
     */
    public io.grpc.stub.StreamObserver<io.envoyproxy.envoy.service.discovery.v3.DiscoveryRequest> streamApplicationMappings(
        io.grpc.stub.StreamObserver<io.envoyproxy.envoy.service.discovery.v3.DiscoveryResponse> responseObserver) {
      return asyncBidiStreamingCall(
          getChannel().newCall(getStreamApplicationMappingsMethod(), getCallOptions()), responseObserver);
    }
  }

  /**
   * <pre>
   * [#protodoc-title: AppMappingDS]
   * </pre>
   */
  public static final class ApplicationMappingDiscoveryServiceBlockingStub extends io.grpc.stub.AbstractBlockingStub<ApplicationMappingDiscoveryServiceBlockingStub> {
    private ApplicationMappingDiscoveryServiceBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected ApplicationMappingDiscoveryServiceBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new ApplicationMappingDiscoveryServiceBlockingStub(channel, callOptions);
    }
  }

  /**
   * <pre>
   * [#protodoc-title: AppMappingDS]
   * </pre>
   */
  public static final class ApplicationMappingDiscoveryServiceFutureStub extends io.grpc.stub.AbstractFutureStub<ApplicationMappingDiscoveryServiceFutureStub> {
    private ApplicationMappingDiscoveryServiceFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected ApplicationMappingDiscoveryServiceFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new ApplicationMappingDiscoveryServiceFutureStub(channel, callOptions);
    }
  }

  private static final int METHODID_STREAM_APPLICATION_MAPPINGS = 0;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final ApplicationMappingDiscoveryServiceImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(ApplicationMappingDiscoveryServiceImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        default:
          throw new AssertionError();
      }
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public io.grpc.stub.StreamObserver<Req> invoke(
        io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_STREAM_APPLICATION_MAPPINGS:
          return (io.grpc.stub.StreamObserver<Req>) serviceImpl.streamApplicationMappings(
              (io.grpc.stub.StreamObserver<io.envoyproxy.envoy.service.discovery.v3.DiscoveryResponse>) responseObserver);
        default:
          throw new AssertionError();
      }
    }
  }

  private static abstract class ApplicationMappingDiscoveryServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    ApplicationMappingDiscoveryServiceBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return org.wso2.apk.enforcer.discovery.service.subscription.AppMappingDSProto.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("ApplicationMappingDiscoveryService");
    }
  }

  private static final class ApplicationMappingDiscoveryServiceFileDescriptorSupplier
      extends ApplicationMappingDiscoveryServiceBaseDescriptorSupplier {
    ApplicationMappingDiscoveryServiceFileDescriptorSupplier() {}
  }

  private static final class ApplicationMappingDiscoveryServiceMethodDescriptorSupplier
      extends ApplicationMappingDiscoveryServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    ApplicationMappingDiscoveryServiceMethodDescriptorSupplier(String methodName) {
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
      synchronized (ApplicationMappingDiscoveryServiceGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new ApplicationMappingDiscoveryServiceFileDescriptorSupplier())
              .addMethod(getStreamApplicationMappingsMethod())
              .build();
        }
      }
    }
    return result;
  }
}
