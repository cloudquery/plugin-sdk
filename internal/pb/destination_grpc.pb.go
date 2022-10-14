// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.6
// source: internal/pb/destination.proto

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

// DestinationClient is the client API for Destination service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DestinationClient interface {
	// Get the current protocol version of the plugin. This helps
	// get the right message about upgrade/downgrade of cli and/or plugin.
	// Also, on the cli side it can try to upgrade/downgrade the protocol if cli supports it.
	GetProtocolVersion(ctx context.Context, in *GetProtocolVersion_Request, opts ...grpc.CallOption) (*GetProtocolVersion_Response, error)
	// Get the name of the plugin
	GetName(ctx context.Context, in *GetName_Request, opts ...grpc.CallOption) (*GetName_Response, error)
	// Get the current version of the plugin
	GetVersion(ctx context.Context, in *GetVersion_Request, opts ...grpc.CallOption) (*GetVersion_Response, error)
	// Configure the plugin with the given credentials and mode
	Configure(ctx context.Context, in *Configure_Request, opts ...grpc.CallOption) (*Configure_Response, error)
	// Migrate tables to the given plugin version
	Migrate(ctx context.Context, in *Migrate_Request, opts ...grpc.CallOption) (*Migrate_Response, error)
	// Write resources
	Write(ctx context.Context, opts ...grpc.CallOption) (Destination_WriteClient, error)
	// Send signal to flush and close open connections
	Close(ctx context.Context, in *Close_Request, opts ...grpc.CallOption) (*Close_Response, error)
	// DeleteStale deletes stale data that was inserted by a given source
	// and is older than the given timestamp
	DeleteStale(ctx context.Context, in *DeleteStale_Request, opts ...grpc.CallOption) (*DeleteStale_Response, error)
}

type destinationClient struct {
	cc grpc.ClientConnInterface
}

func NewDestinationClient(cc grpc.ClientConnInterface) DestinationClient {
	return &destinationClient{cc}
}

