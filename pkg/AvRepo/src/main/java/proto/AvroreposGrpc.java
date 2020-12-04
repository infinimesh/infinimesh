package main.java.proto;

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
    value = "by gRPC proto compiler (version 1.33.1)",
    comments = "Source: avrorepo.proto")
public final class AvroreposGrpc {

  private AvroreposGrpc() {}

  public static final String SERVICE_NAME = "proto.Avrorepos";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<main.java.proto.SaveDeviceStateRequest,
      main.java.proto.SaveDeviceStateResponse> getSetDeviceStateMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "SetDeviceState",
      requestType = main.java.proto.SaveDeviceStateRequest.class,
      responseType = main.java.proto.SaveDeviceStateResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<main.java.proto.SaveDeviceStateRequest,
      main.java.proto.SaveDeviceStateResponse> getSetDeviceStateMethod() {
    io.grpc.MethodDescriptor<main.java.proto.SaveDeviceStateRequest, main.java.proto.SaveDeviceStateResponse> getSetDeviceStateMethod;
    if ((getSetDeviceStateMethod = AvroreposGrpc.getSetDeviceStateMethod) == null) {
      synchronized (AvroreposGrpc.class) {
        if ((getSetDeviceStateMethod = AvroreposGrpc.getSetDeviceStateMethod) == null) {
          AvroreposGrpc.getSetDeviceStateMethod = getSetDeviceStateMethod =
              io.grpc.MethodDescriptor.<main.java.proto.SaveDeviceStateRequest, main.java.proto.SaveDeviceStateResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "SetDeviceState"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  main.java.proto.SaveDeviceStateRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  main.java.proto.SaveDeviceStateResponse.getDefaultInstance()))
              .setSchemaDescriptor(new AvroreposMethodDescriptorSupplier("SetDeviceState"))
              .build();
        }
      }
    }
    return getSetDeviceStateMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static AvroreposStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<AvroreposStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<AvroreposStub>() {
        @java.lang.Override
        public AvroreposStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new AvroreposStub(channel, callOptions);
        }
      };
    return AvroreposStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static AvroreposBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<AvroreposBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<AvroreposBlockingStub>() {
        @java.lang.Override
        public AvroreposBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new AvroreposBlockingStub(channel, callOptions);
        }
      };
    return AvroreposBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static AvroreposFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<AvroreposFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<AvroreposFutureStub>() {
        @java.lang.Override
        public AvroreposFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new AvroreposFutureStub(channel, callOptions);
        }
      };
    return AvroreposFutureStub.newStub(factory, channel);
  }

  /**
   */
  public static abstract class AvroreposImplBase implements io.grpc.BindableService {

    /**
     */
    public void setDeviceState(main.java.proto.SaveDeviceStateRequest request,
        io.grpc.stub.StreamObserver<main.java.proto.SaveDeviceStateResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getSetDeviceStateMethod(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getSetDeviceStateMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                main.java.proto.SaveDeviceStateRequest,
                main.java.proto.SaveDeviceStateResponse>(
                  this, METHODID_SET_DEVICE_STATE)))
          .build();
    }
  }

  /**
   */
  public static final class AvroreposStub extends io.grpc.stub.AbstractAsyncStub<AvroreposStub> {
    private AvroreposStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected AvroreposStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new AvroreposStub(channel, callOptions);
    }

    /**
     */
    public void setDeviceState(main.java.proto.SaveDeviceStateRequest request,
        io.grpc.stub.StreamObserver<main.java.proto.SaveDeviceStateResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getSetDeviceStateMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   */
  public static final class AvroreposBlockingStub extends io.grpc.stub.AbstractBlockingStub<AvroreposBlockingStub> {
    private AvroreposBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected AvroreposBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new AvroreposBlockingStub(channel, callOptions);
    }

    /**
     */
    public main.java.proto.SaveDeviceStateResponse setDeviceState(main.java.proto.SaveDeviceStateRequest request) {
      return blockingUnaryCall(
          getChannel(), getSetDeviceStateMethod(), getCallOptions(), request);
    }
  }

  /**
   */
  public static final class AvroreposFutureStub extends io.grpc.stub.AbstractFutureStub<AvroreposFutureStub> {
    private AvroreposFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected AvroreposFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new AvroreposFutureStub(channel, callOptions);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<main.java.proto.SaveDeviceStateResponse> setDeviceState(
        main.java.proto.SaveDeviceStateRequest request) {
      return futureUnaryCall(
          getChannel().newCall(getSetDeviceStateMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_SET_DEVICE_STATE = 0;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final AvroreposImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(AvroreposImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_SET_DEVICE_STATE:
          serviceImpl.setDeviceState((main.java.proto.SaveDeviceStateRequest) request,
              (io.grpc.stub.StreamObserver<main.java.proto.SaveDeviceStateResponse>) responseObserver);
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

  private static abstract class AvroreposBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    AvroreposBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return main.java.proto.avrorepo.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("Avrorepos");
    }
  }

  private static final class AvroreposFileDescriptorSupplier
      extends AvroreposBaseDescriptorSupplier {
    AvroreposFileDescriptorSupplier() {}
  }

  private static final class AvroreposMethodDescriptorSupplier
      extends AvroreposBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    AvroreposMethodDescriptorSupplier(String methodName) {
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
      synchronized (AvroreposGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new AvroreposFileDescriptorSupplier())
              .addMethod(getSetDeviceStateMethod())
              .build();
        }
      }
    }
    return result;
  }
}
