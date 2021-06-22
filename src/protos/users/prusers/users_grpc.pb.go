// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package prusers

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

// UsersClient is the client API for Users service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UsersClient interface {
	CheckEmail(ctx context.Context, in *CheckEmailRequest, opts ...grpc.CallOption) (*CheckEmailResponse, error)
	CheckPassword(ctx context.Context, in *CheckPasswordRequest, opts ...grpc.CallOption) (*CheckPasswordResponse, error)
	Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error)
	VerifyEmail(ctx context.Context, in *VerifyEmailRequest, opts ...grpc.CallOption) (*VerifyEmailResponse, error)
	ChangePassword(ctx context.Context, in *ChangeRequest, opts ...grpc.CallOption) (*ChangeResponse, error)
	ResetPassword(ctx context.Context, in *UserResetRequest, opts ...grpc.CallOption) (*UserResetResponse, error)
	GetByUsername(ctx context.Context, in *GetByUsernameRequest, opts ...grpc.CallOption) (*GetByUsernameResponse, error)
	GetRole(ctx context.Context, in *RoleRequest, opts ...grpc.CallOption) (*RoleResponse, error)
	UpdateUsername(ctx context.Context, in *UsernameRequest, opts ...grpc.CallOption) (*UsernameResponse, error)
	UpdateProfile(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*UpdateResponse, error)
	GetProfileByUsername(ctx context.Context, in *ProfileRequest, opts ...grpc.CallOption) (*ProfileResponse, error)
	ChangeProfilePublic(ctx context.Context, in *ChangePublicRequest, opts ...grpc.CallOption) (*ChangePublicResponse, error)
	ChangeProfileTaggable(ctx context.Context, in *ChangeTaggableRequest, opts ...grpc.CallOption) (*ChangeTaggableResponse, error)
	Follow(ctx context.Context, in *FollowRequest, opts ...grpc.CallOption) (*FollowRespose, error)
	UnFollow(ctx context.Context, in *FollowRequest, opts ...grpc.CallOption) (*FollowRespose, error)
	GetFollowers(ctx context.Context, in *FollowerRequest, opts ...grpc.CallOption) (Users_GetFollowersClient, error)
	GerFollowing(ctx context.Context, in *FollowerRequest, opts ...grpc.CallOption) (Users_GerFollowingClient, error)
	GetSearchedUsers(ctx context.Context, in *SearchRequest, opts ...grpc.CallOption) (*SearchResponse, error)
	GetFollowRequests(ctx context.Context, in *Profile, opts ...grpc.CallOption) (Users_GetFollowRequestsClient, error)
	SetFollowRequestRespond(ctx context.Context, in *FollowRequestRespond, opts ...grpc.CallOption) (*FollowRequestSet, error)
	GetFollowersDetailed(ctx context.Context, in *ProflieFollowRequest, opts ...grpc.CallOption) (Users_GetFollowersDetailedClient, error)
	GetFollowingDetailed(ctx context.Context, in *ProflieFollowRequest, opts ...grpc.CallOption) (Users_GetFollowingDetailedClient, error)
	CheckIfFollowing(ctx context.Context, in *ProflieFollowRequest, opts ...grpc.CallOption) (*BoolResponse, error)
	CheckForFollowingRequest(ctx context.Context, in *ProflieFollowRequest, opts ...grpc.CallOption) (*BoolResponse, error)
}

type usersClient struct {
	cc grpc.ClientConnInterface
}

func NewUsersClient(cc grpc.ClientConnInterface) UsersClient {
	return &usersClient{cc}
}

