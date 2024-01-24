// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.12.4
// source: publisher.proto

package publisher

import (
	context "context"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Publisher_PushTask_FullMethodName = "/Publisher/PushTask"
	Publisher_PullTask_FullMethodName = "/Publisher/PullTask"
)

// PublisherClient is the client API for Publisher service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PublisherClient interface {
	PushTask(ctx context.Context, in *TaskSpec, opts ...grpc.CallOption) (*NodeSpec, error)
	PullTask(ctx context.Context, in *TaskSpec, opts ...grpc.CallOption) (*empty.Empty, error)
}

type publisherClient struct {
	cc grpc.ClientConnInterface
}

func NewPublisherClient(cc grpc.ClientConnInterface) PublisherClient {
	return &publisherClient{cc}
}

func (c *publisherClient) PushTask(ctx context.Context, in *TaskSpec, opts ...grpc.CallOption) (*NodeSpec, error) {
	out := new(NodeSpec)
	err := c.cc.Invoke(ctx, Publisher_PushTask_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *publisherClient) PullTask(ctx context.Context, in *TaskSpec, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, Publisher_PullTask_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PublisherServer is the server API for Publisher service.
// All implementations must embed UnimplementedPublisherServer
// for forward compatibility
type PublisherServer interface {
	PushTask(context.Context, *TaskSpec) (*NodeSpec, error)
	PullTask(context.Context, *TaskSpec) (*empty.Empty, error)
	mustEmbedUnimplementedPublisherServer()
}

// UnimplementedPublisherServer must be embedded to have forward compatible implementations.
type UnimplementedPublisherServer struct {
}

func (UnimplementedPublisherServer) PushTask(context.Context, *TaskSpec) (*NodeSpec, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PushTask not implemented")
}
func (UnimplementedPublisherServer) PullTask(context.Context, *TaskSpec) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PullTask not implemented")
}
func (UnimplementedPublisherServer) mustEmbedUnimplementedPublisherServer() {}

// UnsafePublisherServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PublisherServer will
// result in compilation errors.
type UnsafePublisherServer interface {
	mustEmbedUnimplementedPublisherServer()
}

func RegisterPublisherServer(s grpc.ServiceRegistrar, srv PublisherServer) {
	s.RegisterService(&Publisher_ServiceDesc, srv)
}

func _Publisher_PushTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TaskSpec)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PublisherServer).PushTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Publisher_PushTask_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PublisherServer).PushTask(ctx, req.(*TaskSpec))
	}
	return interceptor(ctx, in, info, handler)
}

func _Publisher_PullTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TaskSpec)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PublisherServer).PullTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Publisher_PullTask_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PublisherServer).PullTask(ctx, req.(*TaskSpec))
	}
	return interceptor(ctx, in, info, handler)
}

// Publisher_ServiceDesc is the grpc.ServiceDesc for Publisher service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Publisher_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Publisher",
	HandlerType: (*PublisherServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PushTask",
			Handler:    _Publisher_PushTask_Handler,
		},
		{
			MethodName: "PullTask",
			Handler:    _Publisher_PullTask_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "publisher.proto",
}
