// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: internal/pb/base.proto

package pb

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

type GetName struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetName) Reset() {
	*x = GetName{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pb_base_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetName) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetName) ProtoMessage() {}

func (x *GetName) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pb_base_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetName.ProtoReflect.Descriptor instead.
func (*GetName) Descriptor() ([]byte, []int) {
	return file_internal_pb_base_proto_rawDescGZIP(), []int{0}
}

type GetVersion struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetVersion) Reset() {
	*x = GetVersion{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pb_base_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetVersion) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetVersion) ProtoMessage() {}

func (x *GetVersion) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pb_base_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetVersion.ProtoReflect.Descriptor instead.
func (*GetVersion) Descriptor() ([]byte, []int) {
	return file_internal_pb_base_proto_rawDescGZIP(), []int{1}
}

type GetProtocolVersion struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetProtocolVersion) Reset() {
	*x = GetProtocolVersion{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pb_base_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetProtocolVersion) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetProtocolVersion) ProtoMessage() {}

func (x *GetProtocolVersion) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pb_base_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetProtocolVersion.ProtoReflect.Descriptor instead.
func (*GetProtocolVersion) Descriptor() ([]byte, []int) {
	return file_internal_pb_base_proto_rawDescGZIP(), []int{2}
}

type Configure struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Configure) Reset() {
	*x = Configure{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pb_base_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Configure) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Configure) ProtoMessage() {}

func (x *Configure) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pb_base_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Configure.ProtoReflect.Descriptor instead.
func (*Configure) Descriptor() ([]byte, []int) {
	return file_internal_pb_base_proto_rawDescGZIP(), []int{3}
}

type GetName_Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetName_Request) Reset() {
	*x = GetName_Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pb_base_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetName_Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetName_Request) ProtoMessage() {}

func (x *GetName_Request) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pb_base_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetName_Request.ProtoReflect.Descriptor instead.
func (*GetName_Request) Descriptor() ([]byte, []int) {
	return file_internal_pb_base_proto_rawDescGZIP(), []int{0, 0}
}

type GetName_Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *GetName_Response) Reset() {
	*x = GetName_Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pb_base_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetName_Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetName_Response) ProtoMessage() {}

func (x *GetName_Response) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pb_base_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetName_Response.ProtoReflect.Descriptor instead.
func (*GetName_Response) Descriptor() ([]byte, []int) {
	return file_internal_pb_base_proto_rawDescGZIP(), []int{0, 1}
}

func (x *GetName_Response) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type GetVersion_Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetVersion_Request) Reset() {
	*x = GetVersion_Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pb_base_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetVersion_Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetVersion_Request) ProtoMessage() {}

func (x *GetVersion_Request) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pb_base_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetVersion_Request.ProtoReflect.Descriptor instead.
func (*GetVersion_Request) Descriptor() ([]byte, []int) {
	return file_internal_pb_base_proto_rawDescGZIP(), []int{1, 0}
}

type GetVersion_Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Version string `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *GetVersion_Response) Reset() {
	*x = GetVersion_Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pb_base_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetVersion_Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetVersion_Response) ProtoMessage() {}

func (x *GetVersion_Response) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pb_base_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetVersion_Response.ProtoReflect.Descriptor instead.
func (*GetVersion_Response) Descriptor() ([]byte, []int) {
	return file_internal_pb_base_proto_rawDescGZIP(), []int{1, 1}
}

func (x *GetVersion_Response) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

type GetProtocolVersion_Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetProtocolVersion_Request) Reset() {
	*x = GetProtocolVersion_Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pb_base_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetProtocolVersion_Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetProtocolVersion_Request) ProtoMessage() {}

func (x *GetProtocolVersion_Request) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pb_base_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetProtocolVersion_Request.ProtoReflect.Descriptor instead.
func (*GetProtocolVersion_Request) Descriptor() ([]byte, []int) {
	return file_internal_pb_base_proto_rawDescGZIP(), []int{2, 0}
}

type GetProtocolVersion_Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Version uint64 `protobuf:"varint,1,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *GetProtocolVersion_Response) Reset() {
	*x = GetProtocolVersion_Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pb_base_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetProtocolVersion_Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetProtocolVersion_Response) ProtoMessage() {}

func (x *GetProtocolVersion_Response) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pb_base_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetProtocolVersion_Response.ProtoReflect.Descriptor instead.
func (*GetProtocolVersion_Response) Descriptor() ([]byte, []int) {
	return file_internal_pb_base_proto_rawDescGZIP(), []int{2, 1}
}

func (x *GetProtocolVersion_Response) GetVersion() uint64 {
	if x != nil {
		return x.Version
	}
	return 0
}

type Configure_Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Holds information such as credentials, regions, accounts, etc'
	// Marshalled spec.SourceSpec or spec.DestinationSpec
	Config []byte `protobuf:"bytes,1,opt,name=config,proto3" json:"config,omitempty"`
}

