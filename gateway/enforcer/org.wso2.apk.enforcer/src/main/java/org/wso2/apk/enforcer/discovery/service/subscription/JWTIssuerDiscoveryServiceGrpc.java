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
 * [#protodoc-title: JWTIssuerDS]
 * </pre>
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler",
    comments = "Source: wso2/discovery/service/subscription/jwtds.proto")
public final class JWTIssuerDiscoveryServiceGrpc {

  private JWTIssuerDiscoveryServiceGrpc() {}

  public static final String SERVICE_NAME = "discovery.service.subscription.JWTIssuerDiscoveryService";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<io.envoyproxy.envoy.service.discovery.v3.DiscoveryRequest,
      io.envoyproxy.envoy.service.discovery.v3.DiscoveryResponse> getStreamJWTIssuersMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "StreamJWTIssuers",
      requestType = io.envoyproxy.envoy.service.discovery.v3.DiscoveryRequest.class,
      responseType = io.envoyproxy.envoy.service.discovery.v3.DiscoveryResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.BIDI_STREAMING)
  public static io.grpc.MethodDescriptor<io.envoyproxy.envoy.service.discovery.v3.DiscoveryRequest,
      io.envoyproxy.envoy.service.discovery.v3.DiscoveryResponse> getStreamJWTIssuersMethod() {
    io.grpc.MethodDescriptor<io.envoyproxy.envoy.service.discovery.v3.DiscoveryRequest, io.envoyproxy.envoy.service.discovery.v3.DiscoveryResponse> getStreamJWTIssuersMethod;
    if ((getStreamJWTIssuersMethod = JWTIssuerDiscoveryServiceGrpc.getStreamJWTIssuersMethod) == null) {
      synchronized (JWTIssuerDiscoveryServiceGrpc.class) {
        if ((getStreamJWTIssuersMethod = JWTIssuerDiscoveryServiceGrpc.getStreamJWTIssuersMethod) == null) {
          JWTIssuerDiscoveryServiceGrpc.getStreamJWTIssuersMethod = getStreamJWTIssuersMethod =
              io.grpc.MethodDescriptor.<io.envoyproxy.envoy.service.discovery.v3.DiscoveryRequest, io.envoyproxy.envoy.service.discovery.v3.DiscoveryResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.BIDI_STREAMING)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "StreamJWTIssuers"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  io.envoyproxy.envoy.service.discovery.v3.DiscoveryRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  io.envoyproxy.envoy.service.discovery.v3.DiscoveryResponse.getDefaultInstance()))
              .setSchemaDescriptor(new JWTIssuerDiscoveryServiceMethodDescriptorSupplier("StreamJWTIssuers"))
              .build();
        }
      }
    }
    return getStreamJWTIssuersMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static JWTIssuerDiscoveryServiceStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<JWTIssuerDiscoveryServiceStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<JWTIssuerDiscoveryServiceStub>() {
        @java.lang.Override
        public JWTIssuerDiscoveryServiceStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new JWTIssuerDiscoveryServiceStub(channel, callOptions);
        }
      };
    return JWTIssuerDiscoveryServiceStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static JWTIssuerDiscoveryServiceBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<JWTIssuerDiscoveryServiceBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<JWTIssuerDiscoveryServiceBlockingStub>() {
        @java.lang.Override
        public JWTIssuerDiscoveryServiceBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new JWTIssuerDiscoveryServiceBlockingStub(channel, callOptions);
        }
      };
    return JWTIssuerDiscoveryServiceBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static JWTIssuerDiscoveryServiceFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<JWTIssuerDiscoveryServiceFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<JWTIssuerDiscoveryServiceFutureStub>() {
        @java.lang.Override
        public JWTIssuerDiscoveryServiceFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new JWTIssuerDiscoveryServiceFutureStub(channel, callOptions);
        }
      };
    return JWTIssuerDiscoveryServiceFutureStub.newStub(factory, channel);
  }

  /**
   * <pre>
   * [#protodoc-title: JWTIssuerDS]
   * </pre>
   */
  public static abstract class JWTIssuerDiscoveryServiceImplBase implements io.grpc.BindableService {

    /**
     */
    public io.grpc.stub.StreamObserver<io.envoyproxy.envoy.service.discovery.v3.DiscoveryRequest> streamJWTIssuers(
        io.grpc.stub.StreamObserver<io.envoyproxy.envoy.service.discovery.v3.DiscoveryResponse> responseObserver) {
      return asyncUnimplementedStreamingCall(getStreamJWTIssuersMethod(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getStreamJWTIssuersMethod(),
            asyncBidiStreamingCall(
              new MethodHandlers<
                io.envoyproxy.envoy.service.discovery.v3.DiscoveryRequest,
                io.envoyproxy.envoy.service.discovery.v3.DiscoveryResponse>(
                  this, METHODID_STREAM_JWTISSUERS)))
          .build();
    }
  }

  /**
   * <pre>
   * [#protodoc-title: JWTIssuerDS]
   * </pre>
   */
  public static final class JWTIssuerDiscoveryServiceStub extends io.grpc.stub.AbstractAsyncStub<JWTIssuerDiscoveryServiceStub> {
    private JWTIssuerDiscoveryServiceStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected JWTIssuerDiscoveryServiceStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new JWTIssuerDiscoveryServiceStub(channel, callOptions);
    }

    /**
     */
    public io.grpc.stub.StreamObserver<io.envoyproxy.envoy.service.discovery.v3.DiscoveryRequest> streamJWTIssuers(
        io.grpc.stub.StreamObserver<io.envoyproxy.envoy.service.discovery.v3.DiscoveryResponse> responseObserver) {
      return asyncBidiStreamingCall(
          getChannel().newCall(getStreamJWTIssuersMethod(), getCallOptions()), responseObserver);
    }
  }

  /**
   * <pre>
   * [#protodoc-title: JWTIssuerDS]
   * </pre>
   */
  public static final class JWTIssuerDiscoveryServiceBlockingStub extends io.grpc.stub.AbstractBlockingStub<JWTIssuerDiscoveryServiceBlockingStub> {
    private JWTIssuerDiscoveryServiceBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected JWTIssuerDiscoveryServiceBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new JWTIssuerDiscoveryServiceBlockingStub(channel, callOptions);
    }
  }

  /**
   * <pre>
   * [#protodoc-title: JWTIssuerDS]
   * </pre>
   */
  public static final class JWTIssuerDiscoveryServiceFutureStub extends io.grpc.stub.AbstractFutureStub<JWTIssuerDiscoveryServiceFutureStub> {
    private JWTIssuerDiscoveryServiceFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected JWTIssuerDiscoveryServiceFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new JWTIssuerDiscoveryServiceFutureStub(channel, callOptions);
    }
  }

  private static final int METHODID_STREAM_JWTISSUERS = 0;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final JWTIssuerDiscoveryServiceImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(JWTIssuerDiscoveryServiceImplBase serviceImpl, int methodId) {
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
        case METHODID_STREAM_JWTISSUERS:
          return (io.grpc.stub.StreamObserver<Req>) serviceImpl.streamJWTIssuers(
              (io.grpc.stub.StreamObserver<io.envoyproxy.envoy.service.discovery.v3.DiscoveryResponse>) responseObserver);
        default:
          throw new AssertionError();
      }
    }
  }

  private static abstract class JWTIssuerDiscoveryServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    JWTIssuerDiscoveryServiceBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return org.wso2.apk.enforcer.discovery.service.subscription.JWTIssuerDSProto.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("JWTIssuerDiscoveryService");
    }
  }

  private static final class JWTIssuerDiscoveryServiceFileDescriptorSupplier
      extends JWTIssuerDiscoveryServiceBaseDescriptorSupplier {
    JWTIssuerDiscoveryServiceFileDescriptorSupplier() {}
  }

  private static final class JWTIssuerDiscoveryServiceMethodDescriptorSupplier
      extends JWTIssuerDiscoveryServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    JWTIssuerDiscoveryServiceMethodDescriptorSupplier(String methodName) {
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
      synchronized (JWTIssuerDiscoveryServiceGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new JWTIssuerDiscoveryServiceFileDescriptorSupplier())
              .addMethod(getStreamJWTIssuersMethod())
              .build();
        }
      }
    }
    return result;
  }
}
