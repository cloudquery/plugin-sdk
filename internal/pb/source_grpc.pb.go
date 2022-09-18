// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.6
// source: internal/pb/source.proto

package pb

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

// SourceClient is the client API for Source service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SourceClient interface {
	// Get the name of the plugin
	GetName(ctx context.Context, in *GetName_Request, opts ...grpc.CallOption) (*GetName_Response, error)
	// Get an example configuration for the source plugin
	GetExampleConfig(ctx context.Context, in *GetExampleConfig_Request, opts ...grpc.CallOption) (*GetExampleConfig_Response, error)
	// Get all tables the source plugin supports
	GetTables(ctx context.Context, in *GetTables_Request, opts ...grpc.CallOption) (*GetTables_Response, error)
	// Fetch resources
	Sync(ctx context.Context, in *Sync_Request, opts ...grpc.CallOption) (Source_SyncClient, error)
}

type sourceClient struct {
	cc grpc.ClientConnInterface
}

func NewSourceClient(cc grpc.ClientConnInterface) SourceClient {
	return &sourceClient{cc}
}

func (c *sourceClient) GetName(ctx context.Context, in *GetName_Request, opts ...grpc.CallOption) (*GetName_Response, error) {
	out := new(GetName_Response)
	err := c.cc.Invoke(ctx, "/proto.Source/GetName", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sourceClient) GetExampleConfig(ctx context.Context, in *GetExampleConfig_Request, opts ...grpc.CallOption) (*GetExampleConfig_Response, error) {
	out := new(GetExampleConfig_Response)
	err := c.cc.Invoke(ctx, "/proto.Source/GetExampleConfig", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sourceClient) GetTables(ctx context.Context, in *GetTables_Request, opts ...grpc.CallOption) (*GetTables_Response, error) {
	out := new(GetTables_Response)
	err := c.cc.Invoke(ctx, "/proto.Source/GetTables", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sourceClient) Sync(ctx context.Context, in *Sync_Request, opts ...grpc.CallOption) (Source_SyncClient, error) {
	stream, err := c.cc.NewStream(ctx, &Source_ServiceDesc.Streams[0], "/proto.Source/Sync", opts...)
	if err != nil {
		return nil, err
	}
	x := &sourceSyncClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Source_SyncClient interface {
	Recv() (*Sync_Response, error)
	grpc.ClientStream
}

type sourceSyncClient struct {
	grpc.ClientStream
}

func (x *sourceSyncClient) Recv() (*Sync_Response, error) {
	m := new(Sync_Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// SourceServer is the server API for Source service.
// All implementations must embed UnimplementedSourceServer
// for forward compatibility
type SourceServer interface {
	// Get the name of the plugin
	GetName(context.Context, *GetName_Request) (*GetName_Response, error)
	// Get an example configuration for the source plugin
	GetExampleConfig(context.Context, *GetExampleConfig_Request) (*GetExampleConfig_Response, error)
	// Get all tables the source plugin supports
	GetTables(context.Context, *GetTables_Request) (*GetTables_Response, error)
	// Fetch resources
	Sync(*Sync_Request, Source_SyncServer) error
	mustEmbedUnimplementedSourceServer()
}

// UnimplementedSourceServer must be embedded to have forward compatible implementations.
type UnimplementedSourceServer struct {
}

func (UnimplementedSourceServer) GetName(context.Context, *GetName_Request) (*GetName_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetName not implemented")
}
func (UnimplementedSourceServer) GetExampleConfig(context.Context, *GetExampleConfig_Request) (*GetExampleConfig_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetExampleConfig not implemented")
}
func (UnimplementedSourceServer) GetTables(context.Context, *GetTables_Request) (*GetTables_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTables not implemented")
}
func (UnimplementedSourceServer) Sync(*Sync_Request, Source_SyncServer) error {
	return status.Errorf(codes.Unimplemented, "method Sync not implemented")
}
func (UnimplementedSourceServer) mustEmbedUnimplementedSourceServer() {}

// UnsafeSourceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SourceServer will
// result in compilation errors.
type UnsafeSourceServer interface {
	mustEmbedUnimplementedSourceServer()
}

func RegisterSourceServer(s grpc.ServiceRegistrar, srv SourceServer) {
	s.RegisterService(&Source_ServiceDesc, srv)
}

func _Source_GetName_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetName_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SourceServer).GetName(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Source/GetName",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SourceServer).GetName(ctx, req.(*GetName_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Source_GetExampleConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetExampleConfig_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SourceServer).GetExampleConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Source/GetExampleConfig",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SourceServer).GetExampleConfig(ctx, req.(*GetExampleConfig_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Source_GetTables_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTables_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SourceServer).GetTables(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Source/GetTables",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SourceServer).GetTables(ctx, req.(*GetTables_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Source_Sync_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Sync_Request)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SourceServer).Sync(m, &sourceSyncServer{stream})
}

type Source_SyncServer interface {
	Send(*Sync_Response) error
	grpc.ServerStream
}

type sourceSyncServer struct {
	grpc.ServerStream
}

func (x *sourceSyncServer) Send(m *Sync_Response) error {
	return x.ServerStream.SendMsg(m)
}

// Source_ServiceDesc is the grpc.ServiceDesc for Source service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Source_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Source",
	HandlerType: (*SourceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetName",
			Handler:    _Source_GetName_Handler,
		},
		{
			MethodName: "GetExampleConfig",
			Handler:    _Source_GetExampleConfig_Handler,
		},
		{
			MethodName: "GetTables",
			Handler:    _Source_GetTables_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Sync",
			Handler:       _Source_Sync_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "internal/pb/source.proto",
}
