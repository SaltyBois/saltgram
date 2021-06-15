// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.6.1
// source: admin/admin.proto

package pradmin

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type AddVerificationRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FullName string `protobuf:"bytes,1,opt,name=fullName,proto3" json:"fullName,omitempty"`
	UserId   uint64 `protobuf:"varint,2,opt,name=userId,proto3" json:"userId,omitempty"` //Category category = 3;
}

func (x *AddVerificationRequest) Reset() {
	*x = AddVerificationRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_admin_admin_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddVerificationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddVerificationRequest) ProtoMessage() {}

func (x *AddVerificationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_admin_admin_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddVerificationRequest.ProtoReflect.Descriptor instead.
func (*AddVerificationRequest) Descriptor() ([]byte, []int) {
	return file_admin_admin_proto_rawDescGZIP(), []int{0}
}

func (x *AddVerificationRequest) GetFullName() string {
	if x != nil {
		return x.FullName
	}
	return ""
}

func (x *AddVerificationRequest) GetUserId() uint64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

type AddVerificationResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *AddVerificationResponse) Reset() {
	*x = AddVerificationResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_admin_admin_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddVerificationResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddVerificationResponse) ProtoMessage() {}

func (x *AddVerificationResponse) ProtoReflect() protoreflect.Message {
	mi := &file_admin_admin_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddVerificationResponse.ProtoReflect.Descriptor instead.
func (*AddVerificationResponse) Descriptor() ([]byte, []int) {
	return file_admin_admin_proto_rawDescGZIP(), []int{1}
}

type GetVerificationRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetVerificationRequest) Reset() {
	*x = GetVerificationRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_admin_admin_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetVerificationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetVerificationRequest) ProtoMessage() {}

func (x *GetVerificationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_admin_admin_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetVerificationRequest.ProtoReflect.Descriptor instead.
func (*GetVerificationRequest) Descriptor() ([]byte, []int) {
	return file_admin_admin_proto_rawDescGZIP(), []int{2}
}

type GetVerificationResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	FullName string `protobuf:"bytes,2,opt,name=fullName,proto3" json:"fullName,omitempty"` //Category category = 3;
}

func (x *GetVerificationResponse) Reset() {
	*x = GetVerificationResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_admin_admin_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetVerificationResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetVerificationResponse) ProtoMessage() {}

func (x *GetVerificationResponse) ProtoReflect() protoreflect.Message {
	mi := &file_admin_admin_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetVerificationResponse.ProtoReflect.Descriptor instead.
func (*GetVerificationResponse) Descriptor() ([]byte, []int) {
	return file_admin_admin_proto_rawDescGZIP(), []int{3}
}

func (x *GetVerificationResponse) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *GetVerificationResponse) GetFullName() string {
	if x != nil {
		return x.FullName
	}
	return ""
}

type ReviewVerificatonRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"` //VerificationStatus verficationStatus = 2;
}

func (x *ReviewVerificatonRequest) Reset() {
	*x = ReviewVerificatonRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_admin_admin_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReviewVerificatonRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReviewVerificatonRequest) ProtoMessage() {}

func (x *ReviewVerificatonRequest) ProtoReflect() protoreflect.Message {
	mi := &file_admin_admin_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReviewVerificatonRequest.ProtoReflect.Descriptor instead.
func (*ReviewVerificatonRequest) Descriptor() ([]byte, []int) {
	return file_admin_admin_proto_rawDescGZIP(), []int{4}
}

func (x *ReviewVerificatonRequest) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type ReviewVerificatonResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ReviewVerificatonResponse) Reset() {
	*x = ReviewVerificatonResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_admin_admin_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReviewVerificatonResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReviewVerificatonResponse) ProtoMessage() {}

