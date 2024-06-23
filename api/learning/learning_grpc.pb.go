// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: api/learning/learning.proto

package learning

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// LearningServiceClient is the client API for LearningService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LearningServiceClient interface {
	CreateClassRooms(ctx context.Context, in *CreateClassRoomsRequest, opts ...grpc.CallOption) (*OperationStatus, error)
	DeleteClassRoomsByTeacher(ctx context.Context, in *IDRequest, opts ...grpc.CallOption) (*OperationStatus, error)
	InitClassRooms(ctx context.Context, in *IDRequest, opts ...grpc.CallOption) (*Notifications, error)
}

type learningServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewLearningServiceClient(cc grpc.ClientConnInterface) LearningServiceClient {
	return &learningServiceClient{cc}
}

func (c *learningServiceClient) CreateClassRooms(ctx context.Context, in *CreateClassRoomsRequest, opts ...grpc.CallOption) (*OperationStatus, error) {
	out := new(OperationStatus)
	err := c.cc.Invoke(ctx, "/learning.LearningService/CreateClassRooms", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *learningServiceClient) DeleteClassRoomsByTeacher(ctx context.Context, in *IDRequest, opts ...grpc.CallOption) (*OperationStatus, error) {
	out := new(OperationStatus)
	err := c.cc.Invoke(ctx, "/learning.LearningService/DeleteClassRoomsByTeacher", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *learningServiceClient) InitClassRooms(ctx context.Context, in *IDRequest, opts ...grpc.CallOption) (*Notifications, error) {
	out := new(Notifications)
	err := c.cc.Invoke(ctx, "/learning.LearningService/InitClassRooms", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LearningServiceServer is the server API for LearningService service.
// All implementations must embed UnimplementedLearningServiceServer
// for forward compatibility
type LearningServiceServer interface {
	CreateClassRooms(context.Context, *CreateClassRoomsRequest) (*OperationStatus, error)
	DeleteClassRoomsByTeacher(context.Context, *IDRequest) (*OperationStatus, error)
	InitClassRooms(context.Context, *IDRequest) (*Notifications, error)
	mustEmbedUnimplementedLearningServiceServer()
}

// UnimplementedLearningServiceServer must be embedded to have forward compatible implementations.
type UnimplementedLearningServiceServer struct {
}

func (UnimplementedLearningServiceServer) CreateClassRooms(context.Context, *CreateClassRoomsRequest) (*OperationStatus, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateClassRooms not implemented")
}
func (UnimplementedLearningServiceServer) DeleteClassRoomsByTeacher(context.Context, *IDRequest) (*OperationStatus, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteClassRoomsByTeacher not implemented")
}
func (UnimplementedLearningServiceServer) InitClassRooms(context.Context, *IDRequest) (*Notifications, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InitClassRooms not implemented")
}
func (UnimplementedLearningServiceServer) mustEmbedUnimplementedLearningServiceServer() {}

// UnsafeLearningServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LearningServiceServer will
// result in compilation errors.
type UnsafeLearningServiceServer interface {
	mustEmbedUnimplementedLearningServiceServer()
}

func RegisterLearningServiceServer(s grpc.ServiceRegistrar, srv LearningServiceServer) {
	s.RegisterService(&LearningService_ServiceDesc, srv)
}

func _LearningService_CreateClassRooms_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateClassRoomsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LearningServiceServer).CreateClassRooms(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/learning.LearningService/CreateClassRooms",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LearningServiceServer).CreateClassRooms(ctx, req.(*CreateClassRoomsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LearningService_DeleteClassRoomsByTeacher_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LearningServiceServer).DeleteClassRoomsByTeacher(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/learning.LearningService/DeleteClassRoomsByTeacher",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LearningServiceServer).DeleteClassRoomsByTeacher(ctx, req.(*IDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LearningService_InitClassRooms_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LearningServiceServer).InitClassRooms(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/learning.LearningService/InitClassRooms",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LearningServiceServer).InitClassRooms(ctx, req.(*IDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// LearningService_ServiceDesc is the grpc.ServiceDesc for LearningService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LearningService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "learning.LearningService",
	HandlerType: (*LearningServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateClassRooms",
			Handler:    _LearningService_CreateClassRooms_Handler,
		},
		{
			MethodName: "DeleteClassRoomsByTeacher",
			Handler:    _LearningService_DeleteClassRoomsByTeacher_Handler,
		},
		{
			MethodName: "InitClassRooms",
			Handler:    _LearningService_InitClassRooms_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/learning/learning.proto",
}