func (c *usersClient) CheckEmail(ctx context.Context, in *CheckEmailRequest, opts ...grpc.CallOption) (*CheckEmailResponse, error) {
	out := new(CheckEmailResponse)
	err := c.cc.Invoke(ctx, "/Users/CheckEmail", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersClient) CheckPassword(ctx context.Context, in *CheckPasswordRequest, opts ...grpc.CallOption) (*CheckPasswordResponse, error) {
	out := new(CheckPasswordResponse)
	err := c.cc.Invoke(ctx, "/Users/CheckPassword", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersClient) Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error) {
	out := new(RegisterResponse)
	err := c.cc.Invoke(ctx, "/Users/Register", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersClient) VerifyEmail(ctx context.Context, in *VerifyEmailRequest, opts ...grpc.CallOption) (*VerifyEmailResponse, error) {
	out := new(VerifyEmailResponse)
	err := c.cc.Invoke(ctx, "/Users/VerifyEmail", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersClient) ChangePassword(ctx context.Context, in *ChangeRequest, opts ...grpc.CallOption) (*ChangeResponse, error) {
	out := new(ChangeResponse)
	err := c.cc.Invoke(ctx, "/Users/ChangePassword", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersClient) ResetPassword(ctx context.Context, in *UserResetRequest, opts ...grpc.CallOption) (*UserResetResponse, error) {
	out := new(UserResetResponse)
	err := c.cc.Invoke(ctx, "/Users/ResetPassword", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersClient) GetByUsername(ctx context.Context, in *GetByUsernameRequest, opts ...grpc.CallOption) (*GetByUsernameResponse, error) {
	out := new(GetByUsernameResponse)
	err := c.cc.Invoke(ctx, "/Users/GetByUsername", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersClient) GetRole(ctx context.Context, in *RoleRequest, opts ...grpc.CallOption) (*RoleResponse, error) {
	out := new(RoleResponse)
	err := c.cc.Invoke(ctx, "/Users/GetRole", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersClient) UpdateUsername(ctx context.Context, in *UsernameRequest, opts ...grpc.CallOption) (*UsernameResponse, error) {
	out := new(UsernameResponse)
	err := c.cc.Invoke(ctx, "/Users/UpdateUsername", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersClient) UpdateProfile(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*UpdateResponse, error) {
	out := new(UpdateResponse)
	err := c.cc.Invoke(ctx, "/Users/UpdateProfile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersClient) GetProfileByUsername(ctx context.Context, in *ProfileRequest, opts ...grpc.CallOption) (*ProfileResponse, error) {
	out := new(ProfileResponse)
	err := c.cc.Invoke(ctx, "/Users/GetProfileByUsername", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersClient) ChangeProfilePublic(ctx context.Context, in *ChangePublicRequest, opts ...grpc.CallOption) (*ChangePublicResponse, error) {
	out := new(ChangePublicResponse)
	err := c.cc.Invoke(ctx, "/Users/ChangeProfilePublic", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersClient) ChangeProfileTaggable(ctx context.Context, in *ChangeTaggableRequest, opts ...grpc.CallOption) (*ChangeTaggableResponse, error) {
	out := new(ChangeTaggableResponse)
	err := c.cc.Invoke(ctx, "/Users/ChangeProfileTaggable", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersClient) Follow(ctx context.Context, in *FollowRequest, opts ...grpc.CallOption) (*FollowRespose, error) {
	out := new(FollowRespose)
	err := c.cc.Invoke(ctx, "/Users/Follow", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersClient) UnFollow(ctx context.Context, in *FollowRequest, opts ...grpc.CallOption) (*FollowRespose, error) {
	out := new(FollowRespose)
	err := c.cc.Invoke(ctx, "/Users/UnFollow", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersClient) GetFollowers(ctx context.Context, in *FollowerRequest, opts ...grpc.CallOption) (Users_GetFollowersClient, error) {
	stream, err := c.cc.NewStream(ctx, &Users_ServiceDesc.Streams[0], "/Users/GetFollowers", opts...)
	if err != nil {
		return nil, err
	}
	x := &usersGetFollowersClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Users_GetFollowersClient interface {
	Recv() (*ProfileFollower, error)
	grpc.ClientStream
}

type usersGetFollowersClient struct {
	grpc.ClientStream
}

func (x *usersGetFollowersClient) Recv() (*ProfileFollower, error) {
	m := new(ProfileFollower)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *usersClient) GerFollowing(ctx context.Context, in *FollowerRequest, opts ...grpc.CallOption) (Users_GerFollowingClient, error) {
	stream, err := c.cc.NewStream(ctx, &Users_ServiceDesc.Streams[1], "/Users/GerFollowing", opts...)
	if err != nil {
		return nil, err
	}
	x := &usersGerFollowingClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Users_GerFollowingClient interface {
	Recv() (*ProfileFollower, error)
	grpc.ClientStream
}

type usersGerFollowingClient struct {
	grpc.ClientStream
}

func (x *usersGerFollowingClient) Recv() (*ProfileFollower, error) {
	m := new(ProfileFollower)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *usersClient) GetSearchedUsers(ctx context.Context, in *SearchRequest, opts ...grpc.CallOption) (*SearchResponse, error) {
	out := new(SearchResponse)
	err := c.cc.Invoke(ctx, "/Users/GetSearchedUsers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersClient) GetFollowRequests(ctx context.Context, in *Profile, opts ...grpc.CallOption) (Users_GetFollowRequestsClient, error) {
	stream, err := c.cc.NewStream(ctx, &Users_ServiceDesc.Streams[2], "/Users/GetFollowRequests", opts...)
	if err != nil {
		return nil, err
	}
	x := &usersGetFollowRequestsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Users_GetFollowRequestsClient interface {
	Recv() (*FollowingRequest, error)
	grpc.ClientStream
}

type usersGetFollowRequestsClient struct {
	grpc.ClientStream
}

func (x *usersGetFollowRequestsClient) Recv() (*FollowingRequest, error) {
	m := new(FollowingRequest)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *usersClient) SetFollowRequestRespond(ctx context.Context, in *FollowRequestRespond, opts ...grpc.CallOption) (*FollowRequestSet, error) {
	out := new(FollowRequestSet)
	err := c.cc.Invoke(ctx, "/Users/SetFollowRequestRespond", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersClient) GetFollowersDetailed(ctx context.Context, in *ProflieFollowRequest, opts ...grpc.CallOption) (Users_GetFollowersDetailedClient, error) {
	stream, err := c.cc.NewStream(ctx, &Users_ServiceDesc.Streams[3], "/Users/GetFollowersDetailed", opts...)
	if err != nil {
		return nil, err
	}
	x := &usersGetFollowersDetailedClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Users_GetFollowersDetailedClient interface {
	Recv() (*ProfileFollowDetaild, error)
	grpc.ClientStream
}

type usersGetFollowersDetailedClient struct {
	grpc.ClientStream
}

func (x *usersGetFollowersDetailedClient) Recv() (*ProfileFollowDetaild, error) {
	m := new(ProfileFollowDetaild)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *usersClient) GetFollowingDetailed(ctx context.Context, in *ProflieFollowRequest, opts ...grpc.CallOption) (Users_GetFollowingDetailedClient, error) {
	stream, err := c.cc.NewStream(ctx, &Users_ServiceDesc.Streams[4], "/Users/GetFollowingDetailed", opts...)
	if err != nil {
		return nil, err
	}
	x := &usersGetFollowingDetailedClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Users_GetFollowingDetailedClient interface {
	Recv() (*ProfileFollowDetaild, error)
	grpc.ClientStream
}

type usersGetFollowingDetailedClient struct {
	grpc.ClientStream
}

func (x *usersGetFollowingDetailedClient) Recv() (*ProfileFollowDetaild, error) {
	m := new(ProfileFollowDetaild)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *usersClient) CheckIfFollowing(ctx context.Context, in *ProflieFollowRequest, opts ...grpc.CallOption) (*BoolResponse, error) {
	out := new(BoolResponse)
	err := c.cc.Invoke(ctx, "/Users/CheckIfFollowing", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersClient) CheckForFollowingRequest(ctx context.Context, in *ProflieFollowRequest, opts ...grpc.CallOption) (*BoolResponse, error) {
	out := new(BoolResponse)
	err := c.cc.Invoke(ctx, "/Users/CheckForFollowingRequest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UsersServer is the server API for Users service.
// All implementations must embed UnimplementedUsersServer
// for forward compatibility
type UsersServer interface {
	CheckEmail(context.Context, *CheckEmailRequest) (*CheckEmailResponse, error)
	CheckPassword(context.Context, *CheckPasswordRequest) (*CheckPasswordResponse, error)
	Register(context.Context, *RegisterRequest) (*RegisterResponse, error)
	VerifyEmail(context.Context, *VerifyEmailRequest) (*VerifyEmailResponse, error)
	ChangePassword(context.Context, *ChangeRequest) (*ChangeResponse, error)
	ResetPassword(context.Context, *UserResetRequest) (*UserResetResponse, error)
	GetByUsername(context.Context, *GetByUsernameRequest) (*GetByUsernameResponse, error)
	GetRole(context.Context, *RoleRequest) (*RoleResponse, error)
	UpdateUsername(context.Context, *UsernameRequest) (*UsernameResponse, error)
	UpdateProfile(context.Context, *UpdateRequest) (*UpdateResponse, error)
	GetProfileByUsername(context.Context, *ProfileRequest) (*ProfileResponse, error)
	ChangeProfilePublic(context.Context, *ChangePublicRequest) (*ChangePublicResponse, error)
	ChangeProfileTaggable(context.Context, *ChangeTaggableRequest) (*ChangeTaggableResponse, error)
	Follow(context.Context, *FollowRequest) (*FollowRespose, error)
	UnFollow(context.Context, *FollowRequest) (*FollowRespose, error)
	GetFollowers(*FollowerRequest, Users_GetFollowersServer) error
	GerFollowing(*FollowerRequest, Users_GerFollowingServer) error
	GetSearchedUsers(context.Context, *SearchRequest) (*SearchResponse, error)
	GetFollowRequests(*Profile, Users_GetFollowRequestsServer) error
	SetFollowRequestRespond(context.Context, *FollowRequestRespond) (*FollowRequestSet, error)
	GetFollowersDetailed(*ProflieFollowRequest, Users_GetFollowersDetailedServer) error
	GetFollowingDetailed(*ProflieFollowRequest, Users_GetFollowingDetailedServer) error
	CheckIfFollowing(context.Context, *ProflieFollowRequest) (*BoolResponse, error)
	CheckForFollowingRequest(context.Context, *ProflieFollowRequest) (*BoolResponse, error)
	mustEmbedUnimplementedUsersServer()
}

// UnimplementedUsersServer must be embedded to have forward compatible implementations.
type UnimplementedUsersServer struct {
}

func (UnimplementedUsersServer) CheckEmail(context.Context, *CheckEmailRequest) (*CheckEmailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckEmail not implemented")
}
func (UnimplementedUsersServer) CheckPassword(context.Context, *CheckPasswordRequest) (*CheckPasswordResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckPassword not implemented")
}
func (UnimplementedUsersServer) Register(context.Context, *RegisterRequest) (*RegisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedUsersServer) VerifyEmail(context.Context, *VerifyEmailRequest) (*VerifyEmailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VerifyEmail not implemented")
}
func (UnimplementedUsersServer) ChangePassword(context.Context, *ChangeRequest) (*ChangeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChangePassword not implemented")
}
func (UnimplementedUsersServer) ResetPassword(context.Context, *UserResetRequest) (*UserResetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResetPassword not implemented")
}
func (UnimplementedUsersServer) GetByUsername(context.Context, *GetByUsernameRequest) (*GetByUsernameResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetByUsername not implemented")
}
func (UnimplementedUsersServer) GetRole(context.Context, *RoleRequest) (*RoleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRole not implemented")
}
func (UnimplementedUsersServer) UpdateUsername(context.Context, *UsernameRequest) (*UsernameResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUsername not implemented")
}
func (UnimplementedUsersServer) UpdateProfile(context.Context, *UpdateRequest) (*UpdateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateProfile not implemented")
}
func (UnimplementedUsersServer) GetProfileByUsername(context.Context, *ProfileRequest) (*ProfileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProfileByUsername not implemented")
}
func (UnimplementedUsersServer) ChangeProfilePublic(context.Context, *ChangePublicRequest) (*ChangePublicResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChangeProfilePublic not implemented")
}
func (UnimplementedUsersServer) ChangeProfileTaggable(context.Context, *ChangeTaggableRequest) (*ChangeTaggableResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChangeProfileTaggable not implemented")
}
func (UnimplementedUsersServer) Follow(context.Context, *FollowRequest) (*FollowRespose, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Follow not implemented")
}
func (UnimplementedUsersServer) UnFollow(context.Context, *FollowRequest) (*FollowRespose, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnFollow not implemented")
}
func (UnimplementedUsersServer) GetFollowers(*FollowerRequest, Users_GetFollowersServer) error {
	return status.Errorf(codes.Unimplemented, "method GetFollowers not implemented")
}
func (UnimplementedUsersServer) GerFollowing(*FollowerRequest, Users_GerFollowingServer) error {
	return status.Errorf(codes.Unimplemented, "method GerFollowing not implemented")
}
func (UnimplementedUsersServer) GetSearchedUsers(context.Context, *SearchRequest) (*SearchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSearchedUsers not implemented")
}
func (UnimplementedUsersServer) GetFollowRequests(*Profile, Users_GetFollowRequestsServer) error {
	return status.Errorf(codes.Unimplemented, "method GetFollowRequests not implemented")
}
func (UnimplementedUsersServer) SetFollowRequestRespond(context.Context, *FollowRequestRespond) (*FollowRequestSet, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetFollowRequestRespond not implemented")
}
func (UnimplementedUsersServer) GetFollowersDetailed(*ProflieFollowRequest, Users_GetFollowersDetailedServer) error {
	return status.Errorf(codes.Unimplemented, "method GetFollowersDetailed not implemented")
}
func (UnimplementedUsersServer) GetFollowingDetailed(*ProflieFollowRequest, Users_GetFollowingDetailedServer) error {
	return status.Errorf(codes.Unimplemented, "method GetFollowingDetailed not implemented")
}
func (UnimplementedUsersServer) CheckIfFollowing(context.Context, *ProflieFollowRequest) (*BoolResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckIfFollowing not implemented")
}
func (UnimplementedUsersServer) CheckForFollowingRequest(context.Context, *ProflieFollowRequest) (*BoolResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckForFollowingRequest not implemented")
}
func (UnimplementedUsersServer) mustEmbedUnimplementedUsersServer() {}

// UnsafeUsersServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UsersServer will
// result in compilation errors.
type UnsafeUsersServer interface {
	mustEmbedUnimplementedUsersServer()
}

func RegisterUsersServer(s grpc.ServiceRegistrar, srv UsersServer) {
	s.RegisterService(&Users_ServiceDesc, srv)
}

func _Users_CheckEmail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckEmailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServer).CheckEmail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Users/CheckEmail",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServer).CheckEmail(ctx, req.(*CheckEmailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Users_CheckPassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckPasswordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServer).CheckPassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Users/CheckPassword",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServer).CheckPassword(ctx, req.(*CheckPasswordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Users_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Users/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServer).Register(ctx, req.(*RegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Users_VerifyEmail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VerifyEmailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServer).VerifyEmail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Users/VerifyEmail",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServer).VerifyEmail(ctx, req.(*VerifyEmailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Users_ChangePassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChangeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServer).ChangePassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Users/ChangePassword",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServer).ChangePassword(ctx, req.(*ChangeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Users_ResetPassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserResetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServer).ResetPassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Users/ResetPassword",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServer).ResetPassword(ctx, req.(*UserResetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Users_GetByUsername_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetByUsernameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServer).GetByUsername(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Users/GetByUsername",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServer).GetByUsername(ctx, req.(*GetByUsernameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Users_GetRole_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RoleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServer).GetRole(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Users/GetRole",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServer).GetRole(ctx, req.(*RoleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Users_UpdateUsername_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UsernameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServer).UpdateUsername(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Users/UpdateUsername",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServer).UpdateUsername(ctx, req.(*UsernameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Users_UpdateProfile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServer).UpdateProfile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Users/UpdateProfile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServer).UpdateProfile(ctx, req.(*UpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Users_GetProfileByUsername_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProfileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServer).GetProfileByUsername(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Users/GetProfileByUsername",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServer).GetProfileByUsername(ctx, req.(*ProfileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Users_ChangeProfilePublic_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChangePublicRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServer).ChangeProfilePublic(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Users/ChangeProfilePublic",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServer).ChangeProfilePublic(ctx, req.(*ChangePublicRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Users_ChangeProfileTaggable_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChangeTaggableRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServer).ChangeProfileTaggable(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Users/ChangeProfileTaggable",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServer).ChangeProfileTaggable(ctx, req.(*ChangeTaggableRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Users_Follow_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FollowRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServer).Follow(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Users/Follow",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServer).Follow(ctx, req.(*FollowRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Users_UnFollow_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FollowRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServer).UnFollow(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Users/UnFollow",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServer).UnFollow(ctx, req.(*FollowRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Users_GetFollowers_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(FollowerRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(UsersServer).GetFollowers(m, &usersGetFollowersServer{stream})
}

type Users_GetFollowersServer interface {
	Send(*ProfileFollower) error
	grpc.ServerStream
}

type usersGetFollowersServer struct {
	grpc.ServerStream
}

func (x *usersGetFollowersServer) Send(m *ProfileFollower) error {
	return x.ServerStream.SendMsg(m)
}

func _Users_GerFollowing_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(FollowerRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(UsersServer).GerFollowing(m, &usersGerFollowingServer{stream})
}

type Users_GerFollowingServer interface {
	Send(*ProfileFollower) error
	grpc.ServerStream
}

type usersGerFollowingServer struct {
	grpc.ServerStream
}

func (x *usersGerFollowingServer) Send(m *ProfileFollower) error {
	return x.ServerStream.SendMsg(m)
}

func _Users_GetSearchedUsers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServer).GetSearchedUsers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Users/GetSearchedUsers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServer).GetSearchedUsers(ctx, req.(*SearchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Users_GetFollowRequests_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Profile)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(UsersServer).GetFollowRequests(m, &usersGetFollowRequestsServer{stream})
}

type Users_GetFollowRequestsServer interface {
	Send(*FollowingRequest) error
	grpc.ServerStream
}

type usersGetFollowRequestsServer struct {
	grpc.ServerStream
}

func (x *usersGetFollowRequestsServer) Send(m *FollowingRequest) error {
	return x.ServerStream.SendMsg(m)
}

func _Users_SetFollowRequestRespond_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FollowRequestRespond)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServer).SetFollowRequestRespond(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Users/SetFollowRequestRespond",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServer).SetFollowRequestRespond(ctx, req.(*FollowRequestRespond))
	}
	return interceptor(ctx, in, info, handler)
}

func _Users_GetFollowersDetailed_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ProflieFollowRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(UsersServer).GetFollowersDetailed(m, &usersGetFollowersDetailedServer{stream})
}

type Users_GetFollowersDetailedServer interface {
	Send(*ProfileFollowDetaild) error
	grpc.ServerStream
}

type usersGetFollowersDetailedServer struct {
	grpc.ServerStream
}

func (x *usersGetFollowersDetailedServer) Send(m *ProfileFollowDetaild) error {
	return x.ServerStream.SendMsg(m)
}

func _Users_GetFollowingDetailed_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ProflieFollowRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(UsersServer).GetFollowingDetailed(m, &usersGetFollowingDetailedServer{stream})
}

type Users_GetFollowingDetailedServer interface {
	Send(*ProfileFollowDetaild) error
	grpc.ServerStream
}

type usersGetFollowingDetailedServer struct {
	grpc.ServerStream
}

func (x *usersGetFollowingDetailedServer) Send(m *ProfileFollowDetaild) error {
	return x.ServerStream.SendMsg(m)
}

func _Users_CheckIfFollowing_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProflieFollowRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServer).CheckIfFollowing(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Users/CheckIfFollowing",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServer).CheckIfFollowing(ctx, req.(*ProflieFollowRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Users_CheckForFollowingRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProflieFollowRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServer).CheckForFollowingRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Users/CheckForFollowingRequest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServer).CheckForFollowingRequest(ctx, req.(*ProflieFollowRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Users_ServiceDesc is the grpc.ServiceDesc for Users service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Users_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Users",
	HandlerType: (*UsersServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CheckEmail",
			Handler:    _Users_CheckEmail_Handler,
		},
		{
			MethodName: "CheckPassword",
			Handler:    _Users_CheckPassword_Handler,
		},
		{
			MethodName: "Register",
			Handler:    _Users_Register_Handler,
		},
		{
			MethodName: "VerifyEmail",
			Handler:    _Users_VerifyEmail_Handler,
		},
		{
			MethodName: "ChangePassword",
			Handler:    _Users_ChangePassword_Handler,
		},
		{
			MethodName: "ResetPassword",
			Handler:    _Users_ResetPassword_Handler,
		},
		{
			MethodName: "GetByUsername",
			Handler:    _Users_GetByUsername_Handler,
		},
		{
			MethodName: "GetRole",
			Handler:    _Users_GetRole_Handler,
		},
		{
			MethodName: "UpdateUsername",
			Handler:    _Users_UpdateUsername_Handler,
		},
		{
			MethodName: "UpdateProfile",
			Handler:    _Users_UpdateProfile_Handler,
		},
		{
			MethodName: "GetProfileByUsername",
			Handler:    _Users_GetProfileByUsername_Handler,
		},
		{
			MethodName: "ChangeProfilePublic",
			Handler:    _Users_ChangeProfilePublic_Handler,
		},
		{
			MethodName: "ChangeProfileTaggable",
			Handler:    _Users_ChangeProfileTaggable_Handler,
		},
		{
			MethodName: "Follow",
			Handler:    _Users_Follow_Handler,
		},
		{
			MethodName: "UnFollow",
			Handler:    _Users_UnFollow_Handler,
		},
		{
			MethodName: "GetSearchedUsers",
			Handler:    _Users_GetSearchedUsers_Handler,
		},
		{
			MethodName: "SetFollowRequestRespond",
			Handler:    _Users_SetFollowRequestRespond_Handler,
		},
		{
			MethodName: "CheckIfFollowing",
			Handler:    _Users_CheckIfFollowing_Handler,
		},
		{
			MethodName: "CheckForFollowingRequest",
			Handler:    _Users_CheckForFollowingRequest_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetFollowers",
			Handler:       _Users_GetFollowers_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "GerFollowing",
			Handler:       _Users_GerFollowing_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "GetFollowRequests",
			Handler:       _Users_GetFollowRequests_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "GetFollowersDetailed",
			Handler:       _Users_GetFollowersDetailed_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "GetFollowingDetailed",
			Handler:       _Users_GetFollowingDetailed_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "users/users.proto",
}
