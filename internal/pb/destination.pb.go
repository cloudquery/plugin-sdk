// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.19.4
// source: internal/pb/destination.proto

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

type Migrate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Migrate) Reset() {
	*x = Migrate{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pb_destination_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Migrate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Migrate) ProtoMessage() {}

func (x *Migrate) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pb_destination_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Migrate.ProtoReflect.Descriptor instead.
func (*Migrate) Descriptor() ([]byte, []int) {
	return file_internal_pb_destination_proto_rawDescGZIP(), []int{0}
}

type Write struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Write) Reset() {
	*x = Write{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pb_destination_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Write) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Write) ProtoMessage() {}

func (x *Write) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pb_destination_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Write.ProtoReflect.Descriptor instead.
func (*Write) Descriptor() ([]byte, []int) {
	return file_internal_pb_destination_proto_rawDescGZIP(), []int{1}
}

type Migrate_Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name    string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Version string `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
	Tables  []byte `protobuf:"bytes,3,opt,name=tables,proto3" json:"tables,omitempty"`
}

func (x *Migrate_Request) Reset() {
	*x = Migrate_Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pb_destination_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Migrate_Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Migrate_Request) ProtoMessage() {}

func (x *Migrate_Request) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pb_destination_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Migrate_Request.ProtoReflect.Descriptor instead.
func (*Migrate_Request) Descriptor() ([]byte, []int) {
	return file_internal_pb_destination_proto_rawDescGZIP(), []int{0, 0}
}

func (x *Migrate_Request) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Migrate_Request) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *Migrate_Request) GetTables() []byte {
	if x != nil {
		return x.Tables
	}
	return nil
}

type Migrate_Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Migrate_Response) Reset() {
	*x = Migrate_Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pb_destination_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Migrate_Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Migrate_Response) ProtoMessage() {}

func (x *Migrate_Response) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pb_destination_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Migrate_Response.ProtoReflect.Descriptor instead.
func (*Migrate_Response) Descriptor() ([]byte, []int) {
	return file_internal_pb_destination_proto_rawDescGZIP(), []int{0, 1}
}

type Write_Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// marshalled *schema.Resources
	Resource []byte `protobuf:"bytes,1,opt,name=resource,proto3" json:"resource,omitempty"`
}

func (x *Write_Request) Reset() {
	*x = Write_Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pb_destination_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Write_Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Write_Request) ProtoMessage() {}

func (x *Write_Request) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pb_destination_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Write_Request.ProtoReflect.Descriptor instead.
func (*Write_Request) Descriptor() ([]byte, []int) {
	return file_internal_pb_destination_proto_rawDescGZIP(), []int{1, 0}
}

func (x *Write_Request) GetResource() []byte {
	if x != nil {
		return x.Resource
	}
	return nil
}

type Write_Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// error
	Error string `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *Write_Response) Reset() {
	*x = Write_Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pb_destination_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Write_Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Write_Response) ProtoMessage() {}

func (x *Write_Response) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pb_destination_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Write_Response.ProtoReflect.Descriptor instead.
func (*Write_Response) Descriptor() ([]byte, []int) {
	return file_internal_pb_destination_proto_rawDescGZIP(), []int{1, 1}
}

func (x *Write_Response) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

var File_internal_pb_destination_proto protoreflect.FileDescriptor

