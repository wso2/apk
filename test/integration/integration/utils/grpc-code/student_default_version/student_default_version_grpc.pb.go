// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.0
// source: student_default_version.proto

package student_default_version

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	StudentServiceDefaultVersion_GetStudent_FullMethodName              = "/org.apk.student_service_default_version.StudentServiceDefaultVersion/GetStudent"
	StudentServiceDefaultVersion_GetStudentStream_FullMethodName        = "/org.apk.student_service_default_version.StudentServiceDefaultVersion/GetStudentStream"
	StudentServiceDefaultVersion_SendStudentStream_FullMethodName       = "/org.apk.student_service_default_version.StudentServiceDefaultVersion/SendStudentStream"
	StudentServiceDefaultVersion_SendAndGetStudentStream_FullMethodName = "/org.apk.student_service_default_version.StudentServiceDefaultVersion/SendAndGetStudentStream"
)

// StudentServiceDefaultVersionClient is the client API for StudentServiceDefaultVersion service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StudentServiceDefaultVersionClient interface {
	GetStudent(ctx context.Context, in *StudentRequest, opts ...grpc.CallOption) (*StudentResponse, error)
	GetStudentStream(ctx context.Context, in *StudentRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[StudentResponse], error)
	SendStudentStream(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[StudentRequest, StudentResponse], error)
	SendAndGetStudentStream(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[StudentRequest, StudentResponse], error)
}

type studentServiceDefaultVersionClient struct {
	cc grpc.ClientConnInterface
}

func NewStudentServiceDefaultVersionClient(cc grpc.ClientConnInterface) StudentServiceDefaultVersionClient {
	return &studentServiceDefaultVersionClient{cc}
}