func (c *destinationClient) GetProtocolVersion(ctx context.Context, in *GetProtocolVersion_Request, opts ...grpc.CallOption) (*GetProtocolVersion_Response, error) {
	out := new(GetProtocolVersion_Response)
	err := c.cc.Invoke(ctx, "/proto.Destination/GetProtocolVersion", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *destinationClient) GetName(ctx context.Context, in *GetName_Request, opts ...grpc.CallOption) (*GetName_Response, error) {
	out := new(GetName_Response)
	err := c.cc.Invoke(ctx, "/proto.Destination/GetName", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *destinationClient) GetVersion(ctx context.Context, in *GetVersion_Request, opts ...grpc.CallOption) (*GetVersion_Response, error) {
	out := new(GetVersion_Response)
	err := c.cc.Invoke(ctx, "/proto.Destination/GetVersion", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *destinationClient) Configure(ctx context.Context, in *Configure_Request, opts ...grpc.CallOption) (*Configure_Response, error) {
	out := new(Configure_Response)
	err := c.cc.Invoke(ctx, "/proto.Destination/Configure", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *destinationClient) Migrate(ctx context.Context, in *Migrate_Request, opts ...grpc.CallOption) (*Migrate_Response, error) {
	out := new(Migrate_Response)
	err := c.cc.Invoke(ctx, "/proto.Destination/Migrate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *destinationClient) Write(ctx context.Context, opts ...grpc.CallOption) (Destination_WriteClient, error) {
	stream, err := c.cc.NewStream(ctx, &Destination_ServiceDesc.Streams[0], "/proto.Destination/Write", opts...)
	if err != nil {
		return nil, err
	}
	x := &destinationWriteClient{stream}
	return x, nil
}

type Destination_WriteClient interface {
	Send(*Write_Request) error
	CloseAndRecv() (*Write_Response, error)
	grpc.ClientStream
}

type destinationWriteClient struct {
	grpc.ClientStream
}

func (x *destinationWriteClient) Send(m *Write_Request) error {
	return x.ClientStream.SendMsg(m)
}

func (x *destinationWriteClient) CloseAndRecv() (*Write_Response, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(Write_Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *destinationClient) Close(ctx context.Context, in *Close_Request, opts ...grpc.CallOption) (*Close_Response, error) {
	out := new(Close_Response)
	err := c.cc.Invoke(ctx, "/proto.Destination/Close", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *destinationClient) DeleteStale(ctx context.Context, in *DeleteStale_Request, opts ...grpc.CallOption) (*DeleteStale_Response, error) {
	out := new(DeleteStale_Response)
	err := c.cc.Invoke(ctx, "/proto.Destination/DeleteStale", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DestinationServer is the server API for Destination service.
// All implementations must embed UnimplementedDestinationServer
// for forward compatibility
type DestinationServer interface {
	// Get the current protocol version of the plugin. This helps
	// get the right message about upgrade/downgrade of cli and/or plugin.
	// Also, on the cli side it can try to upgrade/downgrade the protocol if cli supports it.
	GetProtocolVersion(context.Context, *GetProtocolVersion_Request) (*GetProtocolVersion_Response, error)
	// Get the name of the plugin
	GetName(context.Context, *GetName_Request) (*GetName_Response, error)
	// Get the current version of the plugin
	GetVersion(context.Context, *GetVersion_Request) (*GetVersion_Response, error)
	// Configure the plugin with the given credentials and mode
	Configure(context.Context, *Configure_Request) (*Configure_Response, error)
	// Migrate tables to the given plugin version
	Migrate(context.Context, *Migrate_Request) (*Migrate_Response, error)
	// Write resources
	Write(Destination_WriteServer) error
	// Send signal to flush and close open connections
	Close(context.Context, *Close_Request) (*Close_Response, error)
	// DeleteStale deletes stale data that was inserted by a given source
	// and is older than the given timestamp
	DeleteStale(context.Context, *DeleteStale_Request) (*DeleteStale_Response, error)
	mustEmbedUnimplementedDestinationServer()
}

// UnimplementedDestinationServer must be embedded to have forward compatible implementations.
type UnimplementedDestinationServer struct {
}

func (UnimplementedDestinationServer) GetProtocolVersion(context.Context, *GetProtocolVersion_Request) (*GetProtocolVersion_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProtocolVersion not implemented")
}
func (UnimplementedDestinationServer) GetName(context.Context, *GetName_Request) (*GetName_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetName not implemented")
}
func (UnimplementedDestinationServer) GetVersion(context.Context, *GetVersion_Request) (*GetVersion_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetVersion not implemented")
}
func (UnimplementedDestinationServer) Configure(context.Context, *Configure_Request) (*Configure_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Configure not implemented")
}
func (UnimplementedDestinationServer) Migrate(context.Context, *Migrate_Request) (*Migrate_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Migrate not implemented")
}
func (UnimplementedDestinationServer) Write(Destination_WriteServer) error {
	return status.Errorf(codes.Unimplemented, "method Write not implemented")
}
func (UnimplementedDestinationServer) Close(context.Context, *Close_Request) (*Close_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Close not implemented")
}
func (UnimplementedDestinationServer) DeleteStale(context.Context, *DeleteStale_Request) (*DeleteStale_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteStale not implemented")
}
func (UnimplementedDestinationServer) mustEmbedUnimplementedDestinationServer() {}

// UnsafeDestinationServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DestinationServer will
// result in compilation errors.
type UnsafeDestinationServer interface {
	mustEmbedUnimplementedDestinationServer()
}

func RegisterDestinationServer(s grpc.ServiceRegistrar, srv DestinationServer) {
	s.RegisterService(&Destination_ServiceDesc, srv)
}

func _Destination_GetProtocolVersion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetProtocolVersion_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DestinationServer).GetProtocolVersion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Destination/GetProtocolVersion",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DestinationServer).GetProtocolVersion(ctx, req.(*GetProtocolVersion_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Destination_GetName_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetName_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DestinationServer).GetName(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Destination/GetName",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DestinationServer).GetName(ctx, req.(*GetName_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Destination_GetVersion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetVersion_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DestinationServer).GetVersion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Destination/GetVersion",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DestinationServer).GetVersion(ctx, req.(*GetVersion_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Destination_Configure_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Configure_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DestinationServer).Configure(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Destination/Configure",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DestinationServer).Configure(ctx, req.(*Configure_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Destination_Migrate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Migrate_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DestinationServer).Migrate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Destination/Migrate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DestinationServer).Migrate(ctx, req.(*Migrate_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Destination_Write_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(DestinationServer).Write(&destinationWriteServer{stream})
}

type Destination_WriteServer interface {
	SendAndClose(*Write_Response) error
	Recv() (*Write_Request, error)
	grpc.ServerStream
}

type destinationWriteServer struct {
	grpc.ServerStream
}

func (x *destinationWriteServer) SendAndClose(m *Write_Response) error {
	return x.ServerStream.SendMsg(m)
}

func (x *destinationWriteServer) Recv() (*Write_Request, error) {
	m := new(Write_Request)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Destination_Close_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Close_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DestinationServer).Close(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Destination/Close",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DestinationServer).Close(ctx, req.(*Close_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Destination_DeleteStale_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteStale_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DestinationServer).DeleteStale(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Destination/DeleteStale",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DestinationServer).DeleteStale(ctx, req.(*DeleteStale_Request))
	}
	return interceptor(ctx, in, info, handler)
}

// Destination_ServiceDesc is the grpc.ServiceDesc for Destination service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Destination_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Destination",
	HandlerType: (*DestinationServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetProtocolVersion",
			Handler:    _Destination_GetProtocolVersion_Handler,
		},
		{
			MethodName: "GetName",
			Handler:    _Destination_GetName_Handler,
		},
		{
			MethodName: "GetVersion",
			Handler:    _Destination_GetVersion_Handler,
		},
		{
			MethodName: "Configure",
			Handler:    _Destination_Configure_Handler,
		},
		{
			MethodName: "Migrate",
			Handler:    _Destination_Migrate_Handler,
		},
		{
			MethodName: "Close",
			Handler:    _Destination_Close_Handler,
		},
		{
			MethodName: "DeleteStale",
			Handler:    _Destination_DeleteStale_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Write",
			Handler:       _Destination_Write_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "internal/pb/destination.proto",
}