var file_internal_pb_destination_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x62, 0x2f, 0x64, 0x65,
	0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x16, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c,
	0x2f, 0x70, 0x62, 0x2f, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x66,
	0x0a, 0x07, 0x4d, 0x69, 0x67, 0x72, 0x61, 0x74, 0x65, 0x1a, 0x4f, 0x0a, 0x07, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73,
	0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x12, 0x16, 0x0a, 0x06, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x73, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x06, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x73, 0x1a, 0x0a, 0x0a, 0x08, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x50, 0x0a, 0x05, 0x57, 0x72, 0x69, 0x74, 0x65, 0x1a,
	0x25, 0x0a, 0x07, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x72, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x1a, 0x20, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x32, 0x9a, 0x02, 0x0a, 0x0b, 0x44, 0x65, 0x73,
	0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x55, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x45,
	0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x1f, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x45, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x43,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x20, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x45, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65,
	0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x40, 0x0a, 0x09, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x65, 0x12, 0x18, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x65, 0x2e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x65, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x3a, 0x0a, 0x07, 0x4d, 0x69, 0x67, 0x72, 0x61, 0x74, 0x65, 0x12, 0x16, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4d, 0x69, 0x67, 0x72, 0x61, 0x74, 0x65, 0x2e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4d, 0x69, 0x67,
	0x72, 0x61, 0x74, 0x65, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x36, 0x0a,
	0x05, 0x57, 0x72, 0x69, 0x74, 0x65, 0x12, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x57,
	0x72, 0x69, 0x74, 0x65, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x57, 0x72, 0x69, 0x74, 0x65, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x28, 0x01, 0x42, 0x05, 0x5a, 0x03, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_pb_destination_proto_rawDescOnce sync.Once
	file_internal_pb_destination_proto_rawDescData = file_internal_pb_destination_proto_rawDesc
)

func file_internal_pb_destination_proto_rawDescGZIP() []byte {
	file_internal_pb_destination_proto_rawDescOnce.Do(func() {
		file_internal_pb_destination_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_pb_destination_proto_rawDescData)
	})
	return file_internal_pb_destination_proto_rawDescData
}

var file_internal_pb_destination_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_internal_pb_destination_proto_goTypes = []interface{}{
	(*Migrate)(nil),                   // 0: proto.Migrate
	(*Write)(nil),                     // 1: proto.Write
	(*Migrate_Request)(nil),           // 2: proto.Migrate.Request
	(*Migrate_Response)(nil),          // 3: proto.Migrate.Response
	(*Write_Request)(nil),             // 4: proto.Write.Request
	(*Write_Response)(nil),            // 5: proto.Write.Response
	(*GetExampleConfig_Request)(nil),  // 6: proto.GetExampleConfig.Request
	(*Configure_Request)(nil),         // 7: proto.Configure.Request
	(*GetExampleConfig_Response)(nil), // 8: proto.GetExampleConfig.Response
	(*Configure_Response)(nil),        // 9: proto.Configure.Response
}
var file_internal_pb_destination_proto_depIdxs = []int32{
	6, // 0: proto.Destination.GetExampleConfig:input_type -> proto.GetExampleConfig.Request
	7, // 1: proto.Destination.Configure:input_type -> proto.Configure.Request
	2, // 2: proto.Destination.Migrate:input_type -> proto.Migrate.Request
	4, // 3: proto.Destination.Write:input_type -> proto.Write.Request
	8, // 4: proto.Destination.GetExampleConfig:output_type -> proto.GetExampleConfig.Response
	9, // 5: proto.Destination.Configure:output_type -> proto.Configure.Response
	3, // 6: proto.Destination.Migrate:output_type -> proto.Migrate.Response
	5, // 7: proto.Destination.Write:output_type -> proto.Write.Response
	4, // [4:8] is the sub-list for method output_type
	0, // [0:4] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_internal_pb_destination_proto_init() }
func file_internal_pb_destination_proto_init() {
	if File_internal_pb_destination_proto != nil {
		return
	}
	file_internal_pb_base_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_internal_pb_destination_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Migrate); i {
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
		file_internal_pb_destination_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Write); i {
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
		file_internal_pb_destination_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Migrate_Request); i {
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
		file_internal_pb_destination_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Migrate_Response); i {
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
		file_internal_pb_destination_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Write_Request); i {
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
		file_internal_pb_destination_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Write_Response); i {
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
			RawDescriptor: file_internal_pb_destination_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_internal_pb_destination_proto_goTypes,
		DependencyIndexes: file_internal_pb_destination_proto_depIdxs,
		MessageInfos:      file_internal_pb_destination_proto_msgTypes,
	}.Build()
	File_internal_pb_destination_proto = out.File
	file_internal_pb_destination_proto_rawDesc = nil
	file_internal_pb_destination_proto_goTypes = nil
	file_internal_pb_destination_proto_depIdxs = nil
}