func (x *Configure_Request) Reset() {
	*x = Configure_Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pb_base_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Configure_Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Configure_Request) ProtoMessage() {}

func (x *Configure_Request) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pb_base_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Configure_Request.ProtoReflect.Descriptor instead.
func (*Configure_Request) Descriptor() ([]byte, []int) {
	return file_internal_pb_base_proto_rawDescGZIP(), []int{3, 0}
}

func (x *Configure_Request) GetConfig() []byte {
	if x != nil {
		return x.Config
	}
	return nil
}

type Configure_Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Error string `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *Configure_Response) Reset() {
	*x = Configure_Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pb_base_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Configure_Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Configure_Response) ProtoMessage() {}

func (x *Configure_Response) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pb_base_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Configure_Response.ProtoReflect.Descriptor instead.
func (*Configure_Response) Descriptor() ([]byte, []int) {
	return file_internal_pb_base_proto_rawDescGZIP(), []int{3, 1}
}

func (x *Configure_Response) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

var File_internal_pb_base_proto protoreflect.FileDescriptor

var file_internal_pb_base_proto_rawDesc = []byte{
	0x0a, 0x16, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x62, 0x2f, 0x62, 0x61,
	0x73, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x34, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x1a, 0x09, 0x0a, 0x07, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x3d, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x56, 0x65, 0x72, 0x73,
	0x69, 0x6f, 0x6e, 0x1a, 0x09, 0x0a, 0x07, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x24,
	0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x22, 0x45, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x50, 0x72, 0x6f, 0x74, 0x6f,
	0x63, 0x6f, 0x6c, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x1a, 0x09, 0x0a, 0x07, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x24, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x50, 0x0a, 0x09, 0x43,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x65, 0x1a, 0x21, 0x0a, 0x07, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x06, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x1a, 0x20, 0x0a, 0x08, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x42, 0x05, 0x5a,
	0x03, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_pb_base_proto_rawDescOnce sync.Once
	file_internal_pb_base_proto_rawDescData = file_internal_pb_base_proto_rawDesc
)

func file_internal_pb_base_proto_rawDescGZIP() []byte {
	file_internal_pb_base_proto_rawDescOnce.Do(func() {
		file_internal_pb_base_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_pb_base_proto_rawDescData)
	})
	return file_internal_pb_base_proto_rawDescData
}

var file_internal_pb_base_proto_msgTypes = make([]protoimpl.MessageInfo, 12)
var file_internal_pb_base_proto_goTypes = []interface{}{
	(*GetName)(nil),                     // 0: proto.GetName
	(*GetVersion)(nil),                  // 1: proto.GetVersion
	(*GetProtocolVersion)(nil),          // 2: proto.GetProtocolVersion
	(*Configure)(nil),                   // 3: proto.Configure
	(*GetName_Request)(nil),             // 4: proto.GetName.Request
	(*GetName_Response)(nil),            // 5: proto.GetName.Response
	(*GetVersion_Request)(nil),          // 6: proto.GetVersion.Request
	(*GetVersion_Response)(nil),         // 7: proto.GetVersion.Response
	(*GetProtocolVersion_Request)(nil),  // 8: proto.GetProtocolVersion.Request
	(*GetProtocolVersion_Response)(nil), // 9: proto.GetProtocolVersion.Response
	(*Configure_Request)(nil),           // 10: proto.Configure.Request
	(*Configure_Response)(nil),          // 11: proto.Configure.Response
}
var file_internal_pb_base_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_internal_pb_base_proto_init() }
func file_internal_pb_base_proto_init() {
	if File_internal_pb_base_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_pb_base_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetName); i {
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
		file_internal_pb_base_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetVersion); i {
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
		file_internal_pb_base_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetProtocolVersion); i {
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
		file_internal_pb_base_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Configure); i {
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
		file_internal_pb_base_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetName_Request); i {
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
		file_internal_pb_base_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetName_Response); i {
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
		file_internal_pb_base_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetVersion_Request); i {
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
		file_internal_pb_base_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetVersion_Response); i {
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
		file_internal_pb_base_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetProtocolVersion_Request); i {
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
		file_internal_pb_base_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetProtocolVersion_Response); i {
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
		file_internal_pb_base_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Configure_Request); i {
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
		file_internal_pb_base_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Configure_Response); i {
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
			RawDescriptor: file_internal_pb_base_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   12,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_internal_pb_base_proto_goTypes,
		DependencyIndexes: file_internal_pb_base_proto_depIdxs,
		MessageInfos:      file_internal_pb_base_proto_msgTypes,
	}.Build()
	File_internal_pb_base_proto = out.File
	file_internal_pb_base_proto_rawDesc = nil
	file_internal_pb_base_proto_goTypes = nil
	file_internal_pb_base_proto_depIdxs = nil
}