func (c *studentServiceDefaultVersionClient) GetStudent(ctx context.Context, in *StudentRequest, opts ...grpc.CallOption) (*StudentResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(StudentResponse)
	err := c.cc.Invoke(ctx, StudentServiceDefaultVersion_GetStudent_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *studentServiceDefaultVersionClient) GetStudentStream(ctx context.Context, in *StudentRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[StudentResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &StudentServiceDefaultVersion_ServiceDesc.Streams[0], StudentServiceDefaultVersion_GetStudentStream_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[StudentRequest, StudentResponse]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type StudentServiceDefaultVersion_GetStudentStreamClient = grpc.ServerStreamingClient[StudentResponse]

func (c *studentServiceDefaultVersionClient) SendStudentStream(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[StudentRequest, StudentResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &StudentServiceDefaultVersion_ServiceDesc.Streams[1], StudentServiceDefaultVersion_SendStudentStream_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[StudentRequest, StudentResponse]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type StudentServiceDefaultVersion_SendStudentStreamClient = grpc.ClientStreamingClient[StudentRequest, StudentResponse]

func (c *studentServiceDefaultVersionClient) SendAndGetStudentStream(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[StudentRequest, StudentResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &StudentServiceDefaultVersion_ServiceDesc.Streams[2], StudentServiceDefaultVersion_SendAndGetStudentStream_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[StudentRequest, StudentResponse]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type StudentServiceDefaultVersion_SendAndGetStudentStreamClient = grpc.BidiStreamingClient[StudentRequest, StudentResponse]

// StudentServiceDefaultVersionServer is the server API for StudentServiceDefaultVersion service.
// All implementations must embed UnimplementedStudentServiceDefaultVersionServer
// for forward compatibility.
type StudentServiceDefaultVersionServer interface {
	GetStudent(context.Context, *StudentRequest) (*StudentResponse, error)
	GetStudentStream(*StudentRequest, grpc.ServerStreamingServer[StudentResponse]) error
	SendStudentStream(grpc.ClientStreamingServer[StudentRequest, StudentResponse]) error
	SendAndGetStudentStream(grpc.BidiStreamingServer[StudentRequest, StudentResponse]) error
	mustEmbedUnimplementedStudentServiceDefaultVersionServer()
}

// UnimplementedStudentServiceDefaultVersionServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedStudentServiceDefaultVersionServer struct{}

func (UnimplementedStudentServiceDefaultVersionServer) GetStudent(context.Context, *StudentRequest) (*StudentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStudent not implemented")
}
func (UnimplementedStudentServiceDefaultVersionServer) GetStudentStream(*StudentRequest, grpc.ServerStreamingServer[StudentResponse]) error {
	return status.Errorf(codes.Unimplemented, "method GetStudentStream not implemented")
}
func (UnimplementedStudentServiceDefaultVersionServer) SendStudentStream(grpc.ClientStreamingServer[StudentRequest, StudentResponse]) error {
	return status.Errorf(codes.Unimplemented, "method SendStudentStream not implemented")
}
func (UnimplementedStudentServiceDefaultVersionServer) SendAndGetStudentStream(grpc.BidiStreamingServer[StudentRequest, StudentResponse]) error {
	return status.Errorf(codes.Unimplemented, "method SendAndGetStudentStream not implemented")
}
func (UnimplementedStudentServiceDefaultVersionServer) mustEmbedUnimplementedStudentServiceDefaultVersionServer() {
}
func (UnimplementedStudentServiceDefaultVersionServer) testEmbeddedByValue() {}

// UnsafeStudentServiceDefaultVersionServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StudentServiceDefaultVersionServer will
// result in compilation errors.
type UnsafeStudentServiceDefaultVersionServer interface {
	mustEmbedUnimplementedStudentServiceDefaultVersionServer()
}

func RegisterStudentServiceDefaultVersionServer(s grpc.ServiceRegistrar, srv StudentServiceDefaultVersionServer) {
	// If the following call pancis, it indicates UnimplementedStudentServiceDefaultVersionServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&StudentServiceDefaultVersion_ServiceDesc, srv)
}

func _StudentServiceDefaultVersion_GetStudent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StudentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StudentServiceDefaultVersionServer).GetStudent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StudentServiceDefaultVersion_GetStudent_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StudentServiceDefaultVersionServer).GetStudent(ctx, req.(*StudentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StudentServiceDefaultVersion_GetStudentStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(StudentRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(StudentServiceDefaultVersionServer).GetStudentStream(m, &grpc.GenericServerStream[StudentRequest, StudentResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type StudentServiceDefaultVersion_GetStudentStreamServer = grpc.ServerStreamingServer[StudentResponse]

func _StudentServiceDefaultVersion_SendStudentStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(StudentServiceDefaultVersionServer).SendStudentStream(&grpc.GenericServerStream[StudentRequest, StudentResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type StudentServiceDefaultVersion_SendStudentStreamServer = grpc.ClientStreamingServer[StudentRequest, StudentResponse]

func _StudentServiceDefaultVersion_SendAndGetStudentStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(StudentServiceDefaultVersionServer).SendAndGetStudentStream(&grpc.GenericServerStream[StudentRequest, StudentResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type StudentServiceDefaultVersion_SendAndGetStudentStreamServer = grpc.BidiStreamingServer[StudentRequest, StudentResponse]

// StudentServiceDefaultVersion_ServiceDesc is the grpc.ServiceDesc for StudentServiceDefaultVersion service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var StudentServiceDefaultVersion_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "org.apk.student_service_default_version.StudentServiceDefaultVersion",
	HandlerType: (*StudentServiceDefaultVersionServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetStudent",
			Handler:    _StudentServiceDefaultVersion_GetStudent_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetStudentStream",
			Handler:       _StudentServiceDefaultVersion_GetStudentStream_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "SendStudentStream",
			Handler:       _StudentServiceDefaultVersion_SendStudentStream_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "SendAndGetStudentStream",
			Handler:       _StudentServiceDefaultVersion_SendAndGetStudentStream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "student_default_version.proto",
}