func (x *ReviewVerificatonResponse) ProtoReflect() protoreflect.Message {
	mi := &file_admin_admin_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReviewVerificatonResponse.ProtoReflect.Descriptor instead.
func (*ReviewVerificatonResponse) Descriptor() ([]byte, []int) {
	return file_admin_admin_proto_rawDescGZIP(), []int{5}
}

var File_admin_admin_proto protoreflect.FileDescriptor

var file_admin_admin_proto_rawDesc = []byte{
	0x0a, 0x11, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x4c, 0x0a, 0x16, 0x41, 0x64, 0x64, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a,
	0x08, 0x66, 0x75, 0x6c, 0x6c, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x66, 0x75, 0x6c, 0x6c, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65,
	0x72, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49,
	0x64, 0x22, 0x19, 0x0a, 0x17, 0x41, 0x64, 0x64, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x18, 0x0a, 0x16,
	0x47, 0x65, 0x74, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x45, 0x0a, 0x17, 0x47, 0x65, 0x74, 0x56, 0x65, 0x72,
	0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x75, 0x6c, 0x6c, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x75, 0x6c, 0x6c, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x2a, 0x0a,
	0x18, 0x52, 0x65, 0x76, 0x69, 0x65, 0x77, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74,
	0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x22, 0x1b, 0x0a, 0x19, 0x52, 0x65, 0x76,
	0x69, 0x65, 0x77, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x6f, 0x6e, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0xef, 0x01, 0x0a, 0x05, 0x41, 0x64, 0x6d, 0x69, 0x6e,
	0x12, 0x47, 0x0a, 0x12, 0x41, 0x64, 0x64, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x12, 0x17, 0x2e, 0x41, 0x64, 0x64, 0x56, 0x65, 0x72, 0x69,
	0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x18, 0x2e, 0x41, 0x64, 0x64, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4e, 0x0a, 0x17, 0x47, 0x65, 0x74,
	0x50, 0x65, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x12, 0x17, 0x2e, 0x47, 0x65, 0x74, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e,
	0x47, 0x65, 0x74, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x30, 0x01, 0x12, 0x4d, 0x0a, 0x14, 0x52, 0x65, 0x76,
	0x69, 0x65, 0x77, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x6f, 0x6e, 0x52, 0x65,
	0x71, 0x12, 0x19, 0x2e, 0x52, 0x65, 0x76, 0x69, 0x65, 0x77, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69,
	0x63, 0x61, 0x74, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x52,
	0x65, 0x76, 0x69, 0x65, 0x77, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x6f, 0x6e,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x11, 0x5a, 0x0f, 0x2e, 0x2f, 0x61, 0x64,
	0x6d, 0x69, 0x6e, 0x2f, 0x70, 0x72, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_admin_admin_proto_rawDescOnce sync.Once
	file_admin_admin_proto_rawDescData = file_admin_admin_proto_rawDesc
)

func file_admin_admin_proto_rawDescGZIP() []byte {
	file_admin_admin_proto_rawDescOnce.Do(func() {
		file_admin_admin_proto_rawDescData = protoimpl.X.CompressGZIP(file_admin_admin_proto_rawDescData)
	})
	return file_admin_admin_proto_rawDescData
}

var file_admin_admin_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_admin_admin_proto_goTypes = []interface{}{
	(*AddVerificationRequest)(nil),    // 0: AddVerificationRequest
	(*AddVerificationResponse)(nil),   // 1: AddVerificationResponse
	(*GetVerificationRequest)(nil),    // 2: GetVerificationRequest
	(*GetVerificationResponse)(nil),   // 3: GetVerificationResponse
	(*ReviewVerificatonRequest)(nil),  // 4: ReviewVerificatonRequest
	(*ReviewVerificatonResponse)(nil), // 5: ReviewVerificatonResponse
}
var file_admin_admin_proto_depIdxs = []int32{
	0, // 0: Admin.AddVerificationReq:input_type -> AddVerificationRequest
	2, // 1: Admin.GetPendingVerifications:input_type -> GetVerificationRequest
	4, // 2: Admin.ReviewVerificatonReq:input_type -> ReviewVerificatonRequest
	1, // 3: Admin.AddVerificationReq:output_type -> AddVerificationResponse
	3, // 4: Admin.GetPendingVerifications:output_type -> GetVerificationResponse
	5, // 5: Admin.ReviewVerificatonReq:output_type -> ReviewVerificatonResponse
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_admin_admin_proto_init() }
func file_admin_admin_proto_init() {
	if File_admin_admin_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_admin_admin_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddVerificationRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_admin_admin_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddVerificationResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_admin_admin_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetVerificationRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_admin_admin_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetVerificationResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_admin_admin_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReviewVerificatonRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_admin_admin_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReviewVerificatonResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_admin_admin_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_admin_admin_proto_goTypes,
		DependencyIndexes: file_admin_admin_proto_depIdxs,
		MessageInfos:      file_admin_admin_proto_msgTypes,
	}.Build()
	File_admin_admin_proto = out.File
	file_admin_admin_proto_rawDesc = nil
	file_admin_admin_proto_goTypes = nil
	file_admin_admin_proto_depIdxs = nil
}
