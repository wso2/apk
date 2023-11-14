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
 * <pre>
 * [#protodoc-title: EventStreamDS]
 * </pre>
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler",
    comments = "Source: wso2/discovery/service/apkmgt/eventds.proto")
public final class EventStreamServiceGrpc {

  private EventStreamServiceGrpc() {}

  public static final String SERVICE_NAME = "discovery.service.apkmgt.EventStreamService";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<org.wso2.apk.enforcer.discovery.service.apkmgt.Request,
      org.wso2.apk.enforcer.discovery.subscription.Event> getStreamEventsMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "StreamEvents",
      requestType = org.wso2.apk.enforcer.discovery.service.apkmgt.Request.class,
      responseType = org.wso2.apk.enforcer.discovery.subscription.Event.class,
      methodType = io.grpc.MethodDescriptor.MethodType.SERVER_STREAMING)
  public static io.grpc.MethodDescriptor<org.wso2.apk.enforcer.discovery.service.apkmgt.Request,
      org.wso2.apk.enforcer.discovery.subscription.Event> getStreamEventsMethod() {
    io.grpc.MethodDescriptor<org.wso2.apk.enforcer.discovery.service.apkmgt.Request, org.wso2.apk.enforcer.discovery.subscription.Event> getStreamEventsMethod;
    if ((getStreamEventsMethod = EventStreamServiceGrpc.getStreamEventsMethod) == null) {
      synchronized (EventStreamServiceGrpc.class) {
        if ((getStreamEventsMethod = EventStreamServiceGrpc.getStreamEventsMethod) == null) {
          EventStreamServiceGrpc.getStreamEventsMethod = getStreamEventsMethod =
              io.grpc.MethodDescriptor.<org.wso2.apk.enforcer.discovery.service.apkmgt.Request, org.wso2.apk.enforcer.discovery.subscription.Event>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.SERVER_STREAMING)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "StreamEvents"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.wso2.apk.enforcer.discovery.service.apkmgt.Request.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.wso2.apk.enforcer.discovery.subscription.Event.getDefaultInstance()))
              .setSchemaDescriptor(new EventStreamServiceMethodDescriptorSupplier("StreamEvents"))
              .build();
        }
      }
    }
    return getStreamEventsMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static EventStreamServiceStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<EventStreamServiceStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<EventStreamServiceStub>() {
        @java.lang.Override
        public EventStreamServiceStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new EventStreamServiceStub(channel, callOptions);
        }
      };
    return EventStreamServiceStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static EventStreamServiceBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<EventStreamServiceBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<EventStreamServiceBlockingStub>() {
        @java.lang.Override
        public EventStreamServiceBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new EventStreamServiceBlockingStub(channel, callOptions);
        }
      };
    return EventStreamServiceBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static EventStreamServiceFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<EventStreamServiceFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<EventStreamServiceFutureStub>() {
        @java.lang.Override
        public EventStreamServiceFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new EventStreamServiceFutureStub(channel, callOptions);
        }
      };
    return EventStreamServiceFutureStub.newStub(factory, channel);
  }

  /**
   * <pre>
   * [#protodoc-title: EventStreamDS]
   * </pre>
   */
  public static abstract class EventStreamServiceImplBase implements io.grpc.BindableService {

    /**
     */
    public void streamEvents(org.wso2.apk.enforcer.discovery.service.apkmgt.Request request,
        io.grpc.stub.StreamObserver<org.wso2.apk.enforcer.discovery.subscription.Event> responseObserver) {
      asyncUnimplementedUnaryCall(getStreamEventsMethod(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getStreamEventsMethod(),
            asyncServerStreamingCall(
              new MethodHandlers<
                org.wso2.apk.enforcer.discovery.service.apkmgt.Request,
                org.wso2.apk.enforcer.discovery.subscription.Event>(
                  this, METHODID_STREAM_EVENTS)))
          .build();
    }
  }

  /**
   * <pre>
   * [#protodoc-title: EventStreamDS]
   * </pre>
   */
  public static final class EventStreamServiceStub extends io.grpc.stub.AbstractAsyncStub<EventStreamServiceStub> {
    private EventStreamServiceStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected EventStreamServiceStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new EventStreamServiceStub(channel, callOptions);
    }

    /**
     */
    public void streamEvents(org.wso2.apk.enforcer.discovery.service.apkmgt.Request request,
        io.grpc.stub.StreamObserver<org.wso2.apk.enforcer.discovery.subscription.Event> responseObserver) {
      asyncServerStreamingCall(
          getChannel().newCall(getStreamEventsMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   * <pre>
   * [#protodoc-title: EventStreamDS]
   * </pre>
   */
  public static final class EventStreamServiceBlockingStub extends io.grpc.stub.AbstractBlockingStub<EventStreamServiceBlockingStub> {
    private EventStreamServiceBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected EventStreamServiceBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new EventStreamServiceBlockingStub(channel, callOptions);
    }

    /**
     */
    public java.util.Iterator<org.wso2.apk.enforcer.discovery.subscription.Event> streamEvents(
        org.wso2.apk.enforcer.discovery.service.apkmgt.Request request) {
      return blockingServerStreamingCall(
          getChannel(), getStreamEventsMethod(), getCallOptions(), request);
    }
  }

  /**
   * <pre>
   * [#protodoc-title: EventStreamDS]
   * </pre>
   */
  public static final class EventStreamServiceFutureStub extends io.grpc.stub.AbstractFutureStub<EventStreamServiceFutureStub> {
    private EventStreamServiceFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected EventStreamServiceFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new EventStreamServiceFutureStub(channel, callOptions);
    }
  }

  private static final int METHODID_STREAM_EVENTS = 0;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final EventStreamServiceImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(EventStreamServiceImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_STREAM_EVENTS:
          serviceImpl.streamEvents((org.wso2.apk.enforcer.discovery.service.apkmgt.Request) request,
              (io.grpc.stub.StreamObserver<org.wso2.apk.enforcer.discovery.subscription.Event>) responseObserver);
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

  private static abstract class EventStreamServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    EventStreamServiceBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return org.wso2.apk.enforcer.discovery.service.apkmgt.EventServiceProto.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("EventStreamService");
    }
  }

  private static final class EventStreamServiceFileDescriptorSupplier
      extends EventStreamServiceBaseDescriptorSupplier {
    EventStreamServiceFileDescriptorSupplier() {}
  }

  private static final class EventStreamServiceMethodDescriptorSupplier
      extends EventStreamServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    EventStreamServiceMethodDescriptorSupplier(String methodName) {
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
      synchronized (EventStreamServiceGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new EventStreamServiceFileDescriptorSupplier())
              .addMethod(getStreamEventsMethod())
              .build();
        }
      }
    }
    return result;
  }
}
