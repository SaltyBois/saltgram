// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.16.0
// source: content/content.proto

package prcontent

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

type SharedMediaRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId uint64 `protobuf:"varint,1,opt,name=userId,proto3" json:"userId,omitempty"`
}

func (x *SharedMediaRequest) Reset() {
	*x = SharedMediaRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_content_content_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SharedMediaRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SharedMediaRequest) ProtoMessage() {}

func (x *SharedMediaRequest) ProtoReflect() protoreflect.Message {
	mi := &file_content_content_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SharedMediaRequest.ProtoReflect.Descriptor instead.
func (*SharedMediaRequest) Descriptor() ([]byte, []int) {
	return file_content_content_proto_rawDescGZIP(), []int{0}
}

func (x *SharedMediaRequest) GetUserId() uint64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

type SharedMediaResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Media []*Media `protobuf:"bytes,1,rep,name=media,proto3" json:"media,omitempty"`
}

func (x *SharedMediaResponse) Reset() {
	*x = SharedMediaResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_content_content_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SharedMediaResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SharedMediaResponse) ProtoMessage() {}

func (x *SharedMediaResponse) ProtoReflect() protoreflect.Message {
	mi := &file_content_content_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SharedMediaResponse.ProtoReflect.Descriptor instead.
func (*SharedMediaResponse) Descriptor() ([]byte, []int) {
	return file_content_content_proto_rawDescGZIP(), []int{1}
}

func (x *SharedMediaResponse) GetMedia() []*Media {
	if x != nil {
		return x.Media
	}
	return nil
}

type Media struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId      uint64    `protobuf:"varint,1,opt,name=userId,proto3" json:"userId,omitempty"`
	Filename    string    `protobuf:"bytes,2,opt,name=filename,proto3" json:"filename,omitempty"`
	Description string    `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	AddedOn     string    `protobuf:"bytes,4,opt,name=addedOn,proto3" json:"addedOn,omitempty"`
	Location    *Location `protobuf:"bytes,5,opt,name=location,proto3" json:"location,omitempty"`
	Tags        []*Tag    `protobuf:"bytes,6,rep,name=tags,proto3" json:"tags,omitempty"`
}

func (x *Media) Reset() {
	*x = Media{}
	if protoimpl.UnsafeEnabled {
		mi := &file_content_content_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Media) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Media) ProtoMessage() {}

func (x *Media) ProtoReflect() protoreflect.Message {
	mi := &file_content_content_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Media.ProtoReflect.Descriptor instead.
func (*Media) Descriptor() ([]byte, []int) {
	return file_content_content_proto_rawDescGZIP(), []int{2}
}

func (x *Media) GetUserId() uint64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *Media) GetFilename() string {
	if x != nil {
		return x.Filename
	}
	return ""
}

func (x *Media) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Media) GetAddedOn() string {
	if x != nil {
		return x.AddedOn
	}
	return ""
}

func (x *Media) GetLocation() *Location {
	if x != nil {
		return x.Location
	}
	return nil
}

func (x *Media) GetTags() []*Tag {
	if x != nil {
		return x.Tags
	}
	return nil
}

type Tag struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value string `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	Id    uint64 `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *Tag) Reset() {
	*x = Tag{}
	if protoimpl.UnsafeEnabled {
		mi := &file_content_content_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Tag) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Tag) ProtoMessage() {}

func (x *Tag) ProtoReflect() protoreflect.Message {
	mi := &file_content_content_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Tag.ProtoReflect.Descriptor instead.
func (*Tag) Descriptor() ([]byte, []int) {
	return file_content_content_proto_rawDescGZIP(), []int{3}
}

func (x *Tag) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

func (x *Tag) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type Location struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Country string `protobuf:"bytes,1,opt,name=Country,proto3" json:"Country,omitempty"`
	State   string `protobuf:"bytes,2,opt,name=State,proto3" json:"State,omitempty"`
	ZipCode string `protobuf:"bytes,3,opt,name=ZipCode,proto3" json:"ZipCode,omitempty"`
	Street  string `protobuf:"bytes,4,opt,name=Street,proto3" json:"Street,omitempty"`
}

func (x *Location) Reset() {
	*x = Location{}
	if protoimpl.UnsafeEnabled {
		mi := &file_content_content_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Location) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Location) ProtoMessage() {}

func (x *Location) ProtoReflect() protoreflect.Message {
	mi := &file_content_content_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Location.ProtoReflect.Descriptor instead.
func (*Location) Descriptor() ([]byte, []int) {
	return file_content_content_proto_rawDescGZIP(), []int{4}
}

func (x *Location) GetCountry() string {
	if x != nil {
		return x.Country
	}
	return ""
}

func (x *Location) GetState() string {
	if x != nil {
		return x.State
	}
	return ""
}

func (x *Location) GetZipCode() string {
	if x != nil {
		return x.ZipCode
	}
	return ""
}

func (x *Location) GetStreet() string {
	if x != nil {
		return x.Street
	}
	return ""
}

type AddSharedMediaRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Media []*Media `protobuf:"bytes,1,rep,name=media,proto3" json:"media,omitempty"`
}

func (x *AddSharedMediaRequest) Reset() {
	*x = AddSharedMediaRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_content_content_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddSharedMediaRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddSharedMediaRequest) ProtoMessage() {}

func (x *AddSharedMediaRequest) ProtoReflect() protoreflect.Message {
	mi := &file_content_content_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddSharedMediaRequest.ProtoReflect.Descriptor instead.
func (*AddSharedMediaRequest) Descriptor() ([]byte, []int) {
	return file_content_content_proto_rawDescGZIP(), []int{5}
}

func (x *AddSharedMediaRequest) GetMedia() []*Media {
	if x != nil {
		return x.Media
	}
	return nil
}

type AddSharedMediaResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *AddSharedMediaResponse) Reset() {
	*x = AddSharedMediaResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_content_content_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddSharedMediaResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddSharedMediaResponse) ProtoMessage() {}

func (x *AddSharedMediaResponse) ProtoReflect() protoreflect.Message {
	mi := &file_content_content_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddSharedMediaResponse.ProtoReflect.Descriptor instead.
func (*AddSharedMediaResponse) Descriptor() ([]byte, []int) {
	return file_content_content_proto_rawDescGZIP(), []int{6}
}

var File_content_content_proto protoreflect.FileDescriptor

var file_content_content_proto_rawDesc = []byte{
	0x0a, 0x15, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e,
	0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x2c, 0x0a, 0x12, 0x53, 0x68, 0x61, 0x72, 0x65,
	0x64, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a,
	0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x75,
	0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x33, 0x0a, 0x13, 0x53, 0x68, 0x61, 0x72, 0x65, 0x64, 0x4d,
	0x65, 0x64, 0x69, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1c, 0x0a, 0x05,
	0x6d, 0x65, 0x64, 0x69, 0x61, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x06, 0x2e, 0x4d, 0x65,
	0x64, 0x69, 0x61, 0x52, 0x05, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x22, 0xb8, 0x01, 0x0a, 0x05, 0x4d,
	0x65, 0x64, 0x69, 0x61, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08,
	0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63,
	0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64,
	0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x64,
	0x64, 0x65, 0x64, 0x4f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x64, 0x64,
	0x65, 0x64, 0x4f, 0x6e, 0x12, 0x25, 0x0a, 0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x09, 0x2e, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x18, 0x0a, 0x04, 0x74,
	0x61, 0x67, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x04, 0x2e, 0x54, 0x61, 0x67, 0x52,
	0x04, 0x74, 0x61, 0x67, 0x73, 0x22, 0x2b, 0x0a, 0x03, 0x54, 0x61, 0x67, 0x12, 0x14, 0x0a, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02,
	0x69, 0x64, 0x22, 0x6c, 0x0a, 0x08, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x18,
	0x0a, 0x07, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x53, 0x74, 0x61, 0x74,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x18,
	0x0a, 0x07, 0x5a, 0x69, 0x70, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x5a, 0x69, 0x70, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x53, 0x74, 0x72, 0x65,
	0x65, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x53, 0x74, 0x72, 0x65, 0x65, 0x74,
	0x22, 0x35, 0x0a, 0x15, 0x41, 0x64, 0x64, 0x53, 0x68, 0x61, 0x72, 0x65, 0x64, 0x4d, 0x65, 0x64,
	0x69, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x05, 0x6d, 0x65, 0x64,
	0x69, 0x61, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x06, 0x2e, 0x4d, 0x65, 0x64, 0x69, 0x61,
	0x52, 0x05, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x22, 0x18, 0x0a, 0x16, 0x41, 0x64, 0x64, 0x53, 0x68,
	0x61, 0x72, 0x65, 0x64, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x32, 0x8b, 0x01, 0x0a, 0x07, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x3d, 0x0a,
	0x0e, 0x47, 0x65, 0x74, 0x53, 0x68, 0x61, 0x72, 0x65, 0x64, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x12,
	0x13, 0x2e, 0x53, 0x68, 0x61, 0x72, 0x65, 0x64, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x53, 0x68, 0x61, 0x72, 0x65, 0x64, 0x4d, 0x65, 0x64,
	0x69, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x30, 0x01, 0x12, 0x41, 0x0a, 0x0e,
	0x41, 0x64, 0x64, 0x53, 0x68, 0x61, 0x72, 0x65, 0x64, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x12, 0x16,
	0x2e, 0x41, 0x64, 0x64, 0x53, 0x68, 0x61, 0x72, 0x65, 0x64, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x41, 0x64, 0x64, 0x53, 0x68, 0x61, 0x72,
	0x65, 0x64, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42,
	0x15, 0x5a, 0x13, 0x2e, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x2f, 0x70, 0x72, 0x63,
	0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_content_content_proto_rawDescOnce sync.Once
	file_content_content_proto_rawDescData = file_content_content_proto_rawDesc
)

func file_content_content_proto_rawDescGZIP() []byte {
	file_content_content_proto_rawDescOnce.Do(func() {
		file_content_content_proto_rawDescData = protoimpl.X.CompressGZIP(file_content_content_proto_rawDescData)
	})
	return file_content_content_proto_rawDescData
}

var file_content_content_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_content_content_proto_goTypes = []interface{}{
	(*SharedMediaRequest)(nil),     // 0: SharedMediaRequest
	(*SharedMediaResponse)(nil),    // 1: SharedMediaResponse
	(*Media)(nil),                  // 2: Media
	(*Tag)(nil),                    // 3: Tag
	(*Location)(nil),               // 4: Location
	(*AddSharedMediaRequest)(nil),  // 5: AddSharedMediaRequest
	(*AddSharedMediaResponse)(nil), // 6: AddSharedMediaResponse
}
var file_content_content_proto_depIdxs = []int32{
	2, // 0: SharedMediaResponse.media:type_name -> Media
	4, // 1: Media.location:type_name -> Location
	3, // 2: Media.tags:type_name -> Tag
	2, // 3: AddSharedMediaRequest.media:type_name -> Media
	0, // 4: Content.GetSharedMedia:input_type -> SharedMediaRequest
	5, // 5: Content.AddSharedMedia:input_type -> AddSharedMediaRequest
	1, // 6: Content.GetSharedMedia:output_type -> SharedMediaResponse
	6, // 7: Content.AddSharedMedia:output_type -> AddSharedMediaResponse
	6, // [6:8] is the sub-list for method output_type
	4, // [4:6] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_content_content_proto_init() }
func file_content_content_proto_init() {
	if File_content_content_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_content_content_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SharedMediaRequest); i {
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
		file_content_content_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SharedMediaResponse); i {
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
		file_content_content_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Media); i {
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
		file_content_content_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Tag); i {
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
		file_content_content_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Location); i {
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
		file_content_content_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddSharedMediaRequest); i {
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
		file_content_content_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddSharedMediaResponse); i {
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
			RawDescriptor: file_content_content_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_content_content_proto_goTypes,
		DependencyIndexes: file_content_content_proto_depIdxs,
		MessageInfos:      file_content_content_proto_msgTypes,
	}.Build()
	File_content_content_proto = out.File
	file_content_content_proto_rawDesc = nil
	file_content_content_proto_goTypes = nil
	file_content_content_proto_depIdxs = nil
}
