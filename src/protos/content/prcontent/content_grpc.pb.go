// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package prcontent

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

// ContentClient is the client API for Content service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ContentClient interface {
	GetSharedMedia(ctx context.Context, in *SharedMediaRequest, opts ...grpc.CallOption) (Content_GetSharedMediaClient, error)
	AddSharedMedia(ctx context.Context, in *AddSharedMediaRequest, opts ...grpc.CallOption) (*AddSharedMediaResponse, error)
	PostProfile(ctx context.Context, in *PostProfileRequest, opts ...grpc.CallOption) (*PostProfileResponse, error)
}

type contentClient struct {
	cc grpc.ClientConnInterface
}

func NewContentClient(cc grpc.ClientConnInterface) ContentClient {
	return &contentClient{cc}
}

func (c *contentClient) GetSharedMedia(ctx context.Context, in *SharedMediaRequest, opts ...grpc.CallOption) (Content_GetSharedMediaClient, error) {
	stream, err := c.cc.NewStream(ctx, &Content_ServiceDesc.Streams[0], "/Content/GetSharedMedia", opts...)
	if err != nil {
		return nil, err
	}
	x := &contentGetSharedMediaClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Content_GetSharedMediaClient interface {
	Recv() (*SharedMediaResponse, error)
	grpc.ClientStream
}

type contentGetSharedMediaClient struct {
	grpc.ClientStream
}

func (x *contentGetSharedMediaClient) Recv() (*SharedMediaResponse, error) {
	m := new(SharedMediaResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *contentClient) AddSharedMedia(ctx context.Context, in *AddSharedMediaRequest, opts ...grpc.CallOption) (*AddSharedMediaResponse, error) {
	out := new(AddSharedMediaResponse)
	err := c.cc.Invoke(ctx, "/Content/AddSharedMedia", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *contentClient) PostProfile(ctx context.Context, in *PostProfileRequest, opts ...grpc.CallOption) (*PostProfileResponse, error) {
	out := new(PostProfileResponse)
	err := c.cc.Invoke(ctx, "/Content/PostProfile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ContentServer is the server API for Content service.
// All implementations must embed UnimplementedContentServer
// for forward compatibility
type ContentServer interface {
	GetSharedMedia(*SharedMediaRequest, Content_GetSharedMediaServer) error
	AddSharedMedia(context.Context, *AddSharedMediaRequest) (*AddSharedMediaResponse, error)
	PostProfile(context.Context, *PostProfileRequest) (*PostProfileResponse, error)
	mustEmbedUnimplementedContentServer()
}

// UnimplementedContentServer must be embedded to have forward compatible implementations.
type UnimplementedContentServer struct {
}

func (UnimplementedContentServer) GetSharedMedia(*SharedMediaRequest, Content_GetSharedMediaServer) error {
	return status.Errorf(codes.Unimplemented, "method GetSharedMedia not implemented")
}
func (UnimplementedContentServer) AddSharedMedia(context.Context, *AddSharedMediaRequest) (*AddSharedMediaResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddSharedMedia not implemented")
}
func (UnimplementedContentServer) PostProfile(context.Context, *PostProfileRequest) (*PostProfileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PostProfile not implemented")
}
func (UnimplementedContentServer) mustEmbedUnimplementedContentServer() {}

// UnsafeContentServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ContentServer will
// result in compilation errors.
type UnsafeContentServer interface {
	mustEmbedUnimplementedContentServer()
}

func RegisterContentServer(s grpc.ServiceRegistrar, srv ContentServer) {
	s.RegisterService(&Content_ServiceDesc, srv)
}

func _Content_GetSharedMedia_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SharedMediaRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ContentServer).GetSharedMedia(m, &contentGetSharedMediaServer{stream})
}

type Content_GetSharedMediaServer interface {
	Send(*SharedMediaResponse) error
	grpc.ServerStream
}

type contentGetSharedMediaServer struct {
	grpc.ServerStream
}

func (x *contentGetSharedMediaServer) Send(m *SharedMediaResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _Content_AddSharedMedia_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddSharedMediaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContentServer).AddSharedMedia(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Content/AddSharedMedia",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContentServer).AddSharedMedia(ctx, req.(*AddSharedMediaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Content_PostProfile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PostProfileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContentServer).PostProfile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Content/PostProfile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContentServer).PostProfile(ctx, req.(*PostProfileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Content_ServiceDesc is the grpc.ServiceDesc for Content service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Content_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Content",
	HandlerType: (*ContentServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddSharedMedia",
			Handler:    _Content_AddSharedMedia_Handler,
		},
		{
			MethodName: "PostProfile",
			Handler:    _Content_PostProfile_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetSharedMedia",
			Handler:       _Content_GetSharedMedia_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "content/content.proto",
}
