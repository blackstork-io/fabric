// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: pluginapi/v1/plugin.proto

package pluginapiv1

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

const (
	PluginService_GetSchema_FullMethodName      = "/pluginapi.v1.PluginService/GetSchema"
	PluginService_RetrieveData_FullMethodName   = "/pluginapi.v1.PluginService/RetrieveData"
	PluginService_ProvideContent_FullMethodName = "/pluginapi.v1.PluginService/ProvideContent"
)

// PluginServiceClient is the client API for PluginService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PluginServiceClient interface {
	GetSchema(ctx context.Context, in *GetSchemaRequest, opts ...grpc.CallOption) (*GetSchemaResponse, error)
	RetrieveData(ctx context.Context, in *RetrieveDataRequest, opts ...grpc.CallOption) (*RetrieveDataResponse, error)
	ProvideContent(ctx context.Context, in *ProvideContentRequest, opts ...grpc.CallOption) (*ProvideContentResponse, error)
}

type pluginServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPluginServiceClient(cc grpc.ClientConnInterface) PluginServiceClient {
	return &pluginServiceClient{cc}
}

func (c *pluginServiceClient) GetSchema(ctx context.Context, in *GetSchemaRequest, opts ...grpc.CallOption) (*GetSchemaResponse, error) {
	out := new(GetSchemaResponse)
	err := c.cc.Invoke(ctx, PluginService_GetSchema_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pluginServiceClient) RetrieveData(ctx context.Context, in *RetrieveDataRequest, opts ...grpc.CallOption) (*RetrieveDataResponse, error) {
	out := new(RetrieveDataResponse)
	err := c.cc.Invoke(ctx, PluginService_RetrieveData_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pluginServiceClient) ProvideContent(ctx context.Context, in *ProvideContentRequest, opts ...grpc.CallOption) (*ProvideContentResponse, error) {
	out := new(ProvideContentResponse)
	err := c.cc.Invoke(ctx, PluginService_ProvideContent_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PluginServiceServer is the server API for PluginService service.
// All implementations must embed UnimplementedPluginServiceServer
// for forward compatibility
type PluginServiceServer interface {
	GetSchema(context.Context, *GetSchemaRequest) (*GetSchemaResponse, error)
	RetrieveData(context.Context, *RetrieveDataRequest) (*RetrieveDataResponse, error)
	ProvideContent(context.Context, *ProvideContentRequest) (*ProvideContentResponse, error)
	mustEmbedUnimplementedPluginServiceServer()
}

// UnimplementedPluginServiceServer must be embedded to have forward compatible implementations.
type UnimplementedPluginServiceServer struct {
}

func (UnimplementedPluginServiceServer) GetSchema(context.Context, *GetSchemaRequest) (*GetSchemaResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSchema not implemented")
}
func (UnimplementedPluginServiceServer) RetrieveData(context.Context, *RetrieveDataRequest) (*RetrieveDataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RetrieveData not implemented")
}
func (UnimplementedPluginServiceServer) ProvideContent(context.Context, *ProvideContentRequest) (*ProvideContentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ProvideContent not implemented")
}
func (UnimplementedPluginServiceServer) mustEmbedUnimplementedPluginServiceServer() {}

// UnsafePluginServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PluginServiceServer will
// result in compilation errors.
type UnsafePluginServiceServer interface {
	mustEmbedUnimplementedPluginServiceServer()
}

func RegisterPluginServiceServer(s grpc.ServiceRegistrar, srv PluginServiceServer) {
	s.RegisterService(&PluginService_ServiceDesc, srv)
}

func _PluginService_GetSchema_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSchemaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginServiceServer).GetSchema(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PluginService_GetSchema_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginServiceServer).GetSchema(ctx, req.(*GetSchemaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PluginService_RetrieveData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RetrieveDataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginServiceServer).RetrieveData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PluginService_RetrieveData_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginServiceServer).RetrieveData(ctx, req.(*RetrieveDataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PluginService_ProvideContent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProvideContentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginServiceServer).ProvideContent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PluginService_ProvideContent_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginServiceServer).ProvideContent(ctx, req.(*ProvideContentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PluginService_ServiceDesc is the grpc.ServiceDesc for PluginService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PluginService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pluginapi.v1.PluginService",
	HandlerType: (*PluginServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetSchema",
			Handler:    _PluginService_GetSchema_Handler,
		},
		{
			MethodName: "RetrieveData",
			Handler:    _PluginService_RetrieveData_Handler,
		},
		{
			MethodName: "ProvideContent",
			Handler:    _PluginService_ProvideContent_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pluginapi/v1/plugin.proto",
}